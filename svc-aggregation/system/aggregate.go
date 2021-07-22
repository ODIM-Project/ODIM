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

// Package system ...
package system

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
	uuid "github.com/satori/go.uuid"
)

//ResetRequest is struct for reset of elements of an aggregate
type ResetRequest struct {
	BatchSize                    int    `json:"BatchSize"`
	DelayBetweenBatchesInSeconds int    `json:"DelayBetweenBatchesInSeconds"`
	ResetType                    string `json:"ResetType"`
}

// CreateAggregate is the handler for creating an aggregate
// check if the elelments/resources added into odimra if not then return an error.
// else add an entry of an aggregayte in db
func (e *ExternalInterface) CreateAggregate(req *aggregatorproto.AggregatorRequest) response.RPC {
	var resp response.RPC
	// parsing the aggregate request
	var createRequest agmodel.Aggregate
	err := json.Unmarshal(req.RequestBody, &createRequest)
	if err != nil {
		errMsg := "unable to parse the create request" + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errMsg, nil, nil)
	}
	//empty request check
	if reflect.DeepEqual(agmodel.Aggregate{}, createRequest) {
		errMsg := "empty request can not be processed"
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"Elements"}, nil)
	}

	statuscode, err := validateElements(createRequest.Elements)
	if err != nil {
		errMsg := "invalid elements for create an aggregate" + err.Error()
		log.Error(errMsg)
		errArgs := []interface{}{"Elements", string(req.RequestBody)}
		return common.GeneralError(statuscode, response.ResourceNotFound, errMsg, errArgs, nil)
	}
	targetURI := "/redfish/v1/AggregationService/Aggregates"
	aggregateUUID := uuid.NewV4().String()
	var aggregateURI = fmt.Sprintf("%s/%s", targetURI, aggregateUUID)

	dbErr := agmodel.CreateAggregate(createRequest, aggregateURI)
	if dbErr != nil {
		errMsg := dbErr.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	commonResponse := response.Response{
		OdataType:    "#Aggregate.v1_0_0.Aggregate",
		OdataID:      aggregateURI,
		OdataContext: "/redfish/v1/$metadata#Aggregate.Aggregate",
		ID:           aggregateUUID,
		Name:         "Aggregate",
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Link":              "<" + aggregateURI + "/>; rel=describedby",
		"Location":          aggregateURI,
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	commonResponse.CreateGenericResponse(response.Created)
	resp.Body = agresponse.AggregateResponse{
		Response: commonResponse,
		Elements: createRequest.Elements,
	}
	resp.StatusCode = http.StatusCreated
	return resp
}

// check if the resource is exist in odim
func validateElements(elements []string) (int32, error) {
	if checkDuplicateElements(elements) {
		return http.StatusBadRequest, errors.PackError(errors.UndefinedErrorType, fmt.Errorf("Duplicate elements present"))
	}
	for _, element := range elements {
		if _, err := agmodel.GetComputerSystem(element); err != nil {
			return http.StatusNotFound, err
		}
	}
	return http.StatusOK, nil
}

//check if the elements have duplicate element
func checkDuplicateElements(elelments []string) bool {
	duplicate := make(map[string]int)
	for _, element := range elelments {
		// check if the item/element exist in the duplicate map
		_, exist := duplicate[element]
		if exist {
			return true
		}
		duplicate[element] = 1

	}
	return false
}

// GetAllAggregates is the handler for getting collection of aggregates
func (e *ExternalInterface) GetAllAggregates(req *aggregatorproto.AggregatorRequest) response.RPC {
	aggregateKeys, err := agmodel.GetAllKeysFromTable("Aggregate")
	if err != nil {
		log.Error("error getting aggregate : " + err.Error())
		errorMessage := err.Error()
		return common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, errorMessage, []interface{}{config.Data.DBConf.OnDiskHost + ":" + config.Data.DBConf.OnDiskPort}, nil)
	}
	var members = make([]agresponse.ListMember, 0)
	for i := 0; i < len(aggregateKeys); i++ {
		members = append(members, agresponse.ListMember{
			OdataID: aggregateKeys[i],
		})
	}
	var resp = response.RPC{
		StatusCode:    http.StatusOK,
		StatusMessage: response.Success,
	}
	commonResponse := response.Response{
		OdataType:    "#AggregateCollection.AggregateCollection",
		OdataID:      "/redfish/v1/AggregationService/Aggregates",
		OdataContext: "/redfish/v1/$metadata#AggregateCollection.AggregateCollection",
		Name:         "Aggregate",
		Description:  "Aggregate collection view",
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}

	resp.Body = agresponse.List{
		Response:     commonResponse,
		MembersCount: len(members),
		Members:      members,
	}
	return resp
}

// GetAggregate is the handler for getting an aggregate
//if the aggregate id is present then return aggregate details else return an error.
func (e *ExternalInterface) GetAggregate(req *aggregatorproto.AggregatorRequest) response.RPC {
	aggregate, err := agmodel.GetAggregate(req.URL)
	if err != nil {
		log.Error("error getting  Aggregate : " + err.Error())
		errorMessage := err.Error()
		if errors.DBKeyNotFound == err.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"Aggregate", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	var data = strings.Split(req.URL, "/redfish/v1/AggregationService/Aggregates/")
	commonResponse := response.Response{
		OdataType:    "#Aggregate.v1_0_0.Aggregate",
		OdataID:      req.URL,
		OdataContext: "/redfish/v1/$metadata#Aggregate.Aggregate",
		ID:           data[1],
		Name:         "Aggregate",
	}
	var resp = response.RPC{
		StatusCode:    http.StatusOK,
		StatusMessage: response.Success,
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	commonResponse.CreateGenericResponse(response.Success)
	resp.Body = agresponse.AggregateResponse{
		Response: commonResponse,
		Elements: aggregate.Elements,
	}
	return resp
}

// DeleteAggregate is the handler for deleting an aggregate
// if the aggregate id is present then delete from the db else return an error.
func (e *ExternalInterface) DeleteAggregate(req *aggregatorproto.AggregatorRequest) response.RPC {
	var resp response.RPC
	_, err := agmodel.GetAggregate(req.URL)
	if err != nil {
		log.Error("error getting  Aggregate : " + err.Error())
		errorMessage := err.Error()
		if errors.DBKeyNotFound == err.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"Aggregate", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	err = agmodel.DeleteAggregate(req.URL)
	if err != nil {
		log.Error("error while deleting an aggregate : " + err.Error())
		errorMessage := err.Error()
		if errors.DBKeyNotFound == err.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"Aggregate", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	resp.StatusCode = http.StatusNoContent
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	return resp
}

// AddElementsToAggregate is the handler for adding elements to an aggregate
func (e *ExternalInterface) AddElementsToAggregate(req *aggregatorproto.AggregatorRequest) response.RPC {
	var resp response.RPC
	// parsing the aggregate request
	var addRequest agmodel.Aggregate
	err := json.Unmarshal(req.RequestBody, &addRequest)
	if err != nil {
		errMsg := "unable to parse the create request" + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errMsg, nil, nil)
	}
	//empty request check
	if reflect.DeepEqual(agmodel.Aggregate{}, addRequest) || reflect.DeepEqual(addRequest.Elements, []string{}) {
		errMsg := "empty request can not be processed"
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"Elements"}, nil)
	}

	statuscode, err := validateElements(addRequest.Elements)
	if err != nil {
		errMsg := "invalid elements for create an aggregate" + err.Error()
		log.Error(errMsg)
		errArgs := []interface{}{"Elements", fmt.Sprintf("%v", addRequest)}
		return common.GeneralError(statuscode, response.ResourceNotFound, errMsg, errArgs, nil)
	}

	if req.URL == "" {
		errMsg := "request uri is not provided"
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"request uri"}, nil)
	}
	url := strings.Split(req.URL, "/redfish/v1/AggregationService/Aggregates/")
	aggregateID := strings.Split(url[1], "/")[0]
	aggregateURL := "/redfish/v1/AggregationService/Aggregates/" + aggregateID
	aggregate, err1 := agmodel.GetAggregate(aggregateURL)
	if err1 != nil {
		log.Error("error getting  Aggregate : " + err1.Error())
		errorMessage := err1.Error()
		if errors.DBKeyNotFound == err1.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"Aggregate", aggregateURL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	if checkElementsPresent(addRequest.Elements, aggregate.Elements) {
		errMsg := "Elements present in aggregate"
		log.Error(errMsg)
		errArgs := []interface{}{"AddElements", "Elements", fmt.Sprintf("%v", addRequest.Elements)}
		return common.GeneralError(http.StatusConflict, response.ResourceAlreadyExists, errMsg, errArgs, nil)
	}

	dbErr := agmodel.AddElementsToAggregate(addRequest, aggregateURL)
	if dbErr != nil {
		errMsg := dbErr.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	commonResponse := response.Response{
		OdataType:    "#Aggregate.v1_0_0.Aggregate",
		OdataID:      aggregateURL,
		OdataContext: "/redfish/v1/$metadata#Aggregate.Aggregate",
		ID:           aggregateID,
		Name:         "Aggregate",
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Link":              "<" + aggregateURL + "/>; rel=describedby",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	aggregate, _ = agmodel.GetAggregate(aggregateURL)
	commonResponse.CreateGenericResponse(response.Success)
	resp.Body = agresponse.AggregateResponse{
		Response: commonResponse,
		Elements: aggregate.Elements,
	}
	resp.StatusCode = http.StatusOK
	return resp
}

// RemoveElementsFromAggregate is the handler for removing elements from an aggregate
func (e *ExternalInterface) RemoveElementsFromAggregate(req *aggregatorproto.AggregatorRequest) response.RPC {
	var resp response.RPC
	// parsing the aggregate request
	var removeRequest agmodel.Aggregate
	err := json.Unmarshal(req.RequestBody, &removeRequest)
	if err != nil {
		errMsg := "unable to parse the create request" + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errMsg, nil, nil)
	}

	//empty request check
	if reflect.DeepEqual(agmodel.Aggregate{}, removeRequest) || reflect.DeepEqual(removeRequest.Elements, []string{}) {
		errMsg := "empty request can not be processed"
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"Elements"}, nil)
	}

	if req.URL == "" {
		errMsg := "request uri is not provided"
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"request uri"}, nil)
	}
	if checkDuplicateElements(removeRequest.Elements) {
		errMsg := "duplicate elements present"
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.ResourceCannotBeDeleted, errMsg, nil, nil)
	}
	url := strings.Split(req.URL, "/redfish/v1/AggregationService/Aggregates/")
	aggregateID := strings.Split(url[1], "/")[0]

	aggregateURL := "/redfish/v1/AggregationService/Aggregates/" + aggregateID
	aggregate, err1 := agmodel.GetAggregate(aggregateURL)
	if err != nil {
		log.Error("error getting aggregate : " + err1.Error())
		errorMessage := err1.Error()
		if errors.DBKeyNotFound == err1.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"Aggregate", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	if !checkRemovingElementsPresent(removeRequest.Elements, aggregate.Elements) {
		errMsg := "Elements not present in aggregate"
		log.Error(errMsg)
		errArgs := []interface{}{"Elements", fmt.Sprintf("%v", removeRequest.Elements)}
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, errArgs, nil)
	}

	dbErr := agmodel.RemoveElementsFromAggregate(removeRequest, aggregateURL)
	if dbErr != nil {
		errMsg := dbErr.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	commonResponse := response.Response{
		OdataType:    "#Aggregate.v1_0_0.Aggregate",
		OdataID:      aggregateURL,
		OdataContext: "/redfish/v1/$metadata#Aggregate.Aggregate",
		ID:           aggregateID,
		Name:         "Aggregate",
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Link":              "<" + aggregateURL + "/>; rel=describedby",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	aggregate, _ = agmodel.GetAggregate(aggregateURL)
	commonResponse.CreateGenericResponse(response.Success)
	resp.Body = agresponse.AggregateResponse{
		Response: commonResponse,
		Elements: aggregate.Elements,
	}
	resp.StatusCode = http.StatusOK
	return resp
}

func checkElementsPresent(requestElements, presentElements []string) bool {
	for _, element := range requestElements {
		front := 0
		rear := len(presentElements) - 1
		for front <= rear {
			if presentElements[front] == element || presentElements[rear] == element {
				return true
			}
			front++
			rear--
		}
	}
	return false
}

func checkRemovingElementsPresent(requestElements, presentElements []string) bool {
	for _, element := range requestElements {
		var present bool
		front := 0
		rear := len(presentElements) - 1
		for front <= rear {
			if presentElements[front] == element || presentElements[rear] == element {
				present = true
			}
			front++
			rear--
		}
		if !present {
			return false
		}
	}
	return true
}

// ResetElementsOfAggregate is the handler for reseting elements of an aggregate
func (e *ExternalInterface) ResetElementsOfAggregate(taskID string, sessionUserName string, req *aggregatorproto.AggregatorRequest) response.RPC {
	var resp response.RPC
	var percentComplete int32
	targetURI := req.URL
	percentComplete = 0

	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: targetURI, UpdateTask: e.UpdateTask, TaskRequest: string(req.RequestBody)}

	var resetRequest ResetRequest
	if err := json.Unmarshal(req.RequestBody, &resetRequest); err != nil {
		errMsg := "error while trying to validate request fields: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errMsg, nil, taskInfo)
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, resetRequest)
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Error(errorMessage)
		resp := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, taskInfo)
		return resp
	}

	missedProperty, err := resetRequest.validateRequestFields()
	if err != nil {
		errMsg := "error while trying to validate request fields: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{missedProperty}, taskInfo)
	}

	url := strings.Split(req.URL, "/redfish/v1/AggregationService/Aggregates/")
	aggregateID := strings.Split(url[1], "/")[0]

	aggregateURL := "/redfish/v1/AggregationService/Aggregates/" + aggregateID
	aggregate, err1 := agmodel.GetAggregate(aggregateURL)
	if err1 != nil {
		errorMessage := err1.Error()
		log.Error("error getting aggregate : " + errorMessage)
		if errors.DBKeyNotFound == err1.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err1.Error(), []interface{}{"Aggregate", req.URL}, taskInfo)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, taskInfo)
	}

	// subTaskChan is a buffered channel with buffer size equal to total number of elements.
	// this also helps while cancelling the task. even if the reader is not available for reading
	// the channel buffer will collect them and allows gracefull exit for already spanned goroutines.
	subTaskChan := make(chan int32, len(aggregate.Elements))
	resp.StatusCode = http.StatusOK
	var cancelled, partialResultFlag bool
	var wg, writeWG sync.WaitGroup
	go func() {
		for i := 0; i < len(aggregate.Elements); i++ {
			if cancelled == false { // task cancelled check to determine whether to collect status codes.
				select {
				case statusCode := <-subTaskChan:
					if statusCode != http.StatusOK {
						partialResultFlag = true
						if resp.StatusCode < statusCode {
							resp.StatusCode = statusCode
						}
					}

					if i < len(aggregate.Elements)-1 {
						percentComplete = int32(((i + 1) / len(aggregate.Elements)) * 100)
						var task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Running, common.OK, percentComplete, http.MethodPost)
						err := e.UpdateTask(task)
						if err != nil && err.Error() == common.Cancelling {
							task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Cancelled, common.OK, percentComplete, http.MethodPost)
							e.UpdateTask(task)
							cancelled = true
						}
					}
				}
			}
			writeWG.Done()
		}
	}()

	var tempIndex int
	for _, element := range aggregate.Elements {
		wg.Add(1)
		writeWG.Add(1)
		// tempIndex is for checking batch size, its increment on each iteration
		// if its equal to batch size then reinitilise.
		// if batch size is 0 then reset all the systems without any kind of batch and ignore the DelayBetweenBatchesInSeconds
		tempIndex = tempIndex + 1
		if resetRequest.BatchSize == 0 || tempIndex <= resetRequest.BatchSize {
			go e.resetSystem(taskID, string(req.RequestBody), subTaskChan, sessionUserName, element, resetRequest.ResetType, &wg)
		}

		if tempIndex == resetRequest.BatchSize && resetRequest.BatchSize != 0 {
			tempIndex = 0
			time.Sleep(time.Second * time.Duration(resetRequest.DelayBetweenBatchesInSeconds))
		}

	}
	wg.Wait()
	writeWG.Wait()
	taskStatus := common.OK
	if partialResultFlag {
		taskStatus = common.Warning
	}
	percentComplete = 100
	var args response.Args
	if resp.StatusCode != http.StatusOK {
		errMsg := "one or more of the reset actions failed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID
		log.Error(errMsg)
		return common.GeneralError(resp.StatusCode, resp.StatusMessage, errMsg, nil, taskInfo)
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	log.Info("all reset actions are successfully completed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID)
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	args = response.Args{
		Code:    resp.StatusMessage,
		Message: "Request completed successfully",
	}
	resp.Body = args.CreateGenericErrorResponse()
	var task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Completed, taskStatus, percentComplete, http.MethodPost)
	err = e.UpdateTask(task)
	if err != nil && err.Error() == common.Cancelling {
		task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Cancelled, common.Critical, percentComplete, http.MethodPost)
		e.UpdateTask(task)
		runtime.Goexit()
	}
	return resp
}

func (e *ExternalInterface) resetSystem(taskID, reqBody string, subTaskChan chan<- int32, sessionUserName, element, resetType string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info("INFO: reset(type: " + resetType + ") of the target " + element + " has been started.")
	var resp response.RPC
	var percentComplete int32
	//Create the child Task
	subTaskURI, err := e.CreateChildTask(sessionUserName, taskID)
	if err != nil {
		subTaskChan <- http.StatusInternalServerError
		log.Error("error while trying to create sub task")
		return
	}
	var subTaskID string
	strArray := strings.Split(subTaskURI, "/")
	if strings.HasSuffix(subTaskURI, "/") {
		subTaskID = strArray[len(strArray)-2]
	} else {
		subTaskID = strArray[len(strArray)-1]
	}
	systemID := element[strings.LastIndexAny(element, "/")+1:]
	var targetURI = element
	taskInfo := &common.TaskUpdateInfo{TaskID: subTaskID, TargetURI: targetURI, UpdateTask: e.UpdateTask, TaskRequest: reqBody}
	data := strings.Split(systemID, ":")
	if len(data) <= 1 {
		subTaskChan <- http.StatusNotFound
		errMsg := "error: SystemUUID not found"
		log.Error(errMsg)
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"SystemUUID", ""}, taskInfo)
		return
	}

	uuid, sysID := data[0], data[1]
	// Get target device Credentials from using device UUID
	target, err := agmodel.GetTarget(uuid)
	if err != nil {
		subTaskChan <- http.StatusNotFound
		errMsg := err.Error()
		log.Error(errMsg)
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"target", uuid}, taskInfo)
		return
	}
	decryptedPasswordByte, err := e.DecryptPassword(target.Password)
	if err != nil {
		errMsg := "error while trying to decrypt device password: " + err.Error()
		subTaskChan <- http.StatusInternalServerError
		log.Error(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		return
	}
	target.Password = decryptedPasswordByte
	// Get the Plugin info
	plugin, errs := agmodel.GetPluginData(target.PluginID)
	if errs != nil {
		subTaskChan <- http.StatusNotFound
		errMsg := errs.Error()
		log.Error(errMsg)
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"plugin", target.PluginID}, taskInfo)
		return
	}

	var pluginContactRequest getResourceRequest
	pluginContactRequest.ContactClient = e.ContactClient
	pluginContactRequest.GetPluginStatus = e.GetPluginStatus
	pluginContactRequest.Plugin = plugin
	pluginContactRequest.StatusPoll = true
	pluginContactRequest.TaskRequest = reqBody

	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		pluginContactRequest.HTTPMethodType = http.MethodPost
		pluginContactRequest.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		pluginContactRequest.OID = "/ODIM/v1/Sessions"
		_, token, getResponse, err := contactPlugin(pluginContactRequest, "error while logging in to plugin: ")
		if err != nil {
			subTaskChan <- getResponse.StatusCode
			errMsg := err.Error()
			log.Error(errMsg)
			common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo)
			return
		}
		pluginContactRequest.Token = token
	} else {
		pluginContactRequest.LoginCredentials = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	// Adding system state entry to db
	postRequest := make(map[string]interface{})
	postRequest["ResetType"] = resetType
	postBody, _ := json.Marshal(postRequest)
	target.PostBody = postBody
	pluginContactRequest.DeviceInfo = target
	pluginContactRequest.OID = "/ODIM/v1/Systems/" + sysID + "/Actions/ComputerSystem.Reset"
	pluginContactRequest.HTTPMethodType = http.MethodPost
	_, _, getResponse, err := contactPlugin(pluginContactRequest, "error while reseting the computer system: ")

	if err != nil {
		subTaskChan <- getResponse.StatusCode
		errMsg := err.Error()
		log.Error(errMsg)
		common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo)
		return
	}

	resp.StatusMessage = response.Success
	resp.Body = response.ErrorClass{
		Code:    resp.StatusMessage,
		Message: "Request completed successfully.",
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
		"Location":          element,
	}
	resp.StatusCode = getResponse.StatusCode
	percentComplete = 100
	subTaskChan <- int32(getResponse.StatusCode)
	var task = fillTaskData(subTaskID, targetURI, reqBody, resp, common.Completed, common.OK, percentComplete, http.MethodPost)
	err = e.UpdateTask(task)
	if err != nil && err.Error() == common.Cancelling {
		var task = fillTaskData(subTaskID, targetURI, reqBody, resp, common.Cancelled, common.Critical, percentComplete, http.MethodPost)
		err = e.UpdateTask(task)
	}
	if getResponse.StatusCode == http.StatusOK {
		agmodel.AddSystemResetInfo(element, resetType)
	}
	return
}

// validateRequestFields validate each field in the request against default value of field type
func (validateReq ResetRequest) validateRequestFields() (string, error) {
	if reflect.DeepEqual(validateReq, ResetRequest{}) {
		return "ResetRequest", fmt.Errorf("ResetRequest is empty")
	}

	if validateReq.ResetType == "" {
		return "ResetType", fmt.Errorf("property ResetType missing in the reset request")
	}
	return "", nil
}

// SetDefaultBootOrderElementsOfAggregate is the handler for set default boot order elements of an aggregate
func (e *ExternalInterface) SetDefaultBootOrderElementsOfAggregate(taskID string, sessionUserName string, req *aggregatorproto.AggregatorRequest) response.RPC {
	var resp response.RPC
	var percentComplete int32 = 100
	targetURI := req.URL

	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: targetURI, UpdateTask: e.UpdateTask}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, taskInfo)
	}
	reqJSON := string(reqBody)
	taskInfo.TaskRequest = reqJSON

	url := strings.Split(req.URL, "/redfish/v1/AggregationService/Aggregates/")
	aggregateID := strings.Split(url[1], "/")[0]

	aggregateURL := "/redfish/v1/AggregationService/Aggregates/" + aggregateID
	aggregate, aggErr := agmodel.GetAggregate(aggregateURL)
	if aggErr != nil {
		errorMessage := aggErr.Error()
		log.Error("error getting aggregate : " + errorMessage)
		if errors.DBKeyNotFound == aggErr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, aggErr.Error(), []interface{}{"Aggregate", req.URL}, taskInfo)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, taskInfo)
	}

	partialResultFlag := false
	subTaskChan := make(chan int32, len(aggregate.Elements))
	for _, element := range aggregate.Elements {
		go e.collectAndSetDefaultOrder(taskID, element, reqJSON, subTaskChan, sessionUserName)
	}
	resp.StatusCode = http.StatusOK
	for i := 0; i < len(aggregate.Elements); i++ {
		select {
		case statusCode := <-subTaskChan:
			if statusCode != http.StatusOK {
				partialResultFlag = true
				if resp.StatusCode < statusCode {
					resp.StatusCode = statusCode
				}
			}
			if i < len(aggregate.Elements)-1 {
				percentComplete := int32(((i + 1) / len(aggregate.Elements)) * 100)
				var task = fillTaskData(taskID, targetURI, reqJSON, resp, common.Running, common.OK, percentComplete, http.MethodPost)
				err := e.UpdateTask(task)
				if err != nil && err.Error() == common.Cancelling {
					task = fillTaskData(taskID, targetURI, reqJSON, resp, common.Cancelled, common.OK, percentComplete, http.MethodPost)
					e.UpdateTask(task)
					runtime.Goexit()
				}

			}
		}
	}

	taskStatus := common.OK
	if partialResultFlag {
		taskStatus = common.Warning
	}
	percentComplete = 100
	if resp.StatusCode != http.StatusOK {
		errMsg := "one or more of the SetDefaultBootOrder requests failed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID
		log.Error(errMsg)
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return common.GeneralError(http.StatusUnauthorized, response.ResourceAtURIUnauthorized, errMsg, []interface{}{aggregate.Elements}, taskInfo)
		case http.StatusNotFound:
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"option", "SetDefaultBootOrder"}, taskInfo)
		default:
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		}
	}

	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	log.Error("all SetDefaultBootOrder requests successfully completed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID)
	resp.StatusMessage = response.Success
	resp.StatusCode = http.StatusOK
	args := response.Args{
		Code:    resp.StatusMessage,
		Message: "Request completed successfully",
	}
	resp.Body = args.CreateGenericErrorResponse()

	var task = fillTaskData(taskID, targetURI, reqJSON, resp, common.Completed, taskStatus, percentComplete, http.MethodPost)
	err = e.UpdateTask(task)
	if err != nil && err.Error() == common.Cancelling {
		task = fillTaskData(taskID, targetURI, reqJSON, resp, common.Cancelled, common.Critical, percentComplete, http.MethodPost)
		e.UpdateTask(task)
		runtime.Goexit()
	}
	return resp
}
