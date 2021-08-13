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
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
	"net/http"
	"regexp"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	accountproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/account"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
)

// Create defines creation of a new account. The function is supposed to be used as part of RPC.
//
// For creating an account, two parameters need to be passed CreateAccountRequest and Session.
// New account UserName, Password and RoleID will be part of CreateAccountRequest,
// and Session parameter will have all session related data, espically the privileges.
// For creating new account the ConfigureUsers privilege is mandatory.
//
// There will be two return values for the fuction. One is the RPC response, which contains the
// status code, status message, headers and body and the second value is error.
func (e *ExternalInterface) Create(req *accountproto.CreateAccountRequest, session *asmodel.Session) (response.RPC, error) {
	// parsing the CreateAccount
	var createAccount asmodel.Account
	err := json.Unmarshal(req.RequestBody, &createAccount)
	if err != nil {
		errMsg := "Unable to parse the create account request" + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), fmt.Errorf(errMsg)
	}

	commonResponse := response.Response{
		OdataType:    common.ManagerAccountType,
		OdataID:      "/redfish/v1/AccountService/Accounts/" + createAccount.UserName,
		OdataContext: "/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
		ID:           createAccount.UserName,
		Name:         "Account Service",
	}
	var resp response.RPC

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, createAccount)
	if err != nil {
		errMsg := "While validating request parameters: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), fmt.Errorf(errMsg)
	} else if invalidProperties != "" {
		errorMessage := "One or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Error(errorMessage)
		resp := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, nil)
		return resp, fmt.Errorf(errorMessage)
	}

	user := asmodel.User{
		UserName: createAccount.UserName,
		Password: createAccount.Password,
		RoleID:   createAccount.RoleID,
	}

	if !(session.Privileges[common.PrivilegeConfigureUsers]) {
		errorMessage := "User does not have the privilege to create a new user"
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
		return resp, fmt.Errorf(errorMessage)
	}
	invalidParams := validateRequest(user)
	if invalidParams != "" {
		errorMessage := "Mandatory fields " + invalidParams + " are empty"
		resp.StatusCode = http.StatusBadRequest
		resp.StatusMessage = response.PropertyMissing
		args := response.Args{
			Code:    response.GeneralError,
			Message: "",
			ErrorArgs: []response.ErrArgs{
				response.ErrArgs{
					StatusMessage: resp.StatusMessage,
					ErrorMessage:  errorMessage,
					MessageArgs:   []interface{}{invalidParams},
				},
			},
		}
		resp.Body = args.CreateGenericErrorResponse()
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		}
		log.Error(errorMessage)
		return resp, fmt.Errorf(errorMessage)
	}
	if _, gerr := e.GetRoleDetailsByID(user.RoleID); gerr != nil {
		errorMessage := "Invalid RoleID present " + gerr.Error()
		log.Error(errorMessage)
		return common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, errorMessage, []interface{}{"Role", user.RoleID}, nil), fmt.Errorf(errorMessage)
	}
	if err := validatePassword(user.UserName, user.Password); err != nil {
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
					MessageArgs:   []interface{}{user.Password, "Password"},
				},
			},
		}
		resp.Body = args.CreateGenericErrorResponse()
		resp.Header = map[string]string{
			"Content-type": "application/json; charset=utf-8", // TODO: add all error headers
		}
		log.Error(errorMessage)
		return resp, err

	}
	hash := sha3.New512()
	hash.Write([]byte(user.Password))
	hashSum := hash.Sum(nil)
	hashedPassword := base64.URLEncoding.EncodeToString(hashSum)
	user.Password = hashedPassword
	user.AccountTypes = []string{"Redfish"}
	if cerr := e.CreateUser(user); cerr != nil {
		errorMessage := "Unable to add new user: " + cerr.Error()
		if errors.DBKeyAlreadyExist == cerr.ErrNo() {
			resp.StatusCode = http.StatusConflict
			resp.StatusMessage = response.ResourceAlreadyExists
			args := response.Args{
				Code:    response.GeneralError,
				Message: "",
				ErrorArgs: []response.ErrArgs{
					response.ErrArgs{
						StatusMessage: response.ResourceAlreadyExists,
						ErrorMessage:  errorMessage,
						MessageArgs:   []interface{}{"ManagerAccount", "Id", user.UserName},
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
		return resp, fmt.Errorf(errorMessage)
	}

	resp.StatusCode = http.StatusCreated
	resp.StatusMessage = response.Created

	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Link":              "</redfish/v1/AccountService/Accounts/" + user.UserName + "/>; rel=describedby",
		"Location":          "/redfish/v1/AccountService/Accounts/" + user.UserName,
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
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

	return resp, nil

}

func validateRequest(user asmodel.User) string {
	param := ""
	if user.UserName == "" {
		param = "UserName "
	}
	if user.Password == "" {
		param = param + "Password "
	}
	if user.RoleID == "" {
		param = param + "RoleID"
	}
	return param
}

func validatePassword(userName, password string) error {

	if strings.Contains(strings.ToLower(password), strings.ToLower(userName)) {
		return fmt.Errorf("error: invalid password, username is present inside the password")
	}
	if len(password) < config.Data.AuthConf.PasswordRules.MinPasswordLength {
		return fmt.Errorf("error: invalid password, password length is less than the minimum length")
	}
	if len(password) > config.Data.AuthConf.PasswordRules.MaxPasswordLength {
		return fmt.Errorf("error: invalid password, password length is greater than the maximum length")
	}
	matched, _ := regexp.Match("[A-Z]+", []byte(password))
	if !matched {
		return fmt.Errorf("error: invalid password, password should contain minimum One Upper case, One Lower case, One Number and One Special character")
	}
	matched, _ = regexp.Match("[a-z]+", []byte(password))
	if !matched {
		return fmt.Errorf("error: invalid password, password should contain minimum One Upper case, One Lower case, One Number and One Special character")
	}
	matched, _ = regexp.Match("[0-9]+", []byte(password))
	if !matched {
		return fmt.Errorf("error: invalid password, password should contain minimum One Upper case, One Lower case, One Number and One Special character")
	}
	matched, _ = regexp.Match("["+config.Data.AuthConf.PasswordRules.AllowedSpecialCharcters+"]+", []byte(password))
	if !matched {
		return fmt.Errorf("error: invalid password, password should contain minimum One Upper case, One Lower case, One Number and One Special character")
	}
	return nil
}
