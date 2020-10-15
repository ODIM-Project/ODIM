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
	"log"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/fsnotify/fsnotify"
	uuid "github.com/satori/go.uuid"
)

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
		log.Println("Error While getting the status for plugin ", plugin.ID, err)
		return status
	}
	log.Println("Status of plugin", plugin.ID, status)
	return status
}

// GetStorageResources will get the resource details from the database for teh given odata id
func GetStorageResources(oid string) map[string]interface{} {
	resourceData := make(map[string]interface{})
	data, dbErr := agmodel.GetResourceDetails(oid)
	if dbErr != nil {
		log.Println("error while getting the system data", dbErr.Error())
		return resourceData
	}
	// unmarshall the resourceData
	err := json.Unmarshal([]byte(data), &resourceData)
	if err != nil {
		log.Println("Error while unmarshaling  the data", err)
		return resourceData
	}
	return resourceData
}

// AddConnectionMethods will add the connection method type and variant into DB
func AddConnectionMethods(connectionMethodConf []config.ConnectionMethodConf) error {
	connectionMethodsKeys, err := agmodel.GetAllKeysFromTable("ConnectionMethod")
	if err != nil {
		log.Printf("error getting connection methods : %v", err.Error())
		return err
	}
	var connectionMethodInfo = make(map[string]agmodel.ConnectionMethod)
	var connectionMehtodIDMap = make(map[string]string)
	// Get all existing connectionmethod info store it in above two map
	for i := 0; i < len(connectionMethodsKeys); i++ {
		connectionmethod, err := agmodel.GetConnectionMethod(connectionMethodsKeys[i])
		if err != nil {
			log.Printf("error getting connection method : %v", err)
			return err
		}
		connectionMethodInfo[connectionMethodsKeys[i]] = connectionmethod
		connectionMehtodIDMap[connectionmethod.ConnectionMethodType+":"+connectionmethod.ConnectionMethodVariant] = connectionMethodsKeys[i]
	}
	for i := 0; i < len(connectionMethodConf); i++ {
		if !SupportedConnectionMethodTypes[connectionMethodConf[i].ConnectionMethodType] {
			log.Printf("Connection method type %v is not supported.", connectionMethodConf[i].ConnectionMethodType)
			continue
		}
		if connectionMethodID, present := connectionMehtodIDMap[connectionMethodConf[i].ConnectionMethodType+":"+connectionMethodConf[i].ConnectionMethodVariant]; present {
			log.Printf("Connection Method Info with Connection Method Type %s and Connection Method Variant %s alredy present in ODIM", connectionMethodConf[i].ConnectionMethodType, connectionMethodConf[i].ConnectionMethodVariant)
			delete(connectionMethodInfo, connectionMethodConf[i].ConnectionMethodType+":"+connectionMethodConf[i].ConnectionMethodVariant)
			delete(connectionMehtodIDMap, connectionMethodID)
		} else {
			connectionMethodURI := "/redfish/v1/AggregationService/ConnectionMethods/" + uuid.NewV4().String()
			connectionMethod := agmodel.ConnectionMethod{
				ConnectionMethodType:    connectionMethodConf[i].ConnectionMethodType,
				ConnectionMethodVariant: connectionMethodConf[i].ConnectionMethodVariant,
				Links: agmodel.Links{
					AggregationSources: []agmodel.OdataID{},
				},
			}
			err := agmodel.AddConnectionMethod(connectionMethod, connectionMethodURI)
			if err != nil {
				log.Printf("error adding connection methods : %v", err.Error())
				return err
			}
		}
	}
	// loop thorugh the remaing connection method data from connectionMethodInfo map
	// delete the connection from ODIM if doesn't manage any aggreation else log the error
	for connectionMethodID, connectionMethodData := range connectionMethodInfo {
		if len(connectionMethodData.Links.AggregationSources) > 0 {
			log.Println("Connection Method ID ", connectionMethodID, " with Connection Method Type", connectionMethodData.ConnectionMethodType,
				" and Connection Method Variant", connectionMethodData.ConnectionMethodType, " managing ", string(len(connectionMethodData.Links.AggregationSources)), " Aggregation Sources it can't be removed")

		} else {
			log.Println("Removing Connection Method ID ", connectionMethodID, " with Connection Method Type", connectionMethodData.ConnectionMethodType,
				" and Connection Method Variant", connectionMethodData.ConnectionMethodType)
			err := agmodel.Delete("ConnectionMethod", connectionMethodID, common.OnDisk)
			if err != nil {
				log.Printf("error removing connection methods : %v", err.Error())
				return err
			}
		}
	}
	return nil
}

// TrackConfigFileChanges monitors the odim config changes using fsnotfiy
// When any files changes events recived it will reload the configuration and verify the existing events
func TrackConfigFileChanges() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	err = watcher.Add(ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			select {
			case fileEvent, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", fileEvent)
				if fileEvent.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", fileEvent.Name)
					// update the odim config
					if err := config.SetConfiguration(); err != nil {
						log.Printf("error while trying to set configuration: %v", err)
					}
					err = AddConnectionMethods(config.Data.ConnectionMethodConf)
					if err != nil {
						log.Printf("error while trying to Add connection methods: %v", err)
					}

				}
			case err, _ := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

}
