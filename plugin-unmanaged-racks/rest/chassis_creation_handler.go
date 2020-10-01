package rest

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"
	"github.com/kataras/iris/v12/context"
	uuid "github.com/satori/go.uuid"
)

func NewCreateChassisHandlerHandler(cm *db.ConnectionManager, pc *config.PluginConfig) context.Handler {
	return (&createChassisHandler{cm, pc}).handle
}

type createChassisHandler struct {
	cm *db.ConnectionManager
	pc *config.PluginConfig
}

func (c *createChassisHandler) createValidator(chassis *redfish.Chassis) *compositeValidator {
	return &compositeValidator{
		validator{
			validationRule: func() bool {
				return len(chassis.Name) == 0
			},
			field:   "Name",
			message: "cannot be empty",
		},
		validator{
			validationRule: func() bool {
				return !strings.Contains(strings.Join([]string{"", "RackGroup", "Rack"}, "#"), chassis.ChassisType)
			},
			field:   "ChassisType",
			message: "supported ChassisTypes are: RackGroup|Rack",
		},
		validator{
			validationRule: func() bool {
				return chassis.ChassisType == "Rack" && len(chassis.Links.ContainedBy) == 0
			},
			field:   "ChassisType",
			message: "Links.ContainedBy is required for creation of \"ChassisType=Rack\"",
		},
		validator{
			validationRule: func() bool {
				return chassis.ChassisType == "Rack" && len(chassis.Links.ContainedBy) != 1
			},
			field:   "ChassisType",
			message: "len(Links.ContainedBy) should equal 1",
		},
		validator{
			validationRule: func() bool {
				if chassis.ChassisType == "Rack" && len(chassis.Links.ContainedBy) == 1 {
					containedByOid := chassis.Links.ContainedBy[0].Oid
					_, err := c.cm.FindByKey("Chassis", containedByOid)
					return err != nil
				}
				return false
			},
			field:   "ChassisType",
			message: "Requested Links.ContainedBy[0] is unknown",
		},
		validator{
			validationRule: func() bool {
				return len(chassis.Links.ManagedBy) == 0
			},
			field:   "Links.ManagedBy",
			message: "cannot be empty",
		},
		validator{
			validationRule: func() bool {
				return len(chassis.Links.ManagedBy) != 0 && chassis.Links.ManagedBy[0].Oid != "/ODIM/v1/Managers/"+c.pc.RootServiceUUID
			},
			field:   "Links.ManagedBy",
			message: "should refer to /ODIM/v1/Managers/" + c.pc.RootServiceUUID,
		},
	}
}

func (c *createChassisHandler) handle(ctx context.Context) {
	requestedChassis := new(redfish.Chassis)
	err := ctx.ReadJSON(requestedChassis)
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		return
	}

	v := c.createValidator(requestedChassis)
	validationResult := v.validate()
	if validationResult.HasErrors() {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(validationResult.Error())
		return
	}

	initializeDefaults(c.pc.RootServiceUUID, requestedChassis)

	value, err := json.Marshal(requestedChassis)
	if err != nil {
		return
	}
	err = c.cm.Create("Chassis", requestedChassis.Oid, value)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		//ctx.JSON(errors.CreateErrResp(errors.InternalError, err.Error()))
		return
	}

	if requestedChassis.ChassisType == "Rack" {
		cbUri := requestedChassis.Links.ContainedBy[0].Oid
		bytes, err := redis.Bytes(c.cm.FindByKey("Chassis", cbUri))
		if err != nil {
			ctx.StatusCode(http.StatusInternalServerError)
			return
		}
		containedByChassis := new(redfish.Chassis)
		err = json.Unmarshal(bytes, containedByChassis)
		if err != nil {
			ctx.StatusCode(http.StatusInternalServerError)
			return
		}

		containedByChassis.Links.Contains = append(containedByChassis.Links.Contains, redfish.Link{Oid: requestedChassis.Oid})
		bytes, err = json.Marshal(containedByChassis)
		if err != nil {
			ctx.StatusCode(http.StatusInternalServerError)
			return
		}
		err = c.cm.Update("Chassis", cbUri, bytes)
		if err != nil {
			ctx.StatusCode(http.StatusInternalServerError)
			return
		}
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.Header("Location", requestedChassis.Oid)
}

func initializeDefaults(rootServiceUUID string, c *redfish.Chassis) {
	//todo: set static part of chassis
	c.ID = uuid.NewV5(uuid.Must(uuid.FromString(rootServiceUUID)), c.Name).String()
	c.Oid = "/ODIM/v1/Chassis/" + c.ID
}
