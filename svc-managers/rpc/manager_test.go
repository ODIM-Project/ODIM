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

package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-managers/managers"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrcommon"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrmodel"
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

func mockGetManagerData(id string) (mgrmodel.RAManager, error) {
	if id == "nonExistingUUID" {
		return mgrmodel.RAManager{}, fmt.Errorf("not found")
	} else if id == "noDevice" {
		return mgrmodel.RAManager{
			Name:            "odimra",
			ManagerType:     "Service",
			FirmwareVersion: "1.0",
			ID:              "noDevice",
			UUID:            "noDevice",
			State:           "Absent",
		}, nil
	}
	return mgrmodel.RAManager{
		Name:            "odimra",
		ManagerType:     "Service",
		FirmwareVersion: "1.0",
		ID:              config.Data.RootServiceUUID,
		UUID:            config.Data.RootServiceUUID,
		State:           "Enabled",
	}, nil
}

func mockGetManagerByURL(url string) (string, *errors.Error) {
	if url == "/redfish/v1/Managers/invalidURL:1" || url == "/redfish/v1/Managers/invalidURL" || url == "/redfish/v1/Managers/invalidID" {
		return "", errors.PackError(errors.DBKeyNotFound, "not found")
	}
	managerData := make(map[string]interface{})
	managerData["ManagerType"] = "BMC"
	managerData["Status"] = `{"State":"Enabled"}}`
	managerData["Name"] = "somePlugin"
	if url == "/redfish/v1/Managers/uuid" {
		managerData["Name"] = "someOtherID"
	} else if url == "/redfish/v1/Managers/noPlugin" {
		managerData["Name"] = "noPlugin"
	} else if url == "/redfish/v1/Managers/noToken" {
		managerData["Name"] = "noToken"
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

func mockUpdateManagersData(key string, managerData map[string]interface{}) error {
	return nil
}

func mockGetResource(table, key string) (string, *errors.Error) {
	if key == "/redfish/v1/Managers/uuid1:1/Ethernet" {
		return "", errors.PackError(errors.DBKeyNotFound, "not found")
	}
	return "body", nil
}

func mockGetDeviceInfo(req mgrcommon.ResourceInfoRequest) (string, error) {
	if req.URL == "/redfish/v1/Managers/deviceAbsent:1" || req.URL == "/redfish/v1/Managers/uuid1:1/Ethernet" {
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

func mockGetExternalInterface() *managers.ExternalInterface {
	return &managers.ExternalInterface{
		Device: managers.Device{
			GetDeviceInfo: mockGetDeviceInfo,
			ContactClient: mockContactClient,
		},
		DB: managers.DB{
			GetAllKeysFromTable: mockGetAllKeysFromTable,
			GetManagerData:      mockGetManagerData,
			GetManagerByURL:     mockGetManagerByURL,
			GetPluginData:       mockGetPluginData,
			UpdateManagersData:  mockUpdateManagersData,
			GetResource:         mockGetResource,
		},
	}
}

func mockGetAllKeysFromTable(table string) ([]string, error) {
	return []string{"/redfish/v1/Managers/uuid:1"}, nil
}

func TestGetManagerCollection(t *testing.T) {
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.EI = mockGetExternalInterface()
	type args struct {
		ctx  context.Context
		req  *managersproto.ManagerRequest
		resp *managersproto.ManagerResponse
	}
	tests := []struct {
		name       string
		mgr        *Managers
		args       args
		StatusCode int
	}{
		{
			name: "Request with valid token",
			mgr:  mgr,
			args: args{
				req: &managersproto.ManagerRequest{
					SessionToken: "validToken",
				},
				resp: &managersproto.ManagerResponse{},
			}, StatusCode: 200,
		},
		{
			name: "Request with invalid token",
			mgr:  mgr,
			args: args{
				req: &managersproto.ManagerRequest{
					SessionToken: "invalidToken",
				},
				resp: &managersproto.ManagerResponse{},
			}, StatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.mgr.GetManagersCollection(tt.args.ctx, tt.args.req, tt.args.resp); err != nil {
				t.Errorf("Manager.GetManagersCollection() got = %v, want %v", tt.args.resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

func TestGetManagerwithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	var ctx context.Context
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.EI = mockGetExternalInterface()
	req := &managersproto.ManagerRequest{
		ManagerID:    "3bd1f589-117a-4cf9-89f2-da44ee8e012b",
		SessionToken: "InvalidToken",
	}
	var resp = &managersproto.ManagerResponse{}
	mgr.GetManager(ctx, req, resp)
	assert.Equal(t, int(resp.StatusCode), http.StatusUnauthorized, "Status code should be StatusOK.")
}
func TestGetManagerwithValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	var ctx context.Context
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.EI = mockGetExternalInterface()
	req := &managersproto.ManagerRequest{
		ManagerID:    config.Data.RootServiceUUID,
		SessionToken: "validToken",
	}
	var resp = &managersproto.ManagerResponse{}
	err := mgr.GetManager(ctx, req, resp)
	assert.Nil(t, err, "There should be no error")

	var manager mgrmodel.Manager
	json.Unmarshal(resp.Body, &manager)

	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "Status code should be StatusOK.")
	assert.Equal(t, manager.Name, "odimra", "Status code should be StatusOK.")
	assert.Equal(t, manager.ManagerType, "Service", "Status code should be StatusOK.")
	assert.Equal(t, manager.ID, req.ManagerID, "Status code should be StatusOK.")
	assert.Equal(t, manager.FirmwareVersion, "1.0", "Status code should be StatusOK.")
}

func TestGetManagerResourcewithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	var ctx context.Context
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.EI = mockGetExternalInterface()
	req := &managersproto.ManagerRequest{
		ManagerID:    "uuid:1",
		SessionToken: "InvalidToken",
	}
	var resp = &managersproto.ManagerResponse{}
	mgr.GetManagersResource(ctx, req, resp)
	assert.Equal(t, int(resp.StatusCode), http.StatusUnauthorized, "Status code should be StatusUnauthorized.")
}
func TestGetManagerResourcewithValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	var ctx context.Context
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.EI = mockGetExternalInterface()

	req := &managersproto.ManagerRequest{
		ManagerID:    "uuid:1",
		SessionToken: "validToken",
		URL:          "/redfish/v1/Managers/uuid:1/EthernetInterfaces/1",
		ResourceID:   "1",
	}
	var resp = &managersproto.ManagerResponse{}
	err := mgr.GetManagersResource(ctx, req, resp)
	assert.Nil(t, err, "The two words should be the same.")
	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "Status code should be StatusOK.")
}
