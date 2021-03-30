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
	log "github.com/sirupsen/logrus"
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

// GetFabricCollection defines the GetFabricCollection iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricCollection(ctx iris.Context) {
	req := fabricsproto.FabricRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}

	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	resp, err := f.GetFabricResourceRPC(req)
	if err != nil && resp == nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
	}
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetFabric defines the GetFabric iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabric(ctx iris.Context) {
	req := fabricsproto.FabricRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}

	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	resp, err := f.GetFabricResourceRPC(req)
	if err != nil && resp == nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
	}
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetFabricSwitchCollection defines the GetFabricSwitchCollection iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricSwitchCollection(ctx iris.Context) {
	req := fabricsproto.FabricRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}

	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	resp, err := f.GetFabricResourceRPC(req)
	if err != nil && resp == nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
	}
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetFabricSwitch defines the GetFabricSwitch iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricSwitch(ctx iris.Context) {
	req := fabricsproto.FabricRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}

	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	resp, err := f.GetFabricResourceRPC(req)
	if err != nil && resp == nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
	}
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetSwitchPortCollection defines the GetSwitchPortCollection iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetSwitchPortCollection(ctx iris.Context) {
	req := fabricsproto.FabricRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}

	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	resp, err := f.GetFabricResourceRPC(req)
	if err != nil && resp == nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
	}
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetSwitchPort defines the GetSwitchPort iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetSwitchPort(ctx iris.Context) {
	req := fabricsproto.FabricRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}

	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	resp, err := f.GetFabricResourceRPC(req)
	if err != nil && resp == nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp.Header = map[string]string{
		"Allow":             `"GET", "PATCH"`,
		"Cache-Control":     "no-cache",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
	}
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetFabricZoneCollection defines the GetFabricZoneCollection iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricZoneCollection(ctx iris.Context) {
	req := fabricsproto.FabricRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}

	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	resp, err := f.GetFabricResourceRPC(req)
	if err != nil && resp == nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp.Header = map[string]string{
		"Allow":             `"GET", "POST"`,
		"Cache-Control":     "no-cache",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
	}
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetFabricZone defines the GetFabricZone iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricZone(ctx iris.Context) {
	req := fabricsproto.FabricRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}

	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	resp, err := f.GetFabricResourceRPC(req)
	if err != nil && resp == nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp.Header = map[string]string{
		"Allow":             `"GET", "PUT", "PATCH", "DELETE"`,
		"Cache-Control":     "no-cache",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
	}
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetFabricEndPointCollection defines the GetFabricEndPointCollection iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricEndPointCollection(ctx iris.Context) {
	req := fabricsproto.FabricRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}

	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	resp, err := f.GetFabricResourceRPC(req)
	if err != nil && resp == nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp.Header = map[string]string{
		"Allow":             `"GET", "POST"`,
		"Cache-Control":     "no-cache",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
	}
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetFabricEndPoints defines the GetFabricEndPoints iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricEndPoints(ctx iris.Context) {
	req := fabricsproto.FabricRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}

	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	resp, err := f.GetFabricResourceRPC(req)
	if err != nil && resp == nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp.Header = map[string]string{
		"Allow":             `"GET", "PUT", "PATCH", "DELETE"`,
		"Cache-Control":     "no-cache",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
	}
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetFabricAddressPoolCollection defines the GetFabricAddressPoolCollection iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricAddressPoolCollection(ctx iris.Context) {
	req := fabricsproto.FabricRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}

	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	resp, err := f.GetFabricResourceRPC(req)
	if err != nil && resp == nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp.Header = map[string]string{
		"Allow":             `"GET",  "POST"`,
		"Cache-Control":     "no-cache",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
	}
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetFabricAddressPool defines the GetFabricAddressPool iris handler.
// The method extracts given Fabric Resource
func (f *FabricRPCs) GetFabricAddressPool(ctx iris.Context) {
	req := fabricsproto.FabricRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}

	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	resp, err := f.GetFabricResourceRPC(req)
	if err != nil && resp == nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp.Header = map[string]string{
		"Allow":             `"GET", "PUT", "PATCH", "DELETE"`,
		"Cache-Control":     "no-cache",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
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
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	var createReq interface{}
	err := ctx.ReadJSON(&createReq)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the  request body: " + err.Error()
		log.Error(errorMessage)
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
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	req.RequestBody = request
	resp, err := f.UpdateFabricResourceRPC(req)
	if err != nil && resp == nil {
		errorMessage := "RPC error: " + err.Error()
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

// DeleteFabricResource defines the DeleteFabricResource iris handler.
// This method is used for deleting requested fabric resource
func (f *FabricRPCs) DeleteFabricResource(ctx iris.Context) {
	req := fabricsproto.FabricRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}

	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	resp, err := f.DeleteFabricResourceRPC(req)
	if err != nil && resp == nil {
		errorMessage := "RPC error: " + err.Error()
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
