package model

import (
	// `fmt`
	counter `naive_reverend/counter`
)

type NaiveBayes struct {
	KeySet *counter.KeySet
	FeatureCategoryCounters map[string]counter.Counter
	ClassCounter            counter.Counter
	FeatureCounter counter.Counter
}

type Datum struct {
	Class    string
	Features []string
}

func New() *NaiveBayes {
	return &NaiveBayes{
		counter.NewKeySet(), 
		make(map[string]counter.Counter), 
		counter.New(), 
		counter.New(),
	}
}

func (nb *NaiveBayes) Train(datum *Datum) {
	nb.ClassCounter.Incr(datum.Class)
	for _, f := range datum.Features {
		dist, ok := nb.FeatureCategoryCounters[f]

		if !ok {
			dist = counter.NewKeySetCounter(nb.KeySet)
			nb.FeatureCategoryCounters[f] = dist
		}

		dist.Incr(datum.Class)

		nb.FeatureCounter.Incr(f)
	}
}

// {
// 	'multiply': {
// 		'prior': {'b': 0.5, 'a': 0.5},
// 		'feature1': {'a': 0.1, 'b': 0.9}
// 	}
// 	'divide': {
// 		'feature1': 
// 	}
// }
func (nb *NaiveBayes) Classify(features []string) (string, float64) {
	var estimator counter.Distribution
	estimator = nb.ClassCounter.Distribution()
	// fmt.Println("Prior:", estimator)

	featureDistribution := nb.FeatureCounter.Distribution()

	for _, f := range features {
		c, ok := nb.FeatureCategoryCounters[f]
		// fmt.Println("Feature:", f, "Counter:", c)
		if ok {
			dist := c.Distribution()
			estimator = estimator.Multiply(dist)
			// fmt.Println("Estimator: ", estimator)
		}
	}

	return estimator.ArgMax()
}
