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

//Package dpmodel ...
package dpmodel

// Volume holds the northbound request body
type Volume struct {
	Name               string        `json:"Name" validate:"required"`
	RAIDType           string        `json:"RAIDType"`
	Drives             []OdataIDLink `json:"Drives"`
	OperationApplyTime string        `json:"@Redfish.OperationApplyTime"`
}

// OdataIDLink contains link to a resource
type OdataIDLink struct {
	OdataID string `json:"@odata.id"`
}

//VolumesCollection data
type VolumesCollection struct {
	OdataContext string        `json:"@odata.context"`
	OdataID      string        `json:"@odata.id"`
	OdataType    string        `json:"@odata.type"`
	Description  string        `json:"Description"`
	Name         string        `json:"Name"`
	Members      []OdataIDLink `json:"Members"`
	MembersCount int           `json:"Members@odata.count"`
}

//FirmwareVersion contains the firmware version of server
type FirmwareVersion struct {
	FirmwareVersion string `json:"FirmwareVersion"`
}
