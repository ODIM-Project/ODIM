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
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

// DeleteRequest is payload of deleting resource
type DeleteRequest struct {
	OdataContext string       `json:"@odata.context"`
	OdataID      string       `json:"@odata.id"`
	Odatatype    string       `json:"@odata.type"`
	ID           string       `json:"Id"`
	Name         string       `json:"Name"`
	OEM          interface{}  `json:"Oem"`
	Parameters   []Parameters `json:"Parameters"`
}

// Parameters is struct to have the delete request parameters
type Parameters struct {
	Name string `json:"Name"`
}

// DeleteCompute is the handler for Deleting system
// at first chech is token authorized and has privileges
// if its verified then delete the compute system and send success response
func (e *ExternalInterface) DeleteCompute(req *aggregatorproto.AggregatorRequest) response.RPC {
	var deleteRequest DeleteRequest
	if err := json.Unmarshal(req.RequestBody, &deleteRequest); err != nil {
		errMsg := "error while trying to delete compute system: " + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errMsg, nil, nil)
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, deleteRequest)
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Println(errorMessage)
		resp := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, nil)
		return resp
	}

	if len(deleteRequest.Parameters) == 1 && deleteRequest.Parameters[0].Name != "" {
		key := deleteRequest.Parameters[0].Name
		// check if given resource is a manager
		if strings.Contains(key, "Managers") {
			return e.deletePlugin(key)
		}
		index := strings.LastIndexAny(key, "/")
		if index > 0 {
			return e.deleteCompute(key, index)
		}
	}

	errMsg := "error while trying to delete compute system: Invalid request"
	log.Println(errMsg)
	return common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errMsg, nil, nil)

}

// deleteplugin removes the given plugin
func (e *ExternalInterface) deletePlugin(oid string) response.RPC {
	var resp response.RPC
	// Get Manager Info
	data, derr := agmodel.GetResource("Managers", oid)
	if derr != nil {
		errMsg := "error while getting Managers data: " + derr.Error()
		log.Println(errMsg)
		if errors.DBKeyNotFound == derr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"Managers", oid}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	var pluginID = resource["Name"].(string)
	plugin, errs := agmodel.GetPluginData(pluginID)
	if errs != nil {
		errMsg := "error while getting plugin data: " + errs.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"Plugin", pluginID}, nil)
	}

	systems, dberr := agmodel.GetAllSystems()
	if dberr != nil {
		errMsg := derr.Error()
		log.Println(errMsg)
		if errors.DBKeyNotFound == derr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"Systems", "everything"}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	// verifying if any device is mapped to plugin
	var systemCnt = 0
	for i := 0; i < len(systems); i++ {
		if systems[i].PluginID == pluginID {
			systemCnt++
		}
	}
	if systemCnt > 0 {
		errMsg := fmt.Sprintf("error: plugin %v can't be removed since it managing some of the devices", pluginID)
		log.Println(errMsg)
		return common.GeneralError(http.StatusNotAcceptable, response.ResourceCannotBeDeleted, errMsg, nil, nil)
	}

	// verifying if plugin is up
	var pluginContactRequest getResourceRequest

	pluginContactRequest.ContactClient = e.ContactClient
	pluginContactRequest.Plugin = plugin
	pluginContactRequest.StatusPoll = false
	pluginContactRequest.HTTPMethodType = http.MethodGet
	pluginContactRequest.LoginCredentials = map[string]string{
		"UserName": plugin.Username,
		"Password": string(plugin.Password),
	}
	pluginContactRequest.OID = "/ODIM/v1/Status"
	_, _, _, err := contactPlugin(pluginContactRequest, "error while getting the details "+pluginContactRequest.OID+": ")
	if err == nil { // no err means plugin is still up, so we can't remove it
		errMsg := "error: plugin is still up, so it cannot be removed."
		log.Println(errMsg)
		return common.GeneralError(http.StatusNotAcceptable, response.ResourceCannotBeDeleted, errMsg, nil, nil)
	}

	// deleting the manager info
	dberr = agmodel.DeleteManagersData(oid)
	if dberr != nil {
		errMsg := derr.Error()
		log.Println(errMsg)
		if errors.DBKeyNotFound == derr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"Managers", oid}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	// deleting the plugin if  zero devices are managed
	dberr = agmodel.DeletePluginData(pluginID)
	if dberr != nil {
		errMsg := derr.Error()
		log.Println(errMsg)
		if errors.DBKeyNotFound == derr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"Plugin", pluginID}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	e.EventNotification(oid, "ResourceRemoved", "ManagerCollection")
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Transfer-Encoding": "chunked",
		"Content-type":      "application/json; charset=utf-8",
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.ResourceRemoved

	args := response.Args{
		Code:    resp.StatusMessage,
		Message: "Request completed successfully",
	}
	resp.Body = args.CreateGenericErrorResponse()
	return resp
}

func (e *ExternalInterface) deleteCompute(key string, index int) response.RPC {
	var resp response.RPC
	// check whether the any system operation is under progress
	systemOperation, dbErr := agmodel.GetSystemOperationInfo(strings.TrimSuffix(key, "/"))
	if dbErr != nil && errors.DBKeyNotFound != dbErr.ErrNo() {
		log.Println(" Delete operation for system  ", key, " can't be processed ", dbErr.Error())
		errMsg := "error while trying to delete compute system: " + dbErr.Error()
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	if systemOperation.Operation != "" {
		log.Println("Delete operation or system  ", key, " can't be processed,", systemOperation.Operation, " operation  is under progress")
		errMsg := systemOperation.Operation + " operation  is under progress"
		return common.GeneralError(http.StatusNotAcceptable, response.ResourceCannotBeDeleted, errMsg, nil, nil)
	}
	systemOperation.Operation = "Delete"
	dbErr = systemOperation.AddSystemOperationInfo(strings.TrimSuffix(key, "/"))
	if dbErr != nil {
		log.Println(" Delete operation for system  ", key, " can't be processed ", dbErr.Error())
		errMsg := "error while trying to delete compute system: " + dbErr.Error()
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	defer func() {
		agmodel.DeleteSystemOperationInfo(strings.TrimSuffix(key, "/"))
	}()
	// Delete Subscription on odimra and also on device
	subResponse, err := e.DeleteEventSubscription(key)
	if err != nil && subResponse == nil {
		errMsg := fmt.Sprintf("error while trying to delete subscriptions: %v", err)
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	// If the DeleteEventSubscription call return status code other than http.StatusNoContent, http.StatusNotFound.
	//Then return with error(delete event subscription failed).
	if subResponse.StatusCode != http.StatusNoContent {
		log.Println("error while deleting the event subscription for ", key, " :", subResponse.Body)
	}
	// Delete Compute System Details from InMemory
	if derr := e.DeleteComputeSystem(index, key); derr != nil {
		errMsg := "error while trying to delete compute system: " + derr.Error()
		log.Println(errMsg)
		if errors.DBKeyNotFound == derr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{index, key}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}

	// Split the key by : (uuid:1) so we will get [uuid 1]
	k := strings.Split(key[index+1:], ":")
	if len(k) < 2 {
		errMsg := fmt.Sprintf("key %v doesn't have system details", key)
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	uuid := k[0]
	// Delete System Details from OnDisk
	if derr := e.DeleteSystem(uuid); derr != nil {
		errMsg := "error while trying to delete system: " + derr.Error()
		log.Println(errMsg)
		if errors.DBKeyNotFound == derr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"System", uuid}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	e.EventNotification(key, "ResourceRemoved", "SystemsCollection")
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Transfer-Encoding": "chunked",
		"Content-type":      "application/json; charset=utf-8",
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.ResourceRemoved
	args := response.Args{
		Code:    resp.StatusMessage,
		Message: "Request completed successfully",
	}
	resp.Body = args.CreateGenericErrorResponse()
	return resp
}
