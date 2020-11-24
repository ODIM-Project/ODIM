package rest

import (
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/logging"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
)

const _PLUGIN_NAME = "URP"

func InitializeAndRun(c *config.PluginConfig, cm *db.ConnectionManager) {
	application := createApplication(c, cm)
	_ = application.Run(
		iris.TLS(
			c.Host+":"+c.Port,
			c.PKICertificatePath,
			c.PKIPrivateKeyPath,
		),
	)
}

func createApplication(c *config.PluginConfig, cm *db.ConnectionManager) *iris.Application {
	application := iris.New()

	application.Logger().Install(logging.Logger())
	application.Logger().SetLevel(c.LogLevel)
	logging.Logger().SetLevel(c.LogLevel)

	//enable request logger
	application.Use(logger.New())

	basicAuthHandler := NewBasicAuthHandler(c.UserName, c.Password)

	pluginRoutes := application.Party("/ODIM/v1")
	pluginRoutes.Post("/EventService/Events", newEventHandler(cm, c.URLTranslation))
	pluginRoutes.Post("/Startup", basicAuthHandler, newStartupHandler(c))
	pluginRoutes.Get("/Status", newPluginStatusController(c))

	managers := pluginRoutes.Party("/Managers", basicAuthHandler)
	managers.Get("", NewGetManagersHandler(c))
	managers.Get("/{id}", NewGetPluginManagerHandler(c))

	chassis := pluginRoutes.Party("/Chassis", basicAuthHandler)
	chassis.Get("", newGetChassisCollectionHandler(cm))
	chassis.Get("/{id}", NewChassisReadingHandler(cm))
	chassis.Delete("/{id}", NewChassisDeletionHandler(cm))
	chassis.Post("", NewCreateChassisHandlerHandler(cm, c))
	chassis.Patch("/{id}", NewChassisUpdateHandler(cm, c))

	return application
}
