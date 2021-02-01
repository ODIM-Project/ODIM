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

// Package session ...
package session

import (
	"fmt"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"time"
)

// UpdateLastUsedTime is supposed to be used whenever there is a session usage.
// The function is for updating the last used time of a session, so that
// the active sessions won't time out and expire. As the input of the function
// we are passing the session token. As return, function give backs the error, if any.
func UpdateLastUsedTime(token string) error {
	session, err := asmodel.GetSession(token)
	if err != nil {
		return fmt.Errorf("error while trying to get the session details with the token %v: %v", token, err)
	}
	session.LastUsedTime = time.Now()
	// Update Session
	err = session.Update()
	if err != nil {
		return fmt.Errorf("error while trying to update session details: %v", err)
	}
	return nil
}
