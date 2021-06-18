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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	phc.StatusRecord.InactiveTime = make(map[string]int)
	for {
		phc.DupPluginConf()
		if pluginList, err := agcommon.GetAllPlugins(); err != nil {
			log.Error("failed to get list of all plugins:", err.Error())
		} else {
			for _, plugin := range pluginList {
				go checkPluginStatus(&phc, plugin)
			}
		}
		time.Sleep(time.Minute * time.Duration(phc.PluginConfig.PollingFrequencyInMins))
	}
}

func checkPluginStatus(phc *agcommon.PluginHealthCheckInterface, plugin agmodel.Plugin) {
	active, topics := phc.GetPluginStatus(plugin)
	if count, exist := phc.GetPluginStatusRecord(plugin.ID); !exist {
		phc.SetPluginStatusRecord(plugin.ID, 0)
	} else {
		switch {
		case count != 0 && active:
			phc.SetPluginStatusRecord(plugin.ID, 0)
			if err := sharePluginInventory(plugin, true); err != nil {
				log.Error("failed to update server inventory of plugin " +
					plugin.ID + ": " + err.Error())
				phc.SetPluginStatusRecord(plugin.ID, count+1)
			}
			PublishPluginStatusOKEvent(plugin.ID, topics)
		case !active:
			phc.SetPluginStatusRecord(plugin.ID, count+1)
		}
	}
}

// PushPluginStartUpData is for sending the plugin startup data
// when the plugin starts or when a server is added or deleted
func PushPluginStartUpData(plugin agmodel.Plugin, startUpData *agmodel.PluginStartUpData) error {
	if startUpData != nil {
		_, err := sendPluginStartupRequest(plugin, *startUpData)
		return err
	}
	return sharePluginInventory(plugin, false)
}

func sharePluginInventory(plugin agmodel.Plugin, resyncSubscription bool) (ret error) {
	phc := agcommon.PluginHealthCheckInterface{
		DecryptPassword: common.DecryptWithPrivateKey,
	}
	phc.DupPluginConf()
	managedServers := phc.GetPluginManagedServers(plugin)
	managedServersCount := len(managedServers)
	pluginStartUpData := agmodel.PluginStartUpData{
		RequestType:           "full",
		ResyncEvtSubscription: resyncSubscription,
	}
	batchedServersData := make([]agmodel.Target, 0)
	startIndex := 0
	for startIndex < managedServersCount {
		endIndex := startIndex + phc.PluginConfig.StartUpResouceBatchSize
		if endIndex > managedServersCount {
			endIndex = managedServersCount
		}
		batchedServersData = append(batchedServersData, managedServers[startIndex:endIndex]...)
		startIndex += phc.PluginConfig.StartUpResouceBatchSize
		pluginStartUpData.Devices = make(map[string]agmodel.DeviceData, phc.PluginConfig.StartUpResouceBatchSize)
		for _, server := range batchedServersData {
			evtSubsInfo := &agmodel.EventSubscriptionInfo{}
			subsID, evtTypes, err := agcommon.GetDeviceSubscriptionDetails(server.ManagerAddress)
			if err != nil {
				log.Error("failed to get event subscription details for " + server.ManagerAddress + ": " + err.Error())
			} else {
				evtSubsInfo.Location = subsID
				evtSubsInfo.EventTypes = append(evtSubsInfo.EventTypes, evtTypes...)
			}
			pluginStartUpData.Devices[server.DeviceUUID] = agmodel.DeviceData{
				Address:               server.ManagerAddress,
				UserName:              server.UserName,
				Password:              server.Password,
				Operation:             "add",
				EventSubscriptionInfo: evtSubsInfo,
			}
		}
		resp, err := sendPluginStartupRequest(plugin, pluginStartUpData)
		if err != nil {
			ret = fmt.Errorf("%w: %w", ret, err)
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ret = fmt.Errorf("%w: %w", ret, err)
			continue
		}
		var subsData map[string]string
		if err := json.Unmarshal(body, &subsData); err != nil {
			ret = fmt.Errorf("%w: %w", ret, err)
			continue
		}
		agcommon.UpdateDeviceSubscriptionDetails(subsData)
	}
	return
}

func sendPluginStartupRequest(plugin agmodel.Plugin, startupData interface{}) (*http.Response, error) {
	contactRequest := agmodel.PluginContactRequest{}
	contactRequest.Plugin = plugin
	contactRequest.URL = "/ODIM/v1/Startup"
	contactRequest.HTTPMethodType = http.MethodPost
	contactRequest.PostBody = startupData
	response, err := agcommon.ContactPlugin(contactRequest)
	if err != nil {
		log.Error("failed to send startup data to " + plugin.ID + ": " + err.Error())
		return nil, err
	}
	log.Info("Successfully sent startup data to " + plugin.ID)
	return response, nil
}
