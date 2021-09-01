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

	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

//GetManagersCollection will do the rpc call to collect Managers
func GetManagersCollection(req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	conn, err := services.ODIMService.Client(services.Managers)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := managersproto.NewManagersClient(conn)
	resp, err := asService.GetManagersCollection(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	return resp, nil
}

// GetManagers will do the rpc calls for the svc-managers
func GetManagers(req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	conn, err := services.ODIMService.Client(services.Managers)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := managersproto.NewManagersClient(conn)
	resp, err := asService.GetManager(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	return resp, nil
}

// GetManagersResource will do the rpc calls for the svc-managers
func GetManagersResource(req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	conn, err := services.ODIMService.Client(services.Managers)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := managersproto.NewManagersClient(conn)
	resp, err := asService.GetManagersResource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	return resp, nil
}

// VirtualMediaInsert will do the rpc calls for the svc-managers
func VirtualMediaInsert(req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	conn, err := services.ODIMService.Client(services.Managers)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := managersproto.NewManagersClient(conn)
	resp, err := asService.VirtualMediaInsert(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// VirtualMediaEject will do the rpc calls for the svc-managers
func VirtualMediaEject(req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	conn, err := services.ODIMService.Client(services.Managers)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := managersproto.NewManagersClient(conn)
	resp, err := asService.VirtualMediaEject(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}
