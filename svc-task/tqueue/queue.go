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
	queue chan *tmodel.Task
}

var TaskQueue *taskQueue

// NewTaskQueue creates an instance of taskQueue if it is not already created.
func NewTaskQueue(size int) {
	if TaskQueue != nil {
		return
	}
	TaskQueue = &taskQueue{
		queue: make(chan *tmodel.Task, size),
	}
}

// EnqueueTasks enqueue the update task requests in the channel which act as a queue.
// if the task status is completed or exception it is also added to the completed tasks queue
func EnqueueTask(task *tmodel.Task) {
	TaskQueue.queue <- task

}

// dequeueTasksInProgress return the task first in and which are in progress
func dequeueTask() (task *tmodel.Task) {
	if len(TaskQueue.queue) <= 0 {
		return nil
	}
	return <-TaskQueue.queue
}

// UpdateTasksStatus will get the task data which needs to be updated in DB from the queue and will update the DB once in an interval.
func UpdateTasksStatus(size int, d time.Duration) {
	table := "task"
	tasks := make(map[string]interface{}, size)
	completedTasks := make(map[string][2]interface{}, size)
	ticker := time.NewTicker(d)

	for {
		select {
		case <-ticker.C:
			if len(tasks) > 0 {
				err := tmodel.UpdateTaskProgress(tasks, common.InMemory)
				if err != nil {
					l.Log.Error(fmt.Sprintf("error : updating tasks failed - %s", err.Error()))
					break
				}
				for k := range tasks {
					delete(tasks, k)
				}
			}
			if len(completedTasks) > 0 {
				err := tmodel.CreateCompletedTasksIndices(completedTasks, common.InMemory)
				if err != nil {
					l.Log.Error(fmt.Sprintf("error : creating indices for completed tasks failed - %s", err.Error()))
					break
				}
				for k := range tasks {
					delete(tasks, k)
				}
			}
		default:
			task := dequeueTask()
			if task != nil {
				saveID := table + ":" + task.ID
				tasks[saveID] = task
				if (task.TaskState == "Completed" || task.TaskState == "Exception") && task.ParentID == "" {
					key := task.UserName + "::" + task.EndTime.String() + "::" + task.ID
					value := [2]interface{}{
						tmodel.CompletedTaskIndex,
						task.EndTime.UnixNano(),
					}
					completedTasks[key] = value
				}
			}
		}
	}
}
