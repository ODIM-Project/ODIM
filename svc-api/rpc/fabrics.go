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

	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

// GetFabricResource defines the RPC call function for
// the GetFabricResource from fabrics micro service
func GetFabricResource(req fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error) {

	conn, err := services.ODIMService.Client(services.Fabrics)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	fab := fabricsproto.NewFabricsClient(conn)
	resp, err := fab.GetFabricResource(context.TODO(), &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}

	return resp, err
}

// UpdateFabricResource defines the RPC call function for creating/updating
// the Fabric Resource such as Endpoints, Zones from fabrics micro service
func UpdateFabricResource(req fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error) {

	conn, err := services.ODIMService.Client(services.Fabrics)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	fab := fabricsproto.NewFabricsClient(conn)

	resp, err := fab.UpdateFabricResource(context.TODO(), &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}

	return resp, err
}

// DeleteFabricResource defines the RPC call function for
// the DeleteFabricResource from fabrics micro service
func DeleteFabricResource(req fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error) {

	conn, err := services.ODIMService.Client(services.Fabrics)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	fab := fabricsproto.NewFabricsClient(conn)

	resp, err := fab.DeleteFabricResource(context.TODO(), &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}

	return resp, err
}
