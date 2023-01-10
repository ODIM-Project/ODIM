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

// Package rpc ...
package rpc

import (
	"context"
	"fmt"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

var (
	NewSystemsClientFunc = systemsproto.NewSystemsClient
)

// GetSystemsCollection will do the rpc call to collect Systems from odimra
func GetSystemsCollection(ctx context.Context, req systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewSystemsClientFunc(conn)
	resp, err := asService.GetSystemsCollection(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// GetSystemRequestRPC will do the rpc calls for the svc-systems
func GetSystemRequestRPC(ctx context.Context, req systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewSystemsClientFunc(conn)
	resp, err := asService.GetSystems(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// GetSystemResource will do the rpc call to collect System Resource
func GetSystemResource(ctx context.Context, req systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewSystemsClientFunc(conn)
	resp, err := asService.GetSystemResource(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// ComputerSystemReset will do the rpc call to reset the computer system
func ComputerSystemReset(ctx context.Context, req systemsproto.ComputerSystemResetRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewSystemsClientFunc(conn)
	resp, err := asService.ComputerSystemReset(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// SetDefaultBootOrder will do the rpc call to set the default boot order of computer system
func SetDefaultBootOrder(ctx context.Context, req systemsproto.DefaultBootOrderRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewSystemsClientFunc(conn)
	resp, err := asService.SetDefaultBootOrder(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// ChangeBiosSettings will do the rpc call to change bios settings
func ChangeBiosSettings(ctx context.Context, req systemsproto.BiosSettingsRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewSystemsClientFunc(conn)
	resp, err := asService.ChangeBiosSettings(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// ChangeBootOrderSettings will do the rpc call to change Boot Order settings
func ChangeBootOrderSettings(ctx context.Context, req systemsproto.BootOrderSettingsRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewSystemsClientFunc(conn)
	resp, err := asService.ChangeBootOrderSettings(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// CreateVolume will do the rpc call to create a volume under storage
func CreateVolume(ctx context.Context, req systemsproto.VolumeRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewSystemsClientFunc(conn)
	resp, err := asService.CreateVolume(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// DeleteVolume will do the rpc call to DeleteVolume a volume under storage
func DeleteVolume(ctx context.Context, req systemsproto.VolumeRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewSystemsClientFunc(conn)
	resp, err := asService.DeleteVolume(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}
