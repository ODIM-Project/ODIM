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
package session

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/account"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
)

func createSession(t *testing.T, role, username string, privileges []string) (string, string) {
	auth.Lock.Lock()
	common.SetUpMockConfig()
	auth.Lock.Unlock()
	errs := createMockRole(role, privileges, []string{})
	if errs != nil {
		t.Fatalf("Error in creating mock admin user %v", errs)
	}
	errs = createMockUser(username, role)
	if errs != nil {
		t.Fatalf("Error in creating mock admin user %v", errs)
	}

	reqBodyBytes, _ := json.Marshal(asmodel.CreateSession{
		UserName: username,
		Password: "P@$$w0rd",
	})
	req := &sessionproto.SessionCreateRequest{
		RequestBody: reqBodyBytes,
	}

	resp, sessionID := CreateNewSession(context.TODO(), req)
	if sessionID == "" {
		t.Fatalf("Session creation failed: %#v", resp)
	}

	return sessionID, resp.Header["X-Auth-Token"]
}
func TestDeleteSession(t *testing.T) {
	config.SetUpMockConfig(t)
	sessionID, sessionToken := createSession(t, common.RoleAdmin, "admin", []string{common.PrivilegeConfigureUsers, common.PrivilegeLogin})
	sessionID2, sessionToken2 := createSession(t, common.RoleClient, "client", []string{common.PrivilegeLogin})
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
	ctx := mockContext()
	type args struct {
		req *sessionproto.SessionRequest
	}

	errArgUnauth := account.GetResponseArgs(response.NoValidSession, "failed to delete session : error while trying to get the session from DB: no data with the with key  found", []interface{}{})

	eArgs := account.GetResponseArgs(response.ResourceNotFound, "failed to delete session : Session ID not found", []interface{}{"Session", "invalid-id"})

	errArgIns := account.GetResponseArgs(response.InsufficientPrivilege, "failed to delete session : Insufficient privileges", []interface{}{})

	tests := []struct {
		name string
		args args
		want response.RPC
	}{

		{
			name: "session deletion with invalid X-Auth-Token",
			args: args{
				req: &sessionproto.SessionRequest{
					SessionId: sessionID,
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusUnauthorized,
				StatusMessage: response.NoValidSession,
				Body:          errArgUnauth.CreateGenericErrorResponse(),
			},
		},
		{
			name: "session deletion with invalid/expired session id",
			args: args{
				req: &sessionproto.SessionRequest{
					SessionId:    "invalid-id",
					SessionToken: sessionToken,
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Body:          eArgs.CreateGenericErrorResponse(),
			},
		},

		{
			name: "successful session deletion",
			args: args{
				req: &sessionproto.SessionRequest{
					SessionId:    sessionID,
					SessionToken: sessionToken,
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusNoContent,
				StatusMessage: response.ResourceRemoved,
			},
		},
		{
			name: "session deletion with insuffecient privileges",
			args: args{
				req: &sessionproto.SessionRequest{
					SessionId:    sessionID2,
					SessionToken: sessionToken2,
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.InsufficientPrivilege,
				Body:          errArgIns.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := DeleteSession(ctx, tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteSession() = %v, want %v", got, tt.want)
			}
		})
	}
}
