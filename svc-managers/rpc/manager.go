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
	"net/http"
	"os"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-managers/managers"
)

// Managers struct helps to register service
type Managers struct {
	IsAuthorizedRPC func(sessionToken string, privileges, oemPrivileges []string) (response.RPC, error)
	EI              *managers.ExternalInterface
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
	l.LogWithFields(ctx).Info("Inside GetManagersCollection function (svc-managers)")
	var resp managersproto.ManagerResponse
	sessionToken := req.SessionToken
	authResp, err := m.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeLogin}, []string{})
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
	data, _ := m.EI.GetManagersCollection(ctx, req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
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
	l.LogWithFields(ctx).Info("Inside GetManager function (svc-managers)")
	var resp managersproto.ManagerResponse
	sessionToken := req.SessionToken
	authResp, err := m.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeLogin}, []string{})
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
	l.LogWithFields(ctx).Info("Inside GetManagersResource function (svc-managers)")
	var resp managersproto.ManagerResponse
	sessionToken := req.SessionToken
	authResp, err := m.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeLogin}, []string{})
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
	return &resp, nil
}

// VirtualMediaInsert defines the operations which handles the RPC request response
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (m *Managers) VirtualMediaInsert(ctx context.Context, req *managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.ManagerService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside VirtualMediaInsert function (svc-managers)")
	sessionToken := req.SessionToken
	resp := &managersproto.ManagerResponse{}
	authResp, err := m.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		resp.StatusCode = authResp.StatusCode
		resp.StatusMessage = authResp.StatusMessage
		resp.Body = generateResponse(ctx, authResp.Body)
		resp.Header = authResp.Header
		return resp, nil
	}
	data := m.EI.VirtualMediaActions(ctx, req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
	return resp, nil
}

// VirtualMediaEject defines the operations which handles the RPC request response
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (m *Managers) VirtualMediaEject(ctx context.Context, req *managersproto.ManagerRequest) (*managersproto.ManagerResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = context.WithValue(ctx, common.ThreadName, common.ManagerService)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	l.LogWithFields(ctx).Info("Inside VirtualMediaEject function (svc-managers)")
	sessionToken := req.SessionToken
	resp := &managersproto.ManagerResponse{}
	authResp, err := m.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		resp.StatusCode = authResp.StatusCode
		resp.StatusMessage = authResp.StatusMessage
		resp.Body = generateResponse(ctx, authResp.Body)
		resp.Header = authResp.Header
		return resp, nil
	}
	data := m.EI.VirtualMediaActions(ctx, req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
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
	l.LogWithFields(ctx).Info("Inside GetRemoteAccountService function (svc-managers)")
	var resp managersproto.ManagerResponse
	sessionToken := req.SessionToken
	authResp, err := m.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeLogin}, []string{})
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
	data := m.EI.GetRemoteAccountService(ctx, req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
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
	l.LogWithFields(ctx).Info("Inside CreateRemoteAccountService function (svc-managers)")
	var resp managersproto.ManagerResponse
	sessionToken := req.SessionToken
	authResp, err := m.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeConfigureUsers}, []string{})
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
	data := m.EI.CreateRemoteAccountService(ctx, req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
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
	l.LogWithFields(ctx).Info("Inside UpdateRemoteAccountService function (svc-managers)")
	var resp managersproto.ManagerResponse
	sessionToken := req.SessionToken
	authResp, err := m.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeConfigureUsers}, []string{})
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
	data := m.EI.UpdateRemoteAccountService(ctx, req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
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
	l.LogWithFields(ctx).Info("Inside DeleteRemoteAccountService function (svc-managers)")
	var resp managersproto.ManagerResponse
	sessionToken := req.SessionToken
	authResp, err := m.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeConfigureUsers}, []string{})
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
	data := m.EI.DeleteRemoteAccountService(ctx, req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
	return &resp, nil
}
