/*
 * Copyright (c) 2020 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package chassis

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/plugin"
	"github.com/ODIM-Project/ODIM/svc-systems/scommon"
	"github.com/ODIM-Project/ODIM/svc-systems/sresponse"
	"github.com/ODIM-Project/ODIM/svc-systems/systems"
)

var (
	// GetResourceInfoFromDeviceFunc function pointer for the scommon.GetResourceInfoFromDevice
	GetResourceInfoFromDeviceFunc = scommon.GetResourceInfoFromDevice
)

// Handle is used to fetch resource data. The function is supposed to be used as part of RPC
// For getting chassis resource information, parameters need to be passed Request .
// Request holds the Uuid and Url ,
// Url will be parsed from that search key will created
// There will be two return values for the function. One is the RPC response, which contains the
// status code, status message, headers and body and the second value is error.
func (h *Get) Handle(ctx context.Context, req *chassisproto.GetChassisRequest) response.RPC {
	//managed chassis lookup
	l.LogWithFields(ctx).Debugln("Inside GetChassisRequest Handle")
	managedChassis := new(dmtf.Chassis)
	e := h.findInMemoryDB("Chassis", req.URL, managedChassis)
	managedChassis.ID = req.RequestParam
	if e == nil {
		requestData := strings.SplitN(req.RequestParam, ".", 2)
		if len(requestData) <= 1 {
			errorMessage := "error: SystemUUID not found"
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"ComputerSystem", req.RequestParam}, nil)
		}
		uuid := requestData[0]

		var pc = systems.PluginContact{
			ContactClient:   pmbhandle.ContactPlugin,
			DevicePassword:  common.DecryptWithPrivateKey,
			GetPluginStatus: scommon.GetPluginStatus,
		}
		var getDeviceInfoRequest = scommon.ResourceInfoRequest{
			URL:             req.URL,
			UUID:            uuid,
			SystemID:        requestData[1],
			ContactClient:   pc.ContactClient,
			DevicePassword:  pc.DevicePassword,
			GetPluginStatus: pc.GetPluginStatus,
			ResourceName:    "Chassis",
		}
		data, err := GetResourceInfoFromDeviceFunc(ctx, getDeviceInfoRequest, true)
		if err != nil {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"ComputerSystem", req.URL}, nil)
		}
		data = strings.Replace(data, `"Id":"`, `"Id":"`+uuid+`.`, -1)
		var resource dmtf.Chassis
		json.Unmarshal([]byte(data), &resource)
		return response.RPC{
			StatusMessage: response.Success,
			StatusCode:    http.StatusOK,
			Body:          resource,
		}
	}

	if e.ErrNo() != errors.DBKeyNotFound {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}
	l.LogWithFields(ctx).Debugln("Built 'Chassis' table information from lib-dmtf chassis model")
	pluginClient, e := h.createPluginClient("URP*")
	if e != nil && e.ErrNo() == errors.DBKeyNotFound {
		//urp plugin is not registered, requested chassis unknown -> status not found
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "", []interface{}{"Chassis", req.URL}, nil)
	}

	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	resp := pluginClient.Get(ctx, "/ODIM/v1/Chassis/"+req.RequestParam)
	if !is2xx(int(resp.StatusCode)) {
		f := h.getFabricFactory(nil)
		r := f.getFabricChassisResource(ctx, req.RequestParam)
		if is2xx(int(r.StatusCode)) {
			return r
		}
	}
	return resp
}

// Get struct helps to get chassis resource information
type Get struct {
	findInMemoryDB     func(table, key string, r interface{}) *errors.Error
	createPluginClient plugin.ClientFactory
	getFabricFactory   func(collection *sresponse.Collection) *fabricFactory
}

// NewGetHandler returns an instance of Get Struct
func NewGetHandler(
	pluginClientCreator plugin.ClientFactory,
	inMemoryDBFinder func(table, key string, r interface{}) *errors.Error) *Get {

	return &Get{
		createPluginClient: pluginClientCreator,
		findInMemoryDB:     inMemoryDBFinder,
		getFabricFactory:   getFabricFactory,
	}
}
