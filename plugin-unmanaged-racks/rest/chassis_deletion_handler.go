package rest

import (
	"encoding/json"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris/v12/context"
	"net/http"
	"strings"
)

func NewChassisDeletionHandler(cm *db.ConnectionManager) context.Handler {
	return (&chassisDeletionHandler{cm}).handle
}

type chassisDeletionHandler struct {
	connectionManager *db.ConnectionManager
}

func (c *chassisDeletionHandler) handle(ctx context.Context) {
	requestedChassis := ctx.Request().RequestURI
	requestedChassisKey := db.CreateKey("Chassis", requestedChassis)

	conn := c.connectionManager.GetConnection()
	defer db.NewConnectionCloser(&conn)

	bytes, err := redis.Bytes(conn.Do("GET", requestedChassisKey))
	if err != nil && err == redis.ErrNil {
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

	validator := c.createValidator(chassisToBeDeleted)
	r := validator.Validate()
	if r.HasErrors() {
		ctx.StatusCode(http.StatusConflict)
		ctx.JSON(r.Error())
		return
	}

	switch chassisToBeDeleted.ChassisType {
	case "RackGroup":
		err = c.connectionManager.DoInTransaction(func(c redis.Conn) error {
			return conn.Send("DEL", requestedChassisKey)
		}, requestedChassis)

	case "Rack":
		err = c.connectionManager.DoInTransaction(func(c redis.Conn) error {
			err := conn.Send("DEL", requestedChassisKey)
			if err != nil {
				return err
			}
			err = conn.Send("DEL", db.CreateContainedInKey("Chassis", requestedChassis))
			if err != nil {
				return err
			}

			return conn.Send("SREM", db.CreateContainsKey("Chassis", chassisToBeDeleted.Links.ContainedBy[0].Oid), requestedChassis)
		},
			requestedChassis,
			db.CreateContainedInKey(requestedChassisKey.String()).String(),
			db.CreateContainsKey(chassisToBeDeleted.Links.ContainedBy[0].Oid).String())
	}

	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
		return
	}

	ctx.StatusCode(http.StatusOK)
}

func (c *chassisDeletionHandler) createValidator(chassis *redfish.Chassis) *redfish.CompositeValidator {
	return &redfish.CompositeValidator{
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
				conn := c.connectionManager.GetConnection()
				defer db.NewConnectionCloser(&conn)
				hasChildren, err := redis.Bool(conn.Do("EXISTS", db.CreateContainsKey("Chassis", chassis.Oid)))
				return err != nil || hasChildren
			},
			ErrorGenerator: func() redfish.MsgExtendedInfo {
				return redfish.NewResourceInUseMsg("there are existing elements(Links.Contains) under requested chassis")
			},
		},
	}
}
