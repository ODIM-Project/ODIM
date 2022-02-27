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

//Package pmbhandle ...
package pmbhandle

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
)

//ContactPlugin is used to send a request to plugin to add a resource
func ContactPlugin(url, method, token string, odataID string, body interface{}, collaboratedInfo map[string]string) (*http.Response, error) {
	jsonStr, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

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
		log.Warn("got " + resp.Status + " while fetching " + url + " with method " + method)
	}

	return resp, nil
}
