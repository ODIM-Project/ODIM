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
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
)

func TestGetAggregationSourceCollection(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	req := agmodel.AggregationSource{
		HostName: "9.9.9.0",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"Oem": map[string]interface{}{
				"PluginID": "ILO",
			},
		},
	}
	err := agmodel.AddAggregationSource(req, "/redfish/v1/AggregationService/AggregationSource/123455")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	commonResponse := response.Response{
		OdataType:    "#AggregationSourceCollection..v1_0_0.AggregationSourceCollection.",
		OdataID:      "/redfish/v1/AggregationService/AggregationSource",
		OdataContext: "/redfish/v1/$metadata#AggregationSourceCollection..AggregationSourceCollection.",
		ID:           "Aggregation Source",
		Name:         "Aggregation Source",
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
		Members:      []agresponse.ListMember{agresponse.ListMember{OdataID: "/redfish/v1/AggregationService/AggregationSource/123455"}},
	}
	tests := []struct {
		name string
		want response.RPC
	}{
		{
			name: "Postive Case",
			want: response.RPC{
				StatusCode: http.StatusNotImplemented, //replace with resp1 after the implementation is completed
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAggregationSourceCollection(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAggregationSourceCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAggregationSource(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	req := agmodel.AggregationSource{
		HostName: "9.9.9.0",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"Oem": map[string]interface{}{
				"PluginID": "ILO",
			},
		},
	}
	err := agmodel.AddAggregationSource(req, "/redfish/v1/AggregationService/AggregationSource/123455")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	commonResponse := response.Response{
		OdataType:    "#AggregationSource.v1_0_0.AggregationSource",
		OdataID:      "/redfish/v1/AggregationService/AggregationSource/123455",
		OdataContext: "/redfish/v1/$metadata#AggregationSource.AggregationSource",
		ID:           "121234",
		Name:         "Aggregation Source",
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

	resp1.Body = agresponse.AggregationSourceResponse{
		Response: commonResponse,
		HostName: req.HostName,
		UserName: req.UserName,
		Links:    req.Links,
	}

	type args struct {
		reqURI string
	}
	tests := []struct {
		name string
		args args
		want response.RPC
	}{
		{
			name: "Postive Case",
			args: args{
				reqURI: "/redfish/v1/AggregationService/AggregationSource/123455",
			},
			want: response.RPC{
				StatusCode: http.StatusNotImplemented, //replace with resp1 after the implementation is completed
			},
		},
		{
			name: "Invalid Aggregation Source URI",
			args: args{
				reqURI: "/redfish/v1/AggregationService/AggregationSource/12355",
			},
			want: response.RPC{
				StatusCode: http.StatusNotImplemented, //update the status code to http.StatusNotFound after implementation added
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAggregationSource(tt.args.reqURI); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}
