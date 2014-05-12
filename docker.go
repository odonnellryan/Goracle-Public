// docker.go
package main

import (
	//"strings"
	"io"
	"io/ioutil"
	//"log"
	"bytes"
	"encoding/json"
	"errors"
	"time"
	//"fmt"
	"net/http"
	//"net/url"
)

//
// For all actions taken on containers through the Luma portal.
//

type DockerServer struct {
	Hostname      string
	Host          Host
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
	ContainerInfo InspectContainerInfo
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
	Running   bool
	Pid       int
	ExitCode  int
	StartedAt time.Time
	Ghost     bool
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
	Config     Config
	State      State
	Image      string
	Id         string
	Created    time.Time
	Warnings   []string
	HostConfig HostConfig
	Exists     bool
}

type DeployedContainerInfo struct {
	Id       string
	Warnings []string
}

type ListContainerInfo struct {
	Id      string
	Image   string
	Command string
	Created int
	Status  string
}

type ListContainers struct {
	ContainerInfoList []ListContainerInfo
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

func (d *DockerServer) BuildDeployment() {
	d.Config = Config{
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
}

func NewDockerServer(d DockerServer) (DockerServer, error) {
	//
	// created a new docker server connection.
	// basically just finds the host details.
	//
	if d.Hostname != "" {
		host, err := GetDockerHostByHostname(d.Hostname)
		if err != nil {
			host, err := SelectHost()
			if err != nil {
				return d, err
			}
			d.Host = host
			d.Hostname = host.Hostname
			return d, nil
		}
		d.Host = host
		d.Hostname = host.Hostname
	} else {
		host, err := SelectHost()
		if err != nil {
			return d, err
		}
		d.Host = host
		d.Hostname = host.Hostname
	}
	return d, nil
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
	// then it tries to take a json action on something odd
	// and it comes back broke (it's not json), throws EOF,
	// then you get an internal server error so you can't start the container.
	// basically, if the body is nil it doesn't matter what
	// the content is, right? so who cares.
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

func (d *DockerServer) DeleteContainer() (http.Response, error) {
	nulResp := http.Response{}
	command := "containers/" + d.DeployedInfo.Id
	resp, err := SendDockerCommand(d.Host, command, "DELETE", nil)
	if err != nil {
		return nulResp, err
	}
	return resp, nil
}

func (d *DockerServer) StartContainer() (http.Response, error) {
	nulResp := http.Response{}
	command := "v1.7/containers/" + d.DeployedInfo.Id + "/start"
	resp, err := SendDockerCommand(d.Host, command, "POST", nil)
	if err != nil {
		return nulResp, err
	}
	return resp, nil
}

func (d *DockerServer) StopContainer() (http.Response, error) {
	nulResp := http.Response{}
	command := "v1.7/containers/" + d.DeployedInfo.Id + "/stop?t=5"
	resp, err := SendDockerCommand(d.Host, command, "POST", nil)
	if err != nil {
		return nulResp, err
	}
	return resp, nil
}

// these functions will probably be refactored and
// thrown into a "docker_actions" file, or something
// basically, these are actions/processes/logic, the above are commands

func (d *DockerServer) DeployNewContainer() (string, error) {

	//
	// order of operations:
	// builds the deployment structure and configs
	// deploys the container
	// gets back container Id
	// logs deployment struct to mongo, including configs and info
	// updates the container count for that docker host
	//

	// build the deployment struct
	d.BuildDeployment()
	deployedInfo := DeployedContainerInfo{}
	body, err := json.Marshal(d.Config)
	if err != nil {
		return "", err
	}
	// send the create command
	resp, err := SendDockerCommand(d.Host, "containers/create", "POST",
		bytes.NewReader(body))
	if err != nil {
		return "Error with send Docker command.", err
	}
	if resp.StatusCode != 201 {
		//fmt.Printf("response status code: %s \n", resp.StatusCode)
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "Error with send Docker command response.", errors.New(string(msg))
	}
	decode := json.NewDecoder(resp.Body)
	err = decode.Decode(&deployedInfo)
	if err != nil {
		return "json decode", err
	}
	d.DeployedInfo = deployedInfo
	d.ContainerName = deployedInfo.Id
	// log it/enter to mongo
	err = MongoInsert(MongoDeployCollection, d)
	if err != nil {
		return "mongo error", err
	}
	// update container count
	err = IncrementContainerCount(d.Host)
	if err != nil {
		return "increment error", err
	}
	return "", nil
}

// for host, not docker really..
func ListAllContainers(host Host) ([]ListContainerInfo, error) {
	containers := []ListContainerInfo{}
	command := "containers/json?all=1"
	resp, err := SendDockerCommand(host, command, "GET", nil)
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

// refactor this. should return an http response
// then there should be the "docker processor" that
// will find the host by the Deployment struct
// and do the processing (find the host once)
// otherwise, we'd have to find it before somewhere
// it'd be nice to do:
// dockerServer := DockerServer(&Deployment)
// containerInfo := dockerAction.inspect()
func (d *DockerServer) InspectContainer() error {
	containerInfo := InspectContainerInfo{}
	command := "containers/" + d.DeployedInfo.Id + "/json"
	resp, err := SendDockerCommand(d.Host, command, "GET", nil)
	if err != nil {
		return err
	}
	if d.DeployedInfo.Id == "" {
		err = errors.New("Containername error: container name does not exist")
		resp.StatusCode = 404
	}
	if resp.StatusCode == 404 {
		d.ContainerInfo.Exists = false
		return err
	}
	if resp.StatusCode == 500 {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		d.ContainerInfo.Exists = false
		return errors.New("Server error when inspecting image: " + string(msg))
	}
	decode := json.NewDecoder(resp.Body)
	err = decode.Decode(&containerInfo)
	if err != nil {
		d.ContainerInfo.Exists = false
		jsonErr := err
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New("json error: " + jsonErr.Error() + string(msg))
	}
	d.ContainerInfo = containerInfo
	d.ContainerInfo.Exists = true
	return nil
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
