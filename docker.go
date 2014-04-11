// docker.go
package main

import (
	//"strings"
	"io"
	"io/ioutil"
	//"log"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	//"net/url"
)

//
// For all actions taken on containers through the Luma portal.
//

type Deployment struct {
	User          string
	ContainerName string
	Image         string
	Memory        int64
	MemorySwap    int64
	CPU           int64
	Command       []string
	IP            string
	ExposedPorts  []string
	Config        CreateContainer
	DeployedInfo  DeployedContainerInfo
}

type CreateContainer struct {
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

type DeployedContainerInfo struct {
	Id       string
	Warnings []string
}

type ContainerInfo struct {
	Id         string
	Image      string
	Command    []string
	Created    string
	Status     string
	Ports      map[string]string
	SizeRw     int
	SizeRootFs int
}

// we aren't deploying containers with an nginx configuration to start.
// the user will have to later choose a hostname and initiate a custom domain deployment.
// deploying the container deploys it using docker's default hostname.
// so nginx deployment will be its own thing

func BuildDeployment(d Deployment) Deployment {
	d.Config = CreateContainer{
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
	ContainerInfoList []ContainerInfo
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
    request.Header.Set("Content-Type", "application/json")
    request.SetBasicAuth(host.User, host.Password)
    response, err = client.Do(request)
    if err != nil {
            return nulResp, err
    }
	return *response, nil
}

// can probably not make this use the http request?

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
		fmt.Printf("response status code: %s \n", resp.StatusCode)
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return deployedInfo, "", err
		}
		fmt.Printf("response error: %s \n", msg)
	}
	if err != nil {
		return deployedInfo, "", err
	}
	decode := json.NewDecoder(resp.Body)
	err = decode.Decode(&deployedInfo)
	if err != nil {
		return deployedInfo, "json decode", err
	}
	deployment.DeployedInfo = deployedInfo
	// log it
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

func SearchForImage(d Deployment, h Host) ([]byte, error) {
	searchString := "images/search?term=" + d.Image
	resp, err := SendDockerCommand(h, searchString, "GET", nil)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func ListAllContainers(h Host) ([]ContainerInfo, error) {
	listString := "containers/json?all=1"
	containers := []ContainerInfo{}
	resp, err := SendDockerCommand(h, listString, "GET", nil)
	if err != nil {
		return nil, err
	}
	decode := json.NewDecoder(resp.Body)
	err = decode.Decode(&containers)
	if err != nil {
		return containers, err
	}
	if err != nil {
		return nil, err
	}
	return containers, nil
}

//func StopContainerRequest() { }
