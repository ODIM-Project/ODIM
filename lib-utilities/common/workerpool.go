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

// Package common ...
package common

import (
	"sync"
)

// RunReadWorkers will create a worker pool for doing a specific task
// which is passed to it as process after reading the data from the channel.
//
// The function will three parameters:
// 1) jobChannel - an out channel of type interface on which the data is passed
// 2) process - a function which will do some logical tasks by taking data from jobChan
// 3) workerCount - number of workers required for doing the process
func RunReadWorkers(jobChannel <-chan interface{}, jobProcess func(interface{}) bool, workerCount int) {
	for w := 0; w < workerCount; w++ {
		go func() {
			for j := range jobChannel {
				jobProcess(j)
			}
		}()
	}
}

type dataBatchStore struct {
	dataBatch []interface{}
	lock      sync.Mutex
}

// RunWriteWorkers will create a worker pool for inserting data into a channel
//
// The function will three parameters:
// 1) jobChannel - an in channel of type interface to which the data is passed
// 2) dataBatch - a slice of data which needed to be passed into the jobChannel
// 3) workerCount - number of workers required for doing the task
// 4) done - is a notification channel that indicates that job is done
// Please make sure all the data from the dataBatch has been wrote to
// the jobChannel before closing it, else it may cause data loss.
func RunWriteWorkers(jobChannel chan<- interface{}, dataBatch []interface{}, workerCount int, done chan bool) {
	var store dataBatchStore
	store.dataBatch = dataBatch
	for w := 0; w < workerCount; w++ {
		go func() {
			for {
				store.lock.Lock()
				if len(store.dataBatch) == 0 {
					done <- true
					return
				}
				data := store.getNextData()
				jobChannel <- data
				store.lock.Unlock()
			}
		}()
	}
}
func (ds *dataBatchStore) getNextData() interface{} {
	data := ds.dataBatch[0]
	ds.dataBatch = ds.dataBatch[1:]
	return data
}
