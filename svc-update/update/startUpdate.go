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
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-update/ucommon"
	"github.com/ODIM-Project/ODIM/svc-update/umodel"
	log "github.com/sirupsen/logrus"
)

// StartUpdate function handler for on start update process
func (e *ExternalInterface) StartUpdate(taskID string, sessionUserName string, req *updateproto.UpdateRequest) response.RPC {
	var resp response.RPC
	var percentComplete int32
	targetURI := "/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate"

	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: targetURI, UpdateTask: e.External.UpdateTask, TaskRequest: string(req.RequestBody)}
	// Read all the requests from database
	targetList, err := umodel.GetAllKeysFromTable("SimpleUpdate", common.OnDisk)
	if err != nil {
		errMsg := "Unable to read SimpleUpdate requests from database: " + err.Error()
		log.Warn(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	}
	partialResultFlag := false
	subTaskChannel := make(chan int32, len(targetList))
	taskStatus := common.OK
	if len(targetList) == 0 {
		resp.StatusCode = http.StatusOK
		resp.StatusMessage = response.Success
		var args response.Args
		args = response.Args{
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
	for _, target := range targetList {
		data, gerr := e.DB.GetResource("SimpleUpdate", target, common.OnDisk)
		if gerr != nil {
			errMsg := "Unable to retrive the start update request" + gerr.Error()
			log.Warn(errMsg)
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		}
		go e.startRequest(target, taskID, data, subTaskChannel, sessionUserName)
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

	if partialResultFlag {
		taskStatus = common.Warning
	}
	percentComplete = 100
	if resp.StatusCode != http.StatusOK {
		errMsg := "One or more of the SimpleUpdate requests failed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID
		log.Warn(errMsg)
		switch resp.StatusCode {
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

	log.Info("All SimpleUpdate requests successfully completed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID)
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

func (e *ExternalInterface) startRequest(uuid, taskID, data string, subTaskChannel chan<- int32, sessionUserName string) {
	var resp response.RPC
	subTaskURI, err := e.External.CreateChildTask(sessionUserName, taskID)
	if err != nil {
		subTaskChannel <- http.StatusInternalServerError
		log.Warn("Unable to create sub task")
		return
	}
	var subTaskID string
	strArray := strings.Split(subTaskURI, "/")
	if strings.HasSuffix(subTaskURI, "/") {
		subTaskID = strArray[len(strArray)-2]
	} else {
		subTaskID = strArray[len(strArray)-1]
	}

	taskInfo := &common.TaskUpdateInfo{TaskID: subTaskID, TargetURI: uuid, UpdateTask: e.External.UpdateTask, TaskRequest: data}

	var percentComplete int32
	updateRequestBody := strings.Replace(data, uuid+":", "", -1)
	//replacing the request url with south bound translation URL
	for key, value := range config.Data.URLTranslation.SouthBoundURL {
		updateRequestBody = strings.Replace(updateRequestBody, key, value, -1)
	}
	target, gerr := e.External.GetTarget(uuid)
	if gerr != nil {
		subTaskChannel <- http.StatusBadRequest
		errMsg := gerr.Error()
		log.Warn(errMsg)
		common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, gerr.Error(), []interface{}{"System", uuid}, taskInfo)
		return
	}

	decryptedPasswordByte, passwdErr := e.External.DevicePassword(target.Password)
	if passwdErr != nil {
		subTaskChannel <- http.StatusInternalServerError
		errMsg := "Unable to decrypt device password: " + passwdErr.Error()
		log.Warn(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		return
	}
	target.Password = decryptedPasswordByte

	// Get the Plugin info
	plugin, gerr := e.External.GetPluginData(target.PluginID)
	if gerr != nil {
		subTaskChannel <- http.StatusNotFound
		errMsg := "Unable to get plugin data: " + gerr.Error()
		log.Warn(errMsg)
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"PluginData", target.PluginID}, taskInfo)
		return
	}
	var contactRequest ucommon.PluginContactRequest
	contactRequest.ContactClient = e.External.ContactClient
	contactRequest.Plugin = plugin

	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
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
			log.Info(errMsg)
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
	contactRequest.OID = "/ODIM/v1/UpdateService/Actions/UpdateService.StartUpdate"
	contactRequest.HTTPMethodType = http.MethodPost
	_, _, getResponse, contactErr := e.External.ContactPlugin(contactRequest, "error while performing simple update action: ")
	if contactErr != nil {
		subTaskChannel <- getResponse.StatusCode
		errMsg := err.Error()
		log.Info(errMsg)
		common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo)
		return
	}

	resp.StatusCode = http.StatusOK
	percentComplete = 100
	subTaskChannel <- int32(getResponse.StatusCode)
	var task = fillTaskData(subTaskID, uuid, data, resp, common.Completed, common.OK, percentComplete, http.MethodPost)
	err = e.External.UpdateTask(task)
	if err != nil && err.Error() == common.Cancelling {
		var task = fillTaskData(subTaskID, uuid, data, resp, common.Cancelled, common.Critical, percentComplete, http.MethodPost)
		e.External.UpdateTask(task)
	}
	return
}
