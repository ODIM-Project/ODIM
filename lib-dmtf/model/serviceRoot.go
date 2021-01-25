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

// ServiceRoot a defines redfish service root
type ServiceRoot struct {
	ODataContext              string                     `json:"@odata.context,omitempty"`
	ODataID                   string                     `json:"@odata.id,omitempty"`
	ODataEtag                 string                     `json:"@odata.etag,omitempty"`
	ODataType                 string                     `json:"@odata.type,omitempty"`
	AccountService            *Link                      `json:"AccountService,omitempty"`
	Chassis                   *Link                      `json:"Chassis,omitempty"`
	CompositionService        *Link                      `json:"CompositionService,omitempty"`
	Description               string                     `json:"Description,omitempty"`
	EventService              *Link                      `json:"EventService,omitempty"`
	Fabrics                   *Link                      `json:"Fabrics,omitempty"`
	ID                        string                     `json:"Id"`
	JSONSchemas               *Link                      `json:"JsonSchemas,omitempty"`
	Links                     Links                      `json:"Links"`
	Managers                  *Link                      `json:"Managers,omitempty"`
	Name                      string                     `json:"Name"`
	OEM                       *Oem                       `json:"Oem,omitempty"`
	Product                   *Product                   `json:"Product,omitempty"`
	ProtocolFeaturesSupported *ProtocolFeaturesSupported `json:"ProtocolFeaturesSupported,omitempty"`
	RedfishVersion            string                     `json:"RedfishVersion,omitempty"`
	Registries                *Link                      `json:"Registries,omitempty"`
	SessionService            *Link                      `json:"SessionService,omitempty"`
	StorageServices           *Link                      `json:"StorageServices,omitempty"`
	StorageSystems            *Link                      `json:"StorageSystems,omitempty"`
	Systems                   *Link                      `json:"Systems,omitempty"`
	Tasks                     *Link                      `json:"Tasks,omitempty"`
	UpdateService             *Link                      `json:"UpdateService,omitempty"`
	UUID                      string                     `json:"UUID,omitempty"`
}

// Product redfish structure
type Product struct{}

// ProtocolFeaturesSupported redfish structure
type ProtocolFeaturesSupported struct{}
