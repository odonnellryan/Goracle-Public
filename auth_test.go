package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	//"log"
	//"encoding/json"
	//"fmt"
)



func TestSendMissingCredentialsHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("TestSendMissingCredentialsHeader error: %s", err)
	}
	w := httptest.NewRecorder()
	SendMissingCredentialsHeader(w, req)
	if w.Code != 401 {
		t.Errorf("TestSendMissingCredentialsHeader Authorization allowed?: %s, %s", w.Code, w.Body)
	}
}

func TestCheckCredentials(t *testing.T) {
	var badPassword = map[string]string{
		"ryanb": "testb",
	}
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("TestCheckCredentials error: %s", err)
	}
	if CheckCredentials(req) {
		t.Errorf("TestCheckCredentials accepted incorrectly?: %s", req)
	}
	// test with a bad username and bad password
	req.SetBasicAuth("ryanb", badPassword["ryanb"])
	if CheckCredentials(req) {
		t.Errorf("TestCheckCredentials accepted incorrectly?: %s", req)
	}
	// test with just a bad password
	req.SetBasicAuth("ryan", badPassword["ryanb"])
	if CheckCredentials(req) {
		t.Errorf("TestCheckCredentials accepted incorrectly?: %s", req)
	}
	// test with the known good username and password
	req.SetBasicAuth("ryan", Password["ryan"])
	if !(CheckCredentials(req)) {
		t.Errorf("TestCheckCredentials good password not accepted?: %s", req)
	}
}

