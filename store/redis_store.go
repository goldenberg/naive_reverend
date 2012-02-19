package store

import (
	godis "github.com/simonz05/godis"
	counter "naive_reverend/counter"
	"fmt"
	"strconv"
)

type RedisStore struct {
	client *godis.Client
}

var _ Interface = new(RedisStore)

func NewRedisStore() Interface {
	c := godis.New("", 0, "")
	return &RedisStore{c}
}

func (s *RedisStore) Fetch(name string) (c counter.Interface, ok bool) {
	r, err := s.client.Hgetall(name)
	ok = (err == nil)
	if ok {
		intMap := stringMapToIntMap(r.StringMap())
		c = counter.MemCounter(intMap)
	}
	return
}

func stringMapToIntMap(strMap map[string]string) (out map[string]int64) {
	out = make(map[string]int64, len(strMap))
	for k, v := range strMap {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			fmt.Println("err: ", err)
		}
		out[k] = i
	}
	return
}
func (s *RedisStore) Incr(name, key string) {
	_, err := s.client.Hincrby(name, key, 1)
	if err != nil {
		panic("err: ")
	}
}
