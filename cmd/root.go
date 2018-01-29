package cmd

import (
	"github.com/aristanetworks/openconfigbeat/beater"

	cmd "github.com/elastic/beats/libbeat/cmd"
)

// Name of this beat
var Name = "openconfigbeat"

var Version string

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmd(Name, Version, beater.New)
