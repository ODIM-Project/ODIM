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
	Oid               string      `json:"@odata.id"`
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
