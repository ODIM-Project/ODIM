//(C) Copyright [2022] American Megatrends International LLC
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

	compositionserviceproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/compositionservice"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

// GetCompositionService will do the rpc call to get Composition Service Information
func GetCompositionService(req compositionserviceproto.GetCompositionServiceRequest) (*compositionserviceproto.CompositionServiceResponse, error) {
	conn, err := services.ODIMService.Client(services.CompositionService)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	csService := compositionserviceproto.NewCompositionClient(conn)
	resp, err := csService.GetCompositionService(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// GetResourceBlockCollection will do the rpc call to get Resource Block collection
func GetResourceBlockCollection(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error) {
	conn, err := services.ODIMService.Client(services.CompositionService)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	csService := compositionserviceproto.NewCompositionClient(conn)
	resp, err := csService.GetCompositionResource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// GetResourceBlock will do the rpc call to get Resource Block Instance
func GetResourceBlock(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error) {
	conn, err := services.ODIMService.Client(services.CompositionService)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	csService := compositionserviceproto.NewCompositionClient(conn)
	resp, err := csService.GetCompositionResource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// CreateResourceBlock will do the rpc call to Create Resource Block
func CreateResourceBlock(req compositionserviceproto.CreateCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error) {
	conn, err := services.ODIMService.Client(services.CompositionService)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	csService := compositionserviceproto.NewCompositionClient(conn)
	resp, err := csService.CreateCompositionResource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// DeleteResourceBlock will do the rpc call to Delete Resource Block
func DeleteResourceBlock(req compositionserviceproto.DeleteCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error) {
	conn, err := services.ODIMService.Client(services.CompositionService)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	csService := compositionserviceproto.NewCompositionClient(conn)
	resp, err := csService.DeleteCompositionResource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// GetResourceZoneCollection will do the rpc call to get Resource Zone Collection
func GetResourceZoneCollection(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error) {
	conn, err := services.ODIMService.Client(services.CompositionService)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	csService := compositionserviceproto.NewCompositionClient(conn)
	resp, err := csService.GetCompositionResource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// GetResourceZone will do the rpc call to get Resource Zone Instance
func GetResourceZone(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error) {
	conn, err := services.ODIMService.Client(services.CompositionService)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	csService := compositionserviceproto.NewCompositionClient(conn)
	resp, err := csService.GetCompositionResource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// CreateResourceZone will do the rpc call to Create Resource Zone
func CreateResourceZone(req compositionserviceproto.CreateCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error) {
	conn, err := services.ODIMService.Client(services.CompositionService)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	csService := compositionserviceproto.NewCompositionClient(conn)
	resp, err := csService.CreateCompositionResource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// DeleteResourceZone will do the rpc call to Delete Resource Zone
func DeleteResourceZone(req compositionserviceproto.DeleteCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error) {
	conn, err := services.ODIMService.Client(services.CompositionService)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	csService := compositionserviceproto.NewCompositionClient(conn)
	resp, err := csService.DeleteCompositionResource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// Compose will do the rpc call to Compose a system
func Compose(req compositionserviceproto.ComposeRequest) (*compositionserviceproto.CompositionServiceResponse, error) {
	conn, err := services.ODIMService.Client(services.CompositionService)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	csService := compositionserviceproto.NewCompositionClient(conn)
	resp, err := csService.Compose(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// GetActivePool will do the rpc call to list out the Active pool Resource block instance collection
func GetActivePool(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error) {
	conn, err := services.ODIMService.Client(services.CompositionService)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	csService := compositionserviceproto.NewCompositionClient(conn)
	resp, err := csService.GetCompositionResource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// GetFreePool will do the rpc call to list out the Free pool Resource block instance collection
func GetFreePool(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error) {
	conn, err := services.ODIMService.Client(services.CompositionService)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	csService := compositionserviceproto.NewCompositionClient(conn)
	resp, err := csService.GetCompositionResource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

// GetCompositionReservations will do the rpc call to list out the Compose action Reservation collection
func GetCompositionReservations(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error) {
	conn, err := services.ODIMService.Client(services.CompositionService)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	csService := compositionserviceproto.NewCompositionClient(conn)
	resp, err := csService.GetCompositionResource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}
