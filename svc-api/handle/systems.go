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
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
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
	UpdateSecureBootRPC        func(ctx context.Context, req systemsproto.SecureBootRequest) (*systemsproto.SystemsResponse, error)
	ResetSecureBootRPC         func(ctx context.Context, req systemsproto.SecureBootRequest) (*systemsproto.SystemsResponse, error)
}

// GetSystemsCollection fetches all systems
func (sys *SystemRPCs) GetSystemsCollection(ctx iris.Context) {
	ctxt := ctx.Request().Context()
	defer ctx.Next()
	req := systemsproto.GetSystemsRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting systems collection %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := sys.GetSystemsCollectionRPC(ctxt, req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting systems collection is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	sendSystemsResponse(ctx, resp)
}

// GetSystem fetches computer system details
func (sys *SystemRPCs) GetSystem(ctx iris.Context) {
	ctxt := ctx.Request().Context()
	defer ctx.Next()
	req := systemsproto.GetSystemsRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		RequestParam: ctx.Params().Get("id"),
		ResourceID:   ctx.Params().Get("rid"),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting system with URL %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := sys.GetSystemRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting system details is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH")
	sendSystemsResponse(ctx, resp)
}

// GetSystemResource defines the GetSystemResource iris handler.
// The method extract the session token,uuid and request url and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (sys *SystemRPCs) GetSystemResource(ctx iris.Context) {
	ctxt := ctx.Request().Context()
	defer ctx.Next()
	req := systemsproto.GetSystemsRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		RequestParam: ctx.Params().Get("id"),
		ResourceID:   ctx.Params().Get("rid"),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting system resources with URL %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := sys.GetSystemResourceRPC(ctxt, req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}

	storageID := ctx.Params().Get("id2")
	switch req.URL {
	case "/redfish/v1/Systems/" + req.RequestParam + "/Bios/Settings":
		ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH")
	case "/redfish/v1/Systems/" + req.RequestParam + "/Storage/" + storageID + "/Volumes":
		ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	case "/redfish/v1/Systems/" + req.RequestParam + "/Storage/" + storageID + "/Volumes/" + req.ResourceID:
		ctx.ResponseWriter().Header().Set("Allow", "GET, DELETE")
	case "/redfish/v1/Systems/" + req.RequestParam + "/SecureBoot":
		ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH")
	default:
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting systems resources is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendSystemsResponse(ctx, resp)
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
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}
	systemID := ctx.Params().Get("id")
	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)
	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	// Marshalling the req to make reset request
	request, err := json.Marshal(req)
	resetRequest := systemsproto.ComputerSystemResetRequest{
		SessionToken: sessionToken,
		SystemID:     systemID,
		RequestBody:  request,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for computer system reset with request body %s", string(request))
	resp, err := sys.SystemResetRPC(ctxt, resetRequest)
	if err != nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for computer system reset is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendSystemsResponse(ctx, resp)
}

// SetDefaultBootOrder is the handler to set default boot order
// from iris context will get the request and check sessiontoken
// and do rpc call and send response back
func (sys *SystemRPCs) SetDefaultBootOrder(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req systemsproto.DefaultBootOrderRequest
	req.SystemID = ctx.Params().Get("id")
	req.SessionToken = ctx.Request().Header.Get(AuthTokenHeader)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for setting default boot order with request id %s", req.SystemID)
	resp, err := sys.SetDefaultBootOrderRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for setting default boot order is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendSystemsResponse(ctx, resp)
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
		errorMessage := "error while trying to get JSON body from bios setting request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}
	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for changing bios setting with request body %s", string(request))
	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)
	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	biosRequest := systemsproto.BiosSettingsRequest{
		SessionToken: sessionToken,
		SystemID:     ctx.Params().Get("id"),
		RequestBody:  request,
	}
	resp, err := sys.ChangeBiosSettingsRPC(ctxt, biosRequest)
	if err != nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for changing bios setting is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendSystemsResponse(ctx, resp)
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
		errorMessage := "error while trying to get JSON body from change boot order setting request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}
	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for changing boot order setting with request body %s", string(request))
	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)
	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	bootOrderRequest := systemsproto.BootOrderSettingsRequest{
		SessionToken: sessionToken,
		SystemID:     ctx.Params().Get("id"),
		RequestBody:  request,
	}
	resp, err := sys.ChangeBootOrderSettingsRPC(ctxt, bootOrderRequest)
	if err != nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting changing boot order setting is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendSystemsResponse(ctx, resp)
}

// UpdateSecureBoot is the handler to set change boot order settings
// from iris context will get the request and check session token
// and do rpc call and send response back
func (sys *SystemRPCs) UpdateSecureBoot(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from update SecureBoot request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}
	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for updating SecureBoot with request body %s", string(request))
	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)
	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	secureBootRequest := systemsproto.SecureBootRequest{
		SessionToken: sessionToken,
		SystemID:     ctx.Params().Get("id"),
		RequestBody:  request,
	}
	resp, err := sys.UpdateSecureBootRPC(ctxt, secureBootRequest)
	if err != nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for updating SecureBoot is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendSystemsResponse(ctx, resp)
}

// ResetSecureBoot shall reset the UEFI Secure Boot key databases.
// The `ResetAllKeysToDefault` value shall reset all UEFI Secure Boot key databases to their default values.
// The `DeleteAllKeys` value shall delete the content of all UEFI Secure Boot key databases.
// The `DeletePK` value shall delete the content of the PK Secure Boot key database.
func (sys *SystemRPCs) ResetSecureBoot(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from reset SecureBoot key databases request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}
	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for resetting SecureBoot with request body %s", string(request))
	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)
	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	secureBootRequest := systemsproto.SecureBootRequest{
		SessionToken: sessionToken,
		SystemID:     ctx.Params().Get("id"),
		RequestBody:  request,
	}
	resp, err := sys.ResetSecureBootRPC(ctxt, secureBootRequest)
	if err != nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for updating SecureBoot is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendSystemsResponse(ctx, resp)
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
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}
	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for creating volume with request body %s", string(request))
	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)
	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	volRequest := systemsproto.VolumeRequest{
		SessionToken:    sessionToken,
		SystemID:        ctx.Params().Get("id"),
		StorageInstance: ctx.Params().Get("id2"),
		RequestBody:     request,
	}
	resp, err := sys.CreateVolumeRPC(ctxt, volRequest)
	if err != nil {
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for creating a volume is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendSystemsResponse(ctx, resp)
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
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for deleting volume with request body %s", string(request))
	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)
	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
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
		errorMessage := rpcFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for deleting a volume is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	sendSystemsResponse(ctx, resp)
}

// sendSystemsResponse writes the systems response to client
func sendSystemsResponse(ctx iris.Context, resp *systemsproto.SystemsResponse) {
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}
