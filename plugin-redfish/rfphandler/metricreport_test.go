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
package rfphandler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/plugin-redfish/config"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpmodel"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfputilities"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func tokenValidationMock(token string) bool {
	return true
}

func getMetricReportMock(uri string, device *rfputilities.RedfishDevice) (int, []byte, map[string]interface{}, error) {
	var data string
	if uri == "/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom1" {
		data = `{
			"@odata.id": "/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom1",
			"@odata.type": "#MetricReport.v1_0_0.MetricReport",
      		"@odata.context": "/redfish/v1/$metadata#MetricReport.MetricReport",
			"Id": "CPUUtilCustom1",
			"MetricReportDefinition": {
			   "@odata.id": "/redfish/v1/TelemetryService/MetricReportDefinitions/CPUUtilCustom1"
			},
			"MetricValues": [
			   {
				  "MetricDefinition": {
					 "@odata.id": "/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"
				  },
				  "MetricId": "CPUUtil",
				  "MetricValue": "0",
				  "Timestamp": "2021-06-16T07:59:43Z"
			   },
			   {
				  "MetricDefinition": {
					 "@odata.id": "/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"
				  },
				  "MetricId": "CPUUtil",
				  "MetricValue": "0",
				  "Timestamp": "2021-06-16T08:00:04Z"
			   }
			],
			"Name": "Metric report of CPU Utilization for 10 minutes with sensing interval of 20 seconds."
		 }`
	}
	if uri == "/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil" {
		data = `{
				"@odata.context": "/redfish/v1/$metadata#MetricDefinition.MetricDefinition",
				"@odata.etag": "W/\"AB720077\"",
				"@odata.id": "/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil",
				"@odata.type": "#MetricDefinition.v1_0_0.MetricDefinition",
				"Id": "CPUUtil",
				"Calculable": "NonSummable",
				"CalculationAlgorithm": "Average",
				"Description": "Metric definition for CPU Utilization",
				"Implementation": "PhysicalSensor",
				"IsLinear": true,
				"MaxReadingRange": 100,
				"MetricDataType": "Decimal",
				"MetricProperties": [
					"/redfish/v1/Systems/1#SystemUsage/CPUUtil"
				],
				"MetricType": "Numeric",
				"MinReadingRange": 0,
				"Name": "Metric definition for CPU Utilization",
				"Units": "%"
		}`
	}
	respMap := make(map[string]interface{})
	json.Unmarshal([]byte(data), &respMap)
	return http.StatusOK, []byte(data), respMap, nil
}

func TestExternalInterface_GetMetricReport(t *testing.T) {
	config.Data.URLTranslation = &config.URLTranslation{
		NorthBoundURL: map[string]string{
			"redfish": "ODIM",
		},
		SouthBoundURL: map[string]string{
			"ODIM": "redfish",
		},
	}
	e := ExternalInterface{
		TokenValidation: tokenValidationMock,
		GetDeviceData:   getMetricReportMock,
	}
	rfpmodel.DeviceInventory.Device["0e343dc6-f5f3-425a-9503-4a3c799579c8"] = rfpmodel.DeviceData{
		Address:  "172.16.1.205",
		UserName: "admin",
		Password: []byte("Admin123"),
	}
	expectedBody := `{"@odata.id":"/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom1","@odata.type":"#MetricReport.v1_0_0.MetricReport","@odata.context":"/redfish/v1/$metadata#MetricReport.MetricReport","Id":"CPUUtilCustom1","Name":"Metric report of CPU Utilization for 10 minutes with sensing interval of 20 seconds.","MetricReportDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricReportDefinitions/CPUUtilCustom1"},"MetricValues":[{"MetricDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"},"MetricId":"CPUUtil","MetricProperty":"/redfish/v1/Systems/0e343dc6-f5f3-425a-9503-4a3c799579c8:1#SystemUsage/CPUUtil","MetricValue":"0","Timestamp":"2021-06-16T07:59:43Z"},{"MetricDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"},"MetricId":"CPUUtil","MetricProperty":"/redfish/v1/Systems/0e343dc6-f5f3-425a-9503-4a3c799579c8:1#SystemUsage/CPUUtil","MetricValue":"0","Timestamp":"2021-06-16T08:00:04Z"}]}`

	body := map[string]interface{}{}

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/ODIM/v1")
	redfishRoutes.Get("/TelemetryService/MetricReports/{id}", e.GetMetricReport)
	app := httptest.New(t, mockApp)
	app.GET(
		"/ODIM/v1/TelemetryService/MetricReports/CPUUtilCustom1",
	).WithHeader("X-Auth-Token", "Token").WithJSON(body).Expect().Status(http.StatusOK).Body().Equal(expectedBody)
}
