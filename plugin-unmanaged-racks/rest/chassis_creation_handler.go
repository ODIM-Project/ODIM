package rest

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	log "log"
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
					v, err := c.cm.FindByKey("Chassis", containedByOid)
					if err != nil {
						log.Println("error:", err)
					}
					return err != nil || v == nil
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

	v := c.createValidator(requestedChassis)
	validationResult := v.Validate()
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

	dbError := c.cm.Create("Chassis", requestedChassis.Oid, value)
	if dbError != nil {
		switch dbError.Code {
		case db.DB_ERR_ALREADY_EXISTS:
			ctx.StatusCode(http.StatusConflict)
			ctx.JSON(redfish.NewError().AddExtendedInfo(redfish.NewResourceAlreadyExistsMsg("Chassis", "Name", requestedChassis.Name, "")))
		default:
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.JSON(redfish.CreateError(redfish.GeneralError, dbError.Error()))
		}
		return
	}

	if requestedChassis.ChassisType == "Rack" {
		cbUri := requestedChassis.Links.ContainedBy[0].Oid
		bytes, err := redis.Bytes(c.cm.FindByKey("Chassis", cbUri))
		if err != nil {
			ctx.StatusCode(http.StatusInternalServerError)
			return
		}
		if bytes == nil {
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
	ctx.JSON(requestedChassis)
}

func initializeDefaults(rootServiceUUID string, c *redfish.Chassis) {
	//todo: set static part of chassis
	c.ID = uuid.NewV5(uuid.Must(uuid.FromString(rootServiceUUID)), c.Name).String()
	c.Oid = "/ODIM/v1/Chassis/" + c.ID
}
