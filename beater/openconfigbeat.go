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
	beatConfig        *config.Config
	done              chan struct{}
	addresses         []string
	paths             []*openconfig.Path
	subscribeClient   openconfig.OpenConfig_SubscribeClient
	subscribeResponse chan map[string]interface{}
	client            publisher.Client
}

// Creates beater
func New() *Openconfigbeat {
	return &Openconfigbeat{
		done:              make(chan struct{}),
		subscribeResponse: make(chan map[string]interface{}),
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

// recvLoop listens for SubscribeResponse notifications on stream, and publishes the
// JSON representation of the notifications it receives on a channel
func (bt *Openconfigbeat) recvLoop() {
	for {
		response, err := bt.subscribeClient.Recv()
		if err != nil {
			logp.Err(err.Error())
			return
		}
		update := response.GetUpdate()
		if update == nil {
			logp.Err("Unhandled subscribe response: %s", response)
			return
		}
		updateMap, err := openconfig.NotificationToMap(update, elasticsearch.EscapeFieldName)
		if err != nil {
			logp.Err(err.Error())
			return
		}
		select {
		case bt.subscribeResponse <- updateMap:
		case <-bt.done:
			return
		}
	}
}

func (bt *Openconfigbeat) Run(b *beat.Beat) error {
	logp.Info("openconfigbeat is running! Hit CTRL-C to stop it.")
	var err error

	// Connect the OpenConfig client
	addr := bt.addresses[0]
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	logp.Info("Connected to %s", addr)
	defer conn.Close()
	client := openconfig.NewOpenConfigClient(conn)

	// Subscribe
	bt.subscribeClient, err = client.Subscribe(context.Background())
	if err != nil {
		return err
	}
	defer bt.subscribeClient.CloseSend()
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
		err = bt.subscribeClient.Send(sub)
		if err != nil {
			return err
		}
	}

	// Main loop
	device, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	go bt.recvLoop()
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case response := <-bt.subscribeResponse:
			event := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"type":       b.Name,
				"counter":    counter,
				device:       response,
			}
			if !bt.client.PublishEvent(event) {
				return fmt.Errorf("Failed to publish %dth event", counter)
			}
			logp.Info("Event sent")
			counter++
		}
	}
}

func (bt *Openconfigbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (bt *Openconfigbeat) Stop() {
	close(bt.done)
}
