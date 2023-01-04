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

	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	"google.golang.org/grpc"
)

func TestGetRole(t *testing.T) {
	type args struct {
		req roleproto.GetRoleRequest
	}
	tests := []struct {
		name               string
		args               args
		ClientFunc         func(clientName string) (*grpc.ClientConn, error)
		NewRolesClientFunc func(cc *grpc.ClientConn) roleproto.RolesClient
		want               *roleproto.RoleResponse
		wantErr            bool
	}{
		{
			name:               "Client func error",
			args:               args{},
			ClientFunc:         func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewRolesClientFunc: func(cc *grpc.ClientConn) roleproto.RolesClient { return nil },
			want:               nil,
			wantErr:            true,
		},
		{
			name:               "GetRole error",
			args:               args{},
			ClientFunc:         func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewRolesClientFunc: func(cc *grpc.ClientConn) roleproto.RolesClient { return fakeStruct{} },
			want:               nil,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewRolesClientFunc = tt.NewRolesClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRole(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllRoles(t *testing.T) {
	type args struct {
		req roleproto.GetRoleRequest
	}
	tests := []struct {
		name               string
		args               args
		ClientFunc         func(clientName string) (*grpc.ClientConn, error)
		NewRolesClientFunc func(cc *grpc.ClientConn) roleproto.RolesClient
		want               *roleproto.RoleResponse
		wantErr            bool
	}{
		{
			name:               "Client func error",
			args:               args{},
			ClientFunc:         func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewRolesClientFunc: func(cc *grpc.ClientConn) roleproto.RolesClient { return nil },
			want:               nil,
			wantErr:            true,
		},
		{
			name:               "GetAllRole error",
			args:               args{},
			ClientFunc:         func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewRolesClientFunc: func(cc *grpc.ClientConn) roleproto.RolesClient { return fakeStruct{} },
			want:               nil,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewRolesClientFunc = tt.NewRolesClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllRoles(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateRole(t *testing.T) {
	type args struct {
		req roleproto.UpdateRoleRequest
	}
	tests := []struct {
		name               string
		args               args
		ClientFunc         func(clientName string) (*grpc.ClientConn, error)
		NewRolesClientFunc func(cc *grpc.ClientConn) roleproto.RolesClient
		want               *roleproto.RoleResponse
		wantErr            bool
	}{
		{
			name:               "Client func error",
			args:               args{},
			ClientFunc:         func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewRolesClientFunc: func(cc *grpc.ClientConn) roleproto.RolesClient { return nil },
			want:               nil,
			wantErr:            true,
		},
		{
			name:               "UpdateRole error",
			args:               args{},
			ClientFunc:         func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewRolesClientFunc: func(cc *grpc.ClientConn) roleproto.RolesClient { return fakeStruct{} },
			want:               nil,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewRolesClientFunc = tt.NewRolesClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateRole(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteRole(t *testing.T) {
	type args struct {
		req roleproto.DeleteRoleRequest
	}
	tests := []struct {
		name               string
		args               args
		ClientFunc         func(clientName string) (*grpc.ClientConn, error)
		NewRolesClientFunc func(cc *grpc.ClientConn) roleproto.RolesClient
		want               *roleproto.RoleResponse
		wantErr            bool
	}{
		{
			name:               "Client func error",
			args:               args{},
			ClientFunc:         func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewRolesClientFunc: func(cc *grpc.ClientConn) roleproto.RolesClient { return nil },
			want:               nil,
			wantErr:            true,
		},
		{
			name:               "DeleteRole error",
			args:               args{},
			ClientFunc:         func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewRolesClientFunc: func(cc *grpc.ClientConn) roleproto.RolesClient { return fakeStruct{} },
			want:               nil,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewRolesClientFunc = tt.NewRolesClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteRole(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
