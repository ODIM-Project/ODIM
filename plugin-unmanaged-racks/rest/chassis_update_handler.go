package rest

import (
	"encoding/json"
	"fmt"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"log"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	mapset "github.com/deckarep/golang-set"
	"github.com/gomodule/redigo/redis"
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

	conn := c.cm.GetConnection()
	defer db.NewConnectionCloser(&conn)

	chassisContainsSetKey := db.CreateContainsKey("Chassis", requestedChassis.Oid)

	existingMembers, err := redis.Strings(conn.Do("SMEMBERS", chassisContainsSetKey))
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

	err = c.cm.DoInTransaction(func(conn redis.Conn) error {
		//remove known but not requested
		knownMembers.Each(func(knownMember interface{}) bool {
			if !requestedMembers.Contains(knownMember) {
				//todo: handle potential errors returned by send
				conn.Send("SREM", chassisContainsSetKey, knownMember)
				conn.Send("DEL", db.CreateContainedInKey("Chassis", knownMember.(string)))
			}
			return false
		})

		//add requested but unknown
		requestedMembers.Each(func(rm interface{}) bool {
			if !knownMembers.Contains(rm) {
				conn.Send("SADD", chassisContainsSetKey, rm)
				conn.Send("SET", db.CreateContainedInKey("Chassis", rm.(string)), requestedChassis.Oid)
			}
			return false
		})

		return nil
	}, chassisContainsSetKey.String())

	if err != nil {
		createInternalError(ctx, fmt.Errorf("cannot commit transaction: %v", err))
		return
	}

	ctx.StatusCode(http.StatusNoContent)
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
				c := redfish.NewRedfishClient(c.config.OdimraNBUrl)
				systemsCollection := new(redfish.Collection)
				err := c.Get("/redfish/v1/Systems", systemsCollection)
				if err != nil {
					log.Print("Couldn't GET systems collection(/redfish/v1/Systems) Error: ", err)
					return true
				}
				existingSystems := map[string]interface{}{}
				for _, m := range systemsCollection.Members {
					existingSystems[m.Oid] = m
				}

				for _, assetUnderChassis := range requestedChange.Links.Contains {
					_, ok := existingSystems[assetUnderChassis.Oid]
					if !ok {
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
	conn := c.cm.GetConnection()
	defer db.NewConnectionCloser(&conn)

	reply, err := conn.Do("GET", db.CreateKey("Chassis", chassisOid))
	if err != nil {
		return nil, fmt.Errorf("%v", redfish.CreateError(redfish.GeneralError, err.Error()))
	}
	if reply == nil {
		return nil, nil
	}

	v, err := redis.Bytes(reply, err)
	requestedChassis := new(redfish.Chassis)
	if err := json.Unmarshal(v, requestedChassis); err != nil {
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
