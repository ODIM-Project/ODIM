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
	"encoding/json"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/logging"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/utils"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/host"
	"github.com/kataras/iris/v12/middleware/logger"
)

const urpPluginName = "URP"

// InitializeAndRun Initializes and runs Unamanged Racks Plugin
func InitializeAndRun(pluginConfiguration *config.PluginConfig) {

	enigma := utils.NewEnigma(pluginConfiguration.RSAPrivateKeyPath, pluginConfiguration.RSAPublicKeyPath)

	odimraHTTPClient := redfish.NewHTTPClient(
		redfish.BaseURL(pluginConfiguration.OdimURL),
		redfish.HTTPTransport(pluginConfiguration),
		redfish.BasicAuth(pluginConfiguration.OdimUserName, enigma.Decrypt(pluginConfiguration.OdimPassword)),
	)

	dao := db.CreateDAO(pluginConfiguration.RedisAddress, pluginConfiguration.SentinelMasterName)

	createApplication(pluginConfiguration, dao, odimraHTTPClient).Run(
		func(app *iris.Application) error {
			supervisor := app.NewHost(&http.Server{Addr: pluginConfiguration.Host + ":" + pluginConfiguration.Port})
			return supervisor.Configure(configureTLS(pluginConfiguration)).ListenAndServe()
		},
	)
}

func configureTLS(c *config.PluginConfig) host.Configurator {
	return func(su *host.Supervisor) {
		cert, err := tls.LoadX509KeyPair(c.PKICertificatePath, c.PKIPrivateKeyPath)
		if err != nil {
			panic(err)
		}

		// 0xc02f #TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
		// 0xc030 #TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
		preferredCipherSuites := []uint16{0xc02f, 0xc030}

		su.Server.TLSConfig = &tls.Config{
			Certificates:             []tls.Certificate{cert},
			MinVersion:               c.TLSConf.MinVersion,
			MaxVersion:               c.TLSConf.MaxVersion,
			CipherSuites:             preferredCipherSuites,
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

		pluginActivator := createPluginActivator(pluginConfig, odimraHTTPClient)
		pluginRoutes.Get("/Status", pluginActivator, newPluginStatusController(pluginConfig))

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

func createPluginActivator(conf *config.PluginConfig, client *redfish.HTTPClient) func(c context.Context) {
	var once sync.Once
	return func(c context.Context) {
		once.Do(func() {
			go newSubscriber(conf, client).Run()
			logging.Debugf("URP plugin has been activated")
		})
		c.Next()
	}
}

func newSubscriber(config *config.PluginConfig, httpClient *redfish.HTTPClient) *subscriber {
	subscriptionTarget, err := url.Parse("https://" + config.Host + ":" + config.Port + "/EventService/Events")
	if err != nil {
		panic(err)
	}

	return &subscriber{
		destinationURL: *subscriptionTarget,
		odimRAClient:   httpClient,
	}
}

type subscriber struct {
	odimRAClient   *redfish.HTTPClient
	destinationURL url.URL
}

func (s *subscriber) Run() {
	logging.Info("Starting EventSubscriber")
	for {
		s.subscribe()
		time.Sleep(time.Second * 15)
	}
}

func (s *subscriber) subscribe() {
	sr := createSubscriptionRequest(s.destinationURL.String())
	bodyBytes, err := json.Marshal(&sr)
	if err != nil {
		logging.Errorf("Unexpected error during Subscription Request serialization: %s", err)
		return
	}

	rsp, err := s.odimRAClient.Post("/redfish/v1/EventService/Subscriptions", bodyBytes)
	if err != nil {
		logging.Errorf("Cannot register subscription: %s", err)
		return
	}
	if rsp.StatusCode != http.StatusAccepted {
		logging.Infof("Registration of subscription has been rejected with code(%s)", rsp.Status)
		return
	}

	monitor := func() (*http.Response, error) {
		return s.odimRAClient.Get(rsp.Header.Get("Location"))
	}

	for {
		r, e := monitor()
		if e != nil {
			logging.Errorf("Task monitoring interrupted by communication error: %s", e)
		}

		switch r.StatusCode {
		case http.StatusOK:
			logging.Infof("URP->ODIMRA event subscription registered successfully")
			return
		case http.StatusAccepted:
			continue
		case http.StatusConflict:
			logging.Info("URP->ODIMRA event subscription is already registered")
			return
		default:
			logging.Infof("Task monitor(%s) reports %s status code", rsp.Header.Get("Location"), r.Status)
			return
		}
	}
}

func createSubscriptionRequest(destination string) redfish.EvtSubPost {
	return redfish.EvtSubPost{
		Name:                 "URP",
		Destination:          destination,
		EventTypes:           []string{"ResourceRemoved"},
		MessageIds:           nil,
		ResourceTypes:        []string{"Chassis"},
		Context:              "ODIMRA_Event",
		Protocol:             "Redfish",
		SubscriptionType:     "RedfishEvent",
		EventFormatType:      "Event",
		SubordinateResources: true,
		OriginResources: []redfish.Link{
			{
				Oid: "/redfish/v1/Chassis",
			},
		},
	}
}
