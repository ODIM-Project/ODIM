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

package rest

import (
	"net/http"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	"github.com/kataras/iris/v12/context"
)

type getPluginManagerHandler struct {
	pluginManager redfish.Manager
}

func (m *getPluginManagerHandler) handle(ctx context.Context) {
	requestedManager := ctx.Request().RequestURI
	if requestedManager == m.pluginManager.OdataID {
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(m.pluginManager)
		return
	}
	ctx.StatusCode(http.StatusNotFound)
}

func createPluginManager(pc *config.PluginConfig) redfish.Manager {
	return redfish.Manager{
		OdataContext:    "/ODIM/v1/$metadata#Manager.Manager",
		OdataID:         "/ODIM/v1/Managers/" + pc.RootServiceUUID,
		OdataType:       "#Manager.v1_3_3.Manager",
		Name:            _PLUGIN_NAME,
		ManagerType:     "Service",
		ID:              pc.RootServiceUUID,
		UUID:            pc.RootServiceUUID,
		FirmwareVersion: pc.FirmwareVersion,
		Status: &redfish.ManagerStatus{
			State: "Enabled",
		},
	}
}

func NewGetPluginManagerHandler(pc *config.PluginConfig) context.Handler {
	return (&getPluginManagerHandler{
		pluginManager: createPluginManager(pc),
	}).handle
}
