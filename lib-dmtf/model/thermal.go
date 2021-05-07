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

//Thermal is the redfish Power model according to the 2020.3 release
type Thermal struct {
	ODataContext string          `json:"@odata.context,omitempty"`
	ODataEtag    string          `json:"@odata.etag,omitempty"`
	ODataID      string          `json:"@odata.id"`
	ODataType    string          `json:"@odata.type"`
	Actions      *OemActions     `json:"Actions,omitempty"`
	Description  string          `json:"Description,omitempty"`
	ID           string          `json:"Id"`
	Name         string          `json:"Name"`
	Oem          interface{}     `json:"Oem,omitempty"`
	Status       *Status         `json:"Status,omitempty"`
	Fans         []*Fans         `json:"Fans,omitempty"`
	Redundancy   []Redundancy    `json:"Redundancy,omitempty"`
	Temperatures []*Temperatures `json:"Temperatures,omitempty"`
}

// Fans redfish model
type Fans struct {
	ODataID                   string       `json:"@odata.id"`
	Actions                   *OemActions  `json:"Actions,omitempty"`
	Assembly                  *Link        `json:"Assembly,omitempty"`
	HotPluggable              bool         `json:"HotPluggable,omitempty"`
	IndicatorLED              string       `json:"IndicatorLED,omitempty"`
	Location                  interface{}  `json:"Location,omitempty"`
	LowerThresholdCritical    int          `json:"LowerThresholdCritical,omitempty"`
	LowerThresholdFatal       int          `json:"LowerThresholdFatal,omitempty"`
	LowerThresholdNonCritical int          `json:"LowerThresholdNonCritical,omitempty"`
	Manufacturer              string       `json:"Manufacturer,omitempty"`
	MaxReadingRange           int          `json:"MaxReadingRange,omitempty"`
	MemberID                  string       `json:"MemberId,omitempty"`
	MinReadingRange           int          `json:"MinReadingRange,omitempty"`
	Model                     string       `json:"Model,omitempty"`
	Name                      string       `json:"Name,omitempty"`
	Oem                       interface{}  `json:"Oem,omitempty"`
	PartNumber                string       `json:"PartNumber,omitempty"`
	PhysicalContext           string       `json:"PhysicalContext,omitempty"`
	Reading                   int          `json:"Reading,omitempty"`
	ReadingUnits              string       `json:"ReadingUnits,omitempty"`
	Redundancy                []Redundancy `json:"Redundancy,omitempty"`
	RelatedItem               []Link       `json:"RelatedItem,omitempty"`
	SensorNumber              int          `json:"SensorNumber,omitempty"`
	SerialNumber              string       `json:"SerialNumber,omitempty"`
	SparePartNumber           string       `json:"SparePartNumber,omitempty"`
	Status                    *Status      `json:"Status,omitempty"`
	UpperThresholdCritical    int          `json:"UpperThresholdCritical,omitempty"`
	UpperThresholdFatal       int          `json:"UpperThresholdFatal,omitempty"`
	UpperThresholdNonCritical int          `json:"UpperThresholdNonCritical,omitempty"`
}

// Temperatures redfish model
type Temperatures struct {
	ODataID                            string      `json:"@odata.id"`
	Actions                            *OemActions `json:"Actions,omitempty"`
	AdjustedMaxAllowableOperatingValue int         `json:"AdjustedMaxAllowableOperatingValue,omitempty"`
	AdjustedMinAllowableOperatingValue int         `json:"AdjustedMinAllowableOperatingValue,omitempty"`
	DeltaPhysicalContext               string      `json:"DeltaPhysicalContext,omitempty"`
	DeltaReadingCelsius                float64     `json:"DeltaReadingCelsius,omitempty"`
	LowerThresholdCritical             float64     `json:"LowerThresholdCritical,omitempty"`
	LowerThresholdFatal                float64     `json:"LowerThresholdFatal,omitempty"`
	LowerThresholdNonCritical          float64     `json:"LowerThresholdNonCritical,omitempty"`
	LowerThresholdUser                 int         `json:"LowerThresholdUser,omitempty"`
	MaxAllowableOperatingValue         int         `json:"MaxAllowableOperatingValue,omitempty"`
	MaxReadingRangeTemp                float64     `json:"MaxReadingRangeTemp,omitempty"`
	MemberID                           string      `json:"MemberId,omitempty"`
	MinAllowableOperatingValue         int         `json:"MinAllowableOperatingValue,omitempty"`
	MinReadingRangeTemp                float64     `json:"MinReadingRangeTemp,omitempty"`
	Name                               string      `json:"Name,omitempty"`
	Oem                                interface{} `json:"Oem,omitempty"`
	PhysicalContext                    string      `json:"PhysicalContext,omitempty"`
	ReadingCelsius                     float64     `json:"ReadingCelsius,omitempty"`
	RelatedItem                        []Link      `json:"RelatedItem,omitempty"`
	SensorNumber                       int         `json:"SensorNumber,omitempty"`
	Status                             *Status     `json:"Status,omitempty"`
	UpperThresholdCritical             float64     `json:"UpperThresholdCritical,omitempty"`
	UpperThresholdFatal                float64     `json:"UpperThresholdFatal,omitempty"`
	UpperThresholdNonCritical          float64     `json:"UpperThresholdNonCritical,omitempty"`
	UpperThresholdUser                 int         `json:"UpperThresholdUser,omitempty"`
}
