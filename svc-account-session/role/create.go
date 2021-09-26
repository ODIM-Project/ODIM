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
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
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
func Create(req *roleproto.RoleRequest, session *asmodel.Session) response.RPC {
	// parsing the request body
	var createRoleReq asmodel.Role
	err := json.Unmarshal(req.RequestBody, &createRoleReq)
	if err != nil {
		errMsg := "unable to parse the add request" + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}

	commonResponse := response.Response{
		OdataType: common.RoleType,
		OdataID:   "/redfish/v1/AccountService/Roles/" + createRoleReq.ID,
		Name:      "User Role",
	}
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, createRoleReq)
	if err != nil {
		errMsg := "Unable to validate request parameters: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	} else if invalidProperties != "" {
		errorMessage := "One or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Error(errorMessage)
		resp := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, nil)
		return resp
	}

	// register the validator package
	validate := validator.New()
	// create validation to check if the ID field contains whitespaces when creating role.
	validate.RegisterValidation("is-empty", CheckWhitespace)
	if err := validate.Var(createRoleReq.ID, "required"); err != nil {
		errorMessage := "error: Mandatory field RoleId is empty"
		response := common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errorMessage, []interface{}{"RoleId"}, nil)
		return response
	} else if err := validate.Var(createRoleReq.ID, "required,is-empty,excludesall=!@#?%&$*"); err != nil {
		errorMessage := "Invalid create role request"
		args := response.Args{
			Code:    response.GeneralError,
			Message: "",
			ErrorArgs: []response.ErrArgs{
				response.ErrArgs{
					StatusMessage: response.PropertyValueNotInList,
					ErrorMessage:  errorMessage,
					MessageArgs:   []interface{}{createRoleReq.ID, "RoleId"},
				},
			},
		}
		resp.StatusCode = http.StatusBadRequest
		resp.StatusMessage = response.PropertyValueNotInList
		resp.Body = args.CreateGenericErrorResponse()
		log.Error(errorMessage)
		return resp
	}

	//check for ConfigureUsers privilege in session object
	status, err := checkForPrivilege(session, "ConfigureUsers")
	if err != nil {
		errorMessage := "User does not have the privilege to create a new role"
		resp.StatusCode = int32(status.Code)
		resp.StatusMessage = status.Message
		args := response.Args{
			Code:    response.GeneralError,
			Message: "",
			ErrorArgs: []response.ErrArgs{
				response.ErrArgs{
					StatusMessage: status.Message,
					ErrorMessage:  errorMessage,
					MessageArgs:   []interface{}{},
				},
			},
		}
		resp.Body = args.CreateGenericErrorResponse()
		log.Error(errorMessage)
		return resp
	}
	if len(createRoleReq.AssignedPrivileges) == 0 && len(createRoleReq.OEMPrivileges) == 0 {
		errorMessage := "Both AssignedPrivileges and OemPrivileges cannot be empty."
		args := response.Args{
			Code:    response.GeneralError,
			Message: "",
			ErrorArgs: []response.ErrArgs{
				response.ErrArgs{
					StatusMessage: response.PropertyMissing,
					ErrorMessage:  errorMessage,
					MessageArgs:   []interface{}{"AssignedPrivileges/OemPrivileges"},
				},
			},
		}
		resp.StatusCode = http.StatusBadRequest
		resp.StatusMessage = response.PropertyMissing
		resp.Body = args.CreateGenericErrorResponse()
		log.Error(errorMessage)
		return resp
	}

	if len(createRoleReq.AssignedPrivileges) != 0 {
		status, messageArgs, err := validateAssignedPrivileges(createRoleReq.AssignedPrivileges)
		if err != nil {
			errorMessage := err.Error()
			resp.StatusCode = int32(status.Code)
			resp.StatusMessage = status.Message
			args := response.Args{
				Code:    response.GeneralError,
				Message: "",
				ErrorArgs: []response.ErrArgs{
					response.ErrArgs{
						StatusMessage: status.Message,
						ErrorMessage:  errorMessage,
						MessageArgs:   messageArgs,
					},
				},
			}
			resp.Body = args.CreateGenericErrorResponse()
			log.Error(errorMessage)
			return resp
		}
	}
	if len(createRoleReq.OEMPrivileges) != 0 {
		status, messageArgs, err := validateOEMPrivileges(createRoleReq.OEMPrivileges)
		if err != nil {
			errorMessage := err.Error()
			resp.StatusCode = int32(status.Code)
			resp.StatusMessage = status.Message
			args := response.Args{
				Code:    response.GeneralError,
				Message: "",
				ErrorArgs: []response.ErrArgs{
					response.ErrArgs{
						StatusMessage: status.Message,
						ErrorMessage:  errorMessage,
						MessageArgs:   messageArgs,
					},
				},
			}
			resp.Body = args.CreateGenericErrorResponse()
			log.Error(errorMessage)
			return resp
		}
	}
	//Get redfish roles from database
	redfishRoles, gerr := asmodel.GetRedfishRoles()
	if gerr != nil {
		log.Error("Unable to get redfish roles: " + gerr.Error())
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
		errorMessage := "Cannot create pre-defined roles"
		args := response.Args{
			Code:    response.GeneralError,
			Message: "",
			ErrorArgs: []response.ErrArgs{
				response.ErrArgs{
					StatusMessage: response.InsufficientPrivilege,
					ErrorMessage:  errorMessage,
					MessageArgs:   []interface{}{},
				},
			},
		}
		resp.StatusCode = http.StatusForbidden
		resp.StatusMessage = response.InsufficientPrivilege
		resp.Body = args.CreateGenericErrorResponse()
		log.Error(errorMessage)
		return resp
	}
	//Response for Create role
	role := asmodel.Role{
		ID:                 createRoleReq.ID,
		IsPredefined:       isPredefined,
		AssignedPrivileges: createRoleReq.AssignedPrivileges,
		OEMPrivileges:      createRoleReq.OEMPrivileges,
	}

	//Persist role in database
	if cerr := role.Create(); cerr != nil {
		if errors.DBKeyAlreadyExist == cerr.ErrNo() {
			log.Error("Unable to create new role: " + cerr.Error())
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
		log.Error("Unable to create new role: " + cerr.Error())
		errorMessage := "Unable to create new role: " + cerr.Error()
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
