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

	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

// DoGetAggregationService defines the RPC call function for
// the GetAggregationService from aggregator micro service
func DoGetAggregationService(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.GetAggregationService(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoResetRequest defines the RPC call function for
// the Reset from aggregator micro service
func DoResetRequest(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.Reset(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoSetDefaultBootOrderRequest defines the RPC call function for
// the SetDefaultBootOrder from aggregator micro service
func DoSetDefaultBootOrderRequest(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.SetDefaultBootOrder(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoAddAggregationSource defines the RPC call function for
// the AddAggregationSource from aggregator micro service
func DoAddAggregationSource(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.AddAggregationSource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoGetAllAggregationSource defines the RPC call function for
// the GetAllAggregationSource from aggregator micro service
func DoGetAllAggregationSource(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.GetAllAggregationSource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoGetAggregationSource defines the RPC call function for
// the GetAggregationSource from aggregator micro service
func DoGetAggregationSource(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.GetAggregationSource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoUpdateAggregationSource defines the RPC call function for
// the UpdateAggregationSource from aggregator micro service
func DoUpdateAggregationSource(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.UpdateAggregationSource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoDeleteAggregationSource defines the RPC call function for
// the DeleteAggregationSource  from aggregator micro service
func DoDeleteAggregationSource(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.DeleteAggregationSource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoCreateAggregate defines the RPC call function for
// the CreateAggregate from aggregator micro service
func DoCreateAggregate(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.CreateAggregate(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoGetAggregateCollection defines the RPC call function for
// the get aggregate collections from aggregator micro service
func DoGetAggregateCollection(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.GetAllAggregates(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoGeteAggregate defines the RPC call function for
// the get aggregate from aggregator micro service
func DoGeteAggregate(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.GetAggregate(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoDeleteAggregate defines the RPC call function for
// the delete aggregate from aggregator micro service
func DoDeleteAggregate(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.DeleteAggregate(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoAddElementsToAggregate defines the RPC call function for
// the add elements to an aggregate from aggregator micro service
func DoAddElementsToAggregate(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.AddElementsToAggregate(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoRemoveElementsFromAggregate defines the RPC call function for
// the remove elements from an aggregate from aggregator micro service
func DoRemoveElementsFromAggregate(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.RemoveElementsFromAggregate(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoResetAggregateElements defines the RPC call function for
// the reset elements of an aggregate from aggregator micro service
func DoResetAggregateElements(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.ResetElementsOfAggregate(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoSetDefaultBootOrderAggregateElements defines the RPC call function for
// the set default boot order elements of an aggregate from aggregator micro service
func DoSetDefaultBootOrderAggregateElements(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.SetDefaultBootOrderElementsOfAggregate(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoGetAllConnectionMethods defines the RPC call function for
// the get connection method collection from aggregator micro service
func DoGetAllConnectionMethods(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.GetAllConnectionMethods(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoGetConnectionMethod defines the RPC call function for
// the get on connection method from aggregator micro service
func DoGetConnectionMethod(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	resp, err := aggregator.GetConnectionMethod(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}
