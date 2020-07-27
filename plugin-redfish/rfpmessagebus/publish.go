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
)

// Publish ...
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
		log.Printf("error: Failed to unmarshal the event: %v", err)
		return false
	}
	for _, eventMessage := range message.Events {
		topic := config.Data.MessageBusConf.EmbQueue[0]
		if err := K.Distribute(topic, event); err != nil {
			fmt.Println("Unable Publish events to kafka", err)
			return false
		}
		fmt.Printf("Event %v Published\n", eventMessage.EventType)
	}
	return true
}
