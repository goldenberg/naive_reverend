package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	model "naive_reverend/model"
	counter "naive_reverend/counter"
	"net/http"
	"os"
	"strings"
	_ "net/http/pprof"
)

func main() {
	flag.Parse()
	data := make(chan *model.Datum, 100)
	quit := make(chan bool)
	quitServer := make(chan bool)

	trainData := make(chan *model.Datum, 100)
	evalData := make(chan *model.Datum, 200)

	nb := model.New()

	go ServeDebug()
	go ReadData(os.Stdin, data, quit)
	go Serve(nb, quitServer)

	for d := range data {
		if rand.Float32() < 0.95 {
			nb.Train(d)
		} else {
			evalData <- d
		}
	}
	close(trainData)
	close(evalData)

	var correct, wrong uint
	for d := range evalData {
		estimator, _ := nb.Classify(d.Features)
		class, _ := counter.ArgMax(estimator)
		fmt.Println("Was:", d.Class, "Got:", class)
		if class == d.Class {
			correct += 1
		} else {
			wrong += 1
		}
	}

	accuracy := float64(correct) / (float64(correct) + float64(wrong))
	fmt.Println(accuracy*100., "Got", correct, "correct and", wrong, "wrong.")

	<-quit
	<-quitServer
}

func Serve(nb *model.NaiveBayes, quit chan bool) {
	var cHandler ClassifyHandler
	cHandler.nb = nb

	http.HandleFunc("/hello", HelloServer)
	http.Handle("/classify", cHandler)

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	quit <- true
}

func ServeDebug() {
	fmt.Println("starting debug server")
	err := http.ListenAndServe(":6060", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!")
}

type ClassifyHandler struct {
	nb *model.NaiveBayes
}

func (h ClassifyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	query := req.FormValue("q")
	features := strings.Split(query, ",")
	estimator, explain := h.nb.Classify(features)
	prediction, _ := counter.ArgMax(estimator)
	output := map[string]interface{}{
		"prediction":  prediction,
		"estimator": counter.JSON(estimator),
		"explain": explain,
	}
	jsonWriter := json.NewEncoder(w)
	jsonWriter.Encode(output)
	return
}

func ReadData(reader io.Reader, out chan *model.Datum, quit chan bool) {
	jsonDecoder := json.NewDecoder(reader)
	i := 0
	for {
		var x model.Datum
		err := jsonDecoder.Decode(&x)
		if err != nil {
			fmt.Println(err)
			break
		}
		out <- &x
		i += 1
	}
	fmt.Println("Processed", i, "lines")
	close(out)
	quit <- true
	return
}
