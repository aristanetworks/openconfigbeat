// Copyright (C) 2016  Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the LICENSE file.

// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"fmt"
	"strconv"
	"strings"
)

type Config struct {
	DefaultPort int      `config:"default_port"`
	Addresses   []string `config:"addresses"`
	Paths       []string `config:"paths"`
	Username    string   `config:"username"`
	Password    string   `config:"password"`
	TLS         bool     `config:"tls"`
}

var DefaultConfig = Config{}

// Split a single option specified via -E into multiple options, if the value
// contains comma-separated values
func split(opts []string) []string {
	if len(opts) == 1 && strings.ContainsRune(opts[0], ',') {
		return strings.Split(opts[0], ",")
	}
	return opts
}

func (c *Config) Validate() error {
	if len(c.Addresses) == 0 {
		return fmt.Errorf("Please specify at least a device to connect to in 'addresses'")
	}
	c.Addresses = split(c.Addresses)
	for i, address := range c.Addresses {
		if !strings.ContainsRune(address, ':') {
			c.Addresses[i] = address + ":" + strconv.Itoa(c.DefaultPort)
		}
	}
	c.Paths = split(c.Paths)
	return nil
}
