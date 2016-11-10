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
	pb "github.com/openconfig/reference/rpc/openconfig"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/aristanetworks/openconfigbeat/config"
)

type Openconfigbeat struct {
	beatConfig       *config.Config
	done             chan struct{}
	paths            []*pb.Path
	subscribeClients map[string]pb.OpenConfig_SubscribeClient
	events           chan common.MapStr
	client           publisher.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	conf := config.DefaultConfig
	err := b.RawConfig.Unpack(&conf)
	if err != nil {
		return nil, err
	}
	var paths []*pb.Path
	if len(conf.Openconfigbeat.Paths) == 0 {
		paths = []*pb.Path{&pb.Path{Element: []string{"/"}}}
	} else {
		for _, path := range conf.Openconfigbeat.Paths {
			paths = append(paths, &pb.Path{Element: strings.Split(path, "/")})
		}
	}
	return &Openconfigbeat{
		beatConfig:       &conf,
		paths:            paths,
		done:             make(chan struct{}),
		subscribeClients: make(map[string]pb.OpenConfig_SubscribeClient),
		events:           make(chan common.MapStr),
	}, nil
}

/// *** Beater interface methods ***///

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
		notifMap["@timestamp"] = common.Time(time.Unix(timestampNs/1e9,
			timestampNs%1e9))
		delete(notifMap, "timestamp")
		select {
		case bt.events <- notifMap:
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
	bt.client = b.Publisher.Connect()

	// Connect the OpenConfig client
	for _, addr := range bt.beatConfig.Openconfigbeat.Addresses {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			logp.Err("Failed to connect to %s: %s", addr, err.Error())
			continue
		}
		logp.Info("Connected to %s", addr)
		defer conn.Close()
		client := pb.NewOpenConfigClient(conn)

		// Subscribe
		s, err := client.Subscribe(context.Background())
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
						&pb.Subscription{
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
			logp.Info("Published: %s", event)
		}
	}
}

func (bt *Openconfigbeat) Stop() {
	close(bt.done)
}
