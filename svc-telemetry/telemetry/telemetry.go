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

//Package telemetry ...
package telemetry

// ---------------------------------------------------------------------------------------
// IMPORT Section
// ---------------------------------------------------------------------------------------
import (
	"net/http"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	tlresp "github.com/ODIM-Project/ODIM/svc-telemetry/tlresponse"
)

// GetTelemetryService defines the functionality for knowing whether
// the telemetry service is enabled or not
//
// As return parameters RPC response, which contains status code, message, headers and data,
// error will be passed back.
func (e *ExternalInterface) GetTelemetryService() response.RPC {
	commonResponse := response.Response{
		OdataType:    "#TelemetryService.v1_2.TelemetryService",
		OdataID:      "/redfish/v1/TelemetryService",
		OdataContext: "/redfish/v1/$metadata#TelemetryService.TelemetryService",
		ID:           "TelemetryService",
		Name:         "Telemetry Service",
	}
	var resp response.RPC

	isServiceEnabled := false
	serviceState := "Disabled"
	//Checks if TelemetryService is enabled and sets the variable isServiceEnabled to true add servicState to enabled
	for _, service := range config.Data.EnabledServices {
		if service == "TelemetryService" {
			isServiceEnabled = true
			serviceState = "Enabled"
		}
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	resp.Header = map[string]string{
		"Allow":         "GET",
		"Cache-Control": "no-cache",
		"Connection":    "Keep-alive",
		"Content-type":  "application/json; charset=utf-8",
		"Link": "	</redfish/v1/SchemaStore/en/TelemetryService.json>; rel=describedby",
		"Transfer-Encoding": "chunked",
		"X-Frame-Options":   "sameorigin",
	}

	commonResponse.CreateGenericResponse(resp.StatusMessage)
	commonResponse.Message = ""
	commonResponse.MessageID = ""
	commonResponse.Severity = ""
	resp.Body = tlresp.TelemetryService{
		Response: commonResponse,
		Status: tlresp.Status{
			State:        serviceState,
			Health:       "OK",
			HealthRollup: "OK",
		},
		ServiceEnabled: isServiceEnabled,
		MetricDefinitions: &dmtf.Link{
			Oid: "/redfish/v1/TelemetryService/MetricDefinitions",
		},
		MetricReportDefinitions: &dmtf.Link{
			Oid: "/redfish/v1/TelemetryService/MetricReportDefinitions",
		},
		MetricReports: &dmtf.Link{
			Oid: "/redfish/v1/TelemetryService/MetricReports",
		},
		Triggers: &dmtf.Link{
			Oid: "/redfish/v1/TelemetryService/Triggers",
		},
	}

	return resp

}
