// Copyright (C) 2017 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the LICENSE file.

package beater

import (
	"os"
	"testing"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/outputs/elasticsearch"
	"github.com/elastic/beats/libbeat/outputs/outil"
	"github.com/elastic/beats/libbeat/publisher"
)

type mockBatch struct {
	updates []map[string]interface{}
}

func (m *mockBatch) Events() []publisher.Event {
	events := []publisher.Event{}
	for _, update := range m.updates {
		fields := common.MapStr{}
		fields.Update(update)
		fields.Update(map[string]interface{}{
			"dataset": "cairo",
			"type":    "Giuseppes-MacBook-Pro.local",
		})
		event := publisher.Event{
			Content: beat.Event{
				Timestamp: time.Now(),
				Fields:    fields,
			},
		}
		events = append(events, event)
	}
	return events
}

func (m *mockBatch) ACK() {
}
func (m *mockBatch) Drop()  {}
func (m *mockBatch) Retry() {}
func (m *mockBatch) RetryEvents(events []publisher.Event) {
	return
}
func (m *mockBatch) Cancelled() {}
func (m *mockBatch) CancelledEvents(events []publisher.Event) {
	return
}

func TestPublish(t *testing.T) {
	const tempSensor1 = "/Sysdb/environment/archer/temperature/status/TempSensor1"
	// TODO: test Notifications
	messages := []map[string]interface{}{
		{
			"path": tempSensor1 + "/maxTemperatureTime",
			"update": map[string]interface{}{
				"float64": 1510780140.838717,
			},
		}, {
			"path": tempSensor1 + "/temperature",
			"update": map[string]interface{}{
				"value": 32.315,
			},
		}, {"path": tempSensor1 + "/hwStatus",
			"update": map[string]interface{}{
				"bool": 1,
			},
		}, {
			"path": tempSensor1 + "/lastAlertRaisedTime",
			"update": map[string]interface{}{
				"float64": 1510779879.456627,
			},
		}, {
			"path": tempSensor1 + "/alertRaisedCount",
			"update": map[string]interface{}{
				"bool": 0,
			},
		}, {
			"path": tempSensor1 + "/name",
			"update": map[string]interface{}{
				"string": "TempSensor1",
			},
		}, {
			"path": tempSensor1 + "/generationId",
			"update": map[string]interface{}{
				"uint64": 0,
			},
		}, {
			"path": tempSensor1 + "/alertRaised",
			"update": map[string]interface{}{
				"bool": false,
			},
		}, {
			"path": tempSensor1 + "/maxTemperature",
			"update": map[string]interface{}{
				"value": 37.79170479593131,
			},
		}, {
			"path": tempSensor1 + "/temperature",
			"update": map[string]interface{}{
				"value": 32.644118526016115,
			},
		}, {
			"path": tempSensor1 + "/temperature",
			"update": map[string]interface{}{
				"value": 33.096784385446064,
			},
		},
	}
	host := os.Getenv("ELASTICSEARCH_HOST")
	if host == "" {
		host = "localhost"
	}
	index := "openconfigbeat-7.0.0-alpha1-2017.11.29"
	settings := elasticsearch.ClientSettings{
		URL:   "http://" + host + ":9200",
		Index: outil.MakeSelector(outil.ConstSelectorExpr(index)),
	}
	client, err := elasticsearch.NewClient(settings, nil)
	if err != nil {
		t.Fatal(err)
	}
	client.Connection.Delete(index, "", "", nil)
	// One by one
	for i, msg := range messages {
		if err = client.Publish(&mockBatch{[]map[string]interface{}{msg}}); err != nil {
			t.Fatalf("failed to publish message #%d: %s", i, err)
		}
		if i == len(messages)/2 {
			break
		}
	}
	// Batch
	if err = client.Publish(&mockBatch{messages[len(messages)/2+1:]}); err != nil {
		t.Fatalf("failed to publish batch: %s", err)
	}
	_, _, err = client.Refresh(index)
	if err != nil {
		t.Fatal(err)
	}
	_, result, err := client.Connection.SearchURI(index, "doc", map[string]string{
		"q": "path:\"/Sysdb/environment/archer/temperature/status/TempSensor1/temperature\"",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Hits.Total != 3 {
		t.Errorf("expected 3 result got %d", result.Hits.Total)
	}
	_, result, err = client.Connection.SearchURI(index, "doc", nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.Hits.Total != 10 {
		t.Errorf("expected 10 result got %d", result.Hits.Total)
	}
}
