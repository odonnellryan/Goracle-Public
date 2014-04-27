package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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
		t.Errorf("TestCheckCredentials http error: %s", err)
	}
	if CheckCredentials(req) {
		t.Errorf("TestCheckCredentials accepted incorrectly?: %s", req)
	}
	// test with a bad username and bad password
	req.SetBasicAuth("ryanb", badPassword["ryanb"])
	if CheckCredentials(req) {
		t.Errorf("TestCheckCredentials accepted bad username and bad password: %s, %s", req, CheckCredentials(req))
	}
	// no username or password
	req.SetBasicAuth("", "")
	if CheckCredentials(req) {
		t.Errorf("TestCheckCredentials accepted no username or password: %s, %s", req, CheckCredentials(req))
	}
	// no password good username
	req.SetBasicAuth("ryan", "")
	if CheckCredentials(req) {
		t.Errorf("TestCheckCredentials accepted no password: %s, %s", req, CheckCredentials(req))
	}
	// test with just a bad password
	req.SetBasicAuth("ryan", badPassword["ryanb"])
	if CheckCredentials(req) {
		t.Errorf("TestCheckCredentials accepted bad password?: %s, %s", req, CheckCredentials(req))
	}
	// test with the known good username and password
	req.SetBasicAuth("ryan", "test")
	if !(CheckCredentials(req)) {
		t.Errorf("TestCheckCredentials good username/password not accepted: %s, %s, %s", req, CheckCredentials(req), CryptToHex(Password["ryan"]))
	}
}

func TestCheckAuthorization(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("TestCheckAuthorization error: %s", err)
	}
	req.RemoteAddr = "1.1.1.1:80"
	checkAuth := CheckAuthorization(req)
	if checkAuth {
		t.Errorf("TestCheckAuthorization should not have passed. req: %+v", req)
	}
	req.RemoteAddr = "127.0.0.1:80"
	checkAuth = CheckAuthorization(req)
	if checkAuth {
		t.Errorf("TestCheckAuthorization http basic not set but passed. req: %+v, %s", req, checkAuth)
	}
	req.SetBasicAuth("ryan", Password["ryan"])
	checkAuth = CheckAuthorization(req)
	if !(checkAuth) {
		t.Errorf("TestCheckAuthorization did not pass. req: %+v, %s", req, checkAuth)
	}
}

// not yet implemented (throws errors)
func TestAuthorizationRequired(t *testing.T) {
	// set up client for control over http headers
	client := &http.Client{}
	// build our test handler to pass to different handler wrapper
	// functions
	testHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("http response"))
		})
	// the specific wrapper we're testing now
	authHandler := AuthorizationRequired(testHandler)
	// we need a test server
	server := httptest.NewServer(authHandler)
	defer server.Close()
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Errorf("TestAuthorizationRequired http client error: %s", err)
	}
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Errorf("TestAuthorizationRequired req build error: %s", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("TestAuthorizationRequired client do error: %s", err)
	}
	actualResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestAuthorizationRequired ioutil error: %s", err)
	}
	if !(resp.StatusCode == 401) {
		t.Errorf("TestAuthorizationRequired passed unexpectedly %d - %s",
			resp.StatusCode, actualResp)
	}
	// use the request we setup to modify http headers so we can
	// properly pass our test and have it pass
	req.RemoteAddr = "127.0.0.1:80"
	req.SetBasicAuth("ryan", Password["ryan"])
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("TestAuthorizationRequired client do error: %s", err)
	}
	actualResp, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("TestAuthorizationRequired ioutil error: %s", err)
	}
	if resp.StatusCode == 401 {
		t.Errorf("TestAuthorizationRequired passed unexpectedly %d - %s",
			resp.StatusCode, actualResp)
	}
}
