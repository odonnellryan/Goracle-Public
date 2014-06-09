package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// 
// this is for actions on the user struct
// anything the user really wants to query for should go through here
// data passed back to the browser should be relatively safe
// no server information, etc..
//

type User struct {
    Username    string
    Containers  map[string]UserContainers
}

type UserContainers struct {
	Image         string
	Memory        int64
	MemorySwap    int64
	CPU           int64
	Command       []string
	IP            string
	ExposedPorts  []string
	DeployedInfo  DeployedContainerInfo
	ContainerInfo InspectContainerInfo
}

//
// gets all the containers associated with a user
//
// NEED TO TEST 
// ALSO NEED TO ENSURE IT DOES NOT RETURN THE ENTIRE DOCKERSERVER (bad)
//
func (u *User) GetAllUserContainers() error {
	dockerFind := []DockerServer{}
	session, err := mgo.Dial(MongoDBAddress)
	if err != nil {
		return err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(MongoDBName).C(MongoDeployCollection)
	err = c.Find(bson.M{"user": u.Username}).All(&dockerFind)
	if err != nil {
		return err
	}
	// initiate the empty map 
	u.Containers = map[string]UserContainers{}
	for index := range(dockerFind) {
	    dockerFind[index].InspectContainer()
	    name := dockerFind[index].DeployedInfo.Id
	    u.Containers[name] = UserContainers {
    	        Image: dockerFind[index].Image,
    	        Memory: dockerFind[index].Memory,
    	        MemorySwap: dockerFind[index].MemorySwap,
    	        CPU: dockerFind[index].CPU,
    	        Command: dockerFind[index].Command,
    	        IP: dockerFind[index].IP,
    	        ExposedPorts: dockerFind[index].ExposedPorts,
    	        DeployedInfo: dockerFind[index].DeployedInfo,
    	        ContainerInfo: dockerFind[index].ContainerInfo,
	    }
	}
	return nil
}