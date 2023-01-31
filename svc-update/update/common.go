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
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-update/ucommon"
	"github.com/ODIM-Project/ODIM/svc-update/umodel"
)

var (
	//ServicesUpdateTaskFunc ...
	ServicesUpdateTaskFunc = services.UpdateTask
)

// Device struct to define the response from plugin for UUID
type Device struct {
	ServerIP   string `json:"ServerIP"`
	Username   string `json:"Username"`
	DeviceUUID string `json:"device_UUID"`
}

// ExternalInterface struct holds the structs to which hold function pointers to outboud calls
type ExternalInterface struct {
	External External
	DB       DB
}

// Plugin is the model for plugin information
type Plugin struct {
	IP                string
	Port              string
	Username          string
	Password          []byte
	ID                string
	PluginType        string
	PreferredAuthType string
}

// External struct holds the function pointers all outboud services
type External struct {
	ContactClient      func(context.Context, string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	Auth               func(string, []string, []string) (response.RPC, error)
	DevicePassword     func([]byte) ([]byte, error)
	GetPluginData      func(string) (umodel.Plugin, *errors.Error)
	ContactPlugin      func(context.Context, ucommon.PluginContactRequest, string) ([]byte, string, ucommon.ResponseStatus, error)
	GetTarget          func(string) (*umodel.Target, *errors.Error)
	CreateChildTask    func(context.Context, string, string) (string, error)
	CreateTask         func(context.Context, string) (string, error)
	UpdateTask         func(context.Context, common.TaskData) error
	GetSessionUserName func(string) (string, error)
	GenericSave        func(context.Context, []byte, string, string) error
}

type responseStatus struct {
	StatusCode    int32
	StatusMessage string
	MsgArgs       []interface{}
}

// DB struct holds the function pointers to database operations
type DB struct {
	GetAllKeysFromTable func(string, common.DbType) ([]string, error)
	GetResource         func(string, string, common.DbType) (string, *errors.Error)
}

// SimpleUpdateRequest struct defines the request body for update action
type SimpleUpdateRequest struct {
	ImageURI                  string   `json:"ImageURI"`
	Password                  string   `json:"Password,omitempty"`
	Targets                   []string `json:"Targets"`
	TransferProtocol          string   `json:"TransferProtocol,omitempty"`
	Username                  string   `json:"Username,omitempty"`
	RedfishOperationApplyTime string   `json:"@Redfish.OperationApplyTime,omitempty"`
}

// monitorTaskRequest hold values required monitorTask function
type monitorTaskRequest struct {
	respBody          []byte
	subTaskID         string
	serverURI         string
	updateRequestBody string
	getResponse       ucommon.ResponseStatus
	location          string
	taskInfo          *common.TaskUpdateInfo
	pluginRequest     ucommon.PluginContactRequest
	resp              response.RPC
}

// GetExternalInterface retrieves all the external connections update package functions uses
func GetExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		External: External{
			ContactClient:      pmbhandle.ContactPlugin,
			Auth:               services.IsAuthorized,
			DevicePassword:     common.DecryptWithPrivateKey,
			GetPluginData:      umodel.GetPluginData,
			ContactPlugin:      ucommon.ContactPlugin,
			GetTarget:          umodel.GetTarget,
			UpdateTask:         TaskData,
			CreateChildTask:    services.CreateChildTask,
			GetSessionUserName: services.GetSessionUserName,
			CreateTask:         services.CreateTask,
			GenericSave:        umodel.GenericSave,
		},
		DB: DB{
			GetAllKeysFromTable: umodel.GetAllKeysFromTable,
			GetResource:         umodel.GetResource,
		},
	}
}

// TaskData update the task with the given data
func TaskData(ctx context.Context, taskData common.TaskData) error {
	respBody, _ := json.Marshal(taskData.Response.Body)
	payLoad := &taskproto.Payload{
		HTTPHeaders:   taskData.Response.Header,
		HTTPOperation: taskData.HTTPMethod,
		JSONBody:      taskData.TaskRequest,
		StatusCode:    taskData.Response.StatusCode,
		TargetURI:     taskData.TargetURI,
		ResponseBody:  respBody,
	}

	err := ServicesUpdateTaskFunc(ctx, taskData.TaskID, taskData.TaskState, taskData.TaskStatus, taskData.PercentComplete, payLoad, time.Now())
	if err != nil && (err.Error() == common.Cancelling) {
		// We cant do anything here as the task has done it work completely, we cant reverse it.
		//Unless if we can do opposite/reverse action for delete server which is add server.
		ServicesUpdateTaskFunc(ctx, taskData.TaskID, common.Cancelled, taskData.TaskStatus, taskData.PercentComplete, payLoad, time.Now())
		if taskData.PercentComplete == 0 {
			return fmt.Errorf("error while starting the task: %v", err)
		}
		runtime.Goexit()
	}
	return nil
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

func (e *ExternalInterface) monitorPluginTask(ctx context.Context, subTaskChannel chan<- int32, monitorTaskData *monitorTaskRequest) (ucommon.ResponseStatus, error) {
	for {

		var task common.TaskData
		if err := json.Unmarshal(monitorTaskData.respBody, &task); err != nil {
			subTaskChannel <- http.StatusInternalServerError
			errMsg := "Unable to parse the simple update respone" + err.Error()
			l.LogWithFields(ctx).Warn(errMsg)
			common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, monitorTaskData.taskInfo)
			return monitorTaskData.getResponse, err
		}
		var updatetask = fillTaskData(monitorTaskData.subTaskID, monitorTaskData.serverURI, monitorTaskData.updateRequestBody, monitorTaskData.resp, task.TaskState, task.TaskStatus, task.PercentComplete, http.MethodPost)
		err := e.External.UpdateTask(ctx, updatetask)
		if err != nil && err.Error() == common.Cancelling {
			var updatetask = fillTaskData(monitorTaskData.subTaskID, monitorTaskData.serverURI, monitorTaskData.updateRequestBody, monitorTaskData.resp, common.Cancelled, common.Critical, 100, http.MethodPost)
			subTaskChannel <- http.StatusInternalServerError
			e.External.UpdateTask(ctx, updatetask)
			return monitorTaskData.getResponse, err
		}
		time.Sleep(time.Second * 5)
		monitorTaskData.pluginRequest.OID = monitorTaskData.location
		monitorTaskData.pluginRequest.HTTPMethodType = http.MethodGet
		monitorTaskData.respBody, _, monitorTaskData.getResponse, err = e.External.ContactPlugin(ctx, monitorTaskData.pluginRequest, "error while performing simple update action: ")
		if err != nil {
			subTaskChannel <- monitorTaskData.getResponse.StatusCode
			errMsg := err.Error()
			l.LogWithFields(ctx).Warn(errMsg)
			common.GeneralError(monitorTaskData.getResponse.StatusCode, monitorTaskData.getResponse.StatusMessage, errMsg, monitorTaskData.getResponse.MsgArgs, monitorTaskData.taskInfo)
			return monitorTaskData.getResponse, err
		}
		if monitorTaskData.getResponse.StatusCode == http.StatusOK {
			break
		}
	}
	return monitorTaskData.getResponse, nil
}
