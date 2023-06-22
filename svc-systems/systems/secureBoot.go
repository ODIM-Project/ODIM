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
	"fmt"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/scommon"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
)

// UpdateSecureBoot defines the logic for updating SecureBoot
func (e *ExternalInterface) UpdateSecureBoot(ctx context.Context, req *systemsproto.SecureBootRequest, pc *PluginContact, taskID string) {
	var resp response.RPC
	var targetURI = "/redfish/v1/Systems/" + req.SystemID + "/SecureBoot"

	taskInfo := &common.TaskUpdateInfo{Context: ctx, TaskID: taskID, TargetURI: targetURI,
		UpdateTask: pc.UpdateTask, TaskRequest: string(req.RequestBody)}

	// spliting the uuid and system id
	requestData := strings.SplitN(req.SystemID, ".", 2)
	if len(requestData) <= 1 {
		errorMessage := "error: SystemUUID not found"
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage,
			[]interface{}{"System", req.SystemID}, taskInfo)
		return
	}

	uuid := requestData[0]
	target, gerr := e.DB.GetTarget(uuid)
	if gerr != nil {
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound,
			gerr.Error(), []interface{}{"System", uuid}, taskInfo)
		return
	}

	var secureBoot smodel.SecureBoot
	// unmarshalling the volume
	err := json.Unmarshal(req.RequestBody, &secureBoot)
	if err != nil {
		errorMessage := "Error while unmarshaling the create volume request: " + err.Error()
		if StringContain(err.Error(), "smodel.OdataIDLink") {
			errorMessage = "Error processing create volume request: @odata.id key(s) is missing in Drives list"
		}
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, []interface{}{}, taskInfo)
		return
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := RequestParamsCaseValidatorFunc(req.RequestBody, secureBoot)
	if err != nil {
		errMsg := "error while validating request parameters for volume creation: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		return
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, taskInfo)
		return
	}

	decryptedPasswordByte, err := e.DevicePassword(target.Password)
	if err != nil {
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, taskInfo)
		return
	}
	target.Password = decryptedPasswordByte
	// Get the Plugin info
	plugin, gerr := e.DB.GetPluginData(target.PluginID)
	if gerr != nil {
		errorMessage := "error while trying to get plugin details"
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, taskInfo)
		return
	}

	var contactRequest scommon.PluginContactRequest
	contactRequest.ContactClient = e.ContactClient
	contactRequest.Plugin = plugin
	contactRequest.GetPluginStatus = e.GetPluginStatus

	if StringsEqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		contactRequest.HTTPMethodType = http.MethodPost
		contactRequest.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		contactRequest.OID = "/ODIM/v1/Sessions"
		_, token, _, getResponse, err := scommon.ContactPlugin(ctx, contactRequest, "error while creating session with the plugin: ")

		if err != nil {
			common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, err.Error(), nil, taskInfo)
			return
		}
		contactRequest.Token = token
	} else {
		contactRequest.BasicAuth = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	target.PostBody = req.RequestBody

	contactRequest.HTTPMethodType = http.MethodPatch
	contactRequest.DeviceInfo = target
	contactRequest.OID = fmt.Sprintf("/ODIM/v1/Systems/%s/SecureBoot", requestData[1])

	body, location, pluginIP, getResponse, err := ContactPluginFunc(ctx, contactRequest, "error while creating a volume: ")
	if err != nil {
		resp.StatusCode = getResponse.StatusCode
		json.Unmarshal(body, &resp.Body)
		errMsg := "error while creating volume: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError,
			errMsg, nil, taskInfo)
		return
	}
	if getResponse.StatusCode == http.StatusAccepted {
		err = pc.SavePluginTaskInfo(ctx, pluginIP, plugin.IP, taskID, location)
		if err != nil {
			l.LogWithFields(ctx).Error(err)
		}
		return
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	err = JSONUnmarshalFunc(body, &resp.Body)
	if err != nil {
		common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, taskInfo)
		return
	}
	task := fillTaskData(taskID, targetURI, string(req.RequestBody), resp,
		common.Completed, common.OK, 100, http.MethodPost)
	pc.UpdateTask(ctx, task)

}

// ResetSecureBoot defines the logic for resetting SecureBoot keys
func (e *ExternalInterface) ResetSecureBoot(ctx context.Context, req *systemsproto.SecureBootRequest, pc *PluginContact, taskID string) {
	var resp response.RPC
	var targetURI = "/redfish/v1/Systems/" + req.SystemID + "/Actions/SecureBoot.ResetKeys"

	taskInfo := &common.TaskUpdateInfo{Context: ctx, TaskID: taskID, TargetURI: targetURI,
		UpdateTask: pc.UpdateTask, TaskRequest: string(req.RequestBody)}

	// spliting the uuid and system id
	requestData := strings.SplitN(req.SystemID, ".", 2)
	if len(requestData) <= 1 {
		errorMessage := "error: SystemUUID not found"
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage,
			[]interface{}{"System", req.SystemID}, taskInfo)
		return
	}

	uuid := requestData[0]
	target, gerr := e.DB.GetTarget(uuid)
	if gerr != nil {
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound,
			gerr.Error(), []interface{}{"System", uuid}, taskInfo)
		return
	}

	var reset smodel.ResetSecureBoot
	// unmarshalling the volume
	err := json.Unmarshal(req.RequestBody, &reset)
	if err != nil {
		errorMessage := "Error while unmarshaling the create volume request: " + err.Error()
		if StringContain(err.Error(), "smodel.OdataIDLink") {
			errorMessage = "Error processing create volume request: @odata.id key(s) is missing in Drives list"
		}
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, []interface{}{}, taskInfo)
		return
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := RequestParamsCaseValidatorFunc(req.RequestBody, reset)
	if err != nil {
		errMsg := "error while validating request parameters for volume creation: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		return
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, taskInfo)
		return
	}

	decryptedPasswordByte, err := e.DevicePassword(target.Password)
	if err != nil {
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, taskInfo)
		return
	}
	target.Password = decryptedPasswordByte
	// Get the Plugin info
	plugin, gerr := e.DB.GetPluginData(target.PluginID)
	if gerr != nil {
		errorMessage := "error while trying to get plugin details"
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, taskInfo)
		return
	}

	var contactRequest scommon.PluginContactRequest
	contactRequest.ContactClient = e.ContactClient
	contactRequest.Plugin = plugin
	contactRequest.GetPluginStatus = e.GetPluginStatus

	if StringsEqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		contactRequest.HTTPMethodType = http.MethodPost
		contactRequest.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		contactRequest.OID = "/ODIM/v1/Sessions"
		_, token, _, getResponse, err := scommon.ContactPlugin(ctx, contactRequest, "error while creating session with the plugin: ")

		if err != nil {
			common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, err.Error(), nil, taskInfo)
			return
		}
		contactRequest.Token = token
	} else {
		contactRequest.BasicAuth = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	target.PostBody = req.RequestBody

	contactRequest.HTTPMethodType = http.MethodPost
	contactRequest.DeviceInfo = target
	contactRequest.OID = fmt.Sprintf("/ODIM/v1/Systems/%s/SecureBoot/Actions/SecureBoot.ResetKeys", requestData[1])

	body, location, pluginIP, getResponse, err := ContactPluginFunc(ctx, contactRequest, "error while resetting secure boot: ")
	if err != nil {
		resp.StatusCode = getResponse.StatusCode
		json.Unmarshal(body, &resp.Body)
		errMsg := "error while resetting secure boot: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError,
			errMsg, nil, taskInfo)
		return
	}
	if getResponse.StatusCode == http.StatusAccepted {
		err = pc.SavePluginTaskInfo(ctx, pluginIP, plugin.IP, taskID, location)
		if err != nil {
			l.LogWithFields(ctx).Error(err)
		}
		return
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	err = JSONUnmarshalFunc(body, &resp.Body)
	if err != nil {
		common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, taskInfo)
		return
	}
	task := fillTaskData(taskID, targetURI, string(req.RequestBody), resp,
		common.Completed, common.OK, 100, http.MethodPost)
	pc.UpdateTask(ctx, task)
}
