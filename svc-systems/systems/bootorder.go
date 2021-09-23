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

//Package systems ...
package systems

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/scommon"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
)

// SetDefaultBootOrder defines the logic for setting the boot order to the default
func (p *PluginContact) SetDefaultBootOrder(systemID string) response.RPC {
	var resp response.RPC

	// spliting the uuid and system id
	requestData := strings.Split(systemID, ":")
	if len(requestData) <= 1 {
		errorMessage := "error: SystemUUID not found"
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"System", systemID}, nil)
	}
	uuid := requestData[0]
	target, gerr := smodel.GetTarget(uuid)
	if gerr != nil {
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, gerr.Error(), []interface{}{"System", uuid}, nil)
	}

	decryptedPasswordByte, err := p.DevicePassword(target.Password)
	if err != nil {
		// Frame the RPC response body and response Header below
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	target.Password = decryptedPasswordByte

	// Get the Plugin info
	plugin, gerr := smodel.GetPluginData(target.PluginID)
	if gerr != nil {
		errorMessage := "error while trying to get plugin details"
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}

	var contactRequest scommon.PluginContactRequest
	contactRequest.ContactClient = p.ContactClient
	contactRequest.Plugin = plugin

	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		contactRequest.HTTPMethodType = http.MethodPost
		contactRequest.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		contactRequest.OID = "/ODIM/v1/Sessions"
		_, token, getResponse, err := scommon.ContactPlugin(contactRequest, "error while creating session with the plugin: ")

		if err != nil {
			return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, err.Error(), nil, nil)
		}
		contactRequest.Token = token
	} else {
		contactRequest.BasicAuth = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	postRequest := make(map[string]interface{})
	postBody, _ := json.Marshal(postRequest)
	target.PostBody = postBody
	contactRequest.DeviceInfo = target
	contactRequest.OID = "/ODIM/v1/Systems/" + requestData[1] + "/Actions/ComputerSystem.SetDefaultBootOrder"
	contactRequest.HTTPMethodType = http.MethodPost

	body, _, getResponse, err := scommon.ContactPlugin(contactRequest, "error while setting the default bootorder of  the computer system: ")
	if err != nil {
		resp.StatusCode = getResponse.StatusCode
		json.Unmarshal(body, &resp.Body)
		resp.Header = map[string]string{"Content-type": "application/json; charset=utf-8"}
		return resp
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	err = json.Unmarshal(body, &resp.Body)
	if err != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
	}
	return resp
}

// ChangeBiosSettings defines the logic for change bios settings
func (p *PluginContact) ChangeBiosSettings(req *systemsproto.BiosSettingsRequest) response.RPC {
	var resp response.RPC

	// spliting the uuid and system id
	requestData := strings.Split(req.SystemID, ":")
	if len(requestData) <= 1 {
		errorMessage := "error: SystemUUID not found"
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"System", req.SystemID}, nil)
	}
	uuid := requestData[0]
	target, gerr := smodel.GetTarget(uuid)
	if gerr != nil {
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, gerr.Error(), []interface{}{"System", uuid}, nil)
	}

	var biosSetting BiosSetting

	// parsing the biosSetting
	err := json.Unmarshal(req.RequestBody, &biosSetting)
	if err != nil {
		errMsg := "unable to parse the BiosSetting request" + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, biosSetting)
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, nil)
		return response
	}

	decryptedPasswordByte, err := p.DevicePassword(target.Password)
	if err != nil {
		// Frame the RPC response body and response Header below
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	target.Password = decryptedPasswordByte
	// Get the Plugin info
	plugin, gerr := smodel.GetPluginData(target.PluginID)
	if gerr != nil {
		errorMessage := "error while trying to get plugin details"
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	var contactRequest scommon.PluginContactRequest
	contactRequest.ContactClient = p.ContactClient
	contactRequest.Plugin = plugin

	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		contactRequest.HTTPMethodType = http.MethodPost
		contactRequest.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		contactRequest.OID = "/ODIM/v1/Sessions"
		_, token, getResponse, err := scommon.ContactPlugin(contactRequest, "error while creating session with the plugin: ")

		if err != nil {
			return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, err.Error(), nil, nil)
		}
		contactRequest.Token = token
	} else {
		contactRequest.BasicAuth = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	target.PostBody = req.RequestBody

	contactRequest.HTTPMethodType = http.MethodPatch
	contactRequest.DeviceInfo = target
	contactRequest.OID = fmt.Sprintf("/ODIM/v1/Systems/%s/Bios/Settings", requestData[1])

	body, _, getResponse, err := scommon.ContactPlugin(contactRequest, "error while changing  bios settings: ")
	if err != nil {
		resp.StatusCode = getResponse.StatusCode
		json.Unmarshal(body, &resp.Body)
		resp.Header = map[string]string{"Content-type": "application/json; charset=utf-8"}
		return resp
	}
	resp.Header = map[string]string{
		"Content-type": "application/json; charset=utf-8",
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	err = json.Unmarshal(body, &resp.Body)
	if err != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
	}

	// Adding Settings URL to the DB to fetch data from device
	URL := fmt.Sprintf("/redfish/v1/Systems/%s/Bios/Settings", req.SystemID)
	smodel.AddSystemResetInfo(URL, "None")
	return resp
}

// ChangeBootOrderSettings defines the logic for change boot order settings
func (p *PluginContact) ChangeBootOrderSettings(req *systemsproto.BootOrderSettingsRequest) response.RPC {
	var resp response.RPC

	// spliting the uuid and system id
	requestData := strings.Split(req.SystemID, ":")
	if len(requestData) <= 1 {
		errorMessage := "error: SystemUUID not found"
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"System", req.SystemID}, nil)
	}

	var bootOrderSettings BootOrderSettings

	// parsing the bootOrderSettings
	err := json.Unmarshal(req.RequestBody, &bootOrderSettings)
	if err != nil {
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			errMsg := fmt.Sprintf("UnmarshalTypeError: Expected field type %v but got %v \n", ute.Type, ute.Value)
			log.Error(errMsg)
			index := strings.LastIndex(string(req.RequestBody[:ute.Offset]), ":")
			if index < 0 {
				index = 0
			}
			return common.GeneralError(http.StatusBadRequest, response.PropertyValueTypeError, errMsg, []interface{}{string(req.RequestBody[index+1 : ute.Offset]), ute.Field}, nil)
		}
		errMsg := "unable to parse the BootOrderSettings request" + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errMsg, []interface{}{}, nil)
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, bootOrderSettings)
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, nil)
		return response
	}

	uuid := requestData[0]
	target, gerr := smodel.GetTarget(uuid)
	if gerr != nil {
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, gerr.Error(), []interface{}{"System", uuid}, nil)
	}
	decryptedPasswordByte, err := p.DevicePassword(target.Password)
	if err != nil {
		// Frame the RPC response body and response Header below
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	target.Password = decryptedPasswordByte
	// Get the Plugin info
	plugin, gerr := smodel.GetPluginData(target.PluginID)
	if gerr != nil {
		errorMessage := "error while trying to get plugin details"
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	var contactRequest scommon.PluginContactRequest
	contactRequest.ContactClient = p.ContactClient
	contactRequest.Plugin = plugin

	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		contactRequest.HTTPMethodType = http.MethodPost
		contactRequest.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		contactRequest.OID = "/ODIM/v1/Sessions"
		_, token, getResponse, err := scommon.ContactPlugin(contactRequest, "error while creating session with the plugin: ")

		if err != nil {
			return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, err.Error(), nil, nil)
		}
		contactRequest.Token = token
	} else {
		contactRequest.BasicAuth = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}

	target.PostBody = req.RequestBody

	contactRequest.HTTPMethodType = http.MethodPatch
	contactRequest.DeviceInfo = target
	contactRequest.OID = fmt.Sprintf("/ODIM/v1/Systems/%s", requestData[1])

	body, _, getResponse, err := scommon.ContactPlugin(contactRequest, "error while changing boot order settings: ")
	if err != nil {
		resp.StatusCode = getResponse.StatusCode
		json.Unmarshal(body, &resp.Body)
		resp.Header = map[string]string{"Content-type": "application/json; charset=utf-8"}
		return resp
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	err = json.Unmarshal(body, &resp.Body)
	if err != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
	}
	smodel.AddSystemResetInfo("/redfish/v1/Systems/"+req.SystemID, "On")
	return resp
}
