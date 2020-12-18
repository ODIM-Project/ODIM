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

	authproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/auth"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// Auth functinality will do the following
// 1. It will check whether the session taken is valid
// 2. fetch the previleges from DB against session token
//    and check the service has the previlege
func Auth(req *authproto.AuthRequest) (int32, string) {
	go expiredSessionCleanUp()
	if req.SessionToken == "" {
		log.Error("Unable to validate the token, is empty")
		return http.StatusUnauthorized, response.NoValidSession
	}
	if len(req.Privileges) == 0 {
		log.Error("Unable to validate the privileges, is empty")
		return http.StatusUnauthorized, response.NoValidSession
	}
	session, err := CheckSessionTimeOut(req.SessionToken)
	if err != nil {
		log.Error("Unable to check session timeout: " + err.Error())
		return err.GetAuthStatusCodeAndMessage()
	}
	session.LastUsedTime = time.Now()
	// Update Session
	if err = session.Update(); err != nil {
		log.Error("Unable to update session: " + err.Error())
		return err.GetAuthStatusCodeAndMessage()
	}

	// if the service has all the privileges then return success
	// if any of the privilege isn't assigned to service then return failure
	for _, privilege := range req.Privileges {
		if !session.Privileges[privilege] {
			return http.StatusForbidden, response.InsufficientPrivilege
		}
	}

	// TODO: Need to check OEM Privileges

	log.Info("Authorization successful")
	return http.StatusOK, response.Success
}
