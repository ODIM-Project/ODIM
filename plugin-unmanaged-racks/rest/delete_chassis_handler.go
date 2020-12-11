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
	"strings"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12/context"
)

func newDeleteChassisHandler(dao *db.DAO) context.Handler {
	return (&deleteChassisHandler{dao}).handle
}

type deleteChassisHandler struct {
	dao *db.DAO
}

func (c *deleteChassisHandler) handle(ctx context.Context) {
	requestedChassis := ctx.Request().RequestURI
	requestedChassisKey := db.CreateKey("Chassis", requestedChassis)

	bytes, err := c.dao.Get(stdCtx.TODO(), requestedChassisKey.String()).Bytes()
	if err != nil && err == redis.Nil {
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(redfish.NewResourceNotFoundMsg("Chassis", requestedChassis, ""))
		return
	}
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
		return
	}

	chassisToBeDeleted := new(redfish.Chassis)
	err = json.Unmarshal(bytes, chassisToBeDeleted)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
		return
	}

	if r := c.createValidator(chassisToBeDeleted).Validate(); r != nil {
		ctx.StatusCode(http.StatusConflict)
		ctx.JSON(redfish.NewError(r...))
		return
	}

	switch chassisToBeDeleted.ChassisType {
	case "RackGroup":
		ctx := stdCtx.TODO()
		err = c.dao.Watch(ctx, func(tx *redis.Tx) error {
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				_, err = pipe.Del(ctx, requestedChassisKey.String()).Result()
				return err
			})
			return err
		}, requestedChassis)

	case "Rack":
		ctx := stdCtx.TODO()

		transactional := func(tx *redis.Tx) error {
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				if _, err = pipe.Del(ctx, requestedChassisKey.String()).Result(); err != nil {
					return fmt.Errorf("del: %s error: %w", requestedChassisKey, err)
				}

				if _, err = pipe.Del(ctx, db.CreateContainedInKey("Chassis", requestedChassis).String()).Result(); err != nil {
					return fmt.Errorf("del: %s error: %w", db.CreateContainedInKey("Chassis", requestedChassis), err)
				}

				parentContainsID := db.CreateContainsKey("Chassis", chassisToBeDeleted.Links.ContainedBy[0].Oid).String()
				_, err = pipe.SRem(ctx, parentContainsID, requestedChassis).Result()
				return err
			})
			return err
		}

		err = c.dao.
			Watch(
				ctx,
				transactional,
				requestedChassis, db.CreateContainedInKey(requestedChassisKey.String()).String(), db.CreateContainsKey(chassisToBeDeleted.Links.ContainedBy[0].Oid).String(),
			)
	}

	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
		return
	}

	ctx.StatusCode(http.StatusNoContent)
}

func (c *deleteChassisHandler) createValidator(chassis *redfish.Chassis) redfish.Validator {
	return redfish.CompositeValidator(
		redfish.NewValidator(
			func() bool {
				return !strings.Contains(strings.Join([]string{"", "RackGroup", "Rack"}, "#"), chassis.ChassisType)
			},
			func() redfish.MsgExtendedInfo {
				return redfish.NewPropertyValueNotInListMsg(chassis.ChassisType, "ChassisType", "supported ChassisTypes are: RackGroup|Rack")
			},
		),
		redfish.NewValidator(
			func() bool {
				hasChildren, err := c.dao.Exists(stdCtx.TODO(), db.CreateContainsKey("Chassis", chassis.Oid).String()).Result()
				return err != nil || hasChildren == 1
			},
			func() redfish.MsgExtendedInfo {
				return redfish.NewResourceInUseMsg("there are existing elements(Links.Contains) under requested chassis")
			},
		),
	)
}
