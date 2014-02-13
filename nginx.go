package main

import (
	"fmt"
)

type nginxConfigValues struct {
	hostname       string
	upstreamServer string
	upstreamPort   string
}

// eventually, will want this to be the proxy deploy (basically just nginx deployment)
// if d.WebPort != "" {
//		nginxConfig := nginxConfigValues{
//			hostname:       d.Hostname,
//			upstreamServer: d.IP,
//			upstreamPort:   d.WebPort,
//		}
//		d.NginxConfig = BuildNginxConfig(nginxConfig)
//	}

// hostname, (upstreamServer + upstreamPort), hostname, hostname.
var HttpConfigFile = `
    upstream %s  {
      %s
    }
    server {
      listen 80; 
      server_name %s;
      location / {
        proxy_pass  http://%s;
      }
    }
`

type NginxConfig struct {
	// nginx is set to only load *.GOOD
	configName   string
	configFile   string
	configValues nginxConfigValues
}

func BuildNginxConfig(values nginxConfigValues) NginxConfig {
	config := fmt.Sprintf(HttpConfigFile, values.hostname, (values.upstreamServer + ":" + values.upstreamPort),
		values.hostname, values.hostname)
	// nginx is set to only load *.GOOD
	return NginxConfig{(values.hostname + ".GOOD"), config, values}
}
