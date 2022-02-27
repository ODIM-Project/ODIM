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

// Package auth ...
package auth

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"

	customLogs "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	authproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/auth"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// Auth functionality will do the following
// 1. It will check whether the session taken is valid
// 2. fetch the privileges from DB against session token
//    and check the service has the previlege
func Auth(req *authproto.AuthRequest) (int32, string) {
	go expiredSessionCleanUp()
	if req.SessionToken == "" {
		CustomAuthLog("", "Invalid session token ", http.StatusUnauthorized)
		return http.StatusUnauthorized, response.NoValidSession
	}
	if len(req.Privileges) == 0 {
		CustomAuthLog(req.SessionToken, "Received empty privileges, unable to proceed ", http.StatusForbidden)
		return http.StatusUnauthorized, response.NoValidSession
	}
	session, err := CheckSessionTimeOut(req.SessionToken)
	if err != nil {
		status, message := err.GetAuthStatusCodeAndMessage()
		if status == http.StatusUnauthorized {
			CustomAuthLog("", "Received invalid session token "+req.SessionToken, http.StatusUnauthorized)
		} else {
			log.Error("SessionToken validation failed, unable to proceed " + err.Error())
		}
		return status, message
	}
	session.LastUsedTime = time.Now()
	// Update Session
	if err = session.Update(); err != nil {
		log.Error("SessionToken update failed with error: " + err.Error())
		return err.GetAuthStatusCodeAndMessage()
	}

	// if the service has all the privileges then return success
	// if any of the privilege isn't assigned to service then return failure
	for _, privilege := range req.Privileges {
		if !session.Privileges[privilege] {
			CustomAuthLog(req.SessionToken, "User does not have sufficient privileges", http.StatusForbidden)
			return http.StatusForbidden, response.InsufficientPrivilege
		}
	}

	// TODO: Need to check OEM Privileges

	CustomAuthLog(req.SessionToken, "Authorization is successful", http.StatusOK)
	return http.StatusOK, response.Success
}

// CustomAuthLog function takes session token, message and response status code
// Gets the user id and role id for the session token provided
// logs the messages in custom log format
func CustomAuthLog(sessionToken, msg string, respStatusCode int32) {
	userID := ""
	roleID := ""
	if sessionToken != "" {
		currentSession, err := CheckSessionTimeOut(sessionToken)
		if err == nil {
			userID = currentSession.UserName
			roleID = currentSession.RoleID
		}
	}
	logProperties := make(map[string]interface{})
	logProperties["SessionToken"] = sessionToken
	logProperties["SessionUserID"] = userID
	logProperties["SessionRoleID"] = roleID
	logProperties["Message"] = msg
	logProperties["ResponseStatusCode"] = respStatusCode
	customLogs.AuthLog(logProperties)
}
