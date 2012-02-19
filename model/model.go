package model

import (
	// `fmt`
	distribution "naive_reverend/distribution"
)

type Datum struct {
	Class    string
	Features []string
	Count int64
}

type Interface interface {
	Train(datum *Datum)
	TrainN(datum *Datum, n uint)
	Classify(features []string) (estimator distribution.Interface, explain map[string]interface{})
}