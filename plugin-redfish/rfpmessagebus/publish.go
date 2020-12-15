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

//Package rfpmessagebus ...
package rfpmessagebus

import (
	"encoding/json"
	"fmt"
	"log"

	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/plugin-redfish/config"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpmodel"
)

// Publish function will handle events request in two originofcondition format
// originofcondition can be with or without @odata.id
func Publish(data interface{}) bool {
	if data == nil {
		log.Printf("Error: Invalid data on publishing events")
		return false
	}
	event := data.(common.Events)

	K, err := dc.Communicator(dc.KAFKA, config.Data.MessageBusConf.MessageQueueConfigFilePath)
	if err != nil {
		fmt.Println("Unable communicate with kafka", err)
		return false
	}
	defer K.Close()
	// Since we are deleting the first event from the eventlist,
	// processing the first event
	var message common.MessageData
	err = json.Unmarshal(event.Request, &message)
	if err != nil {
		var messageData rfpmodel.ForwardEventMessageData
		if err := json.Unmarshal(event.Request, &messageData); err != nil {
			log.Printf("error: Failed to unmarshal the event: %v", err)
			return false
		}
		message.Context = messageData.Context
		message.Name = messageData.Name
		message.OdataType = messageData.OdataType
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
				Oid: messageData.Events[i].OriginOfCondition,
			}
			eventData.MessageArgs = messageData.Events[i].MessageArgs
			message.Events = append(message.Events, eventData)
		}
		event.Request, _ = json.Marshal(message)
	}
	topic := config.Data.MessageBusConf.EmbQueue[0]
	if err := K.Distribute(topic, event); err != nil {
		fmt.Println("Unable Publish events to kafka", err)
		return false
	}
	log.Println("forwarded event", string(event.Request))
	for _, eventMessage := range message.Events {
		fmt.Printf("Event %v Published\n", eventMessage.EventType)
	}
	return true
}
