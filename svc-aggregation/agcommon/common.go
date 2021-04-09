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
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	uuid "github.com/satori/go.uuid"
	"strconv"
)

// DBInterface hold interface for db functions
type DBInterface struct {
	GetAllKeysFromTableInterface func(string) ([]string, error)
	GetConnectionMethodInterface func(string) (agmodel.ConnectionMethod, *errors.Error)
	AddConnectionMethodInterface func(agmodel.ConnectionMethod, string) *errors.Error
	DeleteInterface              func(string, string, common.DbType) *errors.Error
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

// GetPluginStatus checks the status of given plugin in configured interval
func GetPluginStatus(plugin agmodel.Plugin) bool {
	var pluginStatus = common.PluginStatus{
		Method: http.MethodGet,
		RequestBody: common.StatusRequest{
			Comment: "",
		},
		ResponseWaitTime:        config.Data.PluginStatusPolling.ResponseTimeoutInSecs,
		Count:                   config.Data.PluginStatusPolling.MaxRetryAttempt,
		RetryInterval:           config.Data.PluginStatusPolling.RetryIntervalInMins,
		PluginIP:                plugin.IP,
		PluginPort:              plugin.Port,
		PluginUsername:          plugin.Username,
		PluginUserPassword:      string(plugin.Password),
		PluginPrefferedAuthType: plugin.PreferredAuthType,
		CACertificate:           &config.Data.KeyCertConf.RootCACertificate,
	}
	status, _, _, err := pluginStatus.CheckStatus()
	if err != nil && !status {
		log.Error("Unable to get the plugin status for " + plugin.ID + " error : " + err.Error())
		return status
	}
	log.Info("Status of plugin " + plugin.ID + ": " + strconv.FormatBool(status))
	return status
}

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
