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
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

// DeleteAggregationSource is the handler for removing  bmc or manager
func (e *ExternalInterface) DeleteAggregationSource(req *aggregatorproto.AggregatorRequest) response.RPC {
	var resp response.RPC

	aggregationSource, dbErr := agmodel.GetAggregationSourceInfo(req.URL)
	if dbErr != nil {
		errorMessage := dbErr.Error()
		log.Error("Unable to get AggregationSource : " + errorMessage)
		if errors.DBKeyNotFound == dbErr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"AggregationSource", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	// check whether the aggregation source is bmc or manager
	links := aggregationSource.Links.(map[string]interface{})
	connectionMethodLink := links["ConnectionMethod"].(map[string]interface{})
	connectionMethodOdataID := connectionMethodLink["@odata.id"].(string)
	connectionMethod, err := e.GetConnectionMethod(connectionMethodOdataID)
	if err != nil {
		errorMessage := err.Error()
		log.Error("Unable to get connectionmethod : " + errorMessage)
		if errors.DBKeyNotFound == err.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"ConnectionMethod", connectionMethodOdataID}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}

	requestData := strings.Split(req.URL, ":")
	resource := requestData[0]
	uuid := resource[strings.LastIndexByte(resource, '/')+1:]
	target, terr := agmodel.GetTarget(uuid)
	if terr != nil || target == nil {
		cmVariants := getConnectionMethodVariants(connectionMethod.ConnectionMethodVariant)
		if len(connectionMethod.Links.AggregationSources) > 1 {
			errMsg := fmt.Sprintf("Plugin " + cmVariants.PluginID + " can't be removed since it managing devices")
			log.Info(errMsg)
			return common.GeneralError(http.StatusNotAcceptable, response.ResourceCannotBeDeleted, errMsg, nil, nil)
		}
		// Get the plugin
		plugin, errs := agmodel.GetPluginData(cmVariants.PluginID)
		if errs != nil {
			errMsg := errs.Error()
			log.Error(errMsg)
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"plugin", cmVariants.PluginID}, nil)
		}
		// delete the manager
		resp = e.deletePlugin("/redfish/v1/Managers/" + plugin.ManagerUUID)
	} else {
		var data = strings.Split(req.URL, "/redfish/v1/AggregationService/AggregationSources/")
		systemList, dbErr := agmodel.GetAllMatchingDetails("ComputerSystem", data[1], common.InMemory)
		if dbErr != nil {
			errMsg := dbErr.Error()
			log.Error(errMsg)
			if errors.DBKeyNotFound == dbErr.ErrNo() {
				return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"Systems", "everything"}, nil)
			}
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
		}
		for _, systemURI := range systemList {
			index := strings.LastIndexAny(systemURI, "/")
			resp = e.deleteCompute(systemURI, index)
		}
	}
	if resp.StatusCode != http.StatusOK {
		return resp
	}

	if target != nil {
		plugin, errs := agmodel.GetPluginData(target.PluginID)
		if errs != nil {
			log.Error("failed to get " + target.PluginID + " plugin info: " + errs.Error())
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errs.Error(), []interface{}{"plugin", target.PluginID}, nil)
		}
		pluginStartUpData := &agmodel.PluginStartUpData{
			RequestType: "delta",
			Devices: map[string]agmodel.DeviceData{
				target.DeviceUUID: agmodel.DeviceData{
					Operation: "del",
				},
			},
		}
		if err := PushPluginStartUpData(plugin, pluginStartUpData); err != nil {
			log.Error("failed to notify device removal to " + target.PluginID + " plugin: " + err.Error())
		}
	}

	// Delete the Aggregation Source
	dbErr = agmodel.DeleteAggregationSource(req.URL)
	if dbErr != nil {
		errorMessage := "error while trying to delete AggreationSource  " + dbErr.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		log.Error(errorMessage)
		return resp
	}
	connectionMethod.Links.AggregationSources = removeAggregationSource(connectionMethod.Links.AggregationSources, agmodel.OdataID{OdataID: req.URL})
	dbErr = e.UpdateConnectionMethod(connectionMethod, connectionMethodOdataID)
	if dbErr != nil {
		errMsg := dbErr.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}

	resp = response.RPC{
		StatusCode:    http.StatusNoContent,
		StatusMessage: response.ResourceRemoved,
		Header: map[string]string{
			"Content-type":      "application/json; charset=utf-8", // TODO: add all error headers
			"Cache-Control":     "no-cache",
			"Connection":        "keep-alive",
			"Transfer-Encoding": "chunked",
			"OData-Version":     "4.0",
			"X-Frame-Options":   "sameorigin",
		},
	}
	return resp
}

// removeAggregationSource will remove the element from the slice return
// slice of remaining elements
func removeAggregationSource(slice []agmodel.OdataID, element agmodel.OdataID) []agmodel.OdataID {
	var elements []agmodel.OdataID
	for _, val := range slice {
		if val != element {
			elements = append(elements, val)
		}
	}
	return elements
}

// deleteplugin removes the given plugin
func (e *ExternalInterface) deletePlugin(oid string) response.RPC {
	var resp response.RPC
	// Get Manager Info
	data, derr := agmodel.GetResource("Managers", oid)
	if derr != nil {
		errMsg := "error while getting Managers data: " + derr.Error()
		log.Error(errMsg)
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
		log.Error(errMsg)
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"Plugin", pluginID}, nil)
	}

	systems, dberr := agmodel.GetAllSystems()
	if dberr != nil {
		errMsg := derr.Error()
		log.Error(errMsg)
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
		log.Error(errMsg)
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
		log.Error(errMsg)
		return common.GeneralError(http.StatusNotAcceptable, response.ResourceCannotBeDeleted, errMsg, nil, nil)
	}

	// deleting the manager info
	dberr = agmodel.DeleteManagersData(oid)
	if dberr != nil {
		errMsg := derr.Error()
		log.Error(errMsg)
		if errors.DBKeyNotFound == derr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"Managers", oid}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	// deleting the plugin if  zero devices are managed
	dberr = agmodel.DeletePluginData(pluginID)
	if dberr != nil {
		errMsg := derr.Error()
		log.Error(errMsg)
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
		log.Error(" Delete operation for system  " + key + " can't be processed " + dbErr.Error())
		errMsg := "error while trying to delete compute system: " + dbErr.Error()
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	if systemOperation.Operation != "" {
		log.Error("Delete operation or system  " + key + " can't be processed," +
			systemOperation.Operation + " operation  is under progress")
		errMsg := systemOperation.Operation + " operation  is under progress"
		return common.GeneralError(http.StatusNotAcceptable, response.ResourceCannotBeDeleted, errMsg, nil, nil)
	}
	systemOperation.Operation = "Delete"
	dbErr = systemOperation.AddSystemOperationInfo(strings.TrimSuffix(key, "/"))
	if dbErr != nil {
		log.Error(" Delete operation for system  " + key + " can't be processed " + dbErr.Error())
		errMsg := "error while trying to delete compute system: " + dbErr.Error()
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	defer func() {
		if err := agmodel.DeleteSystemOperationInfo(strings.TrimSuffix(key, "/")); err != nil {
			log.Errorf("failed to delete SystemOperation info of %s:%s", key, err.Error())
		}
	}()
	// Delete Subscription on odimra and also on device
	subResponse, err := e.DeleteEventSubscription(key)
	if err != nil && subResponse == nil {
		errMsg := fmt.Sprintf("error while trying to delete subscriptions: %v", err)
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	// If the DeleteEventSubscription call return status code other than http.StatusNoContent, http.StatusNotFound.
	//Then return with error(delete event subscription failed).
	if subResponse.StatusCode != http.StatusNoContent {
		log.Error("error while deleting the event subscription for " + key + " :" + string(subResponse.Body))
	}

	keys := strings.Split(key[index+1:], ":")
	chassisList, derr := agmodel.GetAllMatchingDetails("Chassis", keys[0], common.InMemory)
	if derr != nil {
		log.Error("error while trying to collect the chassis list: " + derr.Error())
	}
	// Delete Compute System Details from InMemory
	if derr := e.DeleteComputeSystem(index, key); derr != nil {
		errMsg := "error while trying to delete compute system: " + derr.Error()
		log.Error(errMsg)
		if errors.DBKeyNotFound == derr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{index, key}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}

	// Split the key by : (uuid:1) so we will get [uuid 1]
	k := strings.Split(key[index+1:], ":")
	if len(k) < 2 {
		errMsg := fmt.Sprintf("key %v doesn't have system details", key)
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	uuid := k[0]
	// Delete System Details from OnDisk
	if derr := e.DeleteSystem(uuid); derr != nil {
		errMsg := "error while trying to delete system: " + derr.Error()
		log.Error(errMsg)
		if errors.DBKeyNotFound == derr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"System", uuid}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	e.deleteWildCardValues(key[index+1:])

	for _, chassis := range chassisList {
		e.EventNotification(chassis, "ResourceRemoved", "ChassisCollection")
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

// deleteWildCardValues will delete the wild card values and
// if all the servers are deleted, then it will delete the telemetry information
func (e *ExternalInterface) deleteWildCardValues(systemID string) {
	telemetryList, dbErr := e.GetAllMatchingDetails("*", "TelemetryService", common.InMemory)
	if dbErr != nil {
		log.Error(dbErr)
		return
	}
	for _, oid := range telemetryList {
		oID := strings.Split(oid, ":")
		if !strings.Contains(oid, "MetricReports") && !strings.Contains(oid, "Collection") {
			odataID := oID[1]
			resourceData := make(map[string]interface{})
			data, dbErr := agmodel.GetResourceDetails(odataID)
			if dbErr != nil {
				log.Error("Unable to get system data : " + dbErr.Error())
				continue
			}
			// unmarshall the resourceData
			err := json.Unmarshal([]byte(data), &resourceData)
			if err != nil {
				log.Error("Unable to unmarshall  the data: " + err.Error())
				continue
			}
			var wildCards []WildCard
			var wildCardPresent bool
			wCards := resourceData["Wildcards"]
			if wCards != nil {
				for _, wCard := range getWildCard(wCards.([]interface{})) {
					wCard.Values = checkAndRemoveWildCardValue(systemID, wCard.Values)
					wildCards = append(wildCards, wCard)
					if len(wCard.Values) > 0 {
						wildCardPresent = true
					}
				}
			}
			if wildCardPresent {
				resourceData["Wildcards"] = wildCards
				resourceDataByte, err := json.Marshal(resourceData)
				if err != nil {
					continue
				}
				e.GenericSave(resourceDataByte, getResourceName(odataID, false), odataID)
			} else {
				exist, dbErr := e.CheckMetricRequest(odataID)
				if exist || dbErr != nil {
					continue
				}
				if derr := e.Delete(oID[0], odataID, common.InMemory); derr != nil {
					log.Error("error while trying to delete data: " + derr.Error())
					continue
				}
				e.updateMemberCollection(oID[0], odataID)
			}
		}
	}
}

// checkAndRemoveWildCardValue will check and remove the wild card value
func checkAndRemoveWildCardValue(val string, values []string) []string {
	var wildCardValues []string
	if len(values) < 1 {
		return wildCardValues
	}
	for _, v := range values {
		if v != val {
			wildCardValues = append(wildCardValues, v)
		}
	}
	return wildCardValues
}

// updateMemberCollection will remove the member from the collection and update into DB
func (e *ExternalInterface) updateMemberCollection(resName, odataID string) {
	resourceName := resName + "Collection"
	collectionOdataID := odataID[:strings.LastIndexByte(odataID, '/')]
	data, dbErr := e.GetResource(resourceName, collectionOdataID)
	if dbErr != nil {
		return
	}
	var telemetryInfo dmtf.Collection
	if err := json.Unmarshal([]byte(data), &telemetryInfo); err != nil {
		return
	}
	result := removeMemberFromCollection(odataID, telemetryInfo.Members)
	telemetryInfo.Members = result
	telemetryInfo.MembersCount = len(result)
	telemetryData, err := json.Marshal(telemetryInfo)
	if err != nil {
		return
	}
	e.GenericSave(telemetryData, resourceName, collectionOdataID)
}

// removeMemberFromCollection will remove the member from the collection
func removeMemberFromCollection(collectionOdataID string, telemetryInfo []*dmtf.Link) []*dmtf.Link {
	result := []*dmtf.Link{}
	for _, v := range telemetryInfo {
		if v.Oid != collectionOdataID {
			result = append(result, v)
		}
	}
	return result
}
