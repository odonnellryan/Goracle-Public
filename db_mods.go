// mongo stuff

package main

import (
	//"log"
	// because this was a pain for me, here is the link to things
	// install in this order
	// mongodb: http://docs.mongodb.org/manual/tutorial/install-mongodb-on-windows/
	// bazaar (need this to install mgo, i don't know man...): http://wiki.bazaar.canonical.com/Download
	// might need the below only if you do the python 2.7 install of bazaar
	// in that case i put the cert file in c:/Python27/
	// cacert: http://curl.haxx.se/ca/cacert.pem
	// mgo: http://labix.org/mgo
	//woooooo
	//"labix.org/v2/mgo/bson"
	"fmt"
	"labix.org/v2/mgo"
	_ "github.com/ziutek/mymysql/mysql"
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
    
}


