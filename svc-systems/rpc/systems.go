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

// Package rpc ...
package rpc

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-systems/scommon"
	"github.com/ODIM-Project/ODIM/svc-systems/systems"
)

// Systems struct helps to register service
type Systems struct {
	IsAuthorizedRPC    func(ctx context.Context, sessionToken string, privileges, oemPrivileges []string) (response.RPC, error)
	GetSessionUserName func(context.Context, string) (string, error)
	CreateTask         func(ctx context.Context, sessionUserName string) (string, error)
	UpdateTask         func(ctx context.Context, task common.TaskData) error
	EI                 *systems.ExternalInterface
}

// GetSystemResource defines the operations which handles the RPC request response
// for the getting the system resource  of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (s *Systems) GetSystemResource(ctx context.Context, req *systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	l.LogWithFields(ctx).Debugf("incoming GetSystemResource request with %s", req.URL)
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp, err := s.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillSystemProtoResponse(ctx, &resp, authResp)
		return &resp, nil
	}
	var pc = systems.PluginContact{
		ContactClient:   pmbhandle.ContactPlugin,
		DevicePassword:  common.DecryptWithPrivateKey,
		GetPluginStatus: scommon.GetPluginStatus,
	}
	data := pc.GetSystemResource(ctx, req)
	fillSystemProtoResponse(ctx, &resp, data)
	return &resp, nil
}

// GetSystemsCollection defines the operation which has the RPC request
// for getting the systems data from odimra.
// Retrieves all the keys with table name systems collection and create the response
// to send back to requested user.
func (s *Systems) GetSystemsCollection(ctx context.Context, req *systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	ctx = context.WithValue(ctx, common.ThreadName, common.SystemService)
	l.LogWithFields(ctx).Debugf("incoming GetSystemsCollection request with %s", req.URL)
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp, err := s.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillSystemProtoResponse(ctx, &resp, authResp)
		return &resp, nil
	}
	data := systems.GetSystemsCollection(ctx, req)
	fillSystemProtoResponse(ctx, &resp, data)
	l.LogWithFields(ctx).Debugf("outgoing response for Get SystemsCollection : %s", string(resp.Body))
	return &resp, nil
}

// GetSystems defines the operations which handles the RPC request response
// for the getting the system resource  of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (s *Systems) GetSystems(ctx context.Context, req *systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	l.LogWithFields(ctx).Debugf("incoming GetSystems request with %s", req.URL)
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp, err := s.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillSystemProtoResponse(ctx, &resp, authResp)
		return &resp, nil
	}
	var pc = systems.PluginContact{
		ContactClient:   pmbhandle.ContactPlugin,
		DevicePassword:  common.DecryptWithPrivateKey,
		GetPluginStatus: scommon.GetPluginStatus,
	}
	data := pc.GetSystems(ctx, req)
	fillSystemProtoResponse(ctx, &resp, data)
	l.LogWithFields(ctx).Debugf("outgoing response for GetSystems : %s", string(resp.Body))
	return &resp, nil
}

// ComputerSystemReset defines the operations which handles the RPC request response
// for the ComputerSystemReset service of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function also checks for the session time out of the token
// which is present in the request.
func (s *Systems) ComputerSystemReset(ctx context.Context, req *systemsproto.ComputerSystemResetRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	l.LogWithFields(ctx).Debugln("incoming ComputerSystemReset request")
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp, err := s.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillSystemProtoResponse(ctx, &resp, authResp)
		return &resp, nil
	}
	sessionUserName, err := s.GetSessionUserName(ctx, req.SessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		fillSystemProtoResponse(ctx, &resp, common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return &resp, nil
	}

	// Task Service using RPC and get the taskID
	taskURI, err := s.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create task: " + err.Error()
		fillSystemProtoResponse(ctx, &resp, common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return &resp, nil
	}
	taskID := strings.TrimPrefix(taskURI, "/redfish/v1/TaskService/Tasks/")
	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Location": "/taskmon/" + taskID,
		},
	}
	generateTaskRespone(taskID, taskURI, &rpcResp)
	l.LogWithFields(ctx).Debugf("response from generateTaskRespone for id: %s , URI: %s , Is.. Response: %s ", string(taskID), string(taskURI), rpcResp.Body)
	fillSystemProtoResponse(ctx, &resp, rpcResp)
	var pc = systems.PluginContact{
		ContactClient:  pmbhandle.ContactPlugin,
		DevicePassword: common.DecryptWithPrivateKey,
		UpdateTask:     s.UpdateTask,
	}
	var threadID int = 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.ComputerSystemReset)
	ctx = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go pc.ComputerSystemReset(ctx, req, taskID, sessionUserName)
	threadID++
	l.LogWithFields(ctx).Debugf("outgoing response for ComputerSystemReset: %s", string(resp.Body))
	return &resp, nil
}

// SetDefaultBootOrder defines the operations which handles the RPC request response
// for the SetDefaultBootOrder service of systems micro service.
// The functionality retrieves the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function also checks for the session time out of the token
// which is present in the request.
func (s *Systems) SetDefaultBootOrder(ctx context.Context, req *systemsproto.DefaultBootOrderRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	l.LogWithFields(ctx).Debugf("incoming SetDefaultBootOrder request")
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp, err := s.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillSystemProtoResponse(ctx, &resp, authResp)
		return &resp, nil
	}
	var pc = systems.PluginContact{
		ContactClient:      pmbhandle.ContactPlugin,
		DevicePassword:     common.DecryptWithPrivateKey,
		SavePluginTaskInfo: services.SavePluginTaskInfo,
	}
	data := pc.SetDefaultBootOrder(ctx, req.SystemID)
	fillSystemProtoResponse(ctx, &resp, data)
	l.LogWithFields(ctx).Debugf("outgoing response for SetDefaultBootOrder : %s", string(resp.Body))
	return &resp, nil
}

// ChangeBiosSettings defines the operations which handles the RPC request response
// for the ChangeBiosSettings service of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function also checks for the session time out of the token
// which is present in the request.
func (s *Systems) ChangeBiosSettings(ctx context.Context, req *systemsproto.BiosSettingsRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	l.LogWithFields(ctx).Debugf("incoming ChangeBiosSettings request")
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp, err := s.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillSystemProtoResponse(ctx, &resp, authResp)
		return &resp, nil
	}
	var pc = systems.PluginContact{
		ContactClient:      pmbhandle.ContactPlugin,
		DevicePassword:     common.DecryptWithPrivateKey,
		UpdateTask:         s.UpdateTask,
		SavePluginTaskInfo: services.SavePluginTaskInfo,
	}
	sessionUserName, err := s.GetSessionUserName(ctx, req.SessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		fillSystemProtoResponse(ctx, &resp, common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return &resp, nil
	}
	// Task Service using RPC and get the taskID
	taskURI, err := s.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create task: " + err.Error()
		fillSystemProtoResponse(ctx, &resp, common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return &resp, nil
	}
	taskID := strings.TrimPrefix(taskURI, "/redfish/v1/TaskService/Tasks/")
	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Location": "/taskmon/" + taskID,
		},
	}
	generateTaskRespone(taskID, taskURI, &rpcResp)
	fillSystemProtoResponse(ctx, &resp, rpcResp)
	go pc.ChangeBiosSettings(ctx, req, taskID)
	l.LogWithFields(ctx).Debugf("outgoing response for ChangeBiosSettings : %s", string(resp.Body))
	return &resp, nil
}

// ChangeBootOrderSettings defines the operations which handles the RPC request response
// for the ChangeBootOrderSettings service of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function also checks for the session time out of the token
// which is present in the request.
func (s *Systems) ChangeBootOrderSettings(ctx context.Context, req *systemsproto.BootOrderSettingsRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	l.LogWithFields(ctx).Debugf("incoming ChangeBootOrderSettings request for SystemID: %s", req.SystemID)
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp, err := s.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillSystemProtoResponse(ctx, &resp, authResp)
		return &resp, nil
	}
	sessionUserName, err := s.GetSessionUserName(ctx, req.SessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		fillSystemProtoResponse(ctx, &resp, common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return &resp, nil
	}

	// Task Service using RPC and get the taskID
	taskURI, err := s.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create task: " + err.Error()
		fillSystemProtoResponse(ctx, &resp, common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return &resp, nil
	}
	taskID := strings.TrimPrefix(taskURI, "/redfish/v1/TaskService/Tasks/")
	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Location": "/taskmon/" + taskID,
		},
	}
	generateTaskRespone(taskID, taskURI, &rpcResp)
	l.LogWithFields(ctx).Debugf("response from generateTaskRespone for id: %s , URI: %s , Is.. Response: %s ", string(taskID), string(taskURI), rpcResp.Body)
	fillSystemProtoResponse(ctx, &resp, rpcResp)

	var pc = systems.PluginContact{
		ContactClient:      pmbhandle.ContactPlugin,
		DevicePassword:     common.DecryptWithPrivateKey,
		UpdateTask:         s.UpdateTask,
		SavePluginTaskInfo: services.SavePluginTaskInfo,
	}

	var threadID int = 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.ChangeBootOrderSettings)
	ctx = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go pc.ChangeBootOrderSettings(ctx, req, taskID)
	threadID++
	l.LogWithFields(ctx).Debugf("outgoing response for ChangeBootOrderSettings : %s", string(resp.Body))
	return &resp, nil
}

// CreateVolume defines the operations which handles the RPC request response
// for the CreateVolume service of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function also checks for the session time out of the token
// which is present in the request.
func (s *Systems) CreateVolume(ctx context.Context, req *systemsproto.VolumeRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	l.LogWithFields(ctx).Debugf("incoming CreateVolume request")
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp, err := s.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillSystemProtoResponse(ctx, &resp, authResp)
		return &resp, nil
	}
	sessionUserName, err := s.GetSessionUserName(ctx, req.SessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		fillSystemProtoResponse(ctx, &resp, common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return &resp, nil
	}

	// Task Service using RPC and get the taskID
	taskURI, err := s.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create task: " + err.Error()
		fillSystemProtoResponse(ctx, &resp, common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return &resp, nil
	}
	taskID := strings.TrimPrefix(taskURI, "/redfish/v1/TaskService/Tasks/")
	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Location": "/taskmon/" + taskID,
		},
	}
	generateTaskRespone(taskID, taskURI, &rpcResp)
	fillSystemProtoResponse(ctx, &resp, rpcResp)
	var pc = systems.PluginContact{
		ContactClient:      pmbhandle.ContactPlugin,
		DevicePassword:     common.DecryptWithPrivateKey,
		UpdateTask:         s.UpdateTask,
		SavePluginTaskInfo: services.SavePluginTaskInfo,
	}

	var threadID int = 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.CreateVolume)
	ctx = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go s.EI.CreateVolume(ctx, req, &pc, taskID)
	threadID++
	l.LogWithFields(ctx).Debugf("outgoing response for CreateVolume: %s", string(resp.Body))
	return &resp, nil
}

// DeleteVolume defines the operations which handles the RPC request response
// for the DeleteVolume service of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function also checks for the session time out of the token
// which is present in the request.
func (s *Systems) DeleteVolume(ctx context.Context, req *systemsproto.VolumeRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	l.LogWithFields(ctx).Debugf("incoming DeleteVolume request")
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp, err := s.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillSystemProtoResponse(ctx, &resp, authResp)
		return &resp, nil
	}
	sessionUserName, err := s.GetSessionUserName(ctx, req.SessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		fillSystemProtoResponse(ctx, &resp, common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return &resp, nil
	}
	// Task Service using RPC and get the taskID
	taskURI, err := s.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create task: " + err.Error()
		fillSystemProtoResponse(ctx, &resp, common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return &resp, nil
	}
	taskID := strings.TrimPrefix(taskURI, "/redfish/v1/TaskService/Tasks/")
	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Location": "/taskmon/" + taskID,
		},
	}
	generateTaskRespone(taskID, taskURI, &rpcResp)
	fillSystemProtoResponse(ctx, &resp, rpcResp)
	var pc = systems.PluginContact{
		ContactClient:      pmbhandle.ContactPlugin,
		DevicePassword:     common.DecryptWithPrivateKey,
		UpdateTask:         s.UpdateTask,
		SavePluginTaskInfo: services.SavePluginTaskInfo,
	}

	var threadID int = 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.DeleteVolume)
	ctx = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go s.EI.DeleteVolume(ctx, req, &pc, taskID)
	threadID++
	l.LogWithFields(ctx).Debugf("outgoing response DeleteVolume: %s", string(resp.Body))
	return &resp, nil
}

// UpdateSecureBoot defines the operations which handles the RPC request response
// for updating SecureBoot service of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function also checks for the session time out of the token
// which is present in the request.
func (s *Systems) UpdateSecureBoot(ctx context.Context, req *systemsproto.SecureBootRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	l.LogWithFields(ctx).Debugf("incoming request to update SecureBoot")
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp, err := s.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillSystemProtoResponse(ctx, &resp, authResp)
		return &resp, nil
	}
	sessionUserName, err := s.GetSessionUserName(ctx, req.SessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		fillSystemProtoResponse(ctx, &resp, common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return &resp, nil
	}
	// Task Service using RPC and get the taskID
	taskURI, err := s.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create task: " + err.Error()
		fillSystemProtoResponse(ctx, &resp, common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return &resp, nil
	}
	taskID := strings.TrimPrefix(taskURI, "/redfish/v1/TaskService/Tasks/")
	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Location": "/taskmon/" + taskID,
		},
	}
	generateTaskRespone(taskID, taskURI, &rpcResp)
	fillSystemProtoResponse(ctx, &resp, rpcResp)
	var pc = systems.PluginContact{
		ContactClient:      pmbhandle.ContactPlugin,
		DevicePassword:     common.DecryptWithPrivateKey,
		UpdateTask:         s.UpdateTask,
		SavePluginTaskInfo: services.SavePluginTaskInfo,
	}

	var threadID int = 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.UpdateSecureBoot)
	ctx = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go s.EI.UpdateSecureBoot(ctx, req, &pc, taskID)
	threadID++
	l.LogWithFields(ctx).Debugf("outgoing response for UpdateSecureBoot: %s", string(resp.Body))
	return &resp, nil
}

// ResetSecureBoot defines the operations which handles the RPC request response
// for resetting SecureBoot keys of systems.
// The `ResetAllKeysToDefault` value shall reset all UEFI Secure Boot key databases to their default values.
// The `DeleteAllKeys` value shall delete the content of all UEFI Secure Boot key databases.
// The `DeletePK` value shall delete the content of the PK Secure Boot key database.
// which is present in the request.
func (s *Systems) ResetSecureBoot(ctx context.Context, req *systemsproto.SecureBootRequest) (*systemsproto.SystemsResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.SystemService, podName)
	l.LogWithFields(ctx).Debugf("incoming request to reset SecureBoot keys")
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken

	authResp, err := s.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillSystemProtoResponse(ctx, &resp, authResp)
		return &resp, nil
	}
	sessionUserName, err := s.GetSessionUserName(ctx, req.SessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		fillSystemProtoResponse(ctx, &resp, common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return &resp, nil
	}

	// Task Service using RPC and get the taskID
	taskURI, err := s.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create task: " + err.Error()
		fillSystemProtoResponse(ctx, &resp, common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return &resp, nil
	}
	taskID := strings.TrimPrefix(taskURI, "/redfish/v1/TaskService/Tasks/")
	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Location": "/taskmon/" + taskID,
		},
	}
	generateTaskRespone(taskID, taskURI, &rpcResp)
	fillSystemProtoResponse(ctx, &resp, rpcResp)
	var pc = systems.PluginContact{
		ContactClient:  pmbhandle.ContactPlugin,
		DevicePassword: common.DecryptWithPrivateKey,
		UpdateTask:     s.UpdateTask,
	}

	var threadID int = 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.ResetSecureBoot)
	ctx = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go s.EI.ResetSecureBoot(ctx, req, &pc, taskID)
	threadID++
	l.LogWithFields(ctx).Debugf("outgoing response ResetSecureBoot: %s", string(resp.Body))
	return &resp, nil
}

func fillSystemProtoResponse(ctx context.Context, resp *systemsproto.SystemsResponse, data response.RPC) {
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
	resp.Header = data.Header
}
