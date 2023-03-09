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
	"encoding/json"
	"os"

	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-fabrics/fabrics"
)

var podName = os.Getenv("POD_NAME")

// Fabrics struct helps to register service
type Fabrics struct {
	IsAuthorizedRPC  func(sessionToken string, privileges []string, oemPrivileges []string) (response.RPC, error)
	ContactClientRPC func(context.Context, string, string, string, string, interface{}, map[string]string) (*http.Response, error)
}

// GetFabricResource defines the operation which handled the RPC request response
// for getting the specified fabric resource information
func (f *Fabrics) GetFabricResource(ctx context.Context, req *fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.FabricService, podName)
	fab := fabrics.Fabrics{
		Auth:          f.IsAuthorizedRPC,
		ContactClient: f.ContactClientRPC,
	}
	resp := &fabricsproto.FabricResponse{}
	data := fab.GetFabricResource(ctx, req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
	l.LogWithFields(ctx).Debugf("final response for get fabric resource request: %s", string(resp.Body))
	return resp, nil
}

// UpdateFabricResource defines  the operation which handles the RPC request response
// for creating/updating the specific fabric resource
func (f *Fabrics) UpdateFabricResource(ctx context.Context, req *fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.FabricService, podName)
	fab := fabrics.Fabrics{
		Auth:          f.IsAuthorizedRPC,
		ContactClient: f.ContactClientRPC,
	}
	resp := &fabricsproto.FabricResponse{}
	data := fab.UpdateFabricResource(ctx, req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.Body = generateResponse(ctx, data.Body)
	resp.StatusMessage = data.StatusMessage
	l.LogWithFields(ctx).Debugf("final response for update fabric resource request: %s", string(resp.Body))
	return resp, nil

}

// AddFabric defines  the operation which handles the RPC request response for Add fabric
func (f *Fabrics) AddFabric(ctx context.Context, req *fabricsproto.AddFabricRequest) (*fabricsproto.FabricResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.FabricService, podName)
	resp := &fabricsproto.FabricResponse{}
	data := fabrics.AddFabric(ctx, req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
	l.LogWithFields(ctx).Debugf("final response for add fabric request: %s", string(resp.Body))
	return resp, nil

}

// RemoveFabric defines  the operation which handles the RPC request response for Remove fabric
func (f *Fabrics) RemoveFabric(ctx context.Context, req *fabricsproto.AddFabricRequest) (*fabricsproto.FabricResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.FabricService, podName)
	resp := &fabricsproto.FabricResponse{}
	data := fabrics.RemoveFabric(ctx, req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
	l.LogWithFields(ctx).Debugf("final response for remove fabric request: %s", string(resp.Body))
	return resp, nil
}

// DeleteFabricResource defines the operation which handled the RPC request response
// for deleting the specified fabric resource information
func (f *Fabrics) DeleteFabricResource(ctx context.Context, req *fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.FabricService, podName)
	fab := fabrics.Fabrics{
		Auth:          f.IsAuthorizedRPC,
		ContactClient: f.ContactClientRPC,
	}
	resp := &fabricsproto.FabricResponse{}
	data := fab.DeleteFabricResource(ctx, req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
	l.LogWithFields(ctx).Debugf("final response for delete fabric resource request: %s", string(resp.Body))
	return resp, nil
}

func generateResponse(ctx context.Context, input interface{}) []byte {
	bytes, err := json.Marshal(input)
	if err != nil {
		l.LogWithFields(ctx).Error("error in unmarshalling response object from util-libs" + err.Error())
	}
	return bytes
}
