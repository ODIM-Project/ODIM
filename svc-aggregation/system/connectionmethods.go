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
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
)

// GetAllConnectionMethods is the handler for getting the connection methods collection
func (e *ExternalInterface) GetAllConnectionMethods(ctx context.Context, req *aggregatorproto.AggregatorRequest) response.RPC {
	l.LogWithFields(ctx).Debugf("get all connection methods request body: %s", string(fmt.Sprintf("%v", req.RequestBody)))
	connectionMethods, err := e.GetAllKeysFromTable(ctx, "ConnectionMethod")
	if err != nil {
		l.LogWithFields(ctx).Error("error getting connection methods : " + err.Error())
		errorMessage := err.Error()
		return common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, errorMessage, []interface{}{config.Data.DBConf.OnDiskHost + ":" + config.Data.DBConf.OnDiskPort}, nil)
	}
	var members = make([]agresponse.ListMember, 0)
	for i := 0; i < len(connectionMethods); i++ {
		members = append(members, agresponse.ListMember{
			OdataID: connectionMethods[i],
		})
	}
	var resp = response.RPC{
		StatusCode:    http.StatusOK,
		StatusMessage: response.Success,
	}
	commonResponse := response.Response{
		OdataType:    "#ConnectionMethodCollection.ConnectionMethodCollection",
		OdataID:      "/redfish/v1/AggregationService/ConnectionMethods",
		OdataContext: "/redfish/v1/$metadata#ConnectionMethodCollection.ConnectionMethodCollection",
		Name:         "Connection Methods",
	}

	resp.Body = agresponse.List{
		Response:     commonResponse,
		MembersCount: len(members),
		Members:      members,
	}
	l.LogWithFields(ctx).Debugf("final response for get all connection methods request: %s", string(fmt.Sprintf("%v", resp.Body)))
	return resp
}

// GetConnectionMethodInfo is the handler for getting the connection method
func (e *ExternalInterface) GetConnectionMethodInfo(ctx context.Context, req *aggregatorproto.AggregatorRequest) response.RPC {
	l.LogWithFields(ctx).Debugf("get connection method info request body: %s", string(fmt.Sprintf("%v", req.RequestBody)))
	connectionmethod, err := e.GetConnectionMethod(ctx, req.URL)
	if err != nil {
		l.LogWithFields(ctx).Error("error getting  connectionmethod : " + err.Error())
		errorMessage := err.Error()
		if errors.DBKeyNotFound == err.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"ConnectionMethod", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	var data = strings.Split(req.URL, "/redfish/v1/AggregationService/ConnectionMethods/")
	commonResponse := response.Response{
		OdataType:    "#ConnectionMethod.v1_1_0.ConnectionMethod",
		OdataID:      req.URL,
		OdataContext: "/redfish/v1/$metadata#ConnectionMethod.v1_1_0.ConnectionMethod",
		ID:           data[1],
		Name:         "Connection Method",
	}
	var resp = response.RPC{
		StatusCode:    http.StatusOK,
		StatusMessage: response.Success,
	}

	links := connectionmethod.Links
	if len(links.AggregationSources) == 0 {
		links = agmodel.Links{
			AggregationSources: []agmodel.OdataID{},
		}
	}
	resp.Body = agresponse.ConnectionMethodResponse{
		Response:                commonResponse,
		ConnectionMethodType:    connectionmethod.ConnectionMethodType,
		ConnectionMethodVariant: connectionmethod.ConnectionMethodVariant,
		Links:                   links,
	}
	l.LogWithFields(ctx).Debugf("final response for get connection method info request: %s", string(fmt.Sprintf("%v", resp.Body)))
	return resp
}
