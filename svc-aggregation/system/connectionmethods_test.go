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

	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
)

func mockGetAllKeysFromTable(table string) ([]string, error) {
	if table == "ConnectionMethod" {
		return []string{"/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73"}, nil
	} else if table == "Plugin" {
		return []string{"/redfish/v1/AggregationService/AggregationSources/5de0bd97-c41c-5de0-937d-85d390691b73"}, nil
	}
	return []string{}, fmt.Errorf("Table not found")
}

func mockGetConnectionMethod(ConnectionMethodURI string) (agmodel.ConnectionMethod, *errors.Error) {
	var connMethod agmodel.ConnectionMethod
	connMethod.ConnectionMethodType = "Redfish"
	switch ConnectionMethodURI {
	case "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73":
		connMethod.ConnectionMethodVariant = "Compute:BasicAuth:GRF_v1.0.0"
		return connMethod, nil
	case "/redfish/v1/AggregationService/ConnectionMethods/c41cbd97-937d-1b73-c41c-1b7385d39069":
		connMethod.ConnectionMethodVariant = "Compute:BasicAuth:ILO_v1.0.0"
		return connMethod, nil
	case "/redfish/v1/AggregationService/ConnectionMethods/6f29f281-f5e2-4873-97b7-376be668f4f4":
		connMethod.ConnectionMethodVariant = "Compute:BasicAuthentication:ILO_v1.0.0"
		return connMethod, nil
	case "/redfish/v1/AggregationService/ConnectionMethods/6456115a-e900-4c11-809f-0957031d2d56":
		connMethod.ConnectionMethodVariant = "plugin:BasicAuth:ILO_v1.0.0"
		return connMethod, nil
	case "/redfish/v1/AggregationService/ConnectionMethods/36474ba4-a201-46aa-badf-d8104da418e8":
		connMethod.ConnectionMethodVariant = "Compute:BasicAuth:PluginWithBadPassword_v1.0.0"
		return connMethod, nil
	case "/redfish/v1/AggregationService/ConnectionMethods/4298f256-c279-44e2-94f2-3987bb7d8f53":
		connMethod.ConnectionMethodVariant = "Compute:BasicAuth:PluginWithBadData_v1.0.0"
		return connMethod, nil
	case "/redfish/v1/AggregationService/ConnectionMethods/058c1876-6f24-439a-8968-2af26154081f":
		connMethod.ConnectionMethodVariant = "Compute:XAuthToken:GRF_v1.0.0"
		return connMethod, nil
	case "/redfish/v1/AggregationService/ConnectionMethods/3489af48-2e99-4d78-a250-b04641e9d98d":
		connMethod.ConnectionMethodVariant = "Compute:XAuthToken:ILO_v1.0.0"
		return connMethod, nil
	case "/redfish/v1/AggregationService/ConnectionMethods/2e99af48-2e99-4d78-a250-b04641e9b046":
		connMethod.ConnectionMethodVariant = "Compute:XAuthToken:IL_v1.0.0"
		return connMethod, nil
	case "/redfish/v1/AggregationService/ConnectionMethods/0a8992dc-8b47-4fe3-b26c-4c34048cf0d2":
		connMethod.ConnectionMethodVariant = "Compute:XAuthToken:XAuthPlugin_v1.0.0"
		return connMethod, nil
	case "/redfish/v1/AggregationService/ConnectionMethods/7551386e-b9d7-4233-a963-3841adc69e17":
		connMethod.ConnectionMethodVariant = "Compute:XAuthToken:XAuthPluginFail_v1.0.0"
		return connMethod, nil
	case "/redfish/v1/AggregationService/ConnectionMethods/e85bd91f-b257-4db8-b049-171099f3beec":
		connMethod.ConnectionMethodVariant = "Compute:BasicAuth:NoStatusPlugin_v1.0.0"
		return connMethod, nil
	}
	return connMethod, errors.PackError(errors.DBKeyNotFound, "error while trying to get compute details: no data with the with key "+ConnectionMethodURI+" found")
}

func TestGetConnectionCollection(t *testing.T) {
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

func TestExternalInterface_GetConnectionMethod(t *testing.T) {
	p := &ExternalInterface{
		Auth:                mockIsAuthorized,
		GetAllKeysFromTable: mockGetAllKeysFromTable,
		GetConnectionMethod: mockGetConnectionMethod,
	}
	type args struct {
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		args args
		want response.RPC
	}{
		{
			name: "Positive case",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
				},
			},
			want: response.RPC{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Invalid conncetion method id",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/ConnectionMethods/1",
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.GetConnectionMethodInfo(tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.GetConnectionMethodInfo() = %v, want %v", got.StatusCode, tt.want.StatusCode)
			}
		})
	}
}
