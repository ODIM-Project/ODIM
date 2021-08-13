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
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

func mockUpdateConnectionMethod(connectionMethod agmodel.ConnectionMethod, cmURI string) *errors.Error {
	return nil
}

func TestExternalInterface_AddBMC(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
	addComputeRetrieval := config.AddComputeSkipResources{
		SkipResourceListUnderSystem: []string{"Chassis", "LogServices"},
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
	mockPluginData(t, "GRF_v1.0.0")
	mockPluginData(t, "XAuthPlugin_v1.0.0")
	mockPluginData(t, "XAuthPluginFail_v1.0.0")

	reqSuccess, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	})
	reqWithoutConnectionMethod, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.11",
		UserName: "admin",
		Password: "password",
	})
	reqPluginID, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/2e99af48-2e99-4d78-a250-b04641e9b046",
			},
		},
	})
	reqSuccessXAuth, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.2",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/0a8992dc-8b47-4fe3-b26c-4c34048cf0d2",
			},
		},
	})
	reqIncorrectDeviceBasicAuth, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin1",
		Password: "incorrectPassword",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	})
	reqIncorrectDeviceXAuth, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.2",
		UserName: "username",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/7551386e-b9d7-4233-a963-3841adc69e17",
			},
		},
	})
	p := getMockExternalInterface()
	type args struct {
		taskID string
		req    *aggregatorproto.AggregatorRequest
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
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqSuccess,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusCreated,
			},
		},
		{
			name: "request without connectionmethod",
			p:    p,
			args: args{
				taskID: "123",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqWithoutConnectionMethod,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "update task failure or invalid taskID",
			p:    p,
			args: args{
				taskID: "invalid",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqSuccess,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "invalid request body format",
			p:    p,
			args: args{
				taskID: "123",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  []byte("some invalid format"),
				},
			},
			want: response.RPC{
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "invalid plugin id",
			p:    p,
			args: args{
				taskID: "123",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqPluginID,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "success: plugin with xauth token",
			p:    p,
			args: args{
				taskID: "123",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqSuccessXAuth,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusCreated,
			},
		},
		{
			name: "with incorrect device credentials and BasicAuth",
			p:    p,
			args: args{
				taskID: "123",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqIncorrectDeviceBasicAuth,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "with incorrect device credentials and XAuth",
			p:    p,
			args: args{
				taskID: "123",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqIncorrectDeviceXAuth,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusUnauthorized,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(2 * time.Second)
			if got := tt.p.AddAggregationSource(tt.args.taskID, "validUserName", tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.AddAggregationSource = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternalInterface_AddBMCForPasswordEncryptFail(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
	addComputeRetrieval := config.AddComputeSkipResources{
		SkipResourceListUnderSystem: []string{"Chassis", "LogServices"},
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
	mockPluginData(t, "GRF_v1.0.0")

	reqEncryptFail, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin",
		Password: "passwordWithInvalidEncryption",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	})
	p := getMockExternalInterface()
	type args struct {
		taskID string
		req    *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name string
		p    *ExternalInterface
		args args
		want response.RPC
	}{
		{
			name: "encryption failure",
			p:    p,
			args: args{
				taskID: "123",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqEncryptFail,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(2 * time.Second)
			if got := tt.p.AddAggregationSource(tt.args.taskID, "validUserName", tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.AddAggregationSource = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestExternalInterface_AddBMCDuplicate handles the test cases for getregistry and duplicate add server
func TestExternalInterface_AddBMCDuplicate(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
	addComputeRetrieval := config.AddComputeSkipResources{
		SkipResourceListUnderSystem: []string{"Chassis", "LogServices"},
	}
	config.Data.AddComputeSkipResources = &addComputeRetrieval
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockPluginData(t, "GRF_v1.0.0")

	reqSuccess, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	})
	p := getMockExternalInterface()
	p.ContactClient = mockContactClientForDuplicate
	type args struct {
		taskID string
		req    *aggregatorproto.AggregatorRequest
	}
	req := &aggregatorproto.AggregatorRequest{
		SessionToken: "validToken",
		RequestBody:  reqSuccess,
	}
	tests := []struct {
		name string
		p    *ExternalInterface
		args args
		want response.RPC
	}{
		{
			name: "success case with registries",
			want: response.RPC{
				StatusCode: http.StatusCreated,
			},
		},
		{
			name: "duplicate case",
			want: response.RPC{
				StatusCode: http.StatusConflict,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := p.AddAggregationSource("123", "validUserName", req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.AddAggregationSource = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternalInterface_Manager(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
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
	reqSuccess, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	})
	reqExistingPlugin, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/c41cbd97-937d-1b73-c41c-1b7385d39069",
			},
		},
	})
	reqInvalidAuthType, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/6f29f281-f5e2-4873-97b7-376be668f4f4",
			},
		},
	})
	reqInvalidPluginType, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/6456115a-e900-4c11-809f-0957031d2d56",
			},
		},
	})
	reqExistingPluginBadPassword, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/36474ba4-a201-46aa-badf-d8104da418e8",
			},
		},
	})
	reqExistingPluginBadData, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/4298f256-c279-44e2-94f2-3987bb7d8f53",
			},
		},
	})

	p := getMockExternalInterface()

	type args struct {
		taskID string
		req    *aggregatorproto.AggregatorRequest
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
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqSuccess,
				},
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
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqExistingPlugin,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusConflict,
			},
		}, {
			name: "Invalid Auth type",
			p:    p,
			args: args{
				taskID: "123",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqInvalidAuthType,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusBadRequest,
			},
		}, {
			name: "Invalid Plugin type",
			p:    p,
			args: args{
				taskID: "123",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqInvalidPluginType,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusBadRequest,
			},
		}, {
			name: "Existing Plugin with bad password",
			p:    p,
			args: args{
				taskID: "123",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqExistingPluginBadPassword,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusConflict,
			},
		}, {
			name: "Existing Plugin with bad data",
			p:    p,
			args: args{
				taskID: "123",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqExistingPluginBadData,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusConflict,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.AddAggregationSource(tt.args.taskID, "validUserName", tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.AddAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternalInterface_ManagerXAuth(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
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
	reqXAuthSuccess, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/058c1876-6f24-439a-8968-2af26154081f",
			},
		},
	})
	reqXAuthFail, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "incorrectusername",
		Password: "incorrectPassword",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/3489af48-2e99-4d78-a250-b04641e9d98d",
			},
		},
	})

	reqStatusFail, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.3:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/3489af48-2e99-4d78-a250-b04641e9d98d",
			},
		},
	})

	reqInvalidStatusBody, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.4:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/3489af48-2e99-4d78-a250-b04641e9d98d",
			},
		},
	})

	reqManagerGetFail, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.5:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/3489af48-2e99-4d78-a250-b04641e9d98d",
			},
		},
	})

	reqInvalidManagerBody, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.6:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/3489af48-2e99-4d78-a250-b04641e9d98d",
			},
		},
	})

	p := getMockExternalInterface()

	type args struct {
		taskID string
		req    *aggregatorproto.AggregatorRequest
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
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqXAuthSuccess,
				},
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
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqXAuthFail,
				},
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
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqStatusFail,
				},
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
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqInvalidStatusBody,
				},
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
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqManagerGetFail,
				},
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
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqInvalidManagerBody,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.AddAggregationSource(tt.args.taskID, "validUserName", tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.AddAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternalInterface_ManagerWithMultipleRequest(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
	addComputeRetrieval := config.AddComputeSkipResources{
		SkipResourceListUnderSystem: []string{"Chassis", "LogServices"},
	}
	config.Data.AddComputeSkipResources = &addComputeRetrieval
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()

	reqSuccess, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	})
	reqWithoutPort, _ := json.Marshal(AggregationSource{
		HostName: "localhost",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			ConnectionMethod: &ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	})

	p := getMockExternalInterface()

	type args struct {
		taskID string
		req    *aggregatorproto.AggregatorRequest
	}

	tests := []struct {
		name string
		p    *ExternalInterface
		args args
		want response.RPC
		req  *aggregatorproto.AggregatorRequest
	}{
		{
			name: "adding same BMC multiple times",
			want: response.RPC{
				StatusCode: http.StatusConflict,
			},
			req: &aggregatorproto.AggregatorRequest{
				SessionToken: "validToken",
				RequestBody:  reqSuccess,
			},
		},
		{
			name: "adding same BMC without port multiple times",
			want: response.RPC{
				StatusCode: http.StatusConflict,
			},
			req: &aggregatorproto.AggregatorRequest{
				SessionToken: "validToken",
				RequestBody:  reqWithoutPort,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go p.AddAggregationSource("123", "validUserName", tt.req)
			time.Sleep(time.Second)
			if got := p.AddAggregationSource("123", "validUserName", tt.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.AddAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

var activeReqFlag bool

func mockGenericSave(data []byte, table, key string) error {
	common.MuxLock.Lock()
	activeReqFlag = true
	common.MuxLock.Unlock()
	return nil
}

func mockCheckActiveRequest(managerAddress string) (bool, *errors.Error) {
	return activeReqFlag, nil
}

func mockDeleteActiveRequest(managerAddress string) *errors.Error {
	activeReqFlag = false
	return nil
}

func mockCheckMetricRequest(managerAddress string) (bool, *errors.Error) {
	return activeReqFlag, nil
}

func mockDeleteMetricRequest(managerAddress string) *errors.Error {
	activeReqFlag = false
	return nil
}

func mockGetAllMatchingDetails(table, pattern string, dbtype common.DbType) ([]string, *errors.Error) {
	return []string{
		"MetricReportDefinitions:/redfish/v1/TelemetryService/MetricReportDefinitions/CPUICUtilCustom1",
		"Triggers:/redfish/v1/TelemetryService/Triggers/CPU0PowerTriggers",
		"MetricReportDefinitions:/redfish/v1/TelemetryService/MetricReportDefinitions/CPUUtilCustom3",
	}, nil
}

func mockDelete(table, key string, dbtype common.DbType) *errors.Error {
	return nil
}

func getMockExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		ContactClient:           mockContactClient,
		Auth:                    mockIsAuthorized,
		CreateChildTask:         mockCreateChildTask,
		UpdateTask:              mockUpdateTask,
		CreateSubcription:       EventFunctionsForTesting,
		PublishEvent:            PostEventFunctionForTesting,
		GetPluginStatus:         GetPluginStatusForTesting,
		SubscribeToEMB:          mockSubscribeEMB,
		EncryptPassword:         stubDevicePassword,
		DecryptPassword:         stubDevicePassword,
		GetConnectionMethod:     mockGetConnectionMethod,
		UpdateConnectionMethod:  mockUpdateConnectionMethod,
		GetAllKeysFromTable:     mockGetAllKeysFromTable,
		GetPluginMgrAddr:        stubPluginMgrAddrData,
		GenericSave:             mockGenericSave,
		CheckActiveRequest:      mockCheckActiveRequest,
		DeleteActiveRequest:     mockDeleteActiveRequest,
		DeleteComputeSystem:     deleteComputeforTest,
		DeleteSystem:            deleteSystemforTest,
		DeleteEventSubscription: mockDeleteSubscription,
		EventNotification:       mockEventNotification,
		GetAllMatchingDetails:   mockGetAllMatchingDetails,
		CheckMetricRequest:      mockCheckMetricRequest,
		DeleteMetricRequest:     mockDeleteMetricRequest,
		GetResource:             mockGetResource,
		Delete:                  mockDelete,
	}
}
