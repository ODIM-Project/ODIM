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
	log "github.com/sirupsen/logrus"
	"net"
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

// addFabric will add the new fabric resource to db when an event is ResourceAdded and
// originofcondition has fabrics odataid.
func addFabric(requestData, host string) {
	var message common.MessageData
	if err := json.Unmarshal([]byte(requestData), &message); err != nil {
		log.Error("failed to unmarshal the incoming event: " + requestData + " with the error: " + err.Error())
		return
	}
	for _, inEvent := range message.Events {
		if inEvent.OriginOfCondition == nil || len(inEvent.OriginOfCondition.Oid) < 1 {
			log.Info("event not forwarded : Originofcondition is empty in incoming event")
			continue
		}
		if strings.EqualFold(inEvent.EventType, "ResourceAdded") &&
			strings.HasPrefix(inEvent.OriginOfCondition.Oid, "/redfish/v1/Fabrics") {
			addFabricRPCCall(inEvent.OriginOfCondition.Oid, host)
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
		log.Info("error: invalid input params")
		return false
	}
	// Extract the Hostname/IP of the event source and Event from input parameter
	event := data.(common.Events)
	host, _, err := net.SplitHostPort(event.IP)
	if err != nil {
		host = event.IP
	}
	log.Info("After splitting host address, IP is: ", host)

	var requestData = string(event.Request)
	//replacing the resposne with north bound translation URL
	for key, value := range config.Data.URLTranslation.NorthBoundURL {
		requestData = strings.Replace(requestData, key, value, -1)
	}

	var flag bool
	var uuid string
	var message common.MessageData

	addFabric(requestData, host)
	searchKey := evcommon.GetSearchKey(host, evmodel.DeviceSubscriptionIndex)
	deviceSubscription, err := evmodel.GetDeviceSubscriptions(searchKey)
	if err != nil {
		log.Error("Failed to get the event destinations: ", err.Error())
		return false
	}

	if len(deviceSubscription.OriginResources) < 1 {
		log.Info("no origin resources found in device subscriptions")
		return false
	}

	requestData, uuid = formatEvent(requestData, deviceSubscription.OriginResources[0], host)

	searchKey = evcommon.GetSearchKey(host, evmodel.SubscriptionIndex)
	subscriptions, err := evmodel.GetEvtSubscriptions(searchKey)
	if err != nil {
		return false
	}

	err = json.Unmarshal([]byte(requestData), &message)
	if err != nil {
		log.Error("failed to unmarshal the incoming event: ", requestData, " with the error: ", err.Error())
		return false
	}

	eventMap := make(map[string][]common.Event)
	for _, inEvent := range message.Events {
		if inEvent.OriginOfCondition == nil {
			log.Info("event not forwarded : Originofcondition is empty in incoming event with body: ", requestData)
			continue
		}

		if len(inEvent.OriginOfCondition.Oid) < 1 {
			log.Info("event not forwarded : Originofcondition is empty in incoming event with body: ", requestData)
			continue
		}

		var resTypePresent bool
		originofCond := strings.Split(strings.TrimSuffix(inEvent.OriginOfCondition.Oid, "/"), "/")
		resType := originofCond[len(originofCond)-2]
		for _, value := range common.ResourceTypes {
			if strings.Contains(resType, value) {
				resTypePresent = true
			}
		}

		if !resTypePresent {
			log.Info("event not forwared: resource type of originofcondition not supported in event with body: ", requestData)
			continue
		}

		for _, sub := range subscriptions {

			// filter and send events to destination if destination is not empty
			// in case of default event subscription destination will be empty
			if sub.Destination != "" {
				// check if hostip present in the hosts slice to make sure that it doesn't filter with the destination ip
				if isHostPresent(sub.Hosts, host) {
					if filterEventsToBeForwarded(sub, inEvent, deviceSubscription.OriginResources) {
						eventMap[sub.Destination] = append(eventMap[sub.Destination], inEvent)
						flag = true
					}
				} else {
					log.Info("event not forwarded : No subscription for the incoming event's originofcondition")
					flag = false
				}

			}
		}

		if strings.EqualFold("Alert", inEvent.EventType) {
			if strings.Contains(inEvent.MessageID, "ServerPostDiscoveryComplete") || strings.Contains(inEvent.MessageID, "ServerPostComplete") {
				go rediscoverSystemInventory(uuid, inEvent.OriginOfCondition.Oid)
				flag = true
			}
			if strings.Contains(inEvent.MessageID, "ServerPoweredOn") || strings.Contains(inEvent.MessageID, "ServerPoweredOff") {
				go updateSystemPowerState(uuid, inEvent.OriginOfCondition.Oid, inEvent.MessageID)
				flag = true
			}
		} else if strings.EqualFold("ResourceAdded", message.Events[0].EventType) || strings.EqualFold("ResourceRemoved", message.Events[0].EventType) {
			if strings.Contains(message.Events[0].OriginOfCondition.Oid, "Volumes") {
				s := strings.Split(message.Events[0].OriginOfCondition.Oid, "/")
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
			log.Error("unable to converts event into bytes: ", err.Error())
			continue
		}
		go postEvent(key, data)
	}
	return flag
}

func filterEventsToBeForwarded(subscription evmodel.Subscription, event common.Event, originResources []string) bool {
	eventTypes := subscription.EventTypes
	messageIds := subscription.MessageIds
	resourceTypes := subscription.ResourceTypes
	originCondition := strings.TrimSuffix(event.OriginOfCondition.Oid, "/")
	if (len(eventTypes) == 0 || isStringPresentInSlice(eventTypes, event.EventType, "event type")) &&
		(len(messageIds) == 0 || isStringPresentInSlice(messageIds, event.MessageID, "message id")) &&
		(len(resourceTypes) == 0 || isResourceTypeSubscribed(resourceTypes, event.OriginOfCondition.Oid, subscription.SubordinateResources)) {
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
	log.Info("Event not forwarded : No subscription for the incoming event's originofcondition")
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
	log.Info("Event not forwarded : No subscription for the incoming event's originofcondition")
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
	log.Info("Event not forwarded : No subscription for the incoming event's ", message)
	return false
}

// postEvent will post the event to destination
func postEvent(destination string, event []byte) {
	httpConf := &config.HTTPConfig{
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
	}
	httpClient, err := httpConf.GetHTTPClientObj()
	if err != nil {
		log.Error("failed to get http client object: ", err.Error())
		return
	}
	req, err := http.NewRequest("POST", destination, bytes.NewBuffer(event))
	if err != nil {
		log.Error("error while getting new http request: ", err.Error())
		return
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	var resp *http.Response
	count := evcommon.DeliveryRetryAttempts + 1
	for i := 0; i < count; i++ {
		config.TLSConfMutex.RLock()
		resp, err = httpClient.Do(req)
		config.TLSConfMutex.RUnlock()
		if err == nil {
			resp.Body.Close()
			log.Printf("event post response: %v", resp)
			return
		}
		log.Println("Retrying event posting")
		time.Sleep(time.Second * evcommon.DeliveryRetryIntervalSeconds)
	}
	log.Error("error while make https call to send the event: ", err.Error())
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
		log.Info("Error while rediscoverSystemInventroy")
		return
	}
	log.Info("info: rediscovery of system and chasis started.")
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
		log.Error("Error while AddFabric ", err.Error())
		return
	}
	p := PluginContact{
		ContactClient: pmbhandle.ContactPlugin,
	}
	p.checkCollectionSubscription(origin, "Redfish")
	log.Info("info: Fabric Added")
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
		log.Info("error: event contains invalid origin of condition - ", systemURI)
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
		log.Error("system power state update failed with ", err.Error())
		return
	}
	log.Info("info: system power state update initiated")
	return
}
