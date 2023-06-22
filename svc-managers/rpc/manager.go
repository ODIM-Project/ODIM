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
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-managers/managers"
)

// Managers struct helps to register service
type Managers struct {
	GetSessionUserName func(context.Context, string) (string, error)
	CreateTask         func(ctx context.Context, sessionUserName string) (string, error)
	SavePluginTaskInfo func(ctx context.Context, pluginIP, pluginServerName, odimTaskID, pluginTaskMonURL string) error
	IsAuthorizedRPC    func(ctx context.Context, sessionToken string, privileges, oemPrivileges []string) (response.RPC, error)
	EI                 *managers.ExternalInterface
}

// podName defines the current name of process
var podName = os.Getenv("POD_NAME")

// GetManagersCollection defines the operation which hasnled the RPC request response
// for getting the odimra systems.
// Retrieves all the keys with table name systems collection and create the response
// to send back to requested user.
func (m *Managers) GetManagersCollection(ctx context.Context, req *managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.ManagerService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	var resp managersproto.ManagerResponse
	sessionToken := req.SessionToken
	authResp, err := m.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("error while authorizing the session token : %s", err.Error())
		}
		resp.StatusCode = authResp.StatusCode
		resp.StatusMessage = authResp.StatusMessage
		resp.Body = generateResponse(ctx, authResp.Body)
		resp.Header = authResp.Header
		return &resp, nil
	}
	data, _ := m.EI.GetManagersCollection(ctx, req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
	l.LogWithFields(ctx).Debugf("Outgoing manager collection response to northbound: %s", string(resp.Body))
	return &resp, nil
}

// GetManager defines the operations which handles the RPC request response
// for the getting the system resource  of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (m *Managers) GetManager(ctx context.Context, req *managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.ManagerService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	var resp managersproto.ManagerResponse
	sessionToken := req.SessionToken
	authResp, err := m.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		resp.StatusCode = authResp.StatusCode
		resp.StatusMessage = authResp.StatusMessage
		resp.Body = generateResponse(ctx, authResp.Body)
		resp.Header = authResp.Header
		return &resp, nil
	}
	data := m.EI.GetManagers(ctx, req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
	return &resp, nil
}

// GetManagersResource defines the operations which handles the RPC request response
// for the getting the system resource  of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (m *Managers) GetManagersResource(ctx context.Context, req *managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.ManagerService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	var resp managersproto.ManagerResponse
	sessionToken := req.SessionToken
	authResp, err := m.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		resp.StatusCode = authResp.StatusCode
		resp.StatusMessage = authResp.StatusMessage
		resp.Body = generateResponse(ctx, authResp.Body)
		resp.Header = authResp.Header
		return &resp, nil
	}
	data := m.EI.GetManagersResource(ctx, req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
	l.LogWithFields(ctx).Debugf("Outgoing manager resource response to northbound: %s", string(resp.Body))
	return &resp, nil
}

// VirtualMediaInsert defines the operations which handles the RPC request response
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (m *Managers) VirtualMediaInsert(ctx context.Context, req *managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.ManagerService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	sessionToken := req.SessionToken
	resp := &managersproto.ManagerResponse{}
	authResp, err := m.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("error while authorizing the session token : %s", err.Error())
		}
		resp.StatusCode = authResp.StatusCode
		resp.StatusMessage = authResp.StatusMessage
		resp.Body = generateResponse(ctx, authResp.Body)
		resp.Header = authResp.Header
		return resp, nil
	}
	sessionUserName, err := m.GetSessionUserName(ctx, req.SessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		fillManagersProtoResponse(ctx, resp, common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
	}
	// Task Service using RPC and get the taskID
	taskURI, err := m.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create task: " + err.Error()
		fillManagersProtoResponse(ctx, resp, common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
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
	generateTaskResponse(taskID, taskURI, &rpcResp)
	fillManagersProtoResponse(ctx, resp, rpcResp)
	go m.EI.VirtualMediaActions(ctx, req, taskID)
	l.LogWithFields(ctx).Debugf("Outgoing virtual media response to northbound: %s", string(resp.Body))
	return resp, nil
}

// VirtualMediaEject defines the operations which handles the RPC request response
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (m *Managers) VirtualMediaEject(ctx context.Context, req *managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.ManagerService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	sessionToken := req.SessionToken
	resp := &managersproto.ManagerResponse{}
	authResp, err := m.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("error while authorizing the session token : %s", err.Error())
		}
		resp.StatusCode = authResp.StatusCode
		resp.StatusMessage = authResp.StatusMessage
		resp.Body = generateResponse(ctx, authResp.Body)
		resp.Header = authResp.Header
		return resp, nil
	}
	sessionUserName, err := m.GetSessionUserName(ctx, req.SessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		fillManagersProtoResponse(ctx, resp, common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
	}
	// Task Service using RPC and get the taskID
	taskURI, err := m.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create task: " + err.Error()
		fillManagersProtoResponse(ctx, resp, common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil))
		l.LogWithFields(ctx).Error(errMsg)
		return resp, nil
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
	generateTaskResponse(taskID, taskURI, &rpcResp)
	fillManagersProtoResponse(ctx, resp, rpcResp)
	go m.EI.VirtualMediaActions(ctx, req, taskID)
	l.LogWithFields(ctx).Debugf("Outgoing virtual media eject response to northbound: %s", string(resp.Body))
	return resp, nil
}

func generateResponse(ctx context.Context, input interface{}) []byte {
	bytes, err := json.Marshal(input)
	if err != nil {
		l.LogWithFields(ctx).Error("error in unmarshalling response object from util-libs" + err.Error())
	}
	return bytes
}

// GetRemoteAccountService defines the operations which handles the RPC request response
// The functionality retrieves the request and return backs the response to
// RPC according to the protoc file defined in the lib-util package.
// The function uses IsAuthorized of lib-util to validate the session token
// which is present in the request.
func (m *Managers) GetRemoteAccountService(ctx context.Context, req *managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.ManagerService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	var resp managersproto.ManagerResponse
	sessionToken := req.SessionToken
	authResp, err := m.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("error while authorizing the session token : %s", err.Error())
		}
		resp.StatusCode = authResp.StatusCode
		resp.StatusMessage = authResp.StatusMessage
		resp.Body = generateResponse(ctx, authResp.Body)
		resp.Header = authResp.Header
		return &resp, nil
	}
	data := m.EI.GetRemoteAccountService(ctx, req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
	l.LogWithFields(ctx).Debugf("Outgoing remote account service response to northbound: %s", string(resp.Body))
	return &resp, nil
}

// CreateRemoteAccountService defines the operations which handles the RPC request response
// The functionality retrieves the request and return backs the response to
// RPC according to the protoc file defined in the lib-util package.
// The function uses IsAuthorized of lib-util to validate the session token
// which is present in the request.
func (m *Managers) CreateRemoteAccountService(ctx context.Context, req *managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.ManagerService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	var resp managersproto.ManagerResponse
	sessionToken := req.SessionToken
	authResp, err := m.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeConfigureUsers}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("error while authorizing the session token : %s", err.Error())
		}
		fillManagersProtoResponse(ctx, &resp, authResp)
		return &resp, nil
	}

	taskID, err := CreateTaskAndResponse(ctx, m, req.SessionToken, &resp)
	if err != nil {
		l.LogWithFields(ctx).Error(err)
		return &resp, nil
	}

	var threadID int = 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.CreateRemoteAccountService)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go m.EI.CreateRemoteAccountService(ctxt, req, taskID)
	l.LogWithFields(ctx).Debugf("Outgoing create remote account service response to northbound: %s ",
		&resp.Body)
	return &resp, nil
}

// UpdateRemoteAccountService defines the operations which handles the RPC request response
// The functionality retrieves the request and return backs the response to
// RPC according to the protoc file defined in the lib-util package.
// The function uses IsAuthorized of lib-util to validate the session token
// which is present in the request.
func (m *Managers) UpdateRemoteAccountService(ctx context.Context, req *managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.ManagerService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	var resp managersproto.ManagerResponse
	sessionToken := req.SessionToken
	authResp, err := m.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeConfigureUsers}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("error while authorizing the session token : %s", err.Error())
		}
		fillManagersProtoResponse(ctx, &resp, authResp)
		return &resp, nil
	}

	taskID, err := CreateTaskAndResponse(ctx, m, req.SessionToken, &resp)
	if err != nil {
		l.LogWithFields(ctx).Error(err)
		return &resp, nil
	}

	var threadID int = 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.UpdateRemoteAccountService)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go m.EI.UpdateRemoteAccountService(ctx, req, taskID)
	l.LogWithFields(ctx).Debugf("Outgoing update remote account service response to northbound: %s", string(resp.Body))
	return &resp, nil
}

// DeleteRemoteAccountService defines the operations which handles the RPC request response
// The functionality retrieves the request and return backs the response to
// RPC according to the protoc file defined in the lib-util package.
// The function uses IsAuthorized of lib-util to validate the session token
// which is present in the request.
func (m *Managers) DeleteRemoteAccountService(ctx context.Context, req *managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.ManagerService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	var resp managersproto.ManagerResponse
	sessionToken := req.SessionToken
	authResp, err := m.IsAuthorizedRPC(ctx, sessionToken, []string{common.PrivilegeConfigureUsers}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("error while authorizing the session token : %s", err.Error())
		}
		fillManagersProtoResponse(ctx, &resp, authResp)
		return &resp, nil
	}
	taskID, err := CreateTaskAndResponse(ctx, m, req.SessionToken, &resp)
	if err != nil {
		l.LogWithFields(ctx).Error(err)
		return &resp, nil
	}

	var threadID int = 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.DeleteRemoteAccountService)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	go m.EI.DeleteRemoteAccountService(ctx, req, taskID)
	l.LogWithFields(ctx).Debugf("Outgoing delete remote account service response to northbound: %s", string(resp.Body))
	return &resp, nil
}

// CreateTaskAndResponse will create the task for corresponding request using
// the RPC call to task service and it will prepare custom task response to the user
// The function returns the ID of created task back.
func CreateTaskAndResponse(ctx context.Context, m *Managers, sessionToken string, resp *managersproto.ManagerResponse) (string, error) {
	sessionUserName, err := m.GetSessionUserName(ctx, sessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		fillManagersProtoResponse(ctx, resp, common.GeneralError(http.StatusUnauthorized,
			response.NoValidSession, errMsg, nil, nil))
		return "", fmt.Errorf(errMsg)
	}

	// Task Service using RPC and get the taskID
	taskURI, err := m.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create task: " + err.Error()
		fillManagersProtoResponse(ctx, resp, common.GeneralError(http.StatusInternalServerError,
			response.InternalError, errMsg, nil, nil))
		return "", fmt.Errorf(errMsg)
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

	generateTaskResponse(taskID, taskURI, &rpcResp)
	fillManagersProtoResponse(ctx, resp, rpcResp)
	return taskID, nil
}

func fillManagersProtoResponse(ctx context.Context, resp *managersproto.ManagerResponse, data response.RPC) {
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
	resp.Header = data.Header
}

// generateTaskResponse is used to generate task response
func generateTaskResponse(taskID, taskURI string, rpcResp *response.RPC) {
	commonResponse := response.Response{
		OdataType:    common.TaskType,
		ID:           taskID,
		Name:         "Task " + taskID,
		OdataContext: "/redfish/v1/$metadata#Task.Task",
		OdataID:      taskURI,
	}
	commonResponse.MessageArgs = []string{taskID}
	commonResponse.CreateGenericResponse(rpcResp.StatusMessage)
	rpcResp.Body = commonResponse
}

// UpdateRemoteAccountPassword defines the operations which handles the RPC request response
// The functionality retrieves the request and return backs the response to
// RPC according to the protoc file defined in the lib-util package.
// The function uses IsAuthorized of lib-util to validate the session token
// which is present in the request.
func (m *Managers) UpdateRemoteAccountPassword(ctx context.Context, req *managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.ManagerService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	var resp managersproto.ManagerResponse
	var threadID int = 1
	ctxt := context.WithValue(ctx, common.ThreadName, common.UpdateRemoteAccountService)
	ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
	ctxt = context.WithValue(ctxt, common.ActionName, "UpdateRemoteAccountPassword")
	go m.EI.UpdateRemoteAccountPasswordService(ctx, req)
	l.LogWithFields(ctx).Debugf("Outgoing update remote account service response to northbound: %s", string(resp.Body))
	return &resp, nil
}
