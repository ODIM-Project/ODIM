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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
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
func (e *ExternalInterface) Reset(ctx context.Context, taskID string, sessionUserName string, req *aggregatorproto.AggregatorRequest) response.RPC {
	var resp response.RPC
	var percentComplete int32
	targetURI := "/redfish/v1/AggregationService/Actions/AggregationService.Reset/" // this will removed later and passed as input param in req struct
	percentComplete = 0

	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: targetURI, UpdateTask: e.UpdateTask, TaskRequest: string(req.RequestBody)}

	var resetRequest AggregationResetRequest
	if err := json.Unmarshal(req.RequestBody, &resetRequest); err != nil {
		errMsg := "Unable to validate request fields: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errMsg, nil, taskInfo)
	}
	missedProperty, err := resetRequest.validateResetRequestFields(req.RequestBody)
	if err != nil {
		errMsg := "Unable to validate request fields: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{missedProperty}, taskInfo)
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, resetRequest)
	if err != nil {
		errMsg := "Unable to validate request fields: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	} else if invalidProperties != "" {
		errorMessage := "One or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		l.LogWithFields(ctx).Error(errorMessage)
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
	threadID := 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.ResetAggregate)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	threadID++
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
						err := e.UpdateTask(ctx, task)
						if err != nil && err.Error() == common.Cancelling {
							task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Cancelled, common.OK, percentComplete, http.MethodPost)
							e.UpdateTask(ctx, task)
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
			if strings.Contains(resource, "/AggregationService/Aggregates") {
				e.aggregateSystems(ctx, resetRequest.ResetType, resource, taskID, string(req.RequestBody), subTaskChan, sessionUserName, resource, resetRequest.ResetType, &wg)
			} else {
				threadID := 1
				resetCtx := context.WithValue(ctxt, common.ThreadName, common.ResetAggregate)
				resetCtx = context.WithValue(resetCtx, common.ThreadID, strconv.Itoa(threadID))
				go e.resetSystem(resetCtx, taskID, string(req.RequestBody), subTaskChan, sessionUserName, resource, resetRequest.ResetType, &wg)
				threadID++
			}
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
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(resp.StatusCode, resp.StatusMessage, errMsg, nil, taskInfo)
	}

	l.LogWithFields(ctx).Info("All reset actions are successfully completed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID)
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	args = response.Args{
		Code:    resp.StatusMessage,
		Message: "Request completed successfully",
	}
	resp.Body = args.CreateGenericErrorResponse()
	var task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Completed, taskStatus, percentComplete, http.MethodPost)
	err = e.UpdateTask(ctx, task)
	if err != nil && err.Error() == common.Cancelling {
		task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Cancelled, common.Critical, percentComplete, http.MethodPost)
		e.UpdateTask(ctx, task)
		runtime.Goexit()
	}
	return resp
}
func (e *ExternalInterface) aggregateSystems(ctx context.Context, requestType, url, taskID, reqBody string, subTaskChan chan<- int32, sessionUserName, element, resetType string, wg *sync.WaitGroup) {
	var resp response.RPC
	var percentComplete int32
	defer wg.Done()
	subTaskURI, err := e.CreateChildTask(ctx, sessionUserName, taskID)
	if err != nil {
		subTaskChan <- http.StatusInternalServerError
		l.LogWithFields(ctx).Error("error while trying to create sub task")
		return
	}
	var subTaskID string
	strArray := strings.Split(subTaskURI, "/")
	if strings.HasSuffix(subTaskURI, "/") {
		subTaskID = strArray[len(strArray)-2]
	} else {
		subTaskID = strArray[len(strArray)-1]
	}
	taskInfo := &common.TaskUpdateInfo{TaskID: subTaskID, TargetURI: url, UpdateTask: e.UpdateTask, TaskRequest: reqBody}
	aggregate, err1 := agmodel.GetAggregate(url)
	if err1 != nil {
		percentComplete = 100
		errorMessage := err1.Error()
		l.LogWithFields(ctx).Error("error getting aggregate : " + errorMessage)
		if errors.DBKeyNotFound == err1.ErrNo() {
			subTaskChan <- http.StatusNotFound
			resp.StatusCode = http.StatusNotFound
			common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"Aggregate", url}, taskInfo)
			return
		}

		subTaskChan <- http.StatusInternalServerError
		resp.StatusCode = http.StatusInternalServerError
		common.GeneralError(http.StatusInternalServerError, response.ResourceNotFound, errorMessage, []interface{}{"Aggregate", url}, taskInfo)
		return
	}

	// subTaskChan is a buffered channel with buffer size equal to total number of elements.
	// this also helps while cancelling the task. even if the reader is not available for reading
	// the channel buffer will collect them and allows gracefull exit for already spanned goroutines.

	subTaskChan1 := make(chan int32, len(aggregate.Elements))
	resp.StatusCode = http.StatusOK
	var cancelled, partialResultFlag bool

	threadID := 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.SubTaskStatusUpdate)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	threadID++
	var wg1, writeWG sync.WaitGroup
	go func() {
		for i := 0; i < len(aggregate.Elements); i++ {
			if !cancelled { // task cancelled check to determine whether to collect status codes.
				select {
				case statusCode := <-subTaskChan1:
					if statusCode != http.StatusOK {
						partialResultFlag = true
						if resp.StatusCode < statusCode {
							resp.StatusCode = statusCode
						}
					}
					if i < len(aggregate.Elements)-1 {
						percentComplete = int32(((i + 1) / len(aggregate.Elements)) * 100)
						var task = fillTaskData(subTaskID, url, reqBody, resp, common.Running, common.OK, percentComplete, http.MethodPost)
						err := e.UpdateTask(ctx, task)
						if err != nil && err.Error() == common.Cancelling {
							task = fillTaskData(subTaskID, url, reqBody, resp, common.Cancelled, common.OK, percentComplete, http.MethodPost)
							e.UpdateTask(ctx, task)
							cancelled = true
						}
					}
				}
			}
			writeWG.Done()
		}
	}()

	for _, element := range aggregate.Elements {
		wg1.Add(1)
		writeWG.Add(1)
		threadID := 1
		resetCtxt := context.WithValue(ctxt, common.ThreadName, common.ResetSystem)
		resetCtxt = context.WithValue(resetCtxt, common.ThreadID, strconv.Itoa(threadID))
		go e.resetSystem(resetCtxt, subTaskID, reqBody, subTaskChan1, sessionUserName, element.OdataID, requestType, &wg1)
		threadID++
	}

	wg1.Wait()
	writeWG.Wait()
	subTaskChan <- int32(resp.StatusCode)
	taskStatus := common.OK
	if partialResultFlag {
		taskStatus = common.Warning
	}
	percentComplete = 100
	var args response.Args
	if resp.StatusCode != http.StatusOK {
		subTaskChan <- resp.StatusCode
		common.GeneralError(resp.StatusCode, resp.StatusMessage, "", []interface{}{"Aggregate", url}, taskInfo)
		return
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	args = response.Args{
		Code:    resp.StatusMessage,
		Message: "Request completed successfully",
	}
	resp.Body = args.CreateGenericErrorResponse()
	var task = fillTaskData(subTaskID, url, reqBody, resp, common.Completed, taskStatus, percentComplete, http.MethodPost)
	err = e.UpdateTask(ctx, task)
	if err != nil && err.Error() == common.Cancelling {
		common.GeneralError(http.StatusNotFound, common.Cancelled, "", []interface{}{"Aggregate", url}, taskInfo)
		return

	}
	return
}
