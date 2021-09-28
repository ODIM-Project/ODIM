/*
 * Copyright (c) 2021 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dputilities

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	pluginConf "github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	// CompletedTaskIndex is a index name which is required for
	// to build index for completed tasks
	CompletedTaskIndex = "CompletedTaskIndex"
	//CompletedTaskTable is a Table name for Completed Task
	CompletedTaskTable = "CompletedTask"
)

type TaskState int32

const (
	Completed TaskState = iota
	Cancelling
	Cancelled
	Exception
	Interrupted
	New
	Pending
	Running
	Service
	Starting
	Stopping
	Suspended
)

func (ts TaskState) String() string {
	return [...]string{"Completed", "Cancelling", "Cancelled", "Exception", "Interrupted", "New", "Pending", "Running",
		"Service", "Starting", "Stopping", "Suspended"}[ts]
}

type TaskStatus int32

const (
	Critical TaskStatus = iota
	Ok
	Warning
)

func (ts TaskStatus) String() string {
	return [...]string{"Critical", "Ok", "Warning"}[ts]
}

// Task struct with TaskState and TaskStatus enums
type Task struct {
	ParentID        string
	ID              string
	URI             string
	UserName        string
	Name            string
	HidePayload     bool
	Payload         Payload
	PercentComplete int32
	TaskMonitor     string
	TaskState       TaskState
	TaskStatus      TaskStatus
	StatusCode      int32
	TaskResponse    []byte
	Messages        []*dpmodel.Message
	StartTime       time.Time
	EndTime         time.Time
}

// Task struct representing Task in database
type TaskDb struct {
	ParentID        string
	ID              string
	URI             string
	UserName        string
	Name            string
	HidePayload     bool
	Payload         Payload
	PercentComplete int32
	TaskMonitor     string
	TaskState       string
	TaskStatus      string
	StatusCode      int32
	TaskResponse    []byte
	Messages        []*dpmodel.Message
	StartTime       time.Time
	EndTime         time.Time
}

// Payload contain information detailing the HTTP and JSON payload
//information for executing the task.
//This object shall not be included in the response if the HidePayload property
// is set to True.
type Payload struct {
	HTTPHeaders   map[string]string `json:"HttpHeaders"`
	HTTPOperation string            `json:"HttpOperation"`
	JSONBody      string            `json:"JsonBody"`
	TargetURI     string            `json:"TargetUri"`
}

type TaskService interface {
	CreateTask() (string, error)
	UpdateTask(taskID, host string, taskState TaskState, taskStatus TaskStatus, percentComplete int32,
		payLoad *taskproto.Payload, endTime time.Time) error
	GetTaskState(state string) (TaskState, error)
	GetTaskStatus(status string) (TaskStatus, error)
}

type TaskServiceImpl struct {
}

func GetTaskService() *TaskServiceImpl {
	return &TaskServiceImpl{}
}

func (ts *TaskServiceImpl) CreateTask() (string, error) {
	userName := pluginConf.Data.PluginConf.UserName

	// Frame the model
	currentTime := time.Now()

	task := TaskDb{
		UserName:        userName,
		ID:              "task" + uuid.New().String(),
		TaskState:       New.String(),
		TaskStatus:      Ok.String(),
		PercentComplete: 0,
		StartTime:       currentTime,
		EndTime:         currentTime,
	}
	task.Name = "Task " + task.ID
	task.TaskMonitor = "/taskmon/" + task.ID
	task.URI = "/redfish/v1/TaskService/Tasks/" + task.ID

	// Persist in the in-memory DB
	err := persistTask(&task)
	if err != nil {
		log.Error("error while trying to insert the task details: " + err.Error())
		return "", err
	}
	// return the Task URI
	return "/redfish/v1/TaskService/Tasks/" + task.ID, err
}

// updateTaskUtil is a function to update the existing task and/or to create sub-task under a parent task.
// This function is to set task status, task end time along with task state based on the task state.
// Takes:
//	taskID - Is of type string, containes task ID of the task to updated
//	taskState - Is of type string, containes new sate of the task
//	taskStatus - Is of type string, containes new status of the task
//	endTime    - Is of type time.Time, containses the endtime of the task
// Retruns:
//	err of type error
//	nil - On Success
//	Non nil - On Failure
func (ts *TaskServiceImpl) UpdateTask(taskID, host string, taskState TaskState, taskStatus TaskStatus, percentComplete int32, payLoad *taskproto.Payload, endTime time.Time) error {
	// Retrieve the task details using taskID
	task, err := ts.getTaskFromDb(taskID)
	if err != nil {
		return fmt.Errorf("error while retrieving the task details from db: " + err.Error())
	}
	//If the task is already in cancelled state, then updates are not allowed to it.
	if task.TaskState == Cancelled || (task.TaskState == Cancelling && taskState != Cancelled) {
		return fmt.Errorf("task is already cancelled or being cancelling")
	}

	task.TaskState = taskState
	task.TaskStatus = taskStatus
	task.EndTime = endTime
	task.PercentComplete = percentComplete
	if payLoad != nil {
		task.Payload.HTTPOperation = payLoad.HTTPOperation
		task.Payload.HTTPHeaders = payLoad.HTTPHeaders
		task.Payload.JSONBody = payLoad.JSONBody
		task.Payload.TargetURI = payLoad.TargetURI
		task.StatusCode = payLoad.StatusCode
		task.TaskResponse = payLoad.ResponseBody
	}
	taskEventMessageID := "TaskEvent.1.0.1.Task" + taskState.String()

	// Update the task data in the InMemory DB
	err = updateTaskStatus(task)
	if err != nil {
		return fmt.Errorf("error while updating the task to In-memory DB: %v" + err.Error())
	}

	event := common.MessageData{
		OdataType: "#Event.v1_2_1.Event",
		Name:      "Task status changed",
		Context:   "/redfish/v1/$metadata#Event.Event",
		Events: []common.Event{
			{
				EventType:      "StatusChange",
				Severity:       "Ok",
				EventTimestamp: time.Now().String(),
				Message:        "Task updated successfully",
				MessageID:      taskEventMessageID,
				OriginOfCondition: &common.Link{
					Oid: task.URI,
				},
			},
		},
	}
	ManualEvents(event, host)
	return err
}

// PersistTask is to store the task data in db
// Takes:
//	t pointer to Task to be stored.
//	db of type common.DbType(int32)
func persistTask(t *TaskDb) error {
	PrepareDbConfig()
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("task", t.ID, t); err != nil {
		return fmt.Errorf("error while trying to create new task: %v", err.Error())
	}
	return nil
}

// UpdateTaskStatus is to update the task data already present in db
// Takes:
//	db of type common.DbType(int32)
//	t of type *Task
// Returns:
//	err of type error
//	On Success - return nil value
//	On Failure - return non nill value
func updateTaskStatus(t *Task) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if _, err = connPool.Update("task", t.ID, convertToDbTask(t)); err != nil {
		return fmt.Errorf("error while trying to update task: %v", err.Error())
	}
	// Build Redis Index here if we dont do it in thandle
	if t.TaskState == Completed && t.ParentID == "" {
		taskIndexErr := buildCompletedTaskIndex(t, CompletedTaskTable)
		if err != nil {
			return taskIndexErr
		}
	}
	return nil
}

// GetTaskStatus is to retrieve the task data already present in db
// Takes:
//	taskID of type string contains the task ID of the task to be retrieved from the db
//	db of type common.DbType(int32)
// Returns:
//	err of type error
//		On Success - return nil value
//		On Failure - return non nill value
//	t of type *Task implicitly valid only when error is nil
func (ts *TaskServiceImpl) getTaskFromDb(taskID string) (*Task, error) {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		log.Error("error while trying to get the db connection")
		return nil, fmt.Errorf("error while trying to connnect to DB: %v", err.Error())
	}
	taskData, err := connPool.Read("task", taskID)
	if err != nil {
		return nil, fmt.Errorf("error while trying to read from DB: %v", err.Error())
	}

	task, errs := ts.unmarshalTask([]byte(taskData))
	if errs != nil {
		return nil, fmt.Errorf("error while trying to unmarshal task data: %v", errs)
	}

	return task, nil
}

//BuildCompletedTaskIndex is used to build the index for Completed Task
func buildCompletedTaskIndex(completedTask *Task, table string) error {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	key := completedTask.UserName + "::" + completedTask.EndTime.String() + "::" + completedTask.ID
	createError := conn.CreateTaskIndex(CompletedTaskIndex, completedTask.EndTime.UnixNano(), key)
	if createError != nil {
		return fmt.Errorf("error while trying to create task index: %v", err)
	}
	return nil
}

func (ts *TaskServiceImpl) unmarshalTask(taskData []byte) (*Task, error) {
	dbTask := new(TaskDb)
	errorMsg := "error while trying to unmarshal task data: %v"
	if err := json.Unmarshal(taskData, &dbTask); err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}

	taskState, err := ts.GetTaskState(dbTask.TaskState)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}

	taskStatus, err := ts.GetTaskStatus(dbTask.TaskStatus)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}

	return &Task{
		ParentID:        dbTask.ParentID,
		ID:              dbTask.ID,
		URI:             dbTask.URI,
		UserName:        dbTask.UserName,
		Name:            dbTask.Name,
		HidePayload:     dbTask.HidePayload,
		Payload:         dbTask.Payload,
		PercentComplete: dbTask.PercentComplete,
		TaskMonitor:     dbTask.TaskMonitor,
		TaskState:       taskState,
		TaskStatus:      taskStatus,
		StatusCode:      dbTask.StatusCode,
		TaskResponse:    dbTask.TaskResponse,
		Messages:        dbTask.Messages,
		StartTime:       dbTask.StartTime,
		EndTime:         dbTask.EndTime,
	}, nil
}

func PrepareDbConfig() {
	config.Data.DBConf = &config.DBConf{
		Protocol:             pluginConf.Data.DBConf.Protocol,
		InMemoryHost:         pluginConf.Data.DBConf.InMemoryHost,
		InMemoryPort:         pluginConf.Data.DBConf.InMemoryPort,
		OnDiskHost:           pluginConf.Data.DBConf.OnDiskHost,
		OnDiskPort:           pluginConf.Data.DBConf.OnDiskPort,
		MaxIdleConns:         pluginConf.Data.DBConf.MaxIdleConns,
		MaxActiveConns:       pluginConf.Data.DBConf.MaxActiveConns,
		RedisHAEnabled:       pluginConf.Data.DBConf.RedisHAEnabled,
		InMemorySentinelPort: pluginConf.Data.DBConf.InMemorySentinelPort,
		InMemoryMasterSet:    pluginConf.Data.DBConf.InMemoryMasterSet,
		OnDiskMasterSet:      pluginConf.Data.DBConf.OnDiskMasterSet,
		OnDiskSentinelPort:   pluginConf.Data.DBConf.OnDiskSentinelPort,
	}
}

func (ts *TaskServiceImpl) GetTaskState(state string) (TaskState, error) {
	switch strings.ToLower(state) {
	case strings.ToLower(Completed.String()):
		return Completed, nil
	case strings.ToLower(Cancelling.String()):
		return Cancelling, nil
	case strings.ToLower(Cancelled.String()):
		return Cancelled, nil
	case strings.ToLower(Exception.String()):
		return Exception, nil
	case strings.ToLower(Interrupted.String()):
		return Interrupted, nil
	case strings.ToLower(New.String()):
		return New, nil
	case strings.ToLower(Pending.String()):
		return Pending, nil
	case strings.ToLower(Running.String()):
		return Running, nil
	case strings.ToLower(Service.String()):
		return Service, nil
	case strings.ToLower(Starting.String()):
		return Starting, nil
	case strings.ToLower(Stopping.String()):
		return Stopping, nil
	case strings.ToLower(Suspended.String()):
		return Suspended, nil
	default:
		return 0, fmt.Errorf("taskState not recognized")
	}
}

func (ts *TaskServiceImpl) GetTaskStatus(status string) (TaskStatus, error) {
	switch strings.ToLower(status) {
	case strings.ToLower(Critical.String()):
		return Critical, nil
	case strings.ToLower(Ok.String()):
		return Ok, nil
	case strings.ToLower(Warning.String()):
		return Warning, nil
	default:
		return 0, fmt.Errorf("taskStatus not recognized")
	}
}

func convertToDbTask(task *Task) *TaskDb {
	return &TaskDb{
		ParentID:        task.ParentID,
		ID:              task.ID,
		URI:             task.URI,
		UserName:        task.UserName,
		Name:            task.Name,
		HidePayload:     task.HidePayload,
		Payload:         task.Payload,
		PercentComplete: task.PercentComplete,
		TaskMonitor:     task.TaskMonitor,
		TaskState:       task.TaskState.String(),
		TaskStatus:      task.TaskStatus.String(),
		StatusCode:      task.StatusCode,
		TaskResponse:    task.TaskResponse,
		Messages:        task.Messages,
		StartTime:       task.StartTime,
		EndTime:         task.EndTime,
	}
}
