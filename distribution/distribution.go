package distribution

import (
	"fmt"
	"math"
	counter "naive_reverend/counter"
)

type Interface interface {
	Get(k string) float64
	Keys() []string
	LogGet(k string) float64
}

type CounterDistribution struct {
	counter counter.Interface
}

var _ Interface = new(CounterDistribution)

func NewCounterDistribution(c counter.Interface) (d Interface) {
	return &CounterDistribution{c}
}

func (d *CounterDistribution) Get(k string) float64 {
	return float64(d.counter.Get(k)) / float64(d.counter.Sum())
}

// Return a list of keys for this counter
func (d *CounterDistribution) Keys() []string {
	return d.counter.Keys()
}

func (d *CounterDistribution) LogGet(k string) float64 {
	return math.Log(float64(d.counter.Get(k)+1)) - math.Log(float64(d.counter.Sum()+1))
}

func Multiply(a, b Interface) Interface {
	logProbs := make(map[string]float64)
	// fmt.Println("Multiply a:", a, "b:", b)
	for k := range mergeKeys(a.Keys(), b.Keys()) {
		// fmt.Println("LogGet key:", k, "d:", d.LogGet(k), "o:", d.LogGet(k))
		logProbs[k] = a.LogGet(k) + b.LogGet(k)
	}
	return &DerivedDistribution{logProbs}
}

func Divide(a, b Interface) Interface {
	logProbs := make(map[string]float64)
	fmt.Println("Divide a:", a, "b:", b)
	for k := range mergeKeys(a.Keys(), b.Keys()) {
		// fmt.Println("LogGet key:", k, "d:", d.LogGet(k), "o:", d.LogGet(k))
		logProbs[k] = a.LogGet(k) - b.LogGet(k)
	}
	return &DerivedDistribution{logProbs}
}

func JSON(d Interface) (out map[string]interface{}) {
	out = make(map[string]interface{})
	for _, k := range d.Keys() {
		out[k] = map[string]float64{
			"p(k)":      d.Get(k),
			"log(p(k))": d.LogGet(k),
		}
	}
	return
}

func ArgMax(d Interface) (maxKey string, maxProb float64) {
	maxProb = math.Inf(-1)
	for _, k := range d.Keys() {
		p := d.Get(k)
		if p > maxProb {
			maxKey = k
			maxProb = p
		}
	}
	return
}

type DerivedDistribution struct {
	logProbabilities map[string]float64
}

var _ Interface = new(DerivedDistribution)

func NewDerivedDistribution() *DerivedDistribution {
	return &DerivedDistribution{make(map[string]float64)}
}

func (d *DerivedDistribution) Get(k string) float64 {
	return math.Exp(d.LogGet(k))
}

// Return a list of keys for this counter
func (d *DerivedDistribution) Keys() []string {
	result := make([]string, 0, len(d.logProbabilities))

	for k, _ := range d.logProbabilities {
		result = append(result, k)
	}

	return result
}

func (d *DerivedDistribution) LogGet(k string) float64 {
	logProb, ok := d.logProbabilities[k]
	if ok {
		return logProb
	}
	return 0.0
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
