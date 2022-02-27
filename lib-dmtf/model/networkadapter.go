//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

// NetworkAdapter is the redfish network adapter model
type NetworkAdapter struct {
	ODataContext           string                 `json:"@odata.context,omitempty"`
	ODataEtag              string                 `json:"@odata.etag,omitempty"`
	ODataID                string                 `json:"@odata.id"`
	ODataType              string                 `json:"@odata.type"`
	ID                     string                 `json:"Id"`
	Name                   string                 `json:"Name"`
	Actions                *NetworkAdapterActions `json:"Actions,omitempty"`
	Assembly               *Link                  `json:"Assembly,omitempty"`
	Certificates           *Link                  `json:"Certificates,omitempty"`
	Controllers            []Controllers          `json:"Controllers,omitempty"`
	Description            string                 `json:"Description,omitempty"`
	EnvironmentMetrics     *Link                  `json:"EnvironmentMetrics,omitempty"`
	Identifiers            *Identifier            `json:"Identifiers,omitempty"`
	Location               *Location              `json:"Location,omitempty"`
	LLDPEnabled            bool                   `json:"LLDPEnabled,omitempty"`
	Manufacturer           interface{}            `json:"Manufacturer"`
	Measurements           []*Link                `json:"Measurements,omitempty"`
	Metrics                *Link                  `json:"Metrics,omitempty"`
	Model                  interface{}            `json:"Model"`
	NetworkDeviceFunctions *Link                  `json:"NetworkDeviceFunctions,omitempty"`
	NetworkPorts           *Link                  `json:"NetworkPorts,omitempty"`
	Oem                    *Oem                   `json:"Oem,omitempty"`
	PartNumber             interface{}            `json:"PartNumber"`
	Ports                  *Link                  `json:"Ports,omitempty"`
	Processors             *Link                  `json:"Processors,omitempty"`
	SerialNumber           interface{}            `json:"SerialNumber"`
	SKU                    interface{}            `json:"SKU"`
	Status                 *Status                `json:"Status,omitempty"`
}

// Controllers are The set of network controllers ASICs that make up this NetworkAdapter
type Controllers struct {
	ControllerCapabilities *ControllerCapabilities `json:"ControllerCapabilities,omitempty"`
	FirmwarePackageVersion interface{}             `json:"FirmwarePackageVersion"`
	Identifiers            *Identifier             `json:"Identifiers,omitempty"`
	Links                  *NLinks                 `json:"Links,omitempty"`
	PCIeInterface          *PCIeInterface          `json:"PCIeInterface,omitempty"`
	Location               *Location               `json:"Location,omitempty"`
}

// ControllerCapabilities are The capabilities of a controller
type ControllerCapabilities struct {
	DataCenterBridging         *DataCenterBridging    `json:"DataCenterBridging,omitempty"`
	NetworkDeviceFunctionCount interface{}            `json:"NetworkDeviceFunctionCount"`
	NetworkPortCount           interface{}            `json:"NetworkPortCount"`
	NPAR                       *NPAR                  `json:"NPAR,omitempty"`
	NPIV                       *NPIV                  `json:"NPIV,omitempty"`
	VirtualizationOffload      *VirtualizationOffload `json:"VirtualizationOffload,omitempty"`
}

// DataCenterBridging are Data center bridging (DCB) for this controller
type DataCenterBridging struct {
	Capable interface{} `json:"Capable"`
}

// NPAR is NIC Partitioning (NPAR) capabilities for this controller
type NPAR struct {
	NparCapable interface{} `json:"NparCapable"`
	NparEnabled interface{} `json:"NparEnabled"`
}

// NPIV is N_Port ID Virtualization (NPIV) capabilities for this controller
type NPIV struct {
	MaxDeviceLogins interface{} `json:"MaxDeviceLogins"`
	MaxPortLogins   interface{} `json:"MaxPortLogins"`
}

// VirtualizationOffload is Virtualization offload for this controller
type VirtualizationOffload struct {
	SRIOV           *SRIOV           `json:"SRIOV,omitempty"`
	VirtualFunction *VirtualFunction `json:"VirtualFunction,omitempty"`
}

// SRIOV is Single-root input/output virtualization (SR-IOV) capabilities
type SRIOV struct {
	SRIOVVEPACapable interface{} `json:"SRIOVVEPACapable"`
}

// VirtualFunction is a virtual function of a controller
type VirtualFunction struct {
	DeviceMaxCount         interface{} `json:"DeviceMaxCount"`
	MinAssignmentGroupSize interface{} `json:"MinAssignmentGroupSize"`
	NetworkPortMaxCount    interface{} `json:"NetworkPortMaxCount,"`
}

// NLinks is links to other resources that are related to this resource
type NLinks struct {
	NetworkDeviceFunctions []*Link `json:"NetworkDeviceFunctions,omitempty"`
	NetworkPorts           []*Link `json:"NetworkPorts,omitempty"`
	Oem                    *Oem    `json:"Oem,omitempty"`
	PCIeDevices            []*Link `json:"PCIeDevices,omitempty"`
	Ports                  []*Link `json:"Ports,omitempty"`
}

// NetworkAdapterActions is available actions for this resource
type NetworkAdapterActions struct {
	ResetSettings interface{} `json:"#NetworkAdapter.ResetSettingsToDefault,omitempty"`
	Oem           *Oem        `json:"Oem,omitempty"`
}
