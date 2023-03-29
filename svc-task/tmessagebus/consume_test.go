// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package tmessagebus

import (
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
)

func TestSubscribeTaskEventsQueue(t *testing.T) {
	err := config.SetUpMockConfig(t)
	if err != nil {
		t.Fatalf("fatal: error while trying to collect mock db config: %v", err)
		return
	}
	type args struct {
		topicName string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "subscribe task event topic",
			args: args{
				topicName: "task-event-topic",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SubscribeTaskEventsQueue(tt.args.topicName)
		})
	}
}

func Test_consumeTaskEvents(t *testing.T) {
	err := config.SetUpMockConfig(t)
	if err != nil {
		t.Fatalf("fatal: error while trying to collect mock db config: %v", err)
		return
	}
	TaskEventRecvQueue, TaskEventProcQueue = common.CreateJobQueue(1)
	type args struct {
		event interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "consume task event",
			args: args{
				event: common.TaskEvent{
					TaskID:          "sample_task_id",
					TaskState:       common.New,
					TaskStatus:      common.OK,
					PercentComplete: 30,
					EndTime:         time.Now(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumeTaskEvents(tt.args.event)
		})
	}
}
