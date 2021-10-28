//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

// Package evmodel have the struct models and DB functionalties
package consumer

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
)

func TestConsume(t *testing.T) {

	tests := []struct {
		name      string
		topicName string
	}{
		{
			name:      "posivite case",
			topicName: "topic",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Consume(tt.topicName)
		})
	}
}

func TestKafkaSubscriber(t *testing.T) {
	In, Out = common.CreateJobQueue(1)
	eventMessage := common.MessageData{
		Name:    "Event",
		Context: "context",
		Events: []common.Event{
			common.Event{
				EventType: "Alert",
				MessageID: "AlertEvent",
				OriginOfCondition: &common.Link{
					Oid: "/redfish/v1/Systems/1",
				},
			},
		},
	}
	event, _ := json.Marshal(eventMessage)
	kafkaMessage := common.Events{
		IP:      "10.1.2.3",
		Request: event,
	}
	KafkaSubscriber(kafkaMessage)

	var currentData int

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for range Out {
			currentData++
		}
		wg.Done()
	}()
	time.Sleep(2 * time.Second)
	close(In)
	wg.Wait()

	if currentData != 1 {
		t.Errorf("error: expected count is 1 but got %v", currentData)
	}
}
