package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

func (e *Events) AuthorizeAndCreateTask(ctx context.Context, sessionToken string,
	resp *eventsproto.EventSubResponse) (string, string, error) {

	var (
		err                              error
		taskID, taskURI, sessionUserName string
	)

	// Athorize the request here
	authResp, err := e.Connector.Auth(ctx, sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("error while trying to authenticate session: status code: %v, status message: %v",
			authResp.StatusCode, authResp.StatusMessage)
		if err != nil {
			errMsg = errMsg + ": " + err.Error()
		}
		l.LogWithFields(ctx).Error(errMsg)
		resp.Body = generateResponse(ctx, authResp.Body)
		resp.StatusCode = authResp.StatusCode
		return sessionUserName, taskID, fmt.Errorf(errMsg)
	}

	sessionUserName, err = e.Connector.GetSessionUserName(ctx, sessionToken)
	if err != nil {
		errorMessage := "error while trying to get the session username: " + err.Error()
		resp.Body = generateResponse(ctx, common.GeneralError(http.StatusUnauthorized,
			response.NoValidSession, errorMessage, nil, nil))
		resp.StatusCode = http.StatusUnauthorized
		l.LogWithFields(ctx).Error(errorMessage)
		return sessionUserName, taskID, err
	}
	// Create the task and get the taskID
	// Contact Task Service using RPC and get the taskID
	taskURI, err = e.Connector.CreateTask(ctx, sessionUserName)
	if err != nil {
		// print err here as we are unbale to contact svc-task service
		errorMessage := "error while trying to create the task: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		resp.Body, _ = json.Marshal(common.GeneralError(http.StatusInternalServerError,
			response.InternalError, errorMessage, nil, nil).Body)
		l.LogWithFields(ctx).Error(errorMessage)
		return sessionUserName, taskID, fmt.Errorf(errorMessage)
	}

	taskID = strings.TrimPrefix(taskURI, "/redfish/v1/TaskService/Tasks/")
	resp.StatusCode = http.StatusAccepted
	resp.Header = map[string]string{
		"Location": "/taskmon/" + taskID,
	}
	resp.StatusMessage = response.TaskStarted
	generateTaskResponse(ctx, taskID, taskURI, resp)
	return sessionUserName, taskID, nil
}
