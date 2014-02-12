package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	//      /deployments/{user_name}/{request_server}/{container_name}/{image}/{GET_Params}
	// {GET_Params} (so far...): memory=(string), hostname=(string), cmd=(string)
	// Ex: deployments/test/server/container/image/?memory=test&hostname=host&cmd=cmd
	Url := strings.Split(r.URL.Path, "/")

	// Ideally, things in the PATH and POST should be REQUIRED and GET values should
	// be optional/return defaults. maybe.

	// Checking...
	if len(Url) != 7 {
		w.Write([]byte(ErrorMessages["UrlError"]))
		return
	}
	host, err := SelectHost()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	d := Deployment{
		Url[2], Url[4], Url[5], r.FormValue("memory"), r.FormValue("hostname"),
		r.FormValue("cmd"), r.FormValue("ip"), host.Address, NginxConfig{}, CreateContainer{},
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
