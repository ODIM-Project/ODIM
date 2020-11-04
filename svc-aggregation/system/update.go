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

package system

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
)

// UpdateAggregationSource defines the  interface for updation of  added Aggregation Source
func (e *ExternalInterface) UpdateAggregationSource(req *aggregatorproto.AggregatorRequest) response.RPC {
	// validate the aggregation source if it's  present in odim
	var resp response.RPC
	aggregationSource, dbErr := agmodel.GetAggregationSourceInfo(req.URL)
	if dbErr != nil {
		log.Printf("error getting  AggregationSource : %v", dbErr)
		errorMessage := dbErr.Error()
		if errors.DBKeyNotFound == dbErr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"AggregationSource", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	// parse the request
	var updateRequest map[string]interface{}
	err := json.Unmarshal(req.RequestBody, &updateRequest)
	if err != nil {
		errMsg := "unable to parse the add request" + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	if len(updateRequest) <= 0 {
		param := "HostName UserName Password "
		errMsg := "error:  field " + param + " Missing"
		log.Printf(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{param}, nil)
	}
	var param string
	for key, value := range updateRequest {
		if value.(string) == "" {
			param = param + key + " "
		}
	}
	if param != "" {
		errMsg := "error:  field " + param + " Missing"
		log.Printf(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{param}, nil)
	}
	if _, ok := updateRequest["UserName"]; !ok {
		updateRequest["UserName"] = aggregationSource.UserName
	}
	var hostNameUpdated bool
	if _, ok := updateRequest["HostName"]; !ok {
		updateRequest["HostName"] = aggregationSource.HostName

	} else {
		err := validateManagerAddress(updateRequest["HostName"].(string))
		if err != nil {
			log.Printf(err.Error())
			return common.GeneralError(http.StatusBadRequest, response.PropertyValueFormatError, err.Error(), []interface{}{updateRequest["HostName"].(string), "HostName"}, nil)

		}
		hostNameUpdated = true
	}
	if _, ok := updateRequest["Password"]; !ok {
		decryptedPasswordByte, err := e.DecryptPassword(aggregationSource.Password)
		if err != nil {
			errMsg := "error while trying to decrypt device password: " + err.Error()
			log.Printf(errMsg)
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
		}
		updateRequest["Password"] = decryptedPasswordByte
	} else {
		bytePassword := []byte(updateRequest["Password"].(string))
		updateRequest["Password"] = bytePassword
	}
	var data = strings.Split(req.URL, "/redfish/v1/AggregationService/AggregationSources/")
	links := aggregationSource.Links.(map[string]interface{})
	if _, ok := links["ConnectionMethod"]; ok {
		resp = e.updateAggregationSourceWithConnectionMethod(req.URL, links["ConnectionMethod"].(map[string]interface{}), updateRequest, hostNameUpdated)
	} else {
		oem := links["Oem"].(map[string]interface{})
		pluginID := oem["PluginID"].(string)
		if _, ok := oem["PluginType"]; ok {
			resp = e.updateManagerAggregationSource(data[1], pluginID, updateRequest, hostNameUpdated)
		} else {
			resp = e.updateBMCAggregationSource(data[1], pluginID, updateRequest, hostNameUpdated)
		}
	}
	if resp.StatusMessage != "" {
		return resp
	}
	// Update the aggregation source info
	aggregationSource.HostName = updateRequest["HostName"].(string)
	aggregationSource.UserName = updateRequest["UserName"].(string)
	aggregationSource.Password = updateRequest["Password"].([]byte)

	dbErr = agmodel.UpdateAggregtionSource(aggregationSource, req.URL)
	if dbErr != nil {
		errMsg := "error while trying to update aggregation source info: " + dbErr.Error()
		fmt.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}

	commonResponse := response.Response{
		OdataType:    "#AggregationSource.v1_0_0.AggregationSource",
		OdataID:      req.URL,
		OdataContext: "/redfish/v1/$metadata#AggregationSource.AggregationSource",
		ID:           data[1],
		Name:         "Aggregation Source",
	}
	resp.Header = map[string]string{
		"Allow":             `"GET","PATCH","DELETE"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	commonResponse.CreateGenericResponse(response.Success)
	commonResponse.Message = ""
	commonResponse.MessageID = ""
	commonResponse.Severity = ""
	resp.Body = agresponse.AggregationSourceResponse{
		Response: commonResponse,
		HostName: updateRequest["HostName"].(string),
		UserName: updateRequest["UserName"].(string),
		Links:    aggregationSource.Links,
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}

func (e *ExternalInterface) updateAggregationSourceWithConnectionMethod(url string, connectionMethodLink, updateRequest map[string]interface{}, hostNameUpdated bool) response.RPC {
	connectionMethodOdataID := connectionMethodLink["@odata.id"].(string)
	connectionMethod, err := e.GetConnectionMethod(connectionMethodOdataID)
	if err != nil {
		log.Printf("error getting  connectionmethod : %v", err)
		errorMessage := err.Error()
		if errors.DBKeyNotFound == err.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"ConnectionMethod", connectionMethodOdataID}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	cmVariants := getConnectionMethodVariants(connectionMethod.ConnectionMethodVariant)
	var data = strings.Split(url, "/redfish/v1/AggregationService/AggregationSources/")
	uuid := url[strings.LastIndexByte(url, '/')+1:]
	target, terr := agmodel.GetTarget(uuid)
	if terr != nil || target == nil {
		return e.updateManagerAggregationSource(data[1], cmVariants.PluginID, updateRequest, hostNameUpdated)
	}
	return e.updateBMCAggregationSource(data[1], cmVariants.PluginID, updateRequest, hostNameUpdated)
}

func (e *ExternalInterface) updateManagerAggregationSource(aggregationSourceID, pluginID string, updateRequest map[string]interface{}, hostNameUpdated bool) response.RPC {
	plugin, errs := agmodel.GetPluginData(pluginID)
	if errs != nil {
		errMsg := errs.Error()
		log.Printf(errMsg)
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"plugin", pluginID}, nil)
	}
	ipData := strings.Split(updateRequest["HostName"].(string), ":")
	plugin.IP = ipData[0]
	plugin.Port = ipData[1]
	plugin.Username = updateRequest["UserName"].(string)
	plugin.Password = updateRequest["Password"].([]byte)
	var pluginContactRequest getResourceRequest
	pluginContactRequest.ContactClient = e.ContactClient
	pluginContactRequest.GetPluginStatus = e.GetPluginStatus

	pluginContactRequest.Plugin = plugin
	pluginContactRequest.StatusPoll = true
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		pluginContactRequest.HTTPMethodType = http.MethodPost
		pluginContactRequest.DeviceInfo = map[string]interface{}{
			"Username": plugin.Username,
			"Password": string(plugin.Password),
		}
		pluginContactRequest.OID = "/ODIM/v1/Sessions"
		_, token, getResponse, err := contactPlugin(pluginContactRequest, "error while creating the session: ")
		if err != nil {
			errMsg := err.Error()
			log.Println(errMsg)
			return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, nil)
		}
		pluginContactRequest.Token = token
	} else {
		pluginContactRequest.LoginCredentials = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
	}

	// Verfiying the plugin Status
	pluginContactRequest.HTTPMethodType = http.MethodGet
	pluginContactRequest.OID = "/ODIM/v1/Status"
	body, _, getResponse, err := contactPlugin(pluginContactRequest, "error while getting the details "+pluginContactRequest.OID+": ")
	if err != nil {
		errMsg := err.Error()
		log.Println(errMsg)
		return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, nil)
	}
	var managerUUID = plugin.ManagerUUID
	if hostNameUpdated {
		var managersMap map[string]interface{}
		// Getting all managers info from plugin
		pluginContactRequest.OID = "/ODIM/v1/Managers"
		body, _, getResponse, err = contactPlugin(pluginContactRequest, "error while getting the details "+pluginContactRequest.OID+": ")
		if err != nil {
			errMsg := err.Error()
			log.Println(errMsg)
			return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, nil)
		}
		//  Extract all managers info and loop  over each members
		err = json.Unmarshal([]byte(body), &managersMap)
		if err != nil {
			errMsg := "unable to parse the managers resposne" + err.Error()
			log.Println(errMsg)
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
		}
		managerMembers := managersMap["Members"]

		// Getting the indivitual managers response
		for _, object := range managerMembers.([]interface{}) {
			pluginContactRequest.OID = object.(map[string]interface{})["@odata.id"].(string)
			body, _, getResponse, err := contactPlugin(pluginContactRequest, "error while getting the details "+pluginContactRequest.OID+": ")
			if err != nil {
				errMsg := err.Error()
				log.Println(errMsg)
				return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, nil)
			}
			var managerData map[string]interface{}
			err = json.Unmarshal([]byte(body), &managerData)
			if err != nil {
				errMsg := "unable to parse the managers resposne" + err.Error()
				log.Println(errMsg)
				return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
			}
			if uuid, ok := managerData["UUID"]; ok {
				managerUUID = uuid.(string)
			}
		}
		if managerUUID != plugin.ManagerUUID {
			errMsg := "error: uuid of the added managers is not matching with given HostName"
			log.Println(errMsg)
			return common.GeneralError(http.StatusBadRequest, response.ResourceInUse, errMsg, nil, nil)
		}
	}

	// encrypt plugin password
	ciphertext, err := e.EncryptPassword(plugin.Password)
	if err != nil {
		errMsg := "error: encryption failed: " + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}

	plugin.Password = ciphertext
	plugin.ManagerUUID = managerUUID
	updateRequest["Password"] = ciphertext
	dbErr := agmodel.UpdatePluginData(plugin, pluginID)
	if dbErr != nil {
		errMsg := "error while trying to update plugin info: " + dbErr.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}

	return response.RPC{
		StatusCode: http.StatusOK,
	}
}

func (e *ExternalInterface) updateBMCAggregationSource(aggregationSourceID, pluginID string, updateRequest map[string]interface{}, hostNameUpdated bool) response.RPC {
	// Get the plugin  from db
	plugin, errs := agmodel.GetPluginData(pluginID)
	if errs != nil {
		errMsg := errs.Error()
		log.Printf(errMsg)
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"plugin", pluginID}, nil)
	}
	var pluginContactRequest getResourceRequest
	pluginContactRequest.ContactClient = e.ContactClient
	pluginContactRequest.GetPluginStatus = e.GetPluginStatus
	pluginContactRequest.Plugin = plugin
	pluginContactRequest.StatusPoll = true

	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		var err error
		pluginContactRequest.HTTPMethodType = http.MethodPost
		pluginContactRequest.DeviceInfo = map[string]interface{}{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
		pluginContactRequest.OID = "/ODIM/v1/Sessions"
		_, token, getResponse, err := contactPlugin(pluginContactRequest, "error while logging in to plugin: ")
		if err != nil {
			errMsg := err.Error()
			log.Println(errMsg)
			return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, nil)
		}
		pluginContactRequest.Token = token
	} else {
		pluginContactRequest.LoginCredentials = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	// validate the device credentials
	var saveSystem = agmodel.SaveSystem{
		ManagerAddress: updateRequest["HostName"].(string),
		UserName:       updateRequest["UserName"].(string),
		Password:       updateRequest["Password"].([]byte),
	}
	pluginContactRequest.DeviceInfo = saveSystem
	pluginContactRequest.OID = "/ODIM/v1/validate"
	pluginContactRequest.HTTPMethodType = http.MethodPost

	body, _, getResponse, err := contactPlugin(pluginContactRequest, "error while trying to authenticate the compute server: ")
	if err != nil {
		errMsg := err.Error()
		log.Println(errMsg)
		return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, nil)
	}

	var commonError errors.CommonError
	err = json.Unmarshal(body, &commonError)
	if err != nil {
		errMsg := err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	if hostNameUpdated {
		// Get All systems
		pluginContactRequest.OID = "/redfish/v1/Systems"
		pluginContactRequest.HTTPMethodType = http.MethodGet
		body, _, getResponse, err = contactPlugin(pluginContactRequest, "error while trying to get system collection details: ")
		if err != nil {
			errMsg := err.Error()
			log.Println(errMsg)
			return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, nil)
		}
		var systemsMap map[string]interface{}
		err = json.Unmarshal([]byte(body), &systemsMap)
		if err != nil {
			errMsg := "error while trying unmarshal systems collection: " + err.Error()
			log.Println(errMsg)
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
		}
		systemMembers := systemsMap["Members"]
		for _, object := range systemMembers.([]interface{}) {
			oDataID := object.(map[string]interface{})["@odata.id"].(string)
			pluginContactRequest.OID = oDataID
			body, _, getResponse, err = contactPlugin(pluginContactRequest, "error while trying to get system details: ")
			if err != nil {
				errMsg := err.Error()
				log.Println(errMsg)
				return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, nil)
			}
			var computeSystem map[string]interface{}
			err = json.Unmarshal(body, &computeSystem)
			if err != nil {
				errMsg := "error while trying unmarshal computer system: " + err.Error()

				log.Println(errMsg)
				return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
			}
			computeSystemID := computeSystem["Id"].(string)
			computeSystemUUID := computeSystem["UUID"].(string)
			oidKey := keyFormation(oDataID, computeSystemID, aggregationSourceID)
			log.Println("Computer SystemUUID", computeSystemUUID)
			indexList, err := agmodel.GetString("UUID", computeSystemUUID)
			if err != nil {
				errMsg := "error while trying get computer system index: " + err.Error()
				log.Println(errMsg)
				return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
			}
			log.Println("Index List", indexList)
			if len(indexList) <= 0 {
				errMsg := "error: uuid of the added bmc is not matching with given HostName"
				log.Println(errMsg)
				return common.GeneralError(http.StatusBadRequest, response.ResourceInUse, errMsg, nil, nil)
			}
			var isPresent bool
			for _, systemID := range indexList {
				if systemID == oidKey {
					isPresent = true
				}
			}
			if !isPresent {
				errMsg := "error: uuid of the added bmc is not matching with given HostName"
				log.Println(errMsg)
				return common.GeneralError(http.StatusBadRequest, response.ResourceInUse, errMsg, nil, nil)
			}

		}
	}
	// update the system
	saveSystem.PluginID = pluginID
	saveSystem.DeviceUUID = aggregationSourceID
	// encrypt the device password
	ciphertext, err := e.EncryptPassword([]byte(saveSystem.Password))
	if err != nil {
		errMsg := "error while trying to encrypt: " + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	saveSystem.Password = ciphertext
	updateRequest["Password"] = ciphertext
	dbErr := agmodel.UpdateSystemData(saveSystem, aggregationSourceID)
	if dbErr != nil {
		errMsg := "error while trying to update system info: " + dbErr.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}

	return response.RPC{
		StatusCode: http.StatusOK,
	}
}

func validateManagerAddress(managerAddress string) error {
	// if the manager address is of the form <IP/FQDN>:<port>
	// will split address to obtain only IP/FQDN. If obtained
	// value is empty, then will use the actual address passed
	addr, _, _ := net.SplitHostPort(managerAddress)
	if addr == "" {
		addr = managerAddress
	}
	if _, err := net.ResolveIPAddr("ip", addr); err != nil {
		return fmt.Errorf("error: failed to resolve ManagerAddress: %v", err)
	}
	return nil
}
