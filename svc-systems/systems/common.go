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

//Package systems ...
package systems

// BiosSetting structure for checking request body case
type BiosSetting struct {
	OdataContext      string      `json:"@odata.context"`
	OdataID           string      `json:"@odata.id"`
	Odatatype         string      `json:"@odata.type"`
	ID                string      `json:"Id"`
	Name              string      `json:"Name"`
	AttributeRegistry string      `json:"AttributeRegistry"`
	Attributes        interface{} `json:"Attributes"`
}

// BootOrderSettings structure for checking request body case
type BootOrderSettings struct {
	Boot Boot `json:"Boot"`
}

// Boot structure for checking request body case in BootOrderSettings
type Boot struct {
	BootOrder                    []string `json:"BootOrder"`
	BootSourceOverrideEnabled    string   `json:"BootSourceOverrideEnabled"`
	BootSourceOverrideMode       string   `json:"BootSourceOverrideMode"`
	BootSourceOverrideTarget     string   `json:"BootSourceOverrideTarget"`
	UefiTargetBootSourceOverride string   `json:"UefiTargetBootSourceOverride"`
}

// ResetComputerSystem structure for checking request body case
type ResetComputerSystem struct {
	ResetType string `json:"ResetType"`
}
