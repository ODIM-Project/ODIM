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
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
)

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
	message := common.Events{
		IP:      "10.1.2.3",
		Request: event,
	}
	EventSubscriber(message)

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
	EventSubscriber("invalidJson")
}

func TestConsume(t *testing.T) {
	config.SetUpMockConfig(t)
	type args struct {
		topicName string
	}
	tests := []struct {
		name      string
		args      args
		setConfig func()
	}{
		{
			name: "Positive Test",
			args: args{
				topicName: "demo",
			},
			setConfig: func() {},
		},
		{
			name: "Negative Invalid Message Type",
			args: args{
				topicName: "demo",
			},
			setConfig: func() {
				config.Data.MessageBusConf.MessageBusType = "Invalid"
			},
		},
	}
	for _, tt := range tests {
		tt.setConfig()
		t.Run(tt.name, func(t *testing.T) {
			Consume(mockContext(), tt.args.topicName)
		})
	}
}

func TestSubscribeCtrlMsgQueue(t *testing.T) {
	config.SetUpMockConfig(t)
	type args struct {
		topicName string
	}
	tests := []struct {
		name      string
		args      args
		setConfig func()
	}{
		{
			name: "Positive Test",
			args: args{
				topicName: "demo",
			},
			setConfig: func() {
			},
		},
		{
			name: "Negative Invalid Message Type",
			args: args{
				topicName: "demo",
			},
			setConfig: func() {
				config.Data.MessageBusConf.MessageBusType = "Invalid"
			},
		},
	}
	for _, tt := range tests {
		tt.setConfig()
		t.Run(tt.name, func(t *testing.T) {
			SubscribeCtrlMsgQueue(tt.args.topicName)
		})
	}
}

func Test_consumeCtrlMsg(t *testing.T) {
	config.SetUpMockConfig(t)
	type args struct {
		event interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Empty Test",
			args: args{
				event: "",
			},
		},
		{
			name: "Positive ",
			args: args{
				event: common.Events{
					IP:        "00.00.00.00",
					Request:   []byte(``),
					EventType: "Alert",
				},
			},
		},
	}
	for _, tt := range tests {
		CtrlMsgRecvQueue, CtrlMsgProcQueue = common.CreateJobQueue(4)
		In, Out = common.CreateJobQueue(1)
		t.Run(tt.name, func(t *testing.T) {
			consumeCtrlMsg(tt.args.event)
		})
	}

}
func mockContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.TransactionID, "xyz")
	ctx = context.WithValue(ctx, common.ActionID, "001")
	ctx = context.WithValue(ctx, common.ActionName, "xyz")
	ctx = context.WithValue(ctx, common.ThreadID, "0")
	ctx = context.WithValue(ctx, common.ThreadName, "xyz")
	ctx = context.WithValue(ctx, common.ProcessName, "xyz")
	return ctx
}
