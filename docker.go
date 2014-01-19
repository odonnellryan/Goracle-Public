// docker.go
package main

import (
	//"strings"
	"io/ioutil"
	//"log"
	"net/http"
)

//
// RO: added some structs and stuff to properly communicate with Docker. Dunno if this is the
// best way to do this. following: http://docs.docker.io/en/latest/api/docker_remote_api_v1.6/
// RO: This is fine.
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
	Param        string
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
func SendDockerCommand(d Deployment, u string) ([]byte, error) {
	client := &http.Client{}
	response, err := client.Get(DockerHost)
	if err != nil {
		return nil, err
	}
	// closes the connection
	defer response.Body.Close()
	request, err := http.NewRequest("GET", (DockerHost + u), nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(DockerUser, DockerPass)
	response, err = client.Do(request)
	if err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return []byte(res), nil
}

//
// For deployments of any *new* container
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

	// create privatekey/nsabackdoor
	returnResult := WriteToGoracleDatabase("deployments", d)
	if returnResult != nil {
		return []byte(ErrorMessages["EncodingError"] + returnResult.Error())
	}
	return []byte(Messages["DeploymentSuccess"])

}

func StopContainerRequest() {

}
