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

// SetDefaultBootOrder defines the logic for setting the boot order to the default
func (p *PluginContact) SetDefaultBootOrder(ctx context.Context, systemID string) response.RPC {
	var resp response.RPC
	l.LogWithFields(ctx).Debugf("incoming SetDefaultBootOrder request for SystemID: %s", systemID)

	// spliting the uuid and system id
	requestData := strings.SplitN(systemID, ".", 2)
	if len(requestData) <= 1 {
		errorMessage := "error: SystemUUID not found"
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"System", systemID}, nil)
	}
	uuid := requestData[0]
	target, gerr := smodel.GetTarget(uuid)
	if gerr != nil {
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, gerr.Error(), []interface{}{"System", uuid}, nil)
	}

	decryptedPasswordByte, err := p.DevicePassword(target.Password)
	if err != nil {
		// Frame the RPC response body and response Header below
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	target.Password = decryptedPasswordByte

	// Get the Plugin info
	plugin, gerr := smodel.GetPluginData(target.PluginID)
	if gerr != nil {
		errorMessage := "error while trying to get plugin details"
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
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
		_, token, _, getResponse, err := ContactPluginFunc(ctx, contactRequest, "error while creating session with the plugin: ")

		if err != nil {
			return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, err.Error(), nil, nil)
		}
		contactRequest.Token = token
	} else {
		contactRequest.BasicAuth = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	postRequest := make(map[string]interface{})
	postBody, _ := json.Marshal(postRequest)
	target.PostBody = postBody
	contactRequest.DeviceInfo = target
	contactRequest.OID = "/ODIM/v1/Systems/" + requestData[1] + "/Actions/ComputerSystem.SetDefaultBootOrder"
	contactRequest.HTTPMethodType = http.MethodPost

	body, _, _, getResponse, err := ContactPluginFunc(ctx, contactRequest, "error while setting the default bootorder of  the computer system: ")
	if err != nil {
		resp.StatusCode = getResponse.StatusCode
		json.Unmarshal(body, &resp.Body)
		return resp
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	err = JSONUnmarshalFunc(body, &resp.Body)
	if err != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
	}
	l.LogWithFields(ctx).Debugf("outgoing response for SetDefaultBootOrder statuscode: %d", resp.StatusCode)
	return resp
}

// ChangeBiosSettings defines the logic for change bios settings
func (p *PluginContact) ChangeBiosSettings(ctx context.Context, req *systemsproto.BiosSettingsRequest, taskID string) {
	var resp response.RPC
	l.LogWithFields(ctx).Debugf("incoming ChangeBiosSettings request for SystemID: %s", req.SystemID)
	var targetURI = "/redfish/v1/Systems/" + req.SystemID + "/Bios/Settings"
	taskInfo := &common.TaskUpdateInfo{Context: ctx, TaskID: taskID, TargetURI: targetURI,
		UpdateTask: p.UpdateTask, TaskRequest: string(req.RequestBody)}
	// spliting the uuid and system id
	requestData := strings.SplitN(req.SystemID, ".", 2)
	if len(requestData) <= 1 {
		errorMessage := "error: SystemUUID not found"
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"System", req.SystemID}, taskInfo)
		return
	}
	uuid := requestData[0]
	target, gerr := smodel.GetTarget(uuid)
	if gerr != nil {
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, gerr.Error(), []interface{}{"System", uuid}, taskInfo)
		return
	}

	var biosSetting BiosSetting

	// parsing the biosSetting
	err := JSONUnmarshalFunc(req.RequestBody, &biosSetting)
	if err != nil {
		errMsg := "unable to parse the BiosSetting request" + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		return
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := RequestParamsCaseValidatorFunc(req.RequestBody, biosSetting)
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		return
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, taskInfo)
		return
	}

	decryptedPasswordByte, err := p.DevicePassword(target.Password)
	if err != nil {
		// Frame the RPC response body and response Header below
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, taskInfo)
		return
	}
	target.Password = decryptedPasswordByte
	// Get the Plugin info
	plugin, gerr := smodel.GetPluginData(target.PluginID)
	if gerr != nil {
		errorMessage := "error while trying to get plugin details"
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, taskInfo)
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
		_, token, _, getResponse, err := ContactPluginFunc(ctx, contactRequest, "error while creating session with the plugin: ")

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
	contactRequest.OID = fmt.Sprintf("/ODIM/v1/Systems/%s/Bios/Settings", requestData[1])

	body, location, pluginIP, getResponse, err := ContactPluginFunc(ctx, contactRequest, "error while changing  bios settings: ")
	if err != nil {
		resp.StatusCode = getResponse.StatusCode
		json.Unmarshal(body, &resp.Body)
		common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, taskInfo)
		return
	}
	if getResponse.StatusCode == http.StatusAccepted {
		scommon.SavePluginTaskInfo(ctx, pluginIP, plugin.IP, taskID, location)
		return
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	err = JSONUnmarshalFunc(body, &resp.Body)
	if err != nil {
		common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, taskInfo)
		return
	}

	// Adding Settings URL to the DB to fetch data from device
	URL := fmt.Sprintf("/redfish/v1/Systems/%s/Bios/Settings", req.SystemID)
	smodel.AddSystemResetInfo(ctx, URL, "None")
	l.LogWithFields(ctx).Debugf("outgoing response for ChangeBiosSettings statuscode: %d", resp.StatusCode)
	task := fillTaskData(taskID, targetURI, string(req.RequestBody), resp,
		common.Completed, common.OK, 100, http.MethodPatch)
	p.UpdateTask(ctx, task)
}

// ChangeBootOrderSettings defines the logic for change boot order settings
func (p *PluginContact) ChangeBootOrderSettings(ctx context.Context, req *systemsproto.BootOrderSettingsRequest,
	taskID string) {
	var targetURI = "/redfish/v1/Systems/" + req.SystemID
	var resp response.RPC

	taskInfo := &common.TaskUpdateInfo{Context: ctx, TaskID: taskID, TargetURI: targetURI,
		UpdateTask: p.UpdateTask, TaskRequest: string(req.RequestBody)}

	// spliting the uuid and system id
	requestData := strings.SplitN(req.SystemID, ".", 2)
	if len(requestData) <= 1 {
		errorMessage := "error: SystemUUID not found"
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound,
			errorMessage, []interface{}{"System", req.SystemID}, taskInfo)
		return
	}

	var bootOrderSettings BootOrderSettings

	// parsing the bootOrderSettings
	err := JSONUnmarshalFunc(req.RequestBody, &bootOrderSettings)
	if err != nil {
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			errMsg := fmt.Sprintf("UnmarshalTypeError: Expected field type %v but got %v \n",
				ute.Type, ute.Value)
			l.LogWithFields(ctx).Error(errMsg)
			index := strings.LastIndex(string(req.RequestBody[:ute.Offset]), ".")
			if index < 0 {
				index = 0
			}
			common.GeneralError(http.StatusBadRequest, response.PropertyValueTypeError,
				errMsg, []interface{}{string(req.RequestBody[index+1 : ute.Offset]),
					ute.Field}, taskInfo)
			return
		}
		errMsg := "unable to parse the BootOrderSettings request" + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusBadRequest, response.MalformedJSON,
			errMsg, []interface{}{}, taskInfo)
		return
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := RequestParamsCaseValidatorFunc(req.RequestBody,
		bootOrderSettings)
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError,
			errMsg, []interface{}{}, taskInfo)
		return
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body" +
			" are not valid, ensure properties are listed in upper camel case "
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(http.StatusBadRequest, response.PropertyUnknown,
			errorMessage, []interface{}{invalidProperties}, taskInfo)
		return
	}

	uuid := requestData[0]
	target, gerr := smodel.GetTarget(uuid)
	if gerr != nil {
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound,
			gerr.Error(), []interface{}{"System", uuid}, taskInfo)
		return
	}
	decryptedPasswordByte, err := p.DevicePassword(target.Password)
	if err != nil {
		// Frame the RPC response body and response Header below
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		common.GeneralError(http.StatusInternalServerError,
			response.InternalError, errorMessage, []interface{}{}, taskInfo)
		return
	}
	target.Password = decryptedPasswordByte
	// Get the Plugin info
	plugin, gerr := smodel.GetPluginData(target.PluginID)
	if gerr != nil {
		errorMessage := "error while trying to get plugin details"
		common.GeneralError(http.StatusInternalServerError, response.InternalError,
			errorMessage, []interface{}{}, taskInfo)
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
				err.Error(), []interface{}{}, taskInfo)
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
	contactRequest.OID = fmt.Sprintf("/ODIM/v1/Systems/%s", requestData[1])

	body, location, pluginIP, getResponse, err := ContactPluginFunc(ctx, contactRequest,
		"error while changing boot order settings: ")
	if err != nil {
		resp.StatusCode = getResponse.StatusCode
		json.Unmarshal(body, &resp.Body)
		return
	}

	if getResponse.StatusCode == http.StatusAccepted {
		scommon.SavePluginTaskInfo(ctx, pluginIP, plugin.IP, taskID, location)
		return
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	err = json.Unmarshal(body, &resp.Body)
	if err != nil {
		common.GeneralError(http.StatusInternalServerError, response.InternalError,
			err.Error(), []interface{}{}, taskInfo)
		return
	}
	smodel.AddSystemResetInfo(ctx, "/redfish/v1/Systems/"+req.SystemID, "On")
	task := fillTaskData(taskID, targetURI, string(req.RequestBody), resp,
		common.Completed, common.OK, 100, http.MethodPatch)
	p.UpdateTask(ctx, task)
	l.LogWithFields(ctx).Debugf("outgoing response for ChangeBootOrderSettings status code: %d", resp.StatusCode)
	return
}
