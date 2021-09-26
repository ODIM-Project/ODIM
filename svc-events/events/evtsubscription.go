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
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	errResponse "github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/ODIM-Project/ODIM/svc-events/evresponse"
	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
)

//PluginContact struct to inject the pmb client function into the handlers
type PluginContact struct {
	ContactClient      func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	Auth               func(string, []string, []string) response.RPC
	UpdateTask         func(common.TaskData) error
	CreateChildTask    func(string, string) (string, error)
	GetSessionUserName func(sessionToken string) (string, error)
}

func fillTaskData(taskID, targetURI, request string, resp errResponse.RPC, taskState string, taskStatus string, percentComplete int32, httpMethod string) common.TaskData {
	return common.TaskData{
		TaskID:          taskID,
		TargetURI:       targetURI,
		Response:        resp,
		TaskRequest:     request,
		TaskState:       taskState,
		TaskStatus:      taskStatus,
		PercentComplete: percentComplete,
		HTTPMethod:      httpMethod,
	}
}

// UpdateTaskData update the task with the given data
func UpdateTaskData(taskData common.TaskData) error {
	respBody, _ := json.Marshal(taskData.Response.Body)
	payLoad := &taskproto.Payload{
		HTTPHeaders:   taskData.Response.Header,
		HTTPOperation: taskData.HTTPMethod,
		JSONBody:      taskData.TaskRequest,
		StatusCode:    taskData.Response.StatusCode,
		TargetURI:     taskData.TargetURI,
		ResponseBody:  respBody,
	}

	err := services.UpdateTask(taskData.TaskID, taskData.TaskState, taskData.TaskStatus, taskData.PercentComplete, payLoad, time.Now())
	if err != nil && (err.Error() == common.Cancelling) {
		// We cant do anything here as the task has done it work completely, we cant reverse it.
		//Unless if we can do opposite/reverse action for delete server which is add server.
		services.UpdateTask(taskData.TaskID, common.Cancelled, taskData.TaskStatus, taskData.PercentComplete, payLoad, time.Now())
		if taskData.PercentComplete == 0 {
			return fmt.Errorf("error while starting the task: %v", err)
		}
		log.Error("error: task update for " + taskData.TaskID + " failed with err: " + err.Error())
		runtime.Goexit()
	}
	return nil
}

func removeOdataIDfromOriginResources(originResources []evmodel.OdataIDLink) []string {
	var originRes []string
	for _, origin := range originResources {
		originRes = append(originRes, origin.OdataID)
	}
	return originRes
}

// CreateEventSubscription is a API to create event subscription
func (p *PluginContact) CreateEventSubscription(taskID string, sessionUserName string, req *eventsproto.EventSubRequest) errResponse.RPC {
	var (
		err             error
		resp            errResponse.RPC
		postRequest     evmodel.RequestBody
		percentComplete int32 = 100
		targetURI             = "/redfish/v1/EventService/Subscriptions"
	)

	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}

	if err = json.Unmarshal(req.PostBody, &postRequest); err != nil {
		// Update the task here with error response
		errorMessage := "Error while Unmarshaling the Request: " + err.Error()
		if strings.Contains(err.Error(), "evmodel.OdataIDLink") {
			errorMessage = "Error processing subscription request: @odata.id key(s) is missing in origin resources list"
		}
		log.Error(errorMessage)

		resp = common.GeneralError(http.StatusBadRequest, errResponse.MalformedJSON, errorMessage, []interface{}{}, nil)
		// Fill task and update
		p.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
		return resp
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.PostBody, postRequest)
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Error(errorMessage)
		resp := common.GeneralError(http.StatusBadRequest, errResponse.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, nil)
		p.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
		return resp
	}

	//check mandatory fields
	statuscode, statusMessage, messageArgs, err := validateFields(&postRequest)
	if err != nil {
		// Update the task here with error response
		errorMessage := "error: request payload validation failed: " + err.Error()
		log.Error(errorMessage)

		resp = common.GeneralError(statuscode, statusMessage, errorMessage, messageArgs, nil)
		// Fill task and update
		p.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
		return resp
	}

	//validate destination URI in the request
	if !common.URIValidator(postRequest.Destination) {
		errorMessage := "error: request body contains invalid value for Destination field, " + postRequest.Destination
		log.Error(errorMessage)

		resp = common.GeneralError(http.StatusBadRequest, errResponse.PropertyValueFormatError, errorMessage, []interface{}{postRequest.Destination, "Destination"}, nil)
		// Fill task and update
		p.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
		return resp
	}

	// check any of the subscription present for the destination from the request
	// if errored out or no subscriptions then add subscriptions else return an error
	subscriptionDetails, err := evmodel.GetEvtSubscriptions(postRequest.Destination)
	if err != nil && !strings.Contains(err.Error(), "No data found for the key") {
		errorMessage := "Error while get subscription details: " + err.Error()
		evcommon.GenErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			[]interface{}{}, &resp)
		log.Error(errorMessage)
		p.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
		return resp
	}
	for _, evtSubscription := range subscriptionDetails {
		if evtSubscription.Destination == postRequest.Destination {
			errorMessage := "Subscription already present for the requested destination"
			evcommon.GenErrorResponse(errorMessage, errResponse.ResourceInUse, http.StatusConflict,
				[]interface{}{}, &resp)
			log.Error(errorMessage)
			p.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
			return resp
		}
	}

	// Get the target device  details from the origin resources
	// Loop through all origin list and form individual event subscription request,
	// Which will then forward to plugin to make subscrption with target device
	var wg, taskCollectionWG sync.WaitGroup
	var result = &evresponse.MutexLock{
		Response: make(map[string]evresponse.EventResponse),
		Lock:     &sync.Mutex{},
	}

	// remove odataid in the originresources
	originResources := removeOdataIDfromOriginResources(postRequest.OriginResources)
	originResourcesCount := len(originResources)

	// check and remove if duplicate OriginResources exist in the request
	removeDuplicatesFromSlice(&originResources, &originResourcesCount)

	// If origin resource is nil then subscribe to all collection
	if originResourcesCount <= 0 {
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
		// Collect the channels and update perentComplete in Task
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
				p.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Running, common.OK, percentComplete, http.MethodPost))
			}
		}
	}()

	for _, origin := range originResources {
		_, _, err := getTargetDetails(origin)
		if err != nil {
			collection, collectionName, collectionFlag, _ := checkCollection(origin)
			wg.Add(1)
			// for origin is collection
			go p.createEventSubscrption(taskID, subTaskChan, sessionUserName, targetURI, postRequest, origin, result, &wg, collectionFlag, collectionName)
			for i := 0; i < len(collection); i++ {
				wg.Add(1)
				// for suboridinate origin
				go p.createEventSubscrption("", subTaskChan, sessionUserName, targetURI, postRequest, collection[i], result, &wg, false, "")
			}
			collectionList = append(collectionList, collection...)
		} else {
			wg.Add(1)
			go p.createEventSubscrption(taskID, subTaskChan, sessionUserName, targetURI, postRequest, origin, result, &wg, false, "")
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
	i := 0
	var resourceID string

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
	successOriginResourceCount := len(successfulSubscriptionList)
	result.Lock.Unlock()
	// remove the underlaying resource uri's from successfulSubscriptionList
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
	successOriginResourceCount = len(successfulSubscriptionList)
	if successOriginResourceCount > 0 {
		subscriptionID := uuid.New().String()
		var hosts []string
		resp, hosts = result.ReadResponse(subscriptionID)
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
		}

		if err = evmodel.SaveEventSubscription(evtSubscription); err != nil {
			// Update the task here with error response
			errorMessage := "error while trying to save event subscription data: " + err.Error()
			log.Error(errorMessage)

			resp = common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, []interface{}{}, nil)
			// Fill task and update
			percentComplete = 100
			p.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
			return resp
		}
		locationHeader = resp.Header["Location"]
	}
	log.Info("Process Count," + strconv.Itoa(originResourceProcessedCount) +
		"successOriginResourceCount" + strconv.Itoa(successOriginResourceCount))
	percentComplete = 100
	if originResourceProcessedCount == successOriginResourceCount {
		p.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Completed, common.OK, percentComplete, http.MethodPost))
	} else {
		args := response.Args{
			Code:    response.GeneralError,
			Message: "event subscription for one or more origin resource(s) failed, check sub tasks for more info.",
		}
		resp.Body = args.CreateGenericErrorResponse()
		resp.StatusCode = bubbleUpStatusCode
		if locationHeader != "" {
			resp.Header["Location"] = locationHeader
		}
		p.UpdateTask(fillTaskData(taskID, targetURI, string(req.PostBody), resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
	}
	return resp
}

// remove duplicate elements in string slice.
// Takes string slice and length, and updates the same with new values
func removeDuplicatesFromSlice(slc *[]string, slcLen *int) {
	if *slcLen > 1 {
		uniqueElementsDs := make(map[string]bool)
		var uniqueElemenstsList []string
		for _, element := range *slc {
			if exist := uniqueElementsDs[element]; !exist {
				uniqueElemenstsList = append(uniqueElemenstsList, element)
				uniqueElementsDs[element] = true
			}
		}
		// length of uniqueElemenstsList will be less than passed string slice,
		// only if duplicates existed, so will assign slc with modified list and update length
		if len(uniqueElemenstsList) < *slcLen {
			*slc = uniqueElemenstsList
			*slcLen = len(*slc)
		}
	}
	return
}

func (p *PluginContact) eventSubscription(postRequest evmodel.RequestBody, origin, collectionName string, collectionFlag bool) (string, evresponse.EventResponse) {
	var resp evresponse.EventResponse
	var err error
	var plugin *evmodel.Plugin
	var contactRequest evcommon.PluginContactRequest
	var target *evmodel.Target
	if !collectionFlag {
		if strings.Contains(origin, "Fabrics") {
			return p.createFabricSubscription(postRequest, origin, collectionName, collectionFlag)
		}
		target, resp, err = getTargetDetails(origin)
		if err != nil {
			return "", resp
		}

		var errs *errors.Error
		plugin, errs = evmodel.GetPluginData(target.PluginID)
		if errs != nil {
			errorMessage := "error while getting plugin data: " + errs.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			log.Error(errorMessage)
			return "", resp
		}

		contactRequest.Plugin = plugin
		if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
			token := p.getPluginToken(plugin)
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
	}
	res, err := p.IsEventsSubscribed("", origin, &subscriptionPost, plugin, target, collectionFlag, collectionName)
	if err != nil {
		resp.Response = res.Body
		resp.StatusCode = int(res.StatusCode)
		return "", resp
	}
	if collectionFlag {
		log.Info("Saving device subscription details of collection subscription")
		err = saveDeviceSubscriptionDetails(evmodel.Subscription{
			Location:       "",
			EventHostIP:    collectionName,
			OriginResource: origin,
		})
		if err != nil {
			errorMessage := "error while trying to save event subscription of device data: " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			log.Error(errorMessage)
			return "", resp
		}

		resp.StatusCode = http.StatusCreated
		resp.Response = createEventSubscriptionResponse()
		return collectionName, resp
	}

	postBody, _ := json.Marshal(subscriptionPost)
	var reqData string
	//replacing the reruest url with south bound translation URL
	for key, value := range config.Data.URLTranslation.SouthBoundURL {
		reqData = strings.Replace(string(postBody), key, value, -1)
	}

	target.PostBody = []byte(reqData)
	contactRequest.URL = "/ODIM/v1/Subscriptions"
	contactRequest.HTTPMethodType = http.MethodPost
	contactRequest.PostBody = target

	log.Info("Subscription Request" + reqData)
	response, err := p.callPlugin(contactRequest)
	if err != nil {
		if evcommon.GetPluginStatus(plugin) {
			response, err = p.callPlugin(contactRequest)
		}
		if err != nil {
			errorMessage := "error while unmarshaling the body : " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			log.Error(errorMessage)
			return "", resp
		}
	}
	defer response.Body.Close()
	log.Info("Subscription Response StatusCode:" + strconv.Itoa(int(response.StatusCode)))
	if response.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			errorMessage := "error while trying to read response body: " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			log.Error(errorMessage)
			return "", resp
		}
		log.Info("Subscription Response:" + string(body))
		var res interface{}
		err = json.Unmarshal(body, &res)
		if err != nil {
			errorMessage := "error while unmarshaling the body : " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			log.Error(errorMessage)
			return "", resp
		}

		errorMessage := "error while trying to create event subscription"
		resp.Response = res
		resp.StatusCode = response.StatusCode
		log.Error(errorMessage)
		return "", resp
	}
	// if Subscription location is empty then don't store event details in DB
	locationHdr := response.Header.Get("location")
	if locationHdr == "" {
		errorMessage := "Subscription Location is missing in the response header"
		evcommon.GenEventErrorResponse(errorMessage, errors.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		log.Error(errorMessage)
		return "", resp
	}
	// get the ip address from the host name
	deviceIPAddress, errorMessage := evcommon.GetIPFromHostName(target.ManagerAddress)
	if errorMessage != "" {
		evcommon.GenEventErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
			&resp, []interface{}{"ManagerAddress", target.ManagerAddress})
		log.Error(errorMessage)
		return "", resp
	}
	log.Info("Saving device subscription details : ", deviceIPAddress)
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
	err = saveDeviceSubscriptionDetails(evtSubscription)
	if err != nil {
		errorMessage := "error while trying to save event subscription of device data: " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		log.Error(errorMessage)
		return "", resp
	}

	var outBody interface{}
	body, err := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &outBody)
	if err != nil {
		errorMessage := "error while unmarshaling the body : " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		log.Error(errorMessage)
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
func (p *PluginContact) IsEventsSubscribed(token, origin string, subscription *evmodel.EvtSubPost, plugin *evmodel.Plugin, target *evmodel.Target, collectionFlag bool, collectionName string) (errResponse.RPC, error) {
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
			log.Error(errorMessage)
			return resp, err
		}
		log.Info("After look up, manager address is: ", host)
		searchKey = evcommon.GetSearchKey(host, evmodel.SubscriptionIndex)
	}
	// uniqueMap is to ignore duplicate eventTypes
	// evevntTypes from request  and eventTypes from the all destinations stored in the DB
	uniqueMap := make(map[string]string)

	// add all events to map to remove duplicate eventTypes
	// this need to be remove after the desination uniquness check added
	for _, eventType := range subscription.EventTypes {
		uniqueMap[eventType] = eventType
	}
	var (
		eventTypes    = subscription.EventTypes
		messageIDs    = subscription.MessageIds
		resourceTypes = subscription.ResourceTypes
	)

	originResource = origin
	subscriptionDetails, err := evmodel.GetEvtSubscriptions(searchKey)
	if err != nil && !strings.Contains(err.Error(), "No data found for the key") {
		errorMessage := "Error while get subscription details: " + err.Error()
		evcommon.GenErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			[]interface{}{}, &resp)
		log.Error(errorMessage)
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

			dest := evtSubscriptions.Destination
			if dest == subscription.Destination {
				eventTypes := evtSubscriptions.EventTypes

				// check if user tries to subscribe for same events
				resp, err = checkEqual(subscription.EventTypes, eventTypes)
				if err != nil {
					return resp, err
				}
				//if EventTypes are not same then delete subscription

				// if there is only one host in Hosts entry then
				// delete the subscription from redis
				if len(evtSubscriptions.Hosts) == 1 {
					err = evmodel.DeleteEvtSubscription(evtSubscriptions.SubscriptionID)
					if err != nil {
						errorMessage := "Error while Updating event subscription : " + err.Error()
						evcommon.GenErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusBadRequest,
							[]interface{}{"Subscription", "invalid value " + origin}, &resp)
						log.Error(errorMessage)
						return resp, err
					}
				} else {
					// Delete the host and origin resource from the respective entry
					evtSubscriptions.Hosts = removeElement(evtSubscriptions.Hosts, host)
					evtSubscriptions.OriginResources = removeElement(evtSubscriptions.OriginResources, originResource)
					err = evmodel.UpdateEventSubscription(evtSubscriptions)
					if err != nil {
						errorMessage := "Error while Updating event subscription : " + err.Error()
						evcommon.GenErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusBadRequest,
							[]interface{}{"Subscription", "invalid value " + origin}, &resp)
						log.Error(errorMessage)
						return resp, err
					}
				}

			} else {
				// Ignore the event types from the same destination, since we are trying modify the subscription
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
	}
	if !subscriptionPresent {
		return resp, nil
	}
	if !collectionFlag {
		log.Info("Delete Subscription from device")
		if strings.Contains(originResource, "Fabrics") {
			resp, err = p.DeleteFabricsSubscription(originResource, plugin)
			if err != nil {
				return resp, err
			}
		} else {
			resp, err = p.DeleteSubscriptions(originResource, token, plugin, target)
			if err != nil {
				return resp, err
			}
		}
	}
	// updating the subscritpion information

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

// removeElement will remove the element from the slice return
// slice of remaining elements
func removeElement(slice []string, element string) []string {
	var elements []string
	for _, val := range slice {
		if val != element {
			elements = append(elements, val)
		}
	}
	return elements
}

// getTypes is to split the string to array
func getTypes(subscription string) []string {
	// array stored in db in string("[alert statuschange]")
	// to convert into an array removing "[" ,"]" and splitting
	events := strings.Replace(subscription, "[", "", -1)
	events = strings.Replace(events, "]", "", -1)
	if len(events) < 1 {
		return []string{}
	}
	return strings.Split(events, " ")
}

//checkequal is to check the previous and new event types are equal
func checkEqual(newEventTypes, prevEventTypes []string) (errResponse.RPC, error) {
	var resp errResponse.RPC
	// if the subscribed events are same as wants to subscribe then return as resource in use
	if reflect.DeepEqual(newEventTypes, prevEventTypes) {
		errorMessage := "Resource already in use"
		evcommon.GenErrorResponse(errorMessage, errResponse.ResourceInUse, http.StatusConflict,
			[]interface{}{}, &resp)
		return resp, fmt.Errorf(errorMessage)
	}
	return resp, nil
}

// PluginCall method is to call to given url and method
// and validate the response and return
func (p *PluginContact) PluginCall(req evcommon.PluginContactRequest) (errResponse.RPC, string, string, error) {
	var resp errResponse.RPC
	response, err := p.callPlugin(req)
	if err != nil {
		if evcommon.GetPluginStatus(req.Plugin) {
			response, err = p.callPlugin(req)
		}
		if err != nil {
			errorMessage := "Error : " + err.Error()
			evcommon.GenErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				[]interface{}{}, &resp)
			log.Error(errorMessage)
			return resp, "", "", err
		}
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errorMessage := "error while trying to read response body: " + err.Error()
		evcommon.GenErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			[]interface{}{}, &resp)
		log.Error(errorMessage)
		return resp, "", "", err
	}
	if !(response.StatusCode == http.StatusCreated || response.StatusCode == http.StatusOK) {
		resp.StatusCode = int32(response.StatusCode)
		resp.Body = string(body)
		return resp, "", "", err
	}
	var outBody interface{}
	json.Unmarshal(body, &outBody)
	resp.StatusCode = int32(response.StatusCode)
	resp.Body = outBody
	return resp, response.Header.Get("location"), response.Header.Get("X-Auth-Token"), nil
}

// CreateDefaultEventSubscription is creates the  subscription with event types which will be required to rediscover the inventory
// after computer system restarts ,This will  triggered from   aggregation service whenever a computer system is added
func (p *PluginContact) CreateDefaultEventSubscription(originResources, eventTypes, messageIDs, resourceTypes []string, protocol string) errResponse.RPC {
	log.Info("Creation of default subscriptions started for:" + strings.Join(originResources, "::"))
	var resp errResponse.RPC
	var response evresponse.EventResponse
	var partialResultFlag bool
	if protocol == "" {
		protocol = "Redfish"
	}
	var hosts []string
	bubbleUpStatusCode := http.StatusCreated
	for i := 0; i < len(originResources); i++ {
		var postRequest evmodel.RequestBody
		postRequest.Destination = ""
		postRequest.EventTypes = eventTypes
		postRequest.MessageIds = messageIDs
		postRequest.ResourceTypes = resourceTypes
		postRequest.Context = "Creating the Default Event Subscription"
		postRequest.Protocol = protocol
		postRequest.SubscriptionType = evmodel.SubscriptionType
		postRequest.SubordinateResources = true
		var host string
		host, response = p.eventSubscription(postRequest, originResources[i], "", false)
		hosts = append(hosts, host)
		go p.checkCollectionSubscription(originResources[i], protocol)
		if response.StatusCode != http.StatusCreated {
			partialResultFlag = true
			if response.StatusCode > bubbleUpStatusCode {
				bubbleUpStatusCode = response.StatusCode
			}
		}

	}
	if !partialResultFlag || len(originResources) == 1 {
		resp.StatusCode = int32(response.StatusCode)
	} else {
		resp.StatusCode = int32(bubbleUpStatusCode)
	}
	subscriptionID := uuid.New().String()
	evtSubscription := evmodel.Subscription{
		SubscriptionID:       subscriptionID,
		EventTypes:           eventTypes,
		MessageIds:           messageIDs,
		ResourceTypes:        resourceTypes,
		OriginResources:      originResources,
		Hosts:                hosts,
		Protocol:             protocol,
		SubscriptionType:     evmodel.SubscriptionType,
		SubordinateResources: true,
	}
	err := evmodel.SaveEventSubscription(evtSubscription)
	if err != nil {
		errorMessage := "error while trying to save event subscription data: " + err.Error()
		evcommon.GenErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			[]interface{}{}, &resp)
		log.Error(errorMessage)
		return resp
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	resp.Body = response.Response
	resp.StatusCode = http.StatusCreated
	log.Info("Creation of default subscriptions completed for :" + strings.Join(originResources, "::"))
	return resp
}

func validateFields(request *evmodel.RequestBody) (int32, string, []interface{}, error) {
	validEventFormatTypes := map[string]bool{"Event": true, "MetricReport": true}
	validEventTypes := map[string]bool{"Alert": true, "MetricReport": true, "ResourceAdded": true, "ResourceRemoved": true, "ResourceUpdated": true, "StatusChange": true, "Other": true}

	validate := validator.New()

	// if any of the mandatory fields missing in the struct, then it return an error
	err := validate.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return http.StatusBadRequest, errResponse.PropertyMissing, []interface{}{err.Field()}, fmt.Errorf(err.Field() + " field is missing")
		}
	}
	if request.EventFormatType == "" {
		request.EventFormatType = "Event"
	}

	if _, ok := validEventFormatTypes[request.EventFormatType]; !ok {
		return http.StatusBadRequest, errResponse.PropertyValueNotInList, []interface{}{request.EventFormatType, "EventFormatType"}, fmt.Errorf("Invalid EventFormatType")
	}

	if len(request.EventTypes) == 0 && request.EventFormatType == "MetricReport" {
		request.EventTypes = []string{"MetricReport"}
	}

	for _, eventType := range request.EventTypes {
		if _, ok := validEventTypes[eventType]; !ok {
			return http.StatusBadRequest, errResponse.PropertyValueNotInList, []interface{}{eventType, "EventTypes"}, fmt.Errorf("Invalid EventTypes")
		}
	}

	if request.EventFormatType == "MetricReport" {
		if len(request.EventTypes) > 1 {
			return http.StatusBadRequest, errResponse.PropertyValueFormatError, []interface{}{request.EventFormatType, "EventTypes"}, fmt.Errorf("Unsupported EventType")
		}
		if request.EventTypes[0] != "MetricReport" {
			return http.StatusBadRequest, errResponse.PropertyValueNotInList, []interface{}{request.EventTypes[0], "EventType"}, fmt.Errorf("Unsupported EventType")
		}
	}

	if request.SubscriptionType == "" {
		request.SubscriptionType = evmodel.SubscriptionType
	} else if request.SubscriptionType == "SSE" || request.SubscriptionType == "SNMPTrap" || request.SubscriptionType == "SNMPInform" {
		return http.StatusBadRequest, errResponse.PropertyMissing, []interface{}{"SubscriptionType"}, fmt.Errorf("Unsupported SubscriptionType")
	} else if request.SubscriptionType != evmodel.SubscriptionType {
		return http.StatusBadRequest, errResponse.PropertyMissing, []interface{}{"SubscriptionType"}, fmt.Errorf("Invalid SubscriptionType")
	}

	if request.Context == "" {
		request.Context = evmodel.Context
	}

	availableProtocols := []string{"Redfish"}
	var validProtocol bool
	validProtocol = false
	for _, protocol := range availableProtocols {
		if request.Protocol == protocol {
			validProtocol = true
		}
	}
	if !validProtocol {
		return http.StatusBadRequest, errResponse.PropertyValueNotInList, []interface{}{request.Protocol, "Protocol"}, fmt.Errorf("Protocol %v is invalid", request.Protocol)
	}

	// check the All ResourceTypes are supported
	for _, resourceType := range request.ResourceTypes {
		if _, ok := common.ResourceTypes[resourceType]; !ok {
			return http.StatusBadRequest, errResponse.PropertyValueNotInList, []interface{}{resourceType, "ResourceType"}, fmt.Errorf("Unsupported ResourceType")
		}
	}

	return http.StatusOK, common.OK, []interface{}{}, nil
}

// saveDeviceSubscriptionDetails will first check if already origin resource details present
// if its present then Update location
// otherwise add an entry to redis
func saveDeviceSubscriptionDetails(evtSubscription evmodel.Subscription) error {
	searchKey := evcommon.GetSearchKey(evtSubscription.EventHostIP, evmodel.DeviceSubscriptionIndex)
	deviceSubscription, _ := evmodel.GetDeviceSubscriptions(searchKey)

	var newDevSubscription = evmodel.DeviceSubscription{
		EventHostIP:     evtSubscription.EventHostIP,
		Location:        evtSubscription.Location,
		OriginResources: []string{evtSubscription.OriginResource},
	}
	// if device subscriptions details for the device is present in db then dont add again
	var save = true
	if deviceSubscription != nil {

		save = true
		// if the origin resource is present in device subscription details then dont add
		for _, originResource := range deviceSubscription.OriginResources {
			if evtSubscription.OriginResource == originResource {

				save = false
			} else {
				newDevSubscription.OriginResources = append(newDevSubscription.OriginResources, originResource)
				save = false
			}
		}
		err := evmodel.UpdateDeviceSubscriptionLocation(newDevSubscription)
		if err != nil {
			return err
		}
	}
	if save {
		return evmodel.SaveDeviceSubscription(newDevSubscription)
	}
	return nil
}

func getTargetDetails(origin string) (*evmodel.Target, evresponse.EventResponse, error) {
	var resp evresponse.EventResponse

	uuid, err := getUUID(origin)
	if err != nil {
		evcommon.GenEventErrorResponse(err.Error(), errResponse.ResourceNotFound, http.StatusNotFound,
			&resp, []interface{}{"System", origin})
		log.Error(err.Error())
		return nil, resp, err
	}

	// Get target device Credentials from using device UUID
	target, err := evmodel.GetTarget(uuid)
	if err != nil {
		// Frame the RPC response body and response Header below

		errorMessage := "error while getting Systems(Target device Credentials) table details: " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
			&resp, []interface{}{"Systems", origin})
		log.Error(errorMessage)
		return nil, resp, err
	}
	decryptedPasswordByte, err := common.DecryptWithPrivateKey(target.Password)
	if err != nil {
		// Frame the RPC response body and response Header below
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		log.Error(errorMessage)
		return nil, resp, err
	}
	target.Password = decryptedPasswordByte
	return target, resp, nil
}

// DeleteSubscriptions will delete subscription from device
func (p *PluginContact) DeleteSubscriptions(originResource, token string, plugin *evmodel.Plugin, target *evmodel.Target) (errResponse.RPC, error) {
	var resp errResponse.RPC
	var err error
	var deviceSubscription *evmodel.DeviceSubscription
	addr, errorMessage := evcommon.GetIPFromHostName(target.ManagerAddress)
	if errorMessage != "" {
		evcommon.GenErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
			[]interface{}{"ManagerAddress", target.ManagerAddress}, &resp)
		log.Error(errorMessage)
		return resp, err
	}
	searchKey := evcommon.GetSearchKey(addr, evmodel.DeviceSubscriptionIndex)
	deviceSubscription, err = evmodel.GetDeviceSubscriptions(searchKey)
	if err != nil {
		// if its first subscription then no need to check events subscribed
		if strings.Contains(err.Error(), "No data found for the key") {
			return resp, nil
		}
		errorMessage := "Error while get subscription details: " + err.Error()
		evcommon.GenErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			[]interface{}{}, &resp)

		log.Error(errorMessage)
		return resp, err
	}

	var contactRequest evcommon.PluginContactRequest

	contactRequest.Plugin = plugin
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token := p.getPluginToken(plugin)
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

	resp, _, _, err = p.PluginCall(contactRequest)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetUUID fetches the UUID from the Origin Resource
func getUUID(origin string) (string, error) {
	var uuid string
	requestData := strings.Split(origin, ":")
	if len(requestData) <= 1 {
		return "", fmt.Errorf("error: SystemUUID not found")
	}
	resource := requestData[0]
	uuid = resource[strings.LastIndexByte(resource, '/')+1:]
	return uuid, nil
}

func (p *PluginContact) createEventSubscrption(taskID string, subTaskChan chan<- int32, reqSessionToken string, targetURI string, request evmodel.RequestBody, originResource string, result *evresponse.MutexLock, wg *sync.WaitGroup, collectionFlag bool, collectionName string) {
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
		log.Error("error while trying to marshal create event request: " + err.Error())
	}
	reqJSON = string(reqBody)

	if taskID != "" {
		subTaskURI, err = p.CreateChildTask(reqSessionToken, taskID)
		if err != nil {
			log.Error("Error while creating the SubTask")
		}
		trimmedURI := strings.TrimSuffix(subTaskURI, "/")
		subTaskID = trimmedURI[strings.LastIndex(trimmedURI, "/")+1:]
		resp.StatusCode = http.StatusAccepted
		p.UpdateTask(fillTaskData(subTaskID, targetURI, reqJSON, resp, common.Running, common.OK, percentComplete, http.MethodPost))
	}

	host, response := p.eventSubscription(request, originResource, collectionName, collectionFlag)
	resp.Body = response.Response
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	resp.StatusCode = int32(response.StatusCode)
	result.AddResponse(originResource, host, response)
	percentComplete = 100
	if subTaskID != "" {
		if response.StatusCode != http.StatusCreated {
			p.UpdateTask(fillTaskData(subTaskID, targetURI, reqJSON, resp, common.Exception, common.Critical, percentComplete, http.MethodPost))
		} else {
			p.UpdateTask(fillTaskData(subTaskID, targetURI, reqJSON, resp, common.Completed, common.OK, percentComplete, http.MethodPost))
		}
		subTaskChan <- int32(response.StatusCode)
	}
}

// checkCollection verifies if the given origin is collection and extracts all the suboridinate resources
func checkCollection(origin string) ([]string, string, bool, error) {
	switch origin {
	case "/redfish/v1/Systems":
		collection, err := evmodel.GetAllKeysFromTable("ComputerSystem")
		return collection, "SystemsCollection", true, err
	case "/redfish/v1/Chassis":
		return []string{}, "ChassisCollection", true, nil
	case "/redfish/v1/Managers":
		//TODO:After Managers implemention need to get all Managers data
		return []string{}, "ManagerCollection", true, nil
	case "/redfish/v1/Fabrics":
		collection, err := evmodel.GetAllFabrics()
		return collection, "FabricsCollection", true, err
	case "/redfish/v1/TaskService/Tasks":
		return []string{}, "TasksCollection", true, nil
	}
	return []string{}, "", false, nil
}

// callPlugin check the given request url and PrefereAuth type plugin
func (p *PluginContact) callPlugin(req evcommon.PluginContactRequest) (*http.Response, error) {
	var reqURL = "https://" + req.Plugin.IP + ":" + req.Plugin.Port + req.URL
	if strings.EqualFold(req.Plugin.PreferredAuthType, "BasicAuth") {
		return p.ContactClient(reqURL, req.HTTPMethodType, "", "", req.PostBody, req.LoginCredential)
	}
	return p.ContactClient(reqURL, req.HTTPMethodType, req.Token, "", req.PostBody, nil)
}

// checkCollectionSubscription checks if any collcetion based subscription exists
// If its' exists it will  update the existing subscription information with newly added server origin
func (p *PluginContact) checkCollectionSubscription(origin, protocol string) {
	//Creating key to get all the System Collection subscription

	var searchKey string
	var bmcFlag bool
	if strings.Contains(origin, "Fabrics") {
		searchKey = "/redfish/v1/Fabrics"
	} else {
		bmcFlag = true
		searchKey = "/redfish/v1/Systems"
	}
	subscriptions, err := evmodel.GetEvtSubscriptions(searchKey)
	if err != nil {
		return
	}
	var chassisSubscriptions, managersSubscriptions []evmodel.Subscription
	if bmcFlag {
		chassisSubscriptions, _ = evmodel.GetEvtSubscriptions("/redfish/v1/Chassis")
		subscriptions = append(subscriptions, chassisSubscriptions...)
		managersSubscriptions, _ = evmodel.GetEvtSubscriptions("/redfish/v1/Managers")
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
		evmodel.OdataIDLink{
			OdataID: origin,
		},
	}

	// Subscribing newly added server with collated event list
	host, response := p.eventSubscription(subscriptionPost, origin, "", false)
	if response.StatusCode != http.StatusCreated {
		return
	}
	for _, evtSubscription := range collectionSubscription {
		evtSubscription.Hosts = append(evtSubscription.Hosts, host)
		err = evmodel.UpdateEventSubscription(evtSubscription)
		if err != nil {
			log.Error("Error while Updating event subscription : " + err.Error())
		}
	}
	// Get Device Subscription Details if collection is bmc and update chassis and managers uri
	if bmcFlag {
		searchKey := evcommon.GetSearchKey(host, evmodel.DeviceSubscriptionIndex)
		deviceSubscription, _ := evmodel.GetDeviceSubscriptions(searchKey)
		data := strings.Split(origin, "/redfish/v1/Systems/")
		chassisList, _ := evmodel.GetAllMatchingDetails("Chassis", data[1], common.InMemory)
		managersList, _ := evmodel.GetAllMatchingDetails("Managers", data[1], common.InMemory)
		var newDevSubscription = evmodel.DeviceSubscription{
			EventHostIP:     deviceSubscription.EventHostIP,
			Location:        deviceSubscription.Location,
			OriginResources: deviceSubscription.OriginResources,
		}
		newDevSubscription.OriginResources = append(newDevSubscription.OriginResources, chassisList...)
		newDevSubscription.OriginResources = append(newDevSubscription.OriginResources, managersList...)

		err := evmodel.UpdateDeviceSubscriptionLocation(newDevSubscription)
		if err != nil {
			log.Error("Error while Updating Device subscription : " + err.Error())
		}
	}

	return
}

func createEventSubscriptionResponse() interface{} {
	return errors.ErrorClass{
		MessageExtendedInfo: []errors.MsgExtendedInfo{
			errors.MsgExtendedInfo{
				MessageID: response.Created,
			},
		},
		Code:    errResponse.Created,
		Message: "See @Message.ExtendedInfo for more information.",
	}
}

func (p *PluginContact) createFabricSubscription(postRequest evmodel.RequestBody, origin, collectionName string, collectionFlag bool) (string, evresponse.EventResponse) {
	var resp evresponse.EventResponse
	var err error
	var plugin *evmodel.Plugin
	var contactRequest evcommon.PluginContactRequest
	log.Info(origin)
	// Extract the fabric id from the Origin
	fabricID := getFabricID(origin)
	fabric, dberr := evmodel.GetFabricData(fabricID)
	if dberr != nil {
		errorMessage := "error while getting fabric data: " + dberr.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
			&resp, []interface{}{"Fabrics", fabricID})
		log.Error(errorMessage)
		return "", resp
	}
	var gerr *errors.Error
	plugin, gerr = evmodel.GetPluginData(fabric.PluginID)
	if gerr != nil {
		errorMessage := "error while getting plugin data: " + gerr.Error() + fabric.PluginID
		evcommon.GenEventErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusNotFound,
			&resp, []interface{}{"Plugin", fabric.PluginID})
		log.Error(errorMessage)
		return "", resp
	}
	contactRequest.Plugin = plugin
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token := p.getPluginToken(plugin)
		if token == "" {
			evcommon.GenEventErrorResponse("error: Unable to create session with plugin "+plugin.ID, errResponse.NoValidSession, http.StatusUnauthorized,
				&resp, []interface{}{})
			log.Error("error: Unable to create session with plugin " + plugin.ID)
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
			evmodel.OdataIDLink{
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
	deviceIPAddress, errorMessage := evcommon.GetIPFromHostName(plugin.IP)
	if errorMessage != "" {
		evcommon.GenEventErrorResponse(errorMessage, errResponse.ResourceNotFound, http.StatusBadRequest,
			&resp, []interface{}{"ManagerAddress", plugin.IP})
		log.Error(errorMessage)
		return "", resp
	}
	var target = evmodel.Target{
		ManagerAddress: deviceIPAddress,
	}
	res, err := p.IsEventsSubscribed("", origin, &subscriptionPost, plugin, &target, collectionFlag, collectionName)
	if err != nil {
		resp.Response = res.Body
		resp.StatusCode = int(res.StatusCode)
		return "", resp
	}

	postBody, _ := json.Marshal(subscriptionPost)
	var reqData string
	//replacing the reruest url with south bound translation URL
	for key, value := range config.Data.URLTranslation.SouthBoundURL {
		reqData = strings.Replace(string(postBody), key, value, -1)
	}
	contactRequest.URL = "/ODIM/v1/Subscriptions"
	contactRequest.HTTPMethodType = http.MethodPost
	err = json.Unmarshal([]byte(reqData), &contactRequest.PostBody)

	response, err := p.callPlugin(contactRequest)
	if err != nil {
		if evcommon.GetPluginStatus(plugin) {
			response, err = p.callPlugin(contactRequest)
		}
		if err != nil {
			evcommon.GenEventErrorResponse(err.Error(), errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			log.Error(err.Error())
			return "", resp
		}
	}
	defer response.Body.Close()
	//retrying the operation if status code is 401
	if response.StatusCode == http.StatusUnauthorized && strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		response, resp, err = p.retryEventSubscriptionOperation(contactRequest)
		if err != nil {
			return "", resp
		}
	}

	log.Error("Subscription Response Status Code" + string(response.StatusCode))
	if response.StatusCode != http.StatusCreated {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			errorMessage := "error while trying to read response body: " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			log.Error(errorMessage)
			return "", resp
		}
		errorMessage := "error while trying to create event subscription"
		var res interface{}
		log.Error("Subscription Response " + string(body))
		err = json.Unmarshal(body, &res)
		if err != nil {
			errorMessage := "error while unmarshaling the body : " + err.Error()
			evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				&resp, []interface{}{})
			log.Error(errorMessage)
			return "", resp
		}

		resp.Response = res
		resp.StatusCode = response.StatusCode
		log.Error(errorMessage)
		return "", resp
	}

	evtSubscription := evmodel.Subscription{
		EventHostIP:    deviceIPAddress,
		OriginResource: origin,
	}

	evtSubscription.Location = response.Header.Get("location")
	err = saveDeviceSubscriptionDetails(evtSubscription)
	if err != nil {
		errorMessage := "error while trying to save event subscription of device data: " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		log.Error(errorMessage)
		return "", resp
	}

	resp.Response = createEventSubscriptionResponse()
	resp.StatusCode = response.StatusCode
	resp.Location = response.Header.Get("location")
	return deviceIPAddress, resp
}

func getFabricID(origin string) string {
	data := strings.Split(origin, "/redfish/v1/Fabrics/")
	if len(data) > 1 {
		fabricData := strings.Split(data[1], "/")
		return fabricData[0]
	}
	return ""
}

// getPluginToken will verify the if any token present to the plugin else it will create token for the new plugin
func (p *PluginContact) getPluginToken(plugin *evmodel.Plugin) string {
	authToken := evcommon.Token.GetToken(plugin.ID)
	if authToken == "" {
		return p.createToken(plugin)
	}
	return authToken
}

func (p *PluginContact) createToken(plugin *evmodel.Plugin) string {
	var contactRequest evcommon.PluginContactRequest

	contactRequest.Plugin = plugin
	contactRequest.HTTPMethodType = http.MethodPost
	contactRequest.PostBody = map[string]interface{}{
		"Username": plugin.Username,
		"Password": string(plugin.Password),
	}
	contactRequest.URL = "/ODIM/v1/Sessions"
	_, _, token, err := p.PluginCall(contactRequest)
	if err != nil {
		log.Error(err.Error())
	}
	pluginToken := evcommon.PluginToken{
		Tokens: make(map[string]string),
	}
	if token != "" {
		pluginToken.StoreToken(plugin.ID, token)
	}
	return token
}

func (p *PluginContact) retryEventOperation(req evcommon.PluginContactRequest) (errResponse.RPC, string, string, error) {
	var resp errResponse.RPC
	var token = p.createToken(req.Plugin)
	if token == "" {
		evcommon.GenErrorResponse("error: Unable to create session with plugin "+req.Plugin.ID, errResponse.NoValidSession, http.StatusUnauthorized,
			[]interface{}{}, &resp)
		return resp, "", "", fmt.Errorf("error: Unable to create session with plugin")
	}
	req.Token = token
	return p.PluginCall(req)
}

func (p *PluginContact) retryEventSubscriptionOperation(req evcommon.PluginContactRequest) (*http.Response, evresponse.EventResponse, error) {
	var resp evresponse.EventResponse
	var token = p.createToken(req.Plugin)
	if token == "" {
		evcommon.GenEventErrorResponse("error: Unable to create session with plugin "+req.Plugin.ID, errResponse.NoValidSession, http.StatusUnauthorized,
			&resp, []interface{}{})
		return nil, resp, fmt.Errorf("error: Unable to create session with plugin")
	}
	req.Token = token

	response, err := p.callPlugin(req)
	if err != nil {
		errorMessage := "error while unmarshaling the body : " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		log.Error(errorMessage)
		return nil, resp, err
	}
	return response, resp, err
}

// isHostPresent will check if hostip present in the hosts slice
func isHostPresent(hosts []string, hostip string) bool {

	if len(hosts) < 1 {
		return false
	}

	front := 0
	rear := len(hosts) - 1
	for front <= rear {
		if hosts[front] == hostip || hosts[rear] == hostip {
			return true
		}
		front++
		rear--
	}
	return false
}
