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

package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
	"github.com/ODIM-Project/ODIM/svc-aggregation/system"
)

var podName = os.Getenv("POD_NAME")

// GetAggregationService is an rpc handler, it gets invoked during GET on AggregationService API (/redfis/v1/AggregationService/)
func (a *Aggregator) GetAggregationService(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	resp := &aggregatorproto.AggregatorResponse{}
	// Fill the response header first
	resp.Header = map[string]string{
		"Date": time.Now().Format(http.TimeFormat),
		"Link": "</redfish/v1/SchemaStore/en/AggregationService.json>; rel=describedby",
	}
	// Validate the token, if user has Login priielege then proceed.
	//Else send 401 Unauthorised
	var oemprivileges []string
	privileges := []string{common.PrivilegeLogin}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	// Check whether the Aggregation Service is enbaled in configuration file.
	//If so set ServiceEnabled to true.
	isServiceEnabled := false
	serviceState := "Disabled"
	for _, service := range config.Data.EnabledServices {
		if service == "AggregationService" {
			isServiceEnabled = true
			serviceState = "Enabled"
			break
		}

	}
	// Construct the response below

	aggregationServiceResponse, _ := json.Marshal(agresponse.AggregationServiceResponse{
		OdataType:    common.AggregationServiceType,
		ID:           "AggregationService",
		Name:         "AggregationService",
		Description:  "AggregationService",
		OdataContext: "/redfish/v1/$metadata#AggregationService.AggregationService",
		OdataID:      "/redfish/v1/AggregationService",
		Actions: agresponse.Actions{
			Reset: agresponse.Action{
				Target:     "/redfish/v1/AggregationService/Actions/AggregationService.Reset/",
				ActionInfo: "/redfish/v1/AggregationService/ResetActionInfo",
			},
			SetDefaultBootOrder: agresponse.Action{
				Target:     "/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder/",
				ActionInfo: "/redfish/v1/AggregationService/SetDefaultBootOrderActionInfo",
			},
		},
		Aggregates: agresponse.OdataID{
			OdataID: "/redfish/v1/AggregationService/Aggregates",
		},
		AggregationSources: agresponse.OdataID{
			OdataID: "/redfish/v1/AggregationService/AggregationSources",
		},
		ConnectionMethods: agresponse.OdataID{
			OdataID: "/redfish/v1/AggregationService/ConnectionMethods",
		},
		ServiceEnabled: isServiceEnabled,
		Status: agresponse.Status{
			State:        serviceState,
			HealthRollup: "OK",
			Health:       "OK",
		},
	})
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	resp.Body = aggregationServiceResponse
	l.LogWithFields(ctx).Debugf("final response for get aggregation service request: %s", string(resp.Body))
	return resp, nil
}

func isIpv4Net(host string) bool {
	return net.ParseIP(host) != nil
}

func validateManagerAddress(managerAddress string) error {
	// if the manager address is of the form <IP/FQDN>:<port>
	// will split address to obtain only IP/FQDN. If obtained
	// value is empty, then will use the actual address passed
	addr, _, _ := net.SplitHostPort(managerAddress)
	if addr == "" {
		addr = managerAddress
	}
	if _, err := net.ResolveIPAddr("ip", addr); err != nil {
		return fmt.Errorf("error: failed to resolve ManagerAddress: %v", err)
	}
	return nil
}

// Reset function is for handling the RPC communication for Reset Action
func (a *Aggregator) Reset(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	var threadID int = 1
	// Verfy the credentials here
	var oemprivileges []string
	privileges := []string{common.PrivilegeConfigureComponents}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	sessionUserName, err := a.connector.GetSessionUserName(ctx, req.SessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		generateResponse(common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
	}

	// Task Service using RPC and get the taskID
	taskURI, err := a.connector.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create task: " + err.Error()
		generateResponse(common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
	}
	taskID := strings.TrimPrefix(taskURI, "/redfish/v1/TaskService/Tasks/")
	ctxt := context.WithValue(ctx, common.ThreadName, common.ResetAggregate)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go a.reset(ctxt, taskID, sessionUserName, req)
	threadID++
	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Location": "/taskmon/" + taskID,
		},
	}
	generateTaskRespone(taskID, taskURI, &rpcResp)
	generateResponse(rpcResp, resp)
	l.LogWithFields(ctx).Debugf("final response for reset request: %s", string(resp.Body))
	return resp, nil

}

func (a *Aggregator) reset(ctx context.Context, taskID string, sessionUserName string, req *aggregatorproto.AggregatorRequest) error {
	// Update the task status here
	// PercentComplete: 0% Completed
	// TaskState: Running - This value shall represent that the operation is executing.
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	err := a.connector.UpdateTask(ctx, common.TaskData{
		TaskID:          taskID,
		TaskState:       common.Running,
		TaskStatus:      common.OK,
		PercentComplete: 0,
		HTTPMethod:      http.MethodPost,
	})
	if err != nil && (err.Error() == common.Cancelling) {
		// We cant do anything here as the task has done it work completely, we cant reverse it.
		//Unless if we can do opposite/reverse action for delete server which is add server.
		a.connector.UpdateTask(ctx, common.TaskData{
			TaskID:          taskID,
			TaskState:       common.Cancelled,
			TaskStatus:      common.OK,
			PercentComplete: 0,
			HTTPMethod:      http.MethodPost,
		})
	}

	a.connector.Reset(ctx, taskID, sessionUserName, req)
	return nil
}

// SetDefaultBootOrder defines the operations which handles the RPC request response
// for the create account service of aggregator micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Aggregator) SetDefaultBootOrder(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	var threadID int = 1
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	var oemprivileges []string
	privileges := []string{common.PrivilegeConfigureComponents}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	sessionUserName, err := a.connector.GetSessionUserName(ctx, req.SessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		generateResponse(common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
	}
	taskURI, err := a.connector.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create task: " + err.Error()
		generateResponse(common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
	}
	strArray := strings.Split(taskURI, "/")
	var taskID string
	if strings.HasSuffix(taskURI, "/") {
		taskID = strArray[len(strArray)-2]
	} else {
		taskID = strArray[len(strArray)-1]
	}
	err = a.connector.UpdateTask(ctx, common.TaskData{
		TaskID:          taskID,
		TargetURI:       taskURI,
		TaskState:       common.Running,
		TaskStatus:      common.OK,
		PercentComplete: 0,
		HTTPMethod:      http.MethodPost,
	})
	if err != nil {
		// print error as we are unable to communicate with svc-task and then return
		l.LogWithFields(ctx).Error("Unable to contact task-service with UpdateTask RPC : " + err.Error())
	}
	ctxt := context.WithValue(ctx, common.ThreadName, common.SetBootOrder)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go a.connector.SetDefaultBootOrder(ctxt, taskID, sessionUserName, req)
	threadID++
	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Location": "/taskmon/" + taskID,
		},
	}
	generateTaskRespone(taskID, taskURI, &rpcResp)
	generateResponse(rpcResp, resp)
	l.LogWithFields(ctx).Debugf("final response for set default boot order request: %s", string(resp.Body))
	return resp, nil
}

// RediscoverSystemInventory defines the operations which handles the RPC request response
// for the RediscoverSystemInventory service of aggregator micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
func (a *Aggregator) RediscoverSystemInventory(ctx context.Context, req *aggregatorproto.RediscoverSystemInventoryRequest) (
	*aggregatorproto.RediscoverSystemInventoryResponse, error) {
	resp := &aggregatorproto.RediscoverSystemInventoryResponse{}
	var threadID int = 1
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	ctx = context.WithValue(ctx, common.ThreadID, strconv.Itoa(threadID))
	go a.connector.RediscoverSystemInventory(ctx, req.SystemID, req.SystemURL, true)
	threadID++
	return resp, nil

}

// UpdateSystemState defines the operations which handles the RPC request response
// for the UpdateSystemState call to aggregator micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
func (a *Aggregator) UpdateSystemState(ctx context.Context, req *aggregatorproto.UpdateSystemStateRequest) (
	*aggregatorproto.UpdateSystemStateResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	resp := &aggregatorproto.UpdateSystemStateResponse{}
	return resp, a.connector.UpdateSystemState(ctx, req)
}

// AddAggregationSource function is for handling the RPC communication for AddAggregationSource
func (a *Aggregator) AddAggregationSource(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	var taskID string
	var oemprivileges []string
	privileges := []string{common.PrivilegeConfigureComponents}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	sessionUserName, err := a.connector.GetSessionUserName(ctx, req.SessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		generateResponse(common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
	}
	// parsing the AggregationSourceRequest
	var addRequest system.AggregationSource
	err = json.Unmarshal(req.RequestBody, &addRequest)
	if err != nil {
		errMsg := "Unable to parse the add request" + err.Error()
		generateResponse(common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
	}
	//password in add request, hence cannot log it
	//validating the AggregationSourceRequest
	invalidParam := validateAggregationSourceRequest(addRequest)
	if invalidParam != "" {
		errMsg := "Mandatory field " + invalidParam + " Missing"
		generateResponse(common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{invalidParam}, nil), resp)
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
	}
	managerAddress := addRequest.HostName
	err = validateManagerAddress(managerAddress)
	if err != nil {
		generateResponse(common.GeneralError(http.StatusBadRequest, response.PropertyValueFormatError, err.Error(), []interface{}{managerAddress, "ManagerAddress"}, nil), resp)
		l.LogWithFields(ctx).Error(err.Error())
		return resp, nil
	}

	// Task Service using RPC and get the taskID
	taskURI, err := a.connector.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create the task: " + err.Error()
		generateResponse(common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
	}
	strArray := strings.Split(taskURI, "/")
	if strings.HasSuffix(taskURI, "/") {
		taskID = strArray[len(strArray)-2]
	} else {
		taskID = strArray[len(strArray)-1]
	}
	// spawn the thread here to process the action asynchronously
	threadID := 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.AddAggregationSource)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go a.connector.AddAggregationSource(ctxt, taskID, sessionUserName, req)
	threadID++

	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Location": "/taskmon/" + taskID,
		},
	}
	generateTaskRespone(taskID, taskURI, &rpcResp)
	generateResponse(rpcResp, resp)
	l.LogWithFields(ctx).Debugf("final response for add aggregation source request: %s", string(resp.Body))
	return resp, nil
}

func validateAggregationSourceRequest(req system.AggregationSource) string {
	param := ""
	if req.HostName == "" {
		param = "HostName "
	}
	if req.Password == "" {
		param = param + "Password "
	}
	if req.UserName == "" {
		param = param + "UserName "
	}
	return param + validateLinks(req.Links)
}

func validateLinks(req *system.Links) string {
	var param = ""
	if req != nil {
		if req.ConnectionMethod != nil {
			if req.ConnectionMethod.OdataID == "" {
				param = param + "ConnectionMethod @odata.id"
			}
		}
	} else {
		param = "Links"
	}
	return param
}

// GetAllAggregationSource defines the operations which handles the RPC request response
// for the GetAllAggregationSource  service of aggregation micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Aggregator) GetAllAggregationSource(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	var oemprivileges []string
	privileges := []string{common.PrivilegeConfigureComponents}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	data := a.connector.GetAggregationSourceCollection(ctx)
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header
	generateResponse(data, resp)
	l.LogWithFields(ctx).Debugf("final response for get all aggregation source request: %s", string(resp.Body))
	return resp, nil
}

// GetAggregationSource defines the operations which handles the RPC request response
// for the GetAggregationSource  service of aggregation micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Aggregator) GetAggregationSource(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	var oemprivileges []string
	privileges := []string{common.PrivilegeConfigureComponents}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	data := a.connector.GetAggregationSource(ctx, req.URL)
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header
	generateResponse(data, resp)
	l.LogWithFields(ctx).Debugf("final response for get aggregation source request: %s", string(resp.Body))
	return resp, nil
}

// UpdateAggregationSource defines the operations which handles the RPC request response
// for the UpdateAggregationSource  service of aggregation micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Aggregator) UpdateAggregationSource(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	var oemprivileges []string
	privileges := []string{common.PrivilegeConfigureComponents}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	data := a.connector.UpdateAggregationSource(ctx, req)
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Header = data.Header
	generateResponse(data, resp)
	l.LogWithFields(ctx).Debugf("final response for update aggregation source request: %s", string(resp.Body))
	return resp, nil
}

// DeleteAggregationSource defines the operations which handles the RPC request response
// for the UpdateAggregationSource  service of aggregation micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Aggregator) DeleteAggregationSource(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	var oemprivileges []string
	// Task Service using RPC and get the taskID
	targetURI := req.URL
	privileges := []string{common.PrivilegeConfigureComponents}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	sessionUserName, err := a.connector.GetSessionUserName(ctx, req.SessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		generateResponse(common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
	}

	// Task Service using RPC and get the taskID
	taskURI, err := a.connector.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create task: " + err.Error()
		generateResponse(common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
	}
	var taskID string
	strArray := strings.Split(taskURI, "/")
	if strings.HasSuffix(taskURI, "/") {
		taskID = strArray[len(strArray)-2]
	} else {
		taskID = strArray[len(strArray)-1]
	}
	var threadID int = 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.DeleteAggregationSource)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go a.connector.DeleteAggregationSources(ctxt, taskID, targetURI, req)
	threadID++
	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Location": "/taskmon/" + taskID,
		},
	}
	generateTaskRespone(taskID, taskURI, &rpcResp)
	generateResponse(rpcResp, resp)
	l.LogWithFields(ctx).Debugf("final response for delete aggregation source request: %s", string(resp.Body))
	return resp, nil
}

// CreateAggregate defines the operations which handles the RPC request response
// for the CreateAggregate service of aggregation micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Aggregator) CreateAggregate(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	var oemprivileges []string
	privileges := []string{common.PrivilegeConfigureComponents}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	rpcResponce := a.connector.CreateAggregate(ctx, req)
	generateResponse(rpcResponce, resp)
	l.LogWithFields(ctx).Debugf("final response for create aggregate request: %s", string(resp.Body))
	return resp, nil
}

// GetAllAggregates defines the operations which handles the RPC request response
// for the GetAllAggregates service of aggregation micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Aggregator) GetAllAggregates(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	var oemprivileges []string
	privileges := []string{common.PrivilegeConfigureComponents}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	rpcResponce := a.connector.GetAllAggregates(ctx, req)
	generateResponse(rpcResponce, resp)
	l.LogWithFields(ctx).Debugf("final response for get all aggregates request: %s", string(resp.Body))
	return resp, nil
}

// GetAggregate defines the operations which handles the RPC request response
// for the GetAggregate service of aggregation micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Aggregator) GetAggregate(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	var oemprivileges []string
	privileges := []string{common.PrivilegeConfigureComponents}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	rpcResponce := a.connector.GetAggregate(ctx, req)
	generateResponse(rpcResponce, resp)
	l.LogWithFields(ctx).Debugf("final response for get aggregate request: %s", string(resp.Body))
	return resp, nil
}

// DeleteAggregate defines the operations which handles the RPC request response
// for the DeleteAggregate service of aggregation micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Aggregator) DeleteAggregate(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	var oemprivileges []string
	privileges := []string{common.PrivilegeConfigureComponents}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	rpcResponce := a.connector.DeleteAggregate(ctx, req)
	generateResponse(rpcResponce, resp)
	l.LogWithFields(ctx).Debugf("final response for delete aggregate request: %s", string(resp.Body))
	return resp, nil
}

// AddElementsToAggregate defines the operations which handles the RPC request response
// for the AddElementsToAggregate service of aggregation micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Aggregator) AddElementsToAggregate(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	var oemprivileges []string
	privileges := []string{common.PrivilegeConfigureComponents}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	rpcResponce := a.connector.AddElementsToAggregate(ctx, req)
	generateResponse(rpcResponce, resp)
	l.LogWithFields(ctx).Debugf("final response for add elements to aggregate request: %s", string(resp.Body))
	return resp, nil
}

// RemoveElementsFromAggregate defines the operations which handles the RPC request response
// for the RemoveElementsFromAggregate service of aggregation micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Aggregator) RemoveElementsFromAggregate(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	var oemprivileges []string
	privileges := []string{common.PrivilegeConfigureComponents}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	rpcResponce := a.connector.RemoveElementsFromAggregate(ctx, req)
	generateResponse(rpcResponce, resp)
	l.LogWithFields(ctx).Debugf("final response for remove elements from aggregate request: %s", string(resp.Body))
	return resp, nil
}

// ResetElementsOfAggregate defines the operations which handles the RPC request response
// for the ResetElementsOfAggregate service of aggregation micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Aggregator) ResetElementsOfAggregate(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	// Verfy the credentials here
	var oemprivileges []string
	privileges := []string{common.PrivilegeConfigureComponents}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	sessionUserName, err := a.connector.GetSessionUserName(ctx, req.SessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		generateResponse(common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
	}

	// Task Service using RPC and get the taskID
	taskURI, err := a.connector.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create task: " + err.Error()
		generateResponse(common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
	}
	taskID := strings.TrimPrefix(taskURI, "/redfish/v1/TaskService/Tasks/")

	threadID := 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.ResetSystem)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go a.resetElements(ctxt, taskID, sessionUserName, req)
	threadID++
	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Location": "/taskmon/" + taskID,
		},
	}
	generateTaskRespone(taskID, taskURI, &rpcResp)
	generateResponse(rpcResp, resp)
	l.LogWithFields(ctx).Debugf("final response for reset elements of aggregate request: %s", string(resp.Body))
	return resp, nil
}

func (a *Aggregator) resetElements(ctx context.Context, taskID string, sessionUserName string, req *aggregatorproto.AggregatorRequest) error {
	// Update the task status here
	// PercentComplete: 0% Completed
	// TaskState: Running - This value shall represent that the operation is executing.
	err := a.connector.UpdateTask(ctx, common.TaskData{
		TaskID:          taskID,
		TaskState:       common.Running,
		TaskStatus:      common.OK,
		PercentComplete: 0,
		HTTPMethod:      http.MethodPost,
	})
	if err != nil && (err.Error() == common.Cancelling) {
		// We cant do anything here as the task has done it work completely, we cant reverse it.
		//Unless if we can do opposite/reverse action for delete server which is add server.
		a.connector.UpdateTask(ctx, common.TaskData{
			TaskID:          taskID,
			TaskState:       common.Cancelled,
			TaskStatus:      common.OK,
			PercentComplete: 0,
			HTTPMethod:      http.MethodPost,
		})
	}

	a.connector.ResetElementsOfAggregate(ctx, taskID, sessionUserName, req)
	return nil
}

// SetDefaultBootOrderElementsOfAggregate defines the operations which handles the RPC request response
// for the SetDefaultBootOrderElementsOfAggregate service of aggregation micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Aggregator) SetDefaultBootOrderElementsOfAggregate(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	var oemprivileges []string
	privileges := []string{common.PrivilegeConfigureComponents}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	sessionUserName, err := a.connector.GetSessionUserName(ctx, req.SessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		generateResponse(common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
	}
	taskURI, err := a.connector.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create task: " + err.Error()
		generateResponse(common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
	}
	strArray := strings.Split(taskURI, "/")
	var taskID string
	if strings.HasSuffix(taskURI, "/") {
		taskID = strArray[len(strArray)-2]
	} else {
		taskID = strArray[len(strArray)-1]
	}
	err = a.connector.UpdateTask(ctx, common.TaskData{
		TaskID:          taskID,
		TargetURI:       taskURI,
		TaskState:       common.Running,
		TaskStatus:      common.OK,
		PercentComplete: 0,
		HTTPMethod:      http.MethodPost,
	})
	if err != nil {
		// print error as we are unable to communicate with svc-task and then return
		l.LogWithFields(ctx).Error("Unable to contact task-service with UpdateTask RPC : " + err.Error())
	}

	threadID := 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.SetDefaultBootOrderElementsOfAggregate)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go a.connector.SetDefaultBootOrderElementsOfAggregate(ctxt, taskID, sessionUserName, req)
	threadID++
	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Location": "/taskmon/" + taskID,
		},
	}
	generateTaskRespone(taskID, taskURI, &rpcResp)
	generateResponse(rpcResp, resp)
	l.LogWithFields(ctx).Debugf("final response for set default boot order elements of aggregate request: %s", string(resp.Body))
	return resp, nil
}

// GetAllConnectionMethods defines the operations which handles the RPC request response
// for the GetAllConnectionMethods service of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Aggregator) GetAllConnectionMethods(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	var oemprivileges []string
	privileges := []string{common.PrivilegeLogin}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	rpcResponce := a.connector.GetAllConnectionMethods(ctx, req)
	generateResponse(rpcResponce, resp)
	l.LogWithFields(ctx).Debugf("final response for get all connection methods request: %s", string(resp.Body))
	return resp, nil
}

// GetConnectionMethod defines the operations which handles the RPC request response
// for the GetConnectionMethod service of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function also checks for the session time out of the token
// which is present in the request.
func (a *Aggregator) GetConnectionMethod(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	var oemprivileges []string
	privileges := []string{common.PrivilegeLogin}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	resp := &aggregatorproto.AggregatorResponse{}
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	rpcResponce := a.connector.GetConnectionMethodInfo(ctx, req)
	generateResponse(rpcResponce, resp)
	l.LogWithFields(ctx).Debugf("final response for get connection method request: %s", string(resp.Body))
	return resp, nil
}

// SendStartUpData defines the operations which handles the RPC request response
// for the SendStartUpData call to aggregator micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function is used for sending plugin start up data to the plugin
// which has restarted.
func (a *Aggregator) SendStartUpData(ctx context.Context, req *aggregatorproto.SendStartUpDataRequest) (
	resp *aggregatorproto.SendStartUpDataResponse, err error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	rpcResponce := a.connector.SendStartUpData(ctx, req)
	bytes, _ := json.Marshal(rpcResponce.Body)
	resp = &aggregatorproto.SendStartUpDataResponse{
		ResponseBody: bytes,
	}
	return resp, nil
}

// GetResetActionInfoService is an rpc handler, it gets invoked during GET on AggregationService API (/redfis/v1/AggregationService/)
func (a *Aggregator) GetResetActionInfoService(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	resp := &aggregatorproto.AggregatorResponse{}
	// Fill the response header first
	resp.Header = map[string]string{
		"Date": time.Now().Format(http.TimeFormat),
		"Link": "</redfish/v1/SchemaStore/en/AggregationService.json>; rel=describedby",
	}
	// Validate the token, if user has Login priielege then proceed.
	//Else send 401 Unauthorised
	var oemprivileges []string
	privileges := []string{common.PrivilegeLogin}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	aggregationServiceResponse, _ := json.Marshal(agresponse.ActionInfo{
		ID:        "ResetActionInfo",
		OdataType: "#ActionInfo.v1_2_0.ActionInfo",
		Name:      "Reset Action Info",
		Parameters: []agresponse.Parameter{
			{
				Name:            "ResetType",
				Required:        true,
				DataType:        "String",
				AllowableValues: []string{"On", "ForceOff", "GracefulShutdown", "GracefulRestart", "ForceRestart", "Nmi", "ForceOn", "PushPowerButton"},
			}, {
				Name:     "TargetURIs",
				Required: true,
				DataType: "ObjectArray",
			}, {
				Name:     "BatchSize",
				Required: false,
				DataType: "Number",
			}, {
				Name:     "DelayBetweenBatchesInSeconds",
				Required: false,
				DataType: "Number",
			},
		},
		OdataID: "/redfish/v1/AggregationService/ResetActionInfo",
	})
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	resp.Body = aggregationServiceResponse
	l.LogWithFields(ctx).Debugf("final response for get reset action info service request: %s", string(resp.Body))
	return resp, nil
}

// GetSetDefaultBootOrderActionInfo is an rpc handler, it gets invoked during GET on AggregationService API (/redfis/v1/AggregationService/)
func (a *Aggregator) GetSetDefaultBootOrderActionInfo(ctx context.Context, req *aggregatorproto.AggregatorRequest) (
	*aggregatorproto.AggregatorResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.AggregationService, podName)
	resp := &aggregatorproto.AggregatorResponse{}
	// Fill the response header first
	resp.Header = map[string]string{
		"Date": time.Now().Format(http.TimeFormat),
		"Link": "</redfish/v1/SchemaStore/en/AggregationService.json>; rel=describedby",
	}
	// Validate the token, if user has Login priielege then proceed.
	//Else send 401 Unauthorised
	var oemprivileges []string
	privileges := []string{common.PrivilegeLogin}
	authResp, err := a.connector.Auth(ctx, req.SessionToken, privileges, oemprivileges)
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		generateResponse(authResp, resp)
		return resp, nil
	}
	setDefaultBootOrderActionInfoResponse, _ := json.Marshal(agresponse.ActionInfo{
		ID:        "SetDefaultBootOrderActionInfo",
		OdataType: "#ActionInfo.v1_2_0.ActionInfo",
		Name:      "SetDefaultBootOrder Action Info",
		Parameters: []agresponse.Parameter{
			{
				Name:     "Systems",
				Required: true,
				DataType: "ObjectArray",
			},
		},
		OdataID: "/redfish/v1/AggregationService/SetDefaultBootOrderActionInfo",
	})
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	resp.Body = setDefaultBootOrderActionInfoResponse
	l.LogWithFields(ctx).Debugf("final response for get set default boot order action info request: %s", string(resp.Body))
	return resp, nil
}
