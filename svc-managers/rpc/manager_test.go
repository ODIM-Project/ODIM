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

	"github.com/stretchr/testify/assert"
	"github.com/bharath-b-hpe/odimra/lib-utilities/common"
	"github.com/bharath-b-hpe/odimra/lib-utilities/config"
	"github.com/bharath-b-hpe/odimra/lib-utilities/response"
	"github.com/bharath-b-hpe/odimra/svc-managers/mgrmodel"

	managersproto "github.com/bharath-b-hpe/odimra/lib-utilities/proto/managers"
)

func mockIsAuthorized(sessionToken string, privileges, oemPrivileges []string) (int32, string) {
	if sessionToken != "validToken" {
		return http.StatusUnauthorized, response.NoValidSession
	}
	return http.StatusOK, response.Success
}

func mockContactClient(url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {
	return nil, fmt.Errorf("InvalidRequest")
}

func TestGetManagerCollection(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
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
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	var ctx context.Context
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.ContactClientRPC = mockContactClient
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
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	mngr := mgrmodel.RAManager{
		Name:            "odimra",
		ManagerType:     "Service",
		FirmwareVersion: config.Data.FirmwareVersion,
		ID:              config.Data.RootServiceUUID,
		UUID:            config.Data.RootServiceUUID,
		State:           "Enabled",
	}
	mngr.AddManagertoDB()

	var ctx context.Context
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.ContactClientRPC = mockContactClient
	req := &managersproto.ManagerRequest{
		ManagerID:    "3bd1f589-117a-4cf9-89f2-da44ee8e012b",
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
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	var ctx context.Context
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.ContactClientRPC = mockContactClient
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
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	var ctx context.Context
	mgr := new(Managers)
	mgr.IsAuthorizedRPC = mockIsAuthorized
	mgr.ContactClientRPC = mockContactClient

	body := []byte(`body`)
	table := "EthernetInterfaces"
	key := "/redfish/v1/Managers/uuid:1/EthernetInterfaces/1"
	mgrmodel.GenericSave(body, table, key)

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
