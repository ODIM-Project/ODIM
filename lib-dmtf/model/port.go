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

//Port is the redfish Port model according to the 2020.3 release
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
	InterfaceEnabled        bool                 `json:"InterfaceEnabled"`
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
	Location                *Link                `json:"Location,omitempty"`
	SFP                     SFP                  `json:"SFP,omitempty"`
}

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

type FunctionBandwidth struct {
	AllocationPercent     int   `json:"AllocationPercent,omitempty"`
	NetworkDeviceFunction *Link `json:"NetworkDeviceFunctions,omitempty"`
}

//PortEthernet redfish model
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

type LLDP struct {
	ChassisId             string `json:"ChassisId,omitempty"`
	ChassisIdSubtype      string `json:"ChassisIdSubtype,omitempty"`
	ManagementAddressIPv4 string `json:"ManagementAddressIPv4,omitempty"`
	ManagementAddressIPv6 string `json:"ManagementAddressIPv6,omitempty"`
	ManagementAddressMAC  string `json:"ManagementAddressMAC,omitempty"`
	ManagementVlanId      int    `json:"ManagementVlanId,omitempty"`
	PortId                string `json:"PortId,omitempty"`
	PortIdSubtype         string `json:"PortIdSubtype,omitempty"`
}

//FibreChannel redfish model
type FibreChannel struct {
	FabricName                  string   `json:"FabricName,omitempty"`
	NumberDiscoveredRemotePorts int      `json:"NumberDiscoveredRemotePorts,omitempty"`
	PortConnectionType          string   `json:"PortConnectionType,omitempty"`
	AssociatedWorldWideNames    []string `json:"AssociatedWorldWideNames,omitempty"`
}

//LinkConfiguration redfish model
type LinkConfiguration struct {
	AutoSpeedNegotiationCapable bool                   `json:"AutoSpeedNegotiationCapable,omitempty"`
	AutoSpeedNegotiationEnabled bool                   `json:"AutoSpeedNegotiationEnabled,omitempty"`
	CapableLinkSpeedGbps        []CapableLinkSpeedGbps `json:"CapableLinkSpeedGbps,omitempty"`
	ConfiguredNetworkLinks      []CapableLinkSpeedGbps `json:"ConfiguredNetworkLinks,omitempty"`
}

//CapableLinkSpeedGbps redfish model
type CapableLinkSpeedGbps struct {
	ConfiguredLinkSpeedGbps []float64 `json:"ConfiguredLinkSpeedGbps,omitempty"`
	ConfiguredWidth         int       `json:"ConfiguredWidth,omitempty"`
}

//PortLinks Port link redfish model
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
}
