package main

import (
    "os/exec"
    )
    
type Pipework struct {
    err             error
    container       string
    host_interface  string
    address         string
    subnet          string
    gateway         string
    finished        error
}

func (p *Pipework) AddInterfaceToContainer() {
    //exec.Command(
}