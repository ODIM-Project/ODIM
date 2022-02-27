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

//Package dpresponse ...
package dpresponse

//ManagersCollection for Plugin
type ManagersCollection struct {
	OdataContext string `json:"@odata.context"`
	Etag         string `json:"@odata.etag,omitempty"`
	OdataID      string `json:"@odata.id"`
	OdataType    string `json:"@odata.type"`
	Description  string `json:"Description"`
	Name         string `json:"Name"`
	Members      []Link `json:"Members"`
	MembersCount int    `json:"Members@odata.count"`
}

// Link stores value of the links
type Link struct {
	Oid string `json:"@odata.id"`
}

// Manager struct for manager deta
type Manager struct {
	OdataContext    string         `json:"@odata.context"`
	Etag            string         `json:"@odata.etag,omitempty"`
	OdataID         string         `json:"@odata.id"`
	OdataType       string         `json:"@odata.type"`
	Name            string         `json:"Name"`
	ManagerType     string         `json:"ManagerType"`
	ID              string         `json:"Id"`
	UUID            string         `json:"UUID"`
	FirmwareVersion string         `json:"FirmwareVersion"`
	Status          *ManagerStatus `json:"Status,omitempty"`
}

// ManagerStatus struct is to define the status of the manager
type ManagerStatus struct {
	State string `json:"State"`
}
