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
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
)

func mockPluginClientData(t *testing.T) error {
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

func contactPluginClient(url, method, token string, odataID string, body interface{}, basicAuth map[string]string) (*http.Response, error) {
	if url == "https://localhost:9091/ODIM/v1/Systems/1/Storage/ArrayControllers-0/Volumes" {
		body := `{"MessageId": "` + response.Success + `"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "https://localhost:9091/ODIM/v1/Systems/1/Storage/1/Volumes/1" {
		body := `{"MessageId": "` + response.Success + `"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "https://localhost:9091/ODIM/v1/Systems/1/Storage/1/Volumes/2" {
		body := `{"MessageId": "` + response.Failure + `"}`
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

func mockGetExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		ContactClient:  contactPluginClient,
		DevicePassword: stubDevicePassword,
		DB: DB{
			GetResource: mockGetResource,
		},
		GetPluginStatus: mockPluginStatus,
	}
}

func mockGetResource(table, key string) (string, *errors.Error) {
	if key == "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a:1/Storage/ArrayControllers-0/Drives/0" {
		return "", nil
	}
	if key == "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a:1/Storage/ArrayControllers-0/Drives/1" {
		return "", nil
	}

	return "body", nil
}

func mockPluginStatus(plugin smodel.Plugin) bool {
	return true
}

func TestPluginContact_CreateVolume(t *testing.T) {
	// Modify the contents with http.StatusNotImplemented to the correct status
	// and modify all other info accordingly after implementations
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
		ManagerAddress: "10.24.0.12",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "54b243cf-f1e3-5319-92d9-2d6737d6b0a",
		PluginID:       "GRF",
	}

	err := mockPluginClientData(t)
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	err = mockDeviceData("54b243cf-f1e3-5319-92d9-2d6737d6b0a", device1)
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}

	var positiveResponse interface{}
	json.Unmarshal([]byte(`{"MessageId": "`+response.Success+`"}`), &positiveResponse)
	pluginContact := mockGetExternalInterface()

	tests := []struct {
		name string
		p    *ExternalInterface
		req  *systemsproto.VolumeRequest
		want response.RPC
	}{
		{
			name: "Valid request",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a:1",
				StorageInstance: "ArrayControllers-0",
				RequestBody: []byte(`{"Name":"Volume1",
										"RAIDType":"RAID0",
										"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a:1/Storage/ArrayControllers-0/Drives/0"}]
										,"@Redfish.OperationApplyTime": "OnReset"
										}`),
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: map[string]interface{}{"MessageId": response.Success},
			},
		}, {
			name: "Valid request with multiple drives",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a:1",
				StorageInstance: "ArrayControllers-0",
				RequestBody: []byte(`{"Name":"Volume1",
										"RAIDType":"RAID0",
										"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a:1/Storage/ArrayControllers-0/Drives/0"},{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a:1/Storage/ArrayControllers-0/Drives/1"}]}`),
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: map[string]interface{}{"MessageId": response.Success},
			},
		}, {
			name: "invalid system id",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0b:1",
				StorageInstance: "ArrayControllers-0",
				RequestBody: []byte(`{"Name":"Volume1",
										"RAIDType":"RAID0",
										"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0b:1/Storage/ArrayControllers-0/Drives/0"}]}`),
			},
			want: common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "error while trying to get compute details: no data with the with key 54b243cf-f1e3-5319-92d9-2d6737d6b0b found", []interface{}{"System", "54b243cf-f1e3-5319-92d9-2d6737d6b0b"}, nil),
		}, {
			name: "invalid storage instance",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a:1",
				StorageInstance: "",
				RequestBody: []byte(`{"Name":"Volume1",
										"RAIDType":"RAID0",
										"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a:1/Storage/ArrayControllers-0/Drives/0"}]}`),
			},
			want: common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, "error: Storage instance is not found", []interface{}{"Storage", ""}, nil),
		}, {
			name: "invalid RaidType",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a:1",
				StorageInstance: "ArrayControllers-0",
				RequestBody: []byte(`{"Name":"Volume1",
										"RAIDType":"Invalid",
										"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a:1/Storage/ArrayControllers-0/Drives/0"}]}`),
			},
			want: common.GeneralError(http.StatusBadRequest, response.PropertyValueNotInList, "error: request payload validation failed: RAIDType Invalid is invalid", []interface{}{"Invalid", "RAIDType"}, nil),
		}, {
			name: "Invalid Drives format",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a:1",
				StorageInstance: "ArrayControllers-0",
				RequestBody: []byte(`{"Name":"Volume1",
										"RaidType":"Invalid",
										"Drives":["/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a:1/Storage/ArrayControllers-0/Drives/12"]`),
			},
			want: common.GeneralError(http.StatusBadRequest, response.MalformedJSON, "Error while unmarshaling the create volume request: unexpected end of JSON input", []interface{}{}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.CreateVolume(tt.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PluginContact.CreateVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluginContact_DeleteVolume(t *testing.T) {
	// Modify the contents with http.StatusNotImplemented to the correct status
	// and modify all other info accordingly after implementations
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
		ManagerAddress: "10.24.0.12",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "54b243cf-f1e3-5319-92d9-2d6737d6b0a",
		PluginID:       "GRF",
	}
	device2 := smodel.Target{
		ManagerAddress: "100.0.0.2",
		Password:       []byte("someValidPassword"),
		UserName:       "admin",
		DeviceUUID:     "8e896459-a8f9-4c83-95b7-7b316b4908e1",
		PluginID:       "Unknown-Plugin",
	}
	err := mockPluginClientData(t)
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	mockDeviceData("54b243cf-f1e3-5319-92d9-2d6737d6b0a", device1)
	mockDeviceData("8e896459-a8f9-4c83-95b7-7b316b4908e1", device2)
	mockSystemData("/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a:1")
	mockSystemData("/redfish/v1/Systems/8e896459-a8f9-4c83-95b7-7b316b4908e1:1")
	var reqData = `{"@odata.id":"/redfish/v1/Systems/1/Storage/1/Volumes/1"}`
	mockSystemResourceData([]byte(reqData), "Volumes", "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a:1/Storage/1/Volumes/1")
	mockSystemResourceData([]byte(reqData), "Volumes", "/redfish/v1/Systems/8e896459-a8f9-4c83-95b7-7b316b4908e1:1/Storage/1/Volumes/1")

	var positiveResponse interface{}
	json.Unmarshal([]byte(`{"MessageId": "`+response.Success+`"}`), &positiveResponse)
	pluginContact := mockGetExternalInterface()

	tests := []struct {
		name           string
		p              *ExternalInterface
		req            *systemsproto.VolumeRequest
		wantStatusCode int32
	}{
		{
			name: "Valid request",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a:1",
				StorageInstance: "1",
				VolumeID:        "1",
				RequestBody:     []byte(`{"@Redfish.OperationApplyTime": "OnReset"}`),
			},
			wantStatusCode: http.StatusNoContent,
		},
		{
			name: "invalid system id",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0b:1",
				StorageInstance: "1",
				VolumeID:        "1",
				RequestBody:     []byte(`{"@Redfish.OperationApplyTime": "OnReset"}`),
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "without system id",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0b:",
				StorageInstance: "1",
				VolumeID:        "1",
				RequestBody:     []byte(`{"@Redfish.OperationApplyTime": "OnReset"}`),
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Invalid volume id",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a:1",
				StorageInstance: "1",
				VolumeID:        "2",
				RequestBody:     []byte(`{"@Redfish.OperationApplyTime": "OnReset"}`),
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "unknown plugin",
			p:    pluginContact,
			req: &systemsproto.VolumeRequest{
				SystemID:        "8e896459-a8f9-4c83-95b7-7b316b4908e1:1",
				StorageInstance: "1",
				VolumeID:        "1",
				RequestBody:     []byte(`{"@Redfish.OperationApplyTime": "OnReset"}`),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.DeleteVolume(tt.req); got.StatusCode != tt.wantStatusCode {
				t.Errorf("PluginContact.DeleteVolume() = %v, want %v", got.StatusCode, tt.wantStatusCode)
			}
		})
	}
}
