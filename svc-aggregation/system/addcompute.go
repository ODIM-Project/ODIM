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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
	uuid "github.com/satori/go.uuid"
)

// AddCompute is the handler for adding system
// Discovers Computersystem, Manager & Chassis and its top level odata.ID links and store them in inmemory db.
// Upon successfull operation this api returns Systems root UUID in the response body with 200 OK.
func (e *ExternalInterface) addCompute(ctx context.Context, taskID, targetURI, pluginID string, percentComplete int32, addResourceRequest AddResourceRequest, pluginContactRequest getResourceRequest) (response.RPC, string, []byte) {
	var resp response.RPC
	l.LogWithFields(ctx).Info("started adding system with manager address " + addResourceRequest.ManagerAddress +
		" using plugin id: " + pluginID)

	taskInfo := &common.TaskUpdateInfo{Context: ctx, TaskID: taskID, TargetURI: targetURI, UpdateTask: e.UpdateTask, TaskRequest: pluginContactRequest.TaskRequest}

	var task = fillTaskData(taskID, targetURI, pluginContactRequest.TaskRequest, resp, common.Running, common.OK, percentComplete, http.MethodPost)

	plugin, errs := agmodel.GetPluginData(pluginID)
	if errs != nil {
		errMsg := "error while getting plugin data: " + errs.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"plugin", pluginID}, taskInfo), "", nil
	}

	var saveSystem agmodel.SaveSystem
	saveSystem.ManagerAddress = strings.ToLower(addResourceRequest.ManagerAddress)
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
		l.LogWithFields(ctx).Debugf("plugin contact request data for %s : %s", pluginContactRequest.OID, string(pluginContactRequest.Data))
		_, token, getResponse, err := contactPlugin(ctx, pluginContactRequest, "error while getting the details "+pluginContactRequest.OID+": ")
		if err != nil {
			errMsg := err.Error()
			l.LogWithFields(ctx).Error(errMsg)
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
	l.LogWithFields(ctx).Debugf("plugin contact request data for %s : %s", pluginContactRequest.OID, string(pluginContactRequest.Data))
	body, _, getResponse, err := contactPlugin(ctx, pluginContactRequest, "error while trying to authenticate the compute server: ")
	if err != nil {
		errMsg := err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo), "", nil
	}

	var commonError errors.CommonError
	err = json.Unmarshal(body, &commonError)
	if err != nil {
		errMsg := err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
	}

	commonError.Error.Code = errors.PropertyValueFormatError
	resp.Body = commonError
	resp.StatusCode = http.StatusCreated
	resp.StatusMessage = getResponse.StatusMessage

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
	pluginContactRequest.BMCAddress = saveSystem.ManagerAddress

	var h respHolder
	h.TraversedLinks = make(map[string]bool)
	h.InventoryData = make(map[string]interface{})
	progress := percentComplete
	systemsEstimatedWork := int32(60)
	var computeSystemID, resourceURI string
	if computeSystemID, resourceURI, progress, err = h.getAllSystemInfo(ctx, taskID, progress, systemsEstimatedWork, pluginContactRequest); err != nil {
		errMsg := "error while trying to add compute: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
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
	e.UpdateTask(ctx, task)
	h.InventoryData = make(map[string]interface{})

	// Populate the resource Firmware inventory for update service
	pluginContactRequest.DeviceInfo = getSystemBody
	pluginContactRequest.OID = "/redfish/v1/UpdateService/FirmwareInventory"
	pluginContactRequest.DeviceUUID = saveSystem.DeviceUUID
	pluginContactRequest.HTTPMethodType = http.MethodGet

	progress = percentComplete
	firmwareEstimatedWork := int32(5)
	progress = h.getAllRootInfo(ctx, taskID, progress, firmwareEstimatedWork, pluginContactRequest, config.Data.AddComputeSkipResources.SkipResourceListUnderOthers)
	percentComplete = progress
	task = fillTaskData(taskID, targetURI, pluginContactRequest.TaskRequest, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	e.UpdateTask(ctx, task)

	// Populate the resource Software inventory for update service
	pluginContactRequest.DeviceInfo = getSystemBody
	pluginContactRequest.OID = "/redfish/v1/UpdateService/SoftwareInventory"
	pluginContactRequest.DeviceUUID = saveSystem.DeviceUUID
	pluginContactRequest.HTTPMethodType = http.MethodGet

	progress = percentComplete
	softwareEstimatedWork := int32(5)
	progress = h.getAllRootInfo(ctx, taskID, progress, softwareEstimatedWork, pluginContactRequest, config.Data.AddComputeSkipResources.SkipResourceListUnderOthers)
	percentComplete = progress
	task = fillTaskData(taskID, targetURI, pluginContactRequest.TaskRequest, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	e.UpdateTask(ctx, task)

	// Discover telemetry service
	percentComplete = e.getTelemetryService(ctx, taskID, targetURI, percentComplete, pluginContactRequest, resp, saveSystem)
	task = fillTaskData(taskID, targetURI, pluginContactRequest.TaskRequest, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	e.UpdateTask(ctx, task)
	// Populate the data for license service
	pluginContactRequest.DeviceInfo = getSystemBody
	pluginContactRequest.OID = "/redfish/v1/LicenseService/Licenses/"
	pluginContactRequest.DeviceUUID = saveSystem.DeviceUUID
	pluginContactRequest.HTTPMethodType = http.MethodGet

	progress = percentComplete
	licenseEstimatedWork := int32(5)
	progress = h.getAllRootInfo(ctx, taskID, progress, licenseEstimatedWork, pluginContactRequest, config.Data.AddComputeSkipResources.SkipResourceListUnderOthers)
	percentComplete = progress
	task = fillTaskData(taskID, targetURI, pluginContactRequest.TaskRequest, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	e.UpdateTask(ctx, task)

	// Lets Discover/gather registry files of this server and store them in DB

	pluginContactRequest.DeviceInfo = getSystemBody
	pluginContactRequest.OID = "/redfish/v1/Registries"
	pluginContactRequest.DeviceUUID = saveSystem.DeviceUUID
	pluginContactRequest.HTTPMethodType = http.MethodGet

	progress = percentComplete
	registriesEstimatedWork := int32(5)
	progress = h.getAllRegistries(ctx, taskID, progress, registriesEstimatedWork, pluginContactRequest)
	percentComplete = progress
	task = fillTaskData(taskID, targetURI, pluginContactRequest.TaskRequest, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	err = e.UpdateTask(ctx, task)
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
	progress = h.getAllRootInfo(ctx, taskID, progress, chassisEstimatedWork, pluginContactRequest, config.Data.AddComputeSkipResources.SkipResourceListUnderChassis)

	percentComplete = progress
	task = fillTaskData(taskID, targetURI, pluginContactRequest.TaskRequest, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	err = e.UpdateTask(ctx, task)
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
	progress = h.getAllRootInfo(ctx, taskID, progress, managerEstimatedWork, pluginContactRequest, config.Data.AddComputeSkipResources.SkipResourceListUnderManager)

	percentComplete = progress
	task = fillTaskData(taskID, targetURI, pluginContactRequest.TaskRequest, resp, common.Running, common.OK, percentComplete, http.MethodPost)
	err = e.UpdateTask(ctx, task)
	if err != nil && (err.Error() == common.Cancelling) {
		go e.rollbackInMemory(resourceURI)
		return resp, "", nil
	}
	if h.ErrorMessage != "" && h.StatusCode != http.StatusServiceUnavailable && h.StatusCode != http.StatusNotFound && h.StatusCode != http.StatusInternalServerError && h.StatusCode != http.StatusBadRequest {
		go e.rollbackInMemory(resourceURI)
		l.LogWithFields(ctx).Error(h.ErrorMessage)
		return common.GeneralError(h.StatusCode, h.StatusMessage, h.ErrorMessage, h.MsgArgs, taskInfo), "", nil
	}
	err = agmodel.SaveBMCInventory(h.InventoryData)
	if err != nil {
		errorMessage := "GenericSave : error while trying to add resource data to DB: " + err.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
			nil, nil), "", nil
	}
	ciphertext, err := e.EncryptPassword([]byte(addResourceRequest.Password))
	if err != nil {
		go e.rollbackInMemory(resourceURI)
		errMsg := "error while trying to encrypt: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
	}
	saveSystem.Password = ciphertext
	aggregationSourceID := saveSystem.DeviceUUID + "." + computeSystemID
	if err := saveSystem.Create(ctx, saveSystem.DeviceUUID); err != nil {
		go e.rollbackInMemory(resourceURI)
		errMsg := "error while trying to add compute: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
	}
	aggSourceIDChassisAndManager := saveSystem.DeviceUUID + "."
	chassisList, _ := agmodel.GetAllMatchingDetails("Chassis", aggSourceIDChassisAndManager, common.InMemory)
	managersList, _ := agmodel.GetAllMatchingDetails("Managers", aggSourceIDChassisAndManager, common.InMemory)
	urlList := h.SystemURL
	urlList = append(urlList, chassisList...)
	urlList = append(urlList, managersList...)
	pluginContactRequest.CreateSubcription(ctx, urlList)

	pluginContactRequest.PublishEvent(ctx, h.SystemURL, "SystemsCollection")

	// get all managers and chassis info
	pluginContactRequest.PublishEvent(ctx, chassisList, "ChassisCollection")
	pluginContactRequest.PublishEvent(ctx, managersList, "ManagerCollection")

	h.PluginResponse = strings.Replace(h.PluginResponse, `/redfish/v1/Systems/`, `/redfish/v1/Systems/`+saveSystem.DeviceUUID+`.`, -1)
	var list agresponse.List
	err = json.Unmarshal([]byte(h.PluginResponse), &list)
	if err != nil {
		l.LogWithFields(ctx).Error(err.Error())
	}

	resp.Header = map[string]string{
		"Location": resourceURI,
	}
	l.LogWithFields(ctx).Info("sucessfully added system with manager address " + addResourceRequest.ManagerAddress +
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
	if err = PushPluginStartUpData(ctx, plugin, pluginStartUpData); err != nil {
		l.LogWithFields(ctx).Error(err.Error())
	}
	managerURI := "/redfish/v1/Managers/" + plugin.ManagerUUID
	var managerData map[string]interface{}
	managerLinks := make(map[string]interface{})
	var chassisLink, serverLink, listOfChassis, listOfServer []interface{}

	data, jerr := agmodel.GetResource(ctx, "Managers", managerURI)
	if jerr != nil {
		errorMessage := "error getting manager details: " + jerr.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
			nil, nil), "", nil
	}

	err = json.Unmarshal([]byte(data), &managerData)
	if err != nil {
		errorMessage := "error unmarshalling manager details: " + err.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
			nil, nil), "", nil
	}

	for _, val := range chassisList {
		listOfChassis = append(listOfChassis, map[string]string{"@odata.id": val})
	}
	for _, val := range h.SystemURL {
		listOfServer = append(listOfServer, map[string]string{"@odata.id": val})
	}
	if links, ok := managerData["Links"].(map[string]interface{}); ok {
		if managerData["Links"].(map[string]interface{})["ManagerForChassis"] != nil {
			chassisLink = links["ManagerForChassis"].([]interface{})
		}
		chassisLink = append(chassisLink, listOfChassis...)
		managerData["Links"].(map[string]interface{})["ManagerForChassis"] = chassisLink

		if managerData["Links"].(map[string]interface{})["ManagerForServers"] != nil {
			serverLink = links["ManagerForServers"].([]interface{})
		}
		serverLink = append(serverLink, listOfServer...)
		managerData["Links"].(map[string]interface{})["ManagerForServers"] = serverLink
	} else {
		chassisLink = append(chassisLink, listOfChassis...)
		serverLink = append(serverLink, listOfServer...)
		managerLinks["ManagerForChassis"] = chassisLink
		managerLinks["ManagerForServers"] = serverLink
		managerData["Links"] = managerLinks
	}
	mgrData, err := json.Marshal(managerData)
	if err != nil {
		errorMessage := "unable to marshal data while updating managers detail: " + err.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
			nil, nil), "", nil
	}
	err = agmodel.GenericSave([]byte(mgrData), "Managers", managerURI)
	if err != nil {
		errorMessage := "GenericSave : error while trying to add resource date to DB: " + err.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
			nil, nil), "", nil
	}
	l.LogWithFields(ctx).Debugf("final response for add compute request: %s", string(fmt.Sprintf("%v", resp.Body)))
	return resp, aggregationSourceID, ciphertext
}
