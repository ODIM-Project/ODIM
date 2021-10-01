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
	"github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dphandler"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmessagebus"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmiddleware"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dputilities"
	iris "github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
)

var subscriptionInfo []dpmodel.Device
var log = logrus.New()

// TokenObject will contains the generated token and public key of odimra
type TokenObject struct {
	AuthToken string `json:"authToken"`
	PublicKey []byte `json:"publicKey"`
}

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
	dphandler.In, dphandler.Out = common.CreateJobQueue(jobQueueSize)

	// RunReadWorkers will create a worker pool for doing a specific task
	// which is passed to it as Publish method after reading the data from the channel.
	go common.RunReadWorkers(dphandler.Out, dpmessagebus.Publish, 5)

	configFilePath := os.Getenv("PLUGIN_CONFIG_FILE_PATH")
	if configFilePath == "" {
		log.Fatal("No value get the environment variable PLUGIN_CONFIG_FILE_PATH")
	}
	// TrackConfigFileChanges monitors the dell config changes using fsnotfiy
	go dputilities.TrackConfigFileChanges(configFilePath)

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
		log.Fatal("While initializing plugin server: " + err.Error())
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
		pluginRoutes.Post("/validate", dpmiddleware.BasicAuth, dphandler.Validate)
		pluginRoutes.Post("/Sessions", dphandler.CreateSession)
		pluginRoutes.Post("/Subscriptions", dpmiddleware.BasicAuth, dphandler.CreateEventSubscription)
		pluginRoutes.Delete("/Subscriptions", dpmiddleware.BasicAuth, dphandler.DeleteEventSubscription)

		//Adding routes related to all system gets
		systems := pluginRoutes.Party("/Systems", dpmiddleware.BasicAuth)
		systems.Get("", dphandler.GetResource)
		systems.Get("/{id}", dphandler.GetResource)
		systems.Get("/{id}/Storage", dphandler.GetResource)
		systems.Get("/{id}/Storage/{id2}", dphandler.GetResource)
		systems.Get("/{id}/Storage/{id2}/Volumes", dphandler.GetResource)
		systems.Post("/{id}/Storage/{id2}/Volumes", dphandler.CreateVolume)
		systems.Get("/{id}/Storage/{id2}/Volumes/{rid}", dphandler.GetResource)
		systems.Delete("/{id}/Storage/{id2}/Volumes/{rid}", dphandler.DeleteVolume)
		systems.Get("/{id}/Storage/{id2}/Drives/{rid}", dphandler.GetResource)
		systems.Get("/{id}/BootOptions", dphandler.GetResource)
		systems.Get("/{id}/BootOptions/{rid}", dphandler.GetResource)
		systems.Get("/{id}/Processors", dphandler.GetResource)
		systems.Get("/{id}/Processors/{rid}", dphandler.GetResource)
		systems.Get("/{id}/LogServices", dphandler.GetResource)
		systems.Get("/{id}/LogServices/{rid}", dphandler.GetResource)
		systems.Get("/{id}/LogServices/{rid}/Entries", dphandler.GetResource)
		systems.Get("/{id}/LogServices/{rid}/Entries/{rid2}", dphandler.GetResource)
		systems.Post("/{id}/LogServices/{rid}/Actions/LogService.ClearLog", dphandler.GetResource)
		systems.Get("/{id}/Memory", dphandler.GetResource)
		systems.Get("/{id}/Memory/{rid}", dphandler.GetResource)
		systems.Get("/{id}/NetworkInterfaces", dphandler.GetResource)
		systems.Get("/{id}/MemoryDomains", dphandler.GetResource)
		systems.Get("/{id}/EthernetInterfaces", dphandler.GetResource)
		systems.Get("/{id}/EthernetInterfaces/{rid}", dphandler.GetResource)
		systems.Get("/{id}/SecureBoot", dphandler.GetResource)
		systems.Get("/{id}/EthernetInterfaces/{id2}/VLANS", dphandler.GetResource)
		systems.Get("/{id}/EthernetInterfaces/{id2}/VLANS/{rid}", dphandler.GetResource)
		systems.Get("/{id}/NetworkInterfaces/{rid}", dphandler.GetResource)
		systems.Get("/{id}/PCIeDevices/{rid}", dphandler.GetResource)
		systems.Patch("/{id}", dphandler.ChangeSettings)

		systemsAction := systems.Party("/{id}/Actions")
		systemsAction.Post("/ComputerSystem.Reset", dphandler.ResetComputerSystem)
		systemsAction.Post("/ComputerSystem.SetDefaultBootOrder", dphandler.SetDefaultBootOrder)

		biosParty := systems.Party("/{id}/Bios")
		biosParty.Get("/", dphandler.GetResource)
		biosParty.Get("/Settings", dphandler.GetResource)
		biosParty.Patch("/Settings", dphandler.ChangeSettings)

		chassis := pluginRoutes.Party("/Chassis")
		chassis.Get("", dphandler.GetResource)
		chassis.Get("/{id}", dphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters", dphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters/{rid}", dphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters/{rid}/NetworkDeviceFunctions", dphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters/{rid}/NetworkPorts", dphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters/{rid}/NetworkDeviceFunctions/{id2}", dphandler.GetResource)
		chassis.Get("/{id}/NetworkAdapters/{rid}/NetworkPorts/{id2}", dphandler.GetResource)
		chassis.Get("/{id}/Assembly", dphandler.GetResource)
		chassis.Get("/{id}/PCIeSlots", dphandler.GetResource)
		chassis.Get("/{id}/PCIeSlots/{rid}", dphandler.GetResource)
		chassis.Get("/{id}/PCIeDevices", dphandler.GetResource)
		chassis.Get("/{id}/PCIeDevices/{rid}", dphandler.GetResource)
		chassis.Get("/{id}/Sensors", dphandler.GetResource)
		chassis.Get("/{id}/Sensors/{rid}", dphandler.GetResource)
		chassis.Get("/{id}/LogServices", dphandler.GetResource)
		chassis.Get("/{id}/LogServices/{rid}", dphandler.GetResource)
		chassis.Get("/{id}/LogServices/{rid}/Entries", dphandler.GetResource)
		chassis.Get("/{id}/LogServices/{rid}/Entries/{rid2}", dphandler.GetResource)

		// Chassis Power URl routes
		chassisPower := chassis.Party("/{id}/Power")
		chassisPower.Get("/", dphandler.GetResource)
		chassisPower.Get("#PowerControl/{id1}", dphandler.GetResource)
		chassisPower.Get("#PowerSupplies/{id1}", dphandler.GetResource)
		chassisPower.Get("#Redundancy/{id1}", dphandler.GetResource)

		// Chassis Thermal Url Routes
		chassisThermal := chassis.Party("/{id}/Thermal")
		chassisThermal.Get("/", dphandler.GetResource)
		chassisThermal.Get("#Fans/{id1}", dphandler.GetResource)
		chassisThermal.Get("#Temperatures/{id1}", dphandler.GetResource)

		// Manager routers
		managers := pluginRoutes.Party("/Managers", dpmiddleware.BasicAuth)
		managers.Get("", dphandler.GetManagersCollection)
		managers.Get("/{id}", dphandler.GetManagersInfo)
		managers.Get("/{id}/EthernetInterfaces", dphandler.GetResource)
		managers.Get("/{id}/EthernetInterfaces/{rid}", dphandler.GetResource)
		managers.Get("/{id}/NetworkProtocol", dphandler.GetResource)
		managers.Get("/{id}/NetworkProtocol/{rid}", dphandler.GetResource)
		managers.Get("/{id}/HostInterfaces", dphandler.GetResource)
		managers.Get("/{id}/HostInterfaces/{rid}", dphandler.GetResource)
		managers.Get("/{id}/VirtualMedia", dphandler.GetResource)
		managers.Get("/{id}/VirtualMedia/{rid}", dphandler.GetResource)
		managers.Post("/{id}/VirtualMedia/{rid}/Actions/VirtualMedia.EjectMedia", dphandler.VirtualMediaActions)
		managers.Post("/{id}/VirtualMedia/{rid}/Actions/VirtualMedia.InsertMedia", dphandler.VirtualMediaActions)
		managers.Get("/{id}/LogServices", dphandler.GetResource)
		managers.Get("/{id}/LogServices/{rid}", dphandler.GetResource)
		managers.Get("/{id}/LogServices/{rid}/Entries", dphandler.GetResource)
		managers.Get("/{id}/LogServices/{rid}/Entries/{rid2}", dphandler.GetResource)
		managers.Post("/{id}/LogServices/{rid}/Actions/LogService.ClearLog", dphandler.GetResource)

		//Registries routers
		registries := pluginRoutes.Party("/Registries", dpmiddleware.BasicAuth)
		registries.Get("", dphandler.GetResource)
		registries.Get("/{id}", dphandler.GetResource)

		registryStore := pluginRoutes.Party("/registrystore", dpmiddleware.BasicAuth)
		registryStore.Get("/registries/en/{id}", dphandler.GetResource)

		registryStoreCap := pluginRoutes.Party("/RegistryStore", dpmiddleware.BasicAuth)
		registryStoreCap.Get("/registries/en/{id}", dphandler.GetResource)

		// Routes related to Update service
		update := pluginRoutes.Party("/UpdateService", dpmiddleware.BasicAuth)
		update.Post("/Actions/UpdateService.SimpleUpdate", dphandler.SimpleUpdate)
		update.Post("/Actions/UpdateService.StartUpdate", dphandler.StartUpdate)
		update.Get("/FirmwareInventory", dphandler.GetResource)
		update.Get("/FirmwareInventory/{id}", dphandler.GetResource)
		update.Get("/SoftwareInventory", dphandler.GetResource)
		update.Get("/SoftwareInventory/{id}", dphandler.GetResource)
	}
	pluginRoutes.Get("/Status", dphandler.GetPluginStatus)
	pluginRoutes.Post("/Startup", dpmiddleware.BasicAuth, dphandler.GetPluginStartup)
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
		log.Fatalf("fatal: error while initializing event server: %v", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc(config.Data.EventConf.DestURI, dphandler.RedfishEvents)
	evtServer.Handler = mux
	log.Fatal(evtServer.ListenAndServeTLS("", ""))
}

// intializePluginStatus sets plugin status
func intializePluginStatus() {
	dputilities.Status.Available = "yes"
	dputilities.Status.Uptime = time.Now().Format(time.RFC3339)
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
	go common.RunWriteWorkers(dphandler.In, events, 1, done)
	log.Info("successfully sent startup event")
}
