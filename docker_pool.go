package main

import (
	"os"
	"encoding/json"
	"fmt"
	)

type Host struct {
	Hostname	string
	Address		string
	User    	string
	Password   	string
}

type DockerHosts struct {
	Host []Host
}

func GetDockerHosts() (dockerHosts, error) {
	config, err := os.Open("dockerhosts.json")
	if err != nil {
		return nil, err
	}
	var dockerHosts DockerHosts
	err = json.Unmarshal(config, &dockerHosts)
    if err != nil {
        return nil, err
    }
    return dockerHosts, nil
}

// update a single host entry in mongo to reflect their container count
func UpdateContainerNumberInHost (host Host) {
	dockerHosts, err := GetDockerHosts()
	if err != nil {
        return nil, err
    }
	containers, err := ListAllContainers(host))
	if err != nil {
    return nil, err
	}
	// not finished
	len(containers)
}

// update all host entries in mongo to reflect their container count
func UpdateTotalContainerNumber(d DockerHosts) {
	dockerHosts, err := GetDockerHosts()
	if err != nil {
        return nil, err
    }
	for index := range(d.Host) {
		containers, err := ListAllContainers(d.Host[index]))
		if err != nil {
        return nil, err
    	}
    	// not finished
    	len(containers)
	}
}

func SelectHost() (Host, error) {
	dockerHosts, err := GetDockerHosts()
	if err != nil {
        return nil, err
    }
	for index := range(dockerHosts.Host) {
		len(ListAllContainers(dockerHosts.Host[index]))
	}
}