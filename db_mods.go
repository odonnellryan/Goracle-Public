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
	"labix.org/v2/mgo"
)

type DockerDatabaseWrite struct {
	//future
	Username string
	UserId   string
}

func WriteToDatabase(c string, d interface{}) error {

	//
	// writes to the database defined in config and the collection defined above (c)
	// uses struct (d interface) to define the data structure
	//

	// mongo db host, set in config.go
	session, errr := mgo.Dial(MongoDBAddress)
	if errr != nil {
		return errr
	}
	// something something cleanup stuff
	defer session.Close()
	container := session.DB(MongoDBName).C(c)
	errr = container.Insert(d)
	if errr != nil {
		return errr
	}
	return nil
}
