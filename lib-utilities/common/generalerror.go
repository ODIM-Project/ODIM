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

package common

import (
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// TaskData holds the response data which is been put into the task
// for task updation
type TaskData struct {
	TaskID          string
	TargetURI       string
	Response        response.RPC
	TaskRequest     string
	TaskState       string
	TaskStatus      string
	PercentComplete int32
	HTTPMethod      string
}

// TaskUpdateInfo holds the info for updating a task during error response
type TaskUpdateInfo struct {
	TaskID      string
	TargetURI   string
	TaskRequest string
	UpdateTask  func(TaskData) error
}

// GeneralError will create the error response and update task if required
// This function can be used only if the expected response have only
// one extended info object. Error code for the response will be GeneralError
// If there is no requirement of task updation pass a nil value for *TaskUpdateInfo
func GeneralError(statusCode int32, statusMsg, errMsg string, msgArgs []interface{}, t *TaskUpdateInfo) response.RPC {
	var resp response.RPC
	resp.StatusCode = statusCode
	resp.StatusMessage = statusMsg
	args := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: resp.StatusMessage,
				ErrorMessage:  errMsg,
				MessageArgs:   msgArgs,
			},
		},
	}
	resp.Body = args.CreateGenericErrorResponse()
	if t != nil && t.TaskID != "" && t.TargetURI != "" && t.UpdateTask != nil {
		task := TaskData{
			TaskID:          t.TaskID,
			TargetURI:       t.TargetURI,
			Response:        resp,
			TaskRequest:     t.TaskRequest,
			TaskState:       Exception,
			TaskStatus:      Critical,
			PercentComplete: 100,
			HTTPMethod:      http.MethodPost,
		}
		t.UpdateTask(task)
	}
	return resp
}
