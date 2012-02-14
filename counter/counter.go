package counter

import (
	"fmt"
)

type Interface interface {
	Get(string) uint
	Set(string, uint)
	Incr(string)
	Keys() []string
	Sum() uint
}

type MemCounter map[string]uint

var _ Interface = new(MemCounter)

func New() *MemCounter {
	return &MemCounter{}
}

func (c MemCounter) Get(k string) uint {
	return c[k]
}

func (c MemCounter) Set(k string, v uint) {
	c[k] = v
}

func (c MemCounter) Incr(k string) {
	c[k] += 1
}

// Return a list of keys for this counter
func (c MemCounter) Keys() []string {
	result := make([]string, 0, len(c))

	for k, _ := range c {
		result = append(result, k)
	}

	return result
}

func (c MemCounter) Sum() (result uint) {
	for _, v := range c {
		result += v
	}
	return
}

func (c MemCounter) String() string {
	s := "Counter: {"

	for _, key := range c.Keys() {
		s += fmt.Sprintf("'%s': %f, ", key, c.Get(key))
	}

	s += "}"

	return s
}
