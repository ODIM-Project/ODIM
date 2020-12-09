/*
 * Copyright (c) 2020 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rest

import (
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/logging"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

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
	odimraHttpClient := redfish.NewHttpClient(
		redfish.BaseURL(c.OdimNBUrl),
		redfish.HttpTransport(c),
	)

	application := iris.New()

	application.Logger().Install(logging.Logger())
	application.Logger().SetLevel(c.LogLevel)
	logging.Logger().SetLevel(c.LogLevel)

	//enable request logger
	application.Use(logger.New())
	application.Post("/EventService/Events", newEventHandler(cm, c.URLTranslation))

	basicAuthHandler := NewBasicAuthHandler(c.UserName, c.Password)

	pluginRoutes := application.Party("/ODIM/v1")
	pluginRoutes.Post("/Startup", basicAuthHandler, newStartupHandler(c, odimraHttpClient))
	pluginRoutes.Get("/Status", newPluginStatusController(c))

	managers := pluginRoutes.Party("/Managers", basicAuthHandler)
	managers.Get("", NewGetManagersHandler(c))
	managers.Get("/{id}", NewGetPluginManagerHandler(c))

	chassis := pluginRoutes.Party("/Chassis", basicAuthHandler)
	chassis.Get("", newGetChassisCollectionHandler(cm))
	chassis.Get("/{id}", NewChassisReadingHandler(cm))
	chassis.Delete("/{id}", NewChassisDeletionHandler(cm))
	chassis.Post("", NewCreateChassisHandlerHandler(cm, c))
	chassis.Patch("/{id}", NewChassisUpdateHandler(cm, odimraHttpClient))

	return application
}
