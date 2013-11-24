// docker.go
package main

import (
	//"strings"
	"encoding/json"
	"labix.org/v2/mgo"
	//"labix.org/v2/mgo/bson"
	//"log"
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

//
// for deployments of any *new* container
//

func DeployNewContainer(d Deployment, w http.ResponseWriter, r *http.Request) string {
	// mongo db host, set in config.go
	session, errr := mgo.Dial(MongoDB)
	if errr != nil {
		return (ErrorMessages["DBConnectionError"] + string(errr.Error()))
	}
	// create privatekey nsabackdoor
	key, errr := GenerateKey(w, r)
	if errr != nil {
		return (ErrorMessages["EncodingError"] + string(key))
	}
	djson, errr := json.Marshal(d)
	if errr != nil {
		return (ErrorMessages["EncodingError"] + string(key))
	}
	c := session.DB("test").C("deployments")
	errr = c.Insert(djson)
	if errr != nil {
		return (ErrorMessages["DBConnectionError"] + string(errr.Error()))
	}
	result := []Deployment{}

	errr = c.Find(nil).All(&result)

	if errr != nil {
		return (ErrorMessages["DBConnectionError"] + string(errr.Error()))
	}

	returnResult, errr := json.Marshal(result)

	if errr != nil {
		return (ErrorMessages["EncodingError"] + string(key))
	}

	return string(returnResult)

}

func StopContainerRequest() {

}
