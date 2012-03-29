package main

import (
	"strings"
	"strconv"
	"encoding/json"
	// corpus "naive_reverend/corpus"
	model "naive_reverend/model"
	store "naive_reverend/store"
	distribution "naive_reverend/distribution"
	"net/http"
)

const (
	CORPUSES_KEY = "__CORPUSES__"
)

type ClassifyHandler struct {
	pool *store.Pool
}

func (h ClassifyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	query := req.FormValue("q")
	corpus := req.FormValue("c")
	n, err := strconv.Atoi(req.FormValue("n"))
	if err != nil {
		n = 2
	}
	m := model.NewNGramModel(h.pool.Get(corpus), n)

	features := strings.Split(query, ",")
	estimator, explain := m.Classify(features)
	prediction, _ := distribution.ArgMax(estimator)
	output := map[string]interface{}{
		"prediction": prediction,
		"estimator":  distribution.JSON(estimator),
		"explain":    explain,
	}
	jsonWriter := json.NewEncoder(w)
	jsonWriter.Encode(output)
	return
}
