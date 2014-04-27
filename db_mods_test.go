package main

import (
	"github.com/ziutek/mymysql/mysql"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"testing"
	// "fmt"
)

var testCollection = "testCollection"

type TestResult struct {
	Testone string
	Testtwo string
}

var testingQuery = TestResult{
	Testone: "1",
	Testtwo: "2",
}

func TestDBConnection(t *testing.T) {
	session, err := mgo.Dial(MongoDBAddress)
	if err != nil {
		t.Errorf("TestDBConnection: %s", err)
	}
	defer session.Close()
}

func TestMongoInsert(t *testing.T) {
	session, err := mgo.Dial(MongoDBAddress)
	if err != nil {
		t.Errorf("Dial: %s", err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(MongoDBName).C(testCollection)
	err = MongoInsert(testCollection, testingQuery)
	if err != nil {
		t.Errorf("Insert: %s", err)
	}
	result := TestResult{}
	err = c.Find(testingQuery).One(&result)
	if err != nil {
		t.Errorf("Find: %s", err)
	}
	if result != testingQuery {
		t.Errorf("Expected: %s, found: %s", testingQuery, result)
	}
}

func TestMongoUpsert(t *testing.T) {
	updateQuery := TestResult{"one", "2"}
	session, err := mgo.Dial(MongoDBAddress)
	if err != nil {
		t.Errorf("Dial: %s", err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(MongoDBName).C(testCollection)
	err = MongoUpsert(testCollection, testingQuery, updateQuery)
	if err != nil {
		t.Errorf("TestMongoUpsertError upsert: %s", err)
	}
	result := TestResult{}
	err = c.Find(bson.M{"testtwo": "2"}).One(&result)
	if err != nil {
		t.Errorf("TestMongoUpsert error: %s", err)
	}
	if result != updateQuery {
		t.Errorf("TestMongoUpsert Expected: %s, found: %s", updateQuery, result)
	}
}

func TestGetDockerHostInformation(t *testing.T) {
	dockerHosts, err := GetDockerHostInformation()
	if err != nil {
		t.Errorf("TestGetDockerHostInformation error: %s", err)
	}
	if len(dockerHosts.Host) != 0 {
		t.Errorf("TestGetDockerHostInformation returned hosts?: %s", dockerHosts)
	}
}

func TestMySQLConnection(t *testing.T) {
	db := mysql.New("tcp", "", (NginxDBAddress + ":" + NginxDBPort),
		NginxDBUser, NginxDBPassword, NginxDBName)
	err := db.Connect()
	if err != nil {
		t.Errorf("TestMySQLConnection connect error: %s", err)
	}
	err = db.Close()
	if err != nil {
		t.Errorf("TestMySQLConnection close error: %s", err)
	}
}

func TestWriteNginxConfig(t *testing.T) {
	testConfigValues := nginxConfigValues{"hostname", "localhost", "9999"}
	testConfig := BuildNginxConfig(testConfigValues)
	err := WriteNginxConfig(testConfig)
	if err != nil {
		t.Errorf("TestWriteNginxConfig write config error: %s", err)
	}
}

func TestGetNginxConfig(t *testing.T) {

}
