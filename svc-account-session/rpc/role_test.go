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

// Package rpc defines the handler for micro services
package rpc

import (
	"context"
	e "errors"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
)

func mockRedfishRoles() error {
	list := asmodel.RedfishRoles{
		List: []string{
			"Administrator",
			"Operator",
			"ReadOnly",
		},
	}
	if err := list.Create(); err != nil {
		return err
	}
	return nil
}

func mockPrivilegeRegistry() error {
	list := asmodel.Privileges{
		List: []string{
			"Login",
			"ConfigureManager",
			"ConfigureUsers",
			"ConfigureSelf",
			"ConfigureComponents",
		},
	}
	if err := list.Create(); err != nil {
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

func TestRole_CreateRole1(t *testing.T) {
	type args struct {
		ctx context.Context
		req *roleproto.RoleRequest
	}
	common.SetUpMockConfig()
	tests := []struct {
		name                    string
		args                    args
		CheckSessionTimeOutFunc func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error)
		UpdateLastUsedTimeFunc  func(ctx context.Context, token string) error
		CreateFunc              func(ctx context.Context, req *roleproto.RoleRequest, session *asmodel.Session) response.RPC
		MarshalFunc             func(v any) ([]byte, error)
		want                    *roleproto.RoleResponse
		wantErr                 bool
	}{
		{
			name: "Session Timeout Error for 401(not valid session)",
			args: args{context.TODO(), &roleproto.RoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(errors.InvalidAuthToken, "error: invalid token ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			CreateFunc: func(ctx context.Context, req *roleproto.RoleRequest, session *asmodel.Session) response.RPC {
				return response.RPC{}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, nil },
			want:        &roleproto.RoleResponse{StatusCode: 401, StatusMessage: "Base.1.13.0.NoValidSession", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.NoValidSession\",\"Message\":\"There is no valid session established with the implementation.error while authorizing session token: error: invalid token \",\"Severity\":\"Critical\",\"Resolution\":\"Establish a session before attempting any operations.\"}]}}")},
			wantErr:     false,
		},
		{
			name: "Session Timeout Error for 504(Service unavailable)",
			args: args{context.TODO(), &roleproto.RoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(5, "error: Service unavailable ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			CreateFunc: func(ctx context.Context, req *roleproto.RoleRequest, session *asmodel.Session) response.RPC {
				return response.RPC{}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, nil },
			want:        &roleproto.RoleResponse{StatusCode: 503, StatusMessage: "Base.1.13.0.CouldNotEstablishConnection", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.CouldNotEstablishConnection\",\"Message\":\"The service failed to establish a connection with the URI 127.0.0.1:6379. error while authorizing session token: error: Service unavailable \",\"Severity\":\"Critical\",\"MessageArgs\":[\"127.0.0.1:6379\"],\"Resolution\":\"Ensure that the URI contains a valid and reachable node name, protocol information and other URI components.\"}]}}")},
			wantErr:     false,
		},
		{
			name: "UpdateLastUsedTime error",
			args: args{context.TODO(), &roleproto.RoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return e.New("fakeError") },
			CreateFunc: func(ctx context.Context, req *roleproto.RoleRequest, session *asmodel.Session) response.RPC {
				return response.RPC{}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, nil },
			want:        &roleproto.RoleResponse{StatusCode: 500, StatusMessage: "Base.1.13.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.error while updating last used time of session with token : fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:     false,
		},
		{
			name: "Marshall error",
			args: args{context.TODO(), &roleproto.RoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			CreateFunc: func(ctx context.Context, req *roleproto.RoleRequest, session *asmodel.Session) response.RPC {
				return response.RPC{StatusCode: 200, StatusMessage: "fakeMsg", Header: map[string]string{"fake": "fake"}}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, e.New("fakeError") },
			want:        &roleproto.RoleResponse{StatusCode: 500, Header: map[string]string{"fake": "fake"}, StatusMessage: "Base.1.13.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.error while trying to marshal the response body of create role API: fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:     false,
		},
		{
			name: "Pass case",
			args: args{context.TODO(), &roleproto.RoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			CreateFunc: func(ctx context.Context, req *roleproto.RoleRequest, session *asmodel.Session) response.RPC {
				return response.RPC{StatusCode: 200, StatusMessage: "fakeMsg", Header: map[string]string{"fake": "fake"}}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, nil },
			want:        &roleproto.RoleResponse{StatusCode: 200, Header: map[string]string{"fake": "fake"}, StatusMessage: "fakeMsg"},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		UpdateLastUsedTimeFunc = tt.UpdateLastUsedTimeFunc
		CreateFunc = tt.CreateFunc
		MarshalFunc = tt.MarshalFunc
		CheckSessionTimeOutFunc = tt.CheckSessionTimeOutFunc
		t.Run(tt.name, func(t *testing.T) {
			r := &Role{}
			got, err := r.CreateRole(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateRole() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRole_GetRole(t *testing.T) {
	type args struct {
		ctx context.Context
		req *roleproto.GetRoleRequest
	}
	common.SetUpMockConfig()
	tests := []struct {
		name                    string
		args                    args
		CheckSessionTimeOutFunc func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error)
		UpdateLastUsedTimeFunc  func(ctx context.Context, token string) error
		GetRoleFunc             func(ctx context.Context, req *roleproto.GetRoleRequest, session *asmodel.Session) response.RPC
		MarshalFunc             func(v any) ([]byte, error)
		want                    *roleproto.RoleResponse
		wantErr                 bool
	}{
		{
			name: "Session Timeout Error for 401(not valid session)",
			args: args{context.TODO(), &roleproto.GetRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(errors.InvalidAuthToken, "error: invalid token ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			GetRoleFunc: func(ctx context.Context, req *roleproto.GetRoleRequest, session *asmodel.Session) response.RPC {
				return response.RPC{}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, nil },
			want:        &roleproto.RoleResponse{StatusCode: 401, StatusMessage: "Base.1.13.0.NoValidSession", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.NoValidSession\",\"Message\":\"There is no valid session established with the implementation.error while authorizing session token: error: invalid token \",\"Severity\":\"Critical\",\"Resolution\":\"Establish a session before attempting any operations.\"}]}}")},
			wantErr:     false,
		},
		{
			name: "Session Timeout Error for 504(Service unavailable)",
			args: args{context.TODO(), &roleproto.GetRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(5, "error: Service unavailable ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			GetRoleFunc: func(ctx context.Context, req *roleproto.GetRoleRequest, session *asmodel.Session) response.RPC {
				return response.RPC{}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, nil },
			want:        &roleproto.RoleResponse{StatusCode: 503, StatusMessage: "Base.1.13.0.CouldNotEstablishConnection", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.CouldNotEstablishConnection\",\"Message\":\"The service failed to establish a connection with the URI 127.0.0.1:6379. error while authorizing session token: error: Service unavailable \",\"Severity\":\"Critical\",\"MessageArgs\":[\"127.0.0.1:6379\"],\"Resolution\":\"Ensure that the URI contains a valid and reachable node name, protocol information and other URI components.\"}]}}")},
			wantErr:     false,
		},
		{
			name: "UpdateLastUsedTime error",
			args: args{context.TODO(), &roleproto.GetRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return e.New("fakeError") },
			GetRoleFunc: func(ctx context.Context, req *roleproto.GetRoleRequest, session *asmodel.Session) response.RPC {
				return response.RPC{}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, nil },
			want:        &roleproto.RoleResponse{StatusCode: 500, StatusMessage: "Base.1.13.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.error while updating last used time of session with token : fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:     false,
		},
		{
			name: "Marshall error",
			args: args{context.TODO(), &roleproto.GetRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			GetRoleFunc: func(ctx context.Context, req *roleproto.GetRoleRequest, session *asmodel.Session) response.RPC {
				return response.RPC{StatusCode: 200, StatusMessage: "fakeMsg", Header: map[string]string{"fake": "fake"}}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, e.New("fakeError") },
			want:        &roleproto.RoleResponse{StatusCode: 500, Header: map[string]string{"fake": "fake"}, StatusMessage: "Base.1.13.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.error while trying to marshal the response body of get role API: fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:     false,
		},
		{
			name: "Pass case",
			args: args{context.TODO(), &roleproto.GetRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			GetRoleFunc: func(ctx context.Context, req *roleproto.GetRoleRequest, session *asmodel.Session) response.RPC {
				return response.RPC{StatusCode: 200, StatusMessage: "fakeMsg", Header: map[string]string{"fake": "fake"}}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, nil },
			want:        &roleproto.RoleResponse{StatusCode: 200, Header: map[string]string{"fake": "fake"}, StatusMessage: "fakeMsg"},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		CheckSessionTimeOutFunc = tt.CheckSessionTimeOutFunc
		UpdateLastUsedTimeFunc = tt.UpdateLastUsedTimeFunc
		GetRoleFunc = tt.GetRoleFunc
		MarshalFunc = tt.MarshalFunc
		t.Run(tt.name, func(t *testing.T) {
			r := &Role{}

			got, err := r.GetRole(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRole() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRole_GetAllRoles(t *testing.T) {
	type args struct {
		ctx context.Context
		req *roleproto.GetRoleRequest
	}
	common.SetUpMockConfig()
	tests := []struct {
		name                    string
		args                    args
		CheckSessionTimeOutFunc func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error)
		UpdateLastUsedTimeFunc  func(ctx context.Context, token string) error
		GetAllRolesFunc         func(ctx context.Context, session *asmodel.Session) response.RPC
		MarshalFunc             func(v any) ([]byte, error)
		want                    *roleproto.RoleResponse
		wantErr                 bool
	}{
		{
			name: "Session Timeout Error for 401(not valid session)",
			args: args{context.TODO(), &roleproto.GetRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(errors.InvalidAuthToken, "error: invalid token ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			GetAllRolesFunc:        func(ctx context.Context, session *asmodel.Session) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &roleproto.RoleResponse{StatusCode: 401, StatusMessage: "Base.1.13.0.NoValidSession", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.NoValidSession\",\"Message\":\"There is no valid session established with the implementation.error while authorizing session token: error: invalid token \",\"Severity\":\"Critical\",\"Resolution\":\"Establish a session before attempting any operations.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "Session Timeout Error for 504(Service unavailable)",
			args: args{context.TODO(), &roleproto.GetRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(5, "error: Service unavailable ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			GetAllRolesFunc:        func(ctx context.Context, session *asmodel.Session) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &roleproto.RoleResponse{StatusCode: 503, StatusMessage: "Base.1.13.0.CouldNotEstablishConnection", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.CouldNotEstablishConnection\",\"Message\":\"The service failed to establish a connection with the URI 127.0.0.1:6379. error while authorizing session token: error: Service unavailable \",\"Severity\":\"Critical\",\"MessageArgs\":[\"127.0.0.1:6379\"],\"Resolution\":\"Ensure that the URI contains a valid and reachable node name, protocol information and other URI components.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "UpdateLastUsedTime error",
			args: args{context.TODO(), &roleproto.GetRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return e.New("fakeError") },
			GetAllRolesFunc:        func(ctx context.Context, session *asmodel.Session) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &roleproto.RoleResponse{StatusCode: 500, StatusMessage: "Base.1.13.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.error while updating last used time of session with token : fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "Marshall error",
			args: args{context.TODO(), &roleproto.GetRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			GetAllRolesFunc: func(ctx context.Context, session *asmodel.Session) response.RPC {
				return response.RPC{StatusCode: 200, StatusMessage: "fakeMsg", Header: map[string]string{"fake": "fake"}}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, e.New("fakeError") },
			want:        &roleproto.RoleResponse{StatusCode: 500, Header: map[string]string{"fake": "fake"}, StatusMessage: "Base.1.13.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.error while trying to marshal the response body of the get all roles API: fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:     false,
		},
		{
			name: "Pass case",
			args: args{context.TODO(), &roleproto.GetRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			GetAllRolesFunc: func(ctx context.Context, session *asmodel.Session) response.RPC {
				return response.RPC{StatusCode: 200, StatusMessage: "fakeMsg", Header: map[string]string{"fake": "fake"}}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, nil },
			want:        &roleproto.RoleResponse{StatusCode: 200, Header: map[string]string{"fake": "fake"}, StatusMessage: "fakeMsg"},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		CheckSessionTimeOutFunc = tt.CheckSessionTimeOutFunc
		UpdateLastUsedTimeFunc = tt.UpdateLastUsedTimeFunc
		GetAllRolesFunc = tt.GetAllRolesFunc
		MarshalFunc = tt.MarshalFunc
		t.Run(tt.name, func(t *testing.T) {
			r := &Role{}
			got, err := r.GetAllRoles(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllRoles() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRole_UpdateRole(t *testing.T) {
	type args struct {
		ctx context.Context
		req *roleproto.UpdateRoleRequest
	}
	common.SetUpMockConfig()
	tests := []struct {
		name                    string
		args                    args
		CheckSessionTimeOutFunc func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error)
		UpdateLastUsedTimeFunc  func(ctx context.Context, token string) error
		UpdateFunc              func(ctx context.Context, req *roleproto.UpdateRoleRequest, session *asmodel.Session) response.RPC
		MarshalFunc             func(v any) ([]byte, error)
		want                    *roleproto.RoleResponse
		wantErr                 bool
	}{
		{
			name: "Session Timeout Error for 401(not valid session)",
			args: args{context.TODO(), &roleproto.UpdateRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(errors.InvalidAuthToken, "error: invalid token ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			UpdateFunc: func(ctx context.Context, req *roleproto.UpdateRoleRequest, session *asmodel.Session) response.RPC {
				return response.RPC{}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, nil },
			want:        &roleproto.RoleResponse{StatusCode: 401, StatusMessage: "Base.1.13.0.NoValidSession", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.NoValidSession\",\"Message\":\"There is no valid session established with the implementation.error while authorizing session token: error: invalid token \",\"Severity\":\"Critical\",\"Resolution\":\"Establish a session before attempting any operations.\"}]}}")},
			wantErr:     false,
		},
		{
			name: "Session Timeout Error for 504(Service unavailable)",
			args: args{context.TODO(), &roleproto.UpdateRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(5, "error: Service unavailable ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			UpdateFunc: func(ctx context.Context, req *roleproto.UpdateRoleRequest, session *asmodel.Session) response.RPC {
				return response.RPC{}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, nil },
			want:        &roleproto.RoleResponse{StatusCode: 503, StatusMessage: "Base.1.13.0.CouldNotEstablishConnection", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.CouldNotEstablishConnection\",\"Message\":\"The service failed to establish a connection with the URI 127.0.0.1:6379. error while authorizing session token: error: Service unavailable \",\"Severity\":\"Critical\",\"MessageArgs\":[\"127.0.0.1:6379\"],\"Resolution\":\"Ensure that the URI contains a valid and reachable node name, protocol information and other URI components.\"}]}}")},
			wantErr:     false,
		},
		{
			name: "UpdateLastUsedTime error",
			args: args{context.TODO(), &roleproto.UpdateRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return e.New("fakeError") },
			UpdateFunc: func(ctx context.Context, req *roleproto.UpdateRoleRequest, session *asmodel.Session) response.RPC {
				return response.RPC{}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, nil },
			want:        &roleproto.RoleResponse{StatusCode: 500, StatusMessage: "Base.1.13.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.error while updating last used time of session with token : fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:     false,
		},
		{
			name: "Marshall error",
			args: args{context.TODO(), &roleproto.UpdateRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			UpdateFunc: func(ctx context.Context, req *roleproto.UpdateRoleRequest, session *asmodel.Session) response.RPC {
				return response.RPC{StatusCode: 200, StatusMessage: "fakeMsg", Header: map[string]string{"fake": "fake"}}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, e.New("fakeError") },
			want:        &roleproto.RoleResponse{StatusCode: 500, Header: map[string]string{"fake": "fake"}, StatusMessage: "Base.1.13.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.error while trying to marshal the response body of the update role API: fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:     false,
		},
		{
			name: "Pass case",
			args: args{context.TODO(), &roleproto.UpdateRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			UpdateFunc: func(ctx context.Context, req *roleproto.UpdateRoleRequest, session *asmodel.Session) response.RPC {
				return response.RPC{StatusCode: 200, StatusMessage: "fakeMsg", Header: map[string]string{"fake": "fake"}}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, nil },
			want:        &roleproto.RoleResponse{StatusCode: 200, Header: map[string]string{"fake": "fake"}, StatusMessage: "fakeMsg"},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		CheckSessionTimeOutFunc = tt.CheckSessionTimeOutFunc
		UpdateLastUsedTimeFunc = tt.UpdateLastUsedTimeFunc
		UpdateFunc = tt.UpdateFunc
		MarshalFunc = tt.MarshalFunc
		t.Run(tt.name, func(t *testing.T) {
			r := &Role{}
			got, err := r.UpdateRole(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateRole() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRole_DeleteRole(t *testing.T) {
	type args struct {
		ctx context.Context
		req *roleproto.DeleteRoleRequest
	}
	common.SetUpMockConfig()
	tests := []struct {
		name                    string
		args                    args
		CheckSessionTimeOutFunc func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error)
		UpdateLastUsedTimeFunc  func(ctx context.Context, token string) error
		DeleteFunc              func(ctx context.Context, req *roleproto.DeleteRoleRequest) *response.RPC
		MarshalFunc             func(v any) ([]byte, error)
		want                    *roleproto.RoleResponse
		wantErr                 bool
	}{
		{
			name: "Marshall error",
			args: args{context.TODO(), &roleproto.DeleteRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			DeleteFunc: func(ctx context.Context, req *roleproto.DeleteRoleRequest) *response.RPC {
				return &response.RPC{StatusCode: 200, StatusMessage: "fakeMsg", Header: map[string]string{"fake": "fake"}}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, e.New("fakeError") },
			want:        &roleproto.RoleResponse{StatusCode: 500, Header: map[string]string{"fake": "fake"}, StatusMessage: "Base.1.13.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.error while trying to marshal the response body of the delete role API: fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:     false,
		},
		{
			name: "Pass case",
			args: args{context.TODO(), &roleproto.DeleteRoleRequest{}},
			CheckSessionTimeOutFunc: func(ctx context.Context, sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(ctx context.Context, token string) error { return nil },
			DeleteFunc: func(ctx context.Context, req *roleproto.DeleteRoleRequest) *response.RPC {
				return &response.RPC{StatusCode: 200, StatusMessage: "fakeMsg", Header: map[string]string{"fake": "fake"}}
			},
			MarshalFunc: func(v any) ([]byte, error) { return nil, nil },
			want:        &roleproto.RoleResponse{StatusCode: 200, Header: map[string]string{"fake": "fake"}, StatusMessage: "fakeMsg"},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		CheckSessionTimeOutFunc = tt.CheckSessionTimeOutFunc
		UpdateLastUsedTimeFunc = tt.UpdateLastUsedTimeFunc
		DeleteFunc = tt.DeleteFunc
		MarshalFunc = tt.MarshalFunc
		t.Run(tt.name, func(t *testing.T) {
			r := &Role{}
			got, err := r.DeleteRole(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteRole() got = %v, want %v", got, tt.want)
			}
		})
	}
}
