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
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

var (
	NewAggregatorClientFunc = aggregatorproto.NewAggregatorClient
)

// DoGetAggregationService defines the RPC call function for
// the GetAggregationService from aggregator micro service
func DoGetAggregationService(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.GetAggregationService(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoResetRequest defines the RPC call function for
// the Reset from aggregator micro service
func DoResetRequest(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.Reset(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoSetDefaultBootOrderRequest defines the RPC call function for
// the SetDefaultBootOrder from aggregator micro service
func DoSetDefaultBootOrderRequest(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.SetDefaultBootOrder(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoAddAggregationSource defines the RPC call function for
// the AddAggregationSource from aggregator micro service
func DoAddAggregationSource(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.AddAggregationSource(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetAllAggregationSource defines the RPC call function for
// the GetAllAggregationSource from aggregator micro service
func DoGetAllAggregationSource(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.GetAllAggregationSource(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetAggregationSource defines the RPC call function for
// the GetAggregationSource from aggregator micro service
func DoGetAggregationSource(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.GetAggregationSource(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoUpdateAggregationSource defines the RPC call function for
// the UpdateAggregationSource from aggregator micro service
func DoUpdateAggregationSource(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.UpdateAggregationSource(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoDeleteAggregationSource defines the RPC call function for
// the DeleteAggregationSource  from aggregator micro service
func DoDeleteAggregationSource(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.DeleteAggregationSource(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoCreateAggregate defines the RPC call function for
// the CreateAggregate from aggregator micro service
func DoCreateAggregate(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.CreateAggregate(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetAggregateCollection defines the RPC call function for
// the get aggregate collections from aggregator micro service
func DoGetAggregateCollection(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.GetAllAggregates(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGeteAggregate defines the RPC call function for
// the get aggregate from aggregator micro service
func DoGeteAggregate(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.GetAggregate(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoDeleteAggregate defines the RPC call function for
// the delete aggregate from aggregator micro service
func DoDeleteAggregate(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.DeleteAggregate(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoAddElementsToAggregate defines the RPC call function for
// the add elements to an aggregate from aggregator micro service
func DoAddElementsToAggregate(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.AddElementsToAggregate(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoRemoveElementsFromAggregate defines the RPC call function for
// the remove elements from an aggregate from aggregator micro service
func DoRemoveElementsFromAggregate(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.RemoveElementsFromAggregate(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoResetAggregateElements defines the RPC call function for
// the reset elements of an aggregate from aggregator micro service
func DoResetAggregateElements(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.ResetElementsOfAggregate(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoSetDefaultBootOrderAggregateElements defines the RPC call function for
// the set default boot order elements of an aggregate from aggregator micro service
func DoSetDefaultBootOrderAggregateElements(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.SetDefaultBootOrderElementsOfAggregate(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetAllConnectionMethods defines the RPC call function for
// the get connection method collection from aggregator micro service
func DoGetAllConnectionMethods(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.GetAllConnectionMethods(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetConnectionMethod defines the RPC call function for
// the get on connection method from aggregator micro service
func DoGetConnectionMethod(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.GetConnectionMethod(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetResetActionInfoService defines the RPC call function for
// the GetResetActionInfoService from aggregator micro service
func DoGetResetActionInfoService(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.GetResetActionInfoService(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetSetDefaultBootOrderActionInfo defines the RPC call function for
// the GetSetDefaultBootOrderActionInfo from aggregator micro service
func DoGetSetDefaultBootOrderActionInfo(ctx context.Context, req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Aggregator)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	aggregator := NewAggregatorClientFunc(conn)

	resp, err := aggregator.GetSetDefaultBootOrderActionInfo(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}
