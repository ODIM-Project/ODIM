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

//Package rfpmessagebus ...
package rfpmessagebus

import (
	"encoding/json"
	"fmt"
	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/plugin-redfish/config"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpmodel"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
)

// Publish function will handle events request in two originofcondition format
// originofcondition can be with or without @odata.id
func Publish(data interface{}) bool {
	if data == nil {
		log.Error("Nil data passed to event publisher")
		return false
	}
	event := data.(common.Events)
	if event.EventType == "PluginStartUp" {
		return writeToMessageBus(event)
	}

	// this map is to check what type of event is recieved.
	// if @odata.type contains MetricReport then its of type MetricReport message objects else event message objects
	var metricReportEventDataMap map[string]interface{}
	err := json.Unmarshal(event.Request, &metricReportEventDataMap)
	if err != nil {
		log.Error("Failed to unmarshal the event, got: " + err.Error())
		return false
	}
	odatatype := metricReportEventDataMap["@odata.type"]
	if odatatype != nil && strings.Contains(odatatype.(string), "MetricReport") {
		event, err = formatMetricReportEventRequest(event)
	} else {
		event, err = formatRedfishEventRequest(event)
	}
	if err != nil {
		return false
	}
	if !writeToMessageBus(event) {
		return false
	}
	log.Info("Forwarded event: " + string(event.Request))
	log.Info("Event Published")
	return true
}

// translateOriginOfCondition changes the ilo specific uri to  redfsih std uri
func translateOriginOfCondition(eventOriginOfCondition string) string {
	originOfCondition := strings.Replace(eventOriginOfCondition, "systems", "Systems", -1)
	originOfCondition = strings.Replace(originOfCondition, "bios", "Bios", -1)
	originOfCondition = strings.Replace(originOfCondition, "settings", "Settings", -1)
	originOfCondition = strings.Replace(originOfCondition, "SmartStorage", "Storage", -1)
	originOfCondition = strings.Replace(originOfCondition, "DiskDrives", "Drives", -1)
	originOfCondition = strings.Replace(originOfCondition, "LogicalDrives", "Volumes", -1)
	originOfCondition = strings.Replace(originOfCondition, "PCIDevices", "PCIeDevices", -1)
	return strings.Replace(originOfCondition, "/ArrayControllers/", "/ArrayControllers-", -1)
}

// translateOriginOfCondition to changes Oem block  ilo specific uri to  redfsih std uri
func translateOemData(oem interface{}) interface{} {
	oemData, _ := json.Marshal(oem)
	updateData := translateOriginOfCondition(string(oemData))
	var updateOemData interface{}
	json.Unmarshal([]byte(updateData), updateOemData)
	return updateOemData
}

func formatMetricReportEventRequest(eventRequest common.Events) (common.Events, error) {
	var event common.Events
	// prepare the device data
	var devices []rfpmodel.Device
	rfpmodel.GetAllDevicesInInventory(&devices)

	systemUUID := getMatchingDeviceUUID(eventRequest.IP, devices)
	if systemUUID == "" {
		return event, fmt.Errorf("System id is empty")
	}

	// to replace id in system
	updatedData := strings.Replace(string(eventRequest.Request), "/redfish/v1/Systems/", "/redfish/v1/Systems/"+systemUUID+":", -1)
	updatedData = strings.Replace(updatedData, "/redfish/v1/systems/", "/redfish/v1/Systems/"+systemUUID+":", -1)
	// to replace id in chassis
	updatedData = strings.Replace(updatedData, "/redfish/v1/Chassis/", "/redfish/v1/Chassis/"+systemUUID+":", -1)
	updatedData = strings.Replace(updatedData, "/redfish/v1/chassis/", "/redfish/v1/Chassis/"+systemUUID+":", -1)

	event.Request, _ = json.Marshal(updatedData)
	event.IP = eventRequest.IP
	event.EventType = "MetricReport"

	return event, nil
}

func formatRedfishEventRequest(eventRequest common.Events) (common.Events, error) {
	// Since we are deleting the first event from the eventlist,
	// processing the first event
	var message common.MessageData
	err := json.Unmarshal(eventRequest.Request, &message)
	if err != nil {
		log.Error("Failed to unmarshal the event: " + err.Error())
		return eventRequest, err
	}
	var messageData dmtf.Event
	message.Context = messageData.Context
	message.Name = messageData.Name
	message.OdataType = messageData.ODataType
	message.Events = make([]common.Event, 0)
	for i := 0; i < len(messageData.Events); i++ {
		var eventData common.Event
		eventData.EventGroupID = messageData.Events[i].EventGroupID
		eventData.EventID = messageData.Events[i].EventID
		eventData.EventTimestamp = messageData.Events[i].EventTimestamp
		eventData.EventType = messageData.Events[i].EventType
		eventData.Message = messageData.Events[i].Message
		eventData.MemberID = messageData.Events[i].MemberID
		eventData.Severity = messageData.Events[i].Severity
		eventData.Oem = messageData.Events[i].Oem
		eventData.MessageID = messageData.Events[i].MessageID
		eventData.OriginOfCondition = &common.Link{
			Oid: messageData.Events[i].OriginOfCondition.Oid,
		}
		eventData.MessageArgs = messageData.Events[i].MessageArgs
		message.Events = append(message.Events, eventData)
	}
	eventRequest.Request, _ = json.Marshal(message)
	return eventRequest, nil
}

func getMatchingDeviceUUID(serverIP string, devices []rfpmodel.Device) string {
	if len(devices) < 1 {
		return ""
	}

	front := 0
	rear := len(devices) - 1
	for front <= rear {
		host := getIPFromHostName(devices[front].Host)
		if host == serverIP {
			return devices[front].SystemID
		}

		host = getIPFromHostName(devices[rear].Host)
		if host == serverIP {
			return devices[rear].SystemID
		}
		front++
		rear--
	}
	return ""
}

// getIPFromHostName - look up the ip from the fqdn
func getIPFromHostName(fqdn string) string {
	host, _, err := net.SplitHostPort(fqdn)
	if err != nil {
		host = fqdn
	}
	addr, err := net.LookupIP(host)
	if err != nil || len(addr) < 1 {
		return ""
	}
	return fmt.Sprintf("%v", addr[0])
}

func writeToMessageBus(events common.Events) bool {
	K, err := dc.Communicator(dc.KAFKA, config.Data.MessageBusConf.MessageQueueConfigFilePath)
	if err != nil {
		log.Error("Unable communicate with kafka, got: " + err.Error())
		return false
	}
	defer K.Close()

	topic := config.Data.MessageBusConf.EmbQueue[0]
	if err := K.Distribute(topic, events); err != nil {
		log.Error("Unable Publish events to kafka, got: " + err.Error())
		return false
	}
	return true
}
