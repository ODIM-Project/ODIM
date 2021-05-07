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
package tmodel

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/satori/uuid"
	"golang.org/x/crypto/sha3"
)

type user struct {
	UserName string `json:"UserName"`
	Password string `json:"Password"`
	RoleID   string `json:"RoleId"`
}

func createMockUser(username, roleID string) error {
	hash := sha3.New512()
	hash.Write([]byte("Password"))
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

func TestPersistTask(t *testing.T) {
	common.SetUpMockConfig()
	defer flushDB(t)

	task := Task{
		UserName:     "admin",
		ParentID:     "",
		ChildTaskIDs: nil,
		ID:           "task" + uuid.NewV4().String(),
		TaskState:    "New",
		TaskStatus:   "OK",
		StartTime:    time.Now(),
		EndTime:      time.Time{},
	}
	task.Name = "Task " + task.ID
	// Negetive Test case
	err := PersistTask(&task, 23)
	if err == nil {
		t.Fatalf("error: expected error here but got no error")
		return
	}
	// Persist in the in-memory DB
	err = PersistTask(&task, common.InMemory)
	if err != nil {
		t.Fatalf("error while trying to insert the task details: %v", err)
		return
	}
}

func TestUpdateTaskStatus(t *testing.T) {
	common.SetUpMockConfig()
	defer flushDB(t)
	task := Task{
		UserName:     "admin",
		ParentID:     "",
		ChildTaskIDs: nil,
		ID:           "task" + uuid.NewV4().String(),
		TaskState:    "New",
		TaskStatus:   "OK",
		StartTime:    time.Now(),
		EndTime:      time.Time{},
	}
	task.Name = "Task " + task.ID
	// Persist in the in-memory DB
	err := PersistTask(&task, common.InMemory)
	if err != nil {
		t.Fatalf("error while trying to insert the task details: %v", err)
		return
	}
	task1 := new(Task)
	task1, err = GetTaskStatus(task.ID, common.InMemory)
	if err != nil {
		t.Fatalf("error while retreving the Task details with Get: %v", err)
		return
	}
	// Positive Test Case
	task1.TaskState = "Running"
	err = UpdateTaskStatus(task1, common.InMemory)
	if err != nil {
		t.Fatalf("error while updating the task details in the db: %v", err)
		return
	}
	task1.TaskState = "Completed"
	// Negetive test case
	err = UpdateTaskStatus(task1, 23)
	if err == nil {
		t.Fatalf("error: expected error here but got no error ")
		return
	}
	// Positive Test case
	err = UpdateTaskStatus(task1, common.InMemory)
	if err != nil {
		t.Fatalf("error while updating the task details in the db: %v", err)
		return
	}
}

func TestGetCompletedTasksIndex(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	task := Task{
		UserName:     "admin",
		ParentID:     "",
		ChildTaskIDs: nil,
		ID:           "task" + uuid.NewV4().String(),
		TaskState:    "Completed",
		TaskStatus:   "OK",
		StartTime:    time.Now(),
		EndTime:      time.Time{},
	}
	task.Name = "Task " + task.ID
	// Persist in the in-memory DB
	err := PersistTask(&task, common.InMemory)
	if err != nil {
		t.Fatalf("error while trying to insert the task details: %v", err)
		return
	}

	task1 := new(Task)
	task1, err = GetTaskStatus(task.ID, common.InMemory)
	if err != nil {
		t.Fatalf("error while retreving the Task details with Get: %v", err)
		return
	}
	err = UpdateTaskStatus(task1, common.InMemory)
	if err != nil {
		t.Fatalf("error while retreving the Task details with Get: %v", err)
		return
	}
	_, err = GetCompletedTasksIndex("task1")
	if err != nil {
		t.Fatalf("error while getting the task details in the db: %v", err)
		return
	}
	task2 := Task{
		UserName:     "admin",
		ParentID:     "",
		ChildTaskIDs: nil,
		ID:           "task" + uuid.NewV4().String(),
		TaskState:    "Exception",
		TaskStatus:   "OK",
		StartTime:    time.Now(),
		EndTime:      time.Time{},
	}
	task2.Name = "Task " + task.ID
	// Persist in the in-memory DB
	err = PersistTask(&task2, common.InMemory)
	if err != nil {
		t.Fatalf("error while trying to insert the task details: %v", err)
		return
	}

	task3 := new(Task)
	task3, err = GetTaskStatus(task.ID, common.InMemory)
	if err != nil {
		t.Fatalf("error while retreving the Task details with Get: %v", err)
		return
	}
	err = UpdateTaskStatus(task3, common.InMemory)
	if err != nil {
		t.Fatalf("error while retreving the Task details with Get: %v", err)
		return
	}
	_, err = GetCompletedTasksIndex("task3")
	if err != nil {
		t.Fatalf("error while getting the task details in the db: %v", err)
		return
	}
}

func TestDeleteTaskFromDB(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	task := Task{
		UserName:     "admin",
		ParentID:     "",
		ChildTaskIDs: nil,
		ID:           "task" + uuid.NewV4().String(),
		TaskState:    "New",
		TaskStatus:   "OK",
		StartTime:    time.Now(),
		EndTime:      time.Time{},
	}
	task.Name = "Task " + task.ID
	// Persist in the in-memory DB
	err := PersistTask(&task, common.InMemory)
	if err != nil {
		t.Fatalf("error while trying to insert the task details: %v", err)
		return
	}
	task1 := new(Task)
	task1, err = GetTaskStatus(task.ID, common.InMemory)
	if err != nil {
		t.Fatalf("error while retreving the Task details with Get: %v", err)
		return
	}
	// Negetive test case
	task2 := new(Task)
	err = DeleteTaskFromDB(task2)
	if err == nil {
		t.Fatalf("error: expected error here but got no error ")
		return
	}
	// Positive Test case
	err = DeleteTaskFromDB(task1)
	if err != nil {
		t.Fatalf("error while deleting the task details in the db: %v", err)
		return
	}
}

func TestDeleteTaskIndex(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	task := Task{
		UserName:     "admin",
		ParentID:     "",
		ChildTaskIDs: nil,
		ID:           "task" + uuid.NewV4().String(),
		TaskState:    "Completed",
		TaskStatus:   "OK",
		StartTime:    time.Now(),
		EndTime:      time.Time{},
	}
	task.Name = "Task " + task.ID
	// Persist in the in-memory DB
	err := PersistTask(&task, common.InMemory)
	if err != nil {
		t.Fatalf("error while trying to insert the task details: %v", err)
		return
	}
	task2 := Task{
		UserName:     "admin",
		ParentID:     "parentID",
		ChildTaskIDs: []string{"chidTask1"},
		ID:           "task" + uuid.NewV4().String(),
		TaskState:    "Completed",
		TaskStatus:   "OK",
		StartTime:    time.Now(),
		EndTime:      time.Time{},
	}
	task2.Name = "Task " + task2.ID
	// Persist in the in-memory DB
	err = PersistTask(&task2, common.InMemory)
	if err != nil {
		t.Fatalf("error while trying to insert the task details: %v", err)
		return
	}
	task1 := new(Task)
	task1, err = GetTaskStatus(task.ID, common.InMemory)
	if err != nil {
		t.Fatalf("error while retreving the Task details with Get: %v", err)
		return
	}
	err = UpdateTaskStatus(task1, common.InMemory)
	if err != nil {
		t.Fatalf("error while retreving the Task details with Get: %v", err)
		return
	}
	// Positive Test case
	err = DeleteTaskIndex(task1.ID)
	if err != nil {
		t.Fatalf("error while deleting the task details in the db: %v", err)
		return
	}
}

func TestGetTaskStatus(t *testing.T) {
	defer flushDB(t)
	common.SetUpMockConfig()
	task := Task{
		UserName:     "admin",
		ParentID:     "",
		ChildTaskIDs: nil,
		ID:           "task" + uuid.NewV4().String(),
		TaskState:    "New",
		TaskStatus:   "OK",
		StartTime:    time.Now(),
		EndTime:      time.Time{},
	}
	task.Name = "Task " + task.ID
	// Persist in the in-memory DB
	err := PersistTask(&task, common.InMemory)
	if err != nil {
		t.Fatalf("error while trying to insert the task details: %v", err)
		return
	}
	// Negetive Test case with wrong task ID
	_, err = GetTaskStatus("", common.InMemory)
	if err == nil {
		t.Fatalf("error: expected error here but got no error")
		return
	}
	// Positive test case
	_, err = GetTaskStatus(task.ID, common.InMemory)
	if err != nil {
		t.Fatalf("error while retreving the Task details with Get: %v", err)
		return
	}
}

func TestGetAllTaskKeys(t *testing.T) {
	count := 10
	flushDB(t)
	defer flushDB(t)
	for i := 0; i < count; i++ {
		createTaskInDB(t)
	}
	taskList, err := GetAllTaskKeys()
	if err != nil {
		t.Fatalf("fatal: error while fetching all task in db: %v", err)
	}
	if len(taskList) != count {
		t.Fatalf("fatal: error failed to got only %v tasks out of %v", len(taskList), count)
	}
}

func createTaskInDB(t *testing.T) {
	err := common.SetUpMockConfig()
	if err != nil {
		t.Fatalf("fatal: error while trying to collect mock db config: %v", err)
		return
	}
	task := Task{
		UserName:     "admin",
		ParentID:     "",
		ChildTaskIDs: nil,
		ID:           "task" + uuid.NewV4().String(),
		TaskState:    "New",
		TaskStatus:   "OK",
		StartTime:    time.Now(),
		EndTime:      time.Time{},
	}
	task.Name = "Task " + task.ID
	// Persist in the in-memory DB
	err = PersistTask(&task, common.InMemory)
	if err != nil {
		t.Fatalf("error while trying to insert the task details: %v", err)
		return
	}
}
func mockCallBack(key string) error {
	if key != "validKey" {
		return fmt.Errorf("error invalid key")
	}
	return nil
}
func TestTransaction(t *testing.T) {
	defer flushDB(t)
	type args struct {
		key string
		cb  func(string) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "Positive cases: All is well, with valid key",
			args: args{
				key: "validKey",
				cb:  mockCallBack,
			},
			wantErr: nil,
		},
		{
			name: "Positive cases: All is well, with inValid key",
			args: args{
				key: "inValidKey",
				cb:  mockCallBack,
			},
			wantErr: fmt.Errorf("error while performing transaction: error invalid key"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Transaction(tt.args.key, tt.args.cb); !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Transaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestValidateTaskUserName(t *testing.T) {

	defer flushDB(t)
	createMockUser("admin", "ADMIN")
	createMockUser("monitor", "MONITOR")
	createMockUser("client", "CLIENT")
	type args struct {
		userName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Positive Case: All is well, with Admin user",
			args: args{
				userName: "admin",
			},
			wantErr: nil,
		},
		{
			name: "Positive Case: All is well, with Admin user",
			args: args{
				userName: "monitor",
			},
			wantErr: nil,
		},
		{
			name: "Positive Case: All is well, with Admin user",
			args: args{
				userName: "client",
			},
			wantErr: nil,
		},
		{
			name: "Negative Case: with non existing user",
			args: args{
				userName: "unknown",
			},
			wantErr: fmt.Errorf("error while trying to read from DB: no data with the with key unknown found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateTaskUserName(tt.args.userName); !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("ValidateTaskUserName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func flushDB(t *testing.T) {
	err := common.TruncateDB(common.OnDisk)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = common.TruncateDB(common.InMemory)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}
