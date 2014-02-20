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
func SendDockerCommand(host Host, command string, method string, body io.Reader) ([]byte, error) {
	client := &http.Client{}
	response, err := client.Get(host.Address)
	if err != nil {
		return nil, err
	}
	// closes the connection
	defer response.Body.Close()
	request, err := http.NewRequest(method, (host.Address + command), body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(host.User, host.Password)
	response, err = client.Do(request)
	if err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// can probably not make this use the http request?

func DeployNewContainer(host Host, d Deployment, r *http.Request) []byte {
	//
	// order of operations:
	// builds the deployment structure and configs
	// deploys the container
	// gets back container Id
	// logs deployment struct to mongo, including configs and info
	// updates the container count for that docker host
	//

	// build the deployment struct
	d = BuildDeployment(d)
	body, err := json.Marshal(d)
	if err != nil {
		return []byte(ErrorMessages["EncodingError"] + err.Error())
	}
	// send the create command
	resp, err := SendDockerCommand(host, "containers/create", "method",
		bytes.NewReader(body))
	if err != nil {
		return []byte(ErrorMessages["EncodingError"] + err.Error())
	}
	deployedInfo := DeployedContainerInfo{}
	decodeResp := json.NewDecoder(bytes.NewReader(resp))
	if err = decodeResp.Decode(&deployedInfo); err != nil {
		return []byte(ErrorMessages["EncodingError"] + err.Error())
	}
	d.DeployedInfo = deployedInfo
	// log it
	err = MongoInsert(MongoDeployCollection, d)
	if err != nil {
		return []byte(ErrorMessages["DBConnectionError"] + err.Error())
	}
	// update container count
	err = IncrementContainerCount(host)
	if err != nil {
		return []byte(ErrorMessages["DBConnectionError"] + err.Error())
	}
	return []byte(Messages["DeploymentSuccess"] + string(resp))
}

func SearchForImage(d Deployment, h Host) ([]byte, error) {
	searchString := "images/search?term=" + d.Image
	resp, err := SendDockerCommand(h, searchString, "GET", nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ListAllContainers(h Host) ([]ContainerInfo, error) {
	listString := "containers/json?all=1"
	containers := &[]ContainerInfo{}
	resp, err := SendDockerCommand(h, listString, "GET", nil)
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
