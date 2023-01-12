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

	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
)

type Logging struct {
	GetUserDetails func(string) (string, string, error)
}

// AuthLog is used for generating security logs in syslog format for each request
// this function logs an info for successful operation and warning for failure auth operation
// properties logged are prival, time, username, roleid and message
func AuthLog(ctx context.Context) *logrus.Entry {
	fields := getProcessLogDetails(ctx)
	if val, ok := ctx.Value("sessiontoken").(string); ok {
		fields["sessiontoken"] = val
	}
	if val, ok := ctx.Value("sessionuserid").(string); ok {
		fields["sessionuserid"] = val
	}
	if val, ok := ctx.Value("sessionroleid").(string); ok {
		fields["sessionroleid"] = val
	}
	if val, ok := ctx.Value("statuscode").(int32); ok {
		fields["statuscode"] = val
	}
	fields["auth"] = true

	return Log.WithFields(fields)
}

// formatAuditStructFields is used to format audit log message with required values
func formatAuditStructFields(entry *logrus.Entry, msg string, priorityNo int8) string {
	var reqStr, logMsg, host, sessionUserName, sessionRoleID, method, rawURI string
	var respStatusCode int32
	if val, ok := entry.Data["reqstr"].(string); ok {
		reqStr = val
	}
	if val, ok := entry.Data["sessionusername"].(string); ok {
		sessionUserName = val
	}
	if val, ok := entry.Data["sessionroleid"].(string); ok {
		sessionRoleID = val
	}
	if val, ok := entry.Data["rawuri"].(string); ok {
		rawURI = val
	}
	if val, ok := entry.Data["host"].(string); ok {
		host = val
	}
	if val, ok := entry.Data["method"].(string); ok {
		method = val
	}
	if val, ok := entry.Data["statuscode"].(int32); ok {
		respStatusCode = val
	}

	// formatting logs in syslog format
	if reqStr == "" {
		logMsg = fmt.Sprintf("%s [account@1 user=\"%s\" roleID=\"%s\"][request@1 method=\"%s\" resource=\"%s\"][response@1 responseCode=%d]", host, sessionUserName, sessionRoleID, method, rawURI, respStatusCode)
	} else {
		logMsg = fmt.Sprintf("%s [account@1 user=\"%s\" roleID=\"%s\"][request@1 method=\"%s\" resource=\"%s\" requestBody= %s][response@1 responseCode=%d]", host, sessionUserName, sessionRoleID, method, rawURI, reqStr, respStatusCode)
	}

	if priorityNo == 110 {
		msg = fmt.Sprintf("%s %s %s", msg, logMsg, "Operation successful")
		return msg
	}
	msg = fmt.Sprintf("%s %s %s", msg, logMsg, "Operation failed")
	return msg
}

// AuditLog is used for generating audit logs in syslog format for each request
// this function logs an info for successful operation and error for failure operation
// properties logged are prival, time, host, username, roleid, request method, resource, requestbody, responsecode and message
func AuditLog(l *Logging, ctx iris.Context, reqBody map[string]interface{}) *logrus.Entry {
	ctxt := ctx.Request().Context()
	fields := getProcessLogDetails(ctxt)
	if val, ok := ctxt.Value("sessiontoken").(string); ok {
		fields["sessiontoken"] = val
	}
	if val, ok := ctxt.Value("sessionuserid").(string); ok {
		fields["sessionuserid"] = val
	}
	if val, ok := ctxt.Value("sessionroleid").(string); ok {
		fields["sessionroleid"] = val
	}
	fields["statuscode"] = int32(ctx.GetStatusCode())
	fields["audit"] = true
	fields, err := l.auditLogEntry(ctx, reqBody, fields)
	if err != nil {
		Log.Error(err)
	}

	return Log.WithFields(fields)
}

// auditLogEntry extracts the required info from context like session token, username, request URI
// and formats in syslog format for audit logs
func (l *Logging) auditLogEntry(ctx iris.Context, reqBody, fields map[string]interface{}) (logrus.Fields, error) {
	// getting the request URI, host and method from context
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	sessionUserName, sessionRoleID, err := l.GetUserDetails(sessionToken)
	fields["sessiontoken"] = sessionToken
	fields["sessionusername"] = sessionUserName
	fields["sessionroleid"] = sessionRoleID
	fields["rawuri"] = ctx.Request().RequestURI
	fields["host"] = ctx.Request().Host
	fields["method"] = ctx.Request().Method
	fields["reqstr"] = MaskRequestBody(reqBody)

	return fields, err
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
