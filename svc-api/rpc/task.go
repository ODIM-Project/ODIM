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

// Package rpc ...
package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

var (
	NewGetTaskServiceClientFunc = taskproto.NewGetTaskServiceClient
)

// DeleteTaskRequest will do the rpc calls for the svc-task DeleteTask
func DeleteTaskRequest(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, connErr := ClientFunc(services.Tasks)
	if connErr != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", connErr)
	}

	asService := NewGetTaskServiceClientFunc(conn)
	// Call the DeleteTask
	rsp, err := asService.DeleteTask(ctx, req)
	if err != nil {
		resp := common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		body, _ := json.Marshal(resp.Body)
		rsp = &taskproto.TaskResponse{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: response.InternalError,
			Body:          body,
		}
		return rsp, fmt.Errorf("error while trying to make DeleteTask rpc call: %v", err)
	}
	defer conn.Close()
	return rsp, nil
}

// GetTaskRequest will do the rpc calls for the svc-task GetTaskStatus
func GetTaskRequest(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, connErr := ClientFunc(services.Tasks)
	if connErr != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", connErr)
	}

	asService := NewGetTaskServiceClientFunc(conn)
	// Call the GetTasks
	rsp, err := asService.GetTasks(ctx, req)
	if err != nil {
		resp := common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		body, _ := json.Marshal(resp.Body)
		rsp = &taskproto.TaskResponse{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: response.InternalError,
			Body:          body,
		}
		return rsp, fmt.Errorf("error while trying to make GetTasks rpc call: %v", err)
	}
	defer conn.Close()
	return rsp, nil
}

// GetSubTasks will do the rpc calls for the svc-task GetSubTasks
func GetSubTasks(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, connErr := ClientFunc(services.Tasks)
	if connErr != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", connErr)
	}

	tService := NewGetTaskServiceClientFunc(conn)
	// Call the GetSubTasks
	rsp, err := tService.GetSubTasks(ctx, req)
	if err != nil {
		resp := common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		body, _ := json.Marshal(resp.Body)
		rsp = &taskproto.TaskResponse{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: response.InternalError,
			Body:          body,
		}
		return rsp, fmt.Errorf("error while trying to make GetSubTasks rpc call: %v", err)
	}
	defer conn.Close()
	return rsp, nil
}

// GetSubTask will do the rpc calls for the svc-task GetSubTask
func GetSubTask(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, connErr := ClientFunc(services.Tasks)
	if connErr != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", connErr)
	}

	tService := NewGetTaskServiceClientFunc(conn)
	// Call the GetSubTask
	rsp, err := tService.GetSubTask(ctx, req)
	if err != nil {
		resp := common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		body, _ := json.Marshal(resp.Body)
		rsp = &taskproto.TaskResponse{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: response.InternalError,
			Body:          body,
		}
		return rsp, fmt.Errorf("error while trying to make GetSubTask rpc call: %v", err)
	}
	defer conn.Close()
	return rsp, nil
}

// GetTaskMonitor will do the rpc calls for the svc-task GetTaskMonitor
func GetTaskMonitor(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, connErr := ClientFunc(services.Tasks)
	if connErr != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", connErr)
	}

	tService := NewGetTaskServiceClientFunc(conn)
	// perform rpc call to svc-task to get TaskMonitor resource
	rsp, err := tService.GetTaskMonitor(ctx, req)
	if err != nil {
		resp := common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		body, _ := json.Marshal(resp.Body)
		rsp = &taskproto.TaskResponse{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: response.InternalError,
			Body:          body,
		}
		return rsp, fmt.Errorf("error while trying to make GetTaskMonitor rpc call: %v", err)
	}
	defer conn.Close()
	return rsp, nil
}

// TaskCollection will perform the rpc call to svc-task TaskCollection
func TaskCollection(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, connErr := ClientFunc(services.Tasks)
	if connErr != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", connErr)
	}

	tService := NewGetTaskServiceClientFunc(conn)
	// perform rpc call to svc-task to get TaskCollection resource
	rsp, err := tService.TaskCollection(ctx, req)
	if err != nil {
		resp := common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		body, _ := json.Marshal(resp.Body)
		rsp = &taskproto.TaskResponse{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: response.InternalError,
			Body:          body,
		}
		return rsp, fmt.Errorf("error while trying to make TaskCollection rpc call: %v", err)
	}
	defer conn.Close()
	return rsp, nil
}

// GetTaskService will perform the rpc call to svc-task GetTaskService
func GetTaskService(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, connErr := ClientFunc(services.Tasks)
	if connErr != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", connErr)
	}

	tService := NewGetTaskServiceClientFunc(conn)
	// perform rpc call to svc-task to get TaskService resource
	rsp, err := tService.GetTaskService(ctx, req)
	if err != nil {
		resp := common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		body, _ := json.Marshal(resp.Body)
		rsp = &taskproto.TaskResponse{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: response.InternalError,
			Body:          body,
		}
		return rsp, fmt.Errorf("error while trying to make GetTaskService rpc call: %v", err)
	}
	defer conn.Close()
	return rsp, nil
}
