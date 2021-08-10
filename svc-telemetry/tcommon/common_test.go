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

package tcommon

import (
	"bytes"
	"fmt"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/svc-telemetry/tmodel"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func stubDevicePassword(password []byte) ([]byte, error) {
	return password, nil
}

func mockContactClient(url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {

	if url == "https://localhost:9091/ODIM/v1/Sessions" {
		body := `{"Token": "12345"}`
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			Header: http.Header{
				"X-Auth-Token": []string{"12345"},
			},
		}, nil
	} else if url == "https://localhost:9092/ODIM/v1/Sessions" {
		body := `{"Token": ""}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "https://localhost:9091/ODIM/v1/TelemetryService/MetricReports/CPUUtilCustom1" && token == "12345" {
		body := `{"@odata.id":"/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom1","@odata.type":"#MetricReport.v1_0_0.MetricReport","Id":"CPUUtilCustom1","Name":"Metric report of CPU Utilization for 10 minutes with sensing interval of 20 seconds.","MetricReportDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricReportDefinitions/CPUUtilCustom1"},"MetricValues":[{"MetricDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"},"MetricId":"CPUUtil","MetricValue":"0","Timestamp":"2021-06-16T07:59:43Z"},{"MetricDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"},"MetricId":"CPUUtil","MetricValue":"0","Timestamp":"2021-06-16T08:00:04Z"}]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9093/ODIM/v1/TelemetryService/MetricReports/CPUUtilCustom1" {
		body := `{"@odata.id":"/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom1","@odata.type":"#MetricReport.v1_0_0.MetricReport","Id":"CPUUtilCustom1","Name":"Metric report of CPU Utilization for 10 minutes with sensing interval of 20 seconds.","MetricReportDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricReportDefinitions/CPUUtilCustom1"},"MetricValues":[{"MetricDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"},"MetricId":"CPUUtil","MetricValue":"0","Timestamp":"2021-06-16T07:59:43Z"},{"MetricDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"},"MetricId":"CPUUtil","MetricValue":"0","Timestamp":"2021-06-16T08:00:04Z"}]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9092/ODIM/v1/TelemetryService/MetricReports/CPUUtilCustom1" && token == "23456" {
		body := `{"data": "ODIM/v1/TelemetryService/MetricReports/CPUUtilCustom1"}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9091/ODIM/v1/TelemetryService/MetricReports/CPUUtilCustom1" {
		body := `{"@odata.id": "/ODIM/v1/TelemetryService/MetricReports/CPUUtilCustom1"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockPluginStatus(plugin tmodel.Plugin) bool {
	return true
}

func mockGetAllKeysFromTable(table string, dbtype common.DbType) ([]string, error) {
	return []string{"ILO", "GRF"}, nil
}

func mockGetPluginData(pluginID string) (tmodel.Plugin, *errors.Error) {
	var t *testing.T
	password := getEncryptedKey(t, []byte("$2a$10$OgSUYvuYdI/7dLL5KkYNp.RCXISefftdj.MjbBTr6vWyNwAvht6ci"))
	plugin := tmodel.Plugin{
		IP:                "localhost",
		Port:              "9091",
		Username:          "admin",
		Password:          password,
		ID:                pluginID,
		PreferredAuthType: "BasisAuth",
		PluginType:        "Compute",
	}
	return plugin, nil
}

func mockGetResource(table, key string, dbType common.DbType) (string, *errors.Error) {
	if key == "/redfish/v1/TelemetryService/MetricReports" {
		return `{
			"@odata.context": "/redfish/v1/$metadata#MetricReportCollection.MetricReportCollection",
			"@odata.id": "/redfish/v1/TelemetryService/MetricReports",
			"@odata.type": "#MetricReportCollection.MetricReportCollection",
			"Description": " Metric Reports view",
			"Name": "Metric Reports",
			"Members": [
			{
				"@odata.id": "/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom1"
			},
			{
				"@odata.id": "/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom2"
			}
			],
			"Members@odata.count": 2
		}`, nil
	}
	return "body", nil
}

func mockGenericSave(body []byte, table string, key string) error {
	return nil
}

func TestGetResourceInfoFromDevice(t *testing.T) {
	config.SetUpMockConfig(t)

	var req = ResourceInfoRequest{
		URL:                 "/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom1",
		ContactClient:       mockContactClient,
		DevicePassword:      stubDevicePassword,
		GetAllKeysFromTable: mockGetAllKeysFromTable,
		GetPluginData:       mockGetPluginData,
		GetResource:         mockGetResource,
		GenericSave:         mockGenericSave,
	}
	_, err := GetResourceInfoFromDevice(req)
	assert.Nil(t, err, "There should be no error getting data")
	req.URL = "/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom1"
	_, err = GetResourceInfoFromDevice(req)
	assert.Nil(t, err, "There should be no error getting data")
}

func TestContactPlugin(t *testing.T) {
	config.SetUpMockConfig(t)
	plugin := tmodel.Plugin{}
	var contactRequest PluginContactRequest

	contactRequest.ContactClient = mockContactClient
	contactRequest.Plugin = plugin
	contactRequest.GetPluginStatus = mockPluginStatus
	_, _, _, err := ContactPlugin(contactRequest, "")
	assert.NotNil(t, err, "There should be an error")
}
