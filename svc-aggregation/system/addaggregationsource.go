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

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
)

// AddAggregationSource is the handler for adding bmc or manager
// Discovers  its top level odata.ID links and store them in inmemory db.
// Upon successfull operation this api returns added AggregationSourceUUID  in the response body with 201 OK.
func (e *ExternalInterface) AddAggregationSource(taskID string, sessionUserName string, req *aggregatorproto.AggregatorRequest) response.RPC {
	targetURI := "/redfish/v1/AggregationService/AggregationSources"
	var resp response.RPC
	var percentComplete int32
	var task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Running, common.OK, percentComplete, http.MethodPost)
	err := e.UpdateTask(task)
	if err != nil {
		errMsg := "error while starting the task: " + err.Error()
		log.Printf("error while starting the task: %v", errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: targetURI, UpdateTask: e.UpdateTask, TaskRequest: string(req.RequestBody)}
	// parsing the request
	var aggregationSourceRequest AggregationSource
	err = json.Unmarshal(req.RequestBody, &aggregationSourceRequest)
	if err != nil {
		errMsg := "unable to parse the add request" + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	}
	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, aggregationSourceRequest)
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

	if aggregationSourceRequest.Links == nil || (aggregationSourceRequest.Links.ConnectionMethod == nil && aggregationSourceRequest.Links.Oem == nil) {
		errMsg := "error: mandatory ConnectionMethod/Oem block missing in the request"
		log.Println(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"ConnectionMethod/Oem"}, taskInfo)
	}
	if aggregationSourceRequest.Links.ConnectionMethod != nil {
		return e.addAggregationSourceWithConnectionMethod(taskID, targetURI, string(req.RequestBody), percentComplete, aggregationSourceRequest, taskInfo)
	}
	var addResourceRequest = AddResourceRequest{
		ManagerAddress: aggregationSourceRequest.HostName,
		UserName:       aggregationSourceRequest.UserName,
		Password:       aggregationSourceRequest.Password,
		Oem:            aggregationSourceRequest.Links.Oem,
	}
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
		e.UpdateTask(fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Exception, common.Warning, percentComplete, http.MethodPost))
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
	var aggregationSourceUUID string
	var cipherText []byte
	if aggregationSourceRequest.Links.Oem.PluginType != "" || aggregationSourceRequest.Links.Oem.PreferredAuthType != "" {
		resp, aggregationSourceUUID, cipherText = e.addPluginData(addResourceRequest, taskID, targetURI, pluginContactRequest)
	} else {
		resp, aggregationSourceUUID, cipherText = e.addCompute(taskID, targetURI, addResourceRequest.Oem.PluginID, percentComplete, addResourceRequest, pluginContactRequest)
	}
	if resp.StatusMessage != "" {
		return resp
	}
	// Adding Aggregation Source to db
	var aggregationSourceData = agmodel.AggregationSource{
		HostName: aggregationSourceRequest.HostName,
		UserName: aggregationSourceRequest.UserName,
		Password: cipherText,
		Links:    aggregationSourceRequest.Links,
	}
	var aggregationSourceURI = fmt.Sprintf("%s/%s", targetURI, aggregationSourceUUID)
	dbErr := agmodel.AddAggregationSource(aggregationSourceData, aggregationSourceURI)
	if dbErr != nil {
		errMsg := dbErr.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	}
	commonResponse := response.Response{
		OdataType:    "#AggregationSource.v1_0_0.AggregationSource",
		OdataID:      aggregationSourceURI,
		OdataContext: "/redfish/v1/$metadata#AggregationSource.AggregationSource",
		ID:           aggregationSourceUUID,
		Name:         "Aggregation Source",
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Link":              "<" + aggregationSourceURI + "/>; rel=describedby",
		"Location":          aggregationSourceURI,
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	commonResponse.CreateGenericResponse(response.Created)
	commonResponse.Message = ""
	commonResponse.MessageID = ""
	commonResponse.Severity = ""
	resp.Body = agresponse.AggregationSourceResponse{
		Response: commonResponse,
		HostName: aggregationSourceRequest.HostName,
		UserName: aggregationSourceRequest.UserName,
		Links:    aggregationSourceRequest.Links,
	}
	resp.StatusCode = http.StatusCreated
	percentComplete = 100

	task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Completed, common.OK, percentComplete, http.MethodPost)
	e.UpdateTask(task)
	return resp
}

func (e *ExternalInterface) addAggregationSourceWithConnectionMethod(taskID, targetURI, reqBody string, percentComplete int32, aggregationSourceRequest AggregationSource, taskInfo *common.TaskUpdateInfo) response.RPC {
	var resp response.RPC
	var addResourceRequest = AddResourceRequest{
		ManagerAddress:   aggregationSourceRequest.HostName,
		UserName:         aggregationSourceRequest.UserName,
		Password:         aggregationSourceRequest.Password,
		ConnectionMethod: aggregationSourceRequest.Links.ConnectionMethod,
	}
	ActiveReqSet.UpdateMu.Lock()
	if _, exist := ActiveReqSet.ReqRecord[addResourceRequest.ManagerAddress]; exist {
		ActiveReqSet.UpdateMu.Unlock()
		var errMsg string
		mIP, _ := getIPAndPortFromAddress(addResourceRequest.ManagerAddress)
		errMsg = fmt.Sprintf("error: An active request already exists for adding aggregation source IP %v", mIP)
		log.Println(errMsg)
		args := response.Args{
			Code:    response.GeneralError,
			Message: errMsg,
		}
		resp.Body = args.CreateGenericErrorResponse()
		resp.Header = map[string]string{"Content-type": "application/json; charset=utf-8"}
		resp.StatusCode = http.StatusConflict
		percentComplete = 100
		e.UpdateTask(fillTaskData(taskID, targetURI, reqBody, resp, common.Exception, common.Warning, percentComplete, http.MethodPost))
		return resp
	}
	ActiveReqSet.ReqRecord[addResourceRequest.ManagerAddress] = addResourceRequest.ConnectionMethod.OdataID
	ActiveReqSet.UpdateMu.Unlock()

	defer func() {
		// check if there is an entry added for the server in ongoing requests tracker and remove it
		ActiveReqSet.UpdateMu.Lock()
		delete(ActiveReqSet.ReqRecord, addResourceRequest.ManagerAddress)
		ActiveReqSet.UpdateMu.Unlock()
	}()

	connectionMethod, err1 := e.GetConnectionMethod(addResourceRequest.ConnectionMethod.OdataID)
	if err1 != nil {
		errMsg := "error while getting connection method id: " + err1.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"connectionmethod id", addResourceRequest.ConnectionMethod.OdataID}, taskInfo)
	}
	cmVariants := getConnectionMethodVariants(connectionMethod.ConnectionMethodVariant)
	var pluginContactRequest getResourceRequest
	pluginContactRequest.ContactClient = e.ContactClient
	pluginContactRequest.GetPluginStatus = e.GetPluginStatus
	pluginContactRequest.TargetURI = targetURI
	pluginContactRequest.UpdateTask = e.UpdateTask
	var aggregationSourceUUID string
	var cipherText []byte

	// check status will do call on the URI /ODIM/v1/Status to the requested manager address
	// if its success then add the plugin, else if its not found then add BMC
	// else return the response
	statusResp, statusCode, queueList := checkStatus(pluginContactRequest, addResourceRequest, cmVariants, taskInfo)
	if statusCode == http.StatusOK {
		// check if AggregationSource has any values, if its there means its managing the bmcs
		if len(connectionMethod.Links.AggregationSources) > 0 {
			errMsg := "Cant proceed to add aggregation source, since connection method is already managing other aggregation sources"
			log.Println(errMsg)
			return common.GeneralError(http.StatusNotAcceptable, response.InternalError, errMsg, nil, taskInfo)
		}
		resp, aggregationSourceUUID, cipherText = e.addPluginDataWIthConnectionMethod(addResourceRequest, taskID, targetURI, pluginContactRequest, queueList, cmVariants)
	} else if statusCode == http.StatusNotFound {
		resp, aggregationSourceUUID, cipherText = e.addCompute(taskID, targetURI, cmVariants.PluginID, percentComplete, addResourceRequest, pluginContactRequest)
	} else {
		return statusResp
	}
	if resp.StatusMessage != "" {
		return resp
	}
	// Adding Aggregation Source to db
	var aggregationSourceData = agmodel.AggregationSource{
		HostName: aggregationSourceRequest.HostName,
		UserName: aggregationSourceRequest.UserName,
		Password: cipherText,
		Links:    aggregationSourceRequest.Links,
	}
	var aggregationSourceURI = fmt.Sprintf("%s/%s", targetURI, aggregationSourceUUID)
	dbErr := agmodel.AddAggregationSource(aggregationSourceData, aggregationSourceURI)
	if dbErr != nil {
		errMsg := dbErr.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	}

	connectionMethod.Links.AggregationSources = append(connectionMethod.Links.AggregationSources, agmodel.OdataID{OdataID: aggregationSourceURI})
	dbErr = e.UpdateConnectionMethod(connectionMethod, addResourceRequest.ConnectionMethod.OdataID)
	if dbErr != nil {
		errMsg := dbErr.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	}
	commonResponse := response.Response{
		OdataType:    "#AggregationSource.v1_0_0.AggregationSource",
		OdataID:      aggregationSourceURI,
		OdataContext: "/redfish/v1/$metadata#AggregationSource.AggregationSource",
		ID:           aggregationSourceUUID,
		Name:         "Aggregation Source",
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Link":              "<" + aggregationSourceURI + "/>; rel=describedby",
		"Location":          aggregationSourceURI,
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	commonResponse.CreateGenericResponse(response.Created)
	commonResponse.Message = ""
	commonResponse.MessageID = ""
	commonResponse.Severity = ""
	resp.Body = agresponse.AggregationSourceResponse{
		Response: commonResponse,
		HostName: aggregationSourceRequest.HostName,
		UserName: aggregationSourceRequest.UserName,
		Links:    aggregationSourceRequest.Links,
	}
	resp.StatusCode = http.StatusCreated
	percentComplete = 100

	task := fillTaskData(taskID, targetURI, reqBody, resp, common.Completed, common.OK, percentComplete, http.MethodPost)
	e.UpdateTask(task)
	return resp
}
