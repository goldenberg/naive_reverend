package main

import ()

type Counter interface {
	Get(string) uint
	Set(string, uint)
	Incr(string)
	Keys() []string
	Sum() uint
	Distribution() *CounterDistribution
}

type MemCounter struct {
	values map[string]uint
}

var _ Counter = new(MemCounter)

func New() *MemCounter {
	return &MemCounter{make(map[string]uint)}
}

func (c *MemCounter) Get(k string) uint {
	return c.values[k]
}

func (c *MemCounter) Set(k string, v uint) {
	c.values[k] = v
}

func (c *MemCounter) Incr(k string) {
	c.values[k] += 1
}

// Return a list of keys for this counter
func (c *MemCounter) Keys() []string {
	result := make([]string, 0, len(c.values))

	for k, _ := range c.values {
		result = append(result, k)
	}

	return result
}

func (c *MemCounter) Sum() (result uint) {
	for _, v := range c.values {
		result += v
	}
	return
}

func (c *MemCounter) Distribution() *CounterDistribution {
	return &CounterDistribution{c}
}
