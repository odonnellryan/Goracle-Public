package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func ReturnDockerHost(w http.ResponseWriter, r *http.Request) {
	host, err := SelectHost()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	respBytes, err := json.Marshal(host)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	w.Write(respBytes)
}

func HandleDeploymentRequest(w http.ResponseWriter, r *http.Request) {
	// 
	// post request:
	// FormValue: memory, memory_swap, user, container_name, image,
	// 				command, exposed_ports, ip
	// 
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
	memory_swap, err := strconv.ParseInt(r.FormValue("memory_swap"),
	10, 64)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	cpu, err := strconv.ParseInt(r.FormValue("cpu"), 10, 64)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	d := Deployment{
		User:          r.FormValue("user"),
		ContainerName: r.FormValue("container_name"),
		Image:         r.FormValue("image"),
		Memory:        memory,      // number in bytes
		MemorySwap:    memory_swap, // number in bytes for memory 
									// + swap, -1 for no swap
		CPU:           cpu,
		Command:       CommaStringToSlice(r.FormValue("command")),
		IP:            r.FormValue("ip"),
		ExposedPorts:  CommaStringToSlice(r.FormValue("exposed_ports")),
	}
	deployment, errMsg, err := DeployNewContainer(host, d)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s errmsg %s", err, errMsg)))
		return
	}
	// Testing! Works kinda.
	resp, err := json.Marshal(deployment)
	w.Write(resp)
}

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
	dispatcher.AddHandler("GET", "/docker_pool/", ReturnDockerHost)

	// Bottom is hit first, then second to last, etc
	secure := AuthorizationRequired(dispatcher.HandleRequest)
	logged := LogRequest(secure)

	http.HandleFunc("/", logged)

	log.Println("Starting")
	log.Fatal(http.ListenAndServe(HostPort, nil))
	// flag.Parse()
}
