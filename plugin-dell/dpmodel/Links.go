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

// Links .. this is Common in all resources
type Links struct {
	Chassis                  []Link `json:",omitempty"`
	ComputerSystems          []Link `json:",omitempty"`
	ConsumingComputerSystems []Link `json:",omitempty"`
	ContainedBy              []Link `json:",omitempty"`
	CooledBy                 []Link `json:",omitempty"`
	Endpoints                []Link `json:",omitempty"`
	Drives                   []Link `json:",omitempty"`
	ManagedBy                []Link `json:",omitempty"`
	Oem                      *Oem   `json:",omitempty"`
	ManagersInChassis        []Link `json:",omitempty"`
	PCIeDevices              []Link `json:",omitempty"`
	PoweredBy                []Link `json:",omitempty"`
	Processors               []Link `json:",omitempty"`
	ResourceBlocks           []Link `json:",omitempty"`
	Storage                  []Link `json:",omitempty"`
	SupplyingComputerSystems []Link `json:",omitempty"`
	Switches                 []Link `json:",omitempty"`
}

// Link get
type Link struct {
	Oid string `json:"@odata.id"`
}

// Oem get
type Oem struct {
}
