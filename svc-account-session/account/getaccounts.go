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
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GetAllAccounts defines the admin functionality of listing of all accounts.
//
// As input parameters we need to pass Session, which contains all session data
// especially configureUsers privilege.
//
// As return parameters RPC response, which contains status code, message, headers and data,
// error will be passed back.
func GetAllAccounts(session *asmodel.Session) response.RPC {
	commonResponse := response.Response{
		OdataType:    "#ManagerAccountCollection.ManagerAccountCollection",
		OdataID:      "/redfish/v1/AccountService/Accounts",
		OdataContext: "/redfish/v1/$metadata#ManagerAccountCollection.ManagerAccountCollection",
		Name:         "Account Service",
	}

	var resp response.RPC

	if !session.Privileges[common.PrivilegeConfigureUsers] {
		errorMessage := "User " + session.UserName + " does not have the privilege to view all users"
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
	//Get all user keys
	users, err := asmodel.GetAllUsers()
	if err != nil {
		errorMessage := "Unable to get users: " + err.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		}
		log.Error(errorMessage)
		return resp
	}
	//Build response body and headers
	var accountLinks []asresponse.ListMember
	for _, key := range users {
		accountLink := asresponse.ListMember{
			OdataID: "/redfish/v1/AccountService/Accounts/" + key.UserName,
		}
		accountLinks = append(accountLinks, accountLink)
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	resp.Header = map[string]string{
		"Allow":             `"GET", "POST", "HEAD"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Link":              "</redfish/v1/SchemaStore/en/ManagerAccountCollection.json/>; rel=describedby",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}

	commonResponse.CreateGenericResponse(resp.StatusMessage)
	commonResponse.Message = ""
	commonResponse.ID = ""
	commonResponse.MessageID = ""
	commonResponse.Severity = ""
	resp.Body = asresponse.List{
		Response:     commonResponse,
		MembersCount: len(users),
		Members:      accountLinks,
	}

	return resp

}

// GetAccount defines the viewing of a particular account which is identified by the account id.
//
// As input parameters we need to pass Session, which contains all session data
// especially configureUsers privilege and the accountID which is used to
// identify the account which is supposed to be viewed.
//
// As return parameters RPC response, which contains status code, message, headers and data,
// error will be passed back.
func GetAccount(session *asmodel.Session, accountID string) response.RPC {
	commonResponse := response.Response{
		OdataType:    common.ManagerAccountType,
		OdataID:      "/redfish/v1/AccountService/Accounts/" + accountID,
		OdataContext: "/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
		ID:           accountID,
		Name:         "Account Service",
	}

	var resp response.RPC

	if !(session.Privileges[common.PrivilegeConfigureUsers]) {
		if accountID != session.UserName || !(session.Privileges[common.PrivilegeConfigureSelf]) {
			errorMessage := session.UserName + " does not have the privilege to view other user's details"
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

	user, err := asmodel.GetUserDetails(accountID)
	if err != nil {
		errorMessage := "Unable to get account: " + err.Error()
		if errors.DBKeyNotFound == err.ErrNo() {
			resp.StatusCode = http.StatusNotFound
			resp.StatusMessage = response.ResourceNotFound
			args := response.Args{
				Code:    response.GeneralError,
				Message: "",
				ErrorArgs: []response.ErrArgs{
					response.ErrArgs{
						StatusMessage: resp.StatusMessage,
						ErrorMessage:  errorMessage,
						MessageArgs:   []interface{}{"Account", accountID},
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

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	resp.Header = map[string]string{
		"Allow":             `"GET", "POST", "HEAD"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Link":              "</redfish/v1/SchemaStore/en/ManagerAccount.json/>; rel=describedby",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}

	commonResponse.CreateGenericResponse(resp.StatusMessage)
	commonResponse.Message = ""
	commonResponse.MessageID = ""
	commonResponse.Severity = ""
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

// GetAccountService defines the functionality for knowing whether
// the account service is enabled or not
//
// As return parameters RPC response, which contains status code, message, headers and data,
// error will be passed back.
func GetAccountService() response.RPC {
	commonResponse := response.Response{
		OdataType:    common.AccountServiceType,
		OdataID:      "/redfish/v1/AccountService",
		OdataContext: "/redfish/v1/$metadata#AccountService.AccountService",
		ID:           "AccountService",
		Name:         "Account Service",
	}
	var resp response.RPC

	isServiceEnabled := false
	serviceState := "Disabled"
	//Checks if AccountService is enabled and sets the variable isServiceEnabled to true add servicState to enabled
	for _, service := range config.Data.EnabledServices {
		if service == "AccountService" {
			isServiceEnabled = true
			serviceState = "Enabled"
		}
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	resp.Header = map[string]string{
		"Allow":         "GET",
		"Cache-Control": "no-cache",
		"Connection":    "Keep-alive",
		"Content-type":  "application/json; charset=utf-8",
		"Link": "	</redfish/v1/SchemaStore/en/AccountService.json>; rel=describedby",
		"Transfer-Encoding": "chunked",
		"X-Frame-Options":   "sameorigin",
	}

	commonResponse.CreateGenericResponse(resp.StatusMessage)
	commonResponse.Message = ""
	commonResponse.MessageID = ""
	commonResponse.Severity = ""
	resp.Body = asresponse.AccountService{
		Response: commonResponse,
		//TODO: Yet to implement AccountService state and health
		Status: asresponse.Status{
			State:  serviceState,
			Health: "OK",
		},
		ServiceEnabled:    isServiceEnabled,
		MinPasswordLength: config.Data.AuthConf.PasswordRules.MinPasswordLength,
		Accounts: asresponse.Accounts{
			OdataID: "/redfish/v1/AccountService/Accounts",
		},
		Roles: asresponse.Accounts{
			OdataID: "/redfish/v1/AccountService/Roles",
		},
	}

	return resp

}
