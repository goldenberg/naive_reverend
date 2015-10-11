package model

import (
	// `fmt`
	distribution "github.com/goldenberg/naive_reverend/distribution"
)

type Datum struct {
	Class    string
	Features []string
	Count    int64
}

type Interface interface {
	Train(datum *Datum)
	Classify(features []string) (estimator distribution.Interface, explain map[string]interface{})
}
