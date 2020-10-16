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

// Package session ...
package session

import (
	"encoding/json"
	"fmt"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/asresponse"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
)

// CreateNewSession is a method to to create a session
// it will accepts the SessionCreateRequest which will have username and password
// and check whether the credentials are correct also it will
// check privileges. and then add the session details in DB
// respond RPC response and error if there is.
func CreateNewSession(req *sessionproto.SessionCreateRequest) (response.RPC, string) {
	commonResponse := response.Response{
		OdataType: "#SessionService.v1_1_6.SessionService",
		OdataID:   "/redfish/v1/SessionService/Sessions",
		ID:        "Sessions",
		Name:      "Session Service",
	}
	var resp response.RPC

	// parsing the CreateSession
	var createSession asmodel.CreateSession
	genErr := json.Unmarshal(req.RequestBody, &createSession)
	if genErr != nil {
		errMsg := "unable to parse the create session request" + genErr.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), ""
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, genErr := common.RequestParamsCaseValidator(req.RequestBody, createSession)
	if genErr != nil {
		errMsg := "error while validating request parameters: " + genErr.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), ""
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Println(errorMessage)
		resp := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, nil)
		return resp, ""
	}

	user, err := auth.CheckSessionCreationCredentials(createSession.UserName, createSession.Password)
	if err != nil {
		errMsg := "error while authorizing session creation credentials: " + err.Error()
		if err.ErrNo() == errors.DBConnFailed {
			msgArgs := []interface{}{fmt.Sprintf("%v:%v", config.Data.DBConf.OnDiskHost, config.Data.DBConf.OnDiskPort)}
			resp = common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, errMsg, msgArgs, nil)
		} else {
			resp = common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil)
		}
		log.Printf(errMsg)
		return resp, ""
	}

	role, err := asmodel.GetRoleDetailsByID(user.RoleID)
	if err != nil {
		errorMessage := "error while trying to get role privileges for session creation: " + err.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		log.Printf(errorMessage)
		return resp, ""
	}
	rolePrivilege := make(map[string]bool)
	for _, privilege := range role.AssignedPrivileges {
		rolePrivilege[privilege] = true
	}
	//User requires Login privelege to create a session
	if _, exist := rolePrivilege[common.PrivilegeLogin]; !exist {
		errorMessage := "user doesn't have required privilege to create a session"
		log.Println(errorMessage)
		return common.GeneralError(http.StatusForbidden, response.InsufficientPrivilege, errorMessage, nil, nil), ""
	}

	currentTime := time.Now()
	sess := asmodel.Session{
		ID:           uuid.NewV4().String(),
		Token:        uuid.NewV4().String(),
		UserName:     user.UserName,
		RoleID:       user.RoleID,
		Privileges:   rolePrivilege,
		CreatedTime:  currentTime,
		LastUsedTime: currentTime,
	}
	auth.Lock.Lock()
	defer auth.Lock.Unlock()
	if err = sess.Persist(); err != nil {
		errMsg := "error while trying to insert session details: " + err.Error()
		if err.ErrNo() == errors.DBConnFailed {
			msgArgs := []interface{}{fmt.Sprintf("%v:%v", config.Data.DBConf.InMemoryHost, config.Data.DBConf.InMemoryPort)}
			resp = common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, errMsg, msgArgs, nil)
		} else {
			resp = common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
		}
		log.Printf(errMsg)
		return resp, ""
	}

	resp.StatusCode = http.StatusCreated
	resp.StatusMessage = response.Created

	resp.StatusCode = http.StatusCreated
	resp.StatusMessage = response.Created
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Link":              "</redfish/v1/SessionService/Sessions/" + sess.ID + "/>; rel=self",
		"Transfer-Encoding": "chunked",
		"X-Auth-Token":      sess.Token,
		"Content-type":      "application/json; charset=utf-8",
	}

	commonResponse.ID = sess.ID
	commonResponse.OdataID = "/redfish/v1/SessionService/Sessions/" + commonResponse.ID
	commonResponse.CreateGenericResponse(resp.StatusMessage)
	resp.Body = asresponse.Session{
		Response: commonResponse,
		UserName: createSession.UserName,
	}

	return resp, commonResponse.ID
}
