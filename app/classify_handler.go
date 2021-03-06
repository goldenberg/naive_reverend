package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"fmt"
	distribution "github.com/goldenberg/naive_reverend/distribution"
	model "github.com/goldenberg/naive_reverend/model"
	store "github.com/goldenberg/naive_reverend/store"
	"net/http"
)

const (
	CORPUSES_KEY = "__CORPUSES__"
)

type ClassifyHandler struct {
	pool *store.Pool
}

// ServeHTTP classifies 
func (h ClassifyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	query := req.FormValue("q")
	corpus := req.FormValue("corpus")
	n, err := strconv.Atoi(req.FormValue("n"))
	if err != nil {
		n = 2
	}
	fmt.Println(req)
	m := model.NewNGramModel(h.pool.Get(corpus), n)

	features := strings.Split(query, " ")
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
