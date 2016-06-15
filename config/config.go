// Copyright (C) 2016  Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the LICENSE file.

// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"errors"
	"net"
)

type Config struct {
	Openconfigbeat OpenconfigbeatConfig
}

type OpenconfigbeatConfig struct {
	Addresses *[]string `config:"addresses"`
	Paths     *[]string `config:"paths"`
}

func (c *OpenconfigbeatConfig) Validate() error {
	if c.Addresses == nil || len(*c.Addresses) == 0 {
		return errors.New("Please specify at least a device to connect to in 'addresses'")
	}
	seen := make(map[string]bool, len(*c.Addresses))
	for _, hostPort := range *c.Addresses {
		host, _, err := net.SplitHostPort(hostPort)
		if err != nil {
			return err
		}
		if _, found := seen[host]; found {
			return errors.New("Duplicate host(s) found in 'addresses'")
		}
		seen[host] = true
	}
	return nil
}
