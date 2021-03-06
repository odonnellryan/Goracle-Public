package main

import (
	"crypto/sha512"
	"fmt"
	"io"
	"strings"
	//"bytes"
	"encoding/hex"
)

type KeyValuePair struct {
	Key   string
	Value string
}

type PortBinding struct {
	HostIp   string
	HostPort string
}

func CommaStringToSlice(s string) []string {
	return strings.Split(s, ",")
}

func ParseMongoEndpoint(endpoint string) (string, string, string, error) {
	portSeperatorIndex := strings.LastIndex(endpoint, ":")
	if portSeperatorIndex < 1 {
		err := fmt.Errorf("Invalid format of mongo endpoint. Correct syntax is host:port/db")
		return "", "0", "", err
	}
	host := endpoint[0:portSeperatorIndex]
	slashIndex := strings.Index(endpoint, "/")
	if slashIndex < 1 {
		err := fmt.Errorf("Invalid format of mongo endpoint. Correct syntax is host:port/db")
		return "", "0", "", err
	}
	port := endpoint[portSeperatorIndex+1 : slashIndex]
	db := endpoint[slashIndex+1:]
	return host, port, db, nil
}

func CryptToHex(s string) string {
	h := sha512.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}
