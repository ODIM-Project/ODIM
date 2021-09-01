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
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/events"
	"github.com/ODIM-Project/ODIM/svc-events/evresponse"
)

//Events struct helps to register service
type Events struct {
	ContactClientRPC      func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	IsAuthorizedRPC       func(sessionToken string, privileges []string, oemPrivileges []string) response.RPC
	GetSessionUserNameRPC func(sessionToken string) (string, error)
	CreateTaskRPC         func(string) (string, error)
	UpdateTaskRPC         func(task common.TaskData) error
	CreateChildTaskRPC    func(sessionid, taskid string) (string, error)
}

func generateResponse(input interface{}) []byte {
	bytes, err := json.Marshal(input)
	if err != nil {
		log.Error("error in unmarshalling response object from util-libs" + err.Error())
	}
	return bytes
}

//GetEventService handles the RPC to get EventService details.
func (e *Events) GetEventService(ctx context.Context, req *eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error) {
	var resp eventsproto.EventSubResponse

	// Fill the response header first
	resp.Header = map[string]string{
		"Allow":             "GET",
		"Cache-Control":     "no-cache",
		"Connection":        "Keep-alive",
		"Link":              "</redfish/v1/SchemaStore/en/EventService.json>; rel=describedby",
		"Transfer-Encoding": "chunked",
		"X-Frame-Options":   "sameorigin",
		"Content-type":      "application/json; charset=utf-8",
		"OData-Version":     "4.0",
	}
	// Validate the token, if user has Login privelege then proceed.
	//Else send 401 Unautherised
	var oemprivileges []string
	privileges := []string{common.PrivilegeLogin}
	authResp := e.IsAuthorizedRPC(req.SessionToken, privileges, oemprivileges)
	if authResp.StatusCode != http.StatusOK {
		log.Error("error while trying to authenticate session: status code: " + string(authResp.StatusCode) +
			", status message: " + authResp.StatusMessage)
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
		DeliveryRetryAttempts:        evcommon.DeliveryRetryAttempts,
		DeliveryRetryIntervalSeconds: evcommon.DeliveryRetryIntervalSeconds,
		EventFormatTypes:             []string{"Event"},
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
	pc := events.PluginContact{
		ContactClient:      e.ContactClientRPC,
		Auth:               e.IsAuthorizedRPC,
		UpdateTask:         e.UpdateTaskRPC,
		CreateChildTask:    e.CreateChildTaskRPC,
		GetSessionUserName: e.GetSessionUserNameRPC,
	}
	// Athorize the request here
	authResp := e.IsAuthorizedRPC(req.SessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Error("error while trying to authenticate session: status code: " +
			string(authResp.StatusCode) + ", status message: " + authResp.StatusMessage)
		resp.Body = generateResponse(authResp.Body)
		resp.StatusCode = authResp.StatusCode
		return &resp, nil
	}
	sessionUserName, err := pc.GetSessionUserName(req.SessionToken)
	if err != nil {
		errorMessage := "error while trying to get the session username: " + err.Error()
		resp.Body = generateResponse(common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil))
		resp.StatusCode = http.StatusUnauthorized
		log.Error(errorMessage)
		return &resp, err
	}
	// Create the task and get the taskID
	// Contact Task Service using RPC and get the taskID
	taskURI, err := e.CreateTaskRPC(sessionUserName)
	if err != nil {
		// print err here as we are unbale to contact svc-task service
		errorMessage := "error while trying to create the task: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		resp.Body, _ = json.Marshal(common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil).Body)
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error   headers
		}
		log.Error(errorMessage)
		return &resp, fmt.Errorf(resp.StatusMessage)
	}
	strArray := strings.Split(taskURI, "/")
	if strings.HasSuffix(taskURI, "/") {
		taskID = strArray[len(strArray)-2]
	} else {
		taskID = strArray[len(strArray)-1]
	}
	//Spawn the thread to process the action asynchronously
	go pc.CreateEventSubscription(taskID, sessionUserName, req)
	// Return 202 accepted
	resp.StatusCode = http.StatusAccepted
	resp.Header = map[string]string{
		"Content-type": "application/json; charset=utf-8",
		"Location":     "/taskmon/" + taskID,
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
	pc := events.PluginContact{
		ContactClient:      e.ContactClientRPC,
		Auth:               e.IsAuthorizedRPC,
		GetSessionUserName: e.GetSessionUserNameRPC,
	}

	data := pc.SubmitTestEvent(req)
	resp.Body, err = json.Marshal(data.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying to marshal the response body for submit test event: " + err.Error()
		log.Error(resp.StatusMessage)
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
	pc := events.PluginContact{
		ContactClient: e.ContactClientRPC,
		Auth:          e.IsAuthorizedRPC,
	}

	data := pc.GetEventSubscriptionsCollection(req)
	resp.Body, err = json.Marshal(data.Body)
	if err != nil {
		errorMessage := "error while trying marshal the response body for get event subsciption : " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		resp.Body, _ = json.Marshal(common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil).Body)
		log.Error(resp.StatusMessage)
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
	pc := events.PluginContact{
		ContactClient: e.ContactClientRPC,
		Auth:          e.IsAuthorizedRPC,
	}

	data := pc.GetEventSubscriptionsDetails(req)
	resp.Body, err = json.Marshal(data.Body)
	if err != nil {
		errorMessage := "error while trying marshal the response body for get event subsciption : " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		resp.Body, _ = json.Marshal(common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil).Body)
		log.Error(resp.StatusMessage)
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
	pc := events.PluginContact{
		ContactClient: e.ContactClientRPC,
		Auth:          e.IsAuthorizedRPC,
	}
	var data response.RPC
	if req.UUID == "" {
		// Delete Event Subscription when admin requested
		data = pc.DeleteEventSubscriptionsDetails(req)
	} else {
		// Delete Event Subscription to Device when Server get Deleted
		data = pc.DeleteEventSubscriptions(req)
	}
	resp.Body, err = json.Marshal(data.Body)
	if err != nil {
		errorMessage := "error while trying marshal the response body for delete event subsciption : " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		resp.Body, _ = json.Marshal(common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil).Body)
		log.Error(resp.StatusMessage)
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
	pc := events.PluginContact{
		ContactClient: e.ContactClientRPC,
		Auth:          e.IsAuthorizedRPC,
	}
	pc.CreateDefaultEventSubscription(req.SystemID, req.EventTypes, req.MessageIDs, req.ResourceTypes, req.Protocol)
	return &resp, nil
}

//SubsribeEMB defines the operations which handles the RPC request response
// it subscribe to the given event message bus queues
func (e *Events) SubsribeEMB(ctx context.Context, req *eventsproto.SubscribeEMBRequest) (*eventsproto.SubscribeEMBResponse, error) {
	var resp eventsproto.SubscribeEMBResponse
	log.Info("Subscribing on emb for plugin " + req.PluginID)
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
