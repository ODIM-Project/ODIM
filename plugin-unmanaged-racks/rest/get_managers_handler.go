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

type getManagersHandler struct {
	managersCollection redfish.Collection
	pc                 *config.PluginConfig
}

func (m *getManagersHandler) handle(ctx context.Context) {
	ctx.JSON(m.managersCollection)
	ctx.StatusCode(http.StatusOK)
}

func newGetManagersHandler(pc *config.PluginConfig) context.Handler {
	collection := redfish.NewCollection("/ODIM/v1/Managers", "#ManagerCollection.ManagerCollection", redfish.Link{Oid: "/ODIM/v1/Managers/" + pc.RootServiceUUID})
	return (&getManagersHandler{collection, pc}).handle
}
