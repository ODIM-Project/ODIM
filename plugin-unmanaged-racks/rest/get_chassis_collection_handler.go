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
	stdCtx "context"
	"net/http"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/logging"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	"github.com/kataras/iris/v12/context"
)

func newGetChassisCollectionHandler(dao *db.DAO) context.Handler {
	return (&getChassisCollectionHandler{dao}).handle
}

type getChassisCollectionHandler struct {
	dao *db.DAO
}

func (c *getChassisCollectionHandler) handle(ctx context.Context) {
	logging.Debug("Process GetChassisCollection request")
	searchKey := db.CreateKey("Chassis")
	keys, err := c.dao.Keys(stdCtx.TODO(), searchKey.WithWildcard().String()).Result()
	logging.Debug("List of chassis entries:", keys)
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
	}

	collection := createChassisCollection()
	for _, key := range keys {
		collection.Members = append(
			collection.Members,
			redfish.Link{
				Oid: db.Key(key).ID(),
			},
		)
		collection.MembersCount++
	}
	logging.Debug("Populated collection with chassis IDs")

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(&collection)
	logging.Debug("Finished processing GetChassisCollection request")
}

func createChassisCollection() redfish.Collection {
	return redfish.NewCollection("/ODIM/v1/Chassis", "#ChassisCollection.ChassisCollection")
}
