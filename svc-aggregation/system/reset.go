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

package system

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

// AggregatorRequest ...
type AggregatorRequest struct {
	OdataContext string    `json:"@odata.context"`
	OdataID      string    `json:"@odata.id"`
	Odatatype    string    `json:"@odata.type"`
	ID           string    `json:"Id"`
	Name         string    `json:"Name"`
	Oem          OEM       `json:"Oem"`
	Parameters   parameter `json:"Parameters"`
}

// OEM is a placeholder for the OEM block
type OEM struct{}

type parameter struct {
	ResetCollection resetCollection `json:"ResetCollection"`
}

// ResetTarget holds the target details for reset
type ResetTarget struct {
	ResetType string `json:"ResetType"`
	TargetURI string `json:"TargetUri"`
	Delay     int    `json:"Delay"`
	Priority  int    `json:"Priority"`
}

type resetCollection struct {
	Description  string        `json:"Description"`
	ResetTargets []ResetTarget `json:"ResetTarget"`
}

// validateRequestFields validate each field in the request against default value of field type
func (validateReq AggregatorRequest) validateRequestFields() (string, error) {
	if reflect.DeepEqual(validateReq.Parameters, parameter{}) {
		return "parameters", fmt.Errorf("property parameters missing in the reset request")
	}

	if reflect.DeepEqual(validateReq.Parameters.ResetCollection, resetCollection{}) {
		return "ResetCollection", fmt.Errorf("property ResetCollection missing in the reset request")
	}

	if len(validateReq.Parameters.ResetCollection.ResetTargets) <= 0 {
		return "ResetTarget", fmt.Errorf("property ResetTarget missing in the reset request")
	}

	for _, resetTarget := range validateReq.Parameters.ResetCollection.ResetTargets {
		if resetTarget.ResetType == "" {
			return "ResetType", fmt.Errorf("property ResetType missing in the reset request")
		}
		if resetTarget.TargetURI == "" {
			return "TargetURI", fmt.Errorf("property TargetURI missing in the reset request")
		}
	}
	return "", nil
}

// Reset is for reseting the computer systems mentioned in the request body
func (e *ExternalInterface) Reset(taskID string, sessionUserName string, req *aggregatorproto.AggregatorRequest) response.RPC {
	var resp response.RPC
	var percentComplete int32
	targetURI := "/redfish/v1/AggregationService/Actions/AggregationService.Reset/" // this will removed later and passed as input param in req struct
	percentComplete = 0

	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: targetURI, UpdateTask: e.UpdateTask}

	var resetRequest AggregatorRequest
	if err := json.Unmarshal(req.RequestBody, &resetRequest); err != nil {
		errMsg := "error while trying to validate request fields: " + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errMsg, nil, taskInfo)
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, resetRequest)
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Println(errorMessage)
		resp := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, taskInfo)
		return resp
	}

	missedProperty, err := resetRequest.validateRequestFields()
	if err != nil {
		errMsg := "error while trying to validate request fields: " + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{missedProperty}, taskInfo)
	}

	partialResultFlag := false

	var wg, writeWG sync.WaitGroup
	prioriyCounter := make(map[int]int)
	resetRequest.Parameters.ResetCollection.ResetTargets, prioriyCounter = sortTargetsWithPriorityAndDelay(
		checkAndCorrectPriorityAndDelay(
			resetRequest.Parameters.ResetCollection.ResetTargets,
		),
	)

	// subTaskChan is a buffered channel with buffer size equal to total number of reset actions.
	// this also helps while cancelling the task. even if the reader is not available for reading
	// the channel buffer will collect them and allows gracefull exit for already spanned goroutines.
	subTaskChan := make(chan int32, len(resetRequest.Parameters.ResetCollection.ResetTargets))
	resp.StatusCode = http.StatusOK
	var cancelled bool
	go func() {
		for i := 0; i < len(resetRequest.Parameters.ResetCollection.ResetTargets); i++ {
			if cancelled == false { // task cancelled check to determine whether to collect status codes.
				select {
				case statusCode := <-subTaskChan:
					if statusCode != http.StatusOK {
						partialResultFlag = true
						if resp.StatusCode < statusCode {
							resp.StatusCode = statusCode
						}
					}

					if i < len(resetRequest.Parameters.ResetCollection.ResetTargets)-1 {
						percentComplete = int32(((i + 1) / len(resetRequest.Parameters.ResetCollection.ResetTargets)) * 100)
						var task = fillTaskData(taskID, targetURI, resp, common.Running, common.OK, percentComplete, http.MethodPost)
						err := e.UpdateTask(task)
						if err != nil && err.Error() == common.Cancelling {
							task = fillTaskData(taskID, targetURI, resp, common.Cancelled, common.OK, percentComplete, http.MethodPost)
							e.UpdateTask(task)
							cancelled = true
						}
					}
				}
			}
			writeWG.Done()
		}
	}()
	executionCounter := make(map[int]int)
	for _, resetTarget := range resetRequest.Parameters.ResetCollection.ResetTargets {
		writeWG.Add(1)
		if cancelled == false { // task cancelled check to determine whether to schedule next reset action
			wg.Add(1)
			executionCounter[resetTarget.Priority]++
			go e.resetComputerSystem(taskID, subTaskChan, sessionUserName, resetTarget, time.NewTimer(time.Duration(resetTarget.Delay)*time.Second), &wg)
			if executionCounter[resetTarget.Priority] == prioriyCounter[resetTarget.Priority] {
				wg.Wait() // waiting for all reset actions in same priority level to get finished
			}
		}
	}

	writeWG.Wait()

	taskStatus := common.OK
	if partialResultFlag {
		taskStatus = common.Warning
	}
	percentComplete = 100
	var args response.Args
	if resp.StatusCode != http.StatusOK {
		errMsg := "one or more of the reset actions failed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID
		log.Printf(errMsg)
		return common.GeneralError(resp.StatusCode, resp.StatusMessage, errMsg, nil, taskInfo)
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	log.Println("all reset actions are successfully completed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID)
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	args = response.Args{
		Code:    resp.StatusMessage,
		Message: "Request completed successfully",
	}
	resp.Body = args.CreateGenericErrorResponse()
	var task = fillTaskData(taskID, targetURI, resp, common.Completed, taskStatus, percentComplete, http.MethodPost)
	err = e.UpdateTask(task)
	if err != nil && err.Error() == common.Cancelling {
		task = fillTaskData(taskID, targetURI, resp, common.Cancelled, common.Critical, percentComplete, http.MethodPost)
		e.UpdateTask(task)
		runtime.Goexit()
	}
	return resp
}

func checkAndCorrectPriorityAndDelay(resetTargets []ResetTarget) []ResetTarget {
	for i := 0; i < len(resetTargets); i++ {
		if resetTargets[i].Priority < config.Data.ExecPriorityDelayConf.MinResetPriority {
			resetTargets[i].Priority = config.Data.ExecPriorityDelayConf.MinResetPriority
		} else if resetTargets[i].Priority > config.Data.ExecPriorityDelayConf.MaxResetPriority {
			resetTargets[i].Priority = config.Data.ExecPriorityDelayConf.MaxResetPriority
		}
		if resetTargets[i].Delay < 0 {
			resetTargets[i].Delay = 0
		} else if resetTargets[i].Delay > config.Data.ExecPriorityDelayConf.MaxResetDelayInSecs {
			resetTargets[i].Delay = config.Data.ExecPriorityDelayConf.MaxResetDelayInSecs
		}
	}
	return resetTargets
}

func sortTargetsWithPriorityAndDelay(resetTargets []ResetTarget) ([]ResetTarget, map[int]int) {
	var orderedTargets []ResetTarget
	priorityCounter := make(map[int]int)

	for i := config.Data.ExecPriorityDelayConf.MaxResetPriority; i >= config.Data.ExecPriorityDelayConf.MinResetPriority; i-- {
		var prioritySlice, tempSlice []ResetTarget
		var nextPriority int
		for _, target := range resetTargets {
			if target.Priority == i { // sorting resets in descending order of Priority
				prioritySlice = append(prioritySlice, target) // collecting all resets of same Priority
				priorityCounter[target.Priority]++            // counting number of reset actions with same priority
			} else {
				tempSlice = append(tempSlice, target)
				if target.Priority >= nextPriority { // collecting the next highest Priority
					nextPriority = target.Priority
				}
			}
		}
		resetTargets = tempSlice // deleting already sorted targets from slice

		// sorting resets of same Priority in ascending order of Delay
		sort.Slice(prioritySlice, func(i, j int) bool { return prioritySlice[i].Delay < prioritySlice[j].Delay })

		// creation of sorted reset targets
		for _, orderedTarget := range prioritySlice {
			orderedTargets = append(orderedTargets, orderedTarget)
		}

		if i != config.Data.ExecPriorityDelayConf.MinResetPriority {
			i = nextPriority + 1 // skipping to next highest priority present. +1 is for neutralizing loop's decrement
		}
	}
	return orderedTargets, priorityCounter
}

func (e *ExternalInterface) resetComputerSystem(taskID string, subTaskChan chan<- int32, sessionUserName string, req ResetTarget, timer *time.Timer, wg *sync.WaitGroup) {
	defer wg.Done()
	<-timer.C
	log.Printf("INFO: reset(type: %v) of the target %v with priority level %v has been started after a delay of %v seconds.", req.ResetType, req.TargetURI, req.Priority, req.Delay)
	var resp response.RPC
	var percentComplete int32
	//Create the child Task
	subTaskURI, err := e.CreateChildTask(sessionUserName, taskID)
	if err != nil {
		subTaskChan <- http.StatusInternalServerError
		log.Println("error while trying to create sub task")
		return
	}
	var subTaskID string
	strArray := strings.Split(subTaskURI, "/")
	if strings.HasSuffix(subTaskURI, "/") {
		subTaskID = strArray[len(strArray)-2]
	} else {
		subTaskID = strArray[len(strArray)-1]
	}
	systemID := req.TargetURI[strings.LastIndexAny(req.TargetURI, "/")+1:]
	var targetURI = req.TargetURI
	taskInfo := &common.TaskUpdateInfo{TaskID: subTaskID, TargetURI: targetURI, UpdateTask: e.UpdateTask}
	data := strings.Split(systemID, ":")
	if len(data) <= 1 {
		subTaskChan <- http.StatusNotFound
		errMsg := "error: SystemUUID not found"
		log.Printf(errMsg)
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"SystemUUID", ""}, taskInfo)
		return
	}

	uuid, sysID := data[0], data[1]
	// Get target device Credentials from using device UUID
	target, err := agmodel.GetTarget(uuid)
	if err != nil {
		subTaskChan <- http.StatusNotFound
		errMsg := err.Error()
		log.Printf(errMsg)
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"target", uuid}, taskInfo)
		return
	}
	decryptedPasswordByte, err := e.DecryptPassword(target.Password)
	if err != nil {
		errMsg := "error while trying to decrypt device password: " + err.Error()
		subTaskChan <- http.StatusInternalServerError
		log.Printf(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		return
	}
	target.Password = decryptedPasswordByte
	// Get the Plugin info
	plugin, errs := agmodel.GetPluginData(target.PluginID)
	if errs != nil {
		subTaskChan <- http.StatusNotFound
		errMsg := errs.Error()
		log.Printf(errMsg)
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"plugin", target.PluginID}, taskInfo)
		return
	}

	var pluginContactRequest getResourceRequest
	pluginContactRequest.ContactClient = e.ContactClient
	pluginContactRequest.GetPluginStatus = e.GetPluginStatus
	pluginContactRequest.Plugin = plugin
	pluginContactRequest.StatusPoll = true

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
			log.Println(errMsg)
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
	postRequest["ResetType"] = req.ResetType
	postBody, _ := json.Marshal(postRequest)
	target.PostBody = postBody
	pluginContactRequest.DeviceInfo = target
	pluginContactRequest.OID = "/ODIM/v1/Systems/" + sysID + "/Actions/ComputerSystem.Reset"
	pluginContactRequest.HTTPMethodType = http.MethodPost
	_, _, getResponse, err := contactPlugin(pluginContactRequest, "error while reseting the computer system: ")

	if err != nil {
		subTaskChan <- getResponse.StatusCode
		errMsg := err.Error()
		log.Println(errMsg)
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
		"Location":          req.TargetURI,
	}
	resp.StatusCode = getResponse.StatusCode
	percentComplete = 100
	subTaskChan <- int32(getResponse.StatusCode)
	var task = fillTaskData(subTaskID, targetURI, resp, common.Completed, common.OK, percentComplete, http.MethodPost)
	err = e.UpdateTask(task)
	if err != nil && err.Error() == common.Cancelling {
		var task = fillTaskData(subTaskID, targetURI, resp, common.Cancelled, common.Critical, percentComplete, http.MethodPost)
		err = e.UpdateTask(task)
	}
	if getResponse.StatusCode == http.StatusOK {
		agmodel.AddSystemResetInfo(req.TargetURI, req.ResetType)
	}
	return
}
