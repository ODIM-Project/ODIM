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
// the service shall use a single element with a default of `Other`.
// Note: This property has been deprecated.  Starting with Redfish Specification
// v1.6 (Event v1.3), subscriptions are based on the RegistryPrefix and ResourceType
// properties and not on the EventType property.  Use EventFormatType to create
// subscriptions for Metric Reports.
type EventType string

// MessageSeverity - The severity of the message in this event. This property shall
// contain the severity of the message in this event. Services can replace the value
// defined in the message registry with a value more applicable to the implementation.
type MessageSeverity string

// SNMPAuthenticationProtocols - The authentication protocol for SNMPv3.
type SNMPAuthenticationProtocols string

// SNMPEncryptionProtocols - the SNMPv3 encryption protocol.
type SNMPEncryptionProtocols string

// SubscriptionType - The subscription type for events.
type SubscriptionType string

// SyslogFacility - This property shall contain the types of programs that can log messages.
// If this property contains an empty array or is absent, all facilities
// shall be indicated. Facility values are described in the RFC5424.
type SyslogFacility string

// This property shall contain the lowest syslog severity level
// that will be forwarded. The service shall forward all messages
// equal to or greater than the value in this property.  The value
// `All` shall indicate all severities.

// DeliveryRetryPolicy - The subscription delivery retry policy for events,
// where the subscription type is RedfishEvent.
type DeliveryRetryPolicy string

// SyslogSeverity - This property shall contain the types of severity for syslog
type SyslogSeverity string

// SMTPAuthentication - This property shall contain the authentication method for the SMTP server.
type SMTPAuthentication string

// SMTPConnectionProtocol - This property shall contain the connection type to the outgoing SMTP server.
type SMTPConnectionProtocol string

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

	// SNMPAuthenticationProtocolsCommunityString - authentication using SNMP community strings and the value of TrapCommunity.
	SNMPAuthenticationProtocolsCommunityString SNMPAuthenticationProtocols = "CommunityString"

	// SNMPAuthenticationProtocolsHMAC128SHA224 - authentication for SNMPv3 access conforms to the RFC7860-defined usmHMAC128SHA224AuthProtocol.
	SNMPAuthenticationProtocolsHMAC128SHA224 SNMPAuthenticationProtocols = "HMAC128_SHA224"

	// SNMPAuthenticationProtocolsHMAC192SHA256 - authentication for SNMPv3 access conforms to the RFC7860-defined usmHMAC192SHA256AuthProtocol.
	SNMPAuthenticationProtocolsHMAC192SHA256 SNMPAuthenticationProtocols = "HMAC192_SHA256"

	// SNMPAuthenticationProtocolsHMAC256SHA384 - authentication for SNMPv3 access conforms to the RFC7860-defined usmHMAC256SHA384AuthProtocol.
	SNMPAuthenticationProtocolsHMAC256SHA384 SNMPAuthenticationProtocols = "HMAC256_SHA384"

	// SNMPAuthenticationProtocolsHMAC384SHA512 - authentication for SNMPv3 access conforms to the RFC7860-defined usmHMAC384SHA512AuthProtocol.
	SNMPAuthenticationProtocolsHMAC384SHA512 SNMPAuthenticationProtocols = "HMAC384_SHA512"

	// SNMPAuthenticationProtocolsHMACMD5 - authentication conforms to the RFC1321-defined HMAC-MD5-96 authentication protocol.
	SNMPAuthenticationProtocolsHMACMD5 SNMPAuthenticationProtocols = "HMAC_MD5"

	// SNMPAuthenticationProtocolsHMACSHA96 - authentication conforms to the RFC3414-defined HMAC-SHA-96 authentication protocol.
	SNMPAuthenticationProtocolsHMACSHA96 SNMPAuthenticationProtocols = "HMAC_SHA96"

	// SNMPAuthenticationProtocolsNone - authentication is not required.
	SNMPAuthenticationProtocolsNone SNMPAuthenticationProtocols = "None"

	// SNMPEncryptionProtocolsCBCDES - encryption conforms to the RFC3414-defined CBC-DES encryption protocol.
	SNMPEncryptionProtocolsCBCDES SNMPEncryptionProtocols = "CBC_DES"

	// SNMPEncryptionProtocolsCFB128AES128 - encryption conforms to the RFC3414-defined CFB128-AES-128 encryption protocol.
	SNMPEncryptionProtocolsCFB128AES128 SNMPEncryptionProtocols = "CFB128_AES128"

	// SNMPEncryptionProtocolsNone - there is no encryption.
	SNMPEncryptionProtocolsNone SNMPEncryptionProtocols = "None"

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

	// SyslogFacilityAuth - Security/authentication messages.
	SyslogFacilityAuth SyslogFacility = "Auth"

	// SyslogFacilityAuthpriv - Security/authentication messages.
	SyslogFacilityAuthpriv SyslogFacility = "Authpriv"

	// SyslogFacilityConsole - Log alert.
	SyslogFacilityConsole SyslogFacility = "Console"

	// SyslogFacilityCron -Clock daemon.
	SyslogFacilityCron SyslogFacility = "Cron"

	// SyslogFacilityDaemon - System daemons.
	SyslogFacilityDaemon SyslogFacility = "Daemon"

	// SyslogFacilityFTP - FTP daemon.
	SyslogFacilityFTP SyslogFacility = "FTP"

	// SyslogFacilityKern - Kernel messages.
	SyslogFacilityKern SyslogFacility = "Kern"

	// SyslogFacilityLPR - Line printer subsystem.
	SyslogFacilityLPR SyslogFacility = "LPR"

	// SyslogFacilityLocal0 - Locally used facility 0.
	SyslogFacilityLocal0 SyslogFacility = "Local0"

	// SyslogFacilityLocal1 Locally used facility 1.
	SyslogFacilityLocal1 SyslogFacility = "Local1"

	// SyslogFacilityLocal2 - Locally used facility 2.
	SyslogFacilityLocal2 SyslogFacility = "Local2"

	// SyslogFacilityLocal3 - Locally used facility 3.
	SyslogFacilityLocal3 SyslogFacility = "Local3"

	// SyslogFacilityLocal4 - Locally used facility 4.
	SyslogFacilityLocal4 SyslogFacility = "Local4"

	// SyslogFacilityLocal5 - Locally used facility 5.
	SyslogFacilityLocal5 SyslogFacility = "Local5"

	// SyslogFacilityLocal6 - Locally used facility 6.
	SyslogFacilityLocal6 SyslogFacility = "Local6"

	// SyslogFacilityLocal7 - Locally used facility 7.
	SyslogFacilityLocal7 SyslogFacility = "Local7"

	// SyslogFacilityMail - Mail system.
	SyslogFacilityMail SyslogFacility = "Mail"

	// SyslogFacilityNTP - NTP subsystem.
	SyslogFacilityNTP SyslogFacility = "NTP"

	// SyslogFacilityNews - Network news subsystem.
	SyslogFacilityNews SyslogFacility = "News"

	// SyslogFacilitySecurity - Log audit.
	SyslogFacilitySecurity SyslogFacility = "Security"

	// SyslogFacilitySolarisCron - Scheduling daemon.
	SyslogFacilitySolarisCron SyslogFacility = "SolarisCron"

	// SyslogFacilitySyslog - Messages generated internally by syslogd.
	SyslogFacilitySyslog SyslogFacility = "Syslog"

	// SyslogFacilityUUCP - UUCP subsystem.
	SyslogFacilityUUCP SyslogFacility = "UUCP"

	// SyslogFacilityUser - User-level messages.
	SyslogFacilityUser SyslogFacility = "User"

	// SyslogSeverityAlert - A condition that should be corrected immediately, such as a corrupted system database.
	SyslogSeverityAlert SyslogSeverity = "Alert"

	// SyslogSeverityAll - A message of any severity.
	SyslogSeverityAll SyslogSeverity = "All"

	// SyslogSeverityCritical - Hard device errors.
	SyslogSeverityCritical SyslogSeverity = "Critical"

	// SyslogSeverityDebug - Messages that contain information normally of use only when debugging a program.
	SyslogSeverityDebug SyslogSeverity = "Debug"

	// SyslogSeverityEmergency - A panic condition.
	SyslogSeverityEmergency SyslogSeverity = "Emergency"

	// SyslogSeverityError - An Error.
	SyslogSeverityError SyslogSeverity = "Error"

	// SyslogSeverityInformational - Informational only.
	SyslogSeverityInformational SyslogSeverity = "Informational"

	// SyslogSeverityNotice - Conditions that are not error conditions, but that may require special handling.
	SyslogSeverityNotice SyslogSeverity = "Notice"

	// SyslogSeverityWarning - A Warning.
	SyslogSeverityWarning SyslogSeverity = "Warning"

	// DeliveryRetryForever - DeliveryRetryPolicy for events. Currently ODIM only support subscriptions
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

	// SMTPAuthenticationAutoDetect - SMTP Authentcation method Auto-detect
	SMTPAuthenticationAutoDetect SMTPAuthentication = "AutoDetect"

	// SMTPAuthenticationCRAMMD5 - SMTPAuthentication method CRAM_MD5
	SMTPAuthenticationCRAMMD5 SMTPAuthentication = "CRAM_MD5"

	// SMTPAuthenticationNone - SMTPAuthentication method None species no authentication
	SMTPAuthenticationNone SMTPAuthentication = "None"

	// SMTPAuthenticationPlain - SMTPAuthentication method Plain
	SMTPAuthenticationPlain SMTPAuthentication = "Plain"

	// SMTPConnectionProtocolAutoDetect - SMTPConnection protocol type Auto-detect
	SMTPConnectionProtocolAutoDetect SMTPConnectionProtocol = "AutoDetect"

	// SMTPConnectionProtocolNone - SMTPConnection Protocol type clear-text
	SMTPConnectionProtocolNone SMTPConnectionProtocol = "None"

	// SMTPConnectionProtocolStartTLS - SMTPConnection Protocol type start-TLS
	SMTPConnectionProtocolStartTLS SMTPConnectionProtocol = "StartTLS"

	// SMTPConnectionProtocolTLSSSL - SMTPConnection Protocol type TLS_SSL
	SMTPConnectionProtocolTLSSSL SMTPConnectionProtocol = "TLS_SSL"
)

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

// EventDestination represents the target of an event subscription,
// including the event types and context to provide to the target
// in the Event payload.
// Reference	                : EventDestination.v1_12_0.json
type EventDestination struct {
	ODataContext                 string              `json:"@odata.context,omitempty"`
	ODataEtag                    string              `json:"@odata.etag,omitempty"`
	ODataID                      string              `json:"@odata.id"`
	ODataType                    string              `json:"@odata.type"`
	Actions                      *Actions            `json:"Actions,omitempty"`
	Certificates                 *Link               `json:"Certificates,omitempty"`
	ClientCertificates           *Link               `json:"ClientCertificates,omitempty"`
	Context                      string              `json:"Context"`
	DeliveryRetryPolicy          DeliveryRetryPolicy `json:"DeliveryRetryPolicy,omitempty"`
	Description                  string              `json:"Description,omitempty"`
	Destination                  string              `json:"Destination"`
	EventFormatType              string              `json:"EventFormatType,omitempty"`
	EventTypes                   []string            `json:"EventTypes,omitempty"`
	ExcludeMessageIds            []string            `json:"ExcludeMessageIds,omitempty"`
	ExcludeRegistryPrefixes      []string            `json:"ExcludeRegistryPrefixes,omitempty"`
	HeartbeatIntervalMins        int                 `json:"HeartbeatIntervalMinutes,omitempty"`
	HTTPHeaders                  []string            `json:"HttpHeaders,omitempty"`
	ID                           string              `json:"Id"`
	IncludeOriginOfCondition     bool                `json:"IncludeOriginOfCondition,omitempty"`
	MessageIds                   []string            `json:"MessageIds,omitempty"`
	MetricReportDefinitions      *Link               `json:"MetricReportDefinitions,omitempty"`
	MetricReportDefinitionsCount int                 `json:"MetricReportDefinitions@odata.count,omitempty"`
	Name                         string              `json:"Name"`
	OEMProtocol                  string              `json:"OEMProtocol,omitempty"`
	OEMSubscriptionType          string              `json:"OEMSubscriptionType"`
	Oem                          interface{}         `json:"Oem,omitempty"`
	OriginResources              []Link              `json:"OriginResources,omitempty"`
	OriginResourcesCount         int                 `json:"OriginResources@odata.count,omitempty"`
	Protocol                     string              `json:"Protocol"`
	RegistryPrefixes             []string            `json:"RegistryPrefixes,omitempty"`
	ResourceTypes                []string            `json:"ResourceTypes,omitempty"`
	SNMP                         *SNMPSettings       `json:"SNMP,omitempty"`
	SendHeartbeat                bool                `json:"SendHeartbeat,omitempty"`
	Status                       *Status             `json:"Status,omitempty"`
	SubordinateResources         bool                `json:"SubordinateResources,omitempty"`
	SubscriptionType             SubscriptionType    `json:"SubscriptionType,omitempty"`
	SyslogFilters                *SyslogFilter       `json:"SyslogFilters,omitempty"`
	VerifyCertificate            bool                `json:"VerifyCertificate,omitempty"`
}

// IsValidSubscriptionType validate subscription type is valid,
func (subscriptionType SubscriptionType) IsValidSubscriptionType() bool {
	switch subscriptionType {
	case SubscriptionTypeRedFishEvent, SubscriptionTypeOEM,
		SubscriptionTypeSNMPInform, SubscriptionTypeSNMPTrap,
		SubscriptionTypeSyslog, SubscriptionTySubscriptionTypeSSE:
		return true
	default:
		return false
	}
}

// IsSubscriptionTypeSupported method return true if subscription type is RedfishEvent
func (subscriptionType SubscriptionType) IsSubscriptionTypeSupported() bool {
	switch subscriptionType {
	case SubscriptionTypeRedFishEvent:
		return true
	default:
		return false
	}
}

// ToString - converts SubscriptionType to string type
func (subscriptionType SubscriptionType) ToString() string {
	return string(subscriptionType)
}

// ToString - converts DeliveryRetryPolicy to string type
func (deliveryRetryPolicy DeliveryRetryPolicy) ToString() string {
	return string(deliveryRetryPolicy)
}

// ToString - converts EventType to string type
func (eventType EventType) ToString() string {
	return string(eventType)
}

// IsValidEventType return true if event type is valid
func (eventType EventType) IsValidEventType() bool {
	switch eventType {
	case EventTypeAlert, EventTypeMetricReport, EventTypeOther,
		EventTypeResourceRemoved, EventTypeResourceAdded,
		EventTypeResourceUpdated, EventTypeStatusChange:
		return true
	default:
		return false
	}
}

// IsValidDeliveryRetryPolicyType is validate DeliveryRetryPolicy value valid or not
func (deliveryRetryPolicy DeliveryRetryPolicy) IsValidDeliveryRetryPolicyType() bool {
	switch deliveryRetryPolicy {
	case DeliveryRetryForever, DeliverySuspendRetries,
		DeliveryTerminateAfterRetries, DeliveryRetryForeverWithBackoff:
		return true
	default:
		return false
	}
}

// IsDeliveryRetryPolicyTypeSupported is return true if DeliveryRetryPolicy
// is RetryForever. Currently ODIM support RetryForever value
func (deliveryRetryPolicy DeliveryRetryPolicy) IsDeliveryRetryPolicyTypeSupported() bool {
	switch deliveryRetryPolicy {
	case DeliveryRetryForever:
		return true
	default:
		return false
	}
}

// EventDestinationAction contain the available
// actions for this resource..
// Reference	                : EventDestination.v1_12_0.json
type EventDestinationAction struct {
	ResumeSubscription  *ResumeSubscription  `json:"ResumeSubscription,omitempty"`
	SuspendSubscription *SuspendSubscription `json:"SuspendSubscription,omitempty"`
	Oem                 *OemActions          `json:"OemActions,omitempty"`
}

// ResumeSubscription - This action shall resume a suspended event subscription,
// which affects the subscription status. The service may deliver
// buffered events when the subscription is resumed.
// Reference	                : EventDestination.v1_12_0.json
type ResumeSubscription struct {
	Target string `json:"target,omitempty"`
	Title  string `json:"title,omitempty"`
}

// SNMPSettings - This type shall contain the settings for an SNMP event destination.
// Reference	                : EventDestination.v1_12_0.json
type SNMPSettings struct {
	AuthenticationKey      string                      `json:"AuthenticationKey,omitempty"`
	AuthenticationKeySet   bool                        `json:"AuthenticationKeySet,omitempty"`
	AuthenticationProtocol SNMPAuthenticationProtocols `json:"AuthenticationProtocol,omitempty"`
	EncryptionKey          string                      `json:"EncryptionKey,omitempty"`
	EncryptionKeySet       bool                        `json:"EncryptionKeySet,omitempty"`
	EncryptionProtocol     SNMPEncryptionProtocols     `json:"EncryptionProtocol,omitempty"`
	TrapCommunity          string                      `json:"TrapCommunity,omitempty"`
}

// SuspendSubscription - This action shall suspend an event subscription.  No events shall be
// sent to the event destination until invocation of the ResumeSubscription
// action.  The value of the State property within Status shall contain
// `Disabled` for a suspended subscription.  The service may buffer
// events while the subscription is suspended.
// Reference	                : EventDestination.v1_12_0.json
type SuspendSubscription struct {
	Target string `json:"target,omitempty"`
	Title  string `json:"title,omitempty"`
}

// SyslogFilter - A list of filters applied to syslog messages before sending
// to a remote syslog server.  An empty list indicates all
// syslog messages are sent.
// Reference	                : EventDestination.v1_12_0.json
type SyslogFilter struct {
	LogFacilities  SyslogFacility `json:"LogFacilities,omitempty"`
	LowestSeverity SyslogSeverity `json:"LowestSeverity,omitempty"`
}

// EventService - The EventService schema contains properties for
// managing event subscriptions and generates the
// events sent to subscribers. The resource has
// links to the actual collection of subscriptions,
// which are called event destinations.
// Reference                    : EventService.v1_8_0.json
type EventService struct {
	ODataContext                      string                        `json:"@odata.context,omitempty"`
	ODataEtag                         string                        `json:"@odata.etag,omitempty"`
	ODataID                           string                        `json:"@odata.id"`
	ODataType                         string                        `json:"@odata.type"`
	Actions                           *EventServiceActions          `json:"Actions,omitempty"`
	ID                                string                        `json:"Id"`
	Name                              string                        `json:"Name"`
	DeliveryRetryAttempts             int                           `json:"DeliveryRetryAttempts,omitempty"`
	DeliveryRetryIntervalSeconds      int                           `json:"DeliveryRetryIntervalSeconds,omitempty"`
	Description                       string                        `json:"Description,omitempty"`
	EventFormatTypes                  []string                      `json:"EventFormatTypes,omitempty"`
	EventTypesForSubscription         []string                      `json:"EventTypesForSubscription,omitempty"`
	ExcludeMessageID                  bool                          `json:"ExcludeMessageId,omitempty"`
	ExcludeRegistryPrefix             bool                          `json:"ExcludeRegistryPrefix,omitempty"`
	IncludeOriginOfConditionSupported bool                          `json:"IncludeOriginOfConditionSupported,omitempty"`
	RegistryPrefixes                  []string                      `json:"RegistryPrefixes,omitempty"`
	ResourceTypes                     []string                      `json:"ResourceTypes,omitempty"`
	ServerSentEventURI                string                        `json:"ServerSentEventUri,omitempty"`
	ServiceEnabled                    bool                          `json:"ServiceEnabled,omitempty"`
	SMTP                              *SMTPSettings                 `json:"SMTP,omitempty"`
	SSEFilterPropertiesSupported      *SSEFilterPropertiesSupported `json:"SSEFilterPropertiesSupported,omitempty"`
	Status                            *Status                       `json:"Status,omitempty"`
	SubordinateResourcesSupported     bool                          `json:"SubordinateResourcesSupported,omitempty"`
	Subscriptions                     *[]EventDestination           `json:"Subscriptions,omitempty"`
	Oem                               interface{}                   `json:"Oem,omitempty"`
}

// SMTPSettings - This type shall contain settings for SMTP event delivery.
// Reference                    : EventService.v1_8_0.json
type SMTPSettings struct {
	Authentication     SMTPAuthentication     `json:"Authentication,omitempty"`
	ConnectionProtocol SMTPConnectionProtocol `json:"ConnectionProtocol,omitempty"`
	FromAddress        string                 `json:"FromAddress,omitempty"`
	Password           string                 `json:"Password,omitempty"`
	Port               int                    `json:"Port,omitempty"`
	ServerAddress      string                 `json:"ServerAddress,omitempty"`
	ServiceEnabled     string                 `json:"ServiceEnabled,omitempty"`
	Username           string                 `json:"Username,omitempty"`
}

// SSEFilterPropertiesSupported - The type shall contain a set of properties that are supported
// in the `$filter` query parameter for the URI indicated by the
// ServerSentEventUri property, as described by the Redfish Specification.
// Reference                    : EventService.v1_8_0.json
type SSEFilterPropertiesSupported struct {
	EventFormatType         bool `json:"EventFormatType,omitempty"`
	EventType               bool `json:"EvenType,omitempty"`
	MessageIds              bool `json:"MessageIds,omitempty"`
	MetricReportDefinitions bool `json:"MetricReportDefinitions,omitempty"`
	OriginResource          bool `json:"OriginResource,omitempty"`
	RegistryPrefix          bool `json:"RegistryPrefix,omitempty"`
	ResourceType            bool `json:"ResourceType,omitempty"`
	SubordinateResources    bool `json:"SubordinateResources,omitempty"`
}

// EventServiceActions - The type shall contain actions on event service
type EventServiceActions struct {
	SubmitTestEvent *SubmitTestEvent `json:"SubmitTestEvent,omitempty"`
	Oem             *OemActions      `json:"OemActions,omitempty"`
}

// SubmitTestEvent - This action shall add a test event to the event service
// with the event data specified in the action parameter
// action parameters.  Then, this message should be sent
// to any appropriate event destinations.
// Reference                    : EventService.v1_8_0.json
type SubmitTestEvent struct {
	EventGroupID      int             `json:"EventGroupId,omitempty"`
	EventID           string          `json:"EventId,omitempty"`
	EventTimestamp    string          `json:"EventTimestamp,omitempty"`
	EventType         EventType       `json:"EventType,omitempty"`
	Message           string          `json:"Message,omitempty"`
	MessageArgs       []string        `json:"MessageArgs,omitempty"`
	MessageID         string          `json:"MessageId"`
	OriginOfCondition string          `json:"OriginOfCondition,omitempty"`
	Severity          MessageSeverity `json:"Severity,omitempty"`
}
