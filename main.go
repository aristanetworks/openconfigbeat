// Copyright (C) 2016  Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/aristanetworks/openconfigbeat/cmd"
	"github.com/aristanetworks/openconfigbeat/gen/fields"
	"github.com/aristanetworks/openconfigbeat/gen/openconfigbeat"
)

const (
	fieldsFile     = "fields.yml"
	openconfigFile = "openconfigbeat.yml"
)

func main() {
	if _, err := os.Stat(openconfigFile); os.IsNotExist(err) {
		fmt.Printf("%s not found, generating the default one\n", openconfigFile)
		if err = openconfigbeat.RestoreAssets(".", openconfigFile); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	if _, err := os.Stat(fieldsFile); os.IsNotExist(err) {
		fmt.Printf("%s not found, generating the default one\n", fieldsFile)
		if err = fields.RestoreAssets(".", fieldsFile); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
