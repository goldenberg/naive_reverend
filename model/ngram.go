package model 

import (
	"fmt"
	"strings"
	counter "naive_reverend/counter"
	store "naive_reverend/store"
	distribution "naive_reverend/distribution"
)

func Generate(terms []string, n int) (ngrams [][]string) {
	ngrams = make([][]string, len(terms) - n + 1)
	for i := 0; i < len(ngrams); i++ {
		if i+n <= len(terms) {
			ngrams = append(ngrams, terms[i:i+n])
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
	// var stores = make([]store.Interface, n)
	// for i := 1; i <= n; i++ {
	// 	stores[i] = storeFactory()
	// }
	return &NGramModel{n, store.NewRedisStore()}
}

func (m *NGramModel) Prior(class string) (c counter.Interface,  ok bool) {
	return nil, false
}

func (m *NGramModel) Lookup(ngram []string) (c counter.Interface,  ok bool) {
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
	for n := 1; n <=m.n; n++ {
		for _, ngram := range Generate(datum.Features, n) {
			fmt.Println("generated", ngram, "for", datum.Features)
			m.incrN(ngram, datum.Count)
		} 		
	}
}

func (m *NGramModel) Classify(features []string) (estimator distribution.Interface, explain map[string]interface{}) {
	return	
}

