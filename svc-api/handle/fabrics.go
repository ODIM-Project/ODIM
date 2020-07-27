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
	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// FabricRPCs defines all the RPC methods in fabric service
type FabricRPCs struct {
	GetFabricResourceRPC    func(fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error)
	UpdateFabricResourceRPC func(fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error)
	DeleteFabricResourceRPC func(fabricsproto.FabricRequest) (*fabricsproto.FabricResponse, error)
}

// GetFabricResource defines the GetFabricResource iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricResource(ctx iris.Context) {
	req := fabricsproto.FabricRequest{
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

	resp, err := f.GetFabricResourceRPC(req)
	if err != nil && resp == nil {
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

// UpdateFabricResource defines the UpdateFabricResource iris handler.
// The method updates if Fabric Resource exists else creates new one.
func (f *FabricRPCs) UpdateFabricResource(ctx iris.Context) {
	req := fabricsproto.FabricRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
		Method:       ctx.Request().Method,
	}

	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	var createReq interface{}
	err := ctx.ReadJSON(&createReq)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the  request body: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusBadRequest) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	// marshalling the req to make fabric UpdateFabricResource request
	// Since fabric FabricRequest accepts []byte stream
	request, err := json.Marshal(createReq)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	req.RequestBody = request
	resp, err := f.UpdateFabricResourceRPC(req)
	if err != nil && resp == nil {
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

// DeleteFabricResource defines the DeleteFabricResource iris handler.
// This method is used for deleting requested fabric resource
func (f *FabricRPCs) DeleteFabricResource(ctx iris.Context) {
	req := fabricsproto.FabricRequest{
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

	resp, err := f.DeleteFabricResourceRPC(req)
	if err != nil && resp == nil {
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
