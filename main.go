//This is a comment

package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
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

//idk how i should include this shit >.> ????

var Password = map[string]string{
	//pass is cleartext but whooocarreeessss
	//hope github doesn't get backed
	//attack vectors = numberOfDevs + copiesOfCode
	"testuser": "hello",
}

var ErrorMessages = map[string]string{
	"addressError": "Incorrect source address.",
}

var AllowedIPs = map[string]bool{
	"127.0.0.1": true,
}

// blame the internet for this

func RequireAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="luma.im"`)
	w.WriteHeader(401)
	w.Write([]byte("401 Unauthorized\n"))
}

func CheckAuth(r *http.Request) bool {
	//gets the auth header and splits it etc.
	authHead := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(authHead) != 2 || authHead[0] != "Basic" {
		return false
	}
	//gets the decoded stuff
	decoded, errr := base64.StdEncoding.DecodeString(authHead[1])
	if errr != nil {
		return false
	}
	//obvious, maybe.
	userPassPair := strings.Split(string(decoded), ":")
	if len(userPassPair) != 2 {
		return false
	}
	//uhhhh....i think this is a good idea....
	passwd := Password[userPassPair[0]]
	if passwd == userPassPair[1] {
		return true
	}
	return false
}

//wrapper to do IP + http basic authentication ;)
func AuthorizationRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//split the string thing
		ipAddress := strings.Split(r.RemoteAddr, ":")
		if AllowedIPs[ipAddress[0]] {
			if CheckAuth(r) {
				h(w, r)
				return
			}
		} else {
			http.Error(w, ErrorMessages["addressError"], http.StatusForbidden)
		}
		RequireAuth(w, r)
	}
}

func (h *MethodServerMux) HandleFunc(action string, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux := h.muxes[action]
	if mux == nil {
		mux = http.NewServeMux()
		h.muxes[action] = mux
	}

	mux.HandleFunc(pattern, handler)
}

var DeploymentRequest = AuthorizationRequired(
	func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Fragment)
		w.Write([]byte("OMG"))
	},
)

func main() {
	mux := &MethodServerMux{make(map[string]*http.ServeMux)}

	// Add handlers here
	mux.HandleFunc("GET", "/deployments", DeploymentRequest)

	http.Handle("/", mux)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
