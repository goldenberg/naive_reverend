package model

import (
	// `fmt`
	counter "naive_reverend/counter"
	distribution "naive_reverend/distribution"
	store "naive_reverend/store"
)

type NaiveBayes struct {
	FeatureCategoryCounters store.Interface
	ClassCounter            counter.Interface
}

var _ Interface = new(NaiveBayes)

func New() *NaiveBayes {
	return &NaiveBayes{store.NewRedisStore(), counter.New()}
}

func (nb *NaiveBayes) Train(datum *Datum) {
	nb.ClassCounter.IncrN(datum.Class, datum.Count)
	for _, f := range datum.Features {
		nb.FeatureCategoryCounters.IncrN(f, datum.Class, datum.Count)
	}
}

func (nb *NaiveBayes) Classify(features []string) (estimator distribution.Interface, explain map[string]interface{}) {
	explain = make(map[string]interface{})
	estimator = distribution.NewCounterDistribution(nb.ClassCounter)

	explain["prior"] = distribution.JSON(estimator)
	// Println("Prior:", estimator)

	for _, f := range features {
		c, ok := nb.FeatureCategoryCounters.Fetch(f)
		// fmt.Println("Feature:", f, "Counter:", c)
		var dist distribution.Interface
		if ok {
			dist = distribution.NewCounterDistribution(c)
		} else {
			dist = distribution.NewDerivedDistribution()
		}
		explain[f] = distribution.JSON(dist)
		estimator = distribution.Multiply(estimator, dist)
		// fmt.Println("Estimator: ", estimator)
	}

	return
}
