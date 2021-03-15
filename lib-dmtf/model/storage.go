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
	Oid                string              `json:"@odata.id"`
	Ocontext           string              `json:"@odata.context,omitempty"`
	Oetag              string              `json:"@odata.etag,omitempty"`
	Otype              string              `json:"@odata.type,omitempty"`
	Description        string              `json:"description,omitempty"`
	ID                 string              `json:"Id,omitempty"`
	Name               string              `json:"Name,omitempty"`
	Oem                Oem                 `json:"Oem,omitempty"`
	Drives             []Link              `json:"Drives,omitempty"`
	Links              StorageLinks        `json:"Links,omitempty"`
	Redundancy         []StorageRedundancy `json:"Redundancy,omitempty"` //TODO
	Status             Status              `json:"Status,omitempty"`
	StorageControllers []StorageController `json:"StorageControllers,omitempty"`
	Volumes            Link                `json:"Volumes,omitempty"`
	ConsistencyGroups  Link                `json:"ConsistencyGroups,omitempty"`
	Controllers        Link                `json:"Controllers,omitempty"`
	EndpointGroups     Link                `json:"EndpointGroups,omitempty"`
	FileSystems        Link                `json:"FileSystems,omitempty"`
	Identifiers        *Identifier         `json:"Identifiers,omitempty"`
	StorageGroups      Link                `json:"StorageGroups,omitempty"`
	StoragePools       Link                `json:"StoragePools,omitempty"`
}

//StorageController in place(it has Oid it may be get, TODO)
type StorageController struct {
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
	Actions                      Actions                   `json:"Actions,omitempty"`
	Assembly                     Assembly                  `json:"Assembly,omitempty"`
	CacheSummary                 CacheSummary              `json:"CacheSummary,omitempty"`
	ControllerRates              ControllerRates           `json:"ControllerRates,omitempty"`
	Identifiers                  []Identifier              `json:"Identifiers,omitempty"`
	Links                        StorageControllersLinks   `json:"Links,omitempty"`
	Location                     StorageControllerLocation `json:"Location,omitempty"`
	NVMeControllerProperties     NVMeControllerProperties  `json:"NVMeControllerProperties,omitempty"`
	Oem                          Oem                       `json:"Oem,omitempty"`
	PCIeInterface                PCIeInterface             `json:"PCIeInterface,omitempty"`
	Ports                        Link                      `json:"Ports,omitempty"`
	Status                       StorageStatus             `json:"Status,omitempty"`
}

//Actions redfish structure
type Actions struct {
}

//Identifier redfish structure
type Identifier struct {
	DurableName       string `json:"DurableName,omitempty"`
	DurableNameFormat string `json:"DurableNameFormat,omitempty"`
}

// StorageLinks struct is for storage Links schema
type StorageLinks struct {
	Enclosures      []Link `json:"Enclosures,omitempty"`
	SimpleStorage   Link   `json:"SimpleStorage,omitempty"`
	StorageServices []Link `json:"StorageServices,omitempty"`
	Oem             *Oem   `json:"Oem,omitempty"`
}

// StorageRedundancy struct is for Redundancy schema
type StorageRedundancy struct {
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
	Endpoints       []Link `json:"Enclosures,omitempty"`
	PCIeFunctions   Link   `json:"PCIeFunctions,omitempty"`
	StorageServices []Link `json:"StorageServices,omitempty"`
	Oem             *Oem   `json:"Oem,omitempty"`
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
	Volume      Link   `json:"Volume,omitempty"`
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

// StorageControllerLocation holds the location information of the storage controller
type StorageControllerLocation struct {
	AltitudeMeters int            `json:"AltitudeMeters,omitempty"`
	Latitude       int            `json:"Latitude,omitempty"`
	Longitude      int            `json:"Longitude,omitempty"`
	Contacts       Contacts       `json:"Contacts,omitempty"`
	Oem            *Oem           `json:"Oem,omitempty"`
	PartLocation   *PartLocation  `json:"PartLocation,omitempty"`
	Placement      *Placement     `json:"Placement,omitempty"`
	PostalAddress  *PostalAddress `json:"PostalAddress,omitempty"`
}

// PartLocation holds the location information of the storage controller
type PartLocation struct {
	Orientation          string `json:"Orientation,omitempty"`
	Reference            string `json:"Reference,omitempty"`
	LocationOrdinalValue int    `json:"LocationOrdinalValue,omitempty"`
	LocationType         string `json:"LocationType,omitempty"`
	ServiceLabel         string `json:"ServiceLabel,omitempty"`
}

// Contacts holds the Contacts information of the storage controller
type Contacts struct {
	ContactName  string `json:"ContactName,omitempty"`
	EmailAddress string `json:"EmailAddress,omitempty"`
	PhoneNumber  string `json:"PhoneNumber,omitempty"`
}

// Placement holds the Placement information of the storage controller
type Placement struct {
	AdditionalInfo  string `json:"AdditionalInfo,omitempty"`
	Rack            string `json:"Rack,omitempty"`
	RackOffset      int    `json:"RackOffset,omitempty"`
	RackOffsetUnits string `json:"RackOffsetUnits,omitempty"`
	Row             string `json:"Row,omitempty"`
}

// PostalAddress holds the PostalAddress information of the storage controller
type PostalAddress struct {
	AdditionalCode         string `json:"AdditionalCode,omitempty"`
	AdditionalInfo         string `json:"AdditionalInfo,omitempty"`
	Building               string `json:"Building,omitempty"`
	City                   string `json:"City,omitempty"`
	Community              string `json:"Community,omitempty"`
	Country                string `json:"Country,omitempty"`
	District               string `json:"District,omitempty"`
	Division               string `json:"Division,omitempty"`
	Floor                  string `json:"Floor,omitempty"`
	GPSCoords              string `json:"GPSCoords,omitempty"`
	HouseNumber            int    `json:"HouseNumber,omitempty"`
	HouseNumberSuffix      string `json:"HouseNumberSuffix,omitempty"`
	Landmark               string `json:"Landmark,omitempty"`
	LeadingStreetDirection string `json:"LeadingStreetDirection,omitempty"`
	Location               string `json:"Location,omitempty"`
	Name                   string `json:"Name,omitempty"`
	Neighborhood           string `json:"Neighborhood,omitempty"`
	PlaceType              string `json:"PlaceType,omitempty"`
	POBox                  string `json:"POBox,omitempty"`
	PostalCode             string `json:"PostalCode,omitempty"`
	Road                   string `json:"Road,omitempty"`
	RoadBranch             string `json:"RoadBranch,omitempty"`
	RoadPostModifier       string `json:"RoadPostModifier,omitempty"`
	RoadPreModifier        string `json:"RoadPreModifier,omitempty"`
	RoadSection            string `json:"RoadSection,omitempty"`
	RoadSubBranch          string `json:"RoadSubBranch,omitempty"`
	Room                   string `json:"Room,omitempty"`
	Seat                   string `json:"Seat,omitempty"`
	Street                 string `json:"Street,omitempty"`
	StreetSuffix           string `json:"StreetSuffix,omitempty"`
	Territory              string `json:"Territory,omitempty"`
	TrailingStreetSuffix   string `json:"TrailingStreetSuffix,omitempty"`
	Unit                   string `json:"Unit,omitempty"`
}
