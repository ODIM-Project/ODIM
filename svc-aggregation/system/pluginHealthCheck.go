/* (C) Copyright [2020] Hewlett Packard Enterprise Development LP
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may
 * not use this file except in compliance with the License. You may obtain
 * a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

// Package system ...
package system

import (
	"net/http"
	"strings"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agcommon"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	log "github.com/sirupsen/logrus"
)

// SendStartUpData is for sending plugin start up data
func (e *ExternalInterface) SendStartUpData(startUpReq *aggregatorproto.SendStartUpDataRequest) response.RPC {
	plugin, err := agcommon.LookupPlugin(startUpReq.PluginAddr)
	if err != nil {
		log.Error("failed to find plugin with address " + startUpReq.PluginAddr + ": " + err.Error())
		return response.RPC{}
	}
	if err = PushPluginStartUpData(plugin, nil); err != nil {
		return response.RPC{}
	}
	return response.RPC{}
}

// PerformPluginHealthCheck is for checking the status of
// all the plugins continuously over a configured interval
func PerformPluginHealthCheck() {
	log.Info("plugins health check routine started")
	phc := agcommon.PluginHealthCheckInterface{
		DecryptPassword: common.DecryptWithPrivateKey,
	}
	for {
		phc.DupPluginConf()
		if pluginList, err := phc.GetAllPlugins(); err != nil {
			log.Error("failed to get list of all plugins:", err.Error())
		} else {
			for i := 0; i < len(pluginList); i++ {
				go phc.GetPluginStatus(pluginList[i])
			}
		}
		time.Sleep(time.Minute * time.Duration(phc.PluginConfig.PollingFrequencyInMins))
	}
}

// PushPluginStartUpData is for sending the plugin startup data
// when the plugin starts or when a server is added or deleted
func PushPluginStartUpData(plugin agmodel.Plugin, startUpData map[string]agmodel.PluginStartUpData) error {
	var serversData []agmodel.StartUpMap
	if startUpData == nil {
		phc := agcommon.PluginHealthCheckInterface{
			DecryptPassword: common.DecryptWithPrivateKey,
		}
		phc.DupPluginConf()
		managedServers := phc.GetPluginManagedServers(plugin)
		startUpMap := agmodel.StartUpMap{}
		startUpMap.PluginStartUpData = make(map[string]agmodel.PluginStartUpData, len(managedServers))
		for _, server := range managedServers {
			startUpMap.PluginStartUpData[server.ManagerAddress] = agmodel.PluginStartUpData{
				UserName:    server.UserName,
				Password:    server.Password,
				DeviceUUID:  server.DeviceUUID,
				Operation:   "add",
				RequestType: "full",
			}
		}
		serversData = append(serversData, startUpMap)
	} else {
		startUpMap := agmodel.StartUpMap{}
		startUpMap.PluginStartUpData = make(map[string]agmodel.PluginStartUpData, len(startUpData))
		for k, v := range startUpData {
			startUpMap.PluginStartUpData[k] = v
		}
		serversData = append(serversData, startUpMap)
	}

	var contactRequest agmodel.PluginContactRequest
	contactRequest.Plugin = plugin
	contactRequest.URL = "/ODIM/v1/Startup"
	contactRequest.HTTPMethodType = http.MethodPost
	contactRequest.PostBody = serversData

	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		contactRequest.HTTPMethodType = http.MethodPost
		contactRequest.PostBody = map[string]interface{}{
			"Username": plugin.Username,
			"Password": string(plugin.Password),
		}
		contactRequest.URL = "/ODIM/v1/Sessions"
		response, err := agcommon.ContactPlugin(contactRequest)
		if err != nil {
			log.Error("failed to get session token from " + plugin.ID + ": " + err.Error())
			return err
		}
		contactRequest.Token = response.Header.Get("X-Auth-Token")
		contactRequest.LoginCredential = nil
	} else {
		contactRequest.LoginCredential = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
	}

	response, err := agcommon.ContactPlugin(contactRequest)
	if err != nil {
		log.Error("failed to send startup data to "+plugin.ID+": "+err.Error()+": ", response)
		return err
	}
	log.Info("Successfully sent startup data to " + plugin.ID)
	return nil
}
