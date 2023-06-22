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
	accountproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/account"
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
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
	}
	l.LogWithFields(ctx).Debug("Incoming request received for the Get Account service")
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := a.GetServiceRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcCallFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}
	l.LogWithFields(ctx).Debugf("Outgoing response for Getting Account service is %s and response status %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	sendAccountResponse(ctx, resp)

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
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debug("Incoming request for create account received")
	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)

	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendInvalidSessionResponse(ctx, errorMessage)
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
		errorMessage := rpcCallFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}

	sendAccountResponse(ctx, resp)
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
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
	}
	l.LogWithFields(ctxt).Debug("Incoming request for get all accounts received")
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := a.GetAllAccountsRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcCallFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	sendAccountResponse(ctx, resp)
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
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		AccountID:    ctx.Params().Get("id"),
	}
	l.LogWithFields(ctxt).Debugf("Incoming request for get account info received for %s", req.AccountID)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := a.GetAccountRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcCallFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH, DELETE")
	sendAccountResponse(ctx, resp)
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
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}

	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)
	accountID := ctx.Params().Get("id")
	l.LogWithFields(ctxt).Debugf("Incoming request for updating account received for %s", accountID)
	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendInvalidSessionResponse(ctx, errorMessage)
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
		errorMessage := rpcCallFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}

	sendAccountResponse(ctx, resp)
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
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		AccountID:    ctx.Params().Get("id"),
	}
	l.LogWithFields(ctxt).Debugf("Incoming request for deleting account received with %s", req.AccountID)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := a.DeleteRPC(ctxt, req)
	if err != nil && resp == nil {
		errorMessage := rpcCallFailedErrMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
		return
	}

	sendAccountResponse(ctx, resp)
	l.LogWithFields(ctxt).Debugf("outgoing response for deleting account with %s and response status %d", req.AccountID, int(resp.StatusCode))
}

// sendAccountResponse writes the account response to client
func sendAccountResponse(ctx iris.Context, resp *accountproto.AccountResponse) {
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}
