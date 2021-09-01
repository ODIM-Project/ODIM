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
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

func EventFunctionsForTesting(s []string)                 {}
func PostEventFunctionForTesting(s []string, name string) {}
func GetPluginStatusForTesting(plugin agmodel.Plugin) bool {
	return true
}
func mockSubscribeEMB(pluginID string, list []string) {
	return
}

func mockContactClientForDuplicate(url, method, token string, odataID string, body interface{}, credentials map[string]string) (*http.Response, error) {
	var bData agmodel.SaveSystem
	bBytes, _ := json.Marshal(body)
	json.Unmarshal(bBytes, &bData)
	host := strings.Split(url, "/ODIM")[0]
	uid := "test1"
	if url == "https://localhost:9091/ODIM/v1/Systems/1/Actions/ComputerSystem.Add" {
		body := `{"MessageId": "` + response.Success + `"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9091/ODIM/v1/Registries" {
		body := `{"Members":[{"@odata.id":"/redfish/v1/Registries/Base"},{"@odata.id":"/redfish/v1/Registries/SomeMember"}]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9091/ODIM/v1/Registries/Base" {
		body := `{"@odata.context":"/redfish/v1/$metadata#MessageRegistryFile.MessageRegistryFile","@odata.etag":"W/\"0DCA67A0\"","@odata.id":"/redfish/v1/Registries/Base","@odata.type":"#MessageRegistryFile.v1_0_4.MessageRegistryFile","Id":"Base","Description":"Registry Definition File for Base","Languages":["en"],"Location":[{"Language":"en","Uri":"/redfish/v1/RegistryStore/registries/en/Base.json"}],"Name":"Base Message Registry File","Registry":"Base.1.4.0"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9091/ODIM/v1/Registries/SomeMember" {
		body := `{"@odata.context":"/redfish/v1/$metadata#MessageRegistryFile.MessageRegistryFile","@odata.etag":"W/\"0DCA67A0\"","@odata.id":"/redfish/v1/Registries/Base","@odata.type":"#MessageRegistryFile.v1_0_4.MessageRegistryFile","Id":"Base","Description":"Registry Definition File for Base","Languages":["en"],"Location":[{"Language":"en","Uri":"/redfish/v1/RegistryStore/registries/en/SomeRegistry.json"}],"Name":"Base Message Registry File","Registry":"SomeRegistry"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if strings.Contains(url, "SomeRegistry.json") {
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
},"SystemType": "Physical",
  "Links":{
	"ManagedBy":[
	{
	"@odata.id": "/redfish/v1/Managers/1"
	}
	]
	},
	"Storage":{
		"@odata.id": "/redfish/v1/Systems/1/Storage"
		}}`
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
		body := `{"Version": "1.0.0","EventMessageBus":{"EmbQueue":[{"EmbQueueName":"GRF"}]}}`
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

	}

	return nil, fmt.Errorf("InvalidRequest")
}

func TestExternalInterface_addcompute(t *testing.T) {
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
	mockPluginData(t, "GRF")
	mockPluginData(t, "XAuthPlugin")
	mockPluginData(t, "XAuthPluginFail")

	reqSuccess := AddResourceRequest{
		ManagerAddress: "100.0.0.1",
		UserName:       "admin",
		Password:       "password",
		ConnectionMethod: &ConnectionMethod{
			OdataID: "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
		},
	}

	reqSuccessXAuth := AddResourceRequest{
		ManagerAddress: "100.0.0.2",
		UserName:       "admin",
		Password:       "password",
		ConnectionMethod: &ConnectionMethod{
			OdataID: "/redfish/v1/AggregationService/ConnectionMethods/0a8992dc-8b47-4fe3-b26c-4c34048cf0d2",
		},
	}
	reqIncorrectDeviceBasicAuth := AddResourceRequest{
		ManagerAddress: "100.0.0.1",
		UserName:       "admin1",
		Password:       "incorrectPassword",
		ConnectionMethod: &ConnectionMethod{
			OdataID: "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
		},
	}

	p := getMockExternalInterface()
	targetURI := "/redfish/v1/AggregationService/AggregationSource"
	var percentComplete int32
	var pluginContactRequest getResourceRequest
	pluginContactRequest.ContactClient = p.ContactClient
	pluginContactRequest.GetPluginStatus = p.GetPluginStatus
	pluginContactRequest.TargetURI = targetURI
	pluginContactRequest.UpdateTask = p.UpdateTask
	type args struct {
		taskID   string
		req      AddResourceRequest
		pluginID string
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
				taskID:   "123",
				req:      reqSuccess,
				pluginID: "GRF",
			},
			want: response.RPC{
				StatusCode: http.StatusCreated,
			},
		},
		{
			name: "success: plugin with xauth token",
			p:    p,
			args: args{
				taskID:   "123",
				req:      reqSuccessXAuth,
				pluginID: "XAuthPlugin",
			},
			want: response.RPC{
				StatusCode: http.StatusCreated,
			},
		},
		{
			name: "with incorrect device credentials and BasicAuth",
			p:    p,
			args: args{
				taskID:   "123",
				req:      reqIncorrectDeviceBasicAuth,
				pluginID: "GRF",
			},
			want: response.RPC{
				StatusCode: http.StatusUnauthorized,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _, _ := tt.p.addCompute(tt.args.taskID, targetURI, tt.args.pluginID, percentComplete, tt.args.req, pluginContactRequest); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.addCompute = %v, want %v", got, tt.want)
			}
		})
	}
}
