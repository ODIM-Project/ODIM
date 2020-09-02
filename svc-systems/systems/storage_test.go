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
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
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
func mockSystemsData(uuid string, device smodel.Target) error {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err)
	}
	if err = connPool.Create("System", uuid, device); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "System", err.Error())
	}
	return nil
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
	err = mockSystemsData("54b243cf-f1e3-5319-92d9-2d6737d6b0a", device1)
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}

	var positiveResponse interface{}
	json.Unmarshal([]byte(`{"MessageId": "Base.1.0.Success"}`), &positiveResponse)
	pluginContact := PluginContact{
		ContactClient:  mockContactClient,
		DevicePassword: stubDevicePassword,
	}

	tests := []struct {
		name string
		p    *PluginContact
		req  *systemsproto.CreateVolumeRequest
		want response.RPC
	}{
		{
			name: "Valid request",
			p:    &pluginContact,
			req: &systemsproto.CreateVolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a:1",
				StorageInstance: "ArrayControllers-0",
				RequestBody: []byte(`{"Name":"Volume1",
										"RaidType":"RAID0",
										"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a:1/Storage/ArrayControllers-0/Drives/0"}]}`),
			},
			want: response.RPC{
				StatusCode:    http.StatusNotImplemented, // change it to StatusOK after implementation
				StatusMessage: response.Success,
				Header: map[string]string{
					"Cache-Control":     "no-cache",
					"Connection":        "keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Transfer-Encoding": "chunked",
					"OData-Version":     "4.0",
				},
				Body: map[string]interface{}{"MessageId": "Base.1.0.Success"},
			},
		}, {
			name: "Valid request with multiple drives",
			p:    &pluginContact,
			req: &systemsproto.CreateVolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a:1",
				StorageInstance: "ArrayControllers-0",
				RequestBody: []byte(`{"Name":"Volume1",
										"RaidType":"RAID0",
										"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a:1/Storage/ArrayControllers-0/Drives/0"},{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a:1/Storage/ArrayControllers-0/Drives/1"}]}`),
			},
			want: response.RPC{
				StatusCode:    http.StatusNotImplemented, // change it to StatusOK after implementation
				StatusMessage: response.Success,
				Header: map[string]string{
					"Cache-Control":     "no-cache",
					"Connection":        "keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Transfer-Encoding": "chunked",
					"OData-Version":     "4.0",
				},
				Body: map[string]interface{}{"MessageId": "Base.1.0.Success"},
			},
		}, {
			name: "invalid system id",
			p:    &pluginContact,
			req: &systemsproto.CreateVolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0b:1",
				StorageInstance: "ArrayControllers-0",
				RequestBody: []byte(`{"Name":"Volume1",
										"RaidType":"RAID0",
										"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0b:1/Storage/ArrayControllers-0/Drives/0"}]}`),
			},
			// change it to StatusNotFound after implementation
			want: common.GeneralError(http.StatusNotImplemented, response.ResourceNotFound, "error while trying to get compute details: no data with the with key 54b243cf-f1e3-5319-92d9-2d6737d6b0b found", []interface{}{"ComputerSystem", "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0b:1"}, nil),
		}, {
			name: "invalid storage instance",
			p:    &pluginContact,
			req: &systemsproto.CreateVolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a:1",
				StorageInstance: "ArrayControllers-XYZ",
				RequestBody: []byte(`{"Name":"Volume1",
										"RaidType":"RAID0",
										"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a:1/Storage/ArrayControllers-XYZ/Drives/0"}]}`),
			},
			// change it to StatusNotFound after implementation
			want: common.GeneralError(http.StatusNotImplemented, response.ResourceNotFound, "error: Storage instance not found", []interface{}{"System", "54b243cf-f1e3-5319-92d9-2d6737d6b0a"}, nil),
		}, {
			name: "invalid RaidType",
			p:    &pluginContact,
			req: &systemsproto.CreateVolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a:1",
				StorageInstance: "ArrayControllers-0",
				RequestBody: []byte(`{"Name":"Volume1",
										"RaidType":"Invalid",
										"Drives":[{"@odata.id": "/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a:1/Storage/ArrayControllers-0/Drives/0"}]}`),
			},
			// change it to StatusBadRequest after implementation
			want: common.GeneralError(http.StatusNotImplemented, response.ResourceNotFound, "error: SystemUUID not found", []interface{}{"System", "54b243cf-f1e3-5319-92d9-2d6737d6b0a"}, nil),
		}, {
			name: "Invalid Drives format",
			p:    &pluginContact,
			req: &systemsproto.CreateVolumeRequest{
				SystemID:        "54b243cf-f1e3-5319-92d9-2d6737d6b0a:1",
				StorageInstance: "ArrayControllers-0",
				RequestBody: []byte(`{"Name":"Volume1",
										"RaidType":"Invalid",
										"Drives":["/redfish/v1/Systems/54b243cf-f1e3-5319-92d9-2d6737d6b0a:1/Storage/ArrayControllers-0/Drives/0"]}`),
			},
			// change it to StatusBadRequest after implementation
			want: common.GeneralError(http.StatusNotImplemented, response.ResourceNotFound, "error: ", []interface{}{"System", "54b243cf-f1e3-5319-92d9-2d6737d6b0a"}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//TODO: after implementation of logic change the statuscode match to entire response match
			if got := tt.p.CreateVolume(tt.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("PluginContact.CreateVolume() = %v, want %v", got.StatusCode, tt.want.StatusCode)
			}
		})
	}
}
