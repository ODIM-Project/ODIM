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
package chassis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

func mockChassisResourceData(body []byte, table, key string) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err := connPool.Create(table, key, string(body)); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", table, err.Error())
	}
	return nil
}

func mockContactClient(url, method, token string, odataID string, body interface{}, basicAuth map[string]string) (*http.Response, error) {
	if url == "http://localhost:9091/redfishplugin/login" {
		header := make(http.Header)
		header.Set("X-Auth-Token", token)
		return &http.Response{
			StatusCode: http.StatusCreated,
			Header:     header,
		}, nil
	}
	if url == "/redfish/v1/Chassis/1/Power2" {
		body := `{"@odata.id": "/redfish/v1/Chassis/1/Power2"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "/ODIM/v1/Chassis/1/NetworkAdapters" {
		body := `{"Name": "NetworkAdapterCollection"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}

	if url == "/ODIM/v1/Chassis/1/NetworkAdapters/1" {
		body := `{"Name": "Network Adapter View"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}

	return nil, fmt.Errorf("InvalidRequest")
}

func TestPluginContact_GetChassisResource(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	reqData, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Power"})

	err := mockChassisResourceData(reqData, "Power", "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Power")

	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	reqData1, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/NetworkAdapters"})
	err1 := mockChassisResourceData(reqData1, "NetworkAdaptersCollection", "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/NetworkAdapters")
	if err1 != nil {
		t.Fatalf("Error in creating mock resource data :%v", err1)
	}

	reqData2, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/NetworkAdapters/1"})
	err2 := mockChassisResourceData(reqData2, "NetworkAdapters", "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/NetworkAdapters/1")
	if err2 != nil {
		t.Fatalf("Error in creating mock resource data :%v", err2)
	}

	header := map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	errArgs := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "error: SystemUUID not found",
				MessageArgs:   []interface{}{"Chassis", "6d4a0a66-7efa-578e-83cf-44dc68d2874e"},
			},
		},
	}
	errArgs1 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "error while trying to get compute details: no data with the with key 6d4a0a66-7efa-578e-83cf-44dc68d2874e found",
				MessageArgs:   []interface{}{"", "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Power1"},
			},
		},
	}
	pluginContact := PluginContact{
		ContactClient: mockContactClient,
	}

	type args struct {
		req *chassisproto.GetChassisRequest
	}
	tests := []struct {
		name    string
		p       *PluginContact
		args    args
		want    response.RPC
		wantErr bool
	}{
		{
			name: "successful get data",
			p:    &pluginContact,
			args: args{
				req: &chassisproto.GetChassisRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					URL:          "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Power",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          map[string]interface{}{"@odata.id": "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Power"},
			},
			wantErr: false,
		},
		{
			name: "invalid request param",
			p:    &pluginContact,
			args: args{
				req: &chassisproto.GetChassisRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
					URL:          "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Power",
				},
			},
			want: response.RPC{
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Body:          errArgs.CreateGenericErrorResponse(),
			},
			wantErr: true,
		},
		{
			name: "invalid url",
			p:    &pluginContact,
			args: args{
				req: &chassisproto.GetChassisRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					URL:          "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Power1",
				},
			},
			want: response.RPC{
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Body:          errArgs1.CreateGenericErrorResponse(),
			},
			wantErr: true,
		},
		{
			name: "successful get NetworkAdapters collection",
			p:    &pluginContact,
			args: args{
				req: &chassisproto.GetChassisRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					URL:          "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/NetworkAdapters",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          map[string]interface{}{"@odata.id": "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/NetworkAdapters"},
			},
			wantErr: false,
		},

		{
			name: "successful get NetworkAdapter",
			p:    &pluginContact,
			args: args{
				req: &chassisproto.GetChassisRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					URL:          "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/NetworkAdapters/1",
					ResourceID:   "1",
				},
			},
			want: response.RPC{
				Header:        header,
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          map[string]interface{}{"@odata.id": "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/NetworkAdapters/1"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.p.GetChassisResource(tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PluginContact.GetChassisResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
