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

//Package rpc ...
package rpc

import (
	"context"
	"fmt"

	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

// DoSessionCreationRequest will do the rpc calls for the auth
func DoSessionCreationRequest(req sessionproto.SessionCreateRequest) (*sessionproto.SessionCreateResponse, error) {
	conn, err := services.ODIMService.Client(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := sessionproto.NewSessionClient(conn)

	// Call the CreateSession
	rsp, err := asService.CreateSession(context.TODO(), &req)
	if err != nil && rsp == nil {
		return nil, fmt.Errorf("error while trying to make create session rpc call: %v", err)
	}
	return rsp, err
}

// DeleteSessionRequest will do the rpc call to delete session
func DeleteSessionRequest(sessionID, sessionToken string) (*sessionproto.SessionResponse, error) {
	conn, err := services.ODIMService.Client(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := sessionproto.NewSessionClient(conn)

	// Call the DeleteSession
	rsp, err := asService.DeleteSession(context.TODO(), &sessionproto.SessionRequest{
		SessionId:    sessionID,
		SessionToken: sessionToken,
	})
	if err != nil && rsp == nil {
		return nil, fmt.Errorf("error while trying to make delete session rpc call: %v", err)
	}

	return rsp, err
}

// GetSessionRequest will do the rpc call to get session
func GetSessionRequest(sessionID, sessionToken string) (*sessionproto.SessionResponse, error) {
	conn, err := services.ODIMService.Client(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := sessionproto.NewSessionClient(conn)

	// Call the GetSession
	rsp, err := asService.GetSession(context.TODO(), &sessionproto.SessionRequest{
		SessionId:    sessionID,
		SessionToken: sessionToken,
	})
	if err != nil && rsp == nil {
		return nil, fmt.Errorf("error while trying to make get session rpc call: %v", err)
	}

	return rsp, err
}

// GetAllActiveSessionRequest will do the rpc call to get session
func GetAllActiveSessionRequest(sessionID, sessionToken string) (*sessionproto.SessionResponse, error) {
	conn, err := services.ODIMService.Client(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := sessionproto.NewSessionClient(conn)

	// Call the GetAllActiveSessions
	rsp, err := asService.GetAllActiveSessions(context.TODO(), &sessionproto.SessionRequest{
		SessionId:    sessionID,
		SessionToken: sessionToken,
	})
	if err != nil && rsp == nil {
		return nil, fmt.Errorf("error while trying to make get session service rpc call: %v", err)
	}

	return rsp, err
}

//GetSessionServiceRequest will do the rpc call to check session
func GetSessionServiceRequest() (*sessionproto.SessionResponse, error) {
	conn, err := services.ODIMService.Client(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := sessionproto.NewSessionClient(conn)

	// Call the GetSessionService
	rsp, err := asService.GetSessionService(context.TODO(), &sessionproto.SessionRequest{})
	if err != nil && rsp == nil {
		return nil, fmt.Errorf("error while trying to make get session service rpc call: %v", err)
	}

	return rsp, err
}
