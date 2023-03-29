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

package managers

import (
	"context"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrcommon"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrmodel"
)

// ExternalInterface holds all the external connections managers package functions uses
type ExternalInterface struct {
	Device Device
	DB     DB
	RPC    RPC
}

// Device struct to inject the contact device function into the handlers
type Device struct {
	GetDeviceInfo         func(context.Context, mgrcommon.ResourceInfoRequest) (string, error)
	DeviceRequest         func(context.Context, mgrcommon.ResourceInfoRequest) (mgrcommon.PluginTaskInfo, response.RPC)
	ContactClient         func(context.Context, string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	DecryptDevicePassword func([]byte) ([]byte, error)
}

// DB struct to inject the contact DB function into the handlers
type DB struct {
	GetAllKeysFromTable func(string) ([]string, error)
	GetManagerByURL     func(string) (string, *errors.Error)
	GetPluginData       func(string) (mgrmodel.Plugin, *errors.Error)
	UpdateData          func(string, map[string]interface{}, string) error
	SavePluginTaskInfo  func(context.Context, string, string, string, string) error
	GetResource         func(string, string) (string, *errors.Error)
}

// RPC struct to inject the rpc call to other services
type RPC struct {
	UpdateTask func(context.Context, common.TaskData) error
}

// GetExternalInterface retrieves all the external connections managers package functions uses
func GetExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		Device: Device{
			GetDeviceInfo:         mgrcommon.GetResourceInfoFromDevice,
			DeviceRequest:         mgrcommon.DeviceCommunication,
			ContactClient:         pmbhandle.ContactPlugin,
			DecryptDevicePassword: common.DecryptWithPrivateKey,
		},
		DB: DB{
			GetAllKeysFromTable: mgrmodel.GetAllKeysFromTable,
			GetManagerByURL:     mgrmodel.GetManagerByURL,
			GetPluginData:       mgrmodel.GetPluginData,
			UpdateData:          mgrmodel.UpdateData,
			SavePluginTaskInfo:  services.SavePluginTaskInfo,
			GetResource:         mgrmodel.GetResource,
		},
		RPC: RPC{
			UpdateTask: mgrcommon.UpdateTask,
		},
	}
}

func fillTaskData(taskID, targetURI, request string, resp response.RPC, taskState string, taskStatus string, percentComplete int32, httpMethod string) common.TaskData {
	return common.TaskData{
		TaskID:          taskID,
		TargetURI:       targetURI,
		TaskRequest:     request,
		Response:        resp,
		TaskState:       taskState,
		TaskStatus:      taskStatus,
		PercentComplete: percentComplete,
		HTTPMethod:      httpMethod,
	}
}
