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
	"fmt"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-telemetry/telemetry"
	"github.com/stretchr/testify/assert"
)

func mockIsAuthorized(sessionToken string, privileges, oemPrivileges []string) response.RPC {
	if sessionToken == "InvalidToken" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "error while trying to authenticate session", nil, nil)
	}
	return common.GeneralError(http.StatusOK, response.Success, "", nil, nil)
}

func mockContactClient(url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {
	return nil, fmt.Errorf("InvalidRequest")
}

func mockGetResource(table, key string, dbType common.DbType) (string, *errors.Error) {
	return "body", nil
}

func mockGetAllKeysFromTable(table string, dbType common.DbType) ([]string, error) {
	return []string{"/redfish/v1/TelemetryService/Triggers/uuid:1"}, nil
}

func mockGetExternalInterface() *telemetry.ExternalInterface {
	return &telemetry.ExternalInterface{
		External: telemetry.External{
			Auth:          mockIsAuthorized,
			ContactClient: mockContactClient,
		},
		DB: telemetry.DB{
			GetAllKeysFromTable: mockGetAllKeysFromTable,
			GetResource:         mockGetResource,
		},
	}
}

func TestTelemetry_GetTelemetryService(t *testing.T) {
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	type args struct {
		ctx  context.Context
		req  *teleproto.TelemetryRequest
		resp *teleproto.TelemetryResponse
	}
	tests := []struct {
		name    string
		a       *Telemetry
		args    args
		wantErr bool
	}{
		{
			name: "positive GetTelemetryService",
			a:    telemetry,
			args: args{
				req:  &teleproto.TelemetryRequest{SessionToken: "validToken"},
				resp: &teleproto.TelemetryResponse{},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    telemetry,
			args: args{
				req:  &teleproto.TelemetryRequest{SessionToken: "invalidToken"},
				resp: &teleproto.TelemetryResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.GetTelemetryService(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Telemetry.GetTelemetryService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTelemetry_GetMetricDefinitionCollection(t *testing.T) {
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	type args struct {
		ctx  context.Context
		req  *teleproto.TelemetryRequest
		resp *teleproto.TelemetryResponse
	}
	tests := []struct {
		name       string
		a          *Telemetry
		args       args
		StatusCode int
	}{
		{
			name: "Request with valid token",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{
					SessionToken: "validToken",
				},
				resp: &teleproto.TelemetryResponse{},
			}, StatusCode: 200,
		},
		{
			name: "Request with invalid token",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{
					SessionToken: "invalidToken",
				},
				resp: &teleproto.TelemetryResponse{},
			}, StatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.GetMetricDefinitionCollection(tt.args.ctx, tt.args.req, tt.args.resp); err != nil {
				t.Errorf("Telemetry.GetMetricDefinitionCollection() got = %v, want %v", tt.args.resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

func TestTelemetry_GetMetricReportDefinitionCollection(t *testing.T) {
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	type args struct {
		ctx  context.Context
		req  *teleproto.TelemetryRequest
		resp *teleproto.TelemetryResponse
	}
	tests := []struct {
		name       string
		a          *Telemetry
		args       args
		StatusCode int
	}{
		{
			name: "Request with valid token",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{
					SessionToken: "validToken",
				},
				resp: &teleproto.TelemetryResponse{},
			}, StatusCode: 200,
		},
		{
			name: "Request with invalid token",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{
					SessionToken: "invalidToken",
				},
				resp: &teleproto.TelemetryResponse{},
			}, StatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.GetMetricReportDefinitionCollection(tt.args.ctx, tt.args.req, tt.args.resp); err != nil {
				t.Errorf("Telemetry.GetMetricReportDefinitionCollection() got = %v, want %v", tt.args.resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

func TestTelemetry_GetMetricReportCollection(t *testing.T) {
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	type args struct {
		ctx  context.Context
		req  *teleproto.TelemetryRequest
		resp *teleproto.TelemetryResponse
	}
	tests := []struct {
		name       string
		a          *Telemetry
		args       args
		StatusCode int
	}{
		{
			name: "Request with valid token",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{
					SessionToken: "validToken",
				},
				resp: &teleproto.TelemetryResponse{},
			}, StatusCode: 200,
		},
		{
			name: "Request with invalid token",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{
					SessionToken: "invalidToken",
				},
				resp: &teleproto.TelemetryResponse{},
			}, StatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.GetMetricReportCollection(tt.args.ctx, tt.args.req, tt.args.resp); err != nil {
				t.Errorf("Telemetry.GetMetricReportCollection() got = %v, want %v", tt.args.resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

func TestTelemetry_GetTriggerCollection(t *testing.T) {
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	type args struct {
		ctx  context.Context
		req  *teleproto.TelemetryRequest
		resp *teleproto.TelemetryResponse
	}
	tests := []struct {
		name       string
		a          *Telemetry
		args       args
		StatusCode int
	}{
		{
			name: "Request with valid token",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{
					SessionToken: "validToken",
				},
				resp: &teleproto.TelemetryResponse{},
			}, StatusCode: 200,
		},
		{
			name: "Request with invalid token",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{
					SessionToken: "invalidToken",
				},
				resp: &teleproto.TelemetryResponse{},
			}, StatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.GetTriggerCollection(tt.args.ctx, tt.args.req, tt.args.resp); err != nil {
				t.Errorf("Telemetry.GetTriggerCollection() got = %v, want %v", tt.args.resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

func TestGetMetricDefinitionwithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	var ctx context.Context
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "InvalidToken",
	}
	var resp = &teleproto.TelemetryResponse{}
	telemetry.GetMetricDefinition(ctx, req, resp)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetMetricDefinitionwithValidtoken(t *testing.T) {
	var ctx context.Context
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "validToken",
	}
	var resp = &teleproto.TelemetryResponse{}
	err := telemetry.GetMetricDefinition(ctx, req, resp)
	assert.Nil(t, err, "There should be no error")
	//To be updated once code is complete
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetMetricReportDefinitionwithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	var ctx context.Context
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "InvalidToken",
	}
	var resp = &teleproto.TelemetryResponse{}
	telemetry.GetMetricReportDefinition(ctx, req, resp)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetMetricReportDefinitionwithValidtoken(t *testing.T) {
	var ctx context.Context
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "validToken",
	}
	var resp = &teleproto.TelemetryResponse{}
	err := telemetry.GetMetricReportDefinition(ctx, req, resp)
	assert.Nil(t, err, "There should be no error")
	//To be updated once code is complete
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetMetricReportwithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	var ctx context.Context
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "InvalidToken",
	}
	var resp = &teleproto.TelemetryResponse{}
	telemetry.GetMetricReport(ctx, req, resp)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetMetricReportwithValidtoken(t *testing.T) {
	var ctx context.Context
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "validToken",
	}
	var resp = &teleproto.TelemetryResponse{}
	err := telemetry.GetMetricReport(ctx, req, resp)
	assert.Nil(t, err, "There should be no error")
	//To be updated once handlers are complete
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetTriggerwithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	var ctx context.Context
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "InvalidToken",
	}
	var resp = &teleproto.TelemetryResponse{}
	telemetry.GetTrigger(ctx, req, resp)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetTriggerwithValidtoken(t *testing.T) {
	var ctx context.Context
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "validToken",
	}
	var resp = &teleproto.TelemetryResponse{}
	err := telemetry.GetTrigger(ctx, req, resp)
	assert.Nil(t, err, "There should be no error")
	//To be updated once handlers are complete
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status code should be StatusOK.")
}
