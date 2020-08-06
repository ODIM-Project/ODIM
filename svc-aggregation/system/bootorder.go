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

// Package system ...
package system

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

// SetBootOrderRequest defines the request for the setting default boot order
type SetBootOrderRequest struct {
	Parameters BootOrderParameters `json:"Parameters"`
}

// BootOrderParameters defines the boot order parameters
type BootOrderParameters struct {
	ServerCollection []string `json:"ServerCollection"`
}

type responseHolder struct {
	response   []interface{}
	anyFailure bool
	lock       sync.Mutex
}

// SetDefaultBootOrder defines the logic for setting the boot order to the default
func (e *ExternalInterface) SetDefaultBootOrder(taskID string, sessionUserName string, req *aggregatorproto.AggregatorRequest) response.RPC {
	var resp response.RPC
	var percentComplete int32
	targetURI := "/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder"

	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: targetURI, UpdateTask: e.UpdateTask}

	var setOrderReq SetBootOrderRequest
	if err := json.Unmarshal(req.RequestBody, &setOrderReq); err != nil {
		errMsg := "error while trying to set default boot order: " + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, &SetBootOrderRequest{})
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Println(errorMessage)
		resp := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, taskInfo)
		return resp
	}

	partialResultFlag := false
	subTaskChan := make(chan int32, len(setOrderReq.Parameters.ServerCollection))
	for _, serverURI := range setOrderReq.Parameters.ServerCollection {
		go e.collectAndSetDefaultOrder(taskID, serverURI, subTaskChan, sessionUserName)
	}
	resp.StatusCode = http.StatusOK
	for i := 0; i < len(setOrderReq.Parameters.ServerCollection); i++ {
		select {
		case statusCode := <-subTaskChan:
			if statusCode != http.StatusOK {
				partialResultFlag = true
				if resp.StatusCode < statusCode {
					resp.StatusCode = statusCode
				}
			}
			if i < len(setOrderReq.Parameters.ServerCollection)-1 {
				percentComplete := int32(((i + 1) / len(setOrderReq.Parameters.ServerCollection)) * 100)
				var task = fillTaskData(taskID, targetURI, resp, common.Running, common.OK, percentComplete, http.MethodPost)
				err := e.UpdateTask(task)
				if err != nil && err.Error() == common.Cancelling {
					task = fillTaskData(taskID, targetURI, resp, common.Cancelled, common.OK, percentComplete, http.MethodPost)
					e.UpdateTask(task)
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
		errMsg := "one or more of the SetDefaultBootOrder requests failed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID
		log.Println(errMsg)
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return common.GeneralError(http.StatusUnauthorized, response.ResourceAtURIUnauthorized, errMsg, []interface{}{setOrderReq.Parameters.ServerCollection}, taskInfo)
		case http.StatusNotFound:
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"option", "SetDefaultBootOrder"}, taskInfo)
		default:
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		}
	}

	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	log.Println("all SetDefaultBootOrder requests successfully completed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID)
	resp.StatusMessage = response.Success
	resp.StatusCode = http.StatusOK
	args := response.Args{
		Code:    resp.StatusMessage,
		Message: "Request completed successfully",
	}
	resp.Body = args.CreateGenericErrorResponse()

	var task = fillTaskData(taskID, targetURI, resp, common.Completed, taskStatus, percentComplete, http.MethodPost)
	err = e.UpdateTask(task)
	if err != nil && err.Error() == common.Cancelling {
		task = fillTaskData(taskID, targetURI, resp, common.Cancelled, common.Critical, percentComplete, http.MethodPost)
		e.UpdateTask(task)
		runtime.Goexit()
	}
	return resp

}

func (e *ExternalInterface) collectAndSetDefaultOrder(taskID, serverURI string, subTaskChan chan<- int32, sessionUserName string) {
	var resp response.RPC
	subTaskURI, err := e.CreateChildTask(sessionUserName, taskID)
	if err != nil {
		subTaskChan <- http.StatusInternalServerError
		log.Println("error while trying to create sub task")
		return
	}
	var subTaskID string
	strArray := strings.Split(subTaskURI, "/")
	if strings.HasSuffix(subTaskURI, "/") {
		subTaskID = strArray[len(strArray)-2]
	} else {
		subTaskID = strArray[len(strArray)-1]
	}

	taskInfo := &common.TaskUpdateInfo{TaskID: subTaskID, TargetURI: serverURI, UpdateTask: e.UpdateTask}

	var percentComplete int32
	uuid, systemID, err := getIDsFromURI(serverURI)
	if err != nil {
		subTaskChan <- http.StatusNotFound
		errMsg := "error while trying to get system ID from " + serverURI + ": " + err.Error()
		log.Println(errMsg)
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"SystemID", serverURI}, taskInfo)
		return
	}
	// Get target device Credentials from using device UUID
	target, err := agmodel.GetTarget(uuid)
	if err != nil {
		subTaskChan <- http.StatusNotFound
		errMsg := err.Error()
		log.Println(errMsg)
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"target", uuid}, taskInfo)
		return
	}
	decryptedPasswordByte, err := e.DecryptPassword(target.Password)
	if err != nil {
		subTaskChan <- http.StatusInternalServerError
		errMsg := "error while trying to decrypt device password: " + err.Error()
		log.Println(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		return
	}
	target.Password = decryptedPasswordByte

	// Get the Plugin info
	plugin, errs := agmodel.GetPluginData(target.PluginID)
	if errs != nil {
		subTaskChan <- http.StatusNotFound
		errMsg := "error while getting plugin data: " + errs.Error()
		log.Println(errMsg)
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"PluginData", target.PluginID}, taskInfo)
		return
	}

	var pluginContactRequest getResourceRequest
	pluginContactRequest.ContactClient = e.ContactClient
	pluginContactRequest.GetPluginStatus = e.GetPluginStatus
	pluginContactRequest.Plugin = plugin
	pluginContactRequest.StatusPoll = true

	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		pluginContactRequest.HTTPMethodType = http.MethodPost
		pluginContactRequest.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		pluginContactRequest.OID = "/ODIM/v1/Sessions"
		_, token, getResponse, err := contactPlugin(pluginContactRequest, "error while logging in to plugin: ")
		if err != nil {
			subTaskChan <- getResponse.StatusCode
			errMsg := err.Error()
			log.Println(errMsg)
			common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo)
			return
		}
		pluginContactRequest.Token = token
	} else {
		pluginContactRequest.LoginCredentials = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	postRequest := make(map[string]interface{})
	postBody, _ := json.Marshal(postRequest)
	target.PostBody = postBody
	pluginContactRequest.DeviceInfo = target
	pluginContactRequest.OID = "/ODIM/v1/Systems/" + systemID + "/Actions/ComputerSystem.SetDefaultBootOrder"
	pluginContactRequest.HTTPMethodType = http.MethodPost
	_, _, getResponse, err := contactPlugin(pluginContactRequest, "error while setting the default boot order: ")
	if err != nil {
		subTaskChan <- getResponse.StatusCode
		errMsg := err.Error()
		log.Println(errMsg)
		common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo)
		return
	}
	// json.Unmarshal(body, &resp.Body)
	resp.StatusMessage = response.Success
	resp.Body = response.ErrorClass{
		Code:    resp.StatusMessage,
		Message: "Request completed successfully.",
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
		"Location":          serverURI,
	}
	resp.StatusCode = getResponse.StatusCode
	percentComplete = 100
	subTaskChan <- int32(getResponse.StatusCode)
	var task = fillTaskData(subTaskID, serverURI, resp, common.Completed, common.OK, percentComplete, http.MethodPost)
	err = e.UpdateTask(task)
	if err != nil && err.Error() == common.Cancelling {
		var task = fillTaskData(subTaskID, serverURI, resp, common.Cancelled, common.Critical, percentComplete, http.MethodPost)
		err = e.UpdateTask(task)
	}
	return
}
