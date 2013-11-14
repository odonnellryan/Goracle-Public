package main

import (
	"log"
	"net/http"
)

type MethodServerMux struct {
	muxes map[string]*http.ServeMux
}

func (h *MethodServerMux) HandleRequest(w http.ResponseWriter, r *http.Request) {
	mux := h.muxes[r.Method]
	if mux == nil {
		http.NotFound(w, r)
		return
	}

	mux.ServeHTTP(w, r)
}

var ErrorMessages = map[string]string{
	"addressError": "Access denied",
}

func (h *MethodServerMux) AddHandler(action string, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux := h.muxes[action]
	if mux == nil {
		mux = http.NewServeMux()
		h.muxes[action] = mux
	}

	mux.HandleFunc(pattern, handler)
}

func LogRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request for '" + r.URL.String() + "' from '" + r.RemoteAddr + "'")
		h(w, r)
	}
}
