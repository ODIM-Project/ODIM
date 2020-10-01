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
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
)

func mockGetAllKeysFromTable(table string) ([]string, error) {
	if table == "ConncectionMethod" {
		return []string{"/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73"}, nil
	}
	return []string{}, fmt.Errorf("Table not found")
}

func mockGetConnectionMethod(ConnectionMethodURI string) (agmodel.ConnectionMethod, *errors.Error) {
	var connMethod agmodel.ConnectionMethod
	if ConnectionMethodURI == "7ff3bd97-c41c-5de0-937d-85d390691b73" {
		connMethod.ConnectionMethodType = "Redfish"
		connMethod.ConnectionMethodVariant = "iLO_v1.0.0"
		return connMethod, nil
	}
	return connMethod, errors.PackError(errors.DBKeyNotFound, "error while trying to get compute details: no data with the with key "+ConnectionMethodURI+" found")
}

func TestGetConnectionCollection(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	config.SetUpMockConfig(t)
	var connMethod agmodel.ConnectionMethod
	connMethod.ConnectionMethodType = "Redfish"
	connMethod.ConnectionMethodVariant = "iLO_v1.0.0"
	err := agmodel.AddConnectionMethod(connMethod, "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	commonResponse := response.Response{
		OdataType:    "#ConnectionMethodCollection.ConnectionMethodCollection",
		OdataID:      "/redfish/v1/AggregationService/ConnectionMethods",
		OdataContext: "/redfish/v1/$metadata#ConnectionMethodCollection.ConnectionMethodCollection",
		Name:         "Connection Methods",
	}
	var resp1 = response.RPC{
		StatusCode:    http.StatusOK,
		StatusMessage: response.Success,
	}
	resp1.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	commonResponse.CreateGenericResponse(response.Success)
	resp1.Body = agresponse.List{
		Response:     commonResponse,
		MembersCount: 1,
		Members:      []agresponse.ListMember{agresponse.ListMember{OdataID: "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73"}},
	}
	p := &ExternalInterface{
		Auth:                mockIsAuthorized,
		GetAllKeysFromTable: mockGetAllKeysFromTable,
		GetConnectionMethod: mockGetConnectionMethod,
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		req  *aggregatorproto.AggregatorRequest
		want response.RPC
	}{
		{
			name: "Postive Case",
			e:    p,
			req: &aggregatorproto.AggregatorRequest{
				SessionToken: "validToken",
			},
			want: resp1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.GetAllConnectionMethods(tt.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllConnectionMethods() = %v, want %v", got, tt.want)
			}
		})
	}
}
