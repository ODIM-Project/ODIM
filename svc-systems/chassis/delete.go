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
	"reflect"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/plugin"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
)

var (
	//SmodelFindAllFunc ...
	SmodelFindAllFunc = smodel.FindAll
	//FindAllPluginsFunc ...
	FindAllPluginsFunc = findAllPlugins
)

// Handle defines the operations which handle the RPC request-response for deleting a chassis
func (d *Delete) Handle(ctx context.Context, req *chassisproto.DeleteChassisRequest) response.RPC {
	e := d.findInMemory("Chassis", req.URL, new(json.RawMessage))
	if e == nil {
		return common.GeneralError(http.StatusMethodNotAllowed, response.ActionNotSupported, "Managed Chassis cannot be deleted", []interface{}{"DELETE"}, nil)
	}

	if e.ErrNo() != errors.DBKeyNotFound {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	//TODO: Handle multiple URP instances
	c, e := d.createPluginClient("URP*")
	if e != nil && e.ErrNo() == errors.DBKeyNotFound {
		return common.GeneralError(http.StatusMethodNotAllowed, response.ActionNotSupported, "", []interface{}{"DELETE"}, nil)
	}
	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	plugins, err := FindAllPluginsFunc("URP*")
	if err != nil {
		errorMessage := "error while getting plugin details: " + err.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
			nil, nil)
	}
	managerURI := "/redfish/v1/Managers/" + plugins[0].ManagerUUID

	data, jerr := GetResourceFunc(ctx, "Managers", managerURI)
	if jerr != nil {
		errorMessage := "error while getting manager details: " + jerr.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
			nil, nil)
	}
	var managerData map[string]interface{}
	err = JSONUnmarshalFunc([]byte(data), &managerData)
	if err != nil {
		errorMessage := "error unmarshalling manager details: " + err.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
			nil, nil)
	}

	if links, ok := managerData["Links"].(map[string]interface{}); ok {
		if managerForChassis, ok := links["ManagerForChassis"].([]interface{}); ok {
			for k, v := range managerForChassis {
				if v.(map[string]interface{})["@odata.id"] != nil {
					if reflect.DeepEqual(v.(map[string]interface{})["@odata.id"], req.URL) {
						managerForChassis = append(managerForChassis[:k], managerForChassis[k+1:]...)
						if len(managerForChassis) != 0 {
							links["ManagerForChassis"] = managerForChassis
						} else {
							delete(links, "ManagerForChassis")
						}
					}
				}
			}
		}
	}
	detail, marshalErr := JSONMarshalFunc(managerData)
	if marshalErr != nil {
		errorMessage := "unable to marshal data for updating: " + marshalErr.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}

	genericErr := GenericSaveFunc(ctx, []byte(detail), "Managers", managerURI)
	if genericErr != nil {
		errorMessage := "GenericSave : error while trying to add resource date to DB: " + genericErr.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	l.LogWithFields(ctx).Debugf("Deleting URI: %s", string(req.URL))
	return c.Delete(ctx, req.URL)
}

// NewDeleteHandler returns an instance of Delete struct
func NewDeleteHandler(createPluginClient plugin.ClientFactory, finder func(Table string, key string, r interface{}) *errors.Error) *Delete {
	return &Delete{
		createPluginClient: createPluginClient,
		findInMemory:       finder,
	}
}

// Delete struct helps to delete chassis
type Delete struct {
	createPluginClient plugin.ClientFactory
	findInMemory       func(Table string, key string, r interface{}) *errors.Error
}

func findAllPlugins(key string) (res []*smodel.Plugin, err error) {
	pluginsAsBytesSlice, err := SmodelFindAllFunc("Plugin", key)
	if err != nil {
		return nil, err
	}

	for _, bytes := range pluginsAsBytesSlice {
		plugin := new(smodel.Plugin)
		err = JSONUnmarshalFunc(bytes, plugin)
		if err != nil {
			return nil, err
		}
		res = append(res, plugin)
	}

	return
}
