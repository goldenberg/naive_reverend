package main

import (
`math`
`fmt`
)
	
type Counter interface {
	Get(string) uint
	Set(string, uint)
	Incr(string)
	Keys() []string
	Sum() uint
	Distribution() *Distribution
}

type MemCounter struct {
	values map[string]uint
}

var _ Counter = new(MemCounter)

func New() *MemCounter {
	return &MemCounter{make(map[string]uint)}
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

func (c *MemCounter) Distribution() (*Distribution) {
	logProbs := make(map[string]float64)
	logSum := math.Log(float64(c.Sum()))

	for _, k := range c.Keys() {
		logProbs[k] = math.Log(float64(c.Get(k))) - logSum
	}

	//fmt.Println("mc.D logProbs:", logProbs)
	d := &Distribution{logProbs}
	//fmt.Println("&Dist:", d)
	return d
}

type Distribution struct {
	logProbabilities map[string]float64
	counter *Counter
}

func NewDistribution() *Distribution {
	return &Distribution{make(map[string]float64)}
}

func (d *Distribution) Get(k string) (float64) {
	return math.Exp(d.LogGet(k))
}

// Return a list of keys for this counter
func (d *Distribution) Keys() []string {
	result := make([]string, 0, len(d.logProbabilities))

	for k, _ := range d.logProbabilities {
		result = append(result, k)
	}

	return result
}

func (d *Distribution) LogGet(k string) float64 {
	logProb, ok := d.logProbabilities[k]
	if ok {
		return logProb
	}
	return 0.0
}

func (d *Distribution) LogSet(k string, v float64) {
	//fmt.Println("k:", k, "v:", v, "d:", d)
	d.logProbabilities[k] = v
}

func (d *Distribution) Multiply(o *Distribution) (result *Distribution) {
	result = NewDistribution()
	fmt.Println("Multiply d:", d, "o:", o)
	for k := range mergeKeys(d.Keys(), o.Keys()) {
		fmt.Println("LogGet d:", d.LogGet(k), "o:", d.LogGet(k))
		result.LogSet(k, d.LogGet(k) + o.LogGet(k))
	}
	return
}

func (d *Distribution) Divide(o *Distribution) (result *Distribution) {
	for k := range mergeKeys(d.Keys(), o.Keys()) {
		result.LogSet(k, d.LogGet(k) - o.LogGet(k))
	}
	return
}

func (d *Distribution) ArgMax() (k string, probability float64) {
	maxLogProb := math.Inf(-1)
	var maxKey string
	for testKey, logProb := range d.logProbabilities {
		if logProb > maxLogProb {
			maxKey = testKey
			maxLogProb = maxLogProb
		} else {
			fmt.Println("Found logProb", logProb, "<", maxLogProb)
		}
	}
	return maxKey, math.Exp(maxLogProb)
}

// Combine two sets of keys w/o duplicates
// borrowed from mattj
func mergeKeys(a, b []string) <-chan string {
	out := make(chan string)

	go func(out chan<- string) {
		defer close(out)

		seen := make(map[string]bool)

		for _, k := range a {
			out <- k
			seen[k] = true
		}

		for _, k := range b {
			if !seen[k] {
				out <- k
			}
		}
	}(out)

	return out
}
