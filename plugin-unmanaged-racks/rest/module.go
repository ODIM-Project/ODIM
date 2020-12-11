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
	"crypto/tls"
	"net/http"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/logging"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/host"
	"github.com/kataras/iris/v12/middleware/logger"
)

const urpPluginName = "URP"

// InitializeAndRun Initializes and runs Unamanged Racks Plugin
func InitializeAndRun(pluginConfiguration *config.PluginConfig) {

	odimraHTTPClient := redfish.NewHTTPClient(
		redfish.BaseURL(pluginConfiguration.OdimNBUrl),
		redfish.HTTPTransport(pluginConfiguration),
	)

	dao := db.CreateDAO(pluginConfiguration.RedisAddress, pluginConfiguration.SentinelMasterName)

	createApplication(pluginConfiguration, dao, odimraHTTPClient).Run(
		func(app *iris.Application) error {
			return app.NewHost(&http.Server{Addr: pluginConfiguration.Host + ":" + pluginConfiguration.Port}).
				Configure(configureTLS(pluginConfiguration)).
				ListenAndServe()
		},
	)
}

func configureTLS(c *config.PluginConfig) host.Configurator {
	return func(su *host.Supervisor) {
		cert, err := tls.LoadX509KeyPair(c.PKICertificatePath, c.PKIPrivateKeyPath)
		if err != nil {
			panic(err)
		}
		su.Server.TLSConfig = &tls.Config{
			Certificates:             []tls.Certificate{cert},
			MinVersion:               c.TLSConf.MinVersion,
			MaxVersion:               c.TLSConf.MaxVersion,
			CipherSuites:             c.TLSConf.PreferredCipherSuites,
			PreferServerCipherSuites: true,
		}
	}
}

func createApplication(pluginConfig *config.PluginConfig, dao *db.DAO, odimraHTTPClient *redfish.HTTPClient) *iris.Application {
	return iris.New().Configure(createLoggersConfigurer(pluginConfig), createAPIHandlersConfigurer(odimraHTTPClient, dao, pluginConfig))
}

func createLoggersConfigurer(pluginConfig *config.PluginConfig) iris.Configurator {
	return func(app *iris.Application) {
		//no startup log
		app.Configure(iris.WithoutStartupLog)
		//iris app attached to application logger
		app.Logger().Install(logging.GetLogger())
		//iris app log level adjusted to configured one
		app.Logger().SetLevel(pluginConfig.LogLevel)
		//app log level adjusted to configured one
		logging.SetLogLevel(pluginConfig.LogLevel)
		//enable request logger
		app.Use(logger.New())
	}
}

func createAPIHandlersConfigurer(odimraHTTPClient *redfish.HTTPClient, dao *db.DAO, pluginConfig *config.PluginConfig) iris.Configurator {
	return func(application *iris.Application) {

		basicAuthHandler := newBasicAuthHandler(pluginConfig.UserName, pluginConfig.Password)

		application.Post("/EventService/Events", newEventHandler(dao))

		pluginRoutes := application.Party("/ODIM/v1")
		pluginRoutes.Post("/Startup", basicAuthHandler, newStartupHandler(pluginConfig, odimraHTTPClient))
		pluginRoutes.Get("/Status", newPluginStatusController(pluginConfig))

		managers := pluginRoutes.Party("/Managers", basicAuthHandler)
		managers.Get("", newGetManagersHandler(pluginConfig))
		managers.Get("/{id}", newGetManagerHandler(pluginConfig))

		chassis := pluginRoutes.Party("/Chassis", basicAuthHandler)
		chassis.Get("", newGetChassisCollectionHandler(dao))
		chassis.Get("/{id}", newGetChassisHandler(dao))
		chassis.Delete("/{id}", newDeleteChassisHandler(dao))
		chassis.Post("", newPostChassisHandler(dao, pluginConfig))
		chassis.Patch("/{id}", newPatchChassisHandler(dao, odimraHTTPClient))
	}
}
