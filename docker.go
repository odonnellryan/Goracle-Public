// docker.go
package main

import (
	//"fmt"
	//"strings"
	//"encoding/json"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
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

func DeployNewContainer(d Deployment, w http.ResponseWriter, r *http.Request) string {
	// create privatekey
	key, kerror := rsa.GenerateKey(rand.Reader, 1024)
	if kerror != nil {
		log.Println(ErrorMessages["EncodingError"])
		response := ErrorMessages["EncodingError"]
		return response
	}
	// save to pem
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	return string(pemdata)
}

func StopContainerRequest() {

}
