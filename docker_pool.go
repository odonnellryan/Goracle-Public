package main

import (
	"encoding/json"
	"os"
	// "fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strconv"
)

type Host struct {
	Hostname   string
	Address    string
	User       string
	Password   string
	Containers string
}

type DockerHosts struct {
	Host []Host
}

// gets the docker configs from the json config file.
func GetDockerHosts() (DockerHosts, error) {
	dockerHosts := DockerHosts{}
	config, err := os.Open("dockerhosts.json")
	if err != nil {
		return dockerHosts, err
	}
	decoder := json.NewDecoder(config)
	decoder.Decode(&dockerHosts)
	if err != nil {
		return dockerHosts, err
	}
	return dockerHosts, nil
}

// updates all hosts in the mongo db
func UpdateAllMongoDockerHostsInCollection() error {
	dockerHosts, err := GetDockerHosts()
	if err != nil {
		return err
	}
	return UpdateTotalContainerNumber(dockerHosts)
}

// update a single host entry in mongo to reflect their container count
func UpdateContainerNumberInHost(host Host) error {
	containers, err := ListAllContainers(host)
	if err != nil {
		return err
	}
	host.Containers = strconv.Itoa(len(containers))
	err = MongoUpsert(MongoDockerHostCollection, host.Hostname, host)
	if err != nil {
		return err
	}
	return nil
}

// update all host entries in mongo to reflect their container count
func UpdateTotalContainerNumber(d DockerHosts) error {
	for index := range d.Host {
		containers, err := ListAllContainers(d.Host[index])
		if err != nil {
			return err
		}
		d.Host[index].Containers = strconv.Itoa(len(containers))
		err = MongoUpsert(MongoDockerHostCollection, d.Host[index].Hostname, d.Host[index])
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

// increments the count of a single container config in mongo
func IncrementContainerCount(update Host) error {
	session, err := mgo.Dial(MongoDBAddress)
	if err != nil {
		return err
	}
	defer session.Close()
	collection := session.DB(MongoDBName).C(MongoDockerHostCollection)
	return collection.Update(update.Hostname, bson.M{"$inc": bson.M{"Containers": 1}})
}

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
	for index := range dockerHosts.Host {
		if dockerHosts.Host[index].Containers < number {
			number = dockerHosts.Host[index].Containers
			hostIndex = index
		}
	}
	host = dockerHosts.Host[hostIndex]
	return host, nil
}
