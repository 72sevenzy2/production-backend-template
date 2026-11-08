package main

import (
	"sync"
	"sync/atomic"
)

type Entry[T any] struct {
	Value T
	Hits  atomic.Uint64 // LFU (Least Frequently Used) eviction.
}

type L1Cache[T any] map[string]*Entry[T]

type Cache[T any] struct {
	mu    sync.RWMutex
	Store L1Cache[T]
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		Store: make(L1Cache[T]),
	}
}

func (c *Cache[T]) Get(k string) (T, bool) {
	c.mu.RLock()
	entry, ok := c.Store[k]
	if !ok {
		c.mu.RUnlock()
		var zero T
		return zero, false
	}
	c.mu.RUnlock()

	entry.Hits.Add(1)
	return c.Store[k].Value, true
}

func (c *Cache[T]) Put(k string, v T) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.Store) >= 10 {
		var (
			victimKey string
			// bitwise ^ operator (equivalent to math.MaxUint64)
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
