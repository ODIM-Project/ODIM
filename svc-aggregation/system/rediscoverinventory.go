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
	"strconv"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agcommon"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/google/uuid"
)

const (
	RediscoverResourcesActionID = "217"

	RediscoverResourcesActionName = "RediscoverResources"
)

// RediscoverSystemInventory  is the handler for redicovering system whenever the restrat event detected in event service
//It deletes old data and  Discovers Computersystem & Chassis and its top level odata.ID links and store them in inmemory db.
func (e *ExternalInterface) RediscoverSystemInventory(ctx context.Context, deviceUUID, systemURL string, updateFlag bool) {
	l.LogWithFields(ctx).Info("Rediscovery of the BMC with ID " + deviceUUID + " is started.")

	var resp response.RPC
	systemURL = strings.TrimSuffix(systemURL, "/")
	data := strings.Split(systemURL, "/")
	// Getting the SystemID from system url
	if len(data) <= 0 {
		genError(ctx, "invalid data ", &resp, http.StatusInternalServerError, errors.InternalError, map[string]string{
			"Content-type": "application/json; charset=utf-8",
		})
		return
	}

	// Getting the device info
	target, err := agmodel.GetTarget(deviceUUID)
	if err != nil {
		genError(ctx, err.Error(), &resp, http.StatusBadRequest, errors.ResourceNotFound, map[string]string{
			"Content-type": "application/json; charset=utf-8",
		})
		l.LogWithFields(ctx).Error("Unable to unmarshal data: " + err.Error())
		return
	}
	decryptedPasswordByte, err := e.DecryptPassword(target.Password)
	if err != nil {
		genError(ctx, "error while trying to decrypt device password: "+err.Error(), &resp, http.StatusInternalServerError, errors.InternalError, map[string]string{
			"Content-type": "application/json; charset=utf-8",
		})
		l.LogWithFields(ctx).Error("Unable to unmarshal data: " + err.Error())
		return
	}
	target.Password = decryptedPasswordByte

	// get the plugin information
	plugin, errs := agmodel.GetPluginData(target.PluginID)
	if errs != nil {
		genError(ctx, errs.Error(), &resp, http.StatusBadRequest, errors.ResourceNotFound, map[string]string{
			"Content-type": "application/json; charset=utf-8",
		})
		l.LogWithFields(ctx).Error(errs.Error())
		return
	}

	var req getResourceRequest
	req.ContactClient = e.ContactClient
	req.GetPluginStatus = e.GetPluginStatus
	req.Plugin = plugin
	req.StatusPoll = true
	req.BMCAddress = target.ManagerAddress
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		req.HTTPMethodType = http.MethodPost
		req.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		req.OID = "/ODIM/v1/Sessions"
		_, token, _, err := contactPlugin(ctx, req, "error while getting the details "+req.OID+": ")
		if err != nil {
			l.LogWithFields(ctx).Error(err.Error())
			return
		}
		req.Token = token
	} else {
		req.LoginCredentials = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	// check whether delete operation for the system is initiated
	if strings.Contains(systemURL, "/Storage") {
		systemURL = strings.Replace(systemURL, "/Storage", "", -1)
	}
	systemOperation, dbErr := agmodel.GetSystemOperationInfo(systemURL)
	if dbErr != nil && errors.DBKeyNotFound != dbErr.ErrNo() {
		l.LogWithFields(ctx).Error("Rediscovery for system: " + systemURL + " can't be processed " + dbErr.Error())
		return
	}
	if systemOperation.Operation == "Delete" {
		l.LogWithFields(ctx).Error("Rediscovery for system: " + systemURL + " can't be processed," +
			systemOperation.Operation + " operation is under progress")
		return
	}

	// Add system operation info to db to block the  delete  request for respective system
	systemOperation.Operation = "InventoryRediscovery"
	dbErr = systemOperation.AddSystemOperationInfo(systemURL)
	if dbErr != nil {
		l.LogWithFields(ctx).Error("Rediscovery for system: " + systemURL + " can't be processed " + dbErr.Error())
		return
	}
	defer func() {
		agmodel.DeleteSystemOperationInfo(systemURL)
		agmodel.DeleteSystemResetInfo(systemURL)
		deleteResourceResetInfo(ctx, systemURL)
	}()

	deleteSubordinateResource(ctx, deviceUUID)

	req.DeviceUUID = deviceUUID
	req.DeviceInfo = target
	req.OID = strings.Replace(systemURL, "/redfish/v1/Systems/"+deviceUUID+".", "/redfish/v1/Systems/", -1)
	l.LogWithFields(ctx).Info("Request oid for rediscovery," + req.OID)
	req.UpdateFlag = updateFlag
	req.UpdateTask = e.UpdateTask
	var h respHolder
	h.TraversedLinks = make(map[string]bool)
	h.InventoryData = make(map[string]interface{})
	progress := int32(100)
	systemsEstimatedWork := int32(75)
	if strings.Contains(systemURL, "/Storage") {
		_, progress, _ = h.getStorageInfo(ctx, progress, systemsEstimatedWork, req)
	} else {
		_, _, progress, _ = h.getSystemInfo(ctx, "", progress, systemsEstimatedWork, req)
		h.InventoryData = make(map[string]interface{})
		//rediscovering the Chassis Information
		req.OID = "/redfish/v1/Chassis"
		chassisEstimatedWork := int32(15)
		progress = h.getAllRootInfo(ctx, "", progress, chassisEstimatedWork, req, config.Data.AddComputeSkipResources.SkipResourceListUnderChassis)

		//rediscovering the Manager Information
		req.OID = "/redfish/v1/Managers"
		managerEstimatedWork := int32(15)
		progress = h.getAllRootInfo(ctx, "", progress, managerEstimatedWork, req, config.Data.AddComputeSkipResources.SkipResourceListUnderManager)
		agmodel.SaveBMCInventory(h.InventoryData)
	}

	var responseBody = map[string]string{
		"UUID": deviceUUID,
	}

	resp.StatusCode = http.StatusCreated
	resp.Body = responseBody

	l.LogWithFields(ctx).Info("Rediscovery of the BMC with ID " + deviceUUID + " is now complete.")
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
	transactionID := uuid.New()
	ctx := agcommon.CreateContext(transactionID.String(), RediscoverResourcesActionID, RediscoverResourcesActionName, "1", common.AggregationService, podName)
	targets, err := agmodel.GetAllSystems()
	if err != nil || len(targets) == 0 {
		// nothing to re-discover
		l.LogWithFields(ctx).Info("Nothing to re-discover.")
		return nil
	}

	serverBatchSize := config.Data.ServerRediscoveryBatchSize
	if config.Data.ServerRediscoveryBatchSize <= 0 {
		serverBatchSize = 1
	}
	threadID := 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.RediscoverSystemInventory)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	threadID++
	var semaphoreChan = make(chan int, serverBatchSize)
	for index := range targets {
		semaphoreChan <- 1
		go func(ctxt context.Context, target agmodel.Target) {
			defer func() {
				<-semaphoreChan
			}()
			// Call the plugin to get the systems collection for this target first
			systemCollectionResponse, err := e.getTargetSystemCollection(ctxt, target)
			if err != nil {
				l.LogWithFields(ctxt).Error("Failed to discover the server: " + err.Error())
				return
			}
			systemsCollection := make(map[string]interface{})
			err = json.Unmarshal(systemCollectionResponse, &systemsCollection)
			if err != nil {
				l.LogWithFields(ctxt).Error("Failed to discover the server: " + err.Error())
				return
			}
			members := systemsCollection["Members"]
			var systemURLArray []string
			for _, member := range members.([]interface{}) {
				systemURL := member.(map[string]interface{})["@odata.id"].(string)
				if e.isServerRediscoveryRequired(ctxt, target.DeviceUUID, systemURL) == true {
					e.RediscoverSystemInventory(ctxt, target.DeviceUUID, systemURL, true)
					systemURLArray = append(systemURLArray, systemURL)
				}
			}
			e.publishResourceUpdatedEvent(ctxt, systemURLArray, "SystemsCollection")
		}(ctxt, targets[index])
	}
	// if everything is OK return success
	return nil

}
func (e *ExternalInterface) getTargetSystemCollection(ctx context.Context, target agmodel.Target) ([]byte, error) {

	decryptedPasswordByte, err := e.DecryptPassword(target.Password)
	if err != nil {
		return nil, err
	}
	target.Password = decryptedPasswordByte
	// get the plugin information
	plugin, errs := agmodel.GetPluginData(target.PluginID)
	if errs != nil {
		l.LogWithFields(ctx).Error(errs.Error())
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
		_, token, _, err := contactPlugin(ctx, req, "error while getting the details "+req.OID+": ")
		if err != nil {
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
	body, _, _, err := contactPlugin(ctx, req, "error while trying to get the system collection details: ")
	if err != nil {
		return nil, fmt.Errorf("error while trying to get the system collection details")
	}
	return body, nil
}

func (e *ExternalInterface) isServerRediscoveryRequired(ctx context.Context, deviceUUID string, systemKey string) bool {
	systemKey = strings.TrimSuffix(systemKey, "/")
	key := strings.Replace(systemKey, "/redfish/v1/Systems/", "/redfish/v1/Systems/"+deviceUUID+".", -1)
	_, err := agmodel.GetResource("ComputerSystem", key)
	if err != nil {
		l.LogWithFields(ctx).Error(err.Error())
		l.LogWithFields(ctx).Info("Rediscovery required for the server with UUID: " + deviceUUID)
		return true

	}

	key = strings.Replace(systemKey, "Systems", "Chassis", -1)
	keys, err := agmodel.GetAllMatchingDetails("Chassis", key, common.InMemory)
	if err != nil || len(keys) == 0 {
		l.LogWithFields(ctx).Info("Rediscovery required for the server with UUID: " + deviceUUID)
		return true
	}
	for _, chassiskey := range keys {
		if _, err = agmodel.GetResource("Chassis", chassiskey); err != nil {
			l.LogWithFields(ctx).Error(err.Error())
			l.LogWithFields(ctx).Info("Rediscovery required for the server with UUID: " + deviceUUID)
			return true
		}
	}

	key = strings.Replace(systemKey, "Systems", "Managers", -1)
	keys, err = agmodel.GetAllMatchingDetails("Managers", key, common.InMemory)
	if err != nil || len(keys) == 0 {
		l.LogWithFields(ctx).Info("Rediscovery required for the server with UUID: " + deviceUUID)
		return true
	}
	for _, managerKey := range keys {
		if _, err = agmodel.GetResource("Managers", managerKey); err != nil {
			l.LogWithFields(ctx).Error(err.Error())
			l.LogWithFields(ctx).Info("Rediscovery required for the server with UUID: " + deviceUUID)
			return true
		}
	}
	l.LogWithFields(ctx).Info("Rediscovery not required for the server with UUID: " + deviceUUID)
	return false
}

// publishResourceUpdatedEvent will publish ResourceUpdated events
func (e *ExternalInterface) publishResourceUpdatedEvent(ctx context.Context, systemIDs []string, collectionName string) {
	for i := 0; i < len(systemIDs); i++ {
		e.PublishEventMB(ctx, systemIDs[i], "ResourceUpdated", collectionName)
	}
}

func deleteResourceResetInfo(ctx context.Context, pattern string) {
	keys, err := agmodel.GetAllMatchingDetails("SystemReset", pattern, common.InMemory)
	if err != nil {
		l.LogWithFields(ctx).Error("Unable to fetch all matching keys from system reset table: " + err.Error())
	}
	for _, key := range keys {
		agmodel.DeleteSystemResetInfo(key)
	}
}

// deleteSubordinateResource will delete all the subordinate resources assosiated with the pattern
func deleteSubordinateResource(ctx context.Context, deviceUUID string) {
	l.LogWithFields(ctx).Info("Initiated removal of subordinate resource for the BMC with ID " +
		deviceUUID + " from the in-memory DB")
	keys, err := agmodel.GetAllMatchingDetails("*", deviceUUID, common.InMemory)
	if err != nil {
		l.LogWithFields(ctx).Error("Unable to fetch all matching keys from system reset table: " + err.Error())
		return
	}
	for _, key := range keys {
		resourceDetails := strings.Split(key, ":")
		switch resourceDetails[0] {
		case "ComputerSystem", "SystemReset", "SystemOperation", "Chassis", "Managers", "FirmwareInventory", "SoftwareInventory":
			continue
		default:
			if err = agmodel.Delete(resourceDetails[0], resourceDetails[1], common.InMemory); err != nil {
				l.LogWithFields(ctx).Error("Delete of " + resourceDetails[1] + " from " + resourceDetails[0] + " in " +
					string(common.InMemory) + " DB failed due to the error: " + err.Error())
			}
		}
	}
	l.LogWithFields(ctx).Info("Removal of subordinate resources for the BMC with ID " + deviceUUID + " from the in-memory DB is now complete.")
}
