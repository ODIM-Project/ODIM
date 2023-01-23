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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agcommon"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/google/uuid"
)

var (
	// GetAllPluginfunc function pointer for the agcommon.GetAllPlugins
	GetAllPluginfunc = agcommon.GetAllPlugins
	// LookupPlugin function pointer for the agcommon.LookupPlugin
	LookupPlugin = agcommon.LookupPlugin
	// DecryptWithPrivateKey function pointer for the agcommon.DecryptWithPrivateKey
	DecryptWithPrivateKey = common.DecryptWithPrivateKey
	// GetPluginStatusRecord function pointer for the agcommon.GetPluginStatusRecord
	GetPluginStatusRecord = agcommon.GetPluginStatusRecord
	podName               = os.Getenv("POD_NAME")
)

const (
	PluginHealthCheckActionID = "216"

	PluginHealthCheckActionName = "PluginHealthCheck"
)

// SendStartUpData is for sending plugin start up data
func (e *ExternalInterface) SendStartUpData(ctx context.Context, startUpReq *aggregatorproto.SendStartUpDataRequest) response.RPC {
	resp := response.RPC{}
	plugin, err := LookupPlugin(ctx, startUpReq.PluginAddr)
	if err != nil {
		l.LogWithFields(ctx).Error("failed to find plugin with address " + startUpReq.PluginAddr + ": " + err.Error())
		return resp
	}

	l.LogWithFields(ctx).Infof("received plugin start up event from %s(%s)", plugin.ID, plugin.PluginType)

	// for plugins managing resources of non Compute type, at present
	// there is no usecase to share inventory, so subscribing to
	// EMB topic of the plugin should be enough
	if plugin.PluginType != "Compute" {
		phc := agcommon.PluginHealthCheckInterface{
			DecryptPassword: DecryptWithPrivateKey,
		}
		phc.DupPluginConf()

		active, topics := phc.GetPluginStatus(ctx, plugin)
		count, exist := GetPluginStatusRecord(plugin.ID)
		if !exist || (active && count != 0) {
			agcommon.SetPluginStatusRecord(plugin.ID, 0)
			PublishPluginStatusOKEvent(ctx, plugin.ID, topics)
			l.LogWithFields(ctx).Infof("subscribing to %s message bus topics of plugin %s", topics, plugin.ID)
		} else {
			agcommon.SetPluginStatusRecord(plugin.ID, count+1)
		}
		return resp
	}

	SendPluginStartUpData(ctx, startUpReq.OriginURI, plugin)
	return resp
}

// PerformPluginHealthCheck is for checking the status of
// all the plugins continuously over a configured interval
func PerformPluginHealthCheck() {
	transactionID := uuid.New()
	ctx := agcommon.CreateContext(transactionID.String(), PluginHealthCheckActionID, PluginHealthCheckActionName, "1", common.AggregationService, podName)
	l.LogWithFields(ctx).Info("plugins health check routine started")
	phc := agcommon.PluginHealthCheckInterface{
		DecryptPassword: DecryptWithPrivateKey,
	}
	for {
		phc.DupPluginConf()
		if pluginList, err := GetAllPluginfunc(ctx); err != nil {
			l.LogWithFields(ctx).Error("failed to get list of all plugins:", err.Error())
		} else {
			for _, plugin := range pluginList {
				threadID := 1
				ctxt := context.WithValue(ctx, common.ThreadName, common.CheckPluginStatus)
				ctxt = context.WithValue(ctxt, common.ThreadID, strconv.Itoa(threadID))
				go checkPluginStatus(ctxt, &phc, plugin)
				threadID++
			}
		}
		time.Sleep(time.Minute * time.Duration(phc.PluginConfig.PollingFrequencyInMins))
	}
}

func checkPluginStatus(ctx context.Context, phc *agcommon.PluginHealthCheckInterface, plugin agmodel.Plugin) {
	active, topics := phc.GetPluginStatus(ctx, plugin)
	if count, exist := GetPluginStatusRecord(plugin.ID); !exist {
		agcommon.SetPluginStatusRecord(plugin.ID, 0)
	} else {
		switch {
		case count != 0 && active:
			agcommon.SetPluginStatusRecord(plugin.ID, 0)
			if plugin.PluginType == "Compute" {
				if err := sharePluginInventory(ctx, plugin, true, plugin.IP); err != nil {
					l.LogWithFields(ctx).Error("failed to update server inventory of plugin " + plugin.ID + ": " + err.Error())
					agcommon.SetPluginStatusRecord(plugin.ID, count+1)
				}
			}
			PublishPluginStatusOKEvent(ctx, plugin.ID, topics)
			l.LogWithFields(ctx).Infof("subscribing to %s message bus topics of plugin %s", topics, plugin.ID)
		case !active:
			agcommon.SetPluginStatusRecord(plugin.ID, count+1)
		}
	}
}

// SendPluginStartUpData is for sending the plugin startup data
// when the plugin requests through an event
func SendPluginStartUpData(ctx context.Context, pluginIP string, plugin agmodel.Plugin) error {
	return sendFullPluginInventory(ctx, pluginIP, plugin)
}

// PushPluginStartUpData is for sending the plugin startup data
// when the plugin starts or when a server is added or deleted
func PushPluginStartUpData(ctx context.Context, plugin agmodel.Plugin, startUpData *agmodel.PluginStartUpData) error {
	if startUpData == nil {
		return fmt.Errorf("received empty startup data to share with %s", plugin.ID)
	}
	return sendPluginInventoryUpdate(ctx, plugin, *startUpData)
}

func sharePluginInventory(ctx context.Context, plugin agmodel.Plugin, resyncSubscription bool, serverName string) (ret error) {
	phc := agcommon.PluginHealthCheckInterface{
		DecryptPassword: DecryptWithPrivateKey,
	}
	phc.DupPluginConf()
	managedServers := phc.GetPluginManagedServers(plugin)
	managedServersCount := len(managedServers)
	if managedServersCount == 0 {
		l.LogWithFields(ctx).Info("plugin " + plugin.ID + " is not managing any server")
		return
	}
	pluginStartUpData := agmodel.PluginStartUpData{
		RequestType:           "full",
		ResyncEvtSubscription: resyncSubscription,
	}
	startIndex := 0
	for startIndex < managedServersCount {
		var batchedServersData []agmodel.Target
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
				l.LogWithFields(ctx).Error("failed to get event subscription details for " + server.ManagerAddress + ": " + err.Error())
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
		resp, err := sendPluginStartupRequest(ctx, plugin, pluginStartUpData, serverName)
		if err != nil {
			ret = fmt.Errorf("%v: %w", ret, err)
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ret = fmt.Errorf("%v: %w", ret, err)
			continue
		}
		var subsData map[string]string
		if err := json.Unmarshal(body, &subsData); err != nil {
			ret = fmt.Errorf("%v: %w", ret, err)
			continue
		}
		agcommon.UpdateDeviceSubscriptionDetails(ctx, subsData)
		batchedServersData = nil
	}
	return
}

func sendPluginInventoryUpdate(ctx context.Context, plugin agmodel.Plugin, startupData interface{}) error {
	if common.IsK8sDeployment() {
		addrList, err := common.GetServiceEndpointAddresses(plugin.IP)
		if err != nil {
			return err
		}
		var ret error
		for _, addr := range addrList {
			serverName := plugin.IP
			plugin.IP = addr
			if _, err := sendPluginStartupRequest(ctx, plugin, startupData, serverName); err != nil {
				ret = fmt.Errorf("%v: %w", ret, err)
			}
		}
		return ret
	}

	if _, err := sendPluginStartupRequest(ctx, plugin, startupData, plugin.IP); err != nil {
		return err
	}

	return nil
}

func sendPluginStartupRequest(ctx context.Context, plugin agmodel.Plugin, startupData interface{}, serverName string) (*http.Response, error) {
	contactRequest := agmodel.PluginContactRequest{}
	contactRequest.Plugin = plugin
	contactRequest.URL = "/ODIM/v1/Startup"
	contactRequest.HTTPMethodType = http.MethodPost
	contactRequest.PostBody = startupData
	response, err := agcommon.ContactPlugin(ctx, contactRequest, serverName)
	if err != nil || (response != nil && response.StatusCode != http.StatusOK) {
		l.LogWithFields(ctx).Errorf("failed to send startup data to %s(%s): %s: %+v", plugin.ID, plugin.IP, err, response)
		return nil, err
	}
	l.LogWithFields(ctx).Infof("Successfully sent startup data to %s(%s)", plugin.ID, plugin.IP)
	return response, nil
}

func sendFullPluginInventory(ctx context.Context, pluginIP string, plugin agmodel.Plugin) error {
	var reSubsEvent bool
	serverName := plugin.IP

	count, exist := GetPluginStatusRecord(plugin.ID)
	if !exist || count > 0 {
		agcommon.SetPluginStatusRecord(plugin.ID, 0)
	}
	if exist && count != 0 {
		reSubsEvent = true
	}
	if pluginIP != "" {
		plugin.IP = pluginIP
	}

	if err := sharePluginInventory(ctx, plugin, reSubsEvent, serverName); err != nil {
		l.LogWithFields(ctx).Errorf("failed to update server inventory of plugin %s(%s): %s", plugin.ID, plugin.IP, err.Error())
		agcommon.SetPluginStatusRecord(plugin.ID, count+1)
		return err
	}
	return nil
}
