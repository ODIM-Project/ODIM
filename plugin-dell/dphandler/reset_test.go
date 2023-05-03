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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func TestResetComputerSystem(t *testing.T) {
	config.SetUpMockConfig(t)
	updateResetResponse()
	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockResetComputerSystem)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()

	mockApp := iris.New()
	pluginRoutes := mockApp.Party("/ODIM/v1")
	pluginRoutes.Post("/Systems/ComputerSystem.Reset", ResetComputerSystem)

	attributes := map[string]interface{}{"Image": "abc", "ResetType": "force"}
	attributeByte, _ := json.Marshal(attributes)
	requestBody := dpmodel.Device{
		Host:     fmt.Sprintf("%s:%s", deviceHost, devicePort),
		Username: "admin",
		Password: []byte("password"),
		Location: "destination",
		PostBody: attributeByte,
	}
	e := httptest.New(t, mockApp)

	// test for success scenario
	e.POST("/ODIM/v1/Systems/ComputerSystem.Reset").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusOK)

	// Invalid Token
	e.POST("/ODIM/v1/Systems/ComputerSystem.Reset").WithHeader("X-Auth-Token", "invalidToken").WithJSON(requestBody).Expect().
		Status(http.StatusUnauthorized)
	// Invalid request body
	requestBody1 := "invalid"
	e.POST("/ODIM/v1/Systems/ComputerSystem.Reset").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody1).Expect().
		Status(http.StatusBadRequest)

	requestBody2 := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       "password",
	}
	e.POST("/ODIM/v1/Systems/ComputerSystem.Reset").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody2).Expect().Status(http.StatusBadRequest)

	QueryDevice = func(ctx context.Context, uri string, device *dputilities.RedfishDevice, method string) (int, http.Header, []byte, error) {
		return http.StatusInternalServerError, httptest.NewRecorder().HeaderMap, nil, fmt.Errorf("fake error ")
	}

	e.POST("/ODIM/v1/Systems/ComputerSystem.Reset").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusInternalServerError)
	QueryDevice = func(ctx context.Context, uri string, device *dputilities.RedfishDevice, method string) (int, http.Header, []byte, error) {
		return http.StatusOK, httptest.NewRecorder().HeaderMap, nil, nil
	}
	//  Invalid reset Type
	attributes = map[string]interface{}{"Image": "abc", "ResetType": "On"}
	attributeByte, _ = json.Marshal(attributes)
	requestBody = dpmodel.Device{
		PostBody: attributeByte,
	}
	e.POST("/ODIM/v1/Systems/ComputerSystem.Reset").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody).Expect().Status(http.StatusBadRequest)

	// Power state on
	attributes1 := map[string]interface{}{"Image": "abc", "PowerState": "On"}
	attributeByte1, _ := json.Marshal(attributes1)
	QueryDevice = func(ctx context.Context, uri string, device *dputilities.RedfishDevice, method string) (int, http.Header, []byte, error) {
		return http.StatusOK, httptest.NewRecorder().HeaderMap, attributeByte1, nil
	}
	e.POST("/ODIM/v1/Systems/ComputerSystem.Reset").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody).Expect().Status(http.StatusOK)

	//	// Power state Off

	attributes1 = map[string]interface{}{"Image": "abc", "PowerState": "Off"}
	attributeByte1, _ = json.Marshal(attributes1)

	QueryDevice = func(ctx context.Context, uri string, device *dputilities.RedfishDevice, method string) (int, http.Header, []byte, error) {
		return http.StatusOK, httptest.NewRecorder().HeaderMap, attributeByte1, nil
	}

	attributes = map[string]interface{}{"Image": "abc", "ResetType": "ForceOff"}
	attributeByte, _ = json.Marshal(attributes)
	requestBody = dpmodel.Device{
		PostBody: attributeByte,
	}
	e.POST("/ODIM/v1/Systems/ComputerSystem.Reset").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody).Expect().Status(http.StatusOK)

	attributes = map[string]interface{}{"Image": "abc", "ResetType": "Invalid"}
	attributeByte, _ = json.Marshal(attributes)
	requestBody = dpmodel.Device{
		PostBody: attributeByte,
	}
	e.POST("/ODIM/v1/Systems/ComputerSystem.Reset").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody).Expect().Status(http.StatusConflict)

	QueryDevice = func(ctx context.Context, uri string, device *dputilities.RedfishDevice, method string) (int, http.Header, []byte, error) {
		return queryDevice(ctx, uri, device, method)
	}
}
func mockResetComputerSystem(username, password, url string, w http.ResponseWriter) {
	if url == "/ODIM/v1/Systems/ComputerSystem.Reset" {
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
