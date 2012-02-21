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
		ngrams = append(ngrams, terms[i:i+n])
	} 
	return
} 

type NGramModel struct {
	n int
	ngramStores []store.Interface
}

var _ Interface = new(NGramModel)

func NewNGramModel(n int, storeFactory func() store.Interface) *NGramModel {
	var stores = make([]store.Interface, n)
	for i := 1; i <= n; i++ {
		stores[i] = storeFactory()
	}
	return &NGramModel{n, stores}
}

func (m *NGramModel) Prior(class string) (c counter.Interface,  ok bool) {
	return nil, false
}

func (m *NGramModel) Lookup(ngram []string) (c counter.Interface,  ok bool) {
	n := len(ngram)
	if n > m.n {
		panic(fmt.Sprintf("ngram must be %d or shorter. Got %v", m.n, ngram))
	}
	c, ok = m.ngramStores[n-1].Fetch(strings.Join(ngram, " "))
	return
}

func (m *NGramModel) Estimate(ngram []string) distribution.Interface {
	return nil
}

func (m *NGramModel) Train(datum *Datum) {
	return
}

func (m *NGramModel) Classify(features []string) (estimator distribution.Interface, explain map[string]interface{}) {
	return	
}

