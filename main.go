//This is a comment

package main

import (
	"encoding/json"
	"fmt"
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

func main() {
	http.HandleFunc("/deployment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			ip = strings.Split(r.RemoteAddr, ":")[0]
			containerValues := Deployment{r.FormValue("name"), r.FormValue("package"), r.FormValue("auth")}
			jsonIze, error := json.Marshal(containerValues)
			if error != nil {
				fmt.Fprintf(w, "Error:, %q", error)
			}
			w.Write([]byte(jsonIze))
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
