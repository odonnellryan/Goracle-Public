package main

import (
	"fmt"
	"testing"
)

var testConfigValues = nginxConfigValues{"hostname", "localhost", "9999"}

var testConfigFile = "upstream hostname { server localhost:9999; } server { listen 80; server_name hostname; location / { proxy_pass http://hostname; }}"
var testConfigHash = CryptToHex(testConfigFile)

var testNginxConfig = NginxConfig{
	"hostname.GOOD",
	testConfigHash,
	testConfigFile,
	testConfigValues,
}

func TestBuildConfig(t *testing.T) {
	compareConfig := BuildNginxConfig(testConfigValues)
	testConfigHash := CryptToHex(testConfigFile)
	configWithHash := fmt.Sprintf("#%s\n%s", testConfigHash, testConfigFile)
	if compareConfig.configFile != configWithHash {
		t.Errorf("Got: \n %+v \n expecting: \n %+v",
			compareConfig.configFile, configWithHash)
	}
}
