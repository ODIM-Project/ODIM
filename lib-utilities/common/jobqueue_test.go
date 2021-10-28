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
)

func TestCreateJobQueue(t *testing.T) {
	in, out := CreateJobQueue(1)
	var currentData int
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for v := range out {
			data := v.(int)
			if currentData+1 != data {
				t.Errorf("unexpected value. expected %v, got %v", currentData+1, data)
			}
			currentData = data
		}
		wg.Done()
	}()

	for i := 1; i <= 100; i++ {
		in <- i
	}
	close(in)
	wg.Wait()

	if currentData != 100 {
		t.Errorf("error: only received up to %v", currentData)
	}
}

type counter struct {
	count int
	lock  sync.Mutex
}

func TestCreateJobQueueWithMultipleReadAndWrites(t *testing.T) {
	in, out := CreateJobQueue(1)
	var c counter
	var readWG, writeWG sync.WaitGroup

	c.startReadRoutine(&readWG, out)
	c.startWriteRoutine(&writeWG, in)

	writeWG.Wait()
	close(in)
	readWG.Wait()

	if c.count != 500 {
		t.Errorf("error: expected count is 500 but got %v", c.count)
	}
}

func TestCreateJobQueueWithMultipleInstances(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			in, out := CreateJobQueue(1)
			var (
				c               counter
				readWG, writeWG sync.WaitGroup
			)

			c.startReadRoutine(&readWG, out)
			c.startWriteRoutine(&writeWG, in)

			writeWG.Wait()
			close(in)
			readWG.Wait()

			if c.count != 500 {
				t.Errorf("error: expected count is 500 but got %v", c.count)
			}
		}()
	}
	wg.Wait()
}

func (c *counter) startReadRoutine(readWG *sync.WaitGroup, out <-chan interface{}) {
	for i := 0; i < 5; i++ {
		readWG.Add(1)
		go func() {
			defer readWG.Done()
			for range out {
				c.lock.Lock()
				c.count++
				c.lock.Unlock()
			}
		}()
	}
}

func (c *counter) startWriteRoutine(writeWG *sync.WaitGroup, in chan<- interface{}) {
	for i := 0; i < 5; i++ {
		writeWG.Add(1)
		go func() {
			defer writeWG.Done()
			for i := 1; i <= 100; i++ {
				in <- i
			}
		}()
	}
}
