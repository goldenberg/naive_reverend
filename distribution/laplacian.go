package distribution

import (
	"math"
	counter "naive_reverend/counter"
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
	classes := d.counter.Len()
	// If there aren't any classes, we assume it's a binary classification, which might not be true, but we don't know anything else.
	if classes == 0 {
		classes = 2
	}
	return (float64(d.counter.Get(k)) + d.alpha) / (float64(d.counter.Sum()) + d.alpha*float64(classes))
}

func (d *Laplacian) LogGet(k string) float64 {
	classes := d.counter.Len()
	// If there aren't any classes, we assume it's a binary classification, which might not be true, but we don't know anything else.
	if classes == 0 {
		classes = 2
	}
	return math.Log(float64(d.counter.Get(k))+d.alpha) - math.Log(float64(d.counter.Sum())+d.alpha*float64(classes))
}

func (d *Laplacian) Len() int {
	return d.counter.Len()
}
