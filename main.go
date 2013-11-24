package main

import (
	//"io"
	//"encoding/json"
	"log"
	"net/http"
	"strings"
)

func HandleDeploymentRequest(w http.ResponseWriter, r *http.Request) {

	// request will be structured as such:
	// /deployments/{user_name}/{request_server}/{container_name}/{image}/{GET_Params}
	// {GET_Params} (so far...): memory=(string), hostname=(string), cmd=(string)
	// C: deployments/test/server/container/image/?memory=test&hostname=host&cmd=cmd
	Url := strings.Split(r.URL.Path, "/")

	// ideally, things in the PATH and POST should be REQUIRED and GET values should
	// be optional/return defaults. maybe.

	lUrl := len(Url)

	//checking..
	if lUrl != 7 {
		w.Write([]byte(ErrorMessages["UrlError"]))
		return
	}

	d := Deployment{
		Url[2], Url[4], Url[5], r.FormValue("memory"), r.FormValue("hostname"),
		r.FormValue("cmd"),
	}

	response := DeployNewContainer(d, r)

	//testing! works kinda

	w.Write([]byte(response))

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

	log.Fatal(http.ListenAndServe(HostPort, nil))
}
