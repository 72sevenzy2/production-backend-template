package main

import "time"

// Entry represents an individual cache entry, with its value and number of times accessed.
type Entry[T ParamType] struct {
	// Value represents the cache entry's value (a type to which satisfies ParamType)
	Value T
	// Hits is the number of times an entry has been accessed.
	Hits float64
	// LastDecay represents the last time an entry's access frequency had been decayed.
	LastDecay time.Time
	// LastAccess represents an entry's sequence number to determine whether it has recently been accessed.
	LastAccess uint64
}

// L1Cache represents the underlying type of the cache's data store.
type L1Cache[T ParamType] map[string]*Entry[T]

// DependencyFn is the function executed once in SingeFlight.Do().
type DependencyFn[T ParamType] func() (T, error)

// LoaderFn is a function responsible for retrieving a user from the database incase of nil user ids.
type LoaderFn[T ParamType] func(string) (T, bool)
