package main

import (
	"log"
	"net/http"
)

type Deployment struct {
	ContainerName    string
	ContainerPackage string
	ContainerAuthKey string
	IpAddress        string
}

func HandleDeploymentRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Deployments and stuff"))
}

func main() {
	mux := &MethodServerMux{make(map[string]*http.ServeMux)}

	// Add handlers here
	mux.AddHandler("GET", "/deployments", HandleDeploymentRequest)

	// Bottom is hit first, then second to last, etc
	secure := AuthorizationRequired(mux.HandleRequest)
	logged := LogRequest(secure)

	http.HandleFunc("/", logged)

	log.Println("Starting")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
