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

//Package rpc ...

package rpc

import (
	"context"
	"errors"
	"reflect"
	"testing"

	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"google.golang.org/grpc"
)

func TestDoSessionCreationRequest(t *testing.T) {
	type args struct {
		req sessionproto.SessionCreateRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewSessionClientFunc func(cc *grpc.ClientConn) sessionproto.SessionClient
		want                 *sessionproto.SessionCreateResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewSessionClientFunc: func(cc *grpc.ClientConn) sessionproto.SessionClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "CreateSession error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewSessionClientFunc: func(cc *grpc.ClientConn) sessionproto.SessionClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewSessionClientFunc = tt.NewSessionClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoSessionCreationRequest(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoSessionCreationRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoSessionCreationRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteSessionRequest(t *testing.T) {
	type args struct {
		sessionID    string
		sessionToken string
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewSessionClientFunc func(cc *grpc.ClientConn) sessionproto.SessionClient
		want                 *sessionproto.SessionResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewSessionClientFunc: func(cc *grpc.ClientConn) sessionproto.SessionClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "DeleteSessionRequest error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewSessionClientFunc: func(cc *grpc.ClientConn) sessionproto.SessionClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewSessionClientFunc = tt.NewSessionClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteSessionRequest(context.Background(), tt.args.sessionID, tt.args.sessionToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteSessionRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteSessionRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSessionRequest(t *testing.T) {
	type args struct {
		sessionID    string
		sessionToken string
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewSessionClientFunc func(cc *grpc.ClientConn) sessionproto.SessionClient
		want                 *sessionproto.SessionResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewSessionClientFunc: func(cc *grpc.ClientConn) sessionproto.SessionClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "GetSessionRequest error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewSessionClientFunc: func(cc *grpc.ClientConn) sessionproto.SessionClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewSessionClientFunc = tt.NewSessionClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSessionRequest(context.Background(), tt.args.sessionID, tt.args.sessionToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSessionRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSessionRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllActiveSessionRequest(t *testing.T) {
	type args struct {
		sessionID    string
		sessionToken string
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewSessionClientFunc func(cc *grpc.ClientConn) sessionproto.SessionClient
		want                 *sessionproto.SessionResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewSessionClientFunc: func(cc *grpc.ClientConn) sessionproto.SessionClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "GetAllActiveSessionRequest error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewSessionClientFunc: func(cc *grpc.ClientConn) sessionproto.SessionClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewSessionClientFunc = tt.NewSessionClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllActiveSessionRequest(context.Background(), tt.args.sessionID, tt.args.sessionToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllActiveSessionRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllActiveSessionRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSessionServiceRequest(t *testing.T) {
	tests := []struct {
		name                 string
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewSessionClientFunc func(cc *grpc.ClientConn) sessionproto.SessionClient
		want                 *sessionproto.SessionResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewSessionClientFunc: func(cc *grpc.ClientConn) sessionproto.SessionClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "GetSessionServiceRequest error",
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewSessionClientFunc: func(cc *grpc.ClientConn) sessionproto.SessionClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewSessionClientFunc = tt.NewSessionClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSessionServiceRequest(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSessionServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSessionServiceRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
