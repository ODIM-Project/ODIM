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
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

func TestExternalInterface_DeleteAggregationSourceManager(t *testing.T) {
	d := &ExternalInterface{
		EventNotification: mockEventNotification,
		ContactClient:     mockContactClientForDelete,
		DecryptPassword:   stubDevicePassword,
	}
	config.SetUpMockConfig(t)
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
	device1 := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "24b243cf-f1e3-5318-92d9-2d6737d6b0b9",
		PluginID:       "ILO",
	}

	mockPluginData(t, "GRF")
	mockPluginData(t, "ILO")
	mockPluginData(t, "NoStatusPlugin")
	mockDeviceData("24b243cf-f1e3-5318-92d9-2d6737d6b0b9", device1)
	mockManagersData("/redfish/v1/Managers/1234877451-1234", map[string]interface{}{
		"Name": "GRF",
		"UUID": "1234877451-1234",
	})
	mockManagersData("/redfish/v1/Managers/1234877451-1233", map[string]interface{}{
		"Name": "ILO",
		"UUID": "1234877451-1233",
	})
	mockManagersData("/redfish/v1/Managers/1234877451-1235", map[string]interface{}{
		"Name": "NoStatusPlugin",
		"UUID": "1234877451-1235",
	})
	reqManagerGRF := agmodel.AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"Oem": map[string]interface{}{
				"PluginID":         "GRF",
				"PreffredAuthType": "BasicAuth",
				"PluginType":       "Compute",
			},
		},
	}
	reqManagerILO := agmodel.AggregationSource{
		HostName: "100.0.0.1:50001",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"Oem": map[string]interface{}{
				"PluginID":         "ILO",
				"PreffredAuthType": "BasicAuth",
				"PluginType":       "Compute",
			},
		},
	}
	req1 := agmodel.AggregationSource{
		HostName: "100.0.0.1:50002",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"Oem": map[string]interface{}{
				"PluginID":         "NoStatusPlugin",
				"PreffredAuthType": "BasicAuth",
				"PluginType":       "Compute",
			},
		},
	}
	err := agmodel.AddAggregationSource(reqManagerILO, "/redfish/v1/AggregationService/AggregationSource/123455")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = agmodel.AddAggregationSource(reqManagerGRF, "/redfish/v1/AggregationService/AggregationSource/123456")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = agmodel.AddAggregationSource(req1, "/redfish/v1/AggregationService/AggregationSource/123457")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	type args struct {
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{
			name: "successful plugin deletion",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					URL:          "/redfish/v1/AggregationService/AggregationSource/123456",
				},
			},
			want: http.StatusNotImplemented, // To be changed to http.StatusNoContent
		},
		{
			name: "deletion of plugin with mangaged devices",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					URL:          "/redfish/v1/AggregationService/AggregationSource/123455",
				},
			},
			want: http.StatusNotImplemented, // To be changed to http.StatusNotAcceptable
		},
		{
			name: "deletion of plugin with invalid aggregation source id",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					URL:          "/redfish/v1/AggregationService/AggregationSource/123434",
				},
			},
			want: http.StatusNotImplemented, //  To be changed to http.StatusNotFound
		},
		{
			name: "plugin status check failure",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					URL:          "/redfish/v1/AggregationService/AggregationSource/123457",
				},
			},
			want: http.StatusNotImplemented, // To be  changed StatusNotAcceptable
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.DeleteAggregationSource(tt.args.req)
			if got.StatusCode != tt.want {
				t.Errorf("DeleteAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExternalInterface_DeleteBMC(t *testing.T) {
	d := &ExternalInterface{
		DeleteComputeSystem:     deleteComputeforTest,
		DeleteSystem:            deleteSystemforTest,
		DeleteEventSubscription: mockDeleteSubscription,
		EventNotification:       mockEventNotification,
		DecryptPassword:         stubDevicePassword,
	}
	config.SetUpMockConfig(t)
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
	reqSuccess := agmodel.AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"Oem": map[string]interface{}{
				"PluginID": "GRF",
			},
		},
	}
	reqFailure := agmodel.AggregationSource{
		HostName: "100.0.0.2",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"Oem": map[string]interface{}{
				"PluginID": "GRF",
			},
		},
	}
	err := agmodel.AddAggregationSource(reqSuccess, "/redfish/v1/AggregationService/AggregationSource/ef83e569-7336-492a-aaee-31c02d9db831")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = agmodel.AddAggregationSource(reqFailure, "/redfish/v1/AggregationService/AggregationSource/ef83e569-7336-492a-aaee-31c02d9db832")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	type args struct {
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{
			name: "successful compute deletion",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					URL:          "/redfish/v1/AggregationService/AggregationSource/ef83e569-7336-492a-aaee-31c02d9db831",
				},
			},
			want: http.StatusNotImplemented, // To be changed to http.StatusNoContent
		},
		{
			name: "delete subscription failure",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					URL:          "/redfish/v1/AggregationService/AggregationSource/ef83e569-7336-492a-aaee-31c02d9db832",
				},
			},
			want: http.StatusNotImplemented, // To be changed to http.StatusInternalServerError
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.DeleteAggregationSource(tt.args.req)
			if got.StatusCode != tt.want {
				t.Errorf("DeleteAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteAggregationSourceWithRediscovery(t *testing.T) {
	d := &ExternalInterface{
		DeleteComputeSystem:     deleteComputeforTest,
		DeleteSystem:            deleteSystemforTest,
		DeleteEventSubscription: mockDeleteSubscription,
		EventNotification:       mockEventNotification,
		DecryptPassword:         stubDevicePassword,
	}
	type args struct {
		req *aggregatorproto.AggregatorRequest
	}
	config.SetUpMockConfig(t)
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
	reqSuccess := agmodel.AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"Oem": map[string]interface{}{
				"PluginID": "GRF",
			},
		},
	}
	err := agmodel.AddAggregationSource(reqSuccess, "/redfish/v1/AggregationService/AggregationSource/ef83e569-7336-492a-aaee-31c02d9db831")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = mockSystemOperationInfo()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	tests := []struct {
		name string
		args args
		want int32
	}{

		{
			name: "successful compute deletion",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					URL:          "/redfish/v1/AggregationService/AggregationSource/ef83e569-7336-492a-aaee-31c02d9db831",
				},
			},
			want: http.StatusNotImplemented, // To be  changed StatusNotAcceptable
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.DeleteAggregationSource(tt.args.req)
			if got.StatusCode != tt.want {
				t.Errorf("DeleteAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}
