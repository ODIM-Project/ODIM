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

package rpc

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	licenseproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/licenses"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	licenseService "github.com/ODIM-Project/ODIM/svc-licenses/licenses"
	"github.com/ODIM-Project/ODIM/svc-licenses/model"
)

func mockGetExternalInterface() *licenseService.ExternalInterface {
	return &licenseService.ExternalInterface{
		External: licenseService.External{
			Auth:           mockIsAuthorized,
			ContactClient:  mockContactClient,
			GetTarget:      mockGetTarget,
			GetPluginData:  mockGetPluginData,
			ContactPlugin:  mockContactPlugin,
			DevicePassword: stubDevicePassword,
			GenericSave:    stubGenericSave,
		},
		DB: licenseService.DB{
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

func TestUpdate_GetLicenseService(t *testing.T) {
	license := new(Licenses)
	license.connector = mockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *licenseproto.GetLicenseServiceRequest
	}
	tests := []struct {
		name    string
		a       *Licenses
		args    args
		wantErr bool
	}{
		{
			name: "positive GetLicenseService",
			a:    license,
			args: args{
				req: &licenseproto.GetLicenseServiceRequest{SessionToken: "validToken"},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    license,
			args: args{
				req: &licenseproto.GetLicenseServiceRequest{SessionToken: "invalidToken"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.a.GetLicenseService(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("License.GetLicenseService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdate_GetLicenseCollection(t *testing.T) {
	license := new(Licenses)
	license.connector = mockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *licenseproto.GetLicenseRequest
	}
	tests := []struct {
		name    string
		a       *Licenses
		args    args
		wantErr bool
	}{
		{
			name: "positive GetLicenseCollection",
			a:    license,
			args: args{
				req: &licenseproto.GetLicenseRequest{SessionToken: "validToken"},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    license,
			args: args{
				req: &licenseproto.GetLicenseRequest{SessionToken: "invalidToken"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.a.GetLicenseCollection(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("License.GetLicenseCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdate_GetLicenseResource(t *testing.T) {
	license := new(Licenses)
	license.connector = mockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *licenseproto.GetLicenseResourceRequest
	}
	tests := []struct {
		name    string
		a       *Licenses
		args    args
		wantErr bool
	}{
		{
			name: "positive GetLicenseResource",
			a:    license,
			args: args{
				req: &licenseproto.GetLicenseResourceRequest{SessionToken: "validToken"},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    license,
			args: args{
				req: &licenseproto.GetLicenseResourceRequest{SessionToken: "invalidToken"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.a.GetLicenseResource(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("License.GetLicenseResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdate_InstallLicenseService(t *testing.T) {
	license := new(Licenses)
	license.connector = mockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *licenseproto.InstallLicenseRequest
	}
	tests := []struct {
		name    string
		a       *Licenses
		args    args
		wantErr bool
	}{
		{
			name: "positive InstallLicenseService",
			a:    license,
			args: args{
				req: &licenseproto.InstallLicenseRequest{SessionToken: "validToken"},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    license,
			args: args{
				req: &licenseproto.InstallLicenseRequest{SessionToken: "invalidToken"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.a.InstallLicenseService(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("License.InstallLicenseService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
