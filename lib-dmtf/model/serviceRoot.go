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
	AggregationService        *Link                      `json:"AggregationService,omitempty"`
	Cables                    *Link                      `json:"Cables,omitempty"`
	CertificateService        *Link                      `json:"CertificateService,omitempty"`
	Facilities                *Link                      `json:"Facilities,omitempty"`
	JobService                *Link                      `json:"JobService,omitempty"`
	KeyService                *Link                      `json:"KeyService,omitempty"`
	NVMeDomains               *Link                      `json:"NVMeDomains,omitempty"`
	ResourceBlocks            *Link                      `json:"ResourceBlocks,omitempty"`
	Storage                   Storage                    `json:"Storage,omitempty"`
	TelemetryService          *Link                      `json:"TelemetryService,omitempty"`
	Vendor                    string                     `json:"Vendor,omitempty"`
}

// Product redfish structure
type Product struct{}

// ProtocolFeaturesSupported redfish structure
type ProtocolFeaturesSupported struct {
	ExcerptQuery    bool         `json:"ExcerptQuery"`
	ExpandQuery     *ExpandQuery `json:"ExpandQuery:omitempty"`
	FilterQuery     bool         `json:"FilterQuery"`
	OnlyMemberQuery bool         `json:"OnlyMemberQuery"`
	SelectQuery     bool         `json:"SelectQuery"`
}

type ExpandQuery struct {
	ExpandAll bool `json:"ExpandAll"`
	Levels    bool `json:"Levels"`
	Links     bool `json:"Links"`
	MaxLevels int  `json:"MaxLevels"`
	NoLinks   bool `json:"NoLinks"`
}
