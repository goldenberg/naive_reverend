package model

import (
	"testing"
)

type nGramTest struct {
	n        int
	tokens   []string
	expected []NGram
}

var nGramTests = []nGramTest{
	nGramTest{1, []string{}, []NGram{}},
	nGramTest{1, []string{"a"}, []NGram{
		NGram{"a"},
	}},
	nGramTest{1, []string{"a", "b"}, []NGram{
		NGram{"a"},
		NGram{"b"},
	}},

	nGramTest{2, []string{}, []NGram{}},
	nGramTest{2, []string{"a"}, []NGram{
		NGram{BLANK, "a"},
		NGram{"a", BLANK}}},
	nGramTest{2, []string{"a", "b"}, []NGram{
		NGram{BLANK, "a"},
		NGram{"a", "b"},
		NGram{"b", BLANK}}},
	nGramTest{2, []string{"a", "b", "c"}, []NGram{
		NGram{BLANK, "a"},
		NGram{"a", "b"},
		NGram{"b", "c"},
		NGram{"c", BLANK},
	}},

	nGramTest{3, []string{}, []NGram{}},
	nGramTest{3, []string{"a"}, []NGram{
		NGram{BLANK, BLANK, "a"},
		NGram{BLANK, "a", BLANK},
		NGram{"a", BLANK, BLANK},
	}},
	nGramTest{3, []string{"a", "b"}, []NGram{
		NGram{BLANK, BLANK, "a"},
		NGram{BLANK, "a", "b"},
		NGram{"a", "b", BLANK},
		NGram{"b", BLANK, BLANK},
	}},
	nGramTest{3, []string{"a", "b", "c"}, []NGram{
		NGram{BLANK, BLANK, "a"},
		NGram{BLANK, "a", "b"},
		NGram{"a", "b", "c"},
		NGram{"b", "c", BLANK},
		NGram{"c", BLANK, BLANK},
	}},
	nGramTest{3, []string{"a", "b", "c", "d"}, []NGram{
		NGram{BLANK, BLANK, "a"},
		NGram{BLANK, "a", "b"},
		NGram{"a", "b", "c"},
		NGram{"b", "c", "d"},
		NGram{"c", "d", BLANK},
		NGram{"d", BLANK, BLANK},
	}},
}

func TestGenerate(t *testing.T) {
	for _, ngt := range nGramTests {
		ngrams := Generate(ngt.tokens, ngt.n)

		if !nestedEqual(ngrams, ngt.expected) {
			t.Errorf("Generate(%v, %v) = %v, want %v.", ngt.tokens, ngt.n, ngrams, ngt.expected)
		}
	}
}

func nestedEqual(x, y []NGram) bool {
	if len(x) != len(y) {
		return false
	}
	for i, _ := range x {
		if len(x[i]) != len(y[i]) {
			return false
		}
		for j, _ := range x[i] {
			if x[i][j] != y[i][j] {
				return false
			}
		}
	}
	return true
}
