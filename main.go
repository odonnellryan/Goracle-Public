package main

import (
	"crypto/rsa"
	"encoding/json"
	//"io"
	"crypto/rand"
	"log"
	"net/http"
	"strings"
)

type Deployment struct {
	Username      string
	ContainerName string
	Image         string
	Memory        string
	Hostname      string
	Cmd           string
	KeyPair       *rsa.PrivateKey
}

func HandleDeploymentRequest(w http.ResponseWriter, r *http.Request) {

	// request will be structured as such:
	// /deployments/{user_name}/{request_server}/{container_name}/{image}/{GET_Params}
	// {GET_Params} (so far...): memory=(string), hostname=(string), cmd=(string)
	// C: deployments/test/server/container/image/?memory=test&host=host&cmd=cmd
	Url := strings.Split(r.URL.Path, "/")

	//todo: error checking. currently goes mad if you miss a value, of course

	// :O probably an awful idea huh!
	PrivateKey, Error := rsa.GenerateKey(rand.Reader, 1024)
	// this just matches the dict above, i think it works OK.
	// only problem - index out of range thing on the below IF you don't set
	// all the correct URL params, so we need to turn an error in that case
	// ideally, things in the PATH should be REQUIRED and GET values should
	// be optional/return defaults.
	DeploymentValues := Deployment{
		Url[2], Url[4], Url[5], r.FormValue("memory"), r.FormValue("hostname"),
		r.FormValue("cmd"), PrivateKey,
	}

	DeploymentJson, JError := json.Marshal(DeploymentValues)

	if Error != nil && JError != nil {
		log.Println("Request from user:")
	}
	//testing! works kinda ;)
	w.Write([]byte(DeploymentJson))

}

func main() {
	mux := &MethodServerMux{make(map[string]*http.ServeMux)}

	// Add handlers here
	mux.AddHandler("GET", "/deployments/", HandleDeploymentRequest)

	// Bottom is hit first, then second to last, etc
	secure := AuthorizationRequired(mux.HandleRequest)
	logged := LogRequest(secure)

	http.HandleFunc("/", logged)

	log.Println("Starting")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
