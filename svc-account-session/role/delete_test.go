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
	"encoding/base64"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"golang.org/x/crypto/sha3"
)

func mockSession(token string, roleID string, privilegeRequired bool) error {
	currentTime := time.Now()
	session := asmodel.Session{
		ID:           "id",
		Token:        token,
		UserName:     "admin",
		RoleID:       roleID,
		CreatedTime:  currentTime,
		LastUsedTime: currentTime,
	}
	if privilegeRequired {
		session.Privileges = map[string]bool{
			common.PrivilegeConfigureUsers: true,
		}
	}

	if err := session.Persist(); err != nil {
		return err
	}
	return nil
}

func truncateDB(t *testing.T) {
	err := common.TruncateDB(common.OnDisk)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = common.TruncateDB(common.InMemory)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

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

func TestDelete(t *testing.T) {
	common.SetUpMockConfig()
	defer truncateDB(t)
	token, tokenWithoutPrivilege := "someToken", "tokenWithoutPrivilege"
	err := mockSession(token, common.RoleAdmin, true)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = mockSession(tokenWithoutPrivilege, common.RoleAdmin, false)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	testRole := "someRole"
	err = createMockRole(testRole, []string{common.PrivilegeConfigureUsers, common.PrivilegeLogin}, []string{}, false)
	if err != nil {
		t.Fatalf("Error in creating mock admin user %v", err)
	}
	testUserRole := "someUserDefinedRole"
	err = createMockRole(testUserRole, []string{common.PrivilegeConfigureUsers}, []string{}, false)
	if err != nil {
		t.Fatalf("Error in creating mock admin user %v", err)
	}
	err = createMockUser("testUser1", testUserRole)
	if err != nil {
		t.Fatalf("Error in creating mock admin user %v", err)
	}
	testPredefinedRole := "someOtherRole"
	err = createMockRole(testPredefinedRole, []string{common.PrivilegeConfigureUsers}, []string{}, true)
	if err != nil {
		t.Fatalf("Error in creating mock admin user %v", err)
	}
	errArg := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InsufficientPrivilege,
				ErrorMessage:  "The session token doesn't have required privilege",
				MessageArgs:   []interface{}{},
			},
		},
	}
	errArgs := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "Unable to get role details: error while trying to get role details: no data with the with key xyz found",
				MessageArgs:   []interface{}{"Role", "xyz"},
			},
		},
	}
	errArgsu := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.NoValidSession,
				ErrorMessage:  "Unable to authorize session token: error while trying to get session details with the token invalid-token: error while trying to get the session from DB: no data with the with key invalid-token found",
				MessageArgs:   []interface{}{},
			},
		},
	}
	errArgu := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InsufficientPrivilege,
				ErrorMessage:  "A predefined role cannot be deleted.",
				MessageArgs:   []interface{}{},
			},
		},
	}
	errArgus := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceInUse,
				ErrorMessage:  "Role is assigned to a user",
				MessageArgs:   []interface{}{},
			},
		},
	}
	type args struct {
		req *roleproto.DeleteRoleRequest
	}
	tests := []struct {
		name string
		args args
		want *response.RPC
	}{
		{
			name: "successful deletion of the role",
			args: args{
				req: &roleproto.DeleteRoleRequest{
					SessionToken: token,
					ID:           testRole,
				},
			},
			want: &response.RPC{
				StatusCode:    http.StatusNoContent,
				StatusMessage: response.ResourceRemoved,
			},
		},
		{
			name: "deletion of the role without valid token",
			args: args{
				req: &roleproto.DeleteRoleRequest{
					SessionToken: "invalid-token",
					ID:           testRole,
				},
			},
			want: &response.RPC{
				StatusCode:    http.StatusUnauthorized,
				StatusMessage: response.NoValidSession,
				Body:          errArgsu.CreateGenericErrorResponse(),
			},
		},
		{
			name: "deletion of the role without privileges",
			args: args{
				req: &roleproto.DeleteRoleRequest{
					SessionToken: tokenWithoutPrivilege,
					ID:           testRole,
				},
			},
			want: &response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.InsufficientPrivilege,
				Body:          errArg.CreateGenericErrorResponse(),
			},
		},
		{
			name: "deletion of the invalid role",
			args: args{
				req: &roleproto.DeleteRoleRequest{
					SessionToken: token,
					ID:           "xyz",
				},
			},
			want: &response.RPC{
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Body:          errArgs.CreateGenericErrorResponse(),
			},
		},
		{
			name: "deletion of the predefined role",
			args: args{
				req: &roleproto.DeleteRoleRequest{
					SessionToken: token,
					ID:           testPredefinedRole,
				},
			},
			want: &response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.InsufficientPrivilege,
				Body:          errArgu.CreateGenericErrorResponse(),
			},
		},
		{
			name: "deletion of the assigned role",
			args: args{
				req: &roleproto.DeleteRoleRequest{
					SessionToken: token,
					ID:           testUserRole,
				},
			},
			want: &response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.ResourceInUse,
				Body:          errArgus.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Delete(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}
