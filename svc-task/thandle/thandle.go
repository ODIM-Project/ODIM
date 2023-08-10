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

// Package thandle ...
package thandle

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-task/tcommon"
	"github.com/ODIM-Project/ODIM/svc-task/tmodel"
	"github.com/ODIM-Project/ODIM/svc-task/tresponse"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
)

const (
	authErrorMessage = "error while trying to authenticate session"
)

var podName = os.Getenv("POD_NAME")

// TasksRPC used to register handler used as rpc call
// AuthenticationRPC is used to authorize user and privileges
// GetTaskStatusModel get task status
type TasksRPC struct {
	AuthenticationRPC                func(ctx context.Context, sessionToken string, privileges []string) (response.RPC, error)
	GetSessionUserNameRPC            func(ctx context.Context, sessionToken string) (string, error)
	GetTaskStatusModel               func(ctx context.Context, taskID string, db common.DbType) (*tmodel.Task, error)
	GetMultipleTaskKeysModel         func(ctx context.Context, taskIDs []interface{}, db common.DbType) (*[]tmodel.Task, error)
	GetAllTaskKeysModel              func(ctx context.Context) ([]string, error)
	TransactionModel                 func(ctx context.Context, key string, cb func(context.Context, string) error) error
	OverWriteCompletedTaskUtilHelper func(ctx context.Context, userName string) error
	CreateTaskUtilHelper             func(ctx context.Context, userName string) (string, error)
	DeleteTaskFromDBModel            func(ctx context.Context, t *tmodel.Task) error
	DeleteMultipleTaskFromDBModel    func(ctx context.Context, t []string) error
	UpdateTaskQueue                  func(t *tmodel.Task)
	PersistTaskModel                 func(ctx context.Context, t *tmodel.Task, db common.DbType) error
	ValidateTaskUserNameModel        func(ctx context.Context, userName string) error
	PublishToMessageBus              func(ctx context.Context, taskURI string, taskEvenMessageID string, eventType string, taskMessage string)
}

// TaskCollectionData ....
type TaskCollectionData struct {
	TaskCollection map[string]int32
	Lock           sync.Mutex
}

func (t *TaskCollectionData) getTaskFromCollectionData(taskID string, percentComplete int) bool {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	if prevComplete, ok := t.TaskCollection[fmt.Sprintf("%s:%v", taskID, percentComplete)]; ok {
		if prevComplete == int32(percentComplete) {
			return true
		} else if percentComplete == 100 {
			delete(t.TaskCollection, taskID)
			return false
		}
	}
	t.TaskCollection[taskID] = int32(percentComplete)

	return false
}

var (
	// TaskCollection ...
	TaskCollection TaskCollectionData
)

// CreateTask is a rpc handler which intern call actual CreateTask to create new task
func (ts *TasksRPC) CreateTask(ctx context.Context, req *taskproto.CreateTaskRequest) (*taskproto.CreateTaskResponse, error) {
	var rsp taskproto.CreateTaskResponse
	// Check for completed task if there are any, get the oldest Completed
	//Task and Delete from the db along with it subtask as well.
	// Search for the Completed tasks
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TaskService, podName)
	l.LogWithFields(ctx).Debugf("Incoming request to create task for user %v", req.UserName)
	taskURI, err := ts.CreateTaskUtilHelper(ctx, req.UserName)
	if err != nil {
		l.LogWithFields(ctx).Error("failed to create task: " + err.Error())
	}
	rsp.TaskURI = taskURI
	l.LogWithFields(ctx).Debugf("Outgoing response for create task request : %v", taskURI)
	return &rsp, err
}

func (ts *TasksRPC) deleteCompletedTask(ctx context.Context, taskID string) error {
	task, err := ts.GetTaskStatusModel(ctx, taskID, common.InMemory)
	if err != nil {
		return fmt.Errorf("error getting taskID - " + taskID + " status : " + err.Error())
	}
	var getChildTask []interface{}
	for _, childkey := range task.ChildTaskIDs {
		getChildTask = append(getChildTask, "task:"+childkey)
	}
	subtasks, err := ts.GetMultipleTaskKeysModel(ctx, getChildTask, common.InMemory)
	if err != nil {
		l.LogWithFields(ctx).Errorf("error getting status of subtask: %s", err.Error())
	}
	var taskStrings []string
	for _, t := range *subtasks {
		taskStrings = append(taskStrings, "task:"+t.ID)
	}
	err = ts.DeleteMultipleTaskFromDBModel(ctx, taskStrings)
	if err != nil {
		l.LogWithFields(ctx).Errorf("error while deleting subtask: %s", err.Error())
	}
	err = ts.DeleteTaskFromDBModel(ctx, task)
	if err != nil {
		l.LogWithFields(ctx).Errorf("error while deleting the main task %s: %s", taskID, err.Error())
		return err
	}
	return nil
}

// CreateChildTask is a rpc handler which intern call actual CreateChildTask to create sub task under parent task.
func (ts *TasksRPC) CreateChildTask(ctx context.Context, req *taskproto.CreateTaskRequest) (*taskproto.CreateTaskResponse, error) {
	var rsp taskproto.CreateTaskResponse
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TaskService, podName)
	l.LogWithFields(ctx).Debugf("Incoming request to create child task for the task %v and user %v", req.ParentTaskID, req.UserName)
	taskURI, err := ts.CreateChildTaskUtil(ctx, req.UserName, req.ParentTaskID)
	if err != nil {
		l.LogWithFields(ctx).Errorf("failed to create child task for the task ID %v : %v", req.ParentTaskID, err.Error())
	}
	rsp.TaskURI = taskURI
	l.LogWithFields(ctx).Debugf("Outgoing response for create child task request : %v", taskURI)
	return &rsp, err
}

// UpdateTask is a rpc handler which inter call actual CreateTask to create new task
func (ts *TasksRPC) UpdateTask(ctx context.Context, req *taskproto.UpdateTaskRequest) (*taskproto.UpdateTaskResponse, error) {
	var rsp taskproto.UpdateTaskResponse
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TaskService, podName)
	l.LogWithFields(ctx).Debugf("Incoming request to update task %v", req.TaskID)
	endTime, err := ptypes.Timestamp(req.EndTime)
	if err != nil {
		l.LogWithFields(ctx).Error("failed to update task: error while trying to convert proto-buff timestamp to time.Time: " + err.Error())
		return &rsp, err
	}
	err = ts.updateTaskUtil(ctx, req.TaskID, req.TaskState, req.TaskStatus, req.PercentComplete, req.PayLoad, endTime)
	if err != nil {
		l.LogWithFields(ctx).Error("failed to update task: error while updating task: " + err.Error())
	}
	return &rsp, err
}

// DeleteTask is an API end point to delete the given task.
func (ts *TasksRPC) DeleteTask(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	var rsp taskproto.TaskResponse
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TaskService, podName)
	logPrefix := fmt.Sprintf("failed to delete task : %v", req.TaskID)
	l.LogWithFields(ctx).Debugf("Incoming request to delete task %v", req.TaskID)
	constructCommonResponseHeader(&rsp)
	task, err := ts.validateAndAuthorize(ctx, req, &rsp)
	if err != nil {
		return &rsp, nil
	}
	privileges := []string{common.PrivilegeConfigureManager}
	authResp, err := ts.AuthenticationRPC(ctx, req.SessionToken, privileges)
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf(logPrefix+"Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, &rsp, authResp)
		return &rsp, nil

	}
	rsp.Header["Allow"] = "DELETE"
	if task.PercentComplete == 100 {
		delErr := ts.deleteCompletedTask(ctx, req.TaskID)
		if delErr != nil {
			errorMessage := "Error while deleting the completed task: " + delErr.Error()
			l.LogWithFields(ctx).Error(logPrefix + errorMessage)
			fillProtoResponse(ctx, &rsp, common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil))
			return &rsp, nil
		}
		rsp.StatusCode = http.StatusNoContent
		rsp.Body = nil
		return &rsp, nil
	}
	// Critical Logic follows

	// Cancel the task using Transaction
	iterCount := new(int)
	ctxt := context.WithValue(ctx, tcommon.IterationCount, iterCount)
	for iter := 0; iter < 5; iter++ {
		err = ts.TransactionModel(ctxt, req.TaskID, ts.taskCancelCallBack)
		if err != nil {
			l.LogWithFields(ctx).Error(logPrefix + "error while requesting for task cancellation retrying: " + err.Error())
			*iterCount++
			continue
		}
		break
	}
	if err != nil {
		errorMessage := "error max retries exceeded for TaskCancel Transaction: " + err.Error()
		fillProtoResponse(ctx, &rsp, common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil))
		l.LogWithFields(ctx).Error(logPrefix + errorMessage)
		return &rsp, nil
	}

	// Critical Logic Ends

	// build the response
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
		OdataType:    common.TaskType,
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
			HTTPHeaders:   httpHeaders,
			HTTPOperation: task.Payload.HTTPOperation,
			JSONBody:      string(task.Payload.JSONBody),
			TargetURI:     task.Payload.TargetURI,
		},
		PercentComplete: task.PercentComplete,
	}
	if task.ParentID == "" {
		var subTask = tresponse.ListMember{
			OdataID: "/redfish/v1/TaskService/Tasks/" + task.ID + "/SubTasks",
		}
		taskResponse.SubTasks = &subTask
	}
	//  return tasks in case of Success
	//Frame the response body below to send back to the user
	rsp.Body = generateResponse(ctx, taskResponse) // cannot convert task response directly to []byte that's why it needs to be marshalled and send as response in byte format
	return &rsp, nil
}
func constructCommonResponseHeader(rsp *taskproto.TaskResponse) {
	rsp.Header = map[string]string{
		"Date": time.Now().Format(http.TimeFormat),
	}

}

func (ts *TasksRPC) validateAndAuthorize(ctx context.Context, req *taskproto.GetTaskRequest, rsp *taskproto.TaskResponse) (*tmodel.Task, error) {
	privileges := []string{common.PrivilegeLogin}
	authResp, err := ts.AuthenticationRPC(ctx, req.SessionToken, privileges)
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, rsp, authResp)
		l.LogWithFields(ctx).Error(authErrorMessage)
		return nil, fmt.Errorf(authErrorMessage)
	}
	sessionUserName, err := ts.GetSessionUserNameRPC(ctx, req.SessionToken)
	if err != nil {
		// handle the error case with appropriate response body
		fillProtoResponse(ctx, rsp, common.GeneralError(http.StatusUnauthorized, response.NoValidSession, authErrorMessage, nil, nil))
		l.LogWithFields(ctx).Error(authErrorMessage)
		return nil, fmt.Errorf(authErrorMessage)

	}
	// get task status from database using task id
	task, err := ts.GetTaskStatusModel(ctx, req.TaskID, common.InMemory)
	if err != nil {
		l.LogWithFields(ctx).Error("error getting task status : " + err.Error())
		fillProtoResponse(ctx, rsp, common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"Task", req.TaskID}, nil))
		return nil, err
	}
	//Compare the task username with requesting session user name.
	//If username doesn't match with task username, then check if the user
	//is an Admin(PrivilegeConfigureUsers). If he is admin then proceed.
	if sessionUserName != task.UserName {
		privileges := []string{common.PrivilegeConfigureUsers}
		authResp, err := ts.AuthenticationRPC(ctx, req.SessionToken, privileges)
		if authResp.StatusCode != http.StatusOK {
			if err != nil {
				l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
			}
			fillProtoResponse(ctx, rsp, authResp)
			return nil, fmt.Errorf(authErrorMessage)
		}
	}
	return task, nil
}

func (ts *TasksRPC) taskCancelCallBack(ctx context.Context, taskID string) error {
	task, err := ts.GetTaskStatusModel(ctx, taskID, common.InMemory)
	if err != nil {
		l.LogWithFields(ctx).Error("error getting task status : " + err.Error())
		return nil
	}
	if task.TaskState == common.Completed || task.TaskState == common.Exception || task.TaskState == common.Pending {
		var getChildTask []interface{}
		// check if this task has any child tasks, if so delete them.
		for _, childkey := range task.ChildTaskIDs {
			getChildTask = append(getChildTask, "task:"+childkey)
		}
		subtasks, err := ts.GetMultipleTaskKeysModel(ctx, getChildTask, common.InMemory)
		if err != nil {
			l.LogWithFields(ctx).Errorf("error getting status of subtask: %s", err.Error())
		}
		var taskStrings []string
		for _, t := range *subtasks {
			taskStrings = append(taskStrings, "task:"+t.ID)
		}
		err = ts.DeleteMultipleTaskFromDBModel(ctx, taskStrings)
		if err != nil {
			l.LogWithFields(ctx).Errorf("error while deleting subtask: %s", err.Error())
		}
		ts.DeleteTaskFromDBModel(ctx, task)
		return nil
	}
	threadID := ctx.Value(tcommon.IterationCount).(*int)
	newCtx := context.WithValue(ctx, common.ThreadName, common.AsyncTaskDelete)
	for _, subTaskID := range task.ChildTaskIDs {
		subTask, err := ts.GetTaskStatusModel(ctx, subTaskID, common.InMemory)
		if err != nil {
			l.LogWithFields(ctx).Error("error getting task status : " + err.Error())
			continue
		}
		// Just changing the TaskState to Cancelling state,
		// After this the thread associated with this task, it can be in any service can see this change and
		// mark the taskstate to Cancelled exits.
		if subTask.TaskState == common.Completed || subTask.TaskState == common.Exception || subTask.TaskState == common.Pending {
			ts.DeleteTaskFromDBModel(ctx, subTask)
		} else if subTask.TaskState != common.Cancelling {
			subTask.TaskState = common.Cancelling
			ts.UpdateTaskQueue(subTask)
			newCtx = context.WithValue(newCtx, common.ThreadID, strconv.Itoa(*threadID))
			go ts.asyncTaskDelete(newCtx, subTaskID)
			*threadID++
		}
	}
	// Delete the parent task
	if task.TaskState != common.Cancelling {
		task.TaskState = common.Cancelling
		ts.UpdateTaskQueue(task)
		newCtx = context.WithValue(newCtx, common.ThreadID, strconv.Itoa(*threadID))
		go ts.asyncTaskDelete(newCtx, taskID)
		*threadID++
	}

	return nil
}

func (ts *TasksRPC) asyncTaskDelete(ctx context.Context, taskID string) {
	//Polling for the taskstate.
	//If the taskstate becomes Cancelled, then this means the thread associated with this task exited successfully,
	//so go ahead delete the task from the db

	// Get the task
	for {
		task, err := ts.GetTaskStatusModel(ctx, taskID, common.InMemory)
		if err != nil {
			l.LogWithFields(ctx).Error("error getting task status : " + err.Error())
			return
		}
		if task.TaskState == common.Cancelled {
			err = ts.DeleteTaskFromDBModel(ctx, task)
			if err != nil {
				l.LogWithFields(ctx).Error("error unable to delete the task from db: " + err.Error())
				return
			}
			break
		}
		time.Sleep(5000 * time.Millisecond)
	}
	return
}

// GetSubTasks is an API end point to get all available tasks
func (ts *TasksRPC) GetSubTasks(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	var rsp taskproto.TaskResponse
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TaskService, podName)
	constructCommonResponseHeader(&rsp)

	l.LogWithFields(ctx).Debugf("Incoming request to get all available subtasks of task %v", req.TaskID)
	task, err := ts.validateAndAuthorize(ctx, req, &rsp)
	if err != nil {
		return &rsp, nil
	}
	var listMembers []tresponse.ListMember
	for _, subTaskID := range task.ChildTaskIDs {
		// Build the subtask list with all child tasks
		member := tresponse.ListMember{
			OdataID: "/redfish/v1/TaskService/Tasks/" + task.ID + "/SubTasks/" + subTaskID,
		}
		listMembers = append(listMembers, member)
	}

	commonResponse := response.Response{
		OdataContext: "/redfish/v1/$metadata#SubTasks.SubTasks",
		OdataID:      "/redfish/v1/TaskService/Tasks/" + task.ID + "/SubTasks/",
		OdataType:    "#TaskCollection.TaskCollection",
		Name:         "SubTasks",
		Description:  "SubTasks",
	}

	rsp.StatusCode = http.StatusOK
	rsp.StatusMessage = response.Success

	//Frame the Response to send it back as response body
	taskResp := tresponse.TaskCollectionResponse{
		Response:     commonResponse,
		MembersCount: len(listMembers),
		Members:      listMembers,
	}

	respBody := generateResponse(ctx, taskResp)
	rsp.Body = respBody
	l.LogWithFields(ctx).Debugf("Outgoing response for getting subtasks: %v", string(respBody))
	return &rsp, nil
}

// GetSubTask is an API end point to get the subtask details
func (ts *TasksRPC) GetSubTask(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	var rsp taskproto.TaskResponse
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TaskService, podName)
	constructCommonResponseHeader(&rsp)
	l.LogWithFields(ctx).Debugf("Incoming request to get subtask %v", req.SubTaskID)
	privileges := []string{common.PrivilegeLogin}
	authResp, err := ts.AuthenticationRPC(ctx, req.SessionToken, privileges)
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, &rsp, authResp)
		return &rsp, nil
	}
	sessionUserName, err := ts.GetSessionUserNameRPC(ctx, req.SessionToken)
	if err != nil {
		fillProtoResponse(ctx, &rsp, common.GeneralError(http.StatusUnauthorized,
			response.NoValidSession, authErrorMessage, nil, nil))
		l.LogWithFields(ctx).Error(authErrorMessage)
		return &rsp, nil
	}
	// get task status from database using task id
	task, err := ts.GetTaskStatusModel(ctx, req.SubTaskID, common.InMemory)
	if err != nil {
		l.LogWithFields(ctx).Error("error getting sub task status : " + err.Error())
		fillProtoResponse(ctx, &rsp, common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(),
			[]interface{}{"Task", req.SubTaskID}, nil))
		return &rsp, nil
	}
	//Compare the task username with requesting session user name
	if sessionUserName != task.UserName {
		privileges := []string{common.PrivilegeConfigureUsers}
		authResp, err := ts.AuthenticationRPC(ctx, req.SessionToken, privileges)
		if authResp.StatusCode != http.StatusOK {
			if err != nil {
				l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
			}
			fillProtoResponse(ctx, &rsp, authResp)
			return &rsp, nil
		}
	}

	//Build the respose Body
	var messageList []tresponse.Messages
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

	var listMembers []tresponse.ListMember
	for _, subTaskID := range task.ChildTaskIDs {
		// Build the subtask list with all child tasks
		member := tresponse.ListMember{
			OdataID: "/redfish/v1/TaskService/Tasks/" + task.ID + "/SubTasks/" + subTaskID,
		}
		listMembers = append(listMembers, member)
	}
	rsp.StatusMessage = response.Success
	commonResponse := response.Response{
		ID:           task.ID,
		Name:         task.Name,
		OdataContext: "/redfish/v1/$metadata#SubTask.SubTask",
		OdataID:      "/redfish/v1/TaskService/Tasks/" + req.TaskID + "/SubTasks/" + req.SubTaskID,
		OdataType:    "#SubTask.v1_4_2.SubTask",
	}
	commonResponse.CreateGenericResponse(rsp.StatusMessage)
	httpHeaders := []string{}
	for key, value := range task.Payload.HTTPHeaders {
		httpHeaders = append(httpHeaders, fmt.Sprintf("%v: %v", key, value))
	}

	taskResponse := tresponse.SubTask{
		Response:     commonResponse,
		MembersCount: len(listMembers),
		Members:      listMembers,
		TaskState:    task.TaskState,
		StartTime:    task.StartTime.UTC(),
		EndTime:      task.EndTime.UTC(),
		TaskStatus:   task.TaskStatus,
		Messages:     messageList,
		TaskMonitor:  task.TaskMonitor,
		Payload: tresponse.Payload{
			HTTPHeaders:   httpHeaders,
			HTTPOperation: task.Payload.HTTPOperation,
			JSONBody:      string(task.Payload.JSONBody),
			TargetURI:     task.Payload.TargetURI,
		},
		PercentComplete: task.PercentComplete,
	}

	// Check the state of the task
	if task.TaskState == "Completed" || task.TaskState == "Cancelled" || task.TaskState == "Killed" ||
		task.TaskState == "Exception" {
		// return with the 200 OK, along with response header and response body
		rsp.StatusCode = http.StatusOK
	} else {
		// return 202
		// build response header
		// return with empty response body
		rsp.Header["location"] = task.TaskMonitor
		rsp.StatusCode = http.StatusAccepted
	}
	// cannot convert task response directly to []byte that's why it needs to be marshalled and send as response in byte format
	respBody := generateResponse(ctx, taskResponse)
	rsp.Body = respBody
	l.LogWithFields(ctx).Debugf("Outgoing response for getting subtask: %v", string(respBody))

	return &rsp, nil
}

// TaskCollection is an API end point to get all available tasks
func (ts *TasksRPC) TaskCollection(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	var rsp taskproto.TaskResponse
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TaskService, podName)
	l.LogWithFields(ctx).Debugf("Incoming request to get task collection")
	commonResponse := response.Response{
		Name:         "Task Collection",
		OdataContext: "/redfish/v1/$metadata#TaskCollection.TaskCollection",
		OdataID:      "/redfish/v1/TaskService/Tasks",
		OdataType:    "#TaskCollection.TaskCollection",
	}
	constructCommonResponseHeader(&rsp)
	privileges := []string{common.PrivilegeLogin}
	authResp, err := ts.AuthenticationRPC(ctx, req.SessionToken, privileges)
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, &rsp, authResp)
		return &rsp, nil
	}
	// Get all task in in-memory db
	tasks, err := ts.GetAllTaskKeysModel(ctx)
	if err != nil {
		errorMessage := "error: while trying to get all task keys from db: " + err.Error()
		fillProtoResponse(ctx, &rsp, common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil))
		l.LogWithFields(ctx).Error(errorMessage)
		return &rsp, nil
	}
	statusConfigureUsers, err := ts.AuthenticationRPC(ctx, req.SessionToken, []string{common.PrivilegeConfigureUsers})
	if err != nil {
		l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
	}
	sessionUserName, err := ts.GetSessionUserNameRPC(ctx, req.SessionToken)
	if err != nil {
		fillProtoResponse(ctx, &rsp, common.GeneralError(http.StatusUnauthorized, response.NoValidSession, authErrorMessage, nil, nil))
		l.LogWithFields(ctx).Error(authErrorMessage)
		return &rsp, nil

	}
	var listMembers = []tresponse.ListMember{}
	for _, taskID := range tasks {
		// Check who owns the task before returning, if this can only be done by admin,
		//then its appropriate to give back all the tasks available in the DB
		//If user has just login privilege then return his own task
		if authResp.StatusCode == http.StatusOK && statusConfigureUsers.StatusCode != http.StatusOK {
			task, err := ts.GetTaskStatusModel(ctx, taskID, common.InMemory)
			if err != nil {
				l.LogWithFields(ctx).Error("error getting task status : " + err.Error())
				fillProtoResponse(ctx, &rsp, common.GeneralError(http.StatusNotFound,
					response.ResourceNotFound, authErrorMessage, nil, nil))
				return &rsp, nil
			}
			//Check if the task belongs to user
			if task.UserName == sessionUserName {
				member := tresponse.ListMember{OdataID: "/redfish/v1/TaskService/Tasks/" + taskID}
				listMembers = append(listMembers, member)
			}
		}
		//if user has configureusers privelege then return all tasks
		if statusConfigureUsers.StatusCode == http.StatusOK {
			member := tresponse.ListMember{OdataID: "/redfish/v1/TaskService/Tasks/" + taskID}
			listMembers = append(listMembers, member)
		}
	}

	// return response with status OK
	rsp.StatusCode = http.StatusOK
	rsp.StatusMessage = response.Success

	//Frame the Response to send it back as response body
	taskResp := tresponse.TaskCollectionResponse{
		Response:     commonResponse,
		MembersCount: len(listMembers),
		Members:      listMembers,
	}
	respBody := generateResponse(ctx, taskResp)
	rsp.Body = respBody
	l.LogWithFields(ctx).Debugf("Outgoing response for getting task collection: %v", string(respBody))
	return &rsp, nil
}

// GetTasks is an API end point to get the task status and response body.
// Takes X-Auth-Token and authorize the request.
// If X-Auth-Token is empty or invalid then it returns "StatusUnauthorized".
// If the TaskID is not found then it return "StatusNotFound".
// If the task is still not completed or cancelled or killed then it return with 202
// with empty response body, else it return with "200 OK" with full task info in the
// response body.
// If the Username doesn't match with the task username then it returns with
// StatusForbidden.
func (ts *TasksRPC) GetTasks(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	var rsp taskproto.TaskResponse
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TaskService, podName)

	l.LogWithFields(ctx).Debugf("Incoming request to get task status of task %v", req.TaskID)
	constructCommonResponseHeader(&rsp)
	task, err := ts.validateAndAuthorize(ctx, req, &rsp)
	if err != nil {
		return &rsp, nil
	}
	rsp.Header["Link"] = "</redfish/v1/SchemaStore/en/TaskCollection.json/>; rel=describedby"
	//Build the respose Body
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
		OdataType:    common.TaskType,
		ID:           task.ID,
		Name:         task.Name,
		OdataContext: "/redfish/v1/$metadata#Task.Task",
		OdataID:      "/redfish/v1/TaskService/Tasks/" + task.ID,
	}
	rsp.StatusCode = http.StatusAccepted
	rsp.StatusMessage = response.Success
	commonResponse.CreateGenericResponse(rsp.StatusMessage)
	commonResponse.Message = ""
	commonResponse.MessageID = ""
	commonResponse.Severity = ""

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
			HTTPHeaders:   httpHeaders,
			HTTPOperation: task.Payload.HTTPOperation,
			JSONBody:      string(task.Payload.JSONBody),
			TargetURI:     task.Payload.TargetURI,
		},
		PercentComplete: task.PercentComplete,
	}
	if task.ParentID == "" && len(task.ChildTaskIDs) != 0 {
		var subTask = tresponse.ListMember{
			OdataID: "/redfish/v1/TaskService/Tasks/" + task.ID + "/SubTasks",
		}
		taskResponse.SubTasks = &subTask
	}
	// Check the state of the task
	if task.TaskState == "Completed" || task.TaskState == "Cancelled" || task.TaskState == "Killed" || task.TaskState == "Exception" {
		// return with the 200 OK, along with response header and response body
		rsp.StatusCode = http.StatusOK
	} else {
		// return 202
		// build response header
		// return with empty response body
		rsp.Header["location"] = task.TaskMonitor
		rsp.StatusCode = http.StatusAccepted
	}
	rsp.StatusMessage = "Success"
	// cannot convert task response directly to []byte that's why it needs to be marshalled and send as response in byte format
	respBody := generateResponse(ctx, taskResponse)
	rsp.Body = respBody
	l.LogWithFields(ctx).Debugf("Outgoing response for getting task status: %v", string(respBody))
	return &rsp, nil
}

// GetTaskService is an API handler to get Task service details
// Takes:
//
//	taskproto.GetTaskRequest(extracts SessionToken from it)
//
// Returns:
//
//	401 Unauthorized or 200 OK with respective response body and response header.
func (ts *TasksRPC) GetTaskService(ctx context.Context, req *taskproto.GetTaskRequest) (*taskproto.TaskResponse, error) {
	var rsp taskproto.TaskResponse
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TaskService, podName)

	l.LogWithFields(ctx).Debugf("Incoming request to get task service request")
	// Fill the response header first
	constructCommonResponseHeader(&rsp)
	rsp.Header["Link"] = "</redfish/v1/SchemaStore/en/TaskService.json>; rel=describedby"
	// Validate the token, if user has ConfigureUsers privilege then proceed.
	//Else send 401 Unauthorized
	privileges := []string{common.PrivilegeConfigureUsers}
	authResp, err := ts.AuthenticationRPC(ctx, req.SessionToken, privileges)
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, &rsp, authResp)
		return &rsp, nil
	}

	// Check whether the Task Service is enabled in configuration file.
	//If so set ServiceEnabled to true.
	isServiceEnabled := false
	serviceState := "Disabled"
	for _, service := range config.Data.EnabledServices {
		if service == "TaskService" {
			isServiceEnabled = true
			serviceState = "Enabled"
			break
		}
	}

	rsp.StatusCode = http.StatusOK
	rsp.StatusMessage = response.Success
	commonResponse := response.Response{
		OdataType:    "#TaskService.v1_2_0.TaskService",
		ID:           "TaskService",
		Name:         "TaskService",
		Description:  "TaskService",
		OdataContext: "/redfish/v1/$metadata#TaskService.TaskService",
		OdataID:      "/redfish/v1/TaskService",
	}

	// Construct the response body hear as below
	taskServiceResponse := tresponse.TaskServiceResponse{
		Response:                        commonResponse,
		CompletedTaskOverWritePolicy:    "Oldest",
		DateTime:                        time.Now().UTC(),
		LifeCycleEventOnTaskStateChange: true,
		ServiceEnabled:                  isServiceEnabled,
		Status: tresponse.Status{
			State:        serviceState,
			Health:       "OK",
			HealthRollup: "OK",
		},
		Tasks: tresponse.Tasks{
			OdataID: "/redfish/v1/TaskService/Tasks",
		},
	}
	respBody := generateResponse(ctx, taskServiceResponse)
	rsp.Body = respBody
	l.LogWithFields(ctx).Debugf("Outgoing response for getting task service details: %v", string(respBody))
	return &rsp, nil
}

func generateResponse(ctx context.Context, input interface{}) []byte {
	bytes, err := json.Marshal(input)
	if err != nil {
		l.LogWithFields(ctx).Error("error in unmarshal response object from util-libs" + err.Error())
	}
	return bytes
}

func fillProtoResponse(ctx context.Context, resp *taskproto.TaskResponse, data response.RPC) {
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
	resp.Header = data.Header

}

// CreateTaskUtil Create the New Task and persist in in-memory DB and return task ID and error
// Takes :
//
//	username : Is a Username of type string
//
// Returns:
//
//	New Task URI of Type string
//	err of type error
func (ts *TasksRPC) CreateTaskUtil(ctx context.Context, userName string) (string, error) {
	if userName == "" {
		return "", fmt.Errorf("error invalid input argument for userName")
	}
	// Validate given username exist in the db
	err := ts.ValidateTaskUserNameModel(ctx, userName)
	if err != nil {
		return "", fmt.Errorf("error invalid user: " + err.Error())
	}
	// Frame the model
	currentTime := time.Now()
	task := tmodel.Task{
		UserName:        userName,
		ParentID:        "",
		ChildTaskIDs:    nil,
		ID:              "task" + uuid.New().String(),
		TaskState:       "New",
		TaskStatus:      "OK",
		PercentComplete: 0,
		StartTime:       currentTime,
		EndTime:         time.Time{},
	}
	task.Name = "Task " + task.ID
	task.TaskMonitor = "/taskmon/" + task.ID
	task.URI = "/redfish/v1/TaskService/Tasks/" + task.ID

	// Persist in the in-memory DB
	err = ts.PersistTaskModel(ctx, &task, common.InMemory)
	if err != nil {
		return "", fmt.Errorf("error while trying to insert the task details: " + err.Error())
	}
	// return the Task URI
	return "/redfish/v1/TaskService/Tasks/" + task.ID, err
}

// CreateChildTaskUtil Creates the child task and attaches to the parent task provided.
// Tasks:
//
//	parentTaskID of type string - Contains Parent task ID for Child task yet to be created
//
// Returns:
//
//	err of type error
//	nil - On Success
//	Non nil - On Failure
func (ts *TasksRPC) CreateChildTaskUtil(ctx context.Context, userName string, parentTaskID string) (string, error) {

	var parentTask *tmodel.Task
	var childTask *tmodel.Task
	var taskURI string
	if parentTaskID == "" {
		return "", fmt.Errorf("error empty/invalid input Parent Task ID")
	}
	// Retrieve the task details from db
	parentTask, err := ts.GetTaskStatusModel(ctx, parentTaskID, common.InMemory)
	if err != nil {
		return "", fmt.Errorf("error while retrieving the task details from DB: " + err.Error())
	}
	// Create the child/sub task with parent task's UserName
	taskURI, err = ts.CreateTaskUtilHelper(ctx, parentTask.UserName)
	if err != nil {
		return "", fmt.Errorf("error while creating child/sub task: " + err.Error())
	}
	var childTaskID string
	strArray := strings.Split(taskURI, "/")
	if strings.HasSuffix(taskURI, "/") {
		childTaskID = strArray[len(strArray)-2]
	} else {
		childTaskID = strArray[len(strArray)-1]
	}
	// Get the Child task to update with Parent task ID
	childTask, err = ts.GetTaskStatusModel(ctx, childTaskID, common.InMemory)
	if err != nil {
		return "", fmt.Errorf("error while retrieving the child/sub task from DB: " + err.Error())
	}
	childTask.ParentID = parentTaskID
	childTask.URI = "/redfish/v1/TaskService/Tasks/" + parentTaskID + "/" + childTaskID
	// Store the updated task in to In Memory DB
	ts.UpdateTaskQueue(childTask)
	// Add the child/sub task id in to ChildTaskIDs(array) of the parent task
	parentTask.ChildTaskIDs = append(parentTask.ChildTaskIDs, childTaskID)
	// Update the parent task in to In Memory DB
	ts.UpdateTaskQueue(parentTask)
	return "/redfish/v1/TaskService/Tasks/" + childTaskID, err
}

// getAllChildTasks is used to get All child task ID's associated with parent task
func (ts *TasksRPC) getAllChildTasks(ctx context.Context, parentID string) ([]string, error) {
	task, err := ts.GetTaskStatusModel(ctx, parentID, common.InMemory)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving the task details from db: " + err.Error())
	}
	return task.ChildTaskIDs, nil
}

// updateTaskUtil is a function to update the existing task and/or to create sub-task under a parent task.
// This function is to set task status, task end time along with task state based on the task state.
// Takes:
//
//	taskID - Is of type string, contains task ID of the task to updated
//	taskState - Is of type string, contains new sate of the task
//	taskStatus - Is of type string, contains new status of the task
//	endTime    - Is of type time.Time, containers the endtime of the task
//
// Returns:
//
//	err of type error
//	nil - On Success
//	Non nil - On Failure
func (ts *TasksRPC) updateTaskUtil(ctx context.Context, taskID string, taskState string,
	taskStatus string, percentComplete int32, payLoad *taskproto.Payload, endTime time.Time) error {

	var task *tmodel.Task
	var taskEvenMessageID, taskMessage string
	// Retrieve the task details using taskID
	task, err := ts.GetTaskStatusModel(ctx, taskID, common.InMemory)
	if err != nil {
		return fmt.Errorf("error while retrieving the task details from db: " + err.Error())
	}

	if task.PercentComplete > percentComplete {
		return fmt.Errorf("the task with id %s is already updated with %d percent complete."+
			"skipping the update request with the percent complete %d", taskID,
			task.PercentComplete, percentComplete)
	}

	//If the task is already in cancelled state, then updates are not allowed to it.
	if task.TaskState == common.Cancelled {
		return fmt.Errorf(common.Cancelled)
	}
	if task.TaskState == common.Cancelling && taskState != common.Cancelled {
		return fmt.Errorf(common.Cancelling)
	}
	// Set the task state
	switch taskState {

	case "Completed":
		/* This State shall represent that the operation is complete and completed
		successfully or with warnings.
		*/
		task.TaskState = taskState
		if taskStatus == "Critical" || taskStatus == "Warning" || taskStatus == "OK" {
			task.TaskStatus = taskStatus
		} else {
			return fmt.Errorf("error invalid taskStatus provided as input argument")
		}
		if endTime == (time.Time{}) {
			return fmt.Errorf("error invalid end time for the task")
		}
		task.EndTime = endTime
		if payLoad != nil {
			task.Payload.HTTPOperation = payLoad.HTTPOperation
			task.Payload.HTTPHeaders = payLoad.HTTPHeaders
			task.Payload.JSONBody = payLoad.JSONBody
			task.Payload.TargetURI = payLoad.TargetURI
			task.StatusCode = payLoad.StatusCode
			task.TaskResponse = payLoad.ResponseBody
		}
		task.PercentComplete = percentComplete
		// Construct the appropriate messageID for task status change notification
		taskEvenMessageID = common.TaskEventType + ".Task" + taskState + taskStatus
		taskMessage = fmt.Sprintf("The task with Id %v has completed.", taskID)
	case "Killed":
		/*This state shall represent that the operation is complete because the task
		was killed by an operator. Deprecated v1.2+. This value has been deprecated
		and is being replaced by the value Cancelled which has	more determinate
		semantics.
		*/
		task.TaskState = taskState
		if taskStatus == "Critical" || taskStatus == "Warning" {
			task.TaskStatus = taskStatus
		} else {
			return fmt.Errorf("error invalid taskStatus provided as input argument")
		}
		if endTime == (time.Time{}) {
			return fmt.Errorf("error invalid end time for the task")
		}
		task.PercentComplete = percentComplete
		if payLoad != nil {
			task.StatusCode = payLoad.StatusCode
		}
		task.EndTime = endTime
		// Construct the appropriate messageID for task status change notification
		taskEvenMessageID = common.TaskEventType + ".TaskAborted"
		taskMessage = fmt.Sprintf("The task with Id %v has completed with errors.", taskID)
	case "Cancelled":
		/* This state shall represent that the operation was cancelled either
		through a Delete on a Task Monitor or Task Resource or by an internal
		process.
		*/
		task.TaskState = taskState
		if taskStatus == "Critical" || taskStatus == "Warning" {
			task.TaskStatus = taskStatus
		} else {
			return fmt.Errorf("error invalid taskStatus provided as input argument")
		}
		if endTime == (time.Time{}) {
			return fmt.Errorf("error invalid end time for the task")
		}
		task.PercentComplete = percentComplete
		if payLoad != nil {
			task.StatusCode = payLoad.StatusCode
			task.TaskResponse = payLoad.ResponseBody
		}
		task.EndTime = endTime
		// Construct the appropriate messageID for task status change notification
		taskEvenMessageID = common.TaskEventType + ".Task" + taskState
		taskMessage = fmt.Sprintf("Work on the task with Id %v has been halted prior to completion due to an explicit request.", taskID)
	case "Exception":
		/* This state shall represent that the operation is complete and
		completed with errors.
		*/
		task.TaskState = taskState
		if taskStatus == "Critical" || taskStatus == "Warning" {
			task.TaskStatus = taskStatus
		} else {
			return fmt.Errorf("error invalid taskStatus provided as input argument")
		}
		if endTime == (time.Time{}) {
			return fmt.Errorf("error invalid end time for the task")
		}
		task.EndTime = endTime
		if payLoad != nil {
			task.Payload.HTTPOperation = payLoad.HTTPOperation
			task.Payload.HTTPHeaders = payLoad.HTTPHeaders
			task.Payload.JSONBody = payLoad.JSONBody
			task.Payload.TargetURI = payLoad.TargetURI
			task.StatusCode = payLoad.StatusCode
			task.TaskResponse = payLoad.ResponseBody
		}
		task.PercentComplete = percentComplete
		// Construct the appropriate messageID for task status change notification
		taskEvenMessageID = common.TaskEventType + ".Task" + taskState + taskStatus
		taskMessage = fmt.Sprintf("The task with Id %v has completed with errors.", taskID)
	case "Cancelling":
		/*This state shall represent that the operation is in the process of being
		cancelled.
		*/
		task.TaskState = taskState
		// Construct the appropriate messageID for task status change notification
		taskEvenMessageID = common.TaskEventType + ".Task" + taskState
		taskMessage = fmt.Sprintf("Work on the task with Id %v has been halted prior to completion due to an explicit request.", taskID)
		// TODO
	case "Interrupted":
		/* This state shall represent that the operation has been interrupted but is
		expected to restart and is therefore not complete.
		*/
		task.TaskState = taskState
		if payLoad != nil {
			task.StatusCode = payLoad.StatusCode
		}
		task.PercentComplete = percentComplete
		// Construct the appropriate messageID for task status change nitification
		taskEvenMessageID = common.TaskEventType + ".Task" + taskState
		taskMessage = fmt.Sprintf("The task with Id %v has completed with errors..", taskID)
		// TODO
	case "New":
		/* This state shall represent that this task is newly created but the
		operation has not yet started.
		*/
		task.TaskState = taskState
		task.PercentComplete = percentComplete
		// Construct the appropriate messageID for task status change nitification
		taskEvenMessageID = common.TaskEventType + ".Task" + taskState
		taskMessage = fmt.Sprintf("The task with Id %v has started.", taskID)
		// TODO
	case "Pending":
		/*This state shall represent that the operation is pending some condition and
		has not yet begun to execute.
		*/
		task.TaskState = taskState
		// Construct the appropriate messageID for task status change nitification
		taskEvenMessageID = common.TaskEventType + ".Task" + taskState
		taskMessage = fmt.Sprintf("The task with Id %v has completed with errors.", taskID)
		// TODO
	case "Running":
		// This state shall represent that the operation is executing.
		if payLoad.FinalResponseBody != nil {
			task.TaskFinalResponse = payLoad.FinalResponseBody
		}
		task.TaskState = taskState
		task.PercentComplete = percentComplete
		// Construct the appropriate messageID for task status change notification
		taskEvenMessageID = common.TaskEventType + ".TaskProgressChanged"
		taskMessage = fmt.Sprintf("The task with Id %v has changed to progress %v percent complete.", taskID, percentComplete)
		// TODO
	case "Service":
		/* This state shall represent that the operation is now running as a service
		and expected to continue operation until stopped or killed.
		*/
		task.TaskState = taskState
		// Construct the appropriate messageID for task status change notification
		taskEvenMessageID = common.TaskEventType + ".Task" + taskState
		taskMessage = fmt.Sprintf("The task with Id %v has started.", taskID)
		// TODO
	case "Starting":
		// This state shall represent that the operation is starting.
		task.TaskState = taskState
		// Construct the appropriate messageID for task status change notification
		taskEvenMessageID = common.TaskEventType + ".Task" + taskState
		taskMessage = fmt.Sprintf("The task with Id %v has started.", taskID)
		// TODO
	case "Stopping":
		/* This state shall represent that the operation is stopping but is not yet
		complete.
		*/
		task.TaskState = taskState
		// Construct the appropriate messageID for task status change notification
		taskEvenMessageID = common.TaskEventType + ".Task" + taskState
		taskMessage = fmt.Sprintf("The task with Id %v has been paused.", taskID)
		// TODO
	case "Suspended":
		/*This state shall represent that the operation has been suspended but is
		expected to restart and is therefore not complete.
		*/
		task.TaskState = taskState
		task.PercentComplete = percentComplete
		// Construct the appropriate messageID for task status change notification
		taskEvenMessageID = common.TaskEventType + ".Task" + taskState
		taskMessage = fmt.Sprintf("The task with Id %v has completed with errors.", taskID)
		// TODO
	default:
		return fmt.Errorf("error invalid input argument for taskState")
	}
	// Update the task data in the InMemory DB
	ts.UpdateTaskQueue(task)
	l.LogWithFields(ctx).Debugf("update task request for task id %s is pushed to to queue", taskID)
	// Notify the user about task state change by sending status change event
	//	notifyTaskStateChange(task.URI, taskEvenMessageID)
	eventType := "StatusChange"

	if !TaskCollection.getTaskFromCollectionData(taskID, int(percentComplete)) {
		ts.PublishToMessageBus(ctx, task.URI, taskEvenMessageID, eventType, taskMessage)
	}

	if task.ParentID != "" && (taskState == common.Completed || taskState == common.Exception ||
		taskState == common.Killed || taskState == common.Cancelled || taskState == common.New) {
		err = ts.updateParentTask(ctx, taskID, taskStatus, taskState, task, payLoad)
		if err != nil {
			return err
		}
	}
	return err
}

// updateParentTask is used to update the status of parent task according to the status of child task
func (ts *TasksRPC) updateParentTask(ctx context.Context, taskID, taskStatus, taskState string, task *tmodel.Task, payLoad *taskproto.Payload) error {
	parentTask, err := ts.GetTaskStatusModel(ctx, task.ParentID, common.InMemory)
	if err != nil {
		return fmt.Errorf("error while retrieving the task details from db: " + err.Error())
	}

	if parentTask.PercentComplete == 100 {
		l.LogWithFields(ctx).Infof("Parent Task is already updated to 100 percent with the task state %s", parentTask.TaskState)
		return nil
	}

	if taskState != common.Completed && taskState != common.New {
		if parentTask.TaskFinalResponse != nil {
			var resp response.RPC
			json.Unmarshal(parentTask.TaskFinalResponse, &resp)
			parentTask.Payload.HTTPHeaders = resp.Header

		}
		errMsg := "One or more of the requests failed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + task.ParentID
		parentTask.TaskState = taskState
		parentTask.PercentComplete = 100
		parentTask.StatusCode = payLoad.StatusCode
		parentTask.Messages = []*tmodel.Message{{MessageID: response.Failure, Message: errMsg}}
		parentTask.TaskResponse = payLoad.ResponseBody
		l.LogWithFields(ctx).Debugf("Updating parent task %s with PercentComplete: %d, TaskState: %s and status code: %d",
			parentTask.ID, parentTask.PercentComplete, parentTask.TaskState, parentTask.StatusCode)
		ts.UpdateTaskQueue(parentTask)
		return fmt.Errorf(errMsg)
	}

	childIDs, err := ts.getAllChildTasks(ctx, task.ParentID)
	if err != nil {
		return err
	}
	l.LogWithFields(ctx).Debugf("Child ID's associated with parent task %s: %v", task.ParentID, childIDs)
	if len(childIDs) < 1 || (len(childIDs) == 1 && taskState == common.Completed && payLoad.StatusCode < http.StatusAccepted) {
		l.LogWithFields(ctx).Debugf("All tasks are completed ! Updating Parent task %s to completed state", parentTask.ID)
		parentTask.StatusCode = http.StatusOK
		if parentTask.TaskFinalResponse != nil {
			var resp response.RPC
			json.Unmarshal(parentTask.TaskFinalResponse, &resp)
			parentTask.Payload.HTTPHeaders = resp.Header
			parentTask.TaskFinalResponse = nil
			parentTask.StatusCode = http.StatusCreated
		}
		ts.updateTaskToCompleted(parentTask)
		return nil
	}

	return ts.validateChildTasksAndUpdateParentTask(ctx, childIDs, taskID, parentTask)
}

func (ts *TasksRPC) validateChildTasksAndUpdateParentTask(ctx context.Context, childIDs []string, taskID string, parentTask *tmodel.Task) error {
	s := make([]interface{}, 0, len(childIDs))
	for _, v := range childIDs {
		if v != taskID {
			s = append(s, "task:"+v)
		}
	}

	data, _ := ts.GetMultipleTaskKeysModel(ctx, s, common.InMemory)
	var isSuccess bool = true
	for _, subtask := range *data {
		if subtask.PercentComplete == 100 && subtask.TaskState == common.Suspended {
			l.LogWithFields(ctx).Debugf("updating sub task %s that made to suspended at 100 percent to completed state",
				subtask.ID)
			task := new(tmodel.Task)
			task.ID = subtask.ID
			task.Name = subtask.Name
			task.TaskMonitor = subtask.TaskMonitor
			task.URI = subtask.URI
			task.StartTime = subtask.StartTime
			task.EndTime = subtask.EndTime
			task.StatusCode = http.StatusOK
			ts.updateTaskToCompleted(task)
		} else if subtask.TaskState != common.Completed {
			isSuccess = false
			break
		}
	}

	if isSuccess {
		l.LogWithFields(ctx).Debugf("All tasks are completed ! Updating Parent task %s to completed state",
			parentTask.ID)
		parentTask.StatusCode = http.StatusOK
		if parentTask.TaskFinalResponse != nil {
			var resp response.RPC
			json.Unmarshal(parentTask.TaskFinalResponse, &resp)
			parentTask.Payload.HTTPHeaders = resp.Header
			parentTask.TaskFinalResponse = nil
			parentTask.StatusCode = http.StatusCreated
		}
		ts.updateTaskToCompleted(parentTask)
	}

	return nil
}

// updateTaskToCompleted update the task to completed state with success response
func (ts *TasksRPC) updateTaskToCompleted(task *tmodel.Task) {
	task.TaskState = common.Completed
	task.TaskStatus = common.OK
	task.PercentComplete = 100
	resp := tcommon.GetTaskResponse(task.StatusCode, response.Success)
	body, _ := json.Marshal(resp.Body)
	task.TaskResponse = body
	ts.UpdateTaskQueue(task)
}

// ProcessTaskEvents receive the task event from plugins
// The function will find out the ODIM task corresponding to the plugin task ID
// and task progress from the events
// Then the function update the ODIM task with the task progress received
func (ts *TasksRPC) ProcessTaskEvents(ctx context.Context, data interface{}) bool {
	event := data.(dmtf.EventRecord)
	var taskID string

	if len(event.MessageArgs) == 0 {
		l.LogWithFields(ctx).Error("task id is not present in the task event." +
			"skipping the task update")
		return false
	}

	taskID = event.MessageArgs[0]
	// get the plugin task information from DB which including ODIM task ID
	// plugin IP, and plugin task ID
	pluginTask, err := tmodel.GetPluginTaskInfo(taskID)
	if err != nil {
		l.LogWithFields(ctx).Error("error while processing task event :", err.Error())
		return false
	}

	messageID := event.MessageID
	var message string
	if strings.HasPrefix(messageID, common.TaskEventType) {
		res := strings.Split(messageID, common.TaskEventType+".")
		message = res[1]
	}

	if message == "" {
		l.LogWithFields(ctx).Errorf("Got invalid messageID for task event with task ID %s",
			taskID)
		return false
	}

	taskState := tcommon.TaskStatusMap[message]
	taskStatus := event.Severity

	var percentComplete int32
	switch taskState {
	case dmtf.TaskStateStarting:
		percentComplete = 0
	case dmtf.TaskStateRunning:
		pc, err := strconv.ParseInt(event.MessageArgs[1], 10, 32)
		if err != nil {
			l.LogWithFields(ctx).Errorf("Invalid percent complete received from task event: %v", event.MessageArgs[1])
			return false
		}
		percentComplete = int32(pc)
	case dmtf.TaskStateCompleted, dmtf.TaskStateCancelled,
		dmtf.TaskStateSuspended, dmtf.TaskStateInterrupted,
		dmtf.TaskStateKilled, dmtf.TaskStateException:
		percentComplete = 100
	}

	sc := event.MessageArgs[len(event.MessageArgs)-1]
	statusCode, err := strconv.ParseInt(sc, 10, 32)
	if err != nil {
		l.LogWithFields(ctx).Errorf("Invalid status code received from task event: %v", event.MessageArgs[1])
		return false
	}
	timestamp, err := time.Parse(time.RFC3339, event.EventTimestamp)
	if err != nil {
		timestamp = time.Now()
	}

	responseMessage := event.MessageArgs[len(event.MessageArgs)-2]
	resp := tcommon.GetTaskResponse(int32(statusCode), responseMessage)
	body, _ := json.Marshal(resp.Body)
	payLoad := &taskproto.Payload{
		StatusCode:   int32(statusCode),
		ResponseBody: body,
	}
	if strings.Contains(responseMessage, "location") && strings.Contains(responseMessage, "host") {
		var data tmodel.SubscriptionCreate
		err := json.Unmarshal([]byte(responseMessage), &data)
		if err == nil {
			go tcommon.UpdateSubscriptionLocation(ctx, data.Location, data.Host)
			body, _ = json.Marshal(data.Body)
		}
		header := make(map[string]string)
		header["location"] = data.Location
		payLoad.HTTPHeaders = header
	}
	if strings.Contains(responseMessage, "UpdateRemoteAccount") && strings.Contains(responseMessage, "host") {
		var data tmodel.UpdateAccount
		err := json.Unmarshal([]byte(responseMessage), &data)
		if err == nil {
			go tcommon.UpdateRemoteAccount(ctx, data.Location, data.Host)
			body, _ = json.Marshal(data.Body)
		}
	}

	l.LogWithFields(ctx).Debugf("Received task event from plugin for odim task %s, "+
		"plugin taskID: %s, taskState: %s, taskStatus: %s, percentComplete: %d, "+
		"status code: %d: response body: %s, end time: %v",
		pluginTask.OdimTaskID, taskID, taskState, taskStatus,
		percentComplete, statusCode, string(body), timestamp)

	err = ts.updateTaskUtil(context.TODO(), pluginTask.OdimTaskID,
		string(taskState), taskStatus, percentComplete,
		payLoad, timestamp)
	if err != nil {
		l.Log.Error("failed to update task: error while updating task: " +
			err.Error())
		return false
	}

	if percentComplete == 100 {
		tmodel.RemovePluginTaskID(context.TODO(), taskID)
	}

	return false
}
