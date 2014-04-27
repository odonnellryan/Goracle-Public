package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	//"log"
	"encoding/json"
	//"fmt"
)

// this should probably be expanded on. tests to ensure
// api is returning a docker host
func TestReturnDockerHost(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("TestReturnDockerHost error: %s", err)
	}
	host := Host{}
	w := httptest.NewRecorder()
	ReturnDockerHost(w, req)
	//fmt.Printf("%d - %+v", w.Code, w.Body)
	decode := json.NewDecoder(w.Body)
	err = decode.Decode(&host)
	if err != nil {
		t.Errorf("TestReturnDockerHost json: %s", err)
	}
	if host.Containers != 0 {
		t.Errorf("TestReturnDockerHost Containers are not zero: %d - %+v", w.Code, host)
	}
}
