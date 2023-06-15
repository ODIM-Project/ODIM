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

// Package handle ...
package handle

import (
	"context"
	"encoding/json"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	iris "github.com/kataras/iris/v12"
)

// RoleRPCs defines all the RPC methods in role
type RoleRPCs struct {
	GetAllRolesRPC func(context.Context, roleproto.GetRoleRequest) (*roleproto.RoleResponse, error)
	GetRoleRPC     func(context.Context, roleproto.GetRoleRequest) (*roleproto.RoleResponse, error)
	UpdateRoleRPC  func(context.Context, roleproto.UpdateRoleRequest) (*roleproto.RoleResponse, error)
	DeleteRoleRPC  func(context.Context, roleproto.DeleteRoleRequest) (*roleproto.RoleResponse, error)
}

// GetAllRoles defines the GetAllRoles iris handler.
// The method extract the session token and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (r *RoleRPCs) GetAllRoles(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	l.LogWithFields(ctxt).Debug("Incoming request received for the getting all roles")
	req := roleproto.GetRoleRequest{SessionToken: ctx.Request().Header.Get("X-Auth-Token")}
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := r.GetAllRolesRPC(ctxt, req)
	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctx).Debugf("Outgoing response for Getting all roles is %s and response status %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	sendRoleResponse(ctx, resp)
}

// GetRole defines the GetRole iris handler.
// The method extract the session token, and necessary
// request parameters and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (r *RoleRPCs) GetRole(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := roleproto.GetRoleRequest{SessionToken: ctx.Request().Header.Get("X-Auth-Token")}
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	req.Id = ctx.Params().Get("id")
	l.LogWithFields(ctxt).Debugf("Incoming request received for the getting a role with id %s", req.Id)
	resp, err := r.GetRoleRPC(ctxt, req)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	if req.Id == common.RoleAdmin || req.Id == common.RoleClient || req.Id == common.RoleMonitor {
		ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH")
	} else {
		ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH, DELETE")
	}
	l.LogWithFields(ctx).Debugf("Outgoing response for Getting a role is %s and response status %d", string(resp.Body), int(resp.StatusCode))
	sendRoleResponse(ctx, resp)
}

// UpdateRole defines the UpdateRole iris handler.
// The method extract the session token, and necessary
// request parameters and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (r *RoleRPCs) UpdateRole(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req roleproto.UpdateRoleRequest

	//Read Body from Request
	var roleReq interface{}
	err := ctx.ReadJSON(&roleReq)
	if err != nil {
		l.LogWithFields(ctxt).Error("Error while trying to collect data from request: " + err.Error())
		errorMessage := "error while trying to get JSON body from the role update request body: " + err.Error()
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}

	req.SessionToken = ctx.Request().Header.Get("X-Auth-Token")
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	req.Id = ctx.Params().Get("id")
	l.LogWithFields(ctxt).Debugf("Incoming request received for the updating a role with id %s", req.Id)
	req.UpdateRequest, _ = json.Marshal(&roleReq)
	resp, err := r.UpdateRoleRPC(ctxt, req)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctx).Debugf("Outgoing response for updating a role is %s and response status %d", string(resp.Body), int(resp.StatusCode))
	sendRoleResponse(ctx, resp)
}

// DeleteRole ...
func (r *RoleRPCs) DeleteRole(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req roleproto.DeleteRoleRequest

	req.SessionToken = ctx.Request().Header.Get("X-Auth-Token")
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	req.ID = ctx.Params().Get("id")
	l.LogWithFields(ctxt).Debugf("Incoming request received for the deleting a role with id %s", req.ID)
	resp, err := r.DeleteRoleRPC(ctxt, req)
	if err != nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctx).Debugf("Outgoing response for deleting a role is %s and response status %d", string(resp.Body), int(resp.StatusCode))
	sendRoleResponse(ctx, resp)
}

// sendRoleResponse writes the role response to client
func sendRoleResponse(ctx iris.Context, resp *roleproto.RoleResponse) {
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}
