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

// Package evresponse have the struct models and DB functionalties
package evresponse

import (
	"sync"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"gotest.tools/assert"
)

func TestResponse(t *testing.T) {
	common.SetUpMockConfig()
	var wg sync.WaitGroup
	var originResource = []string{
		"4228c0db-253b-4dc8-93d1-dab9359139ba.1",
		"423e8254-e3ef-42bd-a130-f096c93a6c42.1",
		"37646c88-a7d7-468c-af58-49e8a0adbbb2.1",
	}
	var hosts = []string{"10.10.10.10", "10.10.10.11", "10.10.10.12"}
	var responses = []EventResponse{{StatusCode: 201}, {StatusCode: 400}, {StatusCode: 201}}
	var result = &MutexLock{
		Response: make(map[string]EventResponse),
		Lock:     &sync.Mutex{},
	}
	for i, origin := range originResource {
		wg.Add(1)
		go func(originResource string, result *MutexLock, wg *sync.WaitGroup, i int) {
			defer wg.Done()
			result.AddResponse(originResource, hosts[i], responses[i])
		}(origin, result, &wg, i)

	}
	wg.Wait()
	_, hostAddresses := result.ReadResponse("1")
	assert.Equal(t, len(hostAddresses), 2, "should be 3 document")
}
