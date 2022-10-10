//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-licenses/model"
)

var (
	// ConfigFilePath holds the value of odim config file path
	ConfigFilePath string
)

//GetAllKeysFromTable fetches all keys in a given table
func GetAllKeysFromTable(table string, dbtype persistencemgr.DbType) ([]string, error) {
	conn, err := persistencemgr.GetDBConnection(dbtype)
	if err != nil {
		return nil, err
	}
	keysArray, err := conn.GetAllDetails(table)
	if err != nil {
		return nil, fmt.Errorf("error while trying to get all keys from table - %v: %v", table, err.Error())
	}
	return keysArray, nil
}

//GetResource fetches a resource from database using table and key
func GetResource(Table, key string, dbtype persistencemgr.DbType) (interface{}, *errors.Error) {
	conn, err := persistencemgr.GetDBConnection(dbtype)
	if err != nil {
		return "", err
	}
	resourceData, err := conn.Read(Table, key)
	if err != nil {
		return "", errors.PackError(err.ErrNo(), "error while trying to get resource details: ", err.Error())
	}
	var resource interface{}
	if errs := json.Unmarshal([]byte(resourceData), &resource); errs != nil {
		return "", errors.PackError(errors.UndefinedErrorType, errs)
	}
	return resource, nil
}

//GetTarget fetches the System(Target Device Credentials) table details
func GetTarget(deviceUUID string) (*model.Target, *errors.Error) {
	var target model.Target
	conn, err := persistencemgr.GetDBConnection(persistencemgr.OnDisk)
	if err != nil {
		return nil, err
	}
	data, err := conn.Read("System", deviceUUID)
	if err != nil {
		return nil, errors.PackError(err.ErrNo(), "error while trying to get compute details: ", err.Error())
	}
	if errs := json.Unmarshal([]byte(data), &target); errs != nil {
		return nil, errors.PackError(errors.UndefinedErrorType, errs)
	}
	return &target, nil
}

//GetPluginData will fetch plugin details
func GetPluginData(pluginID string) (*model.Plugin, *errors.Error) {
	var plugin model.Plugin

	conn, err := persistencemgr.GetDBConnection(persistencemgr.OnDisk)
	if err != nil {
		return nil, err
	}

	plugindata, err := conn.Read("Plugin", pluginID)
	if err != nil {
		return nil, errors.PackError(err.ErrNo(), "error while trying to fetch plugin data: ", err.Error())
	}

	if err := json.Unmarshal([]byte(plugindata), &plugin); err != nil {
		return nil, errors.PackError(errors.JSONUnmarshalFailed, err)
	}

	bytepw, errs := common.DecryptWithPrivateKey([]byte(plugin.Password))
	if errs != nil {
		return nil, errors.PackError(errors.DecryptionFailed, "error: "+pluginID+" plugin password decryption failed: "+errs.Error())
	}
	plugin.Password = bytepw

	return &plugin, nil
}

// ContactPlugin is commons which handles the request and response of Contact Plugin usage
func ContactPlugin(req model.PluginContactRequest, errorMessage string) ([]byte, string, model.ResponseStatus, error) {
	var resp model.ResponseStatus
	var err error
	pluginResponse, err := callPlugin(req)
	if err != nil {
		if getPluginStatus(req.Plugin) {
			pluginResponse, err = callPlugin(req)
		}
		if err != nil {
			errorMessage = errorMessage + err.Error()
			resp.StatusCode = http.StatusServiceUnavailable
			resp.StatusMessage = response.CouldNotEstablishConnection
			resp.MsgArgs = []interface{}{"https://" + req.Plugin.IP + ":" + req.Plugin.Port + req.OID}
			return nil, "", resp, fmt.Errorf(errorMessage)
		}
	}
	defer pluginResponse.Body.Close()
	body, err := ioutil.ReadAll(pluginResponse.Body)
	if err != nil {
		errorMessage := "error while trying to read response body: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = errors.InternalError
		l.Log.Warn(errorMessage)
		return nil, "", resp, fmt.Errorf(errorMessage)
	}

	if pluginResponse.StatusCode != http.StatusCreated && pluginResponse.StatusCode != http.StatusOK {
		if pluginResponse.StatusCode == http.StatusUnauthorized {
			errorMessage += "error: invalid resource username/password"
			resp.StatusCode = int32(pluginResponse.StatusCode)
			resp.StatusMessage = response.ResourceAtURIUnauthorized
			resp.MsgArgs = []interface{}{"https://" + req.Plugin.IP + ":" + req.Plugin.Port + req.OID}
			l.Log.Warn(errorMessage)
			return nil, "", resp, fmt.Errorf(errorMessage)
		}
		errorMessage += string(body)
		resp.StatusCode = int32(pluginResponse.StatusCode)
		resp.StatusMessage = response.InternalError
		l.Log.Warn(errorMessage)
		return body, "", resp, fmt.Errorf(errorMessage)
	}

	data := string(body)
	//replacing the resposne with north bound translation URL
	for key, value := range config.Data.URLTranslation.NorthBoundURL {
		data = strings.Replace(data, key, value, -1)
	}
	return []byte(data), pluginResponse.Header.Get("X-Auth-Token"), resp, nil
}

// getPluginStatus checks the status of given plugin in configured interval
func getPluginStatus(plugin model.Plugin) bool {
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
		l.Log.Warn("Error While getting the status for plugin " + plugin.ID + " " + err.Error())
		return status
	}
	l.Log.Info("Status of plugin " + plugin.ID + " " + strconv.FormatBool(status))
	return status
}

func callPlugin(req model.PluginContactRequest) (*http.Response, error) {
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

//GenericSave will save any resource data into the database
func GenericSave(body []byte, table string, key string) error {
	connPool, err := persistencemgr.GetDBConnection(persistencemgr.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.AddResourceData(table, key, string(body)); err != nil {
		if errors.DBKeyAlreadyExist == err.ErrNo() {
			return fmt.Errorf("error while trying to create new %v resource: %v", table, err.Error())
		}
		l.Log.Warn("skipped saving of duplicate data with key " + key)
	}
	return nil
}

func GetIDsFromURI(uri string) (string, string, error) {
	lastChar := uri[len(uri)-1:]
	if lastChar == "/" {
		uri = uri[:len(uri)-1]
	}
	uriParts := strings.Split(uri, "/")
	ids := strings.SplitN(uriParts[len(uriParts)-1], ".", 2)
	if len(ids) != 2 {
		return "", "", fmt.Errorf("error: no id is found in %v", uri)
	}
	return ids[0], ids[1], nil
}

func TrackConfigFileChanges() {
	eventChan := make(chan interface{})
	go common.TrackConfigFileChanges(ConfigFilePath, eventChan)
	for {
		l.Log.Info(<-eventChan) // new data arrives through eventChan channel
		if l.Log.Level != config.Data.LogLevel {
			l.Log.Info("Log level is updated, new log level is ", config.Data.LogLevel)
			l.Log.Logger.SetLevel(config.Data.LogLevel)
		}

	}
}
