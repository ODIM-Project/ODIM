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

package chassis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	dmtfmodel "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	log "github.com/sirupsen/logrus"
)

// getFabricChassisResource will collect the individual
// fabric chassis resourse from all the fabric plugin
func (f *fabricFactory) getFabricChassisResource(rID string) response.RPC {
	var resp response.RPC
	ch := make(chan response.RPC)

	managers, err := f.getFabricManagers()
	if err != nil {
		log.Warn("while trying to collect fabric managers details from DB, got " + err.Error())
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "", []interface{}{"Chassis", rID}, nil)
	}

	for _, manager := range managers {
		go f.getResource(manager, rID, ch)
	}

	for i := 0; i < len(managers); i++ {
		resp = <-ch
		if is2xx(int(resp.StatusCode)) {
			return resp
		}
	}

	return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "", []interface{}{"Chassis", rID}, nil)
}

// getResource is for collecting the fabric chassis from the individual plugin,
// and returns the result through the channel ch
func (f *fabricFactory) getResource(plugin smodel.Plugin, rID string, ch chan response.RPC) {
	req, errResp, err := f.createChassisRequest(plugin, fmt.Sprintf("%s/%s", collectionURL, rID), http.MethodGet, nil)
	if errResp != nil {
		log.Warn("while trying to create fabric plugin request for " + plugin.ID + ", got " + err.Error())
		ch <- *errResp
		return
	}
	ch <- collectChassisResource(f, req)
}

// collectChassisResource contacts the plugin with the details available in the
// pluginContactRequest, and returns the RPC response
func collectChassisResource(f *fabricFactory, pluginRequest *pluginContactRequest) (r response.RPC) {
	body, _, statusCode, _, err := contactPlugin(pluginRequest)
	if statusCode == http.StatusUnauthorized && strings.EqualFold(pluginRequest.Plugin.PreferredAuthType, "XAuthToken") {
		body, _, statusCode, _, err = retryFabricsOperation(f, pluginRequest)
	}
	if err != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
	}
	if !is2xx(int(statusCode)) {
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "", []interface{}{"Chassis", pluginRequest.URL}, nil)
	}

	data := string(body)
	//replacing the resposne with north bound translation URL
	for key, value := range config.Data.URLTranslation.NorthBoundURL {
		data = strings.Replace(data, key, value, -1)
	}

	var resp dmtfmodel.Chassis
	err = json.Unmarshal([]byte(data), &resp)
	if err != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
	}

	initializeRPCResponse(&r, resp)
	return
}
