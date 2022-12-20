//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

package rpc

import (
	"context"
	"fmt"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	licenseproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/licenses"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

// GetLicenseCollection will do the rpc call to get License Service Information
func GetLicenseService(ctx context.Context, req licenseproto.GetLicenseServiceRequest) (*licenseproto.GetLicenseResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := services.ODIMService.Client(services.Licenses)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	licenseService := licenseproto.NewLicensesClient(conn)
	resp, err := licenseService.GetLicenseService(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	return resp, nil
}

// GetLicenseCollection will do the rpc call to get License Service Information
func GetLicenseCollection(ctx context.Context, req licenseproto.GetLicenseRequest) (*licenseproto.GetLicenseResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := services.ODIMService.Client(services.Licenses)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	licenseService := licenseproto.NewLicensesClient(conn)
	resp, err := licenseService.GetLicenseCollection(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	return resp, nil
}

// GetLicenseResource will do the rpc call to get License Service Information
func GetLicenseResource(ctx context.Context, req licenseproto.GetLicenseResourceRequest) (*licenseproto.GetLicenseResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := services.ODIMService.Client(services.Licenses)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	licenseService := licenseproto.NewLicensesClient(conn)
	resp, err := licenseService.GetLicenseResource(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	return resp, nil
}

// InstallLicenseService will do the rpc call to install License
func InstallLicenseService(ctx context.Context, req licenseproto.InstallLicenseRequest) (*licenseproto.GetLicenseResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := services.ODIMService.Client(services.Licenses)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	licenseService := licenseproto.NewLicensesClient(conn)
	resp, err := licenseService.InstallLicenseService(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}
