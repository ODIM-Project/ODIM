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
	"fmt"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-api/ratelimiter"
)

var (
	NewSessionClientFunc = sessionproto.NewSessionClient
)

// DoSessionCreationRequest will do the rpc calls for the auth
func DoSessionCreationRequest(ctx context.Context, req sessionproto.SessionCreateRequest) (*sessionproto.SessionCreateResponse, error) {
	ctx = common.CreateMetadata(ctx)
	if config.Data.SessionLimitCountPerUser > 0 {
		request := make(map[string]interface{})
		err := json.Unmarshal(req.RequestBody, &request)
		if err != nil {
			return nil, err
		}
		rerr := ratelimiter.SessionRateLimiter(ctx, request["UserName"].(string))
		if rerr != nil {
			fmt.Println("Error in session rate limit: ", rerr)
			return nil, rerr
		}
		defer ratelimiter.DecrementCounter(request["UserName"].(string), ratelimiter.UserRateLimit)
	}
	conn, err := ClientFunc(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewSessionClientFunc(conn)

	// Call the CreateSession
	rsp, err := asService.CreateSession(ctx, &req)
	if err != nil && rsp == nil {
		return nil, fmt.Errorf("error while trying to make create session rpc call: %v", err)
	}
	defer conn.Close()
	return rsp, err
}

// DeleteSessionRequest will do the rpc call to delete session
func DeleteSessionRequest(ctx context.Context, sessionID, sessionToken string) (*sessionproto.SessionResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewSessionClientFunc(conn)

	// Call the DeleteSession
	rsp, err := asService.DeleteSession(ctx, &sessionproto.SessionRequest{
		SessionId:    sessionID,
		SessionToken: sessionToken,
	})
	if err != nil && rsp == nil {
		return nil, fmt.Errorf("error while trying to make delete session rpc call: %v", err)
	}
	defer conn.Close()
	return rsp, err
}

// GetSessionRequest will do the rpc call to get session
func GetSessionRequest(ctx context.Context, sessionID, sessionToken string) (*sessionproto.SessionResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewSessionClientFunc(conn)

	// Call the GetSession
	rsp, err := asService.GetSession(ctx, &sessionproto.SessionRequest{
		SessionId:    sessionID,
		SessionToken: sessionToken,
	})
	if err != nil && rsp == nil {
		return nil, fmt.Errorf("error while trying to make get session rpc call: %v", err)
	}
	defer conn.Close()
	return rsp, err
}

// GetAllActiveSessionRequest will do the rpc call to get session
func GetAllActiveSessionRequest(ctx context.Context, sessionID, sessionToken string) (*sessionproto.SessionResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewSessionClientFunc(conn)

	// Call the GetAllActiveSessions
	rsp, err := asService.GetAllActiveSessions(ctx, &sessionproto.SessionRequest{
		SessionId:    sessionID,
		SessionToken: sessionToken,
	})
	if err != nil && rsp == nil {
		return nil, fmt.Errorf("error while trying to make get session service rpc call: %v", err)
	}
	defer conn.Close()
	return rsp, err
}

// GetSessionServiceRequest will do the rpc call to check session
func GetSessionServiceRequest(ctx context.Context) (*sessionproto.SessionResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	asService := NewSessionClientFunc(conn)

	// Call the GetSessionService
	rsp, err := asService.GetSessionService(ctx, &sessionproto.SessionRequest{})
	if err != nil && rsp == nil {
		return nil, fmt.Errorf("error while trying to make get session service rpc call: %v", err)
	}
	defer conn.Close()
	return rsp, err
}
