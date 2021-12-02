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
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

var pluginContact = ExternalInterface{
	ContactClient:   mockContactClient,
	Auth:            mockIsAuthorized,
	CreateChildTask: mockCreateChildTask,
	UpdateTask:      mockUpdateTask,
	DecryptPassword: stubDevicePassword,
	GetPluginStatus: GetPluginStatusForTesting,
}

func mockCreateChildTask(sessionID, taskID string) (string, error) {
	switch taskID {
	case "taskWithoutChild":
		return "", fmt.Errorf("subtask cannot created")
	case "subTaskWithSlash":
		return "someSubTaskID/", nil
	default:
		return "someSubTaskID", nil
	}
}

func mockSystemData(systemID string) error {
	reqData, _ := json.Marshal(&map[string]interface{}{
		"Id": "1",
	})

	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("ComputerSystem", systemID, string(reqData)); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "System", err.Error())
	}
	return nil
}

func mockUpdateTask(task common.TaskData) error {
	if task.TaskID == "invalid" {
		return fmt.Errorf("task with this ID not found")
	}
	return nil
}

func TestPluginContact_SetDefaultBootOrderSystems(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
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
	device1 := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "24b243cf-f1e3-5318-92d9-2d6737d6b0b9",
		PluginID:       "GRF",
	}
	device2 := agmodel.Target{
		ManagerAddress: "100.0.0.2",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "7a2c6100-67da-5fd6-ab82-6870d29c7279",
		PluginID:       "GRF",
	}
	device3 := agmodel.Target{
		ManagerAddress: "100.0.0.2",
		Password:       []byte("passwordWithInvalidEncryption"),
		UserName:       "admin",
		DeviceUUID:     "7a2c6100-67da-5fd6-ab82-6870d29c7279",
		PluginID:       "GRF",
	}
	device4 := agmodel.Target{
		ManagerAddress: "100.0.0.2",
		Password:       []byte("someValidPassword"),
		UserName:       "admin",
		DeviceUUID:     "8e896459-a8f9-4c83-95b7-7b316b4908e1",
		PluginID:       "Unknown-Plugin",
	}

	device5 := agmodel.Target{
		ManagerAddress: "100.0.0.5",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "2aca8daa-9c20-4a5b-9203-469a24f452c8",
		PluginID:       "XAuthPlugin",
	}
	device6 := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "9dd6e488-31b2-475a-9304-d5f193a6a7cd",
		PluginID:       "XAuthPluginFail",
	}
	mockPluginData(t, "GRF")
	mockPluginData(t, "XAuthPlugin")
	mockPluginData(t, "XAuthPluginFail")

	mockDeviceData("24b243cf-f1e3-5318-92d9-2d6737d6b0b9", device1)
	mockDeviceData("7a2c6100-67da-5fd6-ab82-6870d29c7279", device2)
	mockDeviceData("123443cf-f1e3-5318-92d9-2d6737d65678", device3)
	mockDeviceData("8e896459-a8f9-4c83-95b7-7b316b4908e1", device4)
	mockDeviceData("2aca8daa-9c20-4a5b-9203-469a24f452c8", device5)
	mockDeviceData("9dd6e488-31b2-475a-9304-d5f193a6a7cd", device6)

	mockSystemData("/redfish/v1/Systems/7a2c6100-67da-5fd6-ab82-6870d29c7279.1")
	mockSystemData("/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b9.1")
	mockSystemData("/redfish/v1/Systems/s83405033-67da-5fd6-ab82-458292935.1")
	mockSystemData("/redfish/v1/Systems/123443cf-f1e3-5318-92d9-2d6737d65678.1")
	mockSystemData("/redfish/v1/Systems/8e896459-a8f9-4c83-95b7-7b316b4908e1.1")
	mockSystemData("/redfish/v1/Systems/2aca8daa-9c20-4a5b-9203-469a24f452c8.1")
	mockSystemData("/redfish/v1/Systems/9dd6e488-31b2-475a-9304-d5f193a6a7cd.1")

	type args struct {
		taskID, sessionUserName string
		req                     *aggregatorproto.AggregatorRequest
	}
	positiveReqData, _ := json.Marshal(AggregationSetDefaultBootOrderRequest{
		Systems: []OdataID{
			{
				OdataID: "/redfish/v1/Systems/7a2c6100-67da-5fd6-ab82-6870d29c7279.1",
			},
			{
				OdataID: "/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b9.1",
			},
		},
	})
	invalidUUIDReqData, _ := json.Marshal(AggregationSetDefaultBootOrderRequest{
		Systems: []OdataID{
			{
				OdataID: "/redfish/v1/Systems/7a2c6100-67da-5fd6-ab82-6870d29c7279.1",
			},
			{
				OdataID: "/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b.1",
			},
		},
	})

	invalidSystemReqData, _ := json.Marshal(AggregationSetDefaultBootOrderRequest{
		Systems: []OdataID{
			{
				OdataID: "/redfish/v1/Systems/7a2c6100-67da-5fd6-ab82-6870d29c7279.1",
			},
			{
				OdataID: "/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b9",
			},
		},
	})

	noUUIDInDBReqData, _ := json.Marshal(AggregationSetDefaultBootOrderRequest{
		Systems: []OdataID{
			{
				OdataID: "/redfish/v1/Systems/s83405033-67da-5fd6-ab82-458292935.1",
			},
		},
	})

	decryptionFailReqData, _ := json.Marshal(AggregationSetDefaultBootOrderRequest{
		Systems: []OdataID{
			{
				OdataID: "/redfish/v1/Systems/123443cf-f1e3-5318-92d9-2d6737d65678.1",
			},
		},
	})

	unknownPluginReqData, _ := json.Marshal(AggregationSetDefaultBootOrderRequest{
		Systems: []OdataID{
			{
				OdataID: "/redfish/v1/Systems/8e896459-a8f9-4c83-95b7-7b316b4908e1.1",
			},
		},
	})

	positiveXAuthPluginReqData, _ := json.Marshal(AggregationSetDefaultBootOrderRequest{
		Systems: []OdataID{
			{
				OdataID: "/redfish/v1/Systems/2aca8daa-9c20-4a5b-9203-469a24f452c8.1",
			},
		},
	})

	XAuthPluginFailedReqData, _ := json.Marshal(AggregationSetDefaultBootOrderRequest{
		Systems: []OdataID{
			{
				OdataID: "/redfish/v1/Systems/9dd6e488-31b2-475a-9304-d5f193a6a7cd.1",
			},
		},
	})
	tests := []struct {
		name string
		p    *ExternalInterface
		args args
		want response.RPC
	}{
		{
			name: "postive test Case",
			p:    &pluginContact,
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  positiveReqData,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "invalid uuid id",
			p:    &pluginContact,
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  invalidUUIDReqData,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "invalid system id",
			p:    &pluginContact,
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  invalidSystemReqData,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "invalid request body",
			p:    &pluginContact,
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  []byte("invalidData"),
				},
			},
			want: response.RPC{
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "subtask creation failure",
			p:    &pluginContact,
			args: args{
				taskID: "taskWithoutChild", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  positiveReqData,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "no UUID in DB",
			p:    &pluginContact,
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  noUUIDInDBReqData,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "decryption failure",
			p:    &pluginContact,
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  decryptionFailReqData,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "unknown plugin",
			p:    &pluginContact,
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  unknownPluginReqData,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "postive test Case with a slash",
			p:    &pluginContact,
			args: args{
				taskID: "subTaskWithSlash", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  positiveReqData,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "postive test Case with XAuthToken",
			p:    &pluginContact,
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  positiveXAuthPluginReqData,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "XAuthToken failure",
			p:    &pluginContact,
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  XAuthPluginFailedReqData,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusUnauthorized,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.SetDefaultBootOrder(tt.args.taskID, tt.args.sessionUserName, tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.SetDefaultBootOrder() = %v, want %v", got.StatusCode, tt.want.StatusCode)
			}
		})
	}
}
