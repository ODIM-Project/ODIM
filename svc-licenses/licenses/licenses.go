//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

package licenses

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	licenseproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/licenses"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	lcommon "github.com/ODIM-Project/ODIM/svc-licenses/lcommon"
	"github.com/ODIM-Project/ODIM/svc-licenses/model"
)

var (
	jsonUnMarshalFunc = json.Unmarshal
	jsonMarshalFunc   = json.Marshal
)

// GetLicenseService to get license service details
func (e *ExternalInterface) GetLicenseService(req *licenseproto.GetLicenseServiceRequest) response.RPC {
	var resp response.RPC
	license := dmtf.LicenseService{
		OdataContext:   "/redfish/v1/$metadata#LicenseService.LicenseService",
		OdataID:        "/redfish/v1/LicenseService",
		OdataType:      "#LicenseService.v1_0_0.LicenseService",
		ID:             "LicenseService",
		Description:    "License Service",
		Name:           "License Service",
		ServiceEnabled: true,
	}
	license.Licenses = &dmtf.Link{Oid: "/redfish/v1/LicenseService/Licenses"}

	resp.Body = license
	resp.StatusCode = http.StatusOK
	return resp
}

// GetLicenseCollection to get license collection details
func (e *ExternalInterface) GetLicenseCollection(ctx context.Context, req *licenseproto.GetLicenseRequest) response.RPC {
	var resp response.RPC
	licenseCollection := dmtf.LicenseCollection{
		OdataContext: "/redfish/v1/$metadata#LicenseCollection.LicenseCollection",
		OdataID:      "/redfish/v1/LicenseService/Licenses",
		OdataType:    "#LicenseCollection.LicenseCollection",
		Description:  "License Collection",
		Name:         "License Collection",
	}
	var members = make([]*dmtf.Link, 0)

	licenseCollectionKeysArray, err := e.DB.GetAllKeysFromTable(ctx, "Licenses", persistencemgr.InMemory)
	if err != nil {
		l.LogWithFields(ctx).Error("error while getting license collection details from db")
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
	}
	if len(licenseCollectionKeysArray) == 0 {
		l.LogWithFields(ctx).Error("odimra doesnt have Licenses")
	}

	for _, key := range licenseCollectionKeysArray {
		members = append(members, &dmtf.Link{Oid: key})
	}
	licenseCollection.Members = members
	licenseCollection.MembersCount = len(members)
	resp.Body = licenseCollection
	resp.StatusCode = http.StatusOK
	return resp
}

// GetLicenseResource to get individual license resource
func (e *ExternalInterface) GetLicenseResource(ctx context.Context, req *licenseproto.GetLicenseResourceRequest) response.RPC {
	var resp response.RPC
	licenseResp := dmtf.License{}
	uri := req.URL
	ID := strings.Split(uri, "/")

	data, dbErr := e.DB.GetResource("Licenses", uri, persistencemgr.InMemory)
	if dbErr != nil {
		l.LogWithFields(ctx).Error("Unable to get license data : " + dbErr.Error())
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, dbErr.Error(), nil, nil)
	}

	if data != "" {
		err := json.Unmarshal([]byte(data.(string)), &licenseResp)
		if err != nil {
			l.LogWithFields(ctx).Error("Unable to unmarshall  the data: " + err.Error())
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		}
	}
	licenseResp.OdataContext = "/redfish/v1/$metadata#License.License"
	licenseResp.OdataType = "#License.v1_0_0.License"
	licenseResp.OdataID = uri
	licenseResp.ID = ID[len(ID)-1]

	resp.Body = licenseResp
	resp.StatusCode = http.StatusOK
	return resp
}

// InstallLicenseService to install license
func (e *ExternalInterface) InstallLicenseService(ctx context.Context, req *licenseproto.InstallLicenseRequest, sessionUserName, taskID string) {
	var resp response.RPC
	var installreq dmtf.LicenseInstallRequest
	var percentComplete int32

	targetURI := "/redfish/v1/LicenseService/Licenses"

	taskInfo := &common.TaskUpdateInfo{Context: ctx, TaskID: taskID, TargetURI: targetURI,
		UpdateTask: e.External.UpdateTask, TaskRequest: string(req.RequestBody)}

	genErr := jsonUnMarshalFunc(req.RequestBody, &installreq)
	if genErr != nil {
		errMsg := "Unable to unmarshal the install license request" + genErr.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusBadRequest, response.InternalError, errMsg, nil, taskInfo)
		return
	}

	if installreq.Links == nil {
		errMsg := "Invalid request,mandatory field Links missing"
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"Links"}, taskInfo)
		return
	} else if installreq.LicenseString == "" {
		errMsg := "Invalid request, mandatory field LicenseString is missing"
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"LicenseString"}, taskInfo)
		return
	} else if len(installreq.Links.Link) == 0 {
		errMsg := "Invalid request, mandatory field AuthorizedDevices links is missing"
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"LicenseString"}, taskInfo)
		return
	}

	linksMap, errStatusCode, err := e.getManagerLinksMap(ctx, installreq.Links.Link)
	if err != nil {
		l.LogWithFields(ctx).Error(err)
		common.GeneralError(errStatusCode, response.InternalError, err.Error(), nil, taskInfo)
	}
	l.LogWithFields(ctx).Debug("Map with manager Links: ", linksMap)
	partialResultFlag := false
	subTaskChannel := make(chan int32, len(linksMap))
	for serverURI := range linksMap {
		uuid, managerID, err := lcommon.GetIDsFromURI(serverURI)
		if err != nil {
			errMsg := "error while trying to get system ID from " + serverURI + ": " + err.Error()
			l.LogWithFields(ctx).Error(errMsg)
			common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"SystemID", serverURI}, taskInfo)
			return
		}

		encodedKey := base64.StdEncoding.EncodeToString([]byte(installreq.LicenseString))
		managerURI := "/redfish/v1/Managers/" + managerID
		reqPostBody := map[string]interface{}{"LicenseString": encodedKey, "AuthorizedDevices": managerURI}
		reqBody, _ := json.Marshal(reqPostBody)
		var threadID int = 1
		ctxt := context.WithValue(ctx, common.ThreadName, common.SendRequest)
		ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
		go e.sendRequest(ctx, uuid, sessionUserName, taskID, serverURI, reqBody, subTaskChannel)
		threadID++
	}

	resp.StatusCode = http.StatusOK
	for i := 0; i < len(linksMap); i++ {
		select {
		case statusCode := <-subTaskChannel:
			if statusCode != http.StatusOK {
				if statusCode != http.StatusAccepted {
					partialResultFlag = true
				}
				if resp.StatusCode < statusCode {
					resp.StatusCode = statusCode
				}
			}
			if i < len(linksMap)-1 && statusCode != http.StatusAccepted {
				percentComplete := int32(((i + 1) / len(linksMap)) * 100)
				var task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Running, common.OK, percentComplete, http.MethodPost)
				err := e.External.UpdateTask(ctx, task)
				if err != nil && err.Error() == common.Cancelling {
					task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Cancelled, common.OK, percentComplete, http.MethodPost)
					e.External.UpdateTask(ctx, task)
					runtime.Goexit()
				}

			}
		}
	}

	taskStatus := common.OK
	if partialResultFlag {
		taskStatus = common.Warning
	}

	if resp.StatusCode == http.StatusAccepted {
		return
	}

	percentComplete = 100
	if resp.StatusCode != http.StatusOK {
		errMsg := "One or more of the Install License requests failed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID
		l.LogWithFields(ctx).Warn(errMsg)
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			common.GeneralError(http.StatusUnauthorized, response.ResourceAtURIUnauthorized, errMsg, []interface{}{fmt.Sprintf("%v", linksMap)}, taskInfo)
			return
		case http.StatusNotFound:
			common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"option", "Licenses"}, taskInfo)
			return
		case http.StatusBadRequest:
			common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errMsg, []interface{}{"Licenses"}, taskInfo)
			return
		default:
			common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
			return
		}
	}

	l.LogWithFields(ctx).Info("All Install License requests successfully completed. for more information please check SubTasks in URI: /redfish/v1/TaskService/Tasks/" + taskID)
	resp.StatusMessage = response.Success
	resp.StatusCode = http.StatusOK
	args := response.Args{
		Code:    resp.StatusMessage,
		Message: "Request completed successfully",
	}
	resp.Body = args.CreateGenericErrorResponse()

	var task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Completed, taskStatus, percentComplete, http.MethodPost)
	err = e.External.UpdateTask(ctx, task)
	if err != nil && err.Error() == common.Cancelling {
		task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Cancelled, common.Critical, percentComplete, http.MethodPost)
		e.External.UpdateTask(ctx, task)
		runtime.Goexit()
	}
	respBody := fmt.Sprintf("%v", resp.Body)
	l.LogWithFields(ctx).Debugf("final response for install license request: %s", string(respBody))
}

func (e *ExternalInterface) getManagerLinksMap(ctx context.Context, links []*dmtf.Link) (map[string]bool, int32, error) {
	var serverURI string
	var err error
	var managerLink []string
	var errStatusCode int32
	linksMap := make(map[string]bool)
	for _, serverIDs := range links {
		serverURI = serverIDs.Oid
		switch {
		case strings.Contains(serverURI, "Systems"):
			managerLink, err = e.getManagerURL(serverURI)
			if err != nil {
				errMsg := "Unable to get manager link for " + serverURI
				errStatusCode = http.StatusNotFound
				return linksMap, errStatusCode, fmt.Errorf(errMsg)
			}
			for _, link := range managerLink {
				linksMap[link] = true
			}
		case strings.Contains(serverURI, "Managers"):
			linksMap[serverURI] = true
		case strings.Contains(serverURI, "Aggregates"):
			managerLink, err = e.getDetailsFromAggregate(ctx, serverURI)
			if err != nil {
				errMsg := "Unable to get manager link from aggregates for " + serverURI
				errStatusCode = http.StatusNotFound
				return linksMap, errStatusCode, fmt.Errorf(errMsg)
			}
			for _, link := range managerLink {
				linksMap[link] = true
			}
		default:
			errMsg := "Unable to get manager link from aggregates for " + serverURI
			errStatusCode = http.StatusBadRequest
			return linksMap, errStatusCode, fmt.Errorf(errMsg)
		}
	}
	return linksMap, errStatusCode, nil
}

// sendRequest request the plugin to install license and handles the response from plugin
func (e *ExternalInterface) sendRequest(ctx context.Context, uuid, sessionUserName, taskID, serverURI string, requestBody []byte,
	subTaskChannel chan<- int32) {

	var percentComplete int32
	var resp response.RPC
	subTaskURI, err := e.External.CreateChildTask(ctx, sessionUserName, taskID)
	if err != nil {
		subTaskChannel <- http.StatusInternalServerError
		l.LogWithFields(ctx).Warn("Unable to create sub task: " + err.Error())
		return
	}
	var subTaskID string
	strArray := strings.Split(subTaskURI, "/")
	if strings.HasSuffix(subTaskURI, "/") {
		subTaskID = strArray[len(strArray)-2]
	} else {
		subTaskID = strArray[len(strArray)-1]
	}

	taskInfo := &common.TaskUpdateInfo{Context: ctx, TaskID: subTaskID, TargetURI: serverURI, UpdateTask: e.External.UpdateTask, TaskRequest: string(requestBody)}

	var contactRequest model.PluginContactRequest

	// Get target device Credentials from using device UUID
	target, targetErr := e.External.GetTarget(uuid)
	if targetErr != nil {
		subTaskChannel <- http.StatusNotFound
		errMsg := targetErr.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"target", uuid}, taskInfo)
		return
	}

	decryptedPasswordByte, err := e.External.DevicePassword(target.Password)
	if err != nil {
		subTaskChannel <- http.StatusInternalServerError
		errMsg := "error while trying to decrypt device password: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return
	}
	target.Password = decryptedPasswordByte

	// Get the Plugin info
	plugin, errs := e.External.GetPluginData(target.PluginID)
	if errs != nil {
		subTaskChannel <- http.StatusNotFound
		errMsg := "error while getting plugin data: " + errs.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"PluginData", target.PluginID}, taskInfo)
		return
	}
	l.LogWithFields(ctx).Info("Plugin info: ", plugin)

	contactRequest.Plugin = *plugin
	contactRequest.ContactClient = e.External.ContactClient
	contactRequest.Plugin.ID = target.PluginID
	contactRequest.HTTPMethodType = http.MethodPost

	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		contactRequest.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		contactRequest.OID = "/ODIM/v1/Sessions"
		l.LogWithFields(ctx).Debugf("plugin contact request data for %s: %s", contactRequest.OID, string(requestBody))
		_, token, _, getResponse, err := lcommon.ContactPlugin(ctx, contactRequest, "error while logging in to plugin: ")
		if err != nil {
			subTaskChannel <- getResponse.StatusCode
			errMsg := err.Error()
			l.LogWithFields(ctx).Error(errMsg)
			common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo)
			return
		}
		contactRequest.Token = token
	} else {
		contactRequest.LoginCredentials = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	target.PostBody = []byte(requestBody)
	contactRequest.DeviceInfo = target
	contactRequest.OID = "/ODIM/v1/LicenseService/Licenses"
	contactRequest.PostBody = requestBody
	l.LogWithFields(ctx).Debugf("plugin contact request data for %s: %s", contactRequest.OID, string(requestBody))
	_, _, pluginTaskInfo, getResponse, err := e.External.ContactPlugin(ctx, contactRequest, "error while installing license: ")
	if err != nil {
		subTaskChannel <- getResponse.StatusCode
		errMsg := err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo)
		return
	}
	l.LogWithFields(ctx).Debugf("Install license response: status code : %v, message: %s",
		getResponse.StatusCode, getResponse.StatusMessage)

	if getResponse.StatusCode == http.StatusAccepted {
		services.SavePluginTaskInfo(ctx, pluginTaskInfo.PluginIP, plugin.IP,
			subTaskID, pluginTaskInfo.Location)
		subTaskChannel <- http.StatusAccepted
		return
	}

	if getResponse.StatusCode > http.StatusMultipleChoices {
		resp.StatusCode = getResponse.StatusCode
		subTaskChannel <- getResponse.StatusCode
		percentComplete = 100
		task := fillTaskData(subTaskID, serverURI, string(requestBody), resp, common.Completed, common.Warning, percentComplete, http.MethodPost)
		e.External.UpdateTask(ctx, task)
		return
	}

	resp.StatusCode = http.StatusOK
	percentComplete = 100

	subTaskChannel <- http.StatusOK
	var task = fillTaskData(subTaskID, serverURI, string(requestBody), resp, common.Completed, common.OK, percentComplete, http.MethodPost)
	err = e.External.UpdateTask(ctx, task)
	if err != nil && err.Error() == common.Cancelling {
		var task = fillTaskData(subTaskID, serverURI, string(requestBody), resp, common.Cancelled, common.Critical, percentComplete, http.MethodPost)
		e.External.UpdateTask(ctx, task)
	}
	return
}

func (e *ExternalInterface) getDetailsFromAggregate(ctx context.Context, aggregateURI string) ([]string, error) {
	var resource model.Elements
	var links []string
	respData, err := e.DB.GetResource("Aggregate", aggregateURI, persistencemgr.OnDisk)
	if err != nil {
		return nil, err
	}
	jsonStr, jerr := jsonMarshalFunc(respData)
	if jerr != nil {
		return nil, jerr
	}
	jerr = jsonUnMarshalFunc([]byte(jsonStr), &resource)
	if jerr != nil {
		return nil, jerr
	}
	l.LogWithFields(ctx).Debug("System URL's from agrregate: ", resource)

	for _, key := range resource.Elements {
		res, err := e.getManagerURL(key.OdataID)
		if err != nil {
			errMsg := "Unable to get manager link"
			l.LogWithFields(ctx).Error(errMsg)
			return nil, err
		}
		links = append(links, res...)
	}
	l.LogWithFields(ctx).Debug("manager links: ", links)
	return links, nil
}

func (e *ExternalInterface) getManagerURL(systemURI string) ([]string, error) {
	var resource dmtf.ComputerSystem
	var managerLink string
	var links []string
	respData, err := e.DB.GetResource("ComputerSystem", systemURI, persistencemgr.InMemory)
	if err != nil {
		return nil, err
	}
	jerr := jsonUnMarshalFunc([]byte(respData.(string)), &resource)
	if jerr != nil {
		return nil, jerr
	}
	members := resource.Links.ManagedBy
	for _, member := range members {
		managerLink = member.Oid
	}
	links = append(links, managerLink)
	return links, nil
}
