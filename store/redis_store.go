package store

import (
	"fmt"
	godis "github.com/simonz05/godis"
	counter "naive_reverend/counter"
	"strconv"
)

type RedisStore struct {
	client *godis.Client
}

var _ Interface = new(RedisStore)
var _ KVInterface = new(RedisStore)

func NewRedisStore() *RedisStore {
	c := godis.New("", 0, "")
	return &RedisStore{c}
}

func (s *RedisStore) Fetch(name string) (c counter.Interface, ok bool) {
	r, err := s.client.Hgetall(name)
	ok = (err == nil)
	if ok {
		intMap := stringMapToIntMap(r.StringMap())
		c = counter.NewCounter(intMap)
		if len(intMap) == 0 {
			ok = false
		}
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
func (s *RedisStore) Incr(name, key string) int64 {
	return s.IncrN(name, key, 1)
}

func (s *RedisStore) IncrN(name, key string, n int64) int64 {
	val, err := s.client.Hincrby(name, key, n)
	if err != nil {
		panic(fmt.Sprintf("err: ", err))
	}
	return val
}

func (s *RedisStore) Size() (size int64) {
	size, err := s.client.Dbsize()
	if err != nil {
		panic("couldn't get size")
	}
	return
}

func (s *RedisStore) Get(key string) (val string) {
	r, err := s.client.Get(key)
	if err != nil {
		panic("couldn't Get")
	}
	return r.String()
}

func (s *RedisStore) Set(key, val string) {
	err := s.client.Set(key, val)
	if err != nil {
		panic("couldn't Get")
	}
	return
}
