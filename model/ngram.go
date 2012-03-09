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
	n int
	s store.Interface
}

var _ Interface = new(NGramModel)

func NewNGramModel(n int) *NGramModel {
	return &NGramModel{n, store.NewRedisStore()}
}

func (m *NGramModel) Prior() (d distribution.Interface, ok bool) {
	c, ok := m.fetch("prior", nil)
	if ok {
		d = distribution.NewLaplacian(c)
	} else {
		d = nil
	}
	return
}

func (m *NGramModel) fetch(prefix string, ngram NGram) (c counter.Interface, ok bool) {
	key := fmt.Sprintf("%v:%v", prefix, ngram.String())
	c, ok = m.s.Fetch(key)
	return
}

func (m *NGramModel) incr(prefix, numerator, denominator string, incr int64) {
	key := fmt.Sprintf("%v:%v", prefix, numerator)
	m.s.IncrN(key, denominator, incr)
}

func (m *NGramModel) ngramLookup(ngram NGram) (c counter.Interface, ok bool) {
	n := len(ngram)
	if n > m.n {
		panic(fmt.Sprintf("ngram must be %d or shorter. Got %v", m.n, ngram))
	}
	return m.fetch("ngram", ngram)
}

func (m *NGramModel) classLookup(ngram NGram) (c counter.Interface, ok bool) {
	return m.fetch("class", ngram)
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

func (m *NGramModel) Estimate(ngram NGram) distribution.Interface {
	c, ok := m.classLookup(ngram)
	if !ok {
		c = counter.New()
	}
	return distribution.NewLaplacian(c)
}

func (m *NGramModel) Train(datum *Datum) {
	m.incrPrior(datum.Class, datum.Count)
	for n := 1; n <= m.n; n++ {
		for _, ngram := range Generate(datum.Features, n) {
			m.incrNGram(ngram, datum.Count)
			m.incrClasses(ngram, datum.Class, datum.Count)
		}
	}
}

func (m *NGramModel) Classify(features []string) (estimator distribution.Interface, explain map[string]interface{}) {
	explain = make(map[string]interface{})
	estimator, _ = m.Prior()
	explain["prior"] = distribution.JSON(estimator)
	for _, ngram := range Generate(features, m.n) {
		ngram_est := m.Estimate(ngram)
		estimator = distribution.Multiply(estimator, ngram_est)
		explain[ngram.String()] = distribution.JSON(ngram_est)
	}
	return
}
