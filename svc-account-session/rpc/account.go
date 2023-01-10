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
	"os"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
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

// podName defines the current name of process
var podName = os.Getenv("POD_NAME")

// Create defines the operations which handles the RPC request response
// for the create account service of account-session micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Account) Create(ctx context.Context, req *accountproto.CreateAccountRequest) (*accountproto.AccountResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside Create function (svc-account-session)")
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

	l.LogWithFields(ctx).Info("Validating session and updating the last used time of the session before creating the account")
	sess, errs := CheckSessionTimeOutFunc(ctx, req.SessionToken)
	if errs != nil {
		resp.Body, resp.StatusCode, resp.StatusMessage = validateSessionTimeoutError(ctx, req.SessionToken, errs)
		return &resp, nil
	}

	err := UpdateLastUsedTimeFunc(ctx, req.SessionToken)
	if err != nil {
		errorArgs[0].ErrorMessage, resp.StatusCode, resp.StatusMessage = validateUpdateLastUsedTimeError(ctx, err, req.SessionToken)
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		return &resp, nil
	}

	acc := account.GetExternalInterface()
	data, err := acc.Create(ctx, req, sess)
	var jsonErr error // jsonErr is created to protect the data in err
	body, jsonErr := MarshalFunc(data.Body)
	if jsonErr != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying to marshal the response body of the create account API: " + jsonErr.Error()
		l.LogWithFields(ctx).Error(resp.StatusMessage)
		return &resp, nil
	}
	l.LogWithFields(ctx).Debugf("outgoing response of request to create an account: %s", string(body))
	resp.Body = body
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header

	return &resp, nil
}

// GetAllAccounts defines the operations which handles the RPC request response
// for the list all account service of account-session micro service.
// The functionality retrieves the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Account) GetAllAccounts(ctx context.Context, req *accountproto.AccountRequest) (*accountproto.AccountResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside GetAllAccounts function (svc-account-session)")
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

	l.LogWithFields(ctx).Info("Validating session and updating the last used time of the session before fetching all accounts")
	sess, errs := CheckSessionTimeOutFunc(ctx, req.SessionToken)
	if errs != nil {
		resp.Body, resp.StatusCode, resp.StatusMessage = validateSessionTimeoutError(ctx, req.SessionToken, errs)
		return &resp, nil
	}

	err := UpdateLastUsedTimeFunc(ctx, req.SessionToken)
	if err != nil {
		errorArgs[0].ErrorMessage, resp.StatusCode, resp.StatusMessage = validateUpdateLastUsedTimeError(ctx, err, req.SessionToken)
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		return &resp, nil
	}

	data := GetAllAccountsFunc(ctx, sess)
	body, err := MarshalFunc(data.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying to marshal the response body of the get all accounts API: " + err.Error()
		l.LogWithFields(ctx).Error(resp.StatusMessage)
		return &resp, fmt.Errorf(resp.StatusMessage)
	}
	l.LogWithFields(ctx).Debugf("outgoing response of request to view all accounts: %s", string(body))
	resp.Body = body
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header

	return &resp, err
}

// GetAccount defines the operations which handles the RPC request response
// for the view of a particular account service of account-session micro service.
// The functionality retrieves the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Account) GetAccount(ctx context.Context, req *accountproto.GetAccountRequest) (*accountproto.AccountResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside GetAccount function (svc-account-session)")
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

	l.LogWithFields(ctx).Info("Validating session and updating the last used time of the session before fetching the account")
	sess, errs := CheckSessionTimeOutFunc(ctx, req.SessionToken)
	if errs != nil {
		resp.Body, resp.StatusCode, resp.StatusMessage = validateSessionTimeoutError(ctx, req.SessionToken, errs)
		return &resp, nil
	}

	err := UpdateLastUsedTimeFunc(ctx, req.SessionToken)
	if err != nil {
		errorArgs[0].ErrorMessage, resp.StatusCode, resp.StatusMessage = validateUpdateLastUsedTimeError(ctx, err, req.SessionToken)
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		return &resp, nil
	}

	data := GetAccountFunc(ctx, sess, req.AccountID)
	body, err := MarshalFunc(data.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying to marshal the response body of the get account API: " + err.Error()
		l.LogWithFields(ctx).Error(resp.StatusMessage)
		return &resp, fmt.Errorf(resp.StatusMessage)
	}
	l.LogWithFields(ctx).Debugf("outgoing response of request to view the account: %s", string(body))
	resp.Body = body
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
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside GetAccountService function (svc-account-session)")
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
	l.LogWithFields(ctx).Info("Validating session and updating the last used time of the session before checking the availability of account session")
	_, errs := CheckSessionTimeOutFunc(ctx, req.SessionToken)
	if errs != nil {
		resp.Body, resp.StatusCode, resp.StatusMessage = validateSessionTimeoutError(ctx, req.SessionToken, errs)
		return &resp, nil
	}

	err := UpdateLastUsedTimeFunc(ctx, req.SessionToken)
	if err != nil {
		errorArgs[0].ErrorMessage, resp.StatusCode, resp.StatusMessage = validateUpdateLastUsedTimeError(ctx, err, req.SessionToken)
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		return &resp, nil
	}

	data := GetAccountServiceFunc(ctx)
	body, err := MarshalFunc(data.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying to marshal the response body of the get account service API: " + err.Error()
		l.LogWithFields(ctx).Printf(resp.StatusMessage)
		return &resp, fmt.Errorf(resp.StatusMessage)
	}
	l.LogWithFields(ctx).Debugf("outgoing response of request to view the account session: %s", string(body))
	resp.Body = body
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header

	return &resp, err
}

// Update defines the operations which handles the RPC request response
// for the update of a particular account service of account-session micro service.
// The functionality retrieves the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Account) Update(ctx context.Context, req *accountproto.UpdateAccountRequest) (*accountproto.AccountResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside Update function (svc-account-session)")
	var resp accountproto.AccountResponse
	l.LogWithFields(ctx).Info("Validating session and updating the last used time of the session before updating the account")
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
	sess, errs := CheckSessionTimeOutFunc(ctx, req.SessionToken)
	if errs != nil {
		resp.Body, resp.StatusCode, resp.StatusMessage = validateSessionTimeoutError(ctx, req.SessionToken, errs)
		return &resp, nil
	}

	err := UpdateLastUsedTimeFunc(ctx, req.SessionToken)
	if err != nil {
		errorArgs[0].ErrorMessage, resp.StatusCode, resp.StatusMessage = validateUpdateLastUsedTimeError(ctx, err, req.SessionToken)
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		return &resp, nil
	}

	acc := account.GetExternalInterface()

	data := acc.Update(ctx, req, sess)
	body, err := MarshalFunc(data.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while to trying to marshal the response body of the update account API: " + err.Error()
		l.LogWithFields(ctx).Printf(resp.StatusMessage)
		return &resp, nil
	}
	l.LogWithFields(ctx).Debugf("outgoing response of request to update the account: %s", string(body))
	resp.Body = body
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header

	return &resp, nil
}

// Delete defines the operations which handles the RPC request response
// for the delete of a particular account service of account-session micro service.
// The functionality retrieves the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Account) Delete(ctx context.Context, req *accountproto.DeleteAccountRequest) (*accountproto.AccountResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside Delete function (svc-account-session)")
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
	l.LogWithFields(ctx).Info("Validating session and updating the last used time of the session before deleting the account")
	sess, errs := CheckSessionTimeOutFunc(ctx, req.SessionToken)
	if errs != nil {
		resp.Body, resp.StatusCode, resp.StatusMessage = validateSessionTimeoutError(ctx, req.SessionToken, errs)
		return &resp, nil
	}

	err := UpdateLastUsedTimeFunc(ctx, req.SessionToken)
	if err != nil {
		errorArgs[0].ErrorMessage, resp.StatusCode, resp.StatusMessage = validateUpdateLastUsedTimeError(ctx, err, req.SessionToken)
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		return &resp, nil
	}

	data := AccDeleteFunc(ctx, sess, req.AccountID)
	var jsonErr error // jsonErr is created to protect the data in err
	body, jsonErr := MarshalFunc(data.Body)
	if jsonErr != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying to marshal the response body of the delete account API: " + jsonErr.Error()
		l.LogWithFields(ctx).Error(resp.StatusMessage)
		return &resp, nil
	}
	l.LogWithFields(ctx).Debugf("outgoing response of request to delete the account: %s", string(body))
	resp.Body = body
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header

	return &resp, nil
}

func validateSessionTimeoutError(ctx context.Context, sessionToken string, errs *errors.Error) (body []byte, statusCode int32, statusMessage string) {
	errorMessage := "error while authorizing session token: " + errs.Error()
	statusCode, statusMessage = errs.GetAuthStatusCodeAndMessage()
	if statusCode == http.StatusServiceUnavailable {
		body, _ = json.Marshal(common.GeneralError(statusCode, statusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body)
		l.LogWithFields(ctx).Error(errorMessage)
	} else {
		body, _ = json.Marshal(common.GeneralError(statusCode, statusMessage, errorMessage, nil, nil).Body)
		auth.CustomAuthLog(ctx, sessionToken, "Invalid session token", statusCode)
	}
	return
}

func validateUpdateLastUsedTimeError(ctx context.Context, err error, sessionToken string) (errorMessage string, statusCode int32, statusMessage string) {
	errorMessage = "error while updating last used time of session with token " + sessionToken + ": " + err.Error()
	statusCode = http.StatusInternalServerError
	statusMessage = response.InternalError
	l.LogWithFields(ctx).Error(errorMessage)
	return
}
