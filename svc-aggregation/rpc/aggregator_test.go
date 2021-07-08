// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/system"
)

func TestAggregator_GetAggregationService(t *testing.T) {
	config.SetUpMockConfig(t)
	config.Data.EnabledServices = append(config.Data.EnabledServices, "AggregationService")
	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name           string
		a              *Aggregator
		args           args
		wantStatusCode int32
	}{
		{
			name: "positive GetAggregationService",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken"},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken"},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, _ := tt.a.GetAggregationService(tt.args.ctx, tt.args.req); resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Aggregator.GetAggregationService() got = %v, wantStatusCode %v", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestAggregator_Reset(t *testing.T) {
	successReq, _ := json.Marshal(`map[string]interface{}{"parameters": []Parameters{{Name: "/redfish/v1/systems/ef83e569-7336-492a-aaee-31c02d9db831:1"}}}`)
	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
			},
			wantErr: false,
		},
		{
			name: "get session username fails",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "noDetailsToken", RequestBody: successReq},
			},
			wantErr: false,
		},
		{
			name: "unable to create task",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "noTaskToken", RequestBody: successReq},
			},
			wantErr: false,
		},
		{
			name: "task with slash",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "taskWithSlashToken", RequestBody: successReq},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.a.Reset(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.Reset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_reset(t *testing.T) {
	type args struct {
		ctx             context.Context
		taskID          string
		sessionUserName string
		req             *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				taskID: "someID",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  []byte("someData"),
				},
			},
			wantErr: false,
		},
		{
			name: "task updation fails",
			a:    &Aggregator{connector: connector},
			args: args{
				taskID: "invalid",
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  []byte("successReq"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.reset(tt.args.ctx, tt.args.taskID, tt.args.sessionUserName, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.reset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_SetDefaultBootOrder(t *testing.T) {
	successReq, _ := json.Marshal(`map[string]interface{}{"parameters": []Parameters{{Name: "/redfish/v1/systems/ef83e569-7336-492a-aaee-31c02d9db831:1"}}}`)
	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
			},
			wantErr: false,
		},
		{
			name: "get session username fails",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "noDetailsToken", RequestBody: successReq},
			},
			wantErr: false,
		},
		{
			name: "unable to create task",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "noTaskToken", RequestBody: successReq},
			},
			wantErr: false,
		},
		{
			name: "task with slash",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "taskWithSlashToken", RequestBody: successReq},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.a.SetDefaultBootOrder(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.SetDefaultBootOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_RediscoverSystemInventory(t *testing.T) {
	type args struct {
		ctx context.Context
		req *aggregatorproto.RediscoverSystemInventoryRequest
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.RediscoverSystemInventoryRequest{
					SystemID:  "someSystemID",
					SystemURL: "someURL",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.a.RediscoverSystemInventory(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.RediscoverSystemInventory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_ValidateManagerAddress(t *testing.T) {
	type args struct {
		name    string
		arg     string
		wanterr bool
	}
	tests := []struct {
		name    string
		arg     string
		wanterr bool
	}{
		{
			name:    "Valid manager address - IP",
			arg:     "127.0.0.1",
			wanterr: false,
		},
		{
			name:    "Valid manager address - IP and port",
			arg:     "127.0.0.1:1234",
			wanterr: false,
		},
		{
			name:    "Valid manager address - FQDN",
			arg:     "localhost",
			wanterr: false,
		},
		{
			name:    "Valid manager address - FQDN and Port",
			arg:     "localhost:1234",
			wanterr: false,
		},
		{
			name:    "Invalid manager address - IP",
			arg:     "a.b.c.d",
			wanterr: true,
		},
		{
			name:    "Invalid manager address - FQDN",
			arg:     "unknown",
			wanterr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateManagerAddress(tt.arg); (err != nil) != tt.wanterr {
				t.Errorf("validateManagerAddress error = %v, wantErr %v", err, tt.wanterr)
			}
		})
	}
}

func TestAggregator_AddAggreagationSource(t *testing.T) {
	config.SetUpMockConfig(t)
	addComputeRetrieval := config.AddComputeSkipResources{
		SkipResourceListUnderSystem: []string{"Chassis", "LogServices"},
	}
	mockPluginData(t, "ILO_v1.0.0")

	config.Data.AddComputeSkipResources = &addComputeRetrieval
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	successReq, _ := json.Marshal(system.AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: "password",
		Links: &system.Links{
			ConnectionMethod: &system.ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/c41cbd97-937d-1b73-c41c-1b7385d3906",
			},
		},
	})
	invalidReqBody, _ := json.Marshal(system.AggregationSource{
		HostName: ":50000",
		UserName: "admin",
		Password: "password",
		Links: &system.Links{
			ConnectionMethod: &system.ConnectionMethod{
				OdataID: "/redfish/v1/AggregationService/ConnectionMethods/c41cbd97-937d-1b73-c41c-1b7385d3906",
			},
		},
	})
	missingparamReq, _ := json.Marshal(system.AggregationSource{})
	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
			},
			wantErr: false,
		},
		{
			name: "get session username fails",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "noDetailsToken", RequestBody: successReq},
			},
			wantErr: false,
		},
		{
			name: "unable to create task",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "noTaskToken", RequestBody: successReq},
			},
			wantErr: false,
		},
		{
			name: "task with slash",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "taskWithSlashToken", RequestBody: successReq},
			},
			wantErr: false,
		},
		{
			name: "with invalid request",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: []byte("someData")},
			},
			wantErr: false,
		},
		{
			name: "Invalid Manager Address",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: invalidReqBody},
			},
			wantErr: false,
		},
		{
			name: "with missing parameters",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: missingparamReq},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.a.AddAggregationSource(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.AddAggreagationSource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_GetAllAggregationSource(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	req := agmodel.AggregationSource{
		HostName: "9.9.9.0",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/c41cbd97-937d-1b73-c41c-1b7385d3906",
			},
		},
	}
	err := agmodel.AddAggregationSource(req, "/redfish/v1/AggregationService/AggregationSources/123455")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "Positive cases",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken"},
			},
			wantErr: false,
		},
		{
			name: "Invalid Token",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.a.GetAllAggregationSource(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.GetAllAggregationSource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_GetAggregationSource(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	req := agmodel.AggregationSource{
		HostName: "9.9.9.0",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/c41cbd97-937d-1b73-c41c-1b7385d3906",
			},
		},
	}
	err := agmodel.AddAggregationSource(req, "/redfish/v1/AggregationService/AggregationSources/123455")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "Positive cases",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", URL: "/redfish/v1/AggregationService/AggregationSources/123454564"},
			},
			wantErr: false,
		},
		{
			name: "Invalid Token",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", URL: "/redfish/v1/AggregationService/AggregationSources/123454564"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.a.GetAggregationSource(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.GetAggregationSource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_UpdateAggreagationSource(t *testing.T) {
	config.SetUpMockConfig(t)

	mockPluginData(t, "ILO_v1.0.0")
	req := agmodel.AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/c41cbd97-937d-1b73-c41c-1b7385d3906",
			},
		},
	}
	err := agmodel.AddAggregationSource(req, "/redfish/v1/AggregationService/AggregationSources/123455")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	successReq, _ := json.Marshal(system.AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: "password",
	})
	invalidReqBody, _ := json.Marshal(system.AggregationSource{
		HostName: ":50000",
		UserName: "admin",
		Password: "password",
	})
	missingparamReq, _ := json.Marshal(system.AggregationSource{})
	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
			},
			wantErr: false,
		},
		{
			name: "with invalid request",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: []byte("someData")},
			},
			wantErr: false,
		},
		{
			name: "Invalid Manager Address",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: invalidReqBody},
			},
			wantErr: false,
		},
		{
			name: "with missing parameters",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: missingparamReq},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.a.UpdateAggregationSource(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.UpdateAggreagationSource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_DeleteAggregationSource(t *testing.T) {
	mockPluginData(t, "ILO_v1.0.0")
	req := agmodel.AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: []byte("admin12345"),
		Links: map[string]interface{}{
			"ConnectionMethod": map[string]interface{}{
				"@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/c41cbd97-937d-1b73-c41c-1b7385d3906",
			},
		},
	}
	err := agmodel.AddAggregationSource(req, "/redfish/v1/AggregationService/AggregationSources/123455")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", URL: "/redfish/v1/AggregationService/AggregationSources/123455"},
			},
			wantErr: false,
		},

		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", URL: "/redfish/v1/AggregationService/AggregationSources/123455"},
			},
			wantErr: false,
		},
		{
			name: "get session username fails",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "noDetailsToken", URL: "/redfish/v1/AggregationService/AggregationSources/123455"},
			},
			wantErr: false,
		},
		{
			name: "unable to create task",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "noTaskToken", URL: "/redfish/v1/AggregationService/AggregationSources/123455"},
			},
			wantErr: false,
		},
		{
			name: "task with slash",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "taskWithSlashToken", URL: "/redfish/v1/AggregationService/AggregationSources/123455"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.a.DeleteAggregationSource(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.DeleteAggregationSource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func mockSystemResourceData(body []byte, table, key string) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create(table, key, string(body)); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", table, err.Error())
	}
	return nil
}

func TestAggregator_CreateAggregate(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
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
	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name           string
		a              *Aggregator
		args           args
		wantStatusCode int32
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "positive case with empty elements",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq1},
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "with invalid request",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: []byte("someData")},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid System",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: invalidReqBody},
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "with missing parameters",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: missingparamReq},
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, _ := tt.a.CreateAggregate(tt.args.ctx, tt.args.req); resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Aggregator.CreateAggregate() status code = %v, wantStatusCode %v", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestAggregator_GetAllAggregates(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name           string
		a              *Aggregator
		args           args
		wantStatusCode int32
	}{
		{
			name: "Positive cases",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken"},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Invalid Token",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken"},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, _ := tt.a.GetAllAggregates(tt.args.ctx, tt.args.req); resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Aggregator.GetAllAggregates() error = %v, wantStatusCode %v", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestAggregator_GetAggregate(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name           string
		a              *Aggregator
		args           args
		wantStatusCode int32
	}{
		{
			name: "Positive cases",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", URL: "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73"},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Invalid Token",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", URL: "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73"},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Invalid aggregate id",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", URL: "/redfish/v1/AggregationService/Aggregates/1"},
			},
			wantStatusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, _ := tt.a.GetAggregate(tt.args.ctx, tt.args.req); resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Aggregator.GetAggregate() error = %v, wantStatusCode %v", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestAggregator_DeleteAggregate(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name           string
		a              *Aggregator
		args           args
		wantStatusCode int32
	}{
		{
			name: "Positive cases",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", URL: "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73"},
			},
			wantStatusCode: http.StatusNoContent,
		},
		{
			name: "Invalid Token",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", URL: "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73"},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Invalid aggregate id",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", URL: "/redfish/v1/AggregationService/Aggregates/1"},
			},
			wantStatusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, _ := tt.a.DeleteAggregate(tt.args.ctx, tt.args.req); resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Aggregator.DeleteAggregate() error = %v, wantStatusCode %v", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestAggregator_AddElementsToAggregate(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	reqData, _ := json.Marshal(map[string]interface{}{"@odata.id": "/redfish/v1/Systems/8c624444-87f4-4cfa-b5f9-074cd8cd114d:1"})
	err1 := mockSystemResourceData(reqData, "ComputerSystem", "/redfish/v1/Systems/8c624444-87f4-4cfa-b5f9-074cd8cd114d:1")
	if err1 != nil {
		t.Fatalf("Error in creating mock resource data :%v", err1)
	}

	successReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/8c624444-87f4-4cfa-b5f9-074cd8cd114d:1",
		},
	})

	badReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/8c624444-87f4-4cfa-b5f9-074cd8cd114d:1",
		},
	})

	duplicateReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/8c624444-87f4-4cfa-b5f9-074cd8cd114d:1",
			"/redfish/v1/Systems/8c624444-87f4-4cfa-b5f9-074cd8cd114d:1",
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

	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name           string
		a              *Aggregator
		args           args
		wantStatusCode int32
	}{
		{
			name: "Positive cases",
			a:    &Aggregator{connector: connector},
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
			name: "Invalid Token",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "invalidToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.AddElements",
					RequestBody:  successReq,
				},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Adding elements already present",
			a:    &Aggregator{connector: connector},
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
			a:    &Aggregator{connector: connector},
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
			a:    &Aggregator{connector: connector},
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
			a:    &Aggregator{connector: connector},
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
			name: "Invalid aggregate id ",
			a:    &Aggregator{connector: connector},
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
			a:    &Aggregator{connector: connector},
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
			if resp, _ := tt.a.AddElementsToAggregate(tt.args.ctx, tt.args.req); resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Aggregator.AddElementsToAggregate() error = %v, wantStatusCode %v", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestAggregator_RemoveElementsFromAggregate(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	successReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1",
		},
	})

	badReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1",
		},
	})

	duplicateReq, _ := json.Marshal(agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
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

	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name           string
		a              *Aggregator
		args           args
		wantStatusCode int32
	}{
		{
			name: "Positive cases",
			a:    &Aggregator{connector: connector},
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
			name: "Invalid Token",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "invalidToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.RemoveElements",
					RequestBody:  successReq,
				},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Removing elements not present in aggregate",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.RemoveElements",
					RequestBody:  badReq},
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Removing duplicate elements",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.RemoveElements",
					RequestBody:  duplicateReq},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Removing without elements",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.RemoveElements",
					RequestBody:  emptyReq},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Invalid element",
			a:    &Aggregator{connector: connector},
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
			name: "Invalid aggregate id ",
			a:    &Aggregator{connector: connector},
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
			a:    &Aggregator{connector: connector},
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
			if resp, _ := tt.a.RemoveElementsFromAggregate(tt.args.ctx, tt.args.req); resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Aggregator.RemoveElementsFromAggregate() error = %v, wantStatusCode %v", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestAggregator_ResetElementsOfAggregate(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	successReq, _ := json.Marshal(system.ResetRequest{
		BatchSize:                    2,
		DelayBetweenBatchesInSeconds: 2,
		ResetType:                    "ForceOff",
	})

	badReq, _ := json.Marshal(system.ResetRequest{
		BatchSize:                    2,
		DelayBetweenBatchesInSeconds: 2,
		ResetType:                    "",
	})

	missingparamReq, _ := json.Marshal(system.ResetRequest{})

	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name           string
		a              *Aggregator
		args           args
		wantStatusCode int32
	}{
		{
			name: "Positive cases",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.Reset",
					RequestBody:  successReq,
				},
			},
			wantStatusCode: http.StatusAccepted,
		},
		{
			name: "Invalid Token",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "invalidToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.Reset",
					RequestBody:  successReq,
				},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Invalid aggregate id ",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/12345/Actions/Aggregate.Reset",
					RequestBody:  successReq,
				},
			},
			wantStatusCode: http.StatusAccepted,
		},
		{
			name: "get session username fails",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "noDetailsToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/12345/Actions/Aggregate.Reset",
					RequestBody:  successReq,
				},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "unable to create task",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "noTaskToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/12345/Actions/Aggregate.Reset",
					RequestBody:  successReq,
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "task with slash",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "taskWithSlashToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/12345/Actions/Aggregate.Reset",
					RequestBody:  successReq,
				},
			},
			wantStatusCode: http.StatusAccepted,
		},
		{
			name: "Empty Reset Type",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.Reset",
					RequestBody:  badReq,
				},
			},
			wantStatusCode: http.StatusAccepted,
		},
		{
			name: "with missing parameters",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.Reset",
					RequestBody:  missingparamReq,
				},
			},
			wantStatusCode: http.StatusAccepted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, _ := tt.a.ResetElementsOfAggregate(tt.args.ctx, tt.args.req); resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Aggregator.ResetElementsOfAggregate() error = %v, wantStatusCode %v", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestAggregator_SetDefaultBootOrderElementsOfAggregate(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	req := agmodel.Aggregate{
		Elements: []string{
			"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
			"/redfish/v1/Systems/c14d91b5-3333-48bb-a7b7-75f74a137d48:1",
		},
	}
	err := agmodel.CreateAggregate(req, "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name           string
		a              *Aggregator
		args           args
		wantStatusCode int32
	}{
		{
			name: "Positive cases",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.SetDefaultBootOrder",
				},
			},
			wantStatusCode: http.StatusAccepted,
		},
		{
			name: "Invalid Token",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "invalidToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.SetDefaultBootOrder",
				},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "get session username fails",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "noDetailsToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.SetDefaultBootOrder",
				},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "unable to create task",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "noTaskToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.SetDefaultBootOrder",
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "task with slash",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "taskWithSlashToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/7ff3bd97-c41c-5de0-937d-85d390691b73/Actions/Aggregate.SetDefaultBootOrder",
				},
			},
			wantStatusCode: http.StatusAccepted,
		},
		{
			name: "Invalid aggregate id ",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					URL:          "/redfish/v1/AggregationService/Aggregates/12345/Actions/Aggregate.SetDefaultBootOrder",
				},
			},
			wantStatusCode: http.StatusAccepted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, _ := tt.a.SetDefaultBootOrderElementsOfAggregate(tt.args.ctx, tt.args.req); resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Aggregator.SetDefaultBootOrderElementsOfAggregate() error = %v, wantStatusCode %v", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestAggregator_GetAllConnectionMethods(t *testing.T) {
	config.Data.EnabledServices = append(config.Data.EnabledServices, "AggregationService")
	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name           string
		a              *Aggregator
		args           args
		wantStatusCode int32
	}{
		{
			name: "positive case",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken"},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken"},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, _ := tt.a.GetAllConnectionMethods(tt.args.ctx, tt.args.req); resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Aggregator.GetAllConnectionMethods() got = %v, wantStatusCode %v", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestAggregator_GetConnectionMethod(t *testing.T) {
	type args struct {
		ctx context.Context
		req *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name           string
		a              *Aggregator
		args           args
		wantStatusCode int32
	}{
		{
			name: "Positive cases",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", URL: "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73"},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Invalid Token",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", URL: "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73"},
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Invalid aggregate id",
			a:    &Aggregator{connector: connector},
			args: args{
				req: &aggregatorproto.AggregatorRequest{SessionToken: "validToken", URL: "/redfish/v1/AggregationService/ConnectionMethods/1"},
			},
			wantStatusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, _ := tt.a.GetConnectionMethod(tt.args.ctx, tt.args.req); resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Aggregator.GetConnectionMethod() error = %v, wantStatusCode %v", resp.StatusCode, tt.wantStatusCode)
			}
		})
	}
}
