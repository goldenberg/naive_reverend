package main

import (
	"strings"
	"strconv"
	"encoding/json"
	// corpus "naive_reverend/corpus"
	model "naive_reverend/model"
	// store "naive_reverend/store"
	distribution "naive_reverend/distribution"
	"net/http"
)

const (
	CORPUSES_KEY = "__CORPUSES__"
)

type TrainHandler struct {

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

	count := strconv.Atoi(req.FormValue("count"))
	if err != nil {
		count = 1
	}

	m := model.NewNGramModel(n, corpus)

	features := strings.Split(query, ",")
	d := model.Datum{class, features, count}
	m.Train(d)
	return
}
