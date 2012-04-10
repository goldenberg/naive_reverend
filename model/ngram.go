package model

import (
	"fmt"
	counter "naive_reverend/counter"
	distribution "naive_reverend/distribution"
	store "naive_reverend/store"
)

const (
	PRIOR = "prior"
	CLASS = "class"
)

type NGramModel struct {
	N int
	s store.Interface
}

var _ Interface = new(NGramModel)

func NewNGramModel(s store.Interface, n int) *NGramModel {
	return &NGramModel{n, s}
}

/*
 * Counter of number of instances in training set. i.e. N_c
 */
func (m *NGramModel) priorCounter() (c counter.Interface, ok bool) {
	return m.fetch(PRIOR, "")
}

/*
 * Distribution P(class)
 */
func (m *NGramModel) Prior() (d distribution.Interface, ok bool) {
	c, ok := m.priorCounter()
	d = distribution.NewLaplacian(c)
	return
}

/*
 * Number of Bins (B)
 */
func (m *NGramModel) Bins() int {
	c, ok := m.priorCounter()
	if ok {
		return len(c.Keys())
	}
	return 0
}

func (m *NGramModel) fetch(prefix, ngram string) (c counter.Interface, ok bool) {
	key := fmt.Sprintf("%v:%v", prefix, ngram)
	fmt.Println("looking up", key)
	c, ok = m.s.Fetch(key)
	return
}

func (m *NGramModel) incr(prefix, numerator, denominator string, incr int64) int64 {
	key := fmt.Sprintf("%v:%v", prefix, numerator)
	return m.s.IncrN(key, denominator, incr)
}

<<<<<<< HEAD
=======
/*
 * Lookup an n-gram's frequency, i.e. C(w_1 ... w_n)
 */
>>>>>>> df7ecd4... initial work toward pipelining. might be misguided
func (m *NGramModel) classLookup(ngram NGram) (c counter.Interface, ok bool) {
	return m.fetch(CLASS, ngram.String())
}

func (m *NGramModel) incrPrior(class string, incr int64) {
	// Increment "prior:", class
	m.incr(PRIOR, "", class, incr)
}

func (m *NGramModel) incrClasses(ngram NGram, class string, incr int64) {
	m.incr(CLASS, ngram.String(), class, incr)
}

/*
 * Estimate P(w_1 ... w_n | C) for all C
 */
func (m *NGramModel) Estimate(ngram NGram) distribution.Interface {
	c, ok := m.classLookup(ngram)
	if !ok {
		c = counter.New()
	}
	return distribution.NewLaplacian(c)
}

/*
 * Number of values in the multinomial target feature distribution (B)
 */
func (m *NGramModel) ClassCount() int {
	d, ok := m.Prior()
	if !ok {
		return 0
	}
	return d.Len()
}

func (m *NGramModel) Train(datum *Datum) {
	m.incrPrior(datum.Class, datum.Count)
	for n := 1; n <= m.N; n++ {
		for _, ngram := range Generate(datum.Features, n) {
			m.incrClasses(ngram, datum.Class, datum.Count)
		}
	}
}

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
