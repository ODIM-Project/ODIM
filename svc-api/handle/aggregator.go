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
	GetAggregationServiceRPC                func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	ResetRPC                                func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	SetDefaultBootOrderRPC                  func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	AddAggregationSourceRPC                 func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	GetAllAggregationSourceRPC              func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	GetAggregationSourceRPC                 func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	UpdateAggregationSourceRPC              func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	DeleteAggregationSourceRPC              func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	CreateAggregateRPC                      func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	GetAggregateCollectionRPC               func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	GetAggregateRPC                         func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	DeleteAggregateRPC                      func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	AddElementsToAggregateRPC               func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	RemoveElementsFromAggregateRPC          func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	ResetAggregateElementsRPC               func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	SetDefaultBootOrderAggregateElementsRPC func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	GetAllConnectionMethodsRPC              func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	GetConnectionMethodRPC                  func(aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
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
		URL:          ctx.Request().RequestURI,
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
	resp, err := a.DeleteAggregationSourceRPC(req)
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

// CreateAggregate is the handler for creating an aggregate
func (a *AggregatorRPCs) CreateAggregate(ctx iris.Context) {
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

	// marshalling the req to make aggregator create request
	request, err := json.Marshal(req)

	createRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
	}
	resp, err := a.CreateAggregateRPC(createRequest)
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

// GetAggregateCollection is the handler for getting collection of aggregates
func (a *AggregatorRPCs) GetAggregateCollection(ctx iris.Context) {
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
	resp, err := a.GetAggregateCollectionRPC(req)
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

// GetAggregate is the handler for getting an aggregate
func (a *AggregatorRPCs) GetAggregate(ctx iris.Context) {
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
	resp, err := a.GetAggregateRPC(req)
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

// DeleteAggregate is the handler for deleting an aggregate
func (a *AggregatorRPCs) DeleteAggregate(ctx iris.Context) {
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
	resp, err := a.DeleteAggregateRPC(req)
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

// AddElementsToAggregate is the handler for adding elements to an aggregate
func (a *AggregatorRPCs) AddElementsToAggregate(ctx iris.Context) {
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
	request, _ := json.Marshal(req)

	addRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		URL:          ctx.Request().RequestURI,
		RequestBody:  request,
	}

	resp, err := a.AddElementsToAggregateRPC(addRequest)
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

// RemoveElementsFromAggregate is the handler for removing elements from an aggregate
func (a *AggregatorRPCs) RemoveElementsFromAggregate(ctx iris.Context) {
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

	// marshalling the req to make aggregator remove request
	request, _ := json.Marshal(req)

	removeRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		URL:          ctx.Request().RequestURI,
		RequestBody:  request,
	}

	resp, err := a.RemoveElementsFromAggregateRPC(removeRequest)
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

// ResetAggregateElements is the handler for resetting elements of an aggregate
func (a *AggregatorRPCs) ResetAggregateElements(ctx iris.Context) {
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

	// marshalling the req to make aggregator reset elements request
	request, _ := json.Marshal(req)

	resetRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		URL:          ctx.Request().RequestURI,
		RequestBody:  request,
	}

	resp, err := a.ResetAggregateElementsRPC(resetRequest)
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

// SetDefaultBootOrderAggregateElements is the handler for SetDefaultBootOrder elements of an aggregate
func (a *AggregatorRPCs) SetDefaultBootOrderAggregateElements(ctx iris.Context) {

	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	bootOrderRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		URL:          ctx.Request().RequestURI,
	}

	resp, err := a.SetDefaultBootOrderAggregateElementsRPC(bootOrderRequest)
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

// GetAllConnectionMethods is the handler for get all connection methods
func (a *AggregatorRPCs) GetAllConnectionMethods(ctx iris.Context) {
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
	resp, err := a.GetAllConnectionMethodsRPC(req)
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

// GetConnectionMethod is the handler for get connection method
func (a *AggregatorRPCs) GetConnectionMethod(ctx iris.Context) {
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
	resp, err := a.GetConnectionMethodRPC(req)
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
