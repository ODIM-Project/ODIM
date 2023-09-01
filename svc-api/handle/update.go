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
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// UpdateRPCs used to define the service RPC function
type UpdateRPCs struct {
	GetUpdateServiceRPC               func(context.Context, updateproto.UpdateRequest) (*updateproto.UpdateResponse, error)
	SimpleUpdateRPC                   func(context.Context, updateproto.UpdateRequest) (*updateproto.UpdateResponse, error)
	StartUpdateRPC                    func(context.Context, updateproto.UpdateRequest) (*updateproto.UpdateResponse, error)
	GetFirmwareInventoryRPC           func(context.Context, updateproto.UpdateRequest) (*updateproto.UpdateResponse, error)
	GetFirmwareInventoryCollectionRPC func(context.Context, updateproto.UpdateRequest) (*updateproto.UpdateResponse, error)
	GetSoftwareInventoryRPC           func(context.Context, updateproto.UpdateRequest) (*updateproto.UpdateResponse, error)
	GetSoftwareInventoryCollectionRPC func(context.Context, updateproto.UpdateRequest) (*updateproto.UpdateResponse, error)
}

// GetUpdateService is the handler for getting UpdateService details
func (a *UpdateRPCs) GetUpdateService(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := updateproto.UpdateRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
	}
	l.LogWithFields(ctxt).Debug("Incoming request received for getting update service details")
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetUpdateServiceRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrorMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting update service is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	sendUpdateResponse(ctx, resp)
}

// GetFirmwareInventoryCollection is a handler for firmware inventory collection
func (a *UpdateRPCs) GetFirmwareInventoryCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := updateproto.UpdateRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
	}
	l.LogWithFields(ctxt).Debug("Incoming request received for getting firmware inventory collection")
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetFirmwareInventoryCollectionRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrorMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting firmware inventory collection is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	sendUpdateResponse(ctx, resp)

}

// GetSoftwareInventoryCollection is a handler for software inventory collection
func (a *UpdateRPCs) GetSoftwareInventoryCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := updateproto.UpdateRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
	}
	l.LogWithFields(ctxt).Debug("Incoming request received for software inventory collection")
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetSoftwareInventoryCollectionRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrorMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting software inventory collection is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	sendUpdateResponse(ctx, resp)

}

// GetFirmwareInventory is a handler for firmware inventory
func (a *UpdateRPCs) GetFirmwareInventory(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := updateproto.UpdateRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		ResourceID:   ctx.Params().Get("firmwareInventory_id"),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting firmware inventory with url %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetFirmwareInventoryRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrorMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting firmware inventory details is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	sendUpdateResponse(ctx, resp)

}

// GetSoftwareInventory is a handler for firmware inventory
func (a *UpdateRPCs) GetSoftwareInventory(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := updateproto.UpdateRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		ResourceID:   ctx.Params().Get("softwareInventory_id"),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting software inventory with url %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetSoftwareInventoryRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrorMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting software inventory details is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	sendUpdateResponse(ctx, resp)

}

// SimpleUpdate is a handler for simple update action
func (a *UpdateRPCs) SimpleUpdate(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the simple update request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}
	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)
	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	// Marshalling the req to make reset request
	request, err := json.Marshal(req)
	l.LogWithFields(ctxt).Debugf("Incoming request received for performing simple update with request body %s", string(request))
	updateRequest := updateproto.UpdateRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
	}
	errResp := validateSimpleUpdateRequest(ctxt, updateRequest.RequestBody)
	if errResp.StatusCode != http.StatusOK {
		common.SetResponseHeader(ctx, errResp.Header)
		ctx.StatusCode(int(errResp.StatusCode))
		ctx.JSON(&errResp.Body)
		return
	}
	resp, err := a.SimpleUpdateRPC(ctxt, updateRequest)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for simple update action is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendUpdateResponse(ctx, resp)
}

func validateSimpleUpdateRequest(ctx context.Context, requestBody []byte) response.RPC {
	var request map[string]interface{}
	err := json.Unmarshal(requestBody, &request)
	if err != nil {
		errMsg := "Unable to parse the simple update request" + err.Error()
		l.LogWithFields(ctx).Warn(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errMsg, nil, nil)
	}
	if request["Targets"] != nil {
		if reflect.TypeOf(request["Targets"]).Kind() != reflect.Slice {
			errMsg := "'Targets' parameter should be of type string array"
			l.LogWithFields(ctx).Warn(errMsg)
			return common.GeneralError(http.StatusBadRequest, response.PropertyValueTypeError, errMsg, []interface{}{"", "Targets"}, nil)
		}
		target := request["Targets"].([]interface{})
		for _, k := range target {
			if reflect.TypeOf(k).Kind() != reflect.String {
				errMsg := "'Targets' parameter should be of type string array"
				l.LogWithFields(ctx).Warn(errMsg)
				return common.GeneralError(http.StatusBadRequest, response.PropertyValueTypeError, errMsg, []interface{}{fmt.Sprintf("%v", k), "Targets"}, nil)
			}
		}
	}
	if request["ImageURI"] == nil {
		errMsg := "'ImageURI' parameter cannot be empty"
		l.LogWithFields(ctx).Warn(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"ImageURI"}, nil)
	}
	if reflect.TypeOf(request["ImageURI"]).Kind() != reflect.String {
		errMsg := "'ImageURI' parameter should be of type string"
		l.LogWithFields(ctx).Warn(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyValueTypeError, errMsg, []interface{}{"", "ImageURI"}, nil)
	}
	if request["ImageURI"] != nil {
		URI := request["ImageURI"]
		_, err = url.ParseRequestURI(URI.(string))
		if err != nil {
			errMsg := "Provided ImageURI is Invalid"
			l.LogWithFields(ctx).Warn(errMsg)
			return common.GeneralError(http.StatusBadRequest, response.PropertyValueTypeError, errMsg, []interface{}{fmt.Sprintf("%v", err), "ImageURI"}, nil)
		}
	}
	return response.RPC{StatusCode: http.StatusOK}
}

// StartUpdate is a handler for start update action
func (a *UpdateRPCs) StartUpdate(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	l.LogWithFields(ctxt).Debug("Incoming request received for start update action")
	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)
	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	updateRequest := updateproto.UpdateRequest{
		SessionToken: sessionToken,
	}
	resp, err := a.StartUpdateRPC(ctxt, updateRequest)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for start update action is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendUpdateResponse(ctx, resp)
}

// sendUpdateResponse writes the update response to client
func sendUpdateResponse(ctx iris.Context, resp *updateproto.UpdateResponse) {
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}
