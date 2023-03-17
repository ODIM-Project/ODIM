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
	accountproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/account"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// AccountRPCs defines all the RPC methods in account service
type AccountRPCs struct {
	GetServiceRPC     func(context.Context, accountproto.AccountRequest) (*accountproto.AccountResponse, error)
	CreateRPC         func(context.Context, accountproto.CreateAccountRequest) (*accountproto.AccountResponse, error)
	GetAllAccountsRPC func(context.Context, accountproto.AccountRequest) (*accountproto.AccountResponse, error)
	GetAccountRPC     func(context.Context, accountproto.GetAccountRequest) (*accountproto.AccountResponse, error)
	UpdateRPC         func(context.Context, accountproto.UpdateAccountRequest) (*accountproto.AccountResponse, error)
	DeleteRPC         func(context.Context, accountproto.DeleteAccountRequest) (*accountproto.AccountResponse, error)
}

// GetAccountService defines the GetAccountService iris handler.
// The method extract the session token and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (a *AccountRPCs) GetAccountService(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := accountproto.AccountRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	l.LogWithFields(ctx).Debug("Incoming request received for the Get Account service")
	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	resp, err := a.GetServiceRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := "something went wrong with the RPC calls: " + err.Error()
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
	l.LogWithFields(ctx).Debugf("Outgoing response for Getting Account service is %s and response status %d", string(resp.Body), int(resp.StatusCode))
}

// CreateAccount defines the CreateAccount iris handler.
// The method extract the session token, and necessary
// request parameters and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (a *AccountRPCs) CreateAccount(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the account create request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debug("Incoming request for create account received")
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")

	if sessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	// Marshalling the req to make account request
	// Since create request accepts byte stream
	request, err := json.Marshal(req)
	createRequest := accountproto.CreateAccountRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
	}

	resp, err := a.CreateRPC(ctxt, createRequest)
	if err != nil && resp == nil {
		errorMessage := "something went wrong with the RPC calls: " + err.Error()
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
	l.LogWithFields(ctxt).Debugf("Outgoing response for create account is %s and response status %d", string(resp.Body), int(resp.StatusCode))

}

// GetAllAccounts defines the GetAllAccounts iris handler.
// The method extract the session token and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (a *AccountRPCs) GetAllAccounts(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := accountproto.AccountRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	l.LogWithFields(ctxt).Debug("Incoming request for get all accounts received")
	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	resp, err := a.GetAllAccountsRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := "something went wrong with the RPC calls: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
	l.LogWithFields(ctxt).Debugf("outgoing response for get all accounts is %s and response status %d", string(resp.Body), int(resp.StatusCode))

}

// GetAccount defines the GetAccount iris handler.
// The method extract the session token, and necessary
// request parameters and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (a *AccountRPCs) GetAccount(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := accountproto.GetAccountRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		AccountID:    ctx.Params().Get("id"),
	}
	l.LogWithFields(ctxt).Debug("Incoming request for get account info received for %s", req.AccountID)
	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	resp, err := a.GetAccountRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := "something went wrong with the RPC calls: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH, DELETE")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
	l.LogWithFields(ctxt).Debugf("outgoing response for get account is %s and response status %d", string(resp.Body), int(resp.StatusCode))

}

// UpdateAccount defines the UpdateAccount iris handler.
// The method extract the session token, and necessary
// request parameters and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (a *AccountRPCs) UpdateAccount(ctx iris.Context) {
	defer ctx.Next()
	var req interface{}
	ctxt := ctx.Request().Context()

	err := ctx.ReadJSON(&req)

	if err != nil {
		errorMessage := "error while trying to get JSON body from the account update request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}

	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	accountID := ctx.Params().Get("id")
	l.LogWithFields(ctxt).Debugf("Incoming request for updating account received for %s", accountID)
	if sessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	// Marshalling the req to make account request
	// Since account update request accepts byte stream
	request, err := json.Marshal(req)
	updateRequest := accountproto.UpdateAccountRequest{
		SessionToken: sessionToken,
		AccountID:    accountID,
		RequestBody:  request,
	}

	resp, err := a.UpdateRPC(ctxt, updateRequest)
	if err != nil && resp == nil {
		errorMessage := "something went wrong with the RPC calls: " + err.Error()
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
	l.LogWithFields(ctxt).Debugf("outgoing response for updating account is %s and response status %d", string(resp.Body), int(resp.StatusCode))

}

// DeleteAccount defines the DeleteAccount iris handler.
// The method extract the session token, and necessary
// request parameters and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (a *AccountRPCs) DeleteAccount(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := accountproto.DeleteAccountRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		AccountID:    ctx.Params().Get("id"),
	}
	l.LogWithFields(ctxt).Debugf("Incoming request for deleting account received with %s", req.AccountID)
	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	resp, err := a.DeleteRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := "something went wrong with the RPC calls: " + err.Error()
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
	l.LogWithFields(ctxt).Debugf("outgoing response for deleting account with %s and response status %d", req.AccountID, int(resp.StatusCode))

}
