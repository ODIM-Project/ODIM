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

// Package handle ...
package handle

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// AggregatorRPCs defines all the RPC methods in aggregator service
type AggregatorRPCs struct {
	GetAggregationServiceRPC                func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	ResetRPC                                func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	SetDefaultBootOrderRPC                  func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	AddAggregationSourceRPC                 func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	GetAllAggregationSourceRPC              func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	GetAggregationSourceRPC                 func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	UpdateAggregationSourceRPC              func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	DeleteAggregationSourceRPC              func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	CreateAggregateRPC                      func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	GetAggregateCollectionRPC               func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	GetAggregateRPC                         func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	DeleteAggregateRPC                      func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	AddElementsToAggregateRPC               func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	RemoveElementsFromAggregateRPC          func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	ResetAggregateElementsRPC               func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	SetDefaultBootOrderAggregateElementsRPC func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	GetAllConnectionMethodsRPC              func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	GetConnectionMethodRPC                  func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	GetResetActionInfoServiceRPC            func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
	GetSetDefaultBootOrderActionInfoRPC     func(context.Context, aggregatorproto.AggregatorRequest) (*aggregatorproto.AggregatorResponse, error)
}

const (
	// AuthTokenHeader holds the key of X-Auth-Token header
	AuthTokenHeader        = "X-Auth-Token"
	rpcCallFailedErrMsg    = "something went wrong with the RPC calls: "
	invalidAuthTokenErrMsg = "no X-Auth-Token found in request header"
)

// GetAggregationService is the handler for getting AggregationService details
func (a *AggregatorRPCs) GetAggregationService(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	l.LogWithFields(ctxt).Debugf("Incoming request received for get aggregation service")
	req := aggregatorproto.AggregatorRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
	}
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetAggregationServiceRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET")
	sendAggregatorResponse(ctx, resp)
	l.LogWithFields(ctxt).Debugf("outgoing response for get aggregation service is %s with response code %d", string(resp.Body), int(resp.StatusCode))

}

// Reset is the handler to reset compute system
// from iris context will get the request and check sessiontoken
// and do rpc call and send response back
func (a *AggregatorRPCs) Reset(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	l.LogWithFields(ctxt).Debugf("Incoming request received for Resetting compute system")
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}

	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)
	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	// marshalling the req to make aggregator delete request
	// Since aggregator deleterequest accepts []byte stream
	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}

	resetRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
	}
	resp, err := a.ResetRPC(ctxt, resetRequest)
	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}

	sendAggregatorResponse(ctx, resp)
	l.LogWithFields(ctxt).Debugf("outgoing response for resetting compute system is %s and response code %d", string(resp.Body), int(resp.StatusCode))
}

// SetDefaultBootOrder is the handler to set default boot order
// from iris context will get the request and check sessiontoken
// and do rpc call and send response back
func (a *AggregatorRPCs) SetDefaultBootOrder(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	l.LogWithFields(ctxt).Debugf("Incoming request received for setting default boot order")
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}
	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)

	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	// marshalling the req to make aggregator SetDefaultBootOrder request
	// Since aggregator SetDefaultBootOrderRequest accepts []byte stream
	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}

	resetRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
	}
	resp, err := a.SetDefaultBootOrderRPC(ctxt, resetRequest)
	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}

	sendAggregatorResponse(ctx, resp)
}

// AddAggregationSource is the handler for adding  AggregationSource details
func (a *AggregatorRPCs) AddAggregationSource(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req interface{}
	l.LogWithFields(ctxt).Debugf("Incoming request received for Adding aggregationsource")
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the aggregator request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}

	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)

	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	// marshalling the req to make aggregator add request
	// Since aggregator add request accepts []byte stream
	request, err := json.Marshal(req)

	addRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
	}
	resp, err := a.AddAggregationSourceRPC(ctxt, addRequest)
	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}

	sendAggregatorResponse(ctx, resp)
}

// GetAllAggregationSource is the handler for getting all  AggregationSource details
func (a *AggregatorRPCs) GetAllAggregationSource(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting all aggregationsources")
	req := aggregatorproto.AggregatorRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
	}
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetAllAggregationSourceRPC(ctxt, req)
	if err != nil {
		errorMessage := " RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	sendAggregatorResponse(ctx, resp)
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting all aggregation sources is %s with response code %d", string(resp.Body), int(resp.StatusCode))
}

// GetAggregationSource is the handler for getting  AggregationSource details
func (a *AggregatorRPCs) GetAggregationSource(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := aggregatorproto.AggregatorRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting aggregationsource with uri %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetAggregationSourceRPC(ctxt, req)
	if err != nil {
		errorMessage := " RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH, DELETE")
	sendAggregatorResponse(ctx, resp)
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting aggregation source is %s with response code %d", string(resp.Body), int(resp.StatusCode))
}

// UpdateAggregationSource is the handler for updating  AggregationSource details
func (a *AggregatorRPCs) UpdateAggregationSource(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the aggregator request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}

	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)

	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	// marshalling the req to make aggregator add request
	// Since aggregator add request accepts []byte stream
	request, err := json.Marshal(req)

	updateRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for updating aggregationsource with uri %s", updateRequest.URL)
	resp, err := a.UpdateAggregationSourceRPC(ctxt, updateRequest)
	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}

	sendAggregatorResponse(ctx, resp)
	l.LogWithFields(ctxt).Debugf("Outgoing response for updating aggregation source is %s with response code %d", string(resp.Body), int(resp.StatusCode))

}

// DeleteAggregationSource is the handler for updating  AggregationSource details
func (a *AggregatorRPCs) DeleteAggregationSource(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := aggregatorproto.AggregatorRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for deleting aggregationsource with uri %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.DeleteAggregationSourceRPC(ctxt, req)
	if err != nil {
		errorMessage := " RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for deleting aggregation source is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	sendAggregatorResponse(ctx, resp)
}

// CreateAggregate is the handler for creating an aggregate
func (a *AggregatorRPCs) CreateAggregate(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the aggregator request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}

	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)

	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	// marshalling the req to make aggregator create request
	request, err := json.Marshal(req)
	l.LogWithFields(ctxt).Debugf("Incoming request received for creating aggregate with request body %s", string(request))
	createRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
	}
	resp, err := a.CreateAggregateRPC(ctxt, createRequest)
	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for creating aggregate is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	sendAggregatorResponse(ctx, resp)

}

// GetAggregateCollection is the handler for getting collection of aggregates
func (a *AggregatorRPCs) GetAggregateCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	l.LogWithFields(ctxt).Debug("Incoming request received for getting all aggregate collections")
	req := aggregatorproto.AggregatorRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
	}
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetAggregateCollectionRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting all aggregate collections is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	sendAggregatorResponse(ctx, resp)

}

// GetAggregate is the handler for getting an aggregate
func (a *AggregatorRPCs) GetAggregate(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := aggregatorproto.AggregatorRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting aggregate with request uri %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetAggregateRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting aggregate is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET, DELETE")
	sendAggregatorResponse(ctx, resp)
}

// DeleteAggregate is the handler for deleting an aggregate
func (a *AggregatorRPCs) DeleteAggregate(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := aggregatorproto.AggregatorRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for deleting aggregate with request uri %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.DeleteAggregateRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for deleting aggregate is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	sendAggregatorResponse(ctx, resp)
}

// AddElementsToAggregate is the handler for adding elements to an aggregate
func (a *AggregatorRPCs) AddElementsToAggregate(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the aggregator request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}

	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)

	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	// marshalling the req to make aggregator add request
	request, _ := json.Marshal(req)

	addRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		URL:          ctx.Request().RequestURI,
		RequestBody:  request,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for adding element with uri %s to the aggregation collections with request body %s", addRequest.URL, string(addRequest.RequestBody))

	resp, err := a.AddElementsToAggregateRPC(ctxt, addRequest)
	if err != nil {
		errorMessage := rpcCallFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for adding an element to an aggregate is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	sendAggregatorResponse(ctx, resp)
}

// RemoveElementsFromAggregate is the handler for removing elements from an aggregate
func (a *AggregatorRPCs) RemoveElementsFromAggregate(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the aggregator request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}

	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)

	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	// marshalling the req to make aggregator remove request
	request, _ := json.Marshal(req)

	removeRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		URL:          ctx.Request().RequestURI,
		RequestBody:  request,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for deleting element with uri %s from the aggregation collections with request body %s", removeRequest.URL, string(removeRequest.RequestBody))
	resp, err := a.RemoveElementsFromAggregateRPC(ctxt, removeRequest)
	if err != nil {
		errorMessage := rpcCallFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for removing an element from an aggregate is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	sendAggregatorResponse(ctx, resp)
}

// ResetAggregateElements is the handler for resetting elements of an aggregate
func (a *AggregatorRPCs) ResetAggregateElements(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the aggregator request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}

	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)

	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	// marshalling the req to make aggregator reset elements request
	request, _ := json.Marshal(req)

	resetRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		URL:          ctx.Request().RequestURI,
		RequestBody:  request,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for resetting aggregate elements with uri %s with request body %s", resetRequest.URL, string(resetRequest.RequestBody))
	resp, err := a.ResetAggregateElementsRPC(ctxt, resetRequest)
	if err != nil {
		errorMessage := rpcCallFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for resetting aggregate elements is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	sendAggregatorResponse(ctx, resp)
}

// SetDefaultBootOrderAggregateElements is the handler for SetDefaultBootOrder elements of an aggregate
func (a *AggregatorRPCs) SetDefaultBootOrderAggregateElements(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)
	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	bootOrderRequest := aggregatorproto.AggregatorRequest{
		SessionToken: sessionToken,
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for setting default boot order for all the element of an aggregate with uri %s and request body %s", bootOrderRequest.URL, string(bootOrderRequest.RequestBody))
	resp, err := a.SetDefaultBootOrderAggregateElementsRPC(ctxt, bootOrderRequest)
	if err != nil {
		errorMessage := rpcCallFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for setting default boot order for all the elements of an aggregate is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	sendAggregatorResponse(ctx, resp)
}

// GetAllConnectionMethods is the handler for get all connection methods
func (a *AggregatorRPCs) GetAllConnectionMethods(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting all the connection methods")
	req := aggregatorproto.AggregatorRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
	}
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := a.GetAllConnectionMethodsRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting all the connection methods is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	sendAggregatorResponse(ctx, resp)
}

// GetConnectionMethod is the handler for get connection method
func (a *AggregatorRPCs) GetConnectionMethod(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := aggregatorproto.AggregatorRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting the connection methods with the request URL %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetConnectionMethodRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting a collection method is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	sendAggregatorResponse(ctx, resp)
}

// GetResetActionInfoService is the handler for getting GetResetActionInfoService details
func (a *AggregatorRPCs) GetResetActionInfoService(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting reset action info")
	req := aggregatorproto.AggregatorRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
	}
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetResetActionInfoServiceRPC(ctxt, req)
	if err != nil {
		errorMessage := "RPC call error: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting a reset action information is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	sendAggregatorResponse(ctx, resp)

}

// GetSetDefaultBootOrderActionInfo is the handler for getting GetSetDefaultBootOrderActionInfo details
func (a *AggregatorRPCs) GetSetDefaultBootOrderActionInfo(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting the action info for set default boot order")
	req := aggregatorproto.AggregatorRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
	}
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetSetDefaultBootOrderActionInfoRPC(ctxt, req)
	if err != nil {
		errorMessage := "RPC call error: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting a set default boot order action information is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	sendAggregatorResponse(ctx, resp)

}

// sendSystemsResponse writes the aggregator response to client
func sendAggregatorResponse(ctx iris.Context, resp *aggregatorproto.AggregatorResponse) {
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}
