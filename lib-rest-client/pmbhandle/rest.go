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

// Package pmbhandle ...
package pmbhandle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
)

// Plugin is the model for plugin information
type Plugin struct {
	IP                string
	Port              string
	Username          string
	Password          []byte
	ID                string
	PluginType        string
	PreferredAuthType string
	ManagerUUID       string
}

// PluginContactRequest holds the details required to contact the plugin
type PluginContactRequest struct {
	URL            string
	HTTPMethodType string
	ContactClient  func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	PostBody       interface{}
	BasicAuth      map[string]string
	Token          string
	Plugin         Plugin
}

// ContactPlugin is used to send a request to plugin to add a resource
func ContactPlugin(ctx context.Context, url, method, token string, odataID string, body interface{}, collaboratedInfo map[string]string) (*http.Response, error) {
	jsonStr, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		l.Log.Error(err.Error())
		return nil, err
	}
	req = CreateHeader(req, ctx)
	// indicate to close the request created
	req.Close = true

	// TODO: it can be saved inside inMemory db for use
	req.Header.Set("Content-Type", "application/json")
	if collaboratedInfo != nil {
		req.SetBasicAuth(collaboratedInfo["UserName"], collaboratedInfo["Password"])
	}
	if token != "" {
		req.Header.Set("X-Auth-Token", token)
	}
	if odataID != "" {
		req.Header.Set("OdataID", odataID)
	}
	httpConf := &config.HTTPConfig{
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
	}
	httpClient, err := httpConf.GetHTTPClientObj()
	if err != nil {
		return nil, err
	}
	config.TLSConfMutex.RLock()
	httpClient.Transport.(*http.Transport).TLSClientConfig.ServerName = collaboratedInfo["ServerName"]
	resp, err := httpClient.Do(req)
	config.TLSConfMutex.RUnlock()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		l.Log.Warn("got " + resp.Status + " while fetching " + url + " with method " + method)
	}

	return resp, nil
}

// ContactPluginWithAuth will check the preferred auth for the plugin and send the request with the preferred auth.
func ContactPluginWithAuth(ctx context.Context, req PluginContactRequest, serverName string) (*http.Response, error) {
	req.BasicAuth = map[string]string{}
	req.BasicAuth["ServerName"] = serverName
	if strings.EqualFold(req.Plugin.PreferredAuthType, "XAuthToken") {
		payload := map[string]interface{}{
			"Username": req.Plugin.Username,
			"Password": string(req.Plugin.Password),
		}
		reqURL := fmt.Sprintf("https://%s/ODIM/v1/Sessions", net.JoinHostPort(req.Plugin.IP, req.Plugin.Port))
		response, err := ContactPlugin(ctx, reqURL, http.MethodPost, "", "", payload, nil)
		if err != nil || (response != nil && response.StatusCode != http.StatusOK) {
			return nil,
				fmt.Errorf("failed to get session token from %s: %s: %+v", req.Plugin.ID, err.Error(), response)
		}
		req.Token = response.Header.Get("X-Auth-Token")
	} else {
		req.BasicAuth["UserName"] = req.Plugin.Username
		req.BasicAuth["Password"] = string(req.Plugin.Password)
	}
	reqURL := fmt.Sprintf("https://%s%s", net.JoinHostPort(req.Plugin.IP, req.Plugin.Port), req.URL)
	return ContactPlugin(ctx, reqURL, req.HTTPMethodType, req.Token, "", req.PostBody, req.BasicAuth)
}

// CreateHeader is used to get data from context and set it to header for http request call
func CreateHeader(req *http.Request, ctx context.Context) *http.Request {
	if ctx.Value("transactionid") != nil {
		transactionId := ctx.Value("transactionid").(string)
		actionId := ctx.Value("actionid").(string)
		actionName := ctx.Value("actionname").(string)
		threadId := ctx.Value("threadid").(string)
		threadName := ctx.Value("threadname").(string)
		processName := ctx.Value("processname").(string)
		req.Header.Set(common.TransactionID, transactionId)
		req.Header.Set(common.ActionID, actionId)
		req.Header.Set(common.ActionName, actionName)
		req.Header.Set(common.ThreadID, threadId)
		req.Header.Set(common.ThreadName, threadName)
		req.Header.Set(common.ProcessName, processName)
	}
	return req
}
