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
	"github.com/sirupsen/logrus"
)

var (
	// SendEventFunc function  pointer for calling the files
	SendEventFunc = sendEvent
	//ServiceDiscoveryFunc func pointer for calling the files
	ServiceDiscoveryFunc = services.ODIMService.Client
)

// PublishEventsToDestination This method sends the event/alert to subscriber's destination
// Takes:
//
//	data of type interface{}
//
// Returns:
//
//	bool: return false if any error occurred during execution, else returns true
func (e *ExternalInterfaces) PublishEventsToDestination(ctx context.Context, data interface{}) bool {
	eventUniqueID := uuid.NewV4().String()
	logging = logging.WithFields(logrus.Fields{"transactionid": eventUniqueID})
	if data == nil {
		logging.Info("invalid input params")
		return false
	}
	event := data.(common.Events)
	if event.EventType == "PluginStartUp" {
		logging.Info("received plugin started event from ", event.IP)
		go callPluginStartUp(ctx, event)
		return true
	}

	// Extract the Hostname/IP of the event source and Event from input parameter
	host, _, err := net.SplitHostPort(event.IP)
	if err != nil {
		host = event.IP
	}
	logging.Info("After splitting host address, IP is: ", host)

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
		logging.Error("failed to unmarshal the incoming event: ", requestData, " with the error: ", err.Error())
		return false
	}
	systemID, err := getSourceID(host)
	if err != nil {
		logging.Info("no origin resources found in device subscriptions")
		return false
	}
	message, deviceUUID = formatEvent(rawMessage, systemID, host)
	eventMap := make(map[string][]common.Event)

	for index, inEvent := range message.Events {
		subscriptions := getSubscriptions(inEvent.OriginOfCondition.Oid, systemID, host)
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
			if strings.HasPrefix(inEvent.OriginOfCondition.Oid, "/redfish/v1/Fabrics") {
				if strings.EqualFold(inEvent.EventType, "ResourceAdded") {
					e.addFabricRPCCall(ctx, rawMessage.Events[index].OriginOfCondition.Oid, host)
				}
				if strings.EqualFold(inEvent.EventType, "ResourceRemoved") {
					e.removeFabricRPCCall(ctx, rawMessage.Events[index].OriginOfCondition.Oid, host)
				}
			}
		}
	}

	for key, value := range eventMap {
		message.Events = value
		data, err := json.Marshal(message)
		if err != nil {
			logging.Error("unable to converts event into bytes: ", err.Error())
			continue
		}
		eventForwardingChanel <- evmodel.EventPost{Destination: key, EventID: eventUniqueID, Message: data}
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
		eventForwardingChanel <- evmodel.EventPost{Destination: sub.EventDestination.Destination, EventID: eventUniqueID, Message: []byte(requestData)}
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
		// if SubordinateResources is true then check if originOfresource is top level of originofcondition
		// if SubordinateResources is false then check originofresource is same as originofcondition

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
// add uuid:systemid/chassisid in place of systemid/chassisid
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

			// if subordinateResources is true then first check the child resourceType is present in db.
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
	logging.Info("Event not forwarded : No subscription for the incoming event's originofcondition")
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
	logging.Info("Event not forwarded : No subscription for the incoming event's ", message)
	return false
}

// postEvent will post the event to destination
func (e *ExternalInterfaces) postEvent(eventMessage evmodel.EventPost) {
	resp, err := SendEventFunc(eventMessage.Destination, eventMessage.Message)
	if err == nil {
		resp.Body.Close()
		logging.Info("Event is successfully forwarded 1 ")
		return
	}
	undeliveredEventID := eventMessage.Destination + ":" + eventMessage.EventID
	eventMessage.UndeliveredEventID = undeliveredEventID
	saveEventChanel <- eventMessage
	if reAttemptInQueue[eventMessage.Destination] <= 5 {
		e.reAttemptEvents(eventMessage)
	}
}

// sendEvent function is forward data to destination
func sendEvent(destination string, event []byte) (*http.Response, error) {
	httpConf := &config.HTTPConfig{
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
	}
	httpClient, err := httpConf.GetHTTPClientObj()
	if err != nil {
		l.Log.Error("failed to get http client object: ", err.Error())
		return &http.Response{}, err
	}
	req, err := http.NewRequest("POST", destination, bytes.NewBuffer(event))
	if err != nil {
		l.Log.Error("error while getting new http request: ", err.Error())
		return &http.Response{}, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	config.TLSConfMutex.RLock()
	defer config.TLSConfMutex.RUnlock()
	return httpClient.Do(req)
}
func (e *ExternalInterfaces) reAttemptEvents(eventMessage evmodel.EventPost) {
	reattemptLock.Lock()
	reAttemptInQueue[eventMessage.Destination] = reAttemptInQueue[eventMessage.Destination] + 1
	reattemptLock.Unlock()
	var resp *http.Response
	var err error
	count := config.Data.EventConf.DeliveryRetryAttempts

	defer func() {
		reattemptLock.Lock()
		reAttemptInQueue[eventMessage.Destination] = reAttemptInQueue[eventMessage.Destination] - 1
		reattemptLock.Unlock()
	}()
	for i := 0; i < count; i++ {
		logging.Debug("Retry event forwarding on destination: ")
		time.Sleep(time.Second * time.Duration(config.Data.EventConf.DeliveryRetryIntervalSeconds))
		// if undelivered event already published then ignore retrying
		eventString, err := e.GetUndeliveredEvents(eventMessage.UndeliveredEventID)
		if err != nil || len(eventString) < 1 {
			l.Log.Debug("Event is forwarded to destination")
			return
		}
		resp, err = SendEventFunc(eventMessage.Destination, eventMessage.Message)
		if err == nil {
			resp.Body.Close()
			logging.Info("Event is successfully forwarded after reattempt ")
			err = e.DeleteUndeliveredEvents(eventMessage.UndeliveredEventID)
			if err != nil {
				logging.Error("error while deleting undelivered events: ", err.Error())
			}
			return
		}
	}
	if err != nil {
		logging.Error("error while make https call to send the event: ", err.Error())
	}

}

// rediscoverSystemInventory will be triggered when ever the System Restart or Power On
// event is detected it will create a rpc for aggregation which will delete all system inventory //
// and rediscover all of them
func rediscoverSystemInventory(ctx context.Context, systemID, systemURL string) {
	systemURL = strings.TrimSuffix(systemURL, "/")

	conn, err := ServiceDiscoveryFunc(services.Aggregator)
	if err != nil {
		logging.Error("failed to get client connection object for aggregator service")
		return
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	_, err = aggregator.RediscoverSystemInventory(context.TODO(), &aggregatorproto.RediscoverSystemInventoryRequest{
		SystemID:  systemID,
		SystemURL: systemURL,
	})
	if err != nil {
		logging.Info("Error while rediscoverSystemInventory")
		return
	}
	logging.Info("rediscovery of system and chassis started.")

}

func (e *ExternalInterfaces) addFabricRPCCall(ctx context.Context, origin, address string) {
	if strings.Contains(origin, "Zones") || strings.Contains(origin, "Endpoints") || strings.Contains(origin, "AddressPools") {
		return
	}
	conn, err := ServiceDiscoveryFunc(services.Fabrics)
	if err != nil {
		logging.Error("Error while AddFabric ", err.Error())
		return
	}
	defer conn.Close()
	fab := fabricproto.NewFabricsClient(conn)
	ctxt := common.CreateNewRequestContext(ctx)
	ctxt = common.CreateMetadata(ctxt)
	_, err = fab.AddFabric(ctxt, &fabricproto.AddFabricRequest{
		OriginResource: origin,
		Address:        address,
	})
	if err != nil {
		logging.Error("Error while AddFabric ", err.Error())
		return
	}
	e.checkCollectionSubscription(ctx, origin, "Redfish")
	logging.Info("Fabric Added")
}
func (e *ExternalInterfaces) removeFabricRPCCall(ctx context.Context, origin, address string) {
	if strings.Contains(origin, "Zones") || strings.Contains(origin, "Endpoints") || strings.Contains(origin, "AddressPools") {
		return
	}
	conn, err := ServiceDiscoveryFunc(services.Fabrics)
	if err != nil {
		logging.Error("Error while Remove Fabric ", err.Error())
		return
	}
	defer conn.Close()
	fab := fabricproto.NewFabricsClient(conn)
	ctxt := common.CreateNewRequestContext(ctx)
	ctxt = common.CreateMetadata(ctxt)
	_, err = fab.RemoveFabric(ctxt, &fabricproto.AddFabricRequest{
		OriginResource: origin,
		Address:        address,
	})
	if err != nil {
		logging.Error("Error while RemoveFabric ", err.Error())
		return
	}
	logging.Info("Fabric Removed")
}

// updateSystemPowerState will be triggered when ever the System Powered Off event is received
// When event is detected a rpc is created for aggregation which will update the system inventory
func updateSystemPowerState(ctx context.Context, systemUUID, systemURI, state string) {

	systemURI = strings.TrimSuffix(systemURI, "/")

	index := strings.LastIndex(systemURI, "/")
	uri := systemURI[:index]
	id := systemURI[index+1:]

	if strings.ContainsAny(id, ":/-") {
		logging.Error("event contains invalid origin of condition - ", systemURI)
		return
	}
	if strings.Contains(state, "ServerPoweredOn") {
		state = "On"
	} else {
		state = "Off"
	}

	conn, err := ServiceDiscoveryFunc(services.Aggregator)
	if err != nil {
		logging.Error("failed to get client connection object for aggregator service")
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
		logging.Error("system power state update failed with ", err.Error())
		return
	}
	logging.Info("system power state update initiated")
}

func callPluginStartUp(ctx context.Context, event common.Events) {
	var message common.PluginStatusEvent
	if err := JSONUnmarshal([]byte(event.Request), &message); err != nil {
		logging.Error("failed to unmarshal the plugin startup event from "+event.IP+
			" with the error: ", err.Error())
		return
	}

	conn, err := ServiceDiscoveryFunc(services.Aggregator)
	if err != nil {
		logging.Error("failed to get client connection object for aggregator service")
		return
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)
	if _, err = aggregator.SendStartUpData(context.TODO(), &aggregatorproto.SendStartUpDataRequest{
		PluginAddr: event.IP,
		OriginURI:  message.OriginatorID,
	}); err != nil {
		logging.Error("failed to send plugin startup data to " + event.IP + ": " + err.Error())
		return
	}
	logging.Debug("successfully sent plugin startup data to " + event.IP)
}

func (e *ExternalInterfaces) checkUndeliveredEvents(destination string) {
	// first check any of the instance have already picked up for publishing
	// undelivered events for the destination
	flag, _ := e.GetUndeliveredEventsFlag(destination)
	if !flag {
		// if flag is false then set the flag true, so other instance shouldn't have to read the undelivered events and publish
		err := e.SetUndeliveredEventsFlag(destination)
		if err != nil {
			logging.Error("error while setting undelivered events flag: ", err.Error())
		}
		defer func() {
			err := e.DeleteUndeliveredEventsFlag(destination)
			if err != nil {
				logging.Error("error while deleting undelivered events flag: ", err.Error())
			}
		}()
		cursorCount := 0
		for {
			destData, tempCount, err := e.GetUndeliveredEventsKeyList(evmodel.UndeliveredEvents, destination, common.OnDisk, cursorCount)
			if err != nil {
				logging.Error("error while getting undelivered events list : ", err.Error())
				return
			}
			cursorCount = tempCount
			for _, dest := range destData {

				event, err := e.GetUndeliveredEvents(dest)
				if err != nil {
					logging.Error("error while getting undelivered events: ", err.Error())
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
					logging.Error("error while make https call to send the event: ", err.Error())
					time.Sleep(100 * time.Millisecond)
					break
				}
				logging.Debug("Event is successfully forwarded")
				err = e.DeleteUndeliveredEvents(dest)
				if err != nil {
					logging.Error("error while deleting undelivered events: ", err.Error())
				}
			}
			if cursorCount == 0 {
				break
			}
		}

	}
}
