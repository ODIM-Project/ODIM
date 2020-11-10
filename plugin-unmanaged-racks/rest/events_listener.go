package rest

import (
	"encoding/json"
	"net/http"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/logging"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func newEventHandler(cm *db.ConnectionManager, translator *config.URLTranslation) context.Handler {
	return (&eventHandler{
		cm: cm,
		translator: &redfish.Translator{
			Dictionaries: translator,
		},
	}).handleEvent
}

type eventHandler struct {
	cm         *db.ConnectionManager
	translator *redfish.Translator
}

func (eh *eventHandler) handleEvent(c iris.Context) {
	raw := new(json.RawMessage)
	err := c.ReadJSON(raw)
	if err != nil {
		c.StatusCode(http.StatusBadRequest)
		_, _ = c.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
		return
	}

	message := new(redfish.MessageData)
	err = json.Unmarshal([]byte(eh.translator.RedfishToODIM(string(*raw))), message)
	if err != nil {
		c.StatusCode(http.StatusBadRequest)
		_, _ = c.JSON(redfish.CreateError(redfish.GeneralError, err.Error()))
		return
	}

	for _, e := range message.Events {
		containedInKey := eh.cm.CreateContainedInKey(e.OriginOfCondition.Oid)
		rackKey, err := redis.String(eh.cm.FindByKey(containedInKey))
		if err == redis.ErrNil {
			continue
		}
		if err != nil {
			continue
		}

		err = eh.cm.DoInTransaction(rackKey, func(c redis.Conn) {
			_ = c.Send("DEL", containedInKey)
			_ = c.Send("SREM", eh.cm.CreateChassisContainsKey(rackKey), e.OriginOfCondition.Oid)
		})

		if err != nil {
			logging.Error("cannot consume message(%v): %v", message, err)
		}
	}
}
