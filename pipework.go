package main

import (
	"fmt"
)

type Pipework struct {
	Err            error
	Container      string
	HostInterface  string
	Address        string
	Subnet         string
	Gateway        string
	Finished       error
}

func (p *Pipework) AddInterfaceToContainer(host Host) error {
    err := host.SendSSHCommand("CMD")
    if err != nil {
		return err
	}
	fmt.Printf("testing needs implementation")
	return nil
}
