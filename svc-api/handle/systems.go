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
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// SystemRPCs defines all the RPC methods in account service
type SystemRPCs struct {
	GetSystemsCollectionRPC    func(ctx context.Context, req systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error)
	GetSystemRPC               func(ctx context.Context, req systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error)
	GetSystemResourceRPC       func(ctx context.Context, req systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error)
	SystemResetRPC             func(ctx context.Context, req systemsproto.ComputerSystemResetRequest) (*systemsproto.SystemsResponse, error)
	SetDefaultBootOrderRPC     func(ctx context.Context, req systemsproto.DefaultBootOrderRequest) (*systemsproto.SystemsResponse, error)
	ChangeBiosSettingsRPC      func(ctx context.Context, req systemsproto.BiosSettingsRequest) (*systemsproto.SystemsResponse, error)
	ChangeBootOrderSettingsRPC func(ctx context.Context, req systemsproto.BootOrderSettingsRequest) (*systemsproto.SystemsResponse, error)
	CreateVolumeRPC            func(ctx context.Context, req systemsproto.VolumeRequest) (*systemsproto.SystemsResponse, error)
	DeleteVolumeRPC            func(ctx context.Context, req systemsproto.VolumeRequest) (*systemsproto.SystemsResponse, error)
}

// GetSystemsCollection fetches all systems
func (sys *SystemRPCs) GetSystemsCollection(ctx iris.Context) {
	ctxt := ctx.Request().Context()
	l.LogWithFields(ctxt).Info("Inside GetSystemCollection function (handler)")
	defer ctx.Next()
	req := systemsproto.GetSystemsRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
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
	resp, err := sys.GetSystemsCollectionRPC(ctxt, req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
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

// GetSystem fetches computer system details
func (sys *SystemRPCs) GetSystem(ctx iris.Context) {
	ctxt := ctx.Request().Context()
	defer ctx.Next()
	req := systemsproto.GetSystemsRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		RequestParam: ctx.Params().Get("id"),
		ResourceID:   ctx.Params().Get("rid"),
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
	resp, err := sys.GetSystemRPC(ctxt, req)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetSystemResource defines the GetSystemResource iris handler.
// The method extract the session token,uuid and request url and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (sys *SystemRPCs) GetSystemResource(ctx iris.Context) {
	ctxt := ctx.Request().Context()
	defer ctx.Next()
	req := systemsproto.GetSystemsRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		RequestParam: ctx.Params().Get("id"),
		ResourceID:   ctx.Params().Get("rid"),
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
	resp, err := sys.GetSystemResourceRPC(ctxt, req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	storageID := ctx.Params().Get("id2")
	switch req.URL {
	case "/redfish/v1/Systems/" + req.RequestParam + "/Bios/Settings":
		ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH")
	case "/redfish/v1/Systems/" + req.RequestParam + "/Storage/" + storageID + "/Volumes":
		ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	case "/redfish/v1/Systems/" + req.RequestParam + "/Storage/" + storageID + "/Volumes/" + req.ResourceID:
		ctx.ResponseWriter().Header().Set("Allow", "GET, DELETE")
	default:
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	}
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// ComputerSystemReset resets the indivitual computer systems
func (sys *SystemRPCs) ComputerSystemReset(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the system reset request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}
	systemID := ctx.Params().Get("id")
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
	resetRequest := systemsproto.ComputerSystemResetRequest{
		SessionToken: sessionToken,
		SystemID:     systemID,
		RequestBody:  request,
	}

	resp, err := sys.SystemResetRPC(ctxt, resetRequest)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
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

// SetDefaultBootOrder is the handler to set default boot order
// from iris context will get the request and check sessiontoken
// and do rpc call and send response back
func (sys *SystemRPCs) SetDefaultBootOrder(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req systemsproto.DefaultBootOrderRequest
	req.SystemID = ctx.Params().Get("id")
	req.SessionToken = ctx.Request().Header.Get("X-Auth-Token")
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := sys.SetDefaultBootOrderRPC(ctxt, req)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
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

// ChangeBiosSettings is the handler to set change bios settings
// from iris context will get the request and check sessiontoken
// and do rpc call and send response back
func (sys *SystemRPCs) ChangeBiosSettings(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the system reset request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}
	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
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
	biosRequest := systemsproto.BiosSettingsRequest{
		SessionToken: sessionToken,
		SystemID:     ctx.Params().Get("id"),
		RequestBody:  request,
	}
	resp, err := sys.ChangeBiosSettingsRPC(ctxt, biosRequest)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
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

// ChangeBootOrderSettings is the handler to set change boot order settings
// from iris context will get the request and check sessiontoken
// and do rpc call and send response back
func (sys *SystemRPCs) ChangeBootOrderSettings(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the system reset request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}
	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
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
	bootOrderRequest := systemsproto.BootOrderSettingsRequest{
		SessionToken: sessionToken,
		SystemID:     ctx.Params().Get("id"),
		RequestBody:  request,
	}
	resp, err := sys.ChangeBootOrderSettingsRPC(ctxt, bootOrderRequest)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
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

// CreateVolume is the handler to create a volume under storage
// from iris context will get the request and check sessiontoken
// and do rpc call and send response back
func (sys *SystemRPCs) CreateVolume(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the create volume request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}
	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
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
	volRequest := systemsproto.VolumeRequest{
		SessionToken:    sessionToken,
		SystemID:        ctx.Params().Get("id"),
		StorageInstance: ctx.Params().Get("id2"),
		RequestBody:     request,
	}
	resp, err := sys.CreateVolumeRPC(ctxt, volRequest)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
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

// DeleteVolume is the handler to delete a volume under storage
// from iris context will get the request and check sessiontoken
// and do rpc call and send response back
func (sys *SystemRPCs) DeleteVolume(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req interface{}
	ctx.ReadJSON(&req)
	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
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
	volRequest := systemsproto.VolumeRequest{
		SessionToken:    sessionToken,
		SystemID:        ctx.Params().Get("id"),
		StorageInstance: ctx.Params().Get("id2"),
		VolumeID:        ctx.Params().Get("rid"),
		RequestBody:     request,
	}
	resp, err := sys.DeleteVolumeRPC(ctxt, volRequest)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
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
