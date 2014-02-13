// mongo stuff

package main

import (
	//"log"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	"labix.org/v2/mgo"
)

type DockerDatabaseWrite struct {
	// future
	Username string
	UserId   string
}

func MongoInsert(collectionName string, d interface{}) error {

	//
	// writes to the database defined in config and the collection defined above (c)
	// uses struct (d interface) to define the data structure
	//

	// mongo db host, set in config.go
	session, err := mgo.Dial(MongoDBAddress)
	if err != nil {
		return err
	}
	// something something cleanup stuff
	defer session.Close()
	collection := session.DB(MongoDBName).C(collectionName)
	err = collection.Insert(d)
	return err
}

func MongoUpsert(collectionName string, query interface{}, update interface{}) error {
	//
	// is used to update the dockerhosts collection. this is used to keep
	// track of information about various docker hosts.
	//
	// mongo db host, set in config.go
	session, err := mgo.Dial(MongoDBAddress)
	if err != nil {
		return err
	}
	defer session.Close()
	collection := session.DB(MongoDBName).C(collectionName)
	_, err = collection.Upsert(query, update)
	return err
}

func GetDockerHostInformation() (DockerHosts, error) {
	//
	// gets all dockerhost information from the mongo DB..
	//
	hosts := DockerHosts{}
	// mongo db host, set in config.go
	session, err := mgo.Dial(MongoDBAddress)
	if err != nil {
		return hosts, err
	}
	defer session.Close()
	collection := session.DB(MongoDBName).C(MongoDockerHostCollection)
	err = collection.Find(nil).All(&hosts)
	if err != nil {
		return hosts, err
	}
	return hosts, nil
}

func WriteNginxConfig(n NginxConfig) error {
	db := mysql.New("tcp", "", (NginxDBAddress + NginxDBPort), NginxDBUser, NginxDBPassword, NginxDBName)
	err := db.Connect()
	if err != nil {
		fmt.Println(err)
	}
	err = db.Close()
	if err != nil {
		fmt.Println(err)
	}
	stmt, err := db.Prepare("INSERT INTO configs (name, content, write, hash) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE (name, content, write)=VALUES(name, content, write)")
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Run()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func CheckContainerHostnameExists(d Deployment) (bool, error) {
	session, err := mgo.Dial(MongoDBAddress)
	if err != nil {
		return false, err
	}
	defer session.Close()
	collection := session.DB(MongoDBName).C(MongoDeployCollection)
	containerName := d.ContainerName
	exists := collection.Find(d.ContainerName).One(&containerName)
	if exists == nil {
		return false, nil
	}
	if containerName != "" {
		return true, nil
	}
	return true, nil
}
