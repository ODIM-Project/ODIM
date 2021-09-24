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
package update

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-update/ucommon"
	"github.com/ODIM-Project/ODIM/svc-update/umodel"
	"github.com/ODIM-Project/ODIM/svc-update/uresponse"
	"github.com/stretchr/testify/assert"
)

func mockIsAuthorized(sessionToken string, privileges, oemPrivileges []string) response.RPC {
	if sessionToken != "validToken" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "error while trying to authenticate session", nil, nil)
	}
	return common.GeneralError(http.StatusOK, response.Success, "", nil, nil)
}

func mockContactClient(url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {
	return nil, fmt.Errorf("InvalidRequest")
}

func mockGetResource(table, key string, dbType common.DbType) (string, *errors.Error) {
	if (key == "/redfish/v1/UpdateService/FirmwareInentory/3bd1f589-117a-4cf9-89f2-da44ee8e012b:1") || (key == "/redfish/v1/UpdateService/SoftwareInentory/3bd1f589-117a-4cf9-89f2-da44ee8e012b:1") {
		return "", errors.PackError(errors.DBKeyNotFound, "not found")
	}
	return "body", nil
}

func mockGetAllKeysFromTable(table string, dbType common.DbType) ([]string, error) {
	return []string{"/redfish/v1/UpdateService/FirmwareInentory/uuid:1"}, nil
}
func mockGetTarget(id string) (*umodel.Target, *errors.Error) {
	var target umodel.Target
	target.PluginID = id
	target.DeviceUUID = "uuid"
	target.UserName = "admin"
	target.Password = []byte("password")
	target.ManagerAddress = "ip"
	return &target, nil
}

func mockGetPluginData(id string) (umodel.Plugin, *errors.Error) {
	var plugin umodel.Plugin
	plugin.IP = "ip"
	plugin.Port = "port"
	plugin.Username = "plugin"
	plugin.Password = []byte("password")
	plugin.PluginType = "Redfish"
	plugin.PreferredAuthType = "basic"
	return plugin, nil
}

func mockContactPlugin(req ucommon.PluginContactRequest, errorMessage string) ([]byte, string, ucommon.ResponseStatus, error) {
	var responseStatus ucommon.ResponseStatus

	return []byte(`{"Attributes":"sample"}`), "token", responseStatus, nil
}

func stubDevicePassword(password []byte) ([]byte, error) {
	return password, nil
}

func stubGenericSave(reqBody []byte, table string, uuid string) error {
	return nil
}

func mockGetExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		External: External{
			Auth:            mockIsAuthorized,
			ContactClient:   mockContactClient,
			GetTarget:       mockGetTarget,
			GetPluginData:   mockGetPluginData,
			ContactPlugin:   mockContactPlugin,
			DevicePassword:  stubDevicePassword,
			CreateChildTask: mockCreateChildTask,
			UpdateTask:      mockUpdateTask,
			GenericSave:     stubGenericSave,
		},
		DB: DB{
			GetAllKeysFromTable: mockGetAllKeysFromTable,
			GetResource:         mockGetResource,
		},
	}
}

func TestGetUpdateService(t *testing.T) {
	successResponse := response.Response{
		OdataType:    common.UpdateServiceType,
		OdataID:      "/redfish/v1/UpdateService",
		OdataContext: "/redfish/v1/$metadata#UpdateService.UpdateService",
		ID:           "UpdateService",
		Name:         "Update Service",
	}
	successResponse.CreateGenericResponse(response.Success)
	successResponse.Message = ""
	successResponse.MessageID = ""
	successResponse.Severity = ""
	common.SetUpMockConfig()
	tests := []struct {
		name string
		want response.RPC
	}{
		{
			name: "account service enabled",
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Allow":         "GET",
					"Cache-Control": "no-cache",
					"Connection":    "Keep-alive",
					"Content-type":  "application/json; charset=utf-8",
					"Link": "	</redfish/v1/SchemaStore/en/UpdateService.json>; rel=describedby",
					"Transfer-Encoding": "chunked",
					"X-Frame-Options":   "sameorigin",
				},
				Body: uresponse.UpdateService{
					Response: successResponse,
					Status: uresponse.Status{
						State:        "Enabled",
						Health:       "OK",
						HealthRollup: "OK",
					},
					ServiceEnabled: true,
					SoftwareInventory: uresponse.SoftwareInventory{
						OdataID: "/redfish/v1/UpdateService/SoftwareInventory",
					},
					FirmwareInventory: uresponse.FirmwareInventory{
						OdataID: "/redfish/v1/UpdateService/FirmwareInventory",
					},
					Actions: uresponse.Actions{
						UpdateServiceSimpleUpdate: uresponse.UpdateServiceSimpleUpdate{
							Target: "/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate",
							RedfishOperationApplyTimeSupport: uresponse.RedfishOperationApplyTimeSupport{
								OdataType:       common.SettingsType,
								SupportedValues: []string{"OnStartUpdateRequest"},
							},
						},
						UpdateServiceStartUpdate: uresponse.UpdateServiceStartUpdate{
							Target: "/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate",
						},
					},
				},
			},
		},
		{
			name: "account service disabled",
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Allow":         "GET",
					"Cache-Control": "no-cache",
					"Connection":    "Keep-alive",
					"Content-type":  "application/json; charset=utf-8",
					"Link": "	</redfish/v1/SchemaStore/en/UpdateService.json>; rel=describedby",
					"Transfer-Encoding": "chunked",
					"X-Frame-Options":   "sameorigin",
				},
				Body: uresponse.UpdateService{
					Response: successResponse,
					Status: uresponse.Status{
						State:        "Disabled",
						Health:       "OK",
						HealthRollup: "OK",
					},
					ServiceEnabled: false,
					SoftwareInventory: uresponse.SoftwareInventory{
						OdataID: "/redfish/v1/UpdateService/SoftwareInventory",
					},
					FirmwareInventory: uresponse.FirmwareInventory{
						OdataID: "/redfish/v1/UpdateService/FirmwareInventory",
					},
					Actions: uresponse.Actions{
						UpdateServiceSimpleUpdate: uresponse.UpdateServiceSimpleUpdate{
							Target: "/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate",
							RedfishOperationApplyTimeSupport: uresponse.RedfishOperationApplyTimeSupport{
								OdataType:       common.SettingsType,
								SupportedValues: []string{"OnStartUpdateRequest"},
							},
						},
						UpdateServiceStartUpdate: uresponse.UpdateServiceStartUpdate{
							Target: "/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate",
						},
					},
				},
			},
		},
	}
	config.Data.EnabledServices = []string{"UpdateService"}
	e := mockGetExternalInterface()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := e.GetUpdateService()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUpdateService() = %v, want %v", got, tt.want)
			}
		})
		config.Data.EnabledServices = []string{"XXXX"}
	}
}

func TestFirmwareInventoryCollection(t *testing.T) {
	req := &updateproto.UpdateRequest{}
	e := mockGetExternalInterface()
	response := e.GetAllFirmwareInventory(req)

	update := response.Body.(uresponse.Collection)
	assert.Equal(t, int(response.StatusCode), http.StatusOK, "Status code should be StatusOK.")
	assert.Equal(t, update.MembersCount, 1, "Member count does not match")
}

func TestSoftwareInventoryCollection(t *testing.T) {
	req := &updateproto.UpdateRequest{}
	e := mockGetExternalInterface()
	response := e.GetAllSoftwareInventory(req)

	update := response.Body.(uresponse.Collection)
	assert.Equal(t, int(response.StatusCode), http.StatusOK, "Status code should be StatusOK.")
	assert.Equal(t, update.MembersCount, 1, "Member count does not match")
}

func TestFirmwareInventory(t *testing.T) {
	config.SetUpMockConfig(t)
	req := &updateproto.UpdateRequest{
		ResourceID: "3bd1f589-117a-4cf9-89f2-da44ee8e012b:1",
	}
	e := mockGetExternalInterface()
	response := e.GetFirmwareInventory(req)

	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

}

func TestGetFirmwareInventoryInvalidID(t *testing.T) {
	req := &updateproto.UpdateRequest{
		ResourceID: "invalidID",
	}
	e := mockGetExternalInterface()
	response := e.GetFirmwareInventory(req)

	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusNotFound")
}

func TestSoftwareInventory(t *testing.T) {
	config.SetUpMockConfig(t)
	req := &updateproto.UpdateRequest{
		ResourceID: "3bd1f589-117a-4cf9-89f2-da44ee8e012b:1",
	}
	e := mockGetExternalInterface()
	response := e.GetSoftwareInventory(req)

	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

}

func TestGetSoftwareInventoryInvalidID(t *testing.T) {
	req := &updateproto.UpdateRequest{
		ResourceID: "invalidID",
	}
	e := mockGetExternalInterface()
	response := e.GetSoftwareInventory(req)

	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusNotFound")
}
