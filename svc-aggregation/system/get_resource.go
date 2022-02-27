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
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
)

// GetAggregationSourceCollection is to fetch all the AggregationSourceURI uri's and returns with created collection
// of AggregationSource data from odim
func (e *ExternalInterface) GetAggregationSourceCollection() response.RPC {
	aggregationSourceKeys, err := e.GetAllKeysFromTable("AggregationSource")
	if err != nil {
		errorMessage := err.Error()
		log.Error("Unable to get aggregation source : " + errorMessage)
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
		OdataType:    "#AggregationSourceCollection.AggregationSourceCollection",
		OdataID:      "/redfish/v1/AggregationService/AggregationSources",
		OdataContext: "/redfish/v1/$metadata#AggregationSourceCollection.AggregationSourceCollection",
		ID:           "AggregationSource",
		Name:         "Aggregation Source",
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
func (e *ExternalInterface) GetAggregationSource(reqURI string) response.RPC {
	aggregationSource, err := e.GetAggregationSourceInfo(reqURI)
	if err != nil {
		errorMessage := err.Error()
		log.Error("Unable to get aggregation source : " + errorMessage)
		if errors.DBKeyNotFound == err.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"AggregationSource", reqURI}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
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
	name := connectionMethod.ConnectionMethodType + "-" + aggregationSource.HostName
	var data = strings.Split(reqURI, "/redfish/v1/AggregationService/AggregationSources/")
	commonResponse := response.Response{
		OdataType:    "#AggregationSource.v1_1_0.AggregationSource",
		OdataID:      reqURI,
		OdataContext: "/redfish/v1/$metadata#AggregationSource.AggregationSource",
		ID:           data[1],
		Name:         name,
	}
	var resp = response.RPC{
		StatusCode:    http.StatusOK,
		StatusMessage: response.Success,
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
