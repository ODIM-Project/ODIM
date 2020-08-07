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
package auth

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	authproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/auth"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
)

func createSession(token, id string, privileges map[string]bool, createdTime, lastUsedTime time.Time) error {
	session := asmodel.Session{
		Token:        token,
		ID:           id,
		Privileges:   privileges,
		CreatedTime:  createdTime,
		LastUsedTime: lastUsedTime,
	}
	if err := session.Persist(); err != nil {
		return err
	}
	return nil
}

func TestAuth(t *testing.T) {
	Lock.Lock()
	common.SetUpMockConfig()
	Lock.Unlock()
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
	sessionToken, id := "validToken", "someID"
	privilegesMap := map[string]bool{
		common.PrivilegeConfigureUsers: true,
	}
	currentTime := time.Now()
	err := createSession(sessionToken, id, privilegesMap, currentTime, currentTime)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	// positive test case privilege
	privileges := []string{common.PrivilegeConfigureUsers}

	//negative test case privilege
	invalidPrivilege := []string{common.PrivilegeConfigureManager}

	var oemPrivileges []string

	type args struct {
		req *authproto.AuthRequest
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 string
	}{

		{
			name: "Authorized token",
			args: args{
				req: &authproto.AuthRequest{
					SessionToken:  sessionToken,
					Privileges:    privileges,
					Oemprivileges: oemPrivileges,
				},
			},
			want:  http.StatusOK,
			want1: response.Success,
		},
		{
			name: "Invalid token",
			args: args{
				req: &authproto.AuthRequest{
					SessionToken:  "d8a06327-45f0-4176-b688-b80c1849f931",
					Privileges:    privileges,
					Oemprivileges: oemPrivileges,
				},
			},
			want:  http.StatusUnauthorized,
			want1: response.NoValidSession,
		},
		{
			name: "Empty token",
			args: args{
				req: &authproto.AuthRequest{
					SessionToken:  "",
					Privileges:    privileges,
					Oemprivileges: oemPrivileges,
				},
			},
			want:  http.StatusUnauthorized,
			want1: response.NoValidSession,
		},
		{
			name: "Invalid privileges",
			args: args{
				req: &authproto.AuthRequest{
					SessionToken:  sessionToken,
					Privileges:    invalidPrivilege,
					Oemprivileges: oemPrivileges,
				},
			},
			want:  http.StatusForbidden,
			want1: response.InsufficientPrivilege,
		},
		{
			name: "without privileges",
			args: args{
				req: &authproto.AuthRequest{
					SessionToken:  sessionToken,
					Privileges:    []string{},
					Oemprivileges: oemPrivileges,
				},
			},
			want:  http.StatusUnauthorized,
			want1: response.NoValidSession,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Auth(tt.args.req)
			if !reflect.DeepEqual(got, int32(tt.want)) {
				t.Errorf("Auth() = %v, want = %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Auth() = %v, want1 = %v", got1, tt.want1)
			}
		})
	}

}
