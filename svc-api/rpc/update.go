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
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

var (
	NewUpdateClientFunc = updateproto.NewUpdateClient
)

// DoGetUpdateService defines the RPC call function for
// the GetUpdateService from update micro service
func DoGetUpdateService(ctx context.Context, req updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Update)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	update := NewUpdateClientFunc(conn)

	resp, err := update.GetUpdateService(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetFirmwareInventory defines the RPC call function for
// the GetFirmwareInventory from update micro service
func DoGetFirmwareInventory(ctx context.Context, req updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Update)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	update := NewUpdateClientFunc(conn)

	resp, err := update.GetFirmwareInventory(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetSoftwareInventory defines the RPC call function for
// the GetSoftwareInventory from update micro service
func DoGetSoftwareInventory(ctx context.Context, req updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Update)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	update := NewUpdateClientFunc(conn)

	resp, err := update.GetSoftwareInventory(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetFirmwareInventoryCollection defines the RPC call function for
// the GetFirmwareInventory from update micro service
func DoGetFirmwareInventoryCollection(ctx context.Context, req updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Update)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	update := NewUpdateClientFunc(conn)

	resp, err := update.GetFirmwareInventoryCollection(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetSoftwareInventoryCollection defines the RPC call function for
// the GetSoftwareInventory from update micro service
func DoGetSoftwareInventoryCollection(ctx context.Context, req updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Update)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	update := NewUpdateClientFunc(conn)

	resp, err := update.GetSoftwareInventoryCollection(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoSimpleUpdate defines the RPC call for
// SimpleUpdate from update micro service
func DoSimpleUpdate(ctx context.Context, req updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Update)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	update := NewUpdateClientFunc(conn)

	resp, err := update.SimepleUpdate(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoStartUpdate defines the RPC call for
// StartUpdate from update micro service
func DoStartUpdate(ctx context.Context, req updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Update)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	update := NewUpdateClientFunc(conn)

	resp, err := update.StartUpdate(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}
