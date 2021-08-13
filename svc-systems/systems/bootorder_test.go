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
package systems

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"

	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
)

func mockSystemData(systemID string) error {
	reqData, _ := json.Marshal(&map[string]interface{}{
		"Manufacturer": "HPE",
		"Id":           "1",
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

func stubDevicePassword(password []byte) ([]byte, error) {
	return password, nil
}

func TestPluginContact_SetDefaultBootOrder(t *testing.T) {
	config.SetUpMockConfig(t)
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

	device1 := smodel.Target{
		ManagerAddress: "10.24.0.13",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "24b243cf-f1e3-5318-92d9-2d6737d6b0b9",
		PluginID:       "GRF",
	}
	device2 := smodel.Target{
		ManagerAddress: "10.24.0.12",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "7a2c6100-67da-5fd6-ab82-6870d29c7279",
		PluginID:       "GRF",
	}
	device3 := smodel.Target{
		ManagerAddress: "10.24.0.14",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "7ff3bd97-c41c-5de0-937d-85d390691b73",
	}
	errArg1 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "error: SystemUUID not found",
				MessageArgs:   []interface{}{"System", "24b243cf-f1e3-5318-92d9-2d6737d6b0b"},
			},
		},
	}
	errArg2 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "error while trying to get compute details: no data with the with key 24b243cf-f1e3-5318-92d9-2d6737d6b0b found",
				MessageArgs:   []interface{}{"System", "24b243cf-f1e3-5318-92d9-2d6737d6b0b"},
			},
		},
	}
	errArg3 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InternalError,
				ErrorMessage:  "error while trying to get plugin details",
			},
		},
	}
	err := mockPluginData(t)
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	err = mockDeviceData("7a2c6100-67da-5fd6-ab82-6870d29c7279", device2)
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	err = mockDeviceData("24b243cf-f1e3-5318-92d9-2d6737d6b0b9", device1)
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	err = mockDeviceData("7ff3bd97-c41c-5de0-937d-85d390691b73", device3)
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	err = mockSystemData("/redfish/v1/Systems/7a2c6100-67da-5fd6-ab82-6870d29c7279:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	err = mockSystemData("/redfish/v1/Systems/7ff3bd97-c41c-5de0-937d-85d390691b73:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	pluginContact := PluginContact{
		ContactClient:  mockContactClient,
		DevicePassword: stubDevicePassword,
	}

	type args struct {
		systemID string
	}
	tests := []struct {
		name string
		p    *PluginContact
		args args
		want response.RPC
	}{
		{
			name: "invalid uuid",
			p:    &pluginContact,
			args: args{
				systemID: "24b243cf-f1e3-5318-92d9-2d6737d6b0b:1",
			},
			want: response.RPC{
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Header:        map[string]string{"Content-type": "application/json; charset=utf-8"},
				Body:          errArg2.CreateGenericErrorResponse(),
			},
		}, {
			name: "invalid uuid without system id",
			p:    &pluginContact,
			args: args{
				systemID: "24b243cf-f1e3-5318-92d9-2d6737d6b0b",
			},
			want: response.RPC{
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArg1.CreateGenericErrorResponse(),
			},
		},
		{
			name: "if plugin id doesn't there in db",
			p:    &pluginContact,
			args: args{
				systemID: "7ff3bd97-c41c-5de0-937d-85d390691b73:1",
			},
			want: response.RPC{
				StatusCode:    http.StatusInternalServerError,
				StatusMessage: response.InternalError,
				Header:        map[string]string{"Content-type": "application/json; charset=utf-8"},
				Body:          errArg3.CreateGenericErrorResponse(),
			},
		},
		{
			name: "Valid Request",
			p:    &pluginContact,
			args: args{
				systemID: "7a2c6100-67da-5fd6-ab82-6870d29c7279:1",
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Cache-Control":     "no-cache",
					"Connection":        "keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Transfer-Encoding": "chunked",
					"OData-Version":     "4.0",
				},
				Body: map[string]interface{}{"MessageId": response.Success},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.SetDefaultBootOrder(tt.args.systemID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PluginContact.SetDefaultBootOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluginContact_ChangeBiosSettings(t *testing.T) {
	config.SetUpMockConfig(t)
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

	device1 := smodel.Target{
		ManagerAddress: "10.24.0.13",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "24b243cf-f1e3-5318-92d9-2d6737d6b0b9",
		PluginID:       "GRF",
	}
	device2 := smodel.Target{
		ManagerAddress: "10.24.0.12",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "7a2c6100-67da-5fd6-ab82-6870d29c7279",
		PluginID:       "GRF",
	}
	device3 := smodel.Target{
		ManagerAddress: "10.24.0.14",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "7ff3bd97-c41c-5de0-937d-85d390691b73",
	}
	err := mockPluginData(t)
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	err = mockDeviceData("24b243cf-f1e3-5318-92d9-2d6737d6b0b9", device1)
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	err = mockDeviceData("7ff3bd97-c41c-5de0-937d-85d390691b73", device3)
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	err = mockSystemData("/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b9:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	err = mockSystemData("/redfish/v1/Systems/7ff3bd97-c41c-5de0-937d-85d390691b73:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	err = mockDeviceData("7a2c6100-67da-5fd6-ab82-6870d29c7279", device2)
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	err = mockSystemData("/redfish/v1/Systems/7a2c6100-67da-5fd6-ab82-6870d29c7279:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	pluginContact := PluginContact{
		ContactClient:  mockContactClient,
		DevicePassword: stubDevicePassword,
	}

	errArg1 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "error: SystemUUID not found",
				MessageArgs:   []interface{}{"System", "24b243cf-f1e3-5318-92d9-2d6737d6b0b"},
			},
		},
	}
	errArg2 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "error while trying to get compute details: no data with the with key 24b243cf-f1e3-5318-92d9-2d6737d6b0b found",
				MessageArgs:   []interface{}{"System", "24b243cf-f1e3-5318-92d9-2d6737d6b0b"},
			},
		},
	}
	errArg3 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InternalError,
				ErrorMessage:  "error while trying to get plugin details",
			},
		},
	}
	request := []byte(`{"Attributes": {"BootMode": "mode"}}`)
	tests := []struct {
		name string
		p    *PluginContact
		req  *systemsproto.BiosSettingsRequest
		want response.RPC
	}{
		{
			name: "uuid without system id",
			p:    &pluginContact,
			req: &systemsproto.BiosSettingsRequest{
				SystemID:     "24b243cf-f1e3-5318-92d9-2d6737d6b0b",
				RequestBody:  request,
				SessionToken: "token",
			},
			want: response.RPC{
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArg1.CreateGenericErrorResponse(),
			},
		}, {
			name: "invalid uuid",
			p:    &pluginContact,
			req: &systemsproto.BiosSettingsRequest{
				SystemID:     "24b243cf-f1e3-5318-92d9-2d6737d6b0b:1",
				RequestBody:  request,
				SessionToken: "token",
			},
			want: response.RPC{
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Header:        map[string]string{"Content-type": "application/json; charset=utf-8"},
				Body:          errArg2.CreateGenericErrorResponse(),
			},
		},
		{
			name: "if plugin id doesn't there in db",
			p:    &pluginContact,
			req: &systemsproto.BiosSettingsRequest{
				SystemID:     "7ff3bd97-c41c-5de0-937d-85d390691b73:1",
				RequestBody:  request,
				SessionToken: "token",
			},
			want: response.RPC{
				StatusCode:    http.StatusInternalServerError,
				StatusMessage: response.InternalError,
				Header:        map[string]string{"Content-type": "application/json; charset=utf-8"},
				Body:          errArg3.CreateGenericErrorResponse(),
			},
		},
		{
			name: "Valid Request",
			p:    &pluginContact,
			req: &systemsproto.BiosSettingsRequest{
				SystemID:     "7a2c6100-67da-5fd6-ab82-6870d29c7279:1",
				RequestBody:  request,
				SessionToken: "token",
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: map[string]interface{}{"@odata.id": "/redfish/v1/Systems/1/Bios/Settings"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ChangeBiosSettings(tt.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PluginContact.ChangeBiosSettings() = %v, want %v", got, tt.want)
			}
		})
	}
}

// this client is for plugin login returns an error
func mockloginClient(url, method, token string, odataID string, body interface{}, basicAuth map[string]string) (*http.Response, error) {
	if url == "http://localhost:9091/redfishplugin/login" {
		header := make(http.Header)
		header.Set("X-Auth-Token", token)
		body := `{"MessageId": "` + response.Failure + `"}`
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Header:     header,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

// this client is for change bios settings return error returns an error
func mockPluginClient(url, method, token string, odataID string, body interface{}, basicAuth map[string]string) (*http.Response, error) {

	if url == "http://localhost:9091/ODIM/v1/systems/1/bios/settings" {
		body := `{"MessageId": "` + response.Failure + `"}`
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "http://localhost:9091/ODIM/v1/Systems/1" {
		body := `{"MessageId": "` + response.Failure + `"}`
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

func TestPluginContact_ChangeBootOrderSettings(t *testing.T) {
	config.SetUpMockConfig(t)
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

	device1 := smodel.Target{
		ManagerAddress: "10.24.0.13",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "24b243cf-f1e3-5318-92d9-2d6737d6b0b9",
		PluginID:       "GRF",
	}
	device2 := smodel.Target{
		ManagerAddress: "10.24.0.12",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "7a2c6100-67da-5fd6-ab82-6870d29c7279",
		PluginID:       "GRF",
	}
	device3 := smodel.Target{
		ManagerAddress: "10.24.0.14",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "7ff3bd97-c41c-5de0-937d-85d390691b73",
	}
	err := mockPluginData(t)
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	err = mockDeviceData("7a2c6100-67da-5fd6-ab82-6870d29c7279", device2)
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	err = mockDeviceData("24b243cf-f1e3-5318-92d9-2d6737d6b0b9", device1)
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	err = mockDeviceData("7ff3bd97-c41c-5de0-937d-85d390691b73", device3)
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	err = mockSystemData("/redfish/v1/Systems/7ff3bd97-c41c-5de0-937d-85d390691b73:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	err = mockSystemData("/redfish/v1/Systems/7a2c6100-67da-5fd6-ab82-6870d29c7279:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	pluginContact := PluginContact{
		ContactClient:  mockContactClient,
		DevicePassword: stubDevicePassword,
	}
	errArg1 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "error: SystemUUID not found",
				MessageArgs:   []interface{}{"System", "24b243cf-f1e3-5318-92d9-2d6737d6b0b"},
			},
		},
	}
	errArg2 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "error while trying to get compute details: no data with the with key 24b243cf-f1e3-5318-92d9-2d6737d6b0b found",
				MessageArgs:   []interface{}{"System", "24b243cf-f1e3-5318-92d9-2d6737d6b0b"},
			},
		},
	}
	errArg3 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InternalError,
				ErrorMessage:  "error while trying to get plugin details",
			},
		},
	}

	request := []byte(`{"Attributes": {"BootMode": "mode"}}`)
	tests := []struct {
		name string
		p    *PluginContact
		req  *systemsproto.BootOrderSettingsRequest
		want response.RPC
	}{
		{
			name: "uuid without system id",
			p:    &pluginContact,
			req: &systemsproto.BootOrderSettingsRequest{
				SystemID:     "24b243cf-f1e3-5318-92d9-2d6737d6b0b",
				RequestBody:  request,
				SessionToken: "token",
			},
			want: response.RPC{
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArg1.CreateGenericErrorResponse(),
			},
		}, {
			name: "invalid uuid",
			p:    &pluginContact,
			req: &systemsproto.BootOrderSettingsRequest{
				SystemID:     "24b243cf-f1e3-5318-92d9-2d6737d6b0b:1",
				RequestBody:  request,
				SessionToken: "token",
			},
			want: response.RPC{
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Header:        map[string]string{"Content-type": "application/json; charset=utf-8"},
				Body:          errArg2.CreateGenericErrorResponse(),
			},
		},
		{
			name: "if plugin id doesn't there in db",
			p:    &pluginContact,
			req: &systemsproto.BootOrderSettingsRequest{
				SystemID:     "7ff3bd97-c41c-5de0-937d-85d390691b73:1",
				RequestBody:  request,
				SessionToken: "token",
			},
			want: response.RPC{
				StatusCode:    http.StatusInternalServerError,
				StatusMessage: response.InternalError,
				Header:        map[string]string{"Content-type": "application/json; charset=utf-8"},
				Body:          errArg3.CreateGenericErrorResponse(),
			},
		},
		{
			name: "Valid Request",
			p:    &pluginContact,
			req: &systemsproto.BootOrderSettingsRequest{
				SystemID:     "7a2c6100-67da-5fd6-ab82-6870d29c7279:1",
				RequestBody:  request,
				SessionToken: "token",
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Cache-Control":     "no-cache",
					"Connection":        "keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Transfer-Encoding": "chunked",
					"OData-Version":     "4.0",
				},
				Body: map[string]interface{}{"@odata.id": "/redfish/v1/Systems/1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ChangeBootOrderSettings(tt.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PluginContact.ChangeBootOrderSettings() = %v, want %v", got, tt.want)
			}
		})
	}
}
