// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package role

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
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

func mockContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.TransactionID, "xyz")
	ctx = context.WithValue(ctx, common.ActionID, "001")
	ctx = context.WithValue(ctx, common.ActionName, "xyz")
	ctx = context.WithValue(ctx, common.ThreadID, "0")
	ctx = context.WithValue(ctx, common.ThreadName, "xyz")
	ctx = context.WithValue(ctx, common.ProcessName, "xyz")
	return ctx
}

func TestCreate(t *testing.T) {
	config.SetUpMockConfig(t)
	common.SetUpMockConfig()
	defer truncateDB(t)
	commonResponse := response.Response{
		OdataType: common.RoleType,
		OdataID:   "/redfish/v1/AccountService/Roles/testRole",
		Name:      "User Role",
		ID:        "testRole",
	}
	commonResponse.CreateGenericResponse(response.ResourceCreated)
	err := mockPrivilegeRegistry()
	if err != nil {
		t.Fatalf("Error in creating mock privilege registry %v", err)
	}
	err = mockRedfishRoles()
	if err != nil {
		t.Fatalf("Error in creating mock redfish predefined roles %v", err)
	}
	errArgs := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InsufficientPrivilege,
				ErrorMessage:  "failed to create role testRole: User does not have the privilege of creating a new role",
				MessageArgs:   []interface{}{},
			},
		},
	}
	errArgsMiss := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.PropertyMissing,
				ErrorMessage:  "failed to create role testRole: Both AssignedPrivileges and OemPrivileges cannot be empty.",
				MessageArgs:   []interface{}{"AssignedPrivileges/OemPrivileges"},
			},
		},
	}
	errArg := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.PropertyValueNotInList,
				ErrorMessage:  "failed to create role testRole: Requested Redfish predefined privilege is not correct",
				MessageArgs:   []interface{}{"Configure", "AssignedPrivileges"},
			},
		},
	}
	errArgsInvalid := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.PropertyValueNotInList,
				ErrorMessage:  "failed to create role @testRole: Invalid create role request",
				MessageArgs:   []interface{}{"@testRole", "RoleId"},
			},
		},
	}
	errArgsInvalidRole := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InsufficientPrivilege,
				ErrorMessage:  "failed to create role Administrator: Cannot create pre-defined roles",
				MessageArgs:   []interface{}{},
			},
		},
	}

	errArgu := response.Args{
		Code:    response.GeneralError,
		Message: "Role with name testRole already exists",
	}
	reqBodyCreateRole, _ := json.Marshal(asmodel.Role{
		ID:                 "testRole",
		AssignedPrivileges: []string{common.PrivilegeLogin},
		OEMPrivileges:      []string{},
	})
	reqBodyRoleConfigure, _ := json.Marshal(asmodel.Role{
		ID:                 "testRole",
		AssignedPrivileges: []string{"Configure"},
		OEMPrivileges:      []string{},
	})
	reqBodyInvalidRole, _ := json.Marshal(asmodel.Role{
		ID:                 "@testRole",
		AssignedPrivileges: []string{common.PrivilegeLogin},
		OEMPrivileges:      []string{},
	})
	reqBodyRoleEmpPrivilege, _ := json.Marshal(asmodel.Role{
		ID:                 "testRole",
		AssignedPrivileges: []string{},
		OEMPrivileges:      []string{},
	})
	reqBodyCreateAdminRole, _ := json.Marshal(asmodel.Role{
		ID:                 common.RoleAdmin,
		AssignedPrivileges: []string{common.PrivilegeLogin},
		OEMPrivileges:      []string{},
	})
	ctx := mockContext()
	type args struct {
		req     *roleproto.RoleRequest
		session *asmodel.Session
	}
	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "request for successful role creation",
			args: args{
				req: &roleproto.RoleRequest{
					RequestBody: reqBodyCreateRole,
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusCreated,
				StatusMessage: response.ResourceCreated,
				Body: asresponse.UserRole{
					IsPredefined:       false,
					AssignedPrivileges: []string{common.PrivilegeLogin},
					OEMPrivileges:      []string{},
					Response:           commonResponse,
				},
			},
		},
		{
			name: "request with insufficient privilege",
			args: args{
				req: &roleproto.RoleRequest{
					RequestBody: reqBodyCreateRole,
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureManager: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.InsufficientPrivilege,
				Body:          errArgs.CreateGenericErrorResponse(),
			},
		},
		{
			name: "request with invalid assigned privilege",
			args: args{
				req: &roleproto.RoleRequest{
					RequestBody: reqBodyRoleConfigure,
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusBadRequest,
				StatusMessage: response.PropertyValueNotInList,
				Body:          errArg.CreateGenericErrorResponse(),
			},
		},
		{
			name: "request with invalid character in role",
			args: args{
				req: &roleproto.RoleRequest{
					RequestBody: reqBodyInvalidRole,
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusBadRequest,
				StatusMessage: response.PropertyValueNotInList,
				Body:          errArgsInvalid.CreateGenericErrorResponse(),
			},
		},
		{
			name: "request with empty assigned privilege",
			args: args{
				req: &roleproto.RoleRequest{
					RequestBody: reqBodyRoleEmpPrivilege,
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusBadRequest,
				StatusMessage: response.PropertyMissing,
				Body:          errArgsMiss.CreateGenericErrorResponse(),
			},
		},
		{
			name: "request for creating an existing role",
			args: args{
				req: &roleproto.RoleRequest{
					RequestBody: reqBodyCreateRole,
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			}, want: response.RPC{
				StatusCode:    http.StatusConflict,
				StatusMessage: response.GeneralError,
				Body:          errArgu.CreateGenericErrorResponse(),
			},
		},
		{
			name: "request for creating an pre-existing role",
			args: args{
				req: &roleproto.RoleRequest{
					RequestBody: reqBodyCreateAdminRole,
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			}, want: response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.InsufficientPrivilege,
				Body:          errArgsInvalidRole.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Create(ctx, tt.args.req, tt.args.session)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}

}
