package tqueue

import (
	"github.com/ODIM-Project/ODIM/svc-task/tmodel"
)

// taskQueue is a struct that contains channels for updating tasks statuses and creating index for completed tasks
type taskQueue struct {
	queue chan *tmodel.Task
}

// TaskQueue is an instance of taskQueue struct which act as the queue for update task requests
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

// EnqueueTask enqueue the update task requests in the channel which act as a queue.
func EnqueueTask(task *tmodel.Task) {
	TaskQueue.queue <- task

}

// UpdateTasksWorker is a goroutine which always listens to the queue for the task update requests.
// UpdateTasksWorker starts the process of updating DB using pipelined transaction and
// it enqueue a task to the queue once in a millisecond which acts as a signal to the process that transaction should be committed.
/* UpdateTasksWorker takes the following keys as input:
1."tick" is of type Tick struct in tmodel package
*/
func UpdateTasksWorker(tick *tmodel.Tick) {
	go Ticker(tick)

	conn := tmodel.GetWriteConnection()
	for {
		if !tick.Executing {
			tick.ProcessTaskQueue(&TaskQueue.queue, conn)
		}
	}
}

// Ticker is executed as a goroutine. It makes the Commit flag true when the ticker ticks.
func Ticker(tick *tmodel.Tick) {
	for range tick.Ticker.C {
		tick.M.Lock()
		tick.Commit = true
		tick.M.Unlock()
	}
}
