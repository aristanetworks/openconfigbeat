// Copyright (C) 2016  Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/aristanetworks/openconfigbeat/beater"
)

func main() {
	err := beat.Run("openconfigbeat", "", beater.New())
	if err != nil {
		os.Exit(1)
	}
}
