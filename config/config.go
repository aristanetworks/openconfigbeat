// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"errors"
)

type Config struct {
	Openconfigbeat OpenconfigbeatConfig
}

type OpenconfigbeatConfig struct {
	Addresses *[]string `config:"addresses"`
	Paths     *[]string `config:"paths"`
}

func (c *OpenconfigbeatConfig) Validate() error {
	if c.Addresses == nil {
		return errors.New("Please specify at least a device to connect to in 'addresses'")
	}
	// TODO: implement
	if len(*c.Addresses) > 1 {
		return errors.New("Connecting to more than one device not yet supported")
	}
	// TODO: implement
	if c.Paths != nil && len(*c.Paths) > 1 {
		return errors.New("Subscribing to more than one path not yet supported")
	}
	return nil
}
