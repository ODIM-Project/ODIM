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

// Package tmodel ...
package tmodel

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
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

// GetAllPlugins is for fetching all the plugins added andn stored in db.
func GetAllPlugins(ctx context.Context) ([]Plugin, error) {
	keys, err := GetAllKeysFromTable("Plugin")
	if err != nil {
		return nil, err
	}
	var plugins []Plugin
	for _, key := range keys {
		plugin, err := GetPluginData(key)
		if err != nil {
			l.LogWithFields(ctx).Error("failed to get details of " + key + " plugin: " + err.Error())
			continue
		}
		plugins = append(plugins, plugin)
	}
	return plugins, nil
}

// GetPluginData will fetch plugin details
func GetPluginData(pluginID string) (Plugin, *errors.Error) {
	var plugin Plugin

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return plugin, errors.PackError(err.ErrNo(), "error while trying to connect to DB: ", err.Error())
	}

	plugindata, err := conn.Read("Plugin", pluginID)
	if err != nil {
		return plugin, errors.PackError(err.ErrNo(), "error while trying to fetch plugin data: ", err.Error())
	}

	if err := json.Unmarshal([]byte(plugindata), &plugin); err != nil {
		return plugin, errors.PackError(errors.JSONUnmarshalFailed, err)
	}

	bytepw, errs := common.DecryptWithPrivateKey([]byte(plugin.Password))
	if errs != nil {
		return Plugin{}, errors.PackError(errors.DecryptionFailed, "error: "+pluginID+" plugin password decryption failed: "+errs.Error())
	}
	plugin.Password = bytepw

	return plugin, nil
}

// GetTaskMonResponse will request plugin to get plugin task status
func GetTaskMonResponse(ctx context.Context, plugin Plugin, task *common.PluginTask) (*http.Response, error) {
	contactRequest := PluginContactRequest{}
	plugin.IP = task.IP
	contactRequest.Plugin = plugin
	contactRequest.URL = task.PluginTaskMonURL
	contactRequest.HTTPMethodType = http.MethodGet
	response, err := ContactPlugin(ctx, contactRequest, task.PluginServerName)
	if err != nil {
		l.LogWithFields(ctx).Errorf("failed to get taskmon response from %s(%s): %s: %+v",
			plugin.ID, plugin.IP, err.Error(), response)
		return nil, err
	}
	l.LogWithFields(ctx).Infof("Successfully got task response from %s(%s)", plugin.ID, plugin.IP)
	return response, nil
}

// ContactPlugin is for sending requests to a plugin.
func ContactPlugin(ctx context.Context, req PluginContactRequest, serverName string) (*http.Response, error) {
	req.BasicAuth = map[string]string{}
	req.BasicAuth["ServerName"] = serverName
	if strings.EqualFold(req.Plugin.PreferredAuthType, "XAuthToken") {
		payload := map[string]interface{}{
			"Username": req.Plugin.Username,
			"Password": string(req.Plugin.Password),
		}
		reqURL := fmt.Sprintf("https://%s/ODIM/v1/Sessions", net.JoinHostPort(req.Plugin.IP, req.Plugin.Port))
		response, err := pmbhandle.ContactPlugin(ctx, reqURL, http.MethodPost, "", "", payload, nil)
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
	return pmbhandle.ContactPlugin(ctx, reqURL, req.HTTPMethodType, req.Token, "", req.PostBody, req.BasicAuth)
}
