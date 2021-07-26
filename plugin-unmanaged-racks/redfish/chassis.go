/*
 * Copyright (c) 2020 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package redfish

import "github.com/google/uuid"

// StatusEnabledOk is prototype of Redfish Status containing Enabled,OK values.
var StatusEnabledOk = Status{
	State:  "Enabled",
	Health: "OK",
}

// PowerStateOn is an enum for Chassis.PowerState property
var PowerStateOn = "On"

// Chassis is representation of Redfish Chassis asset
type Chassis struct {
	Ocontext     string `json:"@odata.context"`
	Oid          string `json:"@odata.id"`
	Otype        string `json:"@odata.type"`
	ID           string `json:"Id"`
	Description  string `json:"Description"`
	Name         string `json:"Name"`
	ChassisType  string `json:"ChassisType"`
	Links        Links  `json:"Links"`
	Model        string `json:"Model,omitempty"`
	Manufacturer string `json:"Manufacturer,omitempty"`
	PartNumber   string `json:"PartNumber,omitempty"`
	PowerState   string `json:"PowerState,omitempty"`
	SerialNumber string `json:"SerialNumber,omitempty"`
	Status       Status `json:"Status"`
}

// ShapeChassis function initializes given struct(ch) with values for ids, odata related fields
func ShapeChassis(ch *Chassis) *Chassis {
	ch.Otype = "#Chassis.v1_14_0.Chassis"
	ch.Ocontext = "/redfish/v1/$metadata#Chassis.Chassis"
	ch.ID = generateChassisID(ch.Name)
	ch.Oid = "/ODIM/v1/Chassis/" + ch.ID
	ch.Status = StatusEnabledOk
	ch.Links.ComputerSystems = []Link{}
	ch.PowerState = PowerStateOn
	return ch
}

func generateChassisID(name string) string {
	return uuid.NewSHA1(unmanagedChassisBaseUUID, []byte(name)).String()
}

var unmanagedChassisBaseUUID = uuid.Must(uuid.Parse("1bde942f-36f3-4e92-9b3b-4e497092430d"))

// Status is representation of Redfish Status
type Status struct {
	Health       string `json:"Health,omitempty"`
	HealthRollup string `json:"HealthRollup,omitempty"`
	State        string `json:"State,omitempty"`
}

// Links is representation of Redfish Chassis.Links
type Links struct {
	ComputerSystems []Link `json:""`
	ManagedBy       []Link `json:""`
	Contains        []Link `json:",omitempty"`
	ContainedBy     []Link `json:",omitempty"`
}
