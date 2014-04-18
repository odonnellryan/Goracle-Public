// mongo stuff

package main

import (
	//"log"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type DockerDatabaseWrite struct {
	//
	// future
	// RO: no idea what this is for, if you remember comment.
	//
	Username string
	UserId   string
}

func MongoInsert(collectionName string, d interface{}) error {
	//
	// writes to the database defined in config and the collection
	// defined above (collectionName)
	// uses struct (d interface) to define the data structure
	//
	// mongo db host, set in config.go
	session, err := mgo.Dial(MongoDBAddress)
	if err != nil {
		return err
	}
	defer session.Close()
	collection := session.DB(MongoDBName).C(collectionName)
	err = collection.Insert(d)
	return err
}

func MongoUpsert(collectionName string, query interface{},
	update interface{}) error {
	//
	// is used to update the dockerhosts collection. this is used to keep
	// track of information about various docker hosts.
	//
	// mongo db host, set in config.go
	session, err := mgo.Dial(MongoDBAddress)
	session.SetMode(mgo.Monotonic, true)
	if err != nil {
		return err
	}
	defer session.Close()
	collection := session.DB(MongoDBName).C(collectionName)
	// that set thing is needed because Mongo.
	_, err = collection.Upsert(query, bson.M{"$set":update})
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

// not yet implemented
func WriteNginxConfig(n NginxConfig) error {
	db := mysql.New("tcp", "", (NginxDBAddress + NginxDBPort),
		NginxDBUser, NginxDBPassword, NginxDBName)
	err := db.Connect()
	if err != nil {
		fmt.Println(err)
	}
	err = db.Close()
	if err != nil {
		fmt.Println(err)
	}
	stmt, err := db.Prepare(`INSERT INTO configs 
							(name, content, write, hash) 
							VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE 
							(name, content, write) = VALUES 
							(name, content, write)`)
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Run()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
