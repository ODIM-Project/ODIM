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

// Package session ...
package session

import (
	"context"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
)

// DeleteSession is a method to delete a sessiom
// it will accepts the SessionCreateRequest which will have sessionid and sessiontoken
// and it will check privileges to delete session and then delete the session
// respond RPC response and error if there is.
func DeleteSession(ctx context.Context, req *sessionproto.SessionRequest) response.RPC {
	var resp response.RPC
	errorLogPrefix := "failed to delete session : "
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
	l.LogWithFields(ctx).Info("Validating the request to delete the session")
	currentSession, serr := asmodel.GetSession(req.SessionToken)
	if serr != nil {
		errorMessage := errorLogPrefix + serr.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
	}

	sessionTokens, err := asmodel.GetAllSessionKeys()
	if err != nil {
		errorMessage := errorLogPrefix + "Unable to get all session keys : " + err.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		l.LogWithFields(ctx).Error(errorMessage)
		return resp
	}
	for _, token := range sessionTokens {
		session, err := auth.CheckSessionTimeOut(ctx, token)
		if err != nil {
			l.LogWithFields(ctx).Error(errorLogPrefix + "Unable to get session details with the token " + token + ": " + err.Error())
			continue
		}
		if session.ID == req.SessionId {
			hasprivilege := checkPrivilege(req.SessionToken, session, &currentSession)
			if hasprivilege {
				if req.SessionToken != session.Token {
					err := UpdateLastUsedTime(ctx, req.SessionToken)
					if err != nil {
						errorMessage := errorLogPrefix + "Unable to update last used time of session matching token " + req.SessionToken + ": " + err.Error()
						resp.CreateInternalErrorResponse(errorMessage)
						l.LogWithFields(ctx).Error(errorMessage)
						return resp
					}
				}
				if err := session.Delete(); err != nil {
					errorMessage := errorLogPrefix + err.Error()
					resp.CreateInternalErrorResponse(errorMessage)
					l.LogWithFields(ctx).Error(errorMessage)
					return resp
				}
				resp.StatusCode = http.StatusNoContent
				resp.StatusMessage = response.ResourceRemoved
				l.LogWithFields(ctx).Info("Session is deleted")
				return resp
			}
			errorMessage := errorLogPrefix + "Insufficient privileges"
			resp.StatusCode = http.StatusForbidden
			resp.StatusMessage = response.InsufficientPrivilege
			errorArgs[0].ErrorMessage = errorMessage
			errorArgs[0].StatusMessage = resp.StatusMessage
			resp.Body = args.CreateGenericErrorResponse()
			auth.CustomAuthLog(ctx, req.SessionToken, errorMessage, resp.StatusCode)
			return resp
		}
	}
	sessionTokens = nil
	errorMessage := errorLogPrefix + "Session ID not found"
	l.LogWithFields(ctx).Error(errorMessage)
	resp.StatusCode = http.StatusNotFound
	resp.StatusMessage = response.ResourceNotFound
	errorArgs[0].ErrorMessage = errorMessage
	errorArgs[0].StatusMessage = resp.StatusMessage
	errorArgs[0].MessageArgs = []interface{}{"Session", req.SessionId}
	resp.Body = args.CreateGenericErrorResponse()
	return resp
}

func checkPrivilege(sessionToken string, session, currentSession *asmodel.Session) bool {
	if (session.UserName == currentSession.UserName && currentSession.Privileges[common.PrivilegeConfigureSelf]) ||
		currentSession.Privileges[common.PrivilegeConfigureUsers] {
		return true
	}
	return false
}
