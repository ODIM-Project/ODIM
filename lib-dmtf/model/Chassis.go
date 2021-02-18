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

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

// Chassis redfish structure
type Chassis struct {
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
	IndicatorLED       string           `json:"IndicatorLED"`
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

// Location redfish structure
type Location struct {
	Oid string `json:"@odata.id"`
}

// LogServices get
/*
/redfish/v1/Managers/{ManagerId}/LogServices/{LogServiceId}
/redfish/v1/Systems/{ComputerSystemId}/LogServices/{LogServiceId}
*/
type LogServices struct {
	Oid                 string  `json:"@odata.id"`
	Ocontext            string  `json:"@odata.context,omitempty"`
	Otype               string  `json:"@odata.type,omitempty"`
	Oetag               string  `json:"@odata.etag,omitempty"`
	ID                  string  `json:"Id,omitempty"`
	Description         string  `json:"Description,omitempty"`
	Name                string  `json:"Name,omitempty"`
	DateTime            string  `json:"DateTime,omitempty"`
	DateTimeLocalOffset string  `json:"DateTimeLocalOffset,omitempty"`
	Entries             Entries `json:"Entries,omitempty"`
	LogEntryType        string  `json:"LogEntryType,omitempty"`
	MaxNumberOfRecords  int     `json:"MaxNumberOfRecords,omitempty"`
	OverWritePolicy     string  `json:"OverWritePolicy,omitempty"`
	ServiceEnabled      bool    `json:"ServiceEnabled,omitempty"`
	Status              Status  `json:"Status,omitempty"`
}

//Entries redfish structure
type Entries struct {
	Oid string `json:"@odata.id"`
}

// Assembly redfish structure
type Assembly struct {
	Oid string `json:"@odata.id"`
}

// NetworkAdapters redfish structure
type NetworkAdapters struct {
	Oid string `json:"@odata.id"`
}

// PCIeSlots redfish structure
type PCIeSlots struct {
	Oid string `json:"@odata.id"`
}

// PhysicalSecurity redfish structure
type PhysicalSecurity struct {
	IntrusionSensor       string
	IntrusionSensorNumber int
	IntrusionSensorReArm  string
}

// Power redfish structure
type Power struct {
	Oid string `json:"@odata.id"`
}

// Sensors redfish structure
type Sensors struct {
	Oid string `json:"@odata.id"`
}

// Status redfish structure
type Status struct {
	Oid          string `json:"@odata.id,omitempty"`
	Ocontext     string `json:"@odata.context,omitempty"`
	Oetag        string `json:"@odata.etag,omitempty"`
	Otype        string `json:"@odata.type,omitempty"`
	Description  string `json:"description,omitempty"`
	ID           string `json:"Id,omitempty"`
	Name         string `json:"Name,omitempty"`
	Health       string `json:"Health,omitempty"`
	HealthRollup string `json:"HealthRollup,omitempty"`
	State        string `json:"State,omitempty"`
	Oem          *Oem   `json:"Oem,omitempty"`
}

// Thermal redfish structure
type Thermal struct {
	Oid string `json:"@odata.id"`
}

// SaveInMemory will create the Chassis in inmemory DB, with key as UUID
// Takes:
//	none as function parameter, but takes c of type *Chassis as a pointer receiver implicitly.
// Returns:
//	err of type error
//
//	On Sucess  - returns nil value
//	On Failure - returns non nil value
func (c *Chassis) SaveInMemory(deviceUUID string) *errors.Error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to connect to DB: ", err.Error())
	}
	if err := connPool.Create("chassis", deviceUUID, c); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to create new chassis: ", err.Error())
	}
	return nil
}
