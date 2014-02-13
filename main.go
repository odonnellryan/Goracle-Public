package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"strconv"
)

func ReturnDockerHost(w http.ResponseWriter, r *http.Request) {
	host, err := SelectHost()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	w.Write([]byte(fmt.Sprintf("%s", host)))
}

func HandleDeploymentRequest(w http.ResponseWriter, r *http.Request) {
	// Request will be structured as such:
	// Checking...
	host, err := SelectHost()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	// convert to int64
	memory, err := strconv.ParseInt(r.FormValue("memory"), 10, 64)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	memory_swap, err := strconv.ParseInt(r.FormValue("memory_swap"), 10, 64)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	cpu, err := strconv.ParseInt(r.FormValue("cpu"), 10, 64)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	d := Deployment {
		User:        	r.FormValue("user"),
		ContainerName:  r.FormValue("container_name"),
		Image:       	r.FormValue("image"),
		Memory:      	memory, // number in bytes
		MemorySwap:	    memory_swap, // number in bytes for memory + swap, -1 for no swap
		CPU:         	cpu,
		Command:     	CommaStringToSlice(r.FormValue("command")),
		IP:          	r.FormValue("ip"),
		ExposedPorts:   CommaStringToSlice(r.FormValue("exposed_ports")),
	}
	response := DeployNewContainer(host, d, r)

	// Testing! Works kinda.
	w.Write([]byte(response))
}

func ParseMongoEndpoint(endpoint string) (string, string, string, error) {
	portSeperatorIndex := strings.LastIndex(endpoint, ":")
	if portSeperatorIndex < 1 {
		err := fmt.Errorf("Invalid format of mongo endpoint. Correct syntax is host:port/db")
		return "", "0", "", err
	}
	host := endpoint[0 : portSeperatorIndex-1]
	slashIndex := strings.Index(endpoint, "/")
	if slashIndex < 1 {
		err := fmt.Errorf("Invalid format of mongo endpoint. Correct syntax is host:port/db")
		return "", "0", "", err
	}
	port := endpoint[portSeperatorIndex+1 : slashIndex-1]
	db := endpoint[slashIndex+1:]
	return host, port, db, nil
}

func main() {
	var mongoHost, loadBalancerHost string
	flag.StringVar(&mongoHost, "mongoserver", "", "")
	flag.StringVar(&loadBalancerHost, "loadbalancer", "", "")
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
	}

	dispatcher := &RequestDispatcher{make(map[string]*http.ServeMux)}

	// Add handlers here
	dispatcher.AddHandler("GET", "/deployments/", HandleDeploymentRequest)
	dispatcher.AddHandler("GET", "/docker_pool/", ReturnDockerHost)

	// Bottom is hit first, then second to last, etc
	secure := AuthorizationRequired(dispatcher.HandleRequest)
	logged := LogRequest(secure)

	http.HandleFunc("/", logged)

	log.Println("Starting")
	log.Fatal(http.ListenAndServe(HostPort, nil))
	// flag.Parse()
}
