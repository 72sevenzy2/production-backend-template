package main

import (
	"sync"
	"sync/atomic"
)

// MAX_CACHE_SIZE is the max number of entries of the l1 cache.
const MAX_CACHE_SIZE = 10

type Cache[T ParamType] struct {
	// LoaderFunc is a function reponsible for retrieving a user on the database.
	LoaderFunc LoaderFn[T]
	// SfFunc is a function responsible for request coalescing to prevent cache stampedes.
	SfFunc *SingleFlight[string, T]

	mu sync.RWMutex
	// Store represents the l1 map of the cache.
	Store L1Cache[T]

	// clock is a global counter responsible for computing an entry's last-used sequence number.
	clock atomic.Uint64
}

func NewCache[T ParamType](loaderFn LoaderFn[T], sfFn *SingleFlight[string, T]) *Cache[T] {
	return &Cache[T]{
		SfFunc:     sfFn,
		LoaderFunc: loaderFn,
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
				entry.Hits.Add(1)
				entry.LastAccess.Store(c.clock.Add(1))
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

	entry.Hits.Add(1)
	entry.LastAccess.Store(c.clock.Add(1))
	return entry.Value, true
}

func (c *Cache[T]) Put(k string, v T) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// LFU + LRU tie-breaking eviction policy.
	// todo: add exponential decay to cache entry's hit frequency.
	if len(c.Store) >= MAX_CACHE_SIZE {
		var (
			victimKey string
			// bitwise NOT operator (^) for the highest possible uint64 value (equivalent to math.MaxUint64).
			lowestHits uint64 = ^uint64(0)
			oldest     uint64 = ^uint64(0)
		)
		for key, entry := range c.Store {
			c := entry.Hits.Load()
			lastAccess := entry.LastAccess.Load()
			if c < lowestHits || (c == lowestHits && lastAccess < oldest) {
				oldest = lastAccess
				lowestHits = c
				victimKey = key
			}
		}
		delete(c.Store, victimKey)
	}

	if _, ok := c.Store[k]; ok {
		return false
	}

	c.Store[k] = &Entry[T]{
		Hits:  atomic.Uint64{},
		Value: v,
	}
	return true
}
