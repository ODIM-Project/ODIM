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
	"log"
	"net/http"
	"strings"

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
	commonResponse := response.Response{
		OdataType: "#Role.v1_2_4.Role",
		OdataID:   "/redfish/v1/AccountService/Roles/" + req.RoleId,
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

	// register the validator package
	validate := validator.New()
	// create validation to check if the ID field contains whitespaces when creating role.
	validate.RegisterValidation("is-empty", CheckWhitespace)
	if err := validate.Var(req.RoleId, "required,is-empty,excludesall=!@#?%&$*"); err != nil {
		errorMessage := "error: Invalid Request"
		args := response.Args{
			Code:    response.GeneralError,
			Message: "",
			ErrorArgs: []response.ErrArgs{
				response.ErrArgs{
					StatusMessage: response.PropertyValueNotInList,
					ErrorMessage:  errorMessage,
					MessageArgs:   []interface{}{req.RoleId, "RoleId"},
				},
			},
		}
		resp.StatusCode = http.StatusBadRequest
		resp.StatusMessage = response.PropertyValueNotInList
		resp.Body = args.CreateGenericErrorResponse()
		log.Printf(errorMessage)
		return resp
	}

	//check for ConfigureUsers privilege in session object
	status, err := checkForPrivilege(session, "ConfigureUsers")
	if err != nil {
		errorMessage := "error: user does not have the privilege to create a new role"
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
		log.Printf(errorMessage)
		return resp
	}
	if len(req.AssignedPrivileges) == 0 && len(req.OemPrivileges) == 0 {
		errorMessage := "error: Both AssignedPrivileges and OemPrivileges cannot be empty."
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
		log.Printf(errorMessage)
		return resp
	}

	if len(req.AssignedPrivileges) != 0 {
		status, messageArgs, err := validateAssignedPrivileges(req.AssignedPrivileges)
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
			log.Printf(errorMessage)
			return resp
		}
	}
	if len(req.OemPrivileges) != 0 {
		status, messageArgs, err := validateOEMPrivileges(req.OemPrivileges)
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
			log.Printf(errorMessage)
			return resp
		}
	}
	//Get redfish roles from database
	redfishRoles, gerr := asmodel.GetRedfishRoles()
	if gerr != nil {
		log.Println("error getting redfish roles: ", gerr.Error())
		errorMessage := gerr.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		return resp

	}
	//check if request role is predefined redfish role and set isPredfined to true
	isPredefined := false
	for _, redfishrole := range redfishRoles.List {
		if req.RoleId == redfishrole {
			isPredefined = true
		}
	}
	if isPredefined {
		errorMessage := "error: cannot create pre-defined roles"
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
		log.Printf(errorMessage)
		return resp
	}
	//Response for Create role
	role := asmodel.Role{
		ID:                 req.RoleId,
		IsPredefined:       isPredefined,
		AssignedPrivileges: req.AssignedPrivileges,
		OEMPrivileges:      req.OemPrivileges,
	}

	//Persist role in database
	if cerr := role.Create(); cerr != nil {
		if errors.DBKeyAlreadyExist == cerr.ErrNo() {
			log.Println("error while trying to create new role: ", cerr.Error())
			errorMessage := "error: role with name " + role.ID + " already exists"
			args := response.Args{
				Code:    response.GeneralError,
				Message: errorMessage,
			}
			resp.StatusCode = http.StatusConflict
			resp.StatusMessage = response.GeneralError
			resp.Body = args.CreateGenericErrorResponse()
			return resp

		}
		log.Println("error while trying to create new role: ", cerr.Error())
		errorMessage := "error while trying to create new role: " + cerr.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		return resp
	}

	resp.StatusCode = http.StatusCreated
	resp.StatusMessage = response.ResourceCreated

	commonResponse.CreateGenericResponse(resp.StatusMessage)
	commonResponse.ID = req.RoleId

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
