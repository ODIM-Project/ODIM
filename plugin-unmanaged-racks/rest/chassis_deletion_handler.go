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

	bytes, err := redis.Bytes(c.connectionManager.FindByKey("Chassis", requestedChassis))
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(err)
		return
	}

	chassisToBeDeleted := new(redfish.Chassis)
	err = json.Unmarshal(bytes, chassisToBeDeleted)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(err)
		return
	}

	validator := c.createValidator(chassisToBeDeleted)
	r := validator.validate()
	if r.HasErrors() {
		ctx.StatusCode(http.StatusConflict)
		ctx.JSON(r.Error())
		return
	}

	ok, err := c.connectionManager.Delete("Chassis", requestedChassis)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(err)
		return
	}
	if !ok {
		ctx.StatusCode(http.StatusBadRequest)
		return
	}

	ctx.StatusCode(http.StatusOK)
}

func (c *chassisDeletionHandler) createValidator(chassis *redfish.Chassis) *compositeValidator {
	return &compositeValidator{
		validator{
			validationRule: func() bool {
				return !strings.Contains(strings.Join([]string{"", "RackGroup", "Rack"}, "#"), chassis.ChassisType)
			},
			field:   "ChassisType",
			message: "supported ChassisTypes are: RackGroup|Rack",
		},
		validator{
			validationRule: func() bool {
				return len(chassis.Links.Contains) != 0
			},
			field:   "Links.Contains",
			message: "there are existing elements under requested chassis",
		},
	}
}
