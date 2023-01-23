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

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
)

// AddAggregationSource is the handler for adding bmc or manager
// Discovers  its top level odata.ID links and store them in inmemory db.
// Upon successfull operation this api returns added AggregationSourceUUID  in the response body with 201 OK.
func (e *ExternalInterface) AddAggregationSource(ctx context.Context, taskID string, sessionUserName string, req *aggregatorproto.AggregatorRequest) response.RPC {
	targetURI := "/redfish/v1/AggregationService/AggregationSources"
	var resp response.RPC
	var percentComplete int32
	var task = fillTaskData(taskID, targetURI, string(req.RequestBody), resp, common.Running, common.OK, percentComplete, http.MethodPost)
	err := e.UpdateTask(ctx, task)
	if err != nil {
		errMsg := "error while starting the task: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	taskInfo := &common.TaskUpdateInfo{TaskID: taskID, TargetURI: targetURI, UpdateTask: e.UpdateTask, TaskRequest: string(req.RequestBody)}
	// parsing the request
	var aggregationSourceRequest AggregationSource
	err = json.Unmarshal(req.RequestBody, &aggregationSourceRequest)
	if err != nil {
		errMsg := "unable to parse the add request" + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	}
	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, aggregationSourceRequest)
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		l.LogWithFields(ctx).Error(errorMessage)
		resp := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, taskInfo)
		return resp
	}

	if aggregationSourceRequest.Links == nil || aggregationSourceRequest.Links.ConnectionMethod == nil {
		errMsg := "error: mandatory ConnectionMethod block missing in the request"
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"ConnectionMethod"}, taskInfo)
	}
	return e.addAggregationSource(ctx, taskID, targetURI, string(req.RequestBody), percentComplete, aggregationSourceRequest, taskInfo)
}

func (e *ExternalInterface) addAggregationSource(ctx context.Context, taskID, targetURI, reqBody string, percentComplete int32, aggregationSourceRequest AggregationSource, taskInfo *common.TaskUpdateInfo) response.RPC {
	var resp response.RPC
	var addResourceRequest = AddResourceRequest{
		ManagerAddress:   aggregationSourceRequest.HostName,
		UserName:         aggregationSourceRequest.UserName,
		Password:         aggregationSourceRequest.Password,
		ConnectionMethod: aggregationSourceRequest.Links.ConnectionMethod,
	}

	ipAddr := getKeyFromManagerAddress(addResourceRequest.ManagerAddress)
	indexList, err := agmodel.GetString("BMCAddress", ipAddr)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to collect the active request details from DB: %v", err.Error())
		l.LogWithFields(ctx).Errorln(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	}
	if len(indexList) > 0 {
		errMsg := fmt.Sprintf("Manager address already exist %v", ipAddr)
		return common.GeneralError(http.StatusConflict, response.ResourceAlreadyExists, errMsg, []interface{}{"ComputerSystem", "HostName", ipAddr}, taskInfo)
	}
	exist, dErr := e.CheckActiveRequest(ipAddr)
	if dErr != nil {
		errMsg := fmt.Sprintf("Unable to collect the active request details from DB: %v", dErr.Error())
		l.LogWithFields(ctx).Errorln(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	}
	if exist {
		var errMsg string
		mIP, _ := getIPAndPortFromAddress(addResourceRequest.ManagerAddress)
		errMsg = fmt.Sprintf("An active request already exists for adding aggregation source IP %v", mIP)
		l.LogWithFields(ctx).Error(errMsg)
		args := response.Args{
			Code:    response.GeneralError,
			Message: errMsg,
		}
		resp.Body = args.CreateGenericErrorResponse()
		resp.StatusCode = http.StatusConflict
		percentComplete = 100
		e.UpdateTask(ctx, fillTaskData(taskID, targetURI, reqBody, resp, common.Exception, common.Warning, percentComplete, http.MethodPost))
		return resp
	}
	err = e.GenericSave(nil, "ActiveAddBMCRequest", ipAddr)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to save the active request details from DB: %v", err.Error())
		l.LogWithFields(ctx).Errorln(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	}

	defer func() {
		err := e.DeleteActiveRequest(ipAddr)
		if err != nil {
			l.LogWithFields(ctx).Infof("Unable to collect the active request details from DB: %v", err.Error())
		}
	}()

	connectionMethod, err1 := e.GetConnectionMethod(addResourceRequest.ConnectionMethod.OdataID)
	if err1 != nil {
		errMsg := "Unable to get connection method id: " + err1.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"connectionmethod id", addResourceRequest.ConnectionMethod.OdataID}, taskInfo)
	}
	cmVariants := getConnectionMethodVariants(connectionMethod.ConnectionMethodVariant)
	var pluginContactRequest getResourceRequest
	pluginContactRequest.ContactClient = e.ContactClient
	pluginContactRequest.GetPluginStatus = e.GetPluginStatus
	pluginContactRequest.TargetURI = targetURI
	pluginContactRequest.UpdateTask = e.UpdateTask
	pluginContactRequest.TaskRequest = reqBody
	var aggregationSourceUUID string
	var cipherText []byte

	// check status will do call on the URI /ODIM/v1/Status to the requested manager address
	// if its success then add the plugin, else if its not found then add BMC
	// else return the response
	statusResp, statusCode, queueList := checkStatus(ctx, pluginContactRequest, addResourceRequest, cmVariants, taskInfo)
	if statusCode == http.StatusOK {

		// check if AggregationSource has any values, if its there means its managing the bmcs
		if len(connectionMethod.Links.AggregationSources) > 0 {
			errMsg := "Cant proceed to add aggregation source, since connection method is already managing other aggregation sources"
			l.LogWithFields(ctx).Error(errMsg)
			return common.GeneralError(http.StatusConflict, response.ResourceInUse, errMsg, nil, taskInfo)
		}
		resp, aggregationSourceUUID, cipherText = e.addPluginData(ctx, addResourceRequest, taskID, targetURI, pluginContactRequest, queueList, cmVariants)
	} else if statusCode == http.StatusNotFound {
		resp, aggregationSourceUUID, cipherText = e.addCompute(ctx, taskID, targetURI, cmVariants.PluginID, percentComplete, addResourceRequest, pluginContactRequest)
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
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	}

	connectionMethod.Links.AggregationSources = append(connectionMethod.Links.AggregationSources, agmodel.OdataID{OdataID: aggregationSourceURI})
	dbErr = e.UpdateConnectionMethod(connectionMethod, addResourceRequest.ConnectionMethod.OdataID)
	if dbErr != nil {
		errMsg := dbErr.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo)
	}
	commonResponse := response.Response{
		OdataType:    common.AggregationSourceType,
		OdataID:      aggregationSourceURI,
		OdataContext: "/redfish/v1/$metadata#AggregationSource.AggregationSource",
		ID:           aggregationSourceUUID,
		Name:         "Aggregation Source",
	}
	resp.Header = map[string]string{
		"Link":     "<" + aggregationSourceURI + "/>; rel=describedby",
		"Location": aggregationSourceURI,
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
	e.UpdateTask(ctx, task)
	return resp
}
