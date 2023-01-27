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
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	dmtfmodel "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
)

var (
	//JSONUnmarshalFunc ...
	JSONUnmarshalFunc = json.Unmarshal
	//RequestParamsCaseValidatorFunc ...
	RequestParamsCaseValidatorFunc = common.RequestParamsCaseValidator
)

// updateFabricChassisResource will collect the all available fabric plugins available
// in the DB and communicates with each one of them concurrently to update the resource
func (f *fabricFactory) updateFabricChassisResource(ctx context.Context, url string, body *json.RawMessage) response.RPC {
	l.LogWithFields(ctx).Debugf("Inside updateFabricChassisResource for URI: %s", url)
	var resp response.RPC
	ch := make(chan response.RPC)

	managers, err := f.getFabricManagers(ctx)
	if err != nil {
		l.LogWithFields(ctx).Warn("while trying to collect fabric managers details from DB, got " + err.Error())
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "", []interface{}{"Chassis", url}, nil)
	}
	var threadID int = 1
	for _, manager := range managers {
		ctxt := context.WithValue(ctx, common.ThreadName, common.UpdateChassisResource)
		ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
		go f.updateResource(ctxt, manager, url, body, ch)
		threadID++
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
func (f *fabricFactory) updateResource(ctx context.Context, plugin smodel.Plugin, url string, body *json.RawMessage, ch chan response.RPC) {
	l.LogWithFields(ctx).Debugf("Inside updateResource for URI: %s", url)
	req, errResp, err := f.createChassisRequest(ctx, plugin, url, http.MethodPatch, body)
	if errResp != nil {
		l.LogWithFields(ctx).Warn("while trying to create fabric plugin request for " + plugin.ID + ", got " + err.Error())
		ch <- *errResp
		return
	}
	ch <- patchResource(ctx, f, req)
}

// patchResource contacts the plugin with the details available in the
// pluginContactRequest, and returns the RPC response
func patchResource(ctx context.Context, f *fabricFactory, pluginRequest *pluginContactRequest) (r response.RPC) {
	l.LogWithFields(ctx).Debugf("Inside patchResource")
	body, _, statusCode, statusMessage, err := contactPlugin(ctx, pluginRequest)
	if statusCode == http.StatusUnauthorized && strings.EqualFold(pluginRequest.Plugin.PreferredAuthType, "XAuthToken") {
		body, _, statusCode, statusMessage, err = retryFabricsOperation(ctx, f, pluginRequest)
	}
	if err != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
	}
	if !is2xx(statusCode) {
		json.Unmarshal(body, &r.Body)
		r.StatusCode = int32(statusCode)
		r.StatusMessage = statusMessage
		return
	}

	initializeRPCResponse(&r, common.GeneralError(http.StatusOK, response.Success, "", nil, nil))
	return
}

// validating if request properties are in uppercamelcase or not
func validateReqParamsCase(ctx context.Context, req *json.RawMessage) *response.RPC {
	var errResp response.RPC
	var chassisRequest dmtfmodel.Chassis

	// parsing the fabricRequest
	err := JSONUnmarshalFunc(*req, &chassisRequest)
	if err != nil {
		errResp = common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, err.Error(), nil, nil)
		return &errResp
	}

	// validating the request JSON properties for case sensitive
	invalidProperties, err := RequestParamsCaseValidatorFunc(*req, chassisRequest)
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
