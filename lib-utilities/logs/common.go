//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

// Package logs ...
package logs

import (
	"encoding/json"
	srv "github.com/ODIM-Project/ODIM/lib-utilities/services"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// getUserDetails function
// getting the session user id and role id for a given session token
func getUserDetails(sessionToken string) (string, string) {
	var err error
	sessionUserName := "null"
	sessionRoleID := "null"
	if sessionToken != "" {
		sessionUserName, err = srv.GetSessionUserName(sessionToken)
		if err != nil {
			errMsg := "while trying to get session details: " + err.Error()
			log.Error(errMsg)
			return "null", "null"
		}
		sessionRoleID, err = srv.GetSessionUserRoleID(sessionToken)
		if err != nil {
			errMsg := "while trying to get session details: " + err.Error()
			log.Error(errMsg)
			return sessionUserName, "null"
		}
	}
	return sessionUserName, sessionRoleID
}

// maskRequestBody function
// masking the request body, making password as null
func maskRequestBody(reqBody map[string]interface{}) string {
	var jsonStr []byte
	var err error
	if len(reqBody) > 0 {
		reqBody["Password"] = "null"
		jsonStr, err = json.Marshal(reqBody)
		if err != nil {
			log.Error("while marshalling request body", err.Error())
		}
	}
	reqStr := string(jsonStr)
	// adding null to requestbody property if no payload is sent
	if reqStr == "" {
		reqStr = "null"
	}
	return reqStr
}

// getResponseStatus function
// setting operation status flag based on the response code
func getResponseStatus(respStatusCode int32) bool {
	operationStatus := false
	successStatusCodes := []int32{http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent}
	for _, statusCode := range successStatusCodes {
		if statusCode == respStatusCode {
			operationStatus = true
			break
		}
	}
	return operationStatus
}
