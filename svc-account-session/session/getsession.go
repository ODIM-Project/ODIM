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
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
)

//GetSessionUserName is a RPC handle to get the session username from the session Token
func GetSessionUserName(req *sessionproto.SessionRequest) (*sessionproto.SessionUserName, error) {
	var resp sessionproto.SessionUserName
	resp.UserName = ""
	// Validating the session
	currentSession, err := auth.CheckSessionTimeOut(req.SessionToken)
	if err != nil {
		errorMessage := "Unable to authorize session token: " + err.Error()
		log.Error(errorMessage)
		return &resp, err
	}

	if errs := UpdateLastUsedTime(req.SessionToken); errs != nil {
		errorMessage := "Unable to update last used time of session matching token " + req.SessionToken + ": " + errs.Error()
		log.Error(errorMessage)
		return &resp, errs
	}
	resp.UserName = currentSession.UserName
	return &resp, nil
}

// GetSession is a method to get session
// it will accepts the SessionCreateRequest which will have sessionid and sessiontoken
// and it will check privileges to get session and then get the session against the sessionID
// respond RPC response and error if there is.
func GetSession(req *sessionproto.SessionRequest) response.RPC {
	commonResponse := response.Response{
		OdataType: common.SessionType,
		OdataID:   "/redfish/v1/SessionService/Sessions/" + req.SessionId,
		ID:        req.SessionId,
		Name:      "User Session",
	}
	var resp response.RPC

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
	currentSession, err := auth.CheckSessionTimeOut(req.SessionToken)
	if err != nil {
		errorMessage := "Unable to authorize session token: " + err.Error()
		resp.StatusCode, resp.StatusMessage = err.GetAuthStatusCodeAndMessage()
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body = common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body
		} else {
			resp.Body = common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body
		}
		resp.Header = getHeader()
		log.Error(errorMessage)
		return resp
	}

	if errs := UpdateLastUsedTime(req.SessionToken); errs != nil {
		errorMessage := "Unable to update last used time of session matching token " + req.SessionToken + ": " + errs.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		resp.Header = getHeader()
		log.Error(errorMessage)
		return resp
	}

	sessionTokens, errs := asmodel.GetAllSessionKeys()
	if errs != nil {
		errorMessage := "Unable to get all session keys while deleting session: " + errs.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		resp.Header = getHeader()
		log.Error(errorMessage)
		return resp
	}
	for _, token := range sessionTokens {
		session, err := auth.CheckSessionTimeOut(token)
		if err != nil {
			log.Error("Unable to get session details with the token " + token + ": " + err.Error())
			continue
		}
		if session.ID == req.SessionId {
			if checkPrivilege(req.SessionToken, session, currentSession) {

				resp.StatusCode = http.StatusOK
				resp.StatusMessage = response.Success
				resp.Header = map[string]string{
					"Cache-Control":     "no-cache",
					"Link":              "</redfish/v1/SessionService/Sessions/" + req.SessionId + "/>; rel=self",
					"Transfer-Encoding": "chunked",
					"X-Auth-Token":      token,
					"Content-type":      "application/json; charset=utf-8",
				}

				respBody := asresponse.Session{
					Response: commonResponse,
					UserName: session.UserName,
				}

				resp.Body = respBody
				return resp
			}
			errorMessage := "The session doesn't have the requisite privileges for the action"
			resp.StatusCode = http.StatusForbidden
			resp.StatusMessage = response.InsufficientPrivilege
			errorArgs[0].ErrorMessage = errorMessage
			errorArgs[0].StatusMessage = resp.StatusMessage
			resp.Body = args.CreateGenericErrorResponse()
			resp.Header = getHeader()
			log.Error(errorMessage)
			return resp
		}
	}
	sessionTokens = nil
	errorMessage := "No session with id " + req.SessionId + " found."
	log.Error("Status Not Found")
	resp.StatusCode = http.StatusNotFound
	resp.StatusMessage = response.ResourceNotFound
	errorArgs[0].ErrorMessage = errorMessage
	errorArgs[0].StatusMessage = resp.StatusMessage
	errorArgs[0].MessageArgs = []interface{}{"Session", req.SessionId}
	resp.Body = args.CreateGenericErrorResponse()
	resp.Header = getHeader()
	return resp
}

// GetAllActiveSessions is a method to get session
// it will accepts the SessionCreateRequest which will have sessionid and sessiontoken
// and it will check privileges to get session and then get all the active sessions
// respond RPC response and error if there is.
func GetAllActiveSessions(req *sessionproto.SessionRequest) response.RPC {

	commonResponse := response.Response{
		OdataType:    "#SessionCollection.SessionCollection",
		OdataID:      "/redfish/v1/SessionService/Sessions",
		OdataContext: "/redfish/v1/$metadata#SessionCollection.SessionCollection",
		Name:         "Session Service",
	}

	var resp response.RPC
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
	currentSession, gerr := auth.CheckSessionTimeOut(req.SessionToken)
	if gerr != nil {
		errorMessage := "Unable to authorize session token: " + gerr.Error()
		resp.StatusCode, resp.StatusMessage = gerr.GetAuthStatusCodeAndMessage()
		if resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body = common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil).Body
		} else {
			resp.Body = common.GeneralError(resp.StatusCode, resp.StatusMessage, errorMessage, nil, nil).Body
		}
		resp.Header = getHeader()
		log.Error(errorMessage)
		return resp
	}

	err := UpdateLastUsedTime(req.SessionToken)
	if err != nil {
		errorMessage := "Unable to update last used time of session with token " + req.SessionToken + ": " + err.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		resp.Header = getHeader()
		log.Error(errorMessage)
		return resp
	}

	if !currentSession.Privileges[common.PrivilegeConfigureSelf] && !currentSession.Privileges[common.PrivilegeConfigureUsers] {
		errorMessage := "Insufficient privileges: " + err.Error()
		resp.StatusCode = http.StatusForbidden
		resp.StatusMessage = response.InsufficientPrivilege
		errorArgs[0].ErrorMessage = errorMessage
		errorArgs[0].StatusMessage = resp.StatusMessage
		resp.Body = args.CreateGenericErrorResponse()
		resp.Header = getHeader()
		log.Error(errorMessage)
		return resp
	}

	sessionTokens, errs := asmodel.GetAllSessionKeys()
	if errs != nil {
		errorMessage := "Unable to get all session keys in delete session: " + errs.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		resp.Header = getHeader()
		log.Error(errorMessage)
		return resp
	}

	var listMembers []asresponse.ListMember
	for _, token := range sessionTokens {
		session, err := auth.CheckSessionTimeOut(token)
		if err != nil {
			log.Error("Unable to get session details with the token " + token + ": " + err.Error())
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
	resp.Header = getHeader()
	resp.Body = respBody
	return resp
}

// GetSessionService is a method to get session
// it will accepts the SessionCreateRequest which will have sessionid and sessiontoken
// and it will check the session service is enabled or not from the config file.
// respond RPC response and error if there are.
func GetSessionService(req *sessionproto.SessionRequest) response.RPC {
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
		"Allow":         "GET",
		"Cache-Control": "no-cache",
		"Connection":    "Keep-alive",
		"Link": "	</redfish/v1/SchemaStore/en/SessionService.json>; rel=describedby",
		"Transfer-Encoding": "chunked",
		"X-Frame-Options":   "sameorigin",
		"Content-type":      "application/json; charset=utf-8",
	}

	resp.Body = sessionService
	return resp
}
