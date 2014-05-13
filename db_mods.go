// mongo stuff

package main

import (
	//"log"
	//"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type DockerDatabaseWrite struct {
	//
	// future
	// RO: no idea what this is for, if you remember comment.
	// RO: 04/19/2014: still not sure.
	// RO: 05/03/2014: nope
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
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	collection := session.DB(MongoDBName).C(collectionName)
	// that set thing is needed because Mongo.
	_, err = collection.Upsert(query, bson.M{"$set": update})
	return err
}

//
// mysql schema - this is running on the LB so we won't be
// creating it here (it's created with the python script if it doesn't exist)
//
//"CREATE TABLE `configs`
// (`id` INT(11) NOT NULL AUTO_INCREMENT,
// `name` VARCHAR(255) NOT NULL DEFAULT '0',
// `content` TEXT NOT NULL, `write` INT(11) NULL DEFAULT '0',
// `hostname` varchar(255) NOT NULL UNIQUE,
// `hash` VARCHAR(255) NOT NULL DEFAULT '0',
// PRIMARY KEY (`id`), UNIQUE INDEX `hash` (`hash`)) ENGINE=InnoDB;")

// not yet implemented
func WriteNginxConfig(n NginxConfig) error {
	db := mysql.New("tcp", "", (NginxDBAddress + ":" + NginxDBPort),
		NginxDBUser, NginxDBPassword, NginxDBName)
	err := db.Connect()
	if err != nil {
		return err
	}
	// err = CreateMysqlTableConfigs(db)
	// if err != nil {
	// return err
	// }
	// this just inserts. now, the hostname should be a unique field:
	// we won't be able to deploy twice with the same hostname
	stmt, err := db.Prepare("INSERT INTO `configs` (`name`, `content`, `write`, `hostname`, `hash`) VALUES (?, ?, ?, ?, ?)")
	//  ON DUPLICATE KEY UPDATE name=?, content=?, write=?
	if err != nil {
		return err
	}
	// fmt.Printf("name: %s, content: %s, write %s, hash %s", n.configName, n.configFile, 1, n.configHash)
	stmt.Run(n.configName, n.configFile, 1, n.configValues.hostname, n.configHash)
	//
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	return nil
}

//func GetNginxConfig(hostname string) {
//
//}
