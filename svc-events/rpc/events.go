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

// Package rpc have the functionality of rpc caller functions
package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/events"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/ODIM-Project/ODIM/svc-events/evresponse"
)

//Events struct helps to register service
type Events struct {
	Connector *events.ExternalInterfaces
}

var (
	//JSONMarshal ...
	JSONMarshal = json.Marshal
)

// GetPluginContactInitializer intializes all the required connection functions for the events execution
func GetPluginContactInitializer() *Events {
	connector := &events.ExternalInterfaces{
		External: events.External{
			ContactClient:   pmbhandle.ContactPlugin,
			Auth:            services.IsAuthorized,
			CreateTask:      services.CreateTask,
			UpdateTask:      events.UpdateTaskData,
			CreateChildTask: services.CreateChildTask,
		},
		DB: events.DB{
			GetSessionUserName:               services.GetSessionUserName,
			GetEvtSubscriptions:              evmodel.GetEvtSubscriptions,
			SaveEventSubscription:            evmodel.SaveEventSubscription,
			GetPluginData:                    evmodel.GetPluginData,
			GetDeviceSubscriptions:           evmodel.GetDeviceSubscriptions,
			GetTarget:                        evmodel.GetTarget,
			GetAllKeysFromTable:              evmodel.GetAllKeysFromTable,
			GetAllFabrics:                    evmodel.GetAllFabrics,
			GetAllMatchingDetails:            evmodel.GetAllMatchingDetails,
			UpdateDeviceSubscriptionLocation: evmodel.UpdateDeviceSubscriptionLocation,
			GetFabricData:                    evmodel.GetFabricData,
			DeleteEvtSubscription:            evmodel.DeleteEvtSubscription,
			UpdateEventSubscription:          evmodel.UpdateEventSubscription,
			DeleteDeviceSubscription:         evmodel.DeleteDeviceSubscription,
			SaveUndeliveredEvents:            evmodel.SaveUndeliveredEvents,
			SaveDeviceSubscription:           evmodel.SaveDeviceSubscription,
			GetUndeliveredEvents:             evmodel.GetUndeliveredEvents,
			GetUndeliveredEventsFlag:         evmodel.GetUndeliveredEventsFlag,
			SetUndeliveredEventsFlag:         evmodel.SetUndeliveredEventsFlag,
			DeleteUndeliveredEventsFlag:      evmodel.DeleteUndeliveredEventsFlag,
			DeleteUndeliveredEvents:          evmodel.DeleteUndeliveredEvents,
			GetAggregateData:                 evmodel.GetAggregateData,
			SaveAggregateSubscription:        evmodel.SaveAggregateSubscription,
			GetAggregateHosts:                evmodel.GetAggregateHosts,
			UpdateAggregateHosts:             evmodel.UpdateAggregateHosts,
			GetAggregateList:                 evmodel.GetAggregateList,
		},
	}
	return &Events{
		Connector: connector,
	}
}
func generateResponse(input interface{}) []byte {
	bytes, err := json.Marshal(input)
	if err != nil {
		l.Log.Error("error in unmarshalling response object from util-libs" + err.Error())
	}
	return bytes
}

//GetEventService handles the RPC to get EventService details.
func (e *Events) GetEventService(ctx context.Context, req *eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error) {
	var resp eventsproto.EventSubResponse

	// Fill the response header first
	resp.Header = map[string]string{
		"Link": "</redfish/v1/SchemaStore/en/EventService.json>; rel=describedby",
	}
	// Validate the token, if user has Login privelege then proceed.
	//Else send 401 Unautherised
	var oemprivileges []string
	privileges := []string{common.PrivilegeLogin}
	authResp := e.Connector.Auth(req.SessionToken, privileges, oemprivileges)
	if authResp.StatusCode != http.StatusOK {
		resp.Body = generateResponse(authResp.Body)
		resp.StatusMessage = authResp.StatusMessage
		resp.StatusCode = authResp.StatusCode
		return &resp, nil
	}
	// Check whether the Event Service is enbaled in configuration file.
	//If so set ServiceEnabled to true.
	isServiceEnabled := false
	serviceState := "Disabled"
	for _, service := range config.Data.EnabledServices {
		if service == "EventService" {
			isServiceEnabled = true
			serviceState = "Enabled"
			break
		}
	}
	var resourceTypes []string
	for resType := range common.ResourceTypes {
		resourceTypes = append(resourceTypes, resType)
	}
	// Construct the response below
	eventServiceResponse := evresponse.EventServiceResponse{
		OdataType:    common.EventServiceType,
		ID:           "EventService",
		Name:         "EventService",
		Description:  "EventService",
		OdataContext: "/redfish/v1/$metadata#EventService.EventService",
		OdataID:      "/redfish/v1/EventService",
		Actions: evresponse.Actions{
			SubmitTestEvent: evresponse.Action{
				Target: "/redfish/v1/EventService/Actions/EventService.SubmitTestEvent",
				AllowableValues: []string{
					"StatusChange",
					"ResourceUpdated",
					"ResourceAdded",
					"ResourceRemoved",
					"Alert"},
			},
			Oem: evresponse.Oem{},
		},
		DeliveryRetryAttempts:        config.Data.EventConf.DeliveryRetryAttempts,
		DeliveryRetryIntervalSeconds: config.Data.EventConf.DeliveryRetryIntervalSeconds,
		EventFormatTypes:             []string{"Event", "MetricReport"},
		EventTypesForSubscription: []string{
			"StatusChange",
			"ResourceUpdated",
			"ResourceAdded",
			"ResourceRemoved",
			"Alert"},
		RegistryPrefixes: []string{},
		ResourceTypes:    resourceTypes,
		//		ServerSentEventURI:"/redfish/v1/EventService/SSE",
		ServiceEnabled: isServiceEnabled,
		/*
			SSEFilterPropertiesSupported: &evresponse.SSEFilterPropertiesSupported{
				EventFormatType:        true,
				EventType:              true,
				MessageID:              true,
				MetricReportDefinition: false,
				OriginResource:         true,
				RegistryPrefix:         false,
				ResourceType:           true,
				SubordinateResources:   true,
			},
		*/

		Status: evresponse.Status{
			Health:       "OK",
			HealthRollup: "OK",
			State:        serviceState,
		},
		SubordinateResourcesSupported: true,
		Subscriptions: evresponse.Subscriptions{
			OdataID: "/redfish/v1/EventService/Subscriptions",
		},
		Oem: evresponse.Oem{},
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = "Success"
	resp.Body = generateResponse(eventServiceResponse)

	return &resp, nil
}

//CreateEventSubscription defines the operations which handles the RPC request response
// for the Create event subscription RPC call to events micro service.
// The functionality is to create the subscrription with Resource provided in origin resources.
func (e *Events) CreateEventSubscription(ctx context.Context, req *eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error) {
	var resp eventsproto.EventSubResponse
	var err error
	var taskID string
	// Athorize the request here
	authResp := e.Connector.Auth(req.SessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		resp.Body = generateResponse(authResp.Body)
		resp.StatusCode = authResp.StatusCode
		return &resp, nil
	}
	sessionUserName, err := e.Connector.GetSessionUserName(req.SessionToken)
	if err != nil {
		errorMessage := "error while trying to get the session username: " + err.Error()
		resp.Body = generateResponse(common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil))
		resp.StatusCode = http.StatusUnauthorized
		l.Log.Error(errorMessage)
		return &resp, err
	}
	// Create the task and get the taskID
	// Contact Task Service using RPC and get the taskID
	taskURI, err := e.Connector.CreateTask(sessionUserName)
	if err != nil {
		// print err here as we are unbale to contact svc-task service
		errorMessage := "error while trying to create the task: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		resp.Body, _ = json.Marshal(common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil).Body)
		l.Log.Error(errorMessage)
		return &resp, fmt.Errorf(resp.StatusMessage)
	}
	strArray := strings.Split(taskURI, "/")
	if strings.HasSuffix(taskURI, "/") {
		taskID = strArray[len(strArray)-2]
	} else {
		taskID = strArray[len(strArray)-1]
	}
	//Spawn the thread to process the action asynchronously
	go e.Connector.CreateEventSubscription(taskID, sessionUserName, req)
	// Return 202 accepted
	resp.StatusCode = http.StatusAccepted
	resp.Header = map[string]string{
		"Location": "/taskmon/" + taskID,
	}
	resp.StatusMessage = response.TaskStarted
	generateTaskRespone(taskID, taskURI, &resp)
	return &resp, nil
}

//SubmitTestEvent defines the operations which handles the RPC request response
// for the SubmitTestEvent RPC call to events micro service.
// The functionality is to submit a test event.
func (e *Events) SubmitTestEvent(ctx context.Context, req *eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error) {
	var resp eventsproto.EventSubResponse
	var err error
	data := e.Connector.SubmitTestEvent(req)
	resp.Body, err = JSONMarshal(data.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying to marshal the response body for submit test event: " + err.Error()
		l.Log.Error(resp.StatusMessage)
		return &resp, fmt.Errorf(resp.StatusMessage)
	}
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header

	return &resp, nil
}

//GetEventSubscriptionsCollection defines the operations which handles the RPC request response
// for the get event subscriptions collection RPC call to events micro service.
// The functionality is to get the collection of subscrription details.
func (e *Events) GetEventSubscriptionsCollection(ctx context.Context, req *eventsproto.EventRequest) (*eventsproto.EventSubResponse, error) {
	var resp eventsproto.EventSubResponse
	var err error
	data := e.Connector.GetEventSubscriptionsCollection(req)
	resp.Body, err = JSONMarshal(data.Body)
	if err != nil {
		errorMessage := "error while trying marshal the response body for get event subsciption : " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		resp.Body, _ = json.Marshal(common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil).Body)
		l.Log.Error(resp.StatusMessage)
		return &resp, nil
	}
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header

	return &resp, nil
}

//GetEventSubscription defines the operations which handles the RPC request response
// for the get event subscription RPC call to events micro service.
// The functionality is to get the subscrription details.
func (e *Events) GetEventSubscription(ctx context.Context, req *eventsproto.EventRequest) (*eventsproto.EventSubResponse, error) {
	var resp eventsproto.EventSubResponse
	var err error
	data := e.Connector.GetEventSubscriptionsDetails(req)
	resp.Body, err = JSONMarshal(data.Body)
	if err != nil {
		errorMessage := "error while trying marshal the response body for get event subsciption : " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		resp.Body, _ = json.Marshal(common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil).Body)
		l.Log.Error(resp.StatusMessage)
		return &resp, nil
	}
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header

	return &resp, nil
}

// DeleteEventSubscription defines the operations which handles the RPC request response
// for the delete event subscription RPC call to events micro service.
// The functionality is to delete the subscrription details.
func (e *Events) DeleteEventSubscription(ctx context.Context, req *eventsproto.EventRequest) (*eventsproto.EventSubResponse, error) {
	var resp eventsproto.EventSubResponse
	var err error
	var data response.RPC
	if req.UUID == "" {
		// Delete Event Subscription when admin requested
		data = e.Connector.DeleteEventSubscriptionsDetails(req)
	} else {
		// Delete Event Subscription to Device when Server get Deleted
		data = e.Connector.DeleteEventSubscriptions(req)
	}
	resp.Body, err = JSONMarshal(data.Body)
	if err != nil {
		errorMessage := "error while trying marshal the response body for delete event subsciption : " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		resp.Body, _ = json.Marshal(common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil).Body)
		l.Log.Error(resp.StatusMessage)
		return &resp, nil
	}
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header
	return &resp, nil
}

//CreateDefaultEventSubscription defines the operations which handles the RPC request response
// after computer system restarts ,This will  triggered from   aggregation service whenever a computer system is added
func (e *Events) CreateDefaultEventSubscription(ctx context.Context, req *eventsproto.DefaultEventSubRequest) (*eventsproto.DefaultEventSubResponse, error) {
	var resp eventsproto.DefaultEventSubResponse
	e.Connector.CreateDefaultEventSubscription(req.SystemID, req.EventTypes, req.MessageIDs, req.ResourceTypes, req.Protocol)
	return &resp, nil
}

//SubsribeEMB defines the operations which handles the RPC request response
// it subscribe to the given event message bus queues
func (e *Events) SubsribeEMB(ctx context.Context, req *eventsproto.SubscribeEMBRequest) (*eventsproto.SubscribeEMBResponse, error) {
	var resp eventsproto.SubscribeEMBResponse
	l.Log.Info("Subscribing on emb for plugin " + req.PluginID)
	for i := 0; i < len(req.EMBQueueName); i++ {
		evcommon.EMBTopics.ConsumeTopic(req.EMBQueueName[i])
	}
	resp.Status = true
	return &resp, nil
}

func generateTaskRespone(taskID, taskURI string, resp *eventsproto.EventSubResponse) {
	commonResponse := response.Response{
		OdataType:    common.TaskType,
		ID:           taskID,
		Name:         "Task " + taskID,
		OdataContext: "/redfish/v1/$metadata#Task.Task",
		OdataID:      taskURI,
	}
	commonResponse.MessageArgs = []string{taskID}
	commonResponse.CreateGenericResponse(resp.StatusMessage)
	resp.Body = generateResponse(commonResponse)
}

//RemoveEventSubscriptionsRPC defines the operations which handles the RPC request response
// it subscribe to the given event message bus queues
func (e *Events) RemoveEventSubscriptionsRPC(ctx context.Context, req *eventsproto.EventUpdateRequest) (*eventsproto.SubscribeEMBResponse, error) {
	var resp eventsproto.SubscribeEMBResponse
	e.Connector.UpdateEventSubscriptions(req, true)
	resp.Status = true
	return &resp, nil
}

//UpdateEventSubscriptionsRPC defines the operations which handles the RPC request response
// it subscribe to the given event message bus queues
func (e *Events) UpdateEventSubscriptionsRPC(ctx context.Context, req *eventsproto.EventUpdateRequest) (*eventsproto.SubscribeEMBResponse, error) {
	var resp eventsproto.SubscribeEMBResponse
	resp.Status = true
	e.Connector.UpdateEventSubscriptions(req, false)
	return &resp, nil
}

//IsAggregateHaveSubscription defines the operations which handles the RPC request response
func (e *Events) IsAggregateHaveSubscription(ctx context.Context, req *eventsproto.EventUpdateRequest) (*eventsproto.SubscribeEMBResponse, error) {
	var resp eventsproto.SubscribeEMBResponse
	isAvailable := e.Connector.IsAggregateHaveSubscription(req)
	resp.Status = isAvailable
	return &resp, nil
}

//DeleteAggregateSubscriptionsRPC defines the operations which handles the RPC request response
// it remove subscription details
func (e *Events) DeleteAggregateSubscriptionsRPC(ctx context.Context, req *eventsproto.EventUpdateRequest) (*eventsproto.SubscribeEMBResponse, error) {
	var resp eventsproto.SubscribeEMBResponse
	e.Connector.DeleteAggregateSubscriptions(req, true)
	resp.Status = true
	return &resp, nil
}
