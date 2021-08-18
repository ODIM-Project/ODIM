// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package rpc

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-update/update"
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

func mockGetExternalInterface() *update.ExternalInterface {
	return &update.ExternalInterface{
		External: update.External{
			Auth:          mockIsAuthorized,
			ContactClient: mockContactClient,
		},
		DB: update.DB{
			GetAllKeysFromTable: mockGetAllKeysFromTable,
			GetResource:         mockGetResource,
		},
	}
}

func TestUpdate_GetUpdateService(t *testing.T) {
	update := new(Updater)
	update.connector = mockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *updateproto.UpdateRequest
	}
	tests := []struct {
		name    string
		a       *Updater
		args    args
		wantErr bool
	}{
		{
			name: "positive GetAggregationService",
			a:    update,
			args: args{
				req: &updateproto.UpdateRequest{SessionToken: "validToken"},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    update,
			args: args{
				req: &updateproto.UpdateRequest{SessionToken: "invalidToken"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.a.GetUpdateService(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Update.GetUpdateService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdate_GetFirmwareInventoryCollection(t *testing.T) {
	update := new(Updater)
	update.connector = mockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *updateproto.UpdateRequest
	}
	tests := []struct {
		name       string
		a          *Updater
		args       args
		StatusCode int
	}{
		{
			name: "Request with valid token",
			a:    update,
			args: args{
				req: &updateproto.UpdateRequest{
					SessionToken: "validToken",
				},
			}, StatusCode: 200,
		},
		{
			name: "Request with invalid token",
			a:    update,
			args: args{
				req: &updateproto.UpdateRequest{
					SessionToken: "invalidToken",
				},
			}, StatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, err := tt.a.GetFirmwareInventoryCollection(tt.args.ctx, tt.args.req); err != nil {
				t.Errorf("Update.GetFirmwareInventoryCollection() got = %v, want %v", resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

func TestUpdate_GetSoftwareInventoryCollection(t *testing.T) {
	update := new(Updater)
	update.connector = mockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *updateproto.UpdateRequest
	}
	tests := []struct {
		name       string
		a          *Updater
		args       args
		StatusCode int
	}{
		{
			name: "Request with valid token",
			a:    update,
			args: args{
				req: &updateproto.UpdateRequest{
					SessionToken: "validToken",
				},
			}, StatusCode: 200,
		},
		{
			name: "Request with invalid token",
			a:    update,
			args: args{
				req: &updateproto.UpdateRequest{
					SessionToken: "invalidToken",
				},
			}, StatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, err := tt.a.GetSoftwareInventoryCollection(tt.args.ctx, tt.args.req); err != nil {
				t.Errorf("Update.GetSoftwareInventoryCollection() got = %v, want %v", resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

func TestGetFirmwareInventorywithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	var ctx context.Context
	update := new(Updater)
	update.connector = mockGetExternalInterface()
	req := &updateproto.UpdateRequest{
		ResourceID:   "3bd1f589-117a-4cf9-89f2-da44ee8e012b",
		SessionToken: "InvalidToken",
	}
	resp, _ := update.GetFirmwareInventory(ctx, req)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetFirmwareInventorywithValidtoken(t *testing.T) {
	var ctx context.Context
	update := new(Updater)
	update.connector = mockGetExternalInterface()
	req := &updateproto.UpdateRequest{
		ResourceID:   "3bd1f589-117a-4cf9-89f2-da44ee8e012b:1",
		SessionToken: "validToken",
	}
	resp, err := update.GetFirmwareInventory(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetSoftwareInventorywithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	var ctx context.Context
	update := new(Updater)
	update.connector = mockGetExternalInterface()
	req := &updateproto.UpdateRequest{
		ResourceID:   "3bd1f589-117a-4cf9-89f2-da44ee8e012b",
		SessionToken: "InvalidToken",
	}
	resp, _ := update.GetSoftwareInventory(ctx, req)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetSoftwareInventorywithValidtoken(t *testing.T) {
	var ctx context.Context
	update := new(Updater)
	update.connector = mockGetExternalInterface()
	req := &updateproto.UpdateRequest{
		ResourceID:   "3bd1f589-117a-4cf9-89f2-da44ee8e012b:1",
		SessionToken: "validToken",
	}
	resp, err := update.GetSoftwareInventory(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status code should be StatusOK.")
}
