package main

import (
	`fmt`
	`bufio`
	`json`
	`os`
	`flag`
	`io`
	`rand`
	model `naive_reverend/model`
)

func main() {
	flag.Parse()
	data := make(chan *model.Datum, 100)
	quit := make(chan bool)

	trainData := make(chan *model.Datum, 1000000)
	evalData := make(chan *model.Datum, 1000000)

	go ReadData(os.Stdin, data, quit)

	nb := model.New()

	for d := range data {
		if rand.Float32() < 0.9 {
			nb.Train(d)
		} else {
			evalData <- d
		}
	}
	close(trainData)
	close(evalData)

	var correct, wrong uint
	for d := range evalData {
		class, _ := nb.Classify(d.Features)
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
}

func ReadData(reader io.Reader, out chan *model.Datum, quit chan bool) {
	bufReader, _ := bufio.NewReaderSize(reader, 1000000000)
	i := 0
	for {
		line, isPrefix, err := bufReader.ReadLine()
		if err != nil {
			break
		}
		if isPrefix {
			fmt.Print("uh-oh")
			break
		}
		var x model.Datum
		err = json.Unmarshal(line, &x)
		if err != nil {
			fmt.Print(err)
		}
		out <- &x
		i += 1
	}
	fmt.Println("Processed", i, "lines")
	close(out)
	quit <- true
	return
}
