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
package rpc

import (
	"context"
	e "errors"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	accountproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/account"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
)

func TestAccount_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		req *accountproto.CreateAccountRequest
	}
	common.SetUpMockConfig()
	tests := []struct {
		name                    string
		args                    args
		CheckSessionTimeOutFunc func(sessionToken string) (*asmodel.Session, *errors.Error)
		UpdateLastUsedTimeFunc  func(token string) error
		MarshalFunc             func(v any) ([]byte, error)
		want                    *accountproto.AccountResponse
		wantErr                 bool
	}{
		{
			name: "Session Timeout Error for 401(not valid session)",
			args: args{context.TODO(), &accountproto.CreateAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(errors.InvalidAuthToken, "error: invalid token ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 401, StatusMessage: "Base.1.11.0.NoValidSession", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.NoValidSession\",\"Message\":\"There is no valid session established with the implementation.error while authorizing session token: error: invalid token \",\"Severity\":\"Critical\",\"Resolution\":\"Establish a session before attempting any operations.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "Session Timeout Error for 504(Service unavailable)",
			args: args{context.TODO(), &accountproto.CreateAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(5, "error: Service unavailable ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 503, StatusMessage: "Base.1.11.0.CouldNotEstablishConnection", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.CouldNotEstablishConnection\",\"Message\":\"The service failed to establish a connection with the URI 127.0.0.1:6379. error while authorizing session token: error: Service unavailable \",\"Severity\":\"Critical\",\"MessageArgs\":[\"127.0.0.1:6379\"],\"Resolution\":\"Ensure that the URI contains a valid and reachable node name, protocol information and other URI components.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "UpdateLastUsedTime error",
			args: args{context.TODO(), &accountproto.CreateAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return e.New("fakeError") },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 500, StatusMessage: "Base.1.11.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.error while updating last used time of session with token : fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "Marshall error",
			args: args{context.TODO(), &accountproto.CreateAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, e.New("fakeError") },
			want:                   &accountproto.AccountResponse{StatusCode: 500, StatusMessage: "error while trying marshal the response body for create account: fakeError"},
			wantErr:                false,
		},
		{
			name: "Pass case",
			args: args{context.TODO(), &accountproto.CreateAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 500, StatusMessage: "Base.1.11.0.InternalError"},
			wantErr:                false,
		},
	}
	for _, tt := range tests {
		CheckSessionTimeOutFunc = tt.CheckSessionTimeOutFunc
		UpdateLastUsedTimeFunc = tt.UpdateLastUsedTimeFunc
		MarshalFunc = tt.MarshalFunc
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{}
			got, err := a.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_GetAllAccounts(t *testing.T) {
	type args struct {
		ctx context.Context
		req *accountproto.AccountRequest
	}
	common.SetUpMockConfig()
	tests := []struct {
		name                    string
		args                    args
		CheckSessionTimeOutFunc func(sessionToken string) (*asmodel.Session, *errors.Error)
		UpdateLastUsedTimeFunc  func(token string) error
		GetAllAccountsFunc      func(session *asmodel.Session) response.RPC
		MarshalFunc             func(v any) ([]byte, error)
		want                    *accountproto.AccountResponse
		wantErr                 bool
	}{
		{
			name: "Session Timeout Error for 401(not valid session)",
			args: args{context.TODO(), &accountproto.AccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(errors.InvalidAuthToken, "error: invalid token ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			GetAllAccountsFunc:     func(session *asmodel.Session) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 401, StatusMessage: "Base.1.11.0.NoValidSession", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.NoValidSession\",\"Message\":\"There is no valid session established with the implementation.error while authorizing session token: error: invalid token \",\"Severity\":\"Critical\",\"Resolution\":\"Establish a session before attempting any operations.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "Session Timeout Error for 504(Service unavailable)",
			args: args{context.TODO(), &accountproto.AccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(5, "error: Service unavailable ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			GetAllAccountsFunc:     func(session *asmodel.Session) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 503, StatusMessage: "Base.1.11.0.CouldNotEstablishConnection", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.CouldNotEstablishConnection\",\"Message\":\"The service failed to establish a connection with the URI 127.0.0.1:6379. error while authorizing session token: error: Service unavailable \",\"Severity\":\"Critical\",\"MessageArgs\":[\"127.0.0.1:6379\"],\"Resolution\":\"Ensure that the URI contains a valid and reachable node name, protocol information and other URI components.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "UpdateLastUsedTime error",
			args: args{context.TODO(), &accountproto.AccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return e.New("fakeError") },
			GetAllAccountsFunc:     func(session *asmodel.Session) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 500, StatusMessage: "Base.1.11.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.error while updating last used time of session with token : fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "Marshall error",
			args: args{context.TODO(), &accountproto.AccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			GetAllAccountsFunc:     func(session *asmodel.Session) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, e.New("fakeError") },
			want:                   &accountproto.AccountResponse{StatusCode: 500, StatusMessage: "error while trying marshal the response body for get all accounts: fakeError"},
			wantErr:                true,
		},
		{
			name: "Pass case",
			args: args{context.TODO(), &accountproto.AccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			GetAllAccountsFunc:     func(session *asmodel.Session) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{},
			wantErr:                false,
		},
	}
	for _, tt := range tests {
		CheckSessionTimeOutFunc = tt.CheckSessionTimeOutFunc
		UpdateLastUsedTimeFunc = tt.UpdateLastUsedTimeFunc
		GetAllAccountsFunc = tt.GetAllAccountsFunc
		MarshalFunc = tt.MarshalFunc
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{}
			got, err := a.GetAllAccounts(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllAccounts() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_GetAccount(t *testing.T) {
	type args struct {
		ctx context.Context
		req *accountproto.GetAccountRequest
	}
	common.SetUpMockConfig()
	tests := []struct {
		name                    string
		args                    args
		CheckSessionTimeOutFunc func(sessionToken string) (*asmodel.Session, *errors.Error)
		UpdateLastUsedTimeFunc  func(token string) error
		GetAccountFunc          func(session *asmodel.Session, accountID string) response.RPC
		MarshalFunc             func(v any) ([]byte, error)
		want                    *accountproto.AccountResponse
		wantErr                 bool
	}{
		{
			name: "Session Timeout Error for 401(not valid session)",
			args: args{context.TODO(), &accountproto.GetAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(errors.InvalidAuthToken, "error: invalid token ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			GetAccountFunc:         func(session *asmodel.Session, accountID string) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 401, StatusMessage: "Base.1.11.0.NoValidSession", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.NoValidSession\",\"Message\":\"There is no valid session established with the implementation.error while authorizing session token: error: invalid token \",\"Severity\":\"Critical\",\"Resolution\":\"Establish a session before attempting any operations.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "Session Timeout Error for 504(Service unavailable)",
			args: args{context.TODO(), &accountproto.GetAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(5, "error: Service unavailable ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			GetAccountFunc:         func(session *asmodel.Session, accountID string) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 503, StatusMessage: "Base.1.11.0.CouldNotEstablishConnection", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.CouldNotEstablishConnection\",\"Message\":\"The service failed to establish a connection with the URI 127.0.0.1:6379. error while authorizing session token: error: Service unavailable \",\"Severity\":\"Critical\",\"MessageArgs\":[\"127.0.0.1:6379\"],\"Resolution\":\"Ensure that the URI contains a valid and reachable node name, protocol information and other URI components.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "UpdateLastUsedTime error",
			args: args{context.TODO(), &accountproto.GetAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return e.New("fakeError") },
			GetAccountFunc:         func(session *asmodel.Session, accountID string) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 500, StatusMessage: "Base.1.11.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.error while updating last used time of session with token : fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "Marshall error",
			args: args{context.TODO(), &accountproto.GetAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			GetAccountFunc:         func(session *asmodel.Session, accountID string) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, e.New("fakeError") },
			want:                   &accountproto.AccountResponse{StatusCode: 500, StatusMessage: "error while trying marshal the response body for get account details: fakeError"},
			wantErr:                true,
		},
		{
			name: "Pass case",
			args: args{context.TODO(), &accountproto.GetAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			GetAccountFunc:         func(session *asmodel.Session, accountID string) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{},
			wantErr:                false,
		},
	}
	for _, tt := range tests {
		CheckSessionTimeOutFunc = tt.CheckSessionTimeOutFunc
		UpdateLastUsedTimeFunc = tt.UpdateLastUsedTimeFunc
		GetAccountFunc = tt.GetAccountFunc
		MarshalFunc = tt.MarshalFunc
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{}
			got, err := a.GetAccount(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAccount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_GetAccountServices(t *testing.T) {
	type args struct {
		ctx context.Context
		req *accountproto.AccountRequest
	}
	common.SetUpMockConfig()
	tests := []struct {
		name                    string
		args                    args
		CheckSessionTimeOutFunc func(sessionToken string) (*asmodel.Session, *errors.Error)
		UpdateLastUsedTimeFunc  func(token string) error
		GetAccountServiceFunc   func() response.RPC
		MarshalFunc             func(v any) ([]byte, error)
		want                    *accountproto.AccountResponse
		wantErr                 bool
	}{
		{
			name: "Session Timeout Error for 401(not valid session)",
			args: args{context.TODO(), &accountproto.AccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(errors.InvalidAuthToken, "error: invalid token ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			GetAccountServiceFunc:  func() response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 401, StatusMessage: "Base.1.11.0.NoValidSession", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.NoValidSession\",\"Message\":\"There is no valid session established with the implementation.error while authorizing session token: error: invalid token \",\"Severity\":\"Critical\",\"Resolution\":\"Establish a session before attempting any operations.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "Session Timeout Error for 504(Service unavailable)",
			args: args{context.TODO(), &accountproto.AccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(5, "error: Service unavailable ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			GetAccountServiceFunc:  func() response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 503, StatusMessage: "Base.1.11.0.CouldNotEstablishConnection", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.CouldNotEstablishConnection\",\"Message\":\"The service failed to establish a connection with the URI 127.0.0.1:6379. error while authorizing session token: error: Service unavailable \",\"Severity\":\"Critical\",\"MessageArgs\":[\"127.0.0.1:6379\"],\"Resolution\":\"Ensure that the URI contains a valid and reachable node name, protocol information and other URI components.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "UpdateLastUsedTime error",
			args: args{context.TODO(), &accountproto.AccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return e.New("fakeError") },
			GetAccountServiceFunc:  func() response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 500, StatusMessage: "Base.1.11.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.error while updating last used time of session with token : fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "Marshall error",
			args: args{context.TODO(), &accountproto.AccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			GetAccountServiceFunc:  func() response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, e.New("fakeError") },
			want:                   &accountproto.AccountResponse{StatusCode: 500, StatusMessage: "error while trying marshal the response body for get account details: fakeError"},
			wantErr:                true,
		},
		{
			name: "Pass case",
			args: args{context.TODO(), &accountproto.AccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			GetAccountServiceFunc:  func() response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{},
			wantErr:                false,
		},
	}
	for _, tt := range tests {
		CheckSessionTimeOutFunc = tt.CheckSessionTimeOutFunc
		UpdateLastUsedTimeFunc = tt.UpdateLastUsedTimeFunc
		GetAccountServiceFunc = tt.GetAccountServiceFunc
		MarshalFunc = tt.MarshalFunc
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{}
			got, err := a.GetAccountServices(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAccountServices() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		req *accountproto.UpdateAccountRequest
	}
	common.SetUpMockConfig()
	tests := []struct {
		name                    string
		args                    args
		CheckSessionTimeOutFunc func(sessionToken string) (*asmodel.Session, *errors.Error)
		UpdateLastUsedTimeFunc  func(token string) error
		MarshalFunc             func(v any) ([]byte, error)
		want                    *accountproto.AccountResponse
		wantErr                 bool
	}{
		{
			name: "Session Timeout Error for 401(not valid session)",
			args: args{context.TODO(), &accountproto.UpdateAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(errors.InvalidAuthToken, "error: invalid token ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 401, StatusMessage: "Base.1.11.0.NoValidSession", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.NoValidSession\",\"Message\":\"There is no valid session established with the implementation.error while authorizing session token: error: invalid token \",\"Severity\":\"Critical\",\"Resolution\":\"Establish a session before attempting any operations.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "Session Timeout Error for 504(Service unavailable)",
			args: args{context.TODO(), &accountproto.UpdateAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(5, "error: Service unavailable ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 503, StatusMessage: "Base.1.11.0.CouldNotEstablishConnection", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.CouldNotEstablishConnection\",\"Message\":\"The service failed to establish a connection with the URI 127.0.0.1:6379. error while authorizing session token: error: Service unavailable \",\"Severity\":\"Critical\",\"MessageArgs\":[\"127.0.0.1:6379\"],\"Resolution\":\"Ensure that the URI contains a valid and reachable node name, protocol information and other URI components.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "UpdateLastUsedTime error",
			args: args{context.TODO(), &accountproto.UpdateAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return e.New("fakeError") },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 500, StatusMessage: "Base.1.11.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.error while updating last used time of session with token : fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "Marshall error",
			args: args{context.TODO(), &accountproto.UpdateAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, e.New("fakeError") },
			want:                   &accountproto.AccountResponse{StatusCode: 500, StatusMessage: "error while trying to marshal the response body for create account: fakeError"},
			wantErr:                false,
		},
		{
			name: "Pass case",
			args: args{context.TODO(), &accountproto.UpdateAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 500, StatusMessage: "Base.1.11.0.InternalError"},
			wantErr:                false,
		},
	}
	for _, tt := range tests {
		CheckSessionTimeOutFunc = tt.CheckSessionTimeOutFunc
		UpdateLastUsedTimeFunc = tt.UpdateLastUsedTimeFunc
		MarshalFunc = tt.MarshalFunc
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{}
			got, err := a.Update(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccount_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		req *accountproto.DeleteAccountRequest
	}
	common.SetUpMockConfig()
	tests := []struct {
		name                    string
		args                    args
		CheckSessionTimeOutFunc func(sessionToken string) (*asmodel.Session, *errors.Error)
		UpdateLastUsedTimeFunc  func(token string) error
		AccDeleteFunc           func(session *asmodel.Session, accountID string) response.RPC
		MarshalFunc             func(v any) ([]byte, error)
		want                    *accountproto.AccountResponse
		wantErr                 bool
	}{
		{
			name: "Session Timeout Error for 401(not valid session)",
			args: args{context.TODO(), &accountproto.DeleteAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(errors.InvalidAuthToken, "error: invalid token ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			AccDeleteFunc:          func(session *asmodel.Session, accountID string) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 401, StatusMessage: "Base.1.11.0.NoValidSession", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.NoValidSession\",\"Message\":\"There is no valid session established with the implementation.error while authorizing session token: error: invalid token \",\"Severity\":\"Critical\",\"Resolution\":\"Establish a session before attempting any operations.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "Session Timeout Error for 504(Service unavailable)",
			args: args{context.TODO(), &accountproto.DeleteAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, errors.PackError(5, "error: Service unavailable ", sessionToken)
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			AccDeleteFunc:          func(session *asmodel.Session, accountID string) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 503, StatusMessage: "Base.1.11.0.CouldNotEstablishConnection", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.CouldNotEstablishConnection\",\"Message\":\"The service failed to establish a connection with the URI 127.0.0.1:6379. error while authorizing session token: error: Service unavailable \",\"Severity\":\"Critical\",\"MessageArgs\":[\"127.0.0.1:6379\"],\"Resolution\":\"Ensure that the URI contains a valid and reachable node name, protocol information and other URI components.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "UpdateLastUsedTime error",
			args: args{context.TODO(), &accountproto.DeleteAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return e.New("fakeError") },
			AccDeleteFunc:          func(session *asmodel.Session, accountID string) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{StatusCode: 500, StatusMessage: "Base.1.11.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.11.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.11.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.error while updating last used time of session with token : fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:                false,
		},
		{
			name: "Marshall error",
			args: args{context.TODO(), &accountproto.DeleteAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			AccDeleteFunc:          func(session *asmodel.Session, accountID string) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, e.New("fakeError") },
			want:                   &accountproto.AccountResponse{StatusCode: 500, StatusMessage: "error while trying marshal the response body for delete account: fakeError"},
			wantErr:                false,
		},
		{
			name: "Pass case",
			args: args{context.TODO(), &accountproto.DeleteAccountRequest{}},
			CheckSessionTimeOutFunc: func(sessionToken string) (*asmodel.Session, *errors.Error) {
				return nil, nil
			},
			UpdateLastUsedTimeFunc: func(token string) error { return nil },
			AccDeleteFunc:          func(session *asmodel.Session, accountID string) response.RPC { return response.RPC{} },
			MarshalFunc:            func(v any) ([]byte, error) { return nil, nil },
			want:                   &accountproto.AccountResponse{},
			wantErr:                false,
		},
	}
	for _, tt := range tests {
		CheckSessionTimeOutFunc = tt.CheckSessionTimeOutFunc
		UpdateLastUsedTimeFunc = tt.UpdateLastUsedTimeFunc
		AccDeleteFunc = tt.AccDeleteFunc
		MarshalFunc = tt.MarshalFunc
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{}
			got, err := a.Delete(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}
