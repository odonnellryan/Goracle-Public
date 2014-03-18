package main

import (
	"testing"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	//"encoding/json"
)

var testHostFile = DockerHosts{
	Host: []Host{
		{
		Hostname: "local_testing",
		Address:  "http://127.0.0.1:8888/",
		User:     "ryan",
		Password: "test",
		Containers: 0,
		},
	},
}

var testHostTwo = DockerHosts{
	Host: []Host{
		{
		Hostname: "local_testing",
		Address:  "http://127.0.0.1:8889/",
		User:     "ryan",
		Password: "test",
		Containers: 1,
		},
		{
		Hostname: "local_testing3",
		Address:  "http://127.0.0.1:8888/",
		User:     "ryan",
		Password: "test",
		Containers: 2,
		},
	},
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

func TestUpdateAllDockerHostsInMongo(t *testing.T) {
	err := UpdateAllDockerHostsInMongo()
	if err != nil {
		t.Errorf("UpdateAllDockerHosts error: %s", err)
	}
}

// add a test to test and see if it updated
func TestIncrementContainerCount(t *testing.T) {
	err := IncrementContainerCount(testHostFile.Host[0])
	if err != nil {
		t.Errorf("IncrementContainerCount error: %s", err)
	}
	session, err := mgo.Dial(MongoDBAddress)
	if err != nil {
		t.Errorf("IncrementContainerCount error: %s", err)
	}
	c := session.DB(MongoDBName).C(MongoDockerHostCollection)
	result := Host{}
	err = c.Find(bson.M{"Hostname": testHostFile.Host[0].Hostname}).One(&result)
	if err != nil {
		t.Errorf("IncrementContainerCount error: %s", err)
	}
	if result.Containers != testHostFile.Host[0].Containers {
		t.Errorf("Increment Error Expected: %s, found: %s", 
				result.Containers, testHostFile.Host[0].Containers)
	}
}

func TestSelectHostOne(t *testing.T) {
	host, err := SelectHost()
	if err != nil {
		t.Errorf("TestSelectHost error: %s host returned: %s", err, host)
	}
	fmt.Printf("SelectOne Host: %+v \n", host)
}

func TestUpdateMultipleMongo(t *testing.T) {
	for index := range testHostTwo.Host {
		err := MongoUpsert(MongoDockerHostCollection, 
			bson.M{"Hostname": testHostTwo.Host[index].Hostname}, 
			testHostTwo.Host[index])
		if err != nil {
			t.Errorf("TestUpdateMultipleMongo error: %s host: %s", err,
					testHostTwo.Host[index])
		}
	}
}

func TestSelectHostTwo(t *testing.T) {
	dockerHosts, err := GetDockerHostInformation()
	if err != nil {
		t.Errorf("TestSelectHost error: %s host returned: %s", err, dockerHosts)
	}
	host, err := SelectHost()
	if err != nil {
		t.Errorf("TestSelectHost error: %s host returned: %s", err, host)
	}
	fmt.Printf("SelectTwo Host: %+v \n", host)
	fmt.Printf("all dockerhosts: %+v \n", dockerHosts)
}


