// Copyright (C) 2016  Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the LICENSE file.

package beater

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"path"
	"time"

	pb "github.com/openconfig/gnmi/proto/gnmi"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/aristanetworks/goarista/gnmi"
	"github.com/aristanetworks/openconfigbeat/config"
)

type Openconfigbeat struct {
	done      chan struct{}
	config    config.Config
	client    beat.Client
	paths     [][]string
	responses map[string]chan *pb.SubscribeResponse
	errors    map[string]chan error
	events    chan beat.Event
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	strPaths := config.Paths
	if len(strPaths) == 0 {
		strPaths = []string{"/"}
	}
	return &Openconfigbeat{
		done:      make(chan struct{}),
		config:    config,
		paths:     gnmi.SplitPaths(strPaths),
		responses: make(map[string]chan *pb.SubscribeResponse),
		errors:    make(map[string]chan error),
		events:    make(chan beat.Event),
	}, nil
}

// TODO: move these to goarista.git?
func convertDelete(dataset string, prefix string, delete *pb.Path) common.MapStr {
	m := common.MapStr{}
	m.Put("dataset", dataset)
	m.Put("path", prefix)
	// TODO: need to remove path suffixes here?
	m.Put("delete", gnmi.StrPath(delete))
	return m
}

func formatValue(update *pb.Update) (interface{}, error) {
	if update.Value != nil {
		switch update.Value.Type {
		case pb.Encoding_JSON, pb.Encoding_JSON_IETF:
			decoder := json.NewDecoder(bytes.NewReader(update.Value.Value))
			//decoder.UseNumber()
			var output interface{}
			err := decoder.Decode(&output)
			return output, err
		case pb.Encoding_BYTES, pb.Encoding_PROTO:
			return base64.StdEncoding.EncodeToString(update.Value.Value), nil
		case pb.Encoding_ASCII:
			return string(update.Value.Value), nil
		}
	}
	switch v := update.Val.GetValue().(type) {
	case *pb.TypedValue_StringVal:
		return v.StringVal, nil
	case *pb.TypedValue_JsonIetfVal:
		var output interface{}
		err := json.Unmarshal(v.JsonIetfVal, &output)
		return output, err
	case *pb.TypedValue_JsonVal:
		var output interface{}
		err := json.Unmarshal(v.JsonVal, &output)
		return output, err
	case *pb.TypedValue_IntVal:
		return v.IntVal, nil
	case *pb.TypedValue_UintVal:
		return v.UintVal, nil
	case *pb.TypedValue_BoolVal:
		return v.BoolVal, nil
	case *pb.TypedValue_BytesVal:
		return base64.StdEncoding.EncodeToString(v.BytesVal), nil
	case *pb.TypedValue_DecimalVal:
		// TODO: figure out a representation
		return nil, nil
	case *pb.TypedValue_FloatVal:
		return v.FloatVal, nil
	case *pb.TypedValue_LeaflistVal:
		// TODO: figure out a representation
		return nil, nil
	case *pb.TypedValue_AsciiVal:
		return v.AsciiVal, nil
	case *pb.TypedValue_AnyVal:
		// TODO: figure out a representation
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected type: %s", v)
	}
}

func convertUpdate(dataset string, prefix string, update *pb.Update) (common.MapStr,
	error) {
	m := common.MapStr{}
	m.Put("dataset", dataset)
	m.Put("path", path.Join(prefix, gnmi.StrPath(update.Path)))
	outputValue, err := formatValue(update)
	if err != nil {
		return nil, fmt.Errorf("Malformed update value: %s", err)
	}
	if _, ok := outputValue.(map[string]interface{}); !ok {
		k := fmt.Sprintf("%T", outputValue)
		m["update"] = map[string]interface{}{k: outputValue}
	} else {
		m["update"] = outputValue
	}
	return m, nil
}

// recv listens for SubscribeResponse notifications on a stream, and publishes the
// JSON representation of the notifications it receives on a channel
func (bt *Openconfigbeat) recv(host string) {
	respChan := bt.responses[host]
	errChan := bt.errors[host]
	for {
		select {
		case err := <-errChan:
			logp.Err("error from %s: %s", host, err)
		case response := <-respChan:
			update := response.GetUpdate()
			if update == nil {
				continue
			}
			timestamp := time.Unix(update.Timestamp/1e9, update.Timestamp%1e9)
			prefix := update.GetPrefix()
			prefixStr := "/"
			if prefix != nil {
				prefixStr += gnmi.StrPath(prefix)
			}
			events := []beat.Event{}
			fields := common.MapStr{}
			for _, del := range update.GetDelete() {
				output := convertDelete(host, prefixStr, del)
				fields.Update(output)
			}
			flush := func() {
				event := beat.Event{
					Timestamp: timestamp,
					Fields:    fields,
				}
				events = append(events, event)
			}
			if len(fields) > 0 {
				flush()
			}
			fields = common.MapStr{}
			for _, up := range update.GetUpdate() {
				output, err := convertUpdate(host, prefixStr, up)
				if err != nil {
					logp.Err(err.Error())
					continue
				}
				if len(fields) > 0 && fields["path"] != output["path"] {
					flush()
					fields = output
				} else {
					fields.Update(output)
				}
			}
			if len(fields) > 0 {
				flush()
			}
			for _, event := range events {
				select {
				case bt.events <- event:
				case <-bt.done:
					return
				}
			}
		}
	}
}

// recvAll listens for SubscribeResponse notifications on all streams
func (bt *Openconfigbeat) recvAll() {
	for addr := range bt.responses {
		go bt.recv(addr)
	}
}

func (bt *Openconfigbeat) Run(b *beat.Beat) error {
	logp.Info("openconfigbeat is running! Hit CTRL-C to stop it.")

	// Connect to elastisearch
	var err error
	if bt.client, err = b.Publisher.Connect(); err != nil {
		return err
	}

	// Connect the OpenConfig client
	for _, addr := range bt.config.Addresses {
		gnmiConfig := &gnmi.Config{
			Addr:     addr,
			Username: bt.config.Username,
			Password: bt.config.Password,
			TLS:      bt.config.TLS,
		}
		ctx := gnmi.NewContext(context.Background(), gnmiConfig)
		client, err := gnmi.Dial(gnmiConfig)
		if err != nil {
			return err
		}

		logp.Info("Connected to %s", addr)

		// Subscribe
		respChan := make(chan *pb.SubscribeResponse)
		errChan := make(chan error)
		go gnmi.Subscribe(ctx, client, &gnmi.SubscribeOptions{Paths: bt.paths}, respChan, errChan)
		device, _, err := net.SplitHostPort(addr)
		if err != nil {
			return err
		}
		bt.responses[device] = respChan
		bt.errors[device] = errChan
	}

	bt.recvAll()
	for {
		select {
		case <-bt.done:
			return nil
		case event := <-bt.events:
			event.Fields["type"] = b.Info.Name
			logp.Info("Publishing: %s", event)
			bt.client.Publish(event)
		}
	}
}

func (bt *Openconfigbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
