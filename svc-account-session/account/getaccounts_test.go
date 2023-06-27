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
	"context"
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
	config.SetUpMockConfig(t)
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

	successResponse := createMockResponseObject("#ManagerAccountCollection.ManagerAccountCollection", "/redfish/v1/AccountService/Accounts", "/redfish/v1/$metadata#ManagerAccountCollection.ManagerAccountCollection", "Accounts")
	successResponse.ID = ""

	err := createMockUser("testUser1", common.RoleAdmin)
	if err != nil {
		t.Fatalf("Error in creating mock admin user %v", err)
	}

	errArgs := GetResponseArgs(response.InsufficientPrivilege, "failed to fetch accounts : User SomeOne does not have the privilege to view all users", []interface{}{})
	ctx := mockContext()
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
					"Link": "</redfish/v1/SchemaStore/en/ManagerAccountCollection.json/>; rel=describedby",
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
				Body:          errArgs.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAllAccounts(ctx, tt.args.session)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllAccounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAccount(t *testing.T) {
	successResponse := createMockResponseObject(common.ManagerAccountType, "/redfish/v1/AccountService/Accounts/testUser1", "/redfish/v1/$metadata#ManagerAccount.ManagerAccount", "testUser1")
	config.SetUpMockConfig(t)
	err := createMockUser("testUser1", common.RoleAdmin)
	if err != nil {
		t.Fatalf("Error in creating mock admin user %v", err)
	}
	ctx := mockContext()
	type args struct {
		session   *asmodel.Session
		accountID string
	}

	errArg := GetResponseArgs(response.InsufficientPrivilege, "failed to fetch the account testUser1: testUser2 does not have the privilege to view other user's details", []interface{}{})

	errArg1 := GetResponseArgs(response.ResourceNotFound, "failed to fetch the account testUser4: error while trying to get user: no data with the with key testUser4 found", []interface{}{"Account", "testUser4"})

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
					"Link": "</redfish/v1/SchemaStore/en/ManagerAccount.json/>; rel=describedby",
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
					"Link": "</redfish/v1/SchemaStore/en/ManagerAccount.json/>; rel=describedby",
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
				Body:          errArg.CreateGenericErrorResponse(),
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
				Body:          errArg.CreateGenericErrorResponse(),
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
				Body:          errArg1.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAccount(ctx, tt.args.session, tt.args.accountID)
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
	successResponse := createMockResponseObject(common.AccountServiceType, "/redfish/v1/AccountService", "/redfish/v1/$metadata#AccountService.AccountService", "AccountService")
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
					"Link": "	</redfish/v1/SchemaStore/en/AccountService.json>; rel=describedby",
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
					"Link": "	</redfish/v1/SchemaStore/en/AccountService.json>; rel=describedby",
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
			got := GetAccountService(context.TODO())
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAccountService() = %v, want %v", got, tt.want)
			}
		})
		config.Data.EnabledServices = []string{"XXXX"}
	}
}

func createMockResponseObject(odataType, odataID, odataContext, ID string) response.Response {
	successResponse := response.Response{
		OdataType:    odataType,
		OdataID:      odataID,
		OdataContext: odataContext,
		ID:           ID,
		Name:         "Account Service",
	}
	successResponse.CreateGenericResponse(response.Success)
	successResponse.Message = ""
	successResponse.MessageID = ""
	successResponse.Severity = ""
	return successResponse
}
