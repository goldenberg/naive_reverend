package corpus

import (
	model "naive_reverend/model"
)

type Corpus struct {
	Name  string
	Model model.Interface
}

func NewCorpus(name string) *Corpus {
	return *Corpus{name, model.NewNGramModel(3)}
}

/*
 * Train all the models for this corpus
 */
func (c *Corpus) Train(d *model.Datum) {
	c.Model.Train(d)
}

func (c *Corpus) String() string {
	return c.Name
}

type CorpusCollection map[string]Corpus
