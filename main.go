//This is a comment

package main

import (
	"encoding/json"
	//"html"
	"log"
	"net/http"
)

type Deployment struct {
	ContainerName    string
	ContainerPackage string
	ContainerAuthKey string
	IpAddress        string
}

type MethodServerMux struct {
	muxes map[string]*http.ServeMux
}

func (h *MethodServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request '" + r.URL.String() + "'")

	mux := h.muxes[r.Method]
	if mux == nil {
		http.NotFound(w, r)
		return
	}

	mux.ServeHTTP(w, r)
}

func (h *MethodServerMux) HandleFunc(action string, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux := h.muxes[action]
	if mux == nil {
		mux = http.NewServeMux()
		h.muxes[action] = mux
	}

	mux.HandleFunc(pattern, handler)
}

func HandleDeploymentGet(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Fragment)
}

func main() {
	mux := &MethodServerMux{make(map[string]*http.ServeMux)}

	// Add handlers here
	mux.HandleFunc("GET", "/deployments", HandleDeploymentGet)

	http.Handle("/", func(w http.ResponseWriter, r *http.Request) {
		r.
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
