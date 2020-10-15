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
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

func mockData(t *testing.T, dbType common.DbType, table, id string, data interface{}) {
	connPool, err := common.GetDBConnection(dbType)
	if err != nil {
		t.Fatalf("error: mockData() failed to DB connection: %v", err)
	}
	if err = connPool.Create(table, id, data); err != nil {
		t.Fatalf("error: mockData() failed to create entry %s-%s: %v", table, id, err)
	}
}

func stubPluginData (pluginID string) (agmodel.Plugin, *errors.Error) {
        var plugin agmodel.Plugin

        plugin.IP = "localhost"
        plugin.Port = "9091"
        plugin.Username = "admin"
        plugin.Password = []byte("password")
        plugin.ID = "XAuthPlugin"
        plugin.PluginType = "Compute"
        plugin.PreferredAuthType = "BasicAuth"
        plugin.ManagerUUID = "1s7sda8asd-asdas8as0"
        return  plugin, nil
}

func TestExternalInterface_Plugin(t *testing.T) {
	config.SetUpMockConfig(t)
	addComputeRetrieval := config.AddComputeSkipResources{
		SystemCollection: []string{"Chassis", "LogServices"},
	}
	err := mockPluginData(t, "ILO")
	if err != nil {
		t.Fatalf("Error in creating mock PluginData :%v", err)
	}

	// create plugin with bad password for decryption failure
	pluginData := agmodel.Plugin{
		Password: []byte("password"),
		ID:       "PluginWithBadPassword",
	}
	mockData(t, common.OnDisk, "Plugin", "PluginWithBadPassword", pluginData)
	// create plugin with bad data
	mockData(t, common.OnDisk, "Plugin", "PluginWithBadData", "PluginWithBadData")

	config.Data.AddComputeSkipResources = &addComputeRetrieval
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
	reqSuccess := AddResourceRequest{
		ManagerAddress: "localhost:9091",
		UserName:       "admin",
		Password:       "password",
		Oem: &AddOEM{
			PluginID:          "GRF",
			PreferredAuthType: "BasicAuth",
			PluginType:        "Compute",
		},
	}
	reqExistingPlugin := AddResourceRequest{
		ManagerAddress: "localhost:9091",
		UserName:       "admin",
		Password:       "password",

		Oem: &AddOEM{
			PluginID:          "ILO",
			PreferredAuthType: "BasicAuth",
			PluginType:        "Compute",
		},
	}
	reqInvalidAuthType := AddResourceRequest{
		ManagerAddress: "localhost:9091",
		UserName:       "admin",
		Password:       "password",

		Oem: &AddOEM{
			PluginID:          "ILO",
			PreferredAuthType: "BasicAuthentication",
			PluginType:        "Compute",
		},
	}
	reqInvalidPluginType := AddResourceRequest{
		ManagerAddress: "localhost:9091",
		UserName:       "admin",
		Password:       "password",

		Oem: &AddOEM{
			PluginID:          "ILO",
			PreferredAuthType: "BasicAuth",
			PluginType:        "plugin",
		},
	}
	reqExistingPluginBadPassword := AddResourceRequest{
		ManagerAddress: "localhost:9091",
		UserName:       "admin",
		Password:       "password",

		Oem: &AddOEM{
			PluginID:          "PluginWithBadPassword",
			PreferredAuthType: "BasicAuth",
			PluginType:        "Compute",
		},
	}
	reqExistingPluginBadData := AddResourceRequest{
		ManagerAddress: "localhost:9091",
		UserName:       "admin",
		Password:       "password",

		Oem: &AddOEM{
			PluginID:          "PluginWithBadData",
			PreferredAuthType: "BasicAuth",
			PluginType:        "Compute",
		},
	}

	p := &ExternalInterface{
		ContactClient:     mockContactClient,
		Auth:              mockIsAuthorized,
		CreateChildTask:   mockCreateChildTask,
		UpdateTask:        mockUpdateTask,
		CreateSubcription: EventFunctionsForTesting,
		PublishEvent:      PostEventFunctionForTesting,
		GetPluginStatus:   GetPluginStatusForTesting,
		SubscribeToEMB:    mockSubscribeEMB,
		EncryptPassword:   stubDevicePassword,
		DecryptPassword:   stubDevicePassword,
	}
	targetURI := "/redfish/v1/AggregationService/AggregationSource"
	var pluginContactRequest getResourceRequest
	pluginContactRequest.ContactClient = p.ContactClient
	pluginContactRequest.GetPluginStatus = p.GetPluginStatus
	pluginContactRequest.TargetURI = targetURI
	pluginContactRequest.UpdateTask = p.UpdateTask
	type args struct {
		taskID string
		req    AddResourceRequest
	}
	tests := []struct {
		name string
		p    *ExternalInterface
		args args
		want response.RPC
	}{
		{
			name: "posivite case",
			p:    p,
			args: args{
				taskID: "123",
				req:    reqSuccess,
			},
			want: response.RPC{
				StatusCode: http.StatusCreated,
			},
		},
		{
			name: "Existing Plugin",
			p:    p,
			args: args{
				taskID: "123",
				req:    reqExistingPlugin,
			},
			want: response.RPC{
				StatusCode: http.StatusConflict,
			},
		},
		{
			name: "Invalid Auth type",
			p:    p,
			args: args{
				taskID: "123",
				req:    reqInvalidAuthType,
			},
			want: response.RPC{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Invalid Plugin type",
			p:    p,
			args: args{
				taskID: "123",
				req:    reqInvalidPluginType,
			},
			want: response.RPC{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Existing Plugin with bad password",
			p:    p,
			args: args{
				taskID: "123",
				req:    reqExistingPluginBadPassword,
			},
			want: response.RPC{
				StatusCode: http.StatusConflict,
			},
		},
		{
			name: "Existing Plugin with bad data",
			p:    p,
			args: args{
				taskID: "123",
				req:    reqExistingPluginBadData,
			},

			want: response.RPC{
				StatusCode: http.StatusConflict,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _, _ := tt.p.addPluginData(tt.args.req, tt.args.taskID, targetURI, pluginContactRequest); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.addPluginData = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternalInterface_PluginXAuth(t *testing.T) {
	config.SetUpMockConfig(t)
	addComputeRetrieval := config.AddComputeSkipResources{
		SystemCollection: []string{"Chassis", "LogServices"},
	}
	err := mockPluginData(t, "XAuthPlugin")
	if err != nil {
		t.Fatalf("Error in creating mock PluginData :%v", err)
	}

	config.Data.AddComputeSkipResources = &addComputeRetrieval
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

	if err != nil {
		t.Fatalf("error while trying to create schema: %v", err)
	}
	reqXAuthSuccess := AddResourceRequest{
		ManagerAddress: "100.0.0.7:9091",
		UserName:       "admin",
		Password:       "password",

		Oem: &AddOEM{
			PluginID:          "GRF",
			PreferredAuthType: "XAuthToken",
			PluginType:        "Compute",
		},
	}
	reqXAuthFail := AddResourceRequest{
		ManagerAddress: "100.0.0.8:9091",
		UserName:       "incorrectusername",
		Password:       "incorrectPassword",

		Oem: &AddOEM{
			PluginID:          "ILO",
			PreferredAuthType: "XAuthToken",
			PluginType:        "Compute",
		},
	}

	reqStatusFail := AddResourceRequest{
		ManagerAddress: "100.0.0.3:9091",
		UserName:       "admin",
		Password:       "password",

		Oem: &AddOEM{
			PluginID:          "ILO",
			PreferredAuthType: "XAuthToken",
			PluginType:        "Compute",
		},
	}

	reqInvalidStatusBody := AddResourceRequest{
		ManagerAddress: "100.0.0.4:9091",
		UserName:       "admin",
		Password:       "password",

		Oem: &AddOEM{
			PluginID:          "ILO",
			PreferredAuthType: "XAuthToken",
			PluginType:        "Compute",
		},
	}

	reqManagerGetFail := AddResourceRequest{
		ManagerAddress: "100.0.0.5:9091",
		UserName:       "admin",
		Password:       "password",

		Oem: &AddOEM{
			PluginID:          "ILO",
			PreferredAuthType: "XAuthToken",
			PluginType:        "Compute",
		},
	}

	reqInvalidManagerBody := AddResourceRequest{
		ManagerAddress: "100.0.0.6:9091",
		UserName:       "admin",
		Password:       "password",

		Oem: &AddOEM{
			PluginID:          "ILO",
			PreferredAuthType: "XAuthToken",
			PluginType:        "Compute",
		},
	}

        reqDuplicateManagerAddress:= AddResourceRequest{
                ManagerAddress: "localhost:9091",
                UserName: "admin",
                Password: "password",

                        Oem: &AddOEM{
                                PluginID:          "ILO",
                                PreferredAuthType: "XAuthToken",
                                PluginType:        "Compute",
                        },

        }

	p := &ExternalInterface{
		ContactClient:     mockContactClient,
		Auth:              mockIsAuthorized,
		CreateChildTask:   mockCreateChildTask,
		UpdateTask:        mockUpdateTask,
		CreateSubcription: EventFunctionsForTesting,
		PublishEvent:      PostEventFunctionForTesting,
		GetPluginStatus:   GetPluginStatusForTesting,
		SubscribeToEMB:    mockSubscribeEMB,
		EncryptPassword:   stubDevicePassword,
		DecryptPassword:   stubDevicePassword,
		GetPluginData:     stubPluginData,
	}
	targetURI := "/redfish/v1/AggregationService/AggregationSource"
	var pluginContactRequest getResourceRequest
	pluginContactRequest.ContactClient = p.ContactClient
	pluginContactRequest.GetPluginStatus = p.GetPluginStatus
	pluginContactRequest.TargetURI = targetURI
	pluginContactRequest.UpdateTask = p.UpdateTask
	type args struct {
		taskID string
		req    AddResourceRequest
	}
	tests := []struct {
		name string
		p    *ExternalInterface
		args args
		want response.RPC
	}{
		{
			name: "posivite case with XAuthToken",
			p:    p,
			args: args{
				taskID: "123",
				req:    reqXAuthSuccess,
			},

			want: response.RPC{
				StatusCode: http.StatusCreated,
			},
		},
		{
			name: "Failure with XAuthToken",
			p:    p,
			args: args{
				taskID: "123",
				req:    reqXAuthFail,
			},

			want: response.RPC{
				StatusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "Failure with Status Check",
			p:    p,
			args: args{
				taskID: "123",
				req:    reqStatusFail,
			},

			want: response.RPC{
				StatusCode: http.StatusServiceUnavailable,
			},
		},
		{
			name: "incorrect status body",
			p:    p,
			args: args{
				taskID: "123",
				req:    reqInvalidStatusBody,
			},

			want: response.RPC{
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "Failure with Manager Get",
			p:    p,
			args: args{
				taskID: "123",
				req:    reqManagerGetFail,
			},

			want: response.RPC{
				StatusCode: http.StatusServiceUnavailable,
			},
		},
		{
			name: "incorrect manager body",
			p:    p,
			args: args{
				taskID: "123",
				req:    reqInvalidManagerBody,
			},
			want: response.RPC{
				StatusCode: http.StatusInternalServerError,
			},
		},
                {
                        name: "duplicate manager address",
                        p:    p,
                        args: args{
                                taskID: "123",
                                req: reqDuplicateManagerAddress,
                        },
                        want: response.RPC{
                                StatusCode: http.StatusConflict,
                        },
                },

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _, _ := tt.p.addPluginData(tt.args.req, tt.args.taskID, targetURI, pluginContactRequest); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.addPluginData = %v, want %v", got, tt.want)
			}
		})
	}
}
