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

package update

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/stretchr/testify/assert"
)

func TestGetExternalInterface(t *testing.T) {
	res := GetExternalInterface()
	assert.NotNil(t, res, "There should be an error")
	ServicesUpdateTaskFunc = func(ctx context.Context, taskID, taskState, taskStatus string, percentComplete int32, payLoad *task.Payload, endTime time.Time) error {
		return errors.New("Cancelling")
	}
	resp := TaskData(mockContext(), common.TaskData{PercentComplete: 0})
	assert.NotNil(t, resp, "There should be an error")
	ServicesUpdateTaskFunc = func(ctx context.Context, taskID, taskState, taskStatus string, percentComplete int32, payLoad *task.Payload, endTime time.Time) error {
		return nil
	}
	resp = TaskData(mockContext(), common.TaskData{PercentComplete: 0})
	assert.Nil(t, resp, "There should be no error")

}
