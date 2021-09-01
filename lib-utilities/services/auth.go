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
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	authproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/auth"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	errResponse "github.com/ODIM-Project/ODIM/lib-utilities/response"
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
		return common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errMsg, nil, nil)
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
		return common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errMsg, nil, nil)
	}
	var msgArgs []interface{}
	if response.StatusCode == http.StatusServiceUnavailable {
		msgArgs = append(msgArgs, fmt.Sprintf("%v:%v", config.Data.DBConf.InMemoryHost, config.Data.DBConf.InMemoryPort))
	}
	return common.GeneralError(response.StatusCode, response.StatusMessage, "while checking the authorization", msgArgs, nil)
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
