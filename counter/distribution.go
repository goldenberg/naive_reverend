package counter

import (
	`fmt`
	`math`
)

type Distribution interface {
	Get(k string) float64
	Keys() []string
	LogGet(k string) float64
	LogSet(k string, v float64)
}

type CounterDistribution struct {
	counter Counter
}

var _ Distribution = new(CounterDistribution)

// func NewCounterDistribution(c *Counter) (d *Distribution) {
// 	for k
// }

func (d *CounterDistribution) Get(k string) float64 {
	return float64(d.counter.Get(k)+1) / float64(d.counter.Sum()+1)
}

// Return a list of keys for this counter
func (d *CounterDistribution) Keys() []string {
	return d.counter.Keys()
}

func (d *CounterDistribution) LogGet(k string) float64 {
	return math.Log(float64(d.counter.Get(k)+1)) - math.Log(float64(d.counter.Sum()+1))
}

func (d *CounterDistribution) LogSet(k string, v float64) {
	//fmt.Println("k:", k, "v:", v, "d:", d)
	fmt.Println("baaaaaaaaaaaaaaaaaaaaaaaaaaad")
}

func Multiply(a, b Distribution) Distribution {
	result := NewDerivedDistribution()
	fmt.Println("Multiply a:", a, "b:", b)
	for k := range mergeKeys(a.Keys(), b.Keys()) {
		// fmt.Println("LogGet key:", k, "d:", d.LogGet(k), "o:", d.LogGet(k))
		result.LogSet(k, a.LogGet(k) + b.LogGet(k))
	}
	return result
}


func Divide(a, b Distribution) Distribution {
	result := NewDerivedDistribution()
	fmt.Println("Divide a:", a, "b:", b)
	for k := range mergeKeys(a.Keys(), b.Keys()) {
		// fmt.Println("LogGet key:", k, "d:", d.LogGet(k), "o:", d.LogGet(k))
		result.LogSet(k, a.LogGet(k) - b.LogGet(k))
	}
	return result
}

func ArgMax(d Distribution) (key string, probability float64) {
	probability = math.Inf(-1)
	for _, k := range d.Keys() {
		p := d.Get(key)
		if p > probability {
			key = k
			probability = p
		}
	}
	return
}

type DerivedDistribution struct {
	logProbabilities map[string]float64
}

var _ Distribution = new(DerivedDistribution)

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

func (d *DerivedDistribution) LogSet(k string, v float64) {
	d.logProbabilities[k] = v
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
