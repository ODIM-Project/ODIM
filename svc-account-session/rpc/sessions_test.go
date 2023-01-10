package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
)

type any = interface{}

func TestSession_CreateSession(t *testing.T) {
	ctx := mockContext()
	type args struct {
		req *sessionproto.SessionCreateRequest
	}
	tests := []struct {
		name                 string
		args                 args
		CreateNewSessionFunc func(ctx context.Context, req *sessionproto.SessionCreateRequest) (response.RPC, string)
		MarshalFunc          func(v any) ([]byte, error)
		want                 *sessionproto.SessionCreateResponse
		wantErr              bool
	}{
		{
			name: "Marshall error",
			args: args{&sessionproto.SessionCreateRequest{}},
			CreateNewSessionFunc: func(ctx context.Context, req *sessionproto.SessionCreateRequest) (response.RPC, string) {
				return common.GeneralError(400, "fakeStatus", "fakeError", nil, &common.TaskUpdateInfo{TaskID: "1"}), ""
			},
			MarshalFunc: func(v any) ([]byte, error) { return []byte{}, errors.New("fakeError") },
			want:        &sessionproto.SessionCreateResponse{StatusCode: 500, StatusMessage: "error while trying to marshal the response body of the create session API: fakeError"},
			wantErr:     false,
		},
		{
			name: "No error",
			args: args{},
			CreateNewSessionFunc: func(ctx context.Context, req *sessionproto.SessionCreateRequest) (response.RPC, string) {
				return response.RPC{StatusCode: 200, StatusMessage: "Success", Header: map[string]string{"pass": "case"}}, "413"
			},
			MarshalFunc: func(v any) ([]byte, error) { return json.Marshal(v) },
			want:        &sessionproto.SessionCreateResponse{SessionId: "413", StatusCode: 200, StatusMessage: "Success", Body: []byte("null"), Header: map[string]string{"pass": "case"}},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		CreateNewSessionFunc = tt.CreateNewSessionFunc
		MarshalFunc = tt.MarshalFunc
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			got, err := s.CreateSession(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSession() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_DeleteSession(t *testing.T) {
	ctx := mockContext()
	type args struct {
		req *sessionproto.SessionRequest
	}
	tests := []struct {
		name              string
		args              args
		DeleteSessionFunc func(ctx context.Context, req *sessionproto.SessionRequest) response.RPC
		MarshalFunc       func(v any) ([]byte, error)
		want              *sessionproto.SessionResponse
		wantErr           bool
	}{
		{
			name: "Marshall error",
			args: args{&sessionproto.SessionRequest{}},
			DeleteSessionFunc: func(ctx context.Context, req *sessionproto.SessionRequest) response.RPC {
				return common.GeneralError(400, "fakeStatus", "fakeError", nil, &common.TaskUpdateInfo{TaskID: "1"})
			},
			MarshalFunc: func(v any) ([]byte, error) { return []byte{}, errors.New("fakeError") },
			want:        &sessionproto.SessionResponse{StatusCode: 500, StatusMessage: "error while trying to marshal the response body of the delete session API: fakeError"},
			wantErr:     false,
		},
		{
			name: "No error",
			args: args{},
			DeleteSessionFunc: func(ctx context.Context, req *sessionproto.SessionRequest) response.RPC {
				return response.RPC{StatusCode: 200, StatusMessage: "Success", Header: map[string]string{"pass": "case"}}
			},
			MarshalFunc: func(v any) ([]byte, error) { return json.Marshal(v) },
			want:        &sessionproto.SessionResponse{StatusCode: 200, StatusMessage: "Success", Body: []byte("null"), Header: map[string]string{"pass": "case"}},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		DeleteSessionFunc = tt.DeleteSessionFunc
		MarshalFunc = tt.MarshalFunc
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			got, err := s.DeleteSession(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteSession() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_GetAllActiveSessions(t *testing.T) {
	ctx := mockContext()
	type args struct {
		req *sessionproto.SessionRequest
	}
	tests := []struct {
		name                     string
		args                     args
		GetAllActiveSessionsFunc func(ctx context.Context, req *sessionproto.SessionRequest) response.RPC
		MarshalFunc              func(v any) ([]byte, error)
		want                     *sessionproto.SessionResponse
		wantErr                  bool
	}{
		{
			name: "Marshall error",
			args: args{&sessionproto.SessionRequest{}},
			GetAllActiveSessionsFunc: func(ctx context.Context, req *sessionproto.SessionRequest) response.RPC {
				return common.GeneralError(400, "fakeStatus", "fakeError", nil, &common.TaskUpdateInfo{TaskID: "1"})
			},
			MarshalFunc: func(v any) ([]byte, error) { return []byte{}, errors.New("fakeError") },
			want:        &sessionproto.SessionResponse{StatusCode: 500, StatusMessage: "error while trying to marshal the response body of the get all active session API: fakeError"},
			wantErr:     false,
		},
		{
			name: "No error",
			args: args{},
			GetAllActiveSessionsFunc: func(ctx context.Context, req *sessionproto.SessionRequest) response.RPC {
				return response.RPC{StatusCode: 200, StatusMessage: "Success", Header: map[string]string{"pass": "case"}}
			},
			MarshalFunc: func(v any) ([]byte, error) { return json.Marshal(v) },
			want:        &sessionproto.SessionResponse{StatusCode: 200, StatusMessage: "Success", Body: []byte("null"), Header: map[string]string{"pass": "case"}},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		GetAllActiveSessionsFunc = tt.GetAllActiveSessionsFunc
		MarshalFunc = tt.MarshalFunc
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			got, err := s.GetAllActiveSessions(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllActiveSessions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllActiveSessions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_GetSession(t *testing.T) {
	ctx := mockContext()
	type args struct {
		req *sessionproto.SessionRequest
	}
	tests := []struct {
		name           string
		args           args
		GetSessionFunc func(ctx context.Context, req *sessionproto.SessionRequest) response.RPC
		MarshalFunc    func(v any) ([]byte, error)
		want           *sessionproto.SessionResponse
		wantErr        bool
	}{
		{
			name: "Marshall error",
			args: args{&sessionproto.SessionRequest{}},
			GetSessionFunc: func(ctx context.Context, req *sessionproto.SessionRequest) response.RPC {
				return common.GeneralError(400, "fakeStatus", "fakeError", nil, &common.TaskUpdateInfo{TaskID: "1"})
			},
			MarshalFunc: func(v any) ([]byte, error) { return []byte{}, errors.New("fakeError") },
			want:        &sessionproto.SessionResponse{StatusMessage: "error while trying to marshal the response body of the get session API: fakeError"},
			wantErr:     false,
		},
		{
			name: "No error",
			args: args{},
			GetSessionFunc: func(ctx context.Context, req *sessionproto.SessionRequest) response.RPC {
				return response.RPC{StatusCode: 200, StatusMessage: "Success", Header: map[string]string{"pass": "case"}}
			},
			MarshalFunc: func(v any) ([]byte, error) { return json.Marshal(v) },
			want:        &sessionproto.SessionResponse{StatusCode: 200, StatusMessage: "Success", Body: []byte("null"), Header: map[string]string{"pass": "case"}},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		GetSessionFunc = tt.GetSessionFunc
		MarshalFunc = tt.MarshalFunc
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			got, err := s.GetSession(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSession() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_GetSessionService(t *testing.T) {
	ctx := mockContext()
	type args struct {
		req *sessionproto.SessionRequest
	}
	tests := []struct {
		name                  string
		args                  args
		GetSessionServiceFunc func(ctx context.Context, req *sessionproto.SessionRequest) response.RPC
		MarshalFunc           func(v any) ([]byte, error)
		want                  *sessionproto.SessionResponse
		wantErr               bool
	}{
		{
			name: "Marshall error",
			args: args{&sessionproto.SessionRequest{}},
			GetSessionServiceFunc: func(ctx context.Context, req *sessionproto.SessionRequest) response.RPC {
				return common.GeneralError(400, "fakeStatus", "fakeError", nil, &common.TaskUpdateInfo{TaskID: "1"})
			},
			MarshalFunc: func(v any) ([]byte, error) { return []byte{}, errors.New("fakeError") },
			want:        &sessionproto.SessionResponse{StatusCode: 500, StatusMessage: "error while trying to marshal the response body of the get session service API: fakeError"},
			wantErr:     false,
		},
		{
			name: "No error",
			args: args{},
			GetSessionServiceFunc: func(ctx context.Context, req *sessionproto.SessionRequest) response.RPC {
				return response.RPC{StatusCode: 200, StatusMessage: "Success", Header: map[string]string{"pass": "case"}}
			},
			MarshalFunc: func(v any) ([]byte, error) { return json.Marshal(v) },
			want:        &sessionproto.SessionResponse{StatusCode: 200, StatusMessage: "Success", Body: []byte("null"), Header: map[string]string{"pass": "case"}},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		GetSessionServiceFunc = tt.GetSessionServiceFunc
		MarshalFunc = tt.MarshalFunc
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			got, err := s.GetSessionService(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSessionService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSessionService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_GetSessionUserName(t *testing.T) {
	ctx := mockContext()
	type args struct {
		req *sessionproto.SessionRequest
	}
	tests := []struct {
		name                   string
		args                   args
		GetSessionUserNameFunc func(ctx context.Context, req *sessionproto.SessionRequest) (*sessionproto.SessionUserName, error)
		want                   *sessionproto.SessionUserName
		wantErr                bool
	}{
		{
			name: "Pass case",
			args: args{},
			GetSessionUserNameFunc: func(ctx context.Context, req *sessionproto.SessionRequest) (*sessionproto.SessionUserName, error) {
				return &sessionproto.SessionUserName{}, nil
			},
			want:    &sessionproto.SessionUserName{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		GetSessionUserNameFunc = tt.GetSessionUserNameFunc
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			got, err := s.GetSessionUserName(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSessionUserName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSessionUserName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_GetSessionUserRoleID(t *testing.T) {
	ctx := mockContext()
	type args struct {
		req *sessionproto.SessionRequest
	}
	tests := []struct {
		name                     string
		args                     args
		GetSessionUserRoleIDFunc func(ctx context.Context, req *sessionproto.SessionRequest) (*sessionproto.SessionUsersRoleID, error)
		want                     *sessionproto.SessionUsersRoleID
		wantErr                  bool
	}{
		{
			name: "Pass case",
			args: args{},
			GetSessionUserRoleIDFunc: func(ctx context.Context, req *sessionproto.SessionRequest) (*sessionproto.SessionUsersRoleID, error) {
				return &sessionproto.SessionUsersRoleID{}, nil
			},
			want:    &sessionproto.SessionUsersRoleID{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		GetSessionUserRoleIDFunc = tt.GetSessionUserRoleIDFunc
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			got, err := s.GetSessionUserRoleID(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSessionUserRoleID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSessionUserRoleID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCommonResponse(t *testing.T) {
	type args struct {
		statusMessage string
	}
	tests := []struct {
		name string
		args args
		want asresponse.RedfishSessionResponse
	}{
		{
			name: "Pass case",
			args: args{"Success"},
			want: asresponse.RedfishSessionResponse{Error: asresponse.Error{Code: "Base.1.13.0.GeneralError", Message: "See @Message.ExtendedInfo for more information.", ExtendedInfos: []asresponse.ExtendedInfo{asresponse.ExtendedInfo{"Success"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCommonResponse(tt.args.statusMessage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCommonResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
