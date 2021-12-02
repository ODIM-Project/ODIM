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

func GetActivePool(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error) {
	fmt.Errorf("In rpc.GetCompositionService")
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

func GetFreePool(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error) {
	fmt.Errorf("In rpc.GetCompositionService")
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
