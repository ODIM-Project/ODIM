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
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
)

func mockPluginData(t *testing.T) error {
	password := getEncryptedKey(t, []byte("$2a$10$OgSUYvuYdI/7dLL5KkYNp.RCXISefftdj.MjbBTr6vWyNwAvht6ci"))
	plugin := smodel.Plugin{
		IP:                "localhost",
		Port:              "9091",
		Username:          "admin",
		Password:          password,
		ID:                "GRF",
		PreferredAuthType: "BasicAuth",
	}
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err)
	}
	// Clear previously created key, if exists
	connPool.Delete("Plugin", "GRF")
	if err = connPool.Create("Plugin", "GRF", plugin); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "Plugin", err.Error())
	}
	return nil
}
func mockDeviceData(uuid string, device smodel.Target) error {

	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err)
	}
	if err = connPool.Create("System", uuid, device); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "System", err.Error())
	}
	return nil
}

func mockIsAuthorized(sessionToken string, privileges, oemPrivileges []string) response.RPC {
	if sessionToken != "validToken" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "error while trying to authenticate session", nil, nil)
	}
	return common.GeneralError(http.StatusOK, response.Success, "", nil, nil)
}

func mockContactClient(url, method, token string, odataID string, body interface{}, basicAuth map[string]string) (*http.Response, error) {
	if url == "https://localhost:9091/ODIM/v1/Systems/1/Actions/ComputerSystem.Reset" {
		body := `{"MessageId": "` + response.Success + `"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "https://localhost:9091/ODIM/v1/Systems/1/Actions/ComputerSystem.SetDefaultBootOrder" {
		body := `{"MessageId": "` + response.Success + `"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "https://localhost:9091/ODIM/v1/Systems/1/SecureBoot1" {
		body := `{"@odata.id": "/ODIM/v1/Systems/1/SecureBoot1"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}

	if url == "https://localhost:9098/ODIM/v1/Systems/1/Storage" {
		body := `{"@odata.id": "/ODIM/v1/Systems/1/Storage"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "https://localhost:9091/ODIM/v1/Systems/1/Storage" {
		body := `{"@odata.id": "/ODIM/v1/Systems/1/Storage"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "https://localhost:9091/ODIM/v1/Systems/1/Bios/Settings" {
		body := `{"@odata.id": "/ODIM/v1/Systems/1/Bios/Settings"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "https://localhost:9091/ODIM/v1/Systems/1" {
		body := `{"@odata.id": "/ODIM/v1/Systems/1"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

func TestPluginContact_ComputerSystemReset(t *testing.T) {
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
		ManagerAddress: "10.24.0.12",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "7a2c6100-67da-5fd6-ab82-6870d29c727",
		PluginID:       "GR",
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
	err = mockDeviceData("7a2c6100-67da-5fd6-ab82-6870d29c727", device3)
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	err = mockDeviceData("24b243cf-f1e3-5318-92d9-2d6737d6b0b9", device1)
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	err = mockSystemData("/redfish/v1/Systems/7a2c6100-67da-5fd6-ab82-6870d29c7279:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	err = mockSystemData("/redfish/v1/Systems/7a2c6100-67da-5fd6-ab82-6870d29c727:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	var positiveResponse interface{}
	json.Unmarshal([]byte(`{"MessageId": "`+response.Success+`"}`), &positiveResponse)
	pluginContact := PluginContact{
		ContactClient:  mockContactClient,
		DevicePassword: stubDevicePassword,
	}
	type args struct {
		req *systemsproto.ComputerSystemResetRequest
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
				&systemsproto.ComputerSystemResetRequest{
					SystemID:    "24b243cf-f1e3-5318-92d9-2d6737d6b0b:1",
					RequestBody: []byte(`{"ResetType": "ForceRestart"}`),
				},
			},
			want: common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "error while trying to get compute details: no data with the with key 24b243cf-f1e3-5318-92d9-2d6737d6b0b found", []interface{}{"ComputerSystem", "/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b:1"}, nil),
		}, {
			name: "invalid uuid without system id",
			p:    &pluginContact,
			args: args{
				&systemsproto.ComputerSystemResetRequest{
					SystemID:    "24b243cf-f1e3-5318-92d9-2d6737d6b0b",
					RequestBody: []byte(`{"ResetType": "ForceRestart"}`),
				},
			},
			want: common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "error: SystemUUID not found", []interface{}{"System", "24b243cf-f1e3-5318-92d9-2d6737d6b0b"}, nil),
		},
		{
			name: "if plugin id is not there in the db",
			p:    &pluginContact,
			args: args{
				&systemsproto.ComputerSystemResetRequest{
					SystemID:    "7a2c6100-67da-5fd6-ab82-6870d29c727:1",
					RequestBody: []byte(`{"ResetType": "ForceRestart"}`),
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusInternalServerError,
				StatusMessage: response.InternalError,
				Header:        map[string]string{"Content-type": "application/json; charset=utf-8"},
				Body:          errArg3.CreateGenericErrorResponse(),
			},
		},
		// TODO: check this test case
		{
			name: "Valid Request",
			p:    &pluginContact,
			args: args{
				&systemsproto.ComputerSystemResetRequest{
					SystemID:    "7a2c6100-67da-5fd6-ab82-6870d29c7279:1",
					RequestBody: []byte(`{"ResetType": "ForceRestart"}`),
				},
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
			if got := tt.p.ComputerSystemReset(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PluginContact.ComputerSystemReset() = %v, want %v", got, tt.want)
			}
		})
	}
}
