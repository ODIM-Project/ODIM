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
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	accountproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/account"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
)

// AccountRPCs defines all the RPC methods in account service
type AccountRPCs struct {
	GetServiceRPC     func(accountproto.AccountRequest) (*accountproto.AccountResponse, error)
	CreateRPC         func(accountproto.CreateAccountRequest) (*accountproto.AccountResponse, error)
	GetAllAccountsRPC func(accountproto.AccountRequest) (*accountproto.AccountResponse, error)
	GetAccountRPC     func(accountproto.GetAccountRequest) (*accountproto.AccountResponse, error)
	UpdateRPC         func(accountproto.UpdateAccountRequest) (*accountproto.AccountResponse, error)
	DeleteRPC         func(accountproto.DeleteAccountRequest) (*accountproto.AccountResponse, error)
}

// GetAccountService defines the GetAccountService iris handler.
// The method extract the session token and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (a *AccountRPCs) GetAccountService(ctx iris.Context) {
	defer ctx.Next()
	req := accountproto.AccountRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}

	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	resp, err := a.GetServiceRPC(req)
	if err != nil && resp == nil {
		errorMessage := "something went wrong with the RPC calls: " + err.Error()
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

// CreateAccount defines the CreateAccount iris handler.
// The method extract the session token, and necessary
// request parameters and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (a *AccountRPCs) CreateAccount(ctx iris.Context) {
	defer ctx.Next()
	var req interface{}

	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the account create request body: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}

	sessionToken := ctx.Request().Header.Get("X-Auth-Token")

	if sessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
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

	resp, err := a.CreateRPC(createRequest)
	if err != nil && resp == nil {
		errorMessage := "something went wrong with the RPC calls: " + err.Error()
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

// GetAllAccounts defines the GetAllAccounts iris handler.
// The method extract the session token and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (a *AccountRPCs) GetAllAccounts(ctx iris.Context) {
	defer ctx.Next()
	req := accountproto.AccountRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}

	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	resp, err := a.GetAllAccountsRPC(req)
	if err != nil && resp == nil {
		errorMessage := "something went wrong with the RPC calls: " + err.Error()
		log.Error(errorMessage)
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

}

// GetAccount defines the GetAccount iris handler.
// The method extract the session token, and necessary
// request parameters and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (a *AccountRPCs) GetAccount(ctx iris.Context) {
	defer ctx.Next()
	req := accountproto.GetAccountRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		AccountID:    ctx.Params().Get("id"),
	}

	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	resp, err := a.GetAccountRPC(req)
	if err != nil && resp == nil {
		errorMessage := "something went wrong with the RPC calls: " + err.Error()
		log.Error(errorMessage)
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

}

// UpdateAccount defines the UpdateAccount iris handler.
// The method extract the session token, and necessary
// request parameters and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (a *AccountRPCs) UpdateAccount(ctx iris.Context) {
	defer ctx.Next()
	var req interface{}

	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the account update request body: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}

	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	accountID := ctx.Params().Get("id")

	if sessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
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

	resp, err := a.UpdateRPC(updateRequest)
	if err != nil && resp == nil {
		errorMessage := "something went wrong with the RPC calls: " + err.Error()
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

// DeleteAccount defines the DeleteAccount iris handler.
// The method extract the session token, and necessary
// request parameters and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (a *AccountRPCs) DeleteAccount(ctx iris.Context) {
	defer ctx.Next()
	req := accountproto.DeleteAccountRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		AccountID:    ctx.Params().Get("id"),
	}

	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	resp, err := a.DeleteRPC(req)
	if err != nil && resp == nil {
		errorMessage := "something went wrong with the RPC calls: " + err.Error()
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
