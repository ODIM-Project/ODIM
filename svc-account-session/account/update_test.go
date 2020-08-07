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
package account

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	accountproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/account"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
)

func TestUpdate(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	successResponse := response.Response{
		OdataType:    "#ManagerAccount.v1_4_0.ManagerAccount",
		OdataID:      "/redfish/v1/AccountService/Accounts/testUser1",
		OdataContext: "/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
		ID:           "testUser1",
		Name:         "Account Service",
	}
	successResponse2 := response.Response{
		OdataType:    "#ManagerAccount.v1_4_0.ManagerAccount",
		OdataID:      "/redfish/v1/AccountService/Accounts/testUser2",
		OdataContext: "/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
		ID:           "testUser2",
		Name:         "Account Service",
	}
	successResponse.CreateGenericResponse(response.AccountModified)
	successResponse2.CreateGenericResponse(response.AccountModified)
	err := createMockRole("PrivilegeLogin", []string{common.PrivilegeLogin}, []string{}, false)
	if err != nil {
		t.Fatalf("Error in creating mock admin user %v", err)
	}
	err = createMockUser("testUser1", common.RoleAdmin)
	if err != nil {
		t.Fatalf("Error in creating mock admin user1 %v", err)
	}
	err = createMockUser("testUser2", common.RoleAdmin)
	if err != nil {
		t.Fatalf("Error in creating mock admin user2 %v", err)
	}
	err = createMockUser("testUser3", "PrivilegeLogin")
	if err != nil {
		t.Fatalf("Error in creating mock user3 %v", err)
	}

	errArg := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "error while trying to get  account: error while trying to get user: no data with the with key xyz found",
				MessageArgs:   []interface{}{"Account", "xyz"},
			},
		},
	}
	errArgs := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InsufficientPrivilege,
				ErrorMessage:  "error: user does not have the privilege to update other accounts",
				MessageArgs:   []interface{}{},
			},
		},
	}
	errArg4 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InsufficientPrivilege,
				ErrorMessage:  "error: roles, user is associated with, doesn't allow changing own or other users password",
				MessageArgs:   []interface{}{},
			},
		},
	}

	errArgs1 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InsufficientPrivilege,
				ErrorMessage:  "error: user does not have the privilege to update any account role, including his own account",
				MessageArgs:   []interface{}{},
			},
		},
	}

	errArg2 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.PropertyValueNotInList,
				ErrorMessage:  "error: Invalid RoleID xyz present",
				MessageArgs:   []interface{}{"xyz", "RoleID"},
			},
		},
	}

	errArg3 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.PropertyValueFormatError,
				ErrorMessage:  "error: invalid password, password length is less than the minimum length",
				MessageArgs:   []interface{}{"xyz", "Password"},
			},
		},
	}
	genArgs := response.Args{
		Code:    response.GeneralError,
		Message: "error: username cannot be modified",
	}
	type args struct {
		req     *accountproto.UpdateAccountRequest
		session *asmodel.Session
	}

	reqBodyRoleIDOperator, _ := json.Marshal(asmodel.Account{
		RoleID: "Operator",
	})
	reqBodyUpdateUsername, _ := json.Marshal(asmodel.Account{
		UserName: "xyz",
	})
	reqBodyInvalidRole, _ := json.Marshal(asmodel.Account{
		RoleID: "xyz",
	})
	reqBodyInvalidPwd, _ := json.Marshal(asmodel.Account{
		Password: "xyz",
	})
	reqBodyUpdatePwd, _ := json.Marshal(asmodel.Account{
		Password: "P@$$w0rd@123",
	})
	reqBodyRoleIDAdmin, _ := json.Marshal(asmodel.Account{
		RoleID: "Administrator",
	})

	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "successful updation of account as admin",
			args: args{
				req: &accountproto.UpdateAccountRequest{
					RequestBody: reqBodyRoleIDOperator,
					AccountID:   "testUser1",
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.AccountModified,
				Header: map[string]string{
					"Cache-Control":     "no-cache",
					"Connection":        "keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Link":              "</redfish/v1/AccountService/Accounts/testUser1/>; rel=describedby",
					"Location":          "/redfish/v1/AccountService/Accounts/testUser1/",
					"Transfer-Encoding": "chunked",
					"OData-Version":     "4.0",
				},
				Body: asresponse.Account{
					Response: successResponse,
					UserName: "testUser1",
					RoleID:   "Operator",
					Links: asresponse.Links{
						Role: asresponse.Role{
							OdataID: "/redfish/v1/AccountService/Roles/Operator/",
						},
					},
				},
			},
		},
		{
			name: "update role to admin without privilege",
			args: args{
				req: &accountproto.UpdateAccountRequest{
					RequestBody: reqBodyRoleIDAdmin,
					AccountID:   "testUser1",
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureSelf: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.InsufficientPrivilege,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArgs.CreateGenericErrorResponse(),
			},
		},
		{
			name: "update non-existing account",
			args: args{
				req: &accountproto.UpdateAccountRequest{
					RequestBody: reqBodyRoleIDOperator,
					AccountID:   "xyz",
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
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArg.CreateGenericErrorResponse(),
			},
		},
		{
			name: "update account name",
			args: args{
				req: &accountproto.UpdateAccountRequest{
					RequestBody: reqBodyUpdateUsername,
					AccountID:   "testUser1",
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusBadRequest,
				StatusMessage: response.GeneralError,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: genArgs.CreateGenericErrorResponse(),
			},
		},
		{
			name: "update account with invalid role",
			args: args{
				req: &accountproto.UpdateAccountRequest{
					RequestBody: reqBodyInvalidRole,
					AccountID:   "testUser1",
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
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArg2.CreateGenericErrorResponse(),
			},
		},
		{
			name: "update account with invalid password",
			args: args{
				req: &accountproto.UpdateAccountRequest{
					RequestBody: reqBodyInvalidPwd,
					AccountID:   "testUser1",
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusBadRequest,
				StatusMessage: response.PropertyValueFormatError,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArg3.CreateGenericErrorResponse(),
			},
		},
		{
			name: "update own password with ConfigureSelf privilege",
			args: args{
				req: &accountproto.UpdateAccountRequest{
					RequestBody: reqBodyUpdatePwd,
					AccountID:   "testUser1",
				},
				session: &asmodel.Session{
					ID:       "testUser1",
					UserName: "testUser1",
					RoleID:   "Operator",
					Privileges: map[string]bool{
						common.PrivilegeConfigureSelf: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.AccountModified,
				Header: map[string]string{
					"Cache-Control":     "no-cache",
					"Connection":        "keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Link":              "</redfish/v1/AccountService/Accounts/testUser1/>; rel=describedby",
					"Location":          "/redfish/v1/AccountService/Accounts/testUser1/",
					"Transfer-Encoding": "chunked",
					"OData-Version":     "4.0",
				},
				Body: asresponse.Account{
					Response: successResponse,
					UserName: "testUser1",
					RoleID:   "Operator",
					Links: asresponse.Links{
						Role: asresponse.Role{
							OdataID: "/redfish/v1/AccountService/Roles/Operator/",
						},
					},
				},
			},
		},
		{
			name: "update own password with ConfigureUsers privilege",
			args: args{
				req: &accountproto.UpdateAccountRequest{
					RequestBody: reqBodyUpdatePwd,
					AccountID:   "testUser1",
				},
				session: &asmodel.Session{
					ID:       "testUser1",
					UserName: "testUser1",
					RoleID:   "Operator",
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.AccountModified,
				Header: map[string]string{
					"Cache-Control":     "no-cache",
					"Connection":        "keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Link":              "</redfish/v1/AccountService/Accounts/testUser1/>; rel=describedby",
					"Location":          "/redfish/v1/AccountService/Accounts/testUser1/",
					"Transfer-Encoding": "chunked",
					"OData-Version":     "4.0",
				},
				Body: asresponse.Account{
					Response: successResponse,
					UserName: "testUser1",
					RoleID:   "Operator",
					Links: asresponse.Links{
						Role: asresponse.Role{
							OdataID: "/redfish/v1/AccountService/Roles/Operator/",
						},
					},
				},
			},
		},
		{
			name: "update own password with both ConfigureSelf and ConfigureUsers privileges",
			args: args{
				req: &accountproto.UpdateAccountRequest{
					RequestBody: reqBodyUpdatePwd,
					AccountID:   "testUser1",
				},
				session: &asmodel.Session{
					ID:       "testUser1",
					UserName: "testUser1",
					RoleID:   "Operator",
					Privileges: map[string]bool{
						common.PrivilegeConfigureSelf:  true,
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.AccountModified,
				Header: map[string]string{
					"Cache-Control":     "no-cache",
					"Connection":        "keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Link":              "</redfish/v1/AccountService/Accounts/testUser1/>; rel=describedby",
					"Location":          "/redfish/v1/AccountService/Accounts/testUser1/",
					"Transfer-Encoding": "chunked",
					"OData-Version":     "4.0",
				},
				Body: asresponse.Account{
					Response: successResponse,
					UserName: "testUser1",
					RoleID:   "Operator",
					Links: asresponse.Links{
						Role: asresponse.Role{
							OdataID: "/redfish/v1/AccountService/Roles/Operator/",
						},
					},
				},
			},
		},
		{
			name: "update other account password with both ConfigureUsers and ConfigureSelf priveleges",
			args: args{
				req: &accountproto.UpdateAccountRequest{
					RequestBody: reqBodyUpdatePwd,
					AccountID:   "testUser2",
				},
				session: &asmodel.Session{
					ID:       "testUser1",
					UserName: "testUser1",
					RoleID:   "Operator",
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
						common.PrivilegeConfigureSelf:  true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.AccountModified,
				Header: map[string]string{
					"Cache-Control":     "no-cache",
					"Connection":        "keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Link":              "</redfish/v1/AccountService/Accounts/testUser2/>; rel=describedby",
					"Location":          "/redfish/v1/AccountService/Accounts/testUser2/",
					"Transfer-Encoding": "chunked",
					"OData-Version":     "4.0",
				},
				Body: asresponse.Account{
					Response: successResponse2,
					UserName: "testUser2",
					RoleID:   "Administrator",
					Links: asresponse.Links{
						Role: asresponse.Role{
							OdataID: "/redfish/v1/AccountService/Roles/Administrator/",
						},
					},
				},
			},
		},
		{
			name: "update other account password with only ConfigureUsers privelege",
			args: args{
				req: &accountproto.UpdateAccountRequest{
					RequestBody: reqBodyUpdatePwd,
					AccountID:   "testUser2",
				},
				session: &asmodel.Session{
					ID:       "testUser1",
					UserName: "testUser1",
					RoleID:   "Operator",
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.AccountModified,
				Header: map[string]string{
					"Cache-Control":     "no-cache",
					"Connection":        "keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Link":              "</redfish/v1/AccountService/Accounts/testUser2/>; rel=describedby",
					"Location":          "/redfish/v1/AccountService/Accounts/testUser2/",
					"Transfer-Encoding": "chunked",
					"OData-Version":     "4.0",
				},
				Body: asresponse.Account{
					Response: successResponse2,
					UserName: "testUser2",
					RoleID:   "Administrator",
					Links: asresponse.Links{
						Role: asresponse.Role{
							OdataID: "/redfish/v1/AccountService/Roles/Administrator/",
						},
					},
				},
			},
		},
		{
			name: "update other account password with only ConfigureSelf privilege",
			args: args{
				req: &accountproto.UpdateAccountRequest{
					RequestBody: reqBodyUpdatePwd,
					AccountID:   "testUser2",
				},
				session: &asmodel.Session{
					ID:       "testUser1",
					UserName: "testUser1",
					RoleID:   "Operator",
					Privileges: map[string]bool{
						common.PrivilegeConfigureSelf: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.InsufficientPrivilege,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArgs.CreateGenericErrorResponse(),
			},
		},
		{
			name: "update account password with only Login privilege",
			args: args{
				req: &accountproto.UpdateAccountRequest{
					RequestBody: reqBodyUpdatePwd,
					AccountID:   "testUser3",
				},
				session: &asmodel.Session{
					ID:       "testUser3",
					UserName: "testUser3",
					RoleID:   "PrivilegeLogin",
					Privileges: map[string]bool{
						common.PrivilegeLogin: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.InsufficientPrivilege,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArg4.CreateGenericErrorResponse(),
			},
		},
		{
			name: "update account roleid with only Login privilege",
			args: args{
				req: &accountproto.UpdateAccountRequest{
					RequestBody: reqBodyRoleIDAdmin,
					AccountID:   "testUser3",
				},
				session: &asmodel.Session{
					ID:       "testUser3",
					UserName: "testUser3",
					RoleID:   "PrivilegeLogin",
					Privileges: map[string]bool{
						common.PrivilegeLogin: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.InsufficientPrivilege,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArgs1.CreateGenericErrorResponse(),
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
