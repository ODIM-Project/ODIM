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
	"log"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
	"github.com/ODIM-Project/ODIM/svc-account-session/role"
	"github.com/ODIM-Project/ODIM/svc-account-session/session"
)

// Role struct helps to register service
type Role struct {
}

//CreateRole defines the operations which handles the RPC request response
// for the create role of account-session micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (r *Role) CreateRole(ctx context.Context, req *roleproto.RoleRequest) (*roleproto.RoleResponse, error) {
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

	// Validating the session
	sess, errs := auth.CheckSessionTimeOut(req.SessionToken)
	if errs != nil {
		errorMessage := "error while authorizing session token: " + errs.Error()
		resp.StatusCode, resp.StatusMessage = errs.GetAuthStatusCodeAndMessage()
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body)
		} else {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body)
		}
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		}
		log.Printf(errorMessage)
		return &resp, nil
	}

	err := session.UpdateLastUsedTime(req.SessionToken)
	if err != nil {
		errorMessage := "error while updating last used time of session with token " + req.SessionToken + ": " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		}
		log.Printf(errorMessage)
		return &resp, nil
	}

	data := role.Create(req, sess)
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header
	resp.Body, err = json.Marshal(data.Body)
	if err != nil {
		errorMessage := "error while trying marshal the response body for get role: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		log.Printf(resp.StatusMessage)
		return &resp, nil
	}
	return &resp, nil
}

//GetRole defines the operations which handles the RPC request response
// for the view of a role of account-session micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (r *Role) GetRole(ctx context.Context, req *roleproto.GetRoleRequest) (*roleproto.RoleResponse, error) {
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

	// Validating the session
	sess, errs := auth.CheckSessionTimeOut(req.SessionToken)
	if errs != nil {
		errorMessage := "error while authorizing session token: " + errs.Error()
		resp.StatusCode, resp.StatusMessage = errs.GetAuthStatusCodeAndMessage()
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body)
		} else {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body)
		}
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		}
		log.Printf(errorMessage)
		return &resp, nil
	}

	err := session.UpdateLastUsedTime(req.SessionToken)
	if err != nil {
		errorMessage := "error while updating last used time of session with token " + req.SessionToken + ": " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		}
		log.Printf(errorMessage)
		return &resp, nil
	}

	data := role.GetRole(req, sess)
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header
	resp.Body, err = json.Marshal(data.Body)
	if err != nil {
		errorMessage := "error while trying marshal the response body for get role: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		log.Printf(resp.StatusMessage)
		return &resp, nil
	}

	return &resp, nil
}

//GetAllRoles defines the operations which handles the RPC request response
// for the list all roles  of account-session micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (r *Role) GetAllRoles(ctx context.Context, req *roleproto.GetRoleRequest) (*roleproto.RoleResponse, error) {
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

	sess, errs := auth.CheckSessionTimeOut(req.SessionToken)
	if errs != nil {
		errorMessage := "error while authorizing session token: " + errs.Error()
		resp.StatusCode, resp.StatusMessage = errs.GetAuthStatusCodeAndMessage()
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body)
		} else {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body)
		}
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		}
		log.Printf(errorMessage)
		return &resp, nil
	}

	err := session.UpdateLastUsedTime(req.SessionToken)
	if err != nil {
		errorMessage := "error while updating last used time of session with token " + req.SessionToken + ": " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		}
		log.Printf(errorMessage)
		return &resp, nil
	}

	data := role.GetAllRoles(sess)
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header
	resp.Body, err = json.Marshal(data.Body)
	if err != nil {
		errorMessage := "error while trying marshal the response body for get role: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		log.Printf(resp.StatusMessage)
		return &resp, nil
	}

	return &resp, nil
}

//UpdateRole defines the operations which handles the RPC request response
// for the update of a particular role  of account-session micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (r *Role) UpdateRole(ctx context.Context, req *roleproto.UpdateRoleRequest) (*roleproto.RoleResponse, error) {
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

	// Validating the session
	sess, errs := auth.CheckSessionTimeOut(req.SessionToken)
	if errs != nil {
		errorMessage := "error while authorizing session token: " + errs.Error()
		resp.StatusCode, resp.StatusMessage = errs.GetAuthStatusCodeAndMessage()
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body)
		} else {
			resp.Body, _ = json.Marshal(common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body)
		}
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		}
		log.Printf(errorMessage)
		return &resp, nil
	}

	err := session.UpdateLastUsedTime(req.SessionToken)
	if err != nil {
		errorMessage := "error while updating last used time of session with token " + req.SessionToken + ": " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		}
		log.Printf(errorMessage)
		return &resp, nil
	}

	data := role.Update(req, sess)
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header
	resp.Body, err = json.Marshal(data.Body)
	if err != nil {
		errorMessage := "error while trying marshal the response body for get role: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		log.Printf(resp.StatusMessage)
		return &resp, nil
	}

	return &resp, nil
}

// DeleteRole handles the RPC call from the client
func (r *Role) DeleteRole(ctx context.Context, req *roleproto.DeleteRoleRequest) (*roleproto.RoleResponse, error) {
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
	data := role.Delete(req)
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header
	var err error
	resp.Body, err = json.Marshal(data.Body)
	if err != nil {
		errorMessage := "error while trying marshal the response body for get role: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body, _ = json.Marshal(args.CreateGenericErrorResponse())
		log.Printf(resp.StatusMessage)
		return &resp, nil
	}

	return &resp, nil
}
