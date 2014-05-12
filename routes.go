package main

import (
	"encoding/json"
	//"flag"
	"fmt"
	//"log"
	"net/http"
	"strconv"
)

func DockerHost(w http.ResponseWriter, r *http.Request) {
	host, err := SelectHost()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	respBytes, err := json.Marshal(host)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	w.Write(respBytes)
}

func ContainerInfo(w http.ResponseWriter, r *http.Request) {
	host, err := GetDockerHostByHostname(r.FormValue("docker_hostname"))
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	respBytes, err := json.Marshal(host)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	w.Write(respBytes)
}

func HandleDeploymentRequest(w http.ResponseWriter, r *http.Request) {
	//
	// post request:
	// FormValue: memory, memory_swap, user, container_name, image,
	// 				command, exposed_ports, ip
	//
	// convert to int64
	memory, err := strconv.ParseInt(r.FormValue("memory"), 10, 64)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	memory_swap, err := strconv.ParseInt(r.FormValue("memory_swap"),
		10, 64)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	cpu, err := strconv.ParseInt(r.FormValue("cpu"), 10, 64)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	deploymentConfig := DockerServer {
		User:          r.FormValue("user"),
		ContainerName: r.FormValue("container_name"),
		Image:         r.FormValue("image"),
		Memory:        memory,      // number in bytes
		MemorySwap:    memory_swap, // number in bytes for memory + swap, -1 for no swap
		CPU:           cpu,
		Command:       CommaStringToSlice(r.FormValue("command")),
		IP:            r.FormValue("ip"),
		ExposedPorts:  CommaStringToSlice(r.FormValue("exposed_ports")),
	}
	dockerServer, err := NewDockerServer(deploymentConfig)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	errMsg, err := dockerServer.DeployNewContainer()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s errmsg %s", err, errMsg)))
		return
	}
	// Testing! Works kinda.
	resp, err := json.Marshal(dockerServer)
	w.Write(resp)
}