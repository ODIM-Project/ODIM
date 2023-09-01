// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package tmodel

import (
	"context"
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"
	"time"

	db "github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/satori/uuid"
	"golang.org/x/crypto/sha3"
)

type MockRedisConn struct {
	MockClose   func() error
	MockErr     func() error
	MockDo      func(string, ...interface{}) (interface{}, error)
	MockSend    func(string, ...interface{}) error
	MockFlush   func() error
	MockReceive func() (interface{}, error)
}

func (mc MockRedisConn) Close() error {
	return mc.MockClose()
}

func (mc MockRedisConn) Err() error {
	return mc.MockErr()
}

func (mc MockRedisConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	return mc.MockDo(commandName, args...)
}

func (mc MockRedisConn) Send(commandName string, args ...interface{}) error {
	return mc.MockSend(commandName, args...)
}

func (mc MockRedisConn) Flush() error {
	return mc.MockFlush()
}

func (mc MockRedisConn) Receive() (interface{}, error) {
	return mc.MockReceive()
}

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
	config.SetUpMockConfig(t)
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
	err := PersistTask(mockContext(), &task, 23)
	if err == nil {
		t.Fatalf("error: expected error here but got no error")
		return
	}
	// Persist in the in-memory DB
	err = PersistTask(mockContext(), &task, common.InMemory)
	if err != nil {
		t.Fatalf("error while trying to insert the task details: %v", err)
		return
	}
}

func TestProcessTaskQueue(t *testing.T) {
	c, er := db.MockDBConnection(t)
	if er != nil {
		t.Fatal(er)
	}
	queue := make(chan *Task, 10)
	config.SetUpMockConfig(t)
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

	err := PersistTask(mockContext(), &task, common.InMemory)
	if err != nil {
		t.Fatalf("error while trying to insert the task details: %v", err)
		return
	}

	task1, err := GetTaskStatus(mockContext(), task.ID, common.InMemory)
	if err != nil {
		t.Fatalf("error while retrieving the Task details with Get: %v", err)
		return
	}

	type args struct {
		tasks map[string]interface{}
		conn  *db.Conn
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success case",
			args: args{
				tasks: make(map[string]interface{}),
				conn: &db.Conn{
					WriteConn: c.ReadPool,
				},
			},
		},
		{
			name: "error case 1: no retry",
			args: args{
				tasks: make(map[string]interface{}),
				conn: &db.Conn{
					WriteConn: c.ReadPool,
				},
			},
		},
		{
			name: "error case 2 : retry",
			args: args{
				tasks: make(map[string]interface{}),
				conn: &db.Conn{
					WriteConn: c.ReadPool,
				},
			},
		},
		{
			name: "error case 3 : bad connection",
			args: args{
				tasks: make(map[string]interface{}),
				conn: &db.Conn{
					WriteConn: c.ReadPool,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue <- task1
			tick := &Tick{
				Executing: true,
				Commit:    true,
			}
			go tick.ProcessTaskQueue(&queue, tt.args.conn)
			for {
				if !tick.Executing {
					break
				}
			}
		})
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
	err := PersistTask(mockContext(), &task, common.InMemory)
	if err != nil {
		t.Fatalf("error while trying to insert the task details: %v", err)
		return
	}

	task1, err := GetTaskStatus(mockContext(), task.ID, common.InMemory)
	if err != nil {
		t.Fatalf("error while retreving the Task details with Get: %v", err)
		return
	}
	// Negetive test case
	task2 := new(Task)
	err = DeleteTaskFromDB(mockContext(), task2)
	if err == nil {
		t.Fatalf("error: expected error here but got no error ")
		return
	}
	// Positive Test case
	err = DeleteTaskFromDB(mockContext(), task1)
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
	err := PersistTask(mockContext(), &task, common.InMemory)
	if err != nil {
		t.Fatalf("error while trying to insert the task details: %v", err)
		return
	}
	// Negetive Test case with wrong task ID
	_, err = GetTaskStatus(mockContext(), "", common.InMemory)
	if err == nil {
		t.Fatalf("error: expected error here but got no error")
		return
	}
	// Positive test case
	_, err = GetTaskStatus(mockContext(), task.ID, common.InMemory)
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
	taskList, err := GetAllTaskKeys(mockContext())
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
	err = PersistTask(mockContext(), &task, common.InMemory)
	if err != nil {
		t.Fatalf("error while trying to insert the task details: %v", err)
		return
	}
}
func mockCallBack(ctx context.Context, key string) error {
	if key != "validKey" {
		return fmt.Errorf("error invalid key")
	}
	return nil
}
func TestTransaction(t *testing.T) {
	defer flushDB(t)
	type args struct {
		key string
		cb  func(context.Context, string) error
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
			if err := Transaction(mockContext(), tt.args.key, tt.args.cb); !reflect.DeepEqual(err, tt.wantErr) {
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
			if err := ValidateTaskUserName(mockContext(), tt.args.userName); !reflect.DeepEqual(err, tt.wantErr) {
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

func mockPluginTaskInDB(odimTaskID string, pluginTaskMon string) error {
	table := "PluginTask"
	taskData := &common.PluginTask{
		IP:               "127.0.0.1",
		OdimTaskID:       odimTaskID,
		PluginTaskMonURL: pluginTaskMon,
	}
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to connecting"+
			" to DB: ", err.Error())
	}

	if err = connPool.Create(table, pluginTaskMon, taskData); err != nil {
		return errors.PackError(err.ErrNo(), "error wlhile trying to insert"+
			" plugin task: ", err.Error())
	}
	return nil
}

func TestGetPluginTaskInfo(t *testing.T) {
	defer flushDB(t)
	err := config.SetUpMockConfig(t)
	if err != nil {
		t.Fatalf("fatal: error while trying to collect mock db config: %v", err)
		return
	}

	odimTaskID := "task" + uuid.NewV4().String()
	pluginTaskID := "task" + uuid.NewV4().String()
	err = mockPluginTaskInDB(odimTaskID, pluginTaskID)
	if err != nil {
		t.Fatalf("fatal: error while trying to insert plugin task data in DB: %v", err)
		return
	}

	type args struct {
		taskID string
	}
	tests := []struct {
		name    string
		args    args
		want    common.PluginTask
		wantErr bool
	}{
		{

			name: "get plugin task info from db",
			args: args{
				taskID: pluginTaskID,
			},
			want: common.PluginTask{
				IP:               "127.0.0.1",
				OdimTaskID:       odimTaskID,
				PluginTaskMonURL: pluginTaskID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPluginTaskInfo(tt.args.taskID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPluginTaskInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("GetPluginTaskInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
