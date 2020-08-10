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
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

func mockSystemResourceData(body []byte, table, key string) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return err
	}
	if err = connPool.Create(table, key, string(body)); err != nil {
		return err
	}
	return nil
}

func TestExternalInterface_CreateAggregate(t *testing.T) {
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	reqData, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"})
	err := mockSystemResourceData(reqData, "ComputerSystem", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}

	reqData1, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1"})
	err = mockSystemResourceData(reqData1, "ComputerSystem", "/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}

	successReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1",
		},
	})
	successReq1, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{},
	})
	invalidReqBody, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/123456",
		},
	})
	missingparamReq, _ := json.Marshal(agmodel.Aggregate{})

	p := &ExternalInterface{
		Auth:            mockIsAuthorized,
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
					RequestBody: successReq,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusCreated,
			},
		},
		{
			name: "positive case with empty elements",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					RequestBody: successReq1,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusCreated, 
			},
		},
		{
			name: "with invalid request",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					RequestBody: []byte("someData"),
				},
			},
			want: response.RPC{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Invalid System",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					RequestBody: invalidReqBody,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name: "with missing parameters",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					RequestBody: missingparamReq,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusBadRequest,
			},
		},	
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.CreateAggregate(tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.CreateAggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}
