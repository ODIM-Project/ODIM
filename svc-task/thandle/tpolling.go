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

// Package thandle ...
package thandle

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"syscall"
	"time"

	restClient "github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/ODIM-Project/ODIM/svc-task/tcommon"
	"github.com/ODIM-Project/ODIM/svc-task/tmodel"
	"github.com/google/uuid"
)

/*
MonitorPluginTasks will trigger the polling plugin in a configured interval
*/
func (ts *TasksRPC) MonitorPluginTasks() {
	duration := time.Duration(config.Data.PluginTasksConf.MonitorPluginTasksFrequencyInMins)
	ticker := time.NewTicker(duration * time.Minute)
	for range ticker.C {
		ts.PollPlugin(GetContextForPolling())
	}
}

/*
PollPlugin will get the active plugin tasks from DB.
For each plugin tasks, it will poll the plugin instance with IP.
If IP is not not accessible for multiple retry, that instance will be marked
as unavailable and will fail the tasks those has been handled by the plugin
*/
func (ts *TasksRPC) PollPlugin(ctx context.Context) {
	l.LogWithFields(ctx).Info("Started polling plugin to monitor the plugin tasks")
	unAvailableInstances := make(map[string]struct{})
	pluginTaskIDs, err := tmodel.GetAllActivePluginTaskIDs(ctx)
	if err != nil {
		l.Log.Error(err)
	}

	plugins, err := tmodel.GetAllPlugins(ctx)
	pluginsMap := make(map[string]restClient.Plugin, len(plugins))

	for _, plugin := range plugins {
		pluginsMap[plugin.IP] = plugin
	}

	for _, taskID := range pluginTaskIDs {
		task, err := tmodel.GetPluginTaskInfo(taskID)
		if err != nil {
			continue
		}

		if _, ok := unAvailableInstances[task.IP]; ok {
			updateFailedPluginTasks(ctx, ts, taskID, task)
		}

		plugin := pluginsMap[task.PluginServerName]
		retry := 0
		isIPUnavailable := true
		for retry < 3 {
			_, err := tmodel.GetTaskMonResponse(ctx, plugin, task)
			if err == nil {
				isIPUnavailable = false
				break
			} else if isPluginConnectionError(err) {
				l.LogWithFields(ctx).Errorf("Plugin instance with ip %s is not accessible : error : %s. retrying...",
					task.IP, err.Error())
				retry++
			}
		}

		if isIPUnavailable {
			l.LogWithFields(ctx).Warnf("Plugin instance with ip %s is marked as unavailable.", task.IP)
			unAvailableInstances[task.IP] = struct{}{}
			updateFailedPluginTasks(ctx, ts, taskID, task)
		}
	}
	l.LogWithFields(ctx).Infof("Completed polling plugin to monitor the plugin tasks. Monitored %d plugin tasks",
		len(pluginTaskIDs))
}

/*
updateFailedPluginTasks will update the odim task corresponding to the plugin task as failed
and it will remove the plugin task from the active plugin tasks set in DB
*/
func updateFailedPluginTasks(ctx context.Context, ts *TasksRPC, pluginTaskID string, task *common.PluginTask) {
	statusCode := http.StatusServiceUnavailable
	message := errors.InternalError
	resp := tcommon.GetTaskResponse(int32(statusCode), message)
	body, _ := json.Marshal(resp.Body)

	payLoad := &taskproto.Payload{
		StatusCode:   int32(statusCode),
		ResponseBody: body,
	}

	l.LogWithFields(ctx).Warnf("Updating task %s with task state %s as plugin which handles that task is not accessible ",
		task.OdimTaskID, common.Cancelled)
	ts.updateTaskUtil(ctx, task.OdimTaskID, common.Cancelled,
		common.Critical, 100, payLoad, time.Now())
	tmodel.RemovePluginTaskID(ctx, pluginTaskID)
}

/*
isPluginConnectionError check the error returned by the plugin is a connection error or not
*/
func isPluginConnectionError(err error) bool {
	if netError, ok := err.(net.Error); ok && netError.Timeout() {
		return true
	}
	switch t := err.(type) {
	case *net.OpError:
		if t.Op == "dial" || t.Op == "read" {
			return true
		}

	case syscall.Errno:
		if t == syscall.ECONNREFUSED {
			return true
		}
	}
	return false
}

/*
GetContextForPolling create and returns a new context for polling the plugin
*/
func GetContextForPolling() context.Context {
	transactionID := uuid.New().String()
	actionID := "218"
	ctx := common.CreateContext(transactionID, actionID, common.PollPlugin, "0",
		common.PollPlugin, podName)
	return ctx
}
