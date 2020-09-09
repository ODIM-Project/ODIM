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
 "strings"
 "io/ioutil"
 "fmt"
 "bytes"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
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

func testContactClientWithDelay(url, method, token string, odataID string, body interface{}, credentials map[string]string) (*http.Response, error) {
	time.Sleep(4 * time.Second)
	if strings.Contains(url, "/ODIM/v1/Systems/1") {
		uid := "24b243cf-f1e3-5318-92d9-2d6737d6b0b9"
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
		}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	fBody := `{"Members":[{"@odata.id":"/ODIM/v1/Systems/1"}]}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString(fBody)),
	}, nil
}

func mockContactClientForDuplicate(url, method, token string, odataID string, body interface{}, credentials map[string]string) (*http.Response, error) {
	var bData agmodel.SaveSystem
	bBytes, _ := json.Marshal(body)
	json.Unmarshal(bBytes, &bData)
	host := strings.Split(url, "/ODIM")[0]
	uid := "test1"
	if url == "https://localhost:9091/ODIM/v1/Systems/1/Actions/ComputerSystem.Add" {
		body := `{"MessageId": "Base.1.0.Success"}`
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
		body := `{"MessageId": "Base.1.0.Success"}`
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
		body := `{"EventMessageBus":{"EmbQueue":[{"EmbQueueName":"GRF"}]}}`
		if host == "https://100.0.0.3:9091" {
			return nil, fmt.Errorf("plugin not reachable")
		}
		if host == "https://100.0.0.4:9091" {
			body = "incorrectResponse"
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if strings.Contains(url, "/ODIM/v1/validate") || url == "https://localhost:9091/ODIM/v1/Sessions" || url == host+"/ODIM/v1/Sessions" {
		body := `{"MessageId": "Base.1.0.Success"}`
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

func TestExternalInterface_AddBMC(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
	addComputeRetrieval := config.AddComputeSkipResources{
		SystemCollection: []string{"Chassis", "LogServices"},
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

	reqSuccess, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			Oem: &AddOEM{
				PluginID: "GRF",
			},
		},
	})
	reqWithoutOEM, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.11",
		UserName: "admin",
		Password: "password",
	})
	reqPluginID, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			Oem: &AddOEM{
				PluginID: "invalidpluginid",
			},
		},
	})
	reqSuccessXAuth, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.2",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			Oem: &AddOEM{
				PluginID: "XAuthPlugin",
			},
		},
	})
	reqIncorrectDeviceBasicAuth, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin1",
		Password: "incorrectPassword",
		Links: &Links{
			Oem: &AddOEM{
				PluginID: "GRF",
			},
		},
	})
	reqIncorrectDeviceXAuth, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.2",
		UserName: "username",
		Password: "password",
		Links: &Links{
			Oem: &AddOEM{
				PluginID: "XAuthPluginFail",
			},
		},
	})
	p := &ExternalInterface{
		ContactClient:       mockContactClient,
		Auth:                mockIsAuthorized,
		CreateChildTask:     mockCreateChildTask,
		UpdateTask:          mockUpdateTask,
		CreateSubcription:   EventFunctionsForTesting,
		PublishEvent:        PostEventFunctionForTesting,
		GetPluginStatus:     GetPluginStatusForTesting,
		EncryptPassword:     stubDevicePassword,
		DecryptPassword:     stubDevicePassword,
		DeleteComputeSystem: deleteComputeforTest,
	}
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
			name: "request without OEM",
			p:    p,
			args: args{
				taskID: "123",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  reqWithoutOEM,
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
		ActiveReqSet.UpdateMu.Lock()
		ActiveReqSet.ReqRecord = make(map[string]interface{})
		ActiveReqSet.UpdateMu.Unlock()
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(2 * time.Second)
			if got := tt.p.AddAggregationSource(tt.args.taskID, "validUserName", tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.AddAggregationSource = %v, want %v", got, tt.want)
			}
		})
		ActiveReqSet.UpdateMu.Lock()
		ActiveReqSet.ReqRecord = nil
		ActiveReqSet.UpdateMu.Unlock()
	}
}

func TestExternalInterface_AddBMCForPasswordEncryptFail(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
	addComputeRetrieval := config.AddComputeSkipResources{
		SystemCollection: []string{"Chassis", "LogServices"},
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

	reqEncryptFail, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin",
		Password: "passwordWithInvalidEncryption",
		Links: &Links{
			Oem: &AddOEM{
				PluginID: "GRF",
			},
		},
	})
	p := &ExternalInterface{
		ContactClient:       mockContactClient,
		Auth:                mockIsAuthorized,
		CreateChildTask:     mockCreateChildTask,
		UpdateTask:          mockUpdateTask,
		CreateSubcription:   EventFunctionsForTesting,
		PublishEvent:        PostEventFunctionForTesting,
		GetPluginStatus:     GetPluginStatusForTesting,
		EncryptPassword:     stubDevicePassword,
		DecryptPassword:     stubDevicePassword,
		DeleteComputeSystem: deleteComputeforTest,
	}
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
		ActiveReqSet.UpdateMu.Lock()
		ActiveReqSet.ReqRecord = make(map[string]interface{})
		ActiveReqSet.UpdateMu.Unlock()
		t.Run(tt.name, func(t *testing.T) {
			time.Sleep(2 * time.Second)
			if got := tt.p.AddAggregationSource(tt.args.taskID, "validUserName", tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.AddAggregationSource = %v, want %v", got, tt.want)
			}
		})
		ActiveReqSet.UpdateMu.Lock()
		ActiveReqSet.ReqRecord = nil
		ActiveReqSet.UpdateMu.Unlock()
	}
}

// TestExternalInterface_AddBMCDuplicate handles the test cases for getregistry and duplicate add server
func TestExternalInterface_AddBMCDuplicate(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
	addComputeRetrieval := config.AddComputeSkipResources{
		SystemCollection: []string{"Chassis", "LogServices"},
	}
	config.Data.AddComputeSkipResources = &addComputeRetrieval
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockPluginData(t, "GRF")

	reqSuccess, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			Oem: &AddOEM{
				PluginID: "GRF",
			},
		},
	})
	p := &ExternalInterface{
		ContactClient:       mockContactClientForDuplicate,
		Auth:                mockIsAuthorized,
		CreateChildTask:     mockCreateChildTask,
		UpdateTask:          mockUpdateTask,
		CreateSubcription:   EventFunctionsForTesting,
		PublishEvent:        PostEventFunctionForTesting,
		GetPluginStatus:     GetPluginStatusForTesting,
		EncryptPassword:     stubDevicePassword,
		DecryptPassword:     stubDevicePassword,
		DeleteComputeSystem: deleteComputeforTest,
	}
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
		ActiveReqSet.UpdateMu.Lock()
		ActiveReqSet.ReqRecord = make(map[string]interface{})
		ActiveReqSet.UpdateMu.Unlock()
		t.Run(tt.name, func(t *testing.T) {
			if got := p.AddAggregationSource("123", "validUserName", req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.AddAggregationSource = %v, want %v", got, tt.want)
			}
		})
		ActiveReqSet.UpdateMu.Lock()
		ActiveReqSet.ReqRecord = nil
		ActiveReqSet.UpdateMu.Unlock()
	}
}

func TestExternalInterface_Manager(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
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
	reqSuccess, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			Oem: &AddOEM{
				PluginID:          "GRF",
				PreferredAuthType: "BasicAuth",
				PluginType:        "Compute",
			},
		},
	})
	reqExistingPlugin, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			Oem: &AddOEM{
				PluginID:          "ILO",
				PreferredAuthType: "BasicAuth",
				PluginType:        "Compute",
			},
		},
	})
	reqInvalidAuthType, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			Oem: &AddOEM{
				PluginID:          "ILO",
				PreferredAuthType: "BasicAuthentication",
				PluginType:        "Compute",
			},
		},
	})
	reqInvalidPluginType, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			Oem: &AddOEM{
				PluginID:          "ILO",
				PreferredAuthType: "BasicAuth",
				PluginType:        "plugin",
			},
		},
	})
	reqExistingPluginBadPassword, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			Oem: &AddOEM{
				PluginID:          "PluginWithBadPassword",
				PreferredAuthType: "BasicAuth",
				PluginType:        "Compute",
			},
		},
	})
	reqExistingPluginBadData, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			Oem: &AddOEM{
				PluginID:          "PluginWithBadPassword",
				PreferredAuthType: "BasicAuth",
				PluginType:        "Compute",
			},
		},
	})

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
		ActiveReqSet.UpdateMu.Lock()
		ActiveReqSet.ReqRecord = make(map[string]interface{})
		ActiveReqSet.UpdateMu.Unlock()
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.AddAggregationSource(tt.args.taskID, "validUserName", tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.AddAggregationSource() = %v, want %v", got, tt.want)
			}
		})
		ActiveReqSet.UpdateMu.Lock()
		ActiveReqSet.ReqRecord = nil
		ActiveReqSet.UpdateMu.Unlock()
	}
}

func TestExternalInterface_ManagerXAuth(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
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
	reqXAuthSuccess, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			Oem: &AddOEM{
				PluginID:          "GRF",
				PreferredAuthType: "XAuthToken",
				PluginType:        "Compute",
			},
		},
	})
	reqXAuthFail, _ := json.Marshal(AggregationSource{
		HostName: "localhost:9091",
		UserName: "incorrectusername",
		Password: "incorrectPassword",
		Links: &Links{
			Oem: &AddOEM{
				PluginID:          "ILO",
				PreferredAuthType: "XAuthToken",
				PluginType:        "Compute",
			},
		},
	})

	reqStatusFail, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.3:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			Oem: &AddOEM{
				PluginID:          "ILO",
				PreferredAuthType: "XAuthToken",
				PluginType:        "Compute",
			},
		},
	})

	reqInvalidStatusBody, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.4:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			Oem: &AddOEM{
				PluginID:          "ILO",
				PreferredAuthType: "XAuthToken",
				PluginType:        "Compute",
			},
		},
	})

	reqManagerGetFail, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.5:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			Oem: &AddOEM{
				PluginID:          "ILO",
				PreferredAuthType: "XAuthToken",
				PluginType:        "Compute",
			},
		},
	})

	reqInvalidManagerBody, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.6:9091",
		UserName: "admin",
		Password: "password",
		Links: &Links{
			Oem: &AddOEM{
				PluginID:          "ILO",
				PreferredAuthType: "XAuthToken",
				PluginType:        "Compute",
			},
		},
	})

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
		ActiveReqSet.UpdateMu.Lock()
		ActiveReqSet.ReqRecord = make(map[string]interface{})
		ActiveReqSet.UpdateMu.Unlock()
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.AddAggregationSource(tt.args.taskID, "validUserName", tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.AddAggregationSource() = %v, want %v", got, tt.want)
			}
		})
		ActiveReqSet.UpdateMu.Lock()
		ActiveReqSet.ReqRecord = nil
		ActiveReqSet.UpdateMu.Unlock()
	}
}

func TestExternalInterface_ManagerWithMultipleRequest(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
	addComputeRetrieval := config.AddComputeSkipResources{
		SystemCollection: []string{"Chassis", "LogServices"},
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
			Oem: &AddOEM{
				PluginID:          "GRF",
				PreferredAuthType: "BasicAuth",
				PluginType:        "Compute",
			},
		},
	})

	p := &ExternalInterface{
		ContactClient:     testContactClientWithDelay,
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
			name: "multiple request",
			want: response.RPC{
				StatusCode: http.StatusConflict,
			},
		},
	}
	for _, tt := range tests {
		ActiveReqSet.UpdateMu.Lock()
		ActiveReqSet.ReqRecord = make(map[string]interface{})
		ActiveReqSet.UpdateMu.Unlock()
		t.Run(tt.name, func(t *testing.T) {
			go p.AddAggregationSource("123", "validUserName", req)
			time.Sleep(time.Second)
			if got := p.AddAggregationSource("123", "validUserName", req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.AddAggregationSource() = %v, want %v", got, tt.want)
			}
		})
		ActiveReqSet.UpdateMu.Lock()
		ActiveReqSet.ReqRecord = nil
		ActiveReqSet.UpdateMu.Unlock()
	}
}
