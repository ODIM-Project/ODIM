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

//Power is the redfish Power model according to the 2020.3 release
type Power struct {
	ODataContext  string          `json:"@odata.context,omitempty"`
	ODataEtag     string          `json:"@odata.etag,omitempty"`
	ODataID       string          `json:"@odata.id"`
	ODataType     string          `json:"@odata.type"`
	Actions       *OemActions     `json:"Actions,omitempty"`
	Description   string          `json:"Description,omitempty"`
	ID            string          `json:"Id"`
	Name          string          `json:"Name"`
	Oem           interface{}     `json:"Oem,omitempty"`
	Status        *Status         `json:"Status,omitempty"`
	PowerControl  []*PowerControl `json:"PowerControl,omitempty"`
	PowerSupplies []*PowerControl `json:"PowerSupplies,omitempty"`
	Redundancy    []Redundancy    `json:"Redundancy,omitempty"`
	Voltages      []*Voltages     `json:"Voltages,omitempty"`
}

// PowerControl redfish model
type PowerControl struct {
	ODataID             string        `json:"@odata.id"`
	Actions             *OemActions   `json:"Actions,omitempty"`
	MemberID            string        `json:"MemberId,omitempty"`
	Name                string        `json:"Name,omitempty"`
	Oem                 interface{}   `json:"Oem,omitempty"`
	PhysicalContext     string        `json:"PhysicalContext,omitempty"`
	PowerAllocatedWatts float64       `json:"PowerAllocatedWatts,omitempty"`
	PowerAvailableWatts float64       `json:"PowerAvailableWatts,omitempty"`
	PowerCapacityWatts  float64       `json:"PowerCapacityWatts,omitempty"`
	PowerConsumedWatts  float64       `json:"PowerConsumedWatts,omitempty"`
	PowerLimit          *PowerLimit   `json:"PowerLimit,omitempty"`
	PowerMetrics        *PowerMetrics `json:"PowerMetrics,omitempty"`
	PowerRequestedWatts float64       `json:"PowerRequestedWatts,omitempty"`
	RelatedItem         []Link        `json:"RelatedItem,omitempty"`
	Status              *Status       `json:"Status,omitempty"`
}

// PowerLimit redfish model
type PowerLimit struct {
	CorrectionInMs int     `json:"CorrectionInMs,omitempty"`
	LimitException string  `json:"LimitException,omitempty"`
	LimitInWatts   float64 `json:"LimitInWatts,omitempty"`
}

// PowerMetrics redfish model
type PowerMetrics struct {
	AverageConsumedWatts float64 `json:"AverageConsumedWatts,omitempty"`
	IntervalInMin        int     `json:"IntervalInMin,omitempty"`
	MaxConsumedWatts     float64 `json:"MaxConsumedWatts,omitempty"`
	MinConsumedWatts     float64 `json:"MinConsumedWatts,omitempty"`
}

// PowerSupplies redfish model
type PowerSupplies struct {
	ODataID              string         `json:"@odata.id"`
	Actions              *OemActions    `json:"Actions,omitempty"`
	Assembly             *Link          `json:"Assembly,omitempty"`
	Name                 string         `json:"Name,omitempty"`
	Oem                  interface{}    `json:"Oem,omitempty"`
	Status               *Status        `json:"Status,omitempty"`
	EfficiencyPercent    float64        `json:"EfficiencyPercent,omitempty"`
	FirmwareVersion      string         `json:"FirmwareVersion,omitempty"`
	HotPluggable         bool           `json:"HotPluggable,omitempty"`
	IndicatorLED         string         `json:"IndicatorLED,omitempty"`
	InputRanges          []*InputRanges `json:"InputRanges,omitempty"`
	LastPowerOutputWatts float64        `json:"LastPowerOutputWatts,omitempty"`
	LineInputVoltage     float64        `json:"LineInputVoltage,omitempty"`
	LineInputVoltageType string         `json:"LineInputVoltageType,omitempty"`
	Location             interface{}    `json:"Location,omitempty"`
	Manufacturer         string         `json:"Manufacturer,omitempty"`
	MemberID             string         `json:"MemberId,omitempty"`
	Model                string         `json:"Model,omitempty"`
	PartNumber           string         `json:"PartNumber,omitempty"`
	PowerCapacityWatts   float64        `json:"PowerCapacityWatts,omitempty"`
	PowerInputWatts      float64        `json:"PowerInputWatts,omitempty"`
	PowerOutputWatts     float64        `json:"PowerOutputWatts,omitempty"`
	PowerSupplyType      string         `json:"PowerSupplyType,omitempty"`
	Redundancy           []Redundancy   `json:"Redundancy,omitempty"`
	RelatedItem          []Link         `json:"RelatedItem,omitempty"`
	SerialNumber         string         `json:"SerialNumber,omitempty"`
	SparePartNumber      string         `json:"SparePartNumber,omitempty"`
}

// InputRanges redfish model
type InputRanges struct {
	InputType          string      `json:"InputType,omitempty"`
	MaximumFrequencyHz float64     `json:"MaximumFrequencyHz,omitempty"`
	MaximumVoltage     float64     `json:"MaximumVoltage,omitempty"`
	MinimumFrequencyHz float64     `json:"MinimumFrequencyHz,omitempty"`
	MinimumVoltage     float64     `json:"MinimumVoltage,omitempty"`
	Oem                interface{} `json:"Oem,omitempty"`
	OutputWattage      float64     `json:"OutputWattage,omitempty"`
}

// Voltages redfish model
type Voltages struct {
	ODataID                   string      `json:"@odata.id"`
	Actions                   *OemActions `json:"Actions,omitempty"`
	Name                      string      `json:"Name,omitempty"`
	Oem                       interface{} `json:"Oem,omitempty"`
	Status                    *Status     `json:"Status,omitempty"`
	RelatedItem               []Link      `json:"RelatedItem,omitempty"`
	LowerThresholdCritical    float64     `json:"LowerThresholdCritical,omitempty"`
	LowerThresholdFatal       float64     `json:"LowerThresholdFatal,omitempty"`
	LowerThresholdNonCritical float64     `json:"LowerThresholdNonCritical,omitempty"`
	MaxReadingRange           float64     `json:"MaxReadingRange,omitempty"`
	MemberID                  string      `json:"MemberId,omitempty"`
	MinReadingRange           float64     `json:"MinReadingRange,omitempty"`
	PhysicalContext           string      `json:"PhysicalContext,omitempty"`
	ReadingVolts              float64     `json:"ReadingVolts,omitempty"`
	SensorNumber              int         `json:"SensorNumber,omitempty"`
	UpperThresholdCritical    float64     `json:"UpperThresholdCritical,omitempty"`
	UpperThresholdFatal       float64     `json:"UpperThresholdFatal,omitempty"`
	UpperThresholdNonCritical float64     `json:"UpperThresholdNonCritical,omitempty"`
}
