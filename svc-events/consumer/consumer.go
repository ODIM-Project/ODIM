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
	log "github.com/sirupsen/logrus"

	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
)

var (
	// In Channel
	In chan<- interface{}
	// Out Channel
	Out <-chan interface{}
	// CtrlMsgRecvQueue is the channel for receiving
	// internal messages read from ODIM-CONTROL-MESSAGES topic
	CtrlMsgRecvQueue chan<- interface{}
	// CtrlMsgProcQueue is the channel for processing
	// internal messages received from ODIM-CONTROL-MESSAGES topic
	CtrlMsgProcQueue <-chan interface{}
)

var done = make(chan bool)

// KafkaSubscriber consume messages from PMB
func KafkaSubscriber(event interface{}) {
	byteData, _ := json.Marshal(&event)
	var kafkaMessage common.Events

	err := json.Unmarshal(byteData, &kafkaMessage)
	if err != nil {
		log.Error("error while unmarshaling the event" + err.Error())
		return
	}
	writeEventToJobQueue(kafkaMessage)
}

// writeEventToJobQueue align events to job queue
func writeEventToJobQueue(kafkaMessage common.Events) {
	// events contains a slice of event subscribed from kafka
	var events = make([]interface{}, 0)
	events = append(events, kafkaMessage)

	go common.RunWriteWorkers(In, events, 5, done)
}

// Consume create a consumer for message bus
// the topic can be defined inside configuration file config.toml
func Consume(topicName string) {
	config.TLSConfMutex.RLock()
	messageQueueConfigFilePath := config.Data.MessageQueueConfigFilePath
	config.TLSConfMutex.RUnlock()
	// connecting to kafka
	k, err := dc.Communicator(dc.KAFKA, messageQueueConfigFilePath)
	if err != nil {
		log.Error("Unable to connect to kafka" + err.Error())
		return
	}
	// subscribe from message bus
	if err := k.Accept(topicName, KafkaSubscriber); err != nil {
		log.Error(err.Error())
		return
	}
	return
}

// SubscribeCtrlMsgQueue creates a consumer for the kafka topic
func SubscribeCtrlMsgQueue(topicName string) {
	config.TLSConfMutex.RLock()
	messageQueueConfigFilePath := config.Data.MessageQueueConfigFilePath
	config.TLSConfMutex.RUnlock()
	// connecting to kafka
	k, err := dc.Communicator(dc.KAFKA, messageQueueConfigFilePath)
	if err != nil {
		log.Error("Unable to connect to kafka" + err.Error())
		return
	}
	// subscribe from message bus
	if err := k.Accept(topicName, consumeCtrlMsg); err != nil {
		log.Error(err.Error())
		return
	}
	return
}

// consumeCtrlMsg consume control messages
func consumeCtrlMsg(event interface{}) {
	var ctrlMessage common.ControlMessageData
	data, _ := json.Marshal(&event)
	if err := json.Unmarshal(data, &ctrlMessage); err != nil {
		log.Error("error while unmarshaling the event" + err.Error())
		return
	}
	msg := []interface{}{ctrlMessage}
	go common.RunWriteWorkers(CtrlMsgRecvQueue, msg, 5, done)
}
