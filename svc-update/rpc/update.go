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
	"log"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// GetUpdateService is an rpc handler, it gets invoked during GET on UpdateService API (/redfis/v1/UpdateService/)
func (a *Updater) GetUpdateService(ctx context.Context, req *updateproto.UpdateRequest, resp *updateproto.UpdateResponse) error {
	fillProtoResponse(resp, a.connector.GetUpdateService())
	return nil
}

// GetFirmwareInventoryCollection an rpc handler which is invoked during GET on firmware inventory collection
func (a *Updater) GetFirmwareInventoryCollection(ctx context.Context, req *updateproto.UpdateRequest, resp *updateproto.UpdateResponse) error {
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Println("error while trying to authenticate session")
		fillProtoResponse(resp, authResp)
		return nil
	}
	fillProtoResponse(resp, a.connector.GetAllFirmwareInventory(req))
	return nil
}

// GetFirmwareInventory is an rpc handler which is invoked during GET on firmware inventory
func (a *Updater) GetFirmwareInventory(ctx context.Context, req *updateproto.UpdateRequest, resp *updateproto.UpdateResponse) error {
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Println("error while trying to authenticate session")
		fillProtoResponse(resp, authResp)
		return nil
	}
	fillProtoResponse(resp, a.connector.GetFirmwareInventory(req))
	return nil
}

// GetSoftwareInventoryCollection is an rpc handler which is invoked during GET on software inventory collection
func (a *Updater) GetSoftwareInventoryCollection(ctx context.Context, req *updateproto.UpdateRequest, resp *updateproto.UpdateResponse) error {
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Println("error while trying to authenticate session")
		fillProtoResponse(resp, authResp)
		return nil
	}
	fillProtoResponse(resp, a.connector.GetAllSoftwareInventory(req))
	return nil
}

// GetSoftwareInventory is an rpc handler which is invoked during GET on software inventory
func (a *Updater) GetSoftwareInventory(ctx context.Context, req *updateproto.UpdateRequest, resp *updateproto.UpdateResponse) error {
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Println("error while trying to authenticate session")
		fillProtoResponse(resp, authResp)
		return nil
	}
	fillProtoResponse(resp, a.connector.GetSoftwareInventory(req))
	return nil
}

// SimepleUpdate is an rpc handler, it gets involked during POST on UpdateService API actions (/Actions/UpdateService.SimpleUpdate)
func (a *Updater) SimepleUpdate(ctx context.Context, req *updateproto.UpdateRequest, resp *updateproto.UpdateResponse) error {
	authResp := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Println("error while trying to authenticate session")
		fillProtoResponse(resp, authResp)
		return nil
	}
	sessionUserName, err := a.connector.External.GetSessionUserName(req.SessionToken)
	if err != nil {
		errMsg := "error while trying to get the session username: " + err.Error()
		generateRPCResponse(common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil), resp)
		log.Printf(errMsg)
		return nil
	}
	taskURI, err := a.connector.External.CreateTask(sessionUserName)
	if err != nil {
		errMsg := "error while trying to create task: " + err.Error()
		generateRPCResponse(common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), resp)
		log.Printf(errMsg)
		return nil
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
		log.Printf("error while contacting task-service with UpdateTask RPC : %v", err)
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
	return nil
}

// StartUpdate is an rpc handler, it gets involked during POST on UpdateService API actions (/Actions/UpdateService.StartUpdate)
func (a *Updater) StartUpdate(ctx context.Context, req *updateproto.UpdateRequest, resp *updateproto.UpdateResponse) error {
	sessionToken := req.SessionToken
	authResp := a.connector.External.Auth(sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Println("error while trying to authenticate session")
		fillProtoResponse(resp, authResp)
		return nil
	}
	fillProtoResponse(resp, a.connector.StartUpdate(req))
	return nil
}
