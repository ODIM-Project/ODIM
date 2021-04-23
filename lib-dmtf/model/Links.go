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

// Links - this is Common to all resources
type Links struct {
	AddressPools             []*Link     `json:"AddressPools,omitempty"`
	Chassis                  []*Link     `json:"Chassis,omitempty"`
	ComputerSystems          []*Link     `json:"ComputerSystems,omitempty"`
	ConnectedPorts           []*Link     `json:"ConnectedPorts,omitempty"`
	ConsumingComputerSystems []*Link     `json:"ConsumingComputerSystems,omitempty"`
	ContainedBy              *Link       `json:"ContainedBy,omitempty"`
	ContainedByZones         []*Link     `json:"ContainedByZones,omitempty"`
	CooledBy                 []*Link     `json:"CooledBy,omitempty"`
	Endpoints                []*Link     `json:"Endpoints,omitempty"`
	EndpointsCount           int         `json:"Endpoints@odata.count,omitempty"`
	Drives                   []*Link     `json:"Drives,omitempty"`
	ManagedBy                []*Link     `json:"ManagedBy,omitempty"`
	Oem                      interface{} `json:"Oem,omitempty"`
	ManagersInChassis        []*Link     `json:"ManagersInChassis,omitempty"`
	PCIeDevices              []*Link     `json:"PCIeDevices,omitempty"`
	PCIeFunctions            []*Link     `json:"PCIeFunctions,omitempty"`
	PoweredBy                []*Link     `json:"PoweredBy,omitempty"`
	Processors               []*Link     `json:"Processors,omitempty"`
	ResourceBlocks           []*Link     `json:"ResourceBlocks,omitempty"`
	Storage                  []*Link     `json:"Storage,omitempty"`
	SupplyingComputerSystems []*Link     `json:"SupplyingComputerSystems,omitempty"`
	Switches                 []*Link     `json:"Switches,omitempty"`
	Zones                    []*Link     `json:"Zones,omitempty"`
	ZonesCount               int         `json:"Zones@odata.count,omitempty"`
}

// Link holds the odata id redfish links
type Link struct {
	Oid string `json:"@odata.id"`
}

// Oem holds the vendor specific details which is addtional to redfish contents
type Oem interface{}
