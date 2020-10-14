package rest

import (
	"encoding/json"
	"fmt"
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

func NewChassisUpdateHandler(cm *db.ConnectionManager) context.Handler {
	return (&chassisUpdateHandler{cm}).handle
}

type chassisUpdateHandler struct {
	cm *db.ConnectionManager
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

	if vr := validate(requestedChassis, rur); vr.HasErrors() {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(vr.Error())
		return
	}

	conn := c.cm.GetConnection()
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Print("Error: ", err)
		}
	}()

	chassisContainsSetKey := db.CreateChassisContainsKey(requestedChassis.Oid)

	_, err = conn.Do("WATCH", chassisContainsSetKey)
	if err != nil {
		createInternalError(ctx, err)
		return
	}

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

	conn.Send("MULTI")

	//remove known but not requested
	knownMembers.Each(func(knownMember interface{}) bool {
		if !requestedMembers.Contains(knownMember) {
			conn.Send("SREM", chassisContainsSetKey, knownMember)
		}
		return false
	})

	//add requested but unknown
	requestedMembers.Each(func(rm interface{}) bool {
		if !knownMembers.Contains(rm) {
			conn.Send("SADD", chassisContainsSetKey, rm)
		}
		return false
	})

	s, err := conn.Do("EXEC")
	if err != nil {
		createInternalError(ctx, fmt.Errorf("cannot commit transaction: %v", err))
	}
	if s == redis.ErrNil {
		createInternalError(ctx, fmt.Errorf("transaction aborted for unknown reason: %v", err))
	}

	ctx.StatusCode(http.StatusNoContent)
}

func validate(requestedChassis *redfish.Chassis, requestedChange *rackUpdateRequest) redfish.ValidationResult {
	return (&redfish.CompositeValidator{
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
	}).Validate()
}

func (c *chassisUpdateHandler) findRequestedChassis(chassisOid string) (*redfish.Chassis, error) {
	v, err := redis.Bytes(c.cm.FindByKey("Chassis", chassisOid))
	if err != nil {
		switch err {
		case redis.ErrNil:
			//not found
			return nil, nil
		default:
			return nil, fmt.Errorf("%v", redfish.CreateError(redfish.GeneralError, err.Error()))
		}
	}

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
