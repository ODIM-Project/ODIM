// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
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
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-dell/config"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func mockBasicAuthHandler(username, password, url string, w http.ResponseWriter) {
	if url == "/redfish/v1" {
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
	if url == "/ODIM/validate" && username == "admin" && password == "password" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if url == "/ODIM/validate" && (username != "admin" || password != "password") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}

func TestBasicAuth(t *testing.T) {
	config.SetUpMockConfig(t)

	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockBasicAuthHandler)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()

	mockApp := iris.New()
	pluginRoutes := mockApp.Party("/ODIM")
	pluginRoutes.Post("/validate", Validate)

	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       "password",
	}

	e := httptest.New(t, mockApp)

	// test for success scenario
	e.POST("/ODIM/validate").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusOK)

	IoUtilReadAll = func(r io.Reader) ([]byte, error) {
		return []byte{}, fmt.Errorf("fake error")
	}
	e.POST("/ODIM/validate").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusBadRequest)
	IoUtilReadAll = func(r io.Reader) ([]byte, error) {
		return ioutil.ReadAll(r)
	}
	// Invalid Redis client
	config.Data.KeyCertConf.RootCACertificate = nil
	e.POST("/ODIM/validate").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).
		Expect().Status(http.StatusInternalServerError)
	config.SetUpMockConfig(t)

	// Invalid device details
	requestBody1 := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", "deviceHost", "devicePort"),
		"UserName":       "admin",
		"Password":       "password",
	}
	e.POST("/ODIM/validate").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody1).
		Expect().Status(http.StatusBadRequest)

	//Test for Unauthorized scenario: given token is not valid

	e.POST("/ODIM/validate").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody).Expect().
		Status(http.StatusOK)

	e.POST("/ODIM/validate").WithHeader("Authorization", "").WithJSON(requestBody).Expect().Status(http.StatusUnauthorized)
	//Test for the BadRequest: given server details are wrong in Request body
	requestBody["Password"] = "password1"
	e.POST("/ODIM/validate").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQx").WithJSON(requestBody).Expect().
		Status(http.StatusBadRequest)

	//Test for the BadRequest: Request body is not in JSON format
	requestBody2 := "requestbody"
	e.POST("/ODIM/validate").WithHeader("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=").WithJSON(requestBody2).Expect().
		Status(http.StatusBadRequest)

	//Valid token
	tokenDetails = []TokenMap{
		{
			Token:    "validToken",
			LastUsed: time.Now().Add(-5 * time.Minute),
		},
	}
	e.POST("/ODIM/validate").WithHeader("X-Auth-Token", "test").WithJSON(requestBody).Expect().Status(http.StatusUnauthorized)

}

func TestTokenValidation(t *testing.T) {
	config.SetUpMockConfig(t)
	pluginConfig.Data.SessionTimeoutInMinutes = 1
	tokenDetails = []TokenMap{
		{
			Token:    "test",
			LastUsed: time.Now().Add(-5 * time.Minute),
		},
	}
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Positive case ",
			args: args{
				token: "test",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TokenValidation(tt.args.token); got != tt.want {
				t.Errorf("TokenValidation() = %v, want %v", got, tt.want)
			}
		})
	}
}
