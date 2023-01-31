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
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// podName defines the current name of process
var podName = os.Getenv("POD_NAME")

// GetUpdateService is an rpc handler, it gets invoked during GET on UpdateService API (/redfis/v1/UpdateService/)
func (a *Updater) GetUpdateService(ctx context.Context, req *updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.UpdateService, podName)
	l.LogWithFields(ctx).Info("Inside GetUpdateService function (svc-update)")
	resp := &updateproto.UpdateResponse{}
	fillProtoResponse(ctx, resp, a.connector.GetUpdateService(ctx))
	return resp, nil
}

// GetFirmwareInventoryCollection an rpc handler which is invoked during GET on firmware inventory collection
func (a *Updater) GetFirmwareInventoryCollection(ctx context.Context, req *updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.UpdateService, podName)
	l.LogWithFields(ctx).Info("Inside GetFirmwareInventoryCollection function (svc-update)")
	resp := &updateproto.UpdateResponse{}
	authResp, err := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, a.connector.GetAllFirmwareInventory(ctx, req))
	return resp, nil
}

// GetFirmwareInventory is an rpc handler which is invoked during GET on firmware inventory
func (a *Updater) GetFirmwareInventory(ctx context.Context, req *updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.UpdateService, podName)
	l.LogWithFields(ctx).Info("Inside GetFirmwareInventory function (svc-update)")
	resp := &updateproto.UpdateResponse{}
	authResp, err := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, a.connector.GetFirmwareInventory(ctx, req))
	return resp, nil
}

// GetSoftwareInventoryCollection is an rpc handler which is invoked during GET on software inventory collection
func (a *Updater) GetSoftwareInventoryCollection(ctx context.Context, req *updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.UpdateService, podName)
	l.LogWithFields(ctx).Info("Inside GetSoftwareInventoryCollection function (svc-update)")
	resp := &updateproto.UpdateResponse{}
	authResp, err := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, a.connector.GetAllSoftwareInventory(ctx, req))
	return resp, nil
}

// GetSoftwareInventory is an rpc handler which is invoked during GET on software inventory
func (a *Updater) GetSoftwareInventory(ctx context.Context, req *updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.UpdateService, podName)
	l.LogWithFields(ctx).Info("Inside GetSoftwareInventory function (svc-update)")
	resp := &updateproto.UpdateResponse{}
	authResp, err := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, a.connector.GetSoftwareInventory(ctx, req))
	return resp, nil
}

// SimepleUpdate is an rpc handler, it gets involked during POST on UpdateService API actions (/Actions/UpdateService.SimpleUpdate)
func (a *Updater) SimepleUpdate(ctx context.Context, req *updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.UpdateService, podName)
	l.LogWithFields(ctx).Info("Inside SimepleUpdate function (svc-update)")
	resp := &updateproto.UpdateResponse{}
	authResp, err := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	sessionUserName, err := a.connector.External.GetSessionUserName(req.SessionToken)
	if err != nil {
		errMsg := "error while trying to get the session username: " + err.Error()
		generateRPCResponse(common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Warn(errMsg)
		return resp, nil
	}
	taskURI, err := a.connector.External.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "error while trying to create task: " + err.Error()
		generateRPCResponse(common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Warn(errMsg)
		return resp, nil
	}
	strArray := strings.Split(taskURI, "/")
	var taskID string
	if strings.HasSuffix(taskURI, "/") {
		taskID = strArray[len(strArray)-2]
	} else {
		taskID = strArray[len(strArray)-1]
	}
	err = a.connector.External.UpdateTask(ctx, common.TaskData{
		TaskID:          taskID,
		TargetURI:       taskURI,
		TaskState:       common.Running,
		TaskStatus:      common.OK,
		PercentComplete: 0,
		HTTPMethod:      http.MethodPost,
	})
	if err != nil {
		l.LogWithFields(ctx).Warn("error while contacting task-service with UpdateTask RPC : " + err.Error())
	}
	var threadID int = 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.SimpleUpdate)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go a.connector.SimpleUpdate(ctxt, taskID, sessionUserName, req)
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
	generateRPCResponse(rpcResp, resp)
	//fillProtoResponse(ctx, resp, a.connector.SimpleUpdate(req))
	return resp, nil
}

// StartUpdate is an rpc handler, it gets involked during POST on UpdateService API actions (/Actions/UpdateService.StartUpdate)
func (a *Updater) StartUpdate(ctx context.Context, req *updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.UpdateService, podName)
	l.LogWithFields(ctx).Info("Inside StartUpdate function (svc-update)")
	resp := &updateproto.UpdateResponse{}
	sessionToken := req.SessionToken
	authResp, err := a.connector.External.Auth(sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	sessionUserName, err := a.connector.External.GetSessionUserName(req.SessionToken)
	if err != nil {
		errMsg := "error while trying to get the session username: " + err.Error()
		generateRPCResponse(common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Warn(errMsg)
		return resp, nil
	}
	taskURI, err := a.connector.External.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "error while trying to create task: " + err.Error()
		generateRPCResponse(common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), resp)
		l.LogWithFields(ctx).Warn(errMsg)
		return resp, nil
	}
	strArray := strings.Split(taskURI, "/")
	var taskID string
	if strings.HasSuffix(taskURI, "/") {
		taskID = strArray[len(strArray)-2]
	} else {
		taskID = strArray[len(strArray)-1]
	}
	err = a.connector.External.UpdateTask(ctx, common.TaskData{
		TaskID:          taskID,
		TargetURI:       taskURI,
		TaskState:       common.Running,
		TaskStatus:      common.OK,
		PercentComplete: 0,
		HTTPMethod:      http.MethodPost,
	})
	if err != nil {
		// print error as we are unable to communicate with svc-task and then return
		l.LogWithFields(ctx).Warn("error while contacting task-service with UpdateTask RPC : " + err.Error())
	}
	var threadID int = 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.StartUpdate)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go a.connector.StartUpdate(ctx, taskID, sessionUserName, req)
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
	generateRPCResponse(rpcResp, resp)
	//fillProtoResponse(ctx, resp, a.connector.StartUpdate(req))
	return resp, nil
}
