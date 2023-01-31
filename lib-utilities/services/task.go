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

// Package services ...
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/golang/protobuf/ptypes"
)

//CreateTask function is to contact the svc-task through the rpc call
func CreateTask(ctx context.Context, sessionUserName string) (string, error) {
	conn, errConn := ODIMService.Client(Tasks)
	if errConn != nil {
		return "", fmt.Errorf("Failed to create client connection: %s", errConn.Error())
	}
	defer conn.Close()
	taskService := taskproto.NewGetTaskServiceClient(conn)
	reqCtx := CreateNewRequestContext(ctx)
	reqCtx = common.CreateMetadata(reqCtx)
	response, err := taskService.CreateTask(
		reqCtx, &taskproto.CreateTaskRequest{
			UserName: sessionUserName,
		},
	)
	if err != nil && response == nil {
		return "", fmt.Errorf("rpc error while creating the task: %s", err.Error())
	}
	return response.TaskURI, err
}

// CreateChildTask function is to contact the svc-task through the rpc call
func CreateChildTask(ctx context.Context, sessionUserName string, parentTaskID string) (string, error) {
	conn, errConn := ODIMService.Client(Tasks)
	if errConn != nil {
		return "", fmt.Errorf("Failed to create client connection: %s", errConn.Error())
	}
	defer conn.Close()
	taskService := taskproto.NewGetTaskServiceClient(conn)
	reqCtx := CreateNewRequestContext(ctx)
	reqCtx = common.CreateMetadata(reqCtx)
	response, err := taskService.CreateChildTask(
		reqCtx, &taskproto.CreateTaskRequest{
			UserName:     sessionUserName,
			ParentTaskID: parentTaskID,
		},
	)
	if err != nil && response == nil {
		return "", fmt.Errorf("rpc error while creating the child task:  %s", err.Error())
	}
	return response.TaskURI, err
}

//UpdateTask function is to contact the svc-task through the rpc call
func UpdateTask(ctx context.Context, taskID string, taskState string, taskStatus string,
	percentComplete int32, payLoad *taskproto.Payload, endTime time.Time) error {

	tspb, err := ptypes.TimestampProto(endTime)
	if err != nil {
		return fmt.Errorf("Failed to convert the time to protobuff timestamp: %s", err.Error())
	}
	conn, errConn := ODIMService.Client(Tasks)
	if errConn != nil {
		return fmt.Errorf("Failed to create client connection: %s", errConn.Error())
	}
	defer conn.Close()
	reqCtx := CreateNewRequestContext(ctx)
	reqCtx = common.CreateMetadata(reqCtx)
	taskService := taskproto.NewGetTaskServiceClient(conn)
	_, err = taskService.UpdateTask(
		reqCtx, &taskproto.UpdateTaskRequest{
			TaskID:          taskID,
			TaskState:       taskState,
			TaskStatus:      taskStatus,
			PercentComplete: percentComplete,
			PayLoad:         payLoad,
			EndTime:         tspb,
		},
	)
	return err
}

func CreateNewRequestContext(ctx context.Context) context.Context {
	reqCtx := context.Background()
	processName, _ := ctx.Value(common.ProcessName).(string)
	transactionID, _ := ctx.Value(common.TransactionID).(string)
	actionID, _ := ctx.Value(common.ActionID).(string)
	actionName, _ := ctx.Value(common.ActionName).(string)
	threadID, _ := ctx.Value(common.ThreadID).(string)
	threadName, _ := ctx.Value(common.ThreadName).(string)
	reqCtx = context.WithValue(reqCtx, common.ProcessName, processName)
	reqCtx = context.WithValue(reqCtx, common.TransactionID, transactionID)
	reqCtx = context.WithValue(reqCtx, common.ActionID, actionID)
	reqCtx = context.WithValue(reqCtx, common.ActionName, actionName)
	reqCtx = context.WithValue(reqCtx, common.ThreadID, threadID)
	reqCtx = context.WithValue(reqCtx, common.ThreadName, threadName)
	return reqCtx
}
