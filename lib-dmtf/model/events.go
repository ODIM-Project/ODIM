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

package model

// EventType - This property shall indicate the type of event.
// "This property has been deprecated.  Starting with Redfish Specification v1.6 (Event v1.3),
// subscriptions are based on the RegistryPrefix and ResourceType properties and not on the EventType property."
type EventType string

// MessageSeverity - The severity of the message in this event. This property shall contain the severity of the message in this event.
// Services can replace the value defined in the message registry with a value more applicable to the implementation.
type MessageSeverity string

const (
	// EventTypeAlert - "Alert": "A condition requires attention."
	EventTypeAlert EventType = "Alert"
	// EventTypeMetricReport - "MetricReport": "The telemetry service is sending a metric report."
	EventTypeMetricReport EventType = "MetricReport"
	// EventTypeOther - "Other": "Because EventType is deprecated as of Redfish Specification v1.6,
	// the event is based on a registry or resource but not an EventType."
	EventTypeOther EventType = "Other"
	// EventTypeResourceAdded - "ResourceAdded": "A resource has been added."
	EventTypeResourceAdded EventType = "ResourceAdded"
	// EventTypeResourceRemoved - "ResourceRemoved": "A resource has been removed."
	EventTypeResourceRemoved EventType = "ResourceRemoved"
	// EventTypeResourceUpdated - "ResourceUpdated": "A resource has been updated."
	EventTypeResourceUpdated EventType = "ResourceUpdated"
	// EventTypeStatusChange - "StatusChange": "The status of a resource has changed."
	EventTypeStatusChange EventType = "StatusChange"

	// MessageSeverityCritical - "Critical": "A critical condition requires immediate attention."
	MessageSeverityCritical MessageSeverity = "Critical"
	// MessageSeverityOK - "OK": "Normal."
	MessageSeverityOK MessageSeverity = "OK"
	// MessageSeverityWarning - "Warning": "A condition requires attention."
	MessageSeverityWarning MessageSeverity = "Warning"
)

// Event schema describes the JSON payload received by an event destination, which has subscribed to event notification, when events occur.
// This resource contains data about events, including descriptions, severity, and a message identifier to a message registry that can be accessed for further information.
type Event struct {
	ODataContext string        `json:"@odata.context,omitempty"`
	ODataType    string        `json:"@odata.type"`
	Actions      *OemActions   `json:"Actions,omitempty"`
	Context      string        `json:"Context,omitempty"`
	Description  string        `json:"Description,omitempty"`
	Events       []EventRecord `json:"Events"`
	EventsCount  int           `json:"Events@odata.count,omitempty"`
	ID           string        `json:"id"`
	Name         string        `json:"Name"`
	Oem          interface{}   `json:"Oem,omitempty"`
}

// EventRecord - This property shall indicate the type of event.
type EventRecord struct {
	Actions                    *OemActions `json:"Actions,omitempty"`
	Context                    string      `json:"Context,omitempty"`
	EventGroupID               int         `json:"EventGroupId,omitempty"`
	EventID                    string      `json:"EventId"`
	EventTimestamp             string      `json:"EventTimestamp"`
	EventType                  string      `json:"EventType"`
	MemberID                   string      `json:"MemberId,omitempty"`
	Message                    string      `json:"Message"`
	MessageArgs                []string    `json:"MessageArgs,omitempty"`
	MessageID                  string      `json:"MessageId"`
	MessageSeverity            string      `json:"MessageSeverity,omitempty"`
	Oem                        interface{} `json:"Oem,omitempty"`
	OriginOfCondition          *Link       `json:"OriginOfCondition,omitempty"`
	Severity                   string      `json:"Severity,omitempty"`
	SpecificEventExistsInGroup bool        `json:"SpecificEventExistsInGroup,omitempty"`
}
