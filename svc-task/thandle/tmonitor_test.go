//(C) Copyright [2019] Hewlett Packard Enterprise Development LP
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
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-task/tmodel"
)

func mockIsAuthorized(sessionToken string, privileges []string) response.RPC {
	switch sessionToken {
	case "validToken":
		return common.GeneralError(http.StatusOK, response.Success, "", nil, nil)
	case "NotTaskUserToken":
		// this session user does not have ConfigureUses Privilege
		for _, privilege := range privileges {
			if privilege == common.PrivilegeConfigureUsers {
				fmt.Printf("UnAuthorized %v", privileges)
				return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "error while trying to authenticate session", nil, nil)
			}
		}
		fmt.Printf("Autherized %v", privileges)
		return common.GeneralError(http.StatusOK, response.Success, "", nil, nil)
	case "NotTaskUserButAdminToken":
		return common.GeneralError(http.StatusOK, response.Success, "", nil, nil)
	default:
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "error while trying to authenticate session", nil, nil)

	}
}
func mockGetSessionUserName(sessionToken string) (string, error) {
	var user string
	switch sessionToken {
	case "validToken":
		user = "validUser"
	case "NotTaskUserToken":
		user = "NotTaskUser"
	case "NotTaskUserButAdminToken":
		user = "admin"
	default:
		return "", fmt.Errorf("invalid SessionToken")
	}
	return user, nil
}
func mockGetTaskStatusModel(taskID string, db common.DbType) (*tmodel.Task, error) {
	if db != common.InMemory {
		return nil, fmt.Errorf("Resource not found")
	}
	task := tmodel.Task{
		UserName:     "validUser",
		ParentID:     "",
		ChildTaskIDs: nil,
		ID:           "validTaskID",
		TaskState:    "New",
		TaskStatus:   "OK",
		StartTime:    time.Now(),
		EndTime:      time.Time{},
	}
	switch taskID {
	case "validTaskID":
		task.TaskState = "New"
	case "RunningTaskID":
		task.TaskState = "Running"
		task.TaskStatus = "OK"
		task.ChildTaskIDs = []string{"CompletedSubTaskID", "RunningSubTaskID"}
	case "CompletedTaskID":
		task.TaskState = "Completed"
		task.TaskStatus = "OK"
		task.ChildTaskIDs = []string{"CompletedSubTaskID"}
		task.StatusCode = http.StatusOK
		task.EndTime = time.Now()
	case "ExceptionTaskID":
		task.TaskState = "Exception"
		task.TaskStatus = "Critical"
		task.ChildTaskIDs = []string{"ExceptionSubTaskID"}
		task.StatusCode = http.StatusOK
		task.EndTime = time.Now()
	case "RunningSubTaskID":
		task.TaskState = "Running"
		task.TaskStatus = "OK"
	case "CompletedSubTaskID":
		task.TaskState = "Completed"
		task.TaskStatus = "OK"
		task.StatusCode = http.StatusOK
		task.EndTime = time.Now()
	case "ExceptionSubTaskID":
		task.TaskState = "Exception"
		task.TaskStatus = "Critical"
		task.StatusCode = http.StatusOK
		task.EndTime = time.Now()
	default:
		return nil, fmt.Errorf("Resource not found")
	}
	return &task, nil

}
func mockTransactionModel(taskID string, cb func(string) error) error {
	return nil
}
func TestTasksRPC_GetTaskMonitor(t *testing.T) {
	type args struct {
		ctx context.Context
		req *taskproto.GetTaskRequest
		rsp *taskproto.TaskResponse
	}
	tests := []struct {
		name string
		ts   *TasksRPC
		args args
		want taskproto.TaskResponse
	}{
		// TODO: Add test cases.
		{
			name: "positive case, All is well. But task is running inprogress",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				ctx: nil,
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SubTaskID:    "",
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusAccepted,
			},
		},
		{
			name: "positive case, All is well. But task is completed",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				ctx: nil,
				req: &taskproto.GetTaskRequest{
					TaskID:       "CompletedTaskID",
					SubTaskID:    "",
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Negetive test case, with Invalid Task ID",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				ctx: nil,
				req: &taskproto.GetTaskRequest{
					TaskID:       "InvalidTaskID",
					SubTaskID:    "",
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "Negetive test case With Invalid Session token",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				ctx: nil,
				req: &taskproto.GetTaskRequest{
					TaskID:       "CompletedTaskID",
					SubTaskID:    "",
					SessionToken: "InvalidToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusUnauthorized,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if rsp, err := tt.ts.GetTaskMonitor(tt.args.ctx, tt.args.req); err != nil || !reflect.DeepEqual(rsp.StatusCode, tt.want.StatusCode) {
				got := tt.args.rsp
				t.Errorf("TasksRPC.GetTaskMonitor() got = %v, want %v", got.StatusCode, tt.want.StatusCode)
			}
		})
	}
}
