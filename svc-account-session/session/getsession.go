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
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
)

// GetSessionUserName is a RPC handle to get the session username from the session Token
func GetSessionUserName(ctx context.Context, req *sessionproto.SessionRequest) (*sessionproto.SessionUserName, error) {
	var resp sessionproto.SessionUserName
	resp.UserName = ""
	// Validating the session
	currentSession, err := auth.CheckSessionTimeOut(ctx, req.SessionToken)
	if err != nil {
		return &resp, err
	}

	if errs := UpdateLastUsedTime(ctx, req.SessionToken); errs != nil {
		errorMessage := "Unable to update last used time of session matching token " + req.SessionToken + ": " + errs.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		return &resp, errs
	}
	resp.UserName = currentSession.UserName
	l.LogWithFields(ctx).Debugf("outgoing response of request to get session username: %s", currentSession.UserName)
	return &resp, nil
}

// GetSessionUserRoleID is a RPC handle to get the session user's role id from the session Token
func GetSessionUserRoleID(ctx context.Context, req *sessionproto.SessionRequest) (*sessionproto.SessionUsersRoleID, error) {
	var resp sessionproto.SessionUsersRoleID
	resp.RoleID = ""
	// Validating the session
	currentSession, err := auth.CheckSessionTimeOut(ctx, req.SessionToken)
	if err != nil {
		return &resp, err
	}

	if errs := UpdateLastUsedTime(ctx, req.SessionToken); errs != nil {
		errorMessage := "Unable to update last used time of session matching token " + req.SessionToken + ": " + errs.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		return &resp, errs
	}
	resp.RoleID = currentSession.RoleID
	l.LogWithFields(ctx).Debugf("outgoing response of request to get session role id: %s", currentSession.RoleID)
	return &resp, nil
}

// GetSession is a method to get session
// it will accepts the SessionCreateRequest which will have sessionid and sessiontoken
// and it will check privileges to get session and then get the session against the sessionID
// respond RPC response and error if there is.
func GetSession(ctx context.Context, req *sessionproto.SessionRequest) response.RPC {
	commonResponse := response.Response{
		OdataType: common.SessionType,
		OdataID:   "/redfish/v1/SessionService/Sessions/" + req.SessionId,
		ID:        req.SessionId,
		Name:      "User Session",
	}
	var resp response.RPC
	errLogPrefix := "failed to fetch the session : "

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

	l.LogWithFields(ctx).Info("Validating the request to fetch the session")
	// Validating the session
	currentSession, err := auth.CheckSessionTimeOut(ctx, req.SessionToken)
	if err != nil {
		errorMessage := errLogPrefix + "Unable to authorize session token: " + err.Error()
		resp.StatusCode, resp.StatusMessage = err.GetAuthStatusCodeAndMessage()
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body = common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body
			l.LogWithFields(ctx).Error(errorMessage)
		} else {
			resp.Body = common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body
			auth.CustomAuthLog(ctx, req.SessionToken, "Invalid session token", resp.StatusCode)
		}
		return resp
	}

	if errs := UpdateLastUsedTime(ctx, req.SessionToken); errs != nil {
		errorMessage := errLogPrefix + "Unable to update last used time of session matching token " + req.SessionToken + ": " + errs.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		l.LogWithFields(ctx).Error(errorMessage)
		return resp
	}
	auth.CustomAuthLog(ctx, req.SessionToken, "Authorization is successful", http.StatusOK)
	sessionTokens, errs := asmodel.GetAllSessionKeys()
	if errs != nil {
		errorMessage := errLogPrefix + "Unable to get all session keys while deleting session: " + errs.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		l.LogWithFields(ctx).Error(errorMessage)
		return resp
	}
	for _, token := range sessionTokens {
		session, err := auth.CheckSessionTimeOut(ctx, token)
		if err != nil {
			auth.CustomAuthLog(ctx, req.SessionToken, "Invalid session token", resp.StatusCode)
			continue
		}
		if session.ID == req.SessionId {
			if checkPrivilege(req.SessionToken, session, currentSession) {

				resp.StatusCode = http.StatusOK
				resp.StatusMessage = response.Success
				resp.Header = map[string]string{
					"Link":         "</redfish/v1/SessionService/Sessions/" + req.SessionId + "/>; rel=self",
					"X-Auth-Token": token,
				}

				respBody := asresponse.Session{
					Response:    commonResponse,
					UserName:    session.UserName,
					CreatedTime: session.CreatedTime.Format(time.RFC3339),
				}

				resp.Body = respBody
				return resp
			}
			errorMessage := errLogPrefix + "The session doesn't have the requisite privileges for the action"
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
	errorMessage := "No session with id " + req.SessionId + " found."
	l.LogWithFields(ctx).Error(errLogPrefix + errorMessage)
	resp.StatusCode = http.StatusNotFound
	resp.StatusMessage = response.ResourceNotFound
	errorArgs[0].ErrorMessage = errorMessage
	errorArgs[0].StatusMessage = resp.StatusMessage
	errorArgs[0].MessageArgs = []interface{}{"Session", req.SessionId}
	resp.Body = args.CreateGenericErrorResponse()
	return resp
}

// GetAllActiveSessions is a method to get session
// it will accepts the SessionCreateRequest which will have sessionid and sessiontoken
// and it will check privileges to get session and then get all the active sessions
// respond RPC response and error if there is.
func GetAllActiveSessions(ctx context.Context, req *sessionproto.SessionRequest) response.RPC {

	commonResponse := response.Response{
		OdataType:    "#SessionCollection.SessionCollection",
		OdataID:      "/redfish/v1/SessionService/Sessions",
		OdataContext: "/redfish/v1/$metadata#SessionCollection.SessionCollection",
		Name:         "Session Service",
	}

	var resp response.RPC
	errorLogPrefix := "failed to fetch all active sessions : "
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

	l.LogWithFields(ctx).Info("fetching all active sessions")
	// Validating the session
	currentSession, gerr := auth.CheckSessionTimeOut(ctx, req.SessionToken)
	if gerr != nil {
		errorMessage := errorLogPrefix + "Unable to authorize session token: " + gerr.Error()
		resp.StatusCode, resp.StatusMessage = gerr.GetAuthStatusCodeAndMessage()
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body = common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body
			l.LogWithFields(ctx).Error(errorMessage)
		} else {
			resp.Body = common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body
			auth.CustomAuthLog(ctx, req.SessionToken, errorMessage, resp.StatusCode)
		}
		return resp
	}

	err := UpdateLastUsedTime(ctx, req.SessionToken)
	if err != nil {
		errorMessage := errorLogPrefix + "Unable to update last used time of session with token " + req.SessionToken + ": " + err.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		l.LogWithFields(ctx).Error(errorMessage)
		return resp
	}
	auth.CustomAuthLog(ctx, req.SessionToken, "Authorization is successful", http.StatusOK)
	if !currentSession.Privileges[common.PrivilegeConfigureSelf] && !currentSession.Privileges[common.PrivilegeConfigureUsers] {
		errorMessage := errorLogPrefix + "Insufficient privileges: " + err.Error()
		resp.StatusCode = http.StatusForbidden
		resp.StatusMessage = response.InsufficientPrivilege
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body = args.CreateGenericErrorResponse()
		auth.CustomAuthLog(ctx, req.SessionToken, errorMessage, resp.StatusCode)
		return resp
	}

	sessionTokens, errs := asmodel.GetAllSessionKeys()
	if errs != nil {
		errorMessage := errorLogPrefix + "Unable to get all session keys : " + errs.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		l.LogWithFields(ctx).Error(errorMessage)
		return resp
	}

	var listMembers []asresponse.ListMember
	for _, token := range sessionTokens {
		session, err := auth.CheckSessionTimeOut(ctx, token)
		if err != nil {
			l.LogWithFields(ctx).Error(errorLogPrefix + "Unable to get session details with the token " + token + ": " + err.Error())
			continue
		}

		if checkPrivilege(req.SessionToken, session, currentSession) {
			member := asresponse.ListMember{
				OdataID: "/redfish/v1/SessionService/Sessions/" + session.ID,
			}
			listMembers = append(listMembers, member)
		}
	}
	sessionTokens = nil
	respBody := asresponse.List{
		Response:     commonResponse,
		MembersCount: len(listMembers),
		Members:      listMembers,
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	resp.Body = respBody
	return resp
}

// GetSessionService is a method to get session
// it will accepts the SessionCreateRequest which will have sessionid and sessiontoken
// and it will check the session service is enabled or not from the config file.
// respond RPC response and error if there are.
func GetSessionService(ctx context.Context, req *sessionproto.SessionRequest) response.RPC {
	commonResponse := response.Response{
		OdataType: common.SessionServiceType,
		OdataID:   "/redfish/v1/SessionService",
		ID:        "Sessions",
		Name:      "Session Service",
	}
	var resp response.RPC
	IsServiceEnabled := false
	ServiceState := "Disabled"
	//Checks if SessionService is enabled and sets the variable IsServiceEnabled to true abd ServicState to enabled
	for _, service := range config.Data.EnabledServices {
		if service == "SessionService" {
			IsServiceEnabled = true
			ServiceState = "Enabled"
		}
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	commonResponse.CreateGenericResponse(resp.StatusMessage)
	commonResponse.MessageID = ""
	commonResponse.Message = ""
	commonResponse.Severity = ""

	sessionService := asresponse.SessionService{
		Response: commonResponse,
		//TODO: Yet to implement SessionService state and health
		Status: asresponse.Status{
			State:  ServiceState,
			Health: "OK",
		},
		ServiceEnabled: IsServiceEnabled,
		SessionTimeout: config.Data.AuthConf.SessionTimeOutInMins,
		Sessions: asresponse.Sessions{
			OdataID: "/redfish/v1/SessionService/Sessions",
		},
	}
	resp.Header = map[string]string{
		"Link": "	</redfish/v1/SchemaStore/en/SessionService.json>; rel=describedby",
	}

	resp.Body = sessionService
	return resp
}
