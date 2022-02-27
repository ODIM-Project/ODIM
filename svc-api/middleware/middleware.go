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

//Package middleware ...
package middleware

import (
	"github.com/ODIM-Project/ODIM/svc-api/rpc"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
)

//SessionDelMiddleware is used to delete session created for basic auth
func SessionDelMiddleware(ctx iris.Context) {
	ctx.Next()

	sessionID := ctx.Request().Header.Get("Session-ID")
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	if sessionID != "" {
		resp, err := rpc.DeleteSessionRequest(sessionID, sessionToken)
		if err != nil && resp == nil {
			errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
			log.Error(errorMessage)
			return
		}
	}
}
