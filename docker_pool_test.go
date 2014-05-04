package main

import (
	"testing"
	//"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	//"encoding/json"
)

var testHostFile = DockerHosts{
	Host: []Host{
		{
			Hostname:   "local_testing",
			Address:    "http://127.0.0.1:8888/",
			User:       "ryan",
			Password:   "test",
			Containers: 1,
		},
	},
}

var testHostTwo = DockerHosts{
	Host: []Host{
		{
			Hostname:   "local_testing2",
			Address:    "http://127.0.0.1:8889/",
			User:       "ryan",
			Password:   "test",
			Containers: 2,
		},
		{
			Hostname:   "local_testingthree",
			Address:    "http://127.0.0.1:8888/",
			User:       "ryan",
			Password:   "test",
			Containers: 0,
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

func TestIncrementContainerCount(t *testing.T) {
    var firstResult []Host
	var result []Host
	session, err := mgo.Dial(MongoDBAddress)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	if err != nil {
		t.Errorf("test IncrementContainerCount dial error: %s", err)
	}
	c := session.DB(MongoDBName).C(MongoDockerHostCollection)
	err = c.Find(bson.M{"hostname": testHostFile.Host[0].Hostname}).All(&firstResult)
	if err != nil {
		t.Errorf("IncrementContainerCount mongo find error: %s", err)
	}
	err = IncrementContainerCount(testHostFile.Host[0])
	if err != nil {
		t.Errorf("IncrementContainerCount error: %s", err)
	}
	err = c.Find(bson.M{"hostname": testHostFile.Host[0].Hostname}).All(&result)
	if err != nil {
		t.Errorf("IncrementContainerCount mongo find error: %s", err)
	}
	newContainerCount := firstResult[0].Containers + 1
	if result[0].Containers != newContainerCount {
		t.Errorf("IncrementContainerCount error. Expecting: %+v, found: %s", firstResult, result[0].Containers)
	}
}

func TestSelectHostOne(t *testing.T) {
	host, err := SelectHost()
	if err != nil {
		t.Errorf("TestSelectHost error: %s host returned: %s", err, host)
	}
	//fmt.Printf("SelectOne Host: %+v \n", host)
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
	host, err := SelectHost()
	if err != nil {
		t.Errorf("TestSelectHost error: %s host returned: %s", err, host)
	}
	if host != testHostTwo.Host[1] {
		t.Errorf("SelectHost error: expecting %+v got %+v",
			testHostTwo.Host[1], host)
	}
}

func TestUpdateContainerNumberInHost(t *testing.T) {
	var result []Host
	session, err := mgo.Dial(MongoDBAddress)
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	err = UpdateContainerNumberInHost(testHostFile.Host[0])
	if err != nil {
		t.Errorf("TestUpdateContainerNumberInHost error: %s", err)
	}
	c := session.DB(MongoDBName).C(MongoDockerHostCollection)
	err = c.Find(bson.M{"hostname": testHostFile.Host[0].Hostname}).All(&result)
	if err != nil {
		t.Errorf("TestUpdateContainerNumberInHost mongo find error: %s", err)
	}
	//fmt.Printf("result: %+v", result)
	if result[0].Containers < 0 {
		t.Errorf("TestUpdateContainerNumberInHost error. found: %s", result[0].Containers)
	}
}

func TestGetDockerHostByHostname(t *testing.T) {
	host, err := GetDockerHostByHostname("local_testing")
	if err != nil {
		t.Errorf("TestGetDockerHostByHostname error: %s", err)
	}
	if host.Hostname != "local_testing" {
		t.Errorf("TestGetDockerHostByHostname returned wrong host: %s", host)
	}
}
