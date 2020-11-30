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
	"fmt"
	"net/http"
	"strconv"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/plugin"
)

func (h *Create) Handle(req *chassisproto.CreateChassisRequest) response.RPC {
	mbc := new(linksManagedByCollection)
	e := json.Unmarshal(req.RequestBody, mbc)
	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, fmt.Sprintf("cannot deserialize request: %v", e), nil, nil)
	}

	if len(mbc.Links.ManagedBy) == 0 {
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, "", []interface{}{"Links.ManagedBy[0]"}, nil)
	}

	inMemoryConn, dbErr := common.GetDBConnection(common.InMemory)
	if dbErr != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, fmt.Sprintf("cannot acquire database connection: %v", dbErr), nil, nil)
	}

	managingManager, e := findOrNull(inMemoryConn, "Managers", mbc.Links.ManagedBy[0].Oid)
	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, fmt.Sprintf("error occured during database access: %v", e), nil, nil)
	}

	if managingManager == "" {
		return common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, "", []interface{}{"Manager", mbc.Links.ManagedBy[0].Oid}, nil)
	}

	//todo: not sure why manager in redis is quoted
	managingManager, e = strconv.Unquote(managingManager)
	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	pluginManagingManager := new(nameCarrier)
	e = json.Unmarshal([]byte(managingManager), pluginManagingManager)
	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	body := new(json.RawMessage)
	e = json.Unmarshal(req.RequestBody, body)
	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	pc, pe := h.createPluginClient(pluginManagingManager.Name)
	if pe != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, pe.Error(), nil, nil)
	}

	return pc.Post("/redfish/v1/Chassis", body)
}

type Create struct {
	createPluginClient plugin.ClientFactory
}

func NewCreateHandler(createPluginClient plugin.ClientFactory) *Create {
	return &Create{
		createPluginClient: createPluginClient,
	}
}

//{
//	"Links" : {
//		"ManagedBy": [
//			"@odata.id": "/redfish/v1/Managers/1"
//		]
//	}
//}
type linksManagedByCollection struct {
	Links struct {
		ManagedBy []struct {
			Oid string `json:"@odata.id"`
		}
	}
}

//{
//	"Name" : "name"
//}
type nameCarrier struct {
	Name string
}
