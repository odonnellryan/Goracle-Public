// docker.go
package main

import (
	//"strings"
	"encoding/json"
	// because this was a pain for me, here is the link to things
	// install in this order
	// mongodb: http://docs.mongodb.org/manual/tutorial/install-mongodb-on-windows/
	// bazaar (need this to install mgo, i don't know man...): http://wiki.bazaar.canonical.com/Download
	// might need the below only if you do the python 2.7 install of bazaar
	// in that case i put the cert file in c:/Python27/
	// cacert: http://curl.haxx.se/ca/cacert.pem
	// mgo: http://labix.org/mgo
	//woooooo
	"labix.org/v2/mgo"
	//"labix.org/v2/mgo/bson"
	//"log"
	"io/ioutil"
	"net/http"
)

//
// RO: added some structs and stuff to properly communicate with Docker. Dunno if this is the
// best way to do this. following: http://docs.docker.io/en/latest/api/docker_remote_api_v1.6/
//

type ContainerInfo struct {
	Id         string
	Image      string
	Command    string
	Created    string
	Status     string
	Ports      map[string]string
	SizeRw     int
	SizeRootFs int
}

type CreateContainer struct {
	Hostname     string
	User         string
	Memory       string
	MemorySwap   string
	AttachStdin  bool
	AttachStout  bool
	AttachStderr bool
	PortSpecs    string
	Privileged   bool
	Tty          bool
	OpenStdin    bool
	StdinOnce    bool
	Env          string
	Cmd          string
	Dns          string
	Image        string
	Volumes      string
	VolumesFrom  string
	WorkingDir   string
}

type ExportContainer struct {
	ContentType string
	STREAM      byte
}

type StartContainer struct {
	Binds string
	// docker example: {"lxc.utsname":"docker"}
	LxcConf map[string]string
}

type StopContainer struct {
	Id string
}

type RestartContainer struct {
	Id string
}

type KillContainer struct {
	Id string
}

type ListContainers struct {
	ContainerInfoList ContainerInfo
}

type CreateImageFromChanges struct {
	Container       string
	Repo            string
	Tag             string
	M               string
	Author          string
	ContainerParams []struct {
		Params CreateContainer
	}
}

type SearchImages struct {
	SearchTerm string
	Results    string
}

// HTTP client, http basic auth stuff

func DockerHTTPClient(d Deployment, u string) []byte {
	client := &http.Client{}
	response, errr := client.Get(DockerHost)
	if errr != nil {
		return []byte(errr.Error())
	}
	defer response.Body.Close()
	request, errr := http.NewRequest("GET", (DockerHost + u), nil)
	// do some error checking jesus fucking christ
	request.SetBasicAuth(DockerUser, DockerPass)
	response, errr = client.Do(request)
	res, errr := ioutil.ReadAll(response.Body)

	return []byte(res)
}

//
// for deployments of any *new* container
//

type Deployment struct {
	Username      string
	ContainerName string
	Image         string
	Memory        string
	Hostname      string
	Cmd           string
}

func DeployNewContainer(d Deployment, r *http.Request) []byte {
	// mongo db host, set in config.go
	session, errr := mgo.Dial(MongoDBAddress)
	if errr != nil {
		return []byte(ErrorMessages["DBConnectionError"] + string(errr.Error()))
	}
	// something something cleanup stuff lol
	defer session.Close()
	// create privatekey/nsabackdoor
	key, errr := GenerateKey(r)
	if errr != nil {
		return []byte(ErrorMessages["EncodingError"] + /*using key var remove*/ string(key))
	}
	// session for the mongodb thing
	c := session.DB(MongoDBName).C(MongoDeployCollection)
	errr = c.Insert(d)
	if errr != nil {
		return []byte(ErrorMessages["DBConnectionError"] + string(errr.Error()))
	}
	// creates slice for All() deployments to be returned (can just do One(), whatever)
	// look at docs: http://godoc.org/labix.org/v2/mgo
	result := []Deployment{}
	// i don't know why they do this like this but they do
	errr = c.Find(nil).All(&result)
	if errr != nil {
		return []byte(ErrorMessages["DBConnectionError"] + string(errr.Error()))
	}
	returnResult, errr := json.Marshal(result)
	if errr != nil {
		return []byte(ErrorMessages["EncodingError"] + string(key))
	}

	return []byte(returnResult)

}

func StopContainerRequest() {

}
