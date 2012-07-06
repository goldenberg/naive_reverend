package store

import (
	counter "naive_reverend/counter"
)

type MemCounterStore map[string]counter.Interface

var _ Interface = new(MemCounterStore)

func New() Interface {
	return MemCounterStore(make(map[string]counter.Interface, 2))
}

func (s MemCounterStore) Fetch(name string) (c counter.Interface, ok bool) {
	c, ok = s[name]
	return
}

func (s MemCounterStore) Incr(name, key string) int64 {
	return s.IncrN(name, key, 1)
}

func (s MemCounterStore) IncrN(name, key string, n int64) int64 {
	c, ok := s[name]
	if !ok {
		c = counter.New()
		s[name] = c
	}
	return c.IncrN(key, n)
}

func (s MemCounterStore) Size() int64 {
	return int64(len(s))
}

////

type MemStore map[string]string

func (s MemStore) Set(key, val string) {
	s[key] = val
}

func (s MemStore) Get(key string) (val string) {
	return s[key]
}

func (s MemStore) Items() (out chan []string) {
	out = make(chan []string)
	go func(out chan<- []string) {
		defer close(out)
		for k, v := range s {
			out <- []string{k, v}
		}
	}(out)
	return
}
