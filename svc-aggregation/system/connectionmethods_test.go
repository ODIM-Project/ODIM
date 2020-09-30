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
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)
func TestGetConnectionCollection(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	
	var resp1 = response.RPC{
		StatusCode:    http.StatusNotImplemented, // TODO: Need to be change as http.StatusOK
	}

	// TODO: Need to add these lines when GetAllConnectionMethods is implemented
	// commonResponse := response.Response{
	// 	OdataType:    "#ConnectionMethodCollection.ConnectionMethodCollection",
	// 	OdataID:      "/redfish/v1/AggregationService/ConnectionMethods",
	// 	OdataContext: "/redfish/v1/$metadata##ConnectionMethodCollection.ConnectionMethodCollection",
	// 	Name:         "Connection Methods",
	// }
	// var resp1 = response.RPC{
	// 	StatusCode:    http.StatusOK,
	// 	StatusMessage: response.Success,
	// }
	// resp1.Header = map[string]string{
	// 	"Cache-Control":     "no-cache",
	// 	"Connection":        "keep-alive",
	// 	"Content-type":      "application/json; charset=utf-8",
	// 	"Transfer-Encoding": "chunked",
	// 	"OData-Version":     "4.0",
	// }
	// commonResponse.CreateGenericResponse(response.Success)
	// commonResponse.Message = ""
	// commonResponse.ID = ""
	// commonResponse.MessageID = ""
	// commonResponse.Severity = ""
	// resp1.Body = agresponse.List{
	// 	Response:     commonResponse,
	// 	MembersCount: 1,
	// 	Members:      []agresponse.ListMember{agresponse.ListMember{OdataID: "/redfish/v1/AggregationService/ConnectionMethods/12345"}},
	// }
  p := &ExternalInterface{
		Auth: mockIsAuthorized,
	}
	tests := []struct {
		name string
   e    *ExternalInterface
		req *aggregatorproto.AggregatorRequest
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