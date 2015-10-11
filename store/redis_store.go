package store

import (
	"fmt"
	godis "github.com/simonz05/godis"
	counter "github.com/goldenberg/naive_reverend/counter"
	"strconv"
)

type RedisStore struct {
	client    *godis.Client
	keyPrefix string
}

var _ Interface = new(RedisStore)

func NewRedisStore(keyPrefix string) Interface {
	client := godis.New("", 0, "")
	return &RedisStore{client, keyPrefix}
}

func (s *RedisStore) dbKey(k string) string {
	return fmt.Sprintf("%s:%s", s.keyPrefix, k)
}

func (s *RedisStore) Fetch(name string) (c counter.Interface, ok bool) {
	r, err := s.client.Hgetall(s.dbKey(name))
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
			panic(err)
		}
		out[k] = i
	}
	return
}
func (s *RedisStore) Incr(name, key string) int64 {
	return s.IncrN(name, key, 1)
}

func (s *RedisStore) IncrN(name, key string, n int64) int64 {
	val, err := s.client.Hincrby(s.dbKey(name), key, n)
	if err != nil {
		panic(fmt.Sprintf("err: ", err))
	}
	return val
}

/*
 * Now that we have multiple keyspaces, this isn't correct. It gets
 * the size across all keyspaces.
 */
func (s *RedisStore) Size() (size int64) {
	size, err := s.client.Dbsize()
	if err != nil {
		panic("couldn't get size")
	}
	return
}

func (s *RedisStore) Get(key string) (val string) {
	r, err := s.client.Get(s.dbKey(key))
	if err != nil {
		panic("couldn't Get")
	}
	return r.String()
}

func (s *RedisStore) Set(key, val string) {
	err := s.client.Set(s.dbKey(key), val)
	if err != nil {
		panic("couldn't Set")
	}
	return
}
