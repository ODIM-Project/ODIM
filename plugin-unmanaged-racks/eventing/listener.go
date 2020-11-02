package eventing

import (
	"encoding/json"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris/v12"
	"log"
	"net/http"
)

func NewListener(c config.PluginConfig, cm *db.ConnectionManager) Listener {
	return &listener{
		pluginConfig: c,
		cm:           cm,
	}
}

type Listener interface {
	Run()
}

type listener struct {
	pluginConfig config.PluginConfig
	cm           *db.ConnectionManager
}

func (l *listener) Run() {
	app := iris.New()

	eh := eventHandler{
		cm: l.cm,
		translator: &redfish.Translator{
			Dictionaries: l.pluginConfig.URLTranslation,
		},
	}
	app.Post(l.pluginConfig.EventConf.DestURI, eh.handleEvent)

	app.Run(
		iris.TLS(
			l.pluginConfig.EventConf.ListenerHost+":"+l.pluginConfig.EventConf.ListenerPort,
			l.pluginConfig.KeyCertConf.CertificatePath,
			l.pluginConfig.KeyCertConf.PrivateKeyPath,
		),
	)
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
			log.Printf("error: cannot consume message(%v): %v\n", message, err)
		}
	}
}
