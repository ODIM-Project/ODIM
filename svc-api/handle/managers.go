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
	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// ManagersRPCs defines all the RPC methods in account service
type ManagersRPCs struct {
	GetManagersCollectionRPC      func(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error)
	GetManagersRPC                func(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error)
	GetManagersResourceRPC        func(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error)
	VirtualMediaInsertRPC         func(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error)
	VirtualMediaEjectRPC          func(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error)
	GetRemoteAccountServiceRPC    func(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error)
	CreateRemoteAccountServiceRPC func(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error)
	UpdateRemoteAccountServiceRPC func(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error)
	DeleteRemoteAccountServiceRPC func(ctx context.Context, req managersproto.ManagerRequest) (*managersproto.ManagerResponse, error)
}

// GetManagersCollection fetches all managers
func (mgr *ManagersRPCs) GetManagersCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := managersproto.ManagerRequest{
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
	resp, err := mgr.GetManagersCollectionRPC(ctxt, req)
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

// GetManager fetches computer managers details
func (mgr *ManagersRPCs) GetManager(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := managersproto.ManagerRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		ManagerID:    ctx.Params().Get("id"),
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
	resp, err := mgr.GetManagersRPC(ctxt, req)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
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

// GetManagersResource defines the GetManagersResource iris handler.
// The method extract the session token,uuid and request url and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (mgr *ManagersRPCs) GetManagersResource(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := managersproto.ManagerRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		ManagerID:    ctx.Params().Get("id"),
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
	resp, err := mgr.GetManagersResourceRPC(ctxt, req)
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

// VirtualMediaInsert defines the Insert virtual media iris handler
// The method extract the session token,uuid and request url and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (mgr *ManagersRPCs) VirtualMediaInsert(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var reqIn interface{}
	err := ctx.ReadJSON(&reqIn)
	if err != nil {
		errorMessage := "while trying to get JSON body from the virtual media actions request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}
	request, err := json.Marshal(reqIn)
	if err != nil {
		errorMessage := "while trying to create JSON request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	req := managersproto.ManagerRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		ManagerID:    ctx.Params().Get("id"),
		ResourceID:   ctx.Params().Get("rid"),
		URL:          ctx.Request().RequestURI,
		RequestBody:  request,
	}
	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := mgr.VirtualMediaInsertRPC(ctxt, req)
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

// VirtualMediaEject defines the eject virtual media iris handler
// The method extract the session token,uuid and request url and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (mgr *ManagersRPCs) VirtualMediaEject(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := managersproto.ManagerRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		ManagerID:    ctx.Params().Get("id"),
		ResourceID:   ctx.Params().Get("rid"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := mgr.VirtualMediaEjectRPC(ctxt, req)
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

// GetRemoteAccountService defines the GetRemoteAccountService iris handler.
// This method extract the session token,uuid and request url and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (mgr *ManagersRPCs) GetRemoteAccountService(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := managersproto.ManagerRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		ManagerID:    ctx.Params().Get("id"),
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
	resp, err := mgr.GetRemoteAccountServiceRPC(ctxt, req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	switch req.URL {
	case "/redfish/v1/Managers/" + req.ManagerID + "/RemoteAccountService/Accounts":
		ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	case "/redfish/v1/Managers/" + req.ManagerID + "/RemoteAccountService/Accounts/" + req.ResourceID:
		ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH, DELETE")
	default:
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	}
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// CreateRemoteAccountService defines the CreateRemoteAccountService iris handler.
// This method extract the session token,uuid,request payload and url and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (mgr *ManagersRPCs) CreateRemoteAccountService(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var reqIn interface{}
	err := ctx.ReadJSON(&reqIn)
	if err != nil {
		errorMessage := "while trying to get JSON body from the create remote account request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}
	request, err := json.Marshal(reqIn)
	if err != nil {
		errorMessage := "while trying to create JSON request body in create remote account: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	req := managersproto.ManagerRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		ManagerID:    ctx.Params().Get("id"),
		ResourceID:   ctx.Params().Get("rid"),
		URL:          ctx.Request().RequestURI,
		RequestBody:  request,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := mgr.CreateRemoteAccountServiceRPC(ctxt, req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
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

// UpdateRemoteAccountService defines the UpdateRemoteAccountService iris handler.
// This method extract the session token,uuid,request payload and url and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (mgr *ManagersRPCs) UpdateRemoteAccountService(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var reqIn interface{}
	err := ctx.ReadJSON(&reqIn)
	if err != nil {
		errorMessage := "while trying to get JSON body from the update remote account request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}
	request, err := json.Marshal(reqIn)
	if err != nil {
		errorMessage := "while trying to update JSON request body in update remote account: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	req := managersproto.ManagerRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		ManagerID:    ctx.Params().Get("id"),
		ResourceID:   ctx.Params().Get("rid"),
		URL:          ctx.Request().RequestURI,
		RequestBody:  request,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := mgr.UpdateRemoteAccountServiceRPC(ctxt, req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
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

// DeleteRemoteAccountService defines the DeleteRemoteAccountService iris handler.
// This method extract the session token,uuid and url and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (mgr *ManagersRPCs) DeleteRemoteAccountService(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := managersproto.ManagerRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		ManagerID:    ctx.Params().Get("id"),
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
	resp, err := mgr.DeleteRemoteAccountServiceRPC(ctxt, req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
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
