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

// Package role ...
package role

import (
	"context"
	"fmt"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
	"github.com/ODIM-Project/ODIM/svc-account-session/session"

	"net/http"
)

func doSessionAuthAndUpdate(ctx context.Context, resp *response.RPC, sessionToken string) (*asmodel.Session, error) {
	sess, err := auth.CheckSessionTimeOut(ctx, sessionToken)
	if err != nil {
		errorMessage := "Unable to authorize session token: " + err.Error()
		resp.StatusCode, resp.StatusMessage = err.GetAuthStatusCodeAndMessage()
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body = common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body
			l.LogWithFields(ctx).Error(errorMessage)
		} else {
			resp.Body = common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body
			auth.CustomAuthLog(ctx, sessionToken, "Invalid session token", resp.StatusCode)
		}
		return nil, err
	}
	if errs := session.UpdateLastUsedTime(ctx, sessionToken); errs != nil {
		errorMessage := "Unable to update last used time of session" + ": " + errs.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		l.LogWithFields(ctx).Error(errorMessage)
		return nil, errs
	}
	return sess, nil
}

// Delete defines the functionality of deletion of non predefined roles
func Delete(ctx context.Context, req *roleproto.DeleteRoleRequest) *response.RPC {
	var resp response.RPC
	l.LogWithFields(ctx).Info("Validating session and updating the last used time of the session before deleting the role")
	sess, err := doSessionAuthAndUpdate(ctx, &resp, req.SessionToken)
	if err != nil {
		return &resp
	}

	errLogPrefix := fmt.Sprintf("failed to delete role %s: ", req.ID)
	l.LogWithFields(ctx).Infof("Validating the request to delete the role %s", req.ID)
	if !sess.Privileges[common.PrivilegeConfigureUsers] {
		errorMessage := errLogPrefix + "The session token doesn't have required privilege"
		resp.StatusCode = http.StatusForbidden
		resp.StatusMessage = response.InsufficientPrivilege
		args := response.Args{
			Code:    response.GeneralError,
			Message: "",
			ErrorArgs: []response.ErrArgs{
				response.ErrArgs{
					StatusMessage: response.InsufficientPrivilege,
					ErrorMessage:  errorMessage,
					MessageArgs:   []interface{}{},
				},
			},
		}
		resp.Body = args.CreateGenericErrorResponse()
		auth.CustomAuthLog(ctx, req.SessionToken, errorMessage, resp.StatusCode)
		return &resp
	}
	users, uerr := asmodel.GetAllUsers()
	if uerr != nil {
		errorMessage := errLogPrefix + "Unable to get users list: " + uerr.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		resp.CreateInternalErrorResponse(errorMessage)
		return &resp
	}
	for _, key := range users {
		if req.ID == key.RoleID {
			errorMessage := errLogPrefix + "Role is assigned to a user"
			resp.StatusCode = http.StatusForbidden
			resp.StatusMessage = response.ResourceInUse
			args := response.Args{
				Code:    response.GeneralError,
				Message: "",
				ErrorArgs: []response.ErrArgs{
					response.ErrArgs{
						StatusMessage: response.ResourceInUse,
						ErrorMessage:  errorMessage,
						MessageArgs:   []interface{}{},
					},
				},
			}
			resp.Body = args.CreateGenericErrorResponse()
			l.LogWithFields(ctx).Error(errorMessage)
			return &resp
		}
	}
	role, gerr := asmodel.GetRoleDetailsByID(req.ID)
	if gerr != nil {
		errorMessage := errLogPrefix + "Unable to get role details: " + gerr.Error()
		if errors.DBKeyNotFound == gerr.ErrNo() {
			resp.StatusCode = http.StatusNotFound
			resp.StatusMessage = response.ResourceNotFound
			messageArgs := []interface{}{"Role", req.ID}
			args := response.Args{
				Code:    response.GeneralError,
				Message: "",
				ErrorArgs: []response.ErrArgs{
					response.ErrArgs{
						StatusMessage: resp.StatusMessage,
						ErrorMessage:  errorMessage,
						MessageArgs:   messageArgs,
					},
				},
			}
			resp.Body = args.CreateGenericErrorResponse()
		} else {
			resp.CreateInternalErrorResponse(errorMessage)
		}
		l.LogWithFields(ctx).Error(errorMessage)
		return &resp
	}
	if role.IsPredefined {
		errorMessage := errLogPrefix + "A predefined role cannot be deleted."
		resp.StatusCode = http.StatusForbidden
		resp.StatusMessage = response.InsufficientPrivilege
		args := response.Args{
			Code:    response.GeneralError,
			Message: "",
			ErrorArgs: []response.ErrArgs{
				response.ErrArgs{
					StatusMessage: resp.StatusMessage,
					ErrorMessage:  errorMessage,
					MessageArgs:   []interface{}{},
				},
			},
		}
		resp.Body = args.CreateGenericErrorResponse()
		l.LogWithFields(ctx).Error(errorMessage)
		return &resp
	}

	l.LogWithFields(ctx).Infof("Deleting the role %s", req.ID)
	if derr := role.Delete(); derr != nil {
		errorMessage := errLogPrefix + derr.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		l.LogWithFields(ctx).Error(errorMessage)
		return &resp
	}

	resp.StatusCode = http.StatusNoContent
	resp.StatusMessage = response.ResourceRemoved

	return &resp
}
