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

//Package handle ...
package handle

import (
	"encoding/json"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"

	"github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
)

// ChassisRPCs defines all the RPC methods in system service
type ChassisRPCs struct {
	GetChassisCollectionRPC func(req chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error)
	GetChassisResourceRPC   func(req chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error)
	GetChassisRPC           func(req chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error)
	CreateChassisRPC        func(req chassisproto.CreateChassisRequest) (*chassisproto.GetChassisResponse, error)
	DeleteChassisRPC        func(req chassisproto.DeleteChassisRequest) (*chassisproto.GetChassisResponse, error)
	UpdateChassisRPC        func(req chassisproto.UpdateChassisRequest) (*chassisproto.GetChassisResponse, error)
}

func (chassis *ChassisRPCs) CreateChassis(ctx iris.Context) {
	requestBody := new(json.RawMessage)
	e := ctx.ReadJSON(requestBody)
	if e != nil {
		errorMessage := "error while trying to read obligatory json body: " + e.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusBadRequest)
		common.SetResponseHeader(ctx, response.Header)
		ctx.JSON(&response.Body)
		return
	}

	rpcResp, rpcErr := chassis.CreateChassisRPC(
		chassisproto.CreateChassisRequest{
			RequestBody:  *requestBody,
			SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		},
	)

	if rpcErr != nil {
		log.Error("RPC error:" + rpcErr.Error())
		re := common.GeneralError(http.StatusInternalServerError, response.InternalError, rpcErr.Error(), nil, nil)
		writeResponse(ctx, re.Header, re.StatusCode, re.Body)
		return
	}

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

//GetChassisCollection fetches all Chassis
func (chassis *ChassisRPCs) GetChassisCollection(ctx iris.Context) {
	req := chassisproto.GetChassisRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI}
	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp, err := chassis.GetChassisCollectionRPC(req)
	if err != nil {
		errorMessage := " RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetChassisResource defines the GetChassisResource iris handler.
// The method extract the session token,uuid and request url and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (chassis *ChassisRPCs) GetChassisResource(ctx iris.Context) {
	req := chassisproto.GetChassisRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		RequestParam: ctx.Params().Get("id"),
		ResourceID:   ctx.Params().Get("rid"),
		URL:          ctx.Request().RequestURI}
	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp, err := chassis.GetChassisResourceRPC(req)
	if err != nil {
		errorMessage := " RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetChassis defines the GetChassisResource iris handler.
// The method extract the session token,uuid and request url and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (chassis *ChassisRPCs) GetChassis(ctx iris.Context) {
	req := chassisproto.GetChassisRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		RequestParam: ctx.Params().Get("id"),
		URL:          ctx.Request().RequestURI}
	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp, err := chassis.GetChassisRPC(req)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

func (chassis *ChassisRPCs) UpdateChassis(ctx iris.Context) {
	requestBody := new(json.RawMessage)
	e := ctx.ReadJSON(requestBody)
	if e != nil {
		errorMessage := "error while trying to read obligatory json body: " + e.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusBadRequest)
		common.SetResponseHeader(ctx, response.Header)
		ctx.JSON(&response.Body)
		return
	}
	rr, rerr := chassis.UpdateChassisRPC(chassisproto.UpdateChassisRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
		RequestBody:  *requestBody,
	})

	if rerr != nil {
		log.Println("RPC error:" + rerr.Error())
		re := common.GeneralError(http.StatusInternalServerError, response.InternalError, rerr.Error(), nil, nil)
		writeResponse(ctx, re.Header, re.StatusCode, re.Body)
		return
	}

	writeResponse(ctx, rr.Header, rr.StatusCode, rr.Body)

}

func (chassis *ChassisRPCs) DeleteChassis(ctx iris.Context) {
	rpcResp, rpcErr := chassis.DeleteChassisRPC(chassisproto.DeleteChassisRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	})

	if rpcErr != nil {
		log.Println("RPC error:" + rpcErr.Error())
		re := common.GeneralError(http.StatusInternalServerError, response.InternalError, rpcErr.Error(), nil, nil)
		writeResponse(ctx, re.Header, re.StatusCode, re.Body)
		return
	}

	writeResponse(ctx, rpcResp.Header, rpcResp.StatusCode, rpcResp.Body)
}
