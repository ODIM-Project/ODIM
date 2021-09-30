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
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"

	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	lutilconf "github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/config"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lphandler"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lpmessagebus"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lpmiddleware"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lpmodel"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lputilities"
	iris "github.com/kataras/iris/v12"
)

var subscriptionInfo []lpmodel.Device

// TokenObject will contains the generated token and public key of odimra
type TokenObject struct {
	AuthToken string `json:"authToken"`
	PublicKey []byte `json:"publicKey"`
}

var log = logrus.New()

func main() {
	// verifying the uid of the user
	if uid := os.Geteuid(); uid == 0 {
		log.Fatal("Plugin Service should not be run as the root user")
	}

	if err := config.SetConfiguration(); err != nil {
		log.Fatal("While reading from config, got: " + err.Error())
	}

	if err := dc.SetConfiguration(config.Data.MessageBusConf.MessageQueueConfigFilePath); err != nil {
		log.Fatal("While trying to set messagebus configuration, got: " + err.Error())
	}

	// CreateJobQueue defines the queue which will act as an infinite buffer
	// In channel is an entry or input channel and the Out channel is an exit or output channel
	jobQueueSize := 10
	lphandler.In, lphandler.Out = common.CreateJobQueue(jobQueueSize)

	// RunReadWorkers will create a worker pool for doing a specific task
	// which is passed to it as Publish method after reading the data from the channel.
	go common.RunReadWorkers(lphandler.Out, lpmessagebus.Publish, 5)

	configFilePath := os.Getenv("PLUGIN_CONFIG_FILE_PATH")
	if configFilePath == "" {
		log.Fatal("No value get the environment variable PLUGIN_CONFIG_FILE_PATH")
	}
	// TrackConfigFileChanges monitors the odim config changes using fsnotfiy
	go lputilities.TrackConfigFileChanges(configFilePath)

	intializePluginStatus()
	app()
}

func app() {
	app := routers()
	go func() {
		eventsrouters()
	}()
	conf := &lutilconf.HTTPConfig{
		Certificate:   &config.Data.KeyCertConf.Certificate,
		PrivateKey:    &config.Data.KeyCertConf.PrivateKey,
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
		ServerAddress: config.Data.PluginConf.Host,
		ServerPort:    config.Data.PluginConf.Port,
	}
	pluginServer, err := conf.GetHTTPServerObj()
	if err != nil {
		log.Fatal("Unable to initialize plugin : " + err.Error())
	}
	app.Run(iris.Server(pluginServer))
}

func routers() *iris.Application {
	app := iris.New()
	app.WrapRouter(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		path := r.URL.Path
		if len(path) > 1 && path[len(path)-1] == '/' && path[len(path)-2] != '/' {
			path = path[:len(path)-1]
			r.RequestURI = path
			r.URL.Path = path
		}
		next(w, r)
	})

	pluginRoutes := app.Party("/ODIM/v1")
	{
		pluginRoutes.Post("/validate", lpmiddleware.BasicAuth, lphandler.Validate)
		pluginRoutes.Post("/Sessions", lphandler.CreateSession)
		pluginRoutes.Post("/Subscriptions", lpmiddleware.BasicAuth, lphandler.CreateEventSubscription)
		pluginRoutes.Delete("/Subscriptions", lpmiddleware.BasicAuth, lphandler.DeleteEventSubscription)

		//Adding routes related to all system gets
		systems := pluginRoutes.Party("/Systems", lpmiddleware.BasicAuth)
		systems.Get("", lphandler.GetResource)
		systems.Get("/{id}", lphandler.GetResource)
		systems.Get("/{id}/Storage", lphandler.GetResource)
		systems.Get("/{id}/Storage/{rid}", lphandler.GetResource)
		systems.Get("/{id}/Storage/{rid}/Volumes", lphandler.GetResource)
		systems.Get("/{id}/Storage/{rid}/Volumes/{rid}", lphandler.GetResource)
		systems.Post("/{id}/Storage/{rid}/Volumes", lphandler.MethodNotAllowed)
		systems.Delete("/{id}/Storage/{id2}/Volumes/{rid}", lphandler.MethodNotAllowed)
		systems.Get("/{id}/Storage/{id2}/Drives/{rid}", lphandler.GetResource)
		systems.Get("/{id}/BootOptions", lphandler.GetResource)
		systems.Get("/{id}/BootOptions/{rid}", lphandler.GetResource)
		systems.Get("/{id}/Processors", lphandler.GetResource)
		systems.Get("/{id}/Processors/{rid}", lphandler.GetResource)
		systems.Get("/{id}/LogServices", lphandler.GetResource)
		systems.Get("/{id}/LogServices/{rid}", lphandler.GetResource)
		systems.Get("/{id}/LogServices/{rid}/Entries", lphandler.GetResource)
		systems.Get("/{id}/LogServices/{rid}/Entries/{rid2}", lphandler.GetResource)
		systems.Post("/{id}/LogServices/{rid}/Actions/LogService.ClearLog", lphandler.GetResource)
		systems.Get("/{id}/Memory", lphandler.GetResource)
		systems.Get("/{id}/Memory/{rid}", lphandler.GetResource)
		systems.Get("/{id}/NetworkInterfaces", lphandler.GetResource)
		systems.Get("/{id}/MemoryDomains", lphandler.GetResource)
		systems.Get("/{id}/EthernetInterfaces", lphandler.GetResource)
		systems.Get("/{id}/EthernetInterfaces/{rid}", lphandler.GetResource)
		systems.Get("/{id}/SecureBoot", lphandler.GetResource)
		systems.Get("/{id}/EthernetInterfaces/{id2}/VLANS", lphandler.GetResource)
		systems.Get("/{id}/EthernetInterfaces/{id2}/VLANS/{rid}", lphandler.GetResource)
		systems.Get("/{id}/NetworkInterfaces/{rid}", lphandler.GetResource)
		systems.Get("/{id}/PCIeDevices/{rid}", lphandler.GetResource)
		systems.Patch("/{id}", lphandler.ChangeSettings)

		systemsAction := systems.Party("/{id}/Actions")
		systemsAction.Post("/ComputerSystem.Reset", lphandler.ResetComputerSystem)
		systemsAction.Post("/ComputerSystem.SetDefaultBootOrder", lphandler.SetDefaultBootOrder)

		biosParty := systems.Party("/{id}/Bios")
		biosParty.Get("/", lphandler.GetResource)
		biosParty.Get("/Settings", lphandler.GetResource)
		biosParty.Patch("/Settings", lphandler.ChangeSettings)

		chassis := pluginRoutes.Party("/Chassis")
		chassis.Get("", lphandler.GetResource)
		chassis.Get("/{id}", lphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters", lphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters/{rid}", lphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters/{id2}/NetworkDeviceFunctions", lphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters/{id2}/NetworkPorts", lphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters/{id2}/NetworkDeviceFunctions/{rid}", lphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters/{id2}/NetworkPorts/{rid}", lphandler.GetResource)
		chassis.Get("/{id}/Assembly", lphandler.GetResource)
		chassis.Get("/{id}/PCIeSlots", lphandler.GetResource)
		chassis.Get("/{id}/PCIeSlots/{rid}", lphandler.GetResource)
		chassis.Get("/{id}/PCIeDevices", lphandler.GetResource)
		chassis.Get("/{id}/PCIeDevices/{rid}", lphandler.GetResource)
		chassis.Get("/{id}/Sensors", lphandler.GetResource)
		chassis.Get("/{id}/Sensors/{rid}", lphandler.GetResource)
		chassis.Get("/{id}/LogServices", lphandler.GetResource)
		chassis.Get("/{id}/LogServices/{rid}", lphandler.GetResource)
		chassis.Get("/{id}/LogServices/{rid}/Entries", lphandler.GetResource)
		chassis.Get("/{id}/LogServices/{rid}/Entries/{rid2}", lphandler.GetResource)
		// TODO:
		// chassis.Post("/{id}/LogServices/{rid}/Actions/LogService.ClearLog", lphandler.GetResource)

		// Chassis Power URl routes
		chassisPower := chassis.Party("/{id}/Power")
		chassisPower.Get("/", lphandler.GetResource)
		chassisPower.Get("#PowerControl/{id1}", lphandler.GetResource)
		chassisPower.Get("#PowerSupplies/{id1}", lphandler.GetResource)
		chassisPower.Get("#Redundancy/{id1}", lphandler.GetResource)

		// Chassis Thermal Url Routes
		chassisThermal := chassis.Party("/{id}/Thermal")
		chassisThermal.Get("/", lphandler.GetResource)
		chassisThermal.Get("#Fans/{id1}", lphandler.GetResource)
		chassisThermal.Get("#Temperatures/{id1}", lphandler.GetResource)

		// Manager routers
		managers := pluginRoutes.Party("/Managers", lpmiddleware.BasicAuth)
		managers.Get("", lphandler.GetManagersCollection)
		managers.Get("/{id}", lphandler.GetManagersInfo)
		managers.Get("/{id}/EthernetInterfaces", lphandler.GetResource)
		managers.Get("/{id}/EthernetInterfaces/{rid}", lphandler.GetResource)
		managers.Get("/{id}/NetworkProtocol", lphandler.GetResource)
		managers.Get("/{id}/NetworkProtocol/{rid}", lphandler.GetResource)
		managers.Get("/{id}/HostInterfaces", lphandler.GetResource)
		managers.Get("/{id}/HostInterfaces/{rid}", lphandler.GetResource)
		managers.Get("/{id}/SerialInterface", lphandler.GetResource)
		managers.Get("/{id}/SerialInterface/{rid}", lphandler.GetResource)
		managers.Get("/{id}/VirtualMedia", lphandler.GetResource)
		managers.Get("/{id}/VirtualMedia/{rid}", lphandler.GetResource)
		managers.Post("/{id}/VirtualMedia/{rid}/Actions/VirtualMedia.EjectMedia", lphandler.VirtualMediaActions)
		managers.Post("/{id}/VirtualMedia/{rid}/Actions/VirtualMedia.InsertMedia", lphandler.VirtualMediaActions)
		managers.Get("/{id}/LogServices", lphandler.GetResource)
		managers.Get("/{id}/LogServices/{rid}", lphandler.GetResource)
		managers.Get("/{id}/LogServices/{rid}/Entries", lphandler.GetResource)
		managers.Get("/{id}/LogServices/{rid}/Entries/{rid2}", lphandler.GetResource)
		managers.Post("/{id}/LogServices/{rid}/Actions/LogService.ClearLog", lphandler.GetResource)

		//Registries routers
		registries := pluginRoutes.Party("/Registries", lpmiddleware.BasicAuth)
		registries.Get("", lphandler.GetResource)
		registries.Get("/{id}", lphandler.GetResource)

		registryStore := pluginRoutes.Party("/schemas", lpmiddleware.BasicAuth)
		registryStore.Get("/registries/{id}", lphandler.GetResource)

		// Routes related to Update service
		update := pluginRoutes.Party("/UpdateService", lpmiddleware.BasicAuth)
		update.Post("/Actions/UpdateService.SimpleUpdate", lphandler.SimpleUpdate)
		update.Post("/Actions/UpdateService.StartUpdate", lphandler.SimpleUpdate)
		update.Get("/FirmwareInventory", lphandler.GetResource)
		update.Get("/FirmwareInventory/{id}", lphandler.GetResource)
		update.Get("/SoftwareInventory", lphandler.GetResource)
		update.Get("/SoftwareInventory/{id}", lphandler.GetResource)

		//Adding routes related to telemetry service
		telemetry := pluginRoutes.Party("/TelemetryService", lpmiddleware.BasicAuth)
		telemetry.Get("/MetricDefinitions", lphandler.GetResource)
		telemetry.Get("/MetricReportDefinitions", lphandler.GetResource)
		telemetry.Get("/MetricReports", lphandler.GetResource)
		telemetry.Get("/Triggers", lphandler.GetResource)
	}
	pluginRoutes.Get("/Status", lphandler.GetPluginStatus)
	pluginRoutes.Post("/Startup", lpmiddleware.BasicAuth, lphandler.GetPluginStartup)
	return app
}

func eventsrouters() {
	conf := &lutilconf.HTTPConfig{
		Certificate:   &config.Data.KeyCertConf.Certificate,
		PrivateKey:    &config.Data.KeyCertConf.PrivateKey,
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
		ServerAddress: config.Data.EventConf.ListenerHost,
		ServerPort:    config.Data.EventConf.ListenerPort,
	}
	evtServer, err := conf.GetHTTPServerObj()
	if err != nil {
		log.Fatal("Unable to initialize event server: " + err.Error())
	}
	mux := http.NewServeMux()
	mux.HandleFunc(config.Data.EventConf.DestURI, lphandler.RedfishEvents)
	evtServer.Handler = mux
	log.Fatal(evtServer.ListenAndServeTLS("", ""))
}

// intializePluginStatus sets plugin status
func intializePluginStatus() {
	lputilities.Status.Available = "yes"
	lputilities.Status.Uptime = time.Now().Format(time.RFC3339)

}
