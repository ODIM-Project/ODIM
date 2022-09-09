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

// Package rpc ...
package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	accountproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/account"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/account"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
)

// Account struct helps to register service
type Account struct{}

var (
	GetAllAccountsFunc    = account.GetAllAccounts
	GetAccountFunc        = account.GetAccount
	GetAccountServiceFunc = account.GetAccountService
	AccDeleteFunc         = account.Delete
)

// Create defines the operations which handles the RPC request response
// for the create account service of account-session micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Account) Create(ctx context.Context, req *accountproto.CreateAccountRequest) (*accountproto.AccountResponse, error) {
	var resp accountproto.AccountResponse
	errorArgs := []response.ErrArgs{
		response.ErrArgs{
			StatusMessage: "",
			ErrorMessage:  "",
			MessageArgs:   []interface{}{},
		},
	}
	args := &response.Args{
		Code:      response.GeneralError,
		Message:   "",
		ErrorArgs: errorArgs,
	}

	sess, errs := CheckSessionTimeOutFunc(req.SessionToken)
	if errs != nil {
		errorMessage := "error while authorizing session token: " + errs.Error()
		resp.StatusCode, resp.StatusMessage = errs.GetAuthStatusCodeAndMessage()
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body)
			l.Log.Error(errorMessage)
		} else {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body)
			auth.CustomAuthLog(req.SessionToken, "Invalid session token", resp.StatusCode)
		}
		return &resp, nil
	}
	err := UpdateLastUsedTimeFunc(req.SessionToken)
	if err != nil {
		errorMessage := "error while updating last used time of session with token " + req.SessionToken + ": " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		l.Log.Error(errorMessage)
		return &resp, nil
	}

	acc := account.GetExternalInterface()
	data, err := acc.Create(req, sess)
	var jsonErr error // jsonErr is created to protect the data in err
	resp.Body, jsonErr = MarshalFunc(data.Body)
	if jsonErr != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying marshal the response body for create account: " + jsonErr.Error()
		l.Log.Error(resp.StatusMessage)
		return &resp, nil
	}
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header

	return &resp, nil
}

// GetAllAccounts defines the operations which handles the RPC request response
// for the list all account service of account-session micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Account) GetAllAccounts(ctx context.Context, req *accountproto.AccountRequest) (*accountproto.AccountResponse, error) {
	var resp accountproto.AccountResponse
	errorArgs := []response.ErrArgs{
		response.ErrArgs{
			StatusMessage: "",
			ErrorMessage:  "",
			MessageArgs:   []interface{}{},
		},
	}
	args := &response.Args{
		Code:      response.GeneralError,
		Message:   "",
		ErrorArgs: errorArgs,
	}
	sess, errs := CheckSessionTimeOutFunc(req.SessionToken)
	if errs != nil {
		errorMessage := "error while authorizing session token: " + errs.Error()
		resp.StatusCode, resp.StatusMessage = errs.GetAuthStatusCodeAndMessage()
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body)
			l.Log.Error(errorMessage)
		} else {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body)
			auth.CustomAuthLog(req.SessionToken, "Invalid session token", resp.StatusCode)
		}
		return &resp, nil
	}

	err := UpdateLastUsedTimeFunc(req.SessionToken)
	if err != nil {
		errorMessage := "error while updating last used time of session with token " + req.SessionToken + ": " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		l.Log.Error(errorMessage)
		return &resp, nil
	}

	data := GetAllAccountsFunc(sess)
	resp.Body, err = MarshalFunc(data.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying marshal the response body for get all accounts: " + err.Error()
		l.Log.Error(resp.StatusMessage)
		return &resp, fmt.Errorf(resp.StatusMessage)
	}
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header

	return &resp, err
}

// GetAccount defines the operations which handles the RPC request response
// for the view of a particular account service of account-session micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Account) GetAccount(ctx context.Context, req *accountproto.GetAccountRequest) (*accountproto.AccountResponse, error) {
	var resp accountproto.AccountResponse
	errorArgs := []response.ErrArgs{
		response.ErrArgs{
			StatusMessage: "",
			ErrorMessage:  "",
			MessageArgs:   []interface{}{},
		},
	}
	args := &response.Args{
		Code:      response.GeneralError,
		Message:   "",
		ErrorArgs: errorArgs,
	}
	sess, errs := CheckSessionTimeOutFunc(req.SessionToken)
	if errs != nil {
		errorMessage := "error while authorizing session token: " + errs.Error()
		resp.StatusCode, resp.StatusMessage = errs.GetAuthStatusCodeAndMessage()
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body)
			l.Log.Error(errorMessage)
		} else {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body)
			auth.CustomAuthLog(req.SessionToken, "Invalid session token", resp.StatusCode)
		}
		return &resp, nil
	}

	err := UpdateLastUsedTimeFunc(req.SessionToken)
	if err != nil {
		errorMessage := "error while updating last used time of session with token " + req.SessionToken + ": " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		l.Log.Error(errorMessage)
		return &resp, nil
	}

	data := GetAccountFunc(sess, req.AccountID)
	resp.Body, err = MarshalFunc(data.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying marshal the response body for get account details: " + err.Error()
		l.Log.Error(resp.StatusMessage)
		return &resp, fmt.Errorf(resp.StatusMessage)
	}
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header

	return &resp, nil
}

// GetAccountServices defines the operations which handles the RPC request response
// for checking the availability of account-session micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Account) GetAccountServices(ctx context.Context, req *accountproto.AccountRequest) (*accountproto.AccountResponse, error) {
	var resp accountproto.AccountResponse
	errorArgs := []response.ErrArgs{
		response.ErrArgs{
			StatusMessage: "",
			ErrorMessage:  "",
			MessageArgs:   []interface{}{},
		},
	}
	args := &response.Args{
		Code:      response.GeneralError,
		Message:   "",
		ErrorArgs: errorArgs,
	}
	_, errs := CheckSessionTimeOutFunc(req.SessionToken)
	if errs != nil {
		errorMessage := "error while authorizing session token: " + errs.Error()
		resp.StatusCode, resp.StatusMessage = errs.GetAuthStatusCodeAndMessage()
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body)
			l.Log.Error(errorMessage)
		} else {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body)
			auth.CustomAuthLog(req.SessionToken, "Invalid session token", resp.StatusCode)
		}
		return &resp, nil
	}

	err := UpdateLastUsedTimeFunc(req.SessionToken)
	if err != nil {
		errorMessage := "error while updating last used time of session with token " + req.SessionToken + ": " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		l.Log.Printf(errorMessage)
		return &resp, nil
	}

	data := GetAccountServiceFunc()
	resp.Body, err = MarshalFunc(data.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying marshal the response body for get account details: " + err.Error()
		l.Log.Printf(resp.StatusMessage)
		return &resp, fmt.Errorf(resp.StatusMessage)
	}
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header

	return &resp, err
}

// Update defines the operations which handles the RPC request response
// for the update of a particular account service of account-session micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Account) Update(ctx context.Context, req *accountproto.UpdateAccountRequest) (*accountproto.AccountResponse, error) {
	var resp accountproto.AccountResponse
	errorArgs := []response.ErrArgs{
		response.ErrArgs{
			StatusMessage: "",
			ErrorMessage:  "",
			MessageArgs:   []interface{}{},
		},
	}
	args := &response.Args{
		Code:      response.GeneralError,
		Message:   "",
		ErrorArgs: errorArgs,
	}
	sess, errs := CheckSessionTimeOutFunc(req.SessionToken)
	if errs != nil {
		errorMessage := "error while authorizing session token: " + errs.Error()
		resp.StatusCode, resp.StatusMessage = errs.GetAuthStatusCodeAndMessage()
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body)
			l.Log.Error(errorMessage)
		} else {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body)
			auth.CustomAuthLog(req.SessionToken, "Invalid session token", resp.StatusCode)
		}
		return &resp, nil
	}

	err := UpdateLastUsedTimeFunc(req.SessionToken)
	if err != nil {
		errorMessage := "error while updating last used time of session with token " + req.SessionToken + ": " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		l.Log.Error(errorMessage)
		return &resp, nil
	}

	acc := account.GetExternalInterface()

	data := acc.Update(req, sess)
	resp.Body, err = MarshalFunc(data.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying to marshal the response body for create account: " + err.Error()
		l.Log.Printf(resp.StatusMessage)
		return &resp, nil
	}
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header

	return &resp, nil
}

// Delete defines the operations which handles the RPC request response
// for the delete of a particular account service of account-session micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Account) Delete(ctx context.Context, req *accountproto.DeleteAccountRequest) (*accountproto.AccountResponse, error) {
	var resp accountproto.AccountResponse
	errorArgs := []response.ErrArgs{
		response.ErrArgs{
			StatusMessage: "",
			ErrorMessage:  "",
			MessageArgs:   []interface{}{},
		},
	}
	args := &response.Args{
		Code:      response.GeneralError,
		Message:   "",
		ErrorArgs: errorArgs,
	}
	sess, errs := CheckSessionTimeOutFunc(req.SessionToken)
	if errs != nil {
		errorMessage := "error while authorizing session token: " + errs.Error()
		resp.StatusCode, resp.StatusMessage = errs.GetAuthStatusCodeAndMessage()
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body)
			l.Log.Error(errorMessage)
		} else {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body)
			auth.CustomAuthLog(req.SessionToken, "Invalid session token", resp.StatusCode)
		}
		return &resp, nil
	}

	err := UpdateLastUsedTimeFunc(req.SessionToken)
	if err != nil {
		errorMessage := "error while updating last used time of session with token " + req.SessionToken + ": " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		l.Log.Error(errorMessage)
		return &resp, nil
	}

	data := AccDeleteFunc(sess, req.AccountID)
	var jsonErr error // jsonErr is created to protect the data in err
	resp.Body, jsonErr = MarshalFunc(data.Body)
	if jsonErr != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying marshal the response body for delete account: " + jsonErr.Error()
		l.Log.Error(resp.StatusMessage)
		return &resp, nil
	}
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header

	return &resp, nil
}
