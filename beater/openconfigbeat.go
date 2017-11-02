// Copyright (C) 2016  Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the LICENSE file.

package beater

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	pb "github.com/openconfig/reference/rpc/openconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/aristanetworks/goarista/elasticsearch"
	"github.com/aristanetworks/goarista/openconfig"
	"github.com/aristanetworks/openconfigbeat/config"
)

type Openconfigbeat struct {
	done             chan struct{}
	config           config.Config
	client           beat.Client
	paths            []*pb.Path
	subscribeClients map[string]pb.OpenConfig_SubscribeClient
	events           chan beat.Event
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	var paths []*pb.Path
	if len(config.Paths) == 0 {
		paths = []*pb.Path{{Element: []string{"/"}}}
	} else {
		for _, path := range config.Paths {
			paths = append(paths, &pb.Path{Element: strings.Split(path, "/")})
		}
	}
	return &Openconfigbeat{
		done:             make(chan struct{}),
		config:           config,
		paths:            paths,
		subscribeClients: make(map[string]pb.OpenConfig_SubscribeClient),
		events:           make(chan beat.Event),
	}, nil
}

// recv listens for SubscribeResponse notifications on a stream, and publishes the
// JSON representation of the notifications it receives on a channel
func (bt *Openconfigbeat) recv(host string) {
	for {
		response, err := bt.subscribeClients[host].Recv()
		if err != nil {
			logp.Err(err.Error())
			return
		}
		update := response.GetUpdate()
		if update == nil {
			continue
		}
		notifMap, err := openconfig.NotificationToMap(host, update,
			elasticsearch.EscapeFieldName)
		if err != nil {
			logp.Err(err.Error())
			continue
		}
		timestamp, found := notifMap["timestamp"]
		if !found {
			logp.Err("Malformed subscribe response: %s", notifMap)
			return
		}
		timestampNs, ok := timestamp.(int64)
		if !ok {
			logp.Err("Malformed timestamp: %s", timestamp)
			continue
		}
		event := beat.Event{
			Timestamp: time.Unix(timestampNs/1e9, timestampNs%1e9),
			Fields:    notifMap,
		}
		delete(event.Fields, "timestamp")
		select {
		case bt.events <- event:
		case <-bt.done:
			return
		}
	}
}

// recvAll listens for SubscribeResponse notifications on all streams
func (bt *Openconfigbeat) recvAll() {
	for i := range bt.subscribeClients {
		go bt.recv(i)
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
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			logp.Err("Failed to connect to %s: %s", addr, err.Error())
			continue
		}
		logp.Info("Connected to %s", addr)
		defer conn.Close()
		client := pb.NewOpenConfigClient(conn)

		// Subscribe
		ctx := context.Background()
		if bt.config.Username != "" {
			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
				"username", bt.config.Username,
				"password", bt.config.Password))
		}
		s, err := client.Subscribe(ctx)
		if err != nil {
			logp.Err("Failed to subscribe from %s: %s", addr, err.Error())
			continue
		}
		defer s.CloseSend()
		device, _, err := net.SplitHostPort(addr)
		if err != nil {
			return err
		}
		bt.subscribeClients[device] = s
	}
	for _, path := range bt.paths {
		sub := &pb.SubscribeRequest{
			Request: &pb.SubscribeRequest_Subscribe{
				Subscribe: &pb.SubscriptionList{
					Subscription: []*pb.Subscription{
						{
							Path: path,
						},
					},
				},
			},
		}
		for _, s := range bt.subscribeClients {
			logp.Info("Sending subscribe request: %s", sub)
			err := s.Send(sub)
			if err != nil {
				return err
			}
		}
	}

	bt.recvAll()
	for {
		select {
		case <-bt.done:
			return nil
		case event := <-bt.events:
			event.Fields["type"] = b.Info.Name
			bt.client.Publish(event)
			logp.Info("Published: %s", event)
		}
	}
}

func (bt *Openconfigbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
