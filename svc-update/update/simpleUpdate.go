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
//
import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-update/ucommon"
	"github.com/ODIM-Project/ODIM/svc-update/umodel"
)

// SimpleUpdate function handler for simpe update process
func (e *ExternalInterface) SimpleUpdate(req *updateproto.UpdateRequest) response.RPC {
	var resp response.RPC
	var updateRequest UpdateRequestBody
	err := json.Unmarshal(req.RequestBody, &updateRequest)
	if err != nil {
		errMsg := "error: unable to parse the simple update request" + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	if len(updateRequest.Targets) == 0 {
		errMsg := "error: 'Targets' parameter cannot be empty"
		log.Println(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"Targets"}, nil)
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, updateRequest)
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Println(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, nil)
		return response
	}
	targetList := make(map[string][]string)
	var applyTime []string
	if updateRequest.RedfishOperationApplyTimeSupport != nil && updateRequest.RedfishOperationApplyTimeSupport.SupportedValues != nil {
		applyTime = updateRequest.RedfishOperationApplyTimeSupport.SupportedValues
	}
	targetList, err = sortTargetList(updateRequest.Targets)
	if err != nil {
		errorMessage := "error: SystemUUID not found"
		return common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, errorMessage, []interface{}{"System", fmt.Sprintf("%v", updateRequest.Targets)}, nil)
	}
	if len(targetList) > 1 {
		errMsg := "error: 'Targets' parameter cannot have more than one BMC"
		log.Println(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"Targets"}, nil)
	}
	for id, target := range targetList {
		updateRequest.Targets = target
		marshalBody, err := json.Marshal(updateRequest)
		if err != nil {
			errMsg := "error: unable to parse the simple update request" + err.Error()
			log.Println(errMsg)
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
		}
		updateRequestBody := strings.Replace(string(marshalBody), id+":", "", -1)
		//replacing the reruest url with south bound translation URL
		for key, value := range config.Data.URLTranslation.SouthBoundURL {
			updateRequestBody = strings.Replace(updateRequestBody, key, value, -1)
		}
		resp = e.sendRequest(id, updateRequestBody, applyTime, resp)
	}
	return resp
}

func (e *ExternalInterface) sendRequest(uuid, updateRequestBody string, applyTime []string, resp response.RPC) response.RPC {
	target, gerr := e.External.GetTarget(uuid)
	if gerr != nil {
		return common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, gerr.Error(), []interface{}{"System", uuid}, nil)
	}
	if len(applyTime) != 0 {
		err := umodel.GenericSave([]byte(updateRequestBody), "SimpleUpdate", uuid)
		if err != nil {
			errMsg := "error: unable to save the simple update request" + err.Error()
			log.Println(errMsg)
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
		}
	}

	decryptedPasswordByte, err := e.External.DevicePassword(target.Password)
	if err != nil {
		// Frame the RPC response body and response Header below
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	target.Password = decryptedPasswordByte

	// Get the Plugin info
	plugin, gerr := e.External.GetPluginData(target.PluginID)
	if gerr != nil {
		errorMessage := "error while trying to get plugin details"
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	var contactRequest ucommon.PluginContactRequest
	contactRequest.ContactClient = e.External.ContactClient
	contactRequest.Plugin = plugin

	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		contactRequest.HTTPMethodType = http.MethodPost
		contactRequest.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		contactRequest.OID = "/ODIM/v1/Sessions"
		_, token, getResponse, err := e.External.ContactPlugin(contactRequest, "error while creating session with the plugin: ")

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

	target.PostBody = []byte(updateRequestBody)
	contactRequest.DeviceInfo = target
	contactRequest.OID = "/ODIM/v1/UpdateService/Actions/UpdateService.SimpleUpdate"
	contactRequest.HTTPMethodType = http.MethodPost
	body, _, getResponse, err := e.External.ContactPlugin(contactRequest, "error while performing simple update action: ")
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

func sortTargetList(Targets []string) (map[string][]string, error) {
	returnList := make(map[string][]string)
	for _, individualTarget := range Targets {
		// spliting the uuid and system id
		requestData := strings.Split(individualTarget, "/")
		var requestTarget []string
		for _, data := range requestData {
			if strings.Contains(data, ":") {
				requestTarget = strings.Split(data, ":")
			}
		}
		if len(requestTarget) != 2 || requestTarget[1] == "" {
			errorMessage := "error: SystemUUID not found"
			return returnList, errors.New(errorMessage)
		}
		uuid := requestTarget[0]
		returnList[uuid] = append(returnList[uuid], individualTarget)
	}
	return returnList, nil
}
