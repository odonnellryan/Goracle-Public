package main

//
// this is for the Organization structure
// Organizations are groups of Users
// every user has to be a part of an Organization
// Organizations are used to keep a list of all interconnected
// machines/bridges over the vpn
//
//
// containers are tied to users and exist on hosts
// host bridge->ip address are owned by organizations
// Example (by order of mapping to):
// Organization <-> Organization -> Users -> Host -> Containers
// an Organization is created automatically for users when they create
// a new account (Ryan's Containers)
//

type HostAddress struct {
    Bridge string
    Addresses []string
}

type Organization struct {
    Name            string
    Users           []User
    Administrators  []User
    PartnerOrgs     map[string]struct {}
    HostAddress     map[string]HostAddress // hostname(string) : HostAddress
}