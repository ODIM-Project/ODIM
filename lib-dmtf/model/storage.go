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

// Storage redfish structure
type Storage struct {
	Oid                string                `json:"@odata.id"`
	ODataContext       string                `json:"@odata.context,omitempty"`
	ODataEtag          string                `json:"@odata.etag,omitempty"`
	ODataType          string                `json:"@odata.type,omitempty"`
	Description        string                `json:"Description,omitempty"`
	ID                 string                `json:"Id,omitempty"`
	Name               string                `json:"Name,omitempty"`
	Oem                *Oem                  `json:"Oem,omitempty"`
	Drives             []*Link               `json:"Drives,omitempty"`
	Links              *StorageLinks         `json:"Links,omitempty"`
	Redundancy         []*Redundancy         `json:"Redundancy,omitempty"`
	Status             *Status               `json:"Status,omitempty"`
	StorageControllers []*StorageControllers `json:"StorageControllers,omitempty"`
	Volumes            *Link                 `json:"Volumes,omitempty"`
	ConsistencyGroups  *Link                 `json:"ConsistencyGroups,omitempty"`
	Controllers        *Link                 `json:"Controllers,omitempty"`
	EndpointGroups     *Link                 `json:"EndpointGroups,omitempty"`
	FileSystems        *Link                 `json:"FileSystems,omitempty"`
	Identifiers        *Identifier           `json:"Identifiers,omitempty"`
	StorageGroups      *Link                 `json:"StorageGroups,omitempty"`
	StoragePools       *Link                 `json:"StoragePools,omitempty"`
}

//StorageControllers redfish structure
type StorageControllers struct {
	Oid                          string                    `json:"@odata.id"`
	AssetTag                     string                    `json:"AssetTag,omitempty"`
	FirmwareVersion              string                    `json:"FirmwareVersion,omitempty"`
	Manufacturer                 string                    `json:"Manufacturer,omitempty"`
	MemberID                     string                    `json:"MemberId,omitempty"`
	Model                        string                    `json:"Model,omitempty"`
	Name                         string                    `json:"Name,omitempty"`
	PartNumber                   string                    `json:"PartNumber,omitempty"`
	SerialNumber                 string                    `json:"SerialNumber,omitempty"`
	SKU                          string                    `json:"SKU,omitempty"`
	SpeedGbps                    int                       `json:"SpeedGbps,omitempty"`
	SupportedControllerProtocols []string                  `json:"SupportedControllerProtocols,omitempty"` //enum
	SupportedDeviceProtocols     []string                  `json:"SupportedDeviceProtocols,omitempty"`     //enum
	SupportedRAIDTypes           []string                  `json:"SupportedRAIDTypes,omitempty"`           //enum
	Actions                      *Actions                  `json:"Actions,omitempty"`
	Assembly                     *Assembly                 `json:"Assembly,omitempty"`
	CacheSummary                 *CacheSummary             `json:"CacheSummary,omitempty"`
	ControllerRates              *ControllerRates          `json:"ControllerRates,omitempty"`
	Identifiers                  []*Identifier             `json:"Identifiers,omitempty"`
	Links                        *StorageControllersLinks  `json:"Links,omitempty"`
	Location                     *Location                 `json:"Location,omitempty"`
	NVMeControllerProperties     *NVMeControllerProperties `json:"NVMeControllerProperties,omitempty"`
	Oem                          *Oem                      `json:"Oem,omitempty"`
	PCIeInterface                *PCIeInterface            `json:"PCIeInterface,omitempty"`
	Ports                        *Link                     `json:"Ports,omitempty"`
	Status                       *StorageStatus            `json:"Status,omitempty"`
}

//Actions redfish structure
type Actions struct {
}

// StorageLinks struct is for storage Links schema
type StorageLinks struct {
	Enclosures      []*Link `json:"Enclosures,omitempty"`
	SimpleStorage   *Link   `json:"SimpleStorage,omitempty"`
	StorageServices []*Link `json:"StorageServices,omitempty"`
	Oem             *Oem    `json:"Oem,omitempty"`
}

// StorageStatus struct is to define the status of the Storage
type StorageStatus struct {
	State  string `json:"State,omitempty"`
	Health string `json:"Health,omitempty"`
}

// CacheSummary struct is to define the CacheSummary of the Storage
type CacheSummary struct {
	PersistentCacheSizeMiB int            `json:"PersistentCacheSizeMiB,omitempty"`
	Status                 *StorageStatus `json:"Status,omitempty"`
	TotalCacheSizeMiB      int            `json:"TotalCacheSizeMiB,omitempty"`
}

// ControllerRates struct is to define the ControllerRates of the Storage
type ControllerRates struct {
	ConsistencyCheckRatePercent int `json:"ConsistencyCheckRatePercent,omitempty"`
	RebuildRatePercent          int `json:"RebuildRatePercent,omitempty"`
	TransformationRatePercent   int `json:"TransformationRatePercent,omitempty"`
}

// StorageControllersLinks struct is for smart storage Links schema
type StorageControllersLinks struct {
	PCIeFunctions   *Link   `json:"PCIeFunctions,omitempty"`
	Oem             *Oem    `json:"Oem,omitempty"`
	Endpoints       []*Link `json:"Enclosures,omitempty"`
	StorageServices []*Link `json:"StorageServices,omitempty"`
}

// NVMeControllerProperties struct is to define the NVMeControllerProperties of the Storage
type NVMeControllerProperties struct {
	ControllerType            string                     `json:"ControllerType,omitempty"`
	MaxQueueSize              int                        `json:"MaxQueueSize,omitempty"`
	NVMeVersion               string                     `json:"NVMeVersion,omitempty"`
	ANACharacteristics        *ANACharacteristics        `json:"ANACharacteristics,omitempty"`
	NVMeControllerAttributes  *NVMeControllerAttributes  `json:"NVMeControllerAttributes,omitempty"`
	NVMeSMARTCriticalWarnings *NVMeSMARTCriticalWarnings `json:"NVMeSMARTCriticalWarnings,omitempty"`
}

// ANACharacteristics struct is to define the ANACharacteristics of the Storage
type ANACharacteristics struct {
	AccessState string `json:"AccessState,omitempty"`
	Volume      *Link  `json:"Volume,omitempty"`
}

// NVMeControllerAttributes struct is to define the NVMeControllerAttributes of the Storage
type NVMeControllerAttributes struct {
	ReportsNamespaceGranularity                 bool `json:"ReportsNamespaceGranularity,omitempty"`
	ReportsUUIDList                             bool `json:"ReportsUUIDList,omitempty"`
	Supports128BitHostID                        bool `json:"Supports128BitHostId,omitempty"`
	SupportsEnduranceGroups                     bool `json:"SupportsEnduranceGroups,omitempty"`
	SupportsExceedingPowerOfNonOperationalState bool `json:"SupportsExceedingPowerOfNonOperationalState,omitempty"`
	SupportsPredictableLatencyMode              bool `json:"SupportsPredictableLatencyMode,omitempty"`
	SupportsReadRecoveryLevels                  bool `json:"SupportsReadRecoveryLevels,omitempty"`
	SupportsSQAssociations                      bool `json:"SupportsSQAssociations,omitempty"`
	SupportsTrafficBasedKeepAlive               bool `json:"SupportsTrafficBasedKeepAlive,omitempty"`
}

// NVMeSMARTCriticalWarnings struct is to define the NVMeSMARTCriticalWarnings of the Storage
type NVMeSMARTCriticalWarnings struct {
	MediaInReadOnly          bool `json:"MediaInReadOnly,omitempty"`
	OverallSubsystemDegraded bool `json:"OverallSubsystemDegraded,omitempty"`
	PMRUnreliable            bool `json:"PMRUnreliable,omitempty"`
	PowerBackupFailed        bool `json:"PowerBackupFailed,omitempty"`
	SpareCapacityWornOut     bool `json:"SpareCapacityWornOut,omitempty"`
}

// Drive schema represents a single physical drive for a system
type Drive struct {
	Oid                           string             `json:"@odata.id"`
	ODataContext                  string             `json:"@odata.context,omitempty"`
	ODataEtag                     string             `json:"@odata.etag,omitempty"`
	ODataType                     string             `json:"@odata.type"`
	Actions                       string             `json:"Actions,omitempty"`
	Assembly                      *Link              `json:"Assembly,omitempty"`
	AssetTag                      string             `json:"AssetTag,omitempty"`
	BlockSizeBytes                int                `json:"BlockSizeBytes,omitempty"`
	CapableSpeedGbs               float32            `json:"CapableSpeedGbs,omitempty"`
	CapacityBytes                 int                `json:"CapacityBytes,omitempty"`
	Description                   string             `json:"Description,omitempty"`
	EncryptionAbility             string             `json:"EncryptionAbility,omitempty"`
	EncryptionStatus              string             `json:"EncryptionStatus,omitempty"`
	FailurePredicted              bool               `json:"FailurePredicted,omitempty"`
	HotspareReplacementMode       string             `json:"HotspareReplacementMode,omitempty"`
	HotspareType                  string             `json:"HotspareType,omitempty"`
	ID                            string             `json:"Id"`
	Identifiers                   []*Identifier      `json:"Identifiers,omitempty"`
	IndicatorLED                  string             `json:"IndicatorLED,omitempty"`
	Links                         *DriveLinks        `json:"Links,omitempty"`
	LocationIndicatorActive       bool               `json:"LocationIndicatorActive,omitempty"`
	Manufacturer                  string             `json:"Manufacturer,omitempty"`
	MediaType                     string             `json:"MediaType,omitempty"`
	Model                         string             `json:"Model,omitempty"`
	Multipath                     bool               `json:"Multipath,omitempty"`
	Name                          string             `json:"Name"`
	NegotiatedSpeedGbs            float32            `json:"NegotiatedSpeedGbs,omitempty"`
	Oem                           *Oem               `json:"Oem,omitempty"`
	Operations                    []*DriveOperations `json:"Operations,omitempty"`
	PartNumber                    string             `json:"PartNumber,omitempty"`
	PhysicalLocation              *PhysicalLocation  `json:"PhysicalLocation,omitempty"`
	PredictedMediaLifeLeftPercent float32            `json:"PredictedMediaLifeLeftPercent,omitempty"`
	Protocol                      string             `json:"Protocol,omitempty"`
	ReadyToRemove                 bool               `json:"ReadyToRemove,omitempty"`
	Revision                      string             `json:"Revision,omitempty"`
	RotationSpeedRPM              float32            `json:"RotationSpeedRPM,omitempty"`
	SKU                           string             `json:"SKU,omitempty"`
	SerialNumber                  string             `json:"SerialNumber,omitempty"`
	Status                        *StorageStatus     `json:"Status,omitempty"`
	StatusIndicator               string             `json:"StatusIndicator,omitempty"`
	WriteCacheEnabled             bool               `json:"WriteCacheEnabled,omitempty"`
	CapacityMiB                   int                `json:"CapacityMiB,omitempty"`
	Location                      string             `json:"Location,omitempty"`
	RotationalSpeedRpm            int                `json:"RotationalSpeedRpm,omitempty"`
	FirmwareVersion               *FirmwareVersion   `json:"FirmwareVersion,omitempty"`
}

// DriveLinks represents drive links
type DriveLinks struct {
	Chassis       *Link   `json:"Chassis,omitempty"`
	Endpoints     []*Link `json:"Endpoints,omitempty"`
	Oem           *Oem    `json:"Oem,omitempty"`
	PCIeFunctions []*Link `json:"PCIeFunctions,omitempty"`
	StoragePools  []*Link `json:"StoragePools,omitempty"`
	Volumes       []*Link `json:"Volumes,omitempty"`
}

// DriveOperations represents drive operations
type DriveOperations struct {
	AssociatedTask     *Link  `json:"AssociatedTask,omitempty"`
	OperationName      string `json:"OperationName,omitempty"`
	PercentageComplete int    `json:"PercentageComplete,omitempty"`
}

// PhysicalLocation holds the location information of the drive
type PhysicalLocation struct {
	PartLocation *PartLocation `json:"PartLocation,omitempty"`
}

// FirmwareVersion of the drive
type FirmwareVersion struct {
	Current *Current `json:"Current,omitempty"`
}

// Current firmware version of the drive
type Current struct {
	VersionString string `json:"VersionString,omitempty"`
}

// Volume contains the details volume properties
type Volume struct {
	Oid                              string                   `json:"@odata.id"`
	ODataContext                     string                   `json:"@odata.context"`
	ODataEtag                        string                   `json:"@odata.etag"`
	ODataType                        string                   `json:"@odata.type"`
	AccessCapabilities               []string                 `json:"AccessCapabilities,omitempty"`
	Actions                          *Actions                 `json:"Actions,omitempty"`
	AllocatedPools                   *Link                    `json:"AllocatedPools,omitempty"`
	BlockSizeBytes                   int                      `json:"BlockSizeBytes,omitempty"`
	Capacity                         int                      `json:"Capacity,omitempty"`
	CapacityBytes                    int                      `json:"CapacityBytes,omitempty"`
	CapacitySources                  []*Link                  `json:"CapacitySources,omitempty"`
	Compressed                       bool                     `json:"Compressed,omitempty"`
	Deduplicated                     bool                     `json:"Deduplicated,omitempty"`
	Description                      string                   `json:"Description,omitempty"`
	DisplayName                      string                   `json:"DisplayName,omitempty"`
	Encrypted                        bool                     `json:"Encrypted,omitempty"`
	EncryptionTypes                  []string                 `json:"EncryptionTypes,omitempty"`
	ID                               string                   `json:"Id,omitempty"`
	Identifiers                      []*Identifier            `json:"Identifiers,omitempty"`
	IOPerfModeEnabled                bool                     `json:"IOPerfModeEnabled,omitempty"`
	IOStatistics                     *IOStatistics            `json:"AllocatedIOStatisticsPools,omitempty"`
	Links                            VolumeLinks              `json:"Links"`
	LogicalUnitNumber                int                      `json:"LogicalUnitNumber,omitempty"`
	LowSpaceWarningThresholdPercents []int                    `json:"LowSpaceWarningThresholdPercents,omitempty"`
	Manufacturer                     string                   `json:"Manufacturer,omitempty"`
	MaxBlockSizeBytes                int                      `json:"MaxBlockSizeBytes,omitempty"`
	MediaSpanCount                   int                      `json:"MediaSpanCount,omitempty"`
	Model                            string                   `json:"Model,omitempty"`
	NVMeNamespaceProperties          *NVMeNamespaceProperties `json:"NVMeNamespaceProperties,omitempty"`
	Name                             string                   `json:"Name,omitempty"`
	Oem                              *Oem                     `json:"Oem,omitempty"`
	Operations                       []*VolumeOperations      `json:"Operations,omitempty"`
	OptimumIOSizeBytes               int                      `json:"OptimumIOSizeBytes,omitempty"`
	ProvisioningPolicy               string                   `json:"ProvisioningPolicy,omitempty"`
	RAIDType                         string                   `json:"RAIDType,omitempty"`
	ReadCachePolicy                  string                   `json:"ReadCachePolicy,omitempty"`
	RecoverableCapacitySourceCount   int                      `json:"RecoverableCapacitySourceCount,omitempty"`
	RemainingCapacityPercent         int                      `json:"RemainingCapacityPercent,omitempty"`
	ReplicaInfo                      *ReplicaInfo             `json:"ReplicaInfo,omitempty"`
	ReplicaTargets                   []*Link                  `json:"ReplicaTargets,omitempty"`
	Status                           *StorageStatus           `json:"Status,omitempty"`
	StorageGroups                    *Link                    `json:"StorageGroups,omitempty"`
	StripSizeBytes                   int                      `json:"StripSizeBytes,omitempty"`
	VolumeType                       string                   `json:"VolumeType,omitempty"`
	VolumeUsage                      string                   `json:"VolumeUsage,omitempty"`
	WriteCachePolicy                 string                   `json:"WriteCachePolicy,omitempty"`
	WriteCacheState                  string                   `json:"WriteCacheState,omitempty"`
	WriteHoleProtectionPolicy        string                   `json:"WriteHoleProtectionPolicy,omitempty"`
}

// VolumeLinks represents volume links
type VolumeLinks struct {
	ClassOfService        *Link   `json:"ClassOfService,omitempty"`
	ClientEndpoints       []*Link `json:"ClientEndpoints,omitempty"`
	ConsistencyGroups     []*Link `json:"ConsistencyGroups,omitempty"`
	DedicatedSpareDrives  []*Link `json:"DedicatedSpareDrives,omitempty"`
	Drives                []*Link `json:"Drives,omitempty"`
	JournalingMedia       string  `json:"JournalingMedia,omitempty"`
	Oem                   *Oem    `json:"Oem,omitempty"`
	OwningStorageResource *Link   `json:"OwningStorageResource,omitempty"`
	OwningStorageService  *Link   `json:"OwningStorageService,omitempty"`
	ServerEndpoints       []*Link `json:"ServerEndpoints,omitempty"`
	SpareResourceSets     []*Link `json:"SpareResourceSets,omitempty"`
	StorageGroups         []*Link `json:"StorageGroups,omitempty"`
}

// NVMeNamespaceProperties represents volume NVMe Namespace.
type NVMeNamespaceProperties struct {
	FormattedLBASize                  string             `json:"FormattedLBASize,omitempty"`
	IsShareable                       bool               `json:"IsShareable,omitempty"`
	MetadataTransferredAtEndOfDataLBA bool               `json:"MetadataTransferredAtEndOfDataLBA,omitempty"`
	NamespaceFeatures                 *NamespaceFeatures `json:"NamespaceFeatures,omitempty"`
	NamespaceID                       string             `json:"NamespaceId,omitempty"`
	NumberLBAFormats                  int                `json:"NumberLBAFormats,omitempty"`
	NVMeVersion                       string             `json:"NVMeVersion,omitempty"`
}

// NamespaceFeatures property contains a set of Namespace Features
type NamespaceFeatures struct {
	SupportsAtomicTransactionSize         bool `json:"SupportsAtomicTransactionSize,omitempty"`
	SupportsDeallocatedOrUnwrittenLBError bool `json:"SupportsDeallocatedOrUnwrittenLBError,omitempty"`
	SupportsIOPerformanceHints            bool `json:"SupportsIOPerformanceHints,omitempty"`
	SupportsNGUIDReuse                    bool `json:"SupportsNGUIDReuse,omitempty"`
	SupportsThinProvisioning              bool `json:"SupportsThinProvisioning,omitempty"`
}

// VolumeOperations represents operations running on volume
type VolumeOperations struct {
	AssociatedFeaturesRegistry *Link  `json:"AssociatedFeaturesRegistry,omitempty"`
	OperationName              string `json:"OperationName,omitempty"`
	PercentageComplete         int    `json:"PercentageComplete,omitempty"`
}

// ReplicaInfo describes this storage volume in its role as a target replica
type ReplicaInfo struct {
}
