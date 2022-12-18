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
//Package rpc defines the handler for micro services

// Package rpc ...
package rpc

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
	"github.com/ODIM-Project/ODIM/svc-account-session/role"
	"github.com/ODIM-Project/ODIM/svc-account-session/session"
)

// Role struct helps to register service
type Role struct {
}

var (
	CheckSessionTimeOutFunc = auth.CheckSessionTimeOut
	UpdateLastUsedTimeFunc  = session.UpdateLastUsedTime
	CreateFunc              = role.Create
	GetRoleFunc             = role.GetRole
	GetAllRolesFunc         = role.GetAllRoles
	DeleteFunc              = role.Delete
	UpdateFunc              = role.Update
)

// CreateRole defines the operations which handles the RPC request response
// for the create role of account-session micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (r *Role) CreateRole(ctx context.Context, req *roleproto.RoleRequest) (*roleproto.RoleResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside CreateRole function (svc-account-session)")
	var resp roleproto.RoleResponse
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

	l.LogWithFields(ctx).Info("Validating session and updating the last used time of the session before creating the role")
	// Validating the session
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

	data := CreateFunc(ctx, req, sess)
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header
	body, err := MarshalFunc(data.Body)
	if err != nil {
		errorMessage := "error while trying to marshal the response body of create role API: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		l.LogWithFields(ctx).Error(resp.StatusMessage)
		return &resp, nil
	}
	l.LogWithFields(ctx).Debugf("outgoing response of request to create a role: %s", string(body))
	resp.Body = body
	return &resp, nil
}

// GetRole defines the operations which handles the RPC request response
// for the view of a role of account-session micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (r *Role) GetRole(ctx context.Context, req *roleproto.GetRoleRequest) (*roleproto.RoleResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside GetRole function (svc-account-session)")
	var resp roleproto.RoleResponse
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

	l.LogWithFields(ctx).Info("Validating session and updating the last used time of the session before fetching the role details")
	// Validating the session
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

	data := GetRoleFunc(ctx, req, sess)
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header
	body, err := MarshalFunc(data.Body)
	if err != nil {
		errorMessage := "error while trying to marshal the response body of get role API: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		l.LogWithFields(ctx).Error(resp.StatusMessage)
		return &resp, nil
	}
	l.LogWithFields(ctx).Debugf("outgoing response of request to view role details: %s", string(body))
	resp.Body = body

	return &resp, nil
}

// GetAllRoles defines the operations which handles the RPC request response
// for the list all roles  of account-session micro service.
// The functionality retrieves the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (r *Role) GetAllRoles(ctx context.Context, req *roleproto.GetRoleRequest) (*roleproto.RoleResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside GetAllRoles function (svc-account-session)")
	var resp roleproto.RoleResponse
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

	l.LogWithFields(ctx).Info("Validating session and updating the last used time of the session before fetching all roles")
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

	data := GetAllRolesFunc(ctx, sess)
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header
	body, err := MarshalFunc(data.Body)
	if err != nil {
		errorMessage := "error while trying to marshal the response body of the get all roles API: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		l.LogWithFields(ctx).Error(resp.StatusMessage)
		return &resp, nil
	}
	l.LogWithFields(ctx).Debugf("outgoing response of request to view all roles: %s", string(body))
	resp.Body = body

	return &resp, nil
}

// UpdateRole defines the operations which handles the RPC request response
// for the update of a particular role  of account-session micro service.
// The functionality retrieves the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (r *Role) UpdateRole(ctx context.Context, req *roleproto.UpdateRoleRequest) (*roleproto.RoleResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside UpdateRole function (svc-account-session)")
	var resp roleproto.RoleResponse
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

	l.LogWithFields(ctx).Info("Validating session and updating the last used time of the session before updating the role")
	// Validating the session
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

	data := UpdateFunc(ctx, req, sess)
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header
	body, err := MarshalFunc(data.Body)
	if err != nil {
		errorMessage := "error while trying to marshal the response body of the update role API: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		l.LogWithFields(ctx).Error(resp.StatusMessage)
		return &resp, nil
	}
	l.LogWithFields(ctx).Debugf("outgoing response of request to update the role: %s", string(body))
	resp.Body = body

	return &resp, nil
}

// DeleteRole handles the RPC call from the client
func (r *Role) DeleteRole(ctx context.Context, req *roleproto.DeleteRoleRequest) (*roleproto.RoleResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside DeleteRole function (svc-account-session)")
	var resp roleproto.RoleResponse
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
	data := DeleteFunc(ctx, req)
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header
	var err error
	body, err := MarshalFunc(data.Body)
	if err != nil {
		errorMessage := "error while trying to marshal the response body of the delete role API: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		l.LogWithFields(ctx).Error(resp.StatusMessage)
		return &resp, nil
	}
	l.LogWithFields(ctx).Debugf("outgoing response of request to delete the role: %s", string(body))
	resp.Body = body

	return &resp, nil
}
