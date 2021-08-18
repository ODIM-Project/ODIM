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
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
	uuid "github.com/satori/go.uuid"
)

// AddCompute is the handler for adding system
// Discovers Computersystem, Manager & Chassis and its top level odata.ID links and store them in inmemory db.
// Upon successfull operation this api returns Systems root UUID in the response body with 200 OK.
func (e *ExternalInterface) addCompute(taskID, targetURI, pluginID string, percentComplete int32, addResourceRequest AddResourceRequest, pluginContactRequest getResourceRequest) (response.RPC, string, []byte) {
	var resp response.RPC
	log.Info("started adding system with manager address " + addResourceRequest.ManagerAddress +
		" using plugin id: " + pluginID)

	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: targetURI, UpdateTask: e.UpdateTask, TaskRequest: pluginContactRequest.TaskRequest}

	var task = fillTaskData(taskID, targetURI, pluginContactRequest.TaskRequest, resp, common.Running, common.OK, percentComplete, http.MethodPost)

	plugin, errs := agmodel.GetPluginData(pluginID)
	if errs != nil {
		errMsg := "error while getting plugin data: " + errs.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"plugin", pluginID}, taskInfo), "", nil
	}

	var saveSystem agmodel.SaveSystem
	saveSystem.ManagerAddress = addResourceRequest.ManagerAddress
	saveSystem.UserName = addResourceRequest.UserName
	//saveSystem.Password = ciphertext
	saveSystem.Password = []byte(addResourceRequest.Password)
	saveSystem.PluginID = pluginID

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
			log.Error(errMsg)
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
		log.Error(errMsg)
		return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo), "", nil
	}

	var commonError errors.CommonError
	err = json.Unmarshal(body, &commonError)
	if err != nil {
		errMsg := err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
	}

	commonError.Error.Code = errors.PropertyValueFormatError
	resp.Body = commonError
	resp.StatusCode = http.StatusCreated
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
	systemsEstimatedWork := int32(60)
	var computeSystemID, resourceURI string
	if computeSystemID, resourceURI, progress, err = h.getAllSystemInfo(taskID, progress, systemsEstimatedWork, pluginContactRequest); err != nil {
		errMsg := "error while trying to add compute: " + err.Error()
		log.Error(errMsg)
		var msgArg = make([]interface{}, 0)
		var skipFlag bool
		switch h.StatusMessage {
		case response.ResourceAlreadyExists:
			msgArg = append(msgArg, addResourceRequest.ManagerAddress, pluginID, "ComputerSystem")
		case response.ActionParameterNotSupported:
			msgArg = append(msgArg, addResourceRequest.ManagerAddress, pluginID)
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
	task = fillTaskData(taskID, targetURI, pluginContactRequest.TaskRequest, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	e.UpdateTask(task)

	// Populate the resource Firmware inventory for update service
	pluginContactRequest.DeviceInfo = getSystemBody
	pluginContactRequest.OID = "/redfish/v1/UpdateService/FirmwareInventory"
	pluginContactRequest.DeviceUUID = saveSystem.DeviceUUID
	pluginContactRequest.HTTPMethodType = http.MethodGet

	progress = percentComplete
	firmwareEstimatedWork := int32(5)
	progress = h.getAllRootInfo(taskID, progress, firmwareEstimatedWork, pluginContactRequest, config.Data.AddComputeSkipResources.SkipResourceListUnderOthers)
	percentComplete = progress
	task = fillTaskData(taskID, targetURI, pluginContactRequest.TaskRequest, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	e.UpdateTask(task)

	// Populate the resource Software inventory for update service
	pluginContactRequest.DeviceInfo = getSystemBody
	pluginContactRequest.OID = "/redfish/v1/UpdateService/SoftwareInventory"
	pluginContactRequest.DeviceUUID = saveSystem.DeviceUUID
	pluginContactRequest.HTTPMethodType = http.MethodGet

	progress = percentComplete
	softwareEstimatedWork := int32(5)
	progress = h.getAllRootInfo(taskID, progress, softwareEstimatedWork, pluginContactRequest, config.Data.AddComputeSkipResources.SkipResourceListUnderOthers)
	percentComplete = progress
	task = fillTaskData(taskID, targetURI, pluginContactRequest.TaskRequest, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	e.UpdateTask(task)

	// Discover telemetry service
	percentComplete = e.getTelemetryService(taskID, targetURI, percentComplete, pluginContactRequest, resp, saveSystem)

	// Lets Discover/gather registry files of this server and store them in DB

	pluginContactRequest.DeviceInfo = getSystemBody
	pluginContactRequest.OID = "/redfish/v1/Registries"
	pluginContactRequest.DeviceUUID = saveSystem.DeviceUUID
	pluginContactRequest.HTTPMethodType = http.MethodGet

	progress = percentComplete
	registriesEstimatedWork := int32(5)
	progress = h.getAllRegistries(taskID, progress, registriesEstimatedWork, pluginContactRequest)
	percentComplete = progress
	task = fillTaskData(taskID, targetURI, pluginContactRequest.TaskRequest, resp, common.Running, common.OK, percentComplete, http.MethodPost)
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
	progress = h.getAllRootInfo(taskID, progress, chassisEstimatedWork, pluginContactRequest, config.Data.AddComputeSkipResources.SkipResourceListUnderChassis)

	percentComplete = progress
	task = fillTaskData(taskID, targetURI, pluginContactRequest.TaskRequest, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	err = e.UpdateTask(task)
	if err != nil && (err.Error() == common.Cancelling) {
		go e.rollbackInMemory(resourceURI)
		return resp, "", nil
	}

	//Logic for getting the manager information
	// Logic for getting manager information and saving it into the database
	// Discover manager Collection this can be a function later.
	getManagerBody := map[string]interface{}{
		"ManagerAddress": addResourceRequest.ManagerAddress,
		"UserName":       addResourceRequest.UserName,
		"Password":       saveSystem.Password,
	}
	pluginContactRequest.DeviceInfo = getManagerBody
	pluginContactRequest.OID = "/redfish/v1/Managers"
	pluginContactRequest.DeviceUUID = saveSystem.DeviceUUID
	pluginContactRequest.HTTPMethodType = http.MethodGet

	progress = percentComplete
	managerEstimatedWork := int32(15)
	progress = h.getAllRootInfo(taskID, progress, managerEstimatedWork, pluginContactRequest, config.Data.AddComputeSkipResources.SkipResourceListUnderManager)

	percentComplete = progress
	task = fillTaskData(taskID, targetURI, pluginContactRequest.TaskRequest, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	err = e.UpdateTask(task)
	if err != nil && (err.Error() == common.Cancelling) {
		go e.rollbackInMemory(resourceURI)
		return resp, "", nil
	}
	if h.ErrorMessage != "" && h.StatusCode != http.StatusServiceUnavailable && h.StatusCode != http.StatusNotFound && h.StatusCode != http.StatusInternalServerError && h.StatusCode != http.StatusBadRequest {
		go e.rollbackInMemory(resourceURI)
		log.Error(h.ErrorMessage)
		return common.GeneralError(h.StatusCode, h.StatusMessage, h.ErrorMessage, h.MsgArgs, taskInfo), "", nil
	}

	ciphertext, err := e.EncryptPassword([]byte(addResourceRequest.Password))
	if err != nil {
		go e.rollbackInMemory(resourceURI)
		errMsg := "error while trying to encrypt: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
	}
	saveSystem.Password = ciphertext
	aggregationSourceID := saveSystem.DeviceUUID + ":" + computeSystemID
	if err := saveSystem.Create(saveSystem.DeviceUUID); err != nil {
		go e.rollbackInMemory(resourceURI)
		errMsg := "error while trying to add compute: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
	}
	aggSourceIDChassisAndManager := saveSystem.DeviceUUID + ":"
	chassisList, _ := agmodel.GetAllMatchingDetails("Chassis", aggSourceIDChassisAndManager, common.InMemory)
	managersList, _ := agmodel.GetAllMatchingDetails("Managers", aggSourceIDChassisAndManager, common.InMemory)
	urlList := h.SystemURL
	urlList = append(urlList, chassisList...)
	urlList = append(urlList, managersList...)
	pluginContactRequest.CreateSubcription(urlList)

	pluginContactRequest.PublishEvent(h.SystemURL, "SystemsCollection")

	// get all managers and chassis info
	pluginContactRequest.PublishEvent(chassisList, "ChassisCollection")
	pluginContactRequest.PublishEvent(managersList, "ManagerCollection")

	h.PluginResponse = strings.Replace(h.PluginResponse, `/redfish/v1/Systems/`, `/redfish/v1/Systems/`+saveSystem.DeviceUUID+`:`, -1)
	var list agresponse.List
	err = json.Unmarshal([]byte(h.PluginResponse), &list)
	if err != nil {
		log.Error(err.Error())
	}

	resp.Header = map[string]string{
		"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		"Location":     resourceURI,
	}
	log.Info("sucessfully added system with manager address " + addResourceRequest.ManagerAddress +
		" using plugin id: " + pluginID)

	pluginStartUpData := &agmodel.PluginStartUpData{
		RequestType: "delta",
		Devices: map[string]agmodel.DeviceData{
			saveSystem.DeviceUUID: agmodel.DeviceData{
				Address:   addResourceRequest.ManagerAddress,
				UserName:  addResourceRequest.UserName,
				Password:  []byte(addResourceRequest.Password),
				Operation: "add",
			},
		},
	}
	if err = PushPluginStartUpData(plugin, pluginStartUpData); err != nil {
		log.Error(err.Error())
	}

	return resp, aggregationSourceID, ciphertext
}
