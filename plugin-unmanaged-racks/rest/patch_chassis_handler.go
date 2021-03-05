/*
 * Copyright (c) 2020 Intel Corporation
 * (C) Copyright [2020] Hewlett Packard Enterprise Development LP
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

func newPatchChassisHandler(dao *db.DAO, odimHTTPClient *redfish.HTTPClient) context.Handler {
	return (&chassisUpdateHandler{
		dao:           dao,
		redfishClient: redfish.NewResponseWrappingClient(odimHTTPClient),
	}).handle
}

type chassisUpdateHandler struct {
	dao           *db.DAO
	redfishClient *redfish.ResponseWrappingClient
}

func (c *chassisUpdateHandler) handle(ctx context.Context) {
	rur, err := decodeRequestBody(ctx)
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(redfish.NewError().AddExtendedInfo(redfish.NewMalformedJSONMsg(err.Error())))
		return
	}

	chassisURI := ctx.Request().RequestURI
	logging.Debug("finding requested chassis with uri:", chassisURI)
	requestedChassis, err := c.dao.FindChassis(chassisURI)
	if err != nil {
		createInternalError(ctx, err)
		return
	}
	if requestedChassis == nil {
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(redfish.NewError().AddExtendedInfo(redfish.NewResourceNotFoundMsg("Chassis", ctx.Request().RequestURI, "")))
		return
	}
	logging.Debug("found requested chassis with uri:", chassisURI)

	if violation, errCode := c.createValidator(requestedChassis, rur).Validate(); violation != nil {
		ctx.StatusCode(*errCode)
		ctx.JSON(redfish.NewError(*violation))
		return
	}

	chassisContainsSetKey := db.CreateContainsKey("Chassis", requestedChassis.Oid)
	existingMembers, err := c.dao.SMembers(stdCtx.TODO(), chassisContainsSetKey.String()).Result()
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
	err = c.dao.Watch(
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

	updatedChassis, err := c.dao.FindChassis(ctx.Request().RequestURI)
	if err != nil || updatedChassis == nil {
		createInternalError(ctx, err)
		return
	}

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(updatedChassis)
}

func (c *chassisUpdateHandler) createValidator(requestedChassis *redfish.Chassis, requestedChange *rackUpdateRequest) redfish.Validator {
	return redfish.CompositeValidator(
		redfish.NewValidator(
			func() bool {
				return !strings.Contains(strings.Join([]string{"", "Rack"}, "#"), requestedChassis.ChassisType)
			},
			func() (redfish.MsgExtendedInfo, int) {
				return redfish.NewPropertyValueNotInListMsg(requestedChassis.ChassisType, "ChassisType", "supported ChassisTypes are: Rack"), http.StatusBadRequest
			},
		),
		redfish.NewValidator(
			func() bool {
				if len(requestedChassis.Links.ContainedBy) > 0 {
					parent := requestedChassis.Links.ContainedBy[0].Oid
					for _, l := range requestedChange.Links.Contains {
						return l.Oid == parent
					}
				}
				return false
			},
			func() (redfish.MsgExtendedInfo, int) {
				return redfish.NewPropertyValueConflictMsg("Links.Contains", "Links.ContainedBy", "RackGroup cannot be attached under Rack chassis"), http.StatusBadRequest
			},
		),
		redfish.NewValidator(
			func() bool {
				chassisCollection := new(redfish.Collection)
				if err := c.redfishClient.Get("/redfish/v1/Chassis", chassisCollection); err != nil {
					logging.Errorf("cannot read validate https://ODIMRA/redfish/v1/Chassis: %s", err)
					return true
				}
				existingChassis := map[string]interface{}{}
				for _, m := range chassisCollection.Members {
					existingChassis[m.Oid] = m
				}
				for _, assetUnderChassis := range requestedChange.Links.Contains {
					_, ok := existingChassis[assetUnderChassis.Oid]
					if !ok {
						return true
					}
				}
				return false
			},
			func() (redfish.MsgExtendedInfo, int) {
				return redfish.NewPropertyValueNotInListMsg(
						fmt.Sprintf("%s", requestedChange.Links.Contains),
						"Links.Contains",
						"Couldn't confirm existence of one or more requested 'Links.Contains' elements or one of them is not a pointer to Chassis asset"),
					http.StatusBadRequest
			},
		),
		redfish.NewValidator(
			func() bool {
				unmanagedChassisKeys, err := c.dao.Keys(stdCtx.TODO(), db.CreateKey("Chassis").WithWildcard().String()).Result()
				if err != nil {
					logging.Errorf("DB: cannot read info about unmanaged racks: %s", err)
					return true
				}

				for _, assetUnderChassis := range requestedChange.Links.Contains {
					if utils.Collection(unmanagedChassisKeys).Contains(db.CreateKey("Chassis", assetUnderChassis.Oid).String()) {
						return true
					}
				}

				return false
			},
			func() (redfish.MsgExtendedInfo, int) {
				return redfish.NewPropertyValueNotInListMsg(
						fmt.Sprintf("%s", requestedChange.Links.Contains),
						"Links.Contains",
						"Unmanaged Chassis Cannot Be Attached Under Rack"),
					http.StatusBadRequest
			},
		),

		redfish.NewValidator(
			func() bool {

				alreadyAttached := []string{}
				for _, c := range requestedChassis.Links.Contains {
					alreadyAttached = append(alreadyAttached, c.Oid)
				}

				for _, assetUnderChassis := range requestedChange.Links.Contains {
					if utils.Collection(alreadyAttached).Contains(assetUnderChassis.Oid) {
						continue
					}
					if c.dao.Exists(stdCtx.TODO(), db.CreateContainedInKey("Chassis", assetUnderChassis.Oid).String()).Val() == 1 {
						return true
					}
				}

				return false
			},
			func() (redfish.MsgExtendedInfo, int) {
				return redfish.NewPropertyValueNotInListMsg(
						fmt.Sprintf("%s", requestedChange.Links.Contains),
						"Links.Contains",
						"One of requested 'Links.Contains' element is already attached to another rack"),
					http.StatusConflict
			},
		),
	)
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
