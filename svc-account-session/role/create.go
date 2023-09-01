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

// Package role ...
package role

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/account"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
	validator "gopkg.in/go-playground/validator.v9"
)

// Create defines creation of a new role. The function is supposed to be used as part of RPC.
//
// For creating an role, two parameters need to be passed RoleRequest and Session.
// New RoleID,AssignedPrivileges and OemPrivileges will be part of RoleRequest,
// and Session parameter will have all session related data, espically the privileges.
// For creating new role the ConfigureUsers privilege is mandatory.
//
// There will be two return values for the fuction. One is the RPC response, which contains the
// status code, status message, headers and body and the second value is error.
func Create(ctx context.Context, req *roleproto.RoleRequest, session *asmodel.Session) response.RPC {
	// parsing the request body
	var createRoleReq asmodel.Role
	err := json.Unmarshal(req.RequestBody, &createRoleReq)
	if err != nil {
		errMsg := "error while trying to parse the request body of create role API" + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}

	errorLogPrefix := fmt.Sprintf("failed to create role %s: ", createRoleReq.ID)
	commonResponse := response.Response{
		OdataType: common.RoleType,
		OdataID:   "/redfish/v1/AccountService/Roles/" + createRoleReq.ID,
		Name:      "User Role",
	}
	var resp response.RPC

	l.LogWithFields(ctx).Infof("Validating the request to create the role %s", createRoleReq.ID)
	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, createRoleReq)
	if err != nil {
		errMsg := errorLogPrefix + "Unable to validate request parameters: " + err.Error()
		l.LogWithFields(ctx).Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	} else if invalidProperties != "" {
		errorMessage := errorLogPrefix + "One or more properties given in the request body are not valid, ensure properties are listed in upper camel case "
		l.LogWithFields(ctx).Error(errorMessage)
		resp := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, nil)
		return resp
	}

	// register the validator package
	validate := validator.New()
	// create validation to check if the ID field contains whitespaces when creating role.
	validate.RegisterValidation("is-empty", CheckWhitespace)
	if err := validate.Var(createRoleReq.ID, "required"); err != nil {
		errorMessage := errorLogPrefix + "error: Mandatory field RoleId is empty"
		response := common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errorMessage, []interface{}{"RoleId"}, nil)
		return response
	} else if err := validate.Var(createRoleReq.ID, "required,is-empty,excludesall=!@#?%&$*"); err != nil {
		errorMessage := errorLogPrefix + "Invalid create role request"
		args := account.GetResponseArgs(response.PropertyValueNotInList, errorMessage, []interface{}{createRoleReq.ID, "RoleId"})
		resp.StatusCode = http.StatusBadRequest
		resp.StatusMessage = response.PropertyValueNotInList
		resp.Body = args.CreateGenericErrorResponse()
		l.LogWithFields(ctx).Error(errorMessage)
		return resp
	}

	//check for ConfigureUsers privilege in session object
	status, err := checkForPrivilege(session, "ConfigureUsers")
	if err != nil {
		errorMessage := errorLogPrefix + "User does not have the privilege of creating a new role"
		resp.StatusCode = int32(status.Code)
		resp.StatusMessage = status.Message
		args := account.GetResponseArgs(status.Message, errorMessage, []interface{}{})
		resp.Body = args.CreateGenericErrorResponse()
		auth.CustomAuthLog(ctx, session.Token, errorMessage, resp.StatusCode)
		return resp
	}
	if len(createRoleReq.AssignedPrivileges) == 0 && len(createRoleReq.OEMPrivileges) == 0 {
		errorMessage := errorLogPrefix + "Both AssignedPrivileges and OemPrivileges cannot be empty."
		args := account.GetResponseArgs(response.PropertyMissing, errorMessage, []interface{}{"AssignedPrivileges/OemPrivileges"})
		resp.StatusCode = http.StatusBadRequest
		resp.StatusMessage = response.PropertyMissing
		resp.Body = args.CreateGenericErrorResponse()
		l.LogWithFields(ctx).Error(errorMessage)
		return resp
	}

	if len(createRoleReq.AssignedPrivileges) != 0 {
		status, messageArgs, err := validateAssignedPrivileges(ctx, createRoleReq.AssignedPrivileges)
		if err != nil {
			errorMessage := errorLogPrefix + err.Error()
			resp.StatusCode = int32(status.Code)
			resp.StatusMessage = status.Message
			args := account.GetResponseArgs(status.Message, errorMessage, messageArgs)
			resp.Body = args.CreateGenericErrorResponse()
			l.LogWithFields(ctx).Error(errorMessage)
			return resp
		}
	}
	if len(createRoleReq.OEMPrivileges) != 0 {
		status, messageArgs, err := validateOEMPrivileges(ctx, createRoleReq.OEMPrivileges)
		if err != nil {
			errorMessage := errorLogPrefix + err.Error()
			resp.StatusCode = int32(status.Code)
			resp.StatusMessage = status.Message
			args := account.GetResponseArgs(status.Message, errorMessage, messageArgs)
			resp.Body = args.CreateGenericErrorResponse()
			l.LogWithFields(ctx).Error(errorMessage)
			return resp
		}
	}
	//Get redfish roles from database
	redfishRoles, gerr := asmodel.GetRedfishRoles()
	if gerr != nil {
		l.LogWithFields(ctx).Error(errorLogPrefix + "Unable to get redfish roles: " + gerr.Error())
		errorMessage := gerr.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		return resp

	}
	//check if request role is predefined redfish role and set isPredfined to true
	isPredefined := false
	for _, redfishrole := range redfishRoles.List {
		if createRoleReq.ID == redfishrole {
			isPredefined = true
		}
	}
	if isPredefined {
		errorMessage := errorLogPrefix + "Cannot create pre-defined roles"
		args := account.GetResponseArgs(response.InsufficientPrivilege, errorMessage, []interface{}{})
		resp.StatusCode = http.StatusForbidden
		resp.StatusMessage = response.InsufficientPrivilege
		resp.Body = args.CreateGenericErrorResponse()
		l.LogWithFields(ctx).Error(errorMessage)
		return resp
	}
	//Response for Create role
	role := asmodel.Role{
		ID:                 createRoleReq.ID,
		IsPredefined:       isPredefined,
		AssignedPrivileges: createRoleReq.AssignedPrivileges,
		OEMPrivileges:      createRoleReq.OEMPrivileges,
	}

	l.LogWithFields(ctx).Infof("Creating the role %s", createRoleReq.ID)
	//Persist role in database
	if cerr := role.Create(); cerr != nil {
		if errors.DBKeyAlreadyExist == cerr.ErrNo() {
			l.LogWithFields(ctx).Error(errorLogPrefix + cerr.Error())
			errorMessage := "Role with name " + role.ID + " already exists"
			args := response.Args{
				Code:    response.GeneralError,
				Message: errorMessage,
			}
			resp.StatusCode = http.StatusConflict
			resp.StatusMessage = response.GeneralError
			resp.Body = args.CreateGenericErrorResponse()
			return resp

		}
		errorMessage := errorLogPrefix + cerr.Error()
		l.LogWithFields(ctx).Error(errorMessage)
		resp.CreateInternalErrorResponse(errorMessage)
		return resp
	}

	resp.StatusCode = http.StatusCreated
	resp.StatusMessage = response.ResourceCreated

	commonResponse.CreateGenericResponse(resp.StatusMessage)
	commonResponse.ID = createRoleReq.ID

	resp.Body = asresponse.UserRole{
		Response:           commonResponse,
		IsPredefined:       role.IsPredefined,
		AssignedPrivileges: role.AssignedPrivileges,
		OEMPrivileges:      role.OEMPrivileges,
	}

	return resp
}

// CheckWhitespace func is used to check for whitepsace insde the string
// given func trims the spaces from the value and then checks if it is empty
func CheckWhitespace(fl validator.FieldLevel) bool {
	if len(strings.TrimSpace(fl.Field().String())) == 0 {
		return false
	}
	return true
}
