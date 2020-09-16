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
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-update/ucommon"
	"github.com/ODIM-Project/ODIM/svc-update/umodel"
	"log"
	"net/http"
	"strings"
)

// StartUpdate function handler for on start update process
func (e *ExternalInterface) StartUpdate(req *updateproto.UpdateRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	// Read all the requests from database
	targetList, err := umodel.GetAllKeysFromTable("SimpleUpdate", common.OnDisk)
	if err != nil {
		errMsg := "error: unable to read SimpleUpdate requests from database: " + err.Error()
		log.Println(errMsg)
	}
	if len(targetList) == 0 {
		resp.StatusCode = http.StatusOK
		resp.StatusMessage = response.Success
		var args response.Args
		args = response.Args{
			Code:    resp.StatusMessage,
			Message: "Request completed successfully",
		}
		resp.Body = args.CreateGenericErrorResponse()
		return resp
	}
	for _, target := range targetList {
		data, gerr := e.DB.GetResource("SimpleUpdate", target, common.OnDisk)
		if gerr != nil {
			errMsg := "error: unable to retrive the start update request" + gerr.Error()
			log.Println(errMsg)
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
		}
		resp = e.startRequest(target, data, resp)
	}
	return resp
}

func (e *ExternalInterface) startRequest(uuid, data string, resp response.RPC) response.RPC {
	updateRequestBody := strings.Replace(data, uuid+":", "", -1)
	//replacing the reruest url with south bound translation URL
	for key, value := range config.Data.URLTranslation.SouthBoundURL {
		updateRequestBody = strings.Replace(updateRequestBody, key, value, -1)
	}
	target, gerr := e.External.GetTarget(uuid)
	if gerr != nil {
		return common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, gerr.Error(), []interface{}{"System", uuid}, nil)
	}

	decryptedPasswordByte, passwdErr := e.External.DevicePassword(target.Password)
	if passwdErr != nil {
		// Frame the RPC response body and response Header below
		errorMessage := "error while trying to decrypt device password: " + passwdErr.Error()
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
	contactRequest.OID = "/ODIM/v1/UpdateService/Actions/UpdateService.StartUpdate"
	contactRequest.HTTPMethodType = http.MethodPost
	body, _, getResponse, contactErr := e.External.ContactPlugin(contactRequest, "error while performing simple update action: ")
	if contactErr != nil {
		resp.StatusCode = getResponse.StatusCode
		json.Unmarshal(body, &resp.Body)
		resp.Header = map[string]string{"Content-type": "application/json; charset=utf-8"}
		return resp
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	respErr := json.Unmarshal(body, &resp.Body)
	if respErr != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, respErr.Error(), nil, nil)
	}
	return resp
}
