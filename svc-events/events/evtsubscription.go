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

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
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

// ValidateRequest input request for create subscription
func (e *ExternalInterfaces) ValidateRequest(ctx context.Context, req *eventsproto.EventSubRequest,
	postRequest model.EventDestination) (int32, string, []interface{}, error) {
	invalidProperties, err := common.RequestParamsCaseValidator(req.PostBody, postRequest)
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		return http.StatusInternalServerError, errResponse.InternalError, nil, fmt.Errorf(errMsg)
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in upper camel case "
		return http.StatusBadRequest, errResponse.PropertyUnknown, []interface{}{invalidProperties}, fmt.Errorf(errorMessage)
	}

	//check mandatory fields
	statusCode, statusMessage, messageArgs, invalidFieldError := validateFields(&postRequest)
	if invalidFieldError != nil {
		return statusCode, statusMessage, messageArgs, invalidFieldError
	}

	//validate destination URI in the request
	if !common.URIValidator(postRequest.Destination) {
		errorMessage := "error: request body contains invalid value for Destination field, " + postRequest.Destination
		return http.StatusBadRequest, errResponse.PropertyValueFormatError, []interface{}{postRequest.Destination, "Destination"}, fmt.Errorf(errorMessage)
	}

	// check any of the subscription present for the destination from the request
	// if errored out or no subscriptions then add subscriptions else return an error
	subscriptionDetails, _ := e.GetEvtSubscriptions(postRequest.Destination)
	if len(subscriptionDetails) > 0 {
		return http.StatusConflict, errResponse.ResourceInUse, []interface{}{postRequest.Destination, "Destination"}, fmt.Errorf("subscription already present for the requested destination")
	}
	return http.StatusOK, common.OK, []interface{}{}, nil
}

// CreateEventSubscription is a API to create event subscription
func (e *ExternalInterfaces) CreateEventSubscription(ctx context.Context, taskID string, sessionUserName string, req *eventsproto.EventSubRequest) errResponse.RPC {
	var (
		err             error
		resp            errResponse.RPC
		postRequest     model.EventDestination
		percentComplete int32 = 100
		targetURI             = "/redfish/v1/EventService/Subscriptions"
	)
	if err = json.Unmarshal(req.PostBody, &postRequest); err != nil {
		l.LogWithFields(ctx).Error(err.Error())
		evcommon.GenErrorResponse(err.Error(), errResponse.MalformedJSON, http.StatusBadRequest, []interface{}{}, &resp)
		e.UpdateTask(ctx, fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
		return resp
	}
	// ValidateRequest input request for create subscription
	statusCode, statusMessage, messageArgs, validationErr := e.ValidateRequest(ctx, req, postRequest)
	if validationErr != nil {
		evcommon.GenErrorResponse(validationErr.Error(), statusMessage, statusCode,
			messageArgs, &resp)
		l.LogWithFields(ctx).Error(validationErr.Error())
		e.UpdateTask(ctx, fillTaskData(taskID, targetURI, string(req.PostBody),
			resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
		return resp
	}
	// Get the target device  details from the origin resources
	// Loop through all origin list and form individual event subscription request,
	// Which will then forward to plugin to make subscription with target device
	var wg, taskCollectionWG sync.WaitGroup
	var result = &evresponse.MutexLock{
		Response: make(map[string]evresponse.EventResponse),
		Hosts:    make(map[string]string),
		Lock:     &sync.Mutex{},
	}

	// remove odataid in the origin resources
	originResources := removeOdataIDfromOriginResources(postRequest.OriginResources)

	// check and remove if duplicate OriginResources exist in the request
	removeDuplicatesFromSlice(&originResources)

	// If origin resource is nil then subscribe to all collection
	if len(originResources) == 0 {
		originResources = []string{
			"/redfish/v1/Systems",
			"/redfish/v1/Chassis",
			"/redfish/v1/Fabrics",
			"/redfish/v1/Managers",
			"/redfish/v1/TaskService/Tasks",
		}
	}
	var collectionList = make([]string, 0)
	subTaskChan := make(chan int32, len(originResources))
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
			if i <= len(originResources) {
				percentComplete = int32((i*100)/len(originResources) - 1)
				if resp.StatusCode == 0 {
					resp.StatusCode = http.StatusAccepted
				}
				e.UpdateTask(ctx, fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Running, common.OK, percentComplete, http.MethodPost))
			}
		}
	}()

	for _, origin := range originResources {
		_, err := getUUID(origin)
		if err != nil {
			collection, collectionName, collectionFlag, aggregateResource, isAggregate, _ := e.checkCollection(origin)
			wg.Add(1)
			// for origin is collection
			go e.createEventSubscription(ctx, taskID, subTaskChan, sessionUserName, targetURI, postRequest, origin, result, &wg, collectionFlag, collectionName, aggregateResource, isAggregate)
			for i := 0; i < len(collection); i++ {
				wg.Add(1)
				// for subordinate origin
				go e.createEventSubscription(ctx, "", subTaskChan, sessionUserName, targetURI, postRequest, collection[i], result, &wg, false, "", aggregateResource, isAggregate)
			}
			if !isAggregate {
				collectionList = append(collectionList, collection...)
			}
		} else {
			wg.Add(1)
			go e.createEventSubscription(ctx, taskID, subTaskChan, sessionUserName, targetURI, postRequest, origin, result, &wg, false, "", "", false)
		}
	}

	wg.Wait()
	// close channel once all sub-routines created have exited
	close(subTaskChan)
	// wait till all the subtasks are collected and routine is exited
	taskCollectionWG.Wait()

	var (
		locationHeader             string
		successfulSubscriptionList = make([]model.Link, 0)
	)

	result.Lock.Lock()
	originResourceProcessedCount := len(result.Response)
	successfulSubscriptionList, result.Response = getSuccessfulResponse(result.Response)

	result.Lock.Unlock()
	// remove the underlying resource uri's from successfulSubscriptionList
	for i := 0; i < len(collectionList); i++ {
		for j := 0; j < len(successfulSubscriptionList); j++ {
			if collectionList[i] == successfulSubscriptionList[j].Oid {
				originResourceProcessedCount--
				successfulSubscriptionList = append(successfulSubscriptionList[:j], successfulSubscriptionList[j+1:]...)
				break
			}
		}
	}

	if len(successfulSubscriptionList) > 0 {
		subscriptionID := uuid.New().String()
		var hosts []string
		resp, hosts = result.ReadResponse(subscriptionID)
		if len(postRequest.OriginResources) == 0 {
			successfulSubscriptionList = []model.Link{}
			hosts = []string{}
		}
		statusCode, statusMessage, messageArgs, err = e.SaveSubscription(ctx, sessionUserName, subscriptionID,
			hosts, successfulSubscriptionList, postRequest)
		if err != nil {
			l.LogWithFields(ctx).Error(err.Error())
			evcommon.GenErrorResponse(err.Error(), statusMessage, statusCode,
				messageArgs, &resp)
			e.UpdateTask(ctx, fillTaskData(taskID, targetURI, string(req.PostBody),
				resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
			return resp
		}
		locationHeader = resp.Header["Location"]
	}
	l.LogWithFields(ctx).Debug("Process Count,", originResourceProcessedCount,
		" successOriginResourceCount ", len(successfulSubscriptionList))
	percentComplete = 100
	if originResourceProcessedCount == len(successfulSubscriptionList) {
		e.UpdateTask(ctx, fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Completed, common.OK, percentComplete, http.MethodPost))
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
		e.UpdateTask(ctx, fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
	}
	return resp
}

// SaveSubscription function save subscription in db
func (e *ExternalInterfaces) SaveSubscription(ctx context.Context, sessionUserName, subscriptionID string,
	hosts []string, successfulSubscriptionList []model.Link, postRequest model.EventDestination) (int32, string, []interface{}, error) {
	evtSubscription := evmodel.SubscriptionResource{
		UserName:       sessionUserName,
		SubscriptionID: subscriptionID,
		EventDestination: &model.EventDestination{
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
			DeliveryRetryPolicy:  postRequest.DeliveryRetryPolicy,
		},
		Hosts: hosts,
	}

	if err := e.SaveEventSubscription(evtSubscription); err != nil {
		return http.StatusInternalServerError, errResponse.InternalError, []interface{}{}, err
	}
	return http.StatusOK, common.OK, []interface{}{}, nil
}

// eventSubscription method update subscription on device
func (e *ExternalInterfaces) eventSubscription(ctx context.Context, postRequest model.EventDestination, origin, collectionName string, collectionFlag bool) (string, evresponse.EventResponse) {
	var resp evresponse.EventResponse
	var err error
	var plugin *common.Plugin
	var contactRequest evcommon.PluginContactRequest
	var target *common.Target
	if !collectionFlag {
		if strings.Contains(origin, "Fabrics") {
			return e.createFabricSubscription(ctx, postRequest, origin, collectionName, collectionFlag)
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
			l.LogWithFields(ctx).Error(errorMessage)
			return "", resp
		}

		contactRequest.Plugin = plugin
		if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
			token := e.getPluginToken(ctx, plugin)
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

	subscriptionPost := model.EventDestination{
		Name:                 postRequest.Name,
		Destination:          postRequest.Destination,
		EventTypes:           postRequest.EventTypes,
		MessageIds:           postRequest.MessageIds,
		ResourceTypes:        postRequest.ResourceTypes,
		Protocol:             postRequest.Protocol,
		SubscriptionType:     postRequest.SubscriptionType,
		EventFormatType:      postRequest.EventFormatType,
		SubordinateResources: postRequest.SubordinateResources,
		Context:              postRequest.Context,
		DeliveryRetryPolicy:  postRequest.DeliveryRetryPolicy,
	}
	res, err := e.IsEventsSubscribed(ctx, "", origin, &subscriptionPost, plugin, target, collectionFlag, collectionName, false, "", false)
	if err != nil {
		resp.Response = res.Body
		resp.StatusCode = int(res.StatusCode)
		return "", resp
	}
	if collectionFlag {
		l.LogWithFields(ctx).Info("Saving device subscription details of collection subscription")
		if collectionName == "AggregateCollections" {
			resp.StatusCode = http.StatusCreated
			resp.Response = createEventSubscriptionResponse()
			return collectionName, resp
		}
		err = e.saveDeviceSubscriptionDetails(common.DeviceSubscription{
			Location:       "",
			EventHostIP:    collectionName,
			OriginResource: origin,
		})
		if err != nil {
			errorMessage := "error while trying to save event subscription of device data: " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			l.LogWithFields(ctx).Error(errorMessage)
			return "", resp
		}
		resp.StatusCode = http.StatusCreated
		resp.Response = createEventSubscriptionResponse()
		return collectionName, resp
	}
	return e.SaveSubscriptionOnDevice(ctx, origin, target, plugin, contactRequest, subscriptionPost)
}

// SaveSubscriptionOnDevice method update subscription on device
func (e *ExternalInterfaces) SaveSubscriptionOnDevice(ctx context.Context, origin string, target *common.Target, plugin *common.Plugin, contactRequest evcommon.PluginContactRequest, subscriptionPost model.EventDestination) (string, evresponse.EventResponse) {
	var resp evresponse.EventResponse

	postBody, err := json.Marshal(subscriptionPost)
	if err != nil {
		errorMessage := "error while marshaling: " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		l.LogWithFields(ctx).Error(errorMessage)
		return "", resp

	}
	var reqData string
	//replacing the request url with south bound translation URL
	for key, value := range config.Data.URLTranslation.SouthBoundURL {
		reqData = strings.Replace(string(postBody), key, value, -1)
	}

	target.PostBody = []byte(reqData)
	contactRequest.URL = "/ODIM/v1/Subscriptions"
	contactRequest.HTTPMethodType = http.MethodPost
	contactRequest.PostBody = target

	l.LogWithFields(ctx).Debug("Subscription Request: " + reqData)
	response, err := e.callPlugin(context.TODO(), contactRequest)
	if err != nil {
		if evcommon.GetPluginStatus(ctx, plugin) {
			response, err = e.callPlugin(context.TODO(), contactRequest)
		}
		if err != nil {
			errorMessage := "error while contact plugin : " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			l.LogWithFields(ctx).Error(errorMessage)
			return "", resp
		}
	}
	defer response.Body.Close()
	l.LogWithFields(ctx).Debug("Subscription Response StatusCode: " + strconv.Itoa(int(response.StatusCode)))
	if response.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			errorMessage := "error while trying to read response body: " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			l.LogWithFields(ctx).Error(errorMessage)
			return "", resp
		}
		l.LogWithFields(ctx).Info("Subscription Response: " + string(body))
		var res interface{}
		err = json.Unmarshal(body, &res)
		if err != nil {
			errorMessage := "error while unmarshal the body : " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			l.LogWithFields(ctx).Error(errorMessage)
			return "", resp
		}

		errorMessage := "error while trying to create event subscription"
		resp.Response = res
		resp.StatusCode = response.StatusCode
		l.LogWithFields(ctx).Error(errorMessage)
		return "", resp
	}
	// if Subscription location is empty then don't store event details in DB
	locationHdr := response.Header.Get("location")
	if locationHdr == "" {
		errorMessage := "Subscription Location is missing in the response header"
		evcommon.GenEventErrorResponse(errorMessage, errors.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		l.LogWithFields(ctx).Error(errorMessage)
		return "", resp
	}
	// get the ip address from the host name
	deviceIPAddress, errorMessage := evcommon.GetIPFromHostName(target.ManagerAddress)
	if errorMessage != "" {
		evcommon.GenEventErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
			&resp, []interface{}{"ManagerAddress", target.ManagerAddress})
		l.LogWithFields(ctx).Error(errorMessage)
		return "", resp
	}
	l.LogWithFields(ctx).Debug("Saving device subscription details : ", deviceIPAddress)
	evtSubscription := common.DeviceSubscription{
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
		l.LogWithFields(ctx).Error(errorMessage)
		return "", resp
	}
	var outBody interface{}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errorMessage := "error while reading body  : " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		l.LogWithFields(ctx).Error(errorMessage)
		return "", resp
	}
	err = json.Unmarshal(body, &outBody)
	if err != nil {
		errorMessage := "error while unmarshal the body : " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		l.LogWithFields(ctx).Error(errorMessage)
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
func (e *ExternalInterfaces) IsEventsSubscribed(ctx context.Context, token, origin string, subscription *model.EventDestination, plugin *common.Plugin, target *common.Target, collectionFlag bool, collectionName string, isAggregate bool, aggregateID string, isRemove bool) (errResponse.RPC, error) {
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
			l.LogWithFields(ctx).Error(errorMessage)
			return resp, err
		}
		l.LogWithFields(ctx).Info("After look up, manager address is: ", host)
		searchKey = evcommon.GetSearchKey(host, evmodel.SubscriptionIndex)
	}

	var (
		eventTypes    = subscription.EventTypes
		messageIDs    = subscription.MessageIds
		resourceTypes = subscription.ResourceTypes
	)

	originResource = origin
	subscriptionDetails, err := e.GetEvtSubscriptions(searchKey)
	if err != nil && !strings.Contains(err.Error(), "No data found for the key") {
		errorMessage := "error while get subscription details: " + err.Error()
		evcommon.GenErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			[]interface{}{}, &resp)
		l.LogWithFields(ctx).Error(errorMessage)
		return resp, err
	}
	if isAggregate {
		subscriptionDetails = append(subscriptionDetails, e.GetAggregateSubscriptionList(ctx, host, aggregateID, isRemove)...)
	}

	// if there is no subscription happened then create event subscription
	if len(subscriptionDetails) < 1 {
		return resp, nil
	}

	var subscriptionPresent bool
	for index, evtSubscriptions := range subscriptionDetails {

		if isHostPresent(evtSubscriptions.Hosts, host) {
			subscriptionPresent = true

			if len(evtSubscriptions.EventDestination.EventTypes) > 0 && (index == 0 || len(eventTypes) > 0) {
				eventTypes = append(eventTypes, evtSubscriptions.EventDestination.EventTypes...)
			} else {
				eventTypes = []string{}
			}

			if len(evtSubscriptions.EventDestination.MessageIds) > 0 && (index == 0 || len(messageIDs) > 0) {
				messageIDs = append(messageIDs, evtSubscriptions.EventDestination.MessageIds...)
			} else {
				messageIDs = []string{}
			}

			if len(evtSubscriptions.EventDestination.ResourceTypes) > 0 && (index == 0 || len(resourceTypes) > 0) {
				resourceTypes = append(resourceTypes, evtSubscriptions.EventDestination.ResourceTypes...)
			} else {
				resourceTypes = []string{}
			}

		}
	}
	if !subscriptionPresent {
		return resp, nil
	}
	if !collectionFlag {
		l.LogWithFields(ctx).Debug("Delete Subscription from device")
		if strings.Contains(originResource, "Fabrics") {
			resp, err = e.DeleteFabricsSubscription(ctx, originResource, plugin)
			if err != nil {
				return resp, err
			}
		} else {
			resp, err = e.DeleteSubscriptions(ctx, originResource, token, plugin, target)
			if err != nil {
				return resp, err
			}
		}
	}
	// updating the subscription information
	removeDuplicatesFromSlice(&eventTypes)
	removeDuplicatesFromSlice(&messageIDs)
	removeDuplicatesFromSlice(&resourceTypes)
	subscription.EventTypes = eventTypes
	subscription.MessageIds = messageIDs
	subscription.ResourceTypes = resourceTypes
	return resp, nil
}

// CreateDefaultEventSubscription is creates the  subscription with event
// types which will be required to rediscover the inventory after computer
// system restarts ,This will  triggered from   aggregation service whenever
// a computer system is added
func (e *ExternalInterfaces) CreateDefaultEventSubscription(ctx context.Context, originResources, eventTypes, messageIDs, resourceTypes []string, protocol string) errResponse.RPC {
	l.LogWithFields(ctx).Info("Creation of default subscriptions started for: " + strings.Join(originResources, "::"))
	var resp errResponse.RPC
	var response evresponse.EventResponse
	var partialResultFlag bool
	if protocol == "" {
		protocol = "Redfish"
	}
	bubbleUpStatusCode := http.StatusCreated
	var postRequest model.EventDestination
	postRequest.Destination = ""
	postRequest.EventTypes = eventTypes
	postRequest.MessageIds = messageIDs
	postRequest.ResourceTypes = resourceTypes
	postRequest.Context = "Creating the Default Event Subscription"
	postRequest.Protocol = protocol
	postRequest.SubscriptionType = evmodel.SubscriptionType
	postRequest.SubordinateResources = true
	_, response = e.eventSubscription(ctx, postRequest, originResources[0], "", false)
	e.checkCollectionSubscription(ctx, originResources[0], protocol)
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

	resp.Body = response.Response
	resp.StatusCode = http.StatusCreated
	l.LogWithFields(ctx).Info("Creation of default subscriptions completed for : " + strings.Join(originResources, "::"))
	return resp
}

// saveDeviceSubscriptionDetails will first check if already origin resource details present
// if its present then Update location
// otherwise add an entry to redis
func (e *ExternalInterfaces) saveDeviceSubscriptionDetails(evtSubscription common.DeviceSubscription) error {
	searchKey := evcommon.GetSearchKey(evtSubscription.EventHostIP, evmodel.DeviceSubscriptionIndex)
	deviceSubscription, _ := e.GetDeviceSubscriptions(searchKey)
	var newDevSubscription = common.DeviceSubscription{
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

// getTargetDetails return device credentials from using device UUID
func (e *ExternalInterfaces) getTargetDetails(origin string) (*common.Target, evresponse.EventResponse, error) {
	var resp evresponse.EventResponse
	uuid, err := getUUID(origin)
	if err != nil {
		evcommon.GenEventErrorResponse(err.Error(), errResponse.ResourceNotFound, http.StatusNotFound,
			&resp, []interface{}{"System", origin})
		return nil, resp, err
	}

	// Get target device Credentials from using device UUID
	target, err := e.GetTarget(uuid)
	if err != nil {
		// Frame the RPC response body and response Header below
		errorMessage := "error while getting Systems(Target device Credentials) table details: " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
			&resp, []interface{}{"Systems", origin})
		return nil, resp, err
	}
	decryptedPasswordByte, err := DecryptWithPrivateKeyFunc(target.Password)
	if err != nil {
		// Frame the RPC response body and response Header below
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		return nil, resp, err
	}
	target.Password = decryptedPasswordByte
	return target, resp, nil
}

// DeleteSubscriptions will delete subscription from device
func (e *ExternalInterfaces) DeleteSubscriptions(ctx context.Context, originResource, token string, plugin *common.Plugin, target *common.Target) (errResponse.RPC, error) {
	var resp errResponse.RPC
	var err error
	var deviceSubscription *common.DeviceSubscription

	addr, errorMessage := evcommon.GetIPFromHostName(target.ManagerAddress)
	if errorMessage != "" {
		evcommon.GenErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
			[]interface{}{"ManagerAddress", target.ManagerAddress}, &resp)
		l.LogWithFields(ctx).Error(errorMessage)
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

		l.LogWithFields(ctx).Error(errorMessage)
		return resp, err
	}

	var contactRequest evcommon.PluginContactRequest

	contactRequest.Plugin = plugin
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token := e.getPluginToken(ctx, plugin)
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

	resp, _, _, err = e.PluginCall(ctx, contactRequest)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (e *ExternalInterfaces) createEventSubscription(ctx context.Context, taskID string, subTaskChan chan<- int32, reqSessionToken string,
	targetURI string, request model.EventDestination, originResource string, result *evresponse.MutexLock,
	wg *sync.WaitGroup, collectionFlag bool, collectionName string, aggregateResource string, isAggregateCollection bool) {
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
		l.LogWithFields(ctx).Error("error while trying to marshal create event request: " + err.Error())
	}
	reqJSON = string(reqBody)
	if taskID != "" {
		subTaskURI, err = e.CreateChildTask(ctx, reqSessionToken, taskID)
		if err != nil {
			l.LogWithFields(ctx).Error("Error while creating the SubTask")
		}
		trimmedURI := strings.TrimSuffix(subTaskURI, "/")
		subTaskID = trimmedURI[strings.LastIndex(trimmedURI, "/")+1:]
		resp.StatusCode = http.StatusAccepted
		e.UpdateTask(ctx, fillTaskData(subTaskID, targetURI, reqJSON, resp, common.Running, common.OK, percentComplete, http.MethodPost))
	}

	host, response := e.eventSubscription(ctx, request, originResource, collectionName, collectionFlag)
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
			e.UpdateTask(ctx, fillTaskData(subTaskID, targetURI, reqJSON, resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
		} else {
			e.UpdateTask(ctx, fillTaskData(subTaskID, targetURI, reqJSON, resp, common.Completed, common.OK, percentComplete, http.MethodPost))
		}
		subTaskChan <- int32(response.StatusCode)
	}
}

// checkCollectionSubscription checks if any collection based subscription exists
// If its' exists it will  update the existing subscription information with newly added server origin
func (e *ExternalInterfaces) checkCollectionSubscription(ctx context.Context, origin, protocol string) {
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
	var chassisSubscriptions, managersSubscriptions []evmodel.SubscriptionResource
	if bmcFlag {
		chassisSubscriptions, _ = e.GetEvtSubscriptions("/redfish/v1/Chassis")
		subscriptions = append(subscriptions, chassisSubscriptions...)
		managersSubscriptions, _ = e.GetEvtSubscriptions("/redfish/v1/Managers")
		subscriptions = append(subscriptions, managersSubscriptions...)
	}
	// Checking the collection based subscription
	var collectionSubscription = make([]evmodel.SubscriptionResource, 0)
	for _, evtSubscription := range subscriptions {
		for _, originResource := range evtSubscription.EventDestination.OriginResources {
			if strings.Contains(origin, "Systems") && (originResource.Oid == "/redfish/v1/Systems" || originResource.Oid == "/redfish/v1/Chassis" || originResource.Oid == "/redfish/v1/Managers") {
				collectionSubscription = append(collectionSubscription, evtSubscription)
			} else if strings.Contains(origin, "Fabrics") && originResource.Oid == "/redfish/v1/Fabrics" {
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
		destination = evtSubscription.EventDestination.Destination
		if len(evtSubscription.EventDestination.EventTypes) > 0 && (index == 0 || len(eventTypes) > 0) {
			eventTypes = append(eventTypes, evtSubscription.EventDestination.EventTypes...)
		} else {
			eventTypes = []string{}
		}

		if len(evtSubscription.EventDestination.MessageIds) > 0 && (index == 0 || len(messageIDs) > 0) {
			messageIDs = append(messageIDs, evtSubscription.EventDestination.MessageIds...)
		} else {
			messageIDs = []string{}
		}

		if len(evtSubscription.EventDestination.ResourceTypes) > 0 && (index == 0 || len(resourceTypes) > 0) {
			resourceTypes = append(resourceTypes, evtSubscription.EventDestination.ResourceTypes...)
		} else {
			resourceTypes = []string{}
		}
	}

	removeDuplicatesFromSlice(&eventTypes)
	removeDuplicatesFromSlice(&messageIDs)
	removeDuplicatesFromSlice(&resourceTypes)

	subordinateFlag := false
	if strings.Contains(origin, "Fabrics") {
		subordinateFlag = true
	}

	subscriptionPost := model.EventDestination{
		EventTypes:           eventTypes,
		MessageIds:           messageIDs,
		ResourceTypes:        resourceTypes,
		Context:              context,
		Destination:          destination,
		Protocol:             protocol,
		SubordinateResources: subordinateFlag,
	}
	subscriptionPost.OriginResources = []model.Link{
		{
			Oid: origin,
		},
	}

	// Subscribing newly added server with collated event list
	host, response := e.eventSubscription(ctx, subscriptionPost, origin, "", false)
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
		var newDevSubscription = common.DeviceSubscription{
			EventHostIP:     deviceSubscription.EventHostIP,
			Location:        deviceSubscription.Location,
			OriginResources: deviceSubscription.OriginResources,
		}
		newDevSubscription.OriginResources = append(newDevSubscription.OriginResources, chassisList...)
		newDevSubscription.OriginResources = append(newDevSubscription.OriginResources, managersList...)

		err := e.UpdateDeviceSubscriptionLocation(newDevSubscription)
		if err != nil {
			l.LogWithFields(ctx).Error("error while updating device subscription : " + err.Error())
		}
	}
}

func (e *ExternalInterfaces) createFabricSubscription(ctx context.Context, postRequest model.EventDestination, origin, collectionName string, collectionFlag bool) (string, evresponse.EventResponse) {
	var resp evresponse.EventResponse
	var err error
	var plugin *common.Plugin
	var contactRequest evcommon.PluginContactRequest
	// Extract the fabric id from the Origin
	fabricID := getFabricID(origin)
	fabric, dberr := e.GetFabricData(fabricID)
	if dberr != nil {
		errorMessage := "error while getting fabric data: " + dberr.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
			&resp, []interface{}{"Fabrics", fabricID})
		l.LogWithFields(ctx).Error(errorMessage)
		return "", resp
	}
	var gerr *errors.Error
	plugin, gerr = e.GetPluginData(fabric.PluginID)
	if gerr != nil {
		errorMessage := "error while getting plugin data: " + gerr.Error() + fabric.PluginID
		evcommon.GenEventErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
			&resp, []interface{}{"Plugin", fabric.PluginID})
		l.LogWithFields(ctx).Error(errorMessage)
		return "", resp
	}
	contactRequest.Plugin = plugin
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token := e.getPluginToken(ctx, plugin)
		if token == "" {
			evcommon.GenEventErrorResponse("error: Unable to create session with plugin "+plugin.ID, errResponse.NoValidSession, http.StatusUnauthorized,
				&resp, []interface{}{})
			l.LogWithFields(ctx).Error("error: Unable to create session with plugin " + plugin.ID)
			return "", resp
		}
		contactRequest.Token = token
	} else {
		contactRequest.LoginCredential = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
	}

	subscriptionPost := model.EventDestination{
		Name:                 postRequest.Name,
		Destination:          postRequest.Destination,
		EventTypes:           postRequest.EventTypes,
		MessageIds:           postRequest.MessageIds,
		ResourceTypes:        postRequest.ResourceTypes,
		Protocol:             postRequest.Protocol,
		SubscriptionType:     postRequest.SubscriptionType,
		EventFormatType:      postRequest.EventFormatType,
		SubordinateResources: postRequest.SubordinateResources,

		Context: postRequest.Context,
		OriginResources: []model.Link{
			{
				Oid: origin,
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
		l.LogWithFields(ctx).Error(errorMessage)
		return "", resp
	}
	var target = common.Target{
		ManagerAddress: deviceIPAddress,
	}
	res, err := e.IsEventsSubscribed(ctx, "", origin, &subscriptionPost, plugin, &target, collectionFlag, collectionName, false, "", false)
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
	if err != nil {
		errorMessage := "error while unmarshal the body : " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		l.LogWithFields(ctx).Error(errorMessage)
		return "", resp
	}

	response, err := e.callPlugin(context.TODO(), contactRequest)
	if err != nil {
		if evcommon.GetPluginStatus(ctx, plugin) {
			response, err = e.callPlugin(context.TODO(), contactRequest)
		}
		if err != nil {
			evcommon.GenEventErrorResponse(err.Error(), errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			l.LogWithFields(ctx).Error(err.Error())
			return "", resp
		}
	}
	defer response.Body.Close()
	//retrying the operation if status code is 401
	if response.StatusCode == http.StatusUnauthorized && strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		response, resp, err = e.retryEventSubscriptionOperation(ctx, contactRequest)
		if err != nil {
			return "", resp
		}
	}
	l.LogWithFields(ctx).Debug("Subscription Response Status Code: " + string(rune(response.StatusCode)))
	if response.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			errorMessage := "error while trying to read response body: " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			l.LogWithFields(ctx).Error(errorMessage)
			return "", resp
		}
		errorMessage := "error while trying to create event subscription"
		var res interface{}
		l.LogWithFields(ctx).Error("Subscription Response " + string(body))
		err = json.Unmarshal(body, &res)
		if err != nil {
			errorMessage := "error while unmarshal the body : " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			l.LogWithFields(ctx).Error(errorMessage)
			return "", resp
		}

		resp.Response = res
		resp.StatusCode = response.StatusCode
		l.LogWithFields(ctx).Error(errorMessage)
		return "", resp
	}

	evtSubscription := common.DeviceSubscription{
		EventHostIP:    deviceIPAddress,
		OriginResource: origin,
	}
	evtSubscription.Location = response.Header.Get("location")
	err = e.saveDeviceSubscriptionDetails(evtSubscription)
	if err != nil {
		errorMessage := "error while trying to save event subscription of device data: " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		l.LogWithFields(ctx).Error(errorMessage)
		return "", resp
	}

	resp.Response = createEventSubscriptionResponse()
	resp.StatusCode = response.StatusCode
	resp.Location = response.Header.Get("location")
	return deviceIPAddress, resp
}

// UpdateEventSubscriptions it will add subscription for newly Added system in aggregate
func (e *ExternalInterfaces) UpdateEventSubscriptions(ctx context.Context, req *eventsproto.EventUpdateRequest, isRemove bool) (string, evresponse.EventResponse) {
	var resp evresponse.EventResponse
	var plugin *common.Plugin
	var contactRequest evcommon.PluginContactRequest
	var target *common.Target

	target, resp, err := e.getTargetDetails(req.SystemID)
	if err != nil {
		return "", resp
	}
	var errs *errors.Error
	plugin, errs = e.GetPluginData(target.PluginID)
	if errs != nil {
		errorMessage := "error while getting plugin data: " + errs.Error()
		l.LogWithFields(ctx).Info(errorMessage)
		return "", resp
	}

	contactRequest.Plugin = plugin
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token := e.getPluginToken(ctx, plugin)
		if token == "" {
			l.LogWithFields(ctx).Info("error: Unable to create session with plugin " + plugin.ID)
			return "", resp
		}
		contactRequest.Token = token

	} else {
		contactRequest.LoginCredential = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
	}
	subscriptionPost := model.EventDestination{
		EventTypes:    []string{},
		MessageIds:    []string{},
		ResourceTypes: []string{},
		OriginResources: []model.Link{
			{
				Oid: req.SystemID,
			},
		},
		SubordinateResources: true,
		Protocol:             "Redfish",
		SubscriptionType:     evmodel.SubscriptionType,
		Context:              evmodel.Context,
		DeliveryRetryPolicy:  "RetryForever",
		EventFormatType:      "Event",
	}
	res, err := e.IsEventsSubscribed(ctx, "", req.SystemID, &subscriptionPost, plugin, target, false, "", true, req.AggregateId, isRemove)
	if err != nil {
		resp.Response = res.Body
		resp.StatusCode = int(res.StatusCode)
		return "", resp
	}
	return e.SaveSubscriptionOnDevice(ctx, req.SystemID, target, plugin, contactRequest, subscriptionPost)
}

// GetAggregateSubscriptionList return list of subscription corresponding to host
func (e *ExternalInterfaces) GetAggregateSubscriptionList(ctx context.Context, host, aggregateID string, isRemove bool) (data []evmodel.SubscriptionResource) {
	searchKeyAgg := evcommon.GetSearchKey(host, evmodel.SubscriptionIndex)
	aggregateList, err := e.GetAggregateList(searchKeyAgg)
	if err != nil {
		l.LogWithFields(ctx).Info("No Aggregate subscription Found ", err)
	}
	for _, id := range aggregateList {
		if isRemove {
			if id == aggregateID {
				continue
			}
		}
		searchKey := evcommon.GetSearchKey(id, evmodel.SubscriptionIndex)
		aggregateSubscriptionDetails, err := e.GetEvtSubscriptions(searchKey)

		if err != nil && !strings.Contains(err.Error(), "No data found for the key") {
			l.LogWithFields(ctx).Info("Error while get aggregateSubscriptionDetails details: " + err.Error())
			continue
		}
		data = append(data, aggregateSubscriptionDetails...)
	}
	return
}

// getSuccessfulResponse return successful subscription list
func getSuccessfulResponse(response map[string]evresponse.EventResponse) (successfulSubscriptionList []model.Link, successfulResponses map[string]evresponse.EventResponse) {
	var resourceID string
	successfulResponses = make(map[string]evresponse.EventResponse)
	i := 0
	for originResource, evtResponse := range response {
		OriginResource := strings.SplitAfter(originResource, "/")
		originResourceID := OriginResource[len(OriginResource)-1]
		if i == 0 {
			resourceID = originResourceID
		}
		if originResourceID == resourceID && i > 0 {
			successfulSubscriptionList = append(successfulSubscriptionList, model.Link{Oid: originResource})
		}
		if evtResponse.StatusCode == http.StatusCreated {
			successfulSubscriptionList = append(successfulSubscriptionList, model.Link{Oid: originResource})
			successfulResponses[originResource] = evtResponse
		}
		i++
	}
	return
}
