// docker.go
package main

import (
	//"fmt"
)

// RO: added some structs and stuff to properly
// communicate with Docker. Dunno if this is the
// best way to do this. following: 
// http://docs.docker.io/en/latest/api/docker_remote_api_v1.6/

type ContainerInfo struct {
	Id 				string
    Image 			string
    Command 		string
    Created 		int
    Status 			string
	Ports 			map[string]string
	SizeRw 			int
	SizeRootFs 		int
}

type CreateContainer struct {
	Hostname 		string
	User			string
	Memory			int
	MemorySwap 		int
	AttachStdin		bool
	AttachStout		bool
	AttachStderr	bool
	PortSpecs		string
	Privileged		bool
	Tty				bool
	OpenStdin		bool
	StdinOnce		bool
	Env				string
	Cmd				string
	Dns				string
	Image			string
	Volumes			string
	VolumesFrom		string
	WorkingDir		string
}

type ExportContainer struct {
	Content-Type	string
	STREAM			byte
}

type StartContainer struct {
	Binds			string
	// docker example: {"lxc.utsname":"docker"}
	LxcConf			map[string]string
}

type StopContainer struct {
	Id				string
}

type RestartContainer struct {
	Id				string
}

type KillContainer struct {
	Id				string
}

type ListContainers struct {
	ContainerInfoList ContainerInfo
}

type CreateImageFromChanges {
	Container		string
	Repo			string
	Tag				string
	M 				string
	Author			string
	ContainerParams CreateContainer
}

type SearchImages {
	SearchTerm		string
	Results			string
}
