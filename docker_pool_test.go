package main

import (
	"testing"
	//"encoding/json"
)

var testHostFile = DockerHosts{
	Host: []Host{{
		Hostname: "local_testing",
		Address:  "http://127.0.0.1:8888/",
		User:     "ryan",
		Password: "test",
	}},
}

func TestGetDockerHosts(t *testing.T) {
	dockerHostsFromFile, err := GetDockerHosts()
	if err != nil {
		t.Errorf("GetDockerHosts error: %s", err)
	}
	if dockerHostsFromFile.Host[0] != testHostFile.Host[0] {
		t.Errorf("Expected %s got: %s", testHostFile,
			dockerHostsFromFile)
	}
}
