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

	"github.com/bharath-b-hpe/odimra/lib-utilities/common"
	"github.com/bharath-b-hpe/odimra/lib-utilities/config"
	"github.com/bharath-b-hpe/odimra/lib-utilities/errors"
	roleproto "github.com/bharath-b-hpe/odimra/lib-utilities/proto/role"
	"github.com/bharath-b-hpe/odimra/lib-utilities/response"
	"github.com/bharath-b-hpe/odimra/svc-account-session/asmodel"
	"github.com/bharath-b-hpe/odimra/svc-account-session/asresponse"
)

//GetRole defines the viewing of a particular role which is identified by the id.
//
// As input parameters we need to pass Session, which contains all session data
// especially configureUsers privilege and the roleID which is used to
// identify the role which is supposed to be viewed.
//
// As return parameters RPC response, which contains status code, message, headers and data.
func GetRole(req *roleproto.GetRoleRequest, session *asmodel.Session) response.RPC {
	commonResponse := response.Response{
		OdataType: "#Role.v1_2_4.Role",
		OdataID:   "/redfish/v1/AccountService/Roles/" + req.Id,
		Name:      "User Role",
		ID:        req.Id,
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
	//check for ConfigureUsers privilege in session object
	status, perr := checkForPrivilege(session, "ConfigureUsers")
	if perr != nil {
		errorMessage := "error: user does not have the privilege to get the role"
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
	//Get role from database using role ID
	role, err := asmodel.GetRoleDetailsByID(req.Id)
	if err != nil {
		errorMessage := "error while getting the role : " + err.Error()
		log.Printf(errorMessage)
		if errors.DBKeyNotFound == err.ErrNo() {
			resp.StatusCode = http.StatusNotFound
			resp.StatusMessage = response.ResourceNotFound
			messageArgs := []interface{}{"Role", req.Id}
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
		} else {
			resp.CreateInternalErrorResponse(errorMessage)
		}
		return resp
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	commonResponse.CreateGenericResponse(resp.StatusMessage)
	commonResponse.MessageID = ""
	commonResponse.Message = ""
	commonResponse.Severity = ""

	resp.Body = asresponse.UserRole{
		Response:           commonResponse,
		IsPredefined:       role.IsPredefined,
		AssignedPrivileges: role.AssignedPrivileges,
		OEMPrivileges:      role.OEMPrivileges,
	}

	return resp
}

//GetAllRoles defines the  functionality of listing of all roles.
//
// As input parameters we need to pass Session, which contains all session data
// especially configureUsers privilege.
//
// As return parameters RPC response, which contains status code, message, headers and data.
func GetAllRoles(session *asmodel.Session) response.RPC {
	var resp response.RPC
	commonResponse := response.Response{
		OdataType: "#RoleCollection.RoleCollection",
		OdataID:   "/redfish/v1/AccountService/Roles",
		Name:      "Roles Collection",
	}
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	//check for ConfigureUsers privilege in session object
	status, err := checkForPrivilege(session, "ConfigureUsers")
	if err != nil {
		errorMessage := "error: user does not have the privilege to get the roles"
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
		log.Printf(errorMessage)
		resp.Body = args.CreateGenericErrorResponse()
		return resp
	}
	roles, rerr := asmodel.GetAllRoles()
	if rerr != nil {
		log.Printf("error getting role : %v", rerr.Error())
		errorMessage := rerr.Error()
		return common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, errorMessage, []interface{}{config.Data.DBConf.OnDiskHost + ":" + config.Data.DBConf.OnDiskPort}, nil)
	}
	//Build response body and headers
	var roleLinks []asresponse.ListMember
	for _, key := range roles {
		roleLink := asresponse.ListMember{
			OdataID: "/redfish/v1/AccountService/Roles/" + key.ID,
		}
		roleLinks = append(roleLinks, roleLink)
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	commonResponse.CreateGenericResponse(resp.StatusMessage)
	commonResponse.MessageID = ""
	commonResponse.Message = ""
	commonResponse.Severity = ""
	resp.Body = asresponse.List{
		Response:     commonResponse,
		MembersCount: len(roles),
		Members:      roleLinks,
	}

	return resp
}
