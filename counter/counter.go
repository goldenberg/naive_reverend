package counter

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
	return &MemCounter{make(map[string]uint, 2)}
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

type KeySetCounter struct {
	values []uint
	keyset *KeySet
}

var _ Counter = new(KeySetCounter)

func NewKeySetCounter(ks *KeySet) *KeySetCounter {
	return &KeySetCounter{make([]uint, 0), ks}
}

func (c *KeySetCounter) Get(k string) uint {
	idx, ok := c.keyset.Get(k)
	if !ok {
		return 0
	}
	return c.values[idx]
}

func (c *KeySetCounter) Set(k string, v uint) {
	idx, _ := c.keyset.Get(k)
	c.values[idx] = v
}

func (c *KeySetCounter) Incr(k string) {
	c.Set(k, c.Get(k) + 1)
}

// Return a list of keys for this counter
func (c *KeySetCounter) Keys() []string {
	return c.keyset.Keys()
}

func (c *KeySetCounter) Sum() (result uint) {
	for _, v := range c.values {
		result += v
	}
	return
}

func (c *KeySetCounter) Distribution() *CounterDistribution {
	return &CounterDistribution{c}
}
