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

// Package handle ...
package handle

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// SessionRPCs defines all the RPC methods in session service
type SessionRPCs struct {
	CreateSessionRPC        func(context.Context, sessionproto.SessionCreateRequest) (*sessionproto.SessionCreateResponse, error)
	DeleteSessionRPC        func(context.Context, string, string) (*sessionproto.SessionResponse, error)
	GetSessionRPC           func(context.Context, string, string) (*sessionproto.SessionResponse, error)
	GetAllActiveSessionsRPC func(context.Context, string, string) (*sessionproto.SessionResponse, error)
	GetSessionServiceRPC    func(context.Context) (*sessionproto.SessionResponse, error)
}

// CreateSession defines the Create session iris handler
// This method extracts the user name and password
// create a rpc request and send a request to session micro service
// and feed the response to iris
func (s *SessionRPCs) CreateSession(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req interface{}
	err := ctx.ReadJSON(&req)
	l.LogWithFields(ctxt).Debug("Incoming request received for creating a session")
	if err != nil {
		errorMessage := "error while trying to get JSON body from the session create request body: %v" + err.Error()
		l.LogWithFields(ctxt).Printf(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}

	// Marshalling the req to make session request
	// Since session create request accepts []byte stream
	request, err := json.Marshal(req)
	createRequest := sessionproto.SessionCreateRequest{
		RequestBody: request,
	}

	resp, err := s.CreateSessionRPC(ctxt, createRequest)
	if err != nil && resp == nil {
		if strings.Contains(err.Error(), "too many requests") {
			response := common.GeneralError(http.StatusServiceUnavailable, response.SessionLimitExceeded, err.Error(), nil, nil)
			common.SetResponseHeader(ctx, response.Header)
			ctx.StatusCode(http.StatusServiceUnavailable)
			ctx.JSON(&response.Body)
			return
		}
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	if resp.StatusCode == http.StatusCreated {
		resp.Header["Location"] = "/redfish/v1/SessionService/Sessions/" + resp.SessionId
	}
	l.LogWithFields(ctxt).Debugf("response code for creating a session is %d", resp.StatusCode)
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// DeleteSession defines the Delete session iris handler
// This method extracts the sessionid and sessiontoken
// create a rpc request and send a request to session micro service
// and feed the response to iris
func (s *SessionRPCs) DeleteSession(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	sessionID := ctx.Params().Get("sessionID")
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	if sessionToken == "" {
		errorMessage := "error: session token is missing"
		response := common.GeneralError(http.StatusInternalServerError, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response)
		return
	}
	l.LogWithFields(ctxt).Debug("Incoming request received for deleting a session")
	resp, err := s.DeleteSessionRPC(ctxt, sessionID, sessionToken)
	if err != nil && resp == nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("response code for deleting a session is %d", resp.StatusCode)
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetSession  defines the Get session iris handler
// This method extracts the sessionid and sessiontoken
// create a rpc request and send a request to session micro service
// and feed the response to iris
func (s *SessionRPCs) GetSession(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	sessionID := ctx.Params().Get("sessionID")
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	l.LogWithFields(ctxt).Debug("Incoming request received for getting a session")
	resp, err := s.GetSessionRPC(ctxt, sessionID, sessionToken)
	if err != nil && resp == nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	if resp.StatusCode == http.StatusOK {
		resp.Header["Location"] = ctx.Request().Host + "/redfish/v1/SessionService/Sessions/" + sessionID
	}
	l.LogWithFields(ctxt).Debugf("fresponse code for getting a session details is %d", resp.StatusCode)
	ctx.ResponseWriter().Header().Set("Allow", "GET, DELETE")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetAllActiveSessions defines the GetAllActiveSession iris handler
// This method extracts the sessionid and sessiontoken
// create a rpc request and send a request to session micro service
// and feed the response to iris
func (s *SessionRPCs) GetAllActiveSessions(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	sessionID := ctx.Params().Get("sessionID")
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	l.LogWithFields(ctxt).Debug("Incoming request received for getting all active sessions")
	resp, err := s.GetAllActiveSessionsRPC(ctxt, sessionID, sessionToken)
	if err != nil && resp == nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("response code for getting all active sessions is %d", resp.StatusCode)
	ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetSessionService will do the rpc call to get session service
func (s *SessionRPCs) GetSessionService(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	l.LogWithFields(ctxt).Debug("Incoming request received for getting a session service details")
	resp, err := s.GetSessionServiceRPC(ctxt)
	if err != nil && resp == nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting a session service is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}
