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

// Switch is the redfish Switch model according to the 2020.3 release
type Switch struct {
	ODataContext            string       `json:"@odata.context,omitempty"`
	ODataEtag               string       `json:"@odata.etag,omitempty"`
	ODataID                 string       `json:"@odata.id"`
	ODataType               string       `json:"@odata.type"`
	Actions                 *OemActions  `json:"Actions,omitempty"`
	Description             string       `json:"Description,omitempty"`
	ID                      string       `json:"Id"`
	Links                   *SwitchLinks `json:"Links,omitempty"`
	Name                    string       `json:"Name"`
	Oem                     interface{}  `json:"Oem,omitempty"`
	Status                  *Status      `json:"Status,omitempty"`
	AssetTag                string       `json:"AssetTag,omitempty"`
	CurrentBandwidthGbps    float64      `json:"CurrentBandwidthGbps,omitempty"`
	DomainID                string       `json:"DomainID,omitempty"`
	FirmwareVersion         string       `json:"FirmwareVersion,omitempty"`
	IndicatorLED            string       `json:"IndicatorLED,omitempty"`
	IsManaged               bool         `json:"IsManaged"`
	LocationIndicatorActive bool         `json:"LocationIndicatorActive,omitempty"`
	LogServices             *Link        `json:"LogServices,omitempty"`
	Manufacturer            string       `json:"Manufacturer,omitempty"`
	MaxBandwidthGbps        float64      `json:"MaxBandwidthGbps,omitempty"`
	Model                   string       `json:"Model,omitempty"`
	PartNumber              string       `json:"PartNumber,omitempty"`
	Ports                   *Link        `json:"Ports"`
	PowerState              string       `json:"PowerState,omitempty"`
	Redundancy              []Redundancy `json:"Redundancy,omitempty"`
	SerialNumber            string       `json:"SerialNumber,omitempty"`
	SKU                     string       `json:"SKU,omitempty"`
	SupportedProtocols      []string     `json:"SupportedProtocols,omitempty"`
	SwitchType              string       `json:"SwitchType"`
	TotalSwitchWidth        int          `json:"TotalSwitchWidth,omitempty"`
	UUID                    string       `json:"UUID,omitempty"`
}

// SwitchLinks defines the
type SwitchLinks struct {
	Chassis   *Link       `json:"Chassis,omitempty"`
	Endpoints []*Link     `json:"Endpoints,omitempty"`
	ManagedBy []*Link     `json:"ManagedBy,omitempty"`
	Oem       interface{} `json:"Oem,omitempty"`
}
