//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http:#www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License

package tmessagebus

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"

	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	uuid "github.com/satori/go.uuid"
)

//Publish will takes the taskURI, messageID, Event type and publishes the data to message bus
func Publish(taskURI, messageID, eventType, taskMessage string) {
	topicName := config.Data.MessageBusConf.MessageBusQueue[0]
	k, err := dc.Communicator(config.Data.MessageBusConf.MessageBusType, config.Data.MessageBusConf.MessageBusConfigFilePath, topicName)
	if err != nil {
		log.Error("Unable to connect to " + config.Data.MessageBusConf.MessageBusType + " " + err.Error())
		return
	}

	var eventID = uuid.NewV4().String()
	var event = common.Event{
		EventID:        eventID,
		MessageID:      messageID,
		EventTimestamp: time.Now().Format(time.RFC3339),
		EventType:      eventType,
		Message:        taskMessage,
		OriginOfCondition: &common.Link{
			Oid: taskURI,
		},
		Severity: "OK",
	}
	var events = []common.Event{event}
	var messageData = common.MessageData{
		Name:      "Task Event",
		Context:   "/redfish/v1/$metadata#Event.Event",
		OdataType: common.EventType,
		Events:    events,
	}
	data, _ := json.Marshal(messageData)
	var mbevent = common.Events{
		IP:      "TasksCollection",
		Request: data,
	}

	if err := k.Distribute(mbevent); err != nil {
		log.Error("unable to publish the event to message bus: " + err.Error())
		return
	}
	log.Error("info: TaskURI:" + taskURI + ", EventID:" + eventID + ", MessageID:" + messageID)
}
