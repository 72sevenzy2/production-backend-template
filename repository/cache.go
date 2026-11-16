package repository

import (
	"math"
	"sync"
	"sync/atomic"
	"time"

	utils "github.com/72sevenzy2/production-backend-template/internal"
)

// MAX_CACHE_SIZE is the max number of entries of the l1 cache.
const MAX_CACHE_SIZE = 10

type Cache[T utils.ParamType] struct {
	// LoaderFunc is a function reponsible for retrieving a user on the database.
	LoaderFunc utils.LoaderFn[T]
	// SfFunc is a function responsible for request coalescing to prevent cache stampedes.
	SfFunc *utils.SingleFlight[string, T]
	// Store represents the l1 map of the cache.
	mu    sync.RWMutex
	Store utils.L1Cache[T]
	// clock is a global counter responsible for computing an entry's last-used sequence number.
	clock atomic.Uint64
	// halfLife represents the duration until an entry's access frequency is reduced to half its amount.
	halfLife time.Duration
}

func NewCache[T utils.ParamType](loaderFn utils.LoaderFn[T], sfFn *utils.SingleFlight[string, T], halflife time.Duration) *Cache[T] {
	if halflife == time.Duration(0) || halflife < time.Duration(0) {
		halflife = time.Minute * 10
	}

	return &Cache[T]{
		SfFunc:     sfFn,
		LoaderFunc: loaderFn,
		halfLife:   halflife,
		Store:      make(utils.L1Cache[T]),
	}
}

func (c *Cache[T]) Get(k string) (T, bool) {
	c.mu.Lock()
	entry, ok := c.Store[k]
	if !ok {
		c.mu.Unlock()
		val, err := c.SfFunc.Do(k, func() (T, error) {
			c.mu.Lock()
			entry, ok := c.Store[k]

			if ok {
				c.recordAccess(entry)
				c.mu.Unlock()
				return entry.Value, nil
			}
			c.mu.Unlock()

			val, ok := c.LoaderFunc(k)
			if !ok {
				var zero T
				return zero, utils.ErrInternalServerError
			}

			c.Put(k, val)
			return val, nil
		})
		return val, err == nil
	}
	c.recordAccess(entry)
	c.mu.Unlock()
	return entry.Value, true
}

func (c *Cache[T]) Edit(k string, v T) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	initialStore, ok := c.Store[k]
	if !ok {
		return false
	}

	// copying on write.
	newEntry := &utils.Entry[T]{
		Value:      v,
		Hits:       initialStore.Hits,
		LastDecay:  initialStore.LastDecay,
		LastAccess: initialStore.LastAccess,
	}

	c.Store[k] = newEntry
	return true
}

func (c *Cache[T]) Put(k string, v T) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Store[k]; ok {
		return false
	}

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

	c.Store[k] = &utils.Entry[T]{
		LastDecay: time.Now(),
		Value:     v,
	}
	return true
}

func (c *Cache[T]) Del(k string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Store, k)
}

func (c *Cache[T]) decay(entry *utils.Entry[T], now time.Time) {
	elapsed := now.Sub(entry.LastDecay)
	if elapsed <= 0 {
		return
	}
	entry.Hits *= math.Exp(
		-math.Ln2 * elapsed.Seconds() / c.halfLife.Seconds(),
	)
	entry.LastDecay = now
}

func (c *Cache[T]) recordAccess(entry *utils.Entry[T]) {
	now := time.Now()
	c.decay(entry, now)
	entry.Hits++
	entry.LastAccess = c.clock.Add(1)
}
