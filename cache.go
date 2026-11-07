package main

import (
	"sync"
)

// type Cache interface {
// 	Get(v string) any
// 	Put(k string, v io.Reader) bool
// 	PutTTL(k string, v io.Reader, n time.Duration) bool
// }

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

func (c *Cache[V]) Put(k string, v V) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Store[k]; ok {
		return false
	}

	c.Store[k] = v
	return true
}
