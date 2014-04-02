package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"log"
	"encoding/json"
	"fmt"
)

func TestReturnDockerHost(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
    	log.Fatal(err)
	}
	host := Host{}
	w := httptest.NewRecorder()
	ReturnDockerHost(w, req)
	if err != nil {
		t.Errorf("TestReturnDockerHost error: %s", err)
	}
	fmt.Printf("%d - %+v", w.Code, w.Body)
	decode := json.NewDecoder(w.Body)
	err = decode.Decode(&host)
	if err != nil {
		t.Errorf("TestReturnDockerHost json: %s", err)
	}
	fmt.Printf("%d - %+v", w.Code, host)
}