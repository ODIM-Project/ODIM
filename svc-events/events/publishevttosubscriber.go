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

// PublishEventsToDestination This method sends the event/alert to subscriber's destination
// Takes:
// 	data of type interface{}
//Returns:
//	none
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

	var message common.MessageData
	err := json.Unmarshal([]byte(requestData), &message)
	if err != nil {
		log.Printf("error: Failed to unmarshal the event: %v", err)
		log.Println("incoming event:", requestData)
		return false
	}
	if message.Events[0].OriginOfCondition == nil {
		log.Println("Event not forwarded : Originofcondition is empty in incoming event", requestData)
		return false
	}

	if len(message.Events[0].OriginOfCondition.Oid) < 1 {
		log.Println("Event not forwarded : Originofcondition is empty in incoming event", requestData)
		return false
	}
	if strings.EqualFold(message.Events[0].EventType, "ResourceAdded") &&
		strings.HasPrefix(message.Events[0].OriginOfCondition.Oid, "/redfish/v1/Fabrics") {
		addFabricRPCCall(message.Events[0].OriginOfCondition.Oid, host)
	}

	var resTypePresent bool
	originofCond := strings.Split(strings.TrimSuffix(message.Events[0].OriginOfCondition.Oid, "/"), "/")
	resType := originofCond[len(originofCond)-2]
	for _, value := range common.ResourceTypes {
		if strings.Contains(resType, value) {
			resTypePresent = true
		}
	}

	if !resTypePresent {
		log.Println("Event not forwared: resource type of originofcondition not supported in event", message)
		return false
	}

	var flag bool
	var uuid string

	deviceSubscription, err := evmodel.GetDeviceSubscriptions(host)
	if err != nil {
		log.Printf("error: Failed to get the event destinations: %v", err)
		return false
	}

	if len(deviceSubscription.OriginResources) < 1 {
		return false
	}

	originResource := deviceSubscription.OriginResources[0]
	var eventString string
	eventString, uuid = formatEvent(originResource, requestData, host)
	eventRequest := []byte(eventString)

	subscriptions, err := evmodel.GetEvtSubscriptions(host)
	if err != nil {
		return false
	}
	for _, sub := range subscriptions {

		// filter and send evemts to destination if destination is not empty
		// in case of default event subscription destination will be empty
		if sub.Destination != "" {
			// check if hostip present in the hosts slice to make sure that it doesn't filter with the destination ip
			if isHostPresent(sub.Hosts, host) {
				if filterEventsToBeForwarded(sub, eventRequest, deviceSubscription.OriginResources) {
					log.Printf("Destination: %v\n", sub.Destination)
					go postEvent(sub.Destination, eventRequest)
					flag = true
				}
			} else {
				log.Println("Event not forwarded : No subscription for the incoming event's originofcondition")
				flag = false
			}

		}
	}

	if strings.EqualFold("Alert", message.Events[0].EventType) {
		if strings.Contains(message.Events[0].MessageID, "ServerPostDiscoveryComplete") || strings.Contains(message.Events[0].MessageID, "ServerPostComplete") {
			go rediscoverSystemInventory(uuid, message.Events[0].OriginOfCondition.Oid)
			flag = true
		}
		if strings.Contains(message.Events[0].MessageID, "ServerPoweredOn") || strings.Contains(message.Events[0].MessageID, "ServerPoweredOff") {
			go updateSystemPowerState(uuid, message.Events[0].OriginOfCondition.Oid, message.Events[0].MessageID)
			flag = true
		}
	}
	return flag
}

func filterEventsToBeForwarded(subscription evmodel.Subscription, events []byte, originResources []string) bool {
	eventTypes := subscription.EventTypes
	messageIds := subscription.MessageIds
	resourceTypes := subscription.ResourceTypes
	var message common.MessageData
	err := json.Unmarshal(events, &message)
	if err != nil {
		log.Printf("error: Failed to unmarshal the event: %v", err)
		return false
	}
	originCondition := strings.TrimSuffix(message.Events[0].OriginOfCondition.Oid, "/")
	if (len(eventTypes) == 0 || isStringPresentInSlice(eventTypes, message.Events[0].EventType, "event type")) &&
		(len(messageIds) == 0 || isStringPresentInSlice(messageIds, message.Events[0].MessageID, "message id")) &&
		(len(resourceTypes) == 0 || isResourceTypeSubscribed(resourceTypes, message.Events[0].OriginOfCondition.Oid, subscription.SubordinateResources)) {
		// if SubordinateResources is true then check if originofresource is top level of originofcondition
		// if SubordinateResources is flase then check originofresource is same as originofcondition
		for _, origin := range originResources {
			if subscription.SubordinateResources == true {
				if strings.Contains(originCondition, origin) {
					log.Println("Filtered Event =: ", message)
					return true
				}
			} else {
				if origin == originCondition {
					log.Println("Filtered Event: ", message)
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
func formatEvent(originResource, eventstring, hostIP string) (string, string) {
	uuid, _ := getUUID(originResource)
	var eventRequestString = eventstring
	if !strings.Contains(hostIP, "Collection") {
		str := "/redfish/v1/Systems/" + uuid + ":"
		eventRequestString = strings.Replace(eventRequestString, "/redfish/v1/Systems/", str, -1)
		str = "/redfish/v1/systems/" + uuid + ":"
		eventRequestString = strings.Replace(eventRequestString, "/redfish/v1/systems/", str, -1)
		str = "/redfish/v1/Chassis/" + uuid + ":"
		eventRequestString = strings.Replace(eventRequestString, "/redfish/v1/Chassis/", str, -1)
		str = "/redfish/v1/Managers/" + uuid + ":"
		eventRequestString = strings.Replace(eventRequestString, "/redfish/v1/Managers/", str, -1)
	}
	return eventRequestString, uuid
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
	go p.checkCollectionSubscription(origin, "Redfish")
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
