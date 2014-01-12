package main

var (
	Repository            = "REPO_URL"
	DockerHost            = "http://127.0.0.1:8080"
	DockerUser            = "ryan"
	DockerPass            = "test"
	HostPort              = ":6000"
	MongoDBAddress        = "127.0.0.1"
	MongoDBName           = "test"
	MongoDeployCollection = "deployments"
)

//the below is for this Goracle HTTP server
var Password = map[string]string{
	"testuser": "hello",
}

var AllowedIPs = map[string]bool{
	"127.0.0.1": true,
	"[::1]":     true,
}
