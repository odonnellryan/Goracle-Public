package main

import (
	"testing"
	"fmt"
	"net/http"
	"net/http/httptest"
	)


func testRoute(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("testing the middleware"))
}

func TestHandlers(t *testing.T) {
	testDispatcher := &RequestDispatcher{make(map[string]*http.ServeMux)}
	testDispatcher.AddHandler("POST", "/test/", testRoute)
	logged := LogRequest(testDispatcher.HandleRequest)
	server := httptest.NewServer(logged)
	defer server.Close()
	client := &http.Client{}
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Errorf("TestHandlers client Get error: %s", err)
	}
	req, err := http.NewRequest("POST", (server.URL+"/test/"), nil)
	if err != nil {
		t.Errorf("TestHandlers NewRequest error: %s", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("TestHandlers client do error: %s", err)
	}
	fmt.Printf("test handler response %+v", resp)
}

func TestHandleRequest(t *testing.T) {

}

func TestLogRequest(t *testing.T) {

}