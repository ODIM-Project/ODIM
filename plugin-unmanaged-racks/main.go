//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.
package main

import (
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/eventing"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/rest"
	"github.com/kataras/iris/v12"
	"log"
)

var version = "dev"

func main() {
	log.Printf("Running URP v%s\n", version)
	var pc *config.PluginConfig
	var err error
	if pc, err = config.ReadPluginConfiguration(); err != nil {
		log.Fatalln("error while reading from config", err)
	}

	plugin := Plugin{
		connectionManager: db.NewConnectionManager("tcp", "odimra.local.com", "6380"),
		pluginConfig:      pc,
	}
	plugin.Run()
}

type Plugin struct {
	connectionManager *db.ConnectionManager
	pluginConfig      *config.PluginConfig
}

func (p *Plugin) Run() {
	eventing.NewSubscriber(p.pluginConfig).Run()
	go func() {
		eventing.NewListener(*p.pluginConfig, p.connectionManager).Run()
	}()

	basicAuthHandler := rest.NewBasicAuthHandler(p.pluginConfig.UserName, p.pluginConfig.Password)

	application := iris.New()
	pluginRoutes := application.Party("/ODIM/v1")
	pluginRoutes.Get("/Status", rest.NewPluginStatusController(p.pluginConfig))

	managers := pluginRoutes.Party("/Managers", basicAuthHandler)
	managers.Get("", rest.NewGetManagersHandler(p.pluginConfig))
	managers.Get("/{id}", rest.NewGetPluginManagerHandler(p.pluginConfig))

	chassis := pluginRoutes.Party("/Chassis", basicAuthHandler)
	chassis.Get("", rest.NewGetChassisCollectionHandler(p.connectionManager))
	chassis.Get("/{id}", rest.NewChassisReadingHandler(p.connectionManager))
	chassis.Delete("/{id}", rest.NewChassisDeletionHandler(p.connectionManager))
	chassis.Post("", rest.NewCreateChassisHandlerHandler(p.connectionManager, p.pluginConfig))
	chassis.Patch("/{id}", rest.NewChassisUpdateHandler(p.connectionManager))

	application.Run(
		iris.TLS(
			p.pluginConfig.Host+":"+p.pluginConfig.Port,
			p.pluginConfig.KeyCertConf.CertificatePath,
			p.pluginConfig.KeyCertConf.PrivateKeyPath,
		),
	)
}
