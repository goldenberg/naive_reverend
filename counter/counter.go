package counter

import (
	"fmt"
)

type Interface interface {
	Get(string) int64
	Set(string, int64)
	Incr(string)
	IncrN(string, uint)
	Keys() []string
	Sum() int64
	String() string
}

type MemCounter map[string]int64

var _ Interface = new(MemCounter)

func New() Interface {
	return &MemCounter{}
}

func (c MemCounter) Get(k string) int64 {
	return c[k]
}

func (c MemCounter) Set(k string, v int64) {
	c[k] = v
}

func (c MemCounter) Incr(k string) {
	c[k] += 1
}

func (c MemCounter) IncrN(k string, n uint) {
	c[k] += n
}

// Return a list of keys for this counter
func (c MemCounter) Keys() []string {
	result := make([]string, 0, len(c))

	for k := range c {
		result = append(result, k)
	}

	return result
}

func (c MemCounter) Sum() (result int64) {
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
