// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package model

// ManagerType The type of manager that this resource represents.
type ManagerType string

// PowerState This property shall contain the power state of the manager
type PowerState string

// ConnectTypesSupported This property enumerates the graphical console
// connection types that the implementation allows.
type ConnectTypesSupported string

// AuthenticationModes This property shall contain an array consisting of the
// authentication modes allowed on this interface
type AuthenticationModes string

// HostInterfaceType This property shall contain an enumeration that describes
// the type of the interface
type HostInterfaceType string

const (

	// ManagerTypeManagementController - "ManagementController": "A controller that primarily monitors or manages the operation of a device or system."
	ManagerTypeManagementController ManagerType = "ManagementController"

	// ManagerTypeEnclosureManager - "EnclosureManager": "A controller that provides management functions for a chassis or group of devices or systems."
	ManagerTypeEnclosureManager ManagerType = "EnclosureManager"

	// ManagerTypeBMC - "BMC": "A controller that provides management functions for a single computer system."
	ManagerTypeBMC ManagerType = "BMC"

	// ManagerTypeRackManager - "RackManager": "A controller that provides management functions for a whole or part of a rack."
	ManagerTypeRackManager ManagerType = "RackManager"

	//ManagerTypeAuxiliaryController - "AuxiliaryController": "A controller that provides management functions for a particular subsystem or group of devices."
	ManagerTypeAuxiliaryController ManagerType = "AuxiliaryController"

	// ManagerTypeService - "Service": "A software-based service that provides management functions."
	ManagerTypeService ManagerType = "Service"
)

const (
	// PowerStateOn - "On": "The resource is powered on."
	PowerStateOn PowerState = "On"

	// PowerStateOff - "Off": "The resource is powered off.  The components within the resource might continue to have AUX power."
	PowerStateOff PowerState = "Off"

	// PowerStatePoweringOn - "PoweringOn": "A temporary state between off and on.  The components within the resource can take time to process the power on action."
	PowerStatePoweringOn PowerState = "PoweringOn"

	// PowerStatePoweringOff - "PoweringOff": "A temporary state between on and off.  The components within the resource can take time to process the power off action."
	PowerStatePoweringOff PowerState = "PoweringOff"

	// PowerStatePaused - "Paused": "The resource is paused."
	PowerStatePaused PowerState = "Paused"
)

const (
	// ConnectTypesSupportedKVMIP - "KVMIP": "The controller supports a graphical console connection through a KVM-IP (redirection of Keyboard, Video, Mouse over IP) protocol."
	ConnectTypesSupportedKVMIP ConnectTypesSupported = "KVMIP"

	// ConnectTypesSupportedOem - "Oem": "The controller supports a graphical console connection through an OEM-specific protocol."
	ConnectTypesSupportedOem ConnectTypesSupported = "Oem"
)

const (
	// HostInterfaceTypeNetworkHostInterface - "NetworkHostInterface": "This interface is a Network Host Interface."
	HostInterfaceTypeNetworkHostInterface HostInterfaceType = "NetworkHostInterface"
)

const (
	// AuthenticationModesAuthNone - "AuthNone": "Requests without any sort of authentication are allowed."
	AuthenticationModesAuthNone AuthenticationModes = "AuthNone"

	// AuthenticationModesBasicAuth - "BasicAuth": "Requests using HTTP Basic Authentication are allowed."
	AuthenticationModesBasicAuth AuthenticationModes = "BasicAuth"

	// AuthenticationModesRedfishSessionAuth - "RedfishSessionAuth": "Requests using Redfish Session Authentication are allowed."
	AuthenticationModesRedfishSessionAuth AuthenticationModes = "RedfishSessionAuth"

	// AuthenticationModesOemAuth - "OemAuth": "Requests using OEM authentication mechanisms are allowed."
	AuthenticationModesOemAuth AuthenticationModes = "OemAuth"
)

// Manager is the redfish Manager model according to the 2020.3 release
// Refer to Manager.v1_17_0.json of the redfish spec for more details
type Manager struct {
	ODataContext               string                      `json:"@odata.context,omitempty"`
	ODataEtag                  string                      `json:"@odata.etag,omitempty"`
	ODataID                    string                      `json:"@odata.id"`
	ODataType                  string                      `json:"@odata.type"`
	Actions                    *ManagerActions             `json:"Actions,omitempty"`
	Description                string                      `json:"Description,omitempty"`
	ID                         string                      `json:"Id"`
	Links                      *ManagerLinks               `json:"Links,omitempty"`
	Name                       string                      `json:"Name"`
	Oem                        *Oem                        `json:"Oem,omitempty"`
	EthernetInterfaces         *Link                       `json:"EthernetInterfaces,omitempty"`
	FirmwareVersion            string                      `json:"FirmwareVersion,omitempty"`
	Status                     *ManagerStatus              `json:"Status,omitempty"`
	AutoDSTEnabled             bool                        `json:"AutoDSTEnabled,omitempty"`
	DateTime                   string                      `json:"DateTime,omitempty"`
	DateTimeLocalOffset        string                      `json:"DateTimeLocalOffset,omitempty"`
	HostInterfaces             *HostInterface              `json:"HostInterfaces,omitempty"`
	LastResetTime              string                      `json:"LastResetTime,omitempty"`
	LogServices                *Link                       `json:"LogServices,omitempty"`
	ManagerType                string                      `json:"ManagerType,omitempty"` //enum
	Manufacturer               string                      `json:"Manufacturer,omitempty"`
	Model                      string                      `json:"Model,omitempty"`
	NetworkProtocol            *Link                       `json:"NetworkProtocol,omitempty"`
	PartNumber                 string                      `json:"PartNumber,omitempty"`
	PowerState                 string                      `json:"PowerState,omitempty"` //enum
	Redundancy                 []Redundancy                `json:"Redundancy,omitempty"`
	RemoteAccountService       *Link                       `json:"RemoteAccountService,omitempty"`
	RemoteRedfishServiceURI    string                      `json:"RemoteRedfishServiceUri,omitempty"`
	SerialNumber               string                      `json:"SerialNumber,omitempty"`
	ServiceEntryPointUUID      string                      `json:"ServiceEntryPointUUID,omitempty"`
	SerialInterfaces           *Link                       `json:"SerialInterfaces,omitempty"`
	TimeZoneName               string                      `json:"TimeZoneName,omitempty"`
	UUID                       string                      `json:"UUID,omitempty"`
	Measurements               []*Link                     `json:"Measurements,omitempty"` // Deprecated in version v1_14_0
	Certificates               *Certificates               `json:"Certificates,omitempty"`
	CommandShell               *CommandShell               `json:"CommandShell,omitempty"`
	GraphicalConsole           *GraphicalConsole           `json:"GraphicalConsole,omitempty"`
	Location                   *Link                       `json:"Location,omitempty"`
	LocationIndicatorActive    bool                        `json:"LocationIndicatorActive,omitempty"`
	RedundancyCount            int                         `json:"Redundancy@odata.count,omitempty"`
	SparePartNumber            string                      `json:"SparePartNumber,omitempty"`
	ManagerDiagnosticData      *Link                       `json:"ManagerDiagnosticData,omitempty"`
	ServiceIdentification      string                      `json:"ServiceIdentification,omitempty"`
	AdditionalFirmwareVersions *AdditionalFirmwareVersions `json:"AdditionalFirmwareVersions,omitempty"`
}

// AdditionalFirmwareVersions redfish structure
type AdditionalFirmwareVersions struct {
	Bootloader     string `json:"Bootloader,omitempty"`
	Kernel         string `json:"Kernel,omitempty"`
	Microcode      string `json:"Microcode,omitempty"`
	OSDistribution string `json:"OSDistribution,omitempty"`
	Oem            Oem    `json:"Oem,omitempty"`
}

// CommandShell redfish structure
type CommandShell struct {
	ConnectTypesSupported []string `json:"ConnectTypesSupported"`
	MaxConcurrentSessions int      `json:"MaxConcurrentSessions"`
	ServiceEnabled        bool     `json:"ServiceEnabled"`
}

// GraphicalConsole redfish structure
type GraphicalConsole struct {
	ConnectTypesSupported []string `json:"ConnectTypesSupported"` //enum
	MaxConcurrentSessions int      `json:"MaxConcurrentSessions"`
	ServiceEnabled        bool     `json:"ServiceEnabled"`
}

// ManagerLinks ...
type ManagerLinks struct {
	ActiveSoftwareImage     *Link       `json:"ActiveSoftwareImage,omitempty"`
	ManagedBy               []Link      `json:"ManagedBy,omitempty"`
	ManagerByCount          int         `json:"ManagedBy@odata.count,omitempty"`
	ManagerForChassis       []Link      `json:"ManagerForChassis,omitempty"`
	ManagerForChassisCount  int         `json:"ManagerForChassis@odata.count,omitempty"`
	ManagerForManagers      []Link      `json:"ManagerForManagers,omitempty"`
	ManagerForManagersCount int         `json:"ManagerForManagers@odata.count,omitempty"`
	ManagerForServers       []Link      `json:"ManagerForServers,omitempty"`
	ManagerForServersCount  int         `json:"ManagerForServers@odata.count,omitempty"`
	ManagerForSwitches      []Link      `json:"ManagerForSwitches,omitempty"`
	ManagerForSwitchesCount int         `json:"ManagerForSwitches@odata.count,omitempty"`
	ManagerInChassis        []Link      `json:"ManagerInChassis,omitempty"`
	ManagerInChassisCount   int         `json:"ManagerInChassis@odata.count,omitempty"`
	Oem                     interface{} `json:"Oem,omitempty"`
	SoftwareImages          *Link       `json:"SoftwareImages,omitempty"`
}

// VirtualMedia is a redfish virtual media model
type VirtualMedia struct {
	ODataContext         string              `json:"@odata.context,omitempty"`
	ODataEtag            string              `json:"@odata.etag,omitempty"`
	ODataID              string              `json:"@odata.id"`
	ODataType            string              `json:"@odata.type"`
	Actions              VMActions           `json:"Actions,omitempty"`
	ConnectedVia         string              `json:"ConnectedVia,omitempty"`
	Description          string              `json:"Description,omitempty"`
	ID                   string              `json:"Id"`
	Image                string              `json:"Image"`
	ImageName            string              `json:"ImageName,omitempty"`
	Inserted             bool                `json:"Inserted"`
	MediaTypes           []string            `json:"MediaTypes,omitempty"`
	Name                 string              `json:"Name"`
	Oem                  interface{}         `json:"Oem,omitempty"`
	Password             string              `json:"Password,omitempty"`
	TransferMethod       string              `json:"TransferMethod,omitempty"`
	TransferProtocolType string              `json:"TransferProtocolType,omitempty"`
	UserName             string              `json:"UserName,omitempty"`
	VerifyCertificate    bool                `json:"VerifyCertificate,omitempty"`
	WriteProtected       bool                `json:"WriteProtected,omitempty"`
	Status               *Status             `json:"Status,omitempty"`
	ClientCertificates   *ClientCertificates `json:"ClientCertificates,omitempty"`
	Certificates         Certificates        `json:"Certificates,omitempty"`
}

// ClientCertificates redfish structure
type ClientCertificates struct {
	OdataID string `json:"@odata.id"`
}

// VMActions contains the actions property details of virtual media
type VMActions struct {
	EjectMedia  ActionTarget `json:"EjectMedia"`
	InsertMedia ActionTarget `json:"InsertMedia"`
}

// ActionTarget contains the action target
type ActionTarget struct {
	Target string `json:"target"`
}

// ManagerActions property shall contain the available actions for this resource.
type ManagerActions struct {
	ManagerForceFailover       *ForceFailover       `json:"ManagerForceFailover,omitempty"`
	ManagerModifyRedundancySet *ModifyRedundancySet `json:"ManagerModifyRedundancy,omitempty"`
	ManagerReset               *Reset               `json:"ManagerReset,omitempty"`
	ManagerResetToDefaults     *ResetToDefaults     `json:"ManagerResetToDefaults,omitempty"`
	Oem                        *Oem                 `json:"Oem,omitempty"`
}

// ForceFailover action forces a failover of this manager to the manager
// used in the parameter
type ForceFailover struct {
	Target string `json:"Target,omitempty"`
	Title  string `json:"Title,omitempty"`
}

// ModifyRedundancySet operation shall add members to or remove members from
// a redundant group of managers
type ModifyRedundancySet struct {
	Target string `json:"Target,omitempty"`
	Title  string `json:"Title,omitempty"`
}

// ResetToDefaults resets the manager settings to factory defaults.
// This can cause the manager to reset
type ResetToDefaults struct {
	Target string `json:"Target,omitempty"`
	Title  string `json:"Title,omitempty"`
}

// HostInterfaces link to a collection of host interfaces that this manager uses for local host communication.
// Clients can find host interface configuration options and settings in this navigation property
type HostInterfaces struct {
	ODataContext         string   `json:"@odata.context,omitempty"`
	ODataEtag            string   `json:"@odata.etag,omitempty"`
	ODataID              string   `json:"@odata.id"`
	ODataType            string   `json:"@odata.type"`
	Descriptions         string   `json:"Descriptions,omitempty"`
	Members              *Members `json:"Members"`
	MembersOdataCount    int      `json:"MembersOdataCount"`
	MembersOdataNextLink string   `json:"MembersOdataNextLink,omitempty"`
	Name                 string   `json:"Name"`
	Oem                  *Oem     `json:"Oem,omitempty"`
}

// Members property shall contain an array of links to the members of this collection
type Members struct {
	ODataContext            string                   `json:"@odata.context,omitempty"`
	ODataEtag               string                   `json:"@odata.etag,omitempty"`
	ODataID                 string                   `json:"@odata.id"`
	ODataType               string                   `json:"@odata.type"`
	Actions                 *OemActions              `json:"actions,omitempty"`
	AuthNoneRoleID          string                   `json:"authNoneRoleId,omitempty"`
	AuthenticationModes     []string                 `json:"authenticationModes,omitempty"` //enum
	CredentialBootstrapping *CredentialBootstrapping `json:"credentialBootstrapping,omitempty"`
	Description             string                   `json:"description,omitempty"`
	ExternallyAccessible    bool                     `json:"externallyAccessible,omitempty"`
	FirmwareAuthEnabled     bool                     `json:"firmwareAuthEnabled,omitempty"`
	FirmwareAuthRoleID      string                   `json:"firmwareAuthRoleId,omitempty"`
	HostEthernetInterfaces  *HostInterfaces          `json:"hostEthernetInterfaces,omitempty"`
	HostInterfaceType       string                   `json:"hostInterfaceType,omitempty"` //enum
	ID                      string                   `json:"id,omitempty"`
	InterfaceEnabled        bool                     `json:"interfaceEnabled,omitempty"`
	KernelAuthEnabled       bool                     `json:"kernelAuthEnabled,omitempty"`
	KernelAuthRoleID        string                   `json:"kernelAuthRoleId,omitempty"`
	Links                   *Links                   `json:"links,omitempty"`
	Name                    string                   `json:"name,omitempty"`
	Oem                     *Oem                     `json:"Oem,omitempty"`
	Status                  []*Conditions            `json:"status,omitempty"`
}

// CredentialBootstrapping settings for this interface
type CredentialBootstrapping struct {
	EnableAfterReset bool   `json:"enableAfterReset,omitempty"`
	Enabled          bool   `json:"enabled,omitempty"`
	RoleID           string `json:"roleId,omitempty"`
}

// ManagerStatus property shall contain any status or health properties of the resource
type ManagerStatus struct {
	Conditions   []*Conditions `json:"Conditions,omitempty"`
	Name         string        `json:"Name,omitempty"`
	Health       string        `json:"Health,omitempty"`
	HealthRollup string        `json:"HealthRollup,omitempty"`
	State        string        `json:"State,omitempty"`
	Oem          *Oem          `json:"Oem,omitempty"`
}
