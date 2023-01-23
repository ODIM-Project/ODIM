//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

package licenses

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-licenses/model"
)

func TestGetExternalInterface(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetExternalInterface(); got == nil {
				t.Errorf("Result of GetExternalInterface() should not be equal to nil")
			}
		})
	}
}

func mockGetExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		External: External{
			Auth:           mockIsAuthorized,
			ContactClient:  mockContactClient,
			GetTarget:      mockGetTarget,
			GetPluginData:  mockGetPluginData,
			ContactPlugin:  mockContactPlugin,
			DevicePassword: stubDevicePassword,
			GenericSave:    stubGenericSave,
		},
		DB: DB{
			GetAllKeysFromTable: mockGetAllKeysFromTable,
			GetResource:         mockGetResource,
		},
	}
}

func mockContactPlugin(req model.PluginContactRequest, errorMessage string) ([]byte, string, model.ResponseStatus, error) {
	var responseStatus model.ResponseStatus

	return []byte(`{"Attributes":"sample"}`), "token", responseStatus, nil
}

func stubGenericSave(reqBody []byte, table string, uuid string) error {
	return nil
}

func stubDevicePassword(password []byte) ([]byte, error) {
	return password, nil
}

func mockGetTarget(id string) (*model.Target, *errors.Error) {
	var target model.Target
	target.PluginID = id
	target.DeviceUUID = "uuid"
	target.UserName = "admin"
	target.Password = []byte("password")
	target.ManagerAddress = "ip"
	return &target, nil
}

func mockGetPluginData(id string) (*model.Plugin, *errors.Error) {
	var plugin model.Plugin
	plugin.IP = "ip"
	plugin.Port = "port"
	plugin.Username = "plugin"
	plugin.Password = []byte("password")
	plugin.PluginType = "Redfish"
	plugin.PreferredAuthType = "basic"
	return &plugin, nil
}

func mockIsAuthorized(sessionToken string, privileges, oemPrivileges []string) (response.RPC, error) {
	if sessionToken != "validToken" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "error while trying to authenticate session", nil, nil), nil
	}
	return common.GeneralError(http.StatusOK, response.Success, "", nil, nil), nil
}

func mockGetAllKeysFromTable(table string, dbtype persistencemgr.DbType) ([]string, error) {
	return []string{"/redfish/v1/LicenseService/Licenses/uuid.1.1", "/redfish/v1/LicenseService/Licenses/uuid.1.2"}, nil
}

func mockGetResource(table, key string, dbtype persistencemgr.DbType) (interface{}, *errors.Error) {
	if key == "/redfish/v1/LicenseService/Licenses" {
		return "", errors.PackError(errors.DBKeyNotFound, "not found")
	} else if key == "/redfish/v1/LicenseService/Licenses/uuid.1.1" {
		return string(`{"@odata.id":"/redfish/v1/LicenseService/Licenses/1.1","@odata.type":"#HpeiLOLicense.v2_3_0.HpeiLOLicense","Id":"1","Name":"iLO License","LicenseType":"Perpetual"}`), nil
	} else if key == "/redfish/v1/Systems/uuid.1" {
		return string(`{"Id": "uuid.1",
		"IndicatorLED": "Off",
		"Links":{
		"Chassis":[
		{
		"@odata.id": "/redfish/v1/Chassis/uuid.1"
		}
		],
		"ManagedBy":[
		{
		"@odata.id": "/redfish/v1/Managers/uuid.1"
		}
		]
		}}`), nil
	} else if key == "/redfish/v1/AggregationService/Aggregates/uuid" {
		resourceData := "{\"Elements\":[{\"@odata.id\":\"/redfish/v1/Systems/uuid.1\"}]}"
		var resource interface{}
		json.Unmarshal([]byte(resourceData), &resource)
		return resource, nil
	} else if key == "/redfish/v1/AggregationService/Aggregates/uuid2" {
		resourceData := "{\"Elements\":[{\"@odata.id\":\"/redfish/v1/Systems/uuid.2\"}]}"
		var resource interface{}
		json.Unmarshal([]byte(resourceData), &resource)
		return resource, nil
	} else if key == "/redfish/v1/Systems/uuid.2" {
		return string(`{"Id": "uuid.2",
		"IndicatorLED": "Off",
		"Links":{
		"Chassis":[
		{
		"@odata.id": "/redfish/v1/Chassis/uuid.2"
		}
		],
		"ManagedBy":[
		{
		"@odata.id": "/redfish/v1/Managers/uuid.2"
		}
		]
		}}`), nil
	}
	return "body", nil
}

func mockContactClient(ctx context.Context, url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {
	baseURI := "/redfish/v1"

	if url == "https://localhost:9091"+baseURI+"/LicenseService" {
		body := `{"data": "/ODIM/v1/Managers/1/EthernetInterfaces"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9093"+baseURI+"LicenseService/Licenses" {
		body := `{"data": "/redfish/v1/LicenseService/Licenses"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9092"+baseURI+"LicenseService/Licenses/uuid.1.1" {
		body := `{"data": "/ODIM/v1/Managers/uuid/EthernetInterfaces"}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}
