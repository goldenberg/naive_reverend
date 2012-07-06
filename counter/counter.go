package counter

import (
	"sync"
	"fmt"
)

type Interface interface {
	Get(string) int64
	Set(string, int64)
	Incr(string) int64
	IncrN(string, int64) int64
	Keys() []string
	Len() int
	Sum() int64
	String() string
}

type MemCounter struct {
	counts map[string]int64
	sum    int64
	lock   *sync.RWMutex
}

var _ Interface = new(MemCounter)

func New() Interface {
	return &MemCounter{make(map[string]int64), 0, new(sync.RWMutex)}
}

func NewCounter(counts map[string]int64) Interface {
	return &MemCounter{counts, 0, new(sync.RWMutex)}
}

// Get returns the count for the key.
func (c *MemCounter) Get(k string) int64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.counts[k]
}

// Set safely sets the count for the key.
func (c *MemCounter) Set(k string, v int64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	prevCount = c.counts[k]
	c.counts[k] = v
	c.sum += v - prevCount
}

// Incr safely increments the count for the key by 1.
func (c *MemCounter) Incr(k string) int64 {
	return c.IncrN(k, 1)
}

// IncrN safely increments the count for the key by n.
func (c *MemCounter) IncrN(k string, n int64) int64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	c.counts[k] += n
	c.sum += n
	return c.counts[k]
}

// Len returns the number of keys in the counter.
func (c *MemCounter) Len() int {
	return len(c.counts)
}

// Keys returns a list of keys for this counter
func (c *MemCounter) Keys() []string {
	result := make([]string, 0, len(c.counts))

	for k := range c.counts {
		result = append(result, k)
	}

	return result
}

// Sum returns the sum of all of the counts.
func (c *MemCounter) Sum() (result int64) {
	// Theoretically, the counts could sum to 0, in which case we're doing
	// extra work, and not actually memoizing.
	if c.sum == 0 && len(c.counts) > 0 {
		c.sum = c.computeSum()
	}
	return c.sum
}

func (c *MemCounter) computeSum() (result int64) {
	for _, v := range c.counts {
		result += v
	}
	return
}

// String returns a string representation of the counter.
func (c *MemCounter) String() string {
	s := "{"

	for _, key := range c.Keys() {
		s += fmt.Sprintf("'%s': %f, ", key, c.Get(key))
	}

	s += "}"

	return s
}
