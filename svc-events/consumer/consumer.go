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

// Package consumer will have config details of kafka topic
// and also have the functionality of consuming the events from the kafka on
// corresponding topics
package consumer

import (
	"encoding/json"
	"time"

	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
)

var (
	// In Channel
	In chan<- interface{}
	// Out Channel
	Out <-chan interface{}
	// CtrlMsgRecvQueue is the channel for receiving
	// internal messages read from intercomm message bus queue
	CtrlMsgRecvQueue chan<- interface{}
	// CtrlMsgProcQueue is the channel for processing
	// internal messages received from intercomm messae bus queue
	CtrlMsgProcQueue <-chan interface{}
)

// EventSubscriber consume messages from PMB
func EventSubscriber(event interface{}) {
	byteData, _ := json.Marshal(&event)
	var message common.Events

	err := json.Unmarshal(byteData, &message)
	if err != nil {
		l.Log.Error("error while unmarshaling the event" + err.Error())
		return
	}
	writeEventToJobQueue(message)
}

// writeEventToJobQueue align events to job queue
func writeEventToJobQueue(message common.Events) {
	// events contains a slice of event subscribed from kafka
	events := []interface{}{message}
	go func() {
		// Wait for the write workers to finish writing to
		// In buffer and clear the memory assigned to the data
		ticker := time.NewTicker(500 * time.Millisecond)
		done := make(chan bool)
		breakLoop := false
		workerCount := 1
		common.RunWriteWorkers(In, events, workerCount, done)
		for !breakLoop {
			select {
			case <-done:
				workerCount--
				if workerCount == 0 {
					breakLoop = true
					break
				}
			case <-ticker.C:
			}
		}
		// empty the slice passed to RunWriteWorkers for GC
		events = nil
		// empty the slice in the passed message data for GC
		message.Request = nil
		ticker.Stop()
		close(done)
	}()
}

// Consume create a consumer for message bus
// the topic can be defined inside configuration file config.toml
func Consume(topicName string) {
	config.TLSConfMutex.RLock()
	MessageBusConfigFilePath := config.Data.MessageBusConf.MessageBusConfigFilePath
	messagebusType := config.Data.MessageBusConf.MessageBusType
	config.TLSConfMutex.RUnlock()
	// connecting to kafka
	k, err := dc.Communicator(messagebusType, MessageBusConfigFilePath, topicName)
	if err != nil {
		l.Log.Error("Unable to connect to kafka" + err.Error())
		return
	}
	// subscribe from message bus
	if err := k.Accept(EventSubscriber); err != nil {
		l.Log.Error(err.Error())
		return
	}
}

// SubscribeCtrlMsgQueue creates a consumer for the kafka topic
func SubscribeCtrlMsgQueue(topicName string) {
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
	if err := k.Accept(consumeCtrlMsg); err != nil {
		l.Log.Error(err.Error())
		return
	}
}

// consumeCtrlMsg consume control messages
func consumeCtrlMsg(event interface{}) {
	var ctrlMessage common.ControlMessageData
	done := make(chan bool)
	data, _ := json.Marshal(&event)
	var redfishEvent common.Events
	// verifying the incoming event to check whether it's of type common events or control message data
	if err := json.Unmarshal(data, &redfishEvent); err == nil {
		writeEventToJobQueue(redfishEvent)
	} else {
		if err := json.Unmarshal(data, &ctrlMessage); err != nil {
			l.Log.Error("error while unmarshal the event" + err.Error())
			return
		}
	}
	msg := []interface{}{ctrlMessage}
	go common.RunWriteWorkers(CtrlMsgRecvQueue, msg, 1, done)
	for range done {
		break
	}
	msg = nil
	close(done)
}
