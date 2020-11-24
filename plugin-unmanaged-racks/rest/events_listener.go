package rest

import (
	stdCtx "context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/logging"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"
	"github.com/go-redis/redis/v8"
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
		ctx := stdCtx.TODO()
		containedInKey := db.CreateContainedInKey("Chassis", e.OriginOfCondition.Oid)
		rackKey, err := eh.cm.DAO().Get(ctx, containedInKey.String()).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			continue
		}

		err = eh.cm.DAO().Watch(ctx, func(tx *redis.Tx) error {
			_, err := tx.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
				if _, err := pipeliner.Del(
					ctx,
					containedInKey.String(),
				).Result(); err != nil {
					return fmt.Errorf("del: %s error: %w", containedInKey, err)
				}

				if _, err := pipeliner.SRem(
					ctx,
					db.CreateContainsKey("Chassis", rackKey).String(), e.OriginOfCondition.Oid,
				).Result(); err != nil {
					return fmt.Errorf("srem: %s error: %w", db.CreateContainsKey("Chassis", rackKey).String(), err)
				}
				return nil
			})
			return err
		}, rackKey)

		if err != nil {
			logging.Errorf(
				"cannot consume message(%v): %v",
				message,
				fmt.Errorf("couldn't commit transaction: %w", err),
			)
		}
	}
}
