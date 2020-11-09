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

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
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
	fillProtoResponse(resp, a.connector.SimpleUpdate(req))
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
