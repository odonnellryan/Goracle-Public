package main

import (
	//"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	//"strconv"
)

func main() {
	var mongoHost, loadBalancerHost string
	flag.StringVar(&mongoHost, "mongoserver", "", "host:port/db")
	flag.StringVar(&loadBalancerHost, "loadbalancer", "", "host")
	flag.Parse()
	if mongoHost == "" {
		fmt.Println("Mongo endpoint not specified.")
		flag.PrintDefaults()
		return
	}
	if loadBalancerHost == "" {
		fmt.Println("Load balancer endpoint not specified.")
		flag.PrintDefaults()
		return
	}
	var err error
	MongoDBAddress, MongoDBPort, MongoDBName, err = ParseMongoEndpoint(mongoHost)
	if err != nil {
		fmt.Println(err)
		return
	}
	dispatcher := &RequestDispatcher{make(map[string]*http.ServeMux)}
	// Add handlers here
	dispatcher.AddHandler("POST", "/deployments/", HandleDeploymentRequest)
	//dispatcher.AddHandler("GET", "/docker_pool/", DockerHost)
	// Bottom is hit first, then second to last, etc
	secure := AuthorizationRequired(dispatcher.HandleRequest)
	logged := LogRequest(secure)
	http.HandleFunc("/", logged)
	log.Println("Starting")
	log.Fatal(http.ListenAndServe(HostPort, nil))
	// flag.Parse()
}
