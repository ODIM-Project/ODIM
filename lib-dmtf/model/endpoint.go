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

// Endpoint is the redfish Endpoint model according to the 2020.3 release
type Endpoint struct {
	ODataContext               string               `json:"@odata.context,omitempty"`
	ODataEtag                  string               `json:"@odata.etag,omitempty"`
	ODataID                    string               `json:"@odata.id"`
	ODataType                  string               `json:"@odata.type"`
	Actions                    *OemActions          `json:"Actions,omitempty"`
	Description                string               `json:"Description,omitempty"`
	ID                         string               `json:"Id"`
	Links                      *EndpointLinks       `json:"Links,omitempty"`
	Name                       string               `json:"Name"`
	Oem                        interface{}          `json:"Oem,omitempty"`
	Status                     *Status              `json:"Status,omitempty"`
	ConnectedEntites           []*ConnectedEntites  `json:"ConnectedEntites,omitempty"`
	EndpointProtocol           string               `json:"EndpointProtocol,omitempty"`
	HostReservationMemoryBytes int                  `json:"HostReservationMemoryBytes,omitempty"`
	Identifiers                []EndpointIdentifier `json:"Identifiers,omitempty"`
	IPTransportDetails         []IPTransportDetails `json:"IPTransportDetails,omitempty"`
	Ports                      []Link               `json:"Ports,omitempty"`
	PciID                      *PciID               `json:"PciId,omitempty"`
	Redundancy                 []Redundancy         `json:"Redundancy,omitempty"`
}

// ConnectedEntites for Endpoint
type ConnectedEntites struct {
	EntityPciID       *EntityPciID `json:"EntityPciID,omitempty"`
	EntityRole        string       `json:"EntityRole,omitempty"`
	EntityType        string       `json:"EntityType,omitempty"`
	GenZ              GenZ         `json:"GenZ,omitempty"`
	Identifier        *Identifier  `json:"Identifier,omitempty"`
	Oem               interface{}  `json:"Oem,omitempty"`
	PciClassCode      string       `json:"PciClassCode,omitempty"`
	PciFunctionNumber string       `json:"PciFunctionNumber,omitempty"`
}

// EntityPciID for Endpoint
type EntityPciID struct {
	ClassCode         string `json:"ClassCode,omitempty"`
	DeviceID          string `json:"DeviceId,omitempty"`
	FunctionNumber    int    `json:"FunctionNumber,omitempty"`
	SubsystemID       string `json:"SubsystemId,omitempty"`
	SubsystemVendorID string `json:"SubsystemVendorId,omitempty"`
	VendorID          string `json:"VendorId,omitempty"`
}

// EndpointIdentifier for Endpoint
type EndpointIdentifier struct {
	DurableName       string `json:"DurableName,omitempty"`
	DurableNameFormat string `json:"DurableNameFormat,omitempty"`
}

// IPTransportDetails for Endpoint
type IPTransportDetails struct {
	IPv4Address       *IPv4Address `json:"IPv4Address,omitempty"`
	IPv6Address       *IPv6Address `json:"IPv6Address,omitempty"`
	Port              int          `json:"Port,omitempty"`
	TransportProtocol string       `json:"TransportProtocol,omitempty"`
}

// PciID for Endpoint
type PciID struct {
	ClassCode         string `json:"ClassCode,omitempty"`
	DeviceID          string `json:"DeviceId,omitempty"`
	FunctionNumber    int    `json:"FunctionNumber,omitempty"`
	SubsystemID       string `json:"SubsystemId,omitempty"`
	SubsystemVendorID string `json:"SubsystemVendorId,omitempty"`
	VendorID          string `json:"VendorId,omitempty"`
}

// EndpointLinks is the struct to links under a AddressPool
type EndpointLinks struct {
	AddressPools                    []Link      `json:"AddressPools,omitempty"`
	AddressPoolsCount               int         `json:"AddressPools@odata.count,omitempty"`
	ConnectedPorts                  []Link      `json:"ConnectedPorts,omitempty"`
	ConnectedPortsCount             int         `json:"ConnectedPorts@odata.count,omitempty"`
	Connections                     []Link      `json:"Connections,omitempty"`
	ConnectionsCount                int         `json:"Connections@odata.count,omitempty"`
	MutuallyExclusiveEndpoints      []Link      `json:"MutuallyExclusiveEndpoints,omitempty"`
	MutuallyExclusiveEndpointsCount int         `json:"MutuallyExclusiveEndpoints@odata.count,omitempty"`
	NetworkDeviceFunction           []Link      `json:"NetworkDeviceFunction,omitempty"`
	NetworkDeviceFunctionCount      int         `json:"NetworkDeviceFunction@odata.count,omitempty"`
	Ports                           []Link      `json:"Ports,omitempty"`
	PortsCount                      int         `json:"Ports@odata.count,omitempty"`
	Oem                             interface{} `json:"Oem,omitempty"`
}
