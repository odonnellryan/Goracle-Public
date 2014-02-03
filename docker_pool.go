package main

import (
	"os"
	"encoding/json"
	// "fmt"
	"strconv"
	)

type Host struct {
	Hostname	string
	Address		string
	User    	string
	Password   	string
	Containers	string
}

type DockerHosts struct {
	Host []Host
}

func GetDockerHosts() (*DockerHosts, error) {
	config, err := os.Open("dockerhosts.json")
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(config)
	dockerHosts := &DockerHosts{}
	decoder.Decode(&dockerHosts)
    if err != nil {
        return nil, err
    }
    return dockerHosts, nil
}

// update a single host entry in mongo to reflect their container count
func UpdateContainerNumberInHost(host Host) error {
	containers, err := ListAllContainers(host)
	if err != nil {
    return err
	}
	host.Containers = strconv.Itoa(len(containers))
	UpdateContainerCount(host)
	return nil
}

// update all host entries in mongo to reflect their container count
func UpdateTotalContainerNumber(d DockerHosts) error {
	for index := range(d.Host) {
		containers, err := ListAllContainers(d.Host[index])
		if err != nil {
        	return err
    	}
    	d.Host[index].Containers = strconv.Itoa(len(containers))
		UpdateContainerCount(d.Host[index])
		return nil
	}
	return nil
}

// implement...
func SelectHost() (Host, error) {
	dockerHosts, err := GetDockerHostInformation()
	host := Host{}
	if err != nil {
        return host, err
    }
    // start super advanced algorithm
	// ro: don't try to understand this because you won't be able to.
	// don't even read it.
	
	number := dockerHosts.Host[0].Containers
	hostIndex := 0

	for index := range(dockerHosts.Host) {
		if dockerHosts.Host[index].Containers < number {
			number = dockerHosts.Host[index].Containers
			hostIndex = index
			}
	}
	host = dockerHosts.Host[hostIndex]
	return host, nil
}


