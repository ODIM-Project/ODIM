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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/logging"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func newEventHandler(dao *db.DAO) context.Handler {
	return (&eventHandler{
		dao: dao,
	}).handleEvent
}

type eventHandler struct {
	dao *db.DAO
}

func (eh *eventHandler) handleEvent(c iris.Context) {
	raw := new(json.RawMessage)
	err := c.ReadJSON(raw)
	if err != nil {
		c.StatusCode(http.StatusBadRequest)
		_, _ = c.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
		return
	}

	message := new(redfish.Event)
	err = json.Unmarshal([]byte(redfish.Translator.RedfishToODIM(string(*raw))), message)
	if err != nil {
		c.StatusCode(http.StatusBadRequest)
		_, _ = c.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
		return
	}

	for _, e := range message.Events {
		ctx := stdCtx.TODO()
		containedInKey := db.CreateContainedInKey("Chassis", e.OriginOfCondition.Oid)
		rackKey, err := eh.dao.Get(ctx, containedInKey.String()).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			continue
		}

		err = eh.dao.Watch(ctx, func(tx *redis.Tx) error {
			_, err := tx.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
				if _, err := pipeliner.Del(
					ctx,
					containedInKey.String(),
				).Result(); err != nil {
					return fmt.Errorf("del: %s error: %w", containedInKey, err)
				}

				if _, err := pipeliner.SRem(
					ctx,
					db.CreateContainsKey("Chassis", rackKey).String(), e.OriginOfCondition.Oid,
				).Result(); err != nil {
					return fmt.Errorf("srem: %s error: %w", db.CreateContainsKey("Chassis", rackKey).String(), err)
				}
				return nil
			})
			return err
		}, rackKey)

		if err != nil {
			logging.Errorf(
				"cannot consume message(%v): %v",
				message,
				fmt.Errorf("couldn't commit transaction: %w", err),
			)
		}
	}
}
