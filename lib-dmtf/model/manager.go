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
	ODataContext            string       `json:"@odata.context,omitempty"`
	ODataEtag               string       `json:"@odata.etag,omitempty"`
	ODataID                 string       `json:"@odata.id"`
	ODataType               string       `json:"@odata.type"`
	Actions                 *OemActions  `json:"Actions,omitempty"`
	Description             string       `json:"Description,omitempty"`
	ID                      string       `json:"id"`
	Links                   *Links       `json:"Links,omitempty"`
	Name                    string       `json:"Name"`
	Oem                     *Oem         `json:"Oem,omitempty"`
	EthernetInterfaces      *Link        `json:"EthernetInterfaces,omitempty"`
	FirmwareVersion         string       `json:"FirmwareVersion,omitempty"`
	Status                  *Status      `json:"Status,omitempty"`
	AutoDSTEnabled          string       `json:"AutoDSTEnabled,omitempty"`
	DateTime                string       `json:"DateTime,omitempty"`
	DateTimeLocalOffset     string       `json:"DateTimeLocalOffset,omitempty"`
	HostInterfaces          *Link        `json:"HostInterfaces,omitempty"`
	LastResetTime           string       `json:"LastResetTime,omitempty"`
	LogServices             *Link        `json:"LogServices,omitempty"`
	ManagerType             string       `json:"ManagerType,omitempty"`
	Manufacturer            string       `json:"Manufacturer,omitempty"`
	Model                   string       `json:"Model,omitempty"`
	NetworkProtocol         *Link        `json:"NetworkProtocol,omitempty"`
	PartNumber              string       `json:"PartNumber,omitempty"`
	PowerState              string       `json:"PowerState,omitempty"`
	Redundancy              []Redundancy `json:"Redundancy,omitempty"`
	RemoteAccountService    *Link        `json:"RemoteAccountService,omitempty"`
	RemoteRedfishServiceURI string       `json:"RemoteRedfishServiceUri,omitempty"`
	SerialNumber            string       `json:"SerialNumber,omitempty"`
	ServiceEntryPointUUID   string       `json:"ServiceEntryPointUUID,omitempty"`
	SerialInterfaces        *Link        `json:"SerialInterfaces,omitempty"`
	TimeZoneName            string       `json:"TimeZoneName,omitempty"`
	UUID                    string       `json:"UUID,omitempty"`
}
