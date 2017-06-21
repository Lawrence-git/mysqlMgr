package main

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"time"
)

type dataElement struct {
	exampleDataOne   string
	exampleDataTwo   string
	exampleDataThree int32
	exampleDataFour  int32
}

type randomDataPool struct {
	words []string
	sz    int64
}

func (r *randomDataPool) Load() (e error) {
	srcFile := "/usr/share/dict/words"
	fB, e := ioutil.ReadFile(srcFile)
	if e != nil {
		return
	}
	tmpByte := bytes.Split(fB, []byte("\n"))
	for _, next := range tmpByte {
		r.words = append(r.words, string(next))
		r.sz = r.sz + 1
	}
	return
}

func (r *randomDataPool) NewDataElement() (d dataElement) {
	rand.Seed(time.Now().UTC().UnixNano())
	d.exampleDataOne = r.words[rand.Int63n(r.sz-1)]
	d.exampleDataTwo = r.words[rand.Int63n(r.sz-1)]
	d.exampleDataThree = rand.Int31()
	d.exampleDataFour = rand.Int31()
	return
}
