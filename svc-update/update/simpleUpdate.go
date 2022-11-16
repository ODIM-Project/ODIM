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

//Package update ...
package update

// ---------------------------------------------------------------------------------------
// IMPORT Section
//
import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-update/ucommon"
)

var (
	//RequestParamsCaseValidatorFunc ...
	RequestParamsCaseValidatorFunc = common.RequestParamsCaseValidator
	//JSONMarshalFunc ...
	JSONMarshalFunc = json.Marshal
	//StringsEqualFoldFunc ...
	StringsEqualFoldFunc = strings.EqualFold
)

// SimpleUpdate function handler for simpe update process
func (e *ExternalInterface) SimpleUpdate(taskID string, sessionUserName string, req *updateproto.UpdateRequest) response.RPC {
	var resp response.RPC
	var percentComplete int32
	targetURI := "/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate"

	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: targetURI, UpdateTask: e.External.UpdateTask, TaskRequest: string(req.RequestBody)}

	var updateRequest SimpleUpdateRequest
	err := json.Unmarshal(req.RequestBody, &updateRequest)
	if err != nil {
		errMsg := "Unable to parse the simple update request" + err.Error()
		l.Log.Warn(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	}
	if len(updateRequest.Targets) == 0 {
		errMsg := "'Targets' parameter cannot be empty"
		l.Log.Warn(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"Targets"}, taskInfo)
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := RequestParamsCaseValidatorFunc(req.RequestBody, updateRequest)
	if err != nil {
		errMsg := "Unable to validate request parameters: " + err.Error()
		l.Log.Warn(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	} else if invalidProperties != "" {
		errorMessage := "One or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		l.Log.Warn(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, taskInfo)
		return response
	}

	targetList := make(map[string][]string)
	targetList, err = sortTargetList(updateRequest.Targets)
	if err != nil {
		errorMessage := "SystemUUID not found"
		l.Log.Warn(errorMessage)
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"System", fmt.Sprintf("%v", updateRequest.Targets)}, taskInfo)
	}
	partialResultFlag := false
	subTaskChannel := make(chan int32, len(targetList))
	serverURI := ""
	for id, target := range targetList {
		updateRequest.Targets = target
		marshalBody, err := JSONMarshalFunc(updateRequest)
		if err != nil {
			errMsg := "Unable to parse the simple update request" + err.Error()
			l.Log.Warn(errMsg)
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		}
		updateRequestBody := string(marshalBody)
		serverURI = "/redfish/v1/Systems/" + id
		go e.sendRequest(id, taskID, serverURI, updateRequestBody, updateRequest.RedfishOperationApplyTime, subTaskChannel, sessionUserName)
	}

	resp.StatusCode = http.StatusOK
	for i := 0; i < len(targetList); i++ {
		select {
		case statusCode := <-subTaskChannel:
			if statusCode != http.StatusOK {
				partialResultFlag = true
				if resp.StatusCode < statusCode {
					resp.StatusCode = statusCode
				}
			}
			if i < len(targetList)-1 {
				percentComplete := int32(((i + 1) / len(targetList)) * 100)
				var task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Running, common.OK, percentComplete, http.MethodPost)
				err := e.External.UpdateTask(task)
				if err != nil && err.Error() == common.Cancelling {
					task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Cancelled, common.OK, percentComplete, http.MethodPost)
					e.External.UpdateTask(task)
					runtime.Goexit()
				}

			}
		}
	}

	taskStatus := common.OK
	if partialResultFlag {
		taskStatus = common.Warning
	}
	percentComplete = 100
	if resp.StatusCode != http.StatusOK {
		errMsg := "One or more of the SimpleUpdate requests failed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID
		l.Log.Warn(errMsg)
		switch resp.StatusCode {
		case http.StatusAccepted:
			return common.GeneralError(http.StatusAccepted, response.TaskStarted, errMsg, []interface{}{fmt.Sprintf("%v", targetList)}, taskInfo)
		case http.StatusUnauthorized:
			return common.GeneralError(http.StatusUnauthorized, response.ResourceAtURIUnauthorized, errMsg, []interface{}{fmt.Sprintf("%v", targetList)}, taskInfo)
		case http.StatusNotFound:
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"option", "SimpleUpdate"}, taskInfo)
		case http.StatusBadRequest:
			return common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errMsg, []interface{}{"UpdateService.SimpleUpdate"}, taskInfo)
		default:
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		}
	}

	l.Log.Info("All SimpleUpdate requests successfully completed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID)
	resp.StatusMessage = response.Success
	resp.StatusCode = http.StatusOK
	args := response.Args{
		Code:    resp.StatusMessage,
		Message: "Request completed successfully",
	}
	resp.Body = args.CreateGenericErrorResponse()

	var task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Completed, taskStatus, percentComplete, http.MethodPost)
	err = e.External.UpdateTask(task)
	if err != nil && err.Error() == common.Cancelling {
		task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Cancelled, common.Critical, percentComplete, http.MethodPost)
		e.External.UpdateTask(task)
		runtime.Goexit()
	}
	return resp
}

func (e *ExternalInterface) sendRequest(uuid, taskID, serverURI, updateRequestBody string, applyTime string, subTaskChannel chan<- int32, sessionUserName string) {
	var resp response.RPC
	subTaskURI, err := e.External.CreateChildTask(sessionUserName, taskID)
	if err != nil {
		subTaskChannel <- http.StatusInternalServerError
		l.Log.Warn("Unable to create sub task")
		return
	}
	var subTaskID string
	strArray := strings.Split(subTaskURI, "/")
	if strings.HasSuffix(subTaskURI, "/") {
		subTaskID = strArray[len(strArray)-2]
	} else {
		subTaskID = strArray[len(strArray)-1]
	}

	taskInfo := &common.TaskUpdateInfo{TaskID: subTaskID, TargetURI: serverURI, UpdateTask: e.External.UpdateTask, TaskRequest: updateRequestBody}

	var percentComplete int32
	target, gerr := e.External.GetTarget(uuid)
	if gerr != nil {
		subTaskChannel <- http.StatusBadRequest
		errMsg := gerr.Error()
		l.Log.Warn(errMsg)
		common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, gerr.Error(), []interface{}{"System", uuid}, nil)
		return
	}
	if applyTime == "OnStartUpdateRequest" {
		err := e.External.GenericSave([]byte(updateRequestBody), "SimpleUpdate", uuid)
		if err != nil {
			subTaskChannel <- http.StatusInternalServerError
			errMsg := "Unable to save the simple update request" + err.Error()
			l.Log.Warn(errMsg)
			common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
			return
		}
	}
	updateRequestBody = strings.Replace(string(updateRequestBody), uuid+".", "", -1)
	//replacing the reruest url with south bound translation URL
	for key, value := range config.Data.URLTranslation.SouthBoundURL {
		updateRequestBody = strings.Replace(updateRequestBody, key, value, -1)
	}

	decryptedPasswordByte, err := e.External.DevicePassword(target.Password)
	if err != nil {
		subTaskChannel <- http.StatusInternalServerError
		errMsg := "Unable to decrypt device password: " + err.Error()
		l.Log.Warn(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		return
	}
	target.Password = decryptedPasswordByte

	// Get the Plugin info
	plugin, gerr := e.External.GetPluginData(target.PluginID)
	if gerr != nil {
		subTaskChannel <- http.StatusNotFound
		errMsg := "Unable to get plugin data: " + gerr.Error()
		l.Log.Warn(errMsg)
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"PluginData", target.PluginID}, taskInfo)
		return
	}
	var contactRequest ucommon.PluginContactRequest
	contactRequest.ContactClient = e.External.ContactClient
	contactRequest.Plugin = plugin

	if StringsEqualFoldFunc(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		contactRequest.HTTPMethodType = http.MethodPost
		contactRequest.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		contactRequest.OID = "/ODIM/v1/Sessions"
		_, token, getResponse, err := e.External.ContactPlugin(contactRequest, "error while creating session with the plugin: ")
		if err != nil {
			subTaskChannel <- getResponse.StatusCode
			errMsg := err.Error()
			l.Log.Warn(errMsg)
			common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo)
			return
		}

		contactRequest.Token = token
	} else {
		contactRequest.BasicAuth = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}

	target.PostBody = []byte(updateRequestBody)
	contactRequest.DeviceInfo = target
	contactRequest.OID = "/ODIM/v1/UpdateService/Actions/UpdateService.SimpleUpdate"
	contactRequest.HTTPMethodType = http.MethodPost
	respBody, location, getResponse, err := e.External.ContactPlugin(contactRequest, "error while performing simple update action: ")
	if err != nil {
		subTaskChannel <- getResponse.StatusCode
		errMsg := err.Error()
		l.Log.Warn(errMsg)
		common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo)
		return
	}
	if getResponse.StatusCode == http.StatusAccepted {
		getResponse, err = e.monitorPluginTask(subTaskChannel, &monitorTaskRequest{
			subTaskID:         subTaskID,
			serverURI:         serverURI,
			updateRequestBody: updateRequestBody,
			respBody:          respBody,
			getResponse:       getResponse,
			taskInfo:          taskInfo,
			location:          location,
			pluginRequest:     contactRequest,
			resp:              resp,
		})

		if err != nil {
			return
		}
	}
	resp.StatusCode = http.StatusOK
	percentComplete = 100

	subTaskChannel <- int32(getResponse.StatusCode)
	var task = fillTaskData(subTaskID, serverURI, updateRequestBody, resp, common.Completed, common.OK, percentComplete, http.MethodPost)
	err = e.External.UpdateTask(task)
	if err != nil && err.Error() == common.Cancelling {
		var task = fillTaskData(subTaskID, serverURI, updateRequestBody, resp, common.Cancelled, common.Critical, percentComplete, http.MethodPost)
		e.External.UpdateTask(task)
	}
	return
}

func sortTargetList(Targets []string) (map[string][]string, error) {
	returnList := make(map[string][]string)
	for _, individualTarget := range Targets {
		// spliting the uuid and system id
		requestData := strings.Split(individualTarget, "/")
		var requestTarget []string
		for _, data := range requestData {
			if strings.Contains(data, ".") {
				requestTarget = strings.SplitN(data, ".", 2)
			}
		}
		if len(requestTarget) != 2 || requestTarget[1] == "" {
			errorMessage := "error: SystemUUID not found"
			return returnList, errors.New(errorMessage)
		}
		uuid := requestTarget[0]
		returnList[uuid] = append(returnList[uuid], individualTarget)
	}
	return returnList, nil
}
