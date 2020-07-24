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
	"log"
	"net/http"

	"github.com/bharath-b-hpe/odimra/lib-utilities/common"
	chassisproto "github.com/bharath-b-hpe/odimra/lib-utilities/proto/chassis"
	"github.com/bharath-b-hpe/odimra/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// ChassisRPCs defines all the RPC methods in system service
type ChassisRPCs struct {
	GetChassisCollectionRPC func(req chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error)
	GetChassisResourceRPC   func(req chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error)
	GetChassisRPC           func(req chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error)
}

//GetChassisCollection fetches all Chassis
func (chassis *ChassisRPCs) GetChassisCollection(ctx iris.Context) {
	req := chassisproto.GetChassisRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp, err := chassis.GetChassisCollectionRPC(req)
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

// GetChassisResource defines the GetChassisResource iris handler.
// The method extract the session token,uuid and request url and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (chassis *ChassisRPCs) GetChassisResource(ctx iris.Context) {
	req := chassisproto.GetChassisRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		RequestParam: ctx.Params().Get("id"),
		URL:          ctx.Request().RequestURI}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp, err := chassis.GetChassisResourceRPC(req)
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
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp, err := chassis.GetChassisRPC(req)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
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
