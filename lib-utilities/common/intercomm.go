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

// Package common ...
package common

// ControlMessage is the data type for ControlMessage types
type ControlMessage int

const (
	// SubscribeEMB is the ControlMessage type to
	// indicate to subscribe plugin EMB
	SubscribeEMB ControlMessage = iota
)

const (
	// InterCommMsgQueueName is the name of the messagebus
	// queue which is used for inter comm events
	InterCommMsgQueueName = "ODIM-CONTROL-MESSAGES"
)

// ControlMessageData holds the control message data
type ControlMessageData struct {
	MessageType ControlMessage
	Data        interface{}
}

// SubscribeEMBData holds data required for subscribing
// to plugin EMB
type SubscribeEMBData struct {
	PluginID  string
	EMBQueues []string
}

// PluginStatusEvent contains details of the plugin status
type PluginStatusEvent struct {
	Name         string `json:"Name"`
	Type         string `json:"Type"`
	Timestamp    string `json:"Timestamp"`
	OriginatorID string `json:"OriginatorID"`
}
