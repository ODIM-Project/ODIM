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
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
)

func TestGetSession(t *testing.T) {
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

	sessionID, sessionToken := createSession(t, common.RoleAdmin, "admin", []string{common.PrivilegeConfigureUsers, common.PrivilegeLogin})
	commonResponse := response.Response{
		OdataType: common.SessionType,
		OdataID:   "/redfish/v1/SessionService/Sessions/" + sessionID,
		ID:        sessionID,
		Name:      "User Session",
	}
	successHeader := map[string]string{
		"Cache-Control":     "no-cache",
		"Link":              "</redfish/v1/SessionService/Sessions/" + sessionID + "/>; rel=self",
		"Transfer-Encoding": "chunked",
		"X-Auth-Token":      sessionToken,
		"Content-type":      "application/json; charset=utf-8",
	}
	errArgUnauth := &response.Args{
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
	eArgs := &response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "No session with id invalid-sessionID found.",
				MessageArgs:   []interface{}{"Session", "invalid-sessionID"},
			},
		},
	}

	type args struct {
		req *sessionproto.SessionRequest
	}
	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "successful get session",
			args: args{
				req: &sessionproto.SessionRequest{
					SessionId:    sessionID,
					SessionToken: sessionToken,
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header:        successHeader,
				Body: asresponse.Session{
					Response: commonResponse,
					UserName: "admin",
				},
			},
		},
		{
			name: "get session with invalid session token",
			args: args{
				req: &sessionproto.SessionRequest{
					SessionId:    sessionID,
					SessionToken: "invalid-token",
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusUnauthorized,
				StatusMessage: response.NoValidSession,
				Header:        getHeader(),
				Body:          errArgUnauth.CreateGenericErrorResponse(),
			},
		},
		{
			name: "get session with invalid session id",
			args: args{
				req: &sessionproto.SessionRequest{
					SessionId:    "invalid-sessionID",
					SessionToken: sessionToken,
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusNotFound,
				StatusMessage: response.ResourceNotFound,
				Header:        getHeader(),
				Body:          eArgs.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSession(tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllActiveSessions(t *testing.T) {
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
	sessionID, sessionToken := createSession(t, common.RoleAdmin, "admin", []string{common.PrivilegeConfigureUsers, common.PrivilegeLogin})
	commonResponse := response.Response{
		OdataType:    "#SessionCollection.SessionCollection",
		OdataID:      "/redfish/v1/SessionService/Sessions",
		OdataContext: "/redfish/v1/$metadata#SessionCollection.SessionCollection",
		Name:         "Session Service",
	}
	var listMembers []asresponse.ListMember
	listMembers = append(listMembers, asresponse.ListMember{
		OdataID: "/redfish/v1/SessionService/Sessions/" + sessionID,
	})
	eArgs1 := &response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.NoValidSession,
				ErrorMessage:  "Unable to authorize session token: error: no session token found in header",
				MessageArgs:   []interface{}{},
			},
		},
	}
	errArgUnauth2 := &response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.NoValidSession,
				ErrorMessage:  "Unable to authorize session token: error while trying to get session details with the token invalidToken: error while trying to get the session from DB: no data with the with key invalidToken found",
				MessageArgs:   []interface{}{},
			},
		},
	}
	type args struct {
		req *sessionproto.SessionRequest
	}
	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "successful get all active session",
			args: args{
				req: &sessionproto.SessionRequest{
					SessionId:    sessionID,
					SessionToken: sessionToken,
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header:        getHeader(),
				Body: asresponse.List{
					Response:     commonResponse,
					MembersCount: len(listMembers),
					Members:      listMembers,
				},
			},
		},
		{
			name: "get all sessions with no session token",
			args: args{
				req: &sessionproto.SessionRequest{
					SessionId:    sessionID,
					SessionToken: "",
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusUnauthorized,
				StatusMessage: response.NoValidSession,
				Header:        getHeader(),
				Body:          eArgs1.CreateGenericErrorResponse(),
			},
		},
		{
			name: "get all sessions with invalid session token",
			args: args{
				req: &sessionproto.SessionRequest{
					SessionId:    sessionID,
					SessionToken: "invalidToken",
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusUnauthorized,
				StatusMessage: response.NoValidSession,
				Header: map[string]string{
					"Cache-Control":     "no-cache",
					"Transfer-Encoding": "chunked",
					"Content-type":      "application/json; charset=utf-8",
				},
				Body: errArgUnauth2.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAllActiveSessions(tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllActiveSessions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSessionService(t *testing.T) {
	commonResponse := response.Response{
		OdataType: common.SessionServiceType,
		OdataID:   "/redfish/v1/SessionService",
		ID:        "Sessions",
		Name:      "Session Service",
	}
	commonResponse.CreateGenericResponse(response.Success)
	commonResponse.MessageID = ""
	commonResponse.Message = ""
	commonResponse.Severity = ""

	type args struct {
		req *sessionproto.SessionRequest
	}
	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "successful get session service",
			args: args{
				req: &sessionproto.SessionRequest{},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Allow":         "GET",
					"Cache-Control": "no-cache",
					"Connection":    "Keep-alive",
					"Link": "	</redfish/v1/SchemaStore/en/SessionService.json>; rel=describedby",
					"Transfer-Encoding": "chunked",
					"X-Frame-Options":   "sameorigin",
					"Content-type":      "application/json; charset=utf-8",
				},
				Body: asresponse.SessionService{
					Response: commonResponse,

					Status: asresponse.Status{
						State:  "Enabled",
						Health: "OK",
					},
					ServiceEnabled: true,
					SessionTimeout: config.Data.AuthConf.SessionTimeOutInMins,
					Sessions: asresponse.Sessions{
						OdataID: "/redfish/v1/SessionService/Sessions",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSessionService(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSessionService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSessionUserName(t *testing.T) {
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
	sessionID, sessionToken := createSession(t, common.RoleAdmin, "admin", []string{common.PrivilegeConfigureUsers, common.PrivilegeLogin})
	type args struct {
		req  *sessionproto.SessionRequest
		resp *sessionproto.SessionUserName
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    error
	}{
		// TODO: Add test cases.
		{
			name: "successful get session username",
			args: args{
				req: &sessionproto.SessionRequest{
					SessionId:    sessionID,
					SessionToken: sessionToken,
				},
				resp: &sessionproto.SessionUserName{},
			},
			wantErr: false,
		},
		{
			name: "invaid session id",
			args: args{
				req: &sessionproto.SessionRequest{
					SessionId:    sessionID,
					SessionToken: "sessionToken",
				},
				resp: &sessionproto.SessionUserName{},
			},
			wantErr: true,
			want:    fmt.Errorf("error while trying to get session details with the token sessionToken: error while trying to get the session from DB: no data with the with key sessionToken found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetSessionUserName(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSessionUserName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.want.Error() != err.Error() {
				t.Errorf("Expected %v but got %v", tt.want, err)
			}
		})
	}
}
