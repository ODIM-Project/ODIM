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
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/google/uuid"
	uu "github.com/satori/go.uuid"
)

// DBInterface hold interface for db functions
type DBInterface struct {
	GetAllKeysFromTableInterface func(context.Context, string) ([]string, error)
	GetConnectionMethodInterface func(context.Context, string) (agmodel.ConnectionMethod, *errors.Error)
	AddConnectionMethodInterface func(agmodel.ConnectionMethod, string) *errors.Error
	DeleteInterface              func(string, string, common.DbType) *errors.Error
}

// DBDataInterface holds interface for plugin and all db data functions
type DBDataInterface struct {
	GetAllKeysFromTableFunc func(context.Context, string) ([]string, error)
	GetPluginData           func(string, agmodel.DBPluginDataRead) (agmodel.Plugin, *errors.Error)
	GetDatabaseConnection   func(PluginID string) (string, *errors.Error)
}

// PluginHealthCheckInterface holds the methods required for plugin healthcheck
type PluginHealthCheckInterface struct {
	DecryptPassword func([]byte) ([]byte, error)
	PluginConfig    config.PluginStatusPolling
	RootCA          []byte
}

// PluginStatusRecord holds the record of plugins and the
// number of times is has been inactive during periodic health check
type PluginStatusRecord struct {
	InactiveCount map[string]int
	Lock          sync.Mutex
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
var (
	//GetResourceDetailsFunc function pointer for the agmodel.GetResourceDetails
	GetResourceDetailsFunc = agmodel.GetResourceDetails
	//GetResourceDetailsBytableNameFunc function pointer for the agmodel.GetResourceDetailsBytableName
	GetResourceDetailsBytableNameFunc = agmodel.GetResourceDetailsBytableName
	// GetAllKeysFromTableFunc function pointer for the agmodel.GetAllKeysFromTable
	GetAllKeysFromTableFunc = agmodel.GetAllKeysFromTable
	//GetAllSystemsFunc function pointer for the agmodel.GetAllSystems
	GetAllSystemsFunc = agmodel.GetAllSystems
	//GetDeviceSubscriptionsFunc  function pointer for the  agmodel.GetDeviceSubscriptions
	GetDeviceSubscriptionsFunc = agmodel.GetDeviceSubscriptions
	// UpdateDeviceSubscriptionFunc function pointer for the agmodel.UpdateDeviceSubscription
	UpdateDeviceSubscriptionFunc = agmodel.UpdateDeviceSubscription
	// GetEventSubscriptionsFunc function pointer for the agmodel.GetEventSubscriptions
	GetEventSubscriptionsFunc = agmodel.GetEventSubscriptions
	// JSONUnMarshalFunc function pointer for the json.Unmarshal
	JSONUnMarshalFunc = json.Unmarshal
	//LookupIPfunc  function pointer for the  net.LookupIP
	LookupIPfunc = net.LookupIP
	//SplitHostPortfunc  function pointer for the net.SplitHostPort
	SplitHostPortfunc = net.SplitHostPort
)
var (
	// ConfigFilePath holds the value of odim config file path
	ConfigFilePath string
	// PSRecord holds the record of each plugin health check status
	PSRecord PluginStatusRecord
)

// init is for intializing global variables defined in this package
func init() {
	PSRecord = PluginStatusRecord{
		Lock:          sync.Mutex{},
		InactiveCount: make(map[string]int),
	}
}

// GetStorageResources will get the resource details from the database for the given odata id
func GetStorageResources(ctx context.Context, oid string) map[string]interface{} {
	resourceData := make(map[string]interface{})
	data, dbErr := GetResourceDetailsFunc(ctx, oid)
	if dbErr != nil {
		l.LogWithFields(ctx).Error("Unable to get system data : " + dbErr.Error())
		return resourceData
	}
	// unmarshal the resourceData
	err := JSONUnMarshalFunc([]byte(data), &resourceData)
	if err != nil {
		l.LogWithFields(ctx).Error("Unable to unmarshal  the data: " + err.Error())
		return resourceData
	}

	return resourceData
}

// GetStorageResourcesBytableName will get the resource details from the database for the given odata id and table name
func GetStorageResourcesBytableName(ctx context.Context, table, oid string) map[string]interface{} {
	resourceData := make(map[string]interface{})
	data, dbErr := GetResourceDetailsBytableNameFunc(ctx, table, oid)
	if dbErr != nil {
		l.LogWithFields(ctx).Error("Unable to get system data : " + dbErr.Error())
		return resourceData
	}
	// unmarshal the resourceData
	err := JSONUnMarshalFunc([]byte(data), &resourceData)
	if err != nil {
		l.LogWithFields(ctx).Error("Unable to unmarshal the data: " + err.Error())
		return resourceData
	}

	return resourceData
}

// AddConnectionMethods will add the connection method type and variant into DB
func (e *DBInterface) AddConnectionMethods(connectionMethodConf []config.ConnectionMethodConf) error {
	aggTransactionID := uuid.New()
	podName := os.Getenv("POD_NAME")
	actionID := common.Actions[common.ActionKey{Service: "AddConnectionMethods", URI: "ConnectionMethod", Method: "POST-TO-DB"}].ActionID
	actionName := common.Actions[common.ActionKey{Service: "AddConnectionMethods", URI: "ConnectionMethod", Method: "POST-TO-DB"}].ActionName
	ctx := CreateContext(aggTransactionID.String(), actionID, actionName, "1", common.AggregationService, podName)
	connectionMethodsKeys, err := e.GetAllKeysFromTableInterface(ctx, "ConnectionMethod")
	if err != nil {
		l.Log.Error("Unable to get connection methods : " + err.Error())
		return err
	}
	var connectionMethodInfo = make(map[string]agmodel.ConnectionMethod)
	var connectionMehtodIDMap = make(map[string]string)
	// Get all existing connection method info store it in above two map
	for i := 0; i < len(connectionMethodsKeys); i++ {
		connectionmethod, err := e.GetConnectionMethodInterface(ctx, connectionMethodsKeys[i])
		if err != nil {
			l.Log.Error("Unable to get connection method : " + err.Error())
			return err
		}
		connectionMethodInfo[connectionMethodsKeys[i]] = connectionmethod
		connectionMehtodIDMap[connectionmethod.ConnectionMethodType+":"+connectionmethod.ConnectionMethodVariant] = connectionMethodsKeys[i]
	}
	for i := 0; i < len(connectionMethodConf); i++ {
		if !SupportedConnectionMethodTypes[connectionMethodConf[i].ConnectionMethodType] {
			l.Log.Error("Connection method type " + connectionMethodConf[i].ConnectionMethodType + " is not supported")
			continue
		}
		if connectionMethodID, present := connectionMehtodIDMap[connectionMethodConf[i].ConnectionMethodType+":"+connectionMethodConf[i].ConnectionMethodVariant]; present {
			l.Log.Error("Connection Method Info with Connection Method Type " +
				connectionMethodConf[i].ConnectionMethodType + " and Connection Method Variant " +
				connectionMethodConf[i].ConnectionMethodVariant + " already present in ODIM")
			delete(connectionMehtodIDMap,
				connectionMethodConf[i].ConnectionMethodType+":"+connectionMethodConf[i].ConnectionMethodVariant)
			delete(connectionMethodInfo, connectionMethodID)
		} else {
			connectionMethodURI := "/redfish/v1/AggregationService/ConnectionMethods/" + uu.NewV4().String()
			connectionMethod := agmodel.ConnectionMethod{
				ConnectionMethodType:    connectionMethodConf[i].ConnectionMethodType,
				ConnectionMethodVariant: connectionMethodConf[i].ConnectionMethodVariant,
				Links: agmodel.Links{
					AggregationSources: []agmodel.OdataID{},
				},
			}
			err := e.AddConnectionMethodInterface(connectionMethod, connectionMethodURI)
			if err != nil {
				l.Log.Error("Unable to add connection method : " + err.Error())
				return err
			}
			l.LogWithFields(ctx).Info(
				"Connection method info with connection method type " + connectionMethodConf[i].ConnectionMethodType +
					" and connection method variant " + connectionMethodConf[i].ConnectionMethodVariant + " added to ODIM")
		}
	}
	// loop thorugh the remaining connection method data from connectionMethodInfo map
	// delete the connection from ODIM if doesn't manage any aggreation source else log the error
	for connectionMethodID, connectionMethodData := range connectionMethodInfo {
		if len(connectionMethodData.Links.AggregationSources) > 0 {
			l.Log.Error("Connection Method ID: " + connectionMethodID + " with connection method type " +
				connectionMethodData.ConnectionMethodType + " and connection method variant " +
				connectionMethodData.ConnectionMethodVariant + " managing " +
				string(rune(len(connectionMethodData.Links.AggregationSources))) + " aggregation sources it can't be removed")

		} else {
			l.LogWithFields(ctx).Info("Removing connection method id "+connectionMethodID+
				" with Connection Method Type"+connectionMethodData.ConnectionMethodType+
				" and Connection Method Variant", connectionMethodData.ConnectionMethodVariant)
			err := e.DeleteInterface("ConnectionMethod", connectionMethodID, common.OnDisk)
			if err != nil {
				l.Log.Error("Unable to removing connection method : " + err.Error())
				return err
			}
		}
	}
	return nil
}

// TrackConfigFileChanges monitors the odim config changes using fsnotfiy
// Whenever  any config file changes and events  will be  and  reload the configuration and verify the existing connection methods
func TrackConfigFileChanges(dbInterface DBInterface, errChan chan error) {
	trackTransactionID := uuid.New()
	podName := os.Getenv("POD_NAME")
	actionID := common.Actions[common.ActionKey{Service: "TrackConfigFileChanges", URI: "TrackFile", Method: "GET"}].ActionID
	actionName := common.Actions[common.ActionKey{Service: "TrackConfigFileChanges", URI: "TrackFile", Method: "GET"}].ActionName
	ctx := CreateContext(trackTransactionID.String(), actionID, actionName, "1", common.AggregationService, podName)
	eventChan := make(chan interface{})
	format := config.Data.LogFormat
	go common.TrackConfigFileChanges(ConfigFilePath, eventChan, errChan)
	for {
		select {
		case info := <-eventChan:
			l.LogWithFields(ctx).Info(info) // new data arrives through eventChan channel
			config.TLSConfMutex.RLock()
			l.LogWithFields(ctx).Info("Updating connection method ")
			err := dbInterface.AddConnectionMethods(config.Data.ConnectionMethodConf)
			if err != nil {
				l.Log.Error("error while trying to Add connection methods:" + err.Error())
			}
			config.TLSConfMutex.RUnlock()
			l.LogWithFields(ctx).Info("Update connection method completed")
			if l.Log.Level != config.Data.LogLevel {
				l.LogWithFields(ctx).Info("Log level is updated, new log level is ", config.Data.LogLevel)
				l.Log.Logger.SetLevel(config.Data.LogLevel)
			}
			if format != config.Data.LogFormat {
				l.SetFormatter(config.Data.LogFormat)
				format = config.Data.LogFormat
				l.LogWithFields(ctx).Info("Log format is updated, new log format is ", config.Data.LogFormat)
			}
		case err := <-errChan:
			l.Log.Error(err)
		}
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
	phc.PluginConfig.StartUpResourceBatchSize = config.Data.PluginStatusPolling.StartUpResourceBatchSize
	phc.RootCA = make([]byte, len(config.Data.KeyCertConf.RootCACertificate))
	copy(phc.RootCA, config.Data.KeyCertConf.RootCACertificate)
}

// GetPluginStatus checks the status of given plugin
func GetPluginStatus(ctx context.Context, plugin agmodel.Plugin) bool {
	phc := &PluginHealthCheckInterface{}
	phc.DupPluginConf()
	status, _ := phc.GetPluginStatus(ctx, plugin)
	l.LogWithFields(ctx).Debug("Status of plugin" + plugin.ID + strconv.FormatBool(status))
	return status
}

// LookupHost - look up the ip from the host address
func LookupHost(addr string) (ip, host, port string, err error) {
	host, port, err = SplitHostPortfunc(addr)
	if err != nil {
		host = addr
	}

	ips, errs := LookupIPfunc(host)
	switch {
	case errs != nil:
		err = errs
	case len(ips) < 1:
		err = fmt.Errorf("host lookup gave empty list")
	default:
		err = nil
		ip = ips[0].String()
	}
	return
}

// LookupPlugin is for fetching the plugin data
// using the plugin address for lookup
func LookupPlugin(ctx context.Context, addr string) (agmodel.Plugin, error) {
	phc := &PluginHealthCheckInterface{}
	phc.DupPluginConf()
	DatabaseInterface := DBDataInterface{
		GetAllKeysFromTableFunc: GetAllKeysFromTableFunc,
		GetPluginData:           agmodel.GetPluginData,
		GetDatabaseConnection:   agmodel.GetPluginDBConnection,
	}
	plugins, errs := GetAllPlugins(ctx, DatabaseInterface)
	if errs != nil {
		return agmodel.Plugin{}, errs
	}

	resolvedAddr, host, port, err := LookupHost(addr)
	if err != nil {
		l.LogWithFields(ctx).Warn("plugin address lookup failed with " + err.Error())
	}

	for _, plugin := range plugins {
		if (plugin.IP == host || plugin.IP == resolvedAddr) && (plugin.Port == port) {
			l.LogWithFields(ctx).Debug("lookup plugin IP" + plugin.ID)
			return plugin, nil
		}
	}
	return agmodel.Plugin{}, fmt.Errorf(addr + " address does not belong to any of the plugin")
}

// GetAllPlugins is for fetching all the plugins added andn stored in db.
func GetAllPlugins(ctx context.Context, dbData DBDataInterface) ([]agmodel.Plugin, error) {
	keys, err := dbData.GetAllKeysFromTableFunc(ctx, "Plugin")
	if err != nil {
		return nil, err
	}
	var plugins []agmodel.Plugin
	for _, key := range keys {
		readPlugin := agmodel.DBPluginDataRead{
			DBReadclient: agmodel.GetPluginDBConnection,
		}
		plugin, err := dbData.GetPluginData(key, readPlugin)
		if err != nil {
			l.LogWithFields(ctx).Error("failed to get details of " + key + " plugin: " + err.Error())
			continue
		}
		plugins = append(plugins, plugin)
	}
	return plugins, nil
}

// GetPluginStatus is for checking the status of a plugin
func (phc *PluginHealthCheckInterface) GetPluginStatus(ctx context.Context, plugin agmodel.Plugin) (bool, []string) {
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
		PluginPreferredAuthType: plugin.PreferredAuthType,
		CACertificate:           &phc.RootCA,
	}
	status, _, topics, err := pluginStatus.CheckStatus()
	if err != nil {
		l.LogWithFields(ctx).Error("failed to get the status of plugin " + plugin.ID + err.Error())
		return false, nil
	}
	l.LogWithFields(ctx).Info("Status of plugin " + plugin.ID + " is " + strconv.FormatBool(status))
	return status, topics
}

// GetPluginManagedServers is for fetching the list of servers managed by a plugin
func (phc *PluginHealthCheckInterface) GetPluginManagedServers(plugin agmodel.Plugin) []agmodel.Target {
	serversList, err := phc.getAllServers(plugin.ID)
	if err != nil {
		l.Log.Error("failed to get list of servers managed by " + plugin.ID + err.Error())
	}
	return serversList
}

// getAllServers is for fetching the list of all servers added.
func (phc *PluginHealthCheckInterface) getAllServers(pluginID string) ([]agmodel.Target, error) {
	var matchedServers []agmodel.Target
	allServers, err := GetAllSystemsFunc()
	if err != nil {
		l.Log.Error("failed to get the list of all managed servers " + err.Error())
		return matchedServers, err
	}
	for _, server := range allServers {
		if server.PluginID == pluginID {
			decryptedPasswordByte, err := phc.DecryptPassword(server.Password)
			if err != nil {
				l.Log.Error("failed to decrypt device password of the host: " + server.ManagerAddress + ":" + err.Error())
				continue
			}
			server.Password = decryptedPasswordByte
			matchedServers = append(matchedServers, server)
		}
	}
	return matchedServers, nil
}

// ContactPlugin is for sending requests to a plugin.
func ContactPlugin(ctx context.Context, req agmodel.PluginContactRequest, serverName string) (*http.Response, error) {
	req.LoginCredential = map[string]string{}
	//ToDo: Variable "LoginCredentials" to be changed
	req.LoginCredential["ServerName"] = serverName
	if strings.EqualFold(req.Plugin.PreferredAuthType, "XAuthToken") {
		payload := map[string]interface{}{
			"Username": req.Plugin.Username,
			"Password": string(req.Plugin.Password),
		}
		reqURL := fmt.Sprintf("https://%s/ODIM/v1/Sessions", net.JoinHostPort(req.Plugin.IP, req.Plugin.Port))
		response, err := pmbhandle.ContactPlugin(ctx, reqURL, http.MethodPost, "", "", payload, nil)
		if err != nil || (response != nil && response.StatusCode != http.StatusOK) {
			return nil,
				fmt.Errorf("failed to get session token from %s: %s: %+v", req.Plugin.ID, err.Error(), response)
		}
		req.Token = response.Header.Get("X-Auth-Token")
	} else {
		req.LoginCredential["UserName"] = req.Plugin.Username
		req.LoginCredential["Password"] = string(req.Plugin.Password)
	}
	reqURL := fmt.Sprintf("https://%s%s", net.JoinHostPort(req.Plugin.IP, req.Plugin.Port), req.URL)
	return pmbhandle.ContactPlugin(ctx, reqURL, req.HTTPMethodType, req.Token, "", req.PostBody, req.LoginCredential)
}

// GetDeviceSubscriptionDetails is for getting device event susbcription details
func GetDeviceSubscriptionDetails(ctx context.Context, serverAddress string) (string, []string, error) {
	deviceIPAddress, _, _, err := LookupHost(serverAddress)
	if err != nil {
		return "", nil, err
	}

	searchKey := GetSearchKey(deviceIPAddress, common.DeviceSubscriptionIndex)
	deviceSubscription, err := agmodel.GetDeviceSubscriptions(ctx, searchKey)
	if err != nil {
		return "", nil, err
	}

	searchKey = GetSearchKey(deviceIPAddress, common.SubscriptionIndex)
	eventTypes, err := GetSubscribedEvtTypes(ctx, searchKey)
	if err != nil {
		return "", nil, err
	}

	return deviceSubscription.Location, eventTypes, nil
}

func removeDuplicates(elements []string) []string {
	existing := map[string]bool{}
	result := []string{}

	for v := range elements {
		if !existing[elements[v]] {
			existing[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

// GetSearchKey will return search key with regular expression for filtering
func GetSearchKey(key, index string) string {
	searchKey := key
	if index == common.SubscriptionIndex {
		searchKey = `[^0-9]` + key + `[^0-9]`
	} else if index == common.DeviceSubscriptionIndex {
		searchKey = key + `[^0-9]`
	}
	return searchKey
}

// GetSubscribedEvtTypes is to get event subscription details
func GetSubscribedEvtTypes(ctx context.Context, searchKey string) ([]string, error) {
	subscriptions, err := GetEventSubscriptionsFunc("*" + searchKey + "*")
	if err != nil {
		return nil, err
	}
	var eventTypes []string
	for _, sub := range subscriptions {
		var subscription map[string]interface{}
		if err := JSONUnMarshalFunc([]byte(sub), &subscription); err != nil {
			return nil, fmt.Errorf("error while unmarshalling event subscription: %v", err.Error())
		}
		for _, evtTyps := range subscription["EventTypes"].([]interface{}) {
			eventTypes = append(eventTypes, evtTyps.(string))
		}
	}
	eventTypes = removeDuplicates(eventTypes)
	l.LogWithFields(ctx).Debug("subscribed event types:", eventTypes)
	return eventTypes, nil
}

// UpdateDeviceSubscriptionDetails is for updating the event subscription details fo a device
func UpdateDeviceSubscriptionDetails(ctx context.Context, subsData map[string]string) {
	for serverAddress, location := range subsData {
		if location != "" {
			deviceIPAddress, _, _, err := LookupHost(serverAddress)
			if err != nil {
				continue
			}
			searchKey := GetSearchKey(deviceIPAddress, common.DeviceSubscriptionIndex)
			deviceSubscription, err := GetDeviceSubscriptionsFunc(ctx, searchKey)
			if err != nil {
				l.LogWithFields(ctx).Error("Error getting the device event subscription from DB " +
					" for server address : " + serverAddress + err.Error())
				continue
			}
			deviceSubscription.Location = location
			if err = UpdateDeviceSubscriptionFunc(*deviceSubscription); err != nil {
				l.LogWithFields(ctx).Error("Error updating the subscription location in to DB for " +
					"server address : " + serverAddress + err.Error())
				continue
			}
		}
	}
}

// GetPluginStatusRecord is for getting the status record of a plugin
func GetPluginStatusRecord(plugin string) (int, bool) {
	PSRecord.Lock.Lock()
	count, exist := PSRecord.InactiveCount[plugin]
	PSRecord.Lock.Unlock()
	return count, exist
}

// SetPluginStatusRecord is for setting the status record of a plugin
func SetPluginStatusRecord(plugin string, count int) {
	PSRecord.Lock.Lock()
	PSRecord.InactiveCount[plugin] = count
	PSRecord.Lock.Unlock()
}

// CreateContext creates a new context based on transactionId, actionId, actionName, threadId, threadName, ProcessName
func CreateContext(transactionID, actionID, actionName, threadID, threadName, ProcessName string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.TransactionID, transactionID)
	ctx = context.WithValue(ctx, common.ActionID, actionID)
	ctx = context.WithValue(ctx, common.ActionName, actionName)
	ctx = context.WithValue(ctx, common.ThreadID, threadID)
	ctx = context.WithValue(ctx, common.ThreadName, threadName)
	ctx = context.WithValue(ctx, common.ProcessName, ProcessName)
	return ctx
}
