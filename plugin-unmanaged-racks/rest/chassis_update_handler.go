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
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/utils"

	"github.com/deckarep/golang-set"
	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12/context"
)

type rackUpdateRequest struct {
	Links struct {
		Contains []redfish.Link
	}
}

func NewChassisUpdateHandler(cm *db.ConnectionManager, c *config.PluginConfig) context.Handler {
	return (&chassisUpdateHandler{cm: cm, config: c}).handle
}

type chassisUpdateHandler struct {
	cm     *db.ConnectionManager
	config *config.PluginConfig
}

func (c *chassisUpdateHandler) handle(ctx context.Context) {
	rur, err := decodeRequestBody(ctx)
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(redfish.NewError().AddExtendedInfo(redfish.NewMalformedJsonMsg(err.Error())))
		return
	}

	requestedChassis, err := c.findRequestedChassis(ctx.Request().RequestURI)
	if err != nil {
		createInternalError(ctx, err)
		return
	}
	if requestedChassis == nil {
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(redfish.NewError().AddExtendedInfo(redfish.NewResourceNotFoundMsg("Chassis", ctx.Request().RequestURI, "")))
		return
	}

	if vr := c.createValidator(requestedChassis, rur).Validate(); vr.HasErrors() {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(vr.Error())
		return
	}

	chassisContainsSetKey := db.CreateContainsKey("Chassis", requestedChassis.Oid)
	existingMembers, err := c.cm.DAO().SMembers(stdCtx.TODO(), chassisContainsSetKey.String()).Result()
	if err != nil {
		createInternalError(ctx, err)
		return
	}

	requestedMembers := mapset.NewSet()
	for _, e := range rur.Links.Contains {
		requestedMembers.Add(e.Oid)
	}

	knownMembers := mapset.NewSet()
	for _, e := range existingMembers {
		knownMembers.Add(e)
	}

	sctx := stdCtx.TODO()
	err = c.cm.DAO().Watch(
		sctx,
		func(tx *redis.Tx) error {
			//remove known but not requested
			_, err := tx.TxPipelined(sctx, func(pipe redis.Pipeliner) error {
				knownMembers.Each(func(knownMember interface{}) bool {
					if !requestedMembers.Contains(knownMember) {
						if _, err = pipe.SRem(sctx, chassisContainsSetKey.String(), knownMember).Result(); err != nil {
							err = fmt.Errorf("srem: %s error: %w", chassisContainsSetKey.String(), err)
							return true
						}

						if _, err = pipe.Del(
							sctx,
							db.CreateContainedInKey("Chassis", knownMember.(string)).String(),
						).Result(); err != nil {
							err = fmt.Errorf("del: %s error: %w", db.CreateContainedInKey("Chassis", knownMember.(string)).String(), err)
							return true
						}
					}
					return false
				})

				//add requested but unknown
				requestedMembers.Each(func(rm interface{}) bool {
					if !knownMembers.Contains(rm) {
						if _, err = pipe.SAdd(sctx, chassisContainsSetKey.String(), rm).Result(); err != nil {
							err = fmt.Errorf("sadd: %s error: %w", chassisContainsSetKey.String(), err)
							return true
						}

						ckey := db.CreateContainedInKey("Chassis", rm.(string)).String()
						if _, err = pipe.Set(sctx, ckey, requestedChassis.Oid, 0).Result(); err != nil {
							err = fmt.Errorf("set: %s error: %w", ckey, err)
							return true
						}
					}
					return false
				})
				return nil
			})
			return err
		},
		chassisContainsSetKey.String(),
	)

	if err != nil {
		createInternalError(ctx, fmt.Errorf("cannot commit transaction: %v", err))
		return
	}

	updatedChassis, err := c.cm.FindChassis(ctx.Request().RequestURI)
	if err != nil || updatedChassis == nil {
		createInternalError(ctx, err)
		return
	}

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(updatedChassis)
}

func (c *chassisUpdateHandler) createValidator(requestedChassis *redfish.Chassis, requestedChange *rackUpdateRequest) *redfish.CompositeValidator {
	return &redfish.CompositeValidator{
		redfish.Validator{
			ValidationRule: func() bool {
				return !strings.Contains(strings.Join([]string{"", "Rack"}, "#"), requestedChassis.ChassisType)
			},
			ErrorGenerator: func() redfish.MsgExtendedInfo {
				return redfish.NewPropertyValueNotInListMsg(requestedChassis.ChassisType, "ChassisType", "supported ChassisTypes are: Rack")
			},
		},
		redfish.Validator{
			ValidationRule: func() bool {
				if len(requestedChassis.Links.ContainedBy) > 0 {
					parent := requestedChassis.Links.ContainedBy[0].Oid
					for _, l := range requestedChange.Links.Contains {
						return l.Oid == parent
					}
				}
				return false
			},
			ErrorGenerator: func() redfish.MsgExtendedInfo {
				return redfish.NewPropertyValueConflictMsg(
					"Links.Contains", "Links.ContainedBy", "RackGroup cannot be attached under Rack chassis",
				)
			},
		},
		redfish.Validator{
			ValidationRule: func() bool {
				chassisCollection := new(redfish.Collection)
				if err := redfish.NewRedfishClient(c.config.OdimNBUrl).Get("/redfish/v1/Chassis", chassisCollection); err != nil {
					logging.Errorf("cannot validate requested chassis: %s", err)
					return true
				}
				existingSystems := map[string]interface{}{}
				for _, m := range chassisCollection.Members {
					existingSystems[m.Oid] = m
				}
				for _, assetUnderChassis := range requestedChange.Links.Contains {
					_, ok := existingSystems[assetUnderChassis.Oid]
					if !ok {
						return true
					}
				}

				unmanagedChassisKeys, err := c.cm.DAO().Keys(stdCtx.TODO(), db.CreateKey("Chassis").WithWildcard().String()).Result()
				if err != nil {
					logging.Errorf("cannot validate requested chassis: %s", err)
					return true
				}

				for _, assetUnderChassis := range requestedChange.Links.Contains {
					if utils.Collection(unmanagedChassisKeys).Contains(db.CreateKey("Chassis", assetUnderChassis.Oid).String()) {
						return true
					}
				}

				return false
			},
			ErrorGenerator: func() redfish.MsgExtendedInfo {
				return redfish.NewPropertyValueNotInListMsg(
					fmt.Sprintf("%s", requestedChange.Links.Contains),
					"Links.Contains",
					"Couldn't retrieve information about requested links. Make sure that they are existing!")
			},
		},
	}
}

func (c *chassisUpdateHandler) findRequestedChassis(chassisOid string) (*redfish.Chassis, error) {
	reply, err := c.cm.DAO().Get(stdCtx.TODO(), db.CreateKey("Chassis", chassisOid).String()).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("%v", redfish.CreateError(redfish.GeneralError, err.Error()))
	}

	requestedChassis := new(redfish.Chassis)
	if err := json.Unmarshal(reply, requestedChassis); err != nil {
		return nil, fmt.Errorf("%v", redfish.CreateError(redfish.GeneralError, err.Error()))
	}

	return requestedChassis, nil
}

func decodeRequestBody(ctx context.Context) (*rackUpdateRequest, error) {
	rur := new(rackUpdateRequest)
	dec := json.NewDecoder(ctx.Request().Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(rur); err != nil {
		return nil, err
	}
	return rur, nil
}

func createInternalError(ctx context.Context, err error) {
	ctx.StatusCode(http.StatusInternalServerError)
	ctx.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
}
