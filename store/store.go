// Package store provides an interface for fetching and incrementing counters
// in memory or in a persistent K/V store.
package store

import (
	counter "github.com/goldenberg/naive_reverend/counter"
)

// A type that stores a set of counters that can be fetched and incremented by
// name.
type Interface interface {
	// Fetch returns the counter with the given name.
	Fetch(name string) (c counter.Interface, ok bool)
	// Incr increments the key of the named counter by 1.
	Incr(name, key string) int64
	// Incr increments the key of the named counter by n.
	IncrN(name, key string, n int64) int64
	// Size returns the number of counters currently stored.
	Size() int64
}
