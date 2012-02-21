package model

import (
	// `fmt`
	distribution "naive_reverend/distribution"
)

type Datum struct {
	Class    string
	Features []string
}

type Interface interface {
	Train(datum *Datum)
	Classify(features []string) (estimator distribution.Interface, explain map[string]interface{})
}