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
	"encoding/json"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/plugin"
	"github.com/ODIM-Project/ODIM/svc-systems/sresponse"
)

// Handle defines the operations which handle the RPC request-response for updating a chassis
func (h *Update) Handle(req *chassis.UpdateChassisRequest) response.RPC {
	pc, e := h.createPluginClient("URP*")
	if e != nil && e.ErrNo() == errors.DBKeyNotFound {
		return common.GeneralError(http.StatusMethodNotAllowed, response.ActionNotSupported, "", []interface{}{"PATCH"}, nil)
	}
	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	body := new(json.RawMessage)
	ue := json.Unmarshal(req.RequestBody, body)
	if ue != nil {
		l.Log.Error("while trying to unmarshal, got " + ue.Error())
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, "Cannot deserialize request body", nil, nil)
	}

	resp := pc.Patch(req.URL, body)
	if !is2xx(int(resp.StatusCode)) {
		f := h.getFabricFactory(nil)
		r := f.updateFabricChassisResource(req.URL, body)
		if is2xx(int(r.StatusCode)) || r.StatusCode == http.StatusBadRequest || r.StatusCode == http.StatusUnauthorized {
			return r
		}
	}

	return resp
}

// Update struct helps to update chassis
type Update struct {
	createPluginClient plugin.ClientFactory
	getFabricFactory   func(collection *sresponse.Collection) *fabricFactory
}

// NewUpdateHandler returns an instance of Update struct
func NewUpdateHandler(pluginClientFactory plugin.ClientFactory) *Update {
	return &Update{
		createPluginClient: pluginClientFactory,
		getFabricFactory:   getFabricFactory,
	}
}
