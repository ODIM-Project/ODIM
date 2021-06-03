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

//Package dpmessagebus ...
package dpmessagebus

import (
	"encoding/json"
	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/plugin-dell/config"
	log "github.com/sirupsen/logrus"
)

// Publish ...
func Publish(data interface{}) bool {
	if data == nil {
		log.Error("Invalid data on publishing events")
		return false
	}
	event := data.(common.Events)

	K, err := dc.Communicator(dc.KAFKA, config.Data.MessageBusConf.MessageQueueConfigFilePath)
	if err != nil {
		log.Error("Unable communicate with kafka: " + err.Error())
		return false
	}
	defer K.Close()
	// Since we are deleting the first event from the eventlist,
	// processing the first event
	var message common.MessageData
	err = json.Unmarshal(event.Request, &message)
	if err != nil {
		log.Error("Failed to unmarshal the event, got: " + err.Error())
		return false
	}
	topic := config.Data.MessageBusConf.EmbQueue[0]
	if err := K.Distribute(topic, event); err != nil {
		log.Error("Unable Publish events to kafka: " + err.Error())
		return false
	}
	for _, eventMessage := range message.Events {
		log.Info(eventMessage.EventType + " Event Published")
	}
	return true
}
