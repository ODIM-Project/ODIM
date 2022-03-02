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
	log "github.com/sirupsen/logrus"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// AggregationResetRequest struct for reset the BMC
type AggregationResetRequest struct {
	BatchSize                    int      `json:"BatchSize"`
	DelayBetweenBatchesInSeconds int      `json:"DelayBetweenBatchesInSeconds"`
	ResetType                    string   `json:"ResetType"`
	TargetURIs                   []string `json:"TargetURIs"`
}

// validateRequestFields validate each field in the request against default value of field type
func (validateReq AggregationResetRequest) validateResetRequestFields(reqBody []byte) (string, error) {
	if isEmptyRequest(reqBody) {
		return "ResetRequest", fmt.Errorf("ResetRequest is empty")
	}
	if validateReq.ResetType == "" {
		return "ResetType", fmt.Errorf("property ResetType missing in the reset request")
	}
	if len(validateReq.TargetURIs) == 0 {
		return "TargetURIs", fmt.Errorf("property TargetURIs missing/no resources found in the reset request")
	}
	return "", nil
}

func isEmptyRequest(requestBody []byte) bool {
	var updateRequest map[string]interface{}
	json.Unmarshal(requestBody, &updateRequest)
	if len(updateRequest) <= 0 {
		return true
	}
	return false
}

// Reset is for reseting the computer systems mentioned in the request body
func (e *ExternalInterface) Reset(taskID string, sessionUserName string, req *aggregatorproto.AggregatorRequest) response.RPC {
	var resp response.RPC
	var percentComplete int32
	targetURI := "/redfish/v1/AggregationService/Actions/AggregationService.Reset/" // this will removed later and passed as input param in req struct
	percentComplete = 0

	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: targetURI, UpdateTask: e.UpdateTask, TaskRequest: string(req.RequestBody)}

	var resetRequest AggregationResetRequest
	if err := json.Unmarshal(req.RequestBody, &resetRequest); err != nil {
		errMsg := "Unable to validate request fields: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errMsg, nil, taskInfo)
	}
	missedProperty, err := resetRequest.validateResetRequestFields(req.RequestBody)
	if err != nil {
		errMsg := "Unable to validate request fields: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{missedProperty}, taskInfo)
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, resetRequest)
	if err != nil {
		errMsg := "Unable to validate request fields: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	} else if invalidProperties != "" {
		errorMessage := "One or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Error(errorMessage)
		resp := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, taskInfo)
		return resp
	}
	// subTaskChan is a buffered channel with buffer size equal to total number of resources.
	// this also helps while cancelling the task. even if the reader is not available for reading
	// the channel buffer will collect them and allows gracefull exit for already spanned goroutines.
	subTaskChan := make(chan int32, len(resetRequest.TargetURIs))
	resp.StatusCode = http.StatusOK
	var cancelled, partialResultFlag bool
	var wg, writeWG sync.WaitGroup
	go func() {
		for i := 0; i < len(resetRequest.TargetURIs); i++ {
			if cancelled == false { // task cancelled check to determine whether to collect status codes.
				select {
				case statusCode := <-subTaskChan:
					if statusCode != http.StatusOK {
						partialResultFlag = true
						if resp.StatusCode < statusCode {
							resp.StatusCode = statusCode
						}
					}

					if i < len(resetRequest.TargetURIs)-1 {
						percentComplete = int32(((i + 1) / len(resetRequest.TargetURIs)) * 100)
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
	for _, resource := range resetRequest.TargetURIs {
		wg.Add(1)
		writeWG.Add(1)
		// tempIndex is for checking batch size, its increment on each iteration
		// if its equal to batch size then reinitilise.
		// if batch size is 0 then reset all the systems without any kind of batch and ignore the DelayBetweenBatchesInSeconds
		tempIndex = tempIndex + 1
		if resetRequest.BatchSize == 0 || tempIndex <= resetRequest.BatchSize {
			go e.resetSystem(taskID, string(req.RequestBody), subTaskChan, sessionUserName, resource, resetRequest.ResetType, &wg)
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

	log.Info("All reset actions are successfully completed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID)
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
