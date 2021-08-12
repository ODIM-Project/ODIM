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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

func deleteComputeforTest(index int, key string) *errors.Error {
	if key == "/redfish/v1/Systems/del-comp-internal-err:1" {
		return errors.PackError(errors.UndefinedErrorType, "some internal error happed")
	}
	if key != "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831:1" && key != "/redfish/v1/Systems/" &&
		key != "/redfish/v1/Systems/del-sys-internal-err:1" && key != "/redfish/v1/Systems/sys-not-found:1" {
		return errors.PackError(errors.DBKeyNotFound, "error while trying to get compute details: no data with the with key "+key+" found")
	}
	return nil
}

func deleteSystemforTest(key string) *errors.Error {
	if key == "del-sys-internal-err" {
		return errors.PackError(errors.UndefinedErrorType, "some internal error happed")
	}
	if key != "ef83e569-7336-492a-aaee-31c02d9db831" {
		return errors.PackError(errors.DBKeyNotFound, "error while trying to get compute details: no data with the with key "+key+" found")
	}
	return nil
}

func mockDeleteSubscription(uuid string) (*eventsproto.EventSubResponse, error) {
	if uuid == "/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db832:1" {
		return nil, fmt.Errorf("error while trying to delete event subcription")
	} else if uuid == "/redfish/v1/Systems/unexpected-statuscode:1" {
		return &eventsproto.EventSubResponse{
			StatusCode: http.StatusCreated,
		}, nil
	}
	return &eventsproto.EventSubResponse{
		StatusCode: http.StatusNoContent,
	}, nil
}

func mockEventNotification(systemID, eventType, collectionType string) {
	return
}
func mockManagersData(id string, data map[string]interface{}) error {
	reqData, _ := json.Marshal(data)

	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("Managers", id, string(reqData)); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "Managaers", err.Error())
	}
	return nil
}

func mockContactClientForDelete(url, method, token string, odataID string, body interface{}, credentials map[string]string) (*http.Response, error) {
	if url == "https://localhost:9092/ODIM/v1/Status" || (strings.Contains(url, "/ODIM/v1/Status") && credentials["UserName"] == "noStatusUser") {
		body := `{"MessageId": "` + response.Success + `"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

func mockSystemOperationInfo() *errors.Error {
	systemOperation := agmodel.SystemOperation{
		Operation: "InventoryRediscovery ",
	}
	return systemOperation.AddSystemOperationInfo("/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831:1")
}

func TestDeleteAggregationSourceWithRediscovery(t *testing.T) {
	d := getMockExternalInterface()
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

	mockPluginData(t, "GRF_v1.0.0")
	mockManagersData("/redfish/v1/Managers/1234877451-1234", map[string]interface{}{
		"Name": "GRF_v1.0.0",
		"UUID": "1234877451-1234",
	})
	reqManagerGRF := agmodel.AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	}
	reqSuccess := agmodel.AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	}
	device := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "ef83e569-7336-492a-aaee-31c02d9db831",
		PluginID:       "GRF_v1.0.0",
	}
	mockDeviceData("ef83e569-7336-492a-aaee-31c02d9db831", device)
	mockSystemData("/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831:1")

	err := agmodel.AddAggregationSource(reqManagerGRF, "/redfish/v1/AggregationService/AggregationSources/123456")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = agmodel.AddAggregationSource(reqSuccess, "/redfish/v1/AggregationService/AggregationSources/ef83e569-7336-492a-aaee-31c02d9db831")
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
					URL:          "/redfish/v1/AggregationService/AggregationSources/ef83e569-7336-492a-aaee-31c02d9db831",
				},
			},
			want: http.StatusNotAcceptable,
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

func TestExternalInterface_DeleteAggregationSourceManager(t *testing.T) {
	d := getMockExternalInterface()
	d.ContactClient = mockContactClientForDelete
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
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
		PluginID:       "ILO_v1.0.0",
	}

	mockPluginData(t, "GRF_v1.0.0")
	mockPluginData(t, "ILO_v1.0.0")
	mockPluginData(t, "NoStatusPlugin_v1.0.0")
	mockDeviceData("24b243cf-f1e3-5318-92d9-2d6737d6b0b9", device1)
	mockManagersData("/redfish/v1/Managers/1234877451-1234", map[string]interface{}{
		"Name": "GRF_v1.0.0",
		"UUID": "1234877451-1234",
	})
	mockManagersData("/redfish/v1/Managers/1234877451-1233", map[string]interface{}{
		"Name": "ILO_v1.0.0",
		"UUID": "1234877451-1233",
	})
	mockManagersData("/redfish/v1/Managers/1234877451-1235", map[string]interface{}{
		"Name": "NoStatusPlugin_v1.0.0",
		"UUID": "1234877451-1235",
	})
	reqManagerGRF := agmodel.AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	}
	reqManagerILO := agmodel.AggregationSource{
		HostName: "100.0.0.1:50001",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/c41cbd97-937d-1b73-c41c-1b7385d39069",
			},
		},
	}
	req1 := agmodel.AggregationSource{
		HostName: "100.0.0.1:50002",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/e85bd91f-b257-4db8-b049-171099f3beec",
			},
		},
	}
	err := agmodel.AddAggregationSource(reqManagerILO, "/redfish/v1/AggregationService/AggregationSources/123455")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = agmodel.AddAggregationSource(reqManagerGRF, "/redfish/v1/AggregationService/AggregationSources/123456")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = agmodel.AddAggregationSource(req1, "/redfish/v1/AggregationService/AggregationSources/123457")
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
					URL:          "/redfish/v1/AggregationService/AggregationSources/123456",
				},
			},
			want: http.StatusNoContent,
		},
		{
			name: "deletion of plugin with mangaged devices",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					URL:          "/redfish/v1/AggregationService/AggregationSources/123455",
				},
			},
			want: http.StatusNotAcceptable,
		},
		{
			name: "deletion of plugin with invalid aggregation source id",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					URL:          "/redfish/v1/AggregationService/AggregationSources/123434",
				},
			},
			want: http.StatusNotFound,
		},
		{
			name: "plugin status check failure",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					URL:          "/redfish/v1/AggregationService/AggregationSources/123457",
				},
			},
			want: http.StatusNotAcceptable,
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
	d := getMockExternalInterface()
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
	mockPluginData(t, "GRF_v1.0.0")
	mockManagersData("/redfish/v1/Managers/1234877451-1234", map[string]interface{}{
		"Name": "GRF_v1.0.0",
		"UUID": "1234877451-1234",
	})
	reqManagerGRF := agmodel.AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	}
	reqSuccess := agmodel.AggregationSource{
		HostName: "100.0.0.1",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	}
	reqFailure := agmodel.AggregationSource{
		HostName: "100.0.0.2",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73",
			},
		},
	}
	device1 := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "ef83e569-7336-492a-aaee-31c02d9db831",
		PluginID:       "GRF_v1.0.0",
	}
	device2 := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "ef83e569-7336-492a-aaee-31c02d9db832",
		PluginID:       "GRF_v1.0.0",
	}

	mockDeviceData("ef83e569-7336-492a-aaee-31c02d9db831", device1)
	mockDeviceData("ef83e569-7336-492a-aaee-31c02d9db832", device2)
	mockSystemData("/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db831:1")
	mockSystemData("/redfish/v1/Systems/ef83e569-7336-492a-aaee-31c02d9db832:1")

	err := agmodel.AddAggregationSource(reqManagerGRF, "/redfish/v1/AggregationService/AggregationSources/123456")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = agmodel.AddAggregationSource(reqSuccess, "/redfish/v1/AggregationService/AggregationSources/ef83e569-7336-492a-aaee-31c02d9db831")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = agmodel.AddAggregationSource(reqFailure, "/redfish/v1/AggregationService/AggregationSources/ef83e569-7336-492a-aaee-31c02d9db832")
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
					URL:          "/redfish/v1/AggregationService/AggregationSources/ef83e569-7336-492a-aaee-31c02d9db831",
				},
			},
			want: http.StatusNoContent,
		},
		{
			name: "delete subscription failure",
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "SessionToken",
					URL:          "/redfish/v1/AggregationService/AggregationSources/ef83e569-7336-492a-aaee-31c02d9db832",
				},
			},
			want: http.StatusInternalServerError,
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
