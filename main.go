package main

import (
	`fmt`
	`strings`
	`bufio`
	`json`
	`os`
	`flag`
	`io`
)

var (
	features = "features"
	label    = "label"
	ALL      = "_"
)

type DataPoint struct {
	Features []string
	Label    string
}

type Key struct {
	feature string
	label   string
}

func (k *Key) ToStrKey() string {
	return k.feature + "," + k.label
}

func NewFromStrKey(strKey string) Key {
	vals := strings.SplitN(strKey, ",", 1)
	return Key{vals[0], vals[1]}
}

type KeyValue struct {
	key Key
	val int
}

func main() {
	flag.Parse()
	dataChan := make(chan *DataPoint, 100)
	s := NewRAMStore()
	quit := make(chan bool)
	quitIncr := make(chan bool)
	go ReadData(os.Stdin, dataChan, quit)
	go IncrementKeys(dataChan, quitIncr, s)
	<-quit
	<-quitIncr
	fmt.Println("Store :", s)
	fmt.Println("Store len:", len(s.table))
	for {
		
	}
}

func ReadData(reader io.Reader, out chan *DataPoint, quit chan bool) {
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
		var dataPt DataPoint
		err = json.Unmarshal(line, &dataPt)
		if err != nil {
			fmt.Print(err)
		}
		out <- &dataPt
		i += 1
	}
	fmt.Println("Processed", i, "lines")
	close(out)
	quit <- true
	return
}

func IncrementKeys(dataChan chan *DataPoint, quit chan bool, s KVStore) {
	// incrChan := make(chan *KeyValue)
	for {
		x := <- dataChan
		if x == nil {
			fmt.Println("nil")
			break
		}
		for _, feature := range x.Features {
			go s.Incr(&KeyValue{Key{feature, x.Label}, 1})
			go s.Incr(&KeyValue{Key{feature, ALL}, 1})
			go s.Incr(&KeyValue{Key{ALL, ALL}, 1})
		}
		go s.Incr(&KeyValue{Key{ALL, x.Label}, 1})
	}
	quit <- true
}

func KeysToIncr(x *DataPoint, incrChan chan *KeyValue, quit chan bool) {
	for _, feature := range x.Features {
		incrChan <- &KeyValue{Key{feature, x.Label}, 1}
		incrChan <- &KeyValue{Key{feature, ALL}, 1}
		incrChan <- &KeyValue{Key{ALL, ALL}, 1}
	}
	incrChan <- &KeyValue{Key{ALL, x.Label}, 1}
	quit <- true
	return
}

type KVStore interface {
	Incr(kv *KeyValue)
	Lookup(key Key) int
}

type RAMStore struct {
	table map[string]int
}

func NewRAMStore() *RAMStore {
	s := new(RAMStore)
	s.table = make(map[string]int)
	return s
}

func (s RAMStore) Incr(kv *KeyValue) {
	strkey := kv.key.ToStrKey()
	s.table[strkey] = s.table[strkey] + kv.val
}

func (s RAMStore) Lookup(key Key) int {
	return s.table[key.ToStrKey()]
}

// func (s RAMStore) TopN() []KeyValue {
	
// }
