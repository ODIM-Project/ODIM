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
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
	}
	tests := []struct {
		name    string
		a       *Aggregator
		args    args
		wantErr bool
	}{
		{
			name: "positive GetAggregationService",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken"},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken"},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.GetAggregationService(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.GetAggregationService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_AddCompute(t *testing.T) {
	config.SetUpMockConfig(t)
	addComputeRetrieval := config.AddComputeSkipResources{
		SystemCollection: []string{"Chassis", "LogServices"},
	}
	mockPluginData(t, "ILO")

	config.Data.AddComputeSkipResources = &addComputeRetrieval
	system.ActiveReqSet.ReqRecord = make(map[string]interface{})
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()

	successReq, _ := json.Marshal(system.AddResourceRequest{
		ManagerAddress: "100.0.0.1:50000",
		UserName:       "admin",
		Password:       "password",
		Oem: &system.AddOEM{
			PluginID:          "GRF",
			PreferredAuthType: "BasicAuth",
			PluginType:        "RF-GENERIC",
		},
	})
	invalidReqBody, _ := json.Marshal(system.AddResourceRequest{
		ManagerAddress: ":50000",
		UserName:       "admin",
		Password:       "password",
		Oem: &system.AddOEM{
			PluginID:          "GRF",
			PreferredAuthType: "BasicAuth",
			PluginType:        "RF-GENERIC",
		},
	})
	missingparamReq, _ := json.Marshal(system.AddResourceRequest{})
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
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
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "get session username fails",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noDetailsToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "unable to create task",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noTaskToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "task with slash",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "taskWithSlashToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "with invalid request",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: []byte("someData")},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "Invalid Manager Address",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: invalidReqBody},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "with missing parameters",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: missingparamReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.AddCompute(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.AddCompute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_DeleteCompute(t *testing.T) {
	successReq, _ := json.Marshal(`map[string]interface{}{"parameters": []Parameters{{Name: "/redfish/v1/systems/ef83e569-7336-492a-aaee-31c02d9db831:1"}}}`)
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
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
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},

		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "get session username fails",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noDetailsToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "unable to create task",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noTaskToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "task with slash",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "taskWithSlashToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.DeleteCompute(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.DeleteCompute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_deleteServer(t *testing.T) {
	type args struct {
		taskID    string
		targetURI string
		a         *Aggregator
		req       *aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive case",
			args: args{
				taskID:    "someID",
				targetURI: "someURI",
				a:         &Aggregator{connector: connector},
				req: &aggregatorproto.AggregatorRequest{
					SessionToken: "validToken",
					RequestBody:  []byte("someData"),
				},
			},
			wantErr: false,
		},
		{
			name: "task updation fails",
			args: args{
				taskID:    "invalid",
				targetURI: "someURI",
				a:         &Aggregator{connector: connector},
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
			if err := deleteServer(tt.args.taskID, tt.args.targetURI, tt.args.a, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("deleteServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_Reset(t *testing.T) {
	successReq, _ := json.Marshal(`map[string]interface{}{"parameters": []Parameters{{Name: "/redfish/v1/systems/ef83e569-7336-492a-aaee-31c02d9db831:1"}}}`)
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
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
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "get session username fails",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noDetailsToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "unable to create task",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noTaskToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "task with slash",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "taskWithSlashToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.Reset(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
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
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
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
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "get session username fails",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noDetailsToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "unable to create task",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noTaskToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "task with slash",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "taskWithSlashToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.SetDefaultBootOrder(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.SetDefaultBootOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_RediscoverSystemInventory(t *testing.T) {
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.RediscoverSystemInventoryRequest
		resp *aggregatorproto.RediscoverSystemInventoryResponse
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
			if err := tt.a.RediscoverSystemInventory(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
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
		SystemCollection: []string{"Chassis", "LogServices"},
	}
	mockPluginData(t, "ILO")

	config.Data.AddComputeSkipResources = &addComputeRetrieval
	system.ActiveReqSet.UpdateMu.Lock()
	system.ActiveReqSet.ReqRecord = make(map[string]interface{})
	system.ActiveReqSet.UpdateMu.Unlock()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	successReq, _ := json.Marshal(system.AggregationSource{
		HostName: "100.0.0.1:50000",
		UserName: "admin",
		Password: "password",
		Links: &system.Links{
			Oem: &system.AddOEM{
				PluginID:          "ILO",
				PreferredAuthType: "BasicAuth",
				PluginType:        "Compute",
			},
		},
	})
	invalidReqBody, _ := json.Marshal(system.AggregationSource{
		HostName: ":50000",
		UserName: "admin",
		Password: "password",
		Links: &system.Links{
			Oem: &system.AddOEM{
				PluginID:          "ILO",
				PreferredAuthType: "BasicAuth",
				PluginType:        "Compute",
			},
		},
	})
	missingparamReq, _ := json.Marshal(system.AggregationSource{})
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
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
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "get session username fails",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noDetailsToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "unable to create task",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "noTaskToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "task with slash",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "taskWithSlashToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "with invalid request",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: []byte("someData")},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "Invalid Manager Address",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: invalidReqBody},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "with missing parameters",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: missingparamReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.AddAggregationSource(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
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
			"Oem": map[string]interface{}{
				"PluginID": "ILO",
			},
		},
	}
	err := agmodel.AddAggregationSource(req, "/redfish/v1/AggregationService/AggregationSource/123455")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
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
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken"},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "Invalid Token",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken"},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.GetAllAggregationSource(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
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
			"Oem": map[string]interface{}{
				"PluginID": "ILO",
			},
		},
	}
	err := agmodel.AddAggregationSource(req, "/redfish/v1/AggregationService/AggregationSource/123455")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	type args struct {
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
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
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", URL: "/redfish/v1/AggregationService/AggregationSource/123454564"},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "Invalid Token",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", URL: "/redfish/v1/AggregationService/AggregationSource/123454564"},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.GetAggregationSource(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.GetAggregationSource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_UpdateAggreagationSource(t *testing.T) {
	config.SetUpMockConfig(t)

	mockPluginData(t, "ILO")
	req := agmodel.AggregationSource{
		HostName: "100.0.0.1:50000",
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
		ctx  context.Context
		req  *aggregatorproto.AggregatorRequest
		resp *aggregatorproto.AggregatorResponse
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
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "invalidToken", RequestBody: successReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "with invalid request",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: []byte("someData")},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "Invalid Manager Address",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: invalidReqBody},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
		{
			name: "with missing parameters",
			a:    &Aggregator{connector: connector},
			args: args{
				req:  &aggregatorproto.AggregatorRequest{SessionToken: "validToken", RequestBody: missingparamReq},
				resp: &aggregatorproto.AggregatorResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.UpdateAggregationSource(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Aggregator.UpdateAggreagationSource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
