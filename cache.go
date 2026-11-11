package main

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// MAX_CACHE_SIZE is the max number of entries of the l1 cache.
const MAX_CACHE_SIZE = 10

type Cache[T ParamType] struct {
	// LoaderFunc is a function reponsible for retrieving a user on the database.
	LoaderFunc LoaderFn[T]
	// SfFunc is a function responsible for request coalescing to prevent cache stampedes.
	SfFunc *SingleFlight[string, T]
	// Store represents the l1 map of the cache.
	mu    sync.RWMutex
	Store L1Cache[T]
	// clock is a global counter responsible for computing an entry's last-used sequence number.
	clock atomic.Uint64
	// halfLife represents the duration until an entry's access frequency is reduced to half its amount.
	halfLife time.Duration
}

func NewCache[T ParamType](loaderFn LoaderFn[T], sfFn *SingleFlight[string, T]) *Cache[T] {
	return &Cache[T]{
		SfFunc:     sfFn,
		LoaderFunc: loaderFn,
		halfLife:   10 * time.Minute,
		Store:      make(L1Cache[T]),
	}
}

func (c *Cache[T]) Get(k string) (T, bool) {
	c.mu.RLock()
	entry, ok := c.Store[k]
	if !ok {
		c.mu.RUnlock()
		val, err := c.SfFunc.Do(k, func() (T, error) {
			c.mu.RLock()
			entry, ok := c.Store[k]
			c.mu.RUnlock()

			if ok {
				c.mu.Lock()
				c.recordAccess(entry)
				c.mu.Unlock()
				return entry.Value, nil
			}

			val, ok := c.LoaderFunc(k)
			if !ok {
				var zero T
				return zero, ErrInternalServerError
			}

			c.Put(k, val)
			return val, nil
		})
		return val, err == nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	c.recordAccess(entry)
	c.mu.Unlock()
	return entry.Value, true
}

func (c *Cache[T]) Put(k string, v T) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// LFU + exponential frequency decay paired with LRU tie-breaking.
	if len(c.Store) >= MAX_CACHE_SIZE {
		var (
			victimKey  string
			lowestHits float64 = math.MaxFloat64
			// bitwise NOT operator (^) for the highest possible uint64 value (equivalent to math.MaxUint64).
			oldest uint64 = ^uint64(0)
		)
		now := time.Now()
		for key, entry := range c.Store {
			elapsed := now.Sub(entry.LastDecay)

			effectiveHits := entry.Hits * math.Exp(
				-math.Ln2*elapsed.Seconds()/c.halfLife.Seconds(),
			)

			if effectiveHits < lowestHits || (effectiveHits == lowestHits && entry.LastAccess < oldest) {
				oldest = entry.LastAccess
				lowestHits = effectiveHits
				victimKey = key
			}
		}
		delete(c.Store, victimKey)
	}

	if _, ok := c.Store[k]; ok {
		return false
	}

	c.Store[k] = &Entry[T]{
		LastDecay: time.Now(),
		Value:     v,
	}
	return true
}

func (c *Cache[T]) decay(entry *Entry[T], now time.Time) {
	elapsed := now.Sub(entry.LastDecay)
	if elapsed <= 0 {
		return
	}
	entry.Hits *= math.Exp(
		-math.Ln2 * elapsed.Seconds() / c.halfLife.Seconds(),
	)
	entry.LastDecay = now
}

func (c *Cache[T]) recordAccess(entry *Entry[T]) {
	now := time.Now()
	c.decay(entry, now)
	entry.Hits++
	entry.LastAccess = c.clock.Add(1)
}
