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

// Package services ...
package services

import (
	"context"
	"fmt"
	"net/http"

	authproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/auth"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	errResponse "github.com/ODIM-Project/ODIM/lib-utilities/response"
	log "github.com/sirupsen/logrus"
)

// IsAuthorized is used to authorize the services using svc-account-session.
// As parameters session token, privileges and oem privileges are passed.
// A RPC call is made with these parameters to the Account-Session service
// to check whether the session is valid and have all the privileges which are passed to it.
func IsAuthorized(sessionToken string, privileges, oemPrivileges []string) errResponse.RPC {
	conn, err := ODIMService.Client(AccountSession)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create client connection: %v", err)
		log.Error(errMsg)
		return GeneralError(http.StatusInternalServerError, errResponse.InternalError, errMsg, nil)
	}
	defer conn.Close()
	asService := authproto.NewAuthorizationClient(conn)
	response, err := asService.IsAuthorized(
		context.TODO(),
		&authproto.AuthRequest{
			SessionToken:  sessionToken,
			Privileges:    privileges,
			Oemprivileges: oemPrivileges,
		},
	)
	if err != nil && response == nil {
		errMsg := fmt.Sprintf("rpc call failed: %v", err)
		log.Error(errMsg)
		return GeneralError(http.StatusInternalServerError, errResponse.InternalError, errMsg, nil)
	}
	var msgArgs []interface{}
	if response.StatusCode == http.StatusServiceUnavailable {
		msgArgs = append(msgArgs, fmt.Sprintf("%v:%v", "", ""))
	}
	return GeneralError(response.StatusCode, response.StatusMessage, "while checking the authorization", msgArgs)
}

// GetSessionUserName will get user name from the session token by rpc call to account-session service
func GetSessionUserName(sessionToken string) (string, error) {
	conn, err := ODIMService.Client(AccountSession)
	if err != nil {
		return "", fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := sessionproto.NewSessionClient(conn)
	response, err := asService.GetSessionUserName(
		context.TODO(),
		&sessionproto.SessionRequest{
			SessionToken: sessionToken,
		},
	)
	if err != nil && response == nil {
		log.Error("something went wrong with rpc call: " + err.Error())
		return "", err
	}
	return response.UserName, err
}

// GetSessionUserRoleID will get user name from the session token by rpc call to account-session service
func GetSessionUserRoleID(sessionToken string) (string, error) {
	conn, err := ODIMService.Client(AccountSession)
	if err != nil {
		return "", fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	asService := sessionproto.NewSessionClient(conn)
	response, err := asService.GetSessionUserRoleID(
		context.TODO(),
		&sessionproto.SessionRequest{
			SessionToken: sessionToken,
		},
	)
	if err != nil && response == nil {
		log.Error("something went wrong with rpc call: " + err.Error())
		return "", err
	}
	return response.RoleID, err
}

// GetUserDetails function is used to get the session details
func GetUserDetails(sessionToken string) (string, string) {
	var err error
	sessionUserName := "null"
	sessionRoleID := "null"
	if sessionToken != "" {
		sessionUserName, err = GetSessionUserName(sessionToken)
		if err != nil {
			errMsg := "while trying to get session details: " + err.Error()
			log.Error(errMsg)
			return "null", "null"
		}
		sessionRoleID, err = GetSessionUserRoleID(sessionToken)
		if err != nil {
			errMsg := "while trying to get session details: " + err.Error()
			log.Error(errMsg)
			return sessionUserName, "null"
		}
	}
	return sessionUserName, sessionRoleID
}

// GeneralError will create the error response
// This function can be used only if the expected response have only
// one extended info object. Error code for the response will be GeneralError

func GeneralError(statusCode int32, statusMsg, errMsg string, msgArgs []interface{}) response.RPC {
	var resp response.RPC
	resp.StatusCode = statusCode
	resp.StatusMessage = statusMsg
	args := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			{
				StatusMessage: resp.StatusMessage,
				ErrorMessage:  errMsg,
				MessageArgs:   msgArgs,
			},
		},
	}
	resp.Body = args.CreateGenericErrorResponse()

	return resp
}
