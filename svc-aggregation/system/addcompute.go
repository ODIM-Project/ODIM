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
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
	uuid "github.com/satori/go.uuid"
)

// AggregationServiceAdd to add bmc or manger via AggregationService Add action
func (e *ExternalInterface) AggregationServiceAdd(taskID string, sessionUserName string, req *aggregatorproto.AggregatorRequest) response.RPC {
	var resp response.RPC
	var percentComplete int32
	targetURI := "/redfish/v1/AggregationService/Actions/AggregationService.Add"
	var task = fillTaskData(taskID, targetURI, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	err := e.UpdateTask(task)
	if err != nil {
		errMsg := "error while starting the task: " + err.Error()
		log.Printf("error while starting the task: %v", errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}

	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: targetURI, UpdateTask: e.UpdateTask}

	// parsing the AddResourceRequest
	var addResourceRequest AddResourceRequest
	err = json.Unmarshal(req.RequestBody, &addResourceRequest)
	if err != nil {
		errMsg := "unable to parse the add request" + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, addResourceRequest)
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

	if addResourceRequest.Oem == nil {
		errMsg := "error: mandatory Oem block missing in the request"
		log.Println(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"Oem"}, taskInfo)
	}

	// check if there is a request ongoing for the server in payload
	ActiveReqSet.UpdateMu.Lock()
	if pluginID, exist := ActiveReqSet.ReqRecord[addResourceRequest.ManagerAddress]; exist {
		ActiveReqSet.UpdateMu.Unlock()
		var errMsg string
		mIP, mPort := getIPAndPortFromAddress(addResourceRequest.ManagerAddress)
		// checking whether the request is for adding a server or a manager
		if addResourceRequest.Oem.PluginType != "" || addResourceRequest.Oem.PreferredAuthType != "" {
			errMsg = fmt.Sprintf("error: An active request already exists for adding manager %v plugin with IP %v Port %v", pluginID.(string), mIP, mPort)
		} else {
			errMsg = fmt.Sprintf("error: An active request already exists for adding BMC with IP %v through %v plugin", mIP, pluginID.(string))
		}
		log.Println(errMsg)
		args := response.Args{
			Code:    response.GeneralError,
			Message: errMsg,
		}
		resp.Body = args.CreateGenericErrorResponse()
		resp.Header = map[string]string{"Content-type": "application/json; charset=utf-8"}
		resp.StatusCode = http.StatusConflict
		percentComplete = 100
		e.UpdateTask(fillTaskData(taskID, targetURI, resp, common.Exception, common.Warning, percentComplete, http.MethodPost))
		return resp
	}
	ActiveReqSet.ReqRecord[addResourceRequest.ManagerAddress] = addResourceRequest.Oem.PluginID
	ActiveReqSet.UpdateMu.Unlock()

	defer func() {
		// check if there is an entry added for the server in ongoing requests tracker and remove it
		ActiveReqSet.UpdateMu.Lock()
		delete(ActiveReqSet.ReqRecord, addResourceRequest.ManagerAddress)
		ActiveReqSet.UpdateMu.Unlock()
	}()

	var pluginContactRequest getResourceRequest

	pluginContactRequest.ContactClient = e.ContactClient
	pluginContactRequest.GetPluginStatus = e.GetPluginStatus
	pluginContactRequest.TargetURI = targetURI
	pluginContactRequest.UpdateTask = e.UpdateTask

	if addResourceRequest.Oem.PluginType != "" || addResourceRequest.Oem.PreferredAuthType != "" {
		resp, _, _ = e.addPluginData(addResourceRequest, taskID, targetURI, pluginContactRequest)
	} else {

		resp, _, _ = e.addCompute(taskID, targetURI, percentComplete, addResourceRequest, pluginContactRequest)
	}
	if resp.StatusMessage != "" {
		return resp
	}
	resp.StatusMessage = response.Success
	resp.StatusCode = http.StatusOK
	resp.Body = response.ErrorClass{
		Code:    resp.StatusMessage,
		Message: "Request completed successfully.",
	}
	percentComplete = 100

	task = fillTaskData(taskID, targetURI, resp, common.Completed, common.OK, percentComplete, http.MethodPost)
	e.UpdateTask(task)

	return resp
}

// AddCompute is the handler for adding system
// Discovers Computersystem & Chassis and its top level odata.ID links and store them in inmemory db.
// Upon successfull operation this api returns Systems root UUID in the response body with 200 OK.
func (e *ExternalInterface) addCompute(taskID, targetURI string, percentComplete int32, addResourceRequest AddResourceRequest, pluginContactRequest getResourceRequest) (response.RPC, string, []byte) {
	var resp response.RPC
	log.Printf("started adding system with manager address %v using plugin id %v.", addResourceRequest.ManagerAddress, addResourceRequest.Oem.PluginID)

	var task = fillTaskData(taskID, targetURI, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: targetURI, UpdateTask: e.UpdateTask}

	plugin, errs := agmodel.GetPluginData(addResourceRequest.Oem.PluginID)
	if errs != nil {
		errMsg := "error while getting plugin data: " + errs.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"plugin", addResourceRequest.Oem.PluginID}, taskInfo), "", nil
	}

	var saveSystem agmodel.SaveSystem
	saveSystem.ManagerAddress = addResourceRequest.ManagerAddress
	saveSystem.UserName = addResourceRequest.UserName
	//saveSystem.Password = ciphertext
	saveSystem.Password = []byte(addResourceRequest.Password)
	saveSystem.PluginID = addResourceRequest.Oem.PluginID

	pluginContactRequest.Plugin = plugin
	pluginContactRequest.StatusPoll = true
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		pluginContactRequest.HTTPMethodType = http.MethodPost
		pluginContactRequest.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		pluginContactRequest.OID = "/ODIM/v1/Sessions"
		_, token, getResponse, err := contactPlugin(pluginContactRequest, "error while getting the details "+pluginContactRequest.OID+": ")
		if err != nil {
			errMsg := err.Error()
			log.Println(errMsg)
			return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo), "", nil
		}
		pluginContactRequest.Token = token
	} else {
		pluginContactRequest.LoginCredentials = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
	}
	//get token from response and make the next REST call <currently custome url>

	pluginContactRequest.DeviceInfo = saveSystem
	pluginContactRequest.OID = "/ODIM/v1/validate"
	pluginContactRequest.HTTPMethodType = http.MethodPost

	body, _, getResponse, err := contactPlugin(pluginContactRequest, "error while trying to authenticate the compute server: ")
	if err != nil {
		errMsg := err.Error()
		log.Println(errMsg)
		return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo), "", nil
	}

	var commonError errors.CommonError
	err = json.Unmarshal(body, &commonError)
	if err != nil {
		errMsg := err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
	}

	commonError.Error.Code = errors.PropertyValueFormatError
	resp.Body = commonError
	resp.StatusCode = getResponse.StatusCode
	resp.StatusMessage = getResponse.StatusMessage
	resp.Header = map[string]string{"Content-type": "application/json; charset=utf-8"}

	saveSystem.DeviceUUID = uuid.NewV4().String()
	getSystemBody := map[string]interface{}{
		"ManagerAddress": saveSystem.ManagerAddress,
		"UserName":       saveSystem.UserName,
		"Password":       saveSystem.Password,
	}

	//Discover Systems collection this will be moved to a function later if needed
	pluginContactRequest.DeviceInfo = getSystemBody
	pluginContactRequest.OID = "/redfish/v1/Systems"
	pluginContactRequest.DeviceUUID = saveSystem.DeviceUUID
	pluginContactRequest.HTTPMethodType = http.MethodGet
	pluginContactRequest.CreateSubcription = e.CreateSubcription
	pluginContactRequest.PublishEvent = e.PublishEvent

	var h respHolder
	h.TraversedLinks = make(map[string]bool)
	progress := percentComplete
	systemsEstimatedWork := int32(65)
	var resourceURI string
	if resourceURI, progress, err = h.getAllSystemInfo(taskID, progress, systemsEstimatedWork, pluginContactRequest); err != nil {
		errMsg := "error while trying to add compute: " + err.Error()
		log.Println(errMsg)
		var msgArg = make([]interface{}, 0)
		var skipFlag bool
		switch h.StatusMessage {
		case response.ResourceAlreadyExists:
			msgArg = append(msgArg, addResourceRequest.ManagerAddress, addResourceRequest.Oem.PluginID, "ComputerSystem")
		case response.ActionParameterNotSupported:
			msgArg = append(msgArg, addResourceRequest.ManagerAddress, addResourceRequest.Oem.PluginID)
		case response.ResourceAtURIUnauthorized, response.CouldNotEstablishConnection:
			msgArg = append(msgArg, addResourceRequest.ManagerAddress)
		default:
			skipFlag = true
		}
		if !skipFlag {
			go e.rollbackInMemory(resourceURI)
			return common.GeneralError(h.StatusCode, h.StatusMessage, errMsg, msgArg, taskInfo), "", nil
		}
	}
	percentComplete = progress
	task = fillTaskData(taskID, targetURI, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	e.UpdateTask(task)

	// Lets Discover/gather registry files of this server and store them in DB

	pluginContactRequest.DeviceInfo = getSystemBody
	pluginContactRequest.OID = "/redfish/v1/Registries"
	pluginContactRequest.DeviceUUID = saveSystem.DeviceUUID
	pluginContactRequest.HTTPMethodType = http.MethodGet

	progress = percentComplete
	registriesEstimatedWork := int32(15)
	progress = h.getAllRegistries(taskID, progress, registriesEstimatedWork, pluginContactRequest)
	percentComplete = progress
	task = fillTaskData(taskID, targetURI, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	err = e.UpdateTask(task)
	if err != nil && (err.Error() == common.Cancelling) {
		go e.rollbackInMemory(resourceURI)
		return resp, "", nil
	}

	// End of Registry files Discovery

	// Logic for getting chassis information and saving it into the database
	// Discover Chassis Collection this can be a function later.
	getChassisBody := map[string]interface{}{
		"ManagerAddress": addResourceRequest.ManagerAddress,
		"UserName":       addResourceRequest.UserName,
		"Password":       saveSystem.Password,
	}
	pluginContactRequest.DeviceInfo = getChassisBody
	pluginContactRequest.OID = "/redfish/v1/Chassis"
	pluginContactRequest.DeviceUUID = saveSystem.DeviceUUID
	pluginContactRequest.HTTPMethodType = http.MethodGet

	progress = percentComplete
	chassisEstimatedWork := int32(15)
	progress = h.getAllChassisInfo(taskID, progress, chassisEstimatedWork, pluginContactRequest)
	percentComplete = progress
	task = fillTaskData(taskID, targetURI, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	err = e.UpdateTask(task)
	if err != nil && (err.Error() == common.Cancelling) {
		go e.rollbackInMemory(resourceURI)
		return resp, "", nil
	}
	if h.ErrorMessage != "" && h.StatusCode != http.StatusServiceUnavailable && h.StatusCode != http.StatusNotFound && h.StatusCode != http.StatusInternalServerError && h.StatusCode != http.StatusBadRequest {
		go e.rollbackInMemory(resourceURI)
		log.Println(h.ErrorMessage)
		return common.GeneralError(h.StatusCode, h.StatusMessage, h.ErrorMessage, h.MsgArgs, taskInfo), "", nil
	}

	ciphertext, err := e.EncryptPassword([]byte(addResourceRequest.Password))
	if err != nil {
		go e.rollbackInMemory(resourceURI)
		errMsg := "error while trying to encrypt: " + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
	}
	saveSystem.Password = ciphertext
	if err := saveSystem.Create(saveSystem.DeviceUUID); err != nil {
		go e.rollbackInMemory(resourceURI)
		errMsg := "error while trying to add compute: " + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
	}
	pluginContactRequest.CreateSubcription(h.SystemURL)
	pluginContactRequest.PublishEvent(h.SystemURL, "SystemsCollection")

	h.PluginResponse = strings.Replace(h.PluginResponse, `/redfish/v1/Systems/`, `/redfish/v1/Systems/`+saveSystem.DeviceUUID+`:`, -1)
	var list agresponse.List
	err = json.Unmarshal([]byte(h.PluginResponse), &list)
	if err != nil {
		log.Println("Error: ", err)
	}

	resp.Header = map[string]string{
		"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		"Location":     resourceURI,
	}
	log.Printf("sucessfully added system with manager address %v using plugin id %v.", addResourceRequest.ManagerAddress, addResourceRequest.Oem.PluginID)
	return resp, saveSystem.DeviceUUID, ciphertext
}
