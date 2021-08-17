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

//Package auth ...
package auth

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	srv "github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// Authentication is used to authenticate using session token from svc-account-session
func Authentication(sessionToken string, privileges []string) response.RPC {
	oemprivileges := []string{}
	return srv.IsAuthorized(sessionToken, privileges, oemprivileges)
}

// GetSessionUserName is used to authenticate using session token from svc-account-session
func GetSessionUserName(sessionToken string) (string, error) {
	return srv.GetSessionUserName(sessionToken)
}
