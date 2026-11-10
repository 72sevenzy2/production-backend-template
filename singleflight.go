package main

import (
	"sync"
)

type ParamType interface {
	string | int
}

type Call[T ParamType] struct {
	wg  sync.WaitGroup
	val T
	err error
}

type SingleFlight[K comparable, V ParamType] struct {
	mu    sync.Mutex
	calls map[K]*Call[V]
}

func NewSingleFlight[K comparable, V ParamType]() *SingleFlight[K, V] {
	return &SingleFlight[K, V]{
		calls: make(map[K]*Call[V]),
	}
}

// execute-once logic function for cache stampede prevention.
func (sf *SingleFlight[K, V]) Do(key K, fn DependencyFn[V]) (V, error) {
	sf.mu.Lock()

	if c, ok := sf.calls[key]; ok {
		sf.mu.Unlock()

		c.wg.Wait()
		return c.val, c.err
	}

	c := &Call[V]{}
	c.wg.Add(1)
	sf.calls[key] = c

	sf.mu.Unlock()

	c.val, c.err = fn()

	sf.mu.Lock()
	delete(sf.calls, key)
	sf.mu.Unlock()

	c.wg.Done()
	return c.val, c.err
}
