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

// Fabric is the redfish Fabric model according to the 2020.3 release
type Fabric struct {
	ODataContext   string      `json:"@odata.context,omitempty"`
	ODataEtag      string      `json:"@odata.etag,omitempty"`
	ODataID        string      `json:"@odata.id"`
	ODataType      string      `json:"@odata.type"`
	Actions        *OemActions `json:"Actions,omitempty"`
	Description    string      `json:"Description,omitempty"`
	ID             string      `json:"Id"`
	Links          *Links      `json:"Links,omitempty"`
	Name           string      `json:"Name"`
	Oem            *Oem        `json:"Oem,omitempty"`
	Status         *Status     `json:"Status,omitempty"`
	AddressPools   *Link       `json:"AddressPools"`
	Connections    *Link       `json:"Connections,omitempty"`
	EndpointGroups *Link       `json:"EndpointGroups,omitempty"`
	Endpoints      *Link       `json:"Endpoints"`
	Switches       *Link       `json:"Switches"`
	Zones          *Link       `json:"Zones"`
	MaxZones       int         `json:"MaxZones,omitempty"`
	FabricType     string      `json:"FabricType"`
}
