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

//Package rfphandler ...
package rfphandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	dmtfmodel "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	pluginConfig "github.com/ODIM-Project/ODIM/plugin-redfish/config"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpmodel"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfputilities"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
)

// ExternalInterface enables the communicunication with the external functions
type ExternalInterface struct {
	TokenValidation func(string) bool
	GetDeviceData   func(string, *rfputilities.RedfishDevice) (int, []byte, map[string]interface{}, error)
}

// GetMetricReport is for to get metric report from southbound resource
func (e *ExternalInterface) GetMetricReport(ctx iris.Context) {
	var metricData dmtfmodel.MetricReports

	//Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	reqURI := ctx.Request().RequestURI

	//replacing the request url with south bound translation URL
	for key, value := range pluginConfig.Data.URLTranslation.SouthBoundURL {
		reqURI = strings.Replace(reqURI, key, value, -1)
	}
	//Validating the token
	if token != "" {
		flag := e.TokenValidation(token)
		if !flag {
			log.Error("Invalid/Expired X-Auth-Token")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid/Expired X-Auth-Token")
			return
		}
	}

	// prepare the device data
	var devices []rfpmodel.Device
	rfpmodel.GetAllDevicesInInventory(&devices)

	metricReportData := e.getMetricData(reqURI, devices)
	var respMetricValue []dmtfmodel.MetricValue
	for _, metricData = range metricReportData {
		respMetricValue = append(respMetricValue, metricData.MetricValues...)
	}
	metricData.MetricValues = respMetricValue
	ctx.StatusCode(http.StatusOK)

	data, _ := json.Marshal(metricData)
	ctx.Write([]byte(string(data)))
	return
}

// getMetricData collects the metricreport from the BMC for the given list of BMC
func (e *ExternalInterface) getMetricData(uri string, devices []rfpmodel.Device) map[string]dmtfmodel.MetricReports {
	var wg sync.WaitGroup
	data := make(map[string]dmtfmodel.MetricReports)
	var lock sync.Mutex
	for i := 0; i < len(devices); i++ {
		wg.Add(1)
		go e.getMetricReportInfo(uri, devices[i], &wg, data, &lock)
	}
	wg.Wait()
	return data
}

func (e *ExternalInterface) getMetricReportInfo(uri string, device rfpmodel.Device, wg *sync.WaitGroup, data map[string]dmtfmodel.MetricReports, lock *sync.Mutex) {
	defer wg.Done()
	statusCode, body, _, _ := e.GetDeviceData(uri, &rfputilities.RedfishDevice{
		Host:     device.Host,
		Username: device.Username,
		Password: string(device.Password),
		PostBody: nil,
	})
	if statusCode == http.StatusOK {
		var metricReportData dmtfmodel.MetricReports
		json.Unmarshal(body, &metricReportData)
		lock.Lock()
		defer lock.Unlock()
		var respMetricValue []dmtfmodel.MetricValue
		for _, metricVal := range metricReportData.MetricValues {
			uri = metricVal.MetricDefinition.ODataID
			metricID := metricVal.MetricID
			if _, ok := rfpmodel.MetricPropertyData[metricID]; !ok {
				statusCode, body, _, _ := e.GetDeviceData(uri, &rfputilities.RedfishDevice{
					Host:     device.Host,
					Username: device.Username,
					Password: string(device.Password),
					PostBody: nil,
				})
				if statusCode == http.StatusOK {
					var metricDefinitionData dmtfmodel.MetricDefinitions
					json.Unmarshal(body, &metricDefinitionData)
					metricProperty := metricDefinitionData.MetricProperties[0]
					metricPropertyURL := strings.TrimSuffix(metricProperty, "/")
					data := strings.Split(metricPropertyURL, "/")
					metricPropertyID := data[len(data)-1]
					if metricPropertyID == metricID {
						rfpmodel.MetricPropertyData[metricID] = metricProperty
					}
				} else if statusCode != http.StatusNotFound {
					log.Info("Get on individual metric definition failed on device "+device.Host+" with status code: ", statusCode)
				}
			}
			metricProperty := rfpmodel.MetricPropertyData[metricID]
			metricProperty = strings.Replace(metricProperty, "/Systems/", "/Systems/"+device.SystemID+":", -1)
			metricProperty = strings.Replace(metricProperty, "/Chassis/", "/Chassis/"+device.SystemID+":", -1)
			metricVal.MetricProperty = metricProperty
			respMetricValue = append(respMetricValue, metricVal)
		}
		metricReportData.MetricValues = respMetricValue
		data[device.SystemID] = metricReportData
	} else if statusCode != http.StatusNotFound {
		log.Info("Get on individual metric report failed  on device "+device.Host+" with status code: ", statusCode)
	}
	return
}

// GetDeviceData connects with the device and collect the information
func GetDeviceData(uri string, device *rfputilities.RedfishDevice) (int, []byte, map[string]interface{}, error) {
	redfishClient, err := rfputilities.GetRedfishClient()
	if err != nil {
		errMsg := "While trying to get the redfish client, got: " + err.Error()
		log.Error(errMsg)
		return http.StatusInternalServerError, nil, nil, fmt.Errorf(errMsg)
	}
	//Fetching generic resource details from the device
	resp, err := redfishClient.GetWithBasicAuth(device, uri)
	if err != nil {
		errMsg := "Authentication failed: " + err.Error()
		log.Error(errMsg)
		if resp == nil {
			return http.StatusInternalServerError, nil, nil, fmt.Errorf(errMsg)
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return http.StatusUnauthorized, body, nil, fmt.Errorf("Authtication with the device failed")
	}
	if resp.StatusCode >= 300 {
		errMsg := "Could not retrieve generic resource for " + device.Host + ": " + string(body)
		log.Error(errMsg)
		return resp.StatusCode, body, nil, fmt.Errorf(errMsg)
	}

	respMap := make(map[string]interface{})
	if err := json.Unmarshal(body, &respMap); err != nil {
		errMsg := "While unmarshaling the response from device, got:" + err.Error()
		log.Error(errMsg)
		return http.StatusBadRequest, body, nil, fmt.Errorf(errMsg)
	}

	return resp.StatusCode, body, respMap, nil
}
