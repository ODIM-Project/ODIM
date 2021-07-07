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

//Package scommon ...
package scommon

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
)

// Schema is used to define the allowed values for search/filter
type Schema struct {
	SearchKeys    []map[string]map[string]string `json:"searchKeys"`
	ConditionKeys []string                       `json:"conditionKeys"`
	QueryKeys     []string                       `json:"queryKeys"`
}

// SF holds the schema data for search/filter
var SF Schema

//PluginContactRequest  hold the request of contact plugin
type PluginContactRequest struct {
	Token           string
	OID             string
	DeviceInfo      interface{}
	BasicAuth       map[string]string
	ContactClient   func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	GetPluginStatus func(smodel.Plugin) bool
	Plugin          smodel.Plugin
	HTTPMethodType  string
}

//ResponseStatus holds the response of Contact Plugin
type ResponseStatus struct {
	StatusCode    int32
	StatusMessage string
}

//ResourceInfoRequest  hold the request of getting  Resource
type ResourceInfoRequest struct {
	URL             string
	UUID            string
	SystemID        string
	ContactClient   func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	DevicePassword  func([]byte) ([]byte, error)
	GetPluginStatus func(smodel.Plugin) bool
	ResourceName    string
}

// GetResourceInfoFromDevice will contact to the and gets the Particual resource info from device
// If saveRequired is set to true the newly collected data will be saved in the DB.
// Some specific cases may not require the data to be stored in DB,
// eg: Delete volume requires reset of the BMC to take its effect. Before a reset, volumes retrieval
// request can provide the deleted volume. We can avoid storing such a data with the use of saveRequired.
func GetResourceInfoFromDevice(req ResourceInfoRequest, saveRequired bool) (string, error) {
	target, gerr := smodel.GetTarget(req.UUID)
	if gerr != nil {
		return "", gerr
	}
	// Get the Plugin info
	plugin, gerr := smodel.GetPluginData(target.PluginID)
	if gerr != nil {
		return "", gerr
	}
	var contactRequest PluginContactRequest

	contactRequest.ContactClient = req.ContactClient
	contactRequest.Plugin = plugin
	contactRequest.GetPluginStatus = req.GetPluginStatus
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		contactRequest.HTTPMethodType = http.MethodPost
		contactRequest.DeviceInfo = map[string]interface{}{
			"Username": plugin.Username,
			"Password": string(plugin.Password),
		}
		contactRequest.OID = "/ODIM/v1/Sessions"
		_, token, _, err := ContactPlugin(contactRequest, "error while getting the details "+contactRequest.OID+": ")
		if err != nil {

			return "", err
		}
		contactRequest.Token = token
	} else {
		contactRequest.BasicAuth = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	decryptedPasswordByte, err := req.DevicePassword(target.Password)
	if err != nil {
		// Frame the RPC response body and response Header below
		errorMessage := "error while trying to decrypt device password: " + err.Error()

		return "", fmt.Errorf(errorMessage)
	}
	contactRequest.DeviceInfo = map[string]interface{}{
		"ManagerAddress": target.ManagerAddress,
		"UserName":       target.UserName,
		"Password":       decryptedPasswordByte,
	}
	//replace the uuid:system id with the system to the @odata.id from request url
	contactRequest.OID = strings.Replace(req.URL, req.UUID+":"+req.SystemID, req.SystemID, -1)
	contactRequest.HTTPMethodType = http.MethodGet
	body, _, _, err := ContactPlugin(contactRequest, "error while getting the details "+contactRequest.OID+": ")
	if err != nil {
		return "", err
	}

	var resourceData map[string]interface{}
	err = json.Unmarshal(body, &resourceData)
	if err != nil {
		return "", err
	}
	var resourceName, oidKey string

	/* Checking if the retrieved data is allowed to save in DB,
	 * if not allowed we will not save in DB.
	 */
	//replacing the uuid while saving the data
	//to replace the id of system
	var updatedData = strings.Replace(string(body), "/redfish/v1/Systems/", "/redfish/v1/Systems/"+req.UUID+":", -1)
	updatedData = strings.Replace(updatedData, "/redfish/v1/systems/", "/redfish/v1/systems/"+req.UUID+":", -1)
	// to replace the id in managers
	updatedData = strings.Replace(updatedData, "/redfish/v1/Managers/", "/redfish/v1/Managers/"+req.UUID+":", -1)
	// to replace id in chassis
	updatedData = strings.Replace(updatedData, "/redfish/v1/Chassis/", "/redfish/v1/Chassis/"+req.UUID+":", -1)

	if saveRequired && checkRetrievalInfo(contactRequest.OID) {
		oidKey = keyFormation(contactRequest.OID, req.SystemID, req.UUID)
		var memberFlag bool
		if _, ok := resourceData["Members"]; ok {
			memberFlag = true
		}
		if req.ResourceName != "" {
			resourceName = req.ResourceName
		} else {
			// Get the Table name to save the data in db
			resourceName = getResourceName(contactRequest.OID, memberFlag)
		}
		// persist the response with table resourceName and key as system UUID + Oid Needs relook TODO
		err = smodel.GenericSave([]byte(updatedData), resourceName, oidKey)
		if err != nil {
			return "", err
		}
	}
	return updatedData, nil
}

// keyFormation is to form the key to insert in DB
func keyFormation(oid, systemID, DeviceUUID string) string {
	if oid[len(oid)-1:] == "/" {
		oid = oid[:len(oid)-1]
	}
	str := strings.Split(oid, "/")
	var key []string
	for i, id := range str {
		if id == systemID && (strings.EqualFold(str[i-1], "Systems") || strings.EqualFold(str[i-1], "Chassis")) {
			key = append(key, DeviceUUID+":"+id)
			continue
		}
		key = append(key, id)
	}
	return strings.Join(key, "/")
}

//getResourceName fetches the table name for storing the particualar resource
func getResourceName(oDataID string, memberFlag bool) string {
	str := strings.Split(oDataID, "/")
	if memberFlag {
		return str[len(str)-1] + "Collection"
	}
	if _, err := strconv.Atoi(str[len(str)-2]); err == nil {
		return str[len(str)-1]
	}
	return str[len(str)-2]
}

// ContactPlugin is commons which handles the request and response of Contact Plugin usage
func ContactPlugin(req PluginContactRequest, errorMessage string) ([]byte, string, ResponseStatus, error) {
	var resp ResponseStatus
	var response *http.Response
	var err error
	response, err = callPlugin(req)
	if err != nil {
		if req.GetPluginStatus(req.Plugin) {
			response, err = callPlugin(req)
		}
		if err != nil {
			errorMessage = errorMessage + err.Error()
			resp.StatusCode = http.StatusInternalServerError
			resp.StatusMessage = errors.InternalError
			log.Error(errorMessage)
			return nil, "", resp, fmt.Errorf(errorMessage)
		}
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errorMessage := "error while trying to read response body: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = errors.InternalError
		log.Error(errorMessage)
		return nil, "", resp, fmt.Errorf(errorMessage)
	}
	log.Info("Response" + string(body))
	log.Info("response.StatusCode" + string(response.StatusCode))
	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		resp.StatusCode = int32(response.StatusCode)
		log.Println(errorMessage)
		return body, "", resp, fmt.Errorf(errorMessage)
	}

	data := string(body)
	//replacing the resposne with north bound translation URL
	for key, value := range config.Data.URLTranslation.NorthBoundURL {
		data = strings.Replace(data, key, value, -1)
	}
	return []byte(data), response.Header.Get("X-Auth-Token"), resp, nil
}

func checkRetrievalInfo(oid string) bool {
	//skiping the Retrieval if parent oid contains links in other resource of config
	for _, resourceName := range config.Data.AddComputeSkipResources.SkipResourceListUnderOthers {
		if strings.Contains(oid, resourceName) {
			return false
		}
	}
	return true
}

// GetPluginStatus checks the status of given plugin in configured interval
func GetPluginStatus(plugin smodel.Plugin) bool {
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
		log.Error("Error While getting the status for plugin " + plugin.ID + ": " + err.Error())
		return status
	}
	log.Info("Status of plugin" + plugin.ID + strconv.FormatBool(status))
	return status
}

func callPlugin(req PluginContactRequest) (*http.Response, error) {
	var oid string
	for key, value := range config.Data.URLTranslation.SouthBoundURL {
		oid = strings.Replace(req.OID, key, value, -1)
	}
	var reqURL = "https://" + req.Plugin.IP + ":" + req.Plugin.Port + oid
	if strings.EqualFold(req.Plugin.PreferredAuthType, "BasicAuth") {
		return req.ContactClient(reqURL, req.HTTPMethodType, "", oid, req.DeviceInfo, req.BasicAuth)
	}
	return req.ContactClient(reqURL, req.HTTPMethodType, req.Token, oid, req.DeviceInfo, nil)
}

// TrackConfigFileChanges monitors the odim config changes using fsnotfiy
func TrackConfigFileChanges(configFilePath string) {
	eventChan := make(chan interface{})
	go common.TrackConfigFileChanges(configFilePath, eventChan)
	select {
	case <-eventChan: // new data arrives through eventChan channel
		config.TLSConfMutex.RLock()
		schemaFile, err := ioutil.ReadFile(config.Data.SearchAndFilterSchemaPath)
		if err != nil {
			log.Error("error while trying to read search/filter schema json" + err.Error())
		}
		config.TLSConfMutex.RUnlock()
		err = json.Unmarshal(schemaFile, &SF)
		if err != nil {
			log.Error("error while trying to fetch search/filter schema json" + err.Error())
		}
	}
}
