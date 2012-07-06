// Package store provides an interface for fetching and incrementing counters
// in memory or in a persistent K/V store.
package store

import (
	counter "naive_reverend/counter"
)

type Interface interface {
	Fetch(name string) (c counter.Interface, ok bool)
	Incr(name, key string) int64
	IncrN(name, key string, n int64) int64
	Size() int64
}
