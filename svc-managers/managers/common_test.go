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

package managers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrcommon"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrmodel"
	"io/ioutil"
	"net/http"
	"testing"
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
		Device: Device{
			GetDeviceInfo: mockGetDeviceInfo,
			ContactClient: mockContactClient,
			DeviceRequest: mockDeviceRequest,
		},
		DB: DB{
			GetAllKeysFromTable: mockGetAllKeysFromTable,
			GetManagerByURL:     mockGetManagerByURL,
			GetPluginData:       mockGetPluginData,
			UpdateData:          mockUpdateData,
			GetResource:         mockGetResource,
		},
	}
}

func mockGetAllKeysFromTable(table string) ([]string, error) {
	return []string{"/redfish/v1/Managers/uuid:1"}, nil
}

func mockGetManagerByURL(url string) (string, *errors.Error) {
	managerData := make(map[string]interface{})
	managerData["ManagerType"] = "BMC"
	managerData["Status"] = `{"State":"Enabled"}}`
	managerData["Name"] = "somePlugin"
	switch url {
	case "/redfish/v1/Managers/nonExistingUUID", "/redfish/v1/Managers/invalidURL:1", "/redfish/v1/Managers/invalidURL", "/redfish/v1/Managers/invalidID":
		return "", errors.PackError(errors.DBKeyNotFound, "not found")
	case "/redfish/v1/Managers/noDevice":
		managerData["ManagerType"] = "Service"
		managerData["Status"] = `{"State":"Absent"}}`
		managerData["Name"] = "odimra"
		managerData["ID"] = "noDevice"
		managerData["UUID"] = "noDevice"
		managerData["FirmwareVersion"] = "1.0"
	case "/redfish/v1/Managers/uuid":
		managerData["Name"] = "someOtherID"
	case "/redfish/v1/Managers/noPlugin":
		managerData["Name"] = "noPlugin"
	case "/redfish/v1/Managers/noToken":
		managerData["Name"] = "noToken"
	case "/redfish/v1/Managers/" + config.Data.RootServiceUUID:
		managerData["ManagerType"] = "Service"
		managerData["Status"] = `{"State":"Enabled"}}`
		managerData["Name"] = "odimra"
		managerData["ManagerID"] = config.Data.RootServiceUUID
		managerData["UUID"] = config.Data.RootServiceUUID
		managerData["FirmwareVersion"] = "1.0"
	}
	data, _ := json.Marshal(managerData)
	return string(data), nil
}

func mockGetPluginData(pluginID string) (mgrmodel.Plugin, *errors.Error) {
	if pluginID == "someOtherID" {
		return mgrmodel.Plugin{
			IP:                "localhost",
			Port:              "9091",
			Username:          "admin",
			Password:          []byte("password"),
			ID:                "CFM",
			PreferredAuthType: "XAuthToken",
		}, nil
	} else if pluginID == "noToken" {
		return mgrmodel.Plugin{
			IP:                "localhost",
			Port:              "9092",
			Username:          "admin",
			Password:          []byte("password"),
			ID:                "noToken",
			PreferredAuthType: "XAuthToken",
		}, nil
	} else if pluginID == "noPlugin" {
		return mgrmodel.Plugin{}, errors.PackError(errors.DBKeyNotFound, "not found")
	}
	return mgrmodel.Plugin{
		IP:                "localhost",
		Port:              "9093",
		Username:          "admin",
		Password:          []byte("password"),
		ID:                "somePlugin",
		PreferredAuthType: "BasicAuth",
	}, nil
}

func mockUpdateData(key string, updateData map[string]interface{}, table string) error {
	if key == "/redfish/v1/Managers/uuid:1/VirtualMedia/1" {
		return nil
	} else if key == "/redfish/v1/Managers/uuid1:1/VirtualMedia/4" {
		return errors.PackError(errors.DBKeyNotFound, "not found")
	}
	return nil
}

func mockGetResource(table, key string) (string, *errors.Error) {
	if key == "/redfish/v1/Managers/uuid1:1/Ethernet" {
		return "", errors.PackError(errors.DBKeyNotFound, "not found")
	} else if key == "/redfish/v1/Managers/uuid1:1/Virtual" {
		return "", errors.PackError(errors.DBKeyNotFound, "not found")
	} else if key == "/redfish/v1/Managers/uuid1:1/VirtualMedia/4" {
		return "", errors.PackError(errors.DBKeyNotFound, "not found")
	}
	return "body", nil
}

func mockGetDeviceInfo(req mgrcommon.ResourceInfoRequest) (string, error) {
	if req.URL == "/redfish/v1/Managers/deviceAbsent:1" || req.URL == "/redfish/v1/Managers/uuid1:1/Ethernet" {
		return "", fmt.Errorf("error")
	} else if req.URL == "/redfish/v1/Managers/uuid1:1/Virtual" {
		return "", fmt.Errorf("error")
	} else if req.URL == "/redfish/v1/Managers/uuid1:1/VirtualMedia/4" {
		return "", fmt.Errorf("error")
	}
	manager := mgrmodel.Manager{
		Status: &mgrmodel.Status{
			State: "Enabled",
		},
	}
	dataByte, err := json.Marshal(manager)
	return string(dataByte), err
}

func mockDeviceRequest(req mgrcommon.ResourceInfoRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{"Content-type": "application/json; charset=utf-8"}
	if req.URL == "/redfish/v1/Managers/deviceAbsent:1" || req.URL == "/redfish/v1/Managers/uuid1:1/Virtual" {
		resp.StatusCode = http.StatusNotFound
		resp.StatusMessage = response.ResourceNotFound
		return resp
	}
	manager := mgrmodel.Manager{
		Status: &mgrmodel.Status{
			State: "Enabled",
		},
	}
	dataByte, err := json.Marshal(manager)
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	err = json.Unmarshal(dataByte, &resp.Body)
	if err != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
	}
	return resp
}

func mockContactClient(url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {
	baseURI := "/redfish/v1"
	baseURI = mgrcommon.TranslateToSouthBoundURL(baseURI)

	if url == "https://localhost:9091"+baseURI+"/Sessions" {
		body := `{"Token": "12345"}`
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			Header: http.Header{
				"X-Auth-Token": []string{"12345"},
			},
		}, nil
	} else if url == "https://localhost:9092"+baseURI+"/Sessions" {
		body := `{"Token": ""}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "https://localhost:9091"+baseURI+"/Managers/uuid/EthernetInterfaces" && token == "12345" {
		body := `{"data": "/ODIM/v1/Managers/uuid/EthernetInterfaces"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9093"+baseURI+"/Managers/uuid1/EthernetInterfaces" {
		body := `{"data": "/ODIM/v1/Managers/uuid/EthernetInterfaces"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9092"+baseURI+"/Managers/uuid/EthernetInterfaces" && token == "23456" {
		body := `{"data": "/ODIM/v1/Managers/uuid/EthernetInterfaces"}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9091"+baseURI+"/Managers/uuid/VirtualMedia/1/Actions/VirtualMedia.InsertMedia" {
		body := `{"data": "Success"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}
