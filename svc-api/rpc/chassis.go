//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//(C) Copyright 2020 Intel Corporation
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
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

var (
	NewChassisClientFunc = chassisproto.NewChassisClient
)

// GetChassisCollection will do the rpc call to collect all chassis
func GetChassisCollection(ctx context.Context, req chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewChassisClientFunc(conn)
	resp, err := asService.GetChassisCollection(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// GetChassisResource will do the rpc call to collect Chassis Resource
func GetChassisResource(ctx context.Context, req chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewChassisClientFunc(conn)
	resp, err := asService.GetChassisResource(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// GetChassis will do the rpc call to  System Resource
func GetChassis(ctx context.Context, req chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewChassisClientFunc(conn)
	resp, err := asService.GetChassisInfo(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// CreateChassis will do the rpc call to create a Chassis
func CreateChassis(ctx context.Context, req chassisproto.CreateChassisRequest) (*chassisproto.GetChassisResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	service := NewChassisClientFunc(conn)
	resp, err := service.CreateChassis(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// DeleteChassis will do the rpc call to delete a chassis
func DeleteChassis(ctx context.Context, req chassisproto.DeleteChassisRequest) (*chassisproto.GetChassisResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	service := NewChassisClientFunc(conn)
	resp, err := service.DeleteChassis(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}

// UpdateChassis will do the rpc call to update a chassis
func UpdateChassis(ctx context.Context, req chassisproto.UpdateChassisRequest) (*chassisproto.GetChassisResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Systems)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	service := NewChassisClientFunc(conn)
	resp, err := service.UpdateChassis(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, nil
}
