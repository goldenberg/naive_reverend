package model

import (
	"fmt"
	"strings"
	counter "naive_reverend/counter"
	store "naive_reverend/store"
	distribution "naive_reverend/distribution"
)

const (
	BLANK = "_"
	NGRAM = "ngram"
	PRIOR = "prior"
	CLASS = "class"
)

type NGram []string

func Generate(terms []string, n int) (ngrams []NGram) {
	if len(terms) == 0 {
		return []NGram{}
	}
	ngrams = make([]NGram, 0)
	for i := 0; i < len(terms)+n-1; i++ {
		ngrams = append(ngrams, getNgram(terms, i, n))
	}
	return
}

func getNgram(terms []string, pos, n int) (ngram NGram) {
	ngram = make(NGram, 0)
	start := pos - n + 1

	for i := start; i <= pos; i++ {
		if i < 0 || i >= len(terms) {
			ngram = append(ngram, BLANK)
		} else {
			ngram = append(ngram, terms[i])
		}
	}
	return
}

func (ng NGram) String() string {
	return strings.Join(ng, " ")
}

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

/*
 * Lookup an n-gram's frequency, i.e. C(w_1 ... w_n)
 */
func (m *NGramModel) classLookup(ngram NGram) (c counter.Interface, ok bool) {
	return m.fetch(CLASS, ngram.String())
}

func (m *NGramModel) incrPrior(class string, incr int64) {
	// Increment "prior:", class
	m.incr(PRIOR, "", class, incr)
}

func (m *NGramModel) incrNGram(ngram NGram, incr int64) {
	n := len(ngram)
	denominator := ngram.String()
	var numerator string
	if len(ngram) > 1 {
		numerator = ngram[:n-1].String()
	} else {
		numerator = ""
	}
	m.incr(NGRAM, numerator, denominator, incr)
}

func (m *NGramModel) incrClasses(ngram NGram, class string, incr int64) {
	m.incr(CLASS, ngram.String(), class, incr)
}

/*
 * Estimate P(w_1 ... w_n)
 */
func (m *NGramModel) Estimate(ngram NGram) distribution.Interface {
	c, ok := m.classLookup(ngram)
	if !ok {
		c = counter.New()
	}
	return distribution.NewLaplacian(c)
}

type ExplainedDistrib struct {
	distribution distribution.Interface
	ngram        NGram
	explain      map[string]interface{}
}

func (m *NGramModel) EstimateMany(ngrams []NGram) (distribs chan ExplainedDistrib) {
	distribs = make(chan ExplainedDistrib)
	go func() {
		defer close(distribs)
		for _, ngram := range ngrams {
			d := m.Estimate(ngram)
			fmt.Println("estimated", ngram, "as", d.Keys())
			explained := ExplainedDistrib{d, ngram, distribution.JSON(d)}
			distribs <- explained
		}
	}()
	return
}

/*
 * Number of values in the multinomial target feature distribution (B)
 */
func (m *NGramModel) ClassCount() int {
	d, ok := m.Prior()
	if !ok {
		return 0
	}
	return len(d.Keys())
}

func (m *NGramModel) Train(datum *Datum) {
	m.incrPrior(datum.Class, datum.Count)
	for n := 1; n <= m.N; n++ {
		for _, ngram := range Generate(datum.Features, n) {
			// m.incrNGram(ngram, datum.Count)
			m.incrClasses(ngram, datum.Class, datum.Count)
		}
	}
}

func (m *NGramModel) Classify(features []string) (estimator distribution.Interface, explain map[string]interface{}) {
	explain = make(map[string]interface{})
	estimator, _ = m.Prior()
	explain["prior"] = distribution.JSON(estimator)
	ngrams := Generate(features, m.N)
	for ngram_est := range m.EstimateMany(ngrams) {
		estimator = distribution.Multiply(estimator, ngram_est.distribution)
		explain[ngram_est.ngram.String()] = distribution.JSON(ngram_est.distribution)
	}
	estimator = distribution.Normalize(estimator)
	return
}
