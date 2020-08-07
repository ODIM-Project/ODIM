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

	aggregator := aggregatorproto.NewAggregatorService(services.Aggregator, services.Service.Client())

	resp, err := aggregator.GetAggregationService(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoAddComputeRequest defines the RPC call function for
// the AddCompute from aggregator micro service
func DoAddComputeRequest(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	aggregator := aggregatorproto.NewAggregatorService(services.Aggregator, services.Service.Client())

	resp, err := aggregator.AddCompute(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoDeleteComputeRequest defines the RPC call function for
// the DeleteCompute from aggregator micro service
func DoDeleteComputeRequest(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	aggregator := aggregatorproto.NewAggregatorService(services.Aggregator, services.Service.Client())

	resp, err := aggregator.DeleteCompute(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoResetRequest defines the RPC call function for
// the Reset from aggregator micro service
func DoResetRequest(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	aggregator := aggregatorproto.NewAggregatorService(services.Aggregator, services.Service.Client())

	resp, err := aggregator.Reset(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoSetDefaultBootOrderRequest defines the RPC call function for
// the SetDefaultBootOrder from aggregator micro service
func DoSetDefaultBootOrderRequest(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	aggregator := aggregatorproto.NewAggregatorService(services.Aggregator, services.Service.Client())

	resp, err := aggregator.SetDefaultBootOrder(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoAddAggregationSource defines the RPC call function for
// the AddAggregationSource from aggregator micro service
func DoAddAggregationSource(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	aggregator := aggregatorproto.NewAggregatorService(services.Aggregator, services.Service.Client())

	resp, err := aggregator.AddAggregationSource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoGetAllAggregationSource defines the RPC call function for
// the GetAllAggregationSource from aggregator micro service
func DoGetAllAggregationSource(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	aggregator := aggregatorproto.NewAggregatorService(services.Aggregator, services.Service.Client())

	resp, err := aggregator.GetAllAggregationSource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoGetAggregationSource defines the RPC call function for
// the GetAggregationSource from aggregator micro service
func DoGetAggregationSource(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	aggregator := aggregatorproto.NewAggregatorService(services.Aggregator, services.Service.Client())

	resp, err := aggregator.GetAggregationSource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}

// DoUpdateAggregationSource defines the RPC call function for
// the UpdateAggregationSource from aggregator micro service
func DoUpdateAggregationSource(req aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error) {

	aggregator := aggregatorproto.NewAggregatorService(services.Aggregator, services.Service.Client())

	resp, err := aggregator.UpdateAggregationSource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}

	return resp, err
}
