package tqueue

import (
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/svc-task/tmodel"
	"github.com/satori/uuid"
)

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

func TestUpdateTasksStatus(t *testing.T) {
	config.SetUpMockConfig(t)
	defer flushDB(t)
	task := tmodel.Task{
		UserName:        "admin",
		ParentID:        "",
		ChildTaskIDs:    nil,
		ID:              "task" + uuid.NewV4().String(),
		TaskState:       "New",
		TaskStatus:      "OK",
		PercentComplete: 0,
		StartTime:       time.Now(),
		EndTime:         time.Time{},
	}
	task.Name = "Task " + task.ID

	err := tmodel.PersistTask(&task, common.InMemory)
	if err != nil {
		t.Fatalf("error while trying to insert the task details: %v", err)
		return
	}
	task1, err := tmodel.GetTaskStatus(task.ID, common.InMemory)
	if err != nil {
		t.Fatalf("error while retreving the Task details with Get: %v", err)
		return
	}
	// creating channels
	NewTaskQueue(10)

	type args struct {
		size            int
		d               time.Duration
		taskState       string
		percentComplete int32
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "update task progress",
			args: args{
				size:            10,
				d:               time.Millisecond,
				taskState:       "Completed",
				percentComplete: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task1.TaskState = tt.args.taskState
			task1.PercentComplete = tt.args.percentComplete
			EnqueueTask(task1)
			go UpdateTasksStatus(tt.args.size, tt.args.d)
			time.Sleep(time.Second)
			task2, err := tmodel.GetTaskStatus(task.ID, common.InMemory)
			if err != nil {
				t.Fatalf("error while retrieving the Task details with Get: %v", err)
				return
			}
			if task2.TaskState != tt.args.taskState {
				t.Errorf("UpdateTasksStatus() TaskState:  got = %v, want = %v", task2.TaskState, tt.args.taskState)
				return
			}
			if task2.PercentComplete != tt.args.percentComplete {
				t.Errorf("UpdateTasksStatus() PercentComplete:  got = %v, want = %v", task2.PercentComplete, tt.args.percentComplete)
				return
			}
			if task2.TaskState == "Completed" {
				tasks, err := tmodel.GetCompletedTasksIndex(task2.UserName)
				if err != nil {
					t.Fatalf("error while retrieving the completed Task index with Get: %v", err)
					return
				}
				key := task2.UserName + "::" + task2.EndTime.String() + "::" + task2.ID
				if len(tasks) != 1 && tasks[0] != key {
					t.Fatalf("index for completed task is not created")
				}
			}
		})
	}
}
