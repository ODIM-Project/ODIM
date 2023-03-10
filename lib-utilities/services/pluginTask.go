//(C) Copyright [2023] Hewlett Packard Enterprise Development LP
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

// Package services ...
package services

import (
	"context"
	"fmt"
	"strings"

	redis "github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

// SavePluginTaskInfo saves the ip of plugin instance that handle the task,
// task id of task which created in odim, and the taskmon URL returned
// from plugin in DB
func SavePluginTaskInfo(ctx context.Context, pluginIP, pluginServerName,
	odimTaskID, pluginTaskMonURL string) error {

	pluginTaskID := strings.TrimPrefix(pluginTaskMonURL, "/taskmon/")
	pluginTaskInfo := common.PluginTask{
		IP:               pluginIP,
		PluginServerName: pluginServerName,
		OdimTaskID:       odimTaskID,
		PluginTaskMonURL: pluginTaskMonURL,
	}

	err := createPluginTask(ctx, pluginTaskID, pluginTaskInfo)
	if err != nil {
		return fmt.Errorf("Error while saving plugin task info in DB: %s",
			err.Error())
	}
	return nil
}

// createPluginTask will insert plugin task info in DB
func createPluginTask(ctx context.Context, key string,
	value interface{}) *errors.Error {

	table := "PluginTask"
	connPool, err := redis.GetDBConnection(redis.InMemory)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to connecting"+
			" to DB: ", err.Error())
	}

	if err = connPool.Create(table, key, value); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to insert"+
			" plugin task: ", err.Error())
	}

	if err = connPool.AddMemberToSet(common.PluginTaskIndex, key); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to add "+
			" plugin task to set: ", err.Error())
	}

	return nil
}
