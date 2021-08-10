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
func TestCreate(t *testing.T) {
	acc := getMockExternalInterface()
	common.SetUpMockConfig()
	errArgs := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InsufficientPrivilege,
				ErrorMessage:  "User does not have the privilege to create a new user",
				MessageArgs:   []interface{}{},
			},
		},
	}
	errArg := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.PropertyValueFormatError,
				ErrorMessage:  "error: invalid password, password length is less than the minimum length",
				MessageArgs:   []interface{}{"Password", "Password"},
			},
		},
	}
	errArg2 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.PropertyValueFormatError,
				ErrorMessage:  "error: invalid password, username is present inside the password",
				MessageArgs:   []interface{}{"testUser4", "Password"},
			},
		},
	}
	errArg3 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.PropertyValueFormatError,
				ErrorMessage:  "error: invalid password, password length is greater than the maximum length",
				MessageArgs:   []interface{}{"Password1234567890", "Password"},
			},
		},
	}
	errArg4 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.PropertyValueFormatError,
				ErrorMessage:  "error: invalid password, password should contain minimum One Upper case, One Lower case, One Number and One Special character",
				MessageArgs:   []interface{}{"password@123", "Password"},
			},
		},
	}
	errArg5 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.PropertyValueFormatError,
				ErrorMessage:  "error: invalid password, password should contain minimum One Upper case, One Lower case, One Number and One Special character",
				MessageArgs:   []interface{}{"PASSWORD@123", "Password"},
			},
		},
	}
	errArg6 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.PropertyValueFormatError,
				ErrorMessage:  "error: invalid password, password should contain minimum One Upper case, One Lower case, One Number and One Special character",
				MessageArgs:   []interface{}{"Password@ABC", "Password"},
			},
		},
	}
	errArg7 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.PropertyValueFormatError,
				ErrorMessage:  "error: invalid password, password should contain minimum One Upper case, One Lower case, One Number and One Special character",
				MessageArgs:   []interface{}{"P\\assword123", "Password"},
			},
		},
	}
	errArg1 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.PropertyMissing,
				ErrorMessage:  "Mandatory fields UserName Password RoleID are empty",
				MessageArgs:   []interface{}{"UserName Password RoleID"},
			},
		},
	}
	errArgu := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceAlreadyExists,
				ErrorMessage:  "Unable to add new user: error: data with key existingUser already exists",
				MessageArgs:   []interface{}{"ManagerAccount", "Id", "existingUser"},
			},
		},
	}
	errArgum := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "Invalid RoleID present error while trying to get role details: error: Invalid RoleID xyz present",
				MessageArgs:   []interface{}{"Role", "xyz"},
			},
		},
	}
	type args struct {
		req     *accountproto.CreateAccountRequest
		session *asmodel.Session
	}
	reqBodyValidAcc, _ := json.Marshal(asmodel.Account{
		UserName: "testUser",
		Password: "Password@123",
		RoleID:   "Administrator",
	})
	reqBodyInvalidRole, _ := json.Marshal(asmodel.Account{
		UserName: "testUser1",
		Password: "Password@123",
		RoleID:   "xyz",
	})
	reqBodyExistingAcc, _ := json.Marshal(asmodel.Account{
		UserName: "existingUser",
		Password: "Password@123",
		RoleID:   "Administrator",
	})
	reqBodyInvalidPrivilegeAcc, _ := json.Marshal(asmodel.Account{
		UserName: "testUser3",
		Password: "Password@123",
		RoleID:   "Administrator",
	})
	reqBodyInvalidData, _ := json.Marshal(asmodel.Account{
		UserName: "",
		Password: "",
		RoleID:   "",
	})
	reqBodyInvalidPwd, _ := json.Marshal(asmodel.Account{
		UserName: "testUser4",
		Password: "Password",
		RoleID:   "Administrator",
	})
	reqBodySameUserPwd, _ := json.Marshal(asmodel.Account{
		UserName: "testUser4",
		Password: "testUser4",
		RoleID:   "Administrator",
	})
	reqBodyPwdLenExceed, _ := json.Marshal(asmodel.Account{
		UserName: "testUser4",
		Password: "Password1234567890",
		RoleID:   "Administrator",
	})
	reqBodyNoUpperCasePwd, _ := json.Marshal(asmodel.Account{
		UserName: "testUser4",
		Password: "password@123",
		RoleID:   "Administrator",
	})
	reqBodyNoLowerCasePwd, _ := json.Marshal(asmodel.Account{
		UserName: "testUser4",
		Password: "PASSWORD@123",
		RoleID:   "Administrator",
	})
	reqBodyNoNumPwd, _ := json.Marshal(asmodel.Account{
		UserName: "testUser4",
		Password: "Password@ABC",
		RoleID:   "Administrator",
	})
	reqBodyInvalidSpeCharPwd, _ := json.Marshal(asmodel.Account{
		UserName: "testUser4",
		Password: "P\\assword123",
		RoleID:   "Administrator",
	})

	successResponse := response.Response{
		OdataType:    common.ManagerAccountType,
		OdataID:      "/redfish/v1/AccountService/Accounts/testUser",
		OdataContext: "/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
		ID:           "testUser",
		Name:         "Account Service",
	}
	successResponse.CreateGenericResponse(response.Created)
	tests := []struct {
		name    string
		args    args
		want    response.RPC
		wantErr bool
	}{
		{
			name: "successful account creation",
			args: args{
				req: &accountproto.CreateAccountRequest{
					RequestBody: reqBodyValidAcc,
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusCreated,
				StatusMessage: response.Created,
				Header: map[string]string{
					"Cache-Control":     "no-cache",
					"Connection":        "keep-alive",
					"Content-type":      "application/json; charset=utf-8",
					"Link":              "</redfish/v1/AccountService/Accounts/testUser/>; rel=describedby",
					"Location":          "/redfish/v1/AccountService/Accounts/testUser",
					"Transfer-Encoding": "chunked",
					"OData-Version":     "4.0",
				},
				Body: asresponse.Account{
					Response:     successResponse,
					UserName:     "testUser",
					RoleID:       "Administrator",
					AccountTypes: []string{"Redfish"},
					Links: asresponse.Links{
						Role: asresponse.Role{
							OdataID: "/redfish/v1/AccountService/Roles/Administrator",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "request body with invalid role",
			args: args{
				req: &accountproto.CreateAccountRequest{
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
				StatusMessage: response.ResourceNotFound,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArgum.CreateGenericErrorResponse(),
			},
			wantErr: true,
		},
		{
			name: "request for creating an existing user",
			args: args{
				req: &accountproto.CreateAccountRequest{
					RequestBody: reqBodyExistingAcc,
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusConflict,
				StatusMessage: response.ResourceAlreadyExists,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArgu.CreateGenericErrorResponse(),
			},
			wantErr: true,
		},
		{
			name: "create request with invalid privilege",
			args: args{
				req: &accountproto.CreateAccountRequest{
					RequestBody: reqBodyInvalidPrivilegeAcc,
				},
				session: &asmodel.Session{
					Privileges: map[string]bool{
						"ThisIsAnInvalidPrivilege": true,
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
			wantErr: true,
		},
		{
			name: "create request with invalid data",
			args: args{
				req: &accountproto.CreateAccountRequest{
					RequestBody: reqBodyInvalidData,
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
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArg1.CreateGenericErrorResponse(),
			},
			wantErr: true,
		},
		{
			name: "create request with invalid password length",
			args: args{
				req: &accountproto.CreateAccountRequest{
					RequestBody: reqBodyInvalidPwd,
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
				Body: errArg.CreateGenericErrorResponse(),
			},
			wantErr: true,
		},
		{
			name: "create request with invalid password with username in password",
			args: args{
				req: &accountproto.CreateAccountRequest{
					RequestBody: reqBodySameUserPwd,
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
				Body: errArg2.CreateGenericErrorResponse(),
			},
			wantErr: true,
		},
		{
			name: "create request with password exceeding length",
			args: args{
				req: &accountproto.CreateAccountRequest{
					RequestBody: reqBodyPwdLenExceed,
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
			wantErr: true,
		},
		{
			name: "create request with invalid password with no uppercase character",
			args: args{
				req: &accountproto.CreateAccountRequest{
					RequestBody: reqBodyNoUpperCasePwd,
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
				Body: errArg4.CreateGenericErrorResponse(),
			},
			wantErr: true,
		},
		{
			name: "create request with invalid password with no lowercase character",
			args: args{
				req: &accountproto.CreateAccountRequest{
					RequestBody: reqBodyNoLowerCasePwd,
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
				Body: errArg5.CreateGenericErrorResponse(),
			},
			wantErr: true,
		},
		{
			name: "create request with invalid password with no number",
			args: args{
				req: &accountproto.CreateAccountRequest{
					RequestBody: reqBodyNoNumPwd,
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
				Body: errArg6.CreateGenericErrorResponse(),
			},
			wantErr: true,
		},
		{
			name: "create request with invalid password with invalid special character",
			args: args{
				req: &accountproto.CreateAccountRequest{
					RequestBody: reqBodyInvalidSpeCharPwd,
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
				Body: errArg7.CreateGenericErrorResponse(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := acc.Create(tt.args.req, tt.args.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
