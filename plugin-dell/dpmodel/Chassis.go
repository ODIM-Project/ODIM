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

// ChassisDevice struct definition
type ChassisDevice struct {
	Ocontext           string           `json:"@odata.context"`
	Oid                string           `json:"@odata.id"`
	Otype              string           `json:"@odata.type"`
	Oetag              string           `json:"@odata.etag,omitempty"`
	ID                 string           `json:"Id"`
	Description        string           `json:"Description"`
	Name               string           `json:"Name"`
	AssetTag           string           `json:"AssetTag"`
	ChassisType        string           `json:"ChassisType"`
	DepthMm            int              `json:"DepthMm"`
	EnvironmentalClass string           `json:"EnvironmentalClass"`
	HeightMm           int              `json:"HeightMm"`
	IndicatorLED       int              `json:"IndicatorLED"`
	Manufacturer       string           `json:"Manufacturer"`
	Model              string           `json:"Model"`
	PartNumber         string           `json:"PartNumber"`
	PowerState         string           `json:"PowerState"`
	SerialNumber       string           `json:"SerialNumber"`
	SKU                string           `json:"SKU"`
	UUID               string           `json:"UUID"`
	WeightKg           int              `json:"WeightKg"`
	WidthMm            int              `json:"WidthMm"`
	Links              Links            `json:"Links"`
	Location           Location         `json:"Location"`
	LogServices        LogServices      `json:"LogServices"`
	Assembly           Assembly         `json:"Assembly"`
	NetworkAdapters    NetworkAdapters  `json:"NetworkAdapters"`
	PCIeSlots          PCIeSlots        `json:"PCIeSlots"`
	PhysicalSecurity   PhysicalSecurity `json:"PhysicalSecurity"`
	Power              Power            `json:"Power"`
	Sensors            Sensors          `json:"Sensors"`
	Status             Status           `json:"Status"`
	Thermal            Thermal          `json:"Thermal"`
}

/*
type Links struct {
	ComputerSystems    []ComputerSystems
	ContainedBy        []ContainedBy
	Contains           []Contains
	CooledBy           []CooledBy
	Drives             []Drives
	ManagedBy          []ManagedBy
	ManagersInChassis  []ManagersInChassis
	Oem                Oem
	PCIeDevices        []PCIeDevices
	PoweredBy          []PoweredBy
	Processors         []Processors
	ResourceBlocks     []ResourceBlocks
	Storage            []Storage
	Switches           []Switches
}
*/

// Location get ..
type Location struct {
	Oid string `json:"@odata.id"`
}

// LogServices get
type LogServices struct {
	Oid string `json:"@odata.id"`
}

// Assembly get
type Assembly struct {
	Oid string `json:"@odata.id"`
}

// NetworkAdapters get
type NetworkAdapters struct {
	Oid string `json:"@odata.id"`
}

// PCIeSlots get
type PCIeSlots struct {
	Oid string `json:"@odata.id"`
}

// PhysicalSecurity get
type PhysicalSecurity struct {
	IntrusionSensor       string `json:"IntrusionSensor"`
	IntrusionSensorNumber int    `json:"IntrusionSensorNumber"`
	IntrusionSensorReArm  string `json:"IntrusionSensorReArm"`
}

// Power get
type Power struct {
	Oid string `json:"@odata.id"`
}

// Sensors get
type Sensors struct {
	Oid string `json:"@odata.id"`
}

// Status get
type Status struct {
	Health       string  `json:"Health"`
	HealthRollup *string `json:",omitempty"`
	State        string  `json:"State"`
	Oem          Oem     `json:"Oem"`
}

// Thermal get
type Thermal struct {
	Oid string `json:"@odata.id"`
}

/*
type ComputerSystems struct {
        Oid string `json:"@odata.id"`
}
type ContainedBy struct {
        Oid string `json:"@odata.id"`
}
type Contains struct {
        Oid string `json:"@odata.id"`
}
type CooledBy struct {
        Oid string `json:"@odata.id"`
}
type Drives struct {
        Oid string `json:"@odata.id"`
}
type ManagedBy struct {
        Oid string `json:"@odata.id"`
}
type ManagersInChassis struct {
        Oid string `json:"@odata.id"`
}
type Oem struct {
}
type PCIeDevices struct {
        Oid string `json:"@odata.id"`
}
type PoweredBy struct {
        Oid string `json:"@odata.id"`
}
type Processors struct {
        Oid string `json:"@odata.id"`
}
type ResourceBlocks struct {
        Oid string `json:"@odata.id"`
}
type Storage struct {
        Oid string `json:"@odata.id"`
}
type Switches struct {
        Oid string `json:"@odata.id"`
}
*/
