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

package system

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

//UpdateSystemState is used for updating ComputerSystem table
//and also the server search index, if required.
func (e *ExternalInterface) UpdateSystemState(updateReq *aggregatorproto.UpdateSystemStateRequest) error {

	key := fmt.Sprintf("%s/%s.%s", strings.TrimSuffix(updateReq.SystemURI, "/"), updateReq.SystemUUID, updateReq.SystemID)

	// Getting the device info
	target, err := agmodel.GetTarget(updateReq.SystemUUID)
	if err != nil {
		return err

	}
	decryptedPasswordByte, err := e.DecryptPassword(target.Password)
	if err != nil {
		return err
	}
	target.Password = decryptedPasswordByte

	// get the plugin information
	plugin, errs := agmodel.GetPluginData(target.PluginID)
	if errs != nil {
		return errs
	}

	var req getResourceRequest
	req.ContactClient = e.ContactClient
	req.GetPluginStatus = e.GetPluginStatus
	req.Plugin = plugin
	req.StatusPoll = true
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		req.HTTPMethodType = http.MethodPost
		req.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		req.OID = "/ODIM/v1/Sessions"
		_, token, _, err := contactPlugin(req, "error while getting the details "+req.OID+": ")
		if err != nil {
			return err
		}
		req.Token = token
	} else {
		req.LoginCredentials = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}

	req.DeviceUUID = updateReq.SystemUUID
	req.DeviceInfo = target
	req.OID = fmt.Sprintf("%s/%s", strings.TrimSuffix(updateReq.SystemURI, "/"), updateReq.SystemID)

	rawData, _, getResponse, err := contactPlugin(req, "error while trying to get system details ")
	if err != nil || getResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("error: while trying to get system details ")
	}

	//replacing the uuid while saving the data
	updatedResourceData := updateResourceDataWithUUID(string(rawData), req.DeviceUUID)
	var systemInfo map[string]interface{}

	if err := json.Unmarshal([]byte(updatedResourceData), &systemInfo); err != nil {
		return fmt.Errorf("error: failed to unmarshal fetched data with %v", err.Error())
	}

	// Updating the Computer System Data
	if err := agmodel.UpdateComputeSystem(key, systemInfo); err != nil {
		return fmt.Errorf("error: failed to update data with %v", err.Error())
	}
	if _, ok := systemInfo["UUID"]; !ok {
		return fmt.Errorf("error: failed to update the data computer system uuid not found")
	}
	computerSystemUUID := systemInfo["UUID"].(string)
	searchForm := createServerSearchIndex(systemInfo, key, updateReq.SystemUUID)
	if err := agmodel.UpdateIndex(searchForm, key, computerSystemUUID, target.ManagerAddress); err != nil {
		return fmt.Errorf("error: updating server index failed with err %v", err)
	}
	return nil
}
