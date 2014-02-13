package main

import (
	"testing"
	"io/ioutil"
	"fmt"
	"net/http"
)

var testHost = Host {
	Hostname: "local_testing",
    Address: "http://127.0.0.1:8888/",
    User: "ryan",
    Password: "test",
}

func TestSendDockerCommand(t *testing.T) {
	resp, err := SendDockerCommand(testHost, "images/json", "GET")
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	fmt.Printf("Response: %s", resp)
}

func TestHTTPToDocker(t *testing.T) {
	host := testHost
	client := &http.Client{}
	resp, err := client.Get(host.Address)
	if err != nil {
		t.Errorf("Client error: %s", err)
	}
	// closes the connection
	defer resp.Body.Close()
	request, err := http.NewRequest("GET", (host.Address + "images/json"), nil)
	if err != nil {
		t.Errorf("Request build error: %s", err)
	}
	request.SetBasicAuth(host.User, host.Password)
	response, err := client.Do(request)
	if err != nil {
		t.Errorf("Response error: %s", err)
	}
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Read response error: %s", err)
	}
	fmt.Printf("Response: %s", res)
}
