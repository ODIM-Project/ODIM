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
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
)

func TestGetAllAccounts(t *testing.T) {
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
		OdataType:    "#ManagerAccountCollection.ManagerAccountCollection",
		OdataID:      "/redfish/v1/AccountService/Accounts",
		OdataContext: "/redfish/v1/$metadata#ManagerAccountCollection.ManagerAccountCollection",
		ID:           "Accounts",
		Name:         "Account Service",
	}
	successResponse.CreateGenericResponse(response.Success)
	successResponse.Message = ""
	successResponse.ID = ""
	successResponse.MessageID = ""
	successResponse.Severity = ""

	err := createMockUser("testUser1", common.RoleAdmin)
	if err != nil {
		t.Fatalf("Error in creating mock admin user %v", err)
	}

	errArgs := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InsufficientPrivilege,
				ErrorMessage:  "User SomeOne does not have the privilege to view all users",
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
			name: "successful get on all accounts",
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
				Header: map[string]string{
					"Allow":             `"GET", "POST", "HEAD"`,
					"Cache-Control":     "no-cache",
					"Connection":        "keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Link":              "</redfish/v1/SchemaStore/en/ManagerAccountCollection.json/>; rel=describedby",
					"Transfer-Encoding": "chunked",
					"OData-Version":     "4.0",
				},
				Body: asresponse.List{
					Response:     successResponse,
					MembersCount: 1,
					Members: []asresponse.ListMember{
						asresponse.ListMember{
							OdataID: "/redfish/v1/AccountService/Accounts/testUser1",
						},
					},
				},
			},
		},
		{
			name: "get on all accounts without privilege",
			args: args{
				session: &asmodel.Session{
					UserName: "SomeOne",
					Privileges: map[string]bool{
						"ThisIsSomePrivilege": true,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAllAccounts(tt.args.session)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllAccounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAccount(t *testing.T) {
	successResponse := response.Response{
		OdataType:    common.ManagerAccountType,
		OdataID:      "/redfish/v1/AccountService/Accounts/testUser1",
		OdataContext: "/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
		ID:           "testUser1",
		Name:         "Account Service",
	}
	successResponse.CreateGenericResponse(response.Success)
	successResponse.Message = ""
	successResponse.MessageID = ""
	successResponse.Severity = ""
	common.SetUpMockConfig()
	err := createMockUser("testUser1", common.RoleAdmin)
	if err != nil {
		t.Fatalf("Error in creating mock admin user %v", err)
	}
	type args struct {
		session   *asmodel.Session
		accountID string
	}

	errArg := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InsufficientPrivilege,
				ErrorMessage:  "testUser2 does not have the privilege to view other user's details",
				MessageArgs:   []interface{}{},
			},
		},
	}
	errArg1 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "Unable to get account: error while trying to get user: no data with the with key testUser4 found",
				MessageArgs:   []interface{}{"Account", "testUser4"},
			},
		},
	}

	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "successful get account with configureuser privilege",
			args: args{
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
				accountID: "testUser1",
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Allow":             `"GET", "POST", "HEAD"`,
					"Cache-Control":     "no-cache",
					"Connection":        "keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Link":              "</redfish/v1/SchemaStore/en/ManagerAccount.json/>; rel=describedby",
					"Transfer-Encoding": "chunked",
					"OData-Version":     "4.0",
				},
				Body: asresponse.Account{
					Response: successResponse,
					UserName: "testUser1",
					RoleID:   "Administrator",
					Links: asresponse.Links{
						Role: asresponse.Role{
							OdataID: "/redfish/v1/AccountService/Roles/Administrator"},
					},
				},
			},
		},
		{
			name: "successful get account with configureself privilege",
			args: args{
				session: &asmodel.Session{
					UserName: "testUser1",
					Privileges: map[string]bool{
						common.PrivilegeConfigureSelf: true,
					},
				},
				accountID: "testUser1",
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Allow":             `"GET", "POST", "HEAD"`,
					"Cache-Control":     "no-cache",
					"Connection":        "keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Link":              "</redfish/v1/SchemaStore/en/ManagerAccount.json/>; rel=describedby",
					"Transfer-Encoding": "chunked",
					"OData-Version":     "4.0",
				},
				Body: asresponse.Account{
					Response: successResponse,
					UserName: "testUser1",
					RoleID:   "Administrator",
					Links: asresponse.Links{
						Role: asresponse.Role{
							OdataID: "/redfish/v1/AccountService/Roles/Administrator"},
					},
				},
			},
		},
		{
			name: "get account on other account with configureself privilege",
			args: args{
				session: &asmodel.Session{
					UserName: "testUser2",
					Privileges: map[string]bool{
						common.PrivilegeConfigureSelf: true,
					},
				},
				accountID: "testUser1",
			},
			want: response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.InsufficientPrivilege,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArg.CreateGenericErrorResponse(),
			},
		},
		{
			name: "get account without privilege",
			args: args{
				session: &asmodel.Session{
					UserName: "testUser2",
					Privileges: map[string]bool{
						"ThisIsSomePrivilege": true,
					},
				},
				accountID: "testUser1",
			},
			want: response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.InsufficientPrivilege,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArg.CreateGenericErrorResponse(),
			},
		},
		{
			name: "get on non-existing account",
			args: args{
				session: &asmodel.Session{
					UserName: "testUser4",
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
				accountID: "testUser4",
			},
			want: response.RPC{
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArg1.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAccount(tt.args.session, tt.args.accountID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAccount() = %v, want %v", got, tt.want)
			}
		})
	}
	errs := common.TruncateDB(common.OnDisk)
	if errs != nil {
		t.Fatalf("error: %v", errs)
	}
}

func TestGetAccountService(t *testing.T) {
	successResponse := response.Response{
		OdataType:    common.AccountServiceType,
		OdataID:      "/redfish/v1/AccountService",
		OdataContext: "/redfish/v1/$metadata#AccountService.AccountService",
		ID:           "AccountService",
		Name:         "Account Service",
	}
	successResponse.CreateGenericResponse(response.Success)
	successResponse.Message = ""
	successResponse.MessageID = ""
	successResponse.Severity = ""
	common.SetUpMockConfig()
	tests := []struct {
		name string
		want response.RPC
	}{
		{
			name: "account service enabled",
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Allow":         "GET",
					"Cache-Control": "no-cache",
					"Connection":    "Keep-alive",
					"Content-type":  "application/json; charset=utf-8",
					"Link": "	</redfish/v1/SchemaStore/en/AccountService.json>; rel=describedby",
					"Transfer-Encoding": "chunked",
					"X-Frame-Options":   "sameorigin",
				},
				Body: asresponse.AccountService{
					Response: successResponse,
					Status: asresponse.Status{
						State:  "Enabled",
						Health: "OK",
					},
					ServiceEnabled:    true,
					MinPasswordLength: config.Data.AuthConf.PasswordRules.MinPasswordLength,
					Accounts: asresponse.Accounts{
						OdataID: "/redfish/v1/AccountService/Accounts",
					},
					Roles: asresponse.Accounts{
						OdataID: "/redfish/v1/AccountService/Roles",
					},
				},
			},
		},
		{
			name: "account service disabled",
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Allow":         "GET",
					"Cache-Control": "no-cache",
					"Connection":    "Keep-alive",
					"Content-type":  "application/json; charset=utf-8",
					"Link": "	</redfish/v1/SchemaStore/en/AccountService.json>; rel=describedby",
					"Transfer-Encoding": "chunked",
					"X-Frame-Options":   "sameorigin",
				},
				Body: asresponse.AccountService{
					Response: successResponse,
					Status: asresponse.Status{
						State:  "Disabled",
						Health: "OK",
					},
					ServiceEnabled:    false,
					MinPasswordLength: config.Data.AuthConf.PasswordRules.MinPasswordLength,
					Accounts: asresponse.Accounts{
						OdataID: "/redfish/v1/AccountService/Accounts",
					},
					Roles: asresponse.Accounts{
						OdataID: "/redfish/v1/AccountService/Roles",
					},
				},
			},
		},
	}
	config.Data.EnabledServices = []string{"AccountService"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAccountService()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAccountService() = %v, want %v", got, tt.want)
			}
		})
		config.Data.EnabledServices = []string{"XXXX"}
	}
}
