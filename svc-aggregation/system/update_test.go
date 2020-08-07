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
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"net/http"
	"reflect"
	"testing"

	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

func TestExternalInterface_UpdateAggregationSource(t *testing.T) {
	mockPluginData(t, "ILO")
	reqManager := agmodel.AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"Oem": map[string]interface{}{
				"PluginID":         "ILO",
				"PreffredAuthType": "XAuthToken",
				"PluginType":       "Compute",
			},
		},
	}
	reqBMC := agmodel.AggregationSource{
		HostName: "100.0.1.1",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"Oem": map[string]interface{}{
				"PluginID": "ILO",
			},
		},
	}

	err := agmodel.AddAggregationSource(reqManager, "/redfish/v1/AggregationService/AggregationSource/123455")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = agmodel.AddAggregationSource(reqBMC, "/redfish/v1/AggregationService/AggregationSource/123456")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	successReqManager, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: "password",
	})
	invalidReqBody1, _ := json.Marshal(AggregationSource{
		HostName: ":50000",
		UserName: "admin",
		Password: "password",
	})
	invalidReqBody2, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: "",
	})
	successReqBMC, _ := json.Marshal(AggregationSource{
		HostName: "100.0.1.1",
		UserName: "admin",
		Password: "password",
	})
	invalidReqBody3, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "",
		Password: "admin",
	})
	invalidReqBody4, _ := json.Marshal(AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: "admin",
	})

	missingparamReq, _ := json.Marshal(AggregationSource{})

	p := &ExternalInterface{
		ContactClient:   mockContactClient,
		Auth:            mockIsAuthorized,
		GetPluginStatus: GetPluginStatusForTesting,
		EncryptPassword: stubDevicePassword,
		DecryptPassword: stubDevicePassword,
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
			name: "Positive case Manager",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSource/123455",
					RequestBody: successReqManager,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotImplemented, // to be replaced http.StatusOK
			},
		},
		{
			name: "Positive case BMC",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSource/123456",
					RequestBody: successReqBMC,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotImplemented, // to be replaced http.StatusOK
			},
		},
		{
			name: "invalid aggregation source id",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSource/123466",
					RequestBody: successReqManager,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotImplemented, // to be replaced http.StatusNotFound
			},
		},
		{
			name: "invalid request manger address missing",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSource/123455",
					RequestBody: invalidReqBody1,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotImplemented, // to be replaced http.StatusBadRequest
			},
		},
		{
			name: "invalid request UserName missing",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSource/123455",
					RequestBody: invalidReqBody2,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotImplemented, // to be replaced http.StatusBadRequest
			},
		},
		{
			name: "invalid request Password missing",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSource/123455",
					RequestBody: invalidReqBody2,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotImplemented, // to be replaced http.StatusBadRequest
			},
		},
		{
			name: "invalid request Password missing",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSource/123455",
					RequestBody: invalidReqBody3,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotImplemented, // to be replaced http.StatusBadRequest
			},
		},
		{
			name: "invalid request wrong Password",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSource/123455",
					RequestBody: invalidReqBody4,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotImplemented, // to be replaced http.StatusUnauthorized
			},
		},
		{
			name: "invalid request body missing",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					URL:         "/redfish/v1/AggregationService/AggregationSource/123455",
					RequestBody: missingparamReq,
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotImplemented, // to be replaced http.StatusBadRequest
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.UpdateAggregationSource(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExternalInterface.UpdateAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}
