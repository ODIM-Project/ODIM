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

// Package evcommon ...
package evcommon

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-events/consumer"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/ODIM-Project/ODIM/svc-events/evresponse"
	"github.com/google/uuid"
)

// global variables
var (
	DefaultSubscriptionID        = "0"
	SubscriptionChannelKey       = "__key*__:Subscription"
	DeviceSubscriptionChannelKey = "__key*__:DeviceSubscription"
	AggregateToHostChannelKey    = "__key*__:AggregateToHost"
	RedisNotifierType            = "notify-keyspace-events"
	RedisNotifierFilterKey       = "Kz"
)

// StartUpInterface Holds the function pointer of  external interface functions
type StartUpInterface struct {
	DecryptPassword                  func([]byte) ([]byte, error)
	EMBConsume                       func(context.Context, string)
	GetAllPlugins                    func() ([]common.Plugin, *errors.Error)
	GetAllSystems                    func() ([]string, error)
	GetSingleSystem                  func(string) (string, error)
	GetPluginData                    func(string) (*common.Plugin, *errors.Error)
	GetEvtSubscriptions              func(string) ([]evmodel.SubscriptionResource, error)
	GetDeviceSubscriptions           func(string) (*common.DeviceSubscription, error)
	UpdateDeviceSubscriptionLocation func(common.DeviceSubscription) error
}

var (
	//GetAllPluginsFunc is pointer function evmodel.GetAllPlugins
	GetAllPluginsFunc = evmodel.GetAllPlugins
	// ConfigFilePath holds the value of odim config file path
	ConfigFilePath string
)

// EmbTopic hold the list all consuming topics after
type EmbTopic struct {
	TopicsList map[string]bool
	lock       sync.RWMutex
	EMBConsume func(context.Context, string)
}

// SavedSystems holds the resource details of the saved system
type SavedSystems struct {
	ManagerAddress string
	Password       []byte
	UserName       string
	DeviceUUID     string
	PluginID       string
}

// PluginContactRequest holds the details required to contact the plugin
type PluginContactRequest struct {
	URL             string
	HTTPMethodType  string
	ContactClient   func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	PostBody        interface{}
	LoginCredential map[string]string
	Token           string
	Plugin          *common.Plugin
}

// StartUpMap holds required data for plugin startup
type StartUpMap struct {
	Location   string
	EventTypes []string
	Device     SavedSystems
}

// PluginToken interface to hold the token
type PluginToken struct {
	Tokens map[string]string
	lock   sync.Mutex
}

// Token variable hold the all the XAuthToken  against the plguin ID
var Token PluginToken

// StoreToken to store the token ioto the  map
func (p *PluginToken) StoreToken(plguinID, token string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.Tokens[plguinID] = token
}

// GetToken to get the token from map
func (p *PluginToken) GetToken(pluginID string) string {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.Tokens[pluginID]
}

// ConsumeTopic check the existing topic list if it is not present then it will add topic name to list and consume that topic
func (e *EmbTopic) ConsumeTopic(ctx context.Context, topicName string) {
	e.lock.RLock()
	defer e.lock.RUnlock()
	if ok := e.TopicsList[topicName]; !ok {
		go consumer.Consume(ctx, topicName)
		e.TopicsList[topicName] = true
		//consume the topic
	}
}

// EMBTopics used to store the list of all topics
var EMBTopics EmbTopic

// PluginStartUp is used to call plugin "Startup" only on plugin restart and not on every status check
var PluginStartUp = false

// GetAllPluginStatus ...
func (st *StartUpInterface) GetAllPluginStatus(ctx context.Context) {
	for {
		pluginList, err := evmodel.GetAllPlugins()
		if err != nil {
			l.LogWithFields(ctx).Error(err.Error())
			return
		}
		var threadID int = 1
		for i := 0; i < len(pluginList); i++ {
			ctx = context.WithValue(ctx, common.ThreadID, strconv.Itoa(threadID))
			go st.getPluginStatus(ctx, pluginList[i])
			threadID++
		}
		var pollingTime int
		config.TLSConfMutex.RLock()
		pollingTime = config.Data.PluginStatusPolling.PollingFrequencyInMins
		config.TLSConfMutex.RUnlock()
		time.Sleep(time.Minute * time.Duration(pollingTime))
	}

}

func (st *StartUpInterface) getPluginStatus(ctx context.Context, plugin common.Plugin) {
	PluginsMap := make(map[string]bool)
	StartUpResourceBatchSize := config.Data.PluginStatusPolling.StartUpResourceBatchSize
	config.TLSConfMutex.RLock()
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
		PluginPreferredAuthType: plugin.PreferredAuthType,
		CACertificate:           &config.Data.KeyCertConf.RootCACertificate,
	}
	config.TLSConfMutex.RUnlock()
	status, _, topicsList, err := pluginStatus.CheckStatus()
	if err != nil && !status {
		PluginStartUp = false
		l.LogWithFields(ctx).Error("error While getting the status for plugin " + plugin.ID + err.Error())
		return
	}
	l.LogWithFields(ctx).Debug("Status of plugin " + plugin.ID + " is " + strconv.FormatBool(status))
	PluginsMap[plugin.ID] = status
	var allServers []SavedSystems
	for pluginID, status := range PluginsMap {
		if status && !PluginStartUp {
			allServers, err = st.getAllServers(ctx, pluginID)
			if err != nil {
				l.LogWithFields(ctx).Error("Error While getting the servers" + pluginID + err.Error())
				continue
			}
			for {
				if len(allServers) < StartUpResourceBatchSize {
					err = st.callPluginStartUp(ctx, allServers, pluginID)
					if err != nil {
						l.LogWithFields(ctx).Error("Error While trying call plugin startup" +
							pluginID + err.Error())
					}
					break
				}
				batchServers := allServers[:StartUpResourceBatchSize]
				err = st.callPluginStartUp(ctx, batchServers, pluginID)
				if err != nil {
					l.LogWithFields(ctx).Error("Error While trying call plugin startup" + pluginID + err.Error())
					continue
				}
				allServers = allServers[StartUpResourceBatchSize:]
			}
			PluginStartUp = true
		}
	}
	// Adding the topics to the list
	EMBTopics.lock.Lock()
	EMBTopics.EMBConsume = st.EMBConsume
	EMBTopics.lock.Unlock()
	for j := 0; j < len(topicsList); j++ {
		EMBTopics.ConsumeTopic(ctx, topicsList[j])
	}
}

func (st *StartUpInterface) getAllServers(ctx context.Context, pluginID string) ([]SavedSystems, error) {
	var matchedServers []SavedSystems
	allServers, err := st.GetAllSystems()
	if err != nil {
		return matchedServers, err
	}
	for i := 0; i < len(allServers); i++ {
		var s SavedSystems
		singleServer, err := st.GetSingleSystem(allServers[i])
		if err != nil {
			// skip to next member in the array.
			continue
		}
		json.Unmarshal([]byte(singleServer), &s)
		if s.PluginID == pluginID {
			decryptedPasswordByte, err := st.DecryptPassword(s.Password)
			if err != nil {
				// Frame the RPC response body and response Header below
				errorMessage := "error while trying to decrypt device password for the host: " + s.ManagerAddress + ":" + err.Error()
				l.LogWithFields(ctx).Error(errorMessage)
				continue
			}
			s.Password = decryptedPasswordByte
			matchedServers = append(matchedServers, s)
		}
	}
	return matchedServers, err
}

// GetPluginStatus checks the status of given plugin in configured interval
func GetPluginStatus(ctx context.Context, plugin *common.Plugin) bool {
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
		PluginPreferredAuthType: plugin.PreferredAuthType,
		CACertificate:           &config.Data.KeyCertConf.RootCACertificate,
	}
	status, _, _, err := pluginStatus.CheckStatus()
	if err != nil && !status {
		l.LogWithFields(ctx).Error("Error While getting the status for plugin " + plugin.ID + err.Error())
		return status
	}
	l.LogWithFields(ctx).Info("Status of plugin" + plugin.ID + strconv.FormatBool(status))
	return status
}

func (st *StartUpInterface) callPluginStartUp(ctx context.Context, servers []SavedSystems, pluginID string) error {
	var startUpMap []StartUpMap
	plugin, errs := st.GetPluginData(pluginID)
	if errs != nil {
		return errs
	}
	for _, server := range servers {
		var s StartUpMap
		var err error
		s.Location, s.EventTypes, err = st.getSubscribedEventsDetails(server.ManagerAddress)
		if err != nil {
			l.LogWithFields(ctx).Error("Error while retrieving the Subsection details from DB for device: " +
				server.ManagerAddress + " " + err.Error())
			continue
		}
		s.Device = server
		startUpMap = append(startUpMap, s)
	}
	var contactRequest PluginContactRequest

	contactRequest.Plugin = plugin
	contactRequest.URL = "/ODIM/v1/Startup"
	contactRequest.HTTPMethodType = http.MethodPost
	contactRequest.PostBody = startUpMap

	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		contactRequest.HTTPMethodType = http.MethodPost
		contactRequest.PostBody = map[string]interface{}{
			"Username": plugin.Username,
			"Password": string(plugin.Password),
		}
		contactRequest.URL = "/ODIM/v1/Sessions"
		response, err := callPlugin(ctx, contactRequest)
		if err != nil {
			return err
		}
		contactRequest.Token = response.Header.Get("X-Auth-Token")
	} else {
		contactRequest.LoginCredential = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
	}
	response, err := callPlugin(ctx, contactRequest)
	if err != nil {
		return err
	}

	//return updateDeviceSubscriptionLocation(startUpMap[0].Device.ManagerAddress, response.Header.Get("location"))
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		l.LogWithFields(ctx).Error(err.Error())
		return err
	}
	var r map[string]string
	json.Unmarshal(bodyBytes, &r)
	return updateDeviceSubscriptionLocation(ctx, r)
}

func callPlugin(ctx context.Context, req PluginContactRequest) (*http.Response, error) {
	var reqURL = "https://" + req.Plugin.IP + ":" + req.Plugin.Port + req.URL
	if strings.EqualFold(req.Plugin.PreferredAuthType, "XAuthToken") {
		return pmbhandle.ContactPlugin(ctx, reqURL, req.HTTPMethodType, "", "", req.PostBody, nil)
	}
	if strings.EqualFold(req.Plugin.PreferredAuthType, "BasicAuth") {
		return pmbhandle.ContactPlugin(ctx, reqURL, req.HTTPMethodType, "", "", req.PostBody, req.LoginCredential)
	}
	return pmbhandle.ContactPlugin(ctx, reqURL, req.HTTPMethodType, req.Token, "", req.PostBody, nil)

}

func (st *StartUpInterface) getSubscribedEventsDetails(serverAddress string) (string, []string, error) {
	var location string
	var eventTypes []string
	var emptyListFlag bool

	deviceIPAddress, err := common.GetIPFromHostName(serverAddress)
	if err != nil {
		return "", nil, err
	}
	searchKey := GetSearchKey(deviceIPAddress, evmodel.DeviceSubscriptionIndex)
	deviceSubscription, err := st.GetDeviceSubscriptions(searchKey)
	if err != nil {
		return "", nil, err
	}
	location = deviceSubscription.Location

	searchKey = GetSearchKey(deviceIPAddress, evmodel.SubscriptionIndex)
	subscriptionDetails, err := st.GetEvtSubscriptions(searchKey)
	if err != nil {
		return "", nil, err
	}
	for i := 0; i < len(subscriptionDetails); i++ {
		if len(subscriptionDetails[i].EventDestination.EventTypes) == 0 {
			emptyListFlag = true
		} else {
			eventTypes = append(eventTypes, subscriptionDetails[i].EventDestination.EventTypes...)
		}
	}
	if emptyListFlag {
		eventTypes = []string{}
	} else {
		eventTypes = removeDuplicates(eventTypes)
	}
	return location, eventTypes, nil
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

// getTypes is to split the string to array
func getTypes(subscription string) []string {
	// array stored in db in string("[alert statuschange]")
	// to convert into an array removing "[" ,"]" and splitting
	events := strings.Replace(subscription, "[", "", -1)
	events = strings.Replace(events, "]", "", -1)
	if len(events) < 1 {
		return []string{}
	}
	return strings.Split(events, " ")
}

func updateDeviceSubscriptionLocation(ctx context.Context, r map[string]string) error {
	for serverAddress, location := range r {
		if location != "" {
			deviceIPAddress, err := common.GetIPFromHostName(serverAddress)
			if err != nil {
				continue
			}
			searchKey := GetSearchKey(deviceIPAddress, evmodel.DeviceSubscriptionIndex)
			deviceSubscription, err := evmodel.GetDeviceSubscriptions(searchKey)
			if err != nil {
				l.LogWithFields(ctx).Error("Error getting the device event subscription from DB " +
					" for server address : " + serverAddress + err.Error())
				continue
			}
			var updatedDeviceSubscription common.DeviceSubscription

			updatedDeviceSubscription.Location = location
			updatedDeviceSubscription.EventHostIP = deviceSubscription.EventHostIP
			updatedDeviceSubscription.OriginResources = deviceSubscription.OriginResources
			err = evmodel.UpdateDeviceSubscriptionLocation(updatedDeviceSubscription)
			if err != nil {
				l.LogWithFields(ctx).Error("Error updating the subscription location in to DB for " +
					"server address : " + serverAddress + err.Error())
				continue
			}
		}
	}
	return nil
}

// GenErrorResponse generates the error response in event service
func GenErrorResponse(errorMessage string, statusMessage string, httpStatusCode int32, msgArgs []interface{}, respPtr *response.RPC) {
	respPtr.StatusCode = httpStatusCode
	respPtr.StatusMessage = statusMessage
	args := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			{
				StatusMessage: statusMessage,
				ErrorMessage:  errorMessage,
				MessageArgs:   msgArgs,
			},
		},
	}
	respPtr.Body = args.CreateGenericErrorResponse()
}

// GenEventErrorResponse generates the error response in event service
func GenEventErrorResponse(errorMessage string, StatusMessage string, httpStatusCode int, respPtr *evresponse.EventResponse, argsParams []interface{}) {
	respPtr.StatusCode = httpStatusCode
	args := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			{
				StatusMessage: StatusMessage,
				ErrorMessage:  errorMessage,
				MessageArgs:   argsParams,
			},
		},
	}
	respPtr.Response = args.CreateGenericErrorResponse()

}

// GetSearchKey will return search key with regular expression for filtering
func GetSearchKey(key, index string) string {
	searchKey := key
	if index == evmodel.SubscriptionIndex {
		searchKey = `[^0-9]` + key + `[^0-9]`
	} else if index == evmodel.DeviceSubscriptionIndex {
		searchKey = key + `[^0-9]`
	}
	return searchKey
}

// ProcessCtrlMsg is for processing the ODIM control message
// and to perform required action
func ProcessCtrlMsg(ctx context.Context, data interface{}) bool {
	if data == nil {
		l.LogWithFields(ctx).Warn("received control message event with empty data")
		return false
	}
	event := data.(common.ControlMessageData)
	msg, _ := json.Marshal(event.Data)
	l.LogWithFields(ctx).Info("received control message event of type:", event.MessageType)
	if event.MessageType == common.SubscribeEMB {
		var message common.SubscribeEMBData
		if err := json.Unmarshal([]byte(msg), &message); err != nil {
			return false
		}
		for _, topic := range message.EMBQueues {
			EMBTopics.ConsumeTopic(ctx, topic)
		}
	}
	return true
}

// SubscribePluginEMB is for subscribing to plugin EMB
func (st *StartUpInterface) SubscribePluginEMB(ctx context.Context) {
	time.Sleep(time.Second * 2)
	transactionID := uuid.New()
	ctx = context.WithValue(ctx, common.TransactionID, transactionID.String())
	ctx = context.WithValue(ctx, common.ActionName, "SubscribePluginEMB")
	pluginList, err := GetAllPluginsFunc()
	if err != nil {
		l.LogWithFields(ctx).Error(err.Error())
		return
	}
	threadID := 1
	for i := 0; i < len(pluginList); i++ {
		ctx = context.WithValue(ctx, common.ThreadID, strconv.Itoa(threadID))
		go st.getPluginEMB(ctx, pluginList[i])
		threadID++
	}
}

func (st *StartUpInterface) getPluginEMB(ctx context.Context, plugin common.Plugin) {
	config.TLSConfMutex.RLock()
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
		PluginPreferredAuthType: plugin.PreferredAuthType,
		CACertificate:           &config.Data.KeyCertConf.RootCACertificate,
	}
	config.TLSConfMutex.RUnlock()
	status, _, topicsList, err := pluginStatus.CheckStatus()
	if err != nil && !status {
		l.LogWithFields(ctx).Error("status check of plugin " + plugin.ID + " failed: " + err.Error())
		return
	}
	EMBTopics.lock.Lock()
	EMBTopics.EMBConsume = st.EMBConsume
	EMBTopics.lock.Unlock()
	for j := 0; j < len(topicsList); j++ {
		EMBTopics.ConsumeTopic(ctx, topicsList[j])
	}
}

// TrackConfigFileChanges monitors the config changes using fsnotfiy
func TrackConfigFileChanges(ctx context.Context, errChan chan error) {
	eventChan := make(chan interface{})
	format := config.Data.LogFormat
	go common.TrackConfigFileChanges(ConfigFilePath, eventChan, errChan)
	for {
		select {
		case info := <-eventChan:
			transactionID := uuid.New()
			ctx = context.WithValue(ctx, common.TransactionID, transactionID.String())
			ctx = context.WithValue(ctx, common.ActionName, "TrackConfigFileChanges")
			l.LogWithFields(ctx).Info(info) // new data arrives through eventChan channel
			if l.Log.Level != config.Data.LogLevel {
				l.LogWithFields(ctx).Debug("Log level is updated, new log level is ", config.Data.LogLevel)
				l.LogWithFields(ctx).Logger.SetLevel(config.Data.LogLevel)
			}
			if format != config.Data.LogFormat {
				l.SetFormatter(config.Data.LogFormat)
				format = config.Data.LogFormat
				l.LogWithFields(ctx).Debug("Log format is updated, new log format is ", config.Data.LogFormat)
			}
		case err := <-errChan:
			l.LogWithFields(ctx).Error(err)
		}
	}
}
