package main

import (
	"sync"
)

type L1Cache[T any] map[string]T

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
	defer c.mu.RUnlock()
	if _, ok := c.Store[k]; !ok {
		var zero T
		return zero, false
	}

	return c.Store[k], true
}

func (c *Cache[T]) Put(k string, v T) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.Store) >= 10 {
		return false
	}

	if _, ok := c.Store[k]; ok {
		return false
	}

	c.Store[k] = v
	return true
}
