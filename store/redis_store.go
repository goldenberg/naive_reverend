package store

import (
	godis "github.com/simonz05/godis"
	counter "naive_reverend/counter"
	"fmt"
	"strconv"
)

type RedisStore struct {
	client    *godis.PipeClient
	keyPrefix string
}

var _ Interface = new(RedisStore)
var _ KVInterface = new(RedisStore)

func NewRedisStore(keyPrefix string) Interface {
	client := godis.NewPipeClient("", 0, "")
	return &RedisStore{client, keyPrefix}
}

func (s *RedisStore) dbKey(k string) string {
	return fmt.Sprintf("%s:%s", s.keyPrefix, k)
}

func (s *RedisStore) Fetch(name string) (c counter.Interface, ok bool) {
	// r, err := s.client.Hgetall(s.dbKey(name))
	// ok = (err == nil)
	// if ok {
	// 	intMap := stringMapToIntMap(r.StringMap())
	// 	fmt.Println(intMap)
	// 	c = counter.MemCounter(intMap)
	// 	fmt.Println(c)
	// 	if len(intMap) == 0 {
	// 		ok = false
	// 	}
	// }
	// return
	nameChan := make(chan string, 1)
	nameChan <- name
	counters, ok := s.FetchMany(nameChan)
	c = <-counters
	fmt.Println("c", c)
	return c, ok
}

func (s *RedisStore) FetchMany(names chan string) (counters chan counter.Interface, ok bool) {
	counters = make(chan counter.Interface, 100)
	defer close(counters)
	go func() {
		for segment := range Segment(names, 100) {
			s.client.Multi()
			for _, name := range segment {
				s.client.Hgetall(s.dbKey(name))
			}
			replies := s.client.Exec()
			fmt.Println("replies", replies)
			for _, r := range replies {
				intMap := stringMapToIntMap(r.StringMap())
				c := counter.MemCounter(intMap)
				if len(intMap) == 0 {
					ok = false
				}
				if ok {
					counters <- c
				} else {
					return
				}
			}
		}
	}()
	return
}

func Segment(in chan string, segmentSize int) (out chan []string) {
	out = make(chan []string)
	defer close(out)

	go func() {
		segment := make([]string, 0)
		for x := range in {
			segment = append(segment, x)
			if len(segment) == segmentSize {
				out <- segment
				segment = make([]string, 0)
			}
		}
		if len(segment) != 0 {
			out <- segment
		}
	}()
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
		panic("couldn't Get")
	}
	return
}
