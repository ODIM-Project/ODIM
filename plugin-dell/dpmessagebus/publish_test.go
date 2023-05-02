//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

// Package dpmessagebus ...
package dpmessagebus

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/plugin-dell/config"
)

func TestPublish(t *testing.T) {
	config.SetUpMockConfig(t)
	event := common.MessageData{
		Events: []common.Event{
			common.Event{
				EventType: "Alert",
				MemberID:  "Message1.2",
			},
		},
	}
	eventData, _ := json.Marshal(event)
	type args struct {
		ctx  context.Context
		data interface{}
	}
	tests := []struct {
		name       string
		args       args
		want       bool
		mockConfig func()
	}{
		{
			name: "Positive test ",
			args: args{
				ctx: context.TODO(),
				data: common.Events{
					IP:        "10.10.10.10",
					EventType: "",
					Request:   []byte(`{"event":""}`),
				},
			},
			mockConfig: func() {},
			want:       true,
		},
		{
			name: "Negative test - nil ",
			args: args{
				ctx:  context.TODO(),
				data: nil,
			},
			want:       false,
			mockConfig: func() {},
		},
		{
			name: "Negative test - invalid request ",
			args: args{
				ctx: context.TODO(),
				data: common.Events{
					IP:        "10.10.10.10",
					EventType: "",
					Request:   []byte(`invalid`),
				},
			},
			want:       false,
			mockConfig: func() {},
		},
		{
			name: "Negative test - invalid request ",
			args: args{
				ctx: context.TODO(),
				data: common.Events{
					IP:        "10.10.10.10",
					EventType: "",
					Request:   []byte(`invalid`),
				},
			},
			want:       false,
			mockConfig: func() {},
		},
		{
			name: "Negative test - invalid request ",
			args: args{
				ctx: context.TODO(),
				data: common.Events{
					IP:        "10.10.10.10",
					EventType: "",
					Request:   eventData,
				},
			},
			want: true,
			mockConfig: func() {
				config.SetUpMockConfig(t)
				config.Data.MessageBusConf.MessageBusConfigFilePath = ""
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockConfig()
			if got := Publish(tt.args.ctx, tt.args.data); got != tt.want {
				t.Errorf("Publish() = %v, want %v", got, tt.want)
			}
		})
	}
}
