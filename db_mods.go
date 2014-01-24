// mongo stuff

package main

import (
	//"log"
	//"labix.org/v2/mgo/bson"
	"fmt"
	"labix.org/v2/mgo"
	"github.com/ziutek/mymysql/mysql"
    _ "github.com/ziutek/mymysql/native" // Native engine
)

type DockerDatabaseWrite struct {
	// future
	Username string
	UserId   string
}

func WriteToGoracleDatabase(collectionName string, d interface{}) error {

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
	if err != nil {
		return err
	}
	return nil
}

func WriteNginxConfig(n NginxConfig) error {
	
    db := mysql.New("tcp", "", (NginxDBAddress+NginxDBPort), NginxDBUser, NginxDBPassword, NginxDBName)
    err := db.Connect()
    if err != nil {
        fmt.Println(err)
    }
    defer checkError(db.Close()_
    stmt, err := db.Prepare("INSERT INTO configs (name, content, write, hash) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE (name, content, write)=VALUES(name, content, write)")
	checkError(err)
    return nil
}


