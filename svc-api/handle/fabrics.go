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

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	iris "github.com/kataras/iris/v12"
)

// FabricRPCs defines all the RPC methods in fabric service
type FabricRPCs struct {
	GetFabricResourceRPC    func(context.Context, fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error)
	UpdateFabricResourceRPC func(context.Context, fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error)
	DeleteFabricResourceRPC func(context.Context, fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error)
}

const (
	rpcFailedErrMsg = "RPC error: "
)

// GetFabricCollection defines the GetFabricCollection iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := getFabricRequest(ctx)
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting fabric collection with request url %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := f.GetFabricResourceRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	resp.Header = map[string]string{
		"Allow": `"GET"`,
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting all fabric collection is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendFabricResponse(ctx, resp)
}

// GetFabric defines the GetFabric iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabric(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := getFabricRequest(ctx)
	l.LogWithFields(ctxt).Debugf("Incoming request received for creating fabric with request url %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := f.GetFabricResourceRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	resp.Header = map[string]string{
		"Allow": `"GET"`,
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for creating fabric is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendFabricResponse(ctx, resp)
}

// GetFabricSwitchCollection defines the GetFabricSwitchCollection iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricSwitchCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := getFabricRequest(ctx)
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting fabric switch collection with request url %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := f.GetFabricResourceRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	resp.Header = map[string]string{
		"Allow": `"GET"`,
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting all fabric switch collection is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendFabricResponse(ctx, resp)
}

// GetFabricSwitch defines the GetFabricSwitch iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricSwitch(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := getFabricRequest(ctx)
	l.LogWithFields(ctxt).Debugf("Incoming request received for creating fabric switch with request url %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := f.GetFabricResourceRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	resp.Header = map[string]string{
		"Allow": `"GET"`,
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for creating fabric switch is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendFabricResponse(ctx, resp)
}

// GetSwitchPortCollection defines the GetSwitchPortCollection iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetSwitchPortCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := getFabricRequest(ctx)
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting switch port collection with request url %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := f.GetFabricResourceRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	resp.Header = map[string]string{
		"Allow": `"GET"`,
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting switch port collection is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendFabricResponse(ctx, resp)
}

// GetSwitchPort defines the GetSwitchPort iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetSwitchPort(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := getFabricRequest(ctx)
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting switch port with request url %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := f.GetFabricResourceRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	resp.Header = map[string]string{
		"Allow": `"GET", "PATCH"`,
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting switch port is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendFabricResponse(ctx, resp)
}

// GetFabricZoneCollection defines the GetFabricZoneCollection iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricZoneCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := getFabricRequest(ctx)
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting fabric zone collection with request url %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := f.GetFabricResourceRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	resp.Header = map[string]string{
		"Allow": `"GET", "POST"`,
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting fabric zone collection is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendFabricResponse(ctx, resp)
}

// GetFabricZone defines the GetFabricZone iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricZone(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := getFabricRequest(ctx)
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting fabric zone with request url %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := f.GetFabricResourceRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	resp.Header = map[string]string{
		"Allow": `"GET", "PUT", "PATCH", "DELETE"`,
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting fabric zone is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendFabricResponse(ctx, resp)
}

// GetFabricEndPointCollection defines the GetFabricEndPointCollection iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricEndPointCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := getFabricRequest(ctx)
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting fabric endpoint collection with request url %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := f.GetFabricResourceRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	resp.Header = map[string]string{
		"Allow": `"GET", "POST"`,
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting fabric endpoint collection is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendFabricResponse(ctx, resp)
}

// GetFabricEndPoints defines the GetFabricEndPoints iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricEndPoints(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := getFabricRequest(ctx)
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting fabric endpoint with request url %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := f.GetFabricResourceRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	resp.Header = map[string]string{
		"Allow": `"GET", "PUT", "PATCH", "DELETE"`,
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting fabric endpoint is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendFabricResponse(ctx, resp)
}

// GetFabricAddressPoolCollection defines the GetFabricAddressPoolCollection iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricAddressPoolCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := getFabricRequest(ctx)
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting fabric address pool collection with request url %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := f.GetFabricResourceRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	resp.Header = map[string]string{
		"Allow": `"GET",  "POST"`,
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting fabric address pool collection is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendFabricResponse(ctx, resp)
}

// GetFabricAddressPool defines the GetFabricAddressPool iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricAddressPool(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := getFabricRequest(ctx)
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting fabric address pool with request url %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := f.GetFabricResourceRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	resp.Header = map[string]string{
		"Allow": `"GET", "PUT", "PATCH", "DELETE"`,
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting fabric address pool is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendFabricResponse(ctx, resp)
}

// UpdateFabricResource defines the UpdateFabricResource iris handler.
// The method updates if Fabric Resource exists else creates new one.
func (f *FabricRPCs) UpdateFabricResource(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := fabricsproto.FabricRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
		Method:       ctx.Request().Method,
	}

	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	var createReq interface{}
	err := ctx.ReadJSON(&createReq)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the  request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}

	// marshalling the req to make fabric UpdateFabricResource request
	// Since fabric FabricRequest accepts []byte stream
	request, err := json.Marshal(createReq)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for updating fabric resources with request url %s and request body %s", req.URL, string(request))
	req.RequestBody = request
	resp, err := f.UpdateFabricResourceRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for updating fabric resource is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendFabricResponse(ctx, resp)
}

// DeleteFabricResource defines the DeleteFabricResource iris handler.
// This method is used for deleting requested fabric resource
func (f *FabricRPCs) DeleteFabricResource(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := getFabricRequest(ctx)
	l.LogWithFields(ctxt).Debugf("Incoming request received for deleting fabric resources with request url %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := f.DeleteFabricResourceRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for deleting fabric resource is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendFabricResponse(ctx, resp)
}

// sendFabricResponse writes the fabric response to client
func sendFabricResponse(ctx iris.Context, resp *fabricsproto.FabricResponse) {
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// getFabricRequest will extract the request from the context and return
func getFabricRequest(ctx iris.Context) fabricsproto.FabricRequest {
	return fabricsproto.FabricRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
	}
}
