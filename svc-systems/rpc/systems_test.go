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
package rpc

import (
	"context"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
)

func TestSystems_GetSystemResource(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	reqData := []byte(`\"@odata.id\":\"/redfsh/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot\"`)
	err := mockResourceData(reqData, "SecureBoot", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	sys := new(Systems)
	sys.IsAuthorizedRPC = mockIsAuthorized

	type args struct {
		ctx  context.Context
		req  *systemsproto.GetSystemsRequest
		resp *systemsproto.SystemsResponse
	}
	tests := []struct {
		name    string
		s       *Systems
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			s:    sys,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					URL:          "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			s:    sys,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					URL:          "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot",
					SessionToken: "invalidToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid url",
			s:    sys,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					URL:          "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/SecureBoot1",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.GetSystemResource(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Systems.GetSystemResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystems_GetAllSystems(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	reqData := []byte(`\"@odata.id\":\"/redfsh/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1\"`)
	err := mockResourceData(reqData, "ComputerSystem", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	sys := new(Systems)
	sys.IsAuthorizedRPC = mockIsAuthorized

	type args struct {
		ctx  context.Context
		req  *systemsproto.GetSystemsRequest
		resp *systemsproto.SystemsResponse
	}
	tests := []struct {
		name    string
		s       *Systems
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			s:    sys,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL:          "/redfish/v1/Systems",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			s:    sys,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL:          "/redfish/v1/Systems",
					SessionToken: "invalidToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.GetSystemsCollection(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Systems.GetSystemsCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystems_GetSystems(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	reqData := []byte(`\"@odata.id\":\"/redfsh/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1\"`)
	err := mockResourceData(reqData, "ComputerSystem", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	sys := new(Systems)
	sys.IsAuthorizedRPC = mockIsAuthorized

	type args struct {
		ctx  context.Context
		req  *systemsproto.GetSystemsRequest
		resp *systemsproto.SystemsResponse
	}
	tests := []struct {
		name    string
		s       *Systems
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			s:    sys,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL:          "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			s:    sys,
			args: args{
				req: &systemsproto.GetSystemsRequest{
					URL:          "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "invalidToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.GetSystems(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Systems.GetSystems() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystems_ComputerSystemReset(t *testing.T) {
	common.SetUpMockConfig()
	sys := new(Systems)
	sys.IsAuthorizedRPC = mockIsAuthorized

	type args struct {
		ctx  context.Context
		req  *systemsproto.ComputerSystemResetRequest
		resp *systemsproto.SystemsResponse
	}
	tests := []struct {
		name    string
		s       *Systems
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			s:    sys,
			args: args{
				req: &systemsproto.ComputerSystemResetRequest{
					RequestBody:  []byte(`{"ResetType": "ForceRestart"}`),
					SystemID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			s:    sys,
			args: args{
				req: &systemsproto.ComputerSystemResetRequest{
					RequestBody:  []byte(`{"ResetType": "ForceRestart"}`),
					SystemID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "invalidToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.ComputerSystemReset(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Systems.ComputerSystemReset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystems_SetDefaultBootOrder(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	sys := new(Systems)
	sys.IsAuthorizedRPC = mockIsAuthorized

	type args struct {
		ctx  context.Context
		req  *systemsproto.DefaultBootOrderRequest
		resp *systemsproto.SystemsResponse
	}
	tests := []struct {
		name    string
		s       *Systems
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			s:    sys,
			args: args{
				req: &systemsproto.DefaultBootOrderRequest{
					SystemID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			s:    sys,
			args: args{
				req: &systemsproto.DefaultBootOrderRequest{
					SystemID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "invalidToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.SetDefaultBootOrder(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Systems.SetDefaultBootOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystems_ChangeBiosSettings(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	sys := new(Systems)
	sys.IsAuthorizedRPC = mockIsAuthorized

	type args struct {
		ctx  context.Context
		req  *systemsproto.BiosSettingsRequest
		resp *systemsproto.SystemsResponse
	}
	tests := []struct {
		name    string
		s       *Systems
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			s:    sys,
			args: args{
				req: &systemsproto.BiosSettingsRequest{
					SystemID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			s:    sys,
			args: args{
				req: &systemsproto.BiosSettingsRequest{
					SystemID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "invalidToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.ChangeBiosSettings(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Systems.ChangeBiosSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystems_ChangeBootOrderSettings(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	sys := new(Systems)
	sys.IsAuthorizedRPC = mockIsAuthorized

	type args struct {
		ctx  context.Context
		req  *systemsproto.BootOrderSettingsRequest
		resp *systemsproto.SystemsResponse
	}
	tests := []struct {
		name    string
		s       *Systems
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			s:    sys,
			args: args{
				req: &systemsproto.BootOrderSettingsRequest{
					SystemID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "validToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			s:    sys,
			args: args{
				req: &systemsproto.BootOrderSettingsRequest{
					SystemID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
					SessionToken: "invalidToken",
				},
				resp: &systemsproto.SystemsResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.ChangeBootOrderSettings(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Systems.ChangeBootOrderSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
