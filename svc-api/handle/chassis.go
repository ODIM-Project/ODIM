//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//(C) Copyright 2020 Intel Corporation
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
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/kataras/iris/v12"
)

// ChassisRPCs defines all the RPC methods in system service
type ChassisRPCs struct {
	GetChassisCollectionRPC func(ctx context.Context, req chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error)
	GetChassisResourceRPC   func(ctx context.Context, req chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error)
	GetChassisRPC           func(ctx context.Context, req chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error)
	CreateChassisRPC        func(ctx context.Context, req chassisproto.CreateChassisRequest) (*chassisproto.GetChassisResponse, error)
	DeleteChassisRPC        func(ctx context.Context, req chassisproto.DeleteChassisRequest) (*chassisproto.GetChassisResponse, error)
	UpdateChassisRPC        func(ctx context.Context, req chassisproto.UpdateChassisRequest) (*chassisproto.GetChassisResponse, error)
}

// CreateChassis creates a new chassis
func (chassis *ChassisRPCs) CreateChassis(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	requestBody := new(json.RawMessage)
	e := ctx.ReadJSON(requestBody)
	if e != nil {
		errorMessage := "error while trying to read obligatory json body: " + e.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for create chasses with request body %s", string(*requestBody))
	rpcResp, rpcErr := chassis.CreateChassisRPC(
		ctxt,
		chassisproto.CreateChassisRequest{
			RequestBody:  *requestBody,
			SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		},
	)

	if rpcErr != nil {
		l.LogWithFields(ctxt).Error("RPC error:" + rpcErr.Error())
		re := common.GeneralError(http.StatusInternalServerError, response.InternalError, rpcErr.Error(), nil, nil)
		writeResponse(ctx, re.Header, re.StatusCode, re.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for create chassis with response body %s and status code %d", string(rpcResp.Body), rpcResp.StatusCode)
	writeResponse(ctx, rpcResp.Header, rpcResp.StatusCode, rpcResp.Body)
}

func writeResponse(ctx iris.Context, headers map[string]string, status int32, body interface{}) {
	common.SetResponseHeader(ctx, headers)
	ctx.StatusCode(int(status))
	switch v := body.(type) {
	case []byte:
		ctx.Write(v)
	default:
		ctx.JSON(v)
	}
}

// GetChassisCollection fetches all Chassis
func (chassis *ChassisRPCs) GetChassisCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := chassisproto.GetChassisRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting chassis collection with request uri %s", string(req.URL))
	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := chassis.GetChassisCollectionRPC(ctxt, req)
	if err != nil {
		errorMessage := " RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for get chassis collections with response body %s and status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetChassisResource defines the GetChassisResource iris handler.
// The method extract the session token,uuid and request url and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (chassis *ChassisRPCs) GetChassisResource(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := chassisproto.GetChassisRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		RequestParam: ctx.Params().Get("id"),
		ResourceID:   ctx.Params().Get("rid"),
		URL:          ctx.Request().RequestURI}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting chassis resources with request uri %s with request params %s and resource id %s", req.URL, req.RequestParam, req.ResourceID)
	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := chassis.GetChassisResourceRPC(ctxt, req)
	if err != nil {
		errorMessage := " RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for get chassis resources with response body %s and status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetChassis defines the GetChassisResource iris handler.
// The method extract the session token,uuid and request url and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (chassis *ChassisRPCs) GetChassis(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := chassisproto.GetChassisRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		RequestParam: ctx.Params().Get("id"),
		URL:          ctx.Request().RequestURI}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting chassis with request uri %s and request params %s", req.URL, req.RequestParam)
	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := chassis.GetChassisRPC(ctxt, req)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for get chassis with is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH, DELETE")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// UpdateChassis updates an existing chassis
func (chassis *ChassisRPCs) UpdateChassis(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	requestBody := new(json.RawMessage)
	e := ctx.ReadJSON(requestBody)
	if e != nil {
		errorMessage := "error while trying to read obligatory json body: " + e.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for updating chassis with request body %s", string(*requestBody))
	rr, rerr := chassis.UpdateChassisRPC(ctxt, chassisproto.UpdateChassisRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
		RequestBody:  *requestBody,
	})

	if rerr != nil {
		l.LogWithFields(ctxt).Error("RPC error:" + rerr.Error())
		re := common.GeneralError(http.StatusInternalServerError, response.InternalError, rerr.Error(), nil, nil)
		writeResponse(ctx, re.Header, re.StatusCode, re.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for updating chassis is %s with status code %d", string(rr.Body), int(rr.StatusCode))
	writeResponse(ctx, rr.Header, rr.StatusCode, rr.Body)

}

// DeleteChassis deletes a chassis
func (chassis *ChassisRPCs) DeleteChassis(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	l.LogWithFields(ctxt).Debugf("Incoming request received for deleting chassis with request uri %s", ctx.Request().RequestURI)
	rpcResp, rpcErr := chassis.DeleteChassisRPC(ctxt, chassisproto.DeleteChassisRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
	})

	if rpcErr != nil {
		l.LogWithFields(ctxt).Error("RPC error:" + rpcErr.Error())
		re := common.GeneralError(http.StatusInternalServerError, response.InternalError, rpcErr.Error(), nil, nil)
		writeResponse(ctx, re.Header, re.StatusCode, re.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for deleting is %s with status code %d", string(rpcResp.Body), int(rpcResp.StatusCode))
	writeResponse(ctx, rpcResp.Header, rpcResp.StatusCode, rpcResp.Body)
}
