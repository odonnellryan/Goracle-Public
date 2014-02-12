// docker.go
package main

import (
	//"strings"
	"io/ioutil"
	//"log"
	"encoding/json"
	"fmt"
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
	CpuShare     string
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

func BuildDeployment(d Deployment) Deployment {
	if d.WebPort != "" {
		nginxConfig := nginxConfigValues{
			hostname:       d.Hostname,
			upstreamServer: d.IP,
			upstreamPort:   d.WebPort,
		}
		d.NginxConfig = BuildNginxConfig(nginxConfig)
	}
	d.Config = CreateContainer{
		Hostname:     d.Config.Hostname,
		User:         "",
		Memory:       d.Memory,
		MemorySwap:   "0",
		AttachStdin:  false,
		AttachStout:  true,
		AttachStderr: true,
		CpuShare:     d.CPU,
		PortSpecs:    "",
		Tty:          false,
		OpenStdin:    false,
		StdinOnce:    false,
		Env:          "",
		Param:        "",
		Cmd:          d.Command,
		Dns:          "",
		Image:        d.Image,
		Volumes:      "",
		VolumesFrom:  "",
		WorkingDir:   "",
	}
	return d
}

type StartContainer struct {
	Binds string
	// docker example: {"lxc.utsname":"docker"}
	LxcConf map[string]string
}

type ListContainers struct {
	ContainerInfoList []ContainerInfo
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
func SendDockerCommand(host Host, command string, method string) ([]byte, error) {
	client := &http.Client{}
	response, err := client.Get(host.Address)
	if err != nil {
		return nil, err
	}
	// closes the connection
	defer response.Body.Close()
	request, err := http.NewRequest(method, (host.Address + command), nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(host.User, host.Password)
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
// For all actions taken on containers through the Luma portal.
//

type Deployment struct {
	User        string
	Hostname    string
	Image       string
	Memory      string
	CPU         string
	Command     string
	IP          string
	WebPort     string
	NginxConfig NginxConfig
	Config      CreateContainer
}

func DeployNewContainer(host Host, d Deployment, r *http.Request) []byte {
	//
	// order of operations:
	// check if hostname exists
	// builds the deployment structure and configs
	// logs deployment struct to mongo, including configs
	// saves the nginx configuration
	// updates the container count for that docker host
	// deploys the container and returns connection information
	//
	exists, err := CheckContainerHostnameExists(d)
	if err != nil {
		return []byte(ErrorMessages["DBConnectionError"] + err.Error())
	}
	if exists {
		return []byte(ErrorMessages["EncodingError"] + err.Error())
	}
	// build the deployment struct
	d = BuildDeployment(d)
	// log it locally
	err = MongoInsert(MongoDeployCollection, d)
	if err != nil {
		return []byte(ErrorMessages["EncodingError"] + err.Error())
	}
	// writes to the nginx mysql database
	err = WriteNginxConfig(d.NginxConfig)
	if err != nil {
		return []byte(ErrorMessages["DBConnectionError"] + err.Error())
	}
	// update container count
	err = MongoUpsert(MongoDockerHostCollection, host.Hostname, host)
	if err != nil {
		return []byte(ErrorMessages["DBConnectionError"] + err.Error())
	}
	// make this actually do stuff.
	SendDockerCommand(host, "command", "method")

	return []byte(Messages["DeploymentSuccess"])
}

func SearchForImage(d Deployment, h Host) ([]byte, error) {
	searchString := "/images/search?term=" + d.Image
	resp, err := SendDockerCommand(h, searchString, "GET")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ListAllContainers(h Host) ([]ContainerInfo, error) {
	listString := "/containers/json?all=1"
	containers := &[]ContainerInfo{}
	resp, err := SendDockerCommand(h, listString, "GET")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, &containers)
	if err != nil {
		fmt.Println(err)
	}
	return *containers, nil
}

func StopContainerRequest() {

}
