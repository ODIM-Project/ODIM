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

//Package rpc ...
package rpc

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	"github.com/ODIM-Project/ODIM/svc-managers/managers"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrcommon"
)

// Managers struct helps to register service
type Managers struct {
	IsAuthorizedRPC  func(sessionToken string, privileges, oemPrivileges []string) (int32, string)
	ContactClientRPC func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	EI *managers.ExternalInterface
}

//GetManagersCollection defines the operation which hasnled the RPC request response
// for getting the odimra systems.
// Retrieves all the keys with table name systems collection and create the response
// to send back to requested user.
func (m *Managers) GetManagersCollection(ctx context.Context, req *managersproto.ManagerRequest, resp *managersproto.ManagerResponse) error {
	sessionToken := req.SessionToken
	authStatusCode, authStatusMessage := m.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authStatusCode != http.StatusOK {
		errorMessage := "error while trying to authenticate session"
		resp.StatusCode = authStatusCode
		resp.StatusMessage = authStatusMessage
		rpcResp := common.GeneralError(authStatusCode, authStatusMessage, errorMessage, nil, nil)
		resp.Body = generateResponse(rpcResp.Body)
		resp.Header = rpcResp.Header
		log.Printf(errorMessage)
		return nil
	}
	data, _ := m.EI.GetManagersCollection(req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(data.Body)
	return nil
}

//GetManager defines the operations which handles the RPC request response
// for the getting the system resource  of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (m *Managers) GetManager(ctx context.Context, req *managersproto.ManagerRequest, resp *managersproto.ManagerResponse) error {
	sessionToken := req.SessionToken
	authStatusCode, authStatusMessage := m.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authStatusCode != http.StatusOK {
		errorMessage := "error while trying to authenticate session"
		resp.StatusCode = authStatusCode
		resp.StatusMessage = authStatusMessage
		rpcResp := common.GeneralError(authStatusCode, authStatusMessage, errorMessage, nil, nil)
		resp.Body = generateResponse(rpcResp.Body)
		resp.Header = rpcResp.Header
		log.Printf(errorMessage)
		return nil
	}
	var d = managers.DeviceContact{
		GetDeviceInfo:         mgrcommon.GetResourceInfoFromDevice,
		ContactClient:         m.ContactClientRPC,
		DecryptDevicePassword: common.DecryptWithPrivateKey,
	}
	data := d.GetManagers(req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(data.Body)
	return nil
}

//GetManagersResource defines the operations which handles the RPC request response
// for the getting the system resource  of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (m *Managers) GetManagersResource(ctx context.Context, req *managersproto.ManagerRequest, resp *managersproto.ManagerResponse) error {
	sessionToken := req.SessionToken
	authStatusCode, authStatusMessage := m.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authStatusCode != http.StatusOK {
		errorMessage := "error while trying to authenticate session"
		resp.StatusCode = authStatusCode
		resp.StatusMessage = authStatusMessage
		rpcResp := common.GeneralError(authStatusCode, authStatusMessage, errorMessage, nil, nil)
		resp.Body = generateResponse(rpcResp.Body)
		resp.Header = rpcResp.Header
		log.Printf(errorMessage)
		return nil
	}
	var d = managers.DeviceContact{
		GetDeviceInfo:         mgrcommon.GetResourceInfoFromDevice,
		ContactClient:         m.ContactClientRPC,
		DecryptDevicePassword: common.DecryptWithPrivateKey,
	}
	data := d.GetManagersResource(req)
	resp.Header = data.Header
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(data.Body)
	return nil
}

func generateResponse(input interface{}) []byte {
	bytes, err := json.Marshal(input)
	if err != nil {
		log.Println("error in unmarshalling response object from util-libs", err.Error())
	}
	return bytes
}
