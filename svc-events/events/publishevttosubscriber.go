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
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	fabricproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	uuid "github.com/satori/go.uuid"
)

var (
	// SendEventFunc function  pointer for calling the files
	SendEventFunc = sendEvent
	//ServiceDiscoveryFunc func pointer for calling the files
	ServiceDiscoveryFunc = services.ODIMService.Client
)

// addFabric will add the new fabric resource to db when an event is ResourceAdded and
// originofcondition has fabrics odataid.
func (e *ExternalInterfaces) addFabric(ctx context.Context, message common.MessageData, host string) {
	for _, inEvent := range message.Events {
		if inEvent.OriginOfCondition == nil || len(inEvent.OriginOfCondition.Oid) < 1 {
			l.LogWithFields(ctx).Info("event not forwarded : Originofcondition is empty in incoming event")
			continue
		}
		if strings.EqualFold(inEvent.EventType, "ResourceAdded") &&
			strings.HasPrefix(inEvent.OriginOfCondition.Oid, "/redfish/v1/Fabrics") {
			e.addFabricRPCCall(ctx, inEvent.OriginOfCondition.Oid, host)
		}
		if strings.EqualFold(inEvent.EventType, "ResourceRemoved") &&
			strings.HasPrefix(inEvent.OriginOfCondition.Oid, "/redfish/v1/Fabrics") {
			e.removeFabricRPCCall(ctx, inEvent.OriginOfCondition.Oid, host)
		}
	}
}

// PublishEventsToDestination This method sends the event/alert to subscriber's destination
// Takes:
// 	data of type interface{}
//Returns:
//	bool: return false if any error occurred during execution, else returns true
func (e *ExternalInterfaces) PublishEventsToDestination(ctx context.Context, data interface{}) bool {
	subscribeCacheLock.Lock()
	defer subscribeCacheLock.Unlock()

	if data == nil {
		l.LogWithFields(ctx).Info("invalid input params")
		return false
	}
	event := data.(common.Events)
	if event.EventType == "PluginStartUp" {
		l.LogWithFields(ctx).Info("received plugin started event from ", event.IP)
		go callPluginStartUp(ctx, event)
		return true
	}

	// Extract the Hostname/IP of the event source and Event from input parameter
	host, _, err := net.SplitHostPort(event.IP)
	if err != nil {
		host = event.IP
	}
	host = strings.ToLower(host)
	l.LogWithFields(ctx).Info("After splitting host address, IP is: ", host)

	var requestData = string(event.Request)
	//replacing the response with north bound translation URL
	for key, value := range config.Data.URLTranslation.NorthBoundURL {
		requestData = strings.Replace(requestData, key, value, -1)
	}

	if event.EventType == "MetricReport" {
		return e.publishMetricReport(ctx, requestData)
	}

	var flag bool
	var deviceUUID string
	var message, rawMessage common.MessageData

	if err = json.Unmarshal([]byte(requestData), &rawMessage); err != nil {
		l.LogWithFields(ctx).Error("failed to unmarshal the incoming event: ", requestData, " with the error: ", err.Error())
		return false
	}

	e.addFabric(ctx, rawMessage, host)
	systemId, err := getSourceId(host)
	if err != nil {
		l.LogWithFields(ctx).Info("no origin resources found in device subscriptions")
		return false
	}
	message, deviceUUID = formatEvent(rawMessage, systemId, host)
	eventUniqueID := uuid.NewV4().String()
	eventMap := make(map[string][]common.Event)

	for index, inEvent := range message.Events {
		if inEvent.OriginOfCondition == nil || len(inEvent.OriginOfCondition.Oid) < 1 {
			l.LogWithFields(ctx).Info("event not forwarded as Originofcondition is empty in incoming event: ", requestData)
			continue
		}
		var resTypePresent bool
		originofCond := strings.Split(strings.TrimSuffix(inEvent.OriginOfCondition.Oid, "/"), "/")
		if len(originofCond) > 2 {
			resType := originofCond[len(originofCond)-2]
			for _, value := range common.ResourceTypes {
				if strings.Contains(resType, value) {
					resTypePresent = true
				}
			}
		} else {
			l.LogWithFields(ctx).Info("event not forwarded as originofcondition is empty incoming event: ", requestData)
			continue
		}

		if !resTypePresent {
			l.LogWithFields(ctx).Info("event not forwarded as resource type of originofcondition not supported in incoming event: ", requestData)
			continue
		}
		subscriptions := getSubscriptions(inEvent.OriginOfCondition.Oid, systemId, host)
		for _, sub := range subscriptions {
			if filterEventsToBeForwarded(ctx, sub, inEvent, sub.OriginResources) {
				eventMap[sub.Destination] = append(eventMap[sub.Destination], inEvent)
				flag = true
			}
		}
		if strings.EqualFold("Alert", inEvent.EventType) {
			if strings.Contains(inEvent.MessageID, "ServerPostDiscoveryComplete") || strings.Contains(inEvent.MessageID, "ServerPostComplete") {
				go rediscoverSystemInventory(ctx, deviceUUID, inEvent.OriginOfCondition.Oid)
				flag = true
			}
			if strings.Contains(inEvent.MessageID, "ServerPoweredOn") || strings.Contains(inEvent.MessageID, "ServerPoweredOff") {
				go updateSystemPowerState(ctx, deviceUUID, rawMessage.Events[index].OriginOfCondition.Oid, inEvent.MessageID)
				flag = true
			}
		} else if strings.EqualFold("ResourceAdded", message.Events[0].EventType) || strings.EqualFold("ResourceRemoved", message.Events[0].EventType) {
			if strings.Contains(message.Events[0].OriginOfCondition.Oid, "Volumes") {
				s := strings.Split(message.Events[0].OriginOfCondition.Oid, "/")
				storageURI := fmt.Sprintf("/%s/%s/%s/%s/%s/", s[1], s[2], s[3], s[4], s[5])
				go rediscoverSystemInventory(ctx, deviceUUID, storageURI)
				flag = true
			}
		}
	}

	for key, value := range eventMap {
		message.Events = value
		data, err := json.Marshal(message)
		if err != nil {
			l.LogWithFields(ctx).Error("unable to converts event into bytes: ", err.Error())
			continue
		}
		go e.postEvent(ctx, key, eventUniqueID, data)
	}
	return flag
}

func (e *ExternalInterfaces) publishMetricReport(ctx context.Context, requestData string) bool {
	eventUniqueID := uuid.NewV4().String()
	subscriptions, err := e.GetEvtSubscriptions("MetricReport")
	if err != nil {
		return false
	}
	for _, sub := range subscriptions {
		go e.postEvent(ctx, sub.EventDestination.Destination, eventUniqueID, []byte(requestData))
	}
	return true
}

func filterEventsToBeForwarded(ctx context.Context, subscription dmtf.EventDestination, event common.Event, originResources []model.Link) bool {
	eventTypes := subscription.EventTypes
	messageIds := subscription.MessageIds
	resourceTypes := subscription.ResourceTypes
	originCondition := strings.TrimSuffix(event.OriginOfCondition.Oid, "/")
	if (len(eventTypes) == 0 || isStringPresentInSlice(ctx, eventTypes, event.EventType, "event type")) &&
		(len(messageIds) == 0 || isStringPresentInSlice(ctx, messageIds, event.MessageID, "message id")) &&
		(len(resourceTypes) == 0 || isResourceTypeSubscribed(ctx, resourceTypes, event.OriginOfCondition.Oid, subscription.SubordinateResources)) {
		// if SubordinateResources is true then check if originofresource is top level of originofcondition
		// if SubordinateResources is flase then check originofresource is same as originofcondition

		if len(subscription.OriginResources) == 0 {
			return true
		}
		for _, origin := range originResources {
			if subscription.SubordinateResources {
				if strings.Contains(originCondition, origin.Oid) {
					return true
				}
			} else {
				if origin.Oid == originCondition {
					return true
				}
			}
		}
	}
	return false
}

// formatEvent will format the event string according to the odimra
// add uuid:systemid/chassisid inplace of systemid/chassisid
func formatEvent(event common.MessageData, originResource, hostIP string) (common.MessageData, string) {
	deviceUUID, _ := getUUID(originResource)
	if !strings.Contains(hostIP, "Collection") {
		for _, event := range event.Events {
			if event.OriginOfCondition == nil || len(event.OriginOfCondition.Oid) < 1 {
				continue
			}
			str := "/redfish/v1/Systems/" + deviceUUID + "."
			event.OriginOfCondition.Oid = strings.Replace(event.OriginOfCondition.Oid, "/redfish/v1/Systems/", str, -1)
			str = "/redfish/v1/systems/" + deviceUUID + "."
			event.OriginOfCondition.Oid = strings.Replace(event.OriginOfCondition.Oid, "/redfish/v1/systems/", str, -1)
			str = "/redfish/v1/Chassis/" + deviceUUID + "."
			event.OriginOfCondition.Oid = strings.Replace(event.OriginOfCondition.Oid, "/redfish/v1/Chassis/", str, -1)
			str = "/redfish/v1/Managers/" + deviceUUID + "."
			event.OriginOfCondition.Oid = strings.Replace(event.OriginOfCondition.Oid, "/redfish/v1/Managers/", str, -1)
		}

	}
	return event, deviceUUID
}

func isResourceTypeSubscribed(ctx context.Context, resourceTypes []string, originOfCondition string, subordinateResources bool) bool {
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
	l.LogWithFields(ctx).Info("Event not forwarded : No subscription for the incoming event's originofcondition")
	return false
}

func isStringPresentInSlice(ctx context.Context, slice []string, str, message string) bool {
	//If the incoming event fields contains empty values return true
	if str == "" {
		return true
	}
	for _, value := range slice {
		if value == str {
			return true
		}
	}
	l.LogWithFields(ctx).Info("Event not forwarded : No subscription for the incoming event's ", message)
	return false
}

// postEvent will post the event to destination
func (e *ExternalInterfaces) postEvent(ctx context.Context, destination, eventUniqueID string, event []byte) {
	resp, err := SendEventFunc(destination, event)
	if err == nil {
		resp.Body.Close()
		l.LogWithFields(ctx).Info("Event is successfully forwarded")
		// check any undelivered events are present in db for the destination and publish those
		go e.checkUndeliveredEvents(ctx, destination)
		return
	}
	undeliveredEventID := destination + ":" + eventUniqueID
	serr := e.SaveUndeliveredEvents(undeliveredEventID, event)
	if serr != nil {
		l.LogWithFields(ctx).Error("error while saving undelivered event: ", serr.Error())
	}
	go e.reAttemptEvents(ctx, destination, undeliveredEventID, event)

}

func sendEvent(destination string, event []byte) (*http.Response, error) {
	httpConf := &config.HTTPConfig{
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
	}
	httpClient, err := httpConf.GetHTTPClientObj()
	if err != nil {
		return &http.Response{}, err
	}
	req, err := http.NewRequest("POST", destination, bytes.NewBuffer(event))
	if err != nil {
		return &http.Response{}, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	config.TLSConfMutex.RLock()
	defer config.TLSConfMutex.RUnlock()
	return httpClient.Do(req)
}

func (e *ExternalInterfaces) reAttemptEvents(ctx context.Context, destination, undeliveredEventID string, event []byte) {
	var resp *http.Response
	var err error
	count := config.Data.EventConf.DeliveryRetryAttempts
	for i := 0; i < count; i++ {
		l.LogWithFields(ctx).Info("Retry event forwarding on destination: ", destination)
		time.Sleep(time.Second * time.Duration(config.Data.EventConf.DeliveryRetryIntervalSeconds))
		// if undelivered event already published then ignore retrying
		eventString, err := e.GetUndeliveredEvents(undeliveredEventID)
		if err != nil || len(eventString) < 1 {
			l.LogWithFields(ctx).Info("Event is forwarded to destination")
			return
		}
		resp, err = SendEventFunc(destination, event)
		if err == nil {
			resp.Body.Close()
			l.LogWithFields(ctx).Info("Event is successfully forwarded")
			// if event is delivered then delete the same which is saved in 1st attempt
			err = e.DeleteUndeliveredEvents(undeliveredEventID)
			if err != nil {
				l.LogWithFields(ctx).Error("error while deleting undelivered events: ", err.Error())
			}
			// check any undelivered events are present in db for the destination and publish those
			go e.checkUndeliveredEvents(ctx, destination)
			return
		}

	}
	if err != nil {
		l.LogWithFields(ctx).Error("error while make https call to send the event: ", err.Error())
	}
}

// rediscoverSystemInventory will be triggered when ever the System Restart or Power On
// event is detected it will create a rpc for aggregation which will delete all system inventory //
// and rediscover all of them
func rediscoverSystemInventory(ctx context.Context, systemID, systemURL string) {
	systemURL = strings.TrimSuffix(systemURL, "/")

	conn, err := ServiceDiscoveryFunc(services.Aggregator)
	if err != nil {
		l.LogWithFields(ctx).Error("failed to get client connection object for aggregator service")
		return
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	_, err = aggregator.RediscoverSystemInventory(context.TODO(), &aggregatorproto.RediscoverSystemInventoryRequest{
		SystemID:  systemID,
		SystemURL: systemURL,
	})
	if err != nil {
		l.LogWithFields(ctx).Info("Error while rediscoverSystemInventory")
		return
	}
	l.LogWithFields(ctx).Info("rediscovery of system and chassis started.")

}

func (e *ExternalInterfaces) addFabricRPCCall(ctx context.Context, origin, address string) {
	if strings.Contains(origin, "Zones") || strings.Contains(origin, "Endpoints") || strings.Contains(origin, "AddressPools") {
		return
	}
	conn, err := ServiceDiscoveryFunc(services.Fabrics)
	if err != nil {
		l.LogWithFields(ctx).Error("Error while AddFabric ", err.Error())
		return
	}
	defer conn.Close()
	fab := fabricproto.NewFabricsClient(conn)
	_, err = fab.AddFabric(context.TODO(), &fabricproto.AddFabricRequest{
		OriginResource: origin,
		Address:        address,
	})
	if err != nil {
		l.LogWithFields(ctx).Error("Error while AddFabric ", err.Error())
		return
	}
	e.checkCollectionSubscription(ctx, origin, "Redfish")
	l.LogWithFields(ctx).Info("Fabric Added")
}
func (e *ExternalInterfaces) removeFabricRPCCall(ctx context.Context, origin, address string) {
	if strings.Contains(origin, "Zones") || strings.Contains(origin, "Endpoints") || strings.Contains(origin, "AddressPools") {
		return
	}
	conn, err := ServiceDiscoveryFunc(services.Fabrics)
	if err != nil {
		l.LogWithFields(ctx).Error("Error while Remove Fabric ", err.Error())
		return
	}
	defer conn.Close()
	fab := fabricproto.NewFabricsClient(conn)
	_, err = fab.RemoveFabric(context.TODO(), &fabricproto.AddFabricRequest{
		OriginResource: origin,
		Address:        address,
	})
	if err != nil {
		l.LogWithFields(ctx).Error("Error while RemoveFabric ", err.Error())
		return
	}
	l.LogWithFields(ctx).Info("Fabric Removed")
}

// updateSystemPowerState will be triggered when ever the System Powered Off event is received
// When event is detected a rpc is created for aggregation which will update the system inventory
func updateSystemPowerState(ctx context.Context, systemUUID, systemURI, state string) {

	systemURI = strings.TrimSuffix(systemURI, "/")

	index := strings.LastIndex(systemURI, "/")
	uri := systemURI[:index]
	id := systemURI[index+1:]

	if strings.ContainsAny(id, ":/-") {
		l.LogWithFields(ctx).Error("event contains invalid origin of condition - ", systemURI)
		return
	}
	if strings.Contains(state, "ServerPoweredOn") {
		state = "On"
	} else {
		state = "Off"
	}

	conn, err := ServiceDiscoveryFunc(services.Aggregator)
	if err != nil {
		l.LogWithFields(ctx).Error("failed to get client connection object for aggregator service")
		return
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	_, err = aggregator.UpdateSystemState(context.TODO(), &aggregatorproto.UpdateSystemStateRequest{
		SystemUUID: systemUUID,
		SystemID:   id,
		SystemURI:  uri,
		UpdateKey:  "PowerState",
		UpdateVal:  state,
	})
	if err != nil {
		l.LogWithFields(ctx).Error("system power state update failed with ", err.Error())
		return
	}
	l.LogWithFields(ctx).Info("system power state update initiated")
}

func callPluginStartUp(ctx context.Context, event common.Events) {
	var message common.PluginStatusEvent
	if err := JSONUnmarshal([]byte(event.Request), &message); err != nil {
		l.LogWithFields(ctx).Error("failed to unmarshal the plugin startup event from "+event.IP+
			" with the error: ", err.Error())
		return
	}

	conn, err := ServiceDiscoveryFunc(services.Aggregator)
	if err != nil {
		l.LogWithFields(ctx).Error("failed to get client connection object for aggregator service")
		return
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)
	if _, err = aggregator.SendStartUpData(context.TODO(), &aggregatorproto.SendStartUpDataRequest{
		PluginAddr: event.IP,
		OriginURI:  message.OriginatorID,
	}); err != nil {
		l.LogWithFields(ctx).Error("failed to send plugin startup data to " + event.IP + ": " + err.Error())
		return
	}
	l.LogWithFields(ctx).Info("successfully sent plugin startup data to " + event.IP)
}

func (e *ExternalInterfaces) checkUndeliveredEvents(ctx context.Context, destination string) {
	// first check any of the instance have already picked up for publishing
	// undelivered events for the destination
	flag, _ := e.GetUndeliveredEventsFlag(destination)
	if !flag {
		// if flag is false then set the flag true, so other instance shouldnt have to read the undelivered events and publish
		err := e.SetUndeliveredEventsFlag(destination)
		if err != nil {
			l.LogWithFields(ctx).Error("error while setting undelivered events flag: ", err.Error())
		}
		destData, _ := e.GetAllMatchingDetails(evmodel.UndeliveredEvents, destination, common.OnDisk)
		for _, dest := range destData {
			event, err := e.GetUndeliveredEvents(dest)
			if err != nil {
				l.LogWithFields(ctx).Error("error while getting undelivered events: ", err.Error())
				continue
			}
			event = strings.Replace(event, "\\", "", -1)
			event = strings.TrimPrefix(event, "\"")
			event = strings.TrimSuffix(event, "\"")
			resp, err := SendEventFunc(destination, []byte(event))
			if resp != nil {
				defer resp.Body.Close()
			}
			if err != nil {
				l.LogWithFields(ctx).Error("error while make https call to send the event: ", err.Error())
				continue
			}
			l.LogWithFields(ctx).Info("Event is successfully forwarded")
			err = e.DeleteUndeliveredEvents(dest)
			if err != nil {
				l.LogWithFields(ctx).Error("error while deleting undelivered events: ", err.Error())
			}
		}
		// handle logic if inter connection fails
		derr := e.DeleteUndeliveredEventsFlag(destination)
		if derr != nil {
			l.LogWithFields(ctx).Error("error while deleting undelivered events flag: ", derr.Error())
		}
	}
}
