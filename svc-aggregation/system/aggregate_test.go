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
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
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
	reqData, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"})
	err := mockSystemResourceData(reqData, "ComputerSystem", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}

	reqData1, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48.1"})
	err = mockSystemResourceData(reqData1, "ComputerSystem", "/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48.1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}

	successReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48.1",
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

	p := getMockExternalInterface()
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
				StatusCode: http.StatusNotFound,
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

func TestExternalInterface_GetAllAggregates(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()

	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48.1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	p := getMockExternalInterface()
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
				},
			},
			want: response.RPC{
				StatusCode: http.StatusOK, // to be replaced http.StatusOK
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.GetAllAggregates(tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.GetAllAggregates() = %v, want %v", got.StatusCode, tt.want.StatusCode)
			}
		})
	}
}

func TestExternalInterface_GetAggregate(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()

	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48.1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	p := getMockExternalInterface()
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
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73",
				},
			},
			want: response.RPC{
				StatusCode: http.StatusOK,
			},
		},
		{
			name: "Invalid aggregate id",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/1",
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.GetAggregate(tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.GetAggregate() = %v, want %v", got.StatusCode, tt.want.StatusCode)
			}
		})
	}
}

func TestExternalInterface_DeleteAggregate(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()

	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48.1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	p := getMockExternalInterface()
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
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73",
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNoContent,
			},
		},
		{
			name: "Invalid aggregate id",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/1",
				},
			},
			want: response.RPC{
				StatusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.DeleteAggregate(tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
				t.Errorf("ExternalInterface.DeleteAggregate() = %v, want %v", got.StatusCode, tt.want.StatusCode)
			}
		})
	}
}

func TestExternalInterface_AddElementsToAggregate(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()

	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48.1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	reqData, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Systems/8c624444-87f4-4cfa-b5f9-074cd8cd114d.1"})
	err1 := mockSystemResourceData(reqData, "ComputerSystem", "/redfish/v1/Systems/8c624444-87f4-4cfa-b5f9-074cd8cd114d.1")
	if err1 != nil {
		t.Fatalf("Error in creating mock resource data :%v", err1)
	}

	successReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/8c624444-87f4-4cfa-b5f9-074cd8cd114d.1",
		},
	})

	badReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/8c624444-87f4-4cfa-b5f9-074cd8cd114d.1",
		},
	})

	duplicateReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/8c624444-87f4-4cfa-b5f9-074cd8cd114d.1",
			"/redfish/v1/Systems/8c624444-87f4-4cfa-b5f9-074cd8cd114d.1",
		},
	})

	emptyReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{},
	})

	invalidReqBody, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/123456",
		},
	})

	missingparamReq, _ := json.Marshal(agmodel.Aggregate{})

	p := getMockExternalInterface()
	type args struct {
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name           string
		e              *ExternalInterface
		args           args
		wantStatusCode int32
	}{
		{
			name: "Positive case",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.AddElements",
					RequestBody:  successReq,
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Adding elements already present",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.AddElements",
					RequestBody:  badReq},
			},
			wantStatusCode: http.StatusConflict,
		},
		{
			name: "Adding duplicate elements",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.AddElements",
					RequestBody:  duplicateReq},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Adding empty elements",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.AddElements",
					RequestBody:  emptyReq},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid element",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.AddElements",
					RequestBody:  invalidReqBody,
				},
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Invalid aggregate id",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/12345/Actions/Aggregate.AddElements",
					RequestBody:  successReq,
				},
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "with missing parameters",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.RemoveElements",
					RequestBody:  missingparamReq},
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.AddElementsToAggregate(tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.wantStatusCode) {
				t.Errorf("ExternalInterface.AddElementsToAggregate() = %v, want %v", got.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestExternalInterface_RemoveElementsFromAggregate(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()

	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48.1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	successReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48.1",
		},
	})

	badReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/8c624444-87f4-4cfa-b5f9-074cd8cd114d.1",
		},
	})

	duplicateReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"},
	})

	emptyReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{},
	})

	invalidReqBody, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/123456",
		},
	})

	missingparamReq, _ := json.Marshal(agmodel.Aggregate{})

	p := getMockExternalInterface()
	type args struct {
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name           string
		e              *ExternalInterface
		args           args
		wantStatusCode int32
	}{
		{
			name: "Positive case",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.RemoveElements",
					RequestBody:  successReq,
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Removing elements not present in aggregate",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.RemoveElements",
					RequestBody:  badReq,
				},
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Removing duplicate elements",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.RemoveElements",
					RequestBody:  duplicateReq,
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Removing without elements",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.RemoveElements",
					RequestBody:  emptyReq,
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid element",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.RemoveElements",
					RequestBody:  invalidReqBody,
				},
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Invalid aggregate id",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/12345/Actions/Aggregate.RemoveElements",
					RequestBody:  successReq,
				},
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "with missing parameters",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.RemoveElements",
					RequestBody:  missingparamReq,
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.RemoveElementsFromAggregate(tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.wantStatusCode) {
				t.Errorf("ExternalInterface.RemoveElementsFromAggregate() = %v, want %v", got.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestExternalInterface_ResetElementsOfAggregate(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	reqData, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"})
	reqData1, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48.1"})
	device1 := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
		PluginID:       "GRF",
	}
	device2 := agmodel.Target{
		ManagerAddress: "100.0.0.2",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "c14d91b5-3333-48bb-a7b7-75f74a137d48",
		PluginID:       "GRF",
	}

	mockSystemResourceData(reqData, "ComputerSystem", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1")
	mockSystemResourceData(reqData1, "ComputerSystem", "/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48.1")
	mockDeviceData("c14d91b5-3333-48bb-a7b7-75f74a137d48", device2)
	mockDeviceData("6d4a0a66-7efa-578e-83cf-44dc68d2874e", device1)
	mockPluginData(t, "GRF")

	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48.1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	successReq, _ := json.Marshal(ResetRequest{
		BatchSize:                    2,
		DelayBetweenBatchesInSeconds: 2,
		ResetType:                    "ForceOff",
	})

	badReq, _ := json.Marshal(ResetRequest{
		BatchSize:                    2,
		DelayBetweenBatchesInSeconds: 2,
		ResetType:                    "",
	})

	missingparamReq, _ := json.Marshal(ResetRequest{})

	p := getMockExternalInterface()
	type args struct {
		taskID          string
		sessionUserName string
		req             *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name           string
		e              *ExternalInterface
		args           args
		wantStatusCode int32
	}{
		{
			name: "Positive case",
			e:    p,
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.Reset",
					RequestBody:  successReq,
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "reset without child task",
			e:    p,
			args: args{
				taskID: "taskWithoutChild", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.Reset",
					RequestBody:  successReq,
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "reset without slash in subtask",
			e:    p,
			args: args{
				taskID: "subTaskWithSlash", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.Reset",
					RequestBody:  successReq,
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Invalid aggregate id",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/12345/Actions/Aggregate.Reset",
					RequestBody:  successReq,
				},
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Empty Reset Type",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.Reset",
					RequestBody:  badReq,
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "with missing parameters",
			e:    p,
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.Reset",
					RequestBody:  missingparamReq,
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.ResetElementsOfAggregate(tt.args.taskID, tt.args.sessionUserName, tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.wantStatusCode) {
				t.Errorf("ExternalInterface.ResetElementsOfAggregate() = %v, want %v", got.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestExternalInterface_SetDefaultBootOrderElementsOfAggregate(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
	defer func() {
		common.TruncateDB(common.InMemory)
		common.TruncateDB(common.OnDisk)
	}()
	reqData, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"})
	reqData1, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48.1"})
	device1 := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
		PluginID:       "GRF",
	}
	device2 := agmodel.Target{
		ManagerAddress: "100.0.0.2",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "c14d91b5-3333-48bb-a7b7-75f74a137d48",
		PluginID:       "GRF",
	}
	mockSystemResourceData(reqData, "ComputerSystem", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1")
	mockSystemResourceData(reqData1, "ComputerSystem", "/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48.1")
	mockPluginData(t, "GRF")
	mockDeviceData("c14d91b5-3333-48bb-a7b7-75f74a137d48", device2)
	mockDeviceData("6d4a0a66-7efa-578e-83cf-44dc68d2874e", device1)

	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48.1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	p := getMockExternalInterface()
	type args struct {
		taskID          string
		sessionUserName string
		req             *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name           string
		e              *ExternalInterface
		args           args
		wantStatusCode int32
	}{
		{
			name: "Positive case",
			e:    p,
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.Reset",
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "subtask creation failure",
			e:    p,
			args: args{
				taskID: "taskWithoutChild", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.Reset",
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Invalid aggregate id",
			e:    p,
			args: args{
				taskID: "someID", sessionUserName: "someUser",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/12345/Actions/Aggregate.Reset",
				},
			},
			wantStatusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.SetDefaultBootOrderElementsOfAggregate(tt.args.taskID, tt.args.sessionUserName, tt.args.req); !reflect.DeepEqual(got.StatusCode, tt.wantStatusCode) {
				t.Errorf("ExternalInterface.SetDefaultBootOrderElementsOfAggregate() = %v, want %v", got.StatusCode, tt.wantStatusCode)
			}
		})
	}
}
