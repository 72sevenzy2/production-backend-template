package main

import (
	"sync"
	"sync/atomic"
)

type Cache[T ParamType] struct {
	Loader LoaderFn[T]
	SfFunc *SingleFlight[string, T]

	mu    sync.RWMutex
	Store L1Cache[T]
}

func NewCache[T ParamType](loaderFn LoaderFn[T], sfFn *SingleFlight[string, T]) *Cache[T] {
	return &Cache[T]{
		SfFunc: sfFn,
		Loader: loaderFn,
		Store:  make(L1Cache[T]),
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
				return entry.Value, nil
			}

			// todo: add cache penetration prevention.
			val, err := c.Loader(k)
			if err != nil {
				var zero T
				return zero, err
			}

			c.Put(k, val)
			return val, nil
		})
		return val, err == nil
	}
	c.mu.RUnlock()

	entry.Hits.Add(1)
	return entry.Value, true
}

func (c *Cache[T]) Put(k string, v T) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// LFU (Least Frequently Used) eviction policy.
	// todo: add LRU tie-breaking.
	// todo 2: add exponential decay to cache entry's frequency.
	if len(c.Store) >= 10 {
		var (
			victimKey string
			// bitwise NOT operator (^) for the highest possible uint64 value (equivalent to math.MaxUint64).
			lowestHits uint64 = ^uint64(0)
		)
		for key, entry := range c.Store {
			c := entry.Hits.Load()
			if c < lowestHits {
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
