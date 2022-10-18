package tqueue

import (
	"fmt"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/svc-task/tmodel"
)

// taskQueue is a struct that contains channels for updating tasks statuses and creating index for completed tasks
type taskQueue struct {
	tasksInProgressQueue chan *tmodel.Task
	completedTasksQueue  chan *tmodel.Task
}

var TaskQueue *taskQueue

// NewTaskQueue creates an instance of taskQueue if it is not already created.
func NewTaskQueue(size int) {
	if TaskQueue != nil {
		return
	}
	TaskQueue = &taskQueue{
		tasksInProgressQueue: make(chan *tmodel.Task, size),
		completedTasksQueue:  make(chan *tmodel.Task, size),
	}
}

// EnqueueTasks enqueue the update task requests in the channel which act as a queue.
// if the task status is completed or exception it is also added to the completed tasks queue
func EnqueueTasks(task *tmodel.Task) {
	TaskQueue.tasksInProgressQueue <- task
	if (task.TaskState == "Completed" || task.TaskState == "Exception") && task.ParentID == "" {
		TaskQueue.completedTasksQueue <- task
	}
}

// dequeueTasksInProgress return the task first in and which are in progress
func dequeueTasksInProgress() (task *tmodel.Task) {
	return <-TaskQueue.tasksInProgressQueue
}

// dequeueTasksInProgress return the task first in and which are completed
func dequeueCompletedTasks() (task *tmodel.Task) {
	return <-TaskQueue.completedTasksQueue
}

// UpdateTasksStatus will get the task data which needs to be updated in DB from the queue and will update the DB once in an interval.
func UpdateTasksStatus(size int, d time.Duration) {
	table := "task"
	tasks := make(map[string]interface{}, size)
	ticker := time.NewTicker(d)

	for {
		select {
		case <-ticker.C:
			err := tmodel.UpdateTaskProgress(tasks, common.InMemory)
			if err != nil {
				l.Log.Error(fmt.Sprintf("error : updating tasks failed - %s", err.Error()))
				break
			}
			for k := range tasks {
				delete(tasks, k)
			}
		default:
			task := dequeueTasksInProgress()
			saveID := table + ":" + task.ID
			tasks[saveID] = task
		}
	}
}

// BuildCompletedTaskIndex create index for the tasks that are completed in redis DB once in an interval
func BuildCompletedTaskIndex(size int, d time.Duration) {
	tasks := make(map[string][2]interface{}, size)
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			err := tmodel.CreateCompletedTasksIndices(tasks, common.InMemory)
			if err != nil {
				l.Log.Error(fmt.Sprintf("error : creating indices for completed tasks failed - %s", err.Error()))
				break
			}
			for k := range tasks {
				delete(tasks, k)
			}
		default:
			task := dequeueCompletedTasks()
			saveID := task.UserName + "::" + task.EndTime.String() + "::" + task.ID
			value := [2]interface{}{
				tmodel.CompletedTaskIndex,
				task.EndTime.UnixNano(),
			}
			tasks[saveID] = value
		}
	}
}
