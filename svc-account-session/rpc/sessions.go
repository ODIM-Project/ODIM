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

	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
	"github.com/ODIM-Project/ODIM/svc-account-session/session"

	"net/http"
)

// Session struct helps to register service
type Session struct{}

var (
	CreateNewSessionFunc     = session.CreateNewSession
	DeleteSessionFunc        = session.DeleteSession
	GetSessionFunc           = session.GetSession
	GetAllActiveSessionsFunc = session.GetAllActiveSessions
	GetSessionServiceFunc    = session.GetSessionService
	GetSessionUserNameFunc   = session.GetSessionUserName
	GetSessionUserRoleIDFunc = session.GetSessionUserRoleID
	MarshalFunc              = json.Marshal
)

// CreateSession is a rpc call to create session
// and It will check the credentials of user, if user is authorized
// then create session for the same
func (s *Session) CreateSession(ctx context.Context, req *sessionproto.SessionCreateRequest) (*sessionproto.SessionCreateResponse, error) {
	var err error
	var resp sessionproto.SessionCreateResponse
	response, sessionID := CreateNewSessionFunc(req)
	resp.Body, err = MarshalFunc(response.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying marshal the response body for create account: " + err.Error()
		l.Log.Printf(resp.StatusMessage)
		return &resp, nil
	}
	resp.SessionId = sessionID
	resp.StatusCode = response.StatusCode
	resp.StatusMessage = response.StatusMessage
	resp.Header = response.Header

	return &resp, nil
}

// DeleteSession is a rpc call to delete session
// It will get all the session tokens from the db and from the session token get the session details
// if session id is matched with recieved session id ten delete the session
func (s *Session) DeleteSession(ctx context.Context, req *sessionproto.SessionRequest) (*sessionproto.SessionResponse, error) {
	response := DeleteSessionFunc(req)
	var resp sessionproto.SessionResponse
	body, err := MarshalFunc(response.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying marshal the response body for delete : " + err.Error()
		l.Log.Printf(response.StatusMessage)
		return &resp, nil
	}
	resp.StatusCode = response.StatusCode
	resp.StatusMessage = response.StatusMessage
	resp.Header = response.Header
	resp.Body = body
	return &resp, nil
}

// GetSession is a rpc call to get session
// It will get all the session tokens from the db and from the session token get the session details
// if session id is matched with recieved session id then delete the session
func (s *Session) GetSession(ctx context.Context, req *sessionproto.SessionRequest) (*sessionproto.SessionResponse, error) {
	var resp sessionproto.SessionResponse
	response := GetSessionFunc(req)
	body, err := MarshalFunc(response.Body)
	if err != nil {
		resp.StatusMessage = "error while trying marshal the response body for get session: " + err.Error()
		l.Log.Printf(response.StatusMessage)
		return &resp, nil
	}
	resp.StatusCode = response.StatusCode
	resp.StatusMessage = response.StatusMessage
	resp.Header = response.Header
	resp.Body = body
	return &resp, nil
}

// GetSessionUserName is a rpc call to get session username
// It will get all the session username from the session
func (s *Session) GetSessionUserName(ctx context.Context, req *sessionproto.SessionRequest) (*sessionproto.SessionUserName, error) {
	resp, err := GetSessionUserNameFunc(req)
	return resp, err
}

// GetSessionUserRoleID is a rpc call to get session user's role ID
// It will get the session username's role id from the session
func (s *Session) GetSessionUserRoleID(ctx context.Context, req *sessionproto.SessionRequest) (*sessionproto.SessionUsersRoleID, error) {
	resp, err := GetSessionUserRoleIDFunc(req)
	return resp, err
}

// GetAllActiveSessions is a rpc call to get all active sessions
// This method will accepts the sessionrequest which has session id and session token
// and it will call GetAllActiveSessions from the session package
// and respond all the sessionresponse values along with error if there is.
func (s *Session) GetAllActiveSessions(ctx context.Context, req *sessionproto.SessionRequest) (*sessionproto.SessionResponse, error) {
	var resp sessionproto.SessionResponse
	response := GetAllActiveSessionsFunc(req)
	body, err := MarshalFunc(response.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying marshal the response body for get all active session: " + err.Error()
		l.Log.Printf(response.StatusMessage)
		return &resp, nil
	}
	resp.StatusCode = response.StatusCode
	resp.StatusMessage = response.StatusMessage
	resp.Header = response.Header
	resp.Body = body
	return &resp, nil

}

// GetSessionService is a rpc call to get session service
// which basically checks if the session service is enabled or not
func (s *Session) GetSessionService(ctx context.Context, req *sessionproto.SessionRequest) (*sessionproto.SessionResponse, error) {
	var resp sessionproto.SessionResponse
	response := GetSessionServiceFunc(req)
	body, err := MarshalFunc(response.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying marshal the response body for get session service: " + err.Error()
		l.Log.Printf(response.StatusMessage)
		return &resp, nil
	}
	resp.StatusCode = response.StatusCode
	resp.StatusMessage = response.StatusMessage
	resp.Header = response.Header
	resp.Body = body
	return &resp, nil

}

func getCommonResponse(statusMessage string) asresponse.RedfishSessionResponse {
	return asresponse.RedfishSessionResponse{
		Error: asresponse.Error{
			Code:    response.GeneralError,
			Message: "See @Message.ExtendedInfo for more information.",
			ExtendedInfos: []asresponse.ExtendedInfo{
				asresponse.ExtendedInfo{
					MessageID: statusMessage,
				},
			},
		},
	}
}
