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

// BootSource - AliasBootOrder, BootSourceOverrideTarget
// Ordered array of boot source aliases representing the persistent
// boot order associated with this computer system
type BootSource string

// AutomaticRetryConfig - The configuration of how the system retries booting automatically
type AutomaticRetryConfig string

// BootOrderTypes - The name of the boot order property that the system uses for the persistent boot order
type BootOrderTypes string

// BootSourceOverrideEnabled - The state of the boot source override feature
type BootSourceOverrideEnabled string

// BootSourceOverrideMode - The BIOS boot mode to use when the system boots from the BootSourceOverrideTarget boot source
type BootSourceOverrideMode string

// StopBootOnFault - If the boot should stop on a fault
type StopBootOnFault string

// TrustedModuleRequiredToBoot - The Trusted Module boot requirement
type TrustedModuleRequiredToBoot string

// CompositionUseCase - The composition use cases in which this computer system can participate
type CompositionUseCase string

// SystemType -The type of computer system that this resource represents
type SystemType string

// PowerRestorePolicyTypes -The desired power state of the system when power is restored after a power loss
type PowerRestorePolicyTypes string

// InterfaceType - The interface type of the Trusted Module
type InterfaceType string

// InterfaceTypeSelection - The interface type selection supported by this Trusted Module
type InterfaceTypeSelection string

// MemoryMirroring - The ability and type of memory mirroring that this computer system supports
type MemoryMirroring string

// WatchdogWarningActions - The action to perform when the watchdog timer is close to reaching its timeout value
type WatchdogWarningActions string

// WatchdogTimeoutActions - The action to perform when the watchdog timer reaches its timeout value
type WatchdogTimeoutActions string

// PowerMode - The power mode setting of the computer system
type PowerMode string

// CachePolicy - he cache policy to control how KMIP data is cached
type CachePolicy string

// HostingRole - The hosting roles that this computer system supports
type HostingRole string

// IndicatorLED - The state of the indicator LED, which identifies the system
type IndicatorLED string

// BootProgressTypes - The last boot progress state
type BootProgressTypes string

type ConnectedTypesSupported string

const (

	// BootSourceNone - Boot from the normal boot device
	BootSourceNone BootSource = "None"

	//BootSourcePxe - Boot from the Pre-Boot EXecution (PXE) environmen
	BootSourcePxe BootSource = "Pxe"

	//BootSourceFloppy - Boot from the floppy disk drive
	BootSourceFloppy BootSource = "Floppy"

	//BootSourceCd - Boot from the CD or DVD
	BootSourceCd BootSource = "Cd"

	//BootSourceUsb - Boot from a system BIOS-specified USB device
	BootSourceUsb BootSource = "Usb"

	//BootSourceHdd - Boot from a hard drive
	BootSourceHdd BootSource = "Hdd"

	// BootSourceBiosSetup - Boot to the BIOS setup utility
	BootSourceBiosSetup BootSource = "BiosSetup"

	//BootSourceUtilities -
	BootSourceUtilities BootSource = "Utilities"

	//BootSourceDiags - Boot to the manufacturer's diagnostics program
	BootSourceDiags BootSource = "Diags"

	//BootSourceUefiShell - Boot to the UEFI Shell
	BootSourceUefiShell BootSource = "UefiShell"

	//BootSourceUefiTarget - Boot to the UEFI device specified in the UefiTargetBootSourceOverride property
	BootSourceUefiTarget BootSource = "UefiTarget"

	//BootSourceSDCard - Boot from an SD card
	BootSourceSDCard BootSource = "SDCard"

	//BootSourceUefiHTTP - Boot from a UEFI HTTP network location
	BootSourceUefiHTTP BootSource = "UefiHttp"

	//BootSourceRemoteDrive - Boot from a remote drive, such as an iSCSI target
	BootSourceRemoteDrive BootSource = "RemoteDrive"

	//BootSourceUefiBootNext - Boot to the UEFI device that the BootNext property specifies
	BootSourceUefiBootNext BootSource = "UefiBootNext"

	//BootSourceRecovery - Boot to a system-designated recovery process or image
	BootSourceRecovery BootSource = "Recovery"

	//AutomaticRetryConfigDisabled - Disable automatic retrying of booting
	AutomaticRetryConfigDisabled AutomaticRetryConfig = "Disabled"

	//AutomaticRetryConfigRetryAttempts - Always automatically retry booting
	AutomaticRetryConfigRetryAttempts AutomaticRetryConfig = "RetryAttempts"

	//AutomaticRetryConfigRetryAlways - Automatic retrying of booting is based on a specified retry count
	AutomaticRetryConfigRetryAlways AutomaticRetryConfig = "RetryAlways"

	//BootOrderTypesBootOrder - The system uses the BootOrder property to specify the persistent boot order
	BootOrderTypesBootOrder BootOrderTypes = "BootOrder"

	//BootOrderTypesAliasBootOrder - The system uses the AliasBootOrder property to specify the persistent boot order
	BootOrderTypesAliasBootOrder BootOrderTypes = "AliasBootOrder"

	//BootSourceOverrideEnabledDisabled - The system boots normally
	BootSourceOverrideEnabledDisabled BootSourceOverrideEnabled = "Disabled"

	//BootSourceOverrideEnabledOnce - On its next boot cycle, the system boots one time to the boot source override target
	BootSourceOverrideEnabledOnce BootSourceOverrideEnabled = "Once"

	//BootSourceOverrideEnabledContinuous - The system boots to the target specified in the BootSourceOverrideTarget property until this property is `Disabled`
	BootSourceOverrideEnabledContinuous BootSourceOverrideEnabled = "Continuous"

	//BootSourceOverrideModeLegacy - The system boots in non-UEFI boot mode to the boot source override target
	BootSourceOverrideModeLegacy BootSourceOverrideMode = "Legacy"

	//BootSourceOverrideModeUEFI - The system boots in UEFI boot mode to the boot source override target
	BootSourceOverrideModeUEFI BootSourceOverrideMode = "UEFI"

	//StopBootOnFaultNever - The system performs any normal recovery actions during boot if a fault occurs
	StopBootOnFaultNever StopBootOnFault = "Never"

	//StopBootOnFaultAnyFault - The system should stop the boot on any fault
	StopBootOnFaultAnyFault StopBootOnFault = "AnyFault"

	//TrustedModuleRequiredToBootDisabled - No Trusted Module requirement to boot
	TrustedModuleRequiredToBootDisabled TrustedModuleRequiredToBoot = "Disabled"

	//TrustedModuleRequiredToBootRequired - A functional Trusted Module is required to boot
	TrustedModuleRequiredToBootRequired TrustedModuleRequiredToBoot = "Required"

	//BootProgressTypesNone - The system is not booting
	BootProgressTypesNone BootProgressTypes = "None"

	//BootProgressTypesPrimaryProcessorInitializationStarted - The system has started initializing the primary processor
	BootProgressTypesPrimaryProcessorInitializationStarted BootProgressTypes = "PrimaryProcessorInitializationStarted"

	//BootProgressTypesBusInitializationStarted - The system has started initializing the buses
	BootProgressTypesBusInitializationStarted BootProgressTypes = "BusInitializationStarted"

	//BootProgressTypesMemoryInitializationStarted - The system has started initializing the memory
	BootProgressTypesMemoryInitializationStarted BootProgressTypes = "MemoryInitializationStarted"

	//BootProgressTypesSecondaryProcessorInitializationStarted - The system has started initializing the remaining processors.
	BootProgressTypesSecondaryProcessorInitializationStarted BootProgressTypes = "SecondaryProcessorInitializationStarted"

	//BootProgressTypesPCIResourceConfigStarted - The system has started initializing the PCI resources
	BootProgressTypesPCIResourceConfigStarted BootProgressTypes = "PCIResourceConfigStarted"

	//BootProgressTypesSystemHardwareInitializationComplete - The system has completed initializing all hardware
	BootProgressTypesSystemHardwareInitializationComplete BootProgressTypes = "SystemHardwareInitializationComplete"

	//BootProgressTypesSetupEntered - The system has entered the setup utility
	BootProgressTypesSetupEntered BootProgressTypes = "SetupEntered"

	//BootProgressTypesOSBootStarted - The operating system has started booting
	BootProgressTypesOSBootStarted BootProgressTypes = "OSBootStarted"

	//BootProgressTypesOSRunning - The operating system is running
	BootProgressTypesOSRunning BootProgressTypes = "OSRunning"

	//BootProgressTypesOEM - A boot progress state in an OEM-defined format
	BootProgressTypesOEM BootProgressTypes = "OEM"

	//CompositionUseCaseResourceBlockCapable - This computer system supports being registered as a resource block
	//in order for it to participate in composition requests
	//UseCases
	CompositionUseCaseResourceBlockCapable CompositionUseCase = "ResourceBlockCapable"

	//CompositionUseCaseExpandableSystem - This computer system supports expandable system composition and is associated with a resource block
	CompositionUseCaseExpandableSystem CompositionUseCase = "ExpandableSystem"

	//SystemTypePhysical -A computer system
	SystemTypePhysical SystemType = "Physical"

	//SystemTypeVirtual - A virtual machine instance running on this syste
	SystemTypeVirtual SystemType = "Virtual"

	//SystemTypeOS - An operating system instance
	SystemTypeOS SystemType = "OS"

	//SystemTypePhysicallyPartitioned - A hardware-based partition of a computer system
	SystemTypePhysicallyPartitioned SystemType = "PhysicallyPartitioned"

	//SystemTypeVirtuallyPartitioned - A virtual or software-based partition of a computer system
	SystemTypeVirtuallyPartitioned SystemType = "VirtuallyPartitioned"

	//SystemTypeComposed - A computer system constructed by binding resource blocks together
	SystemTypeComposed SystemType = "Composed"

	//SystemTypeDPU - A computer system that performs the functions of a data processing unit, such as a SmartNIC
	SystemTypeDPU SystemType = "DPU"

	//InterfaceTypeTPM1_2 - Trusted Platform Module (TPM) 1.2.
	InterfaceTypeTPM1_2 InterfaceType = "TPM1_2"

	//InterfaceTypeTPM2_0 - Trusted Platform Module (TPM) 2.0
	InterfaceTypeTPM2_0 InterfaceType = "TPM2_0"

	//InterfaceTypeTCM1_0 -Trusted Cryptography Module (TCM) 1.0
	InterfaceTypeTCM1_0 InterfaceType = "TCM1_0"

	//InterfaceTypeSelectionNone - The TrustedModule does not support switching the InterfaceType
	InterfaceTypeSelectionNone InterfaceTypeSelection = "None"

	//InterfaceTypeSelectionFirmwareUpdate - The TrustedModule supports switching InterfaceType through a firmware update
	InterfaceTypeSelectionFirmwareUpdate InterfaceTypeSelection = "FirmwareUpdate"

	//InterfaceTypeSelectionBiosSetting - The TrustedModule supports switching InterfaceType through platform software,
	// such as a BIOS configuration attribute
	InterfaceTypeSelectionBiosSetting InterfaceTypeSelection = "BiosSetting"

	//InterfaceTypeSelectionOemMethod - The TrustedModule supports switching InterfaceType through an OEM proprietary mechanism
	InterfaceTypeSelectionOemMethod InterfaceTypeSelection = "OemMethod"

	//MemoryMirroringSystem - The system supports DIMM mirroring at the system level
	MemoryMirroringSystem MemoryMirroring = "System"

	//MemoryMirroringDIMM -The system supports DIMM mirroring at the DIMM level.  Individual DIMMs can be mirrored
	MemoryMirroringDIMM MemoryMirroring = "DIMM"

	//MemoryMirroringHybrid - The system supports a hybrid mirroring at the system and DIMM levels.
	MemoryMirroringHybrid MemoryMirroring = "Hybrid"

	//MemoryMirroringNone - The system does not support DIMM mirroring
	MemoryMirroringNone MemoryMirroring = "None"

	//WatchdogWarningActionsNone - No action taken
	WatchdogWarningActionsNone WatchdogWarningActions = "None"

	//WatchdogWarningActionsDiagnosticInterrupt -Raise a (typically non-maskable) Diagnostic Interrupt
	WatchdogWarningActionsDiagnosticInterrupt WatchdogWarningActions = "DiagnosticInterrupt"

	//WatchdogWarningActionsSMI - Raise a Systems Management Interrupt (SMI)
	WatchdogWarningActionsSMI WatchdogWarningActions = "SMI"

	//WatchdogWarningActionsMessagingInterrupt - Raise a legacy IPMI messaging interrupt
	WatchdogWarningActionsMessagingInterrupt WatchdogWarningActions = "MessagingInterrupt"

	//WatchdogWarningActionsSCI - Raise an interrupt using the ACPI System Control Interrupt (SCI)
	WatchdogWarningActionsSCI WatchdogWarningActions = "SCI"

	//WatchdogWarningActionsOEM - Perform an OEM-defined action
	WatchdogWarningActionsOEM WatchdogWarningActions = "OEM"

	//WatchdogTimeoutActionsNone - No action taken
	WatchdogTimeoutActionsNone WatchdogTimeoutActions = "None"

	//WatchdogTimeoutActionsResetSystem - Reset the system
	WatchdogTimeoutActionsResetSystem WatchdogTimeoutActions = "ResetSystem"

	//WatchdogTimeoutActionsPowerCycle -Power cycle the system.
	WatchdogTimeoutActionsPowerCycle WatchdogTimeoutActions = "PowerCycle"

	//WatchdogTimeoutActionsPowerDown - Power down the system
	WatchdogTimeoutActionsPowerDown WatchdogTimeoutActions = "PowerDown"

	//WatchdogTimeoutActionsOEM - Perform an OEM-defined action
	WatchdogTimeoutActionsOEM WatchdogTimeoutActions = "OEM"

	//PowerRestorePolicyTypesAlwaysOn - The system always powers on when power is applied
	PowerRestorePolicyTypesAlwaysOn PowerRestorePolicyTypes = "AlwaysOn"

	//PowerRestorePolicyTypesAlwaysOff -The system always remains powered off when power is applied
	PowerRestorePolicyTypesAlwaysOff PowerRestorePolicyTypes = "AlwaysOff"

	//PowerRestorePolicyTypesLastState - The system returns to its last on or off power state when power is applied
	PowerRestorePolicyTypesLastState PowerRestorePolicyTypes = "LastState"

	//PowerModeMaximumPerformance - The system performs at the highest speeds possible
	PowerModeMaximumPerformance PowerMode = "MaximumPerformance"

	//PowerModeBalancedPerformance -The system performs at the highest speeds while utilization is high
	//and performs at reduced speeds when the utilization is low
	PowerModeBalancedPerformance PowerMode = "BalancedPerformance"

	//PowerModePowerSaving -The system performs at reduced speeds to save power
	PowerModePowerSaving PowerMode = "PowerSaving"

	//PowerModeStatic - The system power mode is static
	PowerModeStatic PowerMode = "Static"

	//PowerModeOSControlled - The system power mode is controlled by the operating system
	PowerModeOSControlled PowerMode = "OSControlled"

	//PowerModeOEM - The system power mode is OEM-defined
	PowerModeOEM PowerMode = "OEM"

	//CachePolicyNone - The system does not cache KMIP data
	CachePolicyNone CachePolicy = "None"

	//CachePolicyAfterFirstUse - The system caches KMIP data after first use
	//for the duration specified by the CacheDuration property
	CachePolicyAfterFirstUse CachePolicy = "AfterFirstUse"

	//HostingRoleApplicationServer - The system hosts functionality that supports general purpose applications
	HostingRoleApplicationServer HostingRole = "ApplicationServer"

	//HostingRoleStorageServer - The system hosts functionality that supports the system acting as a storage server
	HostingRoleStorageServer HostingRole = "StorageServer"

	//HostingRoleSwitch - The system hosts functionality that supports the system acting as a switch
	HostingRoleSwitch HostingRole = "Switch"

	//HostingRoleAppliance - The system hosts functionality that supports the system acting as an appliance
	HostingRoleAppliance HostingRole = "Appliance"

	//HostingRoleBareMetalServer - The system hosts functionality that supports the system acting as a bare metal server
	HostingRoleBareMetalServer HostingRole = "BareMetalServer"

	//HostingRoleVirtualMachineServer -The system hosts functionality that supports the system acting as a virtual machine server
	HostingRoleVirtualMachineServer HostingRole = "VirtualMachineServer"

	//HostingRoleContainerServer - The system hosts functionality that supports the system acting as a container server
	HostingRoleContainerServer HostingRole = "ContainerServer"

	//IndicatorLEDUnknown - The state of the indicator LED cannot be determined
	IndicatorLEDUnknown IndicatorLED = "Unknown"

	//IndicatorLEDLit - The indicator LED is lit
	IndicatorLEDLit IndicatorLED = "Lit"

	//IndicatorLEDBlinking -The indicator LED is blinking
	IndicatorLEDBlinking IndicatorLED = "Blinking"

	//IndicatorLEDOff - The indicator LED is off
	IndicatorLEDOff IndicatorLED = "Off"

	ConnectedTypesSupportedKVMIP ConnectedTypesSupported = "KVMIP"

	ConnectedTypesSupportedOEM ConnectedTypesSupported = "OEM"
)

// AddResourceBlock redfish structure
// This action adds a resource block to a system
// This action shall add a resource block to a system
type AddResourceBlock struct {
	Target string `json:"Target,omitempty"`
	Title  string `json:"title"`
}

// ComputerSystem redfish structure
// The ComputerSystem schema represents a computer or system instance and the software-visible resources,
// or items within the data plane, such as memory, CPU, and other devices that it can access.
//
//	Details of those resources or subsystems are also linked through this resource
//
// Reference :ComputerSystem.v1_20_0.json
type ComputerSystem struct {
	Ocontext                        string                 `json:"@odata.context,omitempty"`
	Oid                             string                 `json:"@odata.id"`
	Otype                           string                 `json:"@odata.type"`
	Oetag                           string                 `json:"@odata.etag,omitempty"`
	ID                              string                 `json:"ID"`
	Description                     string                 `json:"Description,omitempty"`
	Name                            string                 `json:"Name"`
	Actions                         *ComputerSystemActions `json:"Actions,omitempty"`
	AssetTag                        string                 `json:"AssetTag,omitempty"`
	BiosVersion                     string                 `json:"BiosVersion,omitempty"`
	HostName                        string                 `json:"HostName,omitempty"`
	IndicatorLED                    string                 `json:"IndicatorLED,omitempty"` //enum
	Manufacturer                    string                 `json:"Manufacturer,omitempty"`
	Model                           string                 `json:"Model,omitempty"`
	PartNumber                      string                 `json:"PartNumber,omitempty"`
	PowerRestorePolicy              string                 `json:"PowerRestorePolicy,omitempty"` //enum
	PowerState                      string                 `json:"PowerState,omitempty"`
	SerialNumber                    string                 `json:"SerialNumber,omitempty"`
	SKU                             string                 `json:"SKU,omitempty"`
	SubModel                        string                 `json:"SubModel,omitempty"`
	SystemType                      string                 `json:"SystemType,omitempty"` //enum
	UUID                            string                 `json:"UUID,omitempty"`
	HostingRoles                    []string               `json:"HostingRoles,omitempty"` //enum
	PCIeDevices                     []PCIeDevice           `json:"PCIeDevices,omitempty"`
	PCIeFunctions                   []PCIeFunction         `json:"PCIeFunctions,omitempty"`
	Bios                            Bios                   `json:"Bios,omitempty"`
	Boot                            Boot                   `json:"Boot,omitempty"`
	EthernetInterfaces              EthernetInterfaces     `json:"EthernetInterfaces,omitempty"`
	HostedServices                  HostedServices         `json:"HostedServices,omitempty"`
	HostWatchdogTimer               HostWatchdogTimer      `json:"HostWatchdogTimer,omitempty"`
	Links                           Links                  `json:"Links,omitempty"`
	LogServices                     LogServices            `json:"LogServices,omitempty"`
	Memory                          Memory                 `json:"Memory,omitempty"`
	MemoryDomains                   MemoryDomains          `json:"MemoryDomains,omitempty"`
	MemorySummary                   MemorySummary          `json:"MemorySummary,omitempty"`
	NetworkInterfaces               NetworkInterfaces      `json:"NetworkInterfaces,omitempty"`
	Processors                      Processors             `json:"Processors,omitempty"`
	ProcessorSummary                ProcessorSummary       `json:"ProcessorSummary,omitempty"`
	Redundancy                      []Redundancy           `json:"Redundancy,omitempty"`
	SecureBoot                      SecureBoot             `json:"SecureBoot,omitempty"`
	SimpleStorage                   SimpleStorage          `json:"SimpleStorage,omitempty"`
	Status                          Status                 `json:"Status,omitempty"`
	Storage                         Storage                `json:"Storage,omitempty"`
	TrustedModules                  []TrustedModule        `json:"TrustedModules,omitempty"`
	Oem                             Oem                    `json:"Oem,omitempty"`
	PCIeDevicesCount                int                    `json:"PCIeDevices@odata.count,omitempty"`
	IdlePowerSaver                  *IdlePowerSaver        `json:"IdlePowerSaver,omitempty"`
	KeyManagement                   KeyManagement          `json:"KeyManagement,omitempty"`
	BootProgress                    BootProgress           `json:"BootProgress,omitempty"`
	Certificates                    Certificates           `json:"Certificates"`
	FabricAdapters                  *Link                  `json:"FabricAdapters,omitempty"`
	GraphicalConsole                *GraphicalConsole      `json:"GraphicalConsole,omitempty"`
	GraphicsControllers             *Link                  `json:"GraphicsControllers,omitempty"`
	LastResetTime                   string                 `json:"LastResetTime,omitempty"`
	LocationIndicatorActive         bool                   `json:"LocationIndicatorActive,omitempty"`
	Measurements                    []*Link                `json:"Measurements,omitempty"` //Deprecated in version v1.17.0
	PCIeFunctionsCount              int                    `json:"PCIeFunctions@odata.count,omitempty"`
	PowerCycleDelaySeconds          float32                `json:"PowerCycleDelaySeconds,omitempty"`
	PowerMode                       string                 `json:"PowerMode,omitempty"` //enum
	PowerOffDelaySeconds            float32                `json:"PowerOffDelaySeconds,omitempty"`
	PowerOnDelaySeconds             float32                `json:"PowerOnDelaySeconds,omitempty"`
	RedundancyCount                 int                    `json:"Redundancy@odata.count,omitempty"`
	SerialConsole                   SerialConsole          `json:"SerialConsole,omitempty"`
	USBControllers                  *Link                  `json:"USBControllers,omitempty"`
	VirtualMedia                    *Link                  `json:"VirtualMedia,omitempty"`
	VirtualMediaConfig              *VirtualMediaConfig    `json:"VirtualMediaConfig,omitempty"`
	OffloadedNetworkDeviceFunctions []*Link                `json:"OffloadedNetworkDeviceFunctions,omitempty"`
	LastBootTimeSeconds             int                    `json:"LastBootTimeSeconds,omitempty"`
	ManufacturingMode               bool                   `json:"ManufacturingMode,omitempty"`
	Composition                     *Composition           `json:"Composition,omitempty"`
}

type HostGraphicalConsole struct {
	ConnectedTypesSupported []string `json:"ConnectedTypesSupported,omitempty"` //enum
	MaxConcurrentSessions   int      `json:"MaxConcurrentSessions,omitempty"`
	Port                    int      `json:"Port,omitempty"`
	ServiceEnabled          bool     `json:"ServiceEnabled,omitempty"`
}

type HostSerialConsole struct {
	IPMI                  *SerialConsoleProtocol `json:"IPMI,omitempty"`
	MaxConcurrentSessions int                    `json:"MaxConcurrentSessions,omitempty"`
	SSH                   *SerialConsoleProtocol `json:"SSH,omitempty"`
	Telnet                *SerialConsoleProtocol `json:"Telnet,omitempty"`
}

// ComputerSystemActions redfish structure
type ComputerSystemActions struct {
	AddResourceBlock    *AddResourceBlock    `json:"AddResourceBlock,omitempty"`
	RemoveResourceBlock *RemoveResourceBlock `json:"RemoveResourceBlock,omitempty"`
	Reset               *Reset               `json:"Reset,omitempty"`
	SetDefaultBootOrder *SetDefaultBootOrder `json:"SetDefaultBootOrder,omitempty"`
	Oem                 *OemActions          `json:"Oem,omitempty"`
}

// RemoveResourceBlock redfish structure
// This action removes a resource block from a system
// This action shall remove a resource block from a system
type RemoveResourceBlock struct {
	Target string `json:"Target,omitempty"`
	Title  string `json:"Title,omitempty"`
}

// Reset redfish Structure
// This action resets the system.
// This action shall reset the system represented by the resource.
//
//	For systems that implement ACPI Power Button functionality, the PushPowerButton value
//
// shall perform or emulate an ACPI Power Button Push, and the ForceOff value shall perform
// an ACPI Power Button Override, commonly known as a four-second hold of the power button
type Reset struct {
	Target string `json:"Target,omitempty"`
	Title  string `json:"Title,omitempty"`
}

// SetDefaultBootOrder redfish Structure
// This action sets the BootOrder to the default settings
// This action shall set the BootOrder array to the default settings
type SetDefaultBootOrder struct {
	Target string `json:"Target,omitempty"`
	Title  string `json:"Title,omitempty"`
}

// LogEntry redfish structure
type LogEntry struct {
	Oid                     string   `json:"@odata.id"`
	Ocontext                string   `json:"@odata.context,omitempty"`
	Oetag                   string   `json:"@odata.etag,omitempty"`
	Otype                   string   `json:"@odata.type"`
	Description             string   `json:"description,omitempty"`
	Actions                 *Actions `json:"Actions,omitempty"`
	ID                      string   `json:"ID"`
	AdditionalDataSizeBytes int      `json:"AdditionalDataSizeBytes,omitempty"`
	AdditionalDataURI       string   `json:"AdditionalDataURI,omitempty"`
	Created                 string   `json:"Created,omitempty"`
	MessageArgs             string   `json:"MessageArgs,omitempty"`
	MessageID               string   `json:"MessageID,omitempty"`
	Name                    string   `json:"Name"`
	OEMDiagnosticDataType   string   `json:"OEMDiagnosticDataType,omitempty"`
	Oem                     *Oem     `json:"Oem,omitempty"`
	Resolved                bool     `json:"Resolved,omitempty"`
	Resolution              string   `json:"Resolution,omitempty"`
	Persistency             bool     `json:"Persistency,omitempty"`
	OverflowErrorCount      int      `json:"OverflowErrorCount,omitempty"`
	Originator              string   `json:"Originator,omitempty"`
	OemSensorType           string   `json:"OemSensorType,omitempty"`
	OemRecordFormat         string   `json:"OemRecordFormat,omitempty"`
}

// Composition ...
// Information about the composition capabilities and state of a computer system.
// This type shall contain information about the composition capabilities and state of a computer system
type Composition struct {
	UseCases []string `json:"UseCases,omitempty"` //enum
}

// VirtualMediaConfig redfish structure
// The information about virtual media service for this system
// This type shall describe a virtual media service service for a computer system
type VirtualMediaConfig struct {
	Port           int  `json:"Port,omitempty"`
	ServiceEnabled bool `json:"ServiceEnabled,omitempty"`
}

// SerialConsole redfish structure
// The information about the serial console services that this system provides.
// This type shall describe the serial console services for a computer system
type SerialConsole struct {
	IPMI                  interface{} `json:"IPMI,omitempty"`
	MaxConcurrentSessions int         `json:"MaxConcurrentSessions,omitempty"`
	SSH                   interface{} `json:"SSH,omitempty"`
	Telnet                interface{} `json:"Telnet,omitempty"`
}

// BootProgress redfish structure
// This object describes the last boot progress state
// This object shall contain the last boot progress state and time
type BootProgress struct {
	LastBootTimeSeconds float32 `json:"LastBootTimeSeconds,omitempty"`
	LastState           string  `json:"LastState,omitempty"` //enum
	LastStateTime       string  `json:"LastStateTime,omitempty"`
	Oem                 *Oem    `json:"Oem,omitempty"`
	OemLastState        string  `json:"OemLastState,omitempty"`
}

// KeyManagement redfish structure
// The key management settings of a computer system
// This object shall contain the key management settings of a computer system
type KeyManagement struct {
	KMIPCertificates *KMIPCertificates `json:"KMIPCertificates,omitempty"`
	KMIPServers      []*KMIPServers    `json:"KMIPServers,omitempty"`
}

// KMIPCertificates redfish structure
type KMIPCertificates struct {
	Oid string `json:"@odata.id"`
}

// KMIPServers redfish structure
// The KMIP server settings for a computer system
// This object shall contain the KMIP server settings for a computer system
type KMIPServers struct {
	Address       string `json:"Address,omitempty"`
	CacheDuration string `json:"CacheDuration,omitempty"`
	CachePolicy   string `json:"CachePolicy,omitempty"` //enum
	Password      string `json:"Password,omitempty"`
	Port          int    `json:"Port,omitempty"`
	Username      string `json:"Username,omitempty"`
}

// IdlePowerSaver redfish structure
// The idle power saver settings of a computer system
// This object shall contain the idle power saver settings of a computer system
type IdlePowerSaver struct {
	Enabled                 bool    `json:"Enabled,omitempty"`
	EnterDwellTimeSeconds   int     `json:"EnterDwellTimeSeconds,omitempty"`
	EnterUtilizationPercent float32 `json:"EnterUtilizationPercent,omitempty"`
	ExitDwellTimeSeconds    int     `json:"ExitDwellTimeSeconds,omitempty"`
	ExitUtilizationPercent  float32 `json:"ExitUtilizationPercent,omitempty"`
}

// Bios redfish structure
type Bios struct {
	Oid         string `json:"@odata.id"`
	Ocontext    string `json:"@odata.context,omitempty"`
	Oetag       string `json:"@odata.etag,omitempty"`
	Otype       string `json:"@odata.type"`
	Description string `json:"description,omitempty"`
	ID          string `json:"ID"`
	Name        string `json:"Name"`
	Oem         Oem    `json:"Oem,omitempty"`
	/*The reference to the Attribute Registry that lists the metadata describing the
	BIOS attribute settings in this resource.
	*/
	AttributeRegistry          string                 `json:"AttributeRegistry,omitempty"` // read-only (null)
	Attributes                 map[string]interface{} `json:"Attributes,omitempty"`        // object
	Links                      Links                  `json:"Links,omitempty"`
	ResetBiosToDefaultsPending bool                   `json:"ResetBiosToDefaultsPending,omitempty"`
}

// Boot redfish structure
// The boot information for this resource.
// This type shall contain properties that describe boot information for a system
type Boot struct {
	AliasBootOrder                  []string     `json:"AliasBootOrder,omitempty"` //enum
	BootNext                        string       `json:"BootNext,omitempty"`
	BootOptions                     *BootOptions `json:"BootOptions,omitempty"`
	BootOrder                       []string     `json:"BootOrder,omitempty"`
	BootOrderPropertySelection      string       `json:"BootOrderPropertySelection,omitempty"` //enum
	BootSourceOverrideEnabled       string       `json:"BootSourceOverrideEnabled,omitempty"`  //enum
	BootSourceOverrideMode          string       `json:"BootSourceOverrideMode,omitempty"`     //enum
	BootSourceOverrideTarget        string       `json:"BootSourceOverrideTarget,omitempty"`   //enum
	Certificates                    Certificates `json:"Certificates,omitempty"`
	UefiTargetBootSourceOverride    string       `json:"UefiTargetBootSourceOverride,omitempty"`
	AutomaticRetryAttempts          int          `json:"AutomaticRetryAttempts,omitempty"`
	AutomaticRetryConfig            string       `json:"AutomaticRetryConfig,omitempty"`
	HTTPBootURI                     string       `json:"HttpBootUri,omitempty"`
	RemainingAutomaticRetryAttempts int          `json:"RemainingAutomaticRetryAttempts,omitempty"`
	StopBootOnFault                 string       `json:"StopBootOnFault,omitempty"`             //enum
	TrustedModuleRequiredToBoot     string       `json:"TrustedModuleRequiredToBoot,omitempty"` //enum
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
	Otype                  string              `json:"@odata.type"`
	Description            string              `json:"description,omitempty"`
	ID                     string              `json:"ID"`
	Name                   string              `json:"Name"`
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

// DHCPv4 in place object
type DHCPv4 struct {
	DHCPEnabled     bool   `json:"DHCPEnabled"`
	FallbackAddress string `json:"FallbackAddress"` //enum
	UseDNSServers   bool   `json:"UseDNSServers"`
	UseDomainName   bool   `json:"UseDomainName"`
	UseGateway      bool   `json:"UseGateway"`
	UseNTPServers   bool   `json:"UseNTPServers"`
	UseStaticRoutes bool   `json:"UseStaticRoutes"`
}

// DHCPv6 in place object
type DHCPv6 struct {
	OperatingMode  string `json:"OperatingMode"` //enum
	UseDNSServers  bool   `json:"UseDNSServers"`
	UseDomainName  bool   `json:"UseDomainName"`
	UseNTPServers  bool   `json:"UseNTPServers"`
	UseRapidCommit bool   `json:"UseRapidCommit"`
}

// IPv4Address in place object
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

// IPv6StaticAddress in place object
type IPv6StaticAddress struct {
	Address      string `json:"Address"`
	Oem          Oem    `json:"Oem"`
	PrefixLength string `json:"PrefixLength"`
}

// IPv6AddressPolicy in place object
type IPv6AddressPolicy struct {
	Label      int    `json:"Label"`
	Precedence int    `json:"Precedence"`
	Prefix     string `json:"Prefix"`
}

// StatelessAddressAutoConfig in place object
type StatelessAddressAutoConfig struct {
	IPv4AutoConfigEnabled bool `json:"IPv4AutoConfigEnabled"`
	IPv6AutoConfigEnabled bool `json:"IPv6AutoConfigEnabled"`
}

// VLAN redfish structure
type VLAN struct {
	Oid string `json:"@odata.id"`
}

// VLANs redfish structure
type VLANs struct {
	Oid string `json:"@odata.id"`
}

// HostedServices redfish structure
// The services that might be running or installed on the system
// This type shall describe services that a computer system supports
type HostedServices struct {
	Oem             Oem             `json:"Oem"`
	StorageServices StorageServices `json:"StorageServices"`
}

// HostWatchdogTimer redfish structure
// This type describes the host watchdog timer functionality for this system
// This type shall contain properties that describe the host watchdog timer functionality for this ComputerSystem
type HostWatchdogTimer struct {
	FunctionEnabled bool   `json:"FunctionEnabled"`
	Oem             Oem    `json:"Oem"`
	Status          Status `json:"Status"`
	TimeoutAction   string `json:"TimeoutAction"` //enum
	WarningAction   string `json:"WarningAction"` //enum
}

// Memory redfish structure
type Memory struct {
	Oid                                     string                  `json:"@odata.id"`
	Ocontext                                string                  `json:"@odata.context,omitempty"`
	Oetag                                   string                  `json:"@odata.etag,omitempty"`
	Otype                                   string                  `json:"@odata.type"`
	Description                             string                  `json:"description,omitempty"`
	ID                                      string                  `json:"ID"`
	Name                                    string                  `json:"Name"`
	Oem                                     Oem                     `json:"Oem,omitempty"`
	AllocationAlignmentMiB                  int                     `json:"AllocationAlignmentMiB,omitempty"`
	AllocationIncrementMiB                  int                     `json:"AllocationIncrementMiB,omitempty"`
	AllowedSpeedsMHz                        []int                   `json:"AllowedSpeedsMHz,omitempty"`
	Assembly                                Assembly                `json:"Assembly,omitempty"`
	BaseModuleType                          string                  `json:"BaseModuleType,omitempty"` //enum
	BusWidthBits                            int                     `json:"BusWidthBits,omitempty"`
	CacheSizeMiB                            int                     `json:"CacheSizeMiB,omitempty"`
	CapacityMiB                             int                     `json:"CapacityMiB,omitempty"`
	ConfigurationLocked                     bool                    `json:"ConfigurationLocked,omitempty"`
	DataWidthBits                           int                     `json:"DataWidthBits,omitempty"`
	DeviceID                                string                  `json:"DeviceID,omitempty"`
	DeviceLocator                           string                  `json:"DeviceLocator,omitempty"`
	ErrorCorrection                         string                  `json:"ErrorCorrection,omitempty"` //enum
	FirmwareAPIVersion                      string                  `json:"FirmwareApiVersion,omitempty"`
	FirmwareRevision                        string                  `json:"FirmwareRevision,omitempty"`
	FunctionClasses                         []string                `json:"FunctionClasses,omitempty"`
	IsRankSpareEnabled                      bool                    `json:"IsRankSpareEnabled,omitempty"`
	IsSpareDeviceEnabled                    bool                    `json:"IsSpareDeviceEnabled,omitempty"`
	Links                                   Links                   `json:"Links,omitempty"`
	Location                                Location                `json:"Location,omitempty"`
	LogicalSizeMiB                          int                     `json:"LogicalSizeMiB,omitempty"`
	Manufacturer                            string                  `json:"Manufacturer,omitempty"`
	MaxTDPMilliWatts                        []int                   `json:"MaxTDPMilliWatts,omitempty"`
	MemoryDeviceType                        string                  `json:"MemoryDeviceType,omitempty"` //enum
	MemoryLocation                          MemoryLocation          `json:"MemoryLocation,omitempty"`
	MemoryMedia                             []string                `json:"MemoryMedia,omitempty"` //enum
	MemorySubsystemControllerManufacturerID string                  `json:"MemorySubsystemControllerManufacturerID,omitempty"`
	MemorySubsystemControllerProductID      string                  `json:"MemorySubsystemControllerProductID,omitempty"`
	MemoryType                              string                  `json:"MemoryType,omitempty"` //enum
	Metrics                                 Metrics                 `json:"Metrics,omitempty"`
	ModuleManufacturerID                    string                  `json:"ModuleManufacturerID,omitempty"`
	ModuleProductID                         string                  `json:"ModuleProductID,omitempty"`
	NonVolatileSizeMiB                      int                     `json:"NonVolatileSizeMiB,omitempty"`
	OperatingMemoryModes                    []string                `json:"OperatingMemoryModes,omitempty"` //enum
	OperatingSpeedMhz                       int                     `json:"OperatingSpeedMhz,omitempty"`
	PartNumber                              string                  `json:"PartNumber,omitempty"`
	PersistentRegionNumberLimit             int                     `json:"PersistentRegionNumberLimit,omitempty"`
	PersistentRegionSizeLimitMiB            int                     `json:"PersistentRegionSizeLimitMiB,omitempty"`
	PersistentRegionSizeMaxMiB              int                     `json:"PersistentRegionSizeMaxMiB,omitempty"`
	PowerManagementPolicy                   PowerManagementPolicy   `json:"PowerManagementPolicy,omitempty"`
	RankCount                               int                     `json:"RankCount,omitempty"`
	Regions                                 []Region                `json:"Regions,omitempty"`
	SecurityCapabilities                    SecurityCapabilities    `json:"SecurityCapabilities,omitempty"`
	SecurityState                           string                  `json:"SecurityState,omitempty"` //enum
	SerialNumber                            string                  `json:"SerialNumber,omitempty"`
	SpareDeviceCount                        int                     `json:"SpareDeviceCount,omitempty"`
	Status                                  Status                  `json:"Status,omitempty"`
	SubsystemDeviceID                       string                  `json:"SubsystemDeviceID,omitempty"`
	SubsystemVendorID                       string                  `json:"SubsystemVendorID,omitempty"`
	VendorID                                string                  `json:"VendorID,omitempty"`
	VolatileRegionNumberLimit               int                     `json:"VolatileRegionNumberLimit,omitempty"`
	VolatileRegionSizeLimitMiB              int                     `json:"VolatileRegionSizeLimitMiB,omitempty"`
	VolatileRegionSizeMaxMiB                int                     `json:"VolatileRegionSizeMaxMiB,omitempty"`
	VolatileSizeMiB                         int                     `json:"VolatileSizeMiB,omitempty"`
	Log                                     *Link                   `json:"Log,omitempty"`
	OperatingSpeedRangeMHz                  *OperatingSpeedRangeMHz `json:"OperatingSpeedRangeMHz,omitempty"`
	Certificates                            Certificates            `json:"Certificates,omitempty"`
	Enabled                                 bool                    `json:"Enabled,omitempty"`
	EnvironmentMetrics                      *Link                   `json:"EnvironmentMetrics,omitempty"`
	LocationIndicatorActive                 bool                    `json:"LocationIndicatorActive,omitempty"`
	Measurements                            []*Link                 `json:"Measurements,omitempty"` // Deprecated in version v1.14.0
	Model                                   string                  `json:"Model,omitempty"`
	SparePartNumber                         string                  `json:"SparePartNumber,omitempty"`
	Batteries                               []*Link                 `json:"Batteries,omitempty"`
}

// OperatingSpeedRangeMHz redfish structure
type OperatingSpeedRangeMHz struct {
	AllowableMax           float32   `json:"AllowableMax,omitempty"`
	AllowableMin           float32   `json:"AllowableMin,omitempty"`
	AllowableNumericValues []float32 `json:"AllowableNumericValues,omitempty"`
	ControlMode            string    `json:"ControlMode,omitempty"`
	DataSourceURI          string    `json:"DataSourceUri,omitempty"`
	Reading                float32   `json:"Reading,omitempty"`
	ReadingUnits           string    `json:"ReadingUnits,omitempty"`
	SettingMax             float32   `json:"SettingMax,omitempty"`
	SettingMin             float32   `json:"SettingMin,omitempty"`
}

// MemoryLocation in place object
type MemoryLocation struct {
	Channel          int `json:"Channel,omitempty"`
	MemoryController int `json:"MemoryController"`
	Slot             int `json:"Slot"`
	Socket           int `json:"Socket"`
}

// PowerManagementPolicy in place object
type PowerManagementPolicy struct {
	AveragePowerBudgetMilliWatts int  `json:"AveragePowerBudgetMilliWatts"`
	MaxTDPMilliWatts             int  `json:"MaxTDPMilliWatts"`
	PeakPowerBudgetMilliWatts    int  `json:"PeakPowerBudgetMilliWatts"`
	PolicyEnabled                bool `json:"PolicyEnabled"`
}

// Region in place object
type Region struct {
	MemoryClassification string `json:"MemoryClassification"` //enum
	OffsetMiB            int    `json:"OffsetMiB"`
	PassphraseEnabled    bool   `json:"PassphraseEnabled"`
	PassphraseState      bool   `json:"PassphraseState"`
	RegionID             string `json:"RegionID"`
	SizeMiB              int    `json:"SizeMiB"`
}

// SecurityCapabilities in place object
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
	ID                        string                   `json:"ID,omitempty"`
	Name                      string                   `json:"Name,omitempty"`
	Oem                       Oem                      `json:"Oem,omitempty"`
	AllowsBlockProvisioning   bool                     `json:"AllowsBlockProvisioning,omitempty"`
	AllowsMemoryChunkCreation bool                     `json:"AllowsMemoryChunkCreation,omitempty"`
	AllowsMirroring           bool                     `json:"AllowsMirroring,omitempty"`
	AllowsSparing             bool                     `json:"AllowsSparing,omitempty"`
	InterleavableMemorySets   []InterleavableMemorySet `json:"InterleavableMemorySets,omitempty"`
	MemoryChunks              MemoryChunks             `json:"MemoryChunks,omitempty"`
	Actions                   *OemActions              `json:"Actions,omitempty"`
	Links                     Link                     `json:"Links,omitempty"`
}

// InterleavableMemorySet in place object
type InterleavableMemorySet struct {
	MemorySet      []Memory `json:"MemorySet,omitempty"`
	MemorySetCount int      `json:"MemorySet@odata.count,omitempty"`
}

// MemoryChunks redfish structure
type MemoryChunks struct {
	Oid string `json:"@odata.id"`
}

// MemorySummary in place object
// The memory of the system in general detail
// This type shall contain properties that describe the central memory for a system
type MemorySummary struct {
	MemoryMirroring                string  `json:"MemoryMirroring"` //enum
	Metrics                        Metrics `json:"Metrics,omitempty"`
	TotalSystemMemoryGiB           int     `json:"TotalSystemMemoryGiB"`
	TotalSystemPersistentMemoryGiB int     `json:"TotalSystemPersistentMemoryGiB"`
	Status                         Status  `json:"Status"` //deprecated
}

// SerialConsoleProtocol redfish structure
// The information about a serial console service that this system provides
// This type shall describe a serial console service for a computer system
type SerialConsoleProtocol struct {
	ConsoleEntryCommand   string `json:"ConsoleEntryCommand,omitempty"`
	HotKeySequenceDisplay string `json:"HotKeySequenceDisplay,omitempty"`
	Port                  int    `json:"Port,omitempty"`
	ServiceEnabled        bool   `json:"ServiceEnabled,omitempty"`
	SharedWithManagerCLI  bool   `json:"SharedWithManagerCLI,omitempty"`
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
	ID                     string                 `json:"ID,omitempty"`
	Name                   string                 `json:"Name,omitempty"`
	Oem                    Oem                    `json:"Oem,omitempty"`
	Links                  Links                  `json:"Links,omitempty"`
	NetworkDeviceFunctions NetworkDeviceFunctions `json:"NetworkDeviceFunctions,omitempty"`
	NetworkPorts           NetworkPorts           `json:"NetworkPorts,omitempty"`
	Status                 Status                 `json:"Status,omitempty"`
	Actions                *OemActions            `json:"Actions,omitempty"`
	Ports                  *Link                  `json:"Ports,omitempty"`
}

// NetworkDeviceFunctions redfish structure
type NetworkDeviceFunctions struct {
	Oid string `json:"@odata.id"`
}

// NetworkPorts redfish structure
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
	Oid                string         `json:"@odata.id"`
	Ocontext           string         `json:"@odata.context,omitempty"`
	Oetag              string         `json:"@odata.etag,omitempty"`
	Otype              string         `json:"@odata.type"`
	Description        string         `json:"description,omitempty"`
	Id                 string         `json:"Id"`
	Name               string         `json:"Name"`
	Oem                Oem            `json:"Oem,omitempty"`
	Assembly           *Assembly      `json:"Assembly,omitempty"`
	AssetTag           string         `json:"AssetTag,omitempty"`
	DeviceType         string         `json:"DeviceType,omitempty"` //enum
	FirmwareVersion    string         `json:"FirmwareVersion,omitempty"`
	Links              *Links         `json:"Links,omitempty"`
	Manufacturer       string         `json:"Manufacturer,omitempty"`
	Model              string         `json:"Model,omitempty"`
	PartNumber         string         `json:"PartNumber,omitempty"`
	PCIeInterface      *PCIeInterface `json:"PCIeInterface,omitempty"`
	SerialNumber       string         `json:"SerialNumber,omitempty"`
	SKU                string         `json:"SKU,omitempty"`
	Status             *Status        `json:"Status,omitempty"`
	Actions            *OemActions    `json:"Actions,omitempty"`
	EnvironmentMetrics *Link          `json:"EnvironmentMetrics,omitempty"`
	PCIeFunctions      *Link          `json:"PCIeFunctions,omitempty"`
	ReadyToRemove      bool           `json:"ReadyToRemove,omitempty"`
	SparePartNumber    string         `json:"SparePartNumber,omitempty"`
	UUID               string         `json:"UUID,omitempty"`
	Slot               *Slot          `json:"Slot,omitempty"`
	PCIeErrors         *PCIeErrors    `json:"PCIeErrors,omitempty"`
}

// Slot Information about the slot for this PCIe device.
type Slot struct {
	LaneSplitting string `json:"LaneSplitting,omitempty"`
	Lanes         int    `json:"Lanes,omitempty"`
	PCIeType      string `json:"PCIeType,omitempty"`
	SlotType      string `json:"SlotType,omitempty"`
}

// PCIeErrors - The PCIe errors associated with this device
type PCIeErrors struct {
	CorrectableErrorCount int `json:"CorrectableErrorCount,omitempty"`
	FatalErrorCount       int `json:"FatalErrorCount,omitempty"`
	L0ToRecoveryCount     int `json:"L0ToRecoveryCount,omitempty"`
	NAKReceivedCount      int `json:"NAKReceivedCount,omitempty"`
	NAKSentCount          int `json:"NAKSentCount,omitempty"`
	NonFatalErrorCount    int `json:"NonFatalErrorCount,omitempty"`
	ReplayCount           int `json:"ReplayCount,omitempty"`
	ReplayRolloverCount   int `json:"ReplayRolloverCount,omitempty"`
}

// PCIeInterface in place object
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
	Oid               string      `json:"@odata.id"`
	Ocontext          string      `json:"@odata.context,omitempty"`
	Oetag             string      `json:"@odata.etag,omitempty"`
	Otype             string      `json:"@odata.type"`
	Description       string      `json:"description,omitempty"`
	ID                string      `json:"ID"`
	Name              string      `json:"Name"`
	Oem               Oem         `json:"Oem,omitempty"`
	ClassCode         string      `json:"ClassCode,omitempty"`
	DeviceClass       string      `json:"DeviceClass,omitempty"` //enum
	DeviceID          string      `json:"DeviceID,omitempty"`
	FunctionID        int         `json:"FunctionID,omitempty"`
	FunctionType      string      `json:"FunctionType,omitempty"` //enum
	Links             *Links      `json:"Links,omitempty"`
	RevisionID        string      `json:"RevisionID,omitempty"`
	Status            *Status     `json:"Status,omitempty"`
	SubsystemID       string      `json:"SubsystemID,omitempty"`
	SubsystemVendorID string      `json:"SubsystemVendorID,omitempty"`
	VendorID          string      `json:"VendorID,omitempty"`
	Actions           *OemActions `json:"Actions,omitempty"`
	Enabled           bool        `json:"Enabled,omitempty"`
}

/*
Processors 1.5.0

This resource shall be used to represent a single processor contained within a system.
URIs:
/redfish/v1/Systems/{ComputerSystemId}/Processors/{ProcessorId}
/redfish/v1/Systems/{ComputerSystemId}/Processors/{ProcessorId}/SubProcessors/{ProcessorId2}
*/
type Processors struct {
	Oid                        string                      `json:"@odata.id"`
	Ocontext                   string                      `json:"@odata.context,omitempty"`
	Oetag                      string                      `json:"@odata.etag,omitempty"`
	Otype                      string                      `json:"@odata.type"`
	Description                string                      `json:"description,omitempty"`
	ID                         string                      `json:"ID"`
	Name                       string                      `json:"Name"`
	Oem                        Oem                         `json:"Oem,omitempty"`
	AccelerationFunctions      AccelerationFunctions       `json:"AccelerationFunctions,omitempty"`
	Assembly                   Assembly                    `json:"Assembly,omitempty"`
	FPGA                       FPGA                        `json:"FPGA,omitempty"`
	InstructionSet             string                      `json:"InstructionSet,omitempty"` //enum
	Links                      Links                       `json:"Links,omitempty"`
	Location                   Location                    `json:"Location,omitempty"`
	Manufacturer               string                      `json:"Manufacturer,omitempty"`
	MaxSpeedMHz                int                         `json:"MaxSpeedMHz,omitempty"`
	MaxTDPWatts                int                         `json:"MaxTDPWatts,omitempty"`
	Metrics                    Metrics                     `json:"Metrics,omitempty"`
	Model                      string                      `json:"Model,omitempty"`
	ProcessorArchitecture      string                      `json:"ProcessorArchitecture,omitempty"` //enum
	ProcessorID                ProcessorID                 `json:"ProcessorID,omitempty"`
	ProcessorMemory            []ProcessorMemory           `json:"ProcessorMemory,omitempty"`
	ProcessorType              string                      `json:"ProcessorType,omitempty"` //enum
	Socket                     string                      `json:"Socket,omitempty"`
	Status                     Status                      `json:"Status,omitempty"`
	SubProcessors              SubProcessors               `json:"SubProcessors,omitempty"`
	TDPWatts                   int                         `json:"TDPWatts,omitempty"`
	TotalCores                 int                         `json:"TotalCores,omitempty"`
	TotalEnabledCores          int                         `json:"TotalEnabledCores,omitempty"`
	TotalThreads               int                         `json:"TotalThreads,omitempty"`
	UUID                       string                      `json:"UUID,omitempty"`
	OperatingSpeedRangeMHz     *OperatingSpeedRangeMHz     `json:"OperatingSpeedRangeMHz,omitempty"`
	Ports                      *Link                       `json:"Ports,omitempty"`
	Actions                    *OemActions                 `json:"Actions,omitempty"`
	BaseSpeedMHz               int                         `json:"BaseSpeedMHz,omitempty"`
	BaseSpeedPriorityState     string                      `json:"BaseSpeedPriorityState,omitempty"`
	Certificates               Certificates                `json:"Certificates,omitempty"`
	Enabled                    bool                        `json:"Enabled,omitempty"`
	EnvironmentMetrics         *Link                       `json:"EnvironmentMetrics,omitempty"`
	FirmwareVersion            string                      `json:"FirmwareVersion,omitempty"`
	HighSpeedCoreIDs           []int                       `json:"HighSpeedCoreIDs,omitempty"`
	LocationIndicatorActive    bool                        `json:"LocationIndicatorActive,omitempty"`
	Measurements               []*Link                     `json:"Measurements,omitempty"`
	MemorySummary              *MemorySummary              `json:"MemorySummary,omitempty"`
	MinSpeedMHz                int                         `json:"MinSpeedMHz,omitempty"`
	OperatingConfigs           *Link                       `json:"OperatingConfigs,omitempty"`
	OperatingSpeedMHz          int                         `json:"OperatingSpeedMHz,omitempty"`
	PartNumber                 string                      `json:"PartNumber,omitempty"`
	SerialNumber               string                      `json:"SerialNumber,omitempty"`
	SparePartNumber            string                      `json:"SparePartNumber,omitempty"`
	SpeedLimitMHz              int                         `json:"SpeedLimitMHz,omitempty"`
	SpeedLocked                bool                        `json:"SpeedLocked,omitempty"`
	SystemInterface            SystemInterface             `json:"SystemInterface,omitempty"`
	TurboState                 string                      `json:"TurboState,omitempty"`
	Version                    string                      `json:"Version,omitempty"`
	AdditionalFirmwareVersions *AdditionalFirmwareVersions `json:"AdditionalFirmwareVersions,omitempty"`
}

// SystemInterface redfish structure
type MemoryMetrics struct {
	Oid                           string         `json:"@odata.id"`
	Ocontext                      string         `json:"@odata.context,omitempty"`
	Oetag                         string         `json:"@odata.etag,omitempty"`
	Otype                         string         `json:"@odata.type"`
	Actions                       *OemActions    `json:"Actions,omitempty"`
	BandwidthPercent              int            `json:"BandwidthPercent,omitempty"`
	BlockSizeBytes                int            `json:"BlockSizeBytes,omitempty"`
	CXL                           CXL            `json:"CXL,omitempty"`
	CapacityUtilizationPercent    int            `json:"CapacityUtilizationPercent,omitempty"`
	CorrectedPersistentErrorCount int            `json:"CorrectedPersistentErrorCount,omitempty"`
	CorrectedVolatileErrorCount   int            `json:"CorrectedVolatileErrorCount,omitempty"`
	CurrentPeriod                 *CurrentPeriod `json:"CurrentPeriod,omitempty"`
	Description                   string         `json:"description,omitempty"`
	DirtyShutdownCount            int            `json:"DirtyShutdownCount,omitempty"`
	ID                            string         `json:"ID"`
	Name                          string         `json:"Name"`
	Oem                           Oem            `json:"Oem,omitempty"`
	OperatingSpeedMHz             int            `json:"OperatingSpeedMHz,omitempty"`
}

// SystemInterface redfish structure
type SystemInterface struct {
	Ethernet      Ethernet      `json:"Ethernet,omitempty"`
	InterfaceType string        `json:"InterfaceType,omitempty"`
	PCIe          PCIeInterface `json:"PCIe,omitempty"`
}

// MemorySummaryDetails in place object
type MemorySummaryDetails struct {
	ECCModeEnabled     bool    `json:"ECCModeEnabled,omitempty"`
	Metrics            Metrics `json:"Metrics,omitempty"`
	TotalCacheSizeMiB  int     `json:"TotalCacheSizeMiB,omitempty"`
	TotalMemorySizeMiB int     `json:"TotalMemorySizeMiB,omitempty"`
}

// AccelerationFunctions redfish structure
type AccelerationFunctions struct {
	Oid string `json:"@odata.id"`
}

// FPGA in place object
type FPGA struct {
	ExternalInterfaces   []HostInterface       `json:"ExternalInterfaces"`
	FirmwareID           string                `json:"FirmwareID"`
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

// HostInterface in place object
type HostInterface struct {
	Ethernet      Ethernet      `json:"Ethernet"`
	InterfaceType string        `json:"InterfaceType"` //enum
	PCIe          PCIeInterface `json:"PCIe"`
}

// Ethernet in place object
type Ethernet struct {
	MaxLanes     int `json:"MaxLanes"`
	MaxSpeedMbps int `json:"MaxSpeedMbps"`
	Oem          Oem `json:"Oem"`
}

// ReconfigurationSlot in place object
type ReconfigurationSlot struct {
	AccelerationFunction AccelerationFunction `json:"AccelerationFunction"`
	ProgrammableFromHost bool                 `json:"ProgrammableFromHost"`
	SlotID               string               `json:"SlotID"`
	UUID                 string               `json:"UUID"`
}

// AccelerationFunction redfish structure
type AccelerationFunction struct {
	Oid string `json:"@odata.id"`
}

// ProcessorID in place object
type ProcessorID struct {
	EffectiveFamily               string `json:"EffectiveFamily,omitempty"`
	EffectiveModel                string `json:"EffectiveModel,omitempty"`
	IdentificationRegisters       string `json:"IdentificationRegisters,omitempty"`
	MicrocodeInfo                 string `json:"MicrocodeInfo,omitempty"`
	Step                          string `json:"Step,omitempty"`
	VendorID                      string `json:"VendorId,omitempty"`
	ProtectedIdentificationNumber string `json:"ProtectedIdentificationNumber,omitempty"`
}

// ProcessorMemory in place object
type ProcessorMemory struct {
	CapacityMiB      int    `json:"CapacityMiB"`
	IntegratedMemory bool   `json:"IntegratedMemory"`
	MemoryType       string `json:"MemoryType"` //enum
	SpeedMHz         int    `json:"SpeedMHz"`
}

// SubProcessors redfish structure
type SubProcessors struct {
	Oid string `json:"@odata.id"`
}

// ProcessorSummary redfish structure
// The central processors of the system in general detail
// This type shall contain properties that describe the central processors for a system.
//
//	Processors described by this type shall be limited to the processors that execute system code,
//
// and shall not include processors used for offload functionality
type ProcessorSummary struct {
	CoreCount             int     `json:"CoreCount,omitempty"`
	Count                 int     `json:"Count"`
	LogicalProcessorCount int     `json:"LogicalProcessorCount"`
	Model                 string  `json:"Model"`
	Metrics               Metrics `json:"Metrics"`
	Status                Status  `json:"Status"` //deprecated
	ThreadingEnabled      bool    `json:"ThreadingEnabled,omitempty"`
}

// SecureBoot redfish structure
type SecureBoot struct {
	Oid                   string      `json:"@odata.id"`
	Ocontext              string      `json:"@odata.context,omitempty"`
	Oetag                 string      `json:"@odata.etag,omitempty"`
	Otype                 string      `json:"@odata.type"`
	Description           string      `json:"description,omitempty"`
	ID                    string      `json:"ID"`
	Name                  string      `json:"Name"`
	Oem                   Oem         `json:"Oem,omitempty"`
	SecureBootCurrentBoot string      `json:"SecureBootCurrentBoot,omitempty"`
	SecureBootEnable      bool        `json:"SecureBootEnable,omitempty"`
	SecureBootMode        string      `json:"SecureBootMode,omitempty"`
	Actions               *OemActions `json:"Actions,omitempty"`
	SecureBootDatabases   *Link       `json:"SecureBootDatabases,omitempty"`
}

// SecureBootDatabaseCollection This resource shall represent a resource collection
// of SecureBootDatabase instances for a Redfish implementation.
type SecureBootDatabaseCollection struct {
	Oid                  string `json:"@odata.id"`
	Ocontext             string `json:"@odata.context,omitempty"`
	Oetag                string `json:"@odata.etag,omitempty"`
	Otype                string `json:"@odata.type"`
	Description          string `json:"Description,omitempty"`
	Members              string `json:"Members"`
	MembersCount         int    `json:"Members@odata.count"`
	MembersODataNextLink string `json:"Members@odata.nextLink,omitempty"`
	Name                 string `json:"Name"`
	Oem                  *Oem   `json:"Oem,omitempty"`
}

// SecureBootDatabase This resource shall be used to represent
// a UEFI Secure Boot database for a Redfish implementation.
type SecureBootDatabase struct {
	Oid          string `json:"@odata.id"`
	Ocontext     string `json:"@odata.context,omitempty"`
	Oetag        string `json:"@odata.etag,omitempty"`
	Otype        string `json:"@odata.type"`
	Actions      string `json:"Actions,omitempty"`
	Certificates string `json:"Certificates,omitempty"`
	DatabaseID   string `json:"DatabaseID,omitempty"`
	Description  string `json:"Description,omitempty"`
	ID           string `json:"ID"`
	Name         string `json:"Name"`
	Oem          string `json:"Oem,omitempty"`
	Signatures   string `json:"Signatures,omitempty"`
}

// SimpleStorage redfish structure
type SimpleStorage struct {
	Oid            string      `json:"@odata.id"`
	Ocontext       string      `json:"@odata.context,omitempty"`
	Oetag          string      `json:"@odata.etag,omitempty"`
	Otype          string      `json:"@odata.type"`
	Description    string      `json:"description,omitempty"`
	ID             string      `json:"ID"`
	Name           string      `json:"Name"`
	Oem            Oem         `json:"Oem,omitempty"`
	Devices        []Device    `json:"Devices,omitempty"`
	Links          Link        `json:"Links,omitempty"`
	Status         Status      `json:"Status,omitempty"`
	UefiDevicePath string      `json:"UefiDevicePath,omitempty"`
	Actions        *OemActions `json:"Actions,omitempty"`
}

// Device in place object
type Device struct {
	CapacityBytes int    `json:"CapacityBytes,omitempty"`
	Manufacturer  string `json:"Manufacturer,omitempty"`
	Model         string `json:"Model,omitempty"`
	Name          string `json:"Name"`
	Oem           Oem    `json:"Oem,omitempty"`
	Status        Status `json:"Status,omitempty"`
}

// TrustedModule redfish structure
// The Trusted Module installed in the system
// This type shall describe a Trusted Module for a system
type TrustedModule struct {
	FirmwareVersion        string `json:"FirmwareVersion"`
	FirmwareVersion2       string `json:"FirmwareVersion2"`
	InterfaceType          string `json:"InterfaceType"`          //enum
	InterfaceTypeSelection string `json:"InterfaceTypeSelection"` //enum
	Oem                    Oem    `json:"Oem"`
	Status                 Status `json:"Status"`
}

// BootOptions redfish structure
type BootOptions struct {
	ODataContext         string   `json:"@odata.context,omitempty"`
	ODataEtag            string   `json:"@odata.etag,omitempty"`
	ODataID              string   `json:"@odata.id"`
	ODataType            string   `json:"@odata.type"`
	Description          string   `json:"Description,omitempty"`
	Members              []string `json:"Members"`
	MembersODataCount    int      `json:"Members@odata.count"`
	MembersODataNextLink string   `json:"Members@odata.nextLink,omitempty"`
	Name                 string   `json:"Name"`
	Oem                  *Oem     `json:"Oem,omitempty"`
}

// StorageServices redfish structure
type StorageServices struct {
	Oid                  string `json:"@odata.id"`
	Ocontext             string `json:"@odata.context,omitempty"`
	Oetag                string `json:"@odata.etag,omitempty"`
	Otype                string `json:"@odata.type"`
	Description          string `json:"Description,omitempty"`
	Members              string `json:"Members"`
	MembersCount         int    `json:"Members@odata.count"`
	MembersODataNextLink string `json:"Members@odata.nextLink,omitempty"`
	Name                 string `json:"Name"`
	Oem                  *Oem   `json:"Oem,omitempty"`
}

// Metrics redfish structure
type Metrics struct {
	Oid                          string             `json:"@odata.id"`
	Ocontext                     string             `json:"@odata.context,omitempty"`
	Oetag                        string             `json:"@odata.etag,omitempty"`
	Otype                        string             `json:"@odata.type"`
	Actions                      Actions            `json:"Actions,omitempty"`
	AverageFrequencyMHz          float32            `json:"AverageFrequencyMHz,omitempty"`
	BandwidthPercent             float32            `json:"BandwidthPercent,omitempty"`
	CacheMetrics                 *CacheMetrics      `json:"CacheMetrics,omitempty"`
	CacheMetricsTotal            *CacheMetricsTotal `json:"CacheMetricsTotal,omitempty"`
	ConsumedPowerWatt            float32            `json:"ConsumedPowerWatt,omitempty"`
	CorrectableCoreErrorCount    int                `json:"CorrectableCoreErrorCount,omitempty"`
	CorrectableOtherErrorCount   int                `json:"CorrectableOtherErrorCount,omitempty"`
	Description                  string             `json:"Description,omitempty"`
	FrequencyRatio               float32            `json:"FrequencyRatio,omitempty"`
	KernelPercent                float32            `json:"KernelPercent,omitempty"`
	LocalMemoryBandwidthBytes    int                `json:"LocalMemoryBandwidthBytes,omitempty"`
	ID                           string             `json:"ID"`
	Name                         string             `json:"Name"`
	Oem                          Oem                `json:"Oem,omitempty"`
	OperatingSpeedMHz            int                `json:"OperatingSpeedMHz,omitempty"`
	PowerLimitThrottleDuration   string             `json:"PowerLimitThrottleDuration,omitempty"`
	RemoteMemoryBandwidthBytes   int                `json:"RemoteMemoryBandwidthBytes,omitempty"`
	TemperatureCelsius           float32            `json:"TemperatureCelsius,omitempty"`
	ThermalLimitThrottleDuration string             `json:"ThermalLimitThrottleDuration,omitempty"`
	ThrottlingCelsius            float32            `json:"ThrottlingCelsius,omitempty"`
	UncorrectableCoreErrorCount  int                `json:"UncorrectableCoreErrorCount,omitempty"`
	UncorrectableOtherErrorCount int                `json:"UncorrectableOtherErrorCount,omitempty"`
	UserPercent                  float32            `json:"UserPercent,omitempty"`
}

// CacheMetrics redfish structure
type CacheMetrics struct {
	CacheMiss                 float32 `json:"CacheMiss,omitempty"`
	CacheMissesPerInstruction float32 `json:"CacheMissesPerInstruction,omitempty"`
	HitRatio                  float32 `json:"HitRatio,omitempty"`
	Level                     string  `json:"Level,omitempty"`
	OccupancyBytes            int     `json:"OccupancyBytes,omitempty"`
	OccupancyPercent          float32 `json:"OccupancyPercent,omitempty"`
}

// CacheMetricsTotal redfish structure
type CacheMetricsTotal struct {
	CurrentPeriod CurrentPeriod `json:"CurrentPeriod,omitempty"`
	LifeTime      LifeTime      `json:"LifeTime,omitempty"`
}

// CurrentPeriod redfish structure
type CurrentPeriod struct {
	CorrectableECCErrorCount   int `json:"CorrectableECCErrorCount,omitempty"`
	UncorrectableECCErrorCount int `json:"UncorrectableECCErrorCount,omitempty"`
}

// LifeTime redfish structure
type LifeTime struct {
	CorrectableECCErrorCount   int `json:"CorrectableECCErrorCount,omitempty"`
	UncorrectableECCErrorCount int `json:"UncorrectableECCErrorCount,omitempty"`
}

// SaveInMemory will create the ComputerSystem data in in-memory DB, with key as UUID
// Takes:
//
//	none as parameter, but recieves c of type *ComputerSystem as pointeer reciever impicitly.
//
// Returns:
//
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
