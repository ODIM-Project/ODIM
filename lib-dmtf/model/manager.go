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

// Manager is the redfish Manager model according to the 2020.3 release
type Manager struct {
	ODataContext            string        `json:"@odata.context,omitempty"`
	ODataEtag               string        `json:"@odata.etag,omitempty"`
	ODataID                 string        `json:"@odata.id"`
	ODataType               string        `json:"@odata.type"`
	Actions                 *OemActions   `json:"Actions,omitempty"`
	Description             string        `json:"Description,omitempty"`
	ID                      string        `json:"Id"`
	Links                   *ManagerLinks `json:"Links,omitempty"`
	Name                    string        `json:"Name"`
	Oem                     *Oem          `json:"Oem,omitempty"`
	EthernetInterfaces      *Link         `json:"EthernetInterfaces,omitempty"`
	FirmwareVersion         string        `json:"FirmwareVersion,omitempty"`
	Status                  *Status       `json:"Status,omitempty"`
	AutoDSTEnabled          string        `json:"AutoDSTEnabled,omitempty"`
	DateTime                string        `json:"DateTime,omitempty"`
	DateTimeLocalOffset     string        `json:"DateTimeLocalOffset,omitempty"`
	HostInterfaces          *Link         `json:"HostInterfaces,omitempty"`
	LastResetTime           string        `json:"LastResetTime,omitempty"`
	LogServices             *Link         `json:"LogServices,omitempty"`
	ManagerType             string        `json:"ManagerType,omitempty"`
	Manufacturer            string        `json:"Manufacturer,omitempty"`
	Model                   string        `json:"Model,omitempty"`
	NetworkProtocol         *Link         `json:"NetworkProtocol,omitempty"`
	PartNumber              string        `json:"PartNumber,omitempty"`
	PowerState              string        `json:"PowerState,omitempty"`
	Redundancy              []Redundancy  `json:"Redundancy,omitempty"`
	RemoteAccountService    *Link         `json:"RemoteAccountService,omitempty"`
	RemoteRedfishServiceURI string        `json:"RemoteRedfishServiceUri,omitempty"`
	SerialNumber            string        `json:"SerialNumber,omitempty"`
	ServiceEntryPointUUID   string        `json:"ServiceEntryPointUUID,omitempty"`
	SerialInterfaces        *Link         `json:"SerialInterfaces,omitempty"`
	TimeZoneName            string        `json:"TimeZoneName,omitempty"`
	UUID                    string        `json:"UUID,omitempty"`
}

// ManagerLinks ...
type ManagerLinks struct {
	ActiveSoftwareImage     *Link       `json:"ActiveSoftwareImage,omitempty"`
	ManagedBy               []Link      `json:"ManagedBy,omitempty"`
	ManagerByCount          int         `json:"ManagedBy@odata.count,omitempty"`
	ManagerForChassis       []Link      `json:"ManagerForChassis,omitempty"`
	ManagerForChassisCount  int         `json:"ManagerForChassis@odata.count,omitempty"`
	ManagerForManagers      []Link      `json:"ManagerForManagers,omitempty"`
	ManagerForManagersCount int         `json:"ManagerForManagers@odata.count,omitempty"`
	ManagerForServers       []Link      `json:"ManagerForServers,omitempty"`
	ManagerForServersCount  int         `json:"ManagerForServers@odata.count,omitempty"`
	ManagerForSwitches      []Link      `json:"ManagerForSwitches,omitempty"`
	ManagerForSwitchesCount int         `json:"ManagerForSwitches@odata.count,omitempty"`
	ManagerInChassis        []Link      `json:"ManagerInChassis,omitempty"`
	ManagerInChassisCount   int         `json:"ManagerInChassis@odata.count,omitempty"`
	Oem                     interface{} `json:"Oem,omitempty"`
	SoftwareImages          *Link       `json:"SoftwareImages,omitempty"`
}

//VirtualMedia is a redfish virtual media model
type VirtualMedia struct {
	ODataContext         string      `json:"@odata.context,omitempty"`
	ODataEtag            string      `json:"@odata.etag,omitempty"`
	ODataID              string      `json:"@odata.id"`
	ODataType            string      `json:"@odata.type"`
	Actions              VMActions   `json:"Actions,omitempty"`
	ConnectedVia         string      `json:"ConnectedVia,omitempty"`
	Description          string      `json:"Description,omitempty"`
	ID                   string      `json:"Id"`
	Image                string      `json:"Image"`
	ImageName            string      `json:"ImageName,omitempty"`
	Inserted             bool        `json:"Inserted"`
	MediaTypes           []string    `json:"MediaTypes,omitempty"`
	Name                 string      `json:"Name"`
	Oem                  interface{} `json:"Oem,omitempty"`
	Password             string      `json:"Password,omitempty"`
	TransferMethod       string      `json:"TransferMethod,omitempty"`
	TransferProtocolType string      `json:"TransferProtocolType,omitempty"`
	UserName             string      `json:"UserName,omitempty"`
	VerifyCertificate    bool        `json:"VerifyCertificate,omitempty"`
	WriteProtected       bool        `json:"WriteProtected,omitempty"`
	Status               *Status     `json:"Status,omitempty"`
}

// VMActions contains the actions property details of virtual media
type VMActions struct {
	EjectMedia  ActionTarget `json:"EjectMedia"`
	InsertMedia ActionTarget `json:"InsertMedia"`
}

// ActionTarget contains the action target
type ActionTarget struct {
	Target string `json:"target"`
}
