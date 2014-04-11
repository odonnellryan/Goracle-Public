package main

import (
	"testing"
	//"bytes"
	//"fmt"
)

var (
	testMongoEndpointErr1 = "1.1.1.199/database"
	testMongoEndpointErr2 = "1.1.1.1:99database"
	testMongoEndpoint = "1.1.1.1:99/database"
	testString = "what,is,this?"
	testSlice = []string{"what","is","this?",}
)

type testMongoHostData struct {
		testHostP string
		testPortP string
		TestDatabaseP string
}

func TestCommaStringToSlice(t *testing.T) {
	returnedSlice := CommaStringToSlice(testString)
	for index := range(returnedSlice) {
		if 	returnedSlice[index] != testSlice[index] {
			t.Errorf("CommaStringToSlice error expected: %+v got: %+v",
						testSlice, returnedSlice)

		}
	}
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
	host, port, db, err = ParseMongoEndpoint(testMongoEndpointErr1)
	if err == nil {
		t.Errorf("ParseMongoEndpoint error not thrown host/port/db: %s / %s / %s", host, port, db)
	}
	host, port, db, err = ParseMongoEndpoint(testMongoEndpointErr2)
	if err == nil {
		t.Errorf("ParseMongoEndpoint error not thrown host/port/db: %s / %s / %s", host, port, db)
	}
}