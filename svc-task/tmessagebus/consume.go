//(C) Copyright [2023] Hewlett Packard Enterprise Development LP
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

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
)

var (
	// TaskEventRecvQueue is a channel where task events are pushed
	TaskEventRecvQueue chan<- interface{}
	// TaskEventProcQueue is a channel from where task events popped and processed
	TaskEventProcQueue <-chan interface{}
)

// SubscribeTaskEventsQueue creates a consumer for the task event topic
func SubscribeTaskEventsQueue(topicName string) {
	config.TLSConfMutex.RLock()
	MessageBusConfigFilePath := config.Data.MessageBusConf.MessageBusConfigFilePath
	messagebusType := config.Data.MessageBusConf.MessageBusType
	config.TLSConfMutex.RUnlock()
	// connecting to messagbus
	k, err := dc.Communicator(messagebusType, MessageBusConfigFilePath, topicName)
	if err != nil {
		l.Log.Error("Unable to connect to kafka" + err.Error())
		return
	}
	// subscribe from message bus
	if err := k.Accept(consumeTaskEvents); err != nil {
		l.Log.Error(err.Error())
		return
	}
}

// consumeTaskEvents consume task event messages
func consumeTaskEvents(event interface{}) {
	data, _ := json.Marshal(&event)
	var eventData common.Events
	err := json.Unmarshal(data, &eventData)
	if err != nil {
		l.Log.Error("Error while consuming task events", err)
		return
	}

	var taskEvent dmtf.Event
	err = json.Unmarshal(eventData.Request, &taskEvent)
	if err != nil {
		l.Log.Error("Error while consuming task events", err)
		return
	}

	for _, eventRecord := range taskEvent.Events {
		TaskEventRecvQueue <- eventRecord
	}
}
