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
					conn := c.cm.GetConnection()
					defer db.NewConnectionCloser(&conn)
					v, err := conn.Do("GET", db.CreateKey("Chassis", containedByOid))
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

	if validationResult := c.createValidator(requestedChassis).Validate(); validationResult.HasErrors() {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(validationResult.Error())
		return
	}

	conn := c.cm.GetConnection()
	defer db.NewConnectionCloser(&conn)

	queryParams, err := prepareChassisCreationScriptParams(requestedChassis.IntializeIds())
	_, err = _CREATE_CHASSIS_SCRIPT.Do(conn, queryParams...)

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

func prepareChassisCreationScriptParams(rc *redfish.Chassis) (res []interface{}, err error) {
	res = append(res, rc.Oid)

	chassisBody, err := json.Marshal(rc)
	if err != nil {
		return
	}
	res = append(res, chassisBody)

	if len(rc.Links.ContainedBy) > 0 {
		res = append(res, rc.Links.ContainedBy[0].Oid)
	} else {
		res = append(res, nil)
	}
	return
}

// args:
//	KEYS[1]: chassis id
//  KEYS[2]: chassis body as []byte
//  KEYS[3]: parent chassis id
var _CREATE_CHASSIS_SCRIPT = redis.NewScript(3, `
		local res, err

		local chassisKey = "Chassis:"..KEYS[1]
		res, err = redis.pcall('setnx', chassisKey, KEYS[2])
		if err ~= nil then
			return  redis.error_reply(err)
		end
		if res == 0 then		
			return redis.error_reply("already exists")
		end

		if KEYS[3] == nil or KEYS[3] == '' then 
			return redis.status_reply("OK")
		end
		
		local containsKey = "CONTAINS:Chassis:"..KEYS[3]
		res, err = redis.pcall("sadd", containsKey, KEYS[1])
		if err ~= nil then
			return redis.error_reply(err)
		end

		local containedinKey = "CONTAINEDIN:Chassis:"..KEYS[1]		
		res, err = redis.pcall("set", containedinKey, KEYS[3])
		if err ~= nil then
			return redis.error_reply(err)
		end
		
		return redis.status_reply("OK")
	`)
