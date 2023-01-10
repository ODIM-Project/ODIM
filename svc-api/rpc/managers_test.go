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

	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	"google.golang.org/grpc"
)

func TestGetManagersCollection(t *testing.T) {
	type args struct {
		req managersproto.ManagerRequest
	}
	tests := []struct {
		name                  string
		args                  args
		ClientFunc            func(clientName string) (*grpc.ClientConn, error)
		NewManagersClientFunc func(cc *grpc.ClientConn) managersproto.ManagersClient
		want                  *managersproto.ManagerResponse
		wantErr               bool
	}{
		{
			name:                  "Client func error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return nil },
			want:                  nil,
			wantErr:               true,
		},
		{
			name:                  "GetManagersCollection error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return fakeStruct{} },
			want:                  nil,
			wantErr:               true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewManagersClientFunc = tt.NewManagersClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetManagersCollection(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetManagersCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetManagersCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetManagers(t *testing.T) {
	type args struct {
		req managersproto.ManagerRequest
	}
	tests := []struct {
		name                  string
		args                  args
		ClientFunc            func(clientName string) (*grpc.ClientConn, error)
		NewManagersClientFunc func(cc *grpc.ClientConn) managersproto.ManagersClient
		want                  *managersproto.ManagerResponse
		wantErr               bool
	}{
		{
			name:                  "Client func error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return nil },
			want:                  nil,
			wantErr:               true,
		},
		{
			name:                  "GetManagers error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return fakeStruct{} },
			want:                  nil,
			wantErr:               true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewManagersClientFunc = tt.NewManagersClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetManagers(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetManagers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetManagers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetManagersResource(t *testing.T) {
	type args struct {
		req managersproto.ManagerRequest
	}
	tests := []struct {
		name                  string
		args                  args
		ClientFunc            func(clientName string) (*grpc.ClientConn, error)
		NewManagersClientFunc func(cc *grpc.ClientConn) managersproto.ManagersClient
		want                  *managersproto.ManagerResponse
		wantErr               bool
	}{
		{
			name:                  "Client func error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return nil },
			want:                  nil,
			wantErr:               true,
		},
		{
			name:                  "GetManagersResource error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return fakeStruct{} },
			want:                  nil,
			wantErr:               true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewManagersClientFunc = tt.NewManagersClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetManagersResource(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetManagersResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetManagersResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVirtualMediaInsert(t *testing.T) {
	type args struct {
		req managersproto.ManagerRequest
	}
	tests := []struct {
		name                  string
		args                  args
		ClientFunc            func(clientName string) (*grpc.ClientConn, error)
		NewManagersClientFunc func(cc *grpc.ClientConn) managersproto.ManagersClient
		want                  *managersproto.ManagerResponse
		wantErr               bool
	}{
		{
			name:                  "Client func error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return nil },
			want:                  nil,
			wantErr:               true,
		},
		{
			name:                  "VirtualMediaInsert error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return fakeStruct{} },
			want:                  nil,
			wantErr:               true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewManagersClientFunc = tt.NewManagersClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := VirtualMediaInsert(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("VirtualMediaInsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VirtualMediaInsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVirtualMediaEject(t *testing.T) {
	type args struct {
		req managersproto.ManagerRequest
	}
	tests := []struct {
		name                  string
		args                  args
		ClientFunc            func(clientName string) (*grpc.ClientConn, error)
		NewManagersClientFunc func(cc *grpc.ClientConn) managersproto.ManagersClient
		want                  *managersproto.ManagerResponse
		wantErr               bool
	}{
		{
			name:                  "Client func error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return nil },
			want:                  nil,
			wantErr:               true,
		},
		{
			name:                  "VirtualMediaEject error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return fakeStruct{} },
			want:                  nil,
			wantErr:               true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewManagersClientFunc = tt.NewManagersClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := VirtualMediaEject(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("VirtualMediaEject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VirtualMediaEject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRemoteAccountService(t *testing.T) {
	type args struct {
		req managersproto.ManagerRequest
	}
	tests := []struct {
		name                  string
		args                  args
		ClientFunc            func(clientName string) (*grpc.ClientConn, error)
		NewManagersClientFunc func(cc *grpc.ClientConn) managersproto.ManagersClient
		want                  *managersproto.ManagerResponse
		wantErr               bool
	}{
		{
			name:                  "Client func error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return nil },
			want:                  nil,
			wantErr:               true,
		},
		{
			name:                  "GetRemoteAccountService error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return fakeStruct{} },
			want:                  nil,
			wantErr:               true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewManagersClientFunc = tt.NewManagersClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRemoteAccountService(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRemoteAccountService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRemoteAccountService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateRemoteAccountService(t *testing.T) {
	type args struct {
		req managersproto.ManagerRequest
	}
	tests := []struct {
		name                  string
		args                  args
		ClientFunc            func(clientName string) (*grpc.ClientConn, error)
		NewManagersClientFunc func(cc *grpc.ClientConn) managersproto.ManagersClient
		want                  *managersproto.ManagerResponse
		wantErr               bool
	}{
		{
			name:                  "Client func error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return nil },
			want:                  nil,
			wantErr:               true,
		},
		{
			name:                  "CreateRemoteAccountService error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return fakeStruct{} },
			want:                  nil,
			wantErr:               true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewManagersClientFunc = tt.NewManagersClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateRemoteAccountService(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRemoteAccountService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateRemoteAccountService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateRemoteAccountService(t *testing.T) {
	type args struct {
		req managersproto.ManagerRequest
	}
	tests := []struct {
		name                  string
		args                  args
		ClientFunc            func(clientName string) (*grpc.ClientConn, error)
		NewManagersClientFunc func(cc *grpc.ClientConn) managersproto.ManagersClient
		want                  *managersproto.ManagerResponse
		wantErr               bool
	}{
		{
			name:                  "Client func error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return nil },
			want:                  nil,
			wantErr:               true,
		},
		{
			name:                  "UpdateRemoteAccountService error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return fakeStruct{} },
			want:                  nil,
			wantErr:               true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewManagersClientFunc = tt.NewManagersClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateRemoteAccountService(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRemoteAccountService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateRemoteAccountService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteRemoteAccountService(t *testing.T) {
	type args struct {
		req managersproto.ManagerRequest
	}
	tests := []struct {
		name                  string
		args                  args
		ClientFunc            func(clientName string) (*grpc.ClientConn, error)
		NewManagersClientFunc func(cc *grpc.ClientConn) managersproto.ManagersClient
		want                  *managersproto.ManagerResponse
		wantErr               bool
	}{
		{
			name:                  "Client func error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return nil },
			want:                  nil,
			wantErr:               true,
		},
		{
			name:                  "DeleteRemoteAccountService error",
			args:                  args{},
			ClientFunc:            func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewManagersClientFunc: func(cc *grpc.ClientConn) managersproto.ManagersClient { return fakeStruct{} },
			want:                  nil,
			wantErr:               true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewManagersClientFunc = tt.NewManagersClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteRemoteAccountService(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteRemoteAccountService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteRemoteAccountService() = %v, want %v", got, tt.want)
			}
		})
	}
}
