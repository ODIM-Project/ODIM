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
}

//PortEthernet redfish model
type PortEthernet struct {
	FlowControlConfiguration      string `json:"FlowControlConfiguration,omitempty"`
	FlowControlStatus             string `json:"FlowControlStatus,omitempty"`
	SupportedEthernetCapabilities string `json:"SupportedEthernetCapabilities,omitempty"`
}

//FibreChannel redfish model
type FibreChannel struct {
	FabricName                  string `json:"FabricName,omitempty"`
	NumberDiscoveredRemotePorts int    `json:"NumberDiscoveredRemotePorts,omitempty"`
	PortConnectionType          string `json:"PortConnectionType,omitempty"`
}

//LinkConfiguration redfish model
type LinkConfiguration struct {
	AutoSpeedNegotiationCapable bool                   `json:"AutoSpeedNegotiationCapable,omitempty"`
	AutoSpeedNegotiationEnabled bool                   `json:"AutoSpeedNegotiationEnabled,omitempty"`
	CapableLinkSpeedGbps        []CapableLinkSpeedGbps `json:"CapableLinkSpeedGbps,omitempty"`
}

//CapableLinkSpeedGbps redfish model
type CapableLinkSpeedGbps struct {
	ConfiguredLinkSpeedGbps []float64 `json:"ConfiguredLinkSpeedGbps,omitempty"`
	ConfiguredWidth         int       `json:"ConfiguredWidth,omitempty"`
}

//PortLinks Port link redfish model
type PortLinks struct {
	AssociatedEndpoints  []Link      `json:"AssociatedEndpoints,omitempty"`
	ConnectedPorts       []Link      `json:"ConnectedPorts,omitempty"`
	ConnectedSwitches    []Link      `json:"ConnectedSwitches,omitempty"`
	ConnectedSwitchPorts []Link      `json:"ConnectedSwitchPorts,omitempty"`
	Oem                  interface{} `json:"Oem,omitempty"`
}
