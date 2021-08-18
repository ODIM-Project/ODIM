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
	"fmt"

	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

// CreateRole defines the RPC call function for
// the CreateRole from account-session micro service
func CreateRole(req roleproto.RoleRequest) (*roleproto.RoleResponse, error) {
	conn, err := services.ODIMService.Client(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := roleproto.NewRolesClient(conn)
	resp, err := asService.CreateRole(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, err
}

// GetRole defines the RPC call function for
// the GetRole from account-session micro service
func GetRole(req roleproto.GetRoleRequest) (*roleproto.RoleResponse, error) {
	conn, err := services.ODIMService.Client(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := roleproto.NewRolesClient(conn)
	resp, err := asService.GetRole(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// GetAllRoles defines the RPC call function for
// the GetAllRoles from account-session micro service
func GetAllRoles(req roleproto.GetRoleRequest) (*roleproto.RoleResponse, error) {
	conn, err := services.ODIMService.Client(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := roleproto.NewRolesClient(conn)
	resp, err := asService.GetAllRoles(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// UpdateRole defines the RPC call function for
// the UpdateRole from account-session micro service
func UpdateRole(req roleproto.UpdateRoleRequest) (*roleproto.RoleResponse, error) {
	conn, err := services.ODIMService.Client(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := roleproto.NewRolesClient(conn)
	resp, err := asService.UpdateRole(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DeleteRole defines the RPC call function for the DeleteRole from account-session microservice
func DeleteRole(req roleproto.DeleteRoleRequest) (*roleproto.RoleResponse, error) {
	conn, err := services.ODIMService.Client(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := roleproto.NewRolesClient(conn)
	resp, err := asService.DeleteRole(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}
