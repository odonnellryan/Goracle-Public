package main

import (
	//"fmt"
	//"strings"
	//"encoding/json"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
)

func GenerateKey() (string, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println("Error generating key: %s", err)
		return "", errors.New(ErrorMessages["EncodingError"])
	}
	// save to pem
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	return string(pemdata), nil
}
