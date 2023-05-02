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
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/plugin-dell/config"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func TestSetDefaultBootOrder(t *testing.T) {
	config.SetUpMockConfig(t)

	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockBootOrderHandler)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()

	mockApp := iris.New()
	pluginRoutes := mockApp.Party("/ODIM")
	pluginRoutes.Post("/ComputerSystem.SetDefaultBootOrder", SetDefaultBootOrder)

	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       "password",
	}

	e := httptest.New(t, mockApp)

	// test for success scenario
	e.POST("/ODIM/ComputerSystem.SetDefaultBootOrder").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody).Expect().Status(http.StatusOK)

	// test for success scenario
	e.POST("/ODIM/ComputerSystem.SetDefaultBootOrder").WithHeader("X-Auth-Token", "Invalid").
		WithJSON(requestBody).Expect().Status(http.StatusUnauthorized)

	// Invalid DeviceDetails
	requestBody1 := "invalid"
	e.POST("/ODIM/ComputerSystem.SetDefaultBootOrder").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody1).Expect().Status(http.StatusBadRequest)

	// Invalid Redfish Client data
	config.Data.KeyCertConf.RootCACertificate = nil
	e.POST("/ODIM/ComputerSystem.SetDefaultBootOrder").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody).Expect().Status(http.StatusInternalServerError)
	config.SetUpMockConfig(t)

	// Invalid device details
	requestBody2 := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", "deviceHost", "devicePort"),
		"UserName":       "admin",
		"Password":       "password",
	}
	e.POST("/ODIM/ComputerSystem.SetDefaultBootOrder").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody2).Expect().Status(http.StatusInternalServerError)

	// Invalid IO Util
	IoUtilReadAll = func(r io.Reader) ([]byte, error) {
		return nil, fmt.Errorf("fake error ")
	}
	e.POST("/ODIM/ComputerSystem.SetDefaultBootOrder").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").
		WithJSON(requestBody).Expect().Status(http.StatusOK)
	IoUtilReadAll = func(r io.Reader) ([]byte, error) {
		return ioutil.ReadAll(r)
	}

}

func mockBootOrderHandler(username, password, url string, w http.ResponseWriter) {
	if url == "/ODIM/ComputerSystem.SetDefaultBootOrder" {
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

	w.WriteHeader(http.StatusInternalServerError)
}
