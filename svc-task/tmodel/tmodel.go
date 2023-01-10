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
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
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

// to be moved to dmtf
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

// Tick struct is used to help the goroutines that process the task queue to communicate effectively
// Tick contains the following attributes
/*
1. Ticker is of type Ticker in time package. it is used to acknowledge
the function that process task queue that it is time to commit the current
pipelined transaction to redis DB
2. M is of type Mutex in sync package. It ensures only one goroutine access
the Commit and Executing flags at the same time.
3. Commit is a flag which is made true when ticker "ticks". when it is made true,
"ProcessTaskQueue" commit the current pipeline to redis.
4. Executing is a flag which is made true when the "ProcessTaskQueue" function is invoked
and made false when it is finished.
*/
type Tick struct {
	Ticker    *time.Ticker
	M         sync.Mutex
	Commit    bool
	Executing bool
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

// GetWriteConnection returns write connection retrieved from the connection pool.
func GetWriteConnection() *db.Conn {
	connPool, err := db.GetDBConnection(db.InMemory)
	if err != nil {
		l.Log.Error(err.Error())
		return nil
	}

	conn, connErr := connPool.GetWriteConnection()
	if connErr != nil {
		l.Log.Error("ProcessTaskQueue : error while trying to get DB write Connection : " + connErr.Error())
		return nil
	}
	return conn
}

func validateDBConnection(conn *db.Conn) *db.Conn {
	if conn.IsBadConn() {
		conn.Close()
		return GetWriteConnection()
	}
	return conn
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
2."conn" is an instance of Conn struct in persistence manager library
*/
func (tick *Tick) ProcessTaskQueue(queue *chan *Task, conn *db.Conn) {

	defer func() {
		tick.M.Lock()
		tick.Commit = false
		tick.Executing = false
		tick.M.Unlock()
	}()

	const (
		MaxRetry int    = 3
		Table    string = "task"
	)

	var (
		i             int           = 0
		updatedTasks  bool          = false
		createdIndex  bool          = false
		mapSize       int           = config.Data.TaskQueueConf.QueueSize
		retryInterval time.Duration = time.Duration(config.Data.TaskQueueConf.RetryInterval) * time.Millisecond
	)

	tasks := make(map[string]interface{}, mapSize)
	completedTasks := make(map[string]int64, mapSize)

	if len(*queue) <= 0 {
		return
	}

	tick.M.Lock()
	tick.Executing = true
	tick.M.Unlock()

	conn = validateDBConnection(conn)

	for {
		task := dequeueTask(queue)

		if task != nil {
			saveID := Table + ":" + task.ID
			tasks[saveID] = task
			if (task.TaskState == "Completed" || task.TaskState == "Exception") && task.ParentID == "" {
				key := task.UserName + "::" + task.EndTime.String() + "::" + task.ID
				completedTasks[key] = task.EndTime.UnixNano()
			}
		}

		if tick.Commit {
			break
		}
	}

	if len(tasks) > 0 {
		for i < MaxRetry {
			if err := conn.UpdateTransaction(tasks); err != nil {
				if err.ErrNo() == errors.TimeoutError || db.IsRetriable(err) {
					time.Sleep(retryInterval)
					conn = validateDBConnection(conn)
				} else {
					l.Log.Error("ProcessTaskQueue() : task update transaction failed : " + err.Error())
					break
				}
				i++
			} else {
				updatedTasks = true
				break
			}
		}

		if !updatedTasks {
			for task := range tasks {
				l.Log.Errorf("Failed to update the task : %s", task)
			}
		}
	}

	if len(completedTasks) > 0 {
		i = 0
		for i < MaxRetry {
			if err := conn.CreateIndexTransaction(CompletedTaskIndex, completedTasks); err != nil {
				if err.ErrNo() == errors.TimeoutError || db.IsRetriable(err) {
					time.Sleep(retryInterval)
					conn = validateDBConnection(conn)
				} else {
					l.Log.Error("ProcessTaskQueue() : create index transaction failed : " + err.Error())
					break
				}
				i++
			} else {
				createdIndex = true
				break
			}
		}

		if !createdIndex {
			for task := range completedTasks {
				l.Log.Errorf("Failed to create index for the task : %s", task)
			}
		}
	}

	tasks = nil
	completedTasks = nil
}

// dequeueTask dequeue a task from channel and returns. If no elements is present in the queue it returns nil.
func dequeueTask(queue *chan *Task) *Task {
	if len(*queue) <= 0 {
		return nil
	}
	return <-*queue
}
