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
	"errors"
	"reflect"
	"testing"

	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"google.golang.org/grpc"
)

func TestGetSystemsCollection(t *testing.T) {
	type args struct {
		req systemsproto.GetSystemsRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewSystemsClientFunc func(cc *grpc.ClientConn) systemsproto.SystemsClient
		want                 *systemsproto.SystemsResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "GetSystemsCollection error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return fakeStruct2{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewSystemsClientFunc = tt.NewSystemsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSystemsCollection(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSystemsCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSystemsCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSystemRequestRPC(t *testing.T) {
	type args struct {
		req systemsproto.GetSystemsRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewSystemsClientFunc func(cc *grpc.ClientConn) systemsproto.SystemsClient
		want                 *systemsproto.SystemsResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "GetSystemRequestRPC error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return fakeStruct2{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewSystemsClientFunc = tt.NewSystemsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSystemRequestRPC(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSystemRequestRPC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSystemRequestRPC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSystemResource(t *testing.T) {
	type args struct {
		req systemsproto.GetSystemsRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewSystemsClientFunc func(cc *grpc.ClientConn) systemsproto.SystemsClient
		want                 *systemsproto.SystemsResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "GetSystemResource error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return fakeStruct2{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewSystemsClientFunc = tt.NewSystemsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSystemResource(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSystemResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSystemResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComputerSystemReset(t *testing.T) {
	type args struct {
		req systemsproto.ComputerSystemResetRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewSystemsClientFunc func(cc *grpc.ClientConn) systemsproto.SystemsClient
		want                 *systemsproto.SystemsResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "ComputerSystemReset error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return fakeStruct2{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewSystemsClientFunc = tt.NewSystemsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := ComputerSystemReset(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComputerSystemReset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComputerSystemReset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetDefaultBootOrder(t *testing.T) {
	type args struct {
		req systemsproto.DefaultBootOrderRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewSystemsClientFunc func(cc *grpc.ClientConn) systemsproto.SystemsClient
		want                 *systemsproto.SystemsResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "SetDefaultBootOrder error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return fakeStruct2{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewSystemsClientFunc = tt.NewSystemsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetDefaultBootOrder(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetDefaultBootOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetDefaultBootOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChangeBiosSettings(t *testing.T) {
	type args struct {
		req systemsproto.BiosSettingsRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewSystemsClientFunc func(cc *grpc.ClientConn) systemsproto.SystemsClient
		want                 *systemsproto.SystemsResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "ChangeBiosSettings error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return fakeStruct2{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewSystemsClientFunc = tt.NewSystemsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChangeBiosSettings(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeBiosSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChangeBiosSettings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChangeBootOrderSettings(t *testing.T) {
	type args struct {
		req systemsproto.BootOrderSettingsRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewSystemsClientFunc func(cc *grpc.ClientConn) systemsproto.SystemsClient
		want                 *systemsproto.SystemsResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "ChangeBootOrderSettings error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return fakeStruct2{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewSystemsClientFunc = tt.NewSystemsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChangeBootOrderSettings(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeBootOrderSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChangeBootOrderSettings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateVolume(t *testing.T) {
	type args struct {
		req systemsproto.VolumeRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewSystemsClientFunc func(cc *grpc.ClientConn) systemsproto.SystemsClient
		want                 *systemsproto.SystemsResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "CreateVolume error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return fakeStruct2{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewSystemsClientFunc = tt.NewSystemsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateVolume(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteVolume(t *testing.T) {
	type args struct {
		req systemsproto.VolumeRequest
	}
	tests := []struct {
		name                 string
		args                 args
		ClientFunc           func(clientName string) (*grpc.ClientConn, error)
		NewSystemsClientFunc func(cc *grpc.ClientConn) systemsproto.SystemsClient
		want                 *systemsproto.SystemsResponse
		wantErr              bool
	}{
		{
			name:                 "Client func error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return nil },
			want:                 nil,
			wantErr:              true,
		},
		{
			name:                 "DeleteVolume error",
			args:                 args{},
			ClientFunc:           func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewSystemsClientFunc: func(cc *grpc.ClientConn) systemsproto.SystemsClient { return fakeStruct2{} },
			want:                 nil,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewSystemsClientFunc = tt.NewSystemsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteVolume(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}
