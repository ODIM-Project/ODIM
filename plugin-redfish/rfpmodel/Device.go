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

//Package rfpmodel ...
package rfpmodel

import "sync"

//Device struct definition
type Device struct {
	Host     string `json:"ManagerAddress"`
	Username string `json:"UserName"`
	Password []byte `json:"Password"`
	PostBody []byte `json:"PostBody"`
	Location string `json:"Location"`
	SystemID string `json:"SystemID"`
}

//EvtSubPost ...
type EvtSubPost struct {
	Destination     string        `json:"Destination"`
	EventTypes      []string      `json:"EventTypes,omitempty"`
	MessageIds      []string      `json:"MessageIds,omitempty"`
	ResourceTypes   []string      `json:"ResourceTypes,omitempty"`
	HTTPHeaders     []HTTPHeaders `json:"HttpHeaders"`
	Context         string        `json:"Context"`
	Protocol        string        `json:"Protocol"`
	EventFormatType string        `json:"EventFormatType"`
}

//HTTPHeaders ...
type HTTPHeaders struct {
	ContentType string `json:"Content-Type"`
}

// StartUpData holds the required data for plugin startup
type StartUpData struct {
	RequestType           string                `json:"RequestType"`
	ResyncEvtSubscription bool                  `json:"ResyncEvtSubscription"`
	Devices               map[string]DeviceData `json:"Devices"`
}

// DeviceInventory is for storing the device inventory
var DeviceInventory *DeviceInventoryData

// DeviceInventoryData holds the list of all managed devices
type DeviceInventoryData struct {
	mutex  *sync.RWMutex
	Device map[string]DeviceData
}

// DeviceData holds device credentials, event subcription and trigger details
type DeviceData struct {
	UserName              string                 `json:"UserName"`
	Password              []byte                 `json:"Password"`
	Address               string                 `json:"Address"`
	Operation             string                 `json:"Operation"`
	EventSubscriptionInfo *EventSubscriptionInfo `json:"EventSubscriptionInfo"`
	TriggerInfo           *TriggerInfo           `json:"TriggerInfo"`
}

// EventSubscriptionInfo holds the event subscription details of a device
type EventSubscriptionInfo struct {
	EventTypes []string `json:"EventTypes"`
	Location   string   `json:"Location"`
}

// TriggerInfo holds the metric trigger info of a device
type TriggerInfo struct {
}

// init is for intializing global variables defined in this package
func init() {
	DeviceInventory = &DeviceInventoryData{
		mutex:  &sync.RWMutex{},
		Device: make(map[string]DeviceData),
	}
}

// AddDeviceToInventory is for adding new device to the inventory
// by acquiring write lock
func AddDeviceToInventory(uuid string, deviceData DeviceData) {
	DeviceInventory.mutex.Lock()
	defer DeviceInventory.mutex.Unlock()
	DeviceInventory.Device[uuid] = deviceData
	return
}

// DeleteDeviceInInventory is for deleting device in the inventory
// by acquiring write lock
func DeleteDeviceInInventory(uuid string) {
	DeviceInventory.mutex.Lock()
	defer DeviceInventory.mutex.Unlock()
	delete(DeviceInventory.Device, uuid)
	return
}

// GetAllDevicesInInventory is for getting list of all devices
// in the inventory by acquiring read lock
// Expects uninitialized slice of type Device
func GetAllDevicesInInventory(devices *[]Device) int {
	DeviceInventory.mutex.RLock()
	defer DeviceInventory.mutex.RUnlock()

	deviceCount := len(DeviceInventory.Device)
	*devices = make([]Device, deviceCount)
	count := 0
	for uuid, device := range DeviceInventory.Device {
		(*devices)[count] = Device{
			Host:     device.Address,
			Username: device.UserName,
			Password: device.Password,
			SystemID: uuid,
		}
		count++
	}
	return deviceCount
}
