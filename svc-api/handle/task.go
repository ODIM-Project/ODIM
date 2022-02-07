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

//Package handle ...
package handle

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	errResponse "github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// TaskRPCs defines all the RPC methods in task service
type TaskRPCs struct {
	DeleteTaskRPC     func(req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error)
	GetTaskRPC        func(req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error)
	GetSubTasksRPC    func(req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error)
	GetSubTaskRPC     func(req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error)
	GetTaskMonitorRPC func(req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error)
	TaskCollectionRPC func(req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error)
	GetTaskServiceRPC func(req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error)
}

// DeleteTask deletes the task with given TaskID
// It takes iris context and extract auth token and TaskID from the context
// Create a request object in task proto request format and pass it to rpc call
func (task *TaskRPCs) DeleteTask(ctx iris.Context) {
	defer ctx.Next()
	req := &taskproto.GetTaskRequest{
		TaskID:       ctx.Params().Get("TaskID"),
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	response, err := task.DeleteTaskRPC(req)
	common.SetResponseHeader(ctx, response.Header)

	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	common.SetResponseHeader(ctx, response.Header)
	ctx.StatusCode(int(response.StatusCode))
	ctx.Write(response.Body)

	return
}

// GetTaskStatus fetches task status
// It takes iris context and extract auth token and TaskID from the context
// Create a request object in task proto request format and pass it to rpc call
func (task *TaskRPCs) GetTaskStatus(ctx iris.Context) {
	defer ctx.Next()
	req := &taskproto.GetTaskRequest{
		TaskID:       ctx.Params().Get("TaskID"),
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	response, err := task.GetTaskRPC(req)
	common.SetResponseHeader(ctx, response.Header)

	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET, DELETE")
	common.SetResponseHeader(ctx, response.Header)
	ctx.StatusCode(int(response.StatusCode))
	ctx.Write(response.Body)

	return
}

// GetSubTasks fetches sub task collection under given umbralla tasks
// It takes iris context and extract auth token and TaskID from the context
// Create a request object in task proto request format and pass it to rpc call
func (task *TaskRPCs) GetSubTasks(ctx iris.Context) {
	defer ctx.Next()
	req := &taskproto.GetTaskRequest{
		TaskID:       ctx.Params().Get("TaskID"),
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	response, err := task.GetSubTasksRPC(req)
	common.SetResponseHeader(ctx, response.Header)

	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, response.Header)
	ctx.StatusCode(int(response.StatusCode))
	ctx.Write(response.Body)

	return
}

// GetSubTask fetches sub task details
// It takes iris context and extract auth token, TaskID and subTasks from the context
// Create a request object in task proto request format and pass it to rpc call
func (task *TaskRPCs) GetSubTask(ctx iris.Context) {
	defer ctx.Next()
	req := &taskproto.GetTaskRequest{
		TaskID:       ctx.Params().Get("TaskID"),
		SubTaskID:    ctx.Params().Get("subTaskID"),
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	response, err := task.GetSubTaskRPC(req)
	common.SetResponseHeader(ctx, response.Header)

	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, response.Header)
	ctx.StatusCode(int(response.StatusCode))
	ctx.Write(response.Body)

	return
}

// GetTaskMonitor fetches task monitor
// It takes iris context and extract auth token and TaskID from the context
// Create a request object in task proto request format and pass it to rpc call
func (task *TaskRPCs) GetTaskMonitor(ctx iris.Context) {
	defer ctx.Next()
	req := &taskproto.GetTaskRequest{
		TaskID:       ctx.Params().Get("TaskID"),
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	response, err := task.GetTaskMonitorRPC(req)
	common.SetResponseHeader(ctx, response.Header)

	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, response.Header)
	ctx.StatusCode(int(response.StatusCode))
	ctx.Write(response.Body)

	return
}

//TaskCollection fetches all tasks available in DB
//It takes iris context and extract auth token from the context
//Create a request object in task proto request format and pass it to rpc call
func (task *TaskRPCs) TaskCollection(ctx iris.Context) {
	defer ctx.Next()
	req := &taskproto.GetTaskRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, errResponse.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		common.SetResponseHeader(ctx, nil)
		return
	}
	response, err := task.TaskCollectionRPC(req)
	common.SetResponseHeader(ctx, response.Header)

	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, response.Header)
	ctx.StatusCode(int(response.StatusCode))
	ctx.Write(response.Body)

	return
}

//GetTaskService fetches Task Service details
//It takes iris context and extract auth token from the context
//Create a request object in task proto request format and pass it to rpc call
func (task *TaskRPCs) GetTaskService(ctx iris.Context) {
	defer ctx.Next()
	req := &taskproto.GetTaskRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	response, err := task.GetTaskServiceRPC(req)
	common.SetResponseHeader(ctx, response.Header)

	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, response.Header)
	ctx.StatusCode(int(response.StatusCode))
	ctx.Write(response.Body)

	return
}
