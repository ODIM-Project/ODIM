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
package account

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	accountproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/account"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
)

func TestUpdate(t *testing.T) {
	config.SetUpMockConfig(t)
	acc := getMockExternalInterface()

	successResponse := createMockUpdateResponseObject(common.ManagerAccountType, "/redfish/v1/AccountService/Accounts/testUser1", "/redfish/v1/$metadata#ManagerAccount.ManagerAccount", "testUser1")

	operatorSuccessResponse := createMockUpdateResponseObject(common.ManagerAccountType, "/redfish/v1/AccountService/Accounts/operatorUser", "/redfish/v1/$metadata#ManagerAccount.ManagerAccount", "operatorUser")

	successResponse2 := createMockUpdateResponseObject(common.ManagerAccountType, "/redfish/v1/AccountService/Accounts/testUser2", "/redfish/v1/$metadata#ManagerAccount.ManagerAccount", "testUser2")

	errArg := GetResponseArgs(response.ResourceNotFound, "failed to update the account xyz: Unable to get account: error while trying to get user: no data with the with key xyz found", []interface{}{"Account", "xyz"})

	errArgs := GetResponseArgs(response.InsufficientPrivilege, "failed to update the account testUser2: User does not have the privilege of updating other accounts", []interface{}{})

	errArgs5 := GetResponseArgs(response.InsufficientPrivilege, "failed to update the account testUser1: User does not have the privilege of updating other accounts", []interface{}{})

	errArg4 := GetResponseArgs(response.InsufficientPrivilege, "failed to update the account testUser3: Roles, user is associated with, doesn't allow changing own or other users password", []interface{}{})

	errArgs1 := GetResponseArgs(response.InsufficientPrivilege, "failed to update the account testUser3: User does not have the privilege of updating any account role, including his own account", []interface{}{})

	errArg2 := GetResponseArgs(response.PropertyValueNotInList, "failed to update the account testUser1: Invalid RoleID xyz present", []interface{}{"xyz", "RoleID"})

	errArg3 := GetResponseArgs(response.PropertyValueFormatError, "error: invalid password, password length is less than the minimum length", []interface{}{"xyz", "Password"})

	errArg5 := GetResponseArgs(response.PropertyMissing, "failed to update the account testUser1: empty request can not be processed", []interface{}{"request body"})

	genArgs := response.Args{
		Code:    response.GeneralError,
		Message: "failed to update the account testUser1: Username cannot be modified",
	}
	ctx := mockContext()
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

	emptyPayload, _ := json.Marshal(map[string]interface{}{})

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
					"Link":     "</redfish/v1/AccountService/Accounts/testUser1/>; rel=describedby",
					"Location": "/redfish/v1/AccountService/Accounts/testUser1",
				},
				Body: asresponse.Account{
					Response: successResponse,
					UserName: "testUser1",
					RoleID:   "Operator",
					Links: asresponse.Links{
						Role: asresponse.Role{
							OdataID: "/redfish/v1/AccountService/Roles/Operator",
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
				Body:          errArgs5.CreateGenericErrorResponse(),
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
				Body:          errArg.CreateGenericErrorResponse(),
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
				Body:          genArgs.CreateGenericErrorResponse(),
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
				Body:          errArg2.CreateGenericErrorResponse(),
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
				Body:          errArg3.CreateGenericErrorResponse(),
			},
		},
		{
			name: "update own password with ConfigureSelf privilege",
			args: args{
				req: &accountproto.UpdateAccountRequest{
					RequestBody: reqBodyUpdatePwd,
					AccountID:   "operatorUser",
				},
				session: &asmodel.Session{
					ID:       "operatorUser",
					UserName: "operatorUser",
					RoleID:   common.RoleMonitor,
					Privileges: map[string]bool{
						common.PrivilegeConfigureSelf: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.AccountModified,
				Header: map[string]string{
					"Link":     "</redfish/v1/AccountService/Accounts/operatorUser/>; rel=describedby",
					"Location": "/redfish/v1/AccountService/Accounts/operatorUser",
				},
				Body: asresponse.Account{
					Response: operatorSuccessResponse,
					UserName: "operatorUser",
					RoleID:   "Operator",
					Links: asresponse.Links{
						Role: asresponse.Role{
							OdataID: "/redfish/v1/AccountService/Roles/Operator",
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
					AccountID:   "operatorUser",
				},
				session: &asmodel.Session{
					ID:       "operatorUser",
					UserName: "operatorUser",
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
					"Link":     "</redfish/v1/AccountService/Accounts/operatorUser/>; rel=describedby",
					"Location": "/redfish/v1/AccountService/Accounts/operatorUser",
				},
				Body: asresponse.Account{
					Response: operatorSuccessResponse,
					UserName: "operatorUser",
					RoleID:   "Operator",
					Links: asresponse.Links{
						Role: asresponse.Role{
							OdataID: "/redfish/v1/AccountService/Roles/Operator",
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
					AccountID:   "operatorUser",
				},
				session: &asmodel.Session{
					ID:       "operatorUser",
					UserName: "operatorUser",
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
					"Link":     "</redfish/v1/AccountService/Accounts/operatorUser/>; rel=describedby",
					"Location": "/redfish/v1/AccountService/Accounts/operatorUser",
				},
				Body: asresponse.Account{
					Response: operatorSuccessResponse,
					UserName: "operatorUser",
					RoleID:   "Operator",
					Links: asresponse.Links{
						Role: asresponse.Role{
							OdataID: "/redfish/v1/AccountService/Roles/Operator",
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
					"Link":     "</redfish/v1/AccountService/Accounts/testUser2/>; rel=describedby",
					"Location": "/redfish/v1/AccountService/Accounts/testUser2",
				},
				Body: asresponse.Account{
					Response: successResponse2,
					UserName: "testUser2",
					RoleID:   "Administrator",
					Links: asresponse.Links{
						Role: asresponse.Role{
							OdataID: "/redfish/v1/AccountService/Roles/Administrator",
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
					"Link":     "</redfish/v1/AccountService/Accounts/testUser2/>; rel=describedby",
					"Location": "/redfish/v1/AccountService/Accounts/testUser2",
				},
				Body: asresponse.Account{
					Response: successResponse2,
					UserName: "testUser2",
					RoleID:   "Administrator",
					Links: asresponse.Links{
						Role: asresponse.Role{
							OdataID: "/redfish/v1/AccountService/Roles/Administrator",
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
				Body:          errArgs.CreateGenericErrorResponse(),
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
				Body:          errArg4.CreateGenericErrorResponse(),
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
				Body:          errArgs1.CreateGenericErrorResponse(),
			},
		},
		{
			name: "update account without payload",
			args: args{
				req: &accountproto.UpdateAccountRequest{
					RequestBody: emptyPayload,
					AccountID:   "testUser1",
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureSelf: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusBadRequest,
				StatusMessage: response.PropertyMissing,
				Body:          errArg5.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := acc.Update(ctx, tt.args.req, tt.args.session)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createMockUpdateResponseObject(odataType, odataID, odataContext, ID string) response.Response {
	successResponse := response.Response{
		OdataType:    odataType,
		OdataID:      odataID,
		OdataContext: odataContext,
		ID:           ID,
		Name:         "Account Service",
	}
	successResponse.CreateGenericResponse(response.AccountModified)
	return successResponse
}
