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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
)

type Logging struct {
	GetUserDetails func(string) (string, string, error)
}

// AuditLog is used for generating audit logs in syslog format for each request
// this function logs an info for successful operation and error for failure operation
// properties logged are prival, time, host, username, roleid, request method, resource, requestbody, responsecode and message

func (l *Logging) AuditLog(ctx iris.Context, reqBody map[string]interface{}) {
	ctxt := ctx.Request().Context()
	logMsg, err := l.auditLogEntry(ctx, reqBody)
	if err != nil {
		Log.Error(err)
	}
	// Get response code
	respStatusCode := int32(ctx.GetStatusCode())
	operationStatus := getResponseStatus(respStatusCode)
	r := getProcessLogDetails(ctxt)
	logMsg = formatStructuredFields(Log.WithFields(r), logMsg)
	// 110 indicates info audit log for successful operation
	// 107 indicates error audit log for failed operation
	if operationStatus {
		successMsg := "<110>1 " + logMsg + " Operation successful"
		fmt.Println(successMsg)
	} else {
		failedMsg := "<107>1 " + logMsg + " Operation failed"
		fmt.Println(failedMsg)
	}
	return
}

// AuthLog is used for generating security logs in syslog format for each request
// this function logs an info for successful operation and warning for failure auth operation
// properties logged are prival, time, username, roleid and message
func AuthLog(ctx context.Context) *logrus.Entry {
	ctx = context.WithValue(ctx, "auth", true)
	ctx = context.WithValue(ctx, "statuscode", ctx.Value("statuscode").(int32))
	r := getProcessLogDetails(ctx)

	return Log.WithFields(r)
}

// auditLogEntry extracts the required info from context like session token, username, request URI
// and formats in syslog format for audit logs
func (l *Logging) auditLogEntry(ctx iris.Context, reqBody map[string]interface{}) (string, error) {
	var logMsg string
	// getting the request URI, host and method from context
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	sessionUserName, sessionRoleID, err := l.GetUserDetails(sessionToken)
	rawURI := ctx.Request().RequestURI
	host := ctx.Request().Host
	method := ctx.Request().Method
	respStatusCode := ctx.GetStatusCode()
	timeNow := time.Now().Format(time.RFC3339)
	reqStr := MaskRequestBody(reqBody)
	ctxt := ctx.Request().Context()
	thread := ctxt.Value("threadname")
	action := ctxt.Value("actionname")
	process := ctxt.Value("processname").(string)
	pid := os.Getpid()
	procid := process + fmt.Sprintf("_%d", pid)
	// formatting logs in syslog format
	if reqStr == "" {
		logMsg = fmt.Sprintf("%s %s %s %s %s [account@1 user=\"%s\" roleID=\"%s\"][request@1 method=\"%s\" resource=\"%s\"][response@1 responseCode=%d]", timeNow, host, thread, procid, action, sessionUserName, sessionRoleID, method, rawURI, respStatusCode)
	} else {
		logMsg = fmt.Sprintf("%s %s %s %s %s [account@1 user=\"%s\" roleID=\"%s\"][request@1 method=\"%s\" resource=\"%s\" requestBody= %s][response@1 responseCode=%d]", timeNow, host, thread, procid, action, sessionUserName, sessionRoleID, method, rawURI, reqStr, respStatusCode)
	}
	return logMsg, err
}

// MaskRequestBody function
// masking the request body, making password as null
func MaskRequestBody(reqBody map[string]interface{}) string {
	var jsonStr []byte
	var err error
	if len(reqBody) > 0 {
		reqBody["Password"] = "null"
		jsonStr, err = json.Marshal(reqBody)
		if err != nil {
			Log.Error("while marshalling request body", err.Error())
		}
	}
	reqStr := string(jsonStr)

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
