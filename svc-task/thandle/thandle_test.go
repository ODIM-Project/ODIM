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
package thandle

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/ODIM-Project/ODIM/svc-task/tcommon"
	"github.com/ODIM-Project/ODIM/svc-task/tmodel"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/crypto/sha3"
)

const errorCollectingData string = "error while trying to collect data: "

type user struct {
	UserName string `json:"UserName"`
	Password string `json:"Password"`
	RoleID   string `json:"RoleId"`
}

func createMockUser(username, roleID string) error {
	hash := sha3.New512()
	hash.Write([]byte("P@$$w0rd"))
	hashSum := hash.Sum(nil)
	hashedPassword := base64.URLEncoding.EncodeToString(hashSum)
	user := user{
		UserName: username,
		Password: hashedPassword,
		RoleID:   roleID,
	}
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	//Create a header for data entry
	const table string = "User"
	//Save data into Database
	if err = conn.Create(table, user.UserName, user); err != nil {
		return err
	}
	return nil
}

func mockGetTaskStatus(taskID string, db common.DbType) (*tmodel.Task, error) {
	var task tmodel.Task
	err := fmt.Errorf("invalid System ID")
	if taskID != "validUUID" {
		return &task, err
	}
	return &task, nil
}

func mockCreateTaskUtil(ctx context.Context, userName string) (string, error) {
	if userName == "" {
		return "", fmt.Errorf("error invalid input argument for userName")
	}
	if userName != "validUser" {
		return "", fmt.Errorf("error invalid user: ")
	}
	return "/redfish/v1/TaskService/Tasks/validTaskID", nil
}

func mockDeleteTaskFromDBModel(ctx context.Context, task *tmodel.Task) error {

	return nil
}

func mockUpdateTaskStatusModel(task *tmodel.Task) {
}

func mockPublishToMessageBus(ctx context.Context, taskURI, taskEvenMessageID, eventType, taskMessage string) {

}
func mockValidateTaskUserNameModel(ctx context.Context, userName string) error {
	if userName != "validUser" {
		return fmt.Errorf("error while trying to read from DB: %v", errors.PackError(errors.DBKeyNotFound, "no data with the with key ", userName, " found").Error())
	}
	return nil
}
func mockPersistTaskModel(ctx context.Context, task *tmodel.Task, db common.DbType) error {
	if db != common.InMemory {
		return fmt.Errorf("error while trying to connecting to DB: error invalid db type selection")
	}
	return nil
}
func TestTasksRPC_GetTasks(t *testing.T) {
	type args struct {
		req *taskproto.GetTaskRequest
		rsp *taskproto.TaskResponse
	}
	tests := []struct {
		name string
		ts   *TasksRPC
		args args
		want taskproto.TaskResponse
	}{
		{
			name: "Positive Case: All is well, RunningTaskID, valid Token",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusAccepted,
			},
		},
		{
			name: "Positive Case: All is well, CompletedTaskID, valid Token",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "CompletedTaskID",
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Positive Case: All is well, ExceptionTaskID, valid Token",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "ExceptionTaskID",
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Negative Case:InvalidToken",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "validTaskID",
					SessionToken: "invalidToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "Negative case: invalidTaskID",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "invalidTaskID",
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "Negative case: Not Task user token",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SessionToken: "NotTaskUserToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "Negative case: Not Task user token, but admin's token",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SessionToken: "NotTaskUserButAdminToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusAccepted,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsp, err := tt.ts.GetTasks(mockContext(), tt.args.req)
			if err != nil || !reflect.DeepEqual(rsp.StatusCode, tt.want.StatusCode) {
				got := tt.args.rsp
				t.Errorf("TasksRPC.GetTasks() got = %v, want: %v", got.StatusCode, tt.want.StatusCode)
			}
		})
	}
}
func mockGetAllTaskKeysModel(ctx context.Context) ([]string, error) {
	keys := []string{"task:key1", "task:key2"}
	return keys, nil
}
func TestTasksRPC_TaskCollection(t *testing.T) {
	type args struct {
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
			name: "Positive test case, all is well.",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetAllTaskKeysModel:   mockGetAllTaskKeysModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Negative test case, Invalid session token.",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetAllTaskKeysModel:   mockGetAllTaskKeysModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
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
			if rsp, err := tt.ts.TaskCollection(mockContext(), tt.args.req); err != nil || !reflect.DeepEqual(rsp.StatusCode, tt.want.StatusCode) {
				got := tt.args.rsp
				t.Errorf("TasksRPC.TaskCollection() got = %v, want %v", got.StatusCode, tt.want.StatusCode)
			}
		})
	}
}

func TestTasksRPC_GetTaskService(t *testing.T) {
	type args struct {
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
			name: "Positive test case, all is well.",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Negative test case, Invalid session token.",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
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
			if rsp, err := tt.ts.GetTaskService(mockContext(), tt.args.req); err != nil || !reflect.DeepEqual(rsp.StatusCode, tt.want.StatusCode) {
				got := tt.args.rsp
				t.Errorf("TasksRPC.GetTaskService() got = %v, want %v", got.StatusCode, tt.want.StatusCode)
			}
		})
	}
}

func TestTasksRPC_GetSubTasks(t *testing.T) {
	type args struct {
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
			name: "Positive test case, all is well.",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Negative test case, Invalid session token.",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SessionToken: "InvalidToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "Negative test case, Invalid TaskID",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "InvalidTaskID",
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "Negative test case, with not task user's Session token",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SessionToken: "NotTaskUserToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "Negative test case. With not a task user's session token, but with Admin user session token",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SessionToken: "NotTaskUserButAdminToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if rsp, err := tt.ts.GetSubTasks(mockContext(), tt.args.req); err != nil || !reflect.DeepEqual(rsp.StatusCode, tt.want.StatusCode) {
				got := tt.args.rsp
				t.Errorf("TasksRPC.GetSubTasks() got = %v, want %v", got.StatusCode, tt.want.StatusCode)
			}
		})
	}
}

func TestTasksRPC_GetSubTask(t *testing.T) {
	type args struct {
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
			name: "Positive test case, all is well. Running SubTask",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SubTaskID:    "RunningSubTaskID",
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusAccepted,
			},
		},
		{
			name: "Positive test case, all is well. SubTask is Completed",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SubTaskID:    "CompletedSubTaskID",
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Negative test case, Invalid session token.",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SubTaskID:    "RunningSubTaskID",
					SessionToken: "InvalidToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "Negative test case, Invalid TaskID",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "InvalidTaskID",
					SubTaskID:    "InvalidSubTaskID",
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "Negative test case, with not task user's Session token",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SubTaskID:    "RunningSubTaskID",
					SessionToken: "NotTaskUserToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "Negative test case. With not a task user's session token, but with Admin user session token",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SubTaskID:    "RunningSubTaskID",
					SessionToken: "NotTaskUserButAdminToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusAccepted,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if rsp, err := tt.ts.GetSubTask(mockContext(), tt.args.req); err != nil || !reflect.DeepEqual(rsp.StatusCode, tt.want.StatusCode) {
				got := tt.args.rsp
				t.Errorf("TasksRPC.GetSubTask() got = %v, want %v", got.StatusCode, tt.want.StatusCode)
			}
		})
	}
}

func TestTasksRPC_DeleteTask(t *testing.T) {
	type args struct {
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
			name: "Positive test case, all is well. Running Task",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
				TransactionModel:      mockTransactionModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusAccepted,
			},
		},
		{
			name: "Positive test case, all is well. Task is Completed",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
				TransactionModel:      mockTransactionModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusAccepted,
			},
		},
		{
			name: "Negative test case, Invalid session token.",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
				TransactionModel:      mockTransactionModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SessionToken: "InvalidToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "Negative test case, Invalid TaskID",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
				TransactionModel:      mockTransactionModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "InvalidTaskID",
					SessionToken: "validToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name: "Negative test case, with not task user's Session token",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
				TransactionModel:      mockTransactionModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SessionToken: "NotTaskUserToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "Negative test case. With not a task user's session token, but with Admin user session token",
			ts: &TasksRPC{
				AuthenticationRPC:     mockIsAuthorized,
				GetSessionUserNameRPC: mockGetSessionUserName,
				GetTaskStatusModel:    mockGetTaskStatusModel,
				TransactionModel:      mockTransactionModel,
			},
			args: args{
				req: &taskproto.GetTaskRequest{
					TaskID:       "RunningTaskID",
					SessionToken: "NotTaskUserButAdminToken",
				},
				rsp: &taskproto.TaskResponse{},
			},
			want: taskproto.TaskResponse{
				StatusCode: http.StatusAccepted,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if rsp, err := tt.ts.DeleteTask(mockContext(), tt.args.req); err != nil || !reflect.DeepEqual(rsp.StatusCode, tt.want.StatusCode) {
				got := tt.args.rsp
				t.Errorf("TasksRPC.DeleteTask() got = %v, want %v", got.StatusCode, tt.want.StatusCode)
			}
		})
	}
}
func TestTasksRPC_CreateTask(t *testing.T) {
	type args struct {
		req *taskproto.CreateTaskRequest
		rsp *taskproto.CreateTaskResponse
	}
	tests := []struct {
		name      string
		ts        *TasksRPC
		args      args
		want      taskproto.CreateTaskResponse
		wantError error
	}{
		// TODO: Add test cases.
		{
			name: "Positive test case, All is well",
			ts: &TasksRPC{
				CreateTaskUtilHelper: mockCreateTaskUtil,
			},
			args: args{
				req: &taskproto.CreateTaskRequest{
					UserName: "validUser",
				},
				rsp: &taskproto.CreateTaskResponse{},
			},
			want: taskproto.CreateTaskResponse{
				TaskURI: "/redfish/v1/TaskService/Tasks/validTaskID",
			},
			wantError: nil,
		},
		{
			name: "Negative test case, userName is empty",
			ts: &TasksRPC{
				CreateTaskUtilHelper: mockCreateTaskUtil,
			},
			args: args{
				req: &taskproto.CreateTaskRequest{
					UserName: "",
				},
				rsp: &taskproto.CreateTaskResponse{},
			},
			want: taskproto.CreateTaskResponse{
				TaskURI: "",
			},
			wantError: fmt.Errorf("error invalid input argument for userName"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if rsp, err := tt.ts.CreateTask(mockContext(), tt.args.req); !reflect.DeepEqual(rsp.TaskURI, tt.want.TaskURI) {
				t.Errorf("TasksRPC.CreateTask() got error = %v, wantError %v", err, tt.wantError)
				t.Errorf("TasksRPC.CreateTask() got = %v, want %v", rsp, tt.want)
			}
		})
	}
}

func TestTasksRPC_CreateChildTask(t *testing.T) {
	type args struct {
		req *taskproto.CreateTaskRequest
		rsp *taskproto.CreateTaskResponse
	}
	tests := []struct {
		name    string
		ts      *TasksRPC
		args    args
		wantErr error
		wantRsp taskproto.CreateTaskResponse
	}{
		// TODO: Add test cases.
		{
			name: "Positive case: All is well",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
			},
			args: args{
				req: &taskproto.CreateTaskRequest{
					UserName:     "validUser",
					ParentTaskID: "validTaskID",
				},
				rsp: &taskproto.CreateTaskResponse{},
			},
			wantErr: nil,
			wantRsp: taskproto.CreateTaskResponse{
				TaskURI: "/redfish/v1/TaskService/Tasks/validTaskID",
			},
		},
		{
			name: "Negative case: ParentTaskID is Empty",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
			},
			args: args{
				req: &taskproto.CreateTaskRequest{
					UserName:     "validUser",
					ParentTaskID: "",
				},
				rsp: &taskproto.CreateTaskResponse{},
			},
			wantErr: fmt.Errorf("error empty/invalid input Parent Task ID"),
			wantRsp: taskproto.CreateTaskResponse{
				TaskURI: "",
			},
		},
		{
			name: "Negative case: Invalid ParentTaskID",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
			},
			args: args{
				req: &taskproto.CreateTaskRequest{
					UserName:     "validUser",
					ParentTaskID: "InvalidTaskID",
				},
				rsp: &taskproto.CreateTaskResponse{},
			},
			wantErr: fmt.Errorf("error while retrieving the task details from DB: Resource not found"),
			wantRsp: taskproto.CreateTaskResponse{
				TaskURI: "",
			},
		},
		{
			name: "Negative case: Invalid UserName",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
			},
			args: args{
				req: &taskproto.CreateTaskRequest{
					UserName:     "InvalidUser",
					ParentTaskID: "validTaskID",
				},
				rsp: &taskproto.CreateTaskResponse{},
			},
			wantErr: nil,
			wantRsp: taskproto.CreateTaskResponse{
				TaskURI: "/redfish/v1/TaskService/Tasks/validTaskID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if rsp, err := tt.ts.CreateChildTask(mockContext(), tt.args.req); !reflect.DeepEqual(err, tt.wantErr) || !reflect.DeepEqual(rsp.TaskURI, tt.wantRsp.TaskURI) {
				t.Errorf("TasksRPC.CreateChildTask() got error = %v, wantErr: %v", err, tt.wantErr)
				t.Errorf("TasksRPC.CreateChildTask() got response = %v, want: %v", rsp, tt.wantRsp)
			}
		})
	}
}
func TestTasksRPC_UpdateTask(t *testing.T) {
	TaskCollection = TaskCollectionData{
		TaskCollection: make(map[string]int32),
		Lock:           sync.Mutex{},
	}
	type args struct {
		req *taskproto.UpdateTaskRequest
		rsp *taskproto.UpdateTaskResponse
	}
	tests := []struct {
		name    string
		ts      *TasksRPC
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "Positive case: All is well with task state as Completed ",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Completed",
					TaskStatus:      "OK",
					PercentComplete: 50,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: nil,
		},
		{
			name: "Positive case: All is well with task state as Killed, status as Critical ",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Killed",
					TaskStatus:      "Critical",
					PercentComplete: 50,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: nil,
		},
		{
			name: "Negative case: task state as Killed, status as Critical, end time as empty ",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Killed",
					TaskStatus:      "Critical",
					PercentComplete: 50,
					PayLoad:         nil,
					EndTime: func() *timestamp.Timestamp {
						t, _ := ptypes.TimestampProto(time.Time{})
						return t
					}(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: fmt.Errorf("error invalid end time for the task"),
		},
		{
			name: "Negative case: task state as Killed, status as Invalid",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Killed",
					TaskStatus:      "Invalid",
					PercentComplete: 50,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: fmt.Errorf("error invalid taskStatus provided as input argument"),
		},
		{
			name: "Positive case: All is well with task state as Cancelled, status as     Critical ",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Cancelled",
					TaskStatus:      "Critical",
					PercentComplete: 50,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: nil,
		},
		{
			name: "Positive case: All is well with task state as Exception, status     as Critical ",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Exception",
					TaskStatus:      "Critical",
					PercentComplete: 50,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: nil,
		},
		{
			name: "Negative case: Invalid Status for Exception state task",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Exception",
					TaskStatus:      "Invalid",
					PercentComplete: 50,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: fmt.Errorf("error invalid taskStatus provided as input argument"),
		},

		{
			name: "Positive case: All is well with task state as Cancelling, status as Ok",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Cancelling",
					TaskStatus:      "OK",
					PercentComplete: 50,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: nil,
		},
		{
			name: "Positive case: All is well with task state as Interrupted, status as Ok",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Interrupted",
					TaskStatus:      "OK",
					PercentComplete: 50,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: nil,
		},
		{
			name: "Positive case: All is well with task state as New, status   as Ok",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "New",
					TaskStatus:      "OK",
					PercentComplete: 0,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: nil,
		},
		{
			name: "Positive case: All is well with task state as Pending, status   as Ok",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Pending",
					TaskStatus:      "OK",
					PercentComplete: 0,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: nil,
		},
		{
			name: "Positive case: All is well with task state as Running, status   as  Ok",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Running",
					TaskStatus:      "OK",
					PercentComplete: 0,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: nil,
		},
		{
			name: "Positive case: All is well with task state as Service, status as Ok",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Service",
					TaskStatus:      "OK",
					PercentComplete: 0,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: nil,
		},
		{
			name: "Positive case: All is well with task state as Starting, status as Ok",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Starting",
					TaskStatus:      "OK",
					PercentComplete: 29,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: nil,
		},
		{
			name: "Positive case: All is well with task state as Stopping, status   as  Ok",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Stopping",
					TaskStatus:      "OK",
					PercentComplete: 0,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: nil,
		},
		{
			name: "Positive case: All is well with task state as Suspended, status as Ok",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Suspended",
					TaskStatus:      "OK",
					PercentComplete: 30,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: nil,
		},
		{
			name: "Negative case: All is well with task state as InvalidState,status   as Ok",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "InvalidState",
					TaskStatus:      "OK",
					PercentComplete: 30,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: fmt.Errorf("error invalid input argument for taskState"),
		},
		{
			name: "Positive case: All is well with task state as Completed,status as Ok, with payload",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Completed",
					TaskStatus:      "OK",
					PercentComplete: 30,
					PayLoad: &taskproto.Payload{
						HTTPHeaders: map[string]string{
							"Content-Type": "application/json",
						},
						HTTPOperation: "POST",
						JSONBody:      "",
						StatusCode:    201,
						TargetURI:     "/redfish/v1/AggregationService/Actions/AggregationService.Add",
					},
					EndTime: ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: nil,
		},
		{
			name: "Nagative case: State as Completed, status is Invalid, with payload",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Completed",
					TaskStatus:      "Invalid",
					PercentComplete: 30,
					PayLoad: &taskproto.Payload{
						HTTPHeaders: map[string]string{
							"Content-Type": "application/json",
						},
						HTTPOperation: "POST",
						JSONBody:      "",
						StatusCode:    201,
						TargetURI:     "/redfish/v1/AggregationService/Actions/AggregationService.Add",
					},
					EndTime: ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: fmt.Errorf("error invalid taskStatus provided as input argument"),
		},
		{
			name: "Nagative case: State as Completed, status is OK,but end time is null with payload",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Completed",
					TaskStatus:      "Critical",
					PercentComplete: 30,
					PayLoad: &taskproto.Payload{
						HTTPHeaders: map[string]string{
							"Content-Type": "application/json",
						},
						HTTPOperation: "POST",
						JSONBody:      "",
						StatusCode:    201,
						TargetURI:     "/redfish/v1/AggregationService/Actions/AggregationService.Add",
					},
					EndTime: func() *timestamp.Timestamp {
						t, _ := ptypes.TimestampProto(time.Time{})
						return t
					}(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: fmt.Errorf("error invalid end time for the task"),
		},
		{
			name: "Negative case: task state as Exception,status as Critical with payload, endTime as empty",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Exception",
					TaskStatus:      "Critical",
					PercentComplete: 30,
					PayLoad: &taskproto.Payload{
						HTTPHeaders: map[string]string{
							"Content-Type": "application/json",
						},
						HTTPOperation: "POST",
						JSONBody:      "",
						StatusCode:    500,
						TargetURI:     "/redfish/v1/AggregationService/Actions/AggregationService.Add",
					},
					EndTime: func() *timestamp.Timestamp {
						t, _ := ptypes.TimestampProto(time.Time{})
						return t
					}(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: fmt.Errorf("error invalid end time for the task"),
		},
		{
			name: "Negative case: state as Exception,status as   Critical with payload",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Exception",
					TaskStatus:      "Critical",
					PercentComplete: 30,
					PayLoad: &taskproto.Payload{
						HTTPHeaders: map[string]string{
							"Content-Type": "application/json",
						},
						HTTPOperation: "POST",
						JSONBody:      "",
						StatusCode:    500,
						TargetURI:     "/redfish/v1/AggregationService/    Actions/AggregationService.Add",
					},
					EndTime: ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: nil,
		},
		{
			name: "Negative case: All is well with task state as cancelled,status as OK, with payload",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Cancelled",
					TaskStatus:      "OK",
					PercentComplete: 30,
					PayLoad: &taskproto.Payload{
						HTTPHeaders: map[string]string{
							"Content-Type": "application/json",
						},
						HTTPOperation: "POST",
						JSONBody:      "",
						StatusCode:    500,
						TargetURI:     "/redfish/v1/AggregationService/Actions/AggregationService.Add",
					},
					EndTime: ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: fmt.Errorf("error invalid taskStatus provided as input argument"),
		},
		{
			name: "Negative case: All is well with task state as cancelled,status as Critical, with payload",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Cancelled",
					TaskStatus:      "Critical",
					PercentComplete: 30,
					PayLoad: &taskproto.Payload{
						HTTPHeaders: map[string]string{
							"Content-Type": "application/json",
						},
						HTTPOperation: "POST",
						JSONBody:      "",
						StatusCode:    500,
						TargetURI:     "/redfish/v1/AggregationService/Actions/AggregationService.Add",
					},
					EndTime: ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: nil,
		},
		{
			name: "Negative case: Task state as cancelled,status as Invalid",
			ts: &TasksRPC{
				GetTaskStatusModel:   mockGetTaskStatusModel,
				CreateTaskUtilHelper: mockCreateTaskUtil,
				UpdateTaskQueue:      mockUpdateTaskStatusModel,
				PublishToMessageBus:  mockPublishToMessageBus,
			},
			args: args{
				req: &taskproto.UpdateTaskRequest{
					TaskID:          "validTaskID",
					TaskState:       "Cancelled",
					TaskStatus:      "Invalid",
					PercentComplete: 30,
					PayLoad:         nil,
					EndTime:         ptypes.TimestampNow(),
				},
				rsp: &taskproto.UpdateTaskResponse{},
			},
			wantErr: fmt.Errorf("error invalid taskStatus provided as input argument"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.ts.UpdateTask(mockContext(), tt.args.req); !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("TasksRPC.UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTasksRPC_CreateTaskUtil(t *testing.T) {
	type args struct {
		userName string
	}
	tests := []struct {
		name    string
		ts      *TasksRPC
		args    args
		want    string
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "Positive case: All is well, valid Username",
			ts: &TasksRPC{
				ValidateTaskUserNameModel: mockValidateTaskUserNameModel,
				PersistTaskModel:          mockPersistTaskModel,
			},
			args: args{
				userName: "validUser",
			},
			wantErr: nil,
		},
		{
			name: "Negative case: empty UserName",
			ts: &TasksRPC{
				ValidateTaskUserNameModel: mockValidateTaskUserNameModel,
				PersistTaskModel:          mockPersistTaskModel,
			},
			args: args{
				userName: "",
			},
			wantErr: fmt.Errorf("error invalid input argument for userName"),
		},
		{
			name: "Negative case: Invalid Username",
			ts: &TasksRPC{
				ValidateTaskUserNameModel: mockValidateTaskUserNameModel,
				PersistTaskModel:          mockPersistTaskModel,
			},
			args: args{
				userName: "InvalidUser",
			},
			wantErr: fmt.Errorf("error invalid user: error while trying to read from DB: %v", errors.PackError(errors.DBKeyNotFound, "no data with the with key ", "InvalidUser", " found").Error()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.ts.CreateTaskUtil(mockContext(), tt.args.userName)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("TasksRPC.CreateTaskUtil() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTasksRPC_taskCancelCallBack(t *testing.T) {
	type args struct {
		taskID string
	}
	tests := []struct {
		name    string
		ts      *TasksRPC
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "Positive case: All is well",
			ts: &TasksRPC{
				GetTaskStatusModel:    mockGetTaskStatusModel,
				UpdateTaskQueue:       mockUpdateTaskStatusModel,
				DeleteTaskFromDBModel: mockDeleteTaskFromDBModel,
			},
			args: args{
				taskID: "validTaskID",
			},
			wantErr: nil,
		},
		{
			name: "Positive case: All is well, But task state is Completed",
			ts: &TasksRPC{
				GetTaskStatusModel:    mockGetTaskStatusModel,
				UpdateTaskQueue:       mockUpdateTaskStatusModel,
				DeleteTaskFromDBModel: mockDeleteTaskFromDBModel,
			},
			args: args{
				taskID: "CompletedTaskID",
			},
			wantErr: nil,
		},
		{
			name: "Positive case: All is well, But task state is Running",
			ts: &TasksRPC{
				GetTaskStatusModel:    mockGetTaskStatusModel,
				UpdateTaskQueue:       mockUpdateTaskStatusModel,
				DeleteTaskFromDBModel: mockDeleteTaskFromDBModel,
			},
			args: args{
				taskID: "RunningTaskID",
			},
			wantErr: nil,
		},
		{
			name: "Negative case: InvalidTaskID",
			ts: &TasksRPC{
				GetTaskStatusModel:    mockGetTaskStatusModel,
				UpdateTaskQueue:       mockUpdateTaskStatusModel,
				DeleteTaskFromDBModel: mockDeleteTaskFromDBModel,
			},
			args: args{
				taskID: "InvalidTaskID",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := mockContext()
			iterCount := new(int)
			ctxt := context.WithValue(ctx, tcommon.IterationCount, iterCount)
			err := tt.ts.taskCancelCallBack(ctxt, tt.args.taskID)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("TasksRPC.taskCancelCallBack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func mockContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.TransactionID, "xyz")
	ctx = context.WithValue(ctx, common.ActionID, "001")
	ctx = context.WithValue(ctx, common.ActionName, "xyz")
	ctx = context.WithValue(ctx, common.ThreadID, "0")
	ctx = context.WithValue(ctx, common.ThreadName, "xyz")
	ctx = context.WithValue(ctx, common.ProcessName, "xyz")
	return ctx
}
