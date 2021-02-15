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

package smodel

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
)

// Fabric is the model to collect fabric plugin id from DB
type Fabric struct {
	PluginID string
}

// GetFabricManagers fetches all the fabrics details from DB
func GetFabricManagers() ([]Plugin, error) {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	keys, err := conn.GetAllDetails("Fabric")
	if err != nil {
		return nil, err
	}
	var managers []Plugin
	for _, key := range keys {
		var fabric Fabric
		fabricData, err := conn.Read("Fabric", key)
		if err != nil {
			log.Warn("while trying to read DB contents for " + key + " in Fabric table, got " + err.Error())
			continue
		}
		if errs := json.Unmarshal([]byte(fabricData), &fabric); errs != nil {
			log.Warn("while trying to unmarshal DB contents for " + key + " in Fabric table, got " + err.Error())
			continue
		}
		manager, err := GetPluginData(fabric.PluginID)
		if err != nil {
			log.Warn("while trying to collect DB contents for " + fabric.PluginID + " in Plugin table, got " + err.Error())
			continue
		}
		managers = append(managers, manager)

	}
	return managers, nil
}
