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
	"net/http"
	"strings"

	dmtfmodel "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	log "github.com/sirupsen/logrus"
)

// updateFabricChassisResource will collect the all available fabric plugins available
// in the DB and communicates with each one of them concurrently to update the resource
func (f *fabricFactory) updateFabricChassisResource(url string, body *json.RawMessage) response.RPC {
	var resp response.RPC
	ch := make(chan response.RPC)

	managers, err := f.getFabricManagers()
	if err != nil {
		log.Warn("while trying to collect fabric managers details from DB, got " + err.Error())
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "", []interface{}{"Chassis", url}, nil)
	}

	for _, manager := range managers {
		go f.updateResource(manager, url, body, ch)
	}

	for i := 0; i < len(managers); i++ {
		resp = <-ch
		if is2xx(int(resp.StatusCode)) {
			return resp
		}
	}
	return resp
}

// updateResource will validate the request body, creates the request model for communicating
// with the plugin and returns the response
func (f *fabricFactory) updateResource(plugin smodel.Plugin, url string, body *json.RawMessage, ch chan response.RPC) {
	req, errResp, err := f.createChassisRequest(plugin, url, http.MethodPatch, body)
	if errResp != nil {
		log.Warn("while trying to create fabric plugin request for " + plugin.ID + ", got " + err.Error())
		ch <- *errResp
		return
	}
	ch <- patchResource(f, req)
}

// patchResource contacts the plugin with the details available in the
// pluginContactRequest, and returns the RPC response
func patchResource(f *fabricFactory, pluginRequest *pluginContactRequest) (r response.RPC) {
	body, _, statusCode, statusMessage, err := contactPlugin(pluginRequest)
	if statusCode == http.StatusUnauthorized && strings.EqualFold(pluginRequest.Plugin.PreferredAuthType, "XAuthToken") {
		body, _, statusCode, statusMessage, err = retryFabricsOperation(f, pluginRequest)
	}
	if err != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
	}
	if !is2xx(statusCode) {
		json.Unmarshal(body, &r.Body)
		r.Header = map[string]string{"Content-type": "application/json; charset=utf-8"}
		r.StatusCode = int32(statusCode)
		r.StatusMessage = statusMessage
		return
	}

	initializeRPCResponse(&r, common.GeneralError(http.StatusOK, response.Success, "", nil, nil))
	return
}

// validating if request properties are in uppercamelcase or not
func validateReqParamsCase(req *json.RawMessage) *response.RPC {
	var errResp response.RPC
	var chassisRequest dmtfmodel.Chassis

	// parsing the fabricRequest
	err := json.Unmarshal(*req, &chassisRequest)
	if err != nil {
		errResp = common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, err.Error(), nil, nil)
		return &errResp
	}

	// validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(*req, chassisRequest)
	if err != nil {
		errResp = common.GeneralError(http.StatusInternalServerError, response.InternalError, "error while validating request parameters: "+err.Error(), nil, nil)
		return &errResp
	}

	if invalidProperties != "" {
		errMsg := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase"
		errResp = common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errMsg, []interface{}{invalidProperties}, nil)
		return &errResp
	}

	return nil
}
