package main

import (
`fmt`
counter `naive_reverend/counter`
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

func Train(data chan *Datum) *NaiveBayes {
	class := counter.New()
	features := make(map[string]*counter.MemCounter)

	for datum := range data {
		class.Incr(datum.Class)
		for _, f := range datum.Features {
			dist, ok := features[f]

			if !ok {
				dist = counter.New()
				features[f] = dist
			}

			dist.Incr(datum.Class)
		}
	}

	return &NaiveBayes{FeatureCategoryCounters: features, ClassCounter: class}
}

func (nb *NaiveBayes) Classify(features []string) (string, float64) {
	var estimator counter.Distribution
	estimator = nb.ClassCounter.Distribution()
	fmt.Println("Prior:", estimator)

	for _, f := range features {
		c, ok := nb.FeatureCategoryCounters[f]
		fmt.Println("Feature:", f, "Counter:", c)
		if ok {
			dist := c.Distribution()
			estimator = estimator.Multiply(dist)
			fmt.Println("Estimator: ", estimator)
		}
	}

	return estimator.ArgMax()
}