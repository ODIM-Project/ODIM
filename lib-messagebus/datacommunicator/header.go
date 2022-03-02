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

// Package datacommunicator ...
package datacommunicator

// -----------------------------------------------------------------------------
// IMPORT Section
// -----------------------------------------------------------------------------
import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// BrokerType defines the underline MQ platform to be selected for the
// messages. KAFKA and RedisStremas platforms are supported.
const (
	KAFKA                = "Kafka"        // KAFKA as Messaging Platform, Please use this ID
	REDISSTREAMS         = "RedisStreams" // REDISSTREAMS as Messaging Platform
	EVENTREADERGROUPNAME = "eventreaders_grp"
)

// MQBus Interface defines the Process interface function (Only function user
// should call). These functions are implemented as part of Packet struct.
// Distribute - API to Publish Messages into specified Pipe (Topic / Subject)
// Accept - Consume the incoming message if subscribed by that component
// Get - Would initiate blocking call to remote process to get response
// Close - Would disconnect the connection with Middleware.
type MQBus interface {
	Distribute(data interface{}) error
	Accept(fn MsgProcess) error
	Get(pipe string, d interface{}) interface{}
	Remove() error
	Close() error
}

// MsgProcess defines the functions for processing accepted messages. Any client
// who wants to accept and handle the events / notifications / messages, should
// implement this function as part of their procedure. That same function should
// be sent to MessageBus as callback for handling the incoming messages.
type MsgProcess func(d interface{})

// Packet defines all the message related information that Producer or Consumer
// should know for message transactions. Both Producer and Consumer use this
// same structure for message transactions.
// BrokerType - Refer above defined Constants for possible values
// DataResponder - Refer HandleResponse Type description
type Packet struct {
	// BrokerType defines the underline MQ platform
	BrokerType string
}

// Communicator defines the Broker platform Middleware selection and corresponding
// communication object would be created to send / receive the messages. Broker
// type would be stored as part of Connection Object "Packet".
// TODO: We would be looking into Kafka Synchronous communication API for providing
// support for Sync Communication Model in MessageBus
func Communicator(bt string, messageQueueConfigPath, pipe string) (MQBus, error) {

	// Defining pointer for KAFKA Connection Objects Based on
	// BrokerType value, Middleware Connection will be created. Also we would be
	// storing maintain the connections as a Map (Between Connection and Pipe)
	var kp *KafkaPacket
	var rp *RedisStreamsPacket
	switch bt {
	case KAFKA:
		kp = new(KafkaPacket)
		kp.BrokerType = bt
		kp.messageBusConfigFile = messageQueueConfigPath
		kp.pipe = pipe
		return kp, nil
	case REDISSTREAMS:
		rp = new(RedisStreamsPacket)
		rp.BrokerType = bt
		rp.pipe = pipe
		return rp, nil
	default:
		return nil, fmt.Errorf("Broker: \"Broker Type\" is not supported - %s", bt)
	}
}

// Encode converts the interface into Byte stream (ENCODE).
func Encode(d interface{}) ([]byte, error) {

	data, err := json.Marshal(d)
	if err != nil {
		log.Error("Failed to encode the given event data: " + err.Error())
		return nil, err
	}
	return data, nil
}

// Decode converts the byte stream into Data (DECODE).
///data will  be masked as Interface before sent to Consumer or Requester.
func Decode(d []byte, a interface{}) error {
	err := json.Unmarshal(d, &a)
	if err != nil {
		log.Error("error: Failed to decode the event data: " + err.Error())
		return err
	}
	return nil
}
