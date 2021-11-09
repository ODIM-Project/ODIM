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

//Package update ...
package update

// ---------------------------------------------------------------------------------------
// IMPORT Section
// ---------------------------------------------------------------------------------------
import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-update/ucommon"
	"github.com/ODIM-Project/ODIM/svc-update/uresponse"
)

// GetUpdateService defines the functionality for knowing whether
// the update service is enabled or not
//
// As return parameters RPC response, which contains status code, message, headers and data,
// error will be passed back.
func (e *ExternalInterface) GetUpdateService() response.RPC {
	commonResponse := response.Response{
		OdataType:    common.UpdateServiceType,
		OdataID:      "/redfish/v1/UpdateService",
		OdataContext: "/redfish/v1/$metadata#UpdateService.UpdateService",
		ID:           "UpdateService",
		Name:         "Update Service",
	}
	var resp response.RPC

	isServiceEnabled := false
	serviceState := "Disabled"
	//Checks if UpdateService is enabled and sets the variable isServiceEnabled to true add servicState to enabled
	for _, service := range config.Data.EnabledServices {
		if service == "UpdateService" {
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
		"Link": "	</redfish/v1/SchemaStore/en/UpdateService.json>; rel=describedby",
		"Transfer-Encoding": "chunked",
		"X-Frame-Options":   "sameorigin",
	}

	commonResponse.CreateGenericResponse(resp.StatusMessage)
	commonResponse.Message = ""
	commonResponse.MessageID = ""
	commonResponse.Severity = ""
	resp.Body = uresponse.UpdateService{
		Response: commonResponse,
		//TODO: Yet to implement UpdateService state and health
		Status: uresponse.Status{
			State:        serviceState,
			Health:       "OK",
			HealthRollup: "OK",
		},
		ServiceEnabled: isServiceEnabled,
		FirmwareInventory: uresponse.FirmwareInventory{
			OdataID: "/redfish/v1/UpdateService/FirmwareInventory",
		},
		SoftwareInventory: uresponse.SoftwareInventory{
			OdataID: "/redfish/v1/UpdateService/SoftwareInventory",
		},
		Actions: uresponse.Actions{
			UpdateServiceSimpleUpdate: uresponse.UpdateServiceSimpleUpdate{
				Target: "/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate",
				RedfishOperationApplyTimeSupport: uresponse.RedfishOperationApplyTimeSupport{
					OdataType:       common.SettingsType,
					SupportedValues: []string{"OnStartUpdateRequest"},
				},
			},
			UpdateServiceStartUpdate: uresponse.UpdateServiceStartUpdate{
				Target: "/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate",
			},
		},
	}

	return resp

}

// GetAllFirmwareInventory is a functioanlity to retrive all the available inventory
// resources from the added BMC's
func (e *ExternalInterface) GetAllFirmwareInventory(req *updateproto.UpdateRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	firmwareCollection := uresponse.Collection{
		OdataContext: "/redfish/v1/$metadata#FirmwareInventoryCollection.FirmwareCollection",
		OdataID:      "/redfish/v1/UpdateService/FirmwareInventory",
		OdataType:    "#FirmwareInventoryCollection.FirmwareInventoryCollection",
		Description:  "FirmwareInventory view",
		Name:         "FirmwareInventory",
	}

	members := []dmtf.Link{}
	firmwareCollectionKeysArray, err := e.DB.GetAllKeysFromTable("FirmwareInventory", common.InMemory)
	if err != nil || len(firmwareCollectionKeysArray) == 0 {
		log.Warn("odimra doesnt have servers")
	}

	for _, key := range firmwareCollectionKeysArray {
		members = append(members, dmtf.Link{Oid: key})
	}
	firmwareCollection.Members = members
	firmwareCollection.MembersCount = len(members)
	resp.Body = firmwareCollection
	resp.StatusCode = http.StatusOK
	return resp
}

// GetFirmwareInventory is used to fetch resource data. The function is supposed to be used as part of RPC
// For getting firmware inventory resource information,  parameters need to be passed Request .
// Request holds the  Uuid and Url ,
// Url will be parsed from that search key will created
// There will be two return values for the fuction. One is the RPC response, which contains the
// status code, status message, headers and body and the second value is error.
func (e *ExternalInterface) GetFirmwareInventory(req *updateproto.UpdateRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}

	requestData := strings.Split(req.ResourceID, ":")
	if len(requestData) <= 1 {
		errorMessage := "error: SystemUUID not found"
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"FirmwareInventory", req.ResourceID}, nil)
	}
	data, gerr := e.DB.GetResource("FirmwareInventory", req.URL, common.InMemory)
	if gerr != nil {
		log.Warn("Unable to get firmware inventory details : " + gerr.Error())
		errorMessage := gerr.Error()
		if errors.DBKeyNotFound == gerr.ErrNo() {
			var getDeviceInfoRequest = ucommon.ResourceInfoRequest{
				URL:            req.URL,
				UUID:           requestData[0],
				SystemID:       requestData[1],
				ContactClient:  e.External.ContactClient,
				DevicePassword: e.External.DevicePassword,
			}
			var err error
			i := ucommon.CommonInterface{
				GetTarget:     e.External.GetTarget,
				GetPluginData: e.External.GetPluginData,
				ContactPlugin: e.External.ContactPlugin,
			}
			if data, err = i.GetResourceInfoFromDevice(getDeviceInfoRequest); err != nil {
				return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"FirmwareInventory", req.URL}, nil)
			}
		} else {
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		}
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	return resp

}

// GetAllSoftwareInventory is a functioanlity to retrive all the available inventory
// resources from the added BMC's
func (e *ExternalInterface) GetAllSoftwareInventory(req *updateproto.UpdateRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	softwareCollection := uresponse.Collection{
		OdataContext: "/redfish/v1/$metadata#SoftwareInventoryCollection.SoftwareCollection",
		OdataID:      "/redfish/v1/UpdateService/SoftwareInventory",
		OdataType:    "#SoftwareInventoryCollection.SoftwareInventoryCollection",
		Description:  "SoftwareInventory view",
		Name:         "SoftwareInventory",
	}

	members := []dmtf.Link{}
	softwareCollectionKeysArray, err := e.DB.GetAllKeysFromTable("SoftwareInventory", common.InMemory)
	if err != nil || len(softwareCollectionKeysArray) == 0 {
		log.Warn("odimra doesnt have servers")
	}

	for _, key := range softwareCollectionKeysArray {
		members = append(members, dmtf.Link{Oid: key})
	}
	softwareCollection.Members = members
	softwareCollection.MembersCount = len(members)
	resp.Body = softwareCollection
	resp.StatusCode = http.StatusOK
	return resp
}

// GetSoftwareInventory is used to fetch resource data. The function is supposed to be used as part of RPC
// For getting software inventory resource information,  parameters need to be passed Request .
// Request holds the  Uuid and Url ,
// Url will be parsed from that search key will created
// There will be two return values for the fuction. One is the RPC response, which contains the
// status code, status message, headers and body and the second value is error.
func (e *ExternalInterface) GetSoftwareInventory(req *updateproto.UpdateRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}

	requestData := strings.Split(req.ResourceID, ":")
	if len(requestData) <= 1 {
		errorMessage := "error: SystemUUID not found"
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"SoftwareInventory", req.ResourceID}, nil)
	}
	data, gerr := e.DB.GetResource("SoftwareInventory", req.URL, common.InMemory)
	if gerr != nil {
		log.Warn("Unable to get software inventory details : " + gerr.Error())
		errorMessage := gerr.Error()
		if errors.DBKeyNotFound == gerr.ErrNo() {
			var getDeviceInfoRequest = ucommon.ResourceInfoRequest{
				URL:            req.URL,
				UUID:           requestData[0],
				SystemID:       requestData[1],
				ContactClient:  e.External.ContactClient,
				DevicePassword: e.External.DevicePassword,
			}
			var err error
			i := ucommon.CommonInterface{
				GetTarget:     e.External.GetTarget,
				GetPluginData: e.External.GetPluginData,
				ContactPlugin: e.External.ContactPlugin,
			}
			if data, err = i.GetResourceInfoFromDevice(getDeviceInfoRequest); err != nil {
				return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"SoftwareInventory", req.URL}, nil)
			}
		} else {
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		}
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	return resp

}
