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

// Package systems ...
package systems

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/scommon"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
)

// PluginContact struct to inject the pmb client function into the handlers
type PluginContact struct {
	ContactClient      func(context.Context, string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	DevicePassword     func([]byte) ([]byte, error)
	GetPluginStatus    func(context.Context, smodel.Plugin) bool
	SavePluginTaskInfo func(ctx context.Context, pluginIP, pluginServerName, odimTaskID, pluginTaskMonURL string) error
	UpdateTask         func(context.Context, common.TaskData) error
}

var (
	// JSONUnMarshal  function pointer for the json.Unmarshal
	JSONUnMarshal = json.Unmarshal
)

// ComputerSystemReset performs a reset action on the requeseted computer system with the specified ResetType
func (p *PluginContact) ComputerSystemReset(ctx context.Context,
	req *systemsproto.ComputerSystemResetRequest, taskID, sessionUserName string) {
	var targetURI = "/redfish/v1/Systems/" + req.SystemID + "/Actions/ComputerSystem.Reset"
	var resp response.RPC
	resp.StatusCode = http.StatusAccepted
	taskInfo := &common.TaskUpdateInfo{Context: ctx, TaskID: taskID, TargetURI: targetURI,
		UpdateTask: p.UpdateTask, TaskRequest: string(req.RequestBody)}

	// parsing the ResetComputerSystem
	var resetCompSys ResetComputerSystem
	err := JSONUnMarshal(req.RequestBody, &resetCompSys)
	if err != nil {
		errMsg := "error: unable to parse the computer system reset request" + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusInternalServerError,
			response.InternalError, errMsg, nil, taskInfo)
		return
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := RequestParamsCaseValidatorFunc(req.RequestBody,
		resetCompSys)
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusInternalServerError,
			response.InternalError, errMsg, nil, taskInfo)
		return
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request" +
			" body are not valid, ensure properties are listed in upper camel case "
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(http.StatusBadRequest, response.PropertyUnknown,
			errorMessage, []interface{}{invalidProperties}, taskInfo)
		return
	}

	// spliting the uuid and system id
	requestData := strings.SplitN(req.SystemID, ".", 2)
	if len(requestData) <= 1 {
		errorMessage := "error: SystemUUID not found"
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound,
			errorMessage, []interface{}{"System", req.SystemID}, taskInfo)
	}

	uuid := requestData[0]
	target, gerr := smodel.GetTarget(uuid)
	if gerr != nil {
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound,
			gerr.Error(), []interface{}{"ComputerSystem", "/redfish/v1/Systems/" +
				req.SystemID}, taskInfo)
		return
	}
	decryptedPasswordByte, err := p.DevicePassword(target.Password)
	if err != nil {
		// Frame the RPC response body and response Header below
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		common.GeneralError(http.StatusInternalServerError, response.InternalError,
			errorMessage, nil, taskInfo)
		return
	}

	target.Password = decryptedPasswordByte
	plugin, gerr := smodel.GetPluginData(target.PluginID)
	if gerr != nil {
		errorMessage := "error while trying to get plugin details"
		common.GeneralError(http.StatusInternalServerError,
			response.InternalError, errorMessage, nil, nil)
		return
	}

	var contactRequest scommon.PluginContactRequest
	contactRequest.ContactClient = p.ContactClient
	contactRequest.Plugin = plugin

	if StringsEqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		contactRequest.HTTPMethodType = http.MethodPost
		contactRequest.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		contactRequest.OID = "/ODIM/v1/Sessions"
		_, token, _, getResponse, err := ContactPluginFunc(ctx, contactRequest,
			"error while creating session with the plugin: ")

		if err != nil {
			common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage,
				err.Error(), nil, nil)
			return
		}
		contactRequest.Token = token
	} else {
		contactRequest.BasicAuth = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}

	postRequest := make(map[string]interface{})
	postRequest["ResetType"] = resetCompSys.ResetType
	postBody, _ := json.Marshal(postRequest)
	target.PostBody = postBody
	contactRequest.HTTPMethodType = http.MethodPost
	contactRequest.DeviceInfo = target
	contactRequest.OID = "/ODIM/v1/Systems/" + requestData[1] +
		"/Actions/ComputerSystem.Reset"
	body, location, pluginIP, getResponse, err := ContactPluginFunc(ctx,
		contactRequest, "error while reseting the computer system: ")
	if err != nil {
		resp.StatusCode = getResponse.StatusCode
		json.Unmarshal(body, &resp.Body)
		task := fillTaskData(taskID, targetURI, string(req.RequestBody), resp,
			common.Exception, common.Critical, 100, http.MethodPost)
		err = p.UpdateTask(ctx, task)
		if err != nil {
			errMsg := "error while starting the task: " + err.Error()
			l.LogWithFields(ctx).Error(errMsg)
			common.GeneralError(http.StatusInternalServerError, response.InternalError,
				errMsg, nil, taskInfo)
		}
		return
	}

	if getResponse.StatusCode == http.StatusAccepted {
		targetIP, err := common.GetIPFromHostName(target.ManagerAddress)
		if err != nil {
			errMsg := "error while saving task data to DB: " + err.Error()
			l.LogWithFields(ctx).Error(errMsg)
			common.GeneralError(http.StatusInternalServerError, response.InternalError,
				errMsg, nil, taskInfo)
			return
		}
		pluginTaskKey := common.ResetTaskIDPrefix + targetIP
		scommon.SavePluginTaskInfoForResetRequest(ctx, pluginIP, plugin.IP,
			taskID, location, pluginTaskKey)
		return
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	err = JSONUnmarshalFunc(body, &resp.Body)
	if err != nil {
		common.GeneralError(http.StatusInternalServerError,
			response.InternalError, err.Error(), nil, taskInfo)
		return
	}
	smodel.AddSystemResetInfo(ctx, "/redfish/v1/Systems/"+req.SystemID,
		resetCompSys.ResetType)
	task := fillTaskData(taskID, targetURI, string(req.RequestBody), resp,
		common.Completed, common.OK, 100, http.MethodPost)
	p.UpdateTask(ctx, task)
}
