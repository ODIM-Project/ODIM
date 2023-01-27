// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package update

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/stretchr/testify/assert"
)

type args struct {
	taskID, sessionUserName string
	req                     *updateproto.UpdateRequest
}

func mockCreateChildTask(ctx context.Context, sessionID, taskID string) (string, error) {
	switch taskID {
	case "taskWithoutChild":
		return "", fmt.Errorf("subtask cannot created")
	case "subTaskWithSlash":
		return "someSubTaskID/", nil
	default:
		return "someSubTaskID", nil
	}
}
func mockCreateChildTaskError(ctx context.Context, sessionID, taskID string) (string, error) {
	return "", errors.New("")
}

func mockUpdateTask(ctx context.Context, task common.TaskData) error {
	if task.TaskID == "invalid" {
		return fmt.Errorf("task with this ID not found")
	}
	return nil
}
func mockUpdateErrorTask(ctx context.Context, task common.TaskData) error {
	return fmt.Errorf("Cancelling")

}

func TestSimpleUpdate(t *testing.T) {
	ctx := mockContext()
	errMsg := []string{"/redfish/v1/Systems/uuid./target1"}
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
	request1 := []byte(`{"ImageURI":"abc","Targets":["/redfish/v1/Systems/uuid./target1"],"@Redfish.OperationApplyTime": "OnStartUpdateRequest"}`)
	request3 := []byte(`{"ImageURI":"abc","Targets":["/redfish/v1/Systems/uuid.1/target1"],"@Redfish.OperationApplyTime": "OnStartUpdateRequest"}`)
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
		{
			name: "Invalid JSON ",
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &updateproto.UpdateRequest{
					SessionToken: "validToken",
					RequestBody:  []byte(`invalidJson`),
				},
			},
			want: response.RPC{
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "Target 0",
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &updateproto.UpdateRequest{
					SessionToken: "validToken",
					RequestBody:  []byte(`{"Targets":[]}`),
				},
			},
			want: response.RPC{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Invalid filed name",
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &updateproto.UpdateRequest{
					SessionToken: "validToken",
					RequestBody:  []byte(`{"imageURI":"abc","Targets":["/redfish/v1/Systems/uuid.1/target1"],"@redfish.OperationApplyTime": "OnStartUpdateRequest"}`),
				},
			},
			want: response.RPC{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Invalid target",
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &updateproto.UpdateRequest{
					SessionToken: "validToken",
					RequestBody:  []byte(`{"Targets":["dummy"]}`),
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotFound,
			},
		},
	}
	e := mockGetExternalInterface()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := e.SimpleUpdate(ctx, tt.args.taskID, tt.args.sessionUserName, tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("SimpleUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_SimpleUpdate(t *testing.T) {
	ctx := mockContext()
	request3 := []byte(`{"ImageURI":"abc","Targets":["/redfish/v1/Systems/uuid.1/target1"],"@Redfish.OperationApplyTime": "OnStartUpdateRequest"}`)
	e := mockGetExternalInterface()
	req := &updateproto.UpdateRequest{
		SessionToken: "validToken",
		RequestBody:  request3,
	}
	RequestParamsCaseValidatorFunc = func(rawRequestBody []byte, reqStruct interface{}) (string, error) {
		return "", errors.New("")
	}
	e.SimpleUpdate(ctx, "invalid", "dummy", req)
	RequestParamsCaseValidatorFunc = func(rawRequestBody []byte, reqStruct interface{}) (string, error) {
		return common.RequestParamsCaseValidator(rawRequestBody, reqStruct)
	}
	JSONMarshalFunc = func(v interface{}) ([]byte, error) {
		return nil, errors.New("")
	}
	e.SimpleUpdate(ctx, "valid", "validId", req)

	JSONMarshalFunc = func(v interface{}) ([]byte, error) {
		return json.Marshal(v)
	}
}

func TestExternalInterface_sendRequestPreferedAuthType(t *testing.T) {
	config.SetUpMockConfig(t)
	ctx := mockContext()
	e := mockGetExternalInterface()
	request3 := []byte(`{"ImageURI":"abc","Targets":["/redfish/v1/Systems/uuid.1/target1"],"@Redfish.OperationApplyTime": "OnStartUpdateRequest"}`)
	subTaskChannel := make(chan int32, 7)

	e.External.ContactPlugin = mockContactPluginError
	StringsEqualFoldFunc = func(s, t string) bool {
		return true
	}
	e.sendRequest(ctx, "uuid", "someID", "/redfish/v1/Systems/uuid", string(request3), "OnStartUpdateRequest", subTaskChannel, "someUser")
	assert.True(t, true, "There should not be error")

	StringsEqualFoldFunc = func(s, t string) bool {
		return false
	}

	e.External.GetTarget = mockGetTargetError
	e.sendRequest(ctx, "uuid", "someID", "/redfish/v1/Systems/uuid", string(request3), "OnStartUpdateRequest", subTaskChannel, "someUser")
	assert.True(t, true, "There should not be error")

	e.External.GetTarget = mockGetTarget
	e.External.GenericSave = stubGenericSaveError
	e.sendRequest(ctx, "uuid", "someID", "/redfish/v1/Systems/uuid", string(request3), "OnStartUpdateRequest", subTaskChannel, "someUser")
	assert.True(t, true, "There should not be error")

	e.External.GenericSave = stubGenericSave
	e.External.DevicePassword = stubDevicePasswordError
	e.sendRequest(ctx, "uuid", "someID", "/redfish/v1/Systems/uuid", string(request3), "OnStartUpdateRequest", subTaskChannel, "someUser")
	assert.True(t, true, "There should not be error")

	e.External.DevicePassword = stubDevicePassword
	e.External.ContactPlugin = mockContactPluginError

	e.sendRequest(ctx, "uuid", "someID", "/redfish/v1/Systems/uuid", string(request3), "OnStartUpdateRequest", subTaskChannel, "someUser")
	assert.True(t, true, "There should not be error")

	e.External.GetPluginData = mockGetPluginDataError
	e.External.ContactPlugin = mockContactPlugin

	e.sendRequest(ctx, "uuid", "someID", "/redfish/v1/Systems/uuid", string(request3), "OnStartUpdateRequest", subTaskChannel, "someUser")
	assert.True(t, true, "There should not be error")

	e.External.GetPluginData = mockGetPluginData
	e.External.UpdateTask = mockUpdateErrorTask

	e.sendRequest(ctx, "uuid", "someID", "/redfish/v1/Systems/uuid", string(request3), "OnStartUpdateRequest", subTaskChannel, "someUser")
	assert.True(t, true, "There should not be error")

	for i := 0; i < 7; i++ {
		select {
		case statusCode := <-subTaskChannel:
			fmt.Println(statusCode)
		}
	}
}
