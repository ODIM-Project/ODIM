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

	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"google.golang.org/grpc"
)

func TestDoGetUpdateService(t *testing.T) {
	type args struct {
		req updateproto.UpdateRequest
	}
	tests := []struct {
		name                string
		args                args
		ClientFunc          func(clientName string) (*grpc.ClientConn, error)
		NewUpdateClientFunc func(cc *grpc.ClientConn) updateproto.UpdateClient
		want                *updateproto.UpdateResponse
		wantErr             bool
	}{
		{
			name:                "Client func error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewUpdateClientFunc: func(cc *grpc.ClientConn) updateproto.UpdateClient { return nil },
			want:                nil,
			wantErr:             true,
		},
		{
			name:                "DoGetUpdateService error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewUpdateClientFunc: func(cc *grpc.ClientConn) updateproto.UpdateClient { return fakeStruct{} },
			want:                nil,
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewUpdateClientFunc = tt.NewUpdateClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetUpdateService(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetUpdateService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetUpdateService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetFirmwareInventory(t *testing.T) {
	type args struct {
		req updateproto.UpdateRequest
	}
	tests := []struct {
		name                string
		args                args
		ClientFunc          func(clientName string) (*grpc.ClientConn, error)
		NewUpdateClientFunc func(cc *grpc.ClientConn) updateproto.UpdateClient
		want                *updateproto.UpdateResponse
		wantErr             bool
	}{
		{
			name:                "Client func error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewUpdateClientFunc: func(cc *grpc.ClientConn) updateproto.UpdateClient { return nil },
			want:                nil,
			wantErr:             true,
		},
		{
			name:                "DoGetFirmwareInventory error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewUpdateClientFunc: func(cc *grpc.ClientConn) updateproto.UpdateClient { return fakeStruct{} },
			want:                nil,
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewUpdateClientFunc = tt.NewUpdateClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetFirmwareInventory(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetFirmwareInventory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetFirmwareInventory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetSoftwareInventory(t *testing.T) {
	type args struct {
		req updateproto.UpdateRequest
	}
	tests := []struct {
		name                string
		args                args
		ClientFunc          func(clientName string) (*grpc.ClientConn, error)
		NewUpdateClientFunc func(cc *grpc.ClientConn) updateproto.UpdateClient
		want                *updateproto.UpdateResponse
		wantErr             bool
	}{
		{
			name:                "Client func error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewUpdateClientFunc: func(cc *grpc.ClientConn) updateproto.UpdateClient { return nil },
			want:                nil,
			wantErr:             true,
		},
		{
			name:                "DoGetSoftwareInventory error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewUpdateClientFunc: func(cc *grpc.ClientConn) updateproto.UpdateClient { return fakeStruct{} },
			want:                nil,
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewUpdateClientFunc = tt.NewUpdateClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetSoftwareInventory(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetSoftwareInventory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetSoftwareInventory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetFirmwareInventoryCollection(t *testing.T) {
	type args struct {
		req updateproto.UpdateRequest
	}
	tests := []struct {
		name                string
		args                args
		ClientFunc          func(clientName string) (*grpc.ClientConn, error)
		NewUpdateClientFunc func(cc *grpc.ClientConn) updateproto.UpdateClient
		want                *updateproto.UpdateResponse
		wantErr             bool
	}{
		{
			name:                "Client func error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewUpdateClientFunc: func(cc *grpc.ClientConn) updateproto.UpdateClient { return nil },
			want:                nil,
			wantErr:             true,
		},
		{
			name:                "DoGetFirmwareInventoryCollection error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewUpdateClientFunc: func(cc *grpc.ClientConn) updateproto.UpdateClient { return fakeStruct{} },
			want:                nil,
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewUpdateClientFunc = tt.NewUpdateClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetFirmwareInventoryCollection(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetFirmwareInventoryCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetFirmwareInventoryCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetSoftwareInventoryCollection(t *testing.T) {
	type args struct {
		req updateproto.UpdateRequest
	}
	tests := []struct {
		name                string
		args                args
		ClientFunc          func(clientName string) (*grpc.ClientConn, error)
		NewUpdateClientFunc func(cc *grpc.ClientConn) updateproto.UpdateClient
		want                *updateproto.UpdateResponse
		wantErr             bool
	}{
		{
			name:                "Client func error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewUpdateClientFunc: func(cc *grpc.ClientConn) updateproto.UpdateClient { return nil },
			want:                nil,
			wantErr:             true,
		},
		{
			name:                "DoGetSoftwareInventoryCollection error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewUpdateClientFunc: func(cc *grpc.ClientConn) updateproto.UpdateClient { return fakeStruct{} },
			want:                nil,
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewUpdateClientFunc = tt.NewUpdateClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetSoftwareInventoryCollection(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetSoftwareInventoryCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetSoftwareInventoryCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoSimpleUpdate(t *testing.T) {
	type args struct {
		req updateproto.UpdateRequest
	}
	tests := []struct {
		name                string
		args                args
		ClientFunc          func(clientName string) (*grpc.ClientConn, error)
		NewUpdateClientFunc func(cc *grpc.ClientConn) updateproto.UpdateClient
		want                *updateproto.UpdateResponse
		wantErr             bool
	}{
		{
			name:                "Client func error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewUpdateClientFunc: func(cc *grpc.ClientConn) updateproto.UpdateClient { return nil },
			want:                nil,
			wantErr:             true,
		},
		{
			name:                "DoSimpleUpdate error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewUpdateClientFunc: func(cc *grpc.ClientConn) updateproto.UpdateClient { return fakeStruct{} },
			want:                nil,
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewUpdateClientFunc = tt.NewUpdateClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoSimpleUpdate(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoSimpleUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoSimpleUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoStartUpdate(t *testing.T) {
	type args struct {
		req updateproto.UpdateRequest
	}
	tests := []struct {
		name                string
		args                args
		ClientFunc          func(clientName string) (*grpc.ClientConn, error)
		NewUpdateClientFunc func(cc *grpc.ClientConn) updateproto.UpdateClient
		want                *updateproto.UpdateResponse
		wantErr             bool
	}{
		{
			name:                "Client func error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewUpdateClientFunc: func(cc *grpc.ClientConn) updateproto.UpdateClient { return nil },
			want:                nil,
			wantErr:             true,
		},
		{
			name:                "DoStartUpdate error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewUpdateClientFunc: func(cc *grpc.ClientConn) updateproto.UpdateClient { return fakeStruct{} },
			want:                nil,
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewUpdateClientFunc = tt.NewUpdateClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoStartUpdate(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoStartUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoStartUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
