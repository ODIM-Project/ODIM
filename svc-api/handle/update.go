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
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	log "github.com/sirupsen/logrus"

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
	defer ctx.Next()
	req := updateproto.UpdateRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := a.GetUpdateServiceRPC(req)
	if err != nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}

//GetFirmwareInventoryCollection is a handler for firmware inventory collection
func (a *UpdateRPCs) GetFirmwareInventoryCollection(ctx iris.Context) {
	defer ctx.Next()
	req := updateproto.UpdateRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := a.GetFirmwareInventoryCollectionRPC(req)
	if err != nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}

// GetSoftwareInventoryCollection is a handler for software inventory collection
func (a *UpdateRPCs) GetSoftwareInventoryCollection(ctx iris.Context) {
	defer ctx.Next()
	req := updateproto.UpdateRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := a.GetSoftwareInventoryCollectionRPC(req)
	if err != nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}

// GetFirmwareInventory is a handler for firmware inventory
func (a *UpdateRPCs) GetFirmwareInventory(ctx iris.Context) {
	defer ctx.Next()
	req := updateproto.UpdateRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		ResourceID:   ctx.Params().Get("firmwareInventory_id"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := a.GetFirmwareInventoryRPC(req)
	if err != nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}

// GetSoftwareInventory is a handler for firmware inventory
func (a *UpdateRPCs) GetSoftwareInventory(ctx iris.Context) {
	defer ctx.Next()
	req := updateproto.UpdateRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		ResourceID:   ctx.Params().Get("softwareInventory_id"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := a.GetSoftwareInventoryRPC(req)
	if err != nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}

//SimpleUpdate is a handler for simple update action
func (a *UpdateRPCs) SimpleUpdate(ctx iris.Context) {
	defer ctx.Next()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the simple update request body: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	// Marshalling the req to make reset request
	request, err := json.Marshal(req)
	updateRequest := updateproto.UpdateRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
	}
	errResp := validateSimpleUpdateRequest(updateRequest.RequestBody)
	if errResp.StatusCode != http.StatusOK {
		common.SetResponseHeader(ctx, errResp.Header)
		ctx.StatusCode(int(errResp.StatusCode))
		ctx.JSON(&errResp.Body)
		return
	}
	resp, err := a.SimpleUpdateRPC(updateRequest)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

func validateSimpleUpdateRequest(requestBody []byte) response.RPC {
	var request map[string]interface{}
	err := json.Unmarshal(requestBody, &request)
	if err != nil {
		errMsg := "Unable to parse the simple update request" + err.Error()
		log.Warn(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errMsg, nil, nil)
	}
	if request["Targets"] != nil {
		if reflect.TypeOf(request["Targets"]).Kind() != reflect.Slice {
			errMsg := "'Targets' parameter should be of type string array"
			log.Warn(errMsg)
			return common.GeneralError(http.StatusBadRequest, response.PropertyValueTypeError, errMsg, []interface{}{"", "Targets"}, nil)
		}
		target := request["Targets"].([]interface{})
		for _, k := range target {
			if reflect.TypeOf(k).Kind() != reflect.String {
				errMsg := "'Targets' parameter should be of type string array"
				log.Warn(errMsg)
				return common.GeneralError(http.StatusBadRequest, response.PropertyValueTypeError, errMsg, []interface{}{fmt.Sprintf("%v", k), "Targets"}, nil)
			}
		}
	}
	if request["ImageURI"] == nil {
		errMsg := "'ImageURI' parameter cannot be empty"
		log.Warn(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"ImageURI"}, nil)
	}
	if reflect.TypeOf(request["ImageURI"]).Kind() != reflect.String {
		errMsg := "'ImageURI' parameter should be of type string"
		log.Warn(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyValueTypeError, errMsg, []interface{}{"", "ImageURI"}, nil)
	}
	if request["ImageURI"] != nil {
		URI := request["ImageURI"]
		_, err = url.ParseRequestURI(URI.(string))
		if err != nil {
			errMsg := "Provided ImageURI is Invalid"
			log.Warn(errMsg)
			return common.GeneralError(http.StatusBadRequest, response.PropertyValueTypeError, errMsg, []interface{}{fmt.Sprintf("%v", err), "ImageURI"}, nil)
		}
	}
	return response.RPC{StatusCode: http.StatusOK}
}

//StartUpdate is a handler for start update action
func (a *UpdateRPCs) StartUpdate(ctx iris.Context) {
	defer ctx.Next()
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	updateRequest := updateproto.UpdateRequest{
		SessionToken: sessionToken,
	}
	resp, err := a.StartUpdateRPC(updateRequest)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}
