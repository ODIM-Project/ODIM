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

package agmessagebus

import (
	"encoding/json"
	"fmt"
	"time"

	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

//Publish will takes the system id,Event type and publishes the data to message bus
func Publish(systemID, eventType, collectionType string) {
	topicName := config.Data.MessageBusConf.MessageBusQueue[0]
	k, err := dc.Communicator(config.Data.MessageBusConf.MessageBusType, config.Data.MessageBusConf.MessageBusConfigFilePath, topicName)
	if err != nil {
		log.Error("Unable to connect to " + config.Data.MessageBusConf.MessageBusType + " " + err.Error())
		return
	}

	var message string
	switch eventType {
	case "ResourceAdded":
		message = "The resource has been created successfully."
	case "ResourceRemoved":
		message = "The resource has been removed successfully."
	}

	var event = common.Event{
		EventID:        uuid.NewV4().String(),
		MessageID:      "ResourceEvent.1.2.0." + eventType,
		EventTimestamp: time.Now().Format(time.RFC3339),
		EventType:      eventType,
		Message:        message,
		OriginOfCondition: &common.Link{
			Oid: systemID,
		},
		Severity: "OK",
	}
	var events = []common.Event{event}
	var messageData = common.MessageData{
		Name:      "Resource Event",
		Context:   "/redfish/v1/$metadata#Event.Event",
		OdataType: common.EventType,
		Events:    events,
	}
	data, _ := json.Marshal(messageData)
	var mbevent = common.Events{
		IP:      collectionType,
		Request: data,
	}

	if err := k.Distribute(mbevent); err != nil {
		log.Error("Unable Publish events to kafka" + err.Error())
		return
	}
	log.Info("Event Published")

}

// PublishCtrlMsg publishes ODIM control messages to the message bus
func PublishCtrlMsg(msgType common.ControlMessage, msg interface{}) error {
	topicName := config.Data.MessageBusConf.MessageBusQueue[0]
	conn, err := dc.Communicator(config.Data.MessageBusConf.MessageBusType, config.Data.MessageBusConf.MessageBusConfigFilePath, topicName)
	if err != nil {
		return fmt.Errorf("failed to get kafka connection: %s", err.Error())
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %s", err.Error())
	}
	ctrlMsg := common.ControlMessageData{
		MessageType: msgType,
		Data:        data,
	}
	if err := conn.Distribute(ctrlMsg); err != nil {
		return fmt.Errorf("failed to write data to kafka: %s", err.Error())
	}
	return nil
}
