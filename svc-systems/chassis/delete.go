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
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/plugin"
	"net/http"
)

func (d *Delete) Handle(req *chassisproto.DeleteChassisRequest) response.RPC {
	e := d.findInMemory("Chassis", req.URL, new(json.RawMessage))
	if e == nil {
		return common.GeneralError(http.StatusMethodNotAllowed, response.ActionNotSupported, "Managed Chassis cannot be deleted", []interface{}{"DELETE"}, nil)
	}

	if e.ErrNo() != errors.DBKeyNotFound {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	c, e := d.createPluginClient("URP_v1.0.0")
	if e != nil && e.ErrNo() == errors.DBKeyNotFound {
		return common.GeneralError(http.StatusMethodNotAllowed, response.ActionNotSupported, "", []interface{}{"DELETE"}, nil)
	}
	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	return c.Delete(req.URL)
}

func NewDeleteHandler(createPluginClient plugin.ClientFactory, finder func(Table string, key string, r interface{}) *errors.Error) *Delete {
	return &Delete{
		createPluginClient: createPluginClient,
		findInMemory:       finder,
	}
}

type Delete struct {
	createPluginClient plugin.ClientFactory
	findInMemory       func(Table string, key string, r interface{}) *errors.Error
}
