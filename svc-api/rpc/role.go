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
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

var (
	NewRolesClientFunc = roleproto.NewRolesClient
)

// GetRole defines the RPC call function for
// the GetRole from account-session micro service
func GetRole(ctx context.Context, req roleproto.GetRoleRequest) (*roleproto.RoleResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewRolesClientFunc(conn)
	resp, err := asService.GetRole(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// GetAllRoles defines the RPC call function for
// the GetAllRoles from account-session micro service
func GetAllRoles(ctx context.Context, req roleproto.GetRoleRequest) (*roleproto.RoleResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewRolesClientFunc(conn)
	resp, err := asService.GetAllRoles(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// UpdateRole defines the RPC call function for
// the UpdateRole from account-session micro service
func UpdateRole(ctx context.Context, req roleproto.UpdateRoleRequest) (*roleproto.RoleResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewRolesClientFunc(conn)
	resp, err := asService.UpdateRole(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DeleteRole defines the RPC call function for the DeleteRole from account-session microservice
func DeleteRole(ctx context.Context, req roleproto.DeleteRoleRequest) (*roleproto.RoleResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewRolesClientFunc(conn)
	resp, err := asService.DeleteRole(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}
