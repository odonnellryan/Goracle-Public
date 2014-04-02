package main

import (
	"testing"
	//"fmt"
)

var (
	testMongoEndpoint = "1.1.1.1:99/database"

)

type testMongoHostData struct {
		testHostP string
		testPortP string
		TestDatabaseP string
}


func TestParseMongoEndpoint(t *testing.T) {
	testData := testMongoHostData{"1.1.1.1", "99", "database"}
	host, port, db, err := ParseMongoEndpoint(testMongoEndpoint)
	if err != nil {
		t.Errorf("ParseMongoEndpoint error: %s host returned: %s", err, host)
	}
	items := testMongoHostData{host,port,db}
	if testData != items {
		t.Errorf("ParseMongoEndpoint error expected: %+v got: %+v", testData, items)
	}
}