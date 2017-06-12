package main

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/koepkeca/mysqlMgr"
)

const (
	maxElements = 8192
)

func main() {
	rdp := randomDataPool{}
	log.Printf("Loading Random Data Pool")
	e := rdp.Load()
	if e != nil {
		panic(e)
	}
	testDataSet := []dataElement{}
	log.Printf("Generating Random Elements")
	for i := 0; i < maxElements; i++ {
		next := rdp.NewDataElement()
		testDataSet = append(testDataSet, next)
	}
	cf, e := ioutil.ReadFile("example.conf")
	if e != nil {
		panic(e)
	}
	cfByte := bytes.NewBuffer(cf)
	s, e := mysqlMgr.New(cfByte)
	if e != nil {
		panic(e)
	}
	e = s.AddStmt("test", "INSERT INTO example_table (example_data_1, example_data_2, example_data_3, example_data_4) VALUES (?,?,?,?)")
	if e != nil {
		panic(e)
	}
	e = s.AddStmt("truncate", "TRUNCATE example_table")
	log.Printf("Truncating table")
	trunc, e := s.GetStmt("truncate")
	trunc.Exec()
	stmt, e := s.GetStmt("test")
	log.Printf("Loading random data into database sequentially")
	for _, next := range testDataSet {
		stmt.Exec(next.example_data_1, next.example_data_2, next.example_data_3, next.example_data_4)
	}
	log.Printf("Done.. Shutting Down.")
	s.Close()
}
