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

//Package rpc defines the handler for micro services
package rpc

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
)

func mockRedfishRoles() error {
	list := asmodel.RedfishRoles{
		List: []string{
			"Administrator",
			"Operator",
			"ReadOnly",
		},
	}
	if err := list.Create(); err != nil {
		return err
	}
	return nil
}

func mockPrivilegeRegistry() error {
	list := asmodel.Privileges{
		List: []string{
			"Login",
			"ConfigureManager",
			"ConfigureUsers",
			"ConfigureSelf",
			"ConfigureComponents",
		},
	}
	if err := list.Create(); err != nil {
		return err
	}
	return nil
}

func createMockRole(roleID string, privileges []string, oemPrivileges []string) error {
	role := asmodel.Role{
		ID:                 roleID,
		AssignedPrivileges: privileges,
		OEMPrivileges:      oemPrivileges,
	}
	if err := role.Create(); err != nil {
		return err
	}
	return nil
}

func TestRole_CreateRole(t *testing.T) {
	defer truncateDB(t)
	auth.Lock.Lock()
	common.SetUpMockConfig()
	auth.Lock.Unlock()
	token := "token"
	err := mockSession(token, common.RoleAdmin)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	mockRedfishRoles()

	reqBodyCreateRole, _ := json.Marshal(asmodel.Role{
		ID:                 "testRole",
		AssignedPrivileges: []string{common.PrivilegeLogin},
		OEMPrivileges:      []string{},
	})

	type args struct {
		ctx  context.Context
		req  *roleproto.RoleRequest
		resp *roleproto.RoleResponse
	}
	tests := []struct {
		name    string
		r       *Role
		args    args
		wantErr bool
	}{
		{
			name: "CreateRole with valid session",
			args: args{
				req: &roleproto.RoleRequest{
					RequestBody:  reqBodyCreateRole,
					SessionToken: token,
				},
				resp: &roleproto.RoleResponse{},
			},
			wantErr: false,
		},
		{
			name: "CreateRole with invalid session",
			args: args{
				req: &roleproto.RoleRequest{
					RequestBody:  reqBodyCreateRole,
					SessionToken: "testToken",
				},
				resp: &roleproto.RoleResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.r.CreateRole(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Role.CreateRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRole_GetRole(t *testing.T) {
	defer truncateDB(t)
	auth.Lock.Lock()
	common.SetUpMockConfig()
	auth.Lock.Unlock()
	token := "token"
	err := mockSession(token, common.RoleAdmin)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	createMockRole("testRole", []string{common.PrivilegeConfigureManager}, []string{})
	type args struct {
		ctx  context.Context
		req  *roleproto.GetRoleRequest
		resp *roleproto.RoleResponse
	}
	tests := []struct {
		name    string
		r       *Role
		args    args
		wantErr bool
	}{
		{
			name: "GetRole with valid session",
			args: args{
				req: &roleproto.GetRoleRequest{
					Id:           "testRole",
					SessionToken: token,
				},
				resp: &roleproto.RoleResponse{},
			},
			wantErr: false,
		},
		{
			name: "GetRole with invalid session",
			args: args{
				req: &roleproto.GetRoleRequest{
					Id:           "testRole",
					SessionToken: "testToken",
				},
				resp: &roleproto.RoleResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.r.GetRole(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Role.GetRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRole_GetAllRoles(t *testing.T) {
	defer truncateDB(t)
	auth.Lock.Lock()
	common.SetUpMockConfig()
	auth.Lock.Unlock()
	token := "token"
	err := mockSession(token, common.RoleAdmin)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	createMockRole("testRole", []string{common.PrivilegeConfigureManager}, []string{})

	type args struct {
		ctx  context.Context
		req  *roleproto.GetRoleRequest
		resp *roleproto.RoleResponse
	}
	tests := []struct {
		name    string
		r       *Role
		args    args
		wantErr bool
	}{
		{
			name: "GetAllRoles with valid session",
			args: args{
				req: &roleproto.GetRoleRequest{
					SessionToken: token,
				},
				resp: &roleproto.RoleResponse{},
			},
			wantErr: false,
		},
		{
			name: "GetAllRoles with invalid session",
			args: args{
				req: &roleproto.GetRoleRequest{
					SessionToken: "testToken",
				},
				resp: &roleproto.RoleResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.r.GetAllRoles(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Role.GetAllRoles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRole_UpdateRole(t *testing.T) {
	defer truncateDB(t)
	auth.Lock.Lock()
	common.SetUpMockConfig()
	auth.Lock.Unlock()
	token := "token"
	err := mockSession(token, common.RoleAdmin)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	createMockRole("testRole", []string{common.PrivilegeConfigureManager}, []string{})
	mockRedfishRoles()
	validRoleReq, _ := json.Marshal(asmodel.Role{
		AssignedPrivileges: []string{common.PrivilegeLogin, common.PrivilegeConfigureUsers},
		OEMPrivileges:      []string{},
	})
	type args struct {
		ctx  context.Context
		req  *roleproto.UpdateRoleRequest
		resp *roleproto.RoleResponse
	}
	tests := []struct {
		name    string
		r       *Role
		args    args
		wantErr bool
	}{
		{
			name: "UpdateRole with valid session",
			args: args{
				req: &roleproto.UpdateRoleRequest{
					Id:            "testRole",
					UpdateRequest: validRoleReq,
					SessionToken:  token,
				},
				resp: &roleproto.RoleResponse{},
			},
			wantErr: false,
		},
		{
			name: "UpdateRole with invalid session",
			args: args{
				req: &roleproto.UpdateRoleRequest{
					Id:            "testRole",
					UpdateRequest: validRoleReq,
					SessionToken:  "testToken",
				},
				resp: &roleproto.RoleResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.r.UpdateRole(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Role.UpdateRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRole_DeleteRole(t *testing.T) {
	defer truncateDB(t)
	auth.Lock.Lock()
	common.SetUpMockConfig()
	auth.Lock.Unlock()
	token, roleID := "token", "testRole"
	err := mockSession(token, common.RoleClient)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	createMockRole(roleID, []string{common.PrivilegeConfigureComponents}, []string{})
	type args struct {
		ctx  context.Context
		req  *roleproto.DeleteRoleRequest
		resp *roleproto.RoleResponse
	}
	tests := []struct {
		name    string
		r       *Role
		args    args
		want    *roleproto.RoleResponse
		wantErr bool
	}{
		{
			name: "successful role deletion",
			args: args{
				req: &roleproto.DeleteRoleRequest{
					SessionToken: token,
					ID:           roleID,
				},
				resp: &roleproto.RoleResponse{},
			},
			want: &roleproto.RoleResponse{
				StatusCode:    http.StatusNoContent,
				StatusMessage: response.ResourceRemoved,
			},
			wantErr: false,
		},
		{
			name: "delete role with invalid session",
			args: args{
				req: &roleproto.DeleteRoleRequest{
					SessionToken: "invalid-token",
					ID:           roleID,
				},
				resp: &roleproto.RoleResponse{},
			},
			want: &roleproto.RoleResponse{
				StatusCode:    http.StatusUnauthorized,
				StatusMessage: response.NoValidSession,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.r.DeleteRole(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Role.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if resp.StatusCode != tt.want.StatusCode {
				t.Errorf("Role.Delete() StatusCode = %v, want = %v", resp.StatusCode, tt.want.StatusCode)
			}
			if resp.StatusMessage != tt.want.StatusMessage {
				t.Errorf("Role.Delete() StatusMessage = %v, want = %v", resp.StatusMessage, tt.want.StatusMessage)
			}
		})
	}
}
