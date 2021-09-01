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

// Packahe rfphandler ...
package rfphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/plugin-redfish/config"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpresponse"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"io/ioutil"
	"net/http"
	"testing"
)

func mockCreateVolume(username, url string) (*http.Response, error) {
	if url == "/ODIM/v1/Systems/1/Storage/1/Volumes" && username == "admin" {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Success")),
		}, nil
	}
	if url == "/ODIM/v1/Systems/1/Storage/1/Volumes" && username != "admin" {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Failed")),
		}, fmt.Errorf("Error")
	}
	if url == "/ODIM/v1/Systems/1/Storage/1/Volumes/1" && username == "admin" {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Success")),
		}, nil
	}
	if url == "/ODIM/v1/Systems/1/Storage/1/Volumes/1" && username != "admin" {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Failed")),
		}, fmt.Errorf("Error")
	}
	return nil, fmt.Errorf("Error")
}

func mockDevice(username, password, url string, w http.ResponseWriter) {
	resp, err := mockCreateVolume(username, url)
	if err != nil && resp == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(resp.StatusCode)
	return
}

func TestCreateVolume(t *testing.T) {
	config.SetUpMockConfig(t)
	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockDevice)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")

	redfishRoutes.Post("/Systems/{id}/Storage/{rid}/Volumes", CreateVolume)

	rfpresponse.PluginToken = "token"

	e := httptest.New(t, mockApp)

	reqPostBody := map[string]interface{}{
		"Name":     "Volume_Test1",
		"RAIDType": "RAID0",
		"Drives":   []dmtf.Link{{Oid: "/ODIM/v1/Systems/5a9e8356-265c-413b-80d2-58210592d931:1/Storage/ArrayControllers-0/Drives/0"}},
	}
	reqBodyBytes, _ := json.Marshal(reqPostBody)
	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       []byte("P@$$w0rd"),
		"PostBody":       reqBodyBytes,
	}

	//Unit Test for success scenario
	e.POST("/redfish/v1/Systems/1/Storage/1/Volumes").WithJSON(requestBody).Expect().Status(http.StatusOK)

	//Case for invalid token
	e.POST("/redfish/v1/Systems/1/Storage/1/Volumes").WithHeader("X-Auth-Token", "token").WithJSON(requestBody).Expect().Status(http.StatusUnauthorized)

	//unittest for bad request scenario
	invalidRequestBody := "invalid"
	e.POST("/redfish/v1/Systems/1/Storage/1/Volumes").WithJSON(invalidRequestBody).Expect().Status(http.StatusBadRequest)
}

func TesDeleteVolume(t *testing.T) {
	config.SetUpMockConfig(t)
	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockDevice)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")

	redfishRoutes.Delete("/Systems/{id}/Storage/{id2}/Volumes/rid", CreateVolume)

	rfpresponse.PluginToken = "token"

	e := httptest.New(t, mockApp)

	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       []byte("P@$$w0rd"),
	}

	//Unit Test for success scenario
	e.DELETE("/redfish/v1/Systems/1/Storage/1/Volumes/1").WithJSON(requestBody).Expect().Status(http.StatusOK)

	//Case for invalid token
	e.DELETE("/redfish/v1/Systems/1/Storage/1/Volumes/1").WithHeader("X-Auth-Token", "token").WithJSON(requestBody).Expect().Status(http.StatusUnauthorized)

}
