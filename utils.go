package main

import (
	"strings"
	"fmt"
)

func CommaStringToSlice(s string) []string {
	return strings.Split(s, ",")
}

func ParseMongoEndpoint(endpoint string) (string, string, string, error) {
	portSeperatorIndex := strings.LastIndex(endpoint, ":")
	if portSeperatorIndex < 1 {
		err := fmt.Errorf("Invalid format of mongo endpoint. Correct syntax is host:port/db")
		return "", "0", "", err
	}
	host := endpoint[0 : portSeperatorIndex-1]
	slashIndex := strings.Index(endpoint, "/")
	if slashIndex < 1 {
		err := fmt.Errorf("Invalid format of mongo endpoint. Correct syntax is host:port/db")
		return "", "0", "", err
	}
	port := endpoint[portSeperatorIndex+1 : slashIndex-1]
	db := endpoint[slashIndex+1:]
	return host, port, db, nil
}