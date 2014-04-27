package main

import (
	"testing"
	//"bytes"
	//"fmt"
)

var (
	testMongoEndpointErr1 = "1.1.1.199/database"
	testMongoEndpointErr2 = "1.1.1.1:99database"
	testMongoEndpoint     = "1.1.1.1:99/database"
	testString            = "what,is,this?"
	testSlice             = []string{"what", "is", "this?"}
	testHexValue          = "ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff"
	testHexString         = "test"
)

type testMongoHostData struct {
	testHostP     string
	testPortP     string
	TestDatabaseP string
}

func TestCommaStringToSlice(t *testing.T) {
	returnedSlice := CommaStringToSlice(testString)
	for index := range returnedSlice {
		if returnedSlice[index] != testSlice[index] {
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
	items := testMongoHostData{host, port, db}
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

func TestCryptToHex(t *testing.T) {
	hex := CryptToHex(testHexString)
	if hex != testHexValue {
		t.Errorf("testCryptToHex error expected: %+v got: %+v",
			testHexValue, hex)
	}
}
