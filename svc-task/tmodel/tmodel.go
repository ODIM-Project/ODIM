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
	"sync"
	"time"

	db "github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
)

const (
	// CompletedTaskIndex is a index name which is required for
	// to build index for completed tasks
	CompletedTaskIndex = "CompletedTaskIndex"
	//CompletedTaskTable is a Table name for Completed Task
	CompletedTaskTable = "CompletedTask"
	SignalTaskName     = "SignalTask"
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

// GetCompletedTasksIndex Searches Complete Tasks in the db using secondary index with provided search Key
func GetCompletedTasksIndex(searchKey string) ([]string, error) {
	var taskData []string
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		l.Log.Error("GetCompletedTasksIndex : error while trying to get DB Connection : " + err.Error())
		return taskData, fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	list, getErr := conn.GetTaskList(CompletedTaskIndex, 0, -1)
	if getErr != nil && getErr.Error() != "no data with ID found" {
		l.Log.Error("GetCompletedTasksIndex : error while trying to get task list : " + getErr.Error())
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
		l.Log.Error("PersistTask : error while trying to get DB Connection : " + err.Error())
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("task", t.ID, t); err != nil {
		l.Log.Error("PersistTask : error while trying to create task : " + err.Error())
		return fmt.Errorf("error while trying to create new task: %v", err.Error())
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
		l.Log.Error("DeleteTaskFromDB : error while trying to get DB Connection : " + err.Error())
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Delete("task", t.ID); err != nil {
		l.Log.Error("DeleteTaskFromDB : Unable to delete task : " + err.Error())
		return fmt.Errorf("error while trying to delete the task: %v", err.Error())
	}
	return nil
}

//DeleteTaskIndex is used to delete the completed task index
//taskID is the ID with which the completed task index is deleted
func DeleteTaskIndex(taskID string) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		l.Log.Error("DeleteTaskIndex : error while trying to get DB Connection : " + err.Error())
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if delErr := connPool.Del(CompletedTaskIndex, taskID); delErr != nil {
		l.Log.Error("DeleteTaskIndex : Unable to delete task index: " + delErr.Error())
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
		l.Log.Error("GetTaskStatus : error while trying to get DB Connection : " + err.Error())
		return task, fmt.Errorf("error while trying to connnect to DB: %v", err.Error())
	}
	taskData, err = connPool.Read("task", taskID)
	if err != nil {
		l.Log.Error("GetTaskStatus : Unable to read taskdata from DB: " + err.Error())
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
		l.Log.Error("GetAllTaskKeys : error while trying to get DB Connection : " + err.Error())
		return nil, fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	taskKeys, err := connPool.GetAllDetails("task")
	if err != nil {
		l.Log.Error("GetAllTaskKeys : error while trying to get task key details from DB  : " + err.Error())
		return nil, fmt.Errorf("error while fetching data: %v", err.Error())
	}
	return taskKeys, nil
}

//Transaction - is for performing atomic oprations using optimitic locking
func Transaction(key string, cb func(string) error) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		l.Log.Error("Transaction : error while trying to get DB Connection : " + err.Error())
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Transaction(key, cb); err != nil {
		l.Log.Error("Transaction : Unable to perform transaction   : " + err.Error())
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
		l.Log.Error("ValidateTaskUserName : error while trying to get DB Connection : " + err.Error())
		return fmt.Errorf("error while trying to connecting to DB: %v", err)
	}
	// If the user not found in the db, below call sets err to non nil value
	if _, err = connPool.Read("User", userName); err != nil {
		l.Log.Error("ValidateTaskUserName : error while trying to read from the db : " + err.Error())
		return fmt.Errorf("error while trying to read from DB: %v", err.Error())
	}
	return nil
}

// ProcessTaskQueue dequeue the tasks details from queue and update DB using pipelined transaction
// the pipeline is committed when signal task is dequeued from the queue
// a signal task is enqueued by the caller once in a millisecond
/* ProcessTaskQueue takes the following keys as input:
1."queue" is a pointer to the channel which acts as the task queue
2."wg" is of type wait group which acknowledges the caller that process is finished
*/
func ProcessTaskQueue(queue *chan *Task, wg *sync.WaitGroup) {
	// validate if queue is empty. check the first item is signal task. Stops the process if it is.
	task, yes := isEmptyQueue(queue)
	if yes {
		wg.Done()
		return
	}

	connPool, err := db.GetDBConnection(db.InMemory)
	if err != nil {
		l.Log.Error("ProcessTaskQueue : error while trying to get DB Connection : " + err.Error())
		wg.Done()
		return
	}
	conn, connErr := connPool.GetWriteConnection()
	if connErr != nil {
		l.Log.Error("ProcessTaskQueue : error while trying to get DB write Connection : " + err.Error())
		wg.Done()
		return
	}

	conn.InitRedisPipeline()
	pipeErr := pipeRequests(conn, task)
	if pipeErr != nil {
		l.Log.Error(pipeErr)
		wg.Done()
		return
	}

	for task := range *queue {
		if task.Name == SignalTaskName {
			break
		}

		err := pipeRequests(conn, task)
		if err != nil {
			l.Log.Error(err)
			wg.Done()
			return
		}
	}

	commitErr := conn.CommitRedisPipeline()
	if err != nil {
		l.Log.Error("ProcessTaskQueue : error while trying to send task data into pipe : " + commitErr.Error())
		wg.Done()
		return
	}
	wg.Done()
}

// pipeRequests pipes the update task operation and create index operation to pipeline
/* pipeRequests takes the following keys as input:
1."conn" is an instance of Conn struct from persistance manager library
2."task" is task data dequeued
*/
func pipeRequests(conn *db.Conn, task *Task) error {
	table := "task"
	key := table + ":" + task.ID
	err := conn.PipeUpdateRequest(key, task)
	if err != nil {
		return fmt.Errorf("ProcessTaskQueue : error while trying to send task data into update task pipe : " + err.Error())
	}

	if (task.TaskState == "Completed" || task.TaskState == "Exception") && task.ParentID == "" {
		key := task.UserName + "::" + task.EndTime.String() + "::" + task.ID
		err := conn.PipeCreateIndex(CompletedTaskIndex, task.EndTime.UnixNano(), key)
		if err != nil {
			return fmt.Errorf("ProcessTaskQueue : error while trying to send task data into create task index pipe : " + err.Error())
		}
	}
	return nil
}

// isEmptyQueue validates the queue is empty by checking the first value in queue is signal task
/* isEmptyQueue takes the following keys as input:
1."queue" is a pointer to the channel which acts as the task queue
*/
func isEmptyQueue(queue *chan *Task) (*Task, bool) {
	task := <-*queue
	if task.Name == SignalTaskName {
		return nil, true
	}
	return task, false
}
