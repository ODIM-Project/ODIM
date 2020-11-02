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

	bytes, err := redis.Bytes(c.connectionManager.FindByKey(c.connectionManager.CreateKey("Chassis", requestedChassis)))
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

	_, err = c.connectionManager.Delete("Chassis", requestedChassis)
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
				return len(chassis.Links.Contains) != 0
			},
			ErrorGenerator: func() redfish.MsgExtendedInfo {
				return redfish.NewResourceInUseMsg("there are existing elements(Links.Contains) under requested chassis")
			},
		},
	}
}
