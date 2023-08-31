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

// The type of SFP device that is attached to this port.
type Type string

const (

	// The SFP conforms to the SFF Specification for SFP.
	TypeSFP Type = "SFP"

	// The SFP conforms to the SFF Specification for SFP+.
	TypeSFPPlus Type = "SFPPlus"

	// The SFP conforms to the SFF Specification for SFP+ and IEEE 802.3by Specification.
	TypeSFP28 Type = "SFP28"

	// The SFP conforms to the CSFP MSA Specification.
	TypecSFP Type = "cSFP"

	// The SFP conforms to the SFP-DD MSA Specification.
	TypeSFPDD Type = "SFPDD"

	// The SFP conforms to the SFF Specification for QSFP.
	TypeQSFP Type = "QSFP"

	// The SFP conforms to the SFF Specification for QSFP+.
	TypeQSFPPlus Type = "QSFPPlus"

	// The SFP conforms to the SFF Specification for QSFP14.
	TypeQSFP14 Type = "QSFP14"

	// The SFP conforms to the SFF Specification for QSFP28.
	TypeQSFP28 Type = "QSFP28"

	// The SFP conforms to the SFF Specification for QSFP56.
	TypeQSFP56 Type = "QSFP56"

	// The SFP conforms to the SFF Specification SFF-8644.
	TypeMiniSASHD Type = "MiniSASHD"

	// The SFP conforms to the QSFP Double Density Specification.
	TypeQSFPDD Type = "QSFPDD"

	// The SFP conforms to the OSFP Specification.
	TypeOSFP Type = "OSFP"
)

// Port is the redfish Port model according to the 2020.3 release
type Port struct {
	ODataContext            string               `json:"@odata.context,omitempty"`
	ODataEtag               string               `json:"@odata.etag,omitempty"`
	ODataID                 string               `json:"@odata.id"`
	ODataType               string               `json:"@odata.type"`
	Actions                 *OemActions          `json:"Actions,omitempty"`
	Description             string               `json:"Description,omitempty"`
	ID                      string               `json:"Id"`
	Links                   *PortLinks           `json:"Links,omitempty"`
	Name                    string               `json:"Name"`
	Oem                     interface{}          `json:"Oem,omitempty"`
	Status                  *Status              `json:"Status,omitempty"`
	ActiveWidth             int                  `json:"ActiveWidth,omitempty"`
	CurrentSpeedGbps        float64              `json:"CurrentSpeedGbps,omitempty"`
	Ethernet                *PortEthernet        `json:"Ethernet,omitempty"`
	FibreChannel            *FibreChannel        `json:"FibreChannel,omitempty"`
	GenZ                    *GenZ                `json:"GenZ,omitempty"`
	InterfaceEnabled        bool                 `json:"InterfaceEnabled,omitempty"`
	LinkConfiguration       []*LinkConfiguration `json:"LinkConfiguration,omitempty"`
	LinkNetworkTechnology   string               `json:"LinkNetworkTechnology,omitempty"`
	LinkState               string               `json:"LinkState,omitempty"`
	LinkStatus              string               `json:"LinkStatus,omitempty"`
	LinkTransitionIndicator string               `json:"LinkTransitionIndicator,omitempty"`
	LocationIndicatorActive bool                 `json:"LocationIndicatorActive,omitempty"`
	MaxFrameSize            int                  `json:"MaxFrameSize,omitempty"`
	MaxSpeedGbps            float64              `json:"MaxSpeedGbps,omitempty"`
	Metrics                 *Link                `json:"Metrics,omitempty"`
	PortMedium              string               `json:"PortMedium,omitempty"`
	PortProtocol            string               `json:"PortProtocol,omitempty"`
	PortType                string               `json:"PortType,omitempty"`
	PortID                  string               `json:"PortId,omitempty"`
	SignalDetected          bool                 `json:"SignalDetected,omitempty"`
	Width                   int                  `json:"Width,omitempty"`
	CapableProtocolVersions []string             `json:"CapableProtocolVersions,omitempty"`
	CurrentProtocolVersion  string               `json:"CurrentProtocolVersion,omitempty"`
	Enabled                 bool                 `json:"Enabled,omitempty"`
	EnvironmentMetrics      *Link                `json:"EnvironmentMetrics,omitempty"`
	FunctionMaxBandwidth    []*FunctionBandwidth `json:"FunctionMaxBandwidth,omitempty"`
	FunctionMinBandwidth    []*FunctionBandwidth `json:"FunctionMinBandwidth,omitempty"`
	Location                interface{}          `json:"Location,omitempty"`
	SFP                     SFP                  `json:"SFP,omitempty"`
	CXL                     *CXL                 `json:"CXL,omitempty"`
	InfiniBand              *InfiniBand          `json:"InfiniBand,omitempty"`
	RemotePortId            string               `json:"RemotePortId,omitempty"`
}

// InfiniBand redfish structure
type InfiniBand struct {
	AssociatedNodeGUIDs   []string `json:"AssociatedNodeGUIDs,omitempty"`
	AssociatedPortGUIDs   []string `json:"AssociatedPortGUIDs,omitempty"`
	AssociatedSystemGUIDs []string `json:"AssociatedSystemGUIDs,omitempty"`
}

// CXL redfish structure
type CXL struct {
	Congestion                          *Congestion               `json:"Congestion,omitempty"`
	ConnectedDeviceMode                 string                    `json:"ConnectedDeviceMode,omitempty"`
	ConnectedDeviceType                 string                    `json:"ConnectedDeviceType,omitempty"`
	CurrentPortConfigurationState       string                    `json:"CurrentPortConfigurationState,omitempty"`
	MaxLogicalDeviceCount               string                    `json:"MaxLogicalDeviceCount,omitempty"`
	QoSTelemetryCapabilities            *QoSTelemetryCapabilities `json:"QoSTelemetryCapabilities,omitempty"`
	TemporaryThroughputReductionEnabled bool                      `json:"TemporaryThroughputReductionEnabled,omitempty"`
	AlertCapabilities                   *AlertCapabilities        `json:"AlertCapabilities,omitempty"`
}

// AlertCapabilities redfish structure
type AlertCapabilities struct {
	CorrectableECCError   bool `json:"CorrectableECCError,omitempty"`
	SpareBlock            bool `json:"SpareBlock,omitempty"`
	Temperature           bool `json:"Temperature,omitempty"`
	UncorrectableECCError bool `json:"UncorrectableECCError,omitempty"`
}

// Congestion redfish structure
type Congestion struct {
	BackpressureSampleInterval   int  `json:"BackpressureSampleInterval,omitempty"`
	CompletionCollectionInterval int  `json:"CompletionCollectionInterval,omitempty"`
	CongestionTelemetryEnabled   bool `json:"CongestionTelemetryEnabled,omitempty"`
	EgressModeratePercentage     int  `json:"EgressModeratePercentage,omitempty"`
	EgressSeverePercentage       int  `json:"EgressSeverePercentage,omitempty"`
	MaxSustainedRequestCmpBias   int  `json:"MaxSustainedRequestCmpBias,omitempty"`
}

// QoSTelemetryCapabilities redfish structure
type QoSTelemetryCapabilities struct {
	EgressPortBackpressureSupported       bool `json:"EgressPortBackpressureSupported,omitempty"`
	TemporaryThroughputReductionSupported bool `json:"TemporaryThroughputReductionSupported,omitempty"`
}

// SFP redfish structure
type SFP struct {
	FiberConnectionType string   `json:"FiberConnectionType,omitempty"`
	Manufacturer        string   `json:"Manufacturer,omitempty"`
	MediumType          string   `json:"MediumType,omitempty"`
	PartNumber          string   `json:"PartNumber,omitempty"`
	SerialNumber        string   `json:"SerialNumber,omitempty"`
	Status              Status   `json:"Status,omitempty"`
	SupportedSFPTypes   []string `json:"SupportedSFPTypes,omitempty"`
	Type                string   `json:"Type,omitempty"`
}

// FunctionBandwidth redfish structure
type FunctionBandwidth struct {
	AllocationPercent     int   `json:"AllocationPercent,omitempty"`
	NetworkDeviceFunction *Link `json:"NetworkDeviceFunctions,omitempty"`
}

// PortEthernet redfish model
type PortEthernet struct {
	FlowControlConfiguration      string   `json:"FlowControlConfiguration,omitempty"`
	FlowControlStatus             string   `json:"FlowControlStatus,omitempty"`
	SupportedEthernetCapabilities string   `json:"SupportedEthernetCapabilities,omitempty"`
	WakeOnLANEnabled              bool     `json:"WakeOnLANEnabled,omitempty"`
	EEEEnabled                    bool     `json:"EEEEnabled,omitempty"`
	LLDPTransmit                  *LLDP    `json:"LLDPTransmit,omitempty"`
	LLDPReceive                   *LLDP    `json:"LLDPReceive,omitempty"`
	LLDPEnabled                   bool     `json:"LLDPEnabled,omitempty"`
	AssociatedMACAddresses        []string `json:"AssociatedMACAddresses,omitempty"`
}

// LLDP redfish structure
type LLDP struct {
	ChassisID             string   `json:"ChassisId,omitempty"`
	ChassisIDSubtype      string   `json:"ChassisIdSubtype,omitempty"`
	ManagementAddressIPv4 string   `json:"ManagementAddressIPv4,omitempty"`
	ManagementAddressIPv6 string   `json:"ManagementAddressIPv6,omitempty"`
	ManagementAddressMAC  string   `json:"ManagementAddressMAC,omitempty"`
	ManagementVlanID      int      `json:"ManagementVlanId,omitempty"`
	PortID                string   `json:"PortId,omitempty"`
	PortIDSubtype         string   `json:"PortIdSubtype,omitempty"`
	SystemCapabilities    []string `json:"SystemCapabilities,omitempty"`
	SystemDescription     string   `json:"SystemDescription,omitempty"`
	SystemName            string   `json:"SystemName,omitempty"`
}

// FibreChannel redfish model
type FibreChannel struct {
	FabricName                  string   `json:"FabricName,omitempty"`
	NumberDiscoveredRemotePorts int      `json:"NumberDiscoveredRemotePorts,omitempty"`
	PortConnectionType          string   `json:"PortConnectionType,omitempty"`
	AssociatedWorldWideNames    []string `json:"AssociatedWorldWideNames,omitempty"`
}

// LinkConfiguration redfish model
type LinkConfiguration struct {
	AutoSpeedNegotiationCapable bool                   `json:"AutoSpeedNegotiationCapable,omitempty"`
	AutoSpeedNegotiationEnabled bool                   `json:"AutoSpeedNegotiationEnabled,omitempty"`
	CapableLinkSpeedGbps        []int                  `json:"CapableLinkSpeedGbps,omitempty"`
	ConfiguredNetworkLinks      []CapableLinkSpeedGbps `json:"ConfiguredNetworkLinks,omitempty"`
}

// CapableLinkSpeedGbps redfish model
type CapableLinkSpeedGbps struct {
	ConfiguredLinkSpeedGbps float64 `json:"ConfiguredLinkSpeedGbps,omitempty"`
	ConfiguredWidth         int     `json:"ConfiguredWidth,omitempty"`
}

// PortLinks Port link redfish model
type PortLinks struct {
	AssociatedEndpoints       []Link      `json:"AssociatedEndpoints,omitempty"`
	ConnectedPorts            []Link      `json:"ConnectedPorts,omitempty"`
	ConnectedSwitches         []Link      `json:"ConnectedSwitches,omitempty"`
	ConnectedSwitchPorts      []Link      `json:"ConnectedSwitchPorts,omitempty"`
	Oem                       interface{} `json:"Oem,omitempty"`
	Cables                    []Link      `json:"Cables,omitempty"`
	AssociatedEndpointsCount  int         `json:"AssociatedEndpoints@odata.count,omitempty"`
	CablesCount               int         `json:"Cables@odata.count,omitempty"`
	ConnectedPortsCount       int         `json:"ConnectedPorts@odata.count,omitempty"`
	ConnectedSwitchPortsCount int         `json:"ConnectedSwitchPorts@odata.count,omitempty"`
	ConnectedSwitchesCount    int         `json:"ConnectedSwitches@odata.count,omitempty"`
	EthernetInterfaces        []Link      `json:"EthernetInterfaces,omitempty"`
}
