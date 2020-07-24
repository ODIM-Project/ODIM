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
	"github.com/bharath-b-hpe/odimra/lib-utilities/common"
	"github.com/bharath-b-hpe/odimra/lib-utilities/config"
	"github.com/bharath-b-hpe/odimra/lib-utilities/errors"
	roleproto "github.com/bharath-b-hpe/odimra/lib-utilities/proto/role"
	"github.com/bharath-b-hpe/odimra/lib-utilities/response"
	"github.com/bharath-b-hpe/odimra/svc-account-session/asmodel"
	"github.com/bharath-b-hpe/odimra/svc-account-session/auth"
	"github.com/bharath-b-hpe/odimra/svc-account-session/session"

	"log"
	"net/http"
)

func doSessionAuthAndUpdate(resp *response.RPC, sessionToken string) (*asmodel.Session, error) {
	sess, err := auth.CheckSessionTimeOut(sessionToken)
	if err != nil {
		errorMessage := "error while authorizing session token: " + err.Error()
		resp.StatusCode, resp.StatusMessage = err.GetAuthStatusCodeAndMessage()
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body = common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body
		} else {
			resp.Body = common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body
		}
		resp.Header = map[string]string{
			"Content-type":      "application/json; charset=utf-8", // TODO: add all error headers
			"Cache-Control":     "no-cache",
			"Connection":        "keep-alive",
			"Transfer-Encoding": "chunked",
			"OData-Version":     "4.0",
			"X-Frame-Options":   "sameorigin",
		}
		log.Printf(errorMessage)
		return nil, err
	}
	if errs := session.UpdateLastUsedTime(sessionToken); errs != nil {
		errorMessage := "error while updating last used time of session with token " + sessionToken + ": " + errs.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		}
		log.Printf(errorMessage)
		return nil, errs
	}
	return sess, nil
}

// Delete defines the functionality of deletion of non predefined roles
func Delete(req *roleproto.DeleteRoleRequest) *response.RPC {
	var resp response.RPC
	sess, err := doSessionAuthAndUpdate(&resp, req.SessionToken)
	if err != nil {
		return &resp
	}
	/* Populate generic headers */
	resp.Header = map[string]string{
		"Content-type":      "application/json; charset=utf-8", // TODO: add all error headers
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
		"X-Frame-Options":   "sameorigin",
	}
	if !sess.Privileges[common.PrivilegeConfigureUsers] {
		errorMessage := "error: the session token didn't have required privilege"
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
		log.Printf(errorMessage)
		return &resp
	}
	users, uerr := asmodel.GetAllUsers()
	if uerr != nil {
		errorMessage := "error: unable to get users list: " + uerr.Error()
		log.Printf(errorMessage)
		resp.CreateInternalErrorResponse(errorMessage)
		return &resp
	}
	for _, key := range users {
		if req.ID == key.RoleID {
			errorMessage := "error: role is assigned to a user"
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
			log.Printf(errorMessage)
			return &resp
		}
	}
	role, gerr := asmodel.GetRoleDetailsByID(req.ID)
	if gerr != nil {
		errorMessage := "error while trying to get role details: " + gerr.Error()
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
		log.Printf(errorMessage)
		return &resp
	}
	if role.IsPredefined {
		errorMessage := "error: the predifined roles cannot be deleted."
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
		log.Printf(errorMessage)
		return &resp
	}

	if derr := role.Delete(); derr != nil {
		errorMessage := "error while trying to delete role: " + derr.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		log.Printf(errorMessage)
		return &resp
	}

	resp.StatusCode = http.StatusNoContent
	resp.StatusMessage = response.ResourceRemoved

	return &resp
}
