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

// Package dpmessagebus ...
package dpmessagebus

import (
	"context"
	"encoding/json"
	"strings"

	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/plugin-dell/config"
	log "github.com/sirupsen/logrus"
)

// Publish ...
func Publish(ctx context.Context, data interface{}) bool {
	if data == nil {
		log.Error("Invalid data on publishing events")
		return false
	}
	event := data.(common.Events)

	topic := config.Data.MessageBusConf.EmbQueue[0]
	K, err := dc.Communicator(dc.KAFKA, config.Data.MessageBusConf.MessageBusConfigFilePath, topic)
	if err != nil {
		log.Error("Unable communicate with kafka: " + err.Error())
		return false
	}

	// Since we are deleting the first event from the event list,
	// processing the first event
	var message common.MessageData
	err = json.Unmarshal(event.Request, &message)
	if err != nil {
		log.Error("Failed to unmarshal the event, got: " + err.Error())
		return false
	}
	if message.OdataType != "" && !strings.Contains(message.OdataType, "MetricReport") {
		event.Request, err = formatEventRequest(ctx, message)
	}

	if err := K.Distribute(event); err != nil {
		log.Error("Unable Publish events to kafka: " + err.Error())
		return false
	}
	for _, eventMessage := range message.Events {
		log.Info(eventMessage.EventType + " Event Published")
	}
	return true
}
func formatEventRequest(ctx context.Context, eventData common.MessageData) ([]byte, error) {
	for _, event := range eventData.Events {
		if event.OriginOfCondition == nil {
			event.OriginOfCondition = &common.Link{
				Oid: "",
			}
		}
	}
	data1, _ := json.Marshal(eventData)
	return data1, nil
}
