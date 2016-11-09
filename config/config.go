// Copyright (C) 2016  Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the LICENSE file.

// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"fmt"
)

type Config struct {
	Openconfigbeat OpenconfigbeatConfig
}

type OpenconfigbeatConfig struct {
	Addresses []string `config:"addresses"`
	Paths     []string `config:"paths"`
}

var DefaultConfig = Config{}

func (c *OpenconfigbeatConfig) Validate() error {
	if len(c.Addresses) == 0 {
		return fmt.Errorf("Please specify at least a device to connect to in 'addresses'")
	}
	return nil
}
