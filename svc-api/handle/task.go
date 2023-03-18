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

// Package handle ...
package handle

import (
	"context"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	errResponse "github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// TaskRPCs defines all the RPC methods in task service
type TaskRPCs struct {
	DeleteTaskRPC     func(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error)
	GetTaskRPC        func(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error)
	GetSubTasksRPC    func(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error)
	GetSubTaskRPC     func(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error)
	GetTaskMonitorRPC func(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error)
	TaskCollectionRPC func(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error)
	GetTaskServiceRPC func(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error)
}

// DeleteTask deletes the task with given TaskID
// It takes iris context and extract auth token and TaskID from the context
// Create a request object in task proto request format and pass it to rpc call
func (task *TaskRPCs) DeleteTask(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := &taskproto.GetTaskRequest{
		TaskID:       ctx.Params().Get("TaskID"),
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for deleting task with taskid %s", req.TaskID)
	response, err := task.DeleteTaskRPC(ctxt, req)
	common.SetResponseHeader(ctx, response.Header)

	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for deleting task is %s with status code %d", string(response.Body), int(response.StatusCode))
	common.SetResponseHeader(ctx, response.Header)
	ctx.StatusCode(int(response.StatusCode))
	ctx.Write(response.Body)
}

// GetTaskStatus fetches task status
// It takes iris context and extract auth token and TaskID from the context
// Create a request object in task proto request format and pass it to rpc call
func (task *TaskRPCs) GetTaskStatus(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := &taskproto.GetTaskRequest{
		TaskID:       ctx.Params().Get("TaskID"),
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting task status with taskid %s", req.TaskID)
	response, err := task.GetTaskRPC(ctxt, req)
	common.SetResponseHeader(ctx, response.Header)

	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting task is %s with status code %d", string(response.Body), int(response.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET, DELETE")
	common.SetResponseHeader(ctx, response.Header)
	ctx.StatusCode(int(response.StatusCode))
	ctx.Write(response.Body)
}

// GetSubTasks fetches sub task collection under given umbralla tasks
// It takes iris context and extract auth token and TaskID from the context
// Create a request object in task proto request format and pass it to rpc call
func (task *TaskRPCs) GetSubTasks(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := &taskproto.GetTaskRequest{
		TaskID:       ctx.Params().Get("TaskID"),
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	response, err := task.GetSubTasksRPC(ctxt, req)
	common.SetResponseHeader(ctx, response.Header)
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting sub tasks collections with taskid %s", req.TaskID)
	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting sub tasks is %s with status code %d", string(response.Body), int(response.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, response.Header)
	ctx.StatusCode(int(response.StatusCode))
	ctx.Write(response.Body)
}

// GetSubTask fetches sub task details
// It takes iris context and extract auth token, TaskID and subTasks from the context
// Create a request object in task proto request format and pass it to rpc call
func (task *TaskRPCs) GetSubTask(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := &taskproto.GetTaskRequest{
		TaskID:       ctx.Params().Get("TaskID"),
		SubTaskID:    ctx.Params().Get("subTaskID"),
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting sub task information with taskid %s and subtaskid %s", req.TaskID, req.SubTaskID)
	response, err := task.GetSubTaskRPC(ctxt, req)
	common.SetResponseHeader(ctx, response.Header)

	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting sub task information is %s with status code %d", string(response.Body), int(response.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, response.Header)
	ctx.StatusCode(int(response.StatusCode))
	ctx.Write(response.Body)
}

// GetTaskMonitor fetches task monitor
// It takes iris context and extract auth token and TaskID from the context
// Create a request object in task proto request format and pass it to rpc call
func (task *TaskRPCs) GetTaskMonitor(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := &taskproto.GetTaskRequest{
		TaskID:       ctx.Params().Get("TaskID"),
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting task mon for taskid %s", req.TaskID)
	response, err := task.GetTaskMonitorRPC(ctxt, req)
	common.SetResponseHeader(ctx, response.Header)

	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting task monitoring information is %s with status code %d", string(response.Body), int(response.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, response.Header)
	ctx.StatusCode(int(response.StatusCode))
	ctx.Write(response.Body)
}

// TaskCollection fetches all tasks available in DB
// It takes iris context and extract auth token from the context
// Create a request object in task proto request format and pass it to rpc call
func (task *TaskRPCs) TaskCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := &taskproto.GetTaskRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting task all the available tasks collection")
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, errResponse.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		common.SetResponseHeader(ctx, nil)
		return
	}
	response, err := task.TaskCollectionRPC(ctxt, req)
	common.SetResponseHeader(ctx, response.Header)

	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting task collections is %s with status code %d", string(response.Body), int(response.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, response.Header)
	ctx.StatusCode(int(response.StatusCode))
	ctx.Write(response.Body)

}

// GetTaskService fetches Task Service details
// It takes iris context and extract auth token from the context
// Create a request object in task proto request format and pass it to rpc call
func (task *TaskRPCs) GetTaskService(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := &taskproto.GetTaskRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	response, err := task.GetTaskServiceRPC(ctxt, req)
	common.SetResponseHeader(ctx, response.Header)
	l.LogWithFields(ctxt).Debug("Incoming request received for getting task service details")
	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting task service is %s with status code %d", string(response.Body), int(response.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, response.Header)
	ctx.StatusCode(int(response.StatusCode))
	ctx.Write(response.Body)
}
