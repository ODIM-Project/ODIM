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
	"log"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
)

// GetUpdateService is an rpc handler, it gets invoked during GET on UpdateService API (/redfis/v1/UpdateService/)
func (a *Updater) GetUpdateService(ctx context.Context, req *updateproto.UpdateRequest, resp *updateproto.UpdateResponse) error {
	response := a.connector.GetUpdateService()
	body, err := json.Marshal(response.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying marshal the response body for get update service: " + err.Error()
		log.Printf(response.StatusMessage)
		return nil
	}
	resp.StatusCode = response.StatusCode
	resp.StatusMessage = response.StatusMessage
	resp.Header = response.Header
	resp.Body = body
	return nil
}

// GetFirmwareInventoryCollection an rpc handler which is invoked during GET on firmware inventory collection
func (a *Updater) GetFirmwareInventoryCollection(ctx context.Context, req *updateproto.UpdateRequest, resp *updateproto.UpdateResponse) error {
	sessionToken := req.SessionToken
	authStatusCode, authStatusMessage := a.connector.External.Auth(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authStatusCode != http.StatusOK {
		errorMessage := "error while trying to authenticate session"
		resp.StatusCode = authStatusCode
		resp.StatusMessage = authStatusMessage
		rpcResp := common.GeneralError(authStatusCode, authStatusMessage, errorMessage, nil, nil)
		bytes, err := json.Marshal(rpcResp.Body)
		if err != nil {
			log.Println("error in unmarshalling response object from util-libs", err.Error())
		}
		resp.Body = bytes
		resp.Header = rpcResp.Header
		log.Printf(errorMessage)
		return nil
	}
	response := a.connector.GetAllFirmwareInventory(req)
	resp.StatusCode = response.StatusCode
	resp.StatusMessage = response.StatusMessage
	resp.Header = response.Header
	bytes, err := json.Marshal(response.Body)
	if err != nil {
		log.Println("error in unmarshalling response object ", err.Error())
	}
	resp.Body = bytes
	return nil
}

// GetFirmwareInventory is an rpc handler which is invoked during GET on firmware inventory
func (a *Updater) GetFirmwareInventory(ctx context.Context, req *updateproto.UpdateRequest, resp *updateproto.UpdateResponse) error {
	sessionToken := req.SessionToken
	authStatusCode, authStatusMessage := a.connector.External.Auth(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authStatusCode != http.StatusOK {
		errorMessage := "error while trying to authenticate session"
		resp.StatusCode = authStatusCode
		resp.StatusMessage = authStatusMessage
		rpcResp := common.GeneralError(authStatusCode, authStatusMessage, errorMessage, nil, nil)
		bytes, err := json.Marshal(rpcResp.Body)
		if err != nil {
			log.Println("error in unmarshalling response object from util-libs", err.Error())
		}
		resp.Body = bytes
		resp.Header = rpcResp.Header
		log.Printf(errorMessage)
		return nil
	}
	response := a.connector.GetFirmwareInventory(req)
	body, err := json.Marshal(response.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying marshal the response body for get firmware inventory: " + err.Error()
		log.Printf(response.StatusMessage)
		return nil
	}
	resp.StatusCode = response.StatusCode
	resp.StatusMessage = response.StatusMessage
	resp.Header = response.Header
	resp.Body = body
	return nil
}

// GetSoftwareInventoryCollection is an rpc handler which is invoked during GET on software inventory collection
func (a *Updater) GetSoftwareInventoryCollection(ctx context.Context, req *updateproto.UpdateRequest, resp *updateproto.UpdateResponse) error {
	sessionToken := req.SessionToken
	authStatusCode, authStatusMessage := a.connector.External.Auth(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authStatusCode != http.StatusOK {
		errorMessage := "error while trying to authenticate session"
		resp.StatusCode = authStatusCode
		resp.StatusMessage = authStatusMessage
		rpcResp := common.GeneralError(authStatusCode, authStatusMessage, errorMessage, nil, nil)
		bytes, err := json.Marshal(rpcResp.Body)
		if err != nil {
			log.Println("error in unmarshalling response object from util-libs", err.Error())
		}
		resp.Body = bytes
		resp.Header = rpcResp.Header
		log.Printf(errorMessage)
		return nil
	}
	response := a.connector.GetAllSoftwareInventory(req)
	body, err := json.Marshal(response.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying marshal the response body for get software inventory collecion: " + err.Error()
		log.Printf(response.StatusMessage)
		return nil
	}
	resp.StatusCode = response.StatusCode
	resp.StatusMessage = response.StatusMessage
	resp.Header = response.Header
	resp.Body = body
	return nil
}

// GetSoftwareInventory is an rpc handler which is invoked during GET on software inventory
func (a *Updater) GetSoftwareInventory(ctx context.Context, req *updateproto.UpdateRequest, resp *updateproto.UpdateResponse) error {
	sessionToken := req.SessionToken
	authStatusCode, authStatusMessage := a.connector.External.Auth(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authStatusCode != http.StatusOK {
		errorMessage := "error while trying to authenticate session"
		resp.StatusCode = authStatusCode
		resp.StatusMessage = authStatusMessage
		rpcResp := common.GeneralError(authStatusCode, authStatusMessage, errorMessage, nil, nil)
		bytes, err := json.Marshal(rpcResp.Body)
		if err != nil {
			log.Println("error in unmarshalling response object from util-libs", err.Error())
		}
		resp.Body = bytes
		resp.Header = rpcResp.Header
		log.Printf(errorMessage)
		return nil
	}
	response := a.connector.GetSoftwareInventory(req)
	body, err := json.Marshal(response.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying marshal the response body for get software inventory: " + err.Error()
		log.Printf(response.StatusMessage)
		return nil
	}
	resp.StatusCode = response.StatusCode
	resp.StatusMessage = response.StatusMessage
	resp.Header = response.Header
	resp.Body = body
	return nil
}

// SimepleUpdate is an rpc handler, it gets involked during POST on UpdateService API actions (/Actions/UpdateService.SimpleUpdate)
func (a *Updater) SimepleUpdate(ctx context.Context, req *updateproto.UpdateRequest, resp *updateproto.UpdateResponse) error {
	sessionToken := req.SessionToken
	authStatusCode, authStatusMessage := a.connector.External.Auth(sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authStatusCode != http.StatusOK {
		errorMessage := "error while trying to authenticate session"
		resp.StatusCode = authStatusCode
		resp.StatusMessage = authStatusMessage
		rpcResp := common.GeneralError(authStatusCode, authStatusMessage, errorMessage, nil, nil)
		bytes, err := json.Marshal(rpcResp.Body)
		if err != nil {
			log.Println("error in unmarshalling response object from util-libs", err.Error())
		}
		resp.Body = bytes
		resp.Header = rpcResp.Header
		log.Printf(errorMessage)
		return nil
	}
	response := a.connector.SimpleUpdate(req)
	body, err := json.Marshal(response.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying marshal the response body for get update service: " + err.Error()
		log.Printf(response.StatusMessage)
		return nil
	}
	resp.StatusCode = response.StatusCode
	resp.StatusMessage = response.StatusMessage
	resp.Header = response.Header
	resp.Body = body
	return nil
}

// StartUpdate is an rpc handler, it gets involked during POST on UpdateService API actions (/Actions/UpdateService.StartUpdate)
func (a *Updater) StartUpdate(ctx context.Context, req *updateproto.UpdateRequest, resp *updateproto.UpdateResponse) error {
	sessionToken := req.SessionToken
	authStatusCode, authStatusMessage := a.connector.External.Auth(sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authStatusCode != http.StatusOK {
		errorMessage := "error while trying to authenticate session"
		resp.StatusCode = authStatusCode
		resp.StatusMessage = authStatusMessage
		rpcResp := common.GeneralError(authStatusCode, authStatusMessage, errorMessage, nil, nil)
		bytes, err := json.Marshal(rpcResp.Body)
		if err != nil {
			log.Println("error in unmarshalling response object from util-libs", err.Error())
		}
		resp.Body = bytes
		resp.Header = rpcResp.Header
		log.Printf(errorMessage)
		return nil
	}
	response := a.connector.StartUpdate(req)
	body, err := json.Marshal(response.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = "error while trying marshal the response body for get update service: " + err.Error()
		log.Printf(response.StatusMessage)
		return nil
	}
	resp.StatusCode = response.StatusCode
	resp.StatusMessage = response.StatusMessage
	resp.Header = response.Header
	resp.Body = body
	return nil
}
