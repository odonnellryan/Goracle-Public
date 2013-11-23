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
	"net/http"
)

func GenerateKey(w http.ResponseWriter, r *http.Request) (string, error) {
	key, kerror := rsa.GenerateKey(rand.Reader, 2048)
	if kerror != nil {
		log.Println(ErrorMessages["EncodingError"])
		return ErrorMessages["EncodingError"], errors.New(ErrorMessages["EncodingError"])
	}
	// save to pem
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	// response, error
	return string(pemdata), nil
}
