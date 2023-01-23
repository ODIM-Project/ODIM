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
	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

var (
	NewFabricsClientFunc = fabricsproto.NewFabricsClient
)

// GetFabricResource defines the RPC call function for
// the GetFabricResource from fabrics micro service
func GetFabricResource(ctx context.Context, req fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Fabrics)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	fab := NewFabricsClientFunc(conn)
	resp, err := fab.GetFabricResource(ctx, &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// UpdateFabricResource defines the RPC call function for creating/updating
// the Fabric Resource such as Endpoints, Zones from fabrics micro service
func UpdateFabricResource(ctx context.Context, req fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Fabrics)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	fab := NewFabricsClientFunc(conn)

	resp, err := fab.UpdateFabricResource(ctx, &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DeleteFabricResource defines the RPC call function for
// the DeleteFabricResource from fabrics micro service
func DeleteFabricResource(ctx context.Context, req fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Fabrics)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	fab := NewFabricsClientFunc(conn)

	resp, err := fab.DeleteFabricResource(ctx, &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}
	defer conn.Close()
	return resp, err
}
