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

//Package tcommon ...
package tcommon

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/svc-telemetry/tmodel"
)

//PluginContactRequest  hold the request of contact plugin
type PluginContactRequest struct {
	Token           string
	OID             string
	DeviceInfo      interface{}
	BasicAuth       map[string]string
	ContactClient   func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	GetPluginStatus func(tmodel.Plugin) bool
	Plugin          tmodel.Plugin
	HTTPMethodType  string
}

//ResponseStatus holds the response of Contact Plugin
type ResponseStatus struct {
	StatusCode    int32
	StatusMessage string
}

//ResourceInfoRequest  hold the request of getting  Resource
type ResourceInfoRequest struct {
	URL                 string
	ContactClient       func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	DevicePassword      func([]byte) ([]byte, error)
	GetPluginStatus     func(tmodel.Plugin) bool
	ResourceName        string
	GetAllKeysFromTable func(string, common.DbType) ([]string, error)
	GetPluginData       func(string) (tmodel.Plugin, *errors.Error)
	GetResource         func(string, string, common.DbType) (string, *errors.Error)
	GenericSave         func([]byte, string, string) error
}

// GetResourceInfoFromDevice will contact to the southbound client and gets the Particual resource info from device
func GetResourceInfoFromDevice(req ResourceInfoRequest) ([]byte, error) {
	var metricReportData dmtf.MetricReports
	plugins, err := req.GetAllKeysFromTable("Plugin", common.OnDisk)
	if err != nil {
		return []byte{}, err
	}
	var wg sync.WaitGroup
	var lock sync.Mutex
	for _, value := range plugins {
		wg.Add(1)
		go getResourceInfo(value, &metricReportData, req, &lock, &wg)
	}
	wg.Wait()
	if reflect.DeepEqual(metricReportData, dmtf.MetricReports{}) {
		removeNonExistingID(req)
		return []byte{}, fmt.Errorf("Metric report not found")
	}
	data, err := json.Marshal(metricReportData)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func getResourceInfo(pluginID string, metricReportData *dmtf.MetricReports, req ResourceInfoRequest, lock *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	// Get the Plugin info
	plugin, gerr := req.GetPluginData(pluginID)
	if gerr != nil || plugin.PluginType != "Compute" {
		return
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
			return
		}
		contactRequest.Token = token
	} else {
		contactRequest.BasicAuth = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
	}
	contactRequest.OID = req.URL
	contactRequest.HTTPMethodType = http.MethodGet
	body, _, _, err := ContactPlugin(contactRequest, "error while getting the details "+contactRequest.OID+": ")
	if err != nil {
		return
	}
	lock.Lock()
	defer lock.Unlock()
	var metrictData dmtf.MetricReports
	if err := json.Unmarshal(body, &metrictData); err != nil {
		return
	}
	metricReportData.ODataID = metrictData.ODataID
	metricReportData.ODataType = metrictData.ODataType
	metricReportData.ODataContext = metrictData.ODataContext
	metricReportData.ID = metrictData.ID
	metricReportData.Name = metrictData.Name
	metricReportData.Description = metrictData.Description
	metricReportData.MetricReportDefinition = metrictData.MetricReportDefinition
	metricReportData.Context = metrictData.Context
	metricReportData.MetricValues = append(metricReportData.MetricValues, metrictData.MetricValues...)
	return
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
	log.Info("Response StatusCode: " + strconv.Itoa(int(response.StatusCode)))
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

// GetPluginStatus checks the status of given plugin in configured interval
func GetPluginStatus(plugin tmodel.Plugin) bool {
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

func removeNonExistingID(req ResourceInfoRequest) {
	collectionURL := "/redfish/v1/TelemetryService/MetricReports"
	data, err := req.GetResource("MetricReportsCollection", collectionURL, common.InMemory)
	if err != nil {
		return
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)

	var result []*dmtf.Link
	if resource["Members"] != nil {
		members := resource["Members"].([]interface{})
		for _, v := range members {
			oid := v.(map[string]interface{})
			if oid["@odata.id"].(string) != req.URL {
				result = append(result, &dmtf.Link{Oid: oid["@odata.id"].(string)})
			}
		}
	}
	if len(result) > 0 {
		resource["Members"] = result
		resource["Members@odata.count"] = len(result)
		reportCollection, jerr := json.Marshal(resource)
		if jerr != nil {
			return
		}
		req.GenericSave(reportCollection, "MetricReportsCollection", collectionURL)
	}
}
