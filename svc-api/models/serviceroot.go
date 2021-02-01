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

//Package models ...
package models

//ServiceRoot struct definition
type ServiceRoot struct {
	OdataContext              string       `json:"@odata.context"`
	Etag                      string       `json:"@odata.etag,omitempty"`
	OdataID                   string       `json:"@odata.id"`
	OdataType                 string       `json:"@odata.type"`
	ID                        string       `json:"Id"`
	ProtocolFeaturesSupported *PFSupported `json:"ProtocolFeaturesSupported,omitempty"`
	Registries                *Service     `json:"Registries,omitempty"`
	SessionService            *Service     `json:"SessionService,omitempty"`
	AccountService            *Service     `json:"AccountService,omitempty"`
	JSONSchemas               *Service     `json:"JsonSchemas,omitempty"`
	EventService              *Service     `json:"EventService,omitempty"`
	JobService                *Service     `json:"JobService,omitempty"`
	Tasks                     *Service     `json:"Tasks,omitempty"`
	AggregationService        *Service     `json:"AggregationService,omitempty"`
	Systems                   *Service     `json:"Systems,omitempty"`
	Chassis                   *Service     `json:"Chassis,omitempty"`
	Fabrics                   *Service     `json:"Fabrics,omitempty"`
	Managers                  *Service     `json:"Managers,omitempty"`
	UpdateService             *Service     `json:"UpdateService,omitempty"`
	Links                     Links        `json:"Links"`
	Name                      string       `json:"Name"`
	OEM                       OEM          `json:"Oem"`
	RedfishVersion            string       `json:"RedfishVersion"`
	UUID                      string       `json:"UUID"`
}

//PFSupported struct definition
type PFSupported struct {
	ExcerptQuery    bool         `json:"ExcerptQuery"`
	ExpandQuery     *ExpandQuery `json:"ExpandQuery:omitempty"`
	FilterQuery     bool         `json:"FilterQuery"`
	OnlyMemberQuery bool         `json:"OnlyMemberQuery"`
	SelectQuery     bool         `json:"SelectQuery"`
}

//ExpandQuery struct definition
type ExpandQuery struct {
	ExpandAll bool `json:"ExpandAll"`
	Levels    bool `json:"Levels"`
	Links     bool `json:"Links"`
	MaxLevels int  `json:"MaxLevels"`
	NoLinks   bool `json:"NoLinks"`
}

//Service struct definition
type Service struct {
	OdataID string `json:"@odata.id"`
}

//Links struct definition
type Links struct {
	Sessions Sessions `json:"Sessions"`
}

//Sessions struct definition
type Sessions struct {
	OdataID string `json:"@odata.id"`
}

//OEM struct definition
type OEM struct {
}
