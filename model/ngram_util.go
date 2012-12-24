package model

import (
	"strings"
)

const BLANK = "_"

type NGram []string

/* Generate computes n-grams using a sliding window of size n. 
The terms are pre-pended and extended with n - 1 BLANKs. */
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

/* String joins the n-gram together with spaces. */
func (ng NGram) String() string {
	return strings.Join(ng, " ")
}
