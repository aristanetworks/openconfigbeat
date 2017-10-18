// Copyright (C) 2017  Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/aristanetworks/glog"
	"github.com/aristanetworks/goarista/gnmi"
)

// TODO: Make this more clear
var help = `Usage of gnmi:
gnmi -addr ADDRESS:PORT [options...]
  capabilities
  get PATH+
  subscribe PATH+
  ((update|replace PATH JSON)|(delete PATH))+
`

func exitWithError(s string) {
	flag.Usage()
	fmt.Fprintln(os.Stderr, s)
	os.Exit(1)
}

func main() {
	cfg := &gnmi.Config{}
	flag.StringVar(&cfg.Addr, "addr", "", "Address of gNMI gRPC server")
	flag.StringVar(&cfg.CAFile, "cafile", "", "Path to server TLS certificate file")
	flag.StringVar(&cfg.CertFile, "certfile", "", "Path to client TLS certificate file")
	flag.StringVar(&cfg.KeyFile, "keyfile", "", "Path to client TLS private key file")
	flag.StringVar(&cfg.Password, "password", "", "Password to authenticate with")
	flag.StringVar(&cfg.Username, "username", "", "Username to authenticate with")
	flag.BoolVar(&cfg.TLS, "tls", false, "Enable TLS")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, help)
		flag.PrintDefaults()
	}
	flag.Parse()
	if cfg.Addr == "" {
		exitWithError("error: address not specified")
	}

	args := flag.Args()

	ctx := gnmi.NewContext(context.Background(), cfg)
	client := gnmi.Dial(cfg)

	var setOps []*gnmi.Operation
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "capabilities":
			if len(setOps) != 0 {
				exitWithError("error: 'capabilities' not allowed after 'merge|replace|delete'")
			}
			err := gnmi.Capabilities(ctx, client)
			if err != nil {
				glog.Fatal(err)
			}
			return
		case "get":
			if len(setOps) != 0 {
				exitWithError("error: 'get' not allowed after 'merge|replace|delete'")
			}
			err := gnmi.Get(ctx, client, gnmi.SplitPaths(args[i+1:]))
			if err != nil {
				glog.Fatal(err)
			}
			return
		case "subscribe":
			if len(setOps) != 0 {
				exitWithError("error: 'subscribe' not allowed after 'merge|replace|delete'")
			}
			err := gnmi.Subscribe(ctx, client, gnmi.SplitPaths(args[i+1:]))
			if err != nil {
				glog.Fatal(err)
			}
			return
		case "update", "replace", "delete":
			if len(args) == i+1 {
				exitWithError("error: missing path")
			}
			op := &gnmi.Operation{
				Type: args[i],
			}
			i++
			op.Path = gnmi.SplitPath(args[i])
			if op.Type != "delete" {
				if len(args) == i+1 {
					exitWithError("error: missing JSON")
				}
				i++
				op.Val = args[i]
			}
			setOps = append(setOps, op)
		default:
			exitWithError(fmt.Sprintf("error: unknown operation %q", args[i]))
		}
	}
	if len(setOps) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	err := gnmi.Set(ctx, client, setOps)
	if err != nil {
		glog.Fatal(err)
	}

}
