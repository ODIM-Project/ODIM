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

//Package thandle ...
package thandle

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-task/tresponse"
)

//GetTaskMonitor is an API end point to get the task details and response body.
// Takes X-Auth-Token and authorize the request.
//If X-Auth-Token is empty or invalid then it returns "StatusUnauthorized".
// If the TaskID is not found then it return "StatusNotFound".
// If the task is still not completed or cancelled or killed then it return with 202
// with empty response body, else it return with "200 OK" with full task info in the
// response body.
func (ts *TasksRPC) GetTaskMonitor(ctx context.Context, req *taskproto.GetTaskRequest, rsp *taskproto.TaskResponse) error {
	rsp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Date":              time.Now().Format(http.TimeFormat),
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	privileges := []string{common.PrivilegeLogin}
	authStatusCode, authStatusMessage := ts.AuthenticationRPC(req.SessionToken, privileges)
	if authStatusCode != http.StatusOK {
		errorMessage := "error while trying to authenticate session"
		rsp.StatusCode = authStatusCode
		rsp.StatusMessage = authStatusMessage
		rpcResp := common.GeneralError(authStatusCode, authStatusMessage, errorMessage, nil, nil)
		rsp.Body = generateResponse(rpcResp.Body)
		rsp.Header = rpcResp.Header
		log.Printf(errorMessage)
		return nil
	}
	sessionUserName, err := ts.GetSessionUserNameRPC(req.SessionToken)
	if err != nil {
		// handle the error case with appropriate response body
		errorMessage := "error while trying to authenticate session"
		rsp.StatusCode = http.StatusUnauthorized
		rsp.StatusMessage = response.NoValidSession
		rsp.Body = generateResponse(common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil).Body)
		log.Printf(errorMessage)
		return nil
	}
	log.Printf(sessionUserName)
	// get task status from database using task id
	task, err := ts.GetTaskStatusModel(req.TaskID, common.InMemory)
	if err != nil {
		log.Printf("error getting task status : %v", err)
		rsp.StatusCode = http.StatusNotFound
		rsp.StatusMessage = response.ResourceNotFound
		rsp.Body = generateResponse(common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"Task", req.TaskID}, nil).Body)
		return nil
	}

	// Check the state of the task
	if task.TaskState == "Completed" || task.TaskState == "Cancelled" || task.TaskState == "Killed" || task.TaskState == "Exception" {
		// return with the actual status code, along with response header and response body
		//Build the respose Body
		rsp.Header = task.Payload.HTTPHeaders
		rsp.Body = task.TaskResponse
		rsp.StatusCode = task.StatusCode
		// Delete the task from db as it is completed and user requested for the details.
		// return the user with task details by deleting the task from db
		// User should be careful as this is the last call to Task monitor API.
		/*
			err := task.Delete()
			if err != nil {
				log.Printf("error while deleting the task from db: %v", err)
			}
		*/
		return nil
	}
	// Construct the Task object to return as long as 202 code is being returned.

	messageList := []tresponse.Messages{}
	for _, element := range task.Messages {
		message := tresponse.Messages{
			MessageID:         element.MessageID,
			RelatedProperties: element.RelatedProperties,
			Message:           element.Message,
			MessageArgs:       element.MessageArgs,
			Severity:          element.Severity,
		}
		messageList = append(messageList, message)
	}

	commonResponse := response.Response{
		OdataType:    "#Task.v1_4_2a.Task",
		ID:           task.ID,
		Name:         task.Name,
		OdataContext: "/redfish/v1/$metadata#Task.Task",
		OdataID:      "/redfish/v1/TaskService/Tasks/" + task.ID,
	}
	rsp.StatusCode = http.StatusAccepted
	rsp.StatusMessage = response.TaskStarted
	commonResponse.MessageArgs = []string{task.ID}
	commonResponse.CreateGenericResponse(rsp.StatusMessage)

	httpHeaders := []string{}
	for key, value := range task.Payload.HTTPHeaders {
		httpHeaders = append(httpHeaders, fmt.Sprintf("%v: %v", key, value))
	}

	taskResponse := tresponse.Task{
		Response:    commonResponse,
		TaskState:   task.TaskState,
		StartTime:   task.StartTime.UTC(),
		EndTime:     task.EndTime.UTC(),
		TaskStatus:  task.TaskStatus,
		Messages:    messageList,
		TaskMonitor: task.TaskMonitor,
		Payload: tresponse.Payload{
			RespHTTPHeaders: httpHeaders,
			HTTPOperation:   task.Payload.HTTPOperation,
			RespJSONBody:    string(task.Payload.JSONBody),
			TargetURI:       task.Payload.TargetURI,
		},
		PercentComplete: task.PercentComplete,
	}
	if task.ParentID == "" {
		taskResponse.SubTasks = "/redfish/v1/TaskService/Tasks/" + task.ID + "/SubTasks"
	}
	rsp.Body = generateResponse(taskResponse)
	//    end
	// return 202
	// build response header
	// return with empty response body

	rsp.Header["location"] = task.TaskMonitor
	return nil
}
