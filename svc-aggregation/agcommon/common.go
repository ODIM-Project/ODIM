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
	for i := 0; i < len(connectionMethodConf); i++ {
		var present bool
		if !SupportedConnectionMethodTypes[connectionMethodConf[i].ConnectionMethodType] {
			log.Printf("Connection method type %v is not supported.", connectionMethodConf[i].ConnectionMethodType)
			continue
		}
		for j := 0; j < len(connectionMethodsKeys); j++ {
			connectionmethod, err := agmodel.GetConnectionMethod(connectionMethodsKeys[j])
			if err != nil {
				log.Printf("error getting connection method : %v", err)
				return err
			}
			if connectionmethod.ConnectionMethodVariant == connectionMethodConf[i].ConnectionMethodVariant {
				present = true
				break
			}
		}
		if !present {
			connectionMethodURI := "/redfish/v1/AggregationService/ConnectionMethods/" + uuid.NewV4().String()
			connectionMethod := agmodel.ConnectionMethod{
				ConnectionMethodType:    connectionMethodConf[i].ConnectionMethodType,
				ConnectionMethodVariant: connectionMethodConf[i].ConnectionMethodVariant,
			}
			err := agmodel.AddConnectionMethod(connectionMethod, connectionMethodURI)
			if err != nil {
				log.Printf("error adding connection methods : %v", err.Error())
				return err
			}
		}
	}
	return nil
}
