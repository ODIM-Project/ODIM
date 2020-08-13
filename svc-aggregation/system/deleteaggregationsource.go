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
	"log"
	"net/http"
	"strings"

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
		log.Printf("error getting  AggregationSource : %v", dbErr)
		errorMessage := dbErr.Error()
		if errors.DBKeyNotFound == dbErr.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"AggregationSource", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	// check whether the aggregation source is bmc or manager
	links := aggregationSource.Links.(map[string]interface{})
	oem := links["Oem"].(map[string]interface{})
	if _, ok := oem["PluginType"]; ok {
		// Get the plugin
		pluginID := oem["PluginID"].(string)
		plugin, errs := agmodel.GetPluginData(pluginID)
		if errs != nil {
			errMsg := errs.Error()
			log.Printf(errMsg)
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"plugin", pluginID}, nil)
		}
		// delete the manager
		resp = e.deletePlugin("/redfish/v1/Managers/" + plugin.ManagerUUID)
	} else {
		var data = strings.Split(req.URL, "/redfish/v1/AggregationService/AggregationSource/")
		systemList, dbErr := agmodel.GetAllMatchingDetails("ComputerSystem", data[1], common.InMemory)
		if dbErr != nil {
			errMsg := dbErr.Error()
			log.Println(errMsg)
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
	// Delete the Aggregation Source
	dbErr = agmodel.DeleteAggregationSource(req.URL)
	if dbErr != nil {
		errorMessage := "error while trying to delete AggreationSource  " + dbErr.Error()
		resp.CreateInternalErrorResponse(errorMessage)
		log.Printf(errorMessage)
		return resp
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
