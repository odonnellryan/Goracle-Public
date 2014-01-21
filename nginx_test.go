package main

import (
    "testing"
    )

var testConfigValues = nginxConfigValues{"hostname", "localhost","9999"}

// hostname, (upstreamServer + upstreamPort), hostname, hostname
var testConfigFile = `
    upstream hostname  {
      localhost:9999
    }
    server {
      listen 80; 
      server_name hostname;
      location / {
        proxy_pass  http://hostname;
      }
    }
`

var testNginxConfig = NginxConfig {
    "hostname.GOOD",
    testConfigFile,
    testConfigValues,
}

func TestBuildConfig(t *testing.T) {
    compareConfig := BuildNginxConfig(testConfigValues)
    if compareConfig != testNginxConfig {
        t.Errorf("Got %+v, expecting %+v", compareConfig, testNginxConfig)
    }
}