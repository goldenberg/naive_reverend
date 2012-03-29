package main

import (
	model "naive_reverend/model"
	store "naive_reverend/store"
	"net/http"
	"strconv"
	"strings"
)

type TrainHandler struct {
	s store.Interface
}

func (h TrainHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	query := req.FormValue("q")
	class := req.FormValue("class")
	corpus := req.FormValue("corpus")

	n, err := strconv.Atoi(req.FormValue("n"))
	if err != nil {
		n = 2
	}

	count, err := strconv.Atoi(req.FormValue("count"))
	if err != nil {
		count = 1
	}

	m := model.NewNGramModel(h.s, n, corpus)

	features := strings.Split(query, ",")
	d := &model.Datum{class, features, int64(count)}
	m.Train(d)
	return
}
