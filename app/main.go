package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	model "naive_reverend/model"
	distribution "naive_reverend/distribution"
	store "naive_reverend/store"
	"net/http"
	"os"
	_ "net/http/pprof"
	"time"
	pprof "runtime/pprof"
	"bufio"
)

func main() {
	profile := flag.Bool("p", false, "write profiles to ./")
	train := flag.String("t", "", "train using the data in this file")
	evaluate := flag.String("e", "", "train using the data in this file")
	ngram := flag.Int("n", 2, "ngram length")
	corpus := flag.String("c", "", "corpus name")

	flag.Parse()

	trainData := make(chan *model.Datum, 100)
	evalData := make(chan *model.Datum, 100)
	quit := make(chan bool)

	pool := store.NewPool(store.NewRedisStore)
	nb := model.NewNGramModel(pool.Get(*corpus), *ngram)

	quitServer := make(chan bool)

	go ServeDebug()
	go Serve(nb, pool, quitServer)

	if *train != "" {
		fmt.Println("Training on", *train)
		f, _ := os.Open(*train)
		go ReadData(bufio.NewReader(f), trainData, quit)
		go func() {
			for d := range trainData {
				nb.Train(d)
			}
		}()
		<-quit
	}

	if *evaluate != "" {
		fmt.Println("Evaluating on", *evaluate)
		f, _ := os.Open(*evaluate)
		go ReadData(bufio.NewReader(f), evalData, quit)
		go Evaluate(evalData, nb)
	}

	if *profile {
		DumpProfiles()
	}

	<-quit
	<-quitServer
}

func Evaluate(evalData chan *model.Datum, nb model.Interface) {
	var correct, wrong uint
	evalStartTime := time.Now()
	for d := range evalData {
		estimator, _ := nb.Classify(d.Features)
		class, _ := distribution.ArgMax(estimator)
		if class == d.Class {
			correct += 1
		} else {
			fmt.Println("For:", d.Features, "Was:", d.Class, "Got:", class)
			wrong += 1
		}
	}

	elapsed := time.Since(evalStartTime).Seconds()
	fmt.Println("Took", elapsed, "for", correct+wrong, "queries.", elapsed/float64(correct+wrong), "sec/query")
	accuracy := float64(correct) / (float64(correct) + float64(wrong))
	fmt.Println(accuracy*100., "Got", correct, "correct and", wrong, "wrong.")
}

func DumpProfiles() {
	f, err := os.Create("./memprofile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(f)
	f.Close()
	fmt.Println("Wrote memprofile")
	return
}

func Serve(nb model.Interface, p *store.Pool, quit chan bool) {
	fmt.Println("serving")
	http.HandleFunc("/hello", HelloServer)

	// Will eventually be /corpus/classify, /corpus/train, and /corpus/params
	http.Handle("/classify", ClassifyHandler{p})
	http.Handle("/train", TrainHandler{p})

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
		if x.Count == 0 {
			x.Count = 1
		}
		out <- &x
		i += 1
	}
	fmt.Println("Processed", i, "lines")
	close(out)
	quit <- true
	return
}
