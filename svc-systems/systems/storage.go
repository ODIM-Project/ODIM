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

	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/scommon"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"gopkg.in/go-playground/validator.v9"
)

var (
	// StringContain function pointer for the strings.Contains
	StringContain = strings.Contains
	// RequestParamsCaseValidatorFunc function pointer for the common.RequestParamsCaseValidator
	RequestParamsCaseValidatorFunc = common.RequestParamsCaseValidator
	// StringsEqualFold function pointer for the strings.EqualFold
	StringsEqualFold = strings.EqualFold
	// ContactPluginFunc  function pointer for the  scommon.ContactPlugin
	ContactPluginFunc = scommon.ContactPlugin
	// JSONUnmarshalFunc function pointer for the json.Unmarshal
	JSONUnmarshalFunc = json.Unmarshal
	// StringTrimSpace function pointer for the  strings.TrimSpace
	StringTrimSpace = strings.TrimSpace
)

// ExternalInterface holds all the external connections managers package functions uses
type ExternalInterface struct {
	ContactClient   func(context.Context, string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	DevicePassword  func([]byte) ([]byte, error)
	DB              DB
	GetPluginStatus func(context.Context, smodel.Plugin) bool
}

// DB struct to inject the contact DB function into the handlers
type DB struct {
	GetResource        func(context.Context, string, string) (string, *errors.Error)
	DeleteVolume       func(context.Context, string) *errors.Error
	AddSystemResetInfo func(context.Context, string, string) *errors.Error
	GetPluginData      func(string) (smodel.Plugin, *errors.Error)
	GetTarget          func(string) (*smodel.Target, *errors.Error)
}

// GetExternalInterface retrieves all the external connections managers package functions uses
func GetExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		ContactClient:  pmbhandle.ContactPlugin,
		DevicePassword: common.DecryptWithPrivateKey,
		DB: DB{
			GetResource:        smodel.GetResource,
			DeleteVolume:       smodel.DeleteVolume,
			AddSystemResetInfo: smodel.AddSystemResetInfo,
			GetPluginData:      smodel.GetPluginData,
			GetTarget:          smodel.GetTarget,
		},
		GetPluginStatus: scommon.GetPluginStatus,
	}
}

// CreateVolume defines the logic for creating a volume under storage
func (e *ExternalInterface) CreateVolume(ctx context.Context, req *systemsproto.VolumeRequest, pc *PluginContact, taskID string) {
	var resp response.RPC
	var targetURI = "/redfish/v1/Systems/" + req.SystemID + "/Storage/" + req.StorageInstance + "/Volumes"
	//create task
	taskInfo := &common.TaskUpdateInfo{Context: ctx, TaskID: taskID, TargetURI: targetURI,
		UpdateTask: pc.UpdateTask, TaskRequest: string(req.RequestBody)}
	// spliting the uuid and system id
	requestData := strings.SplitN(req.SystemID, ".", 2)
	if len(requestData) <= 1 {
		errorMessage := "error: SystemUUID not found"
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"System", req.SystemID}, taskInfo)
		return
	}
	uuid := requestData[0]
	target, gerr := e.DB.GetTarget(uuid)
	if gerr != nil {
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, gerr.Error(), []interface{}{"System", uuid}, taskInfo)
		return
	}
	// Validating the storage instance
	if strings.TrimSpace(req.StorageInstance) == "" {
		errorMessage := "error: Storage instance is not found"
		common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, errorMessage, []interface{}{"Storage", req.StorageInstance}, taskInfo)
		return
	}

	var volume smodel.Volume
	// unmarshalling the volume
	err := json.Unmarshal(req.RequestBody, &volume)
	if err != nil {
		errorMessage := "Error while unmarshaling the create volume request: " + err.Error()
		if StringContain(err.Error(), "smodel.OdataIDLink") {
			errorMessage = "Error processing create volume request: @odata.id key(s) is missing in Drives list"
		}
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, []interface{}{}, taskInfo)
		return
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := RequestParamsCaseValidatorFunc(req.RequestBody, volume)
	if err != nil {
		errMsg := "error while validating request parameters for volume creation: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		return
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, taskInfo)
		return
	}
	//fields validation
	statuscode, statusMessage, messageArgs, err := e.validateProperties(ctx, &volume, req.SystemID)
	if err != nil {
		errorMessage := "error: request payload validation failed: " + err.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(statuscode, statusMessage, errorMessage, messageArgs, taskInfo)
		return
	}
	decryptedPasswordByte, err := e.DevicePassword(target.Password)
	if err != nil {
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, taskInfo)
		return
	}
	target.Password = decryptedPasswordByte
	// Get the Plugin info
	plugin, gerr := e.DB.GetPluginData(target.PluginID)
	if gerr != nil {
		errorMessage := "error while trying to get plugin details"
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, taskInfo)
		return
	}
	var contactRequest scommon.PluginContactRequest
	contactRequest.ContactClient = e.ContactClient
	contactRequest.Plugin = plugin
	contactRequest.GetPluginStatus = e.GetPluginStatus

	if StringsEqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		contactRequest.HTTPMethodType = http.MethodPost
		contactRequest.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		contactRequest.OID = "/ODIM/v1/Sessions"
		_, token, _, getResponse, err := scommon.ContactPlugin(ctx, contactRequest, "error while creating session with the plugin: ")

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

	contactRequest.HTTPMethodType = http.MethodPost
	contactRequest.DeviceInfo = target
	contactRequest.OID = fmt.Sprintf("/ODIM/v1/Systems/%s/Storage/%s/Volumes", requestData[1], req.StorageInstance)

	body, location, pluginIP, getResponse, err := ContactPluginFunc(ctx, contactRequest, "error while creating a volume: ")
	if err != nil {
		resp.StatusCode = getResponse.StatusCode
		json.Unmarshal(body, &resp.Body)
		errMsg := "error while creating volume: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError,
			errMsg, nil, taskInfo)
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
	task := fillTaskData(taskID, targetURI, string(req.RequestBody), resp,
		common.Completed, common.OK, 100, http.MethodPost)
	pc.UpdateTask(ctx, task)
}

// Validates all the input prorperties
func (e *ExternalInterface) validateProperties(ctx context.Context, request *smodel.Volume, systemID string) (int32, string, []interface{}, error) {
	validate := validator.New()
	// if any of the mandatory fields missing in the struct, then it will return an error
	err := validate.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return http.StatusBadRequest, response.PropertyMissing, []interface{}{err.Field()}, fmt.Errorf(err.Field() + " field is missing")
		}
	}

	// Validates OperationApplyTime
	items := []string{"OnReset", "Immediate"}
	if request.OperationApplyTime == "" {
		request.OperationApplyTime = items[0]
	} else if found := searchItem(items, request.OperationApplyTime); !found {
		return http.StatusBadRequest, response.PropertyValueNotInList, []interface{}{request.OperationApplyTime, "OperationApplyTime"}, fmt.Errorf("OperationApplyTime %v is invalid", request.OperationApplyTime)
	}

	// If RAIDType is provided then drives length will be checked
	if request.RAIDType != "" {
		raidTypeWithMinDrives := mapRaidTypesWithMinDrives(request.RAIDType)
		// Validates the RaidType
		if raidTypeWithMinDrives == 0 {
			return http.StatusBadRequest, response.PropertyValueNotInList, []interface{}{request.RAIDType, "RAIDType"}, fmt.Errorf("RAIDType %v is invalid", request.RAIDType)
		}
		if request.Links == nil {
			return http.StatusBadRequest, response.PropertyMissing, []interface{}{"Links"}, fmt.Errorf("Links Property is not present in the request")
		}
		//validates the number of Drives
		if len(request.Links.Drives) < raidTypeWithMinDrives {
			return http.StatusBadRequest, response.PropertyMissing, []interface{}{"Drives"}, fmt.Errorf("Minimum number of Drives not matching for the RAIDType")
		}
		// Validated the contents of Drives array and even checks if the request drive exists or not
		for _, drive := range request.Links.Drives {
			driveURI := drive.OdataID
			if driveURI == "" {
				return http.StatusBadRequest, response.ResourceNotFound, []interface{}{"Drives", drive}, fmt.Errorf("Error processing create volume request: @odata.id key(s) is missing in Drives list")
			}
			_, err := e.DB.GetResource(ctx, "Drives", driveURI)
			if err != nil {
				l.LogWithFields(ctx).Error(err.Error())
				if errors.DBKeyNotFound == err.ErrNo() {
					requestData := strings.SplitN(systemID, ".", 2)
					var getDeviceInfoRequest = scommon.ResourceInfoRequest{
						URL:             driveURI,
						UUID:            requestData[0],
						SystemID:        requestData[1],
						ContactClient:   e.ContactClient,
						DevicePassword:  e.DevicePassword,
						GetPluginStatus: e.GetPluginStatus,
					}
					var err error
					if _, err = scommon.GetResourceInfoFromDevice(ctx, getDeviceInfoRequest, true); err != nil {
						return http.StatusNotFound, response.ResourceNotFound, []interface{}{"Drives", driveURI}, fmt.Errorf("Error while getting drive details for %s", driveURI)
					}
				} else {
					return http.StatusNotFound, response.ResourceNotFound, []interface{}{"Drives", driveURI}, fmt.Errorf("Error while getting drive details for %s", driveURI)
				}
			}
			// Validating if a a drive URI contains correct system id
			driveURISplit := strings.Split(driveURI, "/")
			if len(driveURISplit) > 5 && driveURISplit[4] != systemID {
				errMsg := "Drive URI contains incorrect system id"
				l.LogWithFields(ctx).Error(errMsg)
				return http.StatusBadRequest, response.ResourceNotFound, []interface{}{"Drives", drive}, fmt.Errorf(errMsg)
			}
		}
		// Validated the contents of Drives array and even checks if the request drive exists or not
		for _, drive := range request.Links.DedicatedSpareDrives {
			driveURI := drive.OdataID
			if driveURI == "" {
				return http.StatusBadRequest, response.ResourceNotFound, []interface{}{"Drives", drive}, fmt.Errorf("Error processing create volume request: @odata.id key(s) is missing in Drives list")
			}
			_, err := e.DB.GetResource(ctx, "Drives", driveURI)
			if err != nil {
				l.LogWithFields(ctx).Error(err.Error())
				if errors.DBKeyNotFound == err.ErrNo() {
					requestData := strings.SplitN(systemID, ".", 2)
					var getDeviceInfoRequest = scommon.ResourceInfoRequest{
						URL:             driveURI,
						UUID:            requestData[0],
						SystemID:        requestData[1],
						ContactClient:   e.ContactClient,
						DevicePassword:  e.DevicePassword,
						GetPluginStatus: e.GetPluginStatus,
					}
					var err error
					if _, err = scommon.GetResourceInfoFromDevice(ctx, getDeviceInfoRequest, true); err != nil {
						return http.StatusNotFound, response.ResourceNotFound, []interface{}{"Drives", driveURI}, fmt.Errorf("Error while getting drive details for %s", driveURI)
					}
				} else {
					return http.StatusNotFound, response.ResourceNotFound, []interface{}{"Drives", driveURI}, fmt.Errorf("Error while getting drive details for %s", driveURI)
				}
			}
			// Validating if a a drive URI contains correct system id
			driveURISplit := strings.Split(driveURI, "/")
			if len(driveURISplit) > 5 && driveURISplit[4] != systemID {
				errMsg := "Drive URI contains incorrect system id"
				l.LogWithFields(ctx).Error(errMsg)
				return http.StatusBadRequest, response.ResourceNotFound, []interface{}{"Drives", drive}, fmt.Errorf(errMsg)
			}
		}
	}
	// validate WriteCachePolicy
	writeCachePolicy := map[string]bool{"WriteThrough": true, "ProtectedWriteBack": true, "UnprotectedWriteBack": true, "Off": true}

	if request.WriteCachePolicy != "" {
		_, isExists := writeCachePolicy[request.WriteCachePolicy]
		if !isExists {
			return http.StatusBadRequest, response.PropertyValueNotInList, []interface{}{request.WriteCachePolicy, "WriteCachePolicy"}, fmt.Errorf("WriteCachePolicy %v is invalid", request.WriteCachePolicy)

		}
	}
	//validate ReadCachePolicy
	readCachePolicy := map[string]bool{"ReadAhead": true, "AdaptiveReadAhead": true, "Off": true}
	if request.ReadCachePolicy != "" {
		_, isExists := readCachePolicy[request.ReadCachePolicy]
		if !isExists {
			return http.StatusBadRequest, response.PropertyValueNotInList, []interface{}{request.ReadCachePolicy, "ReadCachePolicy"}, fmt.Errorf("ReadCachePolicy %v is invalid", request.ReadCachePolicy)
		}
	}

	return http.StatusOK, common.OK, []interface{}{}, nil
}

// Mapping the raid types with minimum number of drives
func mapRaidTypesWithMinDrives(req string) int {
	raidTypesWithMinDrives := map[string]int{
		"RAID0":        1,
		"RAID00":       2,
		"RAID01":       2,
		"RAID1":        2,
		"RAID10":       4,
		"RAID10E":      2,
		"RAID10Triple": 6,
		"RAID1E":       2,
		"RAID1Triple":  3,
		"RAID3":        3,
		"RAID4":        3,
		"RAID5":        3,
		"RAID50":       6,
		"RAID6":        4,
		"RAID60":       8,
		"RAID6TP":      4,
	}
	return raidTypesWithMinDrives[req]
}

// searchItem is used to find an item from the slice
func searchItem(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// DeleteVolume defines the logic for deleting a volume under storage
func (e *ExternalInterface) DeleteVolume(ctx context.Context, req *systemsproto.VolumeRequest, pc *PluginContact, taskID string) {
	var resp response.RPC
	var targetURI = "/redfish/v1/Systems/" + req.SystemID + "/Storage/" + req.StorageInstance + "/Volumes" + req.VolumeID

	taskInfo := &common.TaskUpdateInfo{Context: ctx, TaskID: taskID, TargetURI: targetURI,
		UpdateTask: pc.UpdateTask, TaskRequest: string(req.RequestBody)}
	var volume smodel.Volume
	err := JSONUnmarshalFunc(req.RequestBody, &volume)
	if err != nil {
		errorMessage := "Error while unmarshaling the create volume request: " + err.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, []interface{}{}, taskInfo)
		return
	}

	// spliting the uuid and system id
	requestData := strings.SplitN(req.SystemID, ".", 2)
	if len(requestData) != 2 || requestData[1] == "" {
		errorMessage := "error: SystemUUID not found"
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"System", req.SystemID}, taskInfo)
		return
	}
	uuid := requestData[0]
	target, gerr := e.DB.GetTarget(uuid)
	if gerr != nil {
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, gerr.Error(), []interface{}{"System", uuid}, taskInfo)
		return
	}
	// Validating the storage instance
	if StringTrimSpace(req.VolumeID) == "" {
		errorMessage := "error: Volume id is not found"
		common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, errorMessage, []interface{}{"Volume", req.VolumeID}, taskInfo)
		return
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := RequestParamsCaseValidatorFunc(req.RequestBody, volume)
	if err != nil {
		errMsg := "error while validating request parameters for volume creation: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		return
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, taskInfo)
		return
	}

	key := fmt.Sprintf("/redfish/v1/Systems/%s/Storage/%s/Volumes/%s", req.SystemID, req.StorageInstance, req.VolumeID)
	_, dbErr := e.DB.GetResource(ctx, "Volumes", key)
	if dbErr != nil {
		l.LogWithFields(ctx).Error("error getting volumes details : " + dbErr.Error())
		errorMessage := dbErr.Error()
		if errors.DBKeyNotFound == dbErr.ErrNo() {
			var getDeviceInfoRequest = scommon.ResourceInfoRequest{
				URL:             key,
				UUID:            uuid,
				SystemID:        requestData[1],
				ContactClient:   e.ContactClient,
				DevicePassword:  e.DevicePassword,
				GetPluginStatus: e.GetPluginStatus,
			}
			var err error
			if _, err = scommon.GetResourceInfoFromDevice(ctx, getDeviceInfoRequest, true); err != nil {
				common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"Volumes", key}, taskInfo)
				return
			}

		} else {
			common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, taskInfo)
			return
		}
	}

	decryptedPasswordByte, err := e.DevicePassword(target.Password)
	if err != nil {
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, taskInfo)
		return
	}
	target.Password = decryptedPasswordByte
	// Get the Plugin info
	plugin, gerr := e.DB.GetPluginData(target.PluginID)
	if gerr != nil {
		errorMessage := "error while trying to get plugin details"
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, taskInfo)
		return
	}

	var contactRequest scommon.PluginContactRequest
	contactRequest.ContactClient = e.ContactClient
	contactRequest.Plugin = plugin
	contactRequest.GetPluginStatus = e.GetPluginStatus
	if StringsEqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		contactRequest.HTTPMethodType = http.MethodPost
		contactRequest.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		contactRequest.OID = "/ODIM/v1/Sessions"
		_, token, _, getResponse, err := scommon.ContactPlugin(ctx, contactRequest, "error while creating session with the plugin: ")

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

	if string(req.RequestBody) == "null" {
		target.PostBody = []byte{}
	} else {
		target.PostBody = req.RequestBody

	}
	contactRequest.HTTPMethodType = http.MethodDelete
	contactRequest.DeviceInfo = target
	contactRequest.OID = fmt.Sprintf("/ODIM/v1/Systems/%s/Storage/%s/Volumes/%s", requestData[1], req.StorageInstance, req.VolumeID)

	body, location, pluginIP, getResponse, err := scommon.ContactPlugin(ctx, contactRequest, "error while deleting a volume: ")
	if err != nil {
		resp.StatusCode = getResponse.StatusCode
		json.Unmarshal(body, &resp.Body)
		errMsg := "error while deleting volume: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError,
			errMsg, nil, taskInfo)
		return
	}
	if getResponse.StatusCode == http.StatusAccepted {
		scommon.SavePluginTaskInfo(ctx, pluginIP, plugin.IP, taskID, location)
		return
	}

	// delete a volume in db
	if derr := e.DB.DeleteVolume(ctx, key); derr != nil {
		errMsg := "error while trying to delete volume: " + derr.Error()
		l.LogWithFields(ctx).Error(errMsg)
		if errors.DBKeyNotFound == derr.ErrNo() {
			common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"Volumes", key}, taskInfo)
			return
		}
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		return
	}

	// adding volume collection uri and deleted volume uri to the AddSystemResetInfo
	// for avoiding storing or retrieving them from DB before a BMC reset.
	collectionKey := fmt.Sprintf("/redfish/v1/Systems/%s/Storage/%s/Volumes", req.SystemID, req.StorageInstance)
	e.DB.AddSystemResetInfo(ctx, key, "On")
	e.DB.AddSystemResetInfo(ctx, collectionKey, "On")

	resp.StatusCode = http.StatusNoContent
	resp.StatusMessage = response.Success
	task := fillTaskData(taskID, targetURI, string(req.RequestBody), resp,
		common.Completed, common.OK, 100, http.MethodPost)
	pc.UpdateTask(ctx, task)
}
