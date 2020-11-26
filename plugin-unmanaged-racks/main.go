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

package main

import (
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/logging"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/rest"
)

var version = "dev"

func main() {
	logging.Infof("Starting URP v%s\n", version)
	if pc, err := config.ReadPluginConfiguration(); err != nil {
		logging.Fatal("error while reading from config", err)
	} else {
		plugin := Plugin{
			connectionManager: db.NewConnectionManager(pc.RedisAddress, pc.SentinelMasterName),
			pluginConfig:      pc,
		}
		plugin.Run()
	}
}

type Plugin struct {
	connectionManager *db.ConnectionManager
	pluginConfig      *config.PluginConfig
}

func (p *Plugin) Run() {
	rest.InitializeAndRun(p.pluginConfig, p.connectionManager)
}
