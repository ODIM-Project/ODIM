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

// Package account ...
package account

// ---------------------------------------------------------------------------------------
// IMPORT Section
// ---------------------------------------------------------------------------------------
import (
	"encoding/base64"
	"encoding/json"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	accountproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/account"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
	"net/http"
)

// Update defines the updation of the account details. Every account details can be
// updated other than the UserName if the session parameter have sufficient privileges.
//
// For updating an account, two parameters need to be passed UpdateAccountRequest and Session.
// New Password and RoleID will be part of UpdateAccountRequest,
// and Session parameter will have all session related data, espically the privileges.
//
// Output is the RPC response, which contains the status code, status message, headers and body.
func (e *ExternalInterface) Update(req *accountproto.UpdateAccountRequest, session *asmodel.Session) response.RPC {
	commonResponse := response.Response{
		OdataType:    common.ManagerAccountType,
		OdataID:      "/redfish/v1/AccountService/Accounts/" + req.AccountID,
		OdataContext: "/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
		ID:           req.AccountID,
		Name:         "Account Service",
	}

	var (
		resp response.RPC
		err  error
	)

	// parsing the Account
	var updateAccount asmodel.Account
	err = json.Unmarshal(req.RequestBody, &updateAccount)
	if err != nil {
		errMsg := "unable to parse the update account request" + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}

	requestUser := asmodel.User{
		UserName:     updateAccount.UserName,
		Password:     updateAccount.Password,
		RoleID:       updateAccount.RoleID,
		AccountTypes: []string{"Redfish"},
	}

	//empty request check
	if isEmptyRequest(req.RequestBody) {
		errMsg := "empty request can not be processed"
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"request body"}, nil)
	}

	id := req.AccountID
	if requestUser.UserName != "" {
		errorMessage := "Username cannot be modified"
		resp.StatusCode = http.StatusBadRequest
		resp.StatusMessage = response.GeneralError
		args := response.Args{
			Code:    response.GeneralError,
			Message: errorMessage,
		}
		resp.Body = args.CreateGenericErrorResponse()
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		}
		log.Error(errorMessage)
		return resp
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, updateAccount)
	if err != nil {
		errMsg := "Request parameters validaton failed: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	} else if invalidProperties != "" {
		errorMessage := "One or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Error(errorMessage)
		resp := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, nil)
		return resp
	}

	if requestUser.RoleID != "" {
		if requestUser.RoleID != common.RoleAdmin {
			if requestUser.RoleID != common.RoleMonitor {
				if requestUser.RoleID != common.RoleClient {
					_, err := e.GetRoleDetailsByID(requestUser.RoleID)
					if err != nil {
						errorMessage := "Invalid RoleID " + requestUser.RoleID + " present"
						resp.StatusCode = http.StatusBadRequest
						resp.StatusMessage = response.PropertyValueNotInList
						args := response.Args{
							Code:    response.GeneralError,
							Message: "",
							ErrorArgs: []response.ErrArgs{
								response.ErrArgs{
									StatusMessage: resp.StatusMessage,
									ErrorMessage:  errorMessage,
									MessageArgs:   []interface{}{requestUser.RoleID, "RoleID"},
								},
							},
						}
						resp.Body = args.CreateGenericErrorResponse()
						resp.Header = map[string]string{
							"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
						}
						log.Error(errorMessage)
						return resp
					}
				}
			}
		}

	}

	user, gerr := e.GetUserDetails(id)
	if gerr != nil {
		errorMessage := "Unable to get account: " + gerr.Error()
		if errors.DBKeyNotFound == gerr.ErrNo() {
			resp.StatusCode = http.StatusNotFound
			resp.StatusMessage = response.ResourceNotFound
			args := response.Args{
				Code:    response.GeneralError,
				Message: "",
				ErrorArgs: []response.ErrArgs{
					response.ErrArgs{
						StatusMessage: resp.StatusMessage,
						ErrorMessage:  errorMessage,
						MessageArgs:   []interface{}{"Account", id},
					},
				},
			}
			resp.Body = args.CreateGenericErrorResponse()
		} else {
			resp.CreateInternalErrorResponse(errorMessage)
		}
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		}
		log.Error(errorMessage)
		return resp
	}

	if user.UserName != session.UserName && !session.Privileges[common.PrivilegeConfigureUsers] {
		errorMessage := "User does not have the privilege to update other accounts"
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
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		}
		log.Error(errorMessage)
		return resp
	}

	//To be discussed
	// Check if the user trying to update RoleID, if so check if he has PrivilegeConfigureUsers Privilege,
	// if not return 403 forbidden.
	// Without PrivilegeConfigureUsers user is not allowed to update any user account roleID, including his own account roleID
	if requestUser.RoleID != "" {
		if !session.Privileges[common.PrivilegeConfigureUsers] {
			errorMessage := "User does not have the privilege to update any account role, including his own account"
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
			resp.Header = map[string]string{
				"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
			}
			log.Error(errorMessage)
			return resp
		}
	}

	if requestUser.Password != "" {
		// Password modification not allowed, if user doesn't have ConfigureSelf or ConfigureUsers privilege
		if !session.Privileges[common.PrivilegeConfigureSelf] && !session.Privileges[common.PrivilegeConfigureUsers] {
			errorMessage := "Roles, user is associated with, doesn't allow changing own or other users password"
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
			resp.Header = map[string]string{
				"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
			}
			log.Error(errorMessage)
			return resp
		}

		//TODO: handle all the combination of patch roles(admin,non-admin,default admin, non-default admin)
		if err = validatePassword(user.UserName, requestUser.Password); err != nil {
			errorMessage := err.Error()
			resp.StatusCode = http.StatusBadRequest
			resp.StatusMessage = response.PropertyValueFormatError
			args := response.Args{
				Code:    response.GeneralError,
				Message: "",
				ErrorArgs: []response.ErrArgs{
					response.ErrArgs{
						StatusMessage: resp.StatusMessage,
						ErrorMessage:  errorMessage,
						MessageArgs:   []interface{}{requestUser.Password, "Password"},
					},
				},
			}
			resp.Body = args.CreateGenericErrorResponse()
			resp.Header = map[string]string{
				"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
			}
			log.Error(errorMessage)
			return resp
		}
		hash := sha3.New512()
		hash.Write([]byte(requestUser.Password))
		hashSum := hash.Sum(nil)
		hashedPassword := base64.URLEncoding.EncodeToString(hashSum)
		requestUser.Password = hashedPassword
	}

	if uerr := e.UpdateUserDetails(user, requestUser); uerr != nil {
		errorMessage := "Unable to update user: " + uerr.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		}
		log.Error(errorMessage)
		return resp
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.AccountModified

	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Link":              "</redfish/v1/AccountService/Accounts/" + user.UserName + "/>; rel=describedby",
		"Location":          "/redfish/v1/AccountService/Accounts/" + user.UserName,
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	if requestUser.RoleID != "" {
		user.RoleID = requestUser.RoleID
	}
	commonResponse.CreateGenericResponse(resp.StatusMessage)
	resp.Body = asresponse.Account{
		Response:     commonResponse,
		UserName:     user.UserName,
		RoleID:       user.RoleID,
		AccountTypes: user.AccountTypes,
		Links: asresponse.Links{
			Role: asresponse.Role{
				OdataID: "/redfish/v1/AccountService/Roles/" + user.RoleID,
			},
		},
	}

	return resp
}

func isEmptyRequest(requestBody []byte) bool {
	var updateRequest map[string]interface{}
	json.Unmarshal(requestBody, &updateRequest)
	if len(updateRequest) <= 0 {
		return true
	}
	return false
}
