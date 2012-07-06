package model

import (
	counter "naive_reverend/counter"
	distribution "naive_reverend/distribution"
	store "naive_reverend/store"
)

const (
	PRIOR = "__PRIOR__"
)

type NGramModel struct {
	N int
	s store.Interface
}

var _ Interface = new(NGramModel)

func NewNGramModel(s store.Interface, n int) *NGramModel {
	return &NGramModel{n, s}
}

// Counter of number of instances in training set. i.e. N_c
func (m *NGramModel) priorCounter() (c counter.Interface, ok bool) {
	return m.s.Fetch(PRIOR)
}

// Distribution P(class)
func (m *NGramModel) Prior() (d distribution.Interface, ok bool) {
	c, ok := m.priorCounter()
	d = distribution.NewLaplacian(c)
	return
}

// Number of Bins (B)
func (m *NGramModel) Bins() int {
	c, ok := m.priorCounter()
	if ok {
		return len(c.Keys())
	}
	return 0
}

func (m *NGramModel) incr(feature, class string, incr int64) int64 {
	return m.s.IncrN(feature, class, incr)
}

// Lookup an n-gram's frequency across all labels, i.e. C(w_1 ... w_n)
func (m *NGramModel) classLookup(ngram NGram) (c counter.Interface, ok bool) {
	return m.s.Fetch(ngram.String())
}

// Estimate P(w_1 ... w_n | C) for all C
func (m *NGramModel) Estimate(ngram NGram) distribution.Interface {
	c, ok := m.classLookup(ngram)
	if !ok {
		c = counter.New()
	}
	return distribution.NewLaplacian(c)
}

// Number of possible classes.
func (m *NGramModel) ClassCount() int {
	d, ok := m.Prior()
	if !ok {
		return 0
	}
	return d.Len()
}

// Train adds a new datum to the corpus by incrementing counts for all of its
// ngrams.
func (m *NGramModel) Train(datum *Datum) {
	m.incr(PRIOR, datum.Class, datum.Count)
	for n := 1; n <= m.N; n++ {
		for _, ngram := range Generate(datum.Features, n) {
			m.incr(ngram.String(), datum.Class, datum.Count)
		}
	}
}

// Classify estimates the probability distribution given the specified features.
func (m *NGramModel) Classify(features []string) (estimator distribution.Interface, explain map[string]interface{}) {
	explain = make(map[string]interface{})
	estimator, _ = m.Prior()
	explain["prior"] = distribution.JSON(estimator)
	for _, ngram := range Generate(features, m.N) {
		ngram_est := m.Estimate(ngram)
		estimator = distribution.Multiply(estimator, ngram_est)
		explain[ngram.String()] = distribution.JSON(ngram_est)
	}
	estimator = distribution.Normalize(estimator)
	return
}
