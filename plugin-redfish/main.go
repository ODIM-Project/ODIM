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
	"encoding/json"
	"net"
	"net/http"
	"os"
	"time"

	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	lutilconf "github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/plugin-redfish/config"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfphandler"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpmessagebus"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpmiddleware"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpmodel"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfputilities"
	iris "github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
)

var subscriptionInfo []rfpmodel.Device

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
	rfphandler.In, rfphandler.Out = common.CreateJobQueue(jobQueueSize)

	// RunReadWorkers will create a worker pool for doing a specific task
	// which is passed to it as Publish method after reading the data from the channel.
	go common.RunReadWorkers(rfphandler.Out, rfpmessagebus.Publish, 5)

	configFilePath := os.Getenv("PLUGIN_CONFIG_FILE_PATH")
	if configFilePath == "" {
		log.Fatal("No value get the environment variable PLUGIN_CONFIG_FILE_PATH")
	}
	// TrackConfigFileChanges monitors the odim config changes using fsnotfiy
	go rfputilities.TrackConfigFileChanges(configFilePath)

	intializePluginStatus()
	app()
}

func app() {
	go sendStartupEvent()
	go eventsrouters()

	app := routers()
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
	e := rfphandler.ExternalInterface{
		TokenValidation: rfphandler.TokenValidation,
		GetDeviceData:   rfphandler.GetDeviceData,
	}
	pluginRoutes := app.Party("/ODIM/v1")
	{
		pluginRoutes.Post("/validate", rfpmiddleware.BasicAuth, rfphandler.Validate)
		pluginRoutes.Post("/Sessions", rfphandler.CreateSession)
		pluginRoutes.Post("/Subscriptions", rfpmiddleware.BasicAuth, rfphandler.CreateEventSubscription)
		pluginRoutes.Delete("/Subscriptions", rfpmiddleware.BasicAuth, rfphandler.DeleteEventSubscription)

		//Adding routes related to all system gets
		systems := pluginRoutes.Party("/Systems", rfpmiddleware.BasicAuth)
		systems.Get("", rfphandler.GetResource)
		systems.Get("/{id}", rfphandler.GetResource)
		systems.Get("/{id}/Storage", rfphandler.GetResource)
		systems.Get("/{id}/Storage/{rid}", rfphandler.GetResource)
		systems.Get("/{id}/Storage/{rid}/Volumes", rfphandler.GetResource)
		systems.Post("/{id}/Storage/{rid}/Volumes", rfphandler.CreateVolume)
		systems.Get("/{id}/Storage/{rid}/Volumes/{rid}", rfphandler.GetResource)
		systems.Delete("/{id}/Storage/{id2}/Volumes/{rid}", rfphandler.DeleteVolume)
		systems.Get("/{id}/Storage/{id2}/Drives/{rid}", rfphandler.GetResource)
		systems.Get("/{id}/BootOptions", rfphandler.GetResource)
		systems.Get("/{id}/BootOptions/{rid}", rfphandler.GetResource)
		systems.Get("/{id}/Processors", rfphandler.GetResource)
		systems.Get("/{id}/Processors/{rid}", rfphandler.GetResource)
		systems.Get("/{id}/LogServices", rfphandler.GetResource)
		systems.Get("/{id}/LogServices/{rid}", rfphandler.GetResource)
		systems.Get("/{id}/LogServices/{rid}/Entries", rfphandler.GetResource)
		systems.Get("/{id}/LogServices/{rid}/Entries/{rid2}", rfphandler.GetResource)
		systems.Post("/{id}/LogServices/{rid}/Actions/LogService.ClearLog", rfphandler.GetResource)
		systems.Get("/{id}/Memory", rfphandler.GetResource)
		systems.Get("/{id}/Memory/{rid}", rfphandler.GetResource)
		systems.Get("/{id}/NetworkInterfaces", rfphandler.GetResource)
		systems.Get("/{id}/MemoryDomains", rfphandler.GetResource)
		systems.Get("/{id}/EthernetInterfaces", rfphandler.GetResource)
		systems.Get("/{id}/EthernetInterfaces/{rid}", rfphandler.GetResource)
		systems.Get("/{id}/SecureBoot", rfphandler.GetResource)
		systems.Get("/{id}/EthernetInterfaces/{id2}/VLANS", rfphandler.GetResource)
		systems.Get("/{id}/EthernetInterfaces/{id2}/VLANS/{rid}", rfphandler.GetResource)
		systems.Get("/{id}/NetworkInterfaces/{rid}", rfphandler.GetResource)
		systems.Get("/{id}/PCIeDevices/{rid}", rfphandler.GetResource)
		systems.Patch("/{id}", rfphandler.ChangeSettings)

		systemsAction := systems.Party("/{id}/Actions")
		systemsAction.Post("/ComputerSystem.Reset", rfphandler.ResetComputerSystem)
		systemsAction.Post("/ComputerSystem.SetDefaultBootOrder", rfphandler.SetDefaultBootOrder)

		biosParty := systems.Party("/{id}/Bios")
		biosParty.Get("/", rfphandler.GetResource)
		biosParty.Get("/Settings", rfphandler.GetResource)
		biosParty.Patch("/Settings", rfphandler.ChangeSettings)

		chassis := pluginRoutes.Party("/Chassis")
		chassis.Get("", rfphandler.GetResource)
		chassis.Get("/{id}", rfphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters", rfphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters/{rid}", rfphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters/{id2}/NetworkDeviceFunctions", rfphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters/{id2}/NetworkPorts", rfphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters/{id2}/NetworkDeviceFunctions/{rid}", rfphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters/{id2}/NetworkPorts/{rid}", rfphandler.GetResource)
		chassis.Get("/{id}/Assembly", rfphandler.GetResource)
		chassis.Get("/{id}/PCIeSlots", rfphandler.GetResource)
		chassis.Get("/{id}/PCIeSlots/{rid}", rfphandler.GetResource)
		chassis.Get("/{id}/PCIeDevices", rfphandler.GetResource)
		chassis.Get("/{id}/PCIeDevices/{rid}", rfphandler.GetResource)
		chassis.Get("/{id}/Sensors", rfphandler.GetResource)
		chassis.Get("/{id}/Sensors/{rid}", rfphandler.GetResource)
		chassis.Get("/{id}/LogServices", rfphandler.GetResource)
		chassis.Get("/{id}/LogServices/{rid}", rfphandler.GetResource)
		chassis.Get("/{id}/LogServices/{rid}/Entries", rfphandler.GetResource)
		chassis.Get("/{id}/LogServices/{rid}/Entries/{rid2}", rfphandler.GetResource)
		// TODO:
		// chassis.Post("/{id}/LogServices/{rid}/Actions/LogService.ClearLog", rfphandler.GetResource)

		// Chassis Power URl routes
		chassisPower := chassis.Party("/{id}/Power")
		chassisPower.Get("/", rfphandler.GetResource)
		chassisPower.Get("#PowerControl/{id1}", rfphandler.GetResource)
		chassisPower.Get("#PowerSupplies/{id1}", rfphandler.GetResource)
		chassisPower.Get("#Redundancy/{id1}", rfphandler.GetResource)

		// Chassis Thermal Url Routes
		chassisThermal := chassis.Party("/{id}/Thermal")
		chassisThermal.Get("/", rfphandler.GetResource)
		chassisThermal.Get("#Fans/{id1}", rfphandler.GetResource)
		chassisThermal.Get("#Temperatures/{id1}", rfphandler.GetResource)

		// Manager routers
		managers := pluginRoutes.Party("/Managers", rfpmiddleware.BasicAuth)
		managers.Get("", rfphandler.GetManagersCollection)
		managers.Get("/{id}", rfphandler.GetManagersInfo)
		managers.Get("/{id}/EthernetInterfaces", rfphandler.GetResource)
		managers.Get("/{id}/EthernetInterfaces/{rid}", rfphandler.GetResource)
		managers.Get("/{id}/NetworkProtocol", rfphandler.GetResource)
		managers.Get("/{id}/NetworkProtocol/{rid}", rfphandler.GetResource)
		managers.Get("/{id}/HostInterfaces", rfphandler.GetResource)
		managers.Get("/{id}/HostInterfaces/{rid}", rfphandler.GetResource)
		managers.Get("/{id}/SerialInterface", rfphandler.GetResource)
		managers.Get("/{id}/SerialInterface/{rid}", rfphandler.GetResource)
		managers.Get("/{id}/VirtualMedia", rfphandler.GetResource)
		managers.Get("/{id}/VirtualMedia/{rid}", rfphandler.GetResource)
		managers.Post("/{id}/VirtualMedia/{rid}/Actions/VirtualMedia.EjectMedia", rfphandler.VirtualMediaActions)
		managers.Post("/{id}/VirtualMedia/{rid}/Actions/VirtualMedia.InsertMedia", rfphandler.VirtualMediaActions)
		managers.Get("/{id}/LogServices", rfphandler.GetResource)
		managers.Get("/{id}/LogServices/{rid}", rfphandler.GetResource)
		managers.Get("/{id}/LogServices/{rid}/Entries", rfphandler.GetResource)
		managers.Get("/{id}/LogServices/{rid}/Entries/{rid2}", rfphandler.GetResource)
		managers.Post("/{id}/LogServices/{rid}/Actions/LogService.ClearLog", rfphandler.GetResource)

		//Registries routers
		registries := pluginRoutes.Party("/Registries", rfpmiddleware.BasicAuth)
		registries.Get("", rfphandler.GetResource)
		registries.Get("/{id}", rfphandler.GetResource)

		registryStore := pluginRoutes.Party("/registrystore", rfpmiddleware.BasicAuth)
		registryStore.Get("/registries/en/{id}", rfphandler.GetResource)

		registryStoreCap := pluginRoutes.Party("/RegistryStore", rfpmiddleware.BasicAuth)
		registryStoreCap.Get("/registries/en/{id}", rfphandler.GetResource)

		// Routes related to Update service
		update := pluginRoutes.Party("/UpdateService", rfpmiddleware.BasicAuth)
		update.Post("/Actions/UpdateService.SimpleUpdate", rfphandler.SimpleUpdate)
		update.Post("/Actions/UpdateService.StartUpdate", rfphandler.StartUpdate)
		update.Get("/FirmwareInventory", rfphandler.GetResource)
		update.Get("/FirmwareInventory/{id}", rfphandler.GetResource)
		update.Get("/SoftwareInventory", rfphandler.GetResource)
		update.Get("/SoftwareInventory/{id}", rfphandler.GetResource)

		//Adding routes related to telemetry service
		telemetry := pluginRoutes.Party("/TelemetryService", rfpmiddleware.BasicAuth)
		telemetry.Get("/MetricDefinitions", rfphandler.GetResource)
		telemetry.Get("/MetricReportDefinitions", rfphandler.GetResource)
		telemetry.Get("/MetricReports", rfphandler.GetResource)
		telemetry.Get("/Triggers", rfphandler.GetResource)
		telemetry.Get("/MetricReports/{id}", e.GetMetricReport)
		telemetry.Get("/MetricDefinitions/{id}", rfphandler.GetResource)
		telemetry.Get("/MetricReportDefinitions/{id}", rfphandler.GetResource)
		telemetry.Get("/Triggers/{id}", rfphandler.GetResource)

	}
	pluginRoutes.Get("/Status", rfphandler.GetPluginStatus)
	pluginRoutes.Post("/Startup", rfpmiddleware.BasicAuth, rfphandler.GetPluginStartup)
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
	mux.HandleFunc(config.Data.EventConf.DestURI, rfphandler.RedfishEvents)
	evtServer.Handler = mux
	log.Fatal(evtServer.ListenAndServeTLS("", ""))
}

// intializePluginStatus sets plugin status
func intializePluginStatus() {
	rfputilities.Status.Available = "yes"
	rfputilities.Status.Uptime = time.Now().Format(time.RFC3339)
}

// sendStartupEvent is for sending startup event
func sendStartupEvent() {
	// grace wait time for plugin to be functional
	time.Sleep(3 * time.Second)

	var pluginIP string
	if pluginIP = os.Getenv("ASSIGNED_POD_IP"); pluginIP == "" {
		pluginIP = config.Data.PluginConf.Host
	}

	startupEvt := common.PluginStatusEvent{
		Name:         "Plugin startup event",
		Type:         "PluginStarted",
		Timestamp:    time.Now().String(),
		OriginatorID: pluginIP,
	}

	request, _ := json.Marshal(startupEvt)
	event := common.Events{
		IP:        net.JoinHostPort(config.Data.PluginConf.Host, config.Data.PluginConf.Port),
		Request:   request,
		EventType: "PluginStartUp",
	}

	done := make(chan bool)
	events := []interface{}{event}
	go common.RunWriteWorkers(rfphandler.In, events, 1, done)
	log.Info("successfully sent startup event")
}
