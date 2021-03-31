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

// Redundancy redfish Redundancy model according to the 2020.3 release
type Redundancy struct {
	Oid               string      `json:"@odata.id,omitempty"`
	Actions           *OemActions `json:"Actions,omitempty"`
	MaxNumSupported   int         `json:"MaxNumSupported,omitempty"`
	MemberID          string      `json:"MemberId,omitempty"`
	MinNumNeeded      int         `json:"MinNumNeeded,omitempty"`
	Mode              string      `json:"Mode,omitempty"`
	Name              string      `json:"Name,omitempty"`
	Oem               interface{} `json:"Oem,omitempty"`
	RedundancyEnabled bool        `json:"RedundancyEnabled,omitempty"`
	RedundancySet     []*Link     `json:"RedundancySet,omitempty"`
	Status            *Status     `json:"Status,omitempty"`
}

//Identifier redfish structure
type Identifier struct {
	DurableName       string `json:"DurableName,omitempty"`
	DurableNameFormat string `json:"DurableNameFormat,omitempty"`
}

// Location holds the location information
type Location struct {
	AltitudeMeters int            `json:"AltitudeMeters,omitempty"`
	Latitude       int            `json:"Latitude,omitempty"`
	Longitude      int            `json:"Longitude,omitempty"`
	Contacts       *Contacts      `json:"Contacts,omitempty"`
	Oem            *Oem           `json:"Oem,omitempty"`
	PartLocation   *PartLocation  `json:"PartLocation,omitempty"`
	Placement      *Placement     `json:"Placement,omitempty"`
	PostalAddress  *PostalAddress `json:"PostalAddress,omitempty"`
}

// PartLocation holds the part location information
type PartLocation struct {
	Orientation          string `json:"Orientation,omitempty"`
	Reference            string `json:"Reference,omitempty"`
	LocationOrdinalValue int    `json:"LocationOrdinalValue"`
	LocationType         string `json:"LocationType,omitempty"`
	ServiceLabel         string `json:"ServiceLabel,omitempty"`
}

// Contacts holds the Contacts information
type Contacts struct {
	ContactName  string `json:"ContactName,omitempty"`
	EmailAddress string `json:"EmailAddress,omitempty"`
	PhoneNumber  string `json:"PhoneNumber,omitempty"`
}

// Placement holds the Placement information
type Placement struct {
	AdditionalInfo  string `json:"AdditionalInfo,omitempty"`
	Rack            string `json:"Rack,omitempty"`
	RackOffset      int    `json:"RackOffset,omitempty"`
	RackOffsetUnits string `json:"RackOffsetUnits,omitempty"`
	Row             string `json:"Row,omitempty"`
}

// PostalAddress holds the PostalAddress information
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

// IOStatistics represent IO statistics.
type IOStatistics struct {
	NonIORequests      int    `json:"NonIORequests,omitempty"`
	NonIORequestTime   string `json:"NonIORequestTime,omitempty"`
	ReadHitIORequests  int    `json:"ReadHitIORequests,omitempty"`
	ReadIOKiBytes      int    `json:"ReadIOKiBytes,omitempty"`
	ReadIORequests     int    `json:"ReadIORequests,omitempty"`
	ReadIORequestTime  string `json:"ReadIORequestTime,omitempty"`
	WriteHitIORequests int    `json:"WriteHitIORequests,omitempty"`
	WriteIOKiBytes     int    `json:"WriteIOKiBytes,omitempty"`
	WriteIORequests    int    `json:"WriteIORequests,omitempty"`
	WriteIORequestTime string `json:"WriteIORequestTime,omitempty"`
}
