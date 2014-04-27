// auth.go handles the auth code for http basic + ip confirmation

package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	//"fmt"
)

// sends the header asking for http basic
func SendMissingCredentialsHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	http.Error(w, "Not authorized", 401)
	return
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
	// password pair: user:password
	userPassPair := strings.Split(string(decoded), ":")
	if len(userPassPair) != 2 {
		return false
	}
	// currently, password and usernames *are hard-coded
	// into the binary and not encrypted*
	passwd := Password[userPassPair[0]]
	if (len(userPassPair[0]) == 0) || (len(userPassPair[1]) == 0) {
		return false
	}
	return CryptToHex(passwd) == CryptToHex(userPassPair[1])
}

func CheckIPAddress(r *http.Request) bool {
	// Format is ip:port. IP may be IPv6 format, e.g. ::1, which uses
	// colons, so find the right most colon
	portSeperatorIndex := strings.LastIndex(r.RemoteAddr, ":")
	ipAddress := r.RemoteAddr[0:portSeperatorIndex]
	if _, ok := AllowedIPs[ipAddress]; !ok {
		log.Printf("Denied access to: %s", r.RemoteAddr)
		return AllowedIPs[ipAddress]
	}
	return AllowedIPs[ipAddress]
}

func CheckAuthorization(r *http.Request) bool {
	return (CheckIPAddress(r)) && (CheckCredentials(r))
}

// wrapper to do IP + http basic authentication ;)
func AuthorizationRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check authorization (IP + HTTP Basic)
		if CheckAuthorization(r) {
			h(w, r)
			return
		}
		// Send if either credentials are invalid or none set
		SendMissingCredentialsHeader(w, r)
		return
	}
}
