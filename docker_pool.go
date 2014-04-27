package main

import (
	"encoding/json"
	//"fmt"
	"errors"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
	//"strconv"
)

type Host struct {
	Hostname   string
	Address    string
	User       string
	Password   string
	Containers int
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
	err = decoder.Decode(&dockerHosts)
	if err != nil {
		return dockerHosts, err
	}
	return dockerHosts, nil
}

// gets the hosts then sends them to be updated in mongo db
func UpdateAllDockerHostsInMongo() error {
	dockerHosts, err := GetDockerHosts()
	if err != nil {
		return err
	}
	return UpdateTotalContainerNumber(dockerHosts)
}

// update all host entries in mongo to reflect their container count
func UpdateTotalContainerNumber(d DockerHosts) error {
	for index := range d.Host {
		containers, err := ListAllContainers(d.Host[index])
		if err != nil {
			return err
		}
		containerCount := len(containers)
		if containerCount < 1 {
			containerCount = 1
		}
		d.Host[index].Containers = containerCount
		err = MongoUpsert(MongoDockerHostCollection,
			bson.M{"hostname": d.Host[index].Hostname},
			d.Host[index])
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

// update a single host entry in mongo to reflect their container count
func UpdateContainerNumberInHost(host Host) error {
	containers, err := ListAllContainers(host)
	if err != nil {
		return err
	}
	host.Containers = len(containers)
	err = MongoUpsert(MongoDockerHostCollection,
		bson.M{"hostname": host.Hostname},
		host)
	if err != nil {
		return err
	}
	return nil
}

// increments the count of a single container config in mongo
func IncrementContainerCount(update Host) error {
	session, err := mgo.Dial(MongoDBAddress)
	session.SetMode(mgo.Monotonic, true)
	if err != nil {
		return err
	}
	defer session.Close()
	collection := session.DB(MongoDBName).C(MongoDockerHostCollection)
	// that set thing is needed because Mongo.
	_, err = collection.Upsert(bson.M{"hostname": update.Hostname},
		bson.M{"$inc": bson.M{"containers": 1}})
	return err
}

func GetDockerHostInformation() (DockerHosts, error) {
	//
	// gets all dockerhost information from the mongo DB..
	//
	dockerhosts := DockerHosts{}
	host := []Host{}
	// mongo db host, set in config.go
	session, err := mgo.Dial(MongoDBAddress)
	session.SetMode(mgo.Monotonic, true)
	if err != nil {
		return dockerhosts, err
	}
	defer session.Close()
	collection := session.DB(MongoDBName).C(MongoDockerHostCollection)
	err = collection.Find(nil).All(&host)
	if err != nil {
		return dockerhosts, err
	}
	dockerhosts.Host = host
	return dockerhosts, nil
}

func GetDockerHostByHostname(hostname string) (Host, error) {
	host := Host{}
	dockerHosts, err := GetDockerHostInformation()
	if err != nil {
		return host, err
	}
	for index := range dockerHosts.Host {
		if dockerHosts.Host[index].Hostname == hostname {
			return dockerHosts.Host[index], nil
		}
	}
	return host, errors.New("Docker host not found.")
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
	return dockerHosts.Host[hostIndex], nil
}
