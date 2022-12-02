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

	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"google.golang.org/grpc"
)

func TestGetFabricResource(t *testing.T) {
	type args struct {
		req fabricsproto.FabricRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewFabricsClientFunc func(cc *grpc.ClientConn) fabricsproto.FabricsClient
		want                 *fabricsproto.FabricResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewFabricsClientFunc: func(cc *grpc.ClientConn) fabricsproto.FabricsClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "GetFabricResource error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewFabricsClientFunc: func(cc *grpc.ClientConn) fabricsproto.FabricsClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewFabricsClientFunc = tt.NewFabricsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFabricResource(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFabricResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFabricResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateFabricResource(t *testing.T) {
	type args struct {
		req fabricsproto.FabricRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewFabricsClientFunc func(cc *grpc.ClientConn) fabricsproto.FabricsClient
		want                 *fabricsproto.FabricResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewFabricsClientFunc: func(cc *grpc.ClientConn) fabricsproto.FabricsClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "UpdateFabricResource error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewFabricsClientFunc: func(cc *grpc.ClientConn) fabricsproto.FabricsClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewFabricsClientFunc = tt.NewFabricsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateFabricResource(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateFabricResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateFabricResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteFabricResource(t *testing.T) {
	type args struct {
		req fabricsproto.FabricRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewFabricsClientFunc func(cc *grpc.ClientConn) fabricsproto.FabricsClient
		want                 *fabricsproto.FabricResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewFabricsClientFunc: func(cc *grpc.ClientConn) fabricsproto.FabricsClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "DeleteFabricResource error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewFabricsClientFunc: func(cc *grpc.ClientConn) fabricsproto.FabricsClient { return fakeStruct{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewFabricsClientFunc = tt.NewFabricsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteFabricResource(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFabricResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteFabricResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
