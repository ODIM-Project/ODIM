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

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

// ComputerSystem redfish structure
type ComputerSystem struct {
	Ocontext           string             `json:"@odata.context"`
	Oid                string             `json:"@odata.id"`
	Otype              string             `json:"@odata.type"`
	Oetag              string             `json:"@odata.etag,omitempty"`
	ID                 string             `json:"Id"`
	Description        string             `json:"Description"`
	Name               string             `json:"Name"`
	Actions            *OemActions        `json:"Actions,omitempty"`
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
	PCIeDevices        []PCIeDevice       `json:"PCIeDevices"`
	PCIeFunctions      []PCIeFunction     `json:"PCIeFunctions"`
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
	Redundancy         []Redundancy       `json:"Redundancy"`
	SecureBoot         SecureBoot         `json:"SecureBoot"`
	SimpleStorage      SimpleStorage      `json:"SimpleStorage"`
	Status             Status             `json:"Status"`
	Storage            Storage            `json:"Storage"`
	TrustedModules     []TrustedModule    `json:"TrustedModules"`
	Oem                Oem                `json:"Oem,omitempty"`
	PCIeDevicesCount   int                `json:"PCIeDevices@odata.count,omitempty"`
}

// Bios redfish structure
type Bios struct {
	Oid         string `json:"@odata.id"`
	Ocontext    string `json:"@odata.context,omitempty"`
	Oetag       string `json:"@odata.etag,omitempty"`
	Otype       string `json:"@odata.type,omitempty"`
	Description string `json:"description,omitempty"`
	ID          string `json:"Id,omitempty"`
	Name        string `json:"Name,omitempty"`
	Oem         Oem    `json:"Oem,omitempty"`
	/*The reference to the Attribute Registry that lists the metadata describing the
	BIOS attribute settings in this resource.
	*/
	AttributeRegistry string                 `json:"AttributeRegistry,omitempty"` // read-only (null)
	Attributes        map[string]interface{} `json:"Attributes,omitempty"`        // object
}

// Boot redfish structure
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
/*
EthernetInterface 1.5.0
This resource shall be used to represent NIC resources as part of the Redfish specification.
URIs:
/redfish/v1/Managers/{ManagerId}/EthernetInterfaces/{EthernetInterfaceId}
/redfish/v1/Systems/{ComputerSystemId}/EthernetInterfaces/{EthernetInterfaceId}
*/
type EthernetInterfaces struct {
	Oid                    string              `json:"@odata.id"`
	Ocontext               string              `json:"@odata.context,omitempty"`
	Oetag                  string              `json:"@odata.etag,omitempty"`
	Otype                  string              `json:"@odata.type,omitempty"`
	Description            string              `json:"description,omitempty"`
	ID                     string              `json:"Id,omitempty"`
	Name                   string              `json:"Name,omitempty"`
	Oem                    Oem                 `json:"Oem,omitempty"`
	AutoNeg                bool                `json:"AutoNeg,omitempty"`
	DHCPv4                 DHCPv4              `json:"DHCPv4,omitempty"`
	DHCPv6                 DHCPv6              `json:"DHCPv6,omitempty"`
	FQDN                   string              `json:"FQDN,omitempty"`
	FullDuplex             bool                `json:"FullDuplex,omitempty"`
	HostName               string              `json:"HostName,omitempty"`
	InterfaceEnabled       bool                `json:"InterfaceEnabled,omitempty"`
	IPv4Addresses          []IPv4Address       `json:"IPv4Addresses,omitempty"`
	IPv4StaticAddresses    []IPv4Address       `json:"IPv4StaticAddresses,omitempty"`
	IPv6Addresses          []IPv6Address       `json:"IPv6Addresses,omitempty"`
	IPv6AddressPolicyTable []IPv6AddressPolicy `json:"IPv6AddressPolicyTable,omitempty"`
	IPv6DefaultGateway     string              `json:"IPv6DefaultGateway,omitempty"`
	IPv6StaticAddresses    []IPv6StaticAddress `json:"IPv6StaticAddresses,omitempty"`
	/* Note: IPv6GatewayStaticAddress and IPv6StaticAddress objects or exactly same,
	   decided to use IPv6StaticAddress in place of IPv6GatewayStaticAddress to achieving
	   code reusability in below line.
	*/
	IPv6StaticDefaultGateways  []IPv6StaticAddress        `json:"IPv6StaticDefaultGateways,omitempty"`
	Links                      []Link                     `json:"Links,omitempty"`
	LinkStatus                 string                     `json:"LinkStatus,omitempty"`
	MACAddress                 string                     `json:"MACAddress,omitempty"`
	MaxIPv6StaticAddresses     int                        `json:"MaxIPv6StaticAddresses,omitempty"`
	MTUSize                    int                        `json:"MTUSize,omitempty"`
	NameServers                []string                   `json:"NameServers,omitempty"`
	PermanentMACAddress        string                     `json:"PermanentMACAddress,omitempty"`
	SpeedMbps                  int                        `json:"SpeedMbps,omitempty"`
	StatelessAddressAutoConfig StatelessAddressAutoConfig `json:"StatelessAddressAutoConfig,omitempty"`
	StaticNameServers          []string                   `json:"StaticNameServers,omitempty"`
	Status                     Status                     `json:"Status,omitempty"`
	UefiDevicePath             string                     `json:"UefiDevicePath,omitempty"`
	VLAN                       VLAN                       `json:"VLAN,omitempty"`
	VLANs                      VLANs                      `json:"VLANs,omitempty"`
}

//DHCPv4 in place object
type DHCPv4 struct {
	DHCPEnabled     bool   `json:"DHCPEnabled"`
	FallbackAddress string `json:"FallbackAddress"` //enum
	UseDNSServers   bool   `json:"UseDNSServers"`
	UseDomainName   bool   `json:"UseDomainName"`
	UseGateway      bool   `json:"UseGateway"`
	UseNTPServers   bool   `json:"UseNTPServers"`
	UseStaticRoutes bool   `json:"UseStaticRoutes"`
}

//DHCPv6 in place object
type DHCPv6 struct {
	OperatingMode  string `json:"OperatingMode"` //enum
	UseDNSServers  bool   `json:"UseDNSServers"`
	UseDomainName  bool   `json:"UseDomainName"`
	UseNTPServers  bool   `json:"UseNTPServers"`
	UseRapidCommit bool   `json:"UseRapidCommit"`
}

//IPv4Address in place object
type IPv4Address struct {
	Address       string `json:"Address"`
	AddressOrigin string `json:"AddressOrigin"` //enum
	Gateway       string `json:"Gateway"`
	Oem           Oem    `json:"Oem"`
	SubnetMask    string `json:"SubnetMask"`
}

// IPv6Address in place object
type IPv6Address struct {
	Address       string `json:"Address"`
	AddressOrigin string `json:"AddressOrigin"` //enum
	AddressState  string `json:"AddressState"`  //enum
	Oem           Oem    `json:"Oem"`
	PrefixLength  string `json:"PrefixLength"`
}

//IPv6StaticAddress in place object
type IPv6StaticAddress struct {
	Address      string `json:"Address"`
	Oem          Oem    `json:"Oem"`
	PrefixLength string `json:"PrefixLength"`
}

//IPv6AddressPolicy in place object
type IPv6AddressPolicy struct {
	Label      int    `json:"Label"`
	Precedence int    `json:"Precedence"`
	Prefix     string `json:"Prefix"`
}

//StatelessAddressAutoConfig in place object
type StatelessAddressAutoConfig struct {
	IPv4AutoConfigEnabled bool `json:"IPv4AutoConfigEnabled"`
	IPv6AutoConfigEnabled bool `json:"IPv6AutoConfigEnabled"`
}

//VLAN redfish structure
type VLAN struct {
	Oid string `json:"@odata.id"`
}

//VLANs redfish structure
type VLANs struct {
	Oid string `json:"@odata.id"`
}

// HostedServices redfish structure
type HostedServices struct {
	Oem             Oem             `json:"Oem"`
	StorageServices StorageServices `json:"StorageServices"`
}

// HostWatchdogTimer redfish structure
type HostWatchdogTimer struct {
	FunctionEnabled bool   `json:"FunctionEnabled"`
	Oem             Oem    `json:"Oem"`
	Status          Status `json:"Status"`
	TimeoutAction   string `json:"TimeoutAction"`
	WarningAction   string `json:"WarningAction"`
}

// Memory redfish structure
type Memory struct {
	Oid                                     string                `json:"@odata.id"`
	Ocontext                                string                `json:"@odata.context,omitempty"`
	Oetag                                   string                `json:"@odata.etag,omitempty"`
	Otype                                   string                `json:"@odata.type,omitempty"`
	Description                             string                `json:"description,omitempty"`
	ID                                      string                `json:"Id,omitempty"`
	Name                                    string                `json:"Name,omitempty"`
	Oem                                     Oem                   `json:"Oem,omitempty"`
	AllocationAlignmentMiB                  int                   `json:"AllocationAlignmentMiB,omitempty"`
	AllocationIncrementMiB                  int                   `json:"AllocationIncrementMiB,omitempty"`
	AllowedSpeedsMHz                        []int                 `json:"AllowedSpeedsMHz,omitempty"`
	Assembly                                Assembly              `json:"Assembly,omitempty"`
	BaseModuleType                          string                `json:"BaseModuleType,omitempty"` //enum
	BusWidthBits                            int                   `json:"BusWidthBits,omitempty"`
	CacheSizeMiB                            int                   `json:"CacheSizeMiB,omitempty"`
	CapacityMiB                             int                   `json:"CapacityMiB,omitempty"`
	ConfigurationLocked                     bool                  `json:"ConfigurationLocked,omitempty"`
	DataWidthBits                           int                   `json:"DataWidthBits,omitempty"`
	DeviceID                                string                `json:"DeviceID,omitempty"`
	DeviceLocator                           string                `json:"DeviceLocator,omitempty"`
	ErrorCorrection                         string                `json:"ErrorCorrection,omitempty"` //enum
	FirmwareAPIVersion                      string                `json:"FirmwareApiVersion,omitempty"`
	FirmwareRevision                        string                `json:"FirmwareRevision,omitempty"`
	FunctionClasses                         []string              `json:"FunctionClasses,omitempty"`
	IsRankSpareEnabled                      bool                  `json:"IsRankSpareEnabled,omitempty"`
	IsSpareDeviceEnabled                    bool                  `json:"IsSpareDeviceEnabled,omitempty"`
	Links                                   Links                 `json:"Links,omitempty"`
	Location                                Location              `json:"Location,omitempty"`
	LogicalSizeMiB                          int                   `json:"LogicalSizeMiB,omitempty"`
	Manufacturer                            string                `json:"Manufacturer,omitempty"`
	MaxTDPMilliWatts                        []int                 `json:"MaxTDPMilliWatts,omitempty"`
	MemoryDeviceType                        string                `json:"MemoryDeviceType,omitempty"` //enum
	MemoryLocation                          MemoryLocation        `json:"MemoryLocation,omitempty"`
	MemoryMedia                             []string              `json:"MemoryMedia,omitempty"` //enum
	MemorySubsystemControllerManufacturerID string                `json:"MemorySubsystemControllerManufacturerID,omitempty"`
	MemorySubsystemControllerProductID      string                `json:"MemorySubsystemControllerProductID,omitempty"`
	MemoryType                              string                `json:"MemoryType,omitempty"` //enum
	Metrics                                 Metrics               `json:"Metrics,omitempty"`
	ModuleManufacturerID                    string                `json:"ModuleManufacturerID,omitempty"`
	ModuleProductID                         string                `json:"ModuleProductID,omitempty"`
	NonVolatileSizeMiB                      int                   `json:"NonVolatileSizeMiB,omitempty"`
	OperatingMemoryModes                    []string              `json:"OperatingMemoryModes,omitempty"` //enum
	OperatingSpeedMhz                       int                   `json:"OperatingSpeedMhz,omitempty"`
	PartNumber                              string                `json:"PartNumber,omitempty"`
	PersistentRegionNumberLimit             int                   `json:"PersistentRegionNumberLimit,omitempty"`
	PersistentRegionSizeLimitMiB            int                   `json:"PersistentRegionSizeLimitMiB,omitempty"`
	PersistentRegionSizeMaxMiB              int                   `json:"PersistentRegionSizeMaxMiB,omitempty"`
	PowerManagementPolicy                   PowerManagementPolicy `json:"PowerManagementPolicy,omitempty"`
	RankCount                               int                   `json:"RankCount,omitempty"`
	Regions                                 []Region              `json:"Regions,omitempty"`
	SecurityCapabilities                    SecurityCapabilities  `json:"SecurityCapabilities,omitempty"`
	SecurityState                           string                `json:"SecurityState,omitempty"` //enum
	SerialNumber                            string                `json:"SerialNumber,omitempty"`
	SpareDeviceCount                        int                   `json:"SpareDeviceCount,omitempty"`
	Status                                  Status                `json:"Status,omitempty"`
	SubsystemDeviceID                       string                `json:"SubsystemDeviceID,omitempty"`
	SubsystemVendorID                       string                `json:"SubsystemVendorID,omitempty"`
	VendorID                                string                `json:"VendorID,omitempty"`
	VolatileRegionNumberLimit               int                   `json:"VolatileRegionNumberLimit,omitempty"`
	VolatileRegionSizeLimitMiB              int                   `json:"VolatileRegionSizeLimitMiB,omitempty"`
	VolatileRegionSizeMaxMiB                int                   `json:"VolatileRegionSizeMaxMiB,omitempty"`
	VolatileSizeMiB                         int                   `json:"VolatileSizeMiB,omitempty"`
}

//MemoryLocation in place object
type MemoryLocation struct {
	Channel          int `json:"Channel,omitempty"`
	MemoryController int `json:"MemoryController"`
	Slot             int `json:"Slot"`
	Socket           int `json:"Socket"`
}

//PowerManagementPolicy in place object
type PowerManagementPolicy struct {
	AveragePowerBudgetMilliWatts int  `json:"AveragePowerBudgetMilliWatts"`
	MaxTDPMilliWatts             int  `json:"MaxTDPMilliWatts"`
	PeakPowerBudgetMilliWatts    int  `json:"PeakPowerBudgetMilliWatts"`
	PolicyEnabled                bool `json:"PolicyEnabled"`
}

//Region in place object
type Region struct {
	MemoryClassification string `json:"MemoryClassification"` //enum
	OffsetMiB            int    `json:"OffsetMiB"`
	PassphraseEnabled    bool   `json:"PassphraseEnabled"`
	PassphraseState      bool   `json:"PassphraseState"`
	RegionID             string `json:"RegionId"`
	SizeMiB              int    `json:"SizeMiB"`
}

//SecurityCapabilities in place object
type SecurityCapabilities struct {
	ConfigurationLockCapable bool     `json:"ConfigurationLockCapable"`
	DataLockCapable          bool     `json:"DataLockCapable"`
	MaxPassphraseCount       int      `json:"MaxPassphraseCount"`
	PassphraseCapable        bool     `json:"PassphraseCapable"`
	PassphraseLockLimit      int      `json:"PassphraseLockLimit"`
	SecurityStates           []string `json:"SecurityStates"` //enum
}

// MemoryDomains redfish structure
type MemoryDomains struct {
	Oid                       string                   `json:"@odata.id"`
	Ocontext                  string                   `json:"@odata.context,omitempty"`
	Oetag                     string                   `json:"@odata.etag,omitempty"`
	Otype                     string                   `json:"@odata.type,omitempty"`
	Description               string                   `json:"description,omitempty"`
	ID                        string                   `json:"Id,omitempty"`
	Name                      string                   `json:"Name,omitempty"`
	Oem                       Oem                      `json:"Oem,omitempty"`
	AllowsBlockProvisioning   bool                     `json:"AllowsBlockProvisioning,omitempty"`
	AllowsMemoryChunkCreation bool                     `json:"AllowsMemoryChunkCreation,omitempty"`
	AllowsMirroring           bool                     `json:"AllowsMirroring,omitempty"`
	AllowsSparing             bool                     `json:"AllowsSparing,omitempty"`
	InterleavableMemorySets   []InterleavableMemorySet `json:"InterleavableMemorySets,omitempty"`
	MemoryChunks              MemoryChunks             `json:"MemoryChunks,omitempty"`
}

//InterleavableMemorySet in place object
type InterleavableMemorySet struct {
	MemorySet []Memory `json:"MemorySet"`
}

//MemoryChunks redfish structure
type MemoryChunks struct {
	Oid string `json:"@odata.id"`
}

// MemorySummary in place object
type MemorySummary struct {
	MemoryMirroring                string `json:"MemoryMirroring"`
	TotalSystemMemoryGiB           int    `json:"TotalSystemMemoryGiB"`
	TotalSystemPersistentMemoryGiB int    `json:"TotalSystemPersistentMemoryGiB"`
	Status                         Status `json:"Status"`
}

//NetworkInterfaces get
/*
NetworkInterface 1.1.2

A NetworkInterface contains references linking NetworkAdapter, NetworkPort, and NetworkDeviceFunction resources and represents the
functionality available to the containing system.
URIs:
/redfish/v1/Systems/{ComputerSystemId}/NetworkInterfaces/{NetworkInterfaceId}

*/
type NetworkInterfaces struct {
	Oid                    string                 `json:"@odata.id"`
	Ocontext               string                 `json:"@odata.context,omitempty"`
	Oetag                  string                 `json:"@odata.etag,omitempty"`
	Otype                  string                 `json:"@odata.type,omitempty"`
	Description            string                 `json:"description,omitempty"`
	ID                     string                 `json:"Id,omitempty"`
	Name                   string                 `json:"Name,omitempty"`
	Oem                    Oem                    `json:"Oem,omitempty"`
	Links                  Links                  `json:"Links,omitempty"`
	NetworkDeviceFunctions NetworkDeviceFunctions `json:"NetworkDeviceFunctions,omitempty"`
	NetworkPorts           NetworkPorts           `json:"NetworkPorts,omitempty"`
	Status                 Status                 `json:"Status,omitempty"`
}

//NetworkDeviceFunctions redfish structure
type NetworkDeviceFunctions struct {
	Oid string `json:"@odata.id"`
}

//NetworkPorts redfish structure
type NetworkPorts struct {
	Oid string `json:"@odata.id"`
}

/*
PCIeDevice 1.3.1

This resource shall be used to represent a PCIeDevice attached to a System.
URIs:
/redfish/v1/Chassis/{ChassisId}/PCIeDevices/{PCIeDeviceId}
/redfish/v1/Systems/{ComputerSystemId}/PCIeDevices/{PCIeDeviceId}
*/
type PCIeDevice struct {
	Oid             string         `json:"@odata.id"`
	Ocontext        string         `json:"@odata.context,omitempty"`
	Oetag           string         `json:"@odata.etag,omitempty"`
	Otype           string         `json:"@odata.type"`
	Description     string         `json:"description,omitempty"`
	ID              string         `json:"Id"`
	Name            string         `json:"Name"`
	Oem             Oem            `json:"Oem,omitempty"`
	Assembly        *Assembly      `json:"Assembly,omitempty"`
	AssetTag        string         `json:"AssetTag,omitempty"`
	DeviceType      string         `json:"DeviceType,omitempty"` //enum
	FirmwareVersion string         `json:"FirmwareVersion,omitempty"`
	Links           *Links         `json:"Links,omitempty"`
	Manufacturer    string         `json:"Manufacturer,omitempty"`
	Model           string         `json:"Model,omitempty"`
	PartNumber      string         `json:"PartNumber,omitempty"`
	PCIeInterface   *PCIeInterface `json:"PCIeInterface,omitempty"`
	SerialNumber    string         `json:"SerialNumber,omitempty"`
	SKU             string         `json:"SKU,omitempty"`
	Status          *Status        `json:"Status,omitempty"`
}

//PCIeInterface in place object
type PCIeInterface struct {
	LanesInUse  int    `json:"LanesInUse,omitempty"`
	MaxLanes    int    `json:"MaxLanes,omitempty"`
	MaxPCIeType string `json:"MaxPCIeType,omitempty"` //enum
	Oem         Oem    `json:"Oem,omitempty"`
	PCIeType    string `json:"PCIeType,omitempty"` //enum
}

/*
PCIeFunction 1.2.2
This resource shall be used to represent a PCIeFunction attached to a System.
URIs:
/redfish/v1/Chassis/{ChassisId}/PCIeDevices/{PCIeDeviceId}/PCIeFunctions/{PCIeFunctionId}
/redfish/v1/Systems/{ComputerSystemId}/PCIeDevices/{PCIeDeviceId}/PCIeFunctions/{PCIeFunctionId}
*/
type PCIeFunction struct {
	Oid               string `json:"@odata.id"`
	Ocontext          string `json:"@odata.context,omitempty"`
	Oetag             string `json:"@odata.etag,omitempty"`
	Otype             string `json:"@odata.type,omitempty"`
	Description       string `json:"description,omitempty"`
	ID                string `json:"Id,omitempty"`
	Name              string `json:"Name,omitempty"`
	Oem               Oem    `json:"Oem,omitempty"`
	ClassCode         string `json:"ClassCode,omitempty"`
	DeviceClass       string `json:"DeviceClass,omitempty"` //enum
	DeviceID          string `json:"DeviceId,omitempty"`
	FunctionID        int    `json:"FunctionId,omitempty"`
	FunctionType      string `json:"FunctionType,omitempty"` //enum
	Links             Links  `json:"Links,omitempty"`
	RevisionID        string `json:"RevisionId,omitempty"`
	Status            Status `json:"Status,omitempty"`
	SubsystemID       string `json:"SubsystemId,omitempty"`
	SubsystemVendorID string `json:"SubsystemVendorId,omitempty"`
	VendorID          string `json:"VendorId,omitempty"`
}

/*
Processors 1.5.0

This resource shall be used to represent a single processor contained within a system.
URIs:
/redfish/v1/Systems/{ComputerSystemId}/Processors/{ProcessorId}
/redfish/v1/Systems/{ComputerSystemId}/Processors/{ProcessorId}/SubProcessors/{ProcessorId2}

*/
type Processors struct {
	Oid                   string                `json:"@odata.id"`
	Ocontext              string                `json:"@odata.context,omitempty"`
	Oetag                 string                `json:"@odata.etag,omitempty"`
	Otype                 string                `json:"@odata.type,omitempty"`
	Description           string                `json:"description,omitempty"`
	ID                    string                `json:"Id,omitempty"`
	Name                  string                `json:"Name,omitempty"`
	Oem                   Oem                   `json:"Oem,omitempty"`
	AccelerationFunctions AccelerationFunctions `json:"AccelerationFunctions,omitempty"`
	Assembly              Assembly              `json:"Assembly,omitempty"`
	FPGA                  FPGA                  `json:"FPGA,omitempty"`
	InstructionSet        string                `json:"InstructionSet,omitempty"` //enum
	Links                 Links                 `json:"Links,omitempty"`
	Location              Location              `json:"Location,omitempty"`
	Manufacturer          string                `json:"Manufacturer,omitempty"`
	MaxSpeedMHz           int                   `json:"MaxSpeedMHz,omitempty"`
	MaxTDPWatts           int                   `json:"MaxTDPWatts,omitempty"`
	Metrics               Metrics               `json:"Metrics,omitempty"`
	Model                 string                `json:"Model,omitempty"`
	ProcessorArchitecture string                `json:"ProcessorArchitecture,omitempty"` //enum
	ProcessorID           ProcessorID           `json:"ProcessorId,omitempty"`
	ProcessorMemory       []ProcessorMemory     `json:"ProcessorMemory,omitempty"`
	ProcessorType         string                `json:"ProcessorType,omitempty"` //enum
	Socket                string                `json:"Socket,omitempty"`
	Status                Status                `json:"Status,omitempty"`
	SubProcessors         SubProcessors         `json:"SubProcessors,omitempty"`
	TDPWatts              int                   `json:"TDPWatts,omitempty"`
	TotalCores            int                   `json:"TotalCores,omitempty"`
	TotalEnabledCores     int                   `json:"TotalEnabledCores,omitempty"`
	TotalThreads          int                   `json:"TotalThreads,omitempty"`
	UUID                  string                `json:"UUID,omitempty"`
}

//AccelerationFunctions redfish structure
type AccelerationFunctions struct {
	Oid string `json:"@odata.id"`
}

//FPGA in place object
type FPGA struct {
	ExternalInterfaces   []HostInterface       `json:"ExternalInterfaces"`
	FirmwareID           string                `json:"FirmwareId"`
	FirmwareManufacturer string                `json:"FirmwareManufacturer"`
	FirmwareVersion      string                `json:"FirmwareVersion"`
	FpgaType             string                `json:"FpgaType"` //enum
	HostInterface        HostInterface         `json:"HostInterface"`
	Model                string                `json:"Model"`
	Oem                  Oem                   `json:"Oem"`
	PCIeVirtualFunctions int                   `json:"PCIeVirtualFunctions"`
	ProgrammableFromHost bool                  `json:"ProgrammableFromHost"`
	ReconfigurationSlots []ReconfigurationSlot `json:"ReconfigurationSlots"`
}

//HostInterface in place object
type HostInterface struct {
	Ethernet      Ethernet      `json:"Ethernet"`
	InterfaceType string        `json:"InterfaceType"` //enum
	PCIe          PCIeInterface `json:"PCIe"`
}

//Ethernet in place object
type Ethernet struct {
	MaxLanes     int `json:"MaxLanes"`
	MaxSpeedMbps int `json:"MaxSpeedMbps"`
	Oem          Oem `json:"Oem"`
}

//ReconfigurationSlot in place object
type ReconfigurationSlot struct {
	AccelerationFunction AccelerationFunction `json:"AccelerationFunction"`
	ProgrammableFromHost bool                 `json:"ProgrammableFromHost"`
	SlotID               string               `json:"SlotId"`
	UUID                 string               `json:"UUID"`
}

//AccelerationFunction redfish structure
type AccelerationFunction struct {
	Oid string `json:"@odata.id"`
}

//ProcessorID in place object
type ProcessorID struct {
	EffectiveFamily         string `json:"EffectiveFamily"`
	EffectiveModel          string `json:"EffectiveModel"`
	IdentificationRegisters string `json:"IdentificationRegisters"`
	MicrocodeInfo           string `json:"MicrocodeInfo"`
	Step                    string `json:"Step"`
	VendorID                string `json:"VendorId"`
}

//ProcessorMemory in place object
type ProcessorMemory struct {
	CapacityMiB      int    `json:"CapacityMiB"`
	IntegratedMemory bool   `json:"IntegratedMemory"`
	MemoryType       string `json:"MemoryType"` //enum
	SpeedMHz         int    `json:"SpeedMHz"`
}

//SubProcessors redfish structure
type SubProcessors struct {
	Oid string `json:"@odata.id"`
}

// ProcessorSummary redfish structure
type ProcessorSummary struct {
	Count                 int     `json:"Count"`
	LogicalProcessorCount int     `json:"LogicalProcessorCount"`
	Model                 string  `json:"Model"`
	Metrics               Metrics `json:"Metrics"`
	Status                Status  `json:"Status"`
}

// SecureBoot redfish structure
type SecureBoot struct {
	Oid                   string `json:"@odata.id"`
	Ocontext              string `json:"@odata.context,omitempty"`
	Oetag                 string `json:"@odata.etag,omitempty"`
	Otype                 string `json:"@odata.type,omitempty"`
	Description           string `json:"description,omitempty"`
	ID                    string `json:"Id,omitempty"`
	Name                  string `json:"Name,omitempty"`
	Oem                   Oem    `json:"Oem,omitempty"`
	SecureBootCurrentBoot string `json:"SecureBootCurrentBoot,omitempty"`
	SecureBootEnable      bool   `json:"SecureBootEnable,omitempty"`
	SecureBootMode        string `json:"SecureBootMode,omitempty"`
}

// SimpleStorage redfish structure
type SimpleStorage struct {
	Oid            string   `json:"@odata.id"`
	Ocontext       string   `json:"@odata.context,omitempty"`
	Oetag          string   `json:"@odata.etag,omitempty"`
	Otype          string   `json:"@odata.type,omitempty"`
	Description    string   `json:"description,omitempty"`
	ID             string   `json:"Id,omitempty"`
	Name           string   `json:"Name,omitempty"`
	Oem            Oem      `json:"Oem,omitempty"`
	Devices        []Device `json:"Devices,omitempty"`
	Links          Link     `json:"Links,omitempty"`
	Status         Status   `json:"Status,omitempty"`
	UefiDevicePath string   `json:"UefiDevicePath,omitempty"`
}

//Device in place object
type Device struct {
	CapacityBytes int    `json:"CapacityBytes"`
	Manufacturer  string `json:"Manufacturer"`
	Model         string `json:"Model"`
	Name          string `json:"Name"`
	Oem           Oem    `json:"Oem"`
	Status        Status `json:"Status"`
}

// TrustedModule redfish structure
type TrustedModule struct {
	FirmwareVersion        string `json:"FirmwareVersion"`
	FirmwareVersion2       string `json:"FirmwareVersion2"`
	InterfaceType          string `json:"InterfaceType"`
	InterfaceTypeSelection string `json:"InterfaceTypeSelection"`
	Oem                    Oem    `json:"Oem"`
	Status                 Status `json:"Status"`
}

// BootOptions redfish structure
type BootOptions struct {
	Oid string `json:"@odata.id"`
}

//Certificates redfish structure
type Certificates struct {
	Oid string `json:"@odata.id"`
}

// StorageServices redfish structure
type StorageServices struct {
	Oid string `json:"@odata.id"`
}

// Metrics redfish structure
type Metrics struct {
	Oid string `json:"@odata.id"`
}

// SaveInMemory will create the ComputerSystem data in in-memory DB, with key as UUID
// Takes:
//	none as parameter, but recieves c of type *ComputerSystem as pointeer reciever impicitly.
// Returns:
//	err of type error
//
//	On Success - returns nil value
//	On Failure - return non nil value
func (c *ComputerSystem) SaveInMemory(deviceUUID string) *errors.Error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to connect to DB: ", err.Error())
	}
	if err := connPool.Create("computersystem", deviceUUID, c); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to create new computersystem: ", err.Error())
	}
	return nil
}
