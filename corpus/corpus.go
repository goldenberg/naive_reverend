package corpus

import (
	model "naive_reverend/model"
)

type Corpus struct {
	Name   string
	Models []model.Interface
}

func NewCorpus(name string) *Corpus {
	return &Corpus{name, make([]model.Interface, 1)}
}

/*
 * Train all the models for this corpus
 */
func (c *Corpus) Train(d *model.Datum) {
	for _, m := range c.Models {
		m.Train(d)
	}
}

func (c *Corpus) String() string {
	return c.Name
}

type CorpusCollection map[string]Corpus
