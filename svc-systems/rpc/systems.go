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
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/scommon"
	"github.com/ODIM-Project/ODIM/svc-systems/systems"
)

// Systems struct helps to register service
type Systems struct {
	IsAuthorizedRPC func(sessionToken string, privileges, oemPrivileges []string) response.RPC
	EI              *systems.ExternalInterface
}

//GetSystemResource defines the operations which handles the RPC request response
// for the getting the system resource  of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (s *Systems) GetSystemResource(ctx context.Context, req *systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error) {
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp := s.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Error("error while trying to authenticate session")
		fillSystemProtoResponse(&resp, authResp)
		return &resp, nil
	}
	var pc = systems.PluginContact{
		ContactClient:   pmbhandle.ContactPlugin,
		DevicePassword:  common.DecryptWithPrivateKey,
		GetPluginStatus: scommon.GetPluginStatus,
	}
	data := pc.GetSystemResource(req)
	fillSystemProtoResponse(&resp, data)
	return &resp, nil
}

// GetSystemsCollection defines the operation which has the RPC request
// for getting the systems data from odimra.
// Retrieves all the keys with table name systems collection and create the response
// to send back to requested user.
func (s *Systems) GetSystemsCollection(ctx context.Context, req *systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error) {
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp := s.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Error("error while trying to authenticate session")
		fillSystemProtoResponse(&resp, authResp)
		return &resp, nil
	}
	data := systems.GetSystemsCollection(req)
	fillSystemProtoResponse(&resp, data)
	return &resp, nil
}

//GetSystems defines the operations which handles the RPC request response
// for the getting the system resource  of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the util-lib package.
// The function uses IsAuthorized of util-lib to validate the session
// which is present in the request.
func (s *Systems) GetSystems(ctx context.Context, req *systemsproto.GetSystemsRequest) (*systemsproto.SystemsResponse, error) {
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp := s.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Error("error while trying to authenticate session")
		fillSystemProtoResponse(&resp, authResp)
		return &resp, nil
	}
	var pc = systems.PluginContact{
		ContactClient:   pmbhandle.ContactPlugin,
		DevicePassword:  common.DecryptWithPrivateKey,
		GetPluginStatus: scommon.GetPluginStatus,
	}
	data := pc.GetSystems(req)
	fillSystemProtoResponse(&resp, data)
	return &resp, nil
}

// ComputerSystemReset defines the operations which handles the RPC request response
// for the ComputerSystemReset service of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function also checks for the session time out of the token
// which is present in the request.
func (s *Systems) ComputerSystemReset(ctx context.Context, req *systemsproto.ComputerSystemResetRequest) (*systemsproto.SystemsResponse, error) {
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp := s.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Error("error while trying to authenticate session")
		fillSystemProtoResponse(&resp, authResp)
		return &resp, nil
	}
	var pc = systems.PluginContact{
		ContactClient:  pmbhandle.ContactPlugin,
		DevicePassword: common.DecryptWithPrivateKey,
	}
	data := pc.ComputerSystemReset(req)
	fillSystemProtoResponse(&resp, data)
	return &resp, nil
}

// SetDefaultBootOrder defines the operations which handles the RPC request response
// for the SetDefaultBootOrder service of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function also checks for the session time out of the token
// which is present in the request.
func (s *Systems) SetDefaultBootOrder(ctx context.Context, req *systemsproto.DefaultBootOrderRequest) (*systemsproto.SystemsResponse, error) {
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp := s.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Error("error while trying to authenticate session")
		fillSystemProtoResponse(&resp, authResp)
		return &resp, nil
	}
	var pc = systems.PluginContact{
		ContactClient:  pmbhandle.ContactPlugin,
		DevicePassword: common.DecryptWithPrivateKey,
	}
	data := pc.SetDefaultBootOrder(req.SystemID)
	fillSystemProtoResponse(&resp, data)
	return &resp, nil
}

// ChangeBiosSettings defines the operations which handles the RPC request response
// for the ChangeBiosSettings service of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function also checks for the session time out of the token
// which is present in the request.
func (s *Systems) ChangeBiosSettings(ctx context.Context, req *systemsproto.BiosSettingsRequest) (*systemsproto.SystemsResponse, error) {
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp := s.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Error("error while trying to authenticate session")
		fillSystemProtoResponse(&resp, authResp)
		return &resp, nil
	}
	var pc = systems.PluginContact{
		ContactClient:  pmbhandle.ContactPlugin,
		DevicePassword: common.DecryptWithPrivateKey,
	}
	data := pc.ChangeBiosSettings(req)
	fillSystemProtoResponse(&resp, data)
	return &resp, nil
}

// ChangeBootOrderSettings defines the operations which handles the RPC request response
// for the ChangeBootOrderSettings service of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function also checks for the session time out of the token
// which is present in the request.
func (s *Systems) ChangeBootOrderSettings(ctx context.Context, req *systemsproto.BootOrderSettingsRequest) (*systemsproto.SystemsResponse, error) {
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp := s.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Error("error while trying to authenticate session")
		fillSystemProtoResponse(&resp, authResp)
		return &resp, nil
	}
	var pc = systems.PluginContact{
		ContactClient:  pmbhandle.ContactPlugin,
		DevicePassword: common.DecryptWithPrivateKey,
	}
	data := pc.ChangeBootOrderSettings(req)
	fillSystemProtoResponse(&resp, data)
	return &resp, nil
}

// CreateVolume defines the operations which handles the RPC request response
// for the CreateVolume service of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function also checks for the session time out of the token
// which is present in the request.
func (s *Systems) CreateVolume(ctx context.Context, req *systemsproto.VolumeRequest) (*systemsproto.SystemsResponse, error) {
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp := s.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Error("error while trying to authenticate session")
		fillSystemProtoResponse(&resp, authResp)
		return &resp, nil
	}

	data := s.EI.CreateVolume(req)
	fillSystemProtoResponse(&resp, data)
	return &resp, nil
}

// DeleteVolume defines the operations which handles the RPC request response
// for the DeleteVolume service of systems micro service.
// The functionality retrives the request and return backs the response to
// RPC according to the protoc file defined in the lib-utilities package.
// The function also checks for the session time out of the token
// which is present in the request.
func (s *Systems) DeleteVolume(ctx context.Context, req *systemsproto.VolumeRequest) (*systemsproto.SystemsResponse, error) {
	var resp systemsproto.SystemsResponse
	sessionToken := req.SessionToken
	authResp := s.IsAuthorizedRPC(sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Error("error while trying to authenticate session")
		fillSystemProtoResponse(&resp, authResp)
		return &resp, nil
	}

	data := s.EI.DeleteVolume(req)
	fillSystemProtoResponse(&resp, data)
	return &resp, nil
}

func fillSystemProtoResponse(resp *systemsproto.SystemsResponse, data response.RPC) {
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(data.Body)
	resp.Header = data.Header
}
