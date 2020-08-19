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
	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	uuid "github.com/satori/go.uuid"
	"log"
)

//Publish will takes the system id,Event type and publishes the data to message bus
func Publish(systemID, eventType, collectionType string) {
	k, err := dc.Communicator(dc.KAFKA, config.Data.MessageQueueConfigFilePath)
	if err != nil {
		log.Println("Unable to connect to kafka", err)
		return
	}

	defer k.Close()
	var event = common.Event{
		EventID:           uuid.NewV4().String(),
		MessageID:         "ResourceEvent.1.0.1." + eventType,
		EventType:         eventType,
		OriginOfCondition: systemID,
	}
	var events = []common.Event{event}
	var messageData = common.MessageData{
		Name:      "Resource Event",
		Context:   "/redfish/v1/$metadata#Event.Event",
		OdataType: "#Event.v1_0_0.Event",
		Events:    events,
	}
	data, _ := json.Marshal(messageData)
	var mbevent = common.Events{
		IP:      collectionType,
		Request: data,
	}

	if err := k.Distribute("REDFISH-EVENTS-TOPIC", mbevent); err != nil {
		log.Println("Unable Publish events to kafka", err)
		return
	}
	log.Println("Event Published")

}
