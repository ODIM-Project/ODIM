package tcommon

import (
	"context"
	"encoding/json"
	"net/http"

	eventproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	managerproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

var (
	// ConfigFilePath holds the value of odim config file path
	ConfigFilePath string
)

const (
	// IterationCount is a value that needs to be added in context
	// to track the number of threads created
	IterationCount = "IterationCount"
)

// TaskStatusMap holds the task state against the message ID from task event
var TaskStatusMap = map[string]dmtf.TaskState{
	"TaskStarted":          dmtf.TaskStateStarting,
	"TaskProgressChanged":  dmtf.TaskStateRunning,
	"TaskPaused":           dmtf.TaskStateSuspended,
	"TaskAborted":          dmtf.TaskStateInterrupted,
	"TaskCompletedOK":      dmtf.TaskStateCompleted,
	"TaskRemoved":          dmtf.TaskStateKilled,
	"TaskCompletedWarning": dmtf.TaskStateException,
	"TaskCancelled":        dmtf.TaskStateCancelled,
}

// GetStatusCode return status code for the response based on task state and task status
func GetStatusCode(taskState dmtf.TaskState, taskStatus string) int {
	if taskState == dmtf.TaskStateCompleted && taskStatus == common.OK {
		return http.StatusOK
	} else if taskState == dmtf.TaskStateCancelled && taskStatus == common.Critical {
		return http.StatusInternalServerError
	}
	return http.StatusAccepted
}

// GetTaskResponse return status task response using status code and message
func GetTaskResponse(statusCode int32, message string) response.RPC {
	var resp response.RPC

	if statusCode == http.StatusNoContent {
		resp.StatusCode = statusCode
		resp.StatusMessage = response.ResourceRemoved
		return resp
	}

	err := json.Unmarshal([]byte(message), &resp.Body)
	if err == nil {
		resp.StatusCode = statusCode
		resp.StatusMessage = response.ExtendedInfo
		return resp
	}

	switch statusCode {
	case http.StatusOK:
		resp = common.GeneralError(statusCode, response.Success, message, nil, nil)
	case http.StatusInternalServerError:
		resp = common.GeneralError(statusCode, response.InternalError, message, nil, nil)
	default:
		resp.StatusCode = statusCode
		resp.StatusMessage = response.ExtendedInfo
		resp.Body = response.CommonError{
			Error: response.ErrorClass{
				Code:    response.ExtendedInfo,
				Message: message,
			},
		}
	}
	return resp
}

// TrackConfigFileChanges monitors the config changes using fsnotfiy
func TrackConfigFileChanges(errChan chan error) {
	eventChan := make(chan interface{})
	format := config.Data.LogFormat
	go common.TrackConfigFileChanges(ConfigFilePath, eventChan, errChan)
	for {
		select {
		case info := <-eventChan:
			l.Log.Info(info) // new data arrives through eventChan channel
			if l.Log.Level != config.Data.LogLevel {
				l.Log.Info("Log level is updated, new log level is ", config.Data.LogLevel)
				l.Log.Logger.SetLevel(config.Data.LogLevel)
			}
			if format != config.Data.LogFormat {
				l.SetFormatter(config.Data.LogFormat)
				format = config.Data.LogFormat
				l.Log.Info("Log format is updated, new log format is ", config.Data.LogFormat)
			}
		case err := <-errChan:
			l.Log.Error(err)
		}
	}
}

// UpdateSubscriptionLocation do the RPC call to events service to update subscription location
// once a create event subscription is completed
func UpdateSubscriptionLocation(ctx context.Context, location, host string) {
	conn, err := services.ODIMService.Client(services.Events)
	if err != nil {
		l.LogWithFields(ctx).Error("Error while Event ", err.Error())
		return
	}
	defer conn.Close()
	event := eventproto.NewEventsClient(conn)
	isUpdated, err := event.UpdateSubscriptionLocationRPC(ctx, &eventproto.UpdateSubscriptionLocation{
		Location: location,
		Host:     host,
	})
	if err != nil {
		l.LogWithFields(ctx).Info("Error while updating subscription location", err)
		return
	}
	l.LogWithFields(ctx).Debug("Location update status ", isUpdated)
}

// UpdateAccount do the RPC call to events service to update account details
func UpdateRemoteAccount(ctx context.Context, location, host string) {
	conn, err := services.ODIMService.Client(services.Managers)
	if err != nil {
		l.LogWithFields(ctx).Error("Error while Event ", err.Error())
		return
	}
	defer conn.Close()
	manager := managerproto.NewManagersClient(conn)
	isUpdated, err := manager.UpdateRemoteAccountPassword(ctx, &managerproto.ManagerRequest{
		ResourceID: location,
		ManagerID:  host,
	})
	if err != nil {
		l.LogWithFields(ctx).Info("Error while updating subscription location", err)
		return
	}
	l.LogWithFields(ctx).Debug("Location update status ", isUpdated)
}
