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
	"fmt"
	"github.com/kataras/iris/v12"
	"net/http"
	"time"
)

// AuditLog is used for generating audit logs in syslog format for each request
// this function logs an info for successful operation and error for failure operation
// properties logged are prival, time, host, username, roleid, request method, resource, requestbody, responsecode and message
func AuditLog(ctx iris.Context, reqBody map[string]interface{}) {
	logMsg := auditLogEntry(ctx, reqBody)
	// Get response code
	respStatusCode := int32(ctx.GetStatusCode())
	operationStatus := getResponseStatus(respStatusCode)

	// 110 is for audit log info
	// 107 is for audit log error
	if operationStatus {
		successMsg := "<110> " + logMsg + " Operation successful"
		fmt.Println(successMsg)
	} else {
		failedMsg := "<107> " + logMsg + " Operation failed"
		fmt.Println(failedMsg)
	}
}

// AuthLog is used for generating security logs in syslog format for each request
// this function logs an info for successful operation and warning for failure auth operation
// properties logged are prival, time, username, roleid and message
func AuthLog(logProperties map[string]interface{}) {
	sessionToken := "null"
	sessionUserName := "null"
	sessionRoleID := "null"
	msg := "null"
	respStatusCode := int32(http.StatusUnauthorized)
	tokenMsg := ""

	if logProperties["SessionToken"] != nil {
		sessionToken = logProperties["SessionToken"].(string)
	}
	if logProperties["SessionUserID"] != nil {
		sessionUserName = logProperties["SessionUserID"].(string)
	}
	if logProperties["SessionRoleID"] != nil {
		sessionRoleID = logProperties["SessionRoleID"].(string)
	}
	if logProperties["Message"] != nil {
		msg = logProperties["Message"].(string)
	}
	if logProperties["ResponseStatusCode"] != nil {
		respStatusCode = logProperties["ResponseStatusCode"].(int32)
	}

	timeNow := time.Now().Format(time.RFC3339)
	// formatting logs in syslog format
	logMsg := fmt.Sprintf("%s [account@1 user=\"%s\" roleID=\"%s\"]", timeNow, sessionUserName, sessionRoleID)
	// Get response code
	operationStatus := getResponseStatus(respStatusCode)
	if sessionToken != "null" {
		tokenMsg = "for session token " + sessionToken
	}
	// 86 is for auth log info
	// 84 is for auth log warning
	if operationStatus {
		successMsg := fmt.Sprintf("%s %s %s %s", "<86>", logMsg, "Authentication/Authorization successful", tokenMsg)
		fmt.Println(successMsg)
	} else {
		errMsg := "Authentication/Authorization failed"
		if respStatusCode == http.StatusForbidden {
			errMsg = "Authorization failed"
		} else if respStatusCode == http.StatusUnauthorized {
			errMsg = "Authentication failed"
		}
		failedMsg := fmt.Sprintf("%s %s %s %s, %s", "<84>", logMsg, errMsg, tokenMsg, msg)
		fmt.Println(failedMsg)
	}
}

// auditLogEntry extracts the required info from context like session token, username, request URI
// and formats in syslog format for audit logs
func auditLogEntry(ctx iris.Context, reqBody map[string]interface{}) string {
	var logMsg string
	// getting the request URI, host and method from context
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	sessionUserName, sessionRoleID := getUserDetails(sessionToken)
	rawURI := ctx.Request().RequestURI
	host := ctx.Request().Host
	method := ctx.Request().Method
	respStatusCode := ctx.GetStatusCode()
	timeNow := time.Now().Format(time.RFC3339)
	reqStr := maskRequestBody(reqBody)

	// formatting logs in syslog format
	if reqStr == "null" {
		logMsg = fmt.Sprintf("%s %s [account@1 user=\"%s\" roleID=\"%s\"][request@1 method=\"%s\" resource=\"%s\"][response@1 responseCode=%d]", timeNow, host, sessionUserName, sessionRoleID, method, rawURI, respStatusCode)
	} else {
		logMsg = fmt.Sprintf("%s %s [account@1 user=\"%s\" roleID=\"%s\"][request@1 method=\"%s\" resource=\"%s\" requestBody=\"%s\"][response@1 responseCode=%d]", timeNow, host, sessionUserName, sessionRoleID, method, rawURI, reqStr, respStatusCode)
	}
	return logMsg
}
