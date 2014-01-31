package main

import (
	"os"
	"encoding/json"
	"fmt"
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
func SelectHost() (*Host, error) {
	dockerHosts, err := GetDockerHosts()
	if err != nil {
        return nil, err
    }
	for index := range(dockerHosts.Host) {
		containers, err := ListAllContainers(dockerHosts.Host[index])
		if err != nil {
        	return nil, err
    	}
    	fmt.Sprintf("%s", containers)
	}
	return nil,nil
}