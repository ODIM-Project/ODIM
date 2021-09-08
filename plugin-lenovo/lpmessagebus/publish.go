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

//Package lpmessagebus ...
package lpmessagebus

import (
	"encoding/json"
	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/config"
	log "github.com/sirupsen/logrus"
)

// Publish function will handle events request in two originofcondition format
// originofcondition can be with or without @odata.id
func Publish(data interface{}) bool {
	if data == nil {
		log.Error("Nil data passed to event publisher")
		return false
	}
	event := data.(common.Events)

	K, err := dc.Communicator(dc.KAFKA, config.Data.MessageBusConf.MessageQueueConfigFilePath)
	if err != nil {
		log.Error("Unable communicate with kafka, got:" + err.Error())
		return false
	}
	defer K.Close()
	// Since we are deleting the first event from the eventlist,
	// processing the first event
	var message common.MessageData
	err = json.Unmarshal(event.Request, &message)
	if err != nil {
		var messageData dmtf.Event
		if err := json.Unmarshal(event.Request, &messageData); err != nil {
			log.Error("Failed to unmarshal the event: " + err.Error())
			return false
		}
		message.Context = messageData.Context
		message.Name = messageData.Name
		message.OdataType = messageData.ODataType
		message.Events = make([]common.Event, 0)
		for i := 0; i < len(messageData.Events); i++ {
			var eventData common.Event
			eventData.EventGroupID = messageData.Events[i].EventGroupID
			eventData.EventID = messageData.Events[i].EventID
			eventData.EventTimestamp = messageData.Events[i].EventTimestamp
			eventData.EventType = messageData.Events[i].EventType
			eventData.Message = messageData.Events[i].Message
			eventData.MemberID = messageData.Events[i].MemberID
			eventData.Severity = messageData.Events[i].Severity
			eventData.Oem = messageData.Events[i].Oem
			eventData.MessageID = messageData.Events[i].MessageID
			eventData.OriginOfCondition = &common.Link{
				Oid: messageData.Events[i].OriginOfCondition.Oid,
			}
			eventData.MessageArgs = messageData.Events[i].MessageArgs
			message.Events = append(message.Events, eventData)
		}
		event.Request, _ = json.Marshal(message)
	}
	topic := config.Data.MessageBusConf.EmbQueue[0]
	if err := K.Distribute(topic, event); err != nil {
		log.Error("Unable Publish events to kafka, got:" + err.Error())
		return false
	}
	log.Info("forwarded event" + string(event.Request))
	for _, eventMessage := range message.Events {
		log.Info(eventMessage.EventType + " Event published")
	}
	return true
}
