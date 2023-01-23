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

// Package systems ...
package systems

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/logs"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-systems/scommon"
)

// BiosSetting structure for checking request body case
type BiosSetting struct {
	OdataContext      string      `json:"@odata.context"`
	OdataID           string      `json:"@odata.id"`
	Odatatype         string      `json:"@odata.type"`
	ID                string      `json:"Id"`
	Name              string      `json:"Name"`
	AttributeRegistry string      `json:"AttributeRegistry"`
	Attributes        interface{} `json:"Attributes"`
}

// BootOrderSettings structure for checking request body case
type BootOrderSettings struct {
	Boot Boot `json:"Boot"`
}

// Boot structure for checking request body case in BootOrderSettings
type Boot struct {
	BootOrder                    []string `json:"BootOrder"`
	BootSourceOverrideEnabled    string   `json:"BootSourceOverrideEnabled"`
	BootSourceOverrideMode       string   `json:"BootSourceOverrideMode"`
	BootSourceOverrideTarget     string   `json:"BootSourceOverrideTarget"`
	UefiTargetBootSourceOverride string   `json:"UefiTargetBootSourceOverride"`
}

// ResetComputerSystem structure for checking request body case
type ResetComputerSystem struct {
	ResetType string `json:"ResetType"`
}

// monitorTaskRequest hold values required monitorTask function
type monitorTaskRequest struct {
	taskID        string
	respBody      []byte
	serverURI     string
	requestBody   string
	getResponse   scommon.ResponseStatus
	location      string
	taskInfo      *common.TaskUpdateInfo
	pluginRequest scommon.PluginContactRequest
	resp          response.RPC
}

// UpdateTaskData update the task with the given data
func UpdateTaskData(ctx context.Context, taskData common.TaskData) error {
	var res map[string]interface{}
	if err := json.Unmarshal([]byte(taskData.TaskRequest), &res); err != nil {
		l.Log.Error(err)
	}
	reqStr := logs.MaskRequestBody(res)

	respBody, _ := json.Marshal(taskData.Response.Body)
	payLoad := &taskproto.Payload{
		HTTPHeaders:   taskData.Response.Header,
		HTTPOperation: taskData.HTTPMethod,
		JSONBody:      reqStr,
		StatusCode:    taskData.Response.StatusCode,
		TargetURI:     taskData.TargetURI,
		ResponseBody:  respBody,
	}

	err := services.UpdateTask(ctx, taskData.TaskID, taskData.TaskState, taskData.TaskStatus, taskData.PercentComplete, payLoad, time.Now())
	if err != nil && (err.Error() == common.Cancelling) {
		// We cant do anything here as the task has done it work completely, we cant reverse it.
		//Unless if we can do opposite/reverse action for delete server which is add server.
		services.UpdateTask(ctx, taskData.TaskID, common.Cancelled, taskData.TaskStatus, taskData.PercentComplete, payLoad, time.Now())
		if taskData.PercentComplete == 0 {
			return fmt.Errorf("error while starting the task: %v", err)
		}
		runtime.Goexit()
	}
	return nil
}

func fillTaskData(taskID, targetURI, request string, resp response.RPC, taskState string, taskStatus string, percentComplete int32, httpMethod string) common.TaskData {
	return common.TaskData{
		TaskID:          taskID,
		TargetURI:       targetURI,
		TaskRequest:     request,
		Response:        resp,
		TaskState:       taskState,
		TaskStatus:      taskStatus,
		PercentComplete: percentComplete,
		HTTPMethod:      httpMethod,
	}
}

func (e *PluginContact) monitorPluginTask(ctx context.Context, monitorTaskData *monitorTaskRequest) ([]byte, error) {
	// TODO : should be removed when context from svc-api is passed to this function
	for {
		var task common.TaskData
		if err := json.Unmarshal(monitorTaskData.respBody, &task); err != nil {
			errMsg := "Unable to parse the reset respone" + err.Error()
			l.LogWithFields(ctx).Warn(errMsg)
			common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, monitorTaskData.taskInfo)
			return monitorTaskData.respBody, err
		}
		var updatetask = fillTaskData(monitorTaskData.taskID, monitorTaskData.serverURI, monitorTaskData.requestBody, monitorTaskData.resp, task.TaskState, task.TaskStatus, task.PercentComplete, http.MethodPost)
		err := e.UpdateTask(ctx, updatetask)
		if err != nil && err.Error() == common.Cancelling {
			var updatetask = fillTaskData(monitorTaskData.taskID, monitorTaskData.serverURI, monitorTaskData.requestBody, monitorTaskData.resp, common.Cancelled, common.Critical, 100, http.MethodPost)
			e.UpdateTask(ctx, updatetask)
			return monitorTaskData.respBody, err
		}
		time.Sleep(time.Second * 5)
		monitorTaskData.pluginRequest.OID = monitorTaskData.location
		monitorTaskData.pluginRequest.HTTPMethodType = http.MethodGet
		monitorTaskData.respBody, _, monitorTaskData.getResponse, err = ContactPluginFunc(ctx, monitorTaskData.pluginRequest, "error while reseting the computer system: ")
		if err != nil {
			errMsg := err.Error()
			l.LogWithFields(ctx).Warn(errMsg)
			common.GeneralError(monitorTaskData.getResponse.StatusCode, monitorTaskData.getResponse.StatusMessage, errMsg, nil, monitorTaskData.taskInfo)
			return monitorTaskData.respBody, err
		}
		if monitorTaskData.getResponse.StatusCode == http.StatusOK {
			break
		}
	}
	return monitorTaskData.respBody, nil
}
