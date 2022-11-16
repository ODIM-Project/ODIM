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

	licenseproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/licenses"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

// GetLicenseCollection will do the rpc call to get License Service Information
func GetLicenseService(req licenseproto.GetLicenseServiceRequest) (*licenseproto.GetLicenseResponse, error) {
	conn, err := services.ODIMService.Client(services.Licenses)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	licenseService := licenseproto.NewLicensesClient(conn)
	resp, err := licenseService.GetLicenseService(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	return resp, nil
}

// GetLicenseCollection will do the rpc call to get License Service Information
func GetLicenseCollection(req licenseproto.GetLicenseRequest) (*licenseproto.GetLicenseResponse, error) {
	conn, err := services.ODIMService.Client(services.Licenses)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	licenseService := licenseproto.NewLicensesClient(conn)
	resp, err := licenseService.GetLicenseCollection(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	return resp, nil
}

// GetLicenseResource will do the rpc call to get License Service Information
func GetLicenseResource(req licenseproto.GetLicenseResourceRequest) (*licenseproto.GetLicenseResponse, error) {
	conn, err := services.ODIMService.Client(services.Licenses)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	licenseService := licenseproto.NewLicensesClient(conn)
	resp, err := licenseService.GetLicenseResource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	return resp, nil
}

// InstallLicenseService will do the rpc call to install License
func InstallLicenseService(req licenseproto.InstallLicenseRequest) (*licenseproto.GetLicenseResponse, error) {
	conn, err := services.ODIMService.Client(services.Licenses)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	licenseService := licenseproto.NewLicensesClient(conn)
	resp, err := licenseService.InstallLicenseService(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}
