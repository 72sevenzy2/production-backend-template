package main

import (
	"sync/atomic"
)

// Entry represents an individual cache entry, with its value and number of times accessed.
type Entry[T ParamType] struct {
	Value T
	Hits  atomic.Uint64
}

// L1Cache represents the underlying type of the cache's data store.
type L1Cache[T ParamType] map[string]*Entry[T]

// DependencyFn is the function executed once in SingeFlight.Do().
type DependencyFn[T ParamType] func() (T, error)

// LoaderFn is a function responsible for retrieving a user from the database incase of nil user ids.
type LoaderFn[T ParamType] func(string) (T, bool)
