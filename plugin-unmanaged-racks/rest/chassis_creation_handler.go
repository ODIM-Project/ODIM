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

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/logging"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12/context"
)

func NewCreateChassisHandlerHandler(cm *db.ConnectionManager, pc *config.PluginConfig) context.Handler {
	return (&createChassisHandler{cm, pc}).handle
}

type createChassisHandler struct {
	cm *db.ConnectionManager
	pc *config.PluginConfig
}

func (c *createChassisHandler) createValidator(chassis *redfish.Chassis) *redfish.CompositeValidator {
	return &redfish.CompositeValidator{
		redfish.Validator{
			ValidationRule: func() bool {
				return len(chassis.Name) == 0
			},
			ErrorGenerator: func() redfish.MsgExtendedInfo {
				return redfish.NewPropertyMissingMsg("Name", "cannot be empty")
			},
		},
		redfish.Validator{
			ValidationRule: func() bool {
				return !strings.Contains(strings.Join([]string{"", "RackGroup", "Rack"}, "#"), chassis.ChassisType)
			},
			ErrorGenerator: func() redfish.MsgExtendedInfo {
				return redfish.NewPropertyValueNotInListMsg(chassis.ChassisType, "ChassisType", "supported ChassisTypes are: RackGroup|Rack")
			},
		},
		redfish.Validator{
			ValidationRule: func() bool {
				return chassis.ChassisType == "Rack" && len(chassis.Links.ContainedBy) == 0
			},
			ErrorGenerator: func() redfish.MsgExtendedInfo {
				return redfish.NewPropertyValueConflictMsg(
					"ChassisType", "Links.ContainedBy", "Links.ContainedBy is required for creation of \"ChassisType=Rack\"",
				)
			},
		},
		redfish.Validator{
			ValidationRule: func() bool {
				return chassis.ChassisType == "Rack" && len(chassis.Links.ContainedBy) != 1
			},
			ErrorGenerator: func() redfish.MsgExtendedInfo {
				return redfish.NewPropertyValueConflictMsg(
					"ChassisType", "Links.ContainedBy", "len(Links.ContainedBy) should equal 1",
				)
			},
		},
		redfish.Validator{
			ValidationRule: func() bool {
				if chassis.ChassisType == "Rack" && len(chassis.Links.ContainedBy) == 1 {
					containedByOid := chassis.Links.ContainedBy[0].Oid
					_, err := c.cm.DAO().Get(stdCtx.TODO(), db.CreateKey("Chassis", containedByOid).String()).Result()
					if err != nil || err == redis.Nil {
						logging.Errorf("cannot validate requested rack(%s): %s", chassis.Oid, err)
					}
					return err != nil
				}
				return false
			},
			ErrorGenerator: func() redfish.MsgExtendedInfo {
				return redfish.NewResourceNotFoundMsg(
					"Chassis", chassis.Links.ContainedBy[0].Oid,
					"Requested Links.ContainedBy[0] is unknown")
			},
		},
		redfish.Validator{
			ValidationRule: func() bool {
				return len(chassis.Links.ManagedBy) == 0
			},
			ErrorGenerator: func() redfish.MsgExtendedInfo {
				return redfish.NewPropertyMissingMsg("Links.ManagedBy", "cannot be empty")
			},
		},
		redfish.Validator{
			ValidationRule: func() bool {
				return len(chassis.Links.ManagedBy) != 0 && chassis.Links.ManagedBy[0].Oid != "/ODIM/v1/Managers/"+c.pc.RootServiceUUID
			},
			ErrorGenerator: func() redfish.MsgExtendedInfo {
				return redfish.NewPropertyValueNotInListMsg(
					chassis.Links.ManagedBy[0].Oid,
					"Links.ManagedBy", "should refer to /ODIM/v1/Managers/"+c.pc.RootServiceUUID,
				)
			},
		},
	}
}

func (c *createChassisHandler) handle(ctx context.Context) {
	requestedChassis := new(redfish.Chassis)
	err := ctx.ReadJSON(requestedChassis)
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(redfish.NewError().AddExtendedInfo(redfish.NewMalformedJsonMsg(err.Error())))
		return
	}

	if validationResult := c.createValidator(requestedChassis).Validate(); validationResult.HasErrors() {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(validationResult.Error())
		return
	}

	toBeCreatedChassisId, toBeCreatedBody, parentChassisKey, err := prepareChassisCreationParams(redfish.ShapeChassis(requestedChassis))
	if err != nil {
		return
	}

	sctx := stdCtx.TODO()
	err = c.cm.DAO().Watch(sctx, func(tx *redis.Tx) error {
		exists, err := tx.Exists(sctx, toBeCreatedChassisId).Result()
		if err != nil {
			return err
		}
		if exists == 1 {
			return db.ErrAlreadyExists
		}

		_, err = tx.TxPipelined(sctx, func(pipe redis.Pipeliner) error {
			//create chassis
			if _, err = pipe.SetNX(sctx, toBeCreatedChassisId, toBeCreatedBody, 0).Result(); err != nil {
				return fmt.Errorf("setnx: %s error: %w", toBeCreatedChassisId, err)
			}
			//commit transaction if requested chassis has not parent
			if parentChassisKey == nil {
				return nil
			}

			//set relations between requested and parent
			if _, err = pipe.SAdd(sctx, parentChassisKey.String(), requestedChassis.Oid).Result(); err != nil {
				return fmt.Errorf("sadd: %s error: %w", parentChassisKey, err)
			}

			toBeCreatedContainedInId := db.CreateContainedInKey("Chassis", requestedChassis.Oid).String()
			if _, err = pipe.Set(sctx, toBeCreatedContainedInId, parentChassisKey.Id(), 0).Result(); err != nil {
				return fmt.Errorf("set: %s error: %w", toBeCreatedContainedInId, err)
			}

			return nil
		})
		return err
	}, toBeCreatedChassisId)

	switch err {
	case nil:
		ctx.StatusCode(http.StatusCreated)
		ctx.Header("Location", requestedChassis.Oid)
		ctx.JSON(requestedChassis)
	case db.ErrAlreadyExists:
		ctx.StatusCode(http.StatusConflict)
		ctx.JSON(redfish.NewError().AddExtendedInfo(redfish.NewResourceAlreadyExistsMsg("Chassis", "Name", requestedChassis.Name, "")))
	default:
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
	}
}

func prepareChassisCreationParams(rc *redfish.Chassis) (chassisId string, chassisBody []byte, parentChassisId *db.Key, err error) {
	chassisId = db.CreateKey("Chassis", rc.Oid).String()

	chassisBody, err = json.Marshal(rc)
	if err != nil {
		return
	}

	if len(rc.Links.ContainedBy) > 0 {
		k := db.CreateContainsKey("Chassis", rc.Links.ContainedBy[0].Oid)
		parentChassisId = &k
	}

	return
}
