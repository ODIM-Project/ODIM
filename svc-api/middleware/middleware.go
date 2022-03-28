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
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
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

// ResourceRateLimiter will Limit the get on resource untill previous get completed the task
func ResourceRateLimiter(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	for _, val := range config.Data.ResourceRateLimit {
		resourceLimit := strings.Split(val, ":")
		if len(resourceLimit) > 1 && resourceLimit[1] != "" {
			rLimit, _ := strconv.Atoi(resourceLimit[1])
			resource := strings.Replace(resourceLimit[0], "{id}", "[a-zA-Z0-9._-]+", -1)
			regex := regexp.MustCompile(resource)
			if regex.MatchString(uri) {
				conn, err := common.GetDBConnection(common.InMemory)
				if err != nil {
					log.Error(err.Error())
					response := common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
					common.SetResponseHeader(ctx, response.Header)
					ctx.StatusCode(http.StatusInternalServerError)
					ctx.JSON(&response.Body)
					return
				}
				// convert millisecond to second
				expiretime := rLimit / 1000
				if err = conn.SetExpire("ResourceRateLimit", uri, "", expiretime); err != nil {
					errorMessage := "too many requests, retry after some time"
					log.Error(errorMessage)
					response := common.GeneralError(http.StatusTooManyRequests, response.RateLimitExceeded, errorMessage, nil, nil)
					common.SetResponseHeader(ctx, response.Header)
					ctx.StatusCode(http.StatusTooManyRequests)
					ctx.JSON(&response.Body)
					return
				}
			}
		}
	}
	ctx.Next()
}
