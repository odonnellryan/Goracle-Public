package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

var testHost = Host{
	Hostname: "local_testing",
	Address:  "http://127.0.0.1:8888/",
	User:     "ryan",
	Password: "test",
}

var testDeployment = Deployment{
	User:          "testUser",
	ContainerName: "containerName",
	Image:         "docker-test-image",
	Command:       []string{"sh"},
	IP:            "127.0.0.1",
	ExposedPorts:  []string{"88/tcp", "22/tcp"},
}

func TestBuildDeployment(t *testing.T) {

	testBuild := BuildDeployment(testDeployment)

	if testBuild.Config.Memory != testDeployment.Memory {
		t.Errorf("expected %s, got %s", testDeployment.Memory,
			testBuild.Config.Memory)
	}
	if testBuild.Config.MemorySwap != testDeployment.MemorySwap {
		t.Errorf("expected %s, got %s", testDeployment.MemorySwap,
			testBuild.Config.MemorySwap)
	}
	if testBuild.Config.CpuShares != testDeployment.CPU {
		t.Errorf("expected %s, got %s", testDeployment.CPU,
			testBuild.Config.CpuShares)
	}
	if testBuild.Config.Image != testDeployment.Image {
		t.Errorf("expected %s, got %s", testDeployment.Image,
			testBuild.Config.Image)
	}
	if testBuild.Config.Memory != testDeployment.Memory {
		t.Errorf("expected %s, got %s", testDeployment.Memory,
			testBuild.Config.Memory)
	}
	if testBuild.Config.Memory != testDeployment.Memory {
		t.Errorf("expected %s, got %s", testDeployment.Memory,
			testBuild.Config.Memory)
	}
	for index := range testDeployment.ExposedPorts {
		if _, ok := testBuild.Config.ExposedPorts[testDeployment.ExposedPorts[index]]; !ok {
			t.Errorf("ExposedPorts do not match. Deployment: %s, Build: %s",
				testDeployment.ExposedPorts, testBuild.Config.ExposedPorts)
		}
	}
}

func TestHTTPToDocker(t *testing.T) {
	host := testHost
	client := &http.Client{}
	resp, err := client.Get(host.Address)
	if err != nil {
		t.Errorf("Client error: %s", err)
	}
	// closes the connection
	defer resp.Body.Close()
	request, err := http.NewRequest("GET", (host.Address + "images/json"), nil)
	if err != nil {
		t.Errorf("Request build error: %s", err)
	}
	request.SetBasicAuth(host.User, host.Password)
	response, err := client.Do(request)
	if err != nil {
		t.Errorf("Response error: %s", err)
	}
	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Read response error: %s", err)
	}
}

func TestDockerPull(t *testing.T) {
	resp, err := SendDockerCommand(testHost,
		"images/create?fromImage=docker-test-image", "POST", nil)
	if err != nil {
		t.Errorf("Error: %s \n", err)
	}
	msg, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected Status Code: %s \n", resp.StatusCode)
		if err != nil {
			t.Errorf("Error: %s \n", err)
		}
		t.Errorf("Reason: %s \n", msg)
	}
}

func TestListContainers(t *testing.T) {
	cont, err := ListAllContainers(testHost)
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if len(cont) != 0 {
		t.Errorf("Length: %i, Containers are: %s", len(cont), cont)
	}
}

func TestSendDockerCommand(t *testing.T) {
	resp, err := SendDockerCommand(testHost, "images/json", "GET", nil)
	if err != nil {
		t.Errorf("Error: %s \n", err)
	}
	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error: %s \n", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Unexpected Status Code: %s \n", resp.StatusCode)
		t.Errorf("Reason: %s \n", msg)
	}
	fmt.Printf("Images: %s \n", msg)
}

func TestDeployNewContainer(t *testing.T) {
	containerInfo, errFrom, err := DeployNewContainer(testHost, testDeployment)
	if err != nil {
		t.Errorf("TestDeployNewContainer error: %s thrown by: \n", err, errFrom)

	}
	if len(containerInfo.Warnings) != 0 {
		t.Errorf("TestDeployNewContainer warning thrown: %+v \n", containerInfo.Warnings)
	}
}

func TestSearchForImage(t *testing.T) {

}
