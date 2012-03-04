package model

import (
	"testing"
)

type nGramTest struct {
	n        int
	tokens   []string
	expected [][]string
}

var nGramTests = []nGramTest{
	nGramTest{1, []string{}, [][]string{}},
	nGramTest{1, []string{"a"}, [][]string{
		[]string{"a"},
	}},
	nGramTest{1, []string{"a", "b"}, [][]string{
		[]string{"a"},
		[]string{"b"},
	}},

	nGramTest{2, []string{}, [][]string{}},
	nGramTest{2, []string{"a"}, [][]string{
		[]string{BLANK, "a"},
		[]string{"a", BLANK}}},
	nGramTest{2, []string{"a", "b"}, [][]string{
		[]string{BLANK, "a"},
		[]string{"a", "b"},
		[]string{"b", BLANK}}},
	nGramTest{2, []string{"a", "b", "c"}, [][]string{
		[]string{BLANK, "a"},
		[]string{"a", "b"},
		[]string{"b", "c"},
		[]string{"c", BLANK},
	}},

	nGramTest{3, []string{}, [][]string{}},
	nGramTest{3, []string{"a"}, [][]string{
		[]string{BLANK, BLANK, "a"},
		[]string{BLANK, "a", BLANK},
		[]string{"a", BLANK, BLANK},
	}},
	nGramTest{3, []string{"a", "b"}, [][]string{
		[]string{BLANK, BLANK, "a"},
		[]string{BLANK, "a", "b"},
		[]string{"a", "b", BLANK},
		[]string{"b", BLANK, BLANK},
	}},
	nGramTest{3, []string{"a", "b", "c"}, [][]string{
		[]string{BLANK, BLANK, "a"},
		[]string{BLANK, "a", "b"},
		[]string{"a", "b", "c"},
		[]string{"b", "c", BLANK},
		[]string{"c", BLANK, BLANK},
	}},
	nGramTest{3, []string{"a", "b", "c", "d"}, [][]string{
		[]string{BLANK, BLANK, "a"},
		[]string{BLANK, "a", "b"},
		[]string{"a", "b", "c"},
		[]string{"b", "c", "d"},
		[]string{"c", "d", BLANK},
		[]string{"d", BLANK, BLANK},
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

func nestedEqual(x, y [][]string) bool {
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
