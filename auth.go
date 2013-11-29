// auth.go handles the auth code for http basic + ip confirmation

package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

// Blame the internet for this. sends the header asking for http basic
func SendMissingCredentialsHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="luma.im"`)
	w.WriteHeader(401)
	w.Write([]byte("401 Unauthorized\n"))
}

func CheckCredentials(r *http.Request) bool {
	// Gets the auth header bits
	authHead := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(authHead) != 2 || authHead[0] != "Basic" {
		return false
	}

	// Gets the decoded stuff
	decoded, errr := base64.StdEncoding.DecodeString(authHead[1])
	if errr != nil {
		return false
	}
	// Obvious, maybe.
	userPassPair := strings.Split(string(decoded), ":")
	if len(userPassPair) != 2 {
		return false
	}
	passwd := Password[userPassPair[0]]
	if passwd == userPassPair[1] {
		return true
	}
	return false
}

// wrapper to do IP + http basic authentication ;)
func AuthorizationRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Format is ip:port. IP may be IPv6 format, e.g. ::1, which uses
		// colons, so find the right most colon
		portSeperatorIndex := strings.LastIndex(r.RemoteAddr, ":")
		ipAddress := r.RemoteAddr[0:portSeperatorIndex]

		if _, ok := AllowedIPs[ipAddress]; !ok {
			log.Println("Denied access to '" + ipAddress + "'")
			http.Error(w, ErrorMessages["addressError"], http.StatusForbidden)
			return
		}

		if CheckCredentials(r) {
			h(w, r)
			return
		}

		// Send if either credentials are invalid or none set
		SendMissingCredentialsHeader(w, r)
	}
}
