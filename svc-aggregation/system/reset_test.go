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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	uuid "github.com/satori/go.uuid"
)

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockPluginData(t *testing.T, pluginID string) error {
	password := getEncryptedKey(t, []byte("password"))
	plugin := agmodel.Plugin{
		IP:                "localhost",
		Port:              "9091",
		Username:          "admin",
		Password:          password,
		ID:                pluginID,
		PreferredAuthType: "BasicAuth",
		ManagerUUID:       "1s7sda8asd-asdas8as0",
	}
	switch pluginID {
	case "XAuthPlugin":
		plugin.PreferredAuthType = "XAuthToken"
	case "XAuthPluginFail":
		plugin.PreferredAuthType = "XAuthToken"
		plugin.Username = "incorrectusername"
	case "NoStatusPlugin":
		plugin.Username = "noStatusUser"
		plugin.ManagerUUID = "1234877451-1235"
	case "GRF":
		plugin.ManagerUUID = "1234877451-1234"
	case "ILO":
		plugin.ManagerUUID = "1234877451-1233"
	case "XAuthPlugin_v1.0.0":
		plugin.PreferredAuthType = "XAuthToken"
	case "XAuthPluginFail_v1.0.0":
		plugin.PreferredAuthType = "XAuthToken"
		plugin.Username = "incorrectusername"
	case "NoStatusPlugin_v1.0.0":
		plugin.Username = "noStatusUser"
		plugin.ManagerUUID = "1234877451-1235"
	case "GRF_v1.0.0":
		plugin.ManagerUUID = "1234877451-1234"
	case "ILO_v1.0.0":
		plugin.ManagerUUID = "1234877451-1233"
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
func mockDeviceData(uuid string, device agmodel.Target) error {

	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("System", uuid, device); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "System", err.Error())
	}
	return nil
}

func mockIsAuthorized(sessionToken string, privileges, oemPrivileges []string) response.RPC {
	if sessionToken != "validToken" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "", nil, nil)
	}
	return common.GeneralError(http.StatusOK, response.Success, "", nil, nil)
}

func mockContactClient(url, method, token string, odataID string, body interface{}, credentials map[string]string) (*http.Response, error) {
	if url == "" {
		return nil, fmt.Errorf("InvalidRequest")
	}
	var bData agmodel.SaveSystem
	bBytes, _ := json.Marshal(body)
	json.Unmarshal(bBytes, &bData)
	host := strings.Split(url, "/ODIM")[0]
	uid := uuid.NewV4().String()
	if url == "https://localhost:9091/ODIM/v1/Systems/1/Actions/ComputerSystem.Reset" || url == "https://localhost:9091/ODIM/v1/Systems/1/Actions/ComputerSystem.Add" ||
		url == "https://localhost:9091/ODIM/v1/Systems/1/Actions/ComputerSystem.SetDefaultBootOrder" {
		body := `{"MessageId": "` + response.Success + `"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9091/ODIM/v1/Systems" {
		body := `{"Members":[{"@odata.id":"/ODIM/v1/Systems/1"}]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == "https://localhost:9091/ODIM/v1/Systems/1" {
		body := `{"@odata.id":"/ODIM/v1/Systems/1", "UUID": "` + uid + `", "Id": "1",
		    "MemorySummary":{
			"Status":{
			"HealthRollup": "OK"
			},
			"TotalSystemMemoryGiB": 384,
			"TotalSystemPersistentMemoryGiB": 0
			},
			"PowerState": "On",
			"ProcessorSummary":{
				"Count": 2,
				"Model": "Intel(R) Xeon(R) Gold 6152 CPU @ 2.10GHz",
				"Status":{
					"HealthRollup": "OK"
				}
			},
			"SystemType": "Physical",
			"Links":{
				"ManagedBy":[
				{
					"@odata.id": "/redfish/v1/Managers/1"
				}
				]
			},
			"Storage":{
				"@odata.id": "/redfish/v1/Systems/1/Storage"
			}
		}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == "https://localhost:9091/ODIM/v1/Chassis" {
		body := `{"Members":[{"@odata.id":"/ODIM/v1/Chassis/1"}]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == "https://localhost:9091/ODIM/v1/Chassis/1" {
		body := `{"@odata.id":"/ODIM/v1/Chassis/1", "UUID": "` + uid + `", "Id": "1"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == "https://localhost:9091/ODIM/v1/Managers/1" {
		body := `{"@odata.id":"/ODIM/v1/Managers/1", "UUID": "1s7sda8asd-asdas8as0", "Id": "1","FirmwareVersion": "iLO 5 v2.12"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == "https://localhost:9091/ODIM/v1/Systems/1/Storage" {
		body := `{"Members":[{"@odata.id":"/ODIM/v1/Systems/1/Storage/1"}]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == "https://localhost:9091/ODIM/v1/Systems/1/Storage/1" {
		body := `{"Drives":[
			{
			"@odata.id": "/ODIM/v1/Systems/1/Storage/1/Drives/0"
			}
			]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == "https://localhost:9091/ODIM/v1/Systems/1/Storage/1/Drives/0" {
		body := `{"BlockSizeBytes": 512,
		"CapacityBytes": 2179699264,
		"MediaType": "HDD"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == host+"/ODIM/v1/Managers" {
		body := `{"Members":[{"@odata.id":"/ODIM/v1/Managers/1"}]}`
		if host == "https://100.0.0.5:9091" {
			return nil, fmt.Errorf("manager data not available not reachable")
		}
		if host == "https://100.0.0.6:9091" {
			body = "incorrectResponse"
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == host+"/ODIM/v1/Managers/1" {
		body := `{"@odata.id":"/ODIM/v1/Managers/1", "UUID": "1s7sda8asd-asdas8as0", "Id": "1"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == host+"/ODIM/v1/Status" {
		body := `{"Version": "v1.0.0","EventMessageBus":{"EmbQueue":[{"EmbQueueName":"GRF"}]}}`
		if host == "https://100.0.0.3:9091" {
			return nil, fmt.Errorf("plugin not reachable")
		}
		if host == "https://100.0.0.4:9091" {
			body = "incorrectResponse"
		}
		if host == "https://100.0.0.1:" || host == "https://100.0.0.2:" {
			body = "not found"
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if strings.Contains(url, "/ODIM/v1/validate") || url == "https://localhost:9091/ODIM/v1/Sessions" || url == host+"/ODIM/v1/Sessions" {
		body := `{"MessageId": "` + response.Success + `"}`
		if bData.UserName == "incorrectusername" || bytes.Compare(bData.Password, []byte("incorrectPassword")) == 0 {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if strings.Contains(url, "/ODIM/v1/Registries") {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		}, nil
	} else {
		return &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		}, nil
	}
}

func stubDevicePassword(password []byte) ([]byte, error) {
	if bytes.Compare(password, []byte("passwordWithInvalidEncryption")) == 0 {
		return []byte{}, fmt.Errorf("password decryption failed")
	}
	return password, nil
}

func TestPluginContact_ResetComputerSystem(t *testing.T) {
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
		ManagerAddress: "100.0.0.2",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "7a2c6100-67da-5fd6-ab82-6870d29c7279",
		PluginID:       "GRF",
	}
	device3 := agmodel.Target{
		ManagerAddress: "100.0.0.6",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
		PluginID:       "SOME-INVALID-PLUGIN",
	}
	device4 := agmodel.Target{
		ManagerAddress: "100.0.0.5",
		Password:       []byte("passwordWithInvalidEncryption"),
		UserName:       "admin",
		DeviceUUID:     "c14d91b5-3333-48bb-a7b7-75f74a137d48",
		PluginID:       "GRF",
	}

	device5 := agmodel.Target{
		ManagerAddress: "100.0.0.7",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "8e896459-a8f9-4c83-95b7-7b316b4908e1",
		PluginID:       "XAuthPlugin",
	}
	device6 := agmodel.Target{
		ManagerAddress: "100.0.0.8",
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
	mockDeviceData("6d4a0a66-7efa-578e-83cf-44dc68d2874e", device3)
	mockDeviceData("c14d91b5-3333-48bb-a7b7-75f74a137d48", device4)
	mockDeviceData("8e896459-a8f9-4c83-95b7-7b316b4908e1", device5)
	mockDeviceData("9dd6e488-31b2-475a-9304-d5f193a6a7cd", device6)
	mockSystemData("/redfish/v1/Systems/7a2c6100-67da-5fd6-ab82-6870d29c7279:1")
	mockSystemData("/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b9:1")
	mockSystemData("/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1")
	mockSystemData("/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1")
	mockSystemData("/redfish/v1/Systems/8e896459-a8f9-4c83-95b7-7b316b4908e1:1")
	mockSystemData("/redfish/v1/Systems/9dd6e488-31b2-475a-9304-d5f193a6a7cd:1")

	type args struct {
		taskID          string
		sessionUserName string
		req             *aggregatorproto.AggregatorRequest
	}

	pluginContact := ExternalInterface{
		ContactClient:   mockContactClient,
		Auth:            mockIsAuthorized,
		CreateChildTask: mockCreateChildTask,
		UpdateTask:      mockUpdateTask,
		DecryptPassword: stubDevicePassword,
		GetPluginStatus: GetPluginStatusForTesting,
	}

	successReq, _ := json.Marshal(AggregationResetRequest{
		BatchSize:                    1,
		DelayBetweenBatchesInSeconds: 2,
		ResetType:                    "ForceRestart",
		TargetURIs: []string{
			"/redfish/v1/Systems/7a2c6100-67da-5fd6-ab82-6870d29c7279:1",
			"/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b9:1",
		},
	})

	invalidUUIDReq, _ := json.Marshal(AggregationResetRequest{
		BatchSize:                    1,
		DelayBetweenBatchesInSeconds: 2,
		ResetType:                    "ForceRestart",
		TargetURIs: []string{
			"/redfish/v1/Systems/7a2c6100-67da-5fd6-ab82-6870d29c7279:1",
			"/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b:1",
		},
	})

	invalidSysIDReq, _ := json.Marshal(AggregationResetRequest{
		BatchSize:                    1,
		DelayBetweenBatchesInSeconds: 2,
		ResetType:                    "ForceRestart",
		TargetURIs: []string{
			"/redfish/v1/Systems/7a2c6100-67da-5fd6-ab82-6870d29c7279:1",
			"/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b",
		},
	})

	emptyResetTypeReq, _ := json.Marshal(AggregationResetRequest{
		BatchSize:                    1,
		DelayBetweenBatchesInSeconds: 2,
		ResetType:                    "",
		TargetURIs: []string{
			"/redfish/v1/Systems/7a2c6100-67da-5fd6-ab82-6870d29c7279:1",
			"/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b9:1",
		},
	})

	emptyTargetURIsReq, _ := json.Marshal(AggregationResetRequest{
		BatchSize:                    1,
		DelayBetweenBatchesInSeconds: 2,
		ResetType:                    "ForceRestart",
		TargetURIs:                   []string{},
	})

	InvalidPasswordReq, _ := json.Marshal(AggregationResetRequest{
		BatchSize:                    1,
		DelayBetweenBatchesInSeconds: 2,
		ResetType:                    "ForceRestart",
		TargetURIs: []string{
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1",
		},
	})

	invalidPluginReq, _ := json.Marshal(AggregationResetRequest{
		BatchSize:                    1,
		DelayBetweenBatchesInSeconds: 2,
		ResetType:                    "ForceRestart",
		TargetURIs: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
		},
	})

	XAuthPluginReq, _ := json.Marshal(AggregationResetRequest{
		BatchSize:                    1,
		DelayBetweenBatchesInSeconds: 2,
		ResetType:                    "ForceRestart",
		TargetURIs: []string{
			"/redfish/v1/Systems/8e896459-a8f9-4c83-95b7-7b316b4908e1:1",
		},
	})

	XAuthPluginFailedReq, _ := json.Marshal(AggregationResetRequest{
		BatchSize:                    1,
		DelayBetweenBatchesInSeconds: 2,
		ResetType:                    "ForceRestart",
		TargetURIs: []string{
			"/redfish/v1/Systems/9dd6e488-31b2-475a-9304-d5f193a6a7cd:1",
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
					RequestBody:  successReq,
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
					RequestBody:  invalidUUIDReq,
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
					RequestBody:  invalidSysIDReq,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "invalid req body",
			p:    &pluginContact,
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  []byte("some invalid request"),
				},
			},
			want: response.RPC{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "request missing TargetURIs",
			p:    &pluginContact,
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  emptyTargetURIsReq,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "request missisng ResetTYpe",
			p:    &pluginContact,
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  emptyResetTypeReq,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "reset without child task",
			p:    &pluginContact,
			args: args{
				taskID: "taskWithoutChild", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  successReq,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "reset without slash in subtask",
			p:    &pluginContact,
			args: args{
				taskID: "subTaskWithSlash", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  successReq,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "device password decryption fails",
			p:    &pluginContact,
			args: args{
				taskID: "someId", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  InvalidPasswordReq,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "invalid plugin",
			p:    &pluginContact,
			args: args{
				taskID: "someId", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  invalidPluginReq,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "xauth plugin",
			p:    &pluginContact,
			args: args{
				taskID: "someId", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  XAuthPluginReq,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "xauth fails",
			p:    &pluginContact,
			args: args{
				taskID: "someId", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  XAuthPluginFailedReq,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusUnauthorized,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Reset(tt.args.taskID, tt.args.sessionUserName, tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.Reset() = %v, want %v", got, tt.want)
			}
		})
	}
}
