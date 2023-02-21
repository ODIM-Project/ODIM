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

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/ODIM-Project/ODIM/svc-task/tcommon"
	"github.com/ODIM-Project/ODIM/svc-task/tmodel"
)

func (ts *TasksRPC) MonitorPluginTasks() {
	duration := time.Duration(config.Data.PluginTasksConf.MonitorPluginTasksFrequencyInMins)
	ticker := time.NewTicker(duration * time.Minute)
	for range ticker.C {
		ts.PollPlugin(context.TODO())
	}
}

func (ts *TasksRPC) PollPlugin(ctx context.Context) {
	l.LogWithFields(ctx).Info("Started polling plugin to monitor the plugin tasks")
	unAvailableInstances := make(map[string]struct{})
	pluginTaskIDs, err := tmodel.GetAllActivePluginTaskIDs(ctx)
	if err != nil {
		l.Log.Error(err)
	}

	plugins, err := tmodel.GetAllPlugins(ctx)
	pluginsMap := make(map[string]tmodel.Plugin, len(plugins))

	for _, plugin := range plugins {
		pluginsMap[plugin.IP] = plugin
	}

	for _, taskID := range pluginTaskIDs {
		task, err := tmodel.GetPluginTaskInfo(taskID)
		if err != nil {
			continue
		}

		if _, ok := unAvailableInstances[task.IP]; ok {
			l.LogWithFields(ctx).Infof("Updating task %s with task state %s as plugin which handles that task is not accessible ",
				task.OdimTaskID, common.Cancelled)
			payLoad := getPayloadForInternalError()
			ts.updateTaskUtil(ctx, task.OdimTaskID, common.Cancelled,
				common.Critical, 100, payLoad, time.Now())
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
				l.LogWithFields(ctx).Debugf("Plugin instance with ip %s is not accessible. retrying...", task.IP)
				retry++
			}
		}

		if isIPUnavailable {
			l.LogWithFields(ctx).Debugf("Plugin instance with ip %s is marked as unavailable.", task.IP)
			unAvailableInstances[task.IP] = struct{}{}
			payLoad := getPayloadForInternalError()
			l.LogWithFields(ctx).Infof("Updating task %s with task state %s as plugin which handles that task is not accessible ",
				task.OdimTaskID, common.Cancelled)
			ts.updateTaskUtil(ctx, task.OdimTaskID, common.Cancelled,
				common.Critical, 100, payLoad, time.Now())
		}
	}
	l.LogWithFields(ctx).Infof("Completed polling plugin to monitor the plugin tasks. Monitored %d plugin tasks",
		len(pluginTaskIDs))
}

func getPayloadForInternalError() *taskproto.Payload {
	statusCode := http.StatusInternalServerError
	message := errors.InternalError
	resp := tcommon.GetTaskResponse(statusCode, message)
	body, _ := json.Marshal(resp.Body)

	payLoad := &taskproto.Payload{
		StatusCode:   int32(statusCode),
		ResponseBody: body,
	}
	return payLoad
}

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
