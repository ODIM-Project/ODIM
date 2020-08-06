//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
)

// GetAggregationSourceCollection is to fetch all the AggregationSourceURI uri's and returns with created collection
// of AggregationSource data from odim
func GetAggregationSourceCollection() response.RPC {
	aggregationSourceKeys, err := agmodel.GetAllKeysFromTable("AggregationSource")
	if err != nil {
		log.Printf("error getting aggregation source : %v", err.Error())
		errorMessage := err.Error()
		return common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, errorMessage, []interface{}{config.Data.DBConf.OnDiskHost + ":" + config.Data.DBConf.OnDiskPort}, nil)
	}
	var members = make([]agresponse.ListMember, 0)
	for i := 0; i < len(aggregationSourceKeys); i++ {
		members = append(members, agresponse.ListMember{
			OdataID: aggregationSourceKeys[i],
		})
	}
	var resp = response.RPC{
		StatusCode:    http.StatusOK,
		StatusMessage: response.Success,
	}
	commonResponse := response.Response{
		OdataType:    "#AggregationSourceCollection.v1_0_0.AggregationSourceCollection",
		OdataID:      "/redfish/v1/AggregationService/AggregationSource",
		OdataContext: "/redfish/v1/$metadata#AggregationSourceCollection.AggregationSourceCollection",
		ID:           "AggregationSource",
		Name:         "Aggregation Source",
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	commonResponse.CreateGenericResponse(response.Success)
  commonResponse.Message = ""
	commonResponse.ID = ""
	commonResponse.MessageID = ""
	commonResponse.Severity = ""
	resp.Body = agresponse.List{
		Response:     commonResponse,
		MembersCount: len(members),
		Members:      members,
	}
	return resp
}

// GetAggregationSource is used  to fetch the AggregationSource with given aggregation source uri
//and returns AggregationSource
func GetAggregationSource(reqURI string) response.RPC {
	aggregationSource, err := agmodel.GetAggregationSourceInfo(reqURI)
	if err != nil {
		log.Printf("error getting  AggregationSource : %v", err)
		errorMessage := err.Error()
		if errors.DBKeyNotFound == err.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"AggregationSource", reqURI}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	var data = strings.Split(reqURI, "/redfish/v1/AggregationService/AggregationSource/")
	commonResponse := response.Response{
		OdataType:    "#AggregationSource.v1_0_0.AggregationSource",
		OdataID:      reqURI,
		OdataContext: "/redfish/v1/$metadata#AggregationSource.AggregationSource",
		ID:           data[1],
		Name:         "Aggregation Source",
	}
	var resp = response.RPC{
		StatusCode:    http.StatusOK,
		StatusMessage: response.Success,
	}
	resp.Header = map[string]string{
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
		HostName: aggregationSource.HostName,
		UserName: aggregationSource.UserName,
		Links:    aggregationSource.Links,
	}
	return resp
}
