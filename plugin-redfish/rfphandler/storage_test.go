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
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/plugin-redfish/config"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpresponse"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
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
	return nil, fmt.Errorf("Error")
}

func TestCreateVolume(t *testing.T) {
	config.SetUpMockConfig(t)

	ts := startTestServer(mockDeviceHandler)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")

	redfishRoutes.Post("/Systems/{id}/Storage/{rid}/Volumes", CreateVolume)

	rfpresponse.PluginToken = "token"

	e := httptest.New(t, mockApp)

	var requestBody = `{
		"Name":"Volume 1",
		"RAIDType":"RAID0",
		"Drives": [{"@odata.id":"/ODIM/v1/Systems/1/Storage/ArrayControllers-0/Drives/1"}],
		"@Redfish.OperationApplyTime": "OnReset"
	}`

	//Unit Test for success scenario
	e.POST("/redfish/v1/Systems/1/Storage/1/Volumes").WithJSON(requestBody).Expect().Status(http.StatusNotImplemented) //ToDo Change the status to http.StatusOK

	//Case for invalid token
	e.POST("/redfish/v1/Systems/1/Storage/1/Volumes").WithHeader("X-Auth-Token", "token").WithJSON(requestBody).Expect().Status(http.StatusNotImplemented) //ToDo Change the status to http.StatusUnauthorized

	//unittest for bad request scenario
	invalidRequestBody := "invalid"
	e.POST("/redfish/v1/Systems/1/Storage/1/Volumes").WithJSON(invalidRequestBody).Expect().Status(http.StatusNotImplemented) //ToDo Change the status to http.StatusBadRequest
}
