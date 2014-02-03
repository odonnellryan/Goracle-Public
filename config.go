package main

var (
	// host port for Goracle
	HostPort              = ":6000"
	// mongo db
	MongoDBAddress        = "127.0.0.1"
	MongoDBPort           = "27017"
	MongoDBName           = "test"
	MongoDeployCollection = "deployments"
	// nginx db for configs
	NginxDBAddress        = "127.0.0.1"
	NginxDBPort           = "3306"
	NginxDBName           = "nginx"
	NginxDBUser           = "ryan"
	NginxDBPassword       = "test"
)

// the below is for this Goracle HTTP server
var Password = map[string]string{
	"ryan": "test",
}

// only allow local connections
var AllowedIPs = map[string]bool{
	"127.0.0.1": true,
	"[::1]":     true,
}
