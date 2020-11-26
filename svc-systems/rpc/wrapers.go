/*
 * Copyright (c) 2020 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rpc

import (
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

type authenticator func(sessionToken string, privileges, oemPrivileges []string) (int32, string)

func auth(authenticate authenticator, sessionToken string, callback func() response.RPC) response.RPC {
	if sessionToken == "" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "X-Auth-Token header is missing", nil, nil)
	}

	status, msg := authenticate(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if status == http.StatusOK {
		return callback()
	}
	return common.GeneralError(status, response.NoValidSession, msg, nil, nil)
}
