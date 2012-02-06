package main

import counter "gnlp/counter"
import frozencounter "gnlp/frozencounter"

type NaiveBayes struct {
	FeatureLogDistributions map[string]*frozencounter.Counter
	ClassLogPrior           *frozencounter.Counter
}

type Datum struct {
	Class    string
	Features []string
}

func Train(data chan *Datum) *NaiveBayes {
	class := counter.New(0.0)
	features := make(map[string]*counter.Counter)

	for datum := range data {
		class.Incr(datum.Class)
		for _, f := range datum.Features {
			dist, ok := features[f]

			if !ok {
				dist = counter.New(0.0)
				features[f] = dist
			}

			dist.Incr(datum.Class)
		}
	}

	class.LogNormalize()
	for _, dist := range features {
		dist.LogNormalize()
	}

	frozenFeatures := frozencounter.FreezeMap(features)

	var keyset *frozencounter.KeySet
	for _, dist := range frozenFeatures {
		keyset = dist.Keys
	}

	frozenClass := frozencounter.FreezeWithKeySet(class, keyset)

	return &NaiveBayes{FeatureLogDistributions: frozenFeatures, ClassLogPrior: frozenClass}
}

func (nb *NaiveBayes) Classify(features []string) (string, float64) {
	score := nb.ClassLogPrior.Copy()

	for _, f := range features {
		dist, ok := nb.FeatureLogDistributions[f]
		if ok {
			score.Add(dist)
		}
	}

	score.Exp()
	score.Normalize()

	c, probability := score.ArgMax()
	return c, probability
}
