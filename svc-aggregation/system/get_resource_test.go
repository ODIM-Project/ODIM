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
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
)

func mockGetAggregationSourceInfo(reqURI string) (agmodel.AggregationSource, *errors.Error) {
	var aggSource agmodel.AggregationSource
	if reqURI == "/redfish/v1/AggregationService/AggregationSources/36474ba4-a201-46aa-badf-d8104da418e8" {
		aggSource = agmodel.AggregationSource{
			HostName: "9.9.9.0",
			UserName: "admin",
			Password: []byte("admin12345"),
			Links: map[string]interface{}{
				"ConnectionMethod": map[string]interface{}{
					"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/c41cbd97-937d-1b73-c41c-1b7385d39069",
				},
			},
		}
		return aggSource, nil
	}
	return aggSource, errors.PackError(errors.DBKeyNotFound, "error: while trying to fetch Aggregation Source data: no data with the with key "+reqURI+" found")
}
func TestGetAggregationSourceCollection(t *testing.T) {
	commonResponse := response.Response{
		OdataType:    "#AggregationSourceCollection.AggregationSourceCollection",
		OdataID:      "/redfish/v1/AggregationService/AggregationSources",
		OdataContext: "/redfish/v1/$metadata#AggregationSourceCollection.AggregationSourceCollection",
		Name:         "Aggregation Source",
	}
	var resp1 = response.RPC{
		StatusCode:    http.StatusOK,
		StatusMessage: response.Success,
	}

	commonResponse.CreateGenericResponse(response.Success)
	commonResponse.Message = ""
	commonResponse.ID = ""
	commonResponse.MessageID = ""
	commonResponse.Severity = ""
	resp1.Body = agresponse.List{
		Response:     commonResponse,
		MembersCount: 1,
		Members:      []agresponse.ListMember{agresponse.ListMember{OdataID: "/redfish/v1/AggregationService/AggregationSources/058c1876-6f24-439a-8968-2af261540813"}},
	}
	p := &ExternalInterface{
		GetAllKeysFromTable: mockGetAllKeysFromTable,
	}
	tests := []struct {
		name string
		p    *ExternalInterface
		want response.RPC
	}{
		{
			name: "Postive Case",
			p:    p,
			want: resp1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetAggregationSourceCollection(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAggregationSourceCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAggregationSource(t *testing.T) {

	commonResponse := response.Response{
		OdataType:    "#AggregationSource.v1_1_0.AggregationSource",
		OdataID:      "/redfish/v1/AggregationService/AggregationSources/36474ba4-a201-46aa-badf-d8104da418e8",
		OdataContext: "/redfish/v1/$metadata#AggregationSource.AggregationSource",
		ID:           "36474ba4-a201-46aa-badf-d8104da418e8",
		Name:         "Redfish-9.9.9.0",
	}
	var resp1 = response.RPC{
		StatusCode:    http.StatusOK,
		StatusMessage: response.Success,
	}

	commonResponse.CreateGenericResponse(response.Success)
	commonResponse.Message = ""
	commonResponse.MessageID = ""
	commonResponse.Severity = ""
	resp1.Body = agresponse.AggregationSourceResponse{
		Response: commonResponse,
		HostName: "9.9.9.0",
		UserName: "admin",
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/c41cbd97-937d-1b73-c41c-1b7385d39069",
			},
		},
	}
	errMsg := "error: while trying to fetch Aggregation Source data: no data with the with key /redfish/v1/AggregationService/AggregationSources/12355 found"
	resp2 := common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"AggregationSource", "/redfish/v1/AggregationService/AggregationSources/12355"}, nil)

	p := &ExternalInterface{
		GetConnectionMethod:      mockGetConnectionMethod,
		GetAggregationSourceInfo: mockGetAggregationSourceInfo,
	}

	type args struct {
		reqURI string
	}
	tests := []struct {
		name string
		p    *ExternalInterface
		args args
		want response.RPC
	}{
		{
			name: "Postive Case",
			p:    p,
			args: args{
				reqURI: "/redfish/v1/AggregationService/AggregationSources/36474ba4-a201-46aa-badf-d8104da418e8",
			},
			want: resp1,
		},
		{
			name: "Invalid Aggregation Source URI",
			p:    p,
			args: args{
				reqURI: "/redfish/v1/AggregationService/AggregationSources/12355",
			},
			want: resp2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetAggregationSource(tt.args.reqURI); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}
