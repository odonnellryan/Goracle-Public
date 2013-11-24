package main

var (
	Repository = "REPO_URL"
	HostPort   = ":8080"
	MongoDB    = "127.0.0.1"
)

var Password = map[string]string{
	"testuser": "hello",
}

var AllowedIPs = map[string]bool{
	"127.0.0.1": true,
	"[::1]":     true,
}
