package main

import (
	"encoding/json"
	"strconv"
	"strings"
	// corpus "naive_reverend/corpus"
	"fmt"
	distribution "naive_reverend/distribution"
	model "naive_reverend/model"
	store "naive_reverend/store"
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
	corpus := req.FormValue("corpus")
	n, err := strconv.Atoi(req.FormValue("n"))
	if err != nil {
		n = 2
	}
<<<<<<< HEAD
	fmt.Println(req)
	m := model.NewNGramModel(h.s, n, corpus)
=======
	m := model.NewNGramModel(h.pool.Get(corpus), n)
>>>>>>> df7ecd4... initial work toward pipelining. might be misguided

	features := strings.Split(query, ",")
	estimator, explain := m.Classify(features)
	prediction, _ := distribution.ArgMax(estimator)
	output := map[string]interface{}{
		"prediction": prediction,
		"estimator":  distribution.JSON(estimator),
		"explain":    explain,
	}
	jsonBytes, _ := json.MarshalIndent(output, "", "\t")
	w.Write(jsonBytes)
	return
}
