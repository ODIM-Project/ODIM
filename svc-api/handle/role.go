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

//Package handle ...
package handle

import (
	"encoding/json"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
)

// RoleRPCs defines all the RPC methods in role
type RoleRPCs struct {
	GetAllRolesRPC func(roleproto.GetRoleRequest) (*roleproto.RoleResponse, error)
	CreateRoleRPC  func(roleproto.RoleRequest) (*roleproto.RoleResponse, error)
	GetRoleRPC     func(roleproto.GetRoleRequest) (*roleproto.RoleResponse, error)
	UpdateRoleRPC  func(roleproto.UpdateRoleRequest) (*roleproto.RoleResponse, error)
	DeleteRoleRPC  func(roleproto.DeleteRoleRequest) (*roleproto.RoleResponse, error)
}

// GetAllRoles defines the GetAllRoles iris handler.
// The method extract the session token and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (r *RoleRPCs) GetAllRoles(ctx iris.Context) {
	defer ctx.Next()
	req := roleproto.GetRoleRequest{SessionToken: ctx.Request().Header.Get("X-Auth-Token")}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	resp, err := r.GetAllRolesRPC(req)
	if err != nil {
		errorMessage := "RPC error: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// CreateRole defines the CreateRole iris handler.
// The method extract the session token, and necessary
// request parameters and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (r *RoleRPCs) CreateRole(ctx iris.Context) {
	defer ctx.Next()
	var req interface{}
	//Read Body from Request
	err := ctx.ReadJSON(&req)
	if err != nil {
		log.Error("Error while trying to collect data from request: " + err.Error())
		errorMessage := "error while trying to get JSON body from the account create request body: " + err.Error()
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}

	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	// Marshalling the req to make role request
	request, err := json.Marshal(req)
	roleRequest := roleproto.RoleRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
	}

	resp, err := r.CreateRoleRPC(roleRequest)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetRole defines the GetRole iris handler.
// The method extract the session token, and necessary
// request parameters and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (r *RoleRPCs) GetRole(ctx iris.Context) {
	defer ctx.Next()
	req := roleproto.GetRoleRequest{SessionToken: ctx.Request().Header.Get("X-Auth-Token")}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	req.Id = ctx.Params().Get("id")

	resp, err := r.GetRoleRPC(req)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	if req.Id == common.RoleAdmin{
		ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH")
	} else {
		ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH, DELETE")
	}
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// UpdateRole defines the UpdateRole iris handler.
// The method extract the session token, and necessary
// request parameters and creates the RPC request.
// After the RPC call the method will feed the response to the iris
// and gives out a proper response.
func (r *RoleRPCs) UpdateRole(ctx iris.Context) {
	defer ctx.Next()

	var req roleproto.UpdateRoleRequest

	//Read Body from Request
	var roleReq interface{}
	err := ctx.ReadJSON(&roleReq)
	if err != nil {
		log.Error("Error while trying to collect data from request: " + err.Error())
		errorMessage := "error while trying to get JSON body from the role update request body: " + err.Error()
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(response.Body)
		return
	}

	req.SessionToken = ctx.Request().Header.Get("X-Auth-Token")
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	req.Id = ctx.Params().Get("id")
	req.UpdateRequest, _ = json.Marshal(&roleReq)
	resp, err := r.UpdateRoleRPC(req)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// DeleteRole ...
func (r *RoleRPCs) DeleteRole(ctx iris.Context) {
	defer ctx.Next()

	var req roleproto.DeleteRoleRequest

	req.SessionToken = ctx.Request().Header.Get("X-Auth-Token")
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	req.ID = ctx.Params().Get("id")

	resp, err := r.DeleteRoleRPC(req)
	if err != nil {
		errorMessage := "error: something went wrong with the RPC calls: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}
