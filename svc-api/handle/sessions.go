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

//Package handle ...
package handle

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// SessionRPCs defines all the RPC methods in session service
type SessionRPCs struct {
	CreateSessionRPC        func(sessionproto.SessionCreateRequest) (*sessionproto.SessionCreateResponse, error)
	DeleteSessionRPC        func(string, string) (*sessionproto.SessionResponse, error)
	GetSessionRPC           func(string, string) (*sessionproto.SessionResponse, error)
	GetAllActiveSessionsRPC func(string, string) (*sessionproto.SessionResponse, error)
	GetSessionServiceRPC    func() (*sessionproto.SessionResponse, error)
}

// CreateSession defines the Create session iris handler
// This method extracts the user name and password
// create a rpc request and send a request to session micro service
// and feed the response to iris
func (s *SessionRPCs) CreateSession(ctx iris.Context) {
	defer ctx.Next()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the session create request body: %v" + err.Error()
		log.Printf(errorMessage)
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

	resp, err := s.CreateSessionRPC(createRequest)

	if err != nil && resp == nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	if resp.StatusCode == http.StatusCreated {
		resp.Header["Location"] = ctx.Request().Host + "/redfish/v1/SessionService/Sessions/" + resp.SessionId
	}

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

	resp, err := s.DeleteSessionRPC(sessionID, sessionToken)
	if err != nil && resp == nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

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
	sessionID := ctx.Params().Get("sessionID")
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")

	resp, err := s.GetSessionRPC(sessionID, sessionToken)
	if err != nil && resp == nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	if resp.StatusCode == http.StatusOK {
		resp.Header["Location"] = ctx.Request().Host + "/redfish/v1/SessionService/Sessions/" + sessionID
	}
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
	sessionID := ctx.Params().Get("sessionID")
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")

	resp, err := s.GetAllActiveSessionsRPC(sessionID, sessionToken)
	if err != nil && resp == nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetSessionService will do the rpc call to get session service
func (s *SessionRPCs) GetSessionService(ctx iris.Context) {
	defer ctx.Next()
	resp, err := s.GetSessionServiceRPC()
	if err != nil && resp == nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}
