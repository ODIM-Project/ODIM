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

//Package rfpmodel ...
package rfpmodel

// Event stores Published Event request and IP of source
type Event struct {
	IP      string
	Request []byte
}

// ForwardEventMessageData contains information of Events and message details including arguments
// it will be send as byte stream on the wire to/from kafka
type ForwardEventMessageData struct {
	OdataType string         `json:"@odata.type"`
	Name      string         `json:"Name"`
	Context   string         `json:"@odata.context"`
	Events    []ForwardEvent `json:"Events"`
}

// ForwardEvent contains the details of the event subscribed from PMB
type ForwardEvent struct {
	MemberID          string      `json:"MemberId,omitempty"`
	EventType         string      `json:"EventType"`
	EventGroupID      int         `json:"EventGroupId,omitempty"`
	EventID           string      `json:"EventId"`
	Severity          string      `json:"Severity"`
	EventTimestamp    string      `json:"EventTimestamp"`
	Message           string      `json:"Message"`
	MessageArgs       []string    `json:"MessageArgs,omitempty"`
	MessageID         string      `json:"MessageId"`
	Oem               interface{} `json:"Oem,omitempty"`
	OriginOfCondition string      `json:"OriginOfCondition,omitempty"`
}
