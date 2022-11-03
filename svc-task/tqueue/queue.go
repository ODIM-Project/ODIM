package tqueue

import (
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
	tick := &tmodel.Tick{
		Ticker: time.NewTicker(d),
	}

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
