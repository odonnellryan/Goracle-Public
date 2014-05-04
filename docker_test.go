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
	Command:       []string{"cat"},
	IP:            "127.0.0.1",
	ExposedPorts:  []string{"88/tcp", "22/tcp"},
}


// TODO: review...
func TestBuildDeployment(t *testing.T) {
	testBuild := BuildDeployment(testDeployment)
	if testBuild.Config.MemorySwap != testDeployment.MemorySwap {
		t.Errorf("MemorySwap expected %s, got %s", testDeployment.MemorySwap,
			testBuild.Config.MemorySwap)
	}
	if testBuild.Config.CpuShares != testDeployment.CPU {
		t.Errorf("CpuShares expected %s, got %s", testDeployment.CPU,
			testBuild.Config.CpuShares)
	}
	if testBuild.Config.Image != testDeployment.Image {
		t.Errorf("Image expected %s, got %s", testDeployment.Image,
			testBuild.Config.Image)
	}
	if testBuild.Config.Memory != testDeployment.Memory {
		t.Errorf("Memory expected %s, got %s", testDeployment.Memory,
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
	_, err := ListAllContainers(testHost)
	if err != nil {
		t.Errorf("Error: %s", err)
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

func TestInspectContainer(t *testing.T) {
	containerInfo, errFrom, err := DeployNewContainer(testHost, testDeployment)
	if err != nil {
		t.Errorf("TestInspectContainer error: %s thrown by: \n", err, errFrom)

	}
	if len(containerInfo.Warnings) != 0 {
		t.Errorf("TestInspectContainer warning thrown: %+v \n", containerInfo.Warnings)
	}
	info, err := InspectContainer(testHost, containerInfo.Id)
	if err != nil {
		t.Errorf("TestInspectContainer error: %s \n", err)
	}
	if !info.Exists {
		t.Errorf("TestInspectContainer doesn't exist: %+v inspect response %+v \n", containerInfo, info)
	}
	// test with known bad ID
    info, err = InspectContainer(testHost, (containerInfo.Id + "badid"))
    if err != nil {
		t.Errorf("TestInspectContainer bad id error: %s \n", err)
	}
    if info.Exists {
        t.Errorf("TestInspectContainer bad id exists is true: %+v \n", info)
    }
	// test with no id
	info, err = InspectContainer(testHost, "")
	if err != nil {
		t.Errorf("TestInspectContainer no id error: %s \n", err)
	}
    if info.Exists {
        t.Errorf("TestInspectContainer no id exists is true: %+v \n", info)
    }
	//t.Errorf("TestInspectContainer warning thrown: %+v \n", info)
}

// 
// need to implement
// 
func TestStartStopContainer(t *testing.T) {
	containerInfo, errFrom, err := DeployNewContainer(testHost, testDeployment)
	if err != nil {
		t.Errorf("TestStartStopContainer error: %s thrown by: \n", err, errFrom)

	}
	if len(containerInfo.Warnings) != 0 {
		t.Errorf("TestStartStopContainer warning thrown: %+v \n", containerInfo.Warnings)
	}
	resp, err := StartContainer(testHost, containerInfo.Id)
	if err != nil {
		t.Errorf("TestStartStopContainer error: %s \n", err)
	}
	if resp.StatusCode != 204 {
	    msg, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("TestStartStopContainer Unexpected Status Code: %+v, %+v, %+v \n", resp.StatusCode, msg, containerInfo.Id)
	}
	info, err := InspectContainer(testHost, containerInfo.Id)
	if err != nil {
		t.Errorf("TestStartStopContainer error: %s \n", err)
	}
	if !info.State.Running {
	    t.Errorf("TestStartStopContainer expected container to be running: %+v \n", info.State)
	}
	resp, err = StopContainer(testHost, containerInfo.Id)
	if err != nil {
		t.Errorf("TestStartStopContainer error: %s \n", err)
	}
	if resp.StatusCode != 204 {
	    msg, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Unexpected Status Code: %+v, %+v, %+v \n", resp.StatusCode, msg, containerInfo.Id)
	}
	info, err = InspectContainer(testHost, containerInfo.Id)
	if err != nil {
		t.Errorf("TestStartStopContainer error: %s \n", err)
	}
	if info.State.Running {
	    t.Errorf("TestStartStopContainer expected container to be running: %+v \n", info.State)
	}
}

// 
// need to implement
// 
func TestDeleteContainer(t *testing.T) {
	containerInfo, errFrom, err := DeployNewContainer(testHost, testDeployment)
	if err != nil {
		t.Errorf("TestDeleteContainer error: %s thrown by: \n", err, errFrom)

	}
	if len(containerInfo.Warnings) != 0 {
		t.Errorf("TestDeleteContainer warning thrown: %+v \n", containerInfo.Warnings)
	}
	resp, err := DeleteContainer(testHost, containerInfo.Id)
	if err != nil {
		t.Errorf("TestDeleteContainer error: %s \n", err)
	}
	if resp.StatusCode != 204 {
	    msg, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("TestDeleteContainer Unexpected Status Code: %+v, %+v, %+v \n", resp.StatusCode, msg, containerInfo.Id)
	}
	info, err := InspectContainer(testHost, containerInfo.Id)
	if err != nil {
		t.Errorf("TestDeleteContainer error: %s \n", err)
	}
	if info.Exists {
	    t.Errorf("TestDeleteContainer container exists?: %+v \n", info)
	}
}

//
// disabled for now, takes a long time and we probably
// don't really need it really (just searches the docker repo)
// until we get our own private repo (if that's what we want?)
//
//func TestSearchForImage(t *testing.T) {
//resp, err := SearchForImage(testHost, testDeployment)
//if err != nil {
//t.Errorf("TestSearchForImage error: %s resp: %s \n", err, resp)
//}
// for now, we're just testing that something is being returned
// i'll probably just disable this test, because it takes a long time
// (has to query the docker repo, eventually we'll be querying our own repo)
//if resp[0].Name == "" {
//t.Errorf("TestSearchForImage no response?: %s\n", resp[0].Name)
//}
//}
