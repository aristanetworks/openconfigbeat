// Copyright (C) 2016  Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the LICENSE file.

package main

// This file is mandatory as otherwise the openconfigbeat.test binary is not generated correctly.

import (
	"flag"
	"testing"
)

var systemTest *bool

func init() {
	systemTest = flag.Bool("systemTest", false, "Set to true when running system tests")
}

// Test started when the test binary is started. Only calls main.
func TestSystem(t *testing.T) {

	if *systemTest {
		main()
	}
}
