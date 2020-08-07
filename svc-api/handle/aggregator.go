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

//Package handle ...
package handle

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// AggregatorRPCs defines all the RPC methods in aggregator service
type AggregatorRPCs struct {
	GetAggregationServiceRPC   func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	AddComputeRPC              func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	DeleteComputeRPC           func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	ResetRPC                   func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	SetDefaultBootOrderRPC     func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	AddAggregationSourceRPC    func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	GetAllAggregationSourceRPC func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	GetAggregationSourceRPC    func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	UpdateAggregationSourceRPC func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	DeleteAggregationSourceRPC func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
}

// GetAggregationService is the handler for getting AggregationService details
func (a *AggregatorRPCs) GetAggregationService(ctx iris.Context) {
	req := aggregatorproto.AggregatorRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp, err := a.GetAggregationServiceRPC(req)
	if err != nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}

// AddCompute is the handler for adding system
func (a *AggregatorRPCs) AddCompute(ctx iris.Context) {
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the aggregator request body: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusBadRequest) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	sessionToken := ctx.Request().Header.Get("X-Auth-Token")

	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	// marshalling the req to make aggregator add request
	// Since aggregator add request accepts []byte stream
	request, err := json.Marshal(req)

	addRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
	}
	resp, err := a.AddComputeRPC(addRequest)
	if err != nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// DeleteCompute is the handler for deleting compute system
// from iris context will get the request and check sessiontoken
// and do rpc call and send response back
func (a *AggregatorRPCs) DeleteCompute(ctx iris.Context) {
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the aggregator request body: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusBadRequest) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	sessionToken := ctx.Request().Header.Get("X-Auth-Token")

	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	// marshalling the req to make aggregator delete request
	// Since aggregator deleterequest accepts []byte stream
	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	deleteRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
	}
	resp, err := a.DeleteComputeRPC(deleteRequest)
	if err != nil {
		errorMessage := "error: RPC error: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// Reset is the handler to reset compute system
// from iris context will get the request and check sessiontoken
// and do rpc call and send response back
func (a *AggregatorRPCs) Reset(ctx iris.Context) {
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the  request body: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusBadRequest) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	sessionToken := ctx.Request().Header.Get("X-Auth-Token")

	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	// marshalling the req to make aggregator delete request
	// Since aggregator deleterequest accepts []byte stream
	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	resetRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
	}
	resp, err := a.ResetRPC(resetRequest)
	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// SetDefaultBootOrder is the handler to set default boot order
// from iris context will get the request and check sessiontoken
// and do rpc call and send response back
func (a *AggregatorRPCs) SetDefaultBootOrder(ctx iris.Context) {
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the  request body: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusBadRequest) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	sessionToken := ctx.Request().Header.Get("X-Auth-Token")

	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	// marshalling the req to make aggregator SetDefaultBootOrder request
	// Since aggregator SetDefaultBootOrderRequest accepts []byte stream
	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	resetRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
	}
	resp, err := a.SetDefaultBootOrderRPC(resetRequest)
	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// AddAggregationSource is the handler for adding  AggregationSource details
func (a *AggregatorRPCs) AddAggregationSource(ctx iris.Context) {
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the aggregator request body: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusBadRequest) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	sessionToken := ctx.Request().Header.Get("X-Auth-Token")

	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	// marshalling the req to make aggregator add request
	// Since aggregator add request accepts []byte stream
	request, err := json.Marshal(req)

	addRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
	}
	resp, err := a.AddAggregationSourceRPC(addRequest)
	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetAllAggregationSource is the handler for getting all  AggregationSource details
func (a *AggregatorRPCs) GetAllAggregationSource(ctx iris.Context) {
	req := aggregatorproto.AggregatorRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp, err := a.GetAllAggregationSourceRPC(req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetAggregationSource is the handler for getting  AggregationSource details
func (a *AggregatorRPCs) GetAggregationSource(ctx iris.Context) {
	req := aggregatorproto.AggregatorRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp, err := a.GetAggregationSourceRPC(req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// UpdateAggregationSource is the handler for updating  AggregationSource details
func (a *AggregatorRPCs) UpdateAggregationSource(ctx iris.Context) {
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the aggregator request body: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusBadRequest) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	sessionToken := ctx.Request().Header.Get("X-Auth-Token")

	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	// marshalling the req to make aggregator add request
	// Since aggregator add request accepts []byte stream
	request, err := json.Marshal(req)

	updateRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
	}
	resp, err := a.UpdateAggregationSourceRPC(updateRequest)
	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}

// DeleteAggregationSource is the handler for updating  AggregationSource details
func (a *AggregatorRPCs) DeleteAggregationSource(ctx iris.Context) {
	//need to be changed after the code is added
	ctx.StatusCode(http.StatusNotImplemented)
}
