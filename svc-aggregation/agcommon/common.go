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

package agcommon

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// DBInterface hold interface for db functions
type DBInterface struct {
	GetAllKeysFromTableInterface func(string) ([]string, error)
	GetConnectionMethodInterface func(string) (agmodel.ConnectionMethod, *errors.Error)
	AddConnectionMethodInterface func(agmodel.ConnectionMethod, string) *errors.Error
	DeleteInterface              func(string, string, common.DbType) *errors.Error
}

//PluginHealthCheckInterface holds the methods required for plugin healthcheck
type PluginHealthCheckInterface struct {
	DecryptPassword func([]byte) ([]byte, error)
	PluginConfig    config.PluginStatusPolling
	RootCA          []byte
}

// SupportedConnectionMethodTypes is for validating the connection method type
var SupportedConnectionMethodTypes = map[string]bool{
	"Redfish": true,
	"SNMP":    false,
	"OEM":     false,
	"NETCONF": false,
	"IPMI15":  false,
	"IPMI20":  false,
}

// ConfigFilePath holds the value of odim config file path
var ConfigFilePath string

// GetStorageResources will get the resource details from the database for teh given odata id
func GetStorageResources(oid string) map[string]interface{} {
	resourceData := make(map[string]interface{})
	data, dbErr := agmodel.GetResourceDetails(oid)
	if dbErr != nil {
		log.Error("Unable to get system data : " + dbErr.Error())
		return resourceData
	}
	// unmarshall the resourceData
	err := json.Unmarshal([]byte(data), &resourceData)
	if err != nil {
		log.Error("Unable to unmarshall  the data: " + err.Error())
		return resourceData
	}
	return resourceData
}

// AddConnectionMethods will add the connection method type and variant into DB
func (e *DBInterface) AddConnectionMethods(connectionMethodConf []config.ConnectionMethodConf) error {
	connectionMethodsKeys, err := e.GetAllKeysFromTableInterface("ConnectionMethod")
	if err != nil {
		log.Error("Unable to get connection methods : " + err.Error())
		return err
	}
	var connectionMethodInfo = make(map[string]agmodel.ConnectionMethod)
	var connectionMehtodIDMap = make(map[string]string)
	// Get all existing connectionmethod info store it in above two map
	for i := 0; i < len(connectionMethodsKeys); i++ {
		connectionmethod, err := e.GetConnectionMethodInterface(connectionMethodsKeys[i])
		if err != nil {
			log.Error("Unable to get connection method : " + err.Error())
			return err
		}
		connectionMethodInfo[connectionMethodsKeys[i]] = connectionmethod
		connectionMehtodIDMap[connectionmethod.ConnectionMethodType+":"+connectionmethod.ConnectionMethodVariant] = connectionMethodsKeys[i]
	}
	for i := 0; i < len(connectionMethodConf); i++ {
		if !SupportedConnectionMethodTypes[connectionMethodConf[i].ConnectionMethodType] {
			log.Error("Connection method type " + connectionMethodConf[i].ConnectionMethodType + " is not supported")
			continue
		}
		if connectionMethodID, present := connectionMehtodIDMap[connectionMethodConf[i].ConnectionMethodType+":"+connectionMethodConf[i].ConnectionMethodVariant]; present {
			log.Error("Connection Method Info with Connection Method Type " +
				connectionMethodConf[i].ConnectionMethodType + " and Connection Method Variant " +
				connectionMethodConf[i].ConnectionMethodVariant + " already present in ODIM")
			delete(connectionMehtodIDMap,
				connectionMethodConf[i].ConnectionMethodType+":"+connectionMethodConf[i].ConnectionMethodVariant)
			delete(connectionMethodInfo, connectionMethodID)
		} else {
			connectionMethodURI := "/redfish/v1/AggregationService/ConnectionMethods/" + uuid.NewV4().String()
			connectionMethod := agmodel.ConnectionMethod{
				ConnectionMethodType:    connectionMethodConf[i].ConnectionMethodType,
				ConnectionMethodVariant: connectionMethodConf[i].ConnectionMethodVariant,
				Links: agmodel.Links{
					AggregationSources: []agmodel.OdataID{},
				},
			}
			err := e.AddConnectionMethodInterface(connectionMethod, connectionMethodURI)
			if err != nil {
				log.Error("Unable to add connection method : " + err.Error())
				return err
			}
			log.Info(
				"Connection method info with connection method type " + connectionMethodConf[i].ConnectionMethodType +
					" and connection method variant " + connectionMethodConf[i].ConnectionMethodVariant + " added to ODIM")
		}
	}
	// loop thorugh the remaining connection method data from connectionMethodInfo map
	// delete the connection from ODIM if doesn't manage any aggreation source else log the error
	for connectionMethodID, connectionMethodData := range connectionMethodInfo {
		if len(connectionMethodData.Links.AggregationSources) > 0 {
			log.Error("Connection Method ID: " + connectionMethodID + " with connection method type " +
				connectionMethodData.ConnectionMethodType + " and connection method variant " +
				connectionMethodData.ConnectionMethodVariant + " managing " +
				string(len(connectionMethodData.Links.AggregationSources)) + " aggregation sources it can't be removed")

		} else {
			log.Info("Removing connection method id "+connectionMethodID+
				" with Connection Method Type"+connectionMethodData.ConnectionMethodType+
				" and Connection Method Variant", connectionMethodData.ConnectionMethodVariant)
			err := e.DeleteInterface("ConnectionMethod", connectionMethodID, common.OnDisk)
			if err != nil {
				log.Error("Unable to removing connection method : " + err.Error())
				return err
			}
		}
	}
	return nil
}

// TrackConfigFileChanges monitors the odim config changes using fsnotfiy
// Whenever  any config file changes and events  will be  and  reload the configuration and verify the existing connection methods
func TrackConfigFileChanges(dbInterface DBInterface) {
	eventChan := make(chan interface{})
	go common.TrackConfigFileChanges(ConfigFilePath, eventChan)
	for {
		log.Info(<-eventChan) // new data arrives through eventChan channel
		config.TLSConfMutex.RLock()
		err := dbInterface.AddConnectionMethods(config.Data.ConnectionMethodConf)
		if err != nil {
			log.Error("error while trying to Add connection methods:" + err.Error())
		}
		config.TLSConfMutex.RUnlock()
	}
}

// DupPluginConf is for duplicating the plugin status polling config using a lock
// at one place instead of acquiring a lock and reading the config params multiple times
func (phc *PluginHealthCheckInterface) DupPluginConf() {
	config.TLSConfMutex.RLock()
	defer config.TLSConfMutex.RUnlock()
	phc.PluginConfig.PollingFrequencyInMins = config.Data.PluginStatusPolling.PollingFrequencyInMins
	phc.PluginConfig.MaxRetryAttempt = config.Data.PluginStatusPolling.MaxRetryAttempt
	phc.PluginConfig.RetryIntervalInMins = config.Data.PluginStatusPolling.RetryIntervalInMins
	phc.PluginConfig.ResponseTimeoutInSecs = config.Data.PluginStatusPolling.ResponseTimeoutInSecs
	phc.PluginConfig.StartUpResouceBatchSize = config.Data.PluginStatusPolling.StartUpResouceBatchSize
	phc.RootCA = make([]byte, len(config.Data.KeyCertConf.RootCACertificate))
	copy(phc.RootCA, config.Data.KeyCertConf.RootCACertificate)
	return
}

// GetPluginStatus checks the status of given plugin
func GetPluginStatus(plugin agmodel.Plugin) bool {
	phc := &PluginHealthCheckInterface{}
	phc.DupPluginConf()
	return phc.GetPluginStatus(plugin)
}

// LookupPlugin is for fetching the plugin data
// using the plugin address for lookup
func LookupPlugin(addr string) (agmodel.Plugin, error) {
	phc := &PluginHealthCheckInterface{}
	phc.DupPluginConf()
	plugins, errs := phc.GetAllPlugins()
	if errs != nil {
		return agmodel.Plugin{}, errs
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Warn("splitting plugin address failed with " + err.Error())
		host = addr
		port = ""
	}

	var resolvedAddr string
	ip, err := net.LookupIP(host)
	switch {
	case err != nil:
		log.Warn("plugin address lookup failed with " + err.Error())
	case len(addr) < 1:
		log.Warn("plugin address lookup gave empty list")
	default:
		resolvedAddr = ip[0].String()
	}

	for _, plugin := range plugins {
		if (plugin.IP == host || plugin.IP == resolvedAddr) && (plugin.Port == port) {
			return plugin, nil
		}
	}
	return agmodel.Plugin{}, fmt.Errorf(addr + " address does not belong to any of the plugin")
}

// GetAllPlugins is for fetching all the plugins added andn stored in db.
func (phc *PluginHealthCheckInterface) GetAllPlugins() ([]agmodel.Plugin, *errors.Error) {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	keys, err := conn.GetAllDetails("Plugin")
	if err != nil {
		return nil, err
	}
	var plugins []agmodel.Plugin
	for _, key := range keys {
		plugin, err := agmodel.GetPluginData(key)
		if err != nil {
			log.Error("failed to get details of " + key + " plugin: " + err.Error())
			continue
		}
		plugins = append(plugins, plugin)
	}
	return plugins, nil
}

// GetPluginStatus is for checking the status of a plugin
func (phc *PluginHealthCheckInterface) GetPluginStatus(plugin agmodel.Plugin) bool {
	var pluginStatus = common.PluginStatus{
		Method: http.MethodGet,
		RequestBody: common.StatusRequest{
			Comment: "",
		},
		ResponseWaitTime:        phc.PluginConfig.ResponseTimeoutInSecs,
		Count:                   phc.PluginConfig.MaxRetryAttempt,
		RetryInterval:           phc.PluginConfig.RetryIntervalInMins,
		PluginIP:                plugin.IP,
		PluginPort:              plugin.Port,
		PluginUsername:          plugin.Username,
		PluginUserPassword:      string(plugin.Password),
		PluginPrefferedAuthType: plugin.PreferredAuthType,
		CACertificate:           &phc.RootCA,
	}
	status, _, _, err := pluginStatus.CheckStatus()
	if err != nil {
		log.Error("failed to get the status of plugin " + plugin.ID + err.Error())
		return status
	}
	log.Info("Status of plugin " + plugin.ID + " is " + strconv.FormatBool(status))
	return status
}

// GetPluginManagedServers is for fetching the list of servers managed by a plugin
func (phc *PluginHealthCheckInterface) GetPluginManagedServers(plugin agmodel.Plugin) []agmodel.ServerInfo {
	if status := phc.GetPluginStatus(plugin); !status {
		log.Error(plugin.ID + " status check failed")
		return []agmodel.ServerInfo{}
	}
	serversList, err := phc.getAllServers(plugin.ID)
	if err != nil {
		log.Error("failed to get list of servers managed by " + plugin.ID + err.Error())
	}
	return serversList
}

// getAllServers is for fetching the list of all servers added.
func (phc *PluginHealthCheckInterface) getAllServers(pluginID string) ([]agmodel.ServerInfo, error) {
	var matchedServers []agmodel.ServerInfo
	allServers, err := getAllSystems()
	if err != nil {
		log.Error("failed to get the list of all managed servers " + err.Error())
		return matchedServers, err
	}
	for i := 0; i < len(allServers); i++ {
		var server agmodel.ServerInfo
		singleServer, err := getSingleSystem(allServers[i])
		if err != nil {
			log.Error("failed to get info of " + allServers[i] + " system: " + err.Error())
			continue
		}
		json.Unmarshal([]byte(singleServer), &server)
		if server.PluginID == pluginID {
			decryptedPasswordByte, err := phc.DecryptPassword(server.Password)
			if err != nil {
				log.Error("failed to decrypt device password of the host: " + server.ManagerAddress + ":" + err.Error())
				continue
			}
			server.Password = decryptedPasswordByte
			matchedServers = append(matchedServers, server)
		}
	}
	return matchedServers, err
}

// getAllSystems is for fetching all the keys from the System table
func getAllSystems() ([]string, error) {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	keysArray, err := conn.GetAllDetails("System")
	if err != nil {
		return nil, errors.PackError(errors.UndefinedErrorType, err.Error())
	}
	return keysArray, nil
}

// getSingleSystem is for fetching the details of a server from System table
func getSingleSystem(id string) (string, error) {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return "", errors.PackError(errors.UndefinedErrorType, err)
	}

	data, rerr := conn.Read("System", id)
	if rerr != nil {
		return "", errors.PackError(rerr.ErrNo(), rerr.Error())
	}
	return data, nil
}

// ContactPlugin is for sending requests to a plugin.
func ContactPlugin(req agmodel.PluginContactRequest) (*http.Response, error) {
	reqURL := "https://" + req.Plugin.IP + ":" + req.Plugin.Port + req.URL
	return pmbhandle.ContactPlugin(reqURL, req.HTTPMethodType, req.Token, "", req.PostBody, req.LoginCredential)
}
