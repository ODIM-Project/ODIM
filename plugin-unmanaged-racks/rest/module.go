package rest

import (
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/kataras/iris/v12"
)

func InitializeAndRun(c *config.PluginConfig, cm *db.ConnectionManager) {
	application := iris.New()

	basicAuthHandler := NewBasicAuthHandler(c.UserName, c.Password)

	pluginRoutes := application.Party("/ODIM/v1")
	pluginRoutes.Post("/Startup", basicAuthHandler, newStartupHandler(c))
	pluginRoutes.Get("/Status", NewPluginStatusController(c))

	managers := pluginRoutes.Party("/Managers", basicAuthHandler)
	managers.Get("", NewGetManagersHandler(c))
	managers.Get("/{id}", NewGetPluginManagerHandler(c))

	chassis := pluginRoutes.Party("/Chassis", basicAuthHandler)
	chassis.Get("", NewGetChassisCollectionHandler(cm))
	chassis.Get("/{id}", NewChassisReadingHandler(cm))
	chassis.Delete("/{id}", NewChassisDeletionHandler(cm))
	chassis.Post("", NewCreateChassisHandlerHandler(cm, c))
	chassis.Patch("/{id}", NewChassisUpdateHandler(cm))

	application.Post(c.EventConf.DestURI, newEventHandler(cm, c.URLTranslation))

	application.Run(
		iris.TLS(
			c.Host+":"+c.Port,
			c.KeyCertConf.CertificatePath,
			c.KeyCertConf.PrivateKeyPath,
		),
	)
}
