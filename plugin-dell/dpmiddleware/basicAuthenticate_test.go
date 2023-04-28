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

// Package dpmiddleware ...
package dpmiddleware

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/plugin-dell/config"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func TestBasicAuth(t *testing.T) {
	config.SetUpMockConfig(t)

	deviceHost := "localhost"
	devicePort := "1234"
	mockApp := iris.New()
	pluginRoutes := mockApp.Party("/ODIM/v1")
	pluginRoutes.Post("/Systems", BasicAuth)
	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       "password",
	}

	e := httptest.New(t, mockApp)

	// test for success scenario
	e.POST("/ODIM/v1/Systems").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusOK)

	// Negative - Invalid token
	e.POST("/ODIM/v1/Systems").WithHeader("Authorization", "YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusUnauthorized)
	// invalid token id
	e.POST("/ODIM/v1/Systems").WithHeader("Authorization", "Basic 1").WithJSON(requestBody).Expect().
		Status(http.StatusInternalServerError)

	// invalid token id
	e.POST("/ODIM/v1/Systems").WithHeader("Authorization", "Basic ").WithJSON(requestBody).Expect().
		Status(http.StatusUnauthorized)

	// Invalid username

	e.POST("/ODIM/v1/Systems").WithHeader("Authorization", "Basic YWRtaW4xOk9kIW0xMiQ0").WithJSON(requestBody).Expect().
		Status(http.StatusUnauthorized)

	// Invalid password
	config.Data.PluginConf.Password = "YWRtaW46cGFzc3dvcmQ"
	e.POST("/ODIM/v1/Systems").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusUnauthorized)
}
