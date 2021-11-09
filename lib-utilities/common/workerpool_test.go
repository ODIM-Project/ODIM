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
package common

import (
	"sync"
	"testing"
	"time"
)

var testCounter counter

func process(_ interface{}) bool {
	testCounter.lock.Lock()
	testCounter.count++
	testCounter.lock.Unlock()
	return true
}

func getTestCounter() int {
	testCounter.lock.Lock()
	defer testCounter.lock.Unlock()
	return testCounter.count
}

func TestRunReadWorkers(t *testing.T) {
	in, out := CreateJobQueue(1)
	defer close(in)

	RunReadWorkers(out, process, 5)

	var writeWG sync.WaitGroup

	for j := 0; j < 5; j++ {
		writeWG.Add(1)
		go func() {
			for i := 1; i <= 100; i++ {
				in <- i
			}
			writeWG.Done()
		}()
	}
	writeWG.Wait()

	time.Sleep(2 * time.Second)

	if getTestCounter() != 500 {
		t.Errorf("error: expected count is 500 but got %v", testCounter.count)
	}
}

type data struct {
	currentdata int
	lock        sync.Mutex
}

func (d *data) incrementData() {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.currentdata++
}

func (d *data) getData() int {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.currentdata
}

func TestRunWriteWorkers(t *testing.T) {
	in, out := CreateJobQueue(1)
	var dt data
	var currentData int
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for range out {
			dt.incrementData()
		}
		wg.Done()
	}()

	var dataBatch []interface{}
	for i := 1; i <= 5; i++ {
		dataBatch = append(dataBatch, i)
	}
	done := make(chan bool)
	RunWriteWorkers(in, dataBatch, 2, done)
	ok := <-done
	if ok {
		close(in)
	}
	wg.Wait()

	if dt.getData() != 5 {
		t.Errorf("error: expected count is 5 but got %v", currentData)
	}
}
