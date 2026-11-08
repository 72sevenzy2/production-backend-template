package main

import "sync"

type ParamType interface {
	string | int
}

type Call[T ParamType] struct {
	wg  sync.WaitGroup
	val T
	err error
}

type SingleFlight[T ParamType] struct {
	mu    sync.Mutex
	calls map[T]*Call[T]
}

func NewSingleFlight[T ParamType]() *SingleFlight[T] {
	return &SingleFlight[T]{
		calls: make(map[T]*Call[T]),
	}
}

// execute-once logic function for cache stampede prevention.
func (sf *SingleFlight[T]) Do(key T, fn func() (T, error)) (T, error) {
	sf.mu.Lock()

	if c, ok := sf.calls[key]; ok {
		sf.mu.Unlock()

		c.wg.Wait()
		return c.val, c.err
	}

	c := &Call[T]{}
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
