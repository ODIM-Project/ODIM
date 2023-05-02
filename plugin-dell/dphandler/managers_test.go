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

// Package dphandler
package dphandler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpresponse"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func mockManagers(username, password, url string, w http.ResponseWriter) {
	body := `{"data": "success"}`
	insertURL := "/redfish/v1/Managers/1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia"
	ejectURL := "/redfish/v1/Managers/1/VirtualMedia/1/Actions/VirtualMedia.EjectMedia"
	insertURL = replaceURI(insertURL)
	ejectURL = replaceURI(ejectURL)

	if url == insertURL && username == "admin" {
		e, _ := json.Marshal(body)
		w.WriteHeader(http.StatusOK)
		w.Write(e)
		return
	}
	if url == insertURL && username != "admin" {
		e, _ := json.Marshal(body)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(e)
		return
	}
	if url == ejectURL && username == "admin" {
		e, _ := json.Marshal(body)
		w.WriteHeader(http.StatusOK)
		w.Write(e)
		return
	}
	if url == ejectURL && username != "admin" {
		e, _ := json.Marshal(body)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(e)
		return
	}
	if url == "Managers" {
		e, _ := json.Marshal(body)
		w.WriteHeader(http.StatusOK)
		w.Write(e)
		return
	}
}

func TestGetManagerCollection(t *testing.T) {
	deviceHost := "localhost"
	devicePort := "1234"
	url := "/redfish/v1"
	url = replaceURI(url)

	config.SetUpMockConfig(t)
	ts := startTestServer(mockManagers)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()
	time.Sleep(1 * time.Second)
	mockApp := iris.New()

	redfishRoutes := mockApp.Party(url)

	redfishRoutes.Get("/Managers", GetManagersCollection)

	dpresponse.PluginToken = "token"
	e := httptest.New(t, mockApp)

	var deviceDetails = dpmodel.Device{
		Host: "",
	}
	managerURL := url + "/Managers"
	//Unit Test for success scenario
	e.GET(managerURL).WithJSON(deviceDetails).Expect().Status(http.StatusOK)

	//Case for invalid token
	e.GET(managerURL).WithHeader("X-Auth-Token", "Invalidtoken").WithJSON(deviceDetails).Expect().Status(http.StatusUnauthorized)

	// case for positive
	tokenDetails = []TokenMap{
		{
			Token: "valid",
		},
	}
	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       "password",
	}
	e.GET(managerURL).WithHeader("X-Auth-Token", "valid").WithJSON(requestBody).Expect().Status(http.StatusOK)

	config.Data.KeyCertConf.RootCACertificate = nil
	e.GET(managerURL).WithHeader("X-Auth-Token", "valid").WithJSON(requestBody).Expect().Status(http.StatusInternalServerError)
	config.SetUpMockConfig(t)

	requestBody2 := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", "deviceHost", "devicePort"),
		"UserName":       "admin",
		"Password":       "password",
	}
	e.GET(managerURL).WithHeader("X-Auth-Token", "valid").WithJSON(requestBody2).Expect().Status(http.StatusInternalServerError)

	IoUtilReadAll = func(r io.Reader) ([]byte, error) {
		return nil, fmt.Errorf("fake erro ")
	}
	e.GET(managerURL).WithHeader("X-Auth-Token", "valid").WithJSON(requestBody).Expect().Status(http.StatusOK)
	IoUtilReadAll = func(r io.Reader) ([]byte, error) {
		return ioutil.ReadAll(r)
	}

}

func TestGetManager(t *testing.T) {
	url := "/redfish/v1"
	url = replaceURI(url)
	config.SetUpMockConfig(t)
	mockApp := iris.New()
	redfishRoutes := mockApp.Party(url)
	redfishRoutes.Get("/Managers", GetManagersInfo)

	dpresponse.PluginToken = "token"
	e := httptest.New(t, mockApp)

	var deviceDetails = dpmodel.Device{
		Host: "",
	}
	managerURL := url + "/Managers"
	//Unit Test for success scenario
	e.GET(managerURL).WithJSON(deviceDetails).Expect().Status(http.StatusOK)

	//Case for invalid token
	e.GET(managerURL).WithHeader("X-Auth-Token", "Invalidtoken").WithJSON(deviceDetails).Expect().Status(http.StatusUnauthorized)

}

func TestVirtualMediaActions(t *testing.T) {
	config.SetUpMockConfig(t)
	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockManagers)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()
	time.Sleep(1 * time.Second)
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Post("/Managers/1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia", VirtualMediaActions)
	redfishRoutes.Post("/Managers/1/VirtualMedia/1/Actions/VirtualMedia.EjectMedia", VirtualMediaActions)
	dpresponse.PluginToken = "token"

	test := httptest.New(t, mockApp)
	attributes := map[string]interface{}{"Image": "abc"}
	attributeByte, _ := json.Marshal(attributes)
	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       []byte("P@$$w0rd"),
		"PostBody":       attributeByte,
	}
	test.POST("/redfish/v1/Managers/1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia").WithJSON(requestBody).Expect().Status(http.StatusOK)
	requestBody["UserName"] = "invalid"
	test.POST("/redfish/v1/Managers/1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia").WithJSON(requestBody).Expect().Status(http.StatusBadRequest)

	requestBody = map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       []byte("P@$$w0rd"),
	}
	test.POST("/redfish/v1/Managers/1/VirtualMedia/1/Actions/VirtualMedia.EjectMedia").WithJSON(requestBody).Expect().Status(http.StatusOK)
	requestBody["UserName"] = "invalid"
	test.POST("/redfish/v1/Managers/1/VirtualMedia/1/Actions/VirtualMedia.EjectMedia").WithJSON(requestBody).Expect().Status(http.StatusBadRequest)

	// invalid device details
	requestbody1 := "invalid"
	test.POST("/redfish/v1/Managers/1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia").WithJSON(requestbody1).Expect().
		Status(http.StatusBadRequest)

	config.Data.KeyCertConf.RootCACertificate = nil

	test.POST("/redfish/v1/Managers/1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia").WithJSON(requestBody).Expect().
		Status(http.StatusInternalServerError)
	config.SetUpMockConfig(t)

	createVirtualMediaActionResponse()
}

func Test_getInfoFromDevice(t *testing.T) {
	type args struct {
		uri           string
		deviceDetails dpmodel.Device
		ctx           iris.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getInfoFromDevice(tt.args.uri, tt.args.deviceDetails, tt.args.ctx)
		})
	}
}
