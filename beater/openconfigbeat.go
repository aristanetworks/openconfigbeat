// Copyright (C) 2016  Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the LICENSE file.

package beater

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/aristanetworks/goarista/elasticsearch"
	"github.com/aristanetworks/goarista/openconfig"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/aristanetworks/openconfigbeat/config"
)

type Openconfigbeat struct {
	beatConfig       *config.Config
	done             chan struct{}
	addresses        []string
	paths            []*openconfig.Path
	subscribeClients map[string]openconfig.OpenConfig_SubscribeClient
	events           chan common.MapStr
	client           publisher.Client
}

// Creates beater
func New() *Openconfigbeat {
	return &Openconfigbeat{
		done:             make(chan struct{}),
		subscribeClients: make(map[string]openconfig.OpenConfig_SubscribeClient),
		events:           make(chan common.MapStr),
	}
}

/// *** Beater interface methods ***///

func (bt *Openconfigbeat) Config(b *beat.Beat) error {

	// Load beater beatConfig
	err := b.RawConfig.Unpack(&bt.beatConfig)
	if err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}

	config := bt.beatConfig.Openconfigbeat
	bt.addresses = *config.Addresses
	if config.Paths == nil {
		bt.paths = []*openconfig.Path{&openconfig.Path{Element: []string{"/"}}}
	} else {
		for _, path := range *config.Paths {
			bt.paths = append(bt.paths,
				&openconfig.Path{Element: strings.Split(path, "/")})
		}
	}

	return nil
}

func (bt *Openconfigbeat) Setup(b *beat.Beat) error {

	bt.client = b.Publisher.Connect()

	return nil
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
		updateMap, err := openconfig.NotificationToMap(update,
			elasticsearch.EscapeFieldName)
		if err != nil {
			logp.Err(err.Error())
			continue
		}
		updateMap["device"] = host
		timestamp, found := updateMap["_timestamp"]
		if !found {
			logp.Err("Malformed subscribe response: %s", updateMap)
			continue
		}
		timestampNs, ok := timestamp.(int64)
		if !ok {
			logp.Err("Malformed timestamp: %s", timestamp)
			continue
		}
		updateMap["@timestamp"] = common.Time(time.Unix(timestampNs/1e9,
			timestampNs%1e9))
		delete(updateMap, "_timestamp")
		select {
		case bt.events <- updateMap:
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

	// Connect the OpenConfig client
	for _, addr := range bt.addresses {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			return err
		}
		logp.Info("Connected to %s", addr)
		defer conn.Close()
		client := openconfig.NewOpenConfigClient(conn)

		// Subscribe
		s, err := client.Subscribe(context.Background())
		if err != nil {
			return err
		}
		defer s.CloseSend()
		device, _, err := net.SplitHostPort(addr)
		if err != nil {
			return err
		}
		bt.subscribeClients[device] = s
	}
	for _, path := range bt.paths {
		sub := &openconfig.SubscribeRequest{
			Request: &openconfig.SubscribeRequest_Subscribe{
				Subscribe: &openconfig.SubscriptionList{
					Subscription: []*openconfig.Subscription{
						&openconfig.Subscription{
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
			event["type"] = b.Name
			if !bt.client.PublishEvent(event) {
				return fmt.Errorf("Failed to publish event %q", event)
			}
		}
	}
}

func (bt *Openconfigbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (bt *Openconfigbeat) Stop() {
	close(bt.done)
}
