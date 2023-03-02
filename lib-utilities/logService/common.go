//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

// Package logs ...
package logService

import (
	"fmt"

	srv "github.com/ODIM-Project/ODIM/lib-utilities/services"
)

// GetUserDetails function
// getting the session user id and role id for a given session token
func GetUserDetails(sessionToken string) (string, string, error) {
	var err error
	sessionUserName := "null"
	sessionRoleID := "null"
	if sessionToken != "" {
		sessionUserName, err = srv.GetSessionUserName(sessionToken)
		if err != nil {
			return "null", "null", fmt.Errorf("while trying to get session details: %s", err.Error())
		}
		sessionRoleID, err = srv.GetSessionUserRoleID(sessionToken)
		if err != nil {
			return sessionUserName, "null", fmt.Errorf("while trying to get session details: %s", err.Error())
		}
	}
	return sessionUserName, sessionRoleID, nil
}
