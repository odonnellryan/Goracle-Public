// docker.go
package main

import (
	//"strings"
	"io"
	"io/ioutil"
	//"log"
	"time"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	//"net/url"
)

//
// For all actions taken on containers through the Luma portal.
//

type Deployment struct {
	DockerHost    string
	User          string
	ContainerName string
	Image         string
	Memory        int64
	MemorySwap    int64
	CPU           int64
	Command       []string
	IP            string
	ExposedPorts  []string
	Config        Config
	DeployedInfo  DeployedContainerInfo
}

type Config struct {
	Hostname        string
	Domainname      string
	User            string
	Memory          int64 // Memory limit (in bytes)
	MemorySwap      int64 // Total memory usage (memory + swap); set `-1' to disable swap
	CpuShares       int64 // CPU shares (relative weight vs. other containers)
	AttachStdin     bool
	AttachStdout    bool
	AttachStderr    bool
	PortSpecs       []string            // Deprecated - Can be in the format of 8080/tcp
	ExposedPorts    map[string]struct{} // 80/tcp
	Tty             bool                // Attach standard streams to a tty, including stdin if it is not closed.
	OpenStdin       bool                // Open stdin
	StdinOnce       bool                // If true, close stdin after the 1 attached client disconnects.
	Env             []string
	Param           []string
	Cmd             []string
	Dns             []string
	Image           string // Name of the image as it was passed by the operator (eg. could be symbolic)
	Volumes         map[string]struct{}
	VolumesFrom     string
	WorkingDir      string
	Entrypoint      []string
	NetworkDisabled bool
	OnBuild         []string
}

type State struct {
	Running    bool
	Pid        int
	ExitCode   int
	StartedAt  time.Time
    Ghost      bool
}

type HostConfig struct {
	Binds           []string
	ContainerIDFile string
	LxcConf         []KeyValuePair
	Privileged      bool
	PortBindings    map[string][]PortBinding
	Links           []string
	PublishAllPorts bool
	Dns             []string
	DnsSearch       []string
	VolumesFrom     []string
}

type InspectContainerInfo struct {
    Config      Config
    State       State
    Image       string
    Id          string
    Created     time.Time
    Warnings    []string
    HostConfig  HostConfig
    Exists      bool
}

type DeployedContainerInfo struct {
	Id       string
	Warnings []string
}

type ListContainerInfo struct {
	Id         string
	Image      string
	Command    string
	Created    int
	Status     string
}

type SearchResponse struct {
	Description string
	is_official bool
	is_trusted  bool
	Name        string
	star_count  int
}

// we aren't deploying containers with an nginx configuration to start.
// the user will have to later choose a hostname and initiate
// a custom domain deployment.
// deploying the container deploys it using docker's default hostname.
// nginx deployment will be its own thing

func BuildDeployment(d Deployment) Deployment {
	d.Config = Config {
		Memory:       d.Memory,     // Memory limit (in bytes)
		MemorySwap:   d.MemorySwap, // mem + swap, -1 to disable swap.
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		CpuShares:    d.CPU,
		Cmd:          d.Command,
		Image:        d.Image,
		ExposedPorts: make(map[string]struct{}),
		Volumes:      make(map[string]struct{}),
	}
	for index := range d.ExposedPorts {
		d.Config.ExposedPorts[d.ExposedPorts[index]] = struct{}{}
	}
	return d
}

type ListContainers struct {
	ContainerInfoList []ListContainerInfo
}

// HTTP client, http basic auth stuff
func SendDockerCommand(host Host, command string, method string, body io.Reader) (http.Response, error) {
	// url stuff yeah
	// the client woo
	nulResp := http.Response{}
	client := &http.Client{}
	response, err := client.Get(host.Address)
	if err != nil {
		return nulResp, err
	}
	// closes the connection
	defer response.Body.Close()
	request, err := http.NewRequest(method, (host.Address + command), body)
	if err != nil {
		return nulResp, err
	}
	// this is to fix some odd docker bugs.
	// we should probably fix this within docker, but...
	// really lazy. if you care, it's somewhere line 658 in api.go
	// function postContainersStart
	// there's a problem with it not realizing a nil body (who knows)
	// then it tries to json deserialize something odd
	// and it comes back broke, throws EOF, then you get an internal server error
	// so you can't start containers. 
	if body != nil {
	    request.Header.Set("Content-Type", "application/json")
	}
	request.SetBasicAuth(host.User, host.Password)
	response, err = client.Do(request)
	if err != nil {
		return nulResp, err
	}
	return *response, nil
}

func DeployNewContainer(host Host, d Deployment) (DeployedContainerInfo, string, error) {

	//
	// order of operations:
	// builds the deployment structure and configs
	// deploys the container
	// gets back container Id
	// logs deployment struct to mongo, including configs and info
	// updates the container count for that docker host
	//

	// build the deployment struct
	deployment := BuildDeployment(d)
	deployedInfo := DeployedContainerInfo{}
	body, err := json.Marshal(deployment.Config)
	if err != nil {
		return deployedInfo, "", err
	}
	// send the create command
	resp, err := SendDockerCommand(host, "containers/create", "POST",
		bytes.NewReader(body))
	if resp.StatusCode != 201 {
		//fmt.Printf("response status code: %s \n", resp.StatusCode)
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return deployedInfo, "", err
		}
		return deployedInfo, "Error with send Docker command response.", errors.New(string(msg))
	}
	if err != nil {
		return deployedInfo, "Error with send Docker command.", err
	}
	decode := json.NewDecoder(resp.Body)
	err = decode.Decode(&deployedInfo)
	if err != nil {
		return deployedInfo, "json decode", err
	}
	deployment.DeployedInfo = deployedInfo
	// log it/enter to mongo
	err = MongoInsert(MongoDeployCollection, deployment)
	if err != nil {
		return deployedInfo, "mongo error", err
	}
	// update container count
	err = IncrementContainerCount(host)
	if err != nil {
		return deployedInfo, "increment error", err
	}
	return deployedInfo, "", nil
}

func ListAllContainers(h Host) ([]ListContainerInfo, error) {
	containers := []ListContainerInfo{}
	command := "containers/json?all=1"
	resp, err := SendDockerCommand(h, command, "GET", nil)
	if err != nil {
		return containers, err
	}
	decode := json.NewDecoder(resp.Body)
	err = decode.Decode(&containers)
	if err != nil {
		return containers, err
	}
	return containers, nil
}

func InspectContainer(h Host, containerID string) (InspectContainerInfo, error) {
	containerInfo := InspectContainerInfo{}
	command := "containers/" + containerID + "/json"
	resp, err := SendDockerCommand(h, command, "GET", nil)
	if err != nil {
		return containerInfo, err
	}
	if containerID == "" || (len(containerID) < 12) {
	    resp.StatusCode = 404
	}
	if resp.StatusCode == 404 {
	    containerInfo.Exists = false
	    return containerInfo, err
	}
	if resp.StatusCode == 500 {
	    msg, err := ioutil.ReadAll(resp.Body)
	    if err != nil {
		    return containerInfo, err
	    }
	    containerInfo.Exists = false
	    return containerInfo, errors.New("Server error when inspecting image: " + string(msg))
	}
	decode := json.NewDecoder(resp.Body)
	err = decode.Decode(&containerInfo)
	if err != nil {
		containerInfo.Exists = false
		jsonErr := err
		msg, err := ioutil.ReadAll(resp.Body)
	    if err != nil {
		    return containerInfo, err
	    }
	    fmt.Printf("returned from server: %+v", msg)
	    return containerInfo, errors.New("json error: " + jsonErr.Error() + string(msg))
	}
	containerInfo.Exists = true
	return containerInfo, nil
}

func DeleteContainer(h Host, id string) (http.Response, error) {
	nulResp := http.Response{}
	command := "containers/" + id
	resp, err := SendDockerCommand(h, command, "DELETE", nil)
	if err != nil {
		return nulResp, err
	}
	return resp, nil
}

func StartContainer(h Host, id string) (http.Response, error) {
	nulResp := http.Response{}
	command := "v1.7/containers/" + id + "/start"
	resp, err := SendDockerCommand(h, command, "POST", nil)
	if err != nil {
		return nulResp, err
	}
	return resp, nil
}

func StopContainer(h Host, id string) (http.Response, error) {
	nulResp := http.Response{}
	command := "v1.7/containers/" + id + "/stop?t=5"
	resp, err := SendDockerCommand(h, command, "POST", nil)
	if err != nil {
		return nulResp, err
	}
	return resp, nil
}

// ### begin function graveyard ###
//
// not sure if we'd even want to delete images? probably just containers.
//
//func DeleteImage(h Host, id string) (http.Response, error) {
//	nulResp := http.Response{}
//	command := "containers/" + id + "/start"
//	resp, err := SendDockerCommand(h, command, "DELETE", nil)
//	if err != nil {
//		return nulResp, err
//	}
//	return resp, nil
//}
//
//  we don't need this right now. searches the repo for an image.
// 
//func SearchForImage(h Host, d Deployment) ([]SearchResponse, error) {
//  searchResp := []SearchResponse{}
//  command := "images/search?term=" + d.Image
//  resp, err := SendDockerCommand(h, command, "GET", nil)
//  if err != nil {
//	    return nil, err
//  }
//  decode := json.NewDecoder(resp.Body)
//  err = decode.Decode(&searchResp)
//  if err != nil {
//	    return searchResp, err
//  }
//  return searchResp, nil
//}

//func StopContainerRequest() { }
