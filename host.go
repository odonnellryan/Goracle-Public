package main

//
// host is for actions taken on any specific docker host.
//

import (
    "code.google.com/p/go.crypto/ssh"
    "fmt"
    "encoding/json"
)

type Host struct {
	Hostname    string
	Address     string
	User        string
	Password    string
	Containers  int
	SSHUser     string
	SSHPassword string
}

func (h *Host) SendSSHCommand(command string) error {
    // Dial code is taken from the ssh package example
    config := &ssh.ClientConfig{
        User: h.SSHUser,
        Auth: []ssh.AuthMethod{
            ssh.Password(h.SSHPassword),
        },
    }
    client, err := ssh.Dial("tcp", "127.0.0.1:22", config)
    if err != nil {
        panic("Failed to dial: " + err.Error())
    }

    session, err := client.NewSession()
    if err != nil {
        panic("Failed to create session: " + err.Error())
    }
    defer session.Close()

    go func() {
        w, _ := session.StdinPipe()
        defer w.Close()
        content := "123456789\n"
        fmt.Fprintln(w, "C0644", len(content), "testfile")
        fmt.Fprint(w, content)
        fmt.Fprint(w, "\x00")
    }()
    if err := session.Run("/usr/bin/scp -qrt ./"); err != nil {
        panic("Failed to run: " + err.Error())
    }
    return nil
}

// for host, not docker really..
// maybe make a host.go?
func (h *Host) ListAllContainers() ([]ListContainerInfo, error) {
	containers := []ListContainerInfo{}
	command := "containers/json?all=1"
	resp, err := SendDockerCommand(*h, command, "GET", nil)
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