package distribution

import (
	"fmt"
	"math"
	counter "naive_reverend/counter"
	"sort"
)

type Interface interface {
	Get(k string) float64
	Keys() []string
	LogGet(k string) float64
}

type CounterDistribution struct {
	counter counter.Interface
}

// Return a list of keys for this counter
func (d *CounterDistribution) Keys() []string {
	return d.counter.Keys()
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

// func JSON(d Interface) (out map[string]interface{}) {
// 	out = make(map[string]interface{})
// 	for _, k := range d.Keys() {
// 		out[k] = map[string]float64{
// 			"p(k)":      d.Get(k),
// 			"log(p(k))": d.LogGet(k),
// 		}
// 	}
// 	return
// }

func JSON(d Interface) (out map[string]float64) {
	// out = make(map[string]interface{})
	// for _, k := range d.Keys() {
	// 	out[k] = d.Get(k)
	// }
	return TopN(d, 10)
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

type stringFloatPair struct {
	key string
	val float64
}

type StringFloatSlice []*stringFloatPair

func (s StringFloatSlice) Len() int           { return len(s) }
func (s StringFloatSlice) Less(i, j int) bool { return s[i].val < s[j].val }
func (s StringFloatSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

/* Sort the keys and values and return the top N classes and their probabilities. */
func TopN(d Interface, n int) (out map[string]float64) {
	sorted := make(StringFloatSlice, 0)
	for _, k := range d.Keys() {
		sorted = append(sorted, &stringFloatPair{k, d.Get(k)})
	}
	// A heap would be slightly more efficient. eh.
	sort.Sort(sorted)
	out = make(map[string]float64)
	startIdx := len(sorted) - n
	if startIdx < 0 {
		startIdx = 0
	}
	for _, p := range sorted[startIdx:] {
		out[p.key] = p.val
	}
	return
}

func Sum(d Interface) (sum float64) {
	for _, k := range d.Keys() {
		sum += d.Get(k)
	}
	return
}

func Normalize(d Interface) Interface {
	logProbs := make(map[string]float64)
	sum := Sum(d)
	for _, k := range d.Keys() {
		logProbs[k] = math.Log(d.Get(k) / sum)
	}
	return &DerivedDistribution{logProbs}
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
