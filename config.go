package main

var (
	repository = "REPO_URL"
)

var Password = map[string]string{
	"testuser": "hello",
}

var AllowedIPs = map[string]bool{
	"127.0.0.1": true,
	"[::1]":     true,
}
