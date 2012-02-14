package model

import (
	// `fmt`
	counter "naive_reverend/counter"
)

type NaiveBayes struct {
	FeatureCategoryCounters map[string]*counter.MemCounter
	ClassCounter            *counter.MemCounter
}

type Datum struct {
	Class    string
	Features []string
}

func New() *NaiveBayes {
	return &NaiveBayes{make(map[string]*counter.MemCounter), counter.New()}
}

func (nb *NaiveBayes) Train(datum *Datum) {
	nb.ClassCounter.Incr(datum.Class)
	for _, f := range datum.Features {
		dist, ok := nb.FeatureCategoryCounters[f]

		if !ok {
			dist = counter.New()
			nb.FeatureCategoryCounters[f] = dist
		}

		dist.Incr(datum.Class)
	}
}

func (nb *NaiveBayes) Classify(features []string) (estimator counter.Distribution, explain map[string]interface{}) {
	explain = make(map[string]interface{})
	estimator = counter.NewCounterDistribution(nb.ClassCounter)

	explain["prior"] = counter.JSON(estimator)
	// fmt.Println("Prior:", estimator)

	for _, f := range features {
		c, ok := nb.FeatureCategoryCounters[f]
		// fmt.Println("Feature:", f, "Counter:", c)
		var dist counter.Distribution
		if ok {
			dist = counter.NewCounterDistribution(c)
		} else {
			dist = counter.NewDerivedDistribution()
		}
		explain[f] = counter.JSON(dist)
		estimator = counter.Multiply(estimator, dist)
		// fmt.Println("Estimator: ", estimator)
	}

	return 
}