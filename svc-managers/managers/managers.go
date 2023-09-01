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

// Package managers ...
package managers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrcommon"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrmodel"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrresponse"
	"gopkg.in/go-playground/validator.v9"
)

var (
	jsonUnMarshalFunc              = json.Unmarshal
	requestParamsCaseValidatorFunc = common.RequestParamsCaseValidator
)

// GetManagersCollection will get the all the managers(odimra, Plugins, Servers)
func (e *ExternalInterface) GetManagersCollection(ctx context.Context, req *managersproto.ManagerRequest) (response.RPC, error) {
	var resp response.RPC
	managers := mgrresponse.ManagersCollection{
		OdataContext: "/redfish/v1/$metadata#ManagerCollection.ManagerCollection",
		OdataID:      "/redfish/v1/Managers",
		OdataType:    "#ManagerCollection.ManagerCollection",
		Description:  "Managers view",
		Name:         "Managers",
	}
	var members []dmtf.Link

	// Add servers as manager in manager collection
	managersCollectionKeysArray, err := e.DB.GetAllKeysFromTable("Managers")
	if err != nil || len(managersCollectionKeysArray) == 0 {
		l.LogWithFields(ctx).Error("No servers found in odimra")
	}

	for _, key := range managersCollectionKeysArray {
		members = append(members, dmtf.Link{Oid: key})
	}
	managers.Members = members
	managers.MembersCount = len(members)
	resp.Body = managers
	resp.StatusCode = http.StatusOK
	respBody := fmt.Sprintf("%v", resp.Body)
	l.LogWithFields(ctx).Debugf("Outgoing manager collection details response to northbound: %s", string(respBody))
	return resp, nil
}

// GetManagers will fetch individual manager details with the given ID
func (e *ExternalInterface) GetManagers(ctx context.Context, req *managersproto.ManagerRequest) response.RPC {
	var resp response.RPC
	if req.ManagerID == config.Data.RootServiceUUID {
		manager, err := e.getManagerDetails(ctx, req.ManagerID)
		if err != nil {
			l.LogWithFields(ctx).Error("error getting manager details : " + err.Error())
			errArgs := []interface{}{"Managers", req.ManagerID}
			errorMessage := err.Error()
			resp = common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage,
				errArgs, nil)
			return resp
		}
		resp.Body = manager
	} else {
		requestData := strings.SplitN(req.ManagerID, ".", 2)
		if len(requestData) <= 1 {
			resp = e.getPluginManagerResoure(ctx, requestData[0], req.URL)
			return resp
		}
		uuid := requestData[0]
		data, err := e.DB.GetManagerByURL(req.URL)
		if err != nil {
			l.LogWithFields(ctx).Error("error getting manager details : " + err.Error())
			var errArgs = []interface{}{}
			var statusCode int
			var StatusMessage string
			errorMessage := err.Error()
			if errors.DBKeyNotFound == err.ErrNo() {
				errArgs = []interface{}{"Managers", req.ManagerID}

				statusCode = http.StatusNotFound
				StatusMessage = response.ResourceNotFound
			} else {
				statusCode = http.StatusInternalServerError
				StatusMessage = response.InternalError
			}
			resp = common.GeneralError(int32(statusCode), StatusMessage, errorMessage,
				errArgs, nil)
			return resp
		}
		var managerData map[string]interface{}
		jerr := json.Unmarshal([]byte(data), &managerData)
		if jerr != nil {
			errorMessage := "error unmarshalling manager details: " + jerr.Error()
			l.LogWithFields(ctx).Error(errorMessage)
			resp = common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
				nil, nil)
			return resp
		}
		// extracting the Manager Type from the  managerData
		var managerType string
		if val, ok := managerData["ManagerType"]; ok {
			managerType = val.(string)
		}
		//adding default description
		if _, ok := managerData["Description"]; !ok {
			managerData["Description"] = "BMC Manager"
		}
		//adding RemoteAccountService object to manager response
		if _, ok := managerData["RemoteAccountService"]; !ok {
			managerData["RemoteAccountService"] = map[string]string{
				"@odata.id": "/redfish/v1/Managers/" + req.ManagerID + "/RemoteAccountService",
			}
		}
		//adding PowerState
		if _, ok := managerData["PowerState"]; !ok {
			managerData["PowerState"] = "On"
		}
		if managerType != common.ManagerTypeService && managerType != "" {
			deviceData, err := e.getResourceInfoFromDevice(ctx, req.URL, uuid, requestData[1], nil)
			if err != nil {
				l.LogWithFields(ctx).Error("Device " + req.URL + " is unreachable: " + err.Error())
				// Updating the state
				managerData["Status"] = map[string]string{
					"State": "Absent",
				}
			} else {
				jerr := json.Unmarshal([]byte(deviceData), &managerData)
				if jerr != nil {
					errorMessage := "error unmarshaling manager details: " + jerr.Error()
					l.LogWithFields(ctx).Error(errorMessage)
					resp = common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
						nil, nil)
					return resp
				}
			}
			managerData["Id"] = req.ManagerID
			err = e.DB.UpdateData(req.URL, managerData, "Managers")
			if err != nil {
				errorMessage := "error while saving manager details: " + err.Error()
				l.LogWithFields(ctx).Error(errorMessage)
				resp = common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
					nil, nil)
				return resp
			}
			dataBytes, err := json.Marshal(managerData)
			if err != nil {
				errorMessage := "error while marshalling manager details: " + err.Error()
				l.LogWithFields(ctx).Error(errorMessage)
				resp = common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
					nil, nil)
				return resp
			}
			data = string(dataBytes)
		}
		var resource map[string]interface{}
		json.Unmarshal([]byte(data), &resource)
		resp.Body = resource
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	respBody := fmt.Sprintf("%v", resp.Body)
	l.LogWithFields(ctx).Debugf("Outgoing manager details response to northbound: %s", string(respBody))
	return resp
}

func (e *ExternalInterface) getManagerDetails(ctx context.Context, id string) (mgrmodel.Manager, error) {
	var mgr mgrmodel.Manager
	var mgrData mgrmodel.RAManager
	data, err := e.DB.GetManagerByURL("/redfish/v1/Managers/" + id)
	if err != nil {
		return mgr, fmt.Errorf("unable to retrieve manager information: %v", err)
	}

	if err := jsonUnMarshalFunc([]byte(data), &mgrData); err != nil {
		return mgr, fmt.Errorf("unable to marshal manager information: %v", err)
	}

	chassisList, chassisErr := e.DB.GetAllKeysFromTable("Chassis")
	if chassisErr != nil {
		return mgr, fmt.Errorf("unable to retrieve chassis list information: %v", chassisErr)
	}

	serverList, serverErr := e.DB.GetAllKeysFromTable("ComputerSystem")
	if serverErr != nil {
		return mgr, fmt.Errorf("unable to retrieve server list information: %v", serverErr)
	}
	managerList, mgrErr := e.DB.GetAllKeysFromTable("Managers")
	if mgrErr != nil {
		return mgr, fmt.Errorf("unable to retrieve manager list information: %v", mgrErr)
	}
	var chassisLink, serverLink, managerLink []*dmtf.Link
	if len(chassisList) > 0 {
		for _, key := range chassisList {
			chassisLink = append(chassisLink, &dmtf.Link{Oid: key})
		}
	}
	if len(serverList) > 0 {
		for _, key := range serverList {
			serverLink = append(serverLink, &dmtf.Link{Oid: key})
		}
	}
	odimURI := "/redfish/v1/Managers/" + config.Data.RootServiceUUID
	if len(managerList) > 0 {
		for _, key := range managerList {
			if key != odimURI {
				managerLink = append(managerLink, &dmtf.Link{Oid: key})
			}
		}
	}

	return mgrmodel.Manager{
		OdataContext:    "/redfish/v1/$metadata#Manager.Manager",
		OdataID:         "/redfish/v1/Managers/" + id,
		OdataType:       common.ManagerType,
		Name:            mgrData.Name,
		ManagerType:     mgrData.ManagerType,
		ID:              mgrData.ID,
		UUID:            mgrData.UUID,
		FirmwareVersion: mgrData.FirmwareVersion,
		Status: &mgrmodel.Status{
			State:  mgrData.State,
			Health: mgrData.Health,
		},
		Links: &mgrmodel.Links{
			ManagerForChassis:  chassisLink,
			ManagerForServers:  serverLink,
			ManagerForManagers: managerLink,
		},
		Description:         mgrData.Description,
		LogServices:         mgrData.LogServices,
		Model:               mgrData.Model,
		DateTime:            time.Now().Format(time.RFC3339),
		DateTimeLocalOffset: "+00:00",
		PowerState:          mgrData.PowerState,
	}, nil
}

// GetManagersResource is used to fetch resource data. The function is supposed to be used as part of RPC
// For getting system resource information,  parameters need to be passed GetSystemsRequest .
// GetManagersResource holds the  Uuid,Url and Resourceid ,
// Url will be parsed from that search key will created
// There will be two return values for the function. One is the RPC response, which contains the
// status code, status message, headers and body and the second value is error.
func (e *ExternalInterface) GetManagersResource(ctx context.Context, req *managersproto.ManagerRequest) response.RPC {
	var resp response.RPC
	var tableName string
	var resourceName string
	var resource map[string]interface{}
	requestData := strings.SplitN(req.ManagerID, ".", 2)
	urlData := strings.Split(req.URL, "/")
	if len(requestData) <= 1 {
		if req.ResourceID == "" {
			resourceName = urlData[len(urlData)-1]
			tableName = common.ManagersResource[resourceName]
		} else {
			tableName = urlData[len(urlData)-2]
		}
		data, err := e.DB.GetResource(tableName, req.URL)
		l.LogWithFields(ctx).Error("Data>>>>>>>>>>>>>", data)
		if err != nil {
			if req.ManagerID != config.Data.RootServiceUUID {
				return e.getPluginManagerResoure(ctx, requestData[0], req.URL)
			}
			errorMessage := "unable to get odimra managers details: " + err.Error()
			l.LogWithFields(ctx).Error(errorMessage)
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{}, nil)
			// return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, []interface{}{}, nil)
		}

		json.Unmarshal([]byte(data), &resource)
		resp.Body = resource
		resp.StatusCode = http.StatusOK
		resp.StatusMessage = response.Success
		respBody := fmt.Sprintf("%v", resp.Body)
		l.LogWithFields(ctx).Debugf("Outgoing manager resource response to northbound: %s", string(respBody))
		return resp

	}
	uuid := requestData[0]

	if req.ResourceID == "" {
		resourceName := urlData[len(urlData)-1]
		tableName = common.ManagersResource[resourceName]
	} else {
		tableName = urlData[len(urlData)-2]
	}
	data, err := e.DB.GetResource(tableName, req.URL)
	if err != nil {
		if errors.DBKeyNotFound == err.ErrNo() {
			var err error
			if data, err = e.getResourceInfoFromDevice(ctx, req.URL, uuid, requestData[1], nil); err != nil {
				errorMessage := "unable to get resource details from device: " + err.Error()
				l.LogWithFields(ctx).Error(errorMessage)
				errArgs := []interface{}{tableName, req.ManagerID}
				return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, errArgs, nil)
			}
		} else {
			errorMessage := "unable to get managers details: " + err.Error()
			l.LogWithFields(ctx).Error(errorMessage)
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, []interface{}{}, nil)
		}
	}

	json.Unmarshal([]byte(data), &resource)

	if common.Types[tableName] != "" && resource != nil {
		resource["@odata.type"] = common.Types[tableName]
	}

	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	respBody := fmt.Sprintf("%v", resp.Body)
	l.LogWithFields(ctx).Debugf("Outgoing manager resource response to northbound: %s", string(respBody))
	return resp
}

// VirtualMediaActions is used to perform action on VirtualMedia.
// For insert and eject of virtual media this function is used
func (e *ExternalInterface) VirtualMediaActions(ctx context.Context, req *managersproto.ManagerRequest, taskID string) {
	var resp response.RPC
	var requestBody = req.RequestBody
	targetURI := req.GetURL()
	//create task
	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: targetURI,
		UpdateTask: e.RPC.UpdateTask, TaskRequest: string(req.RequestBody)}
	//InsertMedia payload validation
	if strings.Contains(req.URL, "VirtualMedia.InsertMedia") {
		var vmiReq mgrmodel.VirtualMediaInsert
		// Updating the default values
		vmiReq.Inserted = true
		vmiReq.WriteProtected = true
		err := json.Unmarshal(req.RequestBody, &vmiReq)
		if err != nil {
			errorMessage := "while unmarshal the virtual media insert request: " + err.Error()
			l.LogWithFields(ctx).Error(errorMessage)
			common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, []interface{}{}, taskInfo)
			return
		}

		// Validating the request JSON properties for case sensitive
		invalidProperties, err := requestParamsCaseValidatorFunc(req.RequestBody, vmiReq)
		if err != nil {
			errMsg := "while validating request parameters for virtual media insert: " + err.Error()
			l.LogWithFields(ctx).Error(errMsg)
			common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
			return
		} else if invalidProperties != "" {
			errorMessage := "one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
			l.LogWithFields(ctx).Error(errorMessage)
			common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, taskInfo)
			return
		}

		// Check mandatory fields
		statuscode, statusMessage, messageArgs, err := validateFields(ctx, &vmiReq)
		if err != nil {
			errorMessage := "request payload validation failed: " + err.Error()
			l.LogWithFields(ctx).Error(errorMessage)
			common.GeneralError(statuscode, statusMessage, errorMessage, messageArgs, taskInfo)
			return
		}
		requestBody, err = json.Marshal(vmiReq)
		if err != nil {
			l.LogWithFields(ctx).Error("while marshalling the virtual media insert request: " + err.Error())
			common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, taskInfo)
			return
		}
	}
	// splitting managerID to get uuid
	requestData := strings.SplitN(req.ManagerID, ".", 2)
	uuid := requestData[0]
	plugin, resp := e.deviceCommunication(ctx, req.URL, uuid, requestData[1], http.MethodPost, requestBody)

	// If the virtualmedia action is success then updating DB
	if resp.StatusCode == http.StatusAccepted {
		services.SavePluginTaskInfo(ctx, plugin.PluginIP, plugin.PluginServerName, taskID, plugin.Location)
	}
	e.saveMediaDetails(ctx, req)
	respBody := fmt.Sprintf("%v", resp.Body)
	l.LogWithFields(ctx).Debugf("Outgoing virtual media response to northbound: %s", string(respBody))
}

// saveMediaDetails is used to save virtual media data in DB
func (e *ExternalInterface) saveMediaDetails(ctx context.Context, req *managersproto.ManagerRequest) {
	vmURI := strings.Replace(req.URL, "/Actions/VirtualMedia.InsertMedia", "", -1)
	vmURI = strings.Replace(vmURI, "/Actions/VirtualMedia.EjectMedia", "", -1)
	requestData := strings.SplitN(req.ManagerID, ".", 2)
	uuid := requestData[0]
	deviceData, err := e.getResourceInfoFromDevice(ctx, vmURI, uuid, requestData[1], nil)
	if err != nil {
		l.LogWithFields(ctx).Error("while trying get on URI " + vmURI + " : " + err.Error())
	} else {
		var vmData map[string]interface{}
		jerr := json.Unmarshal([]byte(deviceData), &vmData)
		if jerr != nil {
			l.LogWithFields(ctx).Error("while unmarshaling virtual media details: " + jerr.Error())
		} else {
			err = e.DB.UpdateData(vmURI, vmData, "VirtualMedia")
			if err != nil {
				l.LogWithFields(ctx).Error("while saving virtual media details: " + err.Error())
			}
		}
	}
}

// validateFields will validate the request payload, if any mandatory fields are missing then it will generate an error
func validateFields(ctx context.Context, request *mgrmodel.VirtualMediaInsert) (int32, string, []interface{}, error) {
	validate := validator.New()
	// if any of the mandatory fields missing in the struct, then it will return an error
	err := validate.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return http.StatusBadRequest, response.PropertyMissing, []interface{}{err.Field()}, fmt.Errorf(err.Field() + " field is missing")
		}
	}
	return http.StatusOK, common.OK, []interface{}{}, nil
}

func (e *ExternalInterface) getPluginManagerResoure(ctx context.Context, managerID, reqURI string) response.RPC {
	var resp response.RPC
	data, dberr := e.DB.GetManagerByURL("/redfish/v1/Managers/" + managerID)
	if dberr != nil {
		l.LogWithFields(ctx).Error("unable to get manager details : " + dberr.Error())
		var errArgs = []interface{}{"Managers", managerID}
		errorMessage := dberr.Error()
		resp = common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage,
			errArgs, nil)
		return resp
	}
	var managerData map[string]interface{}
	jerr := json.Unmarshal([]byte(data), &managerData)
	if jerr != nil {
		errorMessage := "unable to unmarshal manager details: " + jerr.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		resp = common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		return resp
	}
	var pluginID = managerData["Name"].(string)
	// Get the Plugin info
	plugin, gerr := e.DB.GetPluginData(pluginID)
	if gerr != nil {
		l.LogWithFields(ctx).Error("unable to get manager details : " + gerr.Error())
		var errArgs = []interface{}{"Plugin", pluginID}
		errorMessage := gerr.Error()
		resp = common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage,
			errArgs, nil)
		return resp
	}
	var req mgrcommon.PluginContactRequest

	req.ContactClient = e.Device.ContactClient
	req.Plugin = plugin

	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token := mgrcommon.GetPluginToken(ctx, req)
		if token == "" {
			var errorMessage = "unable to create session with plugin " + plugin.ID
			return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage,
				[]interface{}{}, nil)
		}
		req.Token = token
	} else {
		req.BasicAuth = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}

	req.OID = reqURI
	var errorMessage = "unable to get the details " + reqURI + ": "
	body, _, _, getResponse, err := mgrcommon.ContactPlugin(ctx, req, errorMessage)
	if err != nil {
		if getResponse.StatusCode == http.StatusUnauthorized && strings.EqualFold(req.Plugin.PreferredAuthType, "XAuthToken") {
			if body, _, _, getResponse, err = mgrcommon.RetryManagersOperation(ctx, req, errorMessage); err != nil {
				resp.StatusCode = getResponse.StatusCode
				json.Unmarshal(body, &resp.Body)
				respBody := fmt.Sprintf("%v", resp.Body)
				l.LogWithFields(ctx).Debugf("Outgoing plugin manager resoure response to northbound: %s", string(respBody))
				return resp
			}
		} else {
			resp.StatusCode = getResponse.StatusCode
			json.Unmarshal(body, &resp.Body)
			respBody := fmt.Sprintf("%v", resp.Body)
			l.LogWithFields(ctx).Debugf("Outgoing plugin manager resoure response to northbound: %s", string(respBody))
			return resp
		}
	}

	return fillResponse(ctx, body, managerData)

}

func fillResponse(ctx context.Context, body []byte, managerData map[string]interface{}) response.RPC {
	var resp response.RPC
	data := string(body)
	//replacing the response with north bound translation URL
	for key, value := range config.Data.URLTranslation.NorthBoundURL {
		data = strings.Replace(data, key, value, -1)
	}
	var respData map[string]interface{}
	err := json.Unmarshal([]byte(data), &respData)
	if err != nil {
		l.LogWithFields(ctx).Error(err.Error())
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(),
			[]interface{}{}, nil)
	}
	//To populate current Datetime and DateTimeLocalOffset for Plugin manager
	respData["DateTime"] = time.Now().Format(time.RFC3339)
	respData["DateTimeLocalOffset"] = "+00:00"

	if _, ok := respData["SerialConsole"]; !ok {
		respData["SerialConsole"] = dmtf.SerialConsole{}
	}
	respData["Links"] = managerData["Links"]
	resp.Body = respData
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	respBody := fmt.Sprintf("%v", resp.Body)
	l.LogWithFields(ctx).Debugf("Outgoing plugin manager resoure response to northbound: %s", string(respBody))
	return resp

}

func (e *ExternalInterface) getResourceInfoFromDevice(ctx context.Context, reqURL, uuid, systemID string, bmcCreds *mgrcommon.BmcUpdatedCreds) (string, error) {
	var getDeviceInfoRequest = mgrcommon.ResourceInfoRequest{
		URL:                   reqURL,
		UUID:                  uuid,
		SystemID:              systemID,
		ContactClient:         e.Device.ContactClient,
		DecryptDevicePassword: e.Device.DecryptDevicePassword,
		BmcUpdatedCreds:       bmcCreds,
	}
	return e.Device.GetDeviceInfo(ctx, getDeviceInfoRequest)

}

func (e *ExternalInterface) deviceCommunication(ctx context.Context, reqURL, uuid, systemID, httpMethod string,
	requestBody []byte) (mgrcommon.PluginTaskInfo, response.RPC) {
	var deviceInfoRequest = mgrcommon.ResourceInfoRequest{
		URL:                   reqURL,
		UUID:                  uuid,
		SystemID:              systemID,
		ContactClient:         e.Device.ContactClient,
		DecryptDevicePassword: e.Device.DecryptDevicePassword,
		HTTPMethod:            httpMethod,
		RequestBody:           requestBody,
	}
	return e.Device.DeviceRequest(ctx, deviceInfoRequest)
}

// GetRemoteAccountService is used to fetch resource data for BMC account service.
// ManagerRequest holds the UUID, URL and ResourceId ,
// There will be two return values for the function. One is the RPC response, which contains the
// status code, status message, headers and body and the second value is error.
func (e *ExternalInterface) GetRemoteAccountService(ctx context.Context, req *managersproto.ManagerRequest) response.RPC {
	var resp response.RPC

	requestData := strings.SplitN(req.ManagerID, ".", 2)
	uuid := requestData[0]
	uri := replaceBMCAccReq(req.URL, req.ManagerID)
	data, err := e.getResourceInfoFromDevice(ctx, uri, uuid, requestData[1], nil)
	if err != nil {
		return handleRemoteAccountServiceError(ctx, req.URL, req.ManagerID, err)
	}
	// Replace response body to BMC manager
	data = replaceBMCAccResp(data, req.ManagerID)
	resource := convertToRedfishModel(ctx, req.URL, data)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	respBody := fmt.Sprintf("%v", resp.Body)
	l.LogWithFields(ctx).Debugf("Outgoing remote account service response to northbound: %s", string(respBody))
	return resp
}

func handleRemoteAccountServiceError(ctx context.Context, uri, managerID string, err error) response.RPC {
	errorMessage := "error while trying to get resource details from device: " + err.Error()
	l.LogWithFields(ctx).Error(errorMessage)
	URIRegexAcc := regexp.MustCompile(`^\/redfish\/v1\/Managers\/[a-zA-Z0-9._-]+\/RemoteAccountService\/Accounts\/[a-zA-Z0-9._-]+[\/]?$`)
	URIRegexRoles := regexp.MustCompile(`^\/redfish\/v1\/Managers\/[a-zA-Z0-9._-]+\/RemoteAccountService\/Roles\/[a-zA-Z0-9._-]+[\/]?$`)
	if URIRegexAcc.MatchString(uri) {
		accID := uri[strings.LastIndex(uri, "/")+1:]
		errArgs := []interface{}{"Accounts", accID}
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, errArgs, nil)
	} else if URIRegexRoles.MatchString(uri) {
		roleID := uri[strings.LastIndex(uri, "/")+1:]
		errArgs := []interface{}{"Roles", roleID}
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, errArgs, nil)
	}
	errArgs := []interface{}{"Managers", managerID}
	return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, errArgs, nil)
}

func convertToRedfishModel(ctx context.Context, uri, data string) interface{} {
	URIRegexRemAcc := regexp.MustCompile(`^\/redfish\/v1\/Managers\/[a-zA-Z0-9._-]+\/RemoteAccountService+[\/]?$`)
	URIRegexAcc := regexp.MustCompile(`^\/redfish\/v1\/Managers\/[a-zA-Z0-9._-]+\/RemoteAccountService\/Accounts\/[a-zA-Z0-9._-]+[\/]?$`)
	URIRegexRoles := regexp.MustCompile(`^\/redfish\/v1\/Managers\/[a-zA-Z0-9._-]+\/RemoteAccountService\/Roles\/[a-zA-Z0-9._-]+[\/]?$`)
	if URIRegexRemAcc.MatchString(uri) {
		var resource dmtf.AccountService
		json.Unmarshal([]byte(data), &resource)
		resource.ODataType = common.ManagerAccountServiceType
		return resource
	} else if URIRegexAcc.MatchString(uri) {
		var resource dmtf.ManagerAccount
		json.Unmarshal([]byte(data), &resource)
		return resource
	} else if URIRegexRoles.MatchString(uri) {
		var resource dmtf.Role
		json.Unmarshal([]byte(data), &resource)
		return resource
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	return resource
}

// CreateRemoteAccountService is used to create BMC account user
func (e *ExternalInterface) CreateRemoteAccountService(ctx context.Context, req *managersproto.ManagerRequest,
	taskID string) {
	var requestBody = req.RequestBody
	var bmcAccReq mgrmodel.CreateBMCAccount
	uri := replaceBMCAccReq(req.URL, req.ManagerID)

	taskInfo := &common.TaskUpdateInfo{Context: ctx, TaskID: taskID, TargetURI: uri,
		UpdateTask: e.RPC.UpdateTask, TaskRequest: string(req.RequestBody)}

	// Updating the default values
	err := json.Unmarshal(req.RequestBody, &bmcAccReq)
	if err != nil {
		errorMessage := "error while unmarshaling the create remote account service request: " + err.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, []interface{}{}, taskInfo)
		return
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := requestParamsCaseValidatorFunc(req.RequestBody, bmcAccReq)
	if err != nil {
		errMsg := "error while validating request parameters for creating BMC account: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		return
	} else if invalidProperties != "" {
		errorMessage := "error in one or more properties given in the request body,they are not valid, ensure properties are listed in uppercamelcase "
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, taskInfo)
		return
	}

	// Check mandatory fields
	statuscode, statusMessage, messageArgs, err := validateCreateRemoteAccFields(ctx, &bmcAccReq)
	if err != nil {
		errorMessage := "error in request payload validation : " + err.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(statuscode, statusMessage, errorMessage, messageArgs, taskInfo)
		return
	}
	requestBody, err = json.Marshal(bmcAccReq)
	if err != nil {
		l.LogWithFields(ctx).Error("error while marshalling the create BMC account request: " + err.Error())
		common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, taskInfo)
		return
	}
	// splitting managerID to get uuid
	requestData := strings.SplitN(req.ManagerID, ".", 2)
	uuid := requestData[0]

	plugin, resp := e.deviceCommunication(ctx, uri, uuid, requestData[1], http.MethodPost, requestBody)

	if resp.StatusCode == http.StatusAccepted {
		e.DB.SavePluginTaskInfo(ctx, plugin.PluginIP, plugin.PluginServerName,
			taskID, plugin.Location)
		return
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		body, _ := json.Marshal(resp.Body)
		respBody := replaceBMCAccResp(string(body), req.ManagerID)
		var managerAcc dmtf.ManagerAccount
		json.Unmarshal([]byte(respBody), &managerAcc)
		resp.Body = managerAcc
		resp.StatusCode = http.StatusCreated
		task := fillTaskData(taskID, uri, string(req.RequestBody), resp,
			common.Completed, common.OK, 100, http.MethodPost)
		e.RPC.UpdateTask(ctx, task)
	} else {
		task := fillTaskData(taskID, uri, string(req.RequestBody), resp,
			common.Completed, common.Warning, 100, http.MethodPost)
		e.RPC.UpdateTask(ctx, task)
	}
	respBody := fmt.Sprintf("%v", resp.Body)
	l.LogWithFields(ctx).Debugf("Outgoing remote account service response to northbound: %s", string(respBody))
}

// validateFields will validate the request payload, if any mandatory fields are missing then it will generate an error
func validateCreateRemoteAccFields(ctx context.Context, request *mgrmodel.CreateBMCAccount) (int32, string, []interface{}, error) {
	validate := validator.New()
	// if any of the mandatory fields missing in the struct, then it will return an error
	err := validate.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return http.StatusBadRequest, response.PropertyMissing, []interface{}{err.Field()}, fmt.Errorf(err.Field() + " field is missing")
		}
	}
	return http.StatusOK, common.OK, []interface{}{}, nil
}

func replaceBMCAccReq(uri, managerID string) string {
	uri = strings.Replace(uri, "Managers/"+managerID+"/Remote", "", -1)
	return uri
}

func replaceBMCAccResp(data, managerID string) string {
	data = strings.Replace(data, "v1/AccountService", "v1/Managers/"+managerID+"/RemoteAccountService", -1)
	return data
}

// UpdateRemoteAccountService is used to update BMC account
func (e *ExternalInterface) UpdateRemoteAccountService(ctx context.Context, req *managersproto.ManagerRequest,
	taskID string) {
	var bmcAccReq mgrmodel.UpdateBMCAccount

	uri := replaceBMCAccReq(req.URL, req.ManagerID)
	taskInfo := &common.TaskUpdateInfo{Context: ctx, TaskID: taskID, TargetURI: uri,
		UpdateTask: e.RPC.UpdateTask, TaskRequest: string(req.RequestBody)}

	// Updating the default values
	err := json.Unmarshal(req.RequestBody, &bmcAccReq)
	if err != nil {
		errorMessage := "error while unmarshal the update remote account service request: " + err.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage,
			[]interface{}{}, taskInfo)
		return
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := requestParamsCaseValidatorFunc(req.RequestBody, bmcAccReq)
	if err != nil {
		errMsg := "error while validating request parameters for updating BMC account: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
		return
	} else if invalidProperties != "" {
		errorMessage := "error in one or more properties given in the request body,they are not valid, ensure properties are listed in uppercamelcase "
		l.LogWithFields(ctx).Error(errorMessage)
		common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage,
			[]interface{}{invalidProperties}, taskInfo)
		return
	}

	// splitting managerID to get uuid
	requestData := strings.SplitN(req.ManagerID, ".", 2)
	uuid := requestData[0]
	//do get call with
	data, err := e.getResourceInfoFromDevice(ctx, uri, uuid, requestData[1], nil)
	if err != nil {
		errorMessage := "error while trying to get resource details from device: " + err.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		errArgs := []interface{}{"RemoteAccounts", requestData[1]}
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, errArgs, taskInfo)
		return
	}
	var dataMap map[string]interface{}
	err = json.Unmarshal([]byte(data), &dataMap)
	if err != nil {
		panic(err)
	}

	var username string
	if dataMap["UserName"] != nil {
		username = dataMap["UserName"].(string)
	}
	bmcCreds := mgrcommon.BmcUpdatedCreds{UserName: username, UpdatedPassword: bmcAccReq.Password}
	requestBody, err := json.Marshal(bmcAccReq)
	if err != nil {
		l.LogWithFields(ctx).Error("error while marshalling the update BMC account request: " + err.Error())
		common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, taskInfo)
		return
	}

	plugin, resp := e.deviceCommunication(ctx, uri, uuid, requestData[1], http.MethodPatch, requestBody)
	if resp.StatusCode == http.StatusAccepted {
		e.DB.SavePluginTaskInfo(ctx, plugin.PluginIP, plugin.PluginServerName,
			taskID, plugin.Location)
		e.UpdatePassword(ctx, uuid, bmcAccReq.Password, username, resp.StatusCode)
		return
	}

	if resp.StatusCode == http.StatusOK {
		data, err := e.getResourceInfoFromDevice(ctx, uri, uuid, requestData[1], &bmcCreds)
		if err != nil {
			errorMessage := "error while trying to get resource details from device: " + err.Error()
			l.LogWithFields(ctx).Error(errorMessage)
			errArgs := []interface{}{}
			common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, errArgs, taskInfo)
			return
		}
		e.UpdatePassword(ctx, uuid, bmcAccReq.Password, username, resp.StatusCode)
		// Replace response body to BMC manager
		data = replaceBMCAccResp(data, req.ManagerID)
		resource := convertToRedfishModel(ctx, req.URL, data)
		resp.Body = resource
		resp.StatusCode = http.StatusOK
		resp.StatusMessage = response.Success
		task := fillTaskData(taskID, uri, string(req.RequestBody), resp,
			common.Completed, common.OK, 100, http.MethodPatch)
		e.RPC.UpdateTask(ctx, task)
	} else {
		task := fillTaskData(taskID, uri, string(req.RequestBody), resp,
			common.Completed, common.Warning, 100, http.MethodPatch)
		e.RPC.UpdateTask(ctx, task)
	}
	respBody := fmt.Sprintf("%v", resp.Body)
	l.LogWithFields(ctx).Debugf("Outgoing update remote account service response to northbound: %s", string(respBody))
}

// DeleteRemoteAccountService is used to delete the BMC account user
func (e *ExternalInterface) DeleteRemoteAccountService(ctx context.Context, req *managersproto.ManagerRequest,
	taskID string) {
	// splitting managerID to get uuid
	requestData := strings.SplitN(req.ManagerID, ".", 2)
	uuid := requestData[0]
	uri := replaceBMCAccReq(req.URL, req.ManagerID)

	plugin, resp := e.deviceCommunication(ctx, uri, uuid, requestData[1], http.MethodDelete, nil)
	if resp.StatusCode == http.StatusAccepted {
		e.DB.SavePluginTaskInfo(ctx, plugin.PluginIP, plugin.PluginServerName,
			taskID, plugin.Location)
		return
	}

	if resp.StatusCode == http.StatusOK {
		resp.StatusCode = http.StatusNoContent
	}
	task := fillTaskData(taskID, uri, string(req.RequestBody), resp,
		common.Completed, common.OK, 100, http.MethodDelete)
	e.RPC.UpdateTask(ctx, task)
	respBody := fmt.Sprintf("%v", resp.Body)
	l.LogWithFields(ctx).Debugf("Outgoing delete remote account service response to northbound: %s", string(respBody))
}

// UpdatePassword method used to update system password
func (e ExternalInterface) UpdatePassword(ctx context.Context, uuid, password string, userName string, statusCode int32) {
	target, gerr := mgrmodel.GetTarget(uuid)
	if gerr != nil {
		l.LogWithFields(ctx).Error("error while getting device details :" + gerr.Error())
		return
	}
	if userName != target.UserName {
		return
	}
	newPassword, err := e.Device.EncryptDevicePassword([]byte(password))
	if err != nil {
		errMsg := "error while trying to encrypt: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return
	}
	target.Password = newPassword
	if statusCode == http.StatusOK {
		err1 := mgrmodel.UpdateSystem(target.ManagerAddress, target)
		if err1 != nil {
			errMsg := "error while update password : " + err.Error()
			l.LogWithFields(ctx).Error(errMsg)
		}
		return
	}
	err1 := mgrmodel.AddTempPassword(uuid, target)
	if err1 != nil {
		errMsg := "error while update password : " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return
	}
}

// UpdateRemoteAccountService is used to update BMC account
func (e *ExternalInterface) UpdateRemoteAccountPasswordService(ctx context.Context, req *managersproto.ManagerRequest) {
	target, err := mgrmodel.GetTempPassword(req.ManagerID)
	if err != nil {
		errMsg := "no password found: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return
	}
	if err := mgrmodel.UpdateSystem(target.DeviceUUID, target); err != nil {
		errMsg := "error while update password : " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return
	}
	if err := mgrmodel.DeleteTempPassword(req.ManagerID); err != nil {
		errMsg := "error while delete temp password : " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return
	}
	l.LogWithFields(ctx).Info("Password updated successfully for device " + target.DeviceUUID)
}
