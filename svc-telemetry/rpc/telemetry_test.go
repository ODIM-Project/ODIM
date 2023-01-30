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
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	tm "github.com/ODIM-Project/ODIM/svc-telemetry/telemetry"
	"github.com/stretchr/testify/assert"
)

func TestTelemetry_GetTelemetryService(t *testing.T) {
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *teleproto.TelemetryRequest
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
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{SessionToken: "validToken"},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    telemetry,
			args: args{
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{SessionToken: "InvalidToken"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.a.GetTelemetryService(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Telemetry.GetTelemetryService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTelemetry_GetMetricDefinitionCollection(t *testing.T) {
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *teleproto.TelemetryRequest
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
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{
					SessionToken: "validToken",
				},
			}, StatusCode: 200,
		},
		{
			name: "Request with invalid token",
			a:    telemetry,
			args: args{
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{
					SessionToken: "InvalidToken",
				},
			}, StatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, err := tt.a.GetMetricDefinitionCollection(tt.args.ctx, tt.args.req); err != nil {
				t.Errorf("Telemetry.GetMetricDefinitionCollection() got = %v, want %v", resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

func TestTelemetry_GetMetricReportDefinitionCollection(t *testing.T) {
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *teleproto.TelemetryRequest
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
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{
					SessionToken: "validToken",
				},
			}, StatusCode: 200,
		},
		{
			name: "Request with invalid token",
			a:    telemetry,
			args: args{
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{
					SessionToken: "InvalidToken",
				},
			}, StatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, err := tt.a.GetMetricReportDefinitionCollection(tt.args.ctx, tt.args.req); err != nil {
				t.Errorf("Telemetry.GetMetricReportDefinitionCollection() got = %v, want %v", resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

func TestTelemetry_GetMetricReportCollection(t *testing.T) {
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *teleproto.TelemetryRequest
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
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{
					SessionToken: "validToken",
				},
			}, StatusCode: 200,
		},
		{
			name: "Request with invalid token",
			a:    telemetry,
			args: args{
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{
					SessionToken: "InvalidToken",
				},
			}, StatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, err := tt.a.GetMetricReportCollection(tt.args.ctx, tt.args.req); err != nil {
				t.Errorf("Telemetry.GetMetricReportCollection() got = %v, want %v", resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

func TestTelemetry_GetTriggerCollection(t *testing.T) {
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *teleproto.TelemetryRequest
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
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{
					SessionToken: "validToken",
				},
			}, StatusCode: 200,
		},
		{
			name: "Request with invalid token",
			a:    telemetry,
			args: args{
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{
					SessionToken: "InvalidToken",
				},
			}, StatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, err := tt.a.GetTriggerCollection(tt.args.ctx, tt.args.req); err != nil {
				t.Errorf("Telemetry.GetTriggerCollection() got = %v, want %v", resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

func TestGetMetricDefinitionwithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	ctx := mockContext()
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "InvalidToken",
	}
	resp, err := telemetry.GetMetricDefinition(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetMetricDefinitionwithValidtoken(t *testing.T) {
	ctx := mockContext()
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "validToken",
	}
	resp, err := telemetry.GetMetricDefinition(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetMetricReportDefinitionwithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	ctx := mockContext()
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "InvalidToken",
	}
	resp, err := telemetry.GetMetricReportDefinition(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetMetricReportDefinitionwithValidtoken(t *testing.T) {
	ctx := mockContext()
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "validToken",
	}
	resp, err := telemetry.GetMetricReportDefinition(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetMetricReportwithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	ctx := mockContext()
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		SessionToken: "InvalidToken",
	}
	resp, err := telemetry.GetMetricReport(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetMetricReportwithValidtoken(t *testing.T) {
	config.SetUpMockConfig(t)
	ctx := mockContext()
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		SessionToken: "validToken",
		URL:          "/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom1",
	}
	resp, err := telemetry.GetMetricReport(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetTriggerwithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	ctx := mockContext()
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "InvalidToken",
	}
	resp, err := telemetry.GetTrigger(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetTriggerwithValidtoken(t *testing.T) {
	ctx := mockContext()
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "validToken",
	}
	resp, err := telemetry.GetTrigger(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestTelemetry_UpdateTrigger(t *testing.T) {
	type args struct {
		ctx context.Context
		req *teleproto.TelemetryRequest
	}
	telemetry := new(Telemetry)
	telemetry.connector = tm.MockGetExternalInterface()
	reqValid := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "validToken",
	}
	reqInValid := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "InvalidToken",
	}
	tests := []struct {
		name    string
		a       *Telemetry
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Update trigger with valid token",
			a:    telemetry,
			args: args{
				ctx: context.Background(),
				req: reqValid,
			},
			want:    0, // as update trigger is not implemented. Need to be changed with http.StatusOK once update trigger operation is implemented
			wantErr: false,
		},
		{
			name: "Update trigger with invalid token",
			a:    telemetry,
			args: args{
				ctx: context.Background(),
				req: reqInValid,
			},
			want:    http.StatusUnauthorized,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.UpdateTrigger(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Telemetry.UpdateTrigger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(int(got.StatusCode), tt.want) {
				t.Errorf("Telemetry.UpdateTrigger() = %v, want %v", got, tt.want)
			}
		})
	}
}
