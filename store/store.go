package store 

import (
	counter "naive_reverend/counter"
)

type Interface interface {
	Fetch(name string) (c counter.Interface, ok bool)
	Incr(name, key string)
	IncrN(name, key string, n int64)
	Size() int64
}

type MemCounterStore map[string]counter.Interface

var _ Interface = new(MemCounterStore)

func New() Interface {
	return MemCounterStore(make(map[string]counter.Interface, 2))
}

func (s MemCounterStore) Fetch(name string) (c counter.Interface, ok bool) {
	c, ok = s[name]
	return
}

func (s MemCounterStore) Incr(name, key string) {
	s.IncrN(name, key, 1)
}

func (s MemCounterStore) IncrN(name, key string, n int64) {
	c, ok := s[name]	
	if !ok {
		c = counter.New()
		s[name] = c
	}
	c.IncrN(key, n)
}

func (s MemCounterStore) Size() int64 {
	return int64(len(s))
}