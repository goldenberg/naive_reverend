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

	Multiply(o Distribution) (result Distribution)
	Divide(o Distribution) (result Distribution)
	ArgMax() (k string, probability float64)
}

type CounterDistribution struct {
	counter Counter
}

var _ Distribution = new(CounterDistribution)

// func NewCounterDistribution(c *Counter) {
// 	d = Distribution	
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

func (d *CounterDistribution) Multiply(o Distribution) Distribution {
	result := NewGeneratedDistribution()
	fmt.Println("Multiply d:", d, "o:", o)
	for k := range mergeKeys(d.Keys(), o.Keys()) {
		// fmt.Println("LogGet key:", k, "d:", d.LogGet(k), "o:", d.LogGet(k))
		result.LogSet(k, d.LogGet(k)+o.LogGet(k))
	}
	return result
}

func (d *CounterDistribution) Divide(o Distribution) (result Distribution) {
	result = NewGeneratedDistribution()
	for k := range mergeKeys(d.Keys(), o.Keys()) {
		result.LogSet(k, d.LogGet(k)-o.LogGet(k))
	}
	return
}

func (d *CounterDistribution) ArgMax() (key string, probability float64) {
	probability = math.Inf(-1)
	for _, k := range d.Keys() {
		p := d.Get(key)
		if p > probability {
			key = k
			probability = p
		} else {
		}
	}
	return
}

type GeneratedDistribution struct {
	logProbabilities map[string]float64
}

var _ Distribution = new(GeneratedDistribution)

func NewGeneratedDistribution() *GeneratedDistribution {
	return &GeneratedDistribution{make(map[string]float64)}
}

func (d *GeneratedDistribution) Get(k string) float64 {
	return math.Exp(d.LogGet(k))
}

// Return a list of keys for this counter
func (d *GeneratedDistribution) Keys() []string {
	result := make([]string, 0, len(d.logProbabilities))

	for k, _ := range d.logProbabilities {
		result = append(result, k)
	}

	return result
}

func (d *GeneratedDistribution) LogGet(k string) float64 {
	logProb, ok := d.logProbabilities[k]
	if ok {
		return logProb
	}
	return 0.0
}

func (d *GeneratedDistribution) LogSet(k string, v float64) {
	//fmt.Println("k:", k, "v:", v, "d:", d)
	d.logProbabilities[k] = v
}

func (d *GeneratedDistribution) Multiply(o Distribution) (result Distribution) {
	result = NewGeneratedDistribution()
	// fmt.Println("Multiply d:", d, "o:", o)
	for k := range mergeKeys(d.Keys(), o.Keys()) {
		// fmt.Println("LogGet key:", k, "d:", d.LogGet(k), "o:", d.LogGet(k))
		result.LogSet(k, d.LogGet(k)+o.LogGet(k))
	}
	return
}

func (d *GeneratedDistribution) Divide(o Distribution) (result Distribution) {
	for k := range mergeKeys(d.Keys(), o.Keys()) {
		result.LogSet(k, d.LogGet(k)-o.LogGet(k))
	}
	return
}

func (d *GeneratedDistribution) ArgMax() (k string, probability float64) {
	maxLogProb := math.Inf(-1)
	var maxKey string
	for testKey, logProb := range d.logProbabilities {
		if logProb > maxLogProb {
			// fmt.Println("Found logProb", logProb, ">", maxLogProb)
			maxKey = testKey
			maxLogProb = logProb
		}
	}
	return maxKey, math.Exp(maxLogProb)
}

// borrowed from mattj
func (c *GeneratedDistribution) String() string {
	s := "Counter: {"

	for _, key := range c.Keys() {
		s += fmt.Sprintf("'%s': %f, ", key, c.Get(key))
	}

	s += "}"

	return s
}

func (c *CounterDistribution) String() string {
	s := "Counter: {"

	for _, key := range c.Keys() {
		s += fmt.Sprintf("'%s': %f, ", key, c.Get(key))
	}

	s += "}"

	return s
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
