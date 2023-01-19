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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	errResponse "github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/ODIM-Project/ODIM/svc-events/evresponse"
	"github.com/google/uuid"
)

// CreateEventSubscription is a API to create event subscription
func (e *ExternalInterfaces) CreateEventSubscription(taskID string, sessionUserName string, req *eventsproto.EventSubRequest) errResponse.RPC {
	var (
		err             error
		resp            errResponse.RPC
		postRequest     evmodel.RequestBody
		percentComplete int32 = 100
		targetURI             = "/redfish/v1/EventService/Subscriptions"
	)

	if err = json.Unmarshal(req.PostBody, &postRequest); err != nil {
		// Update the task here with error response
		errorMessage := "Error while Unmarshaling the Request: " + err.Error()
		if strings.Contains(err.Error(), "evmodel.OdataIDLink") {
			errorMessage = "Error processing subscription request: @odata.id key(s) is missing in origin resources list"
		}
		l.Log.Error(errorMessage)

		resp = common.GeneralError(http.StatusBadRequest, errResponse.MalformedJSON, errorMessage, []interface{}{}, nil)
		// Fill task and update
		e.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
		return resp
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.PostBody, postRequest)
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		l.Log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errMsg, nil, nil)
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		l.Log.Error(errorMessage)
		resp := common.GeneralError(http.StatusBadRequest, errResponse.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, nil)
		e.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
		return resp
	}

	//check mandatory fields
	statuscode, statusMessage, messageArgs, err := validateFields(&postRequest)
	if err != nil {
		// Update the task here with error response
		errorMessage := "error: request payload validation failed: " + err.Error()
		l.Log.Error(errorMessage)

		resp = common.GeneralError(statuscode, statusMessage, errorMessage, messageArgs, nil)
		// Fill task and update
		e.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
		return resp
	}

	//validate destination URI in the request
	if !common.URIValidator(postRequest.Destination) {
		errorMessage := "error: request body contains invalid value for Destination field, " + postRequest.Destination
		l.Log.Error(errorMessage)

		resp = common.GeneralError(http.StatusBadRequest, errResponse.PropertyValueFormatError, errorMessage, []interface{}{postRequest.Destination, "Destination"}, nil)
		// Fill task and update
		e.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
		return resp
	}

	// check any of the subscription present for the destination from the request
	// if errored out or no subscriptions then add subscriptions else return an error
	subscriptionDetails, err := e.GetEvtSubscriptions("")
	if err != nil && !strings.Contains(err.Error(), "No data found for the key") {
		errorMessage := "Error while get subscription details: " + err.Error()
		evcommon.GenErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			[]interface{}{}, &resp)
		l.Log.Error(errorMessage)
		e.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
		return resp
	}
	for _, evtSubscription := range subscriptionDetails {
		if evtSubscription.Destination == postRequest.Destination {
			errorMessage := "Subscription already present for the requested destination"
			evcommon.GenErrorResponse(errorMessage, errResponse.ResourceInUse, http.StatusConflict,
				[]interface{}{}, &resp)
			l.Log.Error(errorMessage)
			e.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
			return resp
		}
	}

	// Get the target device  details from the origin resources
	// Loop through all origin list and form individual event subscription request,
	// Which will then forward to plugin to make subscrption with target device
	var wg, taskCollectionWG sync.WaitGroup
	var result = &evresponse.MutexLock{
		Response: make(map[string]evresponse.EventResponse),
		Hosts:    make(map[string]string),
		Lock:     &sync.Mutex{},
	}

	// remove odataid in the originresources
	originResources := removeOdataIDfromOriginResources(postRequest.OriginResources)
	originResourcesCount := len(originResources)

	// check and remove if duplicate OriginResources exist in the request
	removeDuplicatesFromSlice(&originResources, &originResourcesCount)

	// If origin resource is nil then subscribe to all collection
	isDefaultOriginResource := false
	if originResourcesCount == 0 {
		isDefaultOriginResource = true
		originResources = []string{
			"/redfish/v1/Systems",
			"/redfish/v1/Chassis",
			"/redfish/v1/Fabrics",
			"/redfish/v1/Managers",
			"/redfish/v1/TaskService/Tasks",
		}
		originResourcesCount = len(originResources)
	}
	var collectionList = make([]string, 0)
	subTaskChan := make(chan int32, originResourcesCount)
	taskCollectionWG.Add(1)
	bubbleUpStatusCode := int32(http.StatusCreated)
	go func() {
		// Collect the channels and update percentComplete in Task
		for i := 1; ; i++ {
			statusCode, chanActive := <-subTaskChan
			if !chanActive {
				defer taskCollectionWG.Done()
				break
			}
			if statusCode > bubbleUpStatusCode {
				bubbleUpStatusCode = statusCode
			}
			if i <= originResourcesCount {
				percentComplete = int32((i*100)/originResourcesCount - 1)
				if resp.StatusCode == 0 {
					resp.StatusCode = http.StatusAccepted
				}
				e.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Running, common.OK, percentComplete, http.MethodPost))
			}
		}
	}()

	for _, origin := range originResources {
		_, err := getUUID(origin)
		if err != nil {
			collection, collectionName, collectionFlag, aggregateResource, isAggregate, _ := e.checkCollection(origin)
			wg.Add(1)
			// for origin is collection
			go e.createEventSubscription(taskID, subTaskChan, sessionUserName, targetURI, postRequest, origin, result, &wg, collectionFlag, collectionName, aggregateResource, isAggregate)
			for i := 0; i < len(collection); i++ {
				wg.Add(1)
				// for subordinate origin
				go e.createEventSubscription("", subTaskChan, sessionUserName, targetURI, postRequest, collection[i], result, &wg, false, "", aggregateResource, isAggregate)
			}
			if !isAggregate {
				collectionList = append(collectionList, collection...)
			}
		} else {
			wg.Add(1)
			go e.createEventSubscription(taskID, subTaskChan, sessionUserName, targetURI, postRequest, origin, result, &wg, false, "", "", false)
		}
	}

	wg.Wait()
	// close channel once all sub-routines created have exited
	close(subTaskChan)
	// wait till all the subtasks are collected and routine is exited
	taskCollectionWG.Wait()

	var (
		locationHeader             string
		successfulSubscriptionList = make([]string, 0)
		successfulResponses        = make(map[string]evresponse.EventResponse)
	)

	result.Lock.Lock()
	originResourceProcessedCount := len(result.Response)
	var resourceID string
	i := 0
	for originResource, evtResponse := range result.Response {
		OriginResource := strings.SplitAfter(originResource, "/")
		originResourceID := OriginResource[len(OriginResource)-1]
		if i == 0 {
			resourceID = originResourceID
		}
		if originResourceID == resourceID && i > 0 {
			successfulSubscriptionList = append(successfulSubscriptionList, originResource)
		}
		i = i + 1
		if evtResponse.StatusCode == http.StatusCreated {
			successfulSubscriptionList = append(successfulSubscriptionList, originResource)
			successfulResponses[originResource] = evtResponse
		}
	}
	result.Response = successfulResponses

	result.Lock.Unlock()
	// remove the underlying resource uri's from successfulSubscriptionList
	for i := 0; i < len(collectionList); i++ {
		for j := 0; j < len(successfulSubscriptionList); j++ {
			if collectionList[i] == successfulSubscriptionList[j] {
				originResourceProcessedCount--
				successfulSubscriptionList = append(successfulSubscriptionList[:j], successfulSubscriptionList[j+1:]...)
				break
			}
		}
	}
	// if Subscription Name is empty then use default name
	if postRequest.Name == "" {
		postRequest.Name = evmodel.SubscriptionName
	}

	successOriginResourceCount := len(successfulSubscriptionList)
	if successOriginResourceCount > 0 {
		subscriptionID := uuid.New().String()
		var hosts []string
		resp, hosts = result.ReadResponse(subscriptionID)
		if isDefaultOriginResource {
			successfulSubscriptionList = []string{}
			hosts = []string{}
		}
		evtSubscription := evmodel.Subscription{
			UserName:             sessionUserName,
			SubscriptionID:       subscriptionID,
			Destination:          postRequest.Destination,
			Name:                 postRequest.Name,
			Context:              postRequest.Context,
			EventTypes:           postRequest.EventTypes,
			MessageIds:           postRequest.MessageIds,
			ResourceTypes:        postRequest.ResourceTypes,
			EventFormatType:      postRequest.EventFormatType,
			SubordinateResources: postRequest.SubordinateResources,
			Protocol:             postRequest.Protocol,
			SubscriptionType:     postRequest.SubscriptionType,
			OriginResources:      successfulSubscriptionList,
			Hosts:                hosts,
			DeliveryRetryPolicy:  postRequest.DeliveryRetryPolicy,
		}

		if err = e.SaveEventSubscription(evtSubscription); err != nil {
			// Update the task here with error response
			errorMessage := "error while trying to save event subscription data: " + err.Error()
			l.Log.Error(errorMessage)

			resp = common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, []interface{}{}, nil)
			// Fill task and update
			percentComplete = 100
			e.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
			return resp
		}
		locationHeader = resp.Header["Location"]
	}
	l.Log.Info("Process Count," + strconv.Itoa(originResourceProcessedCount) +
		" successOriginResourceCount " + strconv.Itoa(successOriginResourceCount))
	percentComplete = 100
	if originResourceProcessedCount == successOriginResourceCount {
		e.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Completed, common.OK, percentComplete, http.MethodPost))
	} else {
		args := errResponse.Args{
			Code:    errResponse.GeneralError,
			Message: "event subscription for one or more origin resource(s) failed, check sub tasks for more info.",
		}
		resp.Body = args.CreateGenericErrorResponse()
		resp.StatusCode = bubbleUpStatusCode
		if locationHeader != "" {
			resp.Header["Location"] = locationHeader
		}
		e.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
	}
	return resp
}

func (e *ExternalInterfaces) eventSubscription(postRequest evmodel.RequestBody, origin, collectionName string, collectionFlag bool) (string, evresponse.EventResponse) {
	var resp evresponse.EventResponse
	var err error
	var plugin *evmodel.Plugin
	var contactRequest evcommon.PluginContactRequest
	var target *evmodel.Target
	if !collectionFlag {
		if strings.Contains(origin, "Fabrics") {
			return e.createFabricSubscription(postRequest, origin, collectionName, collectionFlag)
		}
		target, resp, err = e.getTargetDetails(origin)
		if err != nil {
			return "", resp
		}
		var errs *errors.Error
		plugin, errs = e.GetPluginData(target.PluginID)
		if errs != nil {
			errorMessage := "error while getting plugin data: " + errs.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			l.Log.Error(errorMessage)
			return "", resp
		}

		contactRequest.Plugin = plugin
		if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
			token := e.getPluginToken(plugin)
			if token == "" {
				evcommon.GenEventErrorResponse("error: Unable to create session with plugin "+plugin.ID, errResponse.NoValidSession, http.StatusUnauthorized,
					&resp, []interface{}{})

				return "", resp
			}
			contactRequest.Token = token

		} else {
			contactRequest.LoginCredential = map[string]string{
				"UserName": plugin.Username,
				"Password": string(plugin.Password),
			}
		}
	}
	var httpHeadersSlice = make([]evmodel.HTTPHeaders, 0)
	httpHeadersSlice = append(httpHeadersSlice, evmodel.HTTPHeaders{ContentType: "application/json"})
	subscriptionPost := evmodel.EvtSubPost{
		Name:                 postRequest.Name,
		Destination:          postRequest.Destination,
		EventTypes:           postRequest.EventTypes,
		MessageIds:           postRequest.MessageIds,
		ResourceTypes:        postRequest.ResourceTypes,
		Protocol:             postRequest.Protocol,
		SubscriptionType:     postRequest.SubscriptionType,
		EventFormatType:      postRequest.EventFormatType,
		SubordinateResources: postRequest.SubordinateResources,
		HTTPHeaders:          httpHeadersSlice,
		Context:              postRequest.Context,
		DeliveryRetryPolicy:  postRequest.DeliveryRetryPolicy,
	}
	res, err := e.IsEventsSubscribed("", origin, &subscriptionPost, plugin, target, collectionFlag, collectionName)
	if err != nil {
		resp.Response = res.Body
		resp.StatusCode = int(res.StatusCode)
		return "", resp
	}
	if collectionFlag {
		l.Log.Info("Saving device subscription details of collection subscription")
		if collectionName == "AggregateCollections" {
			resp.StatusCode = http.StatusCreated
			resp.Response = createEventSubscriptionResponse()
			return collectionName, resp
		}
		err = e.saveDeviceSubscriptionDetails(evmodel.Subscription{
			Location:       "",
			EventHostIP:    collectionName,
			OriginResource: origin,
		})
		if err != nil {
			errorMessage := "error while trying to save event subscription of device data: " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			l.Log.Error(errorMessage)
			return "", resp
		}
		resp.StatusCode = http.StatusCreated
		resp.Response = createEventSubscriptionResponse()
		return collectionName, resp
	}

	postBody, _ := json.Marshal(subscriptionPost)
	var reqData string
	//replacing the request url with south bound translation URL
	for key, value := range config.Data.URLTranslation.SouthBoundURL {
		reqData = strings.Replace(string(postBody), key, value, -1)
	}

	target.PostBody = []byte(reqData)
	contactRequest.URL = "/ODIM/v1/Subscriptions"
	contactRequest.HTTPMethodType = http.MethodPost
	contactRequest.PostBody = target

	l.Log.Debug("Subscription Request: " + reqData)
	response, err := e.callPlugin(context.TODO(), contactRequest)
	if err != nil {
		if evcommon.GetPluginStatus(plugin) {
			response, err = e.callPlugin(context.TODO(), contactRequest)
		}
		if err != nil {
			errorMessage := "error while unmarshal the body : " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			l.Log.Error(errorMessage)
			return "", resp
		}
	}
	defer response.Body.Close()
	l.Log.Debug("Subscription Response StatusCode: " + strconv.Itoa(int(response.StatusCode)))
	if response.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			errorMessage := "error while trying to read response body: " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			l.Log.Error(errorMessage)
			return "", resp
		}
		l.Log.Info("Subscription Response: " + string(body))
		var res interface{}
		err = json.Unmarshal(body, &res)
		if err != nil {
			errorMessage := "error while unmarshal the body : " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			l.Log.Error(errorMessage)
			return "", resp
		}

		errorMessage := "error while trying to create event subscription"
		resp.Response = res
		resp.StatusCode = response.StatusCode
		l.Log.Error(errorMessage)
		return "", resp
	}
	// if Subscription location is empty then don't store event details in DB
	locationHdr := response.Header.Get("location")
	if locationHdr == "" {
		errorMessage := "Subscription Location is missing in the response header"
		evcommon.GenEventErrorResponse(errorMessage, errors.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		l.Log.Error(errorMessage)
		return "", resp
	}
	// get the ip address from the host name
	deviceIPAddress, errorMessage := evcommon.GetIPFromHostName(target.ManagerAddress)
	if errorMessage != "" {
		evcommon.GenEventErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
			&resp, []interface{}{"ManagerAddress", target.ManagerAddress})
		l.Log.Error(errorMessage)
		return "", resp
	}
	l.Log.Debug("Saving device subscription details : ", deviceIPAddress)
	evtSubscription := evmodel.Subscription{
		Location:       locationHdr,
		EventHostIP:    deviceIPAddress,
		OriginResource: origin,
	}

	host, _, err := net.SplitHostPort(target.ManagerAddress)
	if err != nil {
		host = target.ManagerAddress
	}
	if !(strings.Contains(locationHdr, host)) {
		evtSubscription.Location = "https://" + target.ManagerAddress + locationHdr
	}
	err = e.saveDeviceSubscriptionDetails(evtSubscription)
	if err != nil {
		errorMessage := "error while trying to save event subscription of device data: " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		l.Log.Error(errorMessage)
		return "", resp
	}
	var outBody interface{}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errorMessage := "error while reading body  : " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		l.Log.Error(errorMessage)
		return "", resp
	}
	err = json.Unmarshal(body, &outBody)
	if err != nil {
		errorMessage := "error while unmarshal the body : " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		l.Log.Error(errorMessage)
		return "", resp
	}
	resp.Response = outBody
	resp.StatusCode = response.StatusCode
	resp.Location = response.Header.Get("location")
	return deviceIPAddress, resp
}

// IsEventsSubscribed is to check events already subscribed.
// if event already subscribed then will do search the subscription details in db against host IP
// if data found then delete the entry in db and get the event types
// and also delete the subscription on device also
// subscription: New Subscription
// subscriptionDetails : subscription details stored in db for the particular device
func (e *ExternalInterfaces) IsEventsSubscribed(token, origin string, subscription *evmodel.EvtSubPost, plugin *evmodel.Plugin, target *evmodel.Target, collectionFlag bool, collectionName string) (errResponse.RPC, error) {
	var resp errResponse.RPC
	var err error
	var host, originResource, searchKey string
	// if Origin is collection then setting host with collection name
	if collectionFlag {
		host = collectionName
		searchKey = collectionName
	} else {
		host, errorMessage := evcommon.GetIPFromHostName(target.ManagerAddress)
		if errorMessage != "" {
			evcommon.GenErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
				[]interface{}{"ManagerAddress", target.ManagerAddress}, &resp)
			l.Log.Error(errorMessage)
			return resp, err
		}
		l.Log.Info("After look up, manager address is: ", host)
		searchKey = evcommon.GetSearchKey(host, evmodel.SubscriptionIndex)
	}
	// uniqueMap is to ignore duplicate eventTypes
	// eventTypes from request  and eventTypes from the all destinations stored in the DB
	uniqueMap := make(map[string]string)

	// add all events to map to remove duplicate eventTypes
	// this need to be remove after the designation uniqueness check added
	for _, eventType := range subscription.EventTypes {
		uniqueMap[eventType] = eventType
	}
	var (
		eventTypes    = subscription.EventTypes
		messageIDs    = subscription.MessageIds
		resourceTypes = subscription.ResourceTypes
	)

	originResource = origin
	subscriptionDetails, err := e.GetEvtSubscriptions(searchKey)
	if err != nil && !strings.Contains(err.Error(), "No data found for the key") {
		errorMessage := "Error while get subscription details: " + err.Error()
		evcommon.GenErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			[]interface{}{}, &resp)
		l.Log.Error(errorMessage)
		return resp, err
	}
	// if there is no subscription happened then create event subscription
	if len(subscriptionDetails) < 1 {
		return resp, nil
	}

	var subscriptionPresent bool
	for index, evtSubscriptions := range subscriptionDetails {

		if isHostPresent(evtSubscriptions.Hosts, host) {
			subscriptionPresent = true

			if len(evtSubscriptions.EventTypes) > 0 && (index == 0 || len(eventTypes) > 0) {
				eventTypes = append(eventTypes, evtSubscriptions.EventTypes...)
			} else {
				eventTypes = []string{}
			}

			if len(evtSubscriptions.MessageIds) > 0 && (index == 0 || len(messageIDs) > 0) {
				messageIDs = append(messageIDs, evtSubscriptions.MessageIds...)
			} else {
				messageIDs = []string{}
			}

			if len(evtSubscriptions.ResourceTypes) > 0 && (index == 0 || len(resourceTypes) > 0) {
				resourceTypes = append(resourceTypes, evtSubscriptions.ResourceTypes...)
			} else {
				resourceTypes = []string{}
			}

		}
	}
	if !subscriptionPresent {
		return resp, nil
	}
	if !collectionFlag {
		l.Log.Info("Delete Subscription from device")
		if strings.Contains(originResource, "Fabrics") {
			resp, err = e.DeleteFabricsSubscription(originResource, plugin)
			if err != nil {
				return resp, err
			}
		} else {
			resp, err = e.DeleteSubscriptions(originResource, token, plugin, target)
			if err != nil {
				return resp, err
			}
		}
	}
	// updating the subscription information

	eventTypesCount := len(eventTypes)
	messageIDsCount := len(messageIDs)
	resourceTypesCount := len(resourceTypes)
	removeDuplicatesFromSlice(&eventTypes, &eventTypesCount)
	removeDuplicatesFromSlice(&messageIDs, &messageIDsCount)
	removeDuplicatesFromSlice(&resourceTypes, &resourceTypesCount)
	subscription.EventTypes = eventTypes
	subscription.MessageIds = messageIDs
	subscription.ResourceTypes = resourceTypes
	return resp, nil
}

// CreateDefaultEventSubscription is creates the  subscription with event types which will be required to rediscover the inventory
// after computer system restarts ,This will  triggered from   aggregation service whenever a computer system is added
func (e *ExternalInterfaces) CreateDefaultEventSubscription(originResources, eventTypes, messageIDs, resourceTypes []string, protocol string) errResponse.RPC {
	l.Log.Info("Creation of default subscriptions started for: " + strings.Join(originResources, "::"))
	var resp errResponse.RPC
	var response evresponse.EventResponse
	var partialResultFlag bool
	if protocol == "" {
		protocol = "Redfish"
	}
	bubbleUpStatusCode := http.StatusCreated
	var postRequest evmodel.RequestBody
	postRequest.Destination = ""
	postRequest.EventTypes = eventTypes
	postRequest.MessageIds = messageIDs
	postRequest.ResourceTypes = resourceTypes
	postRequest.Context = "Creating the Default Event Subscription"
	postRequest.Protocol = protocol
	postRequest.SubscriptionType = evmodel.SubscriptionType
	postRequest.SubordinateResources = true
	_, response = e.eventSubscription(postRequest, originResources[0], "", false)
	e.checkCollectionSubscription(originResources[0], protocol)
	if response.StatusCode != http.StatusCreated {
		partialResultFlag = true
		if response.StatusCode > bubbleUpStatusCode {
			bubbleUpStatusCode = response.StatusCode
		}
	}

	if !partialResultFlag || len(originResources) == 1 {
		resp.StatusCode = int32(response.StatusCode)
	} else {
		resp.StatusCode = int32(bubbleUpStatusCode)
	}

	// Removed creation of default subscription for each server add, adding only single default subscription at time of deployment with subscriptionID 0

	// subscriptionID := uuid.New().String()
	// evtSubscription := evmodel.Subscription{
	// 	SubscriptionID:       subscriptionID,
	// 	EventTypes:           eventTypes,
	// 	MessageIds:           messageIDs,
	// 	ResourceTypes:        resourceTypes,
	// 	OriginResources:      originResources,
	// 	Hosts:                []string{host},
	// 	Protocol:             protocol,
	// 	SubscriptionType:     evmodel.SubscriptionType,
	// 	SubordinateResources: true,
	// }
	// err := e.SaveEventSubscription(evtSubscription)
	// if err != nil {
	// 	errorMessage := "error while trying to save event subscription data: " + err.Error()
	// 	evcommon.GenErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
	// 		[]interface{}{}, &resp)
	// 	l.Log.Error(errorMessage)
	// 	return resp
	// }

	resp.Body = response.Response
	resp.StatusCode = http.StatusCreated
	l.Log.Info("Creation of default subscriptions completed for : " + strings.Join(originResources, "::"))
	return resp
}

// saveDeviceSubscriptionDetails will first check if already origin resource details present
// if its present then Update location
// otherwise add an entry to redis
func (e *ExternalInterfaces) saveDeviceSubscriptionDetails(evtSubscription evmodel.Subscription) error {
	searchKey := evcommon.GetSearchKey(evtSubscription.EventHostIP, evmodel.DeviceSubscriptionIndex)
	deviceSubscription, _ := e.GetDeviceSubscriptions(searchKey)
	var newDevSubscription = evmodel.DeviceSubscription{
		EventHostIP:     evtSubscription.EventHostIP,
		Location:        evtSubscription.Location,
		OriginResources: []string{evtSubscription.OriginResource},
	}
	// if device subscriptions details for the device is present in db then don't add again
	var save = true
	if deviceSubscription != nil {

		save = true
		// if the origin resource is present in device subscription details then don't add
		for _, originResource := range deviceSubscription.OriginResources {
			if evtSubscription.OriginResource == originResource {
				save = false
			} else {
				newDevSubscription.OriginResources = append(newDevSubscription.OriginResources, originResource)
				save = false
			}
		}
		err := e.UpdateDeviceSubscriptionLocation(newDevSubscription)
		if err != nil {
			return err
		}
	}
	if save {
		return e.SaveDeviceSubscription(newDevSubscription)
	}
	return nil
}

func (e *ExternalInterfaces) getTargetDetails(origin string) (*evmodel.Target, evresponse.EventResponse, error) {
	var resp evresponse.EventResponse
	uuid, err := getUUID(origin)
	if err != nil {
		evcommon.GenEventErrorResponse(err.Error(), errResponse.ResourceNotFound, http.StatusNotFound,
			&resp, []interface{}{"System", origin})
		l.Log.Error(err.Error())
		return nil, resp, err
	}

	// Get target device Credentials from using device UUID
	target, err := e.GetTarget(uuid)
	if err != nil {
		// Frame the RPC response body and response Header below

		errorMessage := "error while getting Systems(Target device Credentials) table details: " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
			&resp, []interface{}{"Systems", origin})
		l.Log.Error(errorMessage)
		return nil, resp, err
	}
	decryptedPasswordByte, err := DecryptWithPrivateKeyFunc(target.Password)
	if err != nil {
		// Frame the RPC response body and response Header below
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		l.Log.Error(errorMessage)
		return nil, resp, err
	}
	target.Password = decryptedPasswordByte
	return target, resp, nil
}

// DeleteSubscriptions will delete subscription from device
func (e *ExternalInterfaces) DeleteSubscriptions(originResource, token string, plugin *evmodel.Plugin, target *evmodel.Target) (errResponse.RPC, error) {
	var resp errResponse.RPC
	var err error
	var deviceSubscription *evmodel.DeviceSubscription

	addr, errorMessage := evcommon.GetIPFromHostName(target.ManagerAddress)
	if errorMessage != "" {
		evcommon.GenErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
			[]interface{}{"ManagerAddress", target.ManagerAddress}, &resp)
		l.Log.Error(errorMessage)
		return resp, err
	}
	searchKey := evcommon.GetSearchKey(addr, evmodel.DeviceSubscriptionIndex)
	deviceSubscription, err = e.GetDeviceSubscriptions(searchKey)

	if err != nil {
		// if its first subscription then no need to check events subscribed
		if strings.Contains(err.Error(), "No data found for the key") {
			return resp, nil
		}
		errorMessage := "Error while get subscription details: " + err.Error()
		evcommon.GenErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			[]interface{}{}, &resp)

		l.Log.Error(errorMessage)
		return resp, err
	}

	var contactRequest evcommon.PluginContactRequest

	contactRequest.Plugin = plugin
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token := e.getPluginToken(plugin)
		if token == "" {
			evcommon.GenErrorResponse("error: Unable to create session with plugin "+plugin.ID, errResponse.NoValidSession, http.StatusUnauthorized,
				[]interface{}{}, &resp)
			return resp, fmt.Errorf("error: Unable to create session with plugin " + plugin.ID)
		}
		contactRequest.Token = token

	} else {
		contactRequest.LoginCredential = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	target.Location = deviceSubscription.Location

	// Call to delete subscription to plugin
	contactRequest.URL = "/ODIM/v1/Subscriptions"
	contactRequest.HTTPMethodType = http.MethodDelete
	contactRequest.PostBody = target

	resp, _, _, err = e.PluginCall(contactRequest)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (e *ExternalInterfaces) createEventSubscription(taskID string, subTaskChan chan<- int32, reqSessionToken string, targetURI string, request evmodel.RequestBody, originResource string, result *evresponse.MutexLock, wg *sync.WaitGroup, collectionFlag bool, collectionName string, aggregateResource string, isAggregateCollection bool) {
	var (
		subTaskURI      string
		subTaskID       string
		reqBody         []byte
		reqJSON         string
		err             error
		resp            errResponse.RPC
		percentComplete int32
	)
	defer wg.Done()

	reqBody, err = json.Marshal(request)
	if err != nil {
		l.Log.Error("error while trying to marshal create event request: " + err.Error())
	}
	reqJSON = string(reqBody)
	if taskID != "" {
		subTaskURI, err = e.CreateChildTask(reqSessionToken, taskID)
		if err != nil {
			l.Log.Error("Error while creating the SubTask")
		}
		trimmedURI := strings.TrimSuffix(subTaskURI, "/")
		subTaskID = trimmedURI[strings.LastIndex(trimmedURI, "/")+1:]
		resp.StatusCode = http.StatusAccepted
		e.UpdateTask(fillTaskData(subTaskID, targetURI, reqJSON, resp, common.Running, common.OK, percentComplete, http.MethodPost))
	}

	host, response := e.eventSubscription(request, originResource, collectionName, collectionFlag)
	resp.Body = response.Response
	resp.StatusCode = int32(response.StatusCode)
	if isAggregateCollection {
		if resp.StatusCode == http.StatusConflict {
			response.StatusCode = http.StatusCreated
		}
		result.AddResponse(aggregateResource, getAggregateID(aggregateResource), response)
	} else {
		result.AddResponse(originResource, host, response)
	}
	percentComplete = 100
	if subTaskID != "" {
		if response.StatusCode != http.StatusCreated {
			e.UpdateTask(fillTaskData(subTaskID, targetURI, reqJSON, resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
		} else {
			e.UpdateTask(fillTaskData(subTaskID, targetURI, reqJSON, resp, common.Completed, common.OK, percentComplete, http.MethodPost))
		}
		subTaskChan <- int32(response.StatusCode)
	}
}

// checkCollectionSubscription checks if any collection based subscription exists
// If its' exists it will  update the existing subscription information with newly added server origin
func (e *ExternalInterfaces) checkCollectionSubscription(origin, protocol string) {
	//Creating key to get all the System Collection subscription
	var searchKey string
	var bmcFlag bool
	if strings.Contains(origin, "Fabrics") {
		searchKey = "/redfish/v1/Fabrics"
	} else {
		bmcFlag = true
		searchKey = "/redfish/v1/Systems"
	}
	subscriptions, err := e.GetEvtSubscriptions(searchKey)
	if err != nil {
		return
	}
	var chassisSubscriptions, managersSubscriptions []evmodel.Subscription
	if bmcFlag {
		chassisSubscriptions, _ = e.GetEvtSubscriptions("/redfish/v1/Chassis")
		subscriptions = append(subscriptions, chassisSubscriptions...)
		managersSubscriptions, _ = e.GetEvtSubscriptions("/redfish/v1/Managers")
		subscriptions = append(subscriptions, managersSubscriptions...)
	}
	// Checking the collection based subscription
	var collectionSubscription = make([]evmodel.Subscription, 0)
	for _, evtSubscription := range subscriptions {
		for _, originResource := range evtSubscription.OriginResources {
			if strings.Contains(origin, "Systems") && (originResource == "/redfish/v1/Systems" || originResource == "/redfish/v1/Chassis" || originResource == "/redfish/v1/Managers") {
				collectionSubscription = append(collectionSubscription, evtSubscription)
			} else if strings.Contains(origin, "Fabrics") && originResource == "/redfish/v1/Fabrics" {
				collectionSubscription = append(collectionSubscription, evtSubscription)
			}
		}
	}
	if len(collectionSubscription) <= 0 {
		return
	}
	// using the one of the destination
	var destination string
	var context string
	var eventTypes, messageIDs, resourceTypes []string
	for index, evtSubscription := range collectionSubscription {
		destination = evtSubscription.Destination
		if len(evtSubscription.EventTypes) > 0 && (index == 0 || len(eventTypes) > 0) {
			eventTypes = append(eventTypes, evtSubscription.EventTypes...)
		} else {
			eventTypes = []string{}
		}

		if len(evtSubscription.MessageIds) > 0 && (index == 0 || len(messageIDs) > 0) {
			messageIDs = append(messageIDs, evtSubscription.MessageIds...)
		} else {
			messageIDs = []string{}
		}

		if len(evtSubscription.ResourceTypes) > 0 && (index == 0 || len(resourceTypes) > 0) {
			resourceTypes = append(resourceTypes, evtSubscription.ResourceTypes...)
		} else {
			resourceTypes = []string{}
		}
	}
	eventTypesCount := len(eventTypes)
	messageIDsCount := len(messageIDs)
	resourceTypesCount := len(resourceTypes)

	removeDuplicatesFromSlice(&eventTypes, &eventTypesCount)
	removeDuplicatesFromSlice(&messageIDs, &messageIDsCount)
	removeDuplicatesFromSlice(&resourceTypes, &resourceTypesCount)

	subordinateFlag := false
	if strings.Contains(origin, "Fabrics") {
		subordinateFlag = true
	}

	subscriptionPost := evmodel.RequestBody{
		EventTypes:           eventTypes,
		MessageIds:           messageIDs,
		ResourceTypes:        resourceTypes,
		Context:              context,
		Destination:          destination,
		Protocol:             protocol,
		SubordinateResources: subordinateFlag,
	}
	subscriptionPost.OriginResources = []evmodel.OdataIDLink{
		{
			OdataID: origin,
		},
	}

	// Subscribing newly added server with collated event list
	host, response := e.eventSubscription(subscriptionPost, origin, "", false)
	if response.StatusCode != http.StatusCreated {
		return
	}

	// Get Device Subscription Details if collection is bmc and update chassis and managers uri
	if bmcFlag {
		searchKey := evcommon.GetSearchKey(host, evmodel.DeviceSubscriptionIndex)
		deviceSubscription, _ := e.GetDeviceSubscriptions(searchKey)
		data := strings.Split(origin, "/redfish/v1/Systems/")
		chassisList, _ := e.GetAllMatchingDetails("Chassis", data[1], common.InMemory)
		managersList, _ := e.GetAllMatchingDetails("Managers", data[1], common.InMemory)
		var newDevSubscription = evmodel.DeviceSubscription{
			EventHostIP:     deviceSubscription.EventHostIP,
			Location:        deviceSubscription.Location,
			OriginResources: deviceSubscription.OriginResources,
		}
		newDevSubscription.OriginResources = append(newDevSubscription.OriginResources, chassisList...)
		newDevSubscription.OriginResources = append(newDevSubscription.OriginResources, managersList...)

		err := e.UpdateDeviceSubscriptionLocation(newDevSubscription)
		if err != nil {
			l.Log.Error("Error while Updating Device subscription : " + err.Error())
		}
	}
}

func (e *ExternalInterfaces) createFabricSubscription(postRequest evmodel.RequestBody, origin, collectionName string, collectionFlag bool) (string, evresponse.EventResponse) {
	var resp evresponse.EventResponse
	var err error
	var plugin *evmodel.Plugin
	var contactRequest evcommon.PluginContactRequest
	// Extract the fabric id from the Origin
	fabricID := getFabricID(origin)
	fabric, dberr := e.GetFabricData(fabricID)
	if dberr != nil {
		errorMessage := "error while getting fabric data: " + dberr.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
			&resp, []interface{}{"Fabrics", fabricID})
		l.Log.Error(errorMessage)
		return "", resp
	}
	var gerr *errors.Error
	plugin, gerr = e.GetPluginData(fabric.PluginID)
	if gerr != nil {
		errorMessage := "error while getting plugin data: " + gerr.Error() + fabric.PluginID
		evcommon.GenEventErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
			&resp, []interface{}{"Plugin", fabric.PluginID})
		l.Log.Error(errorMessage)
		return "", resp
	}
	contactRequest.Plugin = plugin
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token := e.getPluginToken(plugin)
		if token == "" {
			evcommon.GenEventErrorResponse("error: Unable to create session with plugin "+plugin.ID, errResponse.NoValidSession, http.StatusUnauthorized,
				&resp, []interface{}{})
			l.Log.Error("error: Unable to create session with plugin " + plugin.ID)
			return "", resp
		}
		contactRequest.Token = token
	} else {
		contactRequest.LoginCredential = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
	}
	var httpHeadersSlice = make([]evmodel.HTTPHeaders, 0)
	httpHeadersSlice = append(httpHeadersSlice, evmodel.HTTPHeaders{ContentType: "application/json"})
	subscriptionPost := evmodel.EvtSubPost{
		Name:                 postRequest.Name,
		Destination:          postRequest.Destination,
		EventTypes:           postRequest.EventTypes,
		MessageIds:           postRequest.MessageIds,
		ResourceTypes:        postRequest.ResourceTypes,
		Protocol:             postRequest.Protocol,
		SubscriptionType:     postRequest.SubscriptionType,
		EventFormatType:      postRequest.EventFormatType,
		SubordinateResources: postRequest.SubordinateResources,
		HTTPHeaders:          httpHeadersSlice,
		Context:              postRequest.Context,
		OriginResources: []evmodel.OdataIDLink{
			{
				OdataID: origin,
			},
		},
	}
	// Assigning a proper empty slice to slices with nil value.
	// This will make those slices give [] instead of null.
	var emptySlice []string
	if len(subscriptionPost.EventTypes) == 0 {
		subscriptionPost.EventTypes = emptySlice
	}
	if len(subscriptionPost.MessageIds) == 0 {
		subscriptionPost.MessageIds = emptySlice
	}
	if len(subscriptionPost.ResourceTypes) == 0 {
		subscriptionPost.ResourceTypes = emptySlice
	}
	deviceIPAddress, errorMessage := GetIPFromHostNameFunc(plugin.IP)
	if errorMessage != "" {
		evcommon.GenEventErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusBadRequest,
			&resp, []interface{}{"ManagerAddress", plugin.IP})
		l.Log.Error(errorMessage)
		return "", resp
	}
	var target = evmodel.Target{
		ManagerAddress: deviceIPAddress,
	}
	res, err := e.IsEventsSubscribed("", origin, &subscriptionPost, plugin, &target, collectionFlag, collectionName)
	if err != nil {
		resp.Response = res.Body
		resp.StatusCode = int(res.StatusCode)
		return "", resp
	}

	postBody, _ := json.Marshal(subscriptionPost)
	var reqData string
	//replacing the request url with south bound translation URL
	for key, value := range config.Data.URLTranslation.SouthBoundURL {
		reqData = strings.Replace(string(postBody), key, value, -1)
	}
	contactRequest.URL = "/ODIM/v1/Subscriptions"
	contactRequest.HTTPMethodType = http.MethodPost
	err = json.Unmarshal([]byte(reqData), &contactRequest.PostBody)

	response, err := e.callPlugin(context.TODO(), contactRequest)
	if err != nil {
		if evcommon.GetPluginStatus(plugin) {
			response, err = e.callPlugin(context.TODO(), contactRequest)
		}
		if err != nil {
			evcommon.GenEventErrorResponse(err.Error(), errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			l.Log.Error(err.Error())
			return "", resp
		}
	}
	defer response.Body.Close()
	//retrying the operation if status code is 401
	if response.StatusCode == http.StatusUnauthorized && strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		response, resp, err = e.retryEventSubscriptionOperation(contactRequest)
		if err != nil {
			return "", resp
		}
	}

	l.Log.Debug("Subscription Response Status Code: " + string(rune(response.StatusCode)))
	if response.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			errorMessage := "error while trying to read response body: " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			l.Log.Error(errorMessage)
			return "", resp
		}
		errorMessage := "error while trying to create event subscription"
		var res interface{}
		l.Log.Error("Subscription Response " + string(body))
		err = json.Unmarshal(body, &res)
		if err != nil {
			errorMessage := "error while unmarshal the body : " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			l.Log.Error(errorMessage)
			return "", resp
		}

		resp.Response = res
		resp.StatusCode = response.StatusCode
		l.Log.Error(errorMessage)
		return "", resp
	}

	evtSubscription := evmodel.Subscription{
		EventHostIP:    deviceIPAddress,
		OriginResource: origin,
	}

	evtSubscription.Location = response.Header.Get("location")
	err = e.saveDeviceSubscriptionDetails(evtSubscription)
	if err != nil {
		errorMessage := "error while trying to save event subscription of device data: " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		l.Log.Error(errorMessage)
		return "", resp
	}

	resp.Response = createEventSubscriptionResponse()
	resp.StatusCode = response.StatusCode
	resp.Location = response.Header.Get("location")
	return deviceIPAddress, resp
}

// UpdateEventSubscriptions it will add subscription for newly Added system in aggregate
func (e *ExternalInterfaces) UpdateEventSubscriptions(req *eventsproto.EventUpdateRequest, isRemove bool) error {
	// var resp response.RPC
	authResp, err := e.Auth(req.SessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if err != nil {
		l.Log.Errorf("Error while authorizing the session token : %s", err.Error())
	}
	if authResp.StatusCode != http.StatusOK {
		l.Log.Printf("error while trying to authenticate session: status code: %v, status message: %v", authResp.StatusCode, authResp.StatusMessage)
		return nil
	}
	var plugin *evmodel.Plugin
	var contactRequest evcommon.PluginContactRequest
	var target *evmodel.Target

	target, _, err = e.getTargetDetails(req.SystemID)
	if err != nil {
		return err
	}
	var errs *errors.Error
	plugin, errs = e.GetPluginData(target.PluginID)
	if errs != nil {
		errorMessage := "error while getting plugin data: " + errs.Error()
		l.Log.Info(errorMessage)
		return err
	}

	contactRequest.Plugin = plugin
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token := e.getPluginToken(plugin)
		if token == "" {
			l.Log.Info("error: Unable to create session with plugin " + plugin.ID)
			return nil
		}
		contactRequest.Token = token

	} else {
		contactRequest.LoginCredential = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
	}

	var httpHeadersSlice = make([]evmodel.HTTPHeaders, 0)
	httpHeadersSlice = append(httpHeadersSlice, evmodel.HTTPHeaders{ContentType: "application/json"})
	subscriptionPost := evmodel.EvtSubPost{
		EventTypes:    []string{},
		MessageIds:    []string{},
		ResourceTypes: []string{},
		OriginResources: []evmodel.OdataIDLink{
			{
				OdataID: req.SystemID,
			},
		},
		SubordinateResources: true,
		Protocol:             "Redfish",
		SubscriptionType:     evmodel.SubscriptionType,
		HTTPHeaders:          httpHeadersSlice,
		Context:              evmodel.Context,
		DeliveryRetryPolicy:  "RetryForever",
		EventFormatType:      "Event",
	}
	_, err = e.UpdateEventsSubscribed("", req.SystemID, &subscriptionPost, plugin, target, false, "", true, req.AggregateId, isRemove)
	if err != nil {

		return err
	}
	postBody, _ := json.Marshal(subscriptionPost)
	var reqData string
	//replacing the request url with south bound translation URL
	for key, value := range config.Data.URLTranslation.SouthBoundURL {
		reqData = strings.Replace(string(postBody), key, value, -1)
	}

	target.PostBody = []byte(reqData)
	contactRequest.URL = "/ODIM/v1/Subscriptions"
	contactRequest.HTTPMethodType = http.MethodPost
	contactRequest.PostBody = target

	l.Log.Debug("Subscription Request: " + reqData)
	response, err := e.callPlugin(context.TODO(), contactRequest)
	if err != nil {
		if evcommon.GetPluginStatus(plugin) {
			response, err = e.callPlugin(context.TODO(), contactRequest)
		}
		if err != nil {
			errorMessage := "error while unmarshal the body : " + err.Error()
			l.Log.Info(errorMessage)
			return err
		}
	}
	defer response.Body.Close()
	l.Log.Debug("Subscription Response StatusCode: " + strconv.Itoa(int(response.StatusCode)))
	if response.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			errorMessage := "error while trying to read response body: " + err.Error()
			l.Log.Info(errorMessage)
			return nil
		}
		l.Log.Info("Subscription Response: " + string(body))
		var res interface{}
		err = json.Unmarshal(body, &res)
		if err != nil {
			errorMessage := "error while unmarshal the body : " + err.Error()
			l.Log.Error(errorMessage)
			return nil
		}

		errorMessage := "error while trying to create event subscription"
		l.Log.Error(errorMessage)
		return nil
	}
	// if Subscription location is empty then don't store event details in DB
	locationHdr := response.Header.Get("location")
	if locationHdr == "" {
		errorMessage := "Subscription Location is missing in the response header"
		l.Log.Info(errorMessage)
		return nil
	}
	// get the ip address from the host name
	deviceIPAddress, errorMessage := evcommon.GetIPFromHostName(target.ManagerAddress)
	if errorMessage != "" {
		l.Log.Info(errorMessage)
	}
	l.Log.Info("Saving device subscription details : ", deviceIPAddress)
	evtSubscription := evmodel.Subscription{
		Location:       locationHdr,
		EventHostIP:    deviceIPAddress,
		OriginResource: req.SystemID,
	}

	host, _, err := net.SplitHostPort(target.ManagerAddress)
	if err != nil {
		host = target.ManagerAddress
	}
	if !(strings.Contains(locationHdr, host)) {
		evtSubscription.Location = "https://" + target.ManagerAddress + locationHdr
	}
	err = e.saveDeviceSubscriptionDetails(evtSubscription)
	if err != nil {
		errorMessage := "error while trying to save event subscription of device data: " + err.Error()
		l.Log.Error(errorMessage)
		return nil
	}

	return nil
}

// UpdateEventsSubscribed is to check events already subscribed.
// if event already subscribed then will do search the subscription details in db against host IP
// if data found then delete the entry in db and get the event types
// and also delete the subscription on device also
// subscription: New Subscription
// subscriptionDetails : subscription details stored in db for the particular device
func (e *ExternalInterfaces) UpdateEventsSubscribed(token, origin string, subscription *evmodel.EvtSubPost, plugin *evmodel.Plugin, target *evmodel.Target, collectionFlag bool, collectionName string, isAggregate bool, aggregateID string, isRemove bool) (errResponse.RPC, error) {
	var resp errResponse.RPC
	var err error
	var host, originResource, searchKey string
	// if Origin is collection then setting host with collection name
	if collectionFlag {
		host = collectionName
		searchKey = collectionName
	} else {
		host1, errorMessage := GetIPFromHostNameFunc(target.ManagerAddress)
		host = host1
		if errorMessage != "" {
			evcommon.GenErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
				[]interface{}{"ManagerAddress", target.ManagerAddress}, &resp)
			l.Log.Error(errorMessage)
			return resp, err
		}
		l.Log.Info("After look up, manager address is: ", host)
		searchKey = evcommon.GetSearchKey(host, evmodel.SubscriptionIndex)
	}
	// uniqueMap is to ignore duplicate eventTypes
	// eventTypes from request  and eventTypes from the all destinations stored in the DB
	uniqueMap := make(map[string]string)

	// add all events to map to remove duplicate eventTypes
	// this need to be remove after the designation uniqueness check added
	for _, eventType := range subscription.EventTypes {
		uniqueMap[eventType] = eventType
	}
	var (
		eventTypes    = subscription.EventTypes
		messageIDs    = subscription.MessageIds
		resourceTypes = subscription.ResourceTypes
	)
	originResource = origin
	subscriptionDetails, err := e.GetEvtSubscriptions(searchKey)
	if err != nil && !strings.Contains(err.Error(), "No data found for the key") {
		errorMessage := "Error while get subscription details: " + err.Error()
		evcommon.GenErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			[]interface{}{}, &resp)
		l.Log.Error(errorMessage)
		return resp, err
	}
	var subscriptionPresent, isAggregateSubscriptionPresent bool
	var aggregateSubscriptionDetails []evmodel.Subscription
	// get all aggregate subscription
	if isAggregate {
		searchKeyAgg := evcommon.GetSearchKey(host, evmodel.SubscriptionIndex)
		aggregateList, err := e.GetAggregateList(searchKeyAgg)
		if err != nil {
			l.Log.Info("No Aggregate subscription Found ", err)
		}
		for _, id := range aggregateList {
			if isRemove {
				if id == aggregateID {
					continue
				}
			}
			searchKey = evcommon.GetSearchKey(id, evmodel.SubscriptionIndex)
			aggregateSubscriptionDetails, err = e.GetEvtSubscriptions(searchKey)

			if err != nil && !strings.Contains(err.Error(), "No data found for the key") {
				l.Log.Info("Error while get aggregateSubscriptionDetails details: " + err.Error())
			}
			for index, evtSubscriptions := range aggregateSubscriptionDetails {
				if isHostPresent(evtSubscriptions.Hosts, aggregateID) {
					isAggregateSubscriptionPresent = true
					if len(evtSubscriptions.EventTypes) > 0 && (index == 0 || len(eventTypes) > 0) {
						eventTypes = append(eventTypes, evtSubscriptions.EventTypes...)
					}
					if len(evtSubscriptions.MessageIds) > 0 && (index == 0 || len(messageIDs) > 0) {
						messageIDs = append(messageIDs, evtSubscriptions.MessageIds...)
					}
					if len(evtSubscriptions.ResourceTypes) > 0 && (index == 0 || len(resourceTypes) > 0) {
						resourceTypes = append(resourceTypes, evtSubscriptions.ResourceTypes...)
					}
				}
			}

		}
	}

	// if there is no subscription happened then create event subscription

	if len(subscriptionDetails) < 1 && len(aggregateSubscriptionDetails) < 1 {
		return resp, nil
	}
	for index, evtSubscriptions := range subscriptionDetails {
		if isHostPresent(evtSubscriptions.Hosts, host) {
			subscriptionPresent = true
			if len(evtSubscriptions.EventTypes) > 0 && (index == 0 || len(eventTypes) > 0) {
				eventTypes = append(eventTypes, evtSubscriptions.EventTypes...)
			}
			if len(evtSubscriptions.MessageIds) > 0 && (index == 0 || len(messageIDs) > 0) {
				messageIDs = append(messageIDs, evtSubscriptions.MessageIds...)
			}

			if len(evtSubscriptions.ResourceTypes) > 0 && (index == 0 || len(resourceTypes) > 0) {
				resourceTypes = append(resourceTypes, evtSubscriptions.ResourceTypes...)
			}

		}
	}
	if !subscriptionPresent && !isAggregateSubscriptionPresent {
		return resp, nil
	}
	if !collectionFlag {
		l.Log.Info("Delete Subscription from device")
		if strings.Contains(originResource, "Fabrics") {
			resp, err = e.DeleteFabricsSubscription(originResource, plugin)
			if err != nil {
				return resp, err
			}
		} else {
			resp, err = e.DeleteSubscriptions(originResource, token, plugin, target)
			if err != nil {
				return resp, err
			}
		}
	}
	// updating the subscription information
	eventTypesCount := len(eventTypes)
	messageIDsCount := len(messageIDs)
	resourceTypesCount := len(resourceTypes)
	removeDuplicatesFromSlice(&eventTypes, &eventTypesCount)
	removeDuplicatesFromSlice(&messageIDs, &messageIDsCount)
	removeDuplicatesFromSlice(&resourceTypes, &resourceTypesCount)
	subscription.EventTypes = eventTypes
	subscription.MessageIds = messageIDs
	subscription.ResourceTypes = resourceTypes
	return resp, nil
}
