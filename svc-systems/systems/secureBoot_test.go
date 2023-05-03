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

package systems

import (
	"context"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
)

func TestExternalInterface_UpdateSecureBoot(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	mockDeviceAndSystemData(t)
	pluginContact := PluginContact{
		ContactClient:  mockContactClient,
		DevicePassword: stubDevicePassword,
		UpdateTask:     mockUpdateTask,
	}

	type args struct {
		ctx    context.Context
		req    *systemsproto.SecureBootRequest
		pc     *PluginContact
		taskID string
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		args args
	}{
		{
			name: "invalid uuid",
			args: args{
				ctx: mockContext(),
				req: &systemsproto.SecureBootRequest{
					SystemID:    "24b243cf-f1e3-5318-92d9-2d6737d6b0b.1",
					RequestBody: []byte(`{"SecureBootEnable": true}`),
				},
				pc:     &pluginContact,
				taskID: "task24b243cf-f1e3-5318-92d9-2d6737d6b0b",
			},
			e: mockGetExternalInterface(),
		},
		{
			name: "invalid uuid without system id",
			args: args{
				ctx: mockContext(),
				req: &systemsproto.SecureBootRequest{
					SystemID:    "24b243cf-f1e3-5318-92d9-2d6737d6b0b",
					RequestBody: []byte(`{"SecureBootEnable": true}`),
				},
				pc:     &pluginContact,
				taskID: "task24b243cf-f1e3-5318-92d9-2d6737d6b0b",
			},
			e: mockGetExternalInterface(),
		},
		{
			name: "if plugin id is not there in the db",
			args: args{
				ctx: mockContext(),
				req: &systemsproto.SecureBootRequest{
					SystemID:    "7a2c6100-67da-5fd6-ab82-6870d29c727.1",
					RequestBody: []byte(`{"SecureBootEnable": true}`),
				},
				pc:     &pluginContact,
				taskID: "task24b243cf-f1e3-5318-92d9-2d6737d6b0b",
			},
			e: mockGetExternalInterface(),
		},
		{
			name: "invalid request body",
			args: args{
				ctx: mockContext(),
				req: &systemsproto.SecureBootRequest{
					SystemID:    "7a2c6100-67da-5fd6-ab82-6870d29c7279.1",
					RequestBody: []byte(`{"SecureBootMode": "UserMode"}`),
				},
				pc:     &pluginContact,
				taskID: "task24b243cf-f1e3-5318-92d9-2d6737d6b0b",
			},
			e: mockGetExternalInterface(),
		},
		{
			name: "bad request",
			args: args{
				ctx: mockContext(),
				req: &systemsproto.SecureBootRequest{
					SystemID:    "7a2c6100-67da-5fd6-ab82-6870d29c7279.1",
					RequestBody: []byte(`{"secureBootEnable": true}`),
				},
				pc:     &pluginContact,
				taskID: "task24b243cf-f1e3-5318-92d9-2d6737d6b0b",
			},
			e: mockGetExternalInterface(),
		},
		{
			name: "Valid Request",
			args: args{
				ctx: mockContext(),
				req: &systemsproto.SecureBootRequest{
					SystemID:    "7a2c6100-67da-5fd6-ab82-6870d29c7279.1",
					RequestBody: []byte(`{"SecureBootEnable": true}`),
				},
				pc:     &pluginContact,
				taskID: "task24b243cf-f1e3-5318-92d9-2d6737d6b0b",
			},
			e: mockGetExternalInterface(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.UpdateSecureBoot(tt.args.ctx, tt.args.req, tt.args.pc, tt.args.taskID)
		})
	}
}

func TestExternalInterface_ResetSecureBoot(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	mockDeviceAndSystemData(t)
	pluginContact := PluginContact{
		ContactClient:  mockContactClient,
		DevicePassword: stubDevicePassword,
		UpdateTask:     mockUpdateTask,
	}
	type args struct {
		ctx    context.Context
		req    *systemsproto.SecureBootRequest
		pc     *PluginContact
		taskID string
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		args args
	}{
		{
			name: "invalid uuid",
			args: args{
				ctx: mockContext(),
				req: &systemsproto.SecureBootRequest{
					SystemID:    "24b243cf-f1e3-5318-92d9-2d6737d6b0b.1",
					RequestBody: []byte(`{"ResetKeysType": "ResetAllKeysToDefault"}`),
				},
				pc:     &pluginContact,
				taskID: "task24b243cf-f1e3-5318-92d9-2d6737d6b0b",
			},
			e: mockGetExternalInterface(),
		},
		{
			name: "invalid uuid without system id",
			args: args{
				ctx: mockContext(),
				req: &systemsproto.SecureBootRequest{
					SystemID:    "24b243cf-f1e3-5318-92d9-2d6737d6b0b",
					RequestBody: []byte(`{"ResetKeysType": "ResetAllKeysToDefault"}`),
				},
				pc:     &pluginContact,
				taskID: "task24b243cf-f1e3-5318-92d9-2d6737d6b0b",
			},
			e: mockGetExternalInterface(),
		},
		{
			name: "if plugin id is not there in the db",
			args: args{
				ctx: mockContext(),
				req: &systemsproto.SecureBootRequest{
					SystemID:    "7a2c6100-67da-5fd6-ab82-6870d29c727.1",
					RequestBody: []byte(`{"ResetKeysType": "ResetAllKeysToDefault"}`),
				},
				pc:     &pluginContact,
				taskID: "task24b243cf-f1e3-5318-92d9-2d6737d6b0b",
			},
			e: mockGetExternalInterface(),
		},
		{
			name: "bad request",
			args: args{
				ctx: mockContext(),
				req: &systemsproto.SecureBootRequest{
					SystemID:    "7a2c6100-67da-5fd6-ab82-6870d29c7279.1",
					RequestBody: []byte(`{"resetKeysType": "ResetAllKeysToDefault"}`),
				},
				pc:     &pluginContact,
				taskID: "task24b243cf-f1e3-5318-92d9-2d6737d6b0b",
			},
			e: mockGetExternalInterface(),
		},
		{
			name: "Valid Request",
			args: args{
				ctx: mockContext(),
				req: &systemsproto.SecureBootRequest{
					SystemID:    "7a2c6100-67da-5fd6-ab82-6870d29c7279.1",
					RequestBody: []byte(`{"ResetKeysType": "ResetAllKeysToDefault"}`),
				},
				pc:     &pluginContact,
				taskID: "task24b243cf-f1e3-5318-92d9-2d6737d6b0b",
			},
			e: mockGetExternalInterface(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.ResetSecureBoot(tt.args.ctx, tt.args.req, tt.args.pc, tt.args.taskID)
		})
	}
}
