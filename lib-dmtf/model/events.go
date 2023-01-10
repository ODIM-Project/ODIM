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

// EventType - This property shall contain an array that contains the types
// of events that shall be sent to the destination.  To specify that a client
// is subscribing for Metric Reports, the EventTypes property should
// include 'MetricReport'.  If the subscription does not include this property,
//the service shall use a single element with a default of `Other`.
// Note: This property has been deprecated.  Starting with Redfish Specification
// v1.6 (Event v1.3), subscriptions are based on the RegistryPrefix and ResourceType
// properties and not on the EventType property.  Use EventFormatType to create
// subscriptions for Metric Reports.
type EventType string

// MessageSeverity - The severity of the message in this event. This property shall
// contain the severity of the message in this event. Services can replace the value
// defined in the message registry with a value more applicable to the implementation.
type MessageSeverity string

//The subscription type for events.
type SubscriptionType string

//DeliveryRetryPolicy - The subscription delivery retry policy for events, where the subscription type is RedfishEvent.
type DeliveryRetryPolicy string

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

	// Subscription Types for events. Currently ODIM only support subscriptions
	// of type RedFishTypeEvent.
	// SubscriptionTypeRedFishEvent - The subscription follows the Redfish
	// Specification for event notifications. To send an event notification,
	// a service sends an HTTP POST to the subscriber's destination URI.
	SubscriptionTypeRedFishEvent SubscriptionType = "RedfishEvent"

	// SubscriptionTySubscriptionTypeSSE - The subscription follows the HTML5
	// server-sent event definition for event notifications.
	SubscriptionTySubscriptionTypeSSE SubscriptionType = "SSE"

	// SubscriptionTypeSNMPTrap - The subscription follows the various versions
	// of SNMP Traps for event notifications.
	SubscriptionTypeSNMPTrap SubscriptionType = "SNMPTrap"

	// SubscriptionTypeSNMPInform - The subscription follows versions 2 and 3 of
	// SNMP Inform for event notifications.
	SubscriptionTypeSNMPInform SubscriptionType = "SNMPInform"

	// SubscriptionTypeSyslog - The subscription sends Syslog messages for
	// event notifications.
	SubscriptionTypeSyslog SubscriptionType = "Syslog"

	// SubscriptionTypeOEM - The subscription is an OEM subscription.
	SubscriptionTypeOEM SubscriptionType = "OEM"

	// DeliveryRetryPolicy for events. Currently ODIM only support subscriptions
	// of type RetryForever.
	// DeliveryRetryForever - The subscription is not suspended or terminated,
	// and attempts at delivery of future events shall continue regardless of
	// the number of retries.
	DeliveryRetryForever DeliveryRetryPolicy = "RetryForever"

	// DeliveryRetryForeverWithBackoff - The subscription is not suspended or
	// terminated, and attempts at delivery of future events shall continue
	// regardless of the number of retries, but issued over time according to
	// a service-defined backoff algorithm
	DeliveryRetryForeverWithBackoff DeliveryRetryPolicy = "RetryForeverWithBackoff"

	// DeliverySuspendRetries - The subscription is suspended after the maximum
	// number of retries is reached
	DeliverySuspendRetries DeliveryRetryPolicy = "SuspendRetries"

	// DeliveryTerminateAfterRetries : The subscription is terminated after the
	// maximum number of retries is reached.
	DeliveryTerminateAfterRetries DeliveryRetryPolicy = "TerminateAfterRetries"
)

func (subscriptionType SubscriptionType) IsValidSubscriptionType() bool {
	switch subscriptionType {
	case SubscriptionTypeRedFishEvent, SubscriptionTypeOEM, SubscriptionTypeSNMPInform, SubscriptionTypeSNMPTrap, SubscriptionTypeSyslog, SubscriptionTySubscriptionTypeSSE:
		return true
	default:
		return false
	}
}

func (subscriptionType SubscriptionType) IsSubscriptionTypeSupported() bool {
	switch subscriptionType {
	case SubscriptionTypeRedFishEvent:
		return true
	default:
		return false
	}
}

func (eventType EventType) IsValidEventType() bool {
	switch eventType {
	case EventTypeAlert, EventTypeMetricReport, EventTypeOther, EventTypeResourceRemoved,
		EventTypeResourceAdded, EventTypeResourceUpdated, EventTypeStatusChange:
		return true
	default:
		return false
	}
}

func (deliveryRetryPolicy DeliveryRetryPolicy) IsValidDeliveryRetryPolicyType() bool {
	switch deliveryRetryPolicy {
	case DeliveryRetryForever, DeliverySuspendRetries, DeliveryTerminateAfterRetries,
		DeliveryRetryForeverWithBackoff:
		return true
	default:
		return false
	}
}

func (deliveryRetryPolicy DeliveryRetryPolicy) IsDeliveryRetryPolicyTypeSupported() bool {
	switch deliveryRetryPolicy {
	case DeliveryRetryForever:
		return true
	default:
		return false
	}
}

// Event schema describes the JSON payload received by an event destination, which has
// subscribed to event notification, when events occur. This resource contains data
// about events, including descriptions, severity, and a message identifier to a
// message registry that can be accessed for further information.
// Refer to Event.v1_7_0.json of the redfish spec for more details
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

// EventRecord - a single event in the Events array of the Event Resource. This  has a
// set of properties that describe a single event. Because Events is an array, more than
// one EventRecord can be sent simultaneously.
// Refer to Event.v1_7_0.json of the redfish spec for more details
type EventRecord struct {
	Actions                    *OemActions `json:"Actions,omitempty"`
	Context                    string      `json:"Context,omitempty"`
	EventGroupID               int         `json:"EventGroupId,omitempty"`
	EventID                    string      `json:"EventId,omitempty"`
	EventTimestamp             string      `json:"EventTimestamp,omitempty"`
	EventType                  string      `json:"EventType"`
	MemberID                   string      `json:"MemberId"`
	Message                    string      `json:"Message,omitempty"`
	MessageArgs                []string    `json:"MessageArgs,omitempty"`
	MessageID                  string      `json:"MessageId"`
	MessageSeverity            string      `json:"MessageSeverity,omitempty"`
	Oem                        interface{} `json:"Oem,omitempty"`
	OriginOfCondition          *Link       `json:"OriginOfCondition,omitempty"`
	Severity                   string      `json:"Severity,omitempty"`
	SpecificEventExistsInGroup bool        `json:"SpecificEventExistsInGroup,omitempty"`
	LogEntry                   *Link       `json:"LogEntry,omitempty"`
}

// This Resource shall represent the target of an event subscription, including
// the event types and context to provide to the target in the Event payload.
// Refer to EventDestination.v1_11_2.json of the redfish spec for more details
type EventDestination struct {
	Context                 string              `json:"Context"`
	EventTypes              []string            `json:"EventTypes"`
	EventFormatType         string              `json:"EventFormatType"`
	ExcludeMessageIds       []string            `json:"ExcludeMessageIds,omitempty"`
	ExcludeRegistryPrefixes []string            `json:"ExcludeRegistryPrefixes,omitempty"`
	DeliveryRetryPolicy     DeliveryRetryPolicy `json:"DeliveryRetryPolicy"`
	Destination             string              `json:"Destination"`
	Id                      string              `json:"id"`
	MessageIds              []string            `json:"MessageIds"`
	Name                    string              `json:"Name"`
	OriginResources         []string            `json:"OriginResources"`
	Protocol                string              `json:"Protocol"`
	ResourceTypes           []string            `json:"ResourceTypes"`
	SubscriptionType        SubscriptionType    `json:"SubscriptionType"`
	SubordinateResources    bool                `json:"SubordinateResources"`
}
