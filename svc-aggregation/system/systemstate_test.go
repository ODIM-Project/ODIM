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

package system

import (
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

func TestExternalInterface_UpdateSystemState(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	device1 := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "24b243cf-f1e3-5318-92d9-2d6737d6b0b9",
		PluginID:       "GRF",
	}
	device2 := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("passwordWithInvalidEncryption"),
		UserName:       "admin",
		DeviceUUID:     "24b243cf-f1e3-5318-92d9-2d6737d6b0b9",
		PluginID:       "GRF",
	}
	device3 := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "24b243cf-f1e3-5318-92d9-2d6737d6b0b9",
		PluginID:       "unknown-plugin",
	}
	device4 := agmodel.Target{
		ManagerAddress: "100.0.0.7",
		Password:       []byte("some-password"),
		UserName:       "admin",
		DeviceUUID:     "something",
		PluginID:       "XAuthPlugin",
	}
	device5 := agmodel.Target{
		ManagerAddress: "100.0.0.8",
		Password:       []byte("somepassword"),
		UserName:       "admin",
		DeviceUUID:     "something",
		PluginID:       "XAuthPluginFail",
	}
	mockPluginData(t, "GRF")
	mockPluginData(t, "XAuthPlugin")
	mockPluginData(t, "XAuthPluginFail")
	mockDeviceData("24b243cf-f1e3-5318-92d9-2d6737d6b0b9", device1)
	mockDeviceData("device-with-wrong-password-encrypt", device2)
	mockDeviceData("unknown-plugin", device3)
	mockDeviceData("xauth-plugin", device4)
	mockDeviceData("xauth-plugin-fail", device5)
	mockDeviceData("no-system-data", device1)
	mockSystemData("/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b9.1")
	mockSystemData("/redfish/v1/Systems/xauth-plugin.1")
	pluginContact := ExternalInterface{
		ContactClient:   mockContactClient,
		Auth:            mockIsAuthorized,
		CreateChildTask: mockCreateChildTask,
		UpdateTask:      mockUpdateTask,
		DecryptPassword: stubDevicePassword,
		GetPluginStatus: GetPluginStatusForTesting,
	}
	type args struct {
		updateReq *aggregatorproto.UpdateSystemStateRequest
	}
	tests := []struct {
		name    string
		e       *ExternalInterface
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			e:    &pluginContact,
			args: args{
				updateReq: &aggregatorproto.UpdateSystemStateRequest{
					SystemURI:  "/redfish/v1/Systems/",
					SystemUUID: "24b243cf-f1e3-5318-92d9-2d6737d6b0b9",
					SystemID:   "1",
				},
			},
			wantErr: false,
		},
		{
			name: "device not found",
			e:    &pluginContact,
			args: args{
				updateReq: &aggregatorproto.UpdateSystemStateRequest{
					SystemURI:  "/redfish/v1/Systems/",
					SystemUUID: "no-device-uuid",
					SystemID:   "1",
				},
			},
			wantErr: true,
		},
		{
			name: "device with wrongly encrypted password",
			e:    &pluginContact,
			args: args{
				updateReq: &aggregatorproto.UpdateSystemStateRequest{
					SystemURI:  "/redfish/v1/Systems/",
					SystemUUID: "device-with-wrong-password-encrypt",
					SystemID:   "1",
				},
			},
			wantErr: true,
		},
		{
			name: "unknown plugin",
			e:    &pluginContact,
			args: args{
				updateReq: &aggregatorproto.UpdateSystemStateRequest{
					SystemURI:  "/redfish/v1/Systems/",
					SystemUUID: "unknown-plugin",
					SystemID:   "1",
				},
			},
			wantErr: true,
		},
		{
			name: "xauth token",
			e:    &pluginContact,
			args: args{
				updateReq: &aggregatorproto.UpdateSystemStateRequest{
					SystemURI:  "/redfish/v1/Systems/",
					SystemUUID: "xauth-plugin",
					SystemID:   "1",
				},
			},
			wantErr: false,
		},
		{
			name: "positive case",
			e:    &pluginContact,
			args: args{
				updateReq: &aggregatorproto.UpdateSystemStateRequest{
					SystemURI:  "/redfish/v1/Systems/",
					SystemUUID: "xauth-plugin-fail",
					SystemID:   "1",
				},
			},
			wantErr: true,
		},
		{
			name: "no system data",
			e:    &pluginContact,
			args: args{
				updateReq: &aggregatorproto.UpdateSystemStateRequest{
					SystemURI:  "/redfish/v1/Systems/",
					SystemUUID: "no-system-data",
					SystemID:   "1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UpdateSystemState(tt.args.updateReq); (err != nil) != tt.wantErr {
				t.Errorf("ExternalInterface.UpdateSystemState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
