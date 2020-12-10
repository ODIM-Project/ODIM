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

// Package events have the functionality of
// - Create Event Subscription
// - Delete Event Subscription
// - Get Event Subscription
// - Post Event Subscription to destination
// - Post TestEvent (SubmitTestEvent)
// and corresponding unit test cases
package events

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	fabricproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
)

// ForwardEventMessageData contains information of Events and message details including arguments
// it will be send as byte stream on the wire to/from kafka
type ForwardEventMessageData struct {
	OdataType string         `json:"@odata.type"`
	Name      string         `json:"Name"`
	Context   string         `json:"@odata.context"`
	Events    []ForwardEvent `json:"Events"`
}

// ForwardEvent contains the details of the event to be forwarded
type ForwardEvent struct {
	MemberID          string      `json:"MemberId,omitempty"`
	EventType         string      `json:"EventType"`
	EventGroupID      int         `json:"EventGroupId,omitempty"`
	EventID           string      `json:"EventId"`
	Severity          string      `json:"Severity"`
	EventTimestamp    string      `json:"EventTimestamp"`
	Message           string      `json:"Message"`
	MessageArgs       []string    `json:"MessageArgs,omitempty"`
	MessageID         string      `json:"MessageId"`
	Oem               interface{} `json:"Oem,omitempty"`
	OriginOfCondition string      `json:"OriginOfCondition,omitempty"`
}

// checkAndAddFabrics this function is used to cross check if an event is for fabrics resource added
// if yes then it will add the new fabric resource to db
func checkAndAddFabrics(requestData string, host string) {
	if strings.Contains(requestData, "/redfish/v1/Fabrics") && strings.Contains(requestData, "ResourceAdded") {
		message, parseStatus := parseEventData(requestData)
		if !parseStatus {
			log.Println("Error while trying to parse request data ", requestData)
			return
		}
		for _, inEvent := range message.Events {
			if len(inEvent.OriginOfCondition) < 1 {
				log.Printf("event not forwarded : Originofcondition is empty in incoming event with body %v\n", requestData)
				continue
			}
			if strings.EqualFold(inEvent.EventType, "ResourceAdded") &&
				strings.HasPrefix(inEvent.OriginOfCondition, "/redfish/v1/Fabrics") {
				addFabricRPCCall(inEvent.OriginOfCondition, host)
			}
		}
	}
}

// PublishEventsToDestination This method sends the event/alert to subscriber's destination
// Takes:
// 	data of type interface{}
//Returns:
//	bool: return false if any error occurred during execution, else returns true
func PublishEventsToDestination(data interface{}) bool {

	if data == nil {
		log.Printf("error: invalid input params")
		return false
	}
	// Extract the Hostname/IP of the event source and Event from input parameter
	event := data.(common.Events)
	host := event.IP

	var requestData = string(event.Request)
	//replacing the resposne with north bound translation URL
	for key, value := range config.Data.URLTranslation.NorthBoundURL {
		requestData = strings.Replace(requestData, key, value, -1)
	}

	var flag bool
	var uuid string

	checkAndAddFabrics(requestData, host)

	deviceSubscription, err := evmodel.GetDeviceSubscriptions(host)
	if err != nil {
		log.Printf("error: Failed to get the event destinations: %v\n", err)
		return false
	}

	if len(deviceSubscription.OriginResources) < 1 {
		log.Println("error: no origin resources found in device subscriptions")
		return false
	}

	requestData, uuid = formatEvent(requestData, deviceSubscription.OriginResources[0], host)
	message, parseStatus := parseEventData(requestData)
	if !parseStatus {
		log.Println("Error while trying to parse the input request data", message)
		return false
	}

	subscriptions, err := evmodel.GetEvtSubscriptions(host)
	if err != nil {
		return false
	}

	eventMap := make(map[string][]ForwardEvent)
	for _, inEvent := range message.Events {
		if len(inEvent.OriginOfCondition) < 1 {
			log.Printf("event not forwarded : Originofcondition is empty in incoming event with body %v\n", requestData)
			continue
		}

		var resTypePresent bool
		originofCond := strings.Split(strings.TrimSuffix(inEvent.OriginOfCondition, "/"), "/")
		resType := originofCond[len(originofCond)-2]
		for _, value := range common.ResourceTypes {
			if strings.Contains(resType, value) {
				resTypePresent = true
			}
		}

		if !resTypePresent {
			log.Printf("event not forwared: resource type of originofcondition not supported in event with body: %v\n", message)
			continue
		}

		for _, sub := range subscriptions {

			// filter and send evemts to destination if destination is not empty
			// in case of default event subscription destination will be empty
			if sub.Destination != "" {
				// check if hostip present in the hosts slice to make sure that it doesn't filter with the destination ip
				if isHostPresent(sub.Hosts, host) {
					if filterEventsToBeForwarded(sub, inEvent, deviceSubscription.OriginResources) {
						eventMap[sub.Destination] = append(eventMap[sub.Destination], inEvent)
						flag = true
					}
				} else {
					log.Println("event not forwarded : No subscription for the incoming event's originofcondition")
					flag = false
				}

			}
		}

		if strings.EqualFold("Alert", inEvent.EventType) {
			if strings.Contains(inEvent.MessageID, "ServerPostDiscoveryComplete") || strings.Contains(inEvent.MessageID, "ServerPostComplete") {
				go rediscoverSystemInventory(uuid, inEvent.OriginOfCondition)
				flag = true
			}
			if strings.Contains(inEvent.MessageID, "ServerPoweredOn") || strings.Contains(inEvent.MessageID, "ServerPoweredOff") {
				go updateSystemPowerState(uuid, inEvent.OriginOfCondition, inEvent.MessageID)
				flag = true
			}
		} else if strings.EqualFold("ResourceAdded", message.Events[0].EventType) || strings.EqualFold("ResourceRemoved", message.Events[0].EventType) {
			if strings.Contains(message.Events[0].OriginOfCondition, "Volumes") {
				s := strings.Split(message.Events[0].OriginOfCondition, "/")
				storageURI := fmt.Sprintf("/%s/%s/%s/%s/%s/", s[1], s[2], s[3], s[4], s[5])
				go rediscoverSystemInventory(uuid, storageURI)
				flag = true
			}
		}
	}

	for key, value := range eventMap {
		message.Events = value
		data, err := json.Marshal(message)
		if err != nil {
			log.Printf("unable to converts event into bytes: %v", err)
			continue
		}
		go postEvent(key, data)
	}
	return flag
}

func filterEventsToBeForwarded(subscription evmodel.Subscription, event ForwardEvent, originResources []string) bool {
	eventTypes := subscription.EventTypes
	messageIds := subscription.MessageIds
	resourceTypes := subscription.ResourceTypes
	originCondition := strings.TrimSuffix(event.OriginOfCondition, "/")
	if (len(eventTypes) == 0 || isStringPresentInSlice(eventTypes, event.EventType, "event type")) &&
		(len(messageIds) == 0 || isStringPresentInSlice(messageIds, event.MessageID, "message id")) &&
		(len(resourceTypes) == 0 || isResourceTypeSubscribed(resourceTypes, event.OriginOfCondition, subscription.SubordinateResources)) {
		// if SubordinateResources is true then check if originofresource is top level of originofcondition
		// if SubordinateResources is flase then check originofresource is same as originofcondition
		for _, origin := range originResources {
			if subscription.SubordinateResources == true {
				if strings.Contains(originCondition, origin) {
					return true
				}
			} else {
				if origin == originCondition {
					return true
				}
			}
		}
	}
	log.Println("Event not forwarded : No subscription for the incoming event's originofcondition")
	return false
}

// formatEvent will format the event string according to the odimra
// add uuid:systemid/chassisid inplace of systemid/chassisid
func formatEvent(event, originResource, hostIP string) (string, string) {
	uuid, _ := getUUID(originResource)
	if !strings.Contains(hostIP, "Collection") {
		str := "/redfish/v1/Systems/" + uuid + ":"
		event = strings.Replace(event, "/redfish/v1/Systems/", str, -1)
		str = "/redfish/v1/systems/" + uuid + ":"
		event = strings.Replace(event, "/redfish/v1/systems/", str, -1)
		str = "/redfish/v1/Chassis/" + uuid + ":"
		event = strings.Replace(event, "/redfish/v1/Chassis/", str, -1)
		str = "/redfish/v1/Managers/" + uuid + ":"
		event = strings.Replace(event, "/redfish/v1/Managers/", str, -1)
	}
	return event, uuid
}

func isResourceTypeSubscribed(resourceTypes []string, originOfCondition string, subordinateResources bool) bool {
	//If the incoming odata type field empty then return true
	if originOfCondition == "" {
		return true
	}
	originCond := strings.Split(strings.TrimSuffix(originOfCondition, "/"), "/")

	for _, resourceType := range resourceTypes {
		res := common.ResourceTypes[resourceType]
		if subordinateResources {

			// if subordinateResources is true then first check the child resourcetype is present in db.
			// if its there then return true
			// if its not then check collection resource type
			// Ex : originofcondition:/redfish/v1/Systems/uuid:1/processors/1
			// child resource type would be processors (index-2)
			// collection resource type would be Systems (index-4)
			resType := originCond[len(originCond)-2]
			if strings.Contains(res, resType) {
				return true
			}
			resType = originCond[len(originCond)-4]
			if strings.Contains(res, resType) {
				return true
			}
		} else {
			// if subordinateResources is false then check the child/collection resourcetype is present in db.
			resType := originCond[len(originCond)-2]
			if strings.Contains(resType, res) {
				return true
			}
		}
	}
	log.Println("Event not forwarded : No subscription for the incoming event's originofcondition")
	return false
}

func isStringPresentInSlice(slice []string, str, message string) bool {
	//If the incoming event fields contains empty values return true
	if str == "" {
		return true
	}
	for _, value := range slice {
		if value == str {
			return true
		}
	}
	log.Printf("Event not forwarded : No subscription for the incoming event's  %s", message)
	return false
}

// postEvent will post the event to destination
func postEvent(destination string, event []byte) {
	httpConf := &config.HTTPConfig{
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
	}
	httpClient, err := httpConf.GetHTTPClientObj()
	if err != nil {
		log.Println("error: failed to get http client object:", err)
		return
	}
	req, err := http.NewRequest("POST", destination, bytes.NewBuffer(event))
	if err != nil {
		log.Printf("error while getting new http request:%v", err)
		return
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	var resp *http.Response
	count := evcommon.DeliveryRetryAttempts + 1
	for i := 0; i < count; i++ {
		config.TLSConfMutex.Lock()
		resp, err = httpClient.Do(req)
		if err == nil {
			resp.Body.Close()
			log.Printf("event post response: %v", resp)
			config.TLSConfMutex.Unlock()
			return
		}
		config.TLSConfMutex.Unlock()
		log.Println("Retrying event posting")
		time.Sleep(time.Second * evcommon.DeliveryRetryIntervalSeconds)
	}
	log.Printf("error while make https call to send the event:%v", err)
	return
}

// rediscoverSystemInventory will be triggered when ever the System Restart or Power On
// event is detected it will create a rpc for aggregation which will delete all system inventory //
// and rediscover all of them
func rediscoverSystemInventory(systemID, systemURL string) {
	systemURL = strings.TrimSuffix(systemURL, "/")
	aggregator := aggregatorproto.NewAggregatorService(services.Aggregator, services.Service.Client())
	_, err := aggregator.RediscoverSystemInventory(context.TODO(), &aggregatorproto.RediscoverSystemInventoryRequest{
		SystemID:  systemID,
		SystemURL: systemURL,
	})
	if err != nil {
		log.Println("Error while rediscoverSystemInventroy")
		return
	}
	log.Println("info: rediscovery of system and chasis started.")
	return
}

func addFabricRPCCall(origin, address string) {
	if strings.Contains(origin, "Zones") || strings.Contains(origin, "Endpoints") || strings.Contains(origin, "AddressPools") {
		return
	}
	fab := fabricproto.NewFabricsService(services.Fabrics, services.Service.Client())

	_, err := fab.AddFabric(context.TODO(), &fabricproto.AddFabricRequest{
		OriginResource: origin,
		Address:        address,
	})
	if err != nil {
		log.Println("Error while AddFabric", err)
		return
	}
	p := PluginContact{
		ContactClient: pmbhandle.ContactPlugin,
	}
	p.checkCollectionSubscription(origin, "Redfish")
	log.Println("info: Fabric Added")
	return
}

// updateSystemPowerState will be triggered when ever the System Powered Off event is received
// When event is detected a rpc is created for aggregation which will update the system inventory
func updateSystemPowerState(systemUUID, systemURI, state string) {

	systemURI = strings.TrimSuffix(systemURI, "/")

	index := strings.LastIndex(systemURI, "/")
	uri := systemURI[:index]
	id := systemURI[index+1:]

	if strings.ContainsAny(id, ":/-") {
		log.Println("error: event contains invalid origin of condition -", systemURI)
		return
	}
	if strings.Contains(state, "ServerPoweredOn") {
		state = "On"
	} else {
		state = "Off"
	}
	aggregator := aggregatorproto.NewAggregatorService(services.Aggregator, services.Service.Client())
	_, err := aggregator.UpdateSystemState(context.TODO(), &aggregatorproto.UpdateSystemStateRequest{
		SystemUUID: systemUUID,
		SystemID:   id,
		SystemURI:  uri,
		UpdateKey:  "PowerState",
		UpdateVal:  state,
	})
	if err != nil {
		log.Println("error: system power state update failed with", err)
		return
	}
	log.Println("info: system power state update initiated")
	return
}

// parseEventData is used to parse the input request data and based on the input structure
// of originofcondition further requests will be created
func parseEventData(requestData string) (ForwardEventMessageData, bool) {
	var forwardEventData ForwardEventMessageData
	var message common.MessageData
	err := json.Unmarshal([]byte(requestData), &forwardEventData)
	if err != nil {
		if err = json.Unmarshal([]byte(requestData), &message); err != nil {
			log.Printf("error: Failed to unmarshal the event: %v", err)
			log.Println("incoming event:", requestData)
			return forwardEventData, false
		}
		// map the std event format to forwarding event format
		forwardEventData.Context = message.Context
		forwardEventData.Name = message.Name
		forwardEventData.OdataType = message.OdataType
		var events []ForwardEvent
		for i := 0; i < len(forwardEventData.Events); i++ {
			var event ForwardEvent
			event.EventGroupID = message.Events[i].EventGroupID
			event.EventID = message.Events[i].EventID
			event.EventTimestamp = message.Events[i].EventTimestamp
			event.EventType = message.Events[i].EventType
			event.MemberID = message.Events[i].MemberID
			event.Message = message.Events[i].Message
			event.MessageArgs = message.Events[i].MessageArgs
			event.MessageID = message.Events[i].MessageID
			event.Oem = message.Events[i].Oem
			event.Severity = message.Events[i].Severity
			if message.Events[i].OriginOfCondition != nil {
				event.OriginOfCondition = message.Events[i].OriginOfCondition.Oid
			}
			events = append(events, event)
		}
		forwardEventData.Events = events
	}
	return forwardEventData, true
}
