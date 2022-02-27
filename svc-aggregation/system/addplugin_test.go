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
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
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

func stubPluginMgrAddrData(pluginID string) (agmodel.Plugin, *errors.Error) {
	var plugin agmodel.Plugin

	plugin, err := agmodel.GetPluginData(pluginID)
	if err != nil {
		plugin.ID = pluginID
		plugin.ManagerUUID = "dummy-mgr-addr"
		plugin.Port = "9091"
	}
	plugin.IP = "dummyhost"

	if pluginID == "DUPMGRADDRMOCK" {
		plugin.ManagerUUID = "duplicate-mgr-addr"
		plugin.IP = "duphost"
		plugin.Port = "9091"
	}

	return plugin, nil

}

func mockDupMgrAddrPluginData(t *testing.T, pluginID string) error {
	password, _ := stubDevicePassword([]byte("password"))
	plugin := agmodel.Plugin{
		IP:                "duphost",
		Port:              "9091",
		Username:          "admin",
		Password:          password,
		ID:                pluginID,
		PreferredAuthType: "BasicAuth",
		ManagerUUID:       "duplicate-mgr-addr",
	}
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("Plugin", pluginID, plugin); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "Plugin", err.Error())
	}
	return nil
}

func TestExternalInterface_Plugin(t *testing.T) {
	config.SetUpMockConfig(t)
	addComputeRetrieval := config.AddComputeSkipResources{
		SkipResourceListUnderSystem: []string{"Chassis", "LogServices"},
	}
	err := mockPluginData(t, "ILO_v1.0.0")
	if err != nil {
		t.Fatalf("Error in creating mock PluginData :%v", err)
	}

	// create plugin with bad password for decryption failure
	pluginData := agmodel.Plugin{
		Password: []byte("password"),
		ID:       "PluginWithBadPassword",
	}
	mockData(t, common.OnDisk, "Plugin", "PluginWithBadPassword_v1.0.0", pluginData)
	// create plugin with bad data
	mockData(t, common.OnDisk, "Plugin", "PluginWithBadData_v1.0.0", "PluginWithBadData")

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
		ConnectionMethod: &ConnectionMethod{
			OdataID: "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
		},
	}
	reqExistingPlugin := AddResourceRequest{
		ManagerAddress: "localhost:9091",
		UserName:       "admin",
		Password:       "password",
		ConnectionMethod: &ConnectionMethod{
			OdataID: "/redfish/v1/AggregationService/ConnectionMethods/c41cbd97-937d-1b73-c41c-1b7385d39069",
		},
	}
	reqInvalidAuthType := AddResourceRequest{
		ManagerAddress: "localhost:9091",
		UserName:       "admin",
		Password:       "password",
		ConnectionMethod: &ConnectionMethod{
			OdataID: "/redfish/v1/AggregationService/ConnectionMethods/6f29f281-f5e2-4873-97b7-376be668f4f4",
		},
	}
	reqInvalidPluginType := AddResourceRequest{
		ManagerAddress: "localhost:9091",
		UserName:       "admin",
		Password:       "password",
		ConnectionMethod: &ConnectionMethod{
			OdataID: "/redfish/v1/AggregationService/ConnectionMethods/6456115a-e900-4c11-809f-0957031d2d56",
		},
	}
	reqExistingPluginBadPassword := AddResourceRequest{
		ManagerAddress: "localhost:9091",
		UserName:       "admin",
		Password:       "password",
		ConnectionMethod: &ConnectionMethod{
			OdataID: "/redfish/v1/AggregationService/ConnectionMethods/36474ba4-a201-46aa-badf-d8104da418e8",
		},
	}
	reqExistingPluginBadData := AddResourceRequest{
		ManagerAddress: "localhost:9091",
		UserName:       "admin",
		Password:       "password",
		ConnectionMethod: &ConnectionMethod{
			OdataID: "/redfish/v1/AggregationService/ConnectionMethods/4298f256-c279-44e2-94f2-3987bb7d8f53",
		},
	}
	reqPluginWithDuplciateUUID := AddResourceRequest{
		ManagerAddress: "localhost:9091",
		UserName:       "admin",
		Password:       "password",
		ConnectionMethod: &ConnectionMethod{
			OdataID: "/redfish/v1/AggregationService/ConnectionMethods/7e821c50-ecc7-4f45-ba15-2a31b389df1e",
		},
	}

	p := getMockExternalInterface()
	targetURI := "/redfish/v1/AggregationService/AggregationSource"
	var queueList []string
	var pluginContactRequest getResourceRequest
	pluginContactRequest.ContactClient = p.ContactClient
	pluginContactRequest.GetPluginStatus = p.GetPluginStatus
	pluginContactRequest.TargetURI = targetURI
	pluginContactRequest.UpdateTask = p.UpdateTask
	type args struct {
		taskID     string
		req        AddResourceRequest
		cmVariants connectionMethodVariants
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
				taskID:     "123",
				req:        reqSuccess,
				cmVariants: getConnectionMethodVariants("Compute:BasicAuth:GRF_v1.0.0"),
			},
			want: response.RPC{
				StatusCode: http.StatusCreated,
			},
		},
		{
			name: "Existing Plugin",
			p:    p,
			args: args{
				taskID:     "123",
				req:        reqExistingPlugin,
				cmVariants: getConnectionMethodVariants("Compute:BasicAuth:ILO_v1.0.0"),
			},
			want: response.RPC{
				StatusCode: http.StatusConflict,
			},
		},
		{
			name: "Invalid Auth type",
			p:    p,
			args: args{
				taskID:     "123",
				req:        reqInvalidAuthType,
				cmVariants: getConnectionMethodVariants("Compute:BasicAuthentication:ILO_v1.0.0"),
			},
			want: response.RPC{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Invalid Plugin type",
			p:    p,
			args: args{
				taskID:     "123",
				req:        reqInvalidPluginType,
				cmVariants: getConnectionMethodVariants("plugin:BasicAuth:ILO_v1.0.0"),
			},
			want: response.RPC{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Existing Plugin with bad password",
			p:    p,
			args: args{
				taskID:     "123",
				req:        reqExistingPluginBadPassword,
				cmVariants: getConnectionMethodVariants("Compute:BasicAuth:PluginWithBadPassword_v1.0.0"),
			},
			want: response.RPC{
				StatusCode: http.StatusConflict,
			},
		},
		{
			name: "Existing Plugin with bad data",
			p:    p,
			args: args{
				taskID:     "123",
				req:        reqExistingPluginBadData,
				cmVariants: getConnectionMethodVariants("Compute:BasicAuth:PluginWithBadData_v1.0.0"),
			},

			want: response.RPC{
				StatusCode: http.StatusConflict,
			},
		},

		{
			name: "Adding plugin with duplicate managre uuid",
			p:    p,
			args: args{
				taskID:     "123",
				req:        reqPluginWithDuplciateUUID,
				cmVariants: getConnectionMethodVariants("Compute:BasicAuth:STGtest_v1.0.0"),
			},
			want: response.RPC{
				StatusCode: http.StatusConflict,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _, _ := tt.p.addPluginData(tt.args.req, tt.args.taskID, targetURI, pluginContactRequest, queueList, tt.args.cmVariants); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.addPluginData = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternalInterface_PluginXAuth(t *testing.T) {
	config.SetUpMockConfig(t)
	addComputeRetrieval := config.AddComputeSkipResources{
		SkipResourceListUnderSystem: []string{"Chassis", "LogServices"},
	}
	err := mockPluginData(t, "XAuthPlugin_v1.0.0")
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
		ManagerAddress: "localhost:9091",
		UserName:       "admin",
		Password:       "password",
		ConnectionMethod: &ConnectionMethod{
			OdataID: "/redfish/v1/AggregationService/ConnectionMethods/058c1876-6f24-439a-8968-2af26154081f",
		},
	}
	reqXAuthFail := AddResourceRequest{
		ManagerAddress: "localhost:9091",
		UserName:       "incorrectusername",
		Password:       "incorrectPassword",
		ConnectionMethod: &ConnectionMethod{
			OdataID: "/redfish/v1/AggregationService/ConnectionMethods/3489af48-2e99-4d78-a250-b04641e9d98d",
		},
	}

	reqManagerGetFail := AddResourceRequest{
		ManagerAddress: "100.0.0.5:9091",
		UserName:       "admin",
		Password:       "password",
		ConnectionMethod: &ConnectionMethod{
			OdataID: "/redfish/v1/AggregationService/ConnectionMethods/3489af48-2e99-4d78-a250-b04641e9d98d",
		},
	}

	reqInvalidManagerBody := AddResourceRequest{
		ManagerAddress: "100.0.0.6:9091",
		UserName:       "admin",
		Password:       "password",
		ConnectionMethod: &ConnectionMethod{
			OdataID: "/redfish/v1/AggregationService/ConnectionMethods/3489af48-2e99-4d78-a250-b04641e9d98d",
		},
	}

	p := getMockExternalInterface()
	var queueList []string
	targetURI := "/redfish/v1/AggregationService/AggregationSource"
	var pluginContactRequest getResourceRequest
	pluginContactRequest.ContactClient = p.ContactClient
	pluginContactRequest.GetPluginStatus = p.GetPluginStatus
	pluginContactRequest.TargetURI = targetURI
	pluginContactRequest.UpdateTask = p.UpdateTask
	type args struct {
		taskID     string
		req        AddResourceRequest
		cmVariants connectionMethodVariants
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
				taskID:     "123",
				req:        reqXAuthSuccess,
				cmVariants: getConnectionMethodVariants("Compute:XAuthToken:GRF_v1.0.0"),
			},

			want: response.RPC{
				StatusCode: http.StatusCreated,
			},
		},
		{
			name: "Failure with XAuthToken",
			p:    p,
			args: args{
				taskID:     "123",
				req:        reqXAuthFail,
				cmVariants: getConnectionMethodVariants("Compute:XAuthToken:ILO_v1.0.0"),
			},

			want: response.RPC{
				StatusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "Failure with Manager Get",
			p:    p,
			args: args{
				taskID:     "123",
				req:        reqManagerGetFail,
				cmVariants: getConnectionMethodVariants("Compute:XAuthToken:ILO_v1.0.0"),
			},

			want: response.RPC{
				StatusCode: http.StatusServiceUnavailable,
			},
		},
		{
			name: "incorrect manager body",
			p:    p,
			args: args{
				taskID:     "123",
				req:        reqInvalidManagerBody,
				cmVariants: getConnectionMethodVariants("Compute:XAuthToken:ILO_v1.0.0"),
			},
			want: response.RPC{
				StatusCode: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _, _ := tt.p.addPluginData(tt.args.req, tt.args.taskID, targetURI, pluginContactRequest, queueList, tt.args.cmVariants); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.addPluginData = %v, want %v", got, tt.want)
			}
		})
	}
}
