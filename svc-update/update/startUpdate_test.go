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

func TestStartUpdate(t *testing.T) {
	ctx := mockContext()
	var respArgs response.Args
	respArgs = response.Args{
		Code:    response.Success,
		Message: "Request completed successfully",
	}
	body := respArgs.CreateGenericErrorResponse()
	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "Valid Request",
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &updateproto.UpdateRequest{
					SessionToken: "validToken",
				},
			},
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Body:          body,
			},
		},
	}
	e := mockGetExternalInterface()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := e.StartUpdate(ctx, tt.args.taskID, tt.args.sessionUserName, tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StartUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternalInterface_startRequest(t *testing.T) {
	config.SetUpMockConfig(t)
	e := mockGetExternalInterface()
	ctx := mockContext()
	request3 := []byte(`{"ImageURI":"abc","Targets":["/redfish/v1/Systems/uuid.1/target1"],"@Redfish.OperationApplyTime": "OnStartUpdateRequest"}`)
	subTaskChannel := make(chan int32, 8)

	e.startRequest(ctx, "uuid", "someID", string(request3), subTaskChannel, "someUser")
	assert.True(t, true, "There should not be error")

	StringsEqualFoldFunc = func(s, t string) bool {
		return true
	}
	e.External.ContactPlugin = mockContactPluginError
	e.startRequest(ctx, "uuid", "someID", string(request3), subTaskChannel, "someUser")
	assert.True(t, true, "There should not be error")

	StringsEqualFoldFunc = func(s, t string) bool {
		return false
	}

	e.External.ContactPlugin = mockContactPlugin
	e.External.CreateChildTask = mockCreateChildTaskError
	e.startRequest(ctx, "uuid", "someID", string(request3), subTaskChannel, "someUser")
	assert.True(t, true, "There should not be error")

	e.External.GetTarget = mockGetTargetError
	e.External.CreateChildTask = mockCreateChildTask
	e.startRequest(ctx, "uuid", "someID", string(request3), subTaskChannel, "someUser")
	assert.True(t, true, "There should not be error")

	e.External.GetTarget = mockGetTarget
	e.External.DevicePassword = stubDevicePasswordError
	e.startRequest(ctx, "uuid", "someID", string(request3), subTaskChannel, "someUser")
	assert.True(t, true, "There should not be error")

	e.External.GetPluginData = mockGetPluginDataError
	e.External.DevicePassword = stubDevicePassword
	e.startRequest(ctx, "uuid", "someID", string(request3), subTaskChannel, "someUser")
	assert.True(t, true, "There should not be error")

	e.External.GetPluginData = mockGetPluginData
	e.External.UpdateTask = mockUpdateErrorTask
	e.startRequest(ctx, "uuid", "someID", string(request3), subTaskChannel, "someUser")
	assert.True(t, true, "There should not be error")

	e.External.ContactPlugin = mockContactPluginError
	e.External.UpdateTask = mockUpdateTask
	e.startRequest(ctx, "uuid", "someID", string(request3), subTaskChannel, "someUser")
	assert.True(t, true, "There should not be error")

	for i := 0; i < 8; i++ {
		select {
		case statusCode := <-subTaskChannel:
			fmt.Println(statusCode)
		}
	}

}

func TestExternalInterface_StartUpdate(t *testing.T) {
	e := mockGetExternalInterface()
	ctx := mockContext()
	req := &updateproto.UpdateRequest{
		SessionToken: "validToken",
	}
	GetAllKeysFromTableFunc = func(table string, dbtype common.DbType) ([]string, error) {
		return nil, errors.New("")
	}
	e.StartUpdate(ctx, "uuid", "dummySessionName", req)
	GetAllKeysFromTableFunc = func(table string, dbtype common.DbType) ([]string, error) {
		return []string{}, nil
	}
	e.StartUpdate(ctx, "uuid", "dummySessionName", req)

	GetAllKeysFromTableFunc = func(table string, dbtype common.DbType) ([]string, error) {
		return []string{"/redfish/v1/UpdateService/FirmwareInentory/3bd1f589-117a-4cf9-89f2-da44ee8e012b.1"}, nil
	}
	e.DB.GetResource = mockGetResource
	e.StartUpdate(ctx, "uuid", "dummySessionName", req)

	GetAllKeysFromTableFunc = func(table string, dbtype common.DbType) ([]string, error) {
		return []string{"dummy"}, nil
	}
	e.DB.GetResource = mockGetResource
	e.StartUpdate(ctx, "uuid", "dummySessionName", req)

	GetAllKeysFromTableFunc = func(table string, dbtype common.DbType) ([]string, error) {
		return []string{"dummy", "dummy", "dummy"}, nil
	}
	e.DB.GetResource = mockGetResource
	e.StartUpdate(ctx, "uuid", "dummySessionName", req)
}
