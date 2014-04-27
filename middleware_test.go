package main

import (
	"testing"
	//"fmt"
	"net/http"
	"net/http/httptest"
	"io/ioutil"
	"log"
	"bytes"
	"strings"
	)

var testMessage = []byte("testing the middleware")

func testRoute(w http.ResponseWriter, r *http.Request) {
	w.Write(testMessage)
}

func TestHandlers(t *testing.T) {
	// redirect logging to a buffer for testing
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
    // set up the dispatcher testing
	testDispatcher := &RequestDispatcher{make(map[string]*http.ServeMux)}
	testDispatcher.AddHandler("POST", "/test/", testRoute)
	// test the logging
	logged := LogRequest(testDispatcher.HandleRequest)
	server := httptest.NewServer(logged)
	defer server.Close()
	client := &http.Client{}
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Errorf("TestHandlers client Get error: %s \n", err)
	}
	req, err := http.NewRequest("POST", (server.URL+"/test/"), nil)
	if err != nil {
		t.Errorf("TestHandlers NewRequest error: %s \n", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("TestHandlers client.do error: %s \n", err)
	}
	actualResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestHandlers ioutil error: %s \n", err)
	}
	if !(bytes.Equal(actualResp, testMessage)) {
		t.Errorf("TestHandlers got the wrong response: %s \n", actualResp)
	}
	// kind of messy to verify:
	// so for now i'm just going to split and ensure the values are
	// the expected length for the logs
	logString := logBuffer.String()
	logSlice := strings.Split(logString, "Request")
	if !(len(logSlice) == 3) {
		t.Errorf("logTestValue seems to be off : %+v, %s, %s \n", logSlice, len(logSlice), logString)
	}
}