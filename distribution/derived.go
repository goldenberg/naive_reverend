package distribution

import (
	"math"
)

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

	for k := range d.logProbabilities {
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
