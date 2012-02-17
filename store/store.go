package store 

import (
	counter "naive_reverend/counter"
)

type Interface interface {
	Fetch(name string) (c counter.Interface, ok bool)
	Incr(name, key string)
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
	c, ok := s[name]	
	if !ok {
		c = counter.New()
		s[name] = c
	}
	c.Incr(key)
}

