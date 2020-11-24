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
