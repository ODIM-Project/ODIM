//(C) Copyright [2020] Hewlett Packard Enterprise Development LP

//(C) Copyright 2020 Intel Corporation

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

	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"google.golang.org/grpc"
)

func TestGetChassisCollection(t *testing.T) {
	type args struct {
		req chassisproto.GetChassisRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewChassisClientFunc func(cc *grpc.ClientConn) chassisproto.ChassisClient
		want                 *chassisproto.GetChassisResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewChassisClientFunc: func(cc *grpc.ClientConn) chassisproto.ChassisClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "GetChassisCollection error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewChassisClientFunc: func(cc *grpc.ClientConn) chassisproto.ChassisClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewChassisClientFunc = tt.NewChassisClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetChassisCollection(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChassisCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChassisCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetChassisResource(t *testing.T) {
	type args struct {
		req chassisproto.GetChassisRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewChassisClientFunc func(cc *grpc.ClientConn) chassisproto.ChassisClient
		want                 *chassisproto.GetChassisResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewChassisClientFunc: func(cc *grpc.ClientConn) chassisproto.ChassisClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "GetChassisResource error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewChassisClientFunc: func(cc *grpc.ClientConn) chassisproto.ChassisClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewChassisClientFunc = tt.NewChassisClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetChassisResource(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChassisResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChassisResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetChassis(t *testing.T) {
	type args struct {
		req chassisproto.GetChassisRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewChassisClientFunc func(cc *grpc.ClientConn) chassisproto.ChassisClient
		want                 *chassisproto.GetChassisResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewChassisClientFunc: func(cc *grpc.ClientConn) chassisproto.ChassisClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "GetChassis error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewChassisClientFunc: func(cc *grpc.ClientConn) chassisproto.ChassisClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewChassisClientFunc = tt.NewChassisClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetChassis(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChassis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChassis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateChassis(t *testing.T) {
	type args struct {
		req chassisproto.CreateChassisRequest
	}
	tests := []struct {
		name string
		args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewChassisClientFunc func(cc *grpc.ClientConn) chassisproto.ChassisClient
		want                 *chassisproto.GetChassisResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewChassisClientFunc: func(cc *grpc.ClientConn) chassisproto.ChassisClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "CreateChassis error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewChassisClientFunc: func(cc *grpc.ClientConn) chassisproto.ChassisClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewChassisClientFunc = tt.NewChassisClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateChassis(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateChassis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateChassis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteChassis(t *testing.T) {
	type args struct {
		req chassisproto.DeleteChassisRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewChassisClientFunc func(cc *grpc.ClientConn) chassisproto.ChassisClient
		want                 *chassisproto.GetChassisResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewChassisClientFunc: func(cc *grpc.ClientConn) chassisproto.ChassisClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "DeleteChassis error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewChassisClientFunc: func(cc *grpc.ClientConn) chassisproto.ChassisClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewChassisClientFunc = tt.NewChassisClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteChassis(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteChassis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteChassis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateChassis(t *testing.T) {
	type args struct {
		req chassisproto.UpdateChassisRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewChassisClientFunc func(cc *grpc.ClientConn) chassisproto.ChassisClient
		want                 *chassisproto.GetChassisResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewChassisClientFunc: func(cc *grpc.ClientConn) chassisproto.ChassisClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "UpdateChassis error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewChassisClientFunc: func(cc *grpc.ClientConn) chassisproto.ChassisClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewChassisClientFunc = tt.NewChassisClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateChassis(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateChassis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateChassis() = %v, want %v", got, tt.want)
			}
		})
	}
}
