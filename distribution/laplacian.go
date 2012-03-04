package distribution 

import (
	counter "naive_reverend/counter"
	"math"
)

type Laplacian struct {
	*CounterDistribution
}

var _ Interface = new(Laplacian)

func NewLaplacian(c counter.Interface) (d Interface) {
	cd := &CounterDistribution{c}
	return &Laplacian{cd}
}

func (d *Laplacian) Get(k string) float64 {
	return float64(d.counter.Get(k)) / float64(d.counter.Sum())
}


func (d *Laplacian) LogGet(k string) float64 {
	return math.Log(float64(d.counter.Get(k)+1)) - math.Log(float64(d.counter.Sum()+1))
}