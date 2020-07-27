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
	"log"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
)

func getHeader() map[string]string {
	return map[string]string{
		"Cache-Control":     "no-cache",
		"Transfer-Encoding": "chunked",
		"Content-type":      "application/json; charset=utf-8",
	}
}

// DeleteSession is a method to delete a sessiom
// it will accepts the SessionCreateRequest which will have sessionid and sessiontoken
// and it will check privileges to delete session and then delete the session
// respond RPC response and error if there is.
func DeleteSession(req *sessionproto.SessionRequest) response.RPC {
	var resp response.RPC
	errorArgs := []response.ErrArgs{
		response.ErrArgs{
			StatusMessage: "",
			ErrorMessage:  "",
			MessageArgs:   []interface{}{},
		},
	}
	args := &response.Args{
		Code:      response.GeneralError,
		Message:   "",
		ErrorArgs: errorArgs,
	}
	currentSession, serr := asmodel.GetSession(req.SessionToken)
	if serr != nil {
		errorMessage := "error: while trying to delete session: " + serr.Error()
		log.Printf(errorMessage)
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
	}

	sessionTokens, err := asmodel.GetAllSessionKeys()
	if err != nil {
		errorMessage := "error:  while trying to get all session keys in delete session: " + err.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		resp.Header = getHeader()
		log.Printf(errorMessage)
		return resp
	}
	for _, token := range sessionTokens {
		session, err := auth.CheckSessionTimeOut(token)
		if err != nil {
			log.Printf("error while trying to get session details with the token %v: %v", token, err)
			continue
		}
		if session.ID == req.SessionId {
			hasprivilege := checkPrivilege(req.SessionToken, session, &currentSession)
			if hasprivilege {
				if req.SessionToken != session.Token {
					err := UpdateLastUsedTime(req.SessionToken)
					if err != nil {
						errorMessage := "error while updating last used time of session with token " + req.SessionToken + ": " + err.Error()
						resp.CreateInternalErrorResponse(errorMessage)
						resp.Header = getHeader()
						log.Printf(errorMessage)
						return resp
					}
				}
				if err := session.Delete(); err != nil {
					errorMessage := "error:  while trying to get all session keys in delete session: " + err.Error()
					resp.CreateInternalErrorResponse(errorMessage)
					resp.Header = getHeader()
					log.Printf(errorMessage)
					return resp
				}
				log.Printf("Successfully Deleted: %v", err)
				resp.StatusCode = http.StatusNoContent
				resp.StatusMessage = response.ResourceRemoved
				resp.Header = getHeader()
				return resp
			}
			errorMessage := "Insufficient privileges"
			resp.StatusCode = http.StatusForbidden
			resp.StatusMessage = response.InsufficientPrivilege
			errorArgs[0].ErrorMessage = errorMessage
			errorArgs[0].StatusMessage = resp.StatusMessage
			resp.Header = getHeader()
			resp.Body = args.CreateGenericErrorResponse()
			log.Printf(errorMessage)
			return resp
		}
	}
	sessionTokens = nil
	log.Printf("error: Status Not Found")
	errorMessage := "error: Session ID not found"
	resp.StatusCode = http.StatusNotFound
	resp.StatusMessage = response.ResourceNotFound
	errorArgs[0].ErrorMessage = errorMessage
	errorArgs[0].StatusMessage = resp.StatusMessage
	errorArgs[0].MessageArgs = []interface{}{"Session", req.SessionId}
	resp.Body = args.CreateGenericErrorResponse()
	resp.Header = getHeader()
	return resp
}

func checkPrivilege(sessionToken string, session, currentSession *asmodel.Session) bool {
	if (session.UserName == currentSession.UserName && currentSession.Privileges[common.PrivilegeConfigureSelf]) ||
		currentSession.Privileges[common.PrivilegeConfigureUsers] {
		return true
	}
	return false
}
