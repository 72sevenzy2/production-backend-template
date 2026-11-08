package main

import "sync"

type Call struct {
	wg  sync.WaitGroup
	val string
	err error
}

type SingleFlight struct {
	mu    sync.Mutex
	calls map[string]*Call
}

func NewSingleFlight() *SingleFlight {
	return &SingleFlight{
		calls: make(map[string]*Call),
	}
}

// execute-once logic function for cache stampede prevention.
func (sf *SingleFlight) Do(key string, fn func() (string, error)) (string, error) {
	sf.mu.Lock()

	if c, ok := sf.calls[key]; ok {
		sf.mu.Unlock()

		c.wg.Wait()
		return c.val, c.err
	}

	c := &Call{}
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
