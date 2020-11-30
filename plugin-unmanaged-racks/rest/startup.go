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
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/logging"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func newStartupHandler(c *config.PluginConfig) context.Handler {
	return (&startup{
		subscriber: newSubscriber(c),
	}).handle
}

type startup struct {
	subscriber *subscriber
}

func (s *startup) handle(c iris.Context) {
	if !s.subscriber.isRunning {
		go s.subscriber.Run()
	}
	c.StatusCode(http.StatusOK)
}

func newSubscriber(config *config.PluginConfig) *subscriber {

	subscriptionTarget, err := url.Parse("https://" + config.Host + ":" + config.Port + "/EventService/Events")
	if err != nil {
		panic(err)
	}

	return &subscriber{
		destinationURL: *subscriptionTarget,
		odimRAClient:   redfish.NewClient(config.OdimNBUrl),
	}
}

type subscriber struct {
	odimRAClient   redfish.Client
	destinationURL url.URL
	isRunning      bool
}

func (s *subscriber) Run() {
	logging.Info("Starting EventSubscriber")
	s.isRunning = true
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
