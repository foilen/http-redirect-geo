package main

import (
	"encoding/json"
	"io/ioutil"
)

// Example:
// {
//     "port" : 8888,
//     "dbIpFile" : "dbip-city-lite.mmdb",
//     "redirectionUrls" : [
//         "https://tor.cdn.foilen.com",
//         "https://fra.cdn.foilen.com"
//     ]
// }

// RootConfiguration is the json configuration file
type RootConfiguration struct {
	Port            uint16
	DbIPFile        string
	RedirectionUrls []string
}

func getRootConfiguration(path string) (*RootConfiguration, error) {
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rootConfiguration RootConfiguration
	err = json.Unmarshal(jsonBytes, &rootConfiguration)

	return &rootConfiguration, err
}
