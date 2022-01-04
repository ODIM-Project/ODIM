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
	"encoding/base64"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
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

func TestDelete(t *testing.T) {
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

	common.SetUpMockConfig()

	errArgs := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "Unable to delete user: no data with the with key xyz found",
				MessageArgs:   []interface{}{"Account", "xyz"},
			},
		},
	}
	errArg := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.InsufficientPrivilege,
				ErrorMessage:  "SomeOne does not have the privilege to delete user",
				MessageArgs:   []interface{}{},
			},
		},
	}

	type args struct {
		session   *asmodel.Session
		accountID string
	}
	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "successful deletion of account",
			args: args{
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
				accountID: "0",
			},
			want: response.RPC{
				StatusCode:    http.StatusNoContent,
				StatusMessage: response.AccountRemoved,
			},
		},
		{
			name: "delete non-existing account",
			args: args{
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
				accountID: "xyz",
			},
			want: response.RPC{
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Body:          errArgs.CreateGenericErrorResponse(),
			},
		},
		{
			name: "delete account without privileges",
			args: args{
				session: &asmodel.Session{
					UserName: "SomeOne",
					Privileges: map[string]bool{
						"someprivliege": true,
					},
				},
				accountID: "2",
			},
			want: response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.InsufficientPrivilege,
				Body:          errArg.CreateGenericErrorResponse(),
			},
		},
	}
	for index, tt := range tests {
		err := createMockUser(strconv.Itoa(index), common.RoleAdmin)
		if err != nil {
			t.Fatalf("Error in creating mock admin user %v", err)
		}
		t.Run(tt.name, func(t *testing.T) {
			got := Delete(tt.args.session, tt.args.accountID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteDefaultAdminAccount(t *testing.T) {
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	common.SetUpMockConfig()

	errArgs := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceCannotBeDeleted,
				ErrorMessage:  "default user account can not be deleted",
			},
		},
	}
	err := createMockUser(defaultAdminAccount, common.RoleAdmin)
	if err != nil {
		t.Fatalf("Error in creating mock admin user %v", err)
	}
	type args struct {
		session   *asmodel.Session
		accountID string
	}
	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "deletion of default admin user account",
			args: args{
				session: &asmodel.Session{
					Privileges: map[string]bool{
						common.PrivilegeConfigureUsers: true,
					},
				},
				accountID: "admin",
			},
			want: response.RPC{
				StatusCode:    http.StatusBadRequest,
				StatusMessage: response.ResourceCannotBeDeleted,
				Body:          errArgs.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got := Delete(tt.args.session, tt.args.accountID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}
