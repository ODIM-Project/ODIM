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

package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	licenseproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/licenses"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-licenses/licenses"
)

// Licenses struct helps to register service
type Licenses struct {
	connector *licenses.ExternalInterface
}

// GetLicense intializes all the required connection
func GetLicense() *Licenses {
	return &Licenses{
		connector: licenses.GetExternalInterface(),
	}
}

func generateResponse(ctx context.Context, input interface{}) []byte {
	bytes, err := json.Marshal(input)
	if err != nil {
		l.LogWithFields(ctx).Warn("Unable to unmarshall response object from util-libs " + err.Error())
	}
	return bytes
}

func fillProtoResponse(ctx context.Context, resp *licenseproto.GetLicenseResponse, data response.RPC) {
	resp.StatusCode = data.StatusCode
	resp.StatusMessage = data.StatusMessage
	resp.Body = generateResponse(ctx, data.Body)
	resp.Header = data.Header

}

func generateRPCResponse(rpcResp response.RPC, licenseResp *licenseproto.GetLicenseResponse) {
	bytes, _ := json.Marshal(rpcResp.Body)
	*licenseResp = licenseproto.GetLicenseResponse{
		StatusCode:    rpcResp.StatusCode,
		StatusMessage: rpcResp.StatusMessage,
		Header:        rpcResp.Header,
		Body:          bytes,
	}
}

// CreateTaskAndResponse will create the task for corresponding request using
// the RPC call to task service and it will prepare custom task response to the user
// The function returns the ID of created task back.
func CreateTaskAndResponse(ctx context.Context, l *Licenses, sessionToken string, resp *licenseproto.GetLicenseResponse) (string, string, error) {
	sessionUserName, err := l.connector.External.GetSessionUserName(ctx, sessionToken)
	if err != nil {
		errMsg := "Unable to get session username: " + err.Error()
		fillProtoResponse(ctx, resp, common.GeneralError(http.StatusUnauthorized,
			response.NoValidSession, errMsg, nil, nil))
		return "", "", fmt.Errorf(errMsg)
	}

	// Task Service using RPC and get the taskID
	taskURI, err := l.connector.External.CreateTask(ctx, sessionUserName)
	if err != nil {
		errMsg := "Unable to create task: " + err.Error()
		fillProtoResponse(ctx, resp, common.GeneralError(http.StatusInternalServerError,
			response.InternalError, errMsg, nil, nil))
		return "", "", fmt.Errorf(errMsg)
	}
	taskID := strings.TrimPrefix(taskURI, "/redfish/v1/TaskService/Tasks/")
	// return 202 Accepted
	var rpcResp = response.RPC{
		StatusCode:    http.StatusAccepted,
		StatusMessage: response.TaskStarted,
		Header: map[string]string{
			"Location": "/taskmon/" + taskID,
		},
	}

	generateTaskRespone(taskID, taskURI, &rpcResp)
	fillProtoResponse(ctx, resp, rpcResp)
	return sessionUserName, taskID, nil
}

func generateTaskRespone(taskID, taskURI string, rpcResp *response.RPC) {
	commonResponse := response.Response{
		OdataType:    common.TaskType,
		ID:           taskID,
		Name:         "Task " + taskID,
		OdataContext: "/redfish/v1/$metadata#Task.Task",
		OdataID:      taskURI,
	}
	commonResponse.MessageArgs = []string{taskID}
	commonResponse.CreateGenericResponse(rpcResp.StatusMessage)
	rpcResp.Body = commonResponse
}
