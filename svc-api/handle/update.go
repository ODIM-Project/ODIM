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
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// UpdateRPCs used to define the service RPC function
type UpdateRPCs struct {
	GetUpdateServiceRPC               func(updateproto.UpdateRequest) (*updateproto.UpdateResponse, error)
	SimpleUpdateRPC                   func(updateproto.UpdateRequest) (*updateproto.UpdateResponse, error)
	StartUpdateRPC                    func(updateproto.UpdateRequest) (*updateproto.UpdateResponse, error)
	GetFirmwareInventoryRPC           func(updateproto.UpdateRequest) (*updateproto.UpdateResponse, error)
	GetFirmwareInventoryCollectionRPC func(updateproto.UpdateRequest) (*updateproto.UpdateResponse, error)
	GetSoftwareInventoryRPC           func(updateproto.UpdateRequest) (*updateproto.UpdateResponse, error)
	GetSoftwareInventoryCollectionRPC func(updateproto.UpdateRequest) (*updateproto.UpdateResponse, error)
}

// GetUpdateService is the handler for getting UpdateService details
func (a *UpdateRPCs) GetUpdateService(ctx iris.Context) {
	req := updateproto.UpdateRequest{
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
	resp, err := a.GetUpdateServiceRPC(req)
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

//GetFirmwareInventoryCollection is a handler for firmware inventory collection
func (a *UpdateRPCs) GetFirmwareInventoryCollection(ctx iris.Context) {
	req := updateproto.UpdateRequest{
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
	resp, err := a.GetFirmwareInventoryCollectionRPC(req)
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

// GetSoftwareInventoryCollection is a handler for software inventory collection
func (a *UpdateRPCs) GetSoftwareInventoryCollection(ctx iris.Context) {
	req := updateproto.UpdateRequest{
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
	resp, err := a.GetSoftwareInventoryCollectionRPC(req)
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

// GetFirmwareInventory is a handler for firmware inventory
func (a *UpdateRPCs) GetFirmwareInventory(ctx iris.Context) {
	req := updateproto.UpdateRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		ResourceID:   ctx.Params().Get("firmwareInventory_id"),
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
	resp, err := a.GetFirmwareInventoryRPC(req)
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

// GetSoftwareInventory is a handler for firmware inventory
func (a *UpdateRPCs) GetSoftwareInventory(ctx iris.Context) {
	req := updateproto.UpdateRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		ResourceID:   ctx.Params().Get("softwareInventory_id"),
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
	resp, err := a.GetSoftwareInventoryRPC(req)
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

//SimpleUpdate is a handler for simple update action
func (a *UpdateRPCs) SimpleUpdate(ctx iris.Context) {
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the simple update request body: " + err.Error()
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
	// Marshalling the req to make reset request
	request, err := json.Marshal(req)
	updateRequest := updateproto.UpdateRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
	}
	resp, err := a.SimpleUpdateRPC(updateRequest)
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

//StartUpdate is a handler for start update action
func (a *UpdateRPCs) StartUpdate(ctx iris.Context) {
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	updateRequest := updateproto.UpdateRequest{
		SessionToken: sessionToken,
	}
	resp, err := a.StartUpdateRPC(updateRequest)
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
