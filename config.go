//  config.go - functions to (re)load config
// 	Copyright (C) 2020  Joshua French
//
// 	This program is free software: you can redistribute it and/or modify
// 	it under the terms of the GNU Affero General Public License as published
// 	by the Free Software Foundation, either version 3 of the License, or
// 	(at your option) any later version.
//
// 	This program is distributed in the hope that it will be useful,
// 	but WITHOUT ANY WARRANTY; without even the implied warranty of
// 	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// 	GNU Affero General Public License for more details.
//
// 	You should have received a copy of the GNU Affero General Public License
// 	along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config is a struct
type Config struct {
	sync.RWMutex
	Interfaces map[string]IfaceConfig
	Routes     map[int]RouteConfig
}

// IfaceConfig is the config representation of an Iface
type IfaceConfig struct {
	Addr      int
	Protocol  string
	ConnProto string
	Device    string
	IP        string
	Port      int16
}

// RouteConfig is the config representation of a route
type RouteConfig struct {
	Interface string
	NAT       []NATConfig
}

// NATConfig is the config representation of NAT for a Route
type NATConfig struct {
	Internal int
	External int
}

// findConfigFile finds a named file with a yml or yaml extension in either CWD or /etc/
func findConfigFile(name string) (string, error) {
	allMatches := make([]string, 0)
	for _, path := range []string{"", "/etc/"} {
		for _, ext := range []string{".yml", ".yaml"} {
			matches, err := filepath.Glob(path + name + ext)
			if err != nil {
				return "", err
			}
			allMatches = append(allMatches, matches...)
		}
	}

	for _, match := range allMatches {
		return match, nil
	}

	return "", errors.New("Unable to find config file in CWD or /etc/")
}

func loadConfig(configPath string, config *Config) {
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	config.Lock()
	err = yaml.Unmarshal(configFile, &config)
	config.Unlock()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded Config: %+v\n", config)
}

func syncRouterConfig(config *Config, router *Router) {

}
