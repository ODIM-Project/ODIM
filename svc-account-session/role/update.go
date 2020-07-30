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
	"log"
	"net/http"
	"reflect"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
)

// Update defines the updation of the role details. Every role details can be
// updated other than the roleID if the session parameter have sufficient privileges.
//
// For updating an account,  parameters need to be passed are RoleRequest and Session.
// New RoleID,AssignedPrivileges and OEMPrivileges will be part of RoleRequest,
// and Session parameter will have all session related data, espically the privileges.
//
// There will be two return values for the fuction. One is the RPC response, which contains the
// status code, status message, headers and body.
func Update(req *roleproto.UpdateRoleRequest, session *asmodel.Session) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	var updateReq asmodel.Role
	json.Unmarshal(req.UpdateRequest, &updateReq)

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.UpdateRequest, updateReq)
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Println(errorMessage)
		resp := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, nil)
		return resp
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
		if req.Id == redfishrole {
			isPredefined = true
		}
	}

	if isPredefined {
		log.Println("Cannot update predefined role: ")
		resp.StatusCode = http.StatusMethodNotAllowed
		resp.StatusMessage = response.GeneralError
		errorMessage := "Updating predefined role is restricted"
		args := response.Args{
			Code:    response.GeneralError,
			Message: errorMessage,
		}
		resp.Body = args.CreateGenericErrorResponse()
		return resp
	}

	//check for ConfigureUsers privilege in session object
	status, err := checkForPrivilege(session, common.PrivilegeConfigureUsers)
	if err != nil {
		errorMessage := "error: user does not have the privilege to update the role"
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
	role, gerr := asmodel.GetRoleDetailsByID(req.Id)
	if gerr != nil {
		errorMessage := gerr.Error()
		resp.StatusCode = http.StatusBadRequest
		resp.StatusMessage = response.ResourceNotFound
		args := response.Args{
			Code:    response.GeneralError,
			Message: "",
			ErrorArgs: []response.ErrArgs{
				response.ErrArgs{
					StatusMessage: resp.StatusMessage,
					ErrorMessage:  errorMessage,
					MessageArgs:   []interface{}{"Role", req.Id},
				},
			},
		}
		resp.Body = args.CreateGenericErrorResponse()
		log.Printf(errorMessage)
		return resp
	}

	errorMessage := validateUpdateRequest(&updateReq, &role, map[string]bool{
		"AssignedPrivileges": true,
		"OEMPrivileges":      true,
	})
	if errorMessage != "" {
		log.Println(errorMessage)
		resp.StatusCode = http.StatusForbidden
		resp.StatusMessage = response.InsufficientPrivilege
		args := response.Args{
			Code:    response.GeneralError,
			Message: "",
			ErrorArgs: []response.ErrArgs{
				response.ErrArgs{
					StatusMessage: resp.StatusMessage,
					ErrorMessage:  errorMessage,
					MessageArgs:   []interface{}{},
				},
			},
		}
		resp.Body = args.CreateGenericErrorResponse()
		return resp
	}
	if len(updateReq.AssignedPrivileges) == 0 && len(updateReq.OEMPrivileges) == 0 {
		log.Println("Mandatory field is empty")
		errorMessage := "Mandatory field is empty"
		resp.StatusCode = http.StatusBadRequest
		resp.StatusMessage = response.PropertyMissing
		args := response.Args{
			Code:    response.GeneralError,
			Message: "",
			ErrorArgs: []response.ErrArgs{
				response.ErrArgs{
					StatusMessage: resp.StatusMessage,
					ErrorMessage:  errorMessage,
					MessageArgs:   []interface{}{"AssignedPrivileges/OemPrivileges"},
				},
			},
		}
		resp.Body = args.CreateGenericErrorResponse()
		return resp
	}

	if len(updateReq.AssignedPrivileges) != 0 {
		status, messageArgs, err := validateAssignedPrivileges(updateReq.AssignedPrivileges)
		if err != nil {
			errorMessage := err.Error()
			resp.StatusCode = int32(status.Code)
			resp.StatusMessage = status.Message
			args := response.Args{
				Code:    response.GeneralError,
				Message: "",
				ErrorArgs: []response.ErrArgs{
					response.ErrArgs{
						StatusMessage: resp.StatusMessage,
						ErrorMessage:  errorMessage,
						MessageArgs:   messageArgs,
					},
				},
			}
			resp.Body = args.CreateGenericErrorResponse()
			log.Printf(errorMessage)
			return resp
		}
		role.AssignedPrivileges = updateReq.AssignedPrivileges
	}
	if len(updateReq.OEMPrivileges) != 0 {
		status, messageArgs, err := validateOEMPrivileges(updateReq.OEMPrivileges)
		if err != nil {
			errorMessage := err.Error()
			resp.StatusCode = int32(status.Code)
			resp.StatusMessage = status.Message
			args := response.Args{
				Code:    response.GeneralError,
				Message: "",
				ErrorArgs: []response.ErrArgs{
					response.ErrArgs{
						StatusMessage: resp.StatusMessage,
						ErrorMessage:  errorMessage,
						MessageArgs:   messageArgs,
					},
				},
			}
			resp.Body = args.CreateGenericErrorResponse()
			log.Printf(errorMessage)
			return resp
		}
		role.OEMPrivileges = updateReq.OEMPrivileges
	}
	if uerr := role.UpdateRoleDetails(); uerr != nil {
		errorMessage := "error while trying to updating role:" + uerr.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		return resp
	}

	resp.Body = role
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	return resp
}

// validateUpdateRequest validates user update request for role  against store data in database
func validateUpdateRequest(req, data *asmodel.Role, exceptFields map[string]bool) string {
	val := reflect.ValueOf(data).Elem()
	reqFields := reflect.Indirect(reflect.ValueOf(req))
	var field = make([]string, 0)

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		if exceptFields[typeField.Name] {
			continue
		}

		datavalue := val.Field(i)
		reqValue := reqFields.FieldByName(typeField.Name)

		if reqValue.Interface() != "" {
			if datavalue.Interface() != reqValue.Interface() {
				field = append(field, typeField.Name)
			}
		}
	}
	if len(field) <= 0 {
		return ""
	}
	errorMessage := "error: user doesn't have privilege to edit this properties "
	for i := 0; i < len(field); i++ {
		errorMessage = errorMessage + field[i] + " "
	}
	return errorMessage
}
