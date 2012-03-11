package distribution

import (
	counter "naive_reverend/counter"
	"math"
)

type Laplacian struct {
	*CounterDistribution
	alpha float64
}

var _ Interface = new(Laplacian)

func NewLaplacian(c counter.Interface) (d Interface) {
	cd := &CounterDistribution{c}
	return &Laplacian{cd, 1}
}

func NewMLE(c counter.Interface) (d Interface) {
	cd := &CounterDistribution{c}
	return &Laplacian{cd, 0}
}

func (d *Laplacian) Get(k string) float64 {
	return float64(d.counter.Get(k)) / float64(d.counter.Sum())
}

func (d *Laplacian) LogGet(k string) float64 {
	return math.Log(float64(d.counter.Get(k))+d.alpha) - math.Log(float64(d.counter.Sum())+d.alpha*float64(len(d.counter.Keys())))
}
