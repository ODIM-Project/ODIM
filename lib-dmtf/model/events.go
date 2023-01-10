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

// The authentication protocol for SNMPv3.
type SNMPAuthenticationProtocols string

// the SNMPv3 encryption protocol.
type SNMPEncryptionProtocols string

// The subscription type for events.
type SubscriptionType string

// This property shall contain the types of programs that can log messages.  
// If this property contains an empty array or is absent, all facilities 
// shall be indicated. Facility values are described in the RFC5424.
type SyslogFacility string

// This property shall contain the lowest syslog severity level 
// that will be forwarded. The service shall forward all messages 
// equal to or greater than the value in this property.  The value 
// `All` shall indicate all severities.
type SyslogSeverity string

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
	
	// authentication using SNMP community strings and the value of TrapCommunity.
	SNMPAuthenticationProtocolsCommunityString SNMPAuthenticationProtocols = "CommunityString"
	
	// authentication for SNMPv3 access conforms to the RFC7860-defined usmHMAC128SHA224AuthProtocol.
	SNMPAuthenticationProtocolsHMAC128_SHA224 SNMPAuthenticationProtocols = "HMAC128_SHA224"
	
	// authentication for SNMPv3 access conforms to the RFC7860-defined usmHMAC192SHA256AuthProtocol.
	SNMPAuthenticationProtocolsHMAC192_SHA256 SNMPAuthenticationProtocols = "HMAC192_SHA256"
	
	// authentication for SNMPv3 access conforms to the RFC7860-defined usmHMAC256SHA384AuthProtocol.
	SNMPAuthenticationProtocolsHMAC256_SHA384 SNMPAuthenticationProtocols = "HMAC256_SHA384"
	
	// authentication for SNMPv3 access conforms to the RFC7860-defined usmHMAC384SHA512AuthProtocol.
	SNMPAuthenticationProtocolsHMAC384_SHA512 SNMPAuthenticationProtocols = "HMAC384_SHA512"
	
	// authentication conforms to the RFC1321-defined HMAC-MD5-96 authentication protocol.
	SNMPAuthenticationProtocolsHMAC_MD5 SNMPAuthenticationProtocols = "HMAC_MD5"
	
	// authentication conforms to the RFC3414-defined HMAC-SHA-96 authentication protocol.
	SNMPAuthenticationProtocolsHMAC_SHA96 SNMPAuthenticationProtocols = "HMAC_SHA96"
	
	// authentication is not required.
	SNMPAuthenticationProtocolsNone SNMPAuthenticationProtocols = "None"
	
	// encryption conforms to the RFC3414-defined CBC-DES encryption protocol.
	SNMPEncryptionProtocolsCBC_DES SNMPEncryptionProtocols = "CBC_DES"
	
	// encryption conforms to the RFC3414-defined CFB128-AES-128 encryption protocol.
	SNMPEncryptionProtocolsCFB128_AES128 SNMPEncryptionProtocols = "CFB128_AES128"
	
	// there is no encryption.
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

	// Security/authentication messages.
	SyslogFacilityAuth SyslogFacility = "Auth" 

	// Security/authentication messages.
	SyslogFacilityAuthpriv SyslogFacility = "Authpriv" 

	//Log alert.
	SyslogFacilityConsole SyslogFacility = "Console"

	//Clock daemon. 
	SyslogFacilityCron SyslogFacility = "Cron"
 
	//System daemons.
	SyslogFacilityDaemon SyslogFacility = "Daemon"

	//FTP daemon.
	SyslogFacilityFTP SyslogFacility = "FTP" 
 
	//Kernel messages.        
	SyslogFacilityKern SyslogFacility = "Kern" 
                
	//Line printer subsystem.
	SyslogFacilityLPR SyslogFacility = "LPR" 

	//Locally used facility 0.                
	SyslogFacilityLocal0 SyslogFacility = "Local0" 

	//Locally used facility 1.
	SyslogFacilityLocal1 SyslogFacility = "Local1"

	//Locally used facility 2.
	SyslogFacilityLocal2 SyslogFacility = "Local2" 

	//Locally used facility 3.
	SyslogFacilityLocal3 SyslogFacility = "Local3"

	//Locally used facility 4.
	SyslogFacilityLocal4 SyslogFacility = "Local4" 

	//Locally used facility 5.
	SyslogFacilityLocal5 SyslogFacility = "Local5" 

	//Locally used facility 6.
	SyslogFacilityLocal6 SyslogFacility = "Local6" 

	//Locally used facility 7.
	SyslogFacilityLocal7 SyslogFacility = "Local7"

	//Mail system.
	SyslogFacilityMail SyslogFacility = "Mail"

	//NTP subsystem.
	SyslogFacilityNTP SyslogFacility = "NTP"

	//Network news subsystem.
	SyslogFacilityNews SyslogFacility = "News"

	//Log audit.
	SyslogFacilitySecurity SyslogFacility = "Security" 

	//Scheduling daemon.
	SyslogFacilitySolarisCron SyslogFacility = "SolarisCron" 

	//Messages generated internally by syslogd.
	SyslogFacilitySyslog SyslogFacility = "Syslog" 

	//UUCP subsystem.
	SyslogFacilityUUCP SyslogFacility = "UUCP"

	//User-level messages.
	SyslogFacilityUser SyslogFacility = "User"

	// A condition that should be corrected immediately, such as a corrupted system database.
	SyslogSeverityAlert SyslogSeverity = "Alert"

	// A message of any severity.
	SyslogSeverityAll SyslogSeverity = "All"

	// Hard device errors.
	SyslogSeverityCritical SyslogSeverity = "Critical"

	// Messages that contain information normally of use only when debugging a program.
	SyslogSeverityDebug SyslogSeverity = "Debug"

	// A panic condition.
	SyslogSeverityEmergency SyslogSeverity = "Emergency"

	// An Error.
	SyslogSeverityError SyslogSeverity = "Error"

	// Informational only.
	SyslogSeverityInformational SyslogSeverity = "Informational"

	// Conditions that are not error conditions, but that may require special handling.
	SyslogSeverityNotice SyslogSeverity = "Notice"

	// A Warning.
	SyslogSeverityWarning SyslogSeverity = "Warning"
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
	ODataContext                 string           `json:"@odata.context,omitempty"`
	ODataEtag                    string           `json:"@odata.etag,omitempty"`
	ODataId                      string           `json:"@odata.id"`
	ODataType                    string           `json:"@odata.type"`
	Actions                      *Actions         `json:"Actions,omitempty"`
	Certificates                 *Link            `json:"Certificates,omitempty"`
	ClientCertificates           *Link            `json:"ClientCertificates,omitempty"`
	Context                      string           `json:"Context"`
	DeliveryRetryPolicy          string           `json:"DeliveryRetryPolicy,omitempty"`
	Description                  string           `json:"Description,omitempty"`
	Destination                  string           `json:"Destination"`
	EventFormatType              string           `json:"EventFormatType,omitempty"`
	EventTypes                   []string         `json:"EventTypes,omitempty"`
	ExcludeMessageIds            []string         `json:"ExcludeMessageIds,omitempty"`
	ExcludeRegistryPrefixes      []string         `json:"ExcludeRegistryPrefixes,omitempty"`
	HeartbeatIntervalMins        int              `json:"HeartbeatIntervalMinutes,omitempty"`
	HttpHeaders                  []string         `json:"HttpHeaders,omitempty"`
	ID                           string           `json:"Id"`
	IncludeOriginOfCondition     bool             `json:"IncludeOriginOfCondition,omitempty"`
	MessageIds                   []string         `json:"MessageIds,omitempty"`
	MetricReportDefinitions      *Link            `json:"MetricReportDefinitions,omitempty"`
	MetricReportDefinitionsCount int              `json:MetricReportDefinitions@odata.count,omitempty`
	Name                         string           `json:"Name"`
	OEMProtocol                  string           `json:"OEMProtocol,omitempty"`
	OEMSubscriptionType          string           `json:"OEMSubscriptionType"`
	Oem                          interface{}      `json:"Oem,omitempty"`
	OriginResources              []string         `json:"OriginResources,omitempty"`
	OriginResourcesCount         int              `json:MetricReportDefinitions@odata.count,omitempty`
	Protocol                     string           `json:"Protocol"`
	RegistryPrefixes             []string         `json:"RegistryPrefixes,omitempty"`
	ResourceTypes                []string         `json:"ResourceTypes,omitempty"`
	SNMP                         SNMPSettings     `json:"SNMP,omitempty"`
	SendHeartbeat                bool             `json:"SendHeartbeat,omitempty"`
	Status                       Status           `json:"Status,omitempty"`
	SubordinateResources         bool             `json:"SubordinateResources,omitempty"`
	SubscriptionType             SubscriptionType `json:"SubscriptionType,omitempty"`
	SyslogFilters                SyslogFilter     `json:"SyslogFilters,omitempty"`
	VerifyCertificate            bool             `json:"VerifyCertificate,omitempty"`
}

// place holders for attribs needed in EventDestination.
type Actions struct{}


// This type shall contain the settings for an SNMP event destination.
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

// This action shall suspend an event subscription.  No events shall be 
// sent to the event destination until invocation of the ResumeSubscription 
// action.  The value of the State property within Status shall contain 
// `Disabled` for a suspended subscription.  The service may buffer 
// events while the subscription is suspended.
// Reference	                : EventDestination.v1_12_0.json
type SuspendSubscription struct {
	target string `json:"target,omitempty"`
	title  string `json:"target,omitempty"`
}

// A list of filters applied to syslog messages before sending
// to a remote syslog server.  An empty list indicates all
// syslog messages are sent.
// Reference	                : EventDestination.v1_12_0.json
type SyslogFilter struct {
	LogFacilities  SyslogFacility `json:"LogFacilities,omitEmpty"`
	LowestSeverity SyslogSeverity `json:"LogFacilities,omitEmpty"`
}






