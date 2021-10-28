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
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

type args struct {
	taskID, sessionUserName string
	req                     *updateproto.UpdateRequest
}

func mockCreateChildTask(sessionID, taskID string) (string, error) {
	switch taskID {
	case "taskWithoutChild":
		return "", fmt.Errorf("subtask cannot created")
	case "subTaskWithSlash":
		return "someSubTaskID/", nil
	default:
		return "someSubTaskID", nil
	}
}

func mockUpdateTask(task common.TaskData) error {
	if task.TaskID == "invalid" {
		return fmt.Errorf("task with this ID not found")
	}
	return nil
}

func TestSimpleUpdate(t *testing.T) {
	errMsg := []string{"/redfish/v1/Systems/uuid:/target1"}
	errArg1 := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "error: SystemUUID not found",
				MessageArgs:   []interface{}{"System", fmt.Sprintf("%v", errMsg)},
			},
		},
	}
	request1 := []byte(`{"ImageURI":"abc","Targets":["/redfish/v1/Systems/uuid:/target1"],"@Redfish.OperationApplyTime": "OnStartUpdateRequest"}`)
	request3 := []byte(`{"ImageURI":"abc","Targets":["/redfish/v1/Systems/uuid:1/target1"],"@Redfish.OperationApplyTime": "OnStartUpdateRequest"}`)
	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "postive test Case",
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &updateproto.UpdateRequest{
					SessionToken: "validToken",
					RequestBody:  request3,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "invalid system id",
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &updateproto.UpdateRequest{
					SessionToken: "validToken",
					RequestBody:  request1,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotFound,
				Body:       errArg1.CreateGenericErrorResponse(),
			},
		},
		{
			name: "subtask creation failure",
			args: args{
				taskID: "taskWithoutChild", sessionUserName: "someUser",
				req: &updateproto.UpdateRequest{
					SessionToken: "validToken",
					RequestBody:  request3,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusInternalServerError,
			},
		},
	}
	e := mockGetExternalInterface()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := e.SimpleUpdate(tt.args.taskID, tt.args.sessionUserName, tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("SimpleUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
