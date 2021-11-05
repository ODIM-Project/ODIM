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
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-telemetry/tcommon"
	tlresp "github.com/ODIM-Project/ODIM/svc-telemetry/tlresponse"
)

// GetTelemetryService defines the functionality for knowing whether
// the telemetry service is enabled or not
//
// As return parameters RPC response, which contains status code, message, headers and data,
// error will be passed back.
func (e *ExternalInterface) GetTelemetryService() response.RPC {
	commonResponse := response.Response{
		OdataType:    common.TelemetryServiceType,
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

// GetMetricDefinitionCollection is a functioanlity to retrive all the available inventory
// resources from the added BMC's
func (e *ExternalInterface) GetMetricDefinitionCollection(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	resp.Header = header
	data, err := e.DB.GetResource("MetricDefinitionsCollection", req.URL, common.InMemory)
	if err != nil {
		// return empty collection response
		metricDefinitionCollection := tlresp.Collection{
			OdataContext: "/redfish/v1/$metadata#MetricDefinitionCollection.MetricDefinitionCollection",
			OdataID:      "/redfish/v1/TelemetryService/MetricDefinitionCollection",
			OdataType:    "#MetricDefinitionCollection.MetricDefinitionCollection",
			Description:  "MetricDefinition Collection view",
			Name:         "MetricDefinitionCollection",
			Members:      []dmtf.Link{},
			MembersCount: 0,
		}
		resp.Body = metricDefinitionCollection
		resp.StatusCode = http.StatusOK
		return resp
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}

// GetMetricReportDefinitionCollection is a functioanlity to retrive all the available inventory
// resources from the added BMC's
func (e *ExternalInterface) GetMetricReportDefinitionCollection(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	resp.Header = header
	data, err := e.DB.GetResource("MetricReportDefinitionsCollection", req.URL, common.InMemory)
	if err != nil {
		// return empty collection response
		metricReportDefinitionCollection := tlresp.Collection{
			OdataContext: "/redfish/v1/$metadata#MetricReportDefinitionCollection.MetricReportDefinitionCollection",
			OdataID:      "/redfish/v1/TelemetryService/MetricReportDefinition",
			OdataType:    "#MetricReportDefinitionCollection.MetricReportDefinitionCollection",
			Description:  "MetricReportDefinition Collection view",
			Name:         "MetricReportDefinitionCollection",
			Members:      []dmtf.Link{},
			MembersCount: 0,
		}
		resp.Body = metricReportDefinitionCollection
		resp.StatusCode = http.StatusOK
		return resp
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}

// GetMetricReportCollection is a functioanlity to retrive all the available inventory
// resources from the added BMC's
func (e *ExternalInterface) GetMetricReportCollection(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	resp.Header = header
	data, err := e.DB.GetResource("MetricReportsCollection", req.URL, common.InMemory)
	if err != nil {
		// return empty collection response
		metricReportCollection := tlresp.Collection{
			OdataContext: "/redfish/v1/$metadata#MetricReportCollection.MetricReportCollection",
			OdataID:      "/redfish/v1/TelemetryService/MetricReport",
			OdataType:    "#MetricReportCollection.MetricReportCollection",
			Description:  "MetricReport Collection view",
			Name:         "MetricReportCollection",
			Members:      []dmtf.Link{},
			MembersCount: 0,
		}
		resp.Body = metricReportCollection
		resp.StatusCode = http.StatusOK
		return resp
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}

// GetTriggerCollection is a functioanlity to retrive all the available inventory
// resources from the added BMC's
func (e *ExternalInterface) GetTriggerCollection(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	resp.Header = header
	data, err := e.DB.GetResource("TriggersCollection", req.URL, common.InMemory)
	if err != nil {
		// return empty collection response
		triggersCollection := tlresp.Collection{
			OdataContext: "/redfish/v1/$metadata#TriggerCollection.TriggerCollection",
			OdataID:      "/redfish/v1/TelemetryService/Triggers",
			OdataType:    "#TriggerCollection.TriggerCollection",
			Description:  "Triggers Collection view",
			Name:         "Triggers",
			Members:      []dmtf.Link{},
			MembersCount: 0,
		}
		resp.Body = triggersCollection
		resp.StatusCode = http.StatusOK
		return resp
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}

// GetMetricReportDefinition ...
func (e *ExternalInterface) GetMetricReportDefinition(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	resp.Header = header
	data, gerr := e.DB.GetResource("MetricReportDefinitions", req.URL, common.InMemory)
	if gerr != nil {
		log.Warn("Unable to get MetricReportDefinition details : " + gerr.Error())
		errorMessage := gerr.Error()
		if errors.DBKeyNotFound == gerr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"MetricReportDefinition", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	return resp

}

// GetMetricReport is for to get metric report from southbound resource
func (e *ExternalInterface) GetMetricReport(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	resp.Header = header
	var getDeviceInfoRequest = tcommon.ResourceInfoRequest{
		URL:                 req.URL,
		ContactClient:       e.External.ContactClient,
		DevicePassword:      e.External.DevicePassword,
		GetPluginStatus:     e.External.GetPluginStatus,
		GetAllKeysFromTable: e.DB.GetAllKeysFromTable,
		GetPluginData:       e.External.GetPluginData,
		GetResource:         e.DB.GetResource,
		GenericSave:         e.External.GenericSave,
	}
	data, err := tcommon.GetResourceInfoFromDevice(getDeviceInfoRequest)
	if err != nil {
		log.Error(err.Error())
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"MetricReport", req.URL}, nil)
	}
	var resource map[string]interface{}
	json.Unmarshal(data, &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	return resp

}

// GetMetricDefinition ...
func (e *ExternalInterface) GetMetricDefinition(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	resp.Header = header
	data, gerr := e.DB.GetResource("MetricDefinitions", req.URL, common.InMemory)
	if gerr != nil {
		log.Warn("Unable to get MetricDefinition details : " + gerr.Error())
		errorMessage := gerr.Error()
		if errors.DBKeyNotFound == gerr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"MetricDefinition", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	return resp

}

// GetTrigger ...
func (e *ExternalInterface) GetTrigger(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	header["Allow"] = `GET, PATCH`
	resp.Header = header
	data, gerr := e.DB.GetResource("Triggers", req.URL, common.InMemory)
	if gerr != nil {
		log.Warn("Unable to get Triggers details `: " + gerr.Error())
		errorMessage := gerr.Error()
		if errors.DBKeyNotFound == gerr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"Triggers", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	return resp

}

// UpdateTrigger ...
func (e *ExternalInterface) UpdateTrigger(req *teleproto.TelemetryRequest) response.RPC {
	var resp response.RPC
	header["Allow"] = `GET, PATCH`
	resp.Header = header
	// Todo: code for update operation
	return resp
}
