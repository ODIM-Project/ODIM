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

// Package dphandler ...
package dphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/plugin-dell/config"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func TestGetResource(t *testing.T) {

	config.SetUpMockConfig(t)

	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockGetResourceHandler)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()

	mockApp := iris.New()
	pluginRoutes := mockApp.Party("/ODIM")
	pluginRoutes.Post("/Systems", GetResource)
	pluginRoutes.Post("/Systems/{id}", GetResource)
	pluginRoutes.Post("/Systems/{id}/Chassis/{id2}/NetworkAdapters", GetResource)

	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       "password",
	}

	e := httptest.New(t, mockApp)

	// Positive case
	e.POST("/ODIM/Systems").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusOK)
		// Invalid Token
	e.POST("/ODIM/Systems").WithHeader("X-Auth-Token", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusUnauthorized)

		// Invalid device Details
	requestBody1 := "invalid"
	e.POST("/ODIM/Systems").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody1).Expect().
		Status(http.StatusBadRequest)
	// Invalid redis client
	config.Data.KeyCertConf.RootCACertificate = nil
	e.POST("/ODIM/Systems").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusInternalServerError)
	config.SetUpMockConfig(t)

	requestBody2 := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", "deviceHost", "devicePort"),
		"UserName":       "admin",
		"Password":       "password",
	}
	e.POST("/ODIM/Systems").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody2).Expect().
		Status(http.StatusInternalServerError)

		// Positive case
	e.POST("/ODIM/Systems/1").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusUnauthorized)
	e.POST("/ODIM/Systems/1/Chassis/1/NetworkAdapters").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusOK)
}
func mockGetResourceHandler(username, password, url string, w http.ResponseWriter) {
	if url == "/ODIM/Systems" {
		serviceRoot := map[string]interface{}{
			"Systems": map[string]string{
				"@odata.id": "/redfish/v1",
			},
		}
		serviceRootResp, _ := json.Marshal(serviceRoot)
		w.WriteHeader(http.StatusOK)
		w.Write(serviceRootResp)
		return
	}
	if url == "/ODIM/Systems/1" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if url == "/ODIM/Systems/1/Systems/1/NetworkAdapters" {
		w.WriteHeader(http.StatusOK)
	}

	w.WriteHeader(http.StatusInternalServerError)
}
