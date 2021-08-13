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

	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
)

func testSystemIndex(uuid string, indexData map[string]interface{}) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}

	if err := connPool.CreateIndex(indexData, "/redfish/v1/Systems/"+uuid); err != nil {
		return fmt.Errorf("error while creating  the index: %v", err.Error())
	}

	return nil

}
func testUpdateContactClient(url, method, token string, odataID string, body interface{}, credentials map[string]string) (*http.Response, error) {
	if url == "" {
		return nil, fmt.Errorf("InvalidRequest")
	}
	uid := "1s7sda8asd-asdas8as012"
	var bData agmodel.SaveSystem
	bBytes, _ := json.Marshal(body)
	json.Unmarshal(bBytes, &bData)
	host := strings.Split(url, "/ODIM")[0]
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
		body := `{"@odata.id":"/ODIM/v1/Systems/1", "UUID": "1s7sda8asd-asdas8as012", "Id": "1",
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
		body := `{"@odata.id":"/ODIM/v1/Managers/1", "UUID": "1234877451-1234", "Id": "1","FirmwareVersion": "iLO 5 v2.12"}`
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
		body := `{"@odata.id":"/ODIM/v1/Managers/1", "UUID": "1234877451-1234", "Id": "1"}`
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

func TestExternalInterface_UpdateAggregationSource(t *testing.T) {
	mockPluginData(t, "ILO_v1.0.0")
	mockPluginData(t, "GRF_v1.0.0")

	reqManagerGRF := agmodel.AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/058c1876-6f24-439a-8968-2af26154081f",
			},
		},
	}
	reqManagerILO := agmodel.AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/3489af48-2e99-4d78-a250-b04641e9d98d",
			},
		},
	}
	reqBMC := agmodel.AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/c41cbd97-937d-1b73-c41c-1b7385d39069",
			},
		},
	}
	mockDeviceData("123456", agmodel.Target{
		ManagerAddress: "100.0.1.1",
		UserName:       "admin",
		Password:       []byte("admin12345"),
		PluginID:       "ILO_v1.0.0",
	})
	err := agmodel.AddAggregationSource(reqManagerGRF, "/redfish/v1/AggregationService/AggregationSources/123455")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = agmodel.AddAggregationSource(reqManagerILO, "/redfish/v1/AggregationService/AggregationSources/123457")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	err = agmodel.AddAggregationSource(reqBMC, "/redfish/v1/AggregationService/AggregationSources/123456")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	dbErr := testSystemIndex("123456:1", map[string]interface{}{
		"UUID": "1s7sda8asd-asdas8as012",
	})
	if err != nil {
		t.Fatalf("error: %v", dbErr)
	}
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	successReqManager, _ := json.Marshal(map[string]string{
		"HostName": "100.0.0.1:50000",
		"UserName": "admin",
		"Password": "password",
	})
	invalidReqBody1, _ := json.Marshal(map[string]string{
		"HostName": "",
		"UserName": "admin",
		"Password": "password",
	})
	invalidReqBody2, _ := json.Marshal(map[string]string{
		"HostName": "100.0.0.1:50000",
		"UserName": "admin",
		"Password": "",
	})
	successReqBMC, _ := json.Marshal(map[string]string{
		"HostName": "100.0.0.1",
		"UserName": "admin",
		"Password": "password",
	})
	invalidReqBody3, _ := json.Marshal(map[string]string{
		"HostName": "100.0.0.1:50000",
		"UserName": "",
		"Password": "admin",
	})
	invalidReqBody4, _ := json.Marshal(map[string]string{
		"HostName": "100.0.0.1",
		"UserName": "admin1",
		"Password": "incorrectPassword",
	})

	missingparamReq, _ := json.Marshal(map[string]interface{}{})

	commonResponse1 := response.Response{
		OdataType:    "#AggregationSource.v1_0_0.AggregationSource",
		OdataID:      "/redfish/v1/AggregationService/AggregationSources/123455",
		OdataContext: "/redfish/v1/$metadata#AggregationSource.AggregationSource",
		ID:           "123455",
		Name:         "Aggregation Source",
	}
	var resp1 = response.RPC{
		StatusCode:    http.StatusOK,
		StatusMessage: response.Success,
	}
	resp1.Header = map[string]string{
		"Allow":             `"GET","PATCH","DELETE"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	commonResponse1.CreateGenericResponse(response.Success)
	commonResponse1.Message = ""
	commonResponse1.MessageID = ""
	commonResponse1.Severity = ""
	resp1.Body = agresponse.AggregationSourceResponse{
		Response: commonResponse1,
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Links:    reqManagerGRF.Links,
	}
	commonResponse2 := response.Response{
		OdataType:    "#AggregationSource.v1_0_0.AggregationSource",
		OdataID:      "/redfish/v1/AggregationService/AggregationSources/123456",
		OdataContext: "/redfish/v1/$metadata#AggregationSource.AggregationSource",
		ID:           "123456",
		Name:         "Aggregation Source",
	}
	var resp2 = response.RPC{
		StatusCode:    http.StatusOK,
		StatusMessage: response.Success,
	}
	resp2.Header = map[string]string{
		"Allow":             `"GET","PATCH","DELETE"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	commonResponse2.CreateGenericResponse(response.Success)
	commonResponse2.Message = ""
	commonResponse2.MessageID = ""
	commonResponse2.Severity = ""
	resp2.Body = agresponse.AggregationSourceResponse{
		Response: commonResponse2,
		HostName: "100.0.0.1",
		UserName: "admin",
		Links:    reqBMC.Links,
	}
	errMsg := "error: while trying to fetch Aggregation Source data: no data with the with key /redfish/v1/AggregationService/AggregationSources/123466 found"
	resp3 := common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"AggregationSource", "/redfish/v1/AggregationService/AggregationSources/123466"}, nil)
	param := "HostName "
	errMsg = "field " + param + " Missing"
	resp4 := common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{param}, nil)
	param = "UserName "
	errMsg = "field " + param + " Missing"
	resp5 := common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{param}, nil)
	param = "Password "
	errMsg = "field " + param + " Missing"
	resp6 := common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{param}, nil)
	param = "HostName UserName Password "
	errMsg = "error while trying to authenticate the compute server: error: invalid resource username/password"
	resp7 := common.GeneralError(http.StatusUnauthorized, response.ResourceAtURIUnauthorized, errMsg, []interface{}{"https://localhost:9091/ODIM/v1/validate"}, nil)
	errMsg = "field " + param + " Missing"
	resp8 := common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{param}, nil)

	common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{param}, nil)
	p := getMockExternalInterface()
	p.ContactClient = testUpdateContactClient
	type args struct {
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		args args
		want response.RPC
	}{
		{
			name: "Positive case Manager",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSources/123455",
					RequestBody: successReqManager,
				},
			},
			want: resp1,
		},
		{
			name: "Positive case BMC",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSources/123456",
					RequestBody: successReqBMC,
				},
			},
			want: resp2,
		},
		{
			name: "invalid aggregation source id",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSources/123466",
					RequestBody: successReqManager,
				},
			},
			want: resp3,
		},
		{
			name: "invalid request manger address missing",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSources/123457",
					RequestBody: invalidReqBody1,
				},
			},
			want: resp4,
		},
		{
			name: "invalid request UserName missing",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSources/123457",
					RequestBody: invalidReqBody3,
				},
			},
			want: resp5,
		},
		{
			name: "invalid request Password missing",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSources/123457",
					RequestBody: invalidReqBody2,
				},
			},
			want: resp6,
		},
		{
			name: "invalid request wrong Password",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSources/123456",
					RequestBody: invalidReqBody4,
				},
			},
			want: resp7,
		},
		{
			name: "invalid request body missing",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSources/123457",
					RequestBody: missingparamReq,
				},
			},
			want: resp8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.UpdateAggregationSource(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExternalInterface.UpdateAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}
