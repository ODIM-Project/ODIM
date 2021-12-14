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
package role

import (
	"encoding/json"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"

	"net/http"
	"reflect"
	"testing"
)

func TestUpdate(t *testing.T) {
	common.SetUpMockConfig()
	defer truncateDB(t)
	err := createMockRole(common.RoleAdmin, []string{common.PrivilegeConfigureUsers}, []string{}, false)
	if err != nil {
		t.Fatalf("Error in creating mock admin role %v", err)
	}
	err = createMockRole("MockRole", []string{common.PrivilegeConfigureUsers}, []string{}, false)
	if err != nil {
		t.Fatalf("Error in creating mock role MockRole %v", err)
	}
	err = mockPrivilegeRegistry()
	if err != nil {
		t.Fatalf("Error in creating mock privilege registry %v", err)
	}
	err = mockRedfishRoles()
	if err != nil {
		t.Fatalf("Error in creating mock redfish predefined roles %v", err)
	}
	validRoleReq, _ := json.Marshal(asmodel.Role{
		AssignedPrivileges: []string{common.PrivilegeLogin, common.PrivilegeConfigureUsers},
		OEMPrivileges:      []string{},
	})
	invalidRoleReq, _ := json.Marshal(asmodel.Role{
		AssignedPrivileges: []string{"Configue"},
		OEMPrivileges:      []string{},
	})
	duplicateRoleReq, _ := json.Marshal(asmodel.Role{
		AssignedPrivileges: []string{common.PrivilegeLogin, common.PrivilegeLogin},
		OEMPrivileges:      []string{},
	})
	errArg := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InsufficientPrivilege,
				ErrorMessage:  "User does not have the privilege to update the role",
				MessageArgs:   []interface{}{},
			},
		},
	}
	errArgs := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.PropertyValueNotInList,
				ErrorMessage:  "Requested Redfish predefined privilege is not correct",
				MessageArgs:   []interface{}{"Configue", "AssignedPrivileges"},
			},
		},
	}
	errArgu := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "error while trying to get role details: no data with the with key NonExistentRole found",
				MessageArgs:   []interface{}{"Role", "NonExistentRole"},
			},
		},
	}
	errArgGen := response.Args{
		Code:    response.GeneralError,
		Message: "Updating predefined role is restricted",
	}
	errArgGen1 := response.Args{
		Code:    response.GeneralError,
		Message: "Duplicate privileges can not be updated",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.PropertyValueConflict,
				ErrorMessage:  "Duplicate privileges can not be updated",
				MessageArgs:   []interface{}{common.PrivilegeLogin, common.PrivilegeLogin},
			},
		},
	}

	type args struct {
		req     *roleproto.UpdateRoleRequest
		session *asmodel.Session
	}
	tests := []struct {
		name    string
		args    args
		want    response.RPC
		wantErr bool
	}{
		{
			name: "successful updation of role",
			args: args{
				req: &roleproto.UpdateRoleRequest{
					Id:            "MockRole",
					UpdateRequest: validRoleReq,
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body: asmodel.Role{
					ID:                 "MockRole",
					AssignedPrivileges: []string{common.PrivilegeLogin, common.PrivilegeConfigureUsers},
					OEMPrivileges:      []string{},
				},
			},
		},
		{
			name: "request with insufficient privilege",
			args: args{
				req: &roleproto.UpdateRoleRequest{
					Id:            "MockRole",
					UpdateRequest: validRoleReq,
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
				Body:          errArg.CreateGenericErrorResponse(),
			},
		},
		{
			name: " request invalid assigned privileges",
			args: args{
				req: &roleproto.UpdateRoleRequest{
					Id:            "MockRole",
					UpdateRequest: invalidRoleReq,
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
				Body:          errArgs.CreateGenericErrorResponse(),
			},
		},
		{
			name: "update non-existing role",
			args: args{
				req: &roleproto.UpdateRoleRequest{
					Id:            "NonExistentRole",
					UpdateRequest: validRoleReq,
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusBadRequest,
				StatusMessage: response.ResourceNotFound,
				Body:          errArgu.CreateGenericErrorResponse(),
			},
		},
		{
			name: "update predefined role role",
			args: args{
				req: &roleproto.UpdateRoleRequest{
					Id:            common.RoleAdmin,
					UpdateRequest: validRoleReq,
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusMethodNotAllowed,
				StatusMessage: response.GeneralError,
				Body:          errArgGen.CreateGenericErrorResponse(),
			},
		},
		{
			name: "update duplicate privileges",
			args: args{
				req: &roleproto.UpdateRoleRequest{
					Id:            "MockRole",
					UpdateRequest: duplicateRoleReq,
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusBadRequest,
				StatusMessage: response.PropertyValueConflict,
				Body:          errArgGen1.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Update(tt.args.req, tt.args.session)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
