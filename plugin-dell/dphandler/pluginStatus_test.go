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
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func TestGetPluginStatus(t *testing.T) {
	config.SetUpMockConfig(t)

	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockPluginStatus)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()

	mockApp := iris.New()
	pluginRoutes := mockApp.Party("/ODIM/v1")
	pluginRoutes.Get("/Status", GetPluginStatus)
	pluginRoutes.Post("/Startup", GetPluginStartup)
	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       "password",
	}

	e := httptest.New(t, mockApp)

	// test for success scenario
	e.GET("/ODIM/v1/Status").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody).Expect().Status(http.StatusOK)
	// invalid token
	e.GET("/ODIM/v1/Status").WithHeader("X-Auth-Token", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody).Expect().Status(http.StatusUnauthorized)

	// startup case
	e.POST("/ODIM/v1/Startup").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody).Expect().Status(http.StatusOK)
	// invalid token
	e.POST("/ODIM/v1/Startup").WithHeader("X-Auth-Token", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody).Expect().Status(http.StatusUnauthorized)
	//invalid data
	requestBody1 := "invalid"
	e.POST("/ODIM/v1/Startup").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody1).Expect().Status(http.StatusBadRequest)

	devices := map[string]dpmodel.DeviceData{"ilo": {
		Address:   fmt.Sprintf("%s:%s", deviceHost, devicePort),
		UserName:  "admin",
		Password:  []byte("password"),
		Operation: "add",
	}}
	deviceBody := dpmodel.StartUpData{
		Devices: devices,
	}
	e.POST("/ODIM/v1/Startup").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(deviceBody).Expect().Status(http.StatusOK)
	devices = map[string]dpmodel.DeviceData{"ilo": {
		Address:   fmt.Sprintf("%s:%s", deviceHost, devicePort),
		UserName:  "admin",
		Password:  []byte("password"),
		Operation: "del",
		EventSubscriptionInfo: &dpmodel.EventSubscriptionInfo{
			Location: "https://" + fmt.Sprintf("%s:%s", deviceHost, devicePort) + "/ODIM/v1/subscription",
		},
	}}
	deviceBody = dpmodel.StartUpData{
		Devices: devices,
	}
	e.POST("/ODIM/v1/Startup").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(deviceBody).Expect().Status(http.StatusOK)

	deviceBody = dpmodel.StartUpData{

		Devices:               devices,
		ResyncEvtSubscription: true,
		RequestType:           "full",
	}
	e.POST("/ODIM/v1/Startup").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(deviceBody).Expect().Status(http.StatusOK)

	devices = map[string]dpmodel.DeviceData{"ilo": {
		Address:   fmt.Sprintf("%s:%s", deviceHost, devicePort),
		UserName:  "admin",
		Password:  []byte("password"),
		Operation: "del",
		EventSubscriptionInfo: &dpmodel.EventSubscriptionInfo{
			EventTypes: []string{"Alert"},
			Location:   "https://" + fmt.Sprintf("%s:%s", deviceHost, devicePort) + "/ODIM/v1/subscriptions",
		},
	}}
	deviceBody = dpmodel.StartUpData{
		Devices:               devices,
		ResyncEvtSubscription: true,
		RequestType:           "full",
	}
	e.POST("/ODIM/v1/Startup").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(deviceBody).Expect().Status(http.StatusOK)

}
func mockPluginStatus(username, password, url string, w http.ResponseWriter) {
	if url == "/ODIM/Status" || url == "/ODIM/Startup" {
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

	if url == "/ODIM/v1/subscription" {
		evt := dpmodel.EvtSubPost{
			EventTypes: []string{"Alert"},
		}
		serviceRootResp, _ := json.Marshal(evt)
		w.WriteHeader(http.StatusOK)
		w.Write(serviceRootResp)
	}
	if url == "/ODIM/v1/subscriptions" {

		serviceRoot := map[string]interface{}{
			"Systems": map[string]string{
				"@odata.id": "/redfish/v1",
			},
		}
		serviceRootResp, _ := json.Marshal(serviceRoot)
		w.WriteHeader(http.StatusNotFound)
		w.Write(serviceRootResp)
	}

}
