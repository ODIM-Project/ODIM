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
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
)

func createMockRole(roleID string, privileges []string, oemPrivileges []string, predefined bool) error {
	role := asmodel.Role{
		ID:                 roleID,
		AssignedPrivileges: privileges,
		OEMPrivileges:      oemPrivileges,
		IsPredefined:       predefined,
	}

	if err := role.Create(); err != nil {
		return err
	}
	return nil
}

func TestGetRole(t *testing.T) {
	header := map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	commonResponse := response.Response{
		OdataType: common.RoleType,
		OdataID:   "/redfish/v1/AccountService/Roles/" + common.RoleAdmin,
		Name:      "User Role",
		ID:        common.RoleAdmin,
	}
	commonResponse.CreateGenericResponse(response.Success)
	commonResponse.MessageID = ""
	commonResponse.Message = ""
	commonResponse.MessageID = ""
	commonResponse.Severity = ""

	auth.Lock.Lock()
	common.SetUpMockConfig()
	auth.Lock.Unlock()
	defer truncateDB(t)
	err := createMockRole(common.RoleAdmin, []string{common.PrivilegeConfigureUsers}, []string{}, false)
	if err != nil {
		t.Fatalf("Error in creating mock admin user %v", err)
	}
	type args struct {
		req     *roleproto.GetRoleRequest
		session *asmodel.Session
	}
	var errArgs response.Args
	errArgs = response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InsufficientPrivilege,
				ErrorMessage:  "User does not have the privilege to get the role",
				MessageArgs:   []interface{}{},
			},
		},
	}
	errArg := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "Error while getting the role : error while trying to get role details: no data with the with key " + common.RoleClient + " found",
				MessageArgs:   []interface{}{"Role", common.RoleClient},
			},
		},
	}
	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "successful get role with name admin",
			args: args{
				req: &roleproto.GetRoleRequest{
					Id: common.RoleAdmin,
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
				Header:        header,
				Body: asresponse.UserRole{
					Response:           commonResponse,
					AssignedPrivileges: []string{common.PrivilegeConfigureUsers},
					OEMPrivileges:      []string{},
				},
			},
		}, {
			name: "request with insufficient privilege",
			args: args{
				req: &roleproto.GetRoleRequest{
					Id: common.RoleAdmin,
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
				Header:        header,
				Body:          errArgs.CreateGenericErrorResponse(),
			},
		}, {
			name: "get non-existing role",
			args: args{
				req: &roleproto.GetRoleRequest{
					Id: common.RoleClient,
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Header:        header,
				Body:          errArg.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetRole(tt.args.req, tt.args.session)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllRoles(t *testing.T) {
	commonResponse := response.Response{
		OdataType: "#RoleCollection.RoleCollection",
		OdataID:   "/redfish/v1/AccountService/Roles",
		Name:      "Roles Collection",
	}
	header := map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	commonResponse.CreateGenericResponse(response.Success)
	commonResponse.MessageID = ""
	commonResponse.Message = ""
	commonResponse.MessageID = ""
	commonResponse.Severity = ""
	auth.Lock.Lock()
	common.SetUpMockConfig()
	auth.Lock.Unlock()
	defer truncateDB(t)
	err := createMockRole(common.RoleAdmin, []string{common.PrivilegeConfigureUsers}, []string{}, false)
	if err != nil {
		t.Fatalf("Error in creating mock admin role %v", err)
	}
	errArgs := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InsufficientPrivilege,
				ErrorMessage:  "User does not have the privilege to get the roles",
				MessageArgs:   []interface{}{},
			},
		},
	}
	type args struct {
		session *asmodel.Session
	}
	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "successful get all roles",
			args: args{
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header:        header,
				Body: asresponse.List{
					Response:     commonResponse,
					MembersCount: 1,
					Members: []asresponse.ListMember{
						asresponse.ListMember{
							OdataID: "/redfish/v1/AccountService/Roles/" + common.RoleAdmin,
						},
					},
				},
			},
		}, {
			name: "request wth insufficient privilege",
			args: args{
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureManager: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.InsufficientPrivilege,
				Header:        header,
				Body:          errArgs.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAllRoles(tt.args.session)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}
