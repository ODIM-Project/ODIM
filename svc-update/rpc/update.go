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
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// SESSAUTHFAILED string constant to raise errors
const SESSAUTHFAILED string = "Unable to authenticate session"

// GetUpdateService is an rpc handler, it gets invoked during GET on UpdateService API (/redfis/v1/UpdateService/)
func (a *Updater) GetUpdateService(ctx context.Context, req *updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	resp := &updateproto.UpdateResponse{}
	fillProtoResponse(resp, a.connector.GetUpdateService())
	return resp, nil
}

// GetFirmwareInventoryCollection an rpc handler which is invoked during GET on firmware inventory collection
func (a *Updater) GetFirmwareInventoryCollection(ctx context.Context, req *updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	resp := &updateproto.UpdateResponse{}
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Warn(SESSAUTHFAILED)
		fillProtoResponse(resp, authResp)
		return resp, nil
	}
	fillProtoResponse(resp, a.connector.GetAllFirmwareInventory(req))
	return resp, nil
}

// GetFirmwareInventory is an rpc handler which is invoked during GET on firmware inventory
func (a *Updater) GetFirmwareInventory(ctx context.Context, req *updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	resp := &updateproto.UpdateResponse{}
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Warn(SESSAUTHFAILED)
		fillProtoResponse(resp, authResp)
		return resp, nil
	}
	fillProtoResponse(resp, a.connector.GetFirmwareInventory(req))
	return resp, nil
}

// GetSoftwareInventoryCollection is an rpc handler which is invoked during GET on software inventory collection
func (a *Updater) GetSoftwareInventoryCollection(ctx context.Context, req *updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	resp := &updateproto.UpdateResponse{}
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Warn(SESSAUTHFAILED)
		fillProtoResponse(resp, authResp)
		return resp, nil
	}
	fillProtoResponse(resp, a.connector.GetAllSoftwareInventory(req))
	return resp, nil
}

// GetSoftwareInventory is an rpc handler which is invoked during GET on software inventory
func (a *Updater) GetSoftwareInventory(ctx context.Context, req *updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	resp := &updateproto.UpdateResponse{}
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Warn(SESSAUTHFAILED)
		fillProtoResponse(resp, authResp)
		return resp, nil
	}
	fillProtoResponse(resp, a.connector.GetSoftwareInventory(req))
	return resp, nil
}

// SimepleUpdate is an rpc handler, it gets involked during POST on UpdateService API actions (/Actions/UpdateService.SimpleUpdate)
func (a *Updater) SimepleUpdate(ctx context.Context, req *updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {
	resp := &updateproto.UpdateResponse{}
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Warn(SESSAUTHFAILED)
		fillProtoResponse(resp, authResp)
		return resp, nil
	}
	sessionUserName, err := a.connector.External.GetSessionUserName(req.SessionToken)
	if err != nil {
		errMsg := "error while trying to get the session username: " + err.Error()
		generateRPCResponse(common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil), resp)
		log.Warn(errMsg)
		return resp, nil
	}
	taskURI, err := a.connector.External.CreateTask(sessionUserName)
	if err != nil {
		errMsg := "error while trying to create task: " + err.Error()
		generateRPCResponse(common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), resp)
		log.Warn(errMsg)
		return resp, nil
	}
	strArray := strings.Split(taskURI, "/")
	var taskID string
	if strings.HasSuffix(taskURI, "/") {
		taskID = strArray[len(strArray)-2]
	} else {
		taskID = strArray[len(strArray)-1]
	}
	err = a.connector.External.UpdateTask(common.TaskData{
		TaskID:          taskID,
		TargetURI:       taskURI,
		TaskState:       common.Running,
		TaskStatus:      common.OK,
		PercentComplete: 0,
		HTTPMethod:      http.MethodPost,
	})
	if err != nil {
		log.Warn("error while contacting task-service with UpdateTask RPC : " + err.Error())
	}
	go a.connector.SimpleUpdate(taskID, sessionUserName, req)
	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Content-type": "application/json; charset=utf-8",
			"Location":     "/taskmon/" + taskID,
		},
	}
	generateTaskRespone(taskID, taskURI, &rpcResp)
	generateRPCResponse(rpcResp, resp)
	//fillProtoResponse(resp, a.connector.SimpleUpdate(req))
	return resp, nil
}

// StartUpdate is an rpc handler, it gets involked during POST on UpdateService API actions (/Actions/UpdateService.StartUpdate)
func (a *Updater) StartUpdate(ctx context.Context, req *updateproto.UpdateRequest) (*updateproto.UpdateResponse, error) {

	resp := &updateproto.UpdateResponse{}
	sessionToken := req.SessionToken
	authResp := a.connector.External.Auth(sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Warn("Unable to authenticate session")
		fillProtoResponse(resp, authResp)
		return resp, nil
	}
	sessionUserName, err := a.connector.External.GetSessionUserName(req.SessionToken)
	if err != nil {
		errMsg := "error while trying to get the session username: " + err.Error()
		generateRPCResponse(common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil), resp)
		log.Warn(errMsg)
		return resp, nil
	}
	taskURI, err := a.connector.External.CreateTask(sessionUserName)
	if err != nil {
		errMsg := "error while trying to create task: " + err.Error()
		generateRPCResponse(common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), resp)
		log.Warn(errMsg)
		return resp, nil
	}
	strArray := strings.Split(taskURI, "/")
	var taskID string
	if strings.HasSuffix(taskURI, "/") {
		taskID = strArray[len(strArray)-2]
	} else {
		taskID = strArray[len(strArray)-1]
	}
	err = a.connector.External.UpdateTask(common.TaskData{
		TaskID:          taskID,
		TargetURI:       taskURI,
		TaskState:       common.Running,
		TaskStatus:      common.OK,
		PercentComplete: 0,
		HTTPMethod:      http.MethodPost,
	})
	if err != nil {
		// print error as we are unable to communicate with svc-task and then return
		log.Warn("error while contacting task-service with UpdateTask RPC : " + err.Error())
	}
	go a.connector.StartUpdate(taskID, sessionUserName, req)
	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Content-type": "application/json; charset=utf-8",
			"Location":     "/taskmon/" + taskID,
		},
	}
	generateTaskRespone(taskID, taskURI, &rpcResp)
	generateRPCResponse(rpcResp, resp)
	//fillProtoResponse(resp, a.connector.StartUpdate(req))
	return resp, nil
}
