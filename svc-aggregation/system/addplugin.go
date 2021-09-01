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

package system

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agcommon"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
)

func (e *ExternalInterface) addPluginData(req AddResourceRequest, taskID, targetURI string, pluginContactRequest getResourceRequest, queueList []string, cmVariants connectionMethodVariants) (response.RPC, string, []byte) {
	var resp response.RPC
	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: targetURI, UpdateTask: e.UpdateTask, TaskRequest: pluginContactRequest.TaskRequest}

	if !(cmVariants.PreferredAuthType == "BasicAuth" || cmVariants.PreferredAuthType == "XAuthToken") {
		errMsg := "error: incorrect request property value for PreferredAuthType"
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyValueNotInList, errMsg, []interface{}{"PreferredAuthType", "[BasicAuth, XAuthToken]"}, taskInfo), "", nil
	}

	// checking the plugin type
	if !isPluginTypeSupported(cmVariants.PluginType) {
		errMsg := "error: incorrect request property value for PluginType"
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyValueNotInList, errMsg, []interface{}{"PluginType", fmt.Sprintf("%v", config.Data.SupportedPluginTypes)}, taskInfo), "", nil
	}

	// checking whether the Plugin already exists
	// If GetPluginData was successful, it indicates plugin already exists,
	// but it could also return errors, for below reasons, and has to be considered
	// as successful fetch of plugin data
	// error is nil - Plugin was successfully fetched.
	// error is not nil, Plugin data read but JSON unmarshaling failed
	// error is not nil, Plugin data read but decryption of plugin password failed
	// error is not nil, DB query failed, can't say for sure if queried plugin exists,
	// except when read fails with plugin data not found, and will continue with add process,
	// and any other errors, will fail add plugin operation.
	_, errs := agmodel.GetPluginData(cmVariants.PluginID)
	if errs == nil || (errs != nil && (errs.ErrNo() == errors.JSONUnmarshalFailed || errs.ErrNo() == errors.DecryptionFailed)) {
		errMsg := "error:plugin with name " + cmVariants.PluginID + " already exists"
		log.Error(errMsg)
		return common.GeneralError(http.StatusConflict, response.ResourceAlreadyExists, errMsg, []interface{}{"Plugin", "PluginID", cmVariants.PluginID}, taskInfo), "", nil
	}
	if errs != nil && errs.ErrNo() != errors.DBKeyNotFound {
		errMsg := "error: DB lookup failed for " + cmVariants.PluginID + " plugin: " + errs.Error()
		log.Error(errMsg)
		if errs.ErrNo() == errors.DBConnFailed {
			return common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, errMsg,
				[]interface{}{"Backend", config.Data.DBConf.OnDiskHost + ":" + config.Data.DBConf.OnDiskPort}, taskInfo), "", nil
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, []interface{}{}, taskInfo), "", nil
	}

	pluginNameArray, err := agmodel.GetAllKeysFromTable("Plugin")
	if err == nil {
		for _, ID := range pluginNameArray {

			plugin, err := e.GetPluginMgrAddr(ID)

			if err != nil && err.ErrNo() == errors.JSONUnmarshalFailed {
				continue
			}
			if err != nil {
				return common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, err.Error(),
					[]interface{}{"Backend", config.Data.DBConf.OnDiskHost + ":" + config.Data.DBConf.OnDiskPort}, taskInfo), "", nil
			}
			if plugin.IP+":"+plugin.Port == req.ManagerAddress {
				errMsg := "error:plugin with manager adress " + req.ManagerAddress + " already exists with name " + plugin.ID + " and ManagerUUID " + plugin.ManagerUUID
				log.Error(errMsg)
				return common.GeneralError(http.StatusConflict, response.ResourceAlreadyExists, errMsg, []interface{}{"Plugin", "PluginID", ID}, taskInfo), "", nil
			}
		}
	} else {
		return common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, err.Error(),
			[]interface{}{"Backend", config.Data.DBConf.OnDiskHost + ":" + config.Data.DBConf.OnDiskPort}, taskInfo), "", nil
	}
	// encrypt plugin password
	ciphertext, err := e.EncryptPassword([]byte(req.Password))
	if err != nil {
		errMsg := "error: encryption failed: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
	}
	var managerUUID string
	ipData := strings.Split(req.ManagerAddress, ":")
	var plugin = agmodel.Plugin{
		IP:                ipData[0],
		Port:              ipData[1],
		Username:          req.UserName,
		Password:          []byte(req.Password),
		ID:                cmVariants.PluginID,
		PluginType:        cmVariants.PluginType,
		PreferredAuthType: cmVariants.PreferredAuthType,
	}
	pluginContactRequest.Plugin = plugin
	pluginContactRequest.StatusPoll = true
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		pluginContactRequest.HTTPMethodType = http.MethodPost
		pluginContactRequest.DeviceInfo = map[string]interface{}{
			"Username": plugin.Username,
			"Password": string(plugin.Password),
		}
		pluginContactRequest.OID = "/ODIM/v1/Sessions"
		_, token, getResponse, err := contactPlugin(pluginContactRequest, "error while creating the session: ")
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
	// Getting all managers info from plugin
	pluginContactRequest.HTTPMethodType = http.MethodGet
	pluginContactRequest.OID = "/ODIM/v1/Managers"
	body, _, getResponse, err := contactPlugin(pluginContactRequest, "error while getting the details "+pluginContactRequest.OID+": ")
	if err != nil {
		errMsg := err.Error()
		log.Error(errMsg)
		return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo), "", nil
	}
	//  Extract all managers info and loop  over each members
	managersMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(body), &managersMap)
	if err != nil {
		errMsg := "unable to parse the managers resposne" + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
	}
	var managersData = make(map[string][]byte)
	managerMembers := managersMap["Members"]

	// Getting the indivitual managers response
	for _, object := range managerMembers.([]interface{}) {
		pluginContactRequest.OID = object.(map[string]interface{})["@odata.id"].(string)
		body, _, getResponse, err := contactPlugin(pluginContactRequest, "error while getting the details "+pluginContactRequest.OID+": ")
		if err != nil {
			errMsg := err.Error()
			log.Error(errMsg)
			return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo), "", nil
		}
		managerData := make(map[string]interface{})
		err = json.Unmarshal([]byte(body), &managerData)
		if err != nil {
			errMsg := "unable to parse the managers resposne" + err.Error()
			log.Error(errMsg)
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
		}
		if uuid, ok := managerData["UUID"]; ok {
			managerUUID = uuid.(string)
		}

		managersData[pluginContactRequest.OID] = body
	}
	// saving all plugin manager data
	var listMembers = make([]agresponse.ListMember, 0)
	for oid, data := range managersData {

		dbErr := agmodel.SavePluginManagerInfo(updateManagerName(data, plugin.ID), "Managers", oid)
		if dbErr != nil {
			errMsg := dbErr.Error()
			log.Error(errMsg)

			return common.GeneralError(http.StatusConflict, response.ResourceAlreadyExists, errMsg, []interface{}{"Plugin", "PluginID", plugin.ID}, taskInfo), "", nil
		}
		listMembers = append(listMembers, agresponse.ListMember{
			OdataID: oid,
		})
	}
	e.SubscribeToEMB(plugin.ID, queueList)

	// store encrypted password
	plugin.Password = ciphertext
	plugin.ManagerUUID = managerUUID
	// saving the pluginData
	dbErr := agmodel.SavePluginData(plugin)
	if dbErr != nil {
		errMsg := "error: while saving the plugin data: " + dbErr.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), "", nil
	}
	resp.Header = map[string]string{
		"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		"Location":     listMembers[0].OdataID,
	}
	var managersList = make([]string, 0)
	for i := 0; i < len(listMembers); i++ {
		managersList = append(managersList, listMembers[i].OdataID)
	}
	e.PublishEvent(managersList, "ManagerCollection")
	resp.StatusCode = http.StatusCreated
	log.Error("sucessfully added  plugin with the id ", cmVariants.PluginID)

	phc := agcommon.PluginHealthCheckInterface{
		DecryptPassword: common.DecryptWithPrivateKey,
	}
	phc.DupPluginConf()
	_, topics := phc.GetPluginStatus(plugin)
	PublishPluginStatusOKEvent(plugin.ID, topics)

	return resp, managerUUID, ciphertext
}
