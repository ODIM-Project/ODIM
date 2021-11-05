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
package session

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
	"golang.org/x/crypto/sha3"
)

func createMockUser(username, roleID string) error {
	hash := sha3.New512()
	hash.Write([]byte("P@$$w0rd"))
	hashSum := hash.Sum(nil)
	hashedPassword := base64.URLEncoding.EncodeToString(hashSum)
	user := asmodel.User{
		UserName: username,
		Password: hashedPassword,
		RoleID:   roleID,
	}
	if err := asmodel.CreateUser(user); err != nil {
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

func TestCreateSession(t *testing.T) {
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
	auth.Lock.Lock()
	common.SetUpMockConfig()
	auth.Lock.Unlock()
	if err := createMockRole(common.RoleAdmin, []string{common.PrivilegeConfigureManager, common.PrivilegeLogin}, []string{}); err != nil {
		t.Fatalf("Error while creating role: %v", err)
		return
	}
	if err := createMockRole("Sample", []string{common.PrivilegeConfigureManager}, []string{}); err != nil {
		t.Fatalf("Error while creating role: %v", err)
		return
	}
	if err := createMockUser("admin", common.RoleAdmin); err != nil {
		t.Fatalf("Error while creating account: %v", err)
		return
	}
	if err := createMockUser("sample", "Sample"); err != nil {
		t.Fatalf("Error while creating account: %v", err)
		return
	}
	commonResponse := response.Response{
		OdataType: common.SessionServiceType,
		OdataID:   "/redfish/v1/SessionService/Sessions",
		ID:        "Sessions",
		Name:      "Session Service",
	}
	errArgUnauth := &response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.NoValidSession,
				ErrorMessage:  "Unable to authorize session creation credentials: error: username or password missing",
				MessageArgs:   []interface{}{},
			},
		},
	}

	errArgUnauth1 := &response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.NoValidSession,
				ErrorMessage:  "Unable to authorize session creation credentials: error: password mismatch ",
				MessageArgs:   []interface{}{},
			},
		},
	}
	errArg2 := &response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InsufficientPrivilege,
				ErrorMessage:  "User doesn't have required privilege to create a session",
				MessageArgs:   []interface{}{},
			},
		},
	}
	type args struct {
		req *sessionproto.SessionCreateRequest
	}

	reqBodyCreateSession, _ := json.Marshal(asmodel.CreateSession{
		UserName: "admin",
		Password: "P@$$w0rd",
	})
	reqBodyInvalidCred, _ := json.Marshal(asmodel.CreateSession{
		UserName: "admin",
		Password: "HP1",
	})
	reqBodyNoPrivilege, _ := json.Marshal(asmodel.CreateSession{
		UserName: "sample",
		Password: "P@$$w0rd",
	})
	reqBodyEmpty, _ := json.Marshal(asmodel.CreateSession{})

	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "successful creation of session",
			args: args{
				req: &sessionproto.SessionCreateRequest{
					RequestBody: reqBodyCreateSession,
				},
			},
		},
		// TODO: need to correct this test cases when the error response is getting corrected
		{
			name: "Create session without user",
			args: args{
				req: &sessionproto.SessionCreateRequest{
					RequestBody: reqBodyEmpty,
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusUnauthorized,
				StatusMessage: response.NoValidSession,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArgUnauth.CreateGenericErrorResponse(),
			},
		},
		{
			name: "create session with an invalid credentials",
			args: args{
				req: &sessionproto.SessionCreateRequest{
					RequestBody: reqBodyInvalidCred,
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusUnauthorized,
				StatusMessage: response.NoValidSession,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArgUnauth1.CreateGenericErrorResponse(),
			},
		},
		{
			name: "create session with an no login privilege",
			args: args{
				req: &sessionproto.SessionCreateRequest{
					RequestBody: reqBodyNoPrivilege,
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.InsufficientPrivilege,
				Header: map[string]string{
					"Content-type": "application/json; charset=utf-8",
				},
				Body: errArg2.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, tokenID := CreateNewSession(tt.args.req)
			if tokenID != "" { // success case
				commonResponse.ID = tokenID
				commonResponse.OdataID = "/redfish/v1/SessionService/Sessions/" + tokenID
				commonResponse.CreateGenericResponse(response.Created)
				tt.want = response.RPC{
					StatusCode:    http.StatusCreated,
					StatusMessage: response.Created,
					Header:        got.Header,
					Body: asresponse.Session{
						Response: commonResponse,
						UserName: "admin",
					},
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateNewSession() = %v, want %v", got, tt.want)
			}
		})
	}
}
