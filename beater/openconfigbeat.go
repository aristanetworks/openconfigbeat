package beater

import (
	"fmt"
	"strings"
	"time"

	"github.com/aristanetworks/goarista/openconfig"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
	"google.golang.org/grpc"

	"github.com/aristanetworks/openconfigbeat/config"
)

type Openconfigbeat struct {
	beatConfig *config.Config
	done       chan struct{}
	addresses  []string
	paths      []openconfig.Path
	client     publisher.Client
}

// Creates beater
func New() *Openconfigbeat {
	return &Openconfigbeat{
		done: make(chan struct{}),
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
		bt.paths = []openconfig.Path{openconfig.Path{Element: []string{"/"}}}
	} else {
		for _, path := range *config.Paths {
			bt.paths = append(bt.paths,
				openconfig.Path{Element: strings.Split(path, "/")})
		}
	}

	return nil
}

func (bt *Openconfigbeat) Setup(b *beat.Beat) error {

	bt.client = b.Publisher.Connect()

	return nil
}

func (bt *Openconfigbeat) Run(b *beat.Beat) error {
	logp.Info("openconfigbeat is running! Hit CTRL-C to stop it.")

	addr := bt.addresses[0]
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	logp.Info("Connected to %s", addr)
	defer conn.Close()
	openconfig.NewOpenConfigClient(conn)

	// TODO: subscribe

	counter := 1
	for {

		/* TODO: read subscribe responses and publish events
		select {
		case <-bt.done:
			return nil
		}
		*/

		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
			"counter":    counter,
		}
		if !bt.client.PublishEvent(event) {
			return fmt.Errorf("Failed to publish %dth event", counter)
		}
		logp.Info("Event sent")
		counter++
	}
}

func (bt *Openconfigbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (bt *Openconfigbeat) Stop() {
	close(bt.done)
}
