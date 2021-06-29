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

//Package dpmodel ...
package dpmodel

// ComputerSystem ..
type ComputerSystem struct {
	Ocontext           string             `json:"@odata.context"`
	Oid                string             `json:"@odata.id"`
	Otype              string             `json:"@odata.type"`
	Oetag              string             `json:"@odata.etag,omitempty"`
	ID                 string             `json:"Id"`
	Description        string             `json:"Description"`
	Name               string             `json:"Name"`
	AssetTag           string             `json:"AssetTag"`
	BiosVersion        string             `json:"BiosVersion"`
	HostName           string             `json:"HostName"`
	IndicatorLED       string             `json:"IndicatorLED"`
	Manufacturer       string             `json:"Manufacturer"`
	Model              string             `json:"Model"`
	PartNumber         string             `json:"PartNumber"`
	PowerRestorePolicy string             `json:"PowerRestorePolicy"`
	PowerState         string             `json:"PowerState"`
	SerialNumber       string             `json:"SerialNumber"`
	SKU                string             `json:"SKU"`
	SubModel           string             `json:"SubModel"`
	SystemType         string             `json:"SystemType"`
	UUID               string             `json:"UUID"`
	HostingRoles       []string           `json:"HostingRoles"`
	PCIeDevices        []PCIeDevices      `json:"PCIeDevices"`
	PCIeFunctions      []PCIeFunctions    `json:"PCIeFunctions"`
	Bios               Bios               `json:"Bios"`
	Boot               Boot               `json:"Boot"`
	EthernetInterfaces EthernetInterfaces `json:"EthernetInterfaces"`
	HostedServices     HostedServices     `json:"HostedServices"`
	HostWatchdogTimer  HostWatchdogTimer  `json:"HostWatchdogTimer"`
	Links              Links              `json:"Links"`
	LogServices        LogServices        `json:"LogServices"`
	Memory             Memory             `json:"Memory"`
	MemoryDomains      MemoryDomains      `json:"MemoryDomains"`
	MemorySummary      MemorySummary      `json:"MemorySummary"`
	NetworkInterfaces  NetworkInterfaces  `json:"NetworkInterfaces"`
	Processors         Processors         `json:"Processors"`
	ProcessorSummary   ProcessorSummary   `json:"ProcessorSummary"`
	Redundancy         Redundancy         `json:"Redundancy"`
	SecureBoot         SecureBoot         `json:"SecureBoot"`
	SimpleStorage      SimpleStorage      `json:"SimpleStorage"`
	Status             Status             `json:"Status"`
	Storage            Storage            `json:"Storage"`
	TrustedModules     []TrustedModules   `json:"TrustedModules"`
}

// Bios get
type Bios struct {
	Oid string `json:"@odata.id"`
}

// Boot get
type Boot struct {
	AliasBootOrder               []string     `json:"AliasBootOrder"`
	BootNext                     string       `json:"BootNext"`
	BootOptions                  BootOptions  `json:"BootOptions"`
	BootOrder                    []string     `json:"BootOrder"`
	BootOrderPropertySelection   string       `json:"BootOrderPropertySelection"`
	BootSourceOverrideEnabled    string       `json:"BootSourceOverrideEnabled"`
	BootSourceOverrideMode       string       `json:"BootSourceOverrideMode"`
	BootSourceOverrideTarget     string       `json:"BootSourceOverrideTarget"`
	Certificates                 Certificates `json:"Certificates"`
	UefiTargetBootSourceOverride string       `json:"UefiTargetBootSourceOverride"`
}

// EthernetInterfaces get
type EthernetInterfaces struct {
	Oid string `json:"@odata.id"`
}

// HostedServices ..
type HostedServices struct {
	Oem             Oem             `json:"Oem"`
	StorageServices StorageServices `json:"StorageServices"`
}

// HostWatchdogTimer ..
type HostWatchdogTimer struct {
	FunctionEnabled bool   `json:"FunctionEnabled"`
	Oem             Oem    `json:"Oem"`
	Status          Status `json:"Status"`
	TimeoutAction   string `json:"TimeoutAction"`
	WarningAction   string `json:"WarningAction"`
}

/*
type Links struct {
	Chassis                        []Chassis
	ConsumingComputerSystems       []ConsumingComputerSystems
	CooledBy                       []CooledBy
	Endpoints                      []Endpoints
	ManagedBy                      []ManagedBy
	Oem                            Oem
	PoweredBy                      []PoweredBy
	ResourceBlocks                 []ResourceBlocks
	SupplyingComputerSystems       []SupplyingComputerSystems
}

*/

// Memory get
type Memory struct {
	Oid string `json:"@odata.id"`
}

// MemoryDomains get
type MemoryDomains struct {
	Oid string `json:"@odata.id"`
}

// MemorySummary get
type MemorySummary struct {
	MemoryMirroring                string `json:"MemoryMirroring"`
	TotalSystemMemoryGiB           int    `json:"TotalSystemMemoryGiB"`
	TotalSystemPersistentMemoryGiB int    `json:"TotalSystemPersistentMemoryGiB"`
	Status                         Status `json:"Status"`
}

//NetworkInterfaces get
type NetworkInterfaces struct {
	Oid string `json:"@odata.id"`
}

//PCIeDevices get
type PCIeDevices struct {
	Oid string `json:"@odata.id"`
}

//PCIeFunctions get
type PCIeFunctions struct {
	Oid string `json:"@odata.id"`
}

// Processors get
type Processors struct {
	Oid string `json:"@odata.id"`
}

// ProcessorSummary get
type ProcessorSummary struct {
	Count                 int     `json:"Count"`
	LogicalProcessorCount int     `json:"LogicalProcessorCount"`
	Model                 string  `json:"Model"`
	Metrics               Metrics `json:"Metrics"`
	Status                Status  `json:"Status"`
}

// Redundancy get
type Redundancy struct {
	Oid string `json:"@odata.id"`
}

// SecureBoot get
type SecureBoot struct {
	Oid string `json:"@odata.id"`
}

// SimpleStorage get
type SimpleStorage struct {
	Oid string `json:"@odata.id"`
}

// Storage get
type Storage struct {
	Oid string `json:"@odata.id"`
}

// TrustedModules get
type TrustedModules struct {
	FirmwareVersion        string `json:"FirmwareVersion"`
	FirmwareVersion2       string `json:"FirmwareVersion2"`
	InterfaceType          string `json:"InterfaceType"`
	InterfaceTypeSelection string `json:"InterfaceTypeSelection"`
	Oem                    Oem    `json:"Oem"`
	Status                 Status `json:"Status"`
}

// BootOptions get
type BootOptions struct {
	Oid string `json:"@odata.id"`
}

//Certificates get
type Certificates struct {
	Oid string `json:"@odata.id"`
}

// StorageServices get
type StorageServices struct {
	Oid string `json:"@odata.id"`
}

/*
type Chassis struct {
	Oid string `json:"@odata.id"`
}
type ConsumingComputerSystems struct {
	Oid string `json:"@odata.id"`
}
type CooledBy struct {
	Oid string `json:"@odata.id"`
}
type Endpoints struct {
	Oid string `json:"@odata.id"`
}
type ManagedBy struct {
	Oid string `json:"@odata.id"`
}
type PoweredBy struct {
	Oid string `json:"@odata.id"`
}
type ResourceBlocks struct {
	Oid string `json:"@odata.id"`
}
type SupplyingComputerSystems struct {
	Oid string `json:"@odata.id"`
}
*/

// Metrics get
type Metrics struct {
	Oid string `json:"@odata.id"`
}

// BiosSettings get
type BiosSettings struct {
	Oem OemDell `json:"Oem"`
}

//OemDell for bios settings
type OemDell struct {
	Dell Dell `json:"Dell"`
}

//Dell for bios setting's OEM
type Dell struct {
	Jobs OdataID `json:"Jobs"`
}

//OdataID contains link to a resource
type OdataID struct {
	OdataID string `json:"@odata.id"`
}
