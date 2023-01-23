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

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
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
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside CreateSession function (svc-account-session)")
	var err error
	var resp sessionproto.SessionCreateResponse
	response, sessionID := CreateNewSessionFunc(ctx, req)
	body, err := MarshalFunc(response.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying to marshal the response body of the create session API: " + err.Error()
		l.LogWithFields(ctx).Printf(resp.StatusMessage)
		return &resp, nil
	}
	l.LogWithFields(ctx).Debugf("outgoing response of request to create the session: %s", string(body))
	resp.Body = body
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
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside DeleteSession function (svc-account-session)")
	response := DeleteSessionFunc(ctx, req)
	var resp sessionproto.SessionResponse
	body, err := MarshalFunc(response.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying to marshal the response body of the delete session API: " + err.Error()
		l.LogWithFields(ctx).Printf(response.StatusMessage)
		return &resp, nil
	}
	l.LogWithFields(ctx).Debugf("outgoing response of request to delete the session: %s", string(body))
	resp.Body = body
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
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside GetSession function (svc-account-session)")
	var resp sessionproto.SessionResponse
	response := GetSessionFunc(ctx, req)
	body, err := MarshalFunc(response.Body)
	if err != nil {
		resp.StatusMessage = "error while trying to marshal the response body of the get session API: " + err.Error()
		l.LogWithFields(ctx).Printf(response.StatusMessage)
		return &resp, nil
	}
	l.LogWithFields(ctx).Debugf("outgoing response of request to get the session: %s", string(body))
	resp.Body = body
	resp.StatusCode = response.StatusCode
	resp.StatusMessage = response.StatusMessage
	resp.Header = response.Header
	resp.Body = body
	return &resp, nil
}

// GetSessionUserName is a rpc call to get session username
// It will get all the session username from the session
func (s *Session) GetSessionUserName(ctx context.Context, req *sessionproto.SessionRequest) (*sessionproto.SessionUserName, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	resp, err := GetSessionUserNameFunc(ctx, req)
	return resp, err
}

// GetSessionUserRoleID is a rpc call to get session user's role ID
// It will get the session username's role id from the session
func (s *Session) GetSessionUserRoleID(ctx context.Context, req *sessionproto.SessionRequest) (*sessionproto.SessionUsersRoleID, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	resp, err := GetSessionUserRoleIDFunc(ctx, req)
	return resp, err
}

// GetAllActiveSessions is a rpc call to get all active sessions
// This method will accepts the sessionrequest which has session id and session token
// and it will call GetAllActiveSessions from the session package
// and respond all the sessionresponse values along with error if there is.
func (s *Session) GetAllActiveSessions(ctx context.Context, req *sessionproto.SessionRequest) (*sessionproto.SessionResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside GetAllActiveSessions function (svc-account-session)")
	var resp sessionproto.SessionResponse
	response := GetAllActiveSessionsFunc(ctx, req)
	body, err := MarshalFunc(response.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying to marshal the response body of the get all active session API: " + err.Error()
		l.LogWithFields(ctx).Printf(response.StatusMessage)
		return &resp, nil
	}
	l.LogWithFields(ctx).Debugf("outgoing response of request to get all active sessions: %s", string(body))
	resp.Body = body
	resp.StatusCode = response.StatusCode
	resp.StatusMessage = response.StatusMessage
	resp.Header = response.Header
	resp.Body = body
	return &resp, nil

}

// GetSessionService is a rpc call to get session service
// which basically checks if the session service is enabled or not
func (s *Session) GetSessionService(ctx context.Context, req *sessionproto.SessionRequest) (*sessionproto.SessionResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.SessionService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside GetSessionService function (svc-account-session)")
	var resp sessionproto.SessionResponse
	response := GetSessionServiceFunc(ctx, req)
	body, err := MarshalFunc(response.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying to marshal the response body of the get session service API: " + err.Error()
		l.LogWithFields(ctx).Printf(response.StatusMessage)
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
