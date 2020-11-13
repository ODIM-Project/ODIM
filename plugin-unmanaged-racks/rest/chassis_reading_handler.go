package rest

import (
	"encoding/json"
	"net/http"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris/v12/context"
)

func NewChassisReadingHandler(cm *db.ConnectionManager) context.Handler {
	return (&chassisReadingHandler{cm}).handle
}

type chassisReadingHandler struct {
	cm *db.ConnectionManager
}

func (c *chassisReadingHandler) handle(ctx context.Context) {
	conn := c.cm.GetConnection()
	defer db.NewConnectionCloser(&conn)

	requestedChassisOid := ctx.Request().RequestURI
	v, err := redis.String(conn.Do("GET", db.CreateKey("Chassis", requestedChassisOid)))
	if err != nil && err == redis.ErrNil {
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(redfish.NewError().AddExtendedInfo(redfish.NewResourceNotFoundMsg("Chassis", requestedChassisOid, "")))
		return
	}
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
		return
	}

	chassis := new(redfish.Chassis)
	err = json.Unmarshal([]byte(v), chassis)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
		return
	}

	chassisContainsKey := db.CreateContainsKey("Chassis", requestedChassisOid)
	members, err := redis.Strings(conn.Do("SMEMBERS", chassisContainsKey))
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
		return
	}

	for _, m := range members {
		chassis.Links.Contains = append(chassis.Links.Contains, redfish.Link{Oid: m})
	}

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(chassis)
}
