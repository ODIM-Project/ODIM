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

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	"github.com/kataras/iris/v12/context"
)

func newGetChassisHandler(dao *db.DAO) context.Handler {
	return (&getChassisHandler{dao}).handle
}

type getChassisHandler struct {
	dao *db.DAO
}

func (c *getChassisHandler) handle(ctx context.Context) {
	requestedChassisOid := ctx.Request().RequestURI
	chassis, err := c.dao.FindChassis(requestedChassisOid)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
		return
	}

	if chassis == nil {
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(redfish.NewError().AddExtendedInfo(redfish.NewResourceNotFoundMsg("Chassis", requestedChassisOid, "")))
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(chassis)
}
