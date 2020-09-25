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
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

// RediscoverSystemInventory  is the handler for redicovering system whenever the restrat event detected in event service
//It deletes old data and  Discovers Computersystem & Chassis and its top level odata.ID links and store them in inmemory db.
func (e *ExternalInterface) RediscoverSystemInventory(deviceUUID, systemURL string, updateFlag bool) {
	var resp response.RPC
	systemURL = strings.TrimSuffix(systemURL, "/")
	data := strings.Split(systemURL, "/")
	// Getting the SystemID from system url
	if len(data) <= 0 {
		genError("invalid data ", &resp, http.StatusInternalServerError, errors.InternalError, map[string]string{
			"Content-type": "application/json; charset=utf-8",
		})
		return
	}
	sysID := data[len(data)-1]
	log.Println("System URl: ", systemURL, " sysID: ", sysID, " deviceUUID: ", deviceUUID)

	// Getting the device info
	target, err := agmodel.GetTarget(deviceUUID)
	if err != nil {
		genError(err.Error(), &resp, http.StatusBadRequest, errors.ResourceNotFound, map[string]string{
			"Content-type": "application/json; charset=utf-8",
		})
		log.Println("error while unmarshal data: ", err)
		return

	}
	decryptedPasswordByte, err := e.DecryptPassword(target.Password)
	if err != nil {
		genError("error while trying to decrypt device password: "+err.Error(), &resp, http.StatusInternalServerError, errors.InternalError, map[string]string{
			"Content-type": "application/json; charset=utf-8",
		})
		log.Println("error while unmarshal data: ", err)
		return
	}
	target.Password = decryptedPasswordByte

	// get the plugin information
	plugin, errs := agmodel.GetPluginData(target.PluginID)
	if errs != nil {
		genError(errs.Error(), &resp, http.StatusBadRequest, errors.ResourceNotFound, map[string]string{
			"Content-type": "application/json; charset=utf-8",
		})
		log.Printf(errs.Error())
		return
	}

	var req getResourceRequest
	req.ContactClient = e.ContactClient
	req.GetPluginStatus = e.GetPluginStatus
	req.Plugin = plugin
	req.StatusPoll = true
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		req.HTTPMethodType = http.MethodPost
		req.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		req.OID = "/ODIM/v1/Sessions"
		_, token, _, err := contactPlugin(req, "error while getting the details "+req.OID+": ")
		if err != nil {
			log.Println(err)
			return
		}
		req.Token = token
	} else {
		req.LoginCredentials = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	// check whether delete operation for the system is intiated
	udaptedSystemURI := strings.Replace(systemURL, "/redfish/v1/Systems/", "/redfish/v1/Systems/"+deviceUUID+":", -1)
	if strings.Contains(systemURL, "/Storage") {
		udaptedSystemURI = strings.Replace(udaptedSystemURI, "/Storage", "", -1)
	}
	systemOperation, dbErr := agmodel.GetSystemOperationInfo(udaptedSystemURI)
	if dbErr != nil && errors.DBKeyNotFound != dbErr.ErrNo() {
		log.Println("Rediscovery for system: ", udaptedSystemURI, " can't be processed ", dbErr.Error())
		return
	}
	if systemOperation.Operation == "Delete" {
		log.Println("Rediscovery for system: ", udaptedSystemURI, " can't be processed,", systemOperation.Operation, " operation is under progress")
		return
	}

	// Add system operation info to db to block the  delete  request for respective system
	systemOperation.Operation = "InventoryRediscovery"
	dbErr = systemOperation.AddSystemOperationInfo(udaptedSystemURI)
	if dbErr != nil {
		log.Println("Rediscovery for system: ", udaptedSystemURI, " can't be processed ", dbErr.Error())
		return
	}
	defer func() {
		agmodel.DeleteSystemOperationInfo(udaptedSystemURI)
		agmodel.DeleteSystemResetInfo(udaptedSystemURI)
		deleteResourceResetInfo(udaptedSystemURI)
	}()
	req.DeviceUUID = deviceUUID
	req.DeviceInfo = target
	req.OID = systemURL
	req.UpdateFlag = updateFlag
	req.UpdateTask = e.UpdateTask
	var h respHolder
	h.TraversedLinks = make(map[string]bool)
	progress := int32(100)
	systemsEstimatedWork := int32(75)
	if strings.Contains(systemURL, "/Storage") {
		_, progress, _ = h.getStorageInfo(progress, systemsEstimatedWork, req)
	} else {
		_, progress, _ = h.getSystemInfo("", progress, systemsEstimatedWork, req)
		//rediscovering the Chassis Information
		req.OID = "/redfish/v1/Chassis"
		chassisEstimatedWork := int32(15)
		progress = h.getAllRootInfo("", progress, chassisEstimatedWork, req)
	}

	var responseBody = map[string]string{
		"UUID": deviceUUID,
	}

	resp.StatusCode = http.StatusCreated
	resp.Body = responseBody
	resp.Header = map[string]string{
		"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
	}
	log.Println("info: rediscovery of system and chassis completed.")

	return
}

//RediscoverResources is a function to rediscover the server inventory,
// in the event of InMemory DB crashed and/or rebooted all of the content/inventory
// in the Inmemory DB is gone. So to repopulate the inventory of all the added server,
// this function can be used.
//Takes: None
//Returns: error
// On success nil
//On Failure Non nil
func (e *ExternalInterface) RediscoverResources() error {
	// First check if the redicovery requires.
	// InMemory DB is just fine most of the times.
	// Try to get all the systems from InMemory DB, if the collection is not empty
	// then InMemory DB is just fine, so no need for resource inventory rediscovery

	// Get the resources from OnDisk DB

	targets, err := agmodel.GetAllSystems()
	if err != nil || len(targets) == 0 {
		// nothing to re-discover
		log.Printf("Nothing to re-discover.")
		return nil
	}

	serverBatchSize := config.Data.ServerRediscoveryBatchSize
	if config.Data.ServerRediscoveryBatchSize <= 0 {
		serverBatchSize = 1
	}
	var semaphoreChan = make(chan int, serverBatchSize)
	for index := range targets {
		semaphoreChan <- 1
		go func(target agmodel.Target) {
			defer func() {
				<-semaphoreChan
			}()
			// Call the plugin to get the systems collection for this target first
			systemCollectionResponse, err := e.getTargetSystemCollection(target)
			if err != nil {
				log.Printf("error: failed to discover the server: %s ", err.Error())
				return
			}
			systemsCollection := make(map[string]interface{})
			err = json.Unmarshal(systemCollectionResponse, &systemsCollection)
			if err != nil {
				log.Printf("error: failed to discover the server: %s ", err.Error())
				return
			}
			members := systemsCollection["Members"]
			var systemURLArray []string
			for _, member := range members.([]interface{}) {
				systemURL := member.(map[string]interface{})["@odata.id"].(string)
				if e.isServerRediscoveryRequired(target.DeviceUUID, systemURL) == true {
					e.RediscoverSystemInventory(target.DeviceUUID, systemURL, false)
					systemURLArray = append(systemURLArray, systemURL)
				}
			}
			e.publishResourceUpdatedEvent(systemURLArray, "systemsCollection")
		}(targets[index])
	}
	// if everything is OK return success
	return nil

}
func (e *ExternalInterface) getTargetSystemCollection(target agmodel.Target) ([]byte, error) {

	decryptedPasswordByte, err := e.DecryptPassword(target.Password)
	if err != nil {
		return nil, err
	}
	target.Password = decryptedPasswordByte
	// get the plugin information
	plugin, errs := agmodel.GetPluginData(target.PluginID)
	if errs != nil {
		log.Printf(errs.Error())
		return nil, errs
	}

	var req getResourceRequest
	req.ContactClient = e.ContactClient
	req.GetPluginStatus = e.GetPluginStatus
	req.Plugin = plugin
	req.StatusPoll = true
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		req.HTTPMethodType = http.MethodPost
		req.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		req.OID = "/ODIM/v1/Sessions"
		_, token, _, err := contactPlugin(req, "error while getting the details "+req.OID+": ")
		if err != nil {
			log.Println(err)
			return nil, err
		}
		req.Token = token
	} else {
		req.LoginCredentials = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}

	req.DeviceUUID = target.DeviceUUID
	req.DeviceInfo = target
	req.OID = "/redfish/v1/Systems"

	// Make the call to Plugin with above request
	body, _, _, err := contactPlugin(req, "error while trying to get the system collection details: ")
	if err != nil {
		return nil, fmt.Errorf("error while trying to get the system collection details")
	}
	return body, nil
}

func (e *ExternalInterface) isServerRediscoveryRequired(deviceUUID string, systemKey string) bool {
	strArray := strings.Split(systemKey, "/")
	sysID := strArray[len(strArray)-1]
	systemKey = strings.Replace(systemKey, "/"+sysID, "/"+deviceUUID+":"+sysID, -1)
	key := systemKey
	_, err := agmodel.GetResource("ComputerSystem", key)
	if err != nil {
		log.Printf(err.Error())
		log.Printf("Rediscovery required for the server with UUID: %s.", deviceUUID)
		return true

	}
	key = strings.Replace(systemKey, "Systems", "Chassis", -1)
	_, err = agmodel.GetResource("Chassis", key)
	if err != nil {
		log.Printf(err.Error())
		log.Printf("Rediscovery required for the server with UUID: %s.", deviceUUID)
		return true
	}
	key = strings.Replace(systemKey, "Systems", "Managers", -1)
	_, err = agmodel.GetResource("Managers", key)
	if err != nil {
		log.Printf(err.Error())
		log.Printf("Rediscovery required for the server with UUID: %s.", deviceUUID)
		return true
	}
	log.Printf("Rediscovery not required for the server with UUID: %s.", deviceUUID)
	return false
}

// publishResourceUpdatedEvent will publish ResourceUpdated events
func (e *ExternalInterface) publishResourceUpdatedEvent(systemIDs []string, collectionName string) {
	for i := 0; i < len(systemIDs); i++ {
		e.PublishEventMB(systemIDs[i], "ResourceUpdated", collectionName)
	}
}

func deleteResourceResetInfo(pattern string) {
	keys, err := agmodel.GetAllMatchingDetails("SystemReset", pattern, common.InMemory)
	if err != nil {
		log.Printf("error while trying to fetch all matching keys from system reset table: %v", err)
	}
	for _, key := range keys {
		agmodel.DeleteSystemResetInfo(key)
	}
}
