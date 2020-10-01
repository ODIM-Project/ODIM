package rest

import (
	"encoding/json"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris/v12/context"
	"net/http"
)

func NewChassisReadingHandler(cm *db.ConnectionManager) context.Handler {
	return (&chassisReadingHandler{cm}).handle
}

type chassisReadingHandler struct {
	cm *db.ConnectionManager
}

func (c *chassisReadingHandler) handle(ctx context.Context) {
	v, err := redis.String(c.cm.FindByKey("Chassis", ctx.Request().RequestURI))
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		return
	}

	chassis := new(redfish.Chassis)
	err = json.Unmarshal([]byte(v), chassis)
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		return
	}

	ctx.JSON(chassis)
}
