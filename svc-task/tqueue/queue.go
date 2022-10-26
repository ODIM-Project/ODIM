package tqueue

import (
	"sync"
	"time"

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

// UpdateTasksWorker is a goroutine which always listens to the queue for the task update requests.
// UpdateTasksWorker starts the process of updating DB using pipelined transaction and
// it enqueue a task to the queue once in a millisecond which acts as a signal to the process that transaction should be committed.
/* UpdateTasksWorker takes the following keys as input:
1."d" is time duration in which the transaction should be committed. Currently it set as 1 millisecond.
*/
func UpdateTasksWorker(d time.Duration) {
	startBatch := true
	var wg sync.WaitGroup
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			signal := new(tmodel.Task)
			signal.Name = tmodel.SignalTaskName
			TaskQueue.queue <- signal
			wg.Wait()
			startBatch = true
		default:
			if startBatch {
				wg.Add(1)
				go tmodel.ProcessTaskQueue(&TaskQueue.queue, &wg)
				startBatch = false
			}
		}
	}
}
