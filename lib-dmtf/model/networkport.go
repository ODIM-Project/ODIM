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

// NetworkPort is the redfish network port model
type NetworkPort struct {
	ODataContext                  string                    `json:"@odata.context,omitempty"`
	ODataEtag                     string                    `json:"@odata.etag,omitempty"`
	ODataID                       string                    `json:"@odata.id"`
	ODataType                     string                    `json:"@odata.type"`
	ID                            string                    `json:"Id"`
	Name                          string                    `json:"Name"`
	ActiveLinkTechnology          string                    `json:"ActiveLinkTechnology,omitempty"`
	AssociatedNetworkAddresses    []string                  `json:"AssociatedNetworkAddresses",omitempty`
	CurrentLinkSpeedMbps          int                       `json:"CurrentLinkSpeedMbps",omitempty`
	EEEEnabled                    bool                      `json:"EEEEnabled",omitempty`
	FCFabricName                  string                    `json:"FCFabricName",omitempty`
	FCPortConnectionType          string                    `json:"FCPortConnectionType",omitempty`
	FlowControlConfiguration      string                    `json:"FlowControlConfiguration",omitempty`
	FlowControlStatus             string                    `json:"FlowControlStatus",omitempty`
	LinkStatus                    string                    `json:"LinkStatus",omitempty`
	MaxFrameSize                  int                       `json:"MaxFrameSize",omitempty`
	NetDevFuncMaxBWAlloc          []NetDevFuncMaxBWAlloc    `json:"NetDevFuncMaxBWAlloc",omitempty`
	NetDevFuncMinBWAlloc          []NetDevFuncMinBWAlloc    `json:"NetDevFuncMinBWAlloc",omitempty`
	NumberDiscoveredRemotePorts   int                       `json:"NumberDiscoveredRemotePorts",omitempty`
	PhysicalPortNumber            string                    `json:"PhysicalPortNumber",omitempty`
	PortMaximumMTU                int                       `json:"PortMaximumMTU",omitempty`
	SignalDetected                bool                      `json:"SignalDetected",omitempty`
	Status                        *Status                   `json:"Status",omitempty`
	SupportedEthernetCapabilities []string                  `json:"SupportedEthernetCapabilities",omitempty`
	SupportedLinkCapabilities     SupportedLinkCapabilities `json:"SupportedLinkCapabilities",omitempty`
	VendorID                      string                    `json:"VendorId",omitempty`
	WakeOnLANEnabled              bool                      `json:"WakeOnLANEnabled",omitempty`
}

// NetDevFuncMaxBWAlloc contains the information of maximum bandwidth allocation percentages for the
// network device functions associated with the newtwork port.
type NetDevFuncMaxBWAlloc struct {
	MaxBWAllocPercent     int   `json:"MaxBWAllocPercent",omitempty`
	NetworkDeviceFunction *Link `json:"NetworkDeviceFunction",omitempty`
}

// NetDevFuncMinBWAlloc contains the information of minimum bandwidth allocation percentages for the
// network device functions associated with the newtwork port.
type NetDevFuncMinBWAlloc struct {
	MinBWAllocPercent     int   `json:"MinBWAllocPercent",omitempty`
	NetworkDeviceFunction *Link `json:"NetworkDeviceFunction",omitempty`
}

// SupportedLinkCapabilities contains the information of Ethernet capabilities that the network port support
type SupportedLinkCapabilities struct {
	AutoSpeedNegotiation  bool   `json:"AutoSpeedNegotiation",omitempty`
	CapableLinkSpeedMbps  []int  `json:"CapableLinkSpeedMbps",omitempty`
	LinkNetworkTechnology string `json:"LinkNetworkTechnology",omitempty`
	LinkSpeedMbps         int    `json:"LinkSpeedMbps",omitempty`
}
