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
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

//GetSystemsCollection will do the rpc call to collect Systems from odimra
func GetSystemsCollection(req systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error) {
	conn, err := services.ODIMService.Client(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := systemsproto.NewSystemsClient(conn)
	resp, err := asService.GetSystemsCollection(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// GetSystemRequestRPC will do the rpc calls for the svc-systems
func GetSystemRequestRPC(req systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error) {
	conn, err := services.ODIMService.Client(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := systemsproto.NewSystemsClient(conn)
	resp, err := asService.GetSystems(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

//GetSystemResource will do the rpc call to collect System Resource
func GetSystemResource(req systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error) {
	conn, err := services.ODIMService.Client(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := systemsproto.NewSystemsClient(conn)
	resp, err := asService.GetSystemResource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// ComputerSystemReset will do the rpc call to reset the computer system
func ComputerSystemReset(req systemsproto.ComputerSystemResetRequest) (*systemsproto.SystemsResponse, error) {
	conn, err := services.ODIMService.Client(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := systemsproto.NewSystemsClient(conn)
	resp, err := asService.ComputerSystemReset(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// SetDefaultBootOrder will do the rpc call to set the default boot order of computer system
func SetDefaultBootOrder(req systemsproto.DefaultBootOrderRequest) (*systemsproto.SystemsResponse, error) {
	conn, err := services.ODIMService.Client(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := systemsproto.NewSystemsClient(conn)
	resp, err := asService.SetDefaultBootOrder(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// ChangeBiosSettings will do the rpc call to change bios settings
func ChangeBiosSettings(req systemsproto.BiosSettingsRequest) (*systemsproto.SystemsResponse, error) {
	conn, err := services.ODIMService.Client(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := systemsproto.NewSystemsClient(conn)
	resp, err := asService.ChangeBiosSettings(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// ChangeBootOrderSettings will do the rpc call to change Boot Order settings
func ChangeBootOrderSettings(req systemsproto.BootOrderSettingsRequest) (*systemsproto.SystemsResponse, error) {
	conn, err := services.ODIMService.Client(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := systemsproto.NewSystemsClient(conn)
	resp, err := asService.ChangeBootOrderSettings(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// CreateVolume will do the rpc call to create a volume under storage
func CreateVolume(req systemsproto.VolumeRequest) (*systemsproto.SystemsResponse, error) {
	conn, err := services.ODIMService.Client(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := systemsproto.NewSystemsClient(conn)
	resp, err := asService.CreateVolume(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// DeleteVolume will do the rpc call to DeleteVolume a volume under storage
func DeleteVolume(req systemsproto.VolumeRequest) (*systemsproto.SystemsResponse, error) {
	conn, err := services.ODIMService.Client(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := systemsproto.NewSystemsClient(conn)
	resp, err := asService.DeleteVolume(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}
