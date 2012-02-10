package main

import `fmt`

type NaiveBayes struct {
	FeatureCategoryCounters map[string]*MemCounter
	ClassCounter            *MemCounter
}

type Datum struct {
	Class    string
	Features []string
}

func Train(data chan *Datum) *NaiveBayes {
	class := New()
	features := make(map[string]*MemCounter)

	for datum := range data {
		class.Incr(datum.Class)
		for _, f := range datum.Features {
			dist, ok := features[f]

			if !ok {
				dist = New()
				features[f] = dist
			}

			dist.Incr(datum.Class)
		}
	}

	return &NaiveBayes{FeatureCategoryCounters: features, ClassCounter: class}
}

func (nb *NaiveBayes) Classify(features []string) (string, float64) {
	var estimator Distribution
	estimator = nb.ClassCounter.Distribution()
	fmt.Println("Prior:", estimator)

	for _, f := range features {
		counter, ok := nb.FeatureCategoryCounters[f]
		fmt.Println("Feature:", f, "Counter:", counter)
		if ok {
			dist := counter.Distribution()
			estimator = estimator.Multiply(dist)
			fmt.Println("Estimator: ", estimator)
		}
	}

	return estimator.ArgMax()
}
