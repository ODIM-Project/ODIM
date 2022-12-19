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
	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

var (
	NewManagersClientFunc = managersproto.NewManagersClient
)

// GetManagersCollection will do the rpc call to collect Managers
func GetManagersCollection(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Managers)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewManagersClientFunc(conn)
	resp, err := asService.GetManagersCollection(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// GetManagers will do the rpc calls for the svc-managers
func GetManagers(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Managers)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewManagersClientFunc(conn)
	resp, err := asService.GetManager(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// GetManagersResource will do the rpc calls for the svc-managers
func GetManagersResource(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Managers)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewManagersClientFunc(conn)
	resp, err := asService.GetManagersResource(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// VirtualMediaInsert will do the rpc calls for the svc-managers
func VirtualMediaInsert(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Managers)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewManagersClientFunc(conn)
	resp, err := asService.VirtualMediaInsert(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// VirtualMediaEject will do the rpc calls for the svc-managers
func VirtualMediaEject(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Managers)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewManagersClientFunc(conn)
	resp, err := asService.VirtualMediaEject(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// GetRemoteAccountService will do the rpc call to collect BMC accounts
func GetRemoteAccountService(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Managers)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	mService := NewManagersClientFunc(conn)
	resp, err := mService.GetRemoteAccountService(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// CreateRemoteAccountService will do the rpc call to create a new BMC account
func CreateRemoteAccountService(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Managers)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	mService := NewManagersClientFunc(conn)
	resp, err := mService.CreateRemoteAccountService(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// UpdateRemoteAccountService will do rpc call to update BMC account
func UpdateRemoteAccountService(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Managers)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	mService := NewManagersClientFunc(conn)
	resp, err := mService.UpdateRemoteAccountService(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// DeleteRemoteAccountService will do the rpc call to delete an existing BMC account
func DeleteRemoteAccountService(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Managers)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	mService := NewManagersClientFunc(conn)
	resp, err := mService.DeleteRemoteAccountService(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}
