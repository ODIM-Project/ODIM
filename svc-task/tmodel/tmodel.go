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

//Package tmodel ...
package tmodel

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
)

const (
	// CompletedTaskIndex is a index name which is required for
	// to build index for completed tasks
	CompletedTaskIndex = "CompletedTaskIndex"
	//CompletedTaskTable is a Table name for Completed Task
	CompletedTaskTable = "CompletedTask"
)

//CompletedTask is used to build index for redis
type CompletedTask struct {
	UserName string
	ID       string
	EndTime  int64
}

// Task Model
type Task struct {
	ParentID     string
	ChildTaskIDs []string
	ID           string
	URI          string
	UserName     string
	Name         string
	HidePayload  bool
	Payload      Payload
	/*The value of this property shall indicate the completion progress of
	the task, reported in percent of completion.
	If the task has not been started, the value shall be zero.
	*/
	PercentComplete int32
	TaskMonitor     string
	TaskState       string
	TaskStatus      string
	StatusCode      int32
	TaskResponse    []byte
	Messages        []*Message // Its there in the spec, how are we going to use it
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

// Message Model
type Message struct {
	Message           string   `json:"Message"`
	MessageID         string   `json:"MessageId"`
	MessageArgs       []string `json:"MessageArgs"`
	Oem               Oem      `json:"Oem"`
	RelatedProperties []string `json:"RelatedProperties"`
	Resolution        string   `json:"Resolution"`
	Severity          string   `json:"Severity"`
}

//BuildCompletedTaskIndex is used to build the index for Completed Task
func BuildCompletedTaskIndex(completedTask *Task, table string) error {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		log.Error("BuildCompletedTaskIndex : error while trying to get DB Connection : " + err.Error())
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	key := completedTask.UserName + "::" + completedTask.EndTime.String() + "::" + completedTask.ID
	createError := conn.CreateTaskIndex(CompletedTaskIndex, completedTask.EndTime.UnixNano(), key)
	if createError != nil {
		log.Error("BuildCompletedTaskIndex : error while trying to CreateTaskIndex : " + createError.Error())
		return fmt.Errorf("error while trying to create task index: %v", err)
	}
	return nil
}

// GetCompletedTasksIndex Searches Complete Tasks in the db using secondary index with provided search Key
func GetCompletedTasksIndex(searchKey string) ([]string, error) {
	var taskData []string
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		log.Error("GetCompletedTasksIndex : error while trying to get DB Connection : " + err.Error())
		return taskData, fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	list, getErr := conn.GetTaskList(CompletedTaskIndex, 0, -1)
	if getErr != nil && getErr.Error() != "no data with ID found" {
		log.Error("GetCompletedTasksIndex : error while trying to get task list : " + getErr.Error())
		return taskData, nil
	}
	taskData = list
	return taskData, nil
}

// Oem Model
type Oem struct {
}

// PersistTask is to store the task data in db
// Takes:
//	t pointer to Task to be stored.
//	db of type common.DbType(int32)
func PersistTask(t *Task, db common.DbType) error {
	connPool, err := common.GetDBConnection(db)
	if err != nil {
		log.Error("PersistTask : error while trying to get DB Connection : " + err.Error())
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("task", t.ID, t); err != nil {
		log.Error("PersistTask : error while trying to create task : " + err.Error())
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
func UpdateTaskStatus(t *Task, db common.DbType) error {
	connPool, err := common.GetDBConnection(db)
	if err != nil {
		log.Error("UpdateTaskStatus : error while trying to get DB Connection : " + err.Error())
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if _, err = connPool.Update("task", t.ID, t); err != nil {
		log.Error("UpdateTaskStatus : error while trying to updating task status : " + err.Error())
		return fmt.Errorf("error while trying to update task: %v", err.Error())
	}
	// Build Redis Index here if we dont do it in thandle
	if (t.TaskState == "Completed" || t.TaskState == "Exception") && t.ParentID == "" {
		taskIndexErr := BuildCompletedTaskIndex(t, CompletedTaskTable)
		if taskIndexErr != nil {
			log.Error("UpdateTaskStatus : error in creating index for task : " + taskIndexErr.Error())
			return taskIndexErr
		}
	}
	return nil
}

// DeleteTaskFromDB is to delete the task from db
// Takes:
// 	t of type pointer to Task object
// Returns:
//      err of type error
//      On Success - return nil value
//      On Failure - return non nill value
func DeleteTaskFromDB(t *Task) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		log.Error("DeleteTaskFromDB : error while trying to get DB Connection : " + err.Error())
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Delete("task", t.ID); err != nil {
		log.Error("DeleteTaskFromDB : Unable to delete task : " + err.Error())
		return fmt.Errorf("error while trying to delete the task: %v", err.Error())
	}
	return nil
}

//DeleteTaskIndex is used to delete the completed task index
//taskID is the ID with which the completed task index is deleted
func DeleteTaskIndex(taskID string) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		log.Error("DeleteTaskIndex : error while trying to get DB Connection : " + err.Error())
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if delErr := connPool.Del(CompletedTaskIndex, taskID); delErr != nil {
		log.Error("DeleteTaskIndex : Unable to delete task index: " + delErr.Error())
		return fmt.Errorf("error while trying to delete the completed task index: %v", delErr.Error())
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
func GetTaskStatus(taskID string, db common.DbType) (*Task, error) {
	task := new(Task)
	var taskData string
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		log.Error("GetTaskStatus : error while trying to get DB Connection : " + err.Error())
		return task, fmt.Errorf("error while trying to connnect to DB: %v", err.Error())
	}
	taskData, err = connPool.Read("task", taskID)
	if err != nil {
		log.Error("GetTaskStatus : Unable to read taskdata from DB: " + err.Error())
		return task, fmt.Errorf("error while trying to read from DB: %v", err.Error())
	}
	if errs := json.Unmarshal([]byte(taskData), task); errs != nil {
		return task, fmt.Errorf("error while trying to unmarshal task data: %v", errs)
	}
	return task, nil
}

// GetAllTaskKeys will collect all task keys available in the DB
//Takes:
//	None
//Returns:
//	Slice of type strings and error
//	On Success - error is set to nil and returns slice of tasks
//	On Failure - error is set to appropriate reason why it got failed
//	and slice of task is set to nil
func GetAllTaskKeys() ([]string, error) {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		log.Error("GetAllTaskKeys : error while trying to get DB Connection : " + err.Error())
		return nil, fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	taskKeys, err := connPool.GetAllDetails("task")
	if err != nil {
		log.Error("GetAllTaskKeys : error while trying to get task key details from DB  : " + err.Error())
		return nil, fmt.Errorf("error while fetching data: %v", err.Error())
	}
	return taskKeys, nil
}

//Transaction - is for performing atomic oprations using optimitic locking
func Transaction(key string, cb func(string) error) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		log.Error("Transaction : error while trying to get DB Connection : " + err.Error())
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Transaction(key, cb); err != nil {
		log.Error("Transaction : Unable to perform transaction   : " + err.Error())
		return fmt.Errorf("error while performing transaction: %v", err.Error())
	}
	return nil
}

// ValidateTaskUserName validates the username.
// Returns error with non nil value if username is not found in the db,
// if username found in the db error is set to nil.
func ValidateTaskUserName(userName string) error {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		log.Error("ValidateTaskUserName : error while trying to get DB Connection : " + err.Error())
		return fmt.Errorf("error while trying to connecting to DB: %v", err)
	}
	// If the user not found in the db, below call sets err to non nil value
	if _, err = connPool.Read("User", userName); err != nil {
		log.Error("ValidateTaskUserName : error while trying to read from the db : " + err.Error())
		return fmt.Errorf("error while trying to read from DB: %v", err.Error())
	}
	return nil
}
