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

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func TestCreateEventSubscription(t *testing.T) {
	config.SetUpMockConfig(t)

	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockSubscriptionCreation)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()

	mockApp := iris.New()
	pluginRoutes := mockApp.Party("/ODIM/v1")
	pluginRoutes.Post("/Subscriptions", CreateEventSubscription)
	pluginRoutes.Delete("/Subscriptions", DeleteEventSubscription)

	reqPostBody := map[string]interface{}{
		"RAIDType": "RAID0",
		"Links": &dmtf.Links{
			Drives: []*dmtf.Link{{Oid: "/ODIM/v1/Systems/5a9e8356-265c-413b-80d2-58210592d931:2/Storage/ArrayControllers-0/Drives/0"}},
		},
	}
	reqBodyBytes, _ := json.Marshal(reqPostBody)

	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       "password",
		"PostBody":       reqBodyBytes,
	}

	e := httptest.New(t, mockApp)

	e.POST("/ODIM/v1/Subscriptions").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusInternalServerError)

	e.DELETE("/ODIM/v1/Subscriptions").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusInternalServerError)
	// invalid client
	config.Data.KeyCertConf.RootCACertificate = nil

	e.POST("/ODIM/v1/Subscriptions").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusInternalServerError)
	e.DELETE("/ODIM/v1/Subscriptions").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusInternalServerError)
	config.SetUpMockConfig(t)

	// Invalid device details
	requestBody2 := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", "deviceHost", "devicePort"),
		"UserName":       "admin",
		"Password":       "password",
	}
	e.POST("/ODIM/v1/Subscriptions").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody2).Expect().
		Status(http.StatusInternalServerError)

		// Invalid device details
	requestBody2 = map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", "deviceHost", "devicePort"),
		"UserName":       "admin",
		"Password":       "password",
		"PostBody":       reqBodyBytes,
	}
	e.POST("/ODIM/v1/Subscriptions").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody2).Expect().
		Status(http.StatusInternalServerError)

	device := dputilities.RedfishDevice{
		Host:     fmt.Sprintf("%s:%s", "deviceHost", "devicePort"),
		Username: "admin",
		Password: "password",
		PostBody: reqBodyBytes,
	}
	isOurSubscription(&device)

}

func mockSubscriptionCreation(username, password, url string, w http.ResponseWriter) {
	fmt.Println("Url is ", url)
	if url == "/ODIM/v1/Subscriptions" {
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
	if url == "/redfish/v1/EventService/Subscriptions" && password == "password" {
		w.WriteHeader(http.StatusOK)
	}
	if url == "/redfish/v1/EventService/Subscriptions" && password == "password1" {
		w.WriteHeader(http.StatusUnauthorized)
	}
	if url == "/redfish/v1/EventService/Subscriptions" && password == "password2" {
		reqPostBody := map[string]interface{}{
			"RAIDType": "RAID0",
			"Members":  []*dmtf.Link{{Oid: "/ODIM/v1/Systems/1"}},
		}
		serviceRootResp, _ := json.Marshal(reqPostBody)
		w.Write(serviceRootResp)

		w.WriteHeader(http.StatusUnauthorized)
	}
	if url == "/ODIM/v1/Systems/1" {
		reqPostBody := map[string]interface{}{
			"RAIDType":    "RAID0",
			"Destination": "https://127.0.0.1:45006/ODIM/v1/Systems/1",
		}
		serviceRootResp, _ := json.Marshal(reqPostBody)
		w.Write(serviceRootResp)

		w.WriteHeader(http.StatusOK)

	}
	w.WriteHeader(http.StatusInternalServerError)
}

func Test_deleteMatchingSubscriptions(t *testing.T) {
	config.SetUpMockConfig(t)

	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockSubscriptionCreation)
	ts.StartTLS()
	defer ts.Close()

	mockApp := iris.New()
	pluginRoutes := mockApp.Party("/ODIM/v1")
	pluginRoutes.Post("/Subscriptions", CreateEventSubscription)
	pluginRoutes.Delete("/Subscriptions", DeleteEventSubscription)

	device := dputilities.RedfishDevice{
		Host:     fmt.Sprintf("%s:%s", deviceHost, devicePort),
		Username: "admin",
		Password: "password1",
	}
	deleteMatchingSubscriptions(&device)

	device = dputilities.RedfishDevice{
		Host:     fmt.Sprintf("%s:%s", deviceHost, devicePort),
		Username: "admin",
		Password: "password",
	}
	deleteMatchingSubscriptions(&device)
	fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&7  ")

	device = dputilities.RedfishDevice{
		Host:     fmt.Sprintf("%s:%s", deviceHost, devicePort),
		Username: "admin",
		Password: "password2",
	}
	deleteMatchingSubscriptions(&device)

}

func Test_isOurSubscription(t *testing.T) {
	config.SetUpMockConfig(t)
	device := dputilities.RedfishDevice{
		Host:     "localhost",
		Username: "admin",
		Password: "password",
		Token:    "test",
		Location: "empty",
	}
	isOurSubscription(&device)
}
