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
)

func Generate(terms []string, n int) (ngrams [][]string) {
	if len(terms) == 0 {
		return [][]string{}
	}
	ngrams = make([][]string, 0)
	for i := 0; i < len(terms)+n-1; i++ {
		ngrams = append(ngrams, getNgram(terms, i, n))
	}
	return
}

func getNgram(terms []string, pos, n int) (ngram []string) {
	ngram = make([]string, 0)
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

type NGramModel struct {
	n int
	s store.Interface
}

var _ Interface = new(NGramModel)

func NewNGramModel(n int) *NGramModel {
	return &NGramModel{n, store.NewRedisStore()}
}

func (m *NGramModel) Prior(class string) (c counter.Interface, ok bool) {
	return nil, false
}

func (m *NGramModel) Lookup(ngram []string) (c counter.Interface, ok bool) {
	n := len(ngram)
	if n > m.n {
		panic(fmt.Sprintf("ngram must be %d or shorter. Got %v", m.n, ngram))
	}
	c, ok = m.s.Fetch(ngramToStr(ngram))
	return
}

func ngramToStr(ngram []string) string {
	return strings.Join(ngram, " ")
}

func (m *NGramModel) incrN(ngram []string, incr int64) {
	n := len(ngram)
	fmt.Println("ngram: ", ngram)
	denominator := ngramToStr(ngram)
	var numerator string
	if len(ngram) > 1 {
		numerator = ngramToStr(ngram[:n-1])
	} else {
		numerator = ""
	}
	m.s.IncrN(numerator, denominator, incr)
}

func (m *NGramModel) Estimate(ngram []string) distribution.Interface {
	c, ok := m.Lookup(ngram)
	if !ok {
		c = counter.New()
	}
	return distribution.NewLaplacian(c)
}

func (m *NGramModel) Train(datum *Datum) {
	for n := 1; n <= m.n; n++ {
		for _, ngram := range Generate(datum.Features, n) {
			fmt.Println("generated", ngram, "for", datum.Features)
			m.incrN(ngram, datum.Count)
		}
	}
}

func (m *NGramModel) Classify(features []string) (estimator distribution.Interface, explain map[string]interface{}) {
	return
}
