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

// Package common ...
package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	iris "github.com/kataras/iris/v12"
)

const (
	PluginStatusRequestComment = "Plugin Status Request"
	PluginStatusRequestName    = "Common Plugin Status Request"
)

type modifyPluginStatusData func(pluginStatus *PluginStatus)

func TestCheckPluginStatus(t *testing.T) {
	config.SetUpMockConfig(t)
	go mockPlugin(t)
	time.Sleep(2 * time.Second)

	var pluginStatus = &PluginStatus{
		PluginIP:                "localhost",
		PluginPort:              "45100",
		PluginUsername:          "admin",
		PluginUserPassword:      "admin",
		PluginPrefferedAuthType: "BasicAuth",
		CACertificate:           &config.Data.KeyCertConf.RootCACertificate,
		Method:                  http.MethodGet,
		RequestBody: StatusRequest{
			Comment: PluginStatusRequestComment,
			Name:    PluginStatusRequestName,
			Version: "v0.1",
		},
		ResponseWaitTime: 1,
		Count:            3,
	}

	tests := []struct {
		name    string
		exec    modifyPluginStatusData
		p       *PluginStatus
		want    bool
		want1   int
		want2   []string
		wantErr bool
	}{
		{
			name:    "positive case basic auth",
			exec:    nil,
			p:       pluginStatus,
			want:    true,
			want1:   1,
			want2:   []string{"Topic1"},
			wantErr: false,
		},
		{
			name: "positive case auth token",
			exec: func(pluginStatus *PluginStatus) {
				pluginStatus.PluginPrefferedAuthType = "XAuthToken"
			},
			p:       pluginStatus,
			want:    true,
			want1:   1,
			want2:   []string{"Topic1"},
			wantErr: false,
		},
		{
			name: "negative case basic auth",
			exec: func(pluginStatus *PluginStatus) {
				pluginStatus.PluginUserPassword = "password"
			},
			p:       pluginStatus,
			want:    false,
			want1:   3,
			want2:   []string{},
			wantErr: true,
		},
		{
			name: "negative case auth token",
			exec: func(pluginStatus *PluginStatus) {
				pluginStatus.PluginUserPassword = "password"
				pluginStatus.PluginPrefferedAuthType = "XAuthToken"
			},
			p:       pluginStatus,
			want:    false,
			want1:   3,
			want2:   []string{},
			wantErr: true,
		},
		{
			name: "negative case invalid server basic auth",
			exec: func(pluginStatus *PluginStatus) {
				pluginStatus.PluginPort = "8888"
			},
			p:       pluginStatus,
			want:    false,
			want1:   3,
			want2:   []string{},
			wantErr: true,
		},
		{
			name: "negative case invalid server auth token",
			exec: func(pluginStatus *PluginStatus) {
				pluginStatus.PluginPort = "8888"
				pluginStatus.PluginPrefferedAuthType = "XAuthToken"
			},
			p:       pluginStatus,
			want:    false,
			want1:   3,
			want2:   []string{},
			wantErr: true,
		},
		{
			name: "negative case unavailable server",
			exec: func(pluginStatus *PluginStatus) {
				pluginStatus.PluginIP = "abc"
			},
			p:       pluginStatus,
			want:    false,
			want1:   3,
			want2:   []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.exec != nil {
			tt.exec(tt.p)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := tt.p.CheckStatus()
			if (err != nil) != tt.wantErr {
				t.Errorf("PluginStatus.CheckStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PluginStatus.CheckStatus() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PluginStatus.CheckStatus() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("PluginStatus.CheckStatus() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func mockPlugin(t *testing.T) {
	conf := &config.HTTPConfig{
		Certificate:   &config.Data.APIGatewayConf.Certificate,
		PrivateKey:    &config.Data.APIGatewayConf.PrivateKey,
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
		ServerAddress: "localhost",
		ServerPort:    "45100",
	}
	mockServer, err := conf.GetHTTPServerObj()
	if err != nil {
		t.Fatalf("fatal: error while initializing server: %v", err)
	}
	router := iris.New()
	plug := router.Party("/ODIM/v1")
	plug.Get("/Status", mockPluginHandler)
	plug.Post("/Sessions", mockPluginHandler2)
	router.Run(iris.Server(mockServer))
}

func mockPluginHandler(ctx iris.Context) {
	var resp StatusResponse
	token := ctx.GetHeader("X-Auth-Token")
	username, password, _ := ctx.Request().BasicAuth()
	if token == "valid" || (username == "admin" && password == "admin") {
		resp.Status = &PluginResponseStatus{
			Available: "Yes",
		}
		resp.EventMessageBus = &EventMessageBus{
			EmbQueue: []EmbQueue{
				{
					QueueName: "Topic1",
				},
			},
		}
		ctx.StatusCode(http.StatusOK)
	} else {
		resp.Status = &PluginResponseStatus{
			Available: "No",
		}
		ctx.StatusCode(http.StatusBadRequest)
	}

	ctx.JSON(&resp)
}

func mockPluginHandler2(ctx iris.Context) {
	var userCreds map[string]string
	rawBodyAsBytes, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		errorMessage := "Error while trying to validate the credentials: " + err.Error()
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString(errorMessage)
		return
	}
	err = json.Unmarshal(rawBodyAsBytes, &userCreds)
	if err != nil {
		errorMessage := "Error while trying to unmarshal user details: " + err.Error()
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errorMessage)
		return
	}
	//Validate the credentials
	username := userCreds["Username"]
	password := userCreds["Password"]

	if username == "admin" && password == "admin" {
		ctx.Header("X-Auth-Token", "valid")
		ctx.StatusCode(http.StatusCreated)
	} else {
		ctx.StatusCode(http.StatusUnauthorized)
	}
	return
}
